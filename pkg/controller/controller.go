package controller

import (
	"context"
	"log"
	"time"

	"github.com/siredmar/ElwaInTheSun/pkg/mypv"
	"github.com/siredmar/ElwaInTheSun/pkg/sonnen"
)

type Controller struct {
	sonnenClient   *sonnen.Client
	mypvClient     *mypv.Client
	period         time.Duration
	context        context.Context
	wattsReserved  float32
	setPointMemory float32
	maxTemp        float32
}

func New(ctx context.Context, sonnenClient *sonnen.Client, mypvClient *mypv.Client, period time.Duration, wattsReserved float32, maxTemp float32) *Controller {
	return &Controller{
		sonnenClient:   sonnenClient,
		mypvClient:     mypvClient,
		period:         period,
		context:        ctx,
		wattsReserved:  wattsReserved,
		setPointMemory: 0,
		maxTemp:        maxTemp,
	}
}

func (c *Controller) Run() error {
	ticker := time.NewTicker(c.period)
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
		case <-ticker.C:
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
	temp1_f := float32(live.Temp1) * 10.0
	if temp1_f >= c.maxTemp {
		log.Printf("temperature is above %f, turning off ELWA\n", c.maxTemp)
		err = c.mypvClient.SetPowerWithDuration(0, time.Minute)
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
