package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/siredmar/ElwaInTheSun/pkg/mypv"
	"github.com/siredmar/ElwaInTheSun/pkg/server"
	"github.com/siredmar/ElwaInTheSun/pkg/sonnen"
)

type Controller struct {
	sonnenClient   *sonnen.Client
	mypvClient     *mypv.Client
	context        context.Context
	wattsReserved  float32
	setPointMemory float32
	maxTemp        float32
	ticker         *time.Ticker
}

func New(ctx context.Context, sonnenClient *sonnen.Client, mypvClient *mypv.Client, period time.Duration, wattsReserved float32, maxTemp float32) *Controller {
	return &Controller{
		sonnenClient:   sonnenClient,
		mypvClient:     mypvClient,
		context:        ctx,
		wattsReserved:  wattsReserved,
		setPointMemory: 0,
		maxTemp:        maxTemp,
		ticker:         time.NewTicker(period),
	}
}

func (c *Controller) UpdateConfig(new server.Config) error {
	fmt.Println("UpdatingConfig called")
	c.wattsReserved = float32(new.ReservedWatts)
	c.maxTemp = float32(new.MaxTemp)
	c.sonnenClient.SetHost(new.SonnenHost)
	c.sonnenClient.SetToken(new.SonnenToken)
	c.mypvClient.SetToken(new.MypvToken)
	c.mypvClient.SetDevice(new.MypvSerial)
	c.maxTemp = float32(new.MaxTemp)
	period, err := time.ParseDuration(new.Interval)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("new period", period)
	c.ticker.Reset(period)
	return nil
}

func (c *Controller) Run() error {
	go func() {
		configTicker := time.NewTicker(time.Second * 5)
		err := server.LoadConfig()
		if err != nil {
			fmt.Println("Failed to load config:", err)
		}
		oldConfig := server.GetConfig()
		for {
			fmt.Println("checking config")
			select {
			case <-configTicker.C:
				// check if config has changed
				err := server.LoadConfig()
				if err != nil {
					fmt.Println("Failed to load config:", err)
				}
				newConfig := server.GetConfig()
				if !reflect.DeepEqual(oldConfig, newConfig) {
					fmt.Println("config changed")
					oldConfig = newConfig
					err := c.UpdateConfig(newConfig)
					if err != nil {
						log.Println(err)
					}
				}
			case <-c.context.Done():
				return
			}

		}
	}()

	select {
	case <-c.context.Done():
		return nil
	default:
		err := c.doWork()
		if err != nil {
			return err
		}
	}
	for {
		select {
		case <-c.ticker.C:
			err := c.doWork()
			if err != nil {
				return err
			}

		case <-c.context.Done():
			return nil
		}
	}

}

func (c *Controller) doWork() error {
	batteryStatus, err := c.sonnenClient.Status()
	if err != nil {
		return err
	}
	live, err := c.mypvClient.LiveData()
	if err != nil {
		return err
	}

	// Cap temperature to maxTemp
	temp1_f := float32(live.Temp1) / 10.0
	if temp1_f >= c.maxTemp {
		log.Printf("temperature is above %f, turning off ELWA\n", c.maxTemp)
		c.setPointMemory = 0
		err = c.mypvClient.SetPowerWithDuration(int(c.setPointMemory), time.Minute)
		if err != nil {
			return err
		}
		return nil
	}

	gridFeedInW := batteryStatus.GridFeedInW
	if gridFeedInW > c.wattsReserved {
		c.setPointMemory = gridFeedInW - c.wattsReserved
	} else {
		c.setPointMemory = c.setPointMemory / 2
	}
	log.Printf("current battery grid feed in: %.0f Watts\n", batteryStatus.GridFeedInW)
	log.Printf("current ELWA power consumption: %d Watts\n", live.PowerElwa2)
	err = c.mypvClient.SetPowerWithDuration(int(c.setPointMemory), time.Minute)
	if err != nil {
		return err
	}
	log.Printf("new ELWA set point: %.0f Watts\n", c.setPointMemory)
	return nil
}
