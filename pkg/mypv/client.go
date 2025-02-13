package mypv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	api "github.com/siredmar/ElwaInTheSun/pkg/api/mypv/v1"
	"github.com/siredmar/ElwaInTheSun/pkg/client"
)

const (
	// MyPVURI is the URI for the my-pv API
	MyPVURI = "https://api.my-pv.com"
	// LiveDataURI is the URI for the live data endpoint
	LiveDataURIFormatString = "/api/v1/device/%s/data"
)

type Client struct {
	client *client.Client
	device string
}

func New(token string, device string) *Client {
	httpClient := client.NewClient(MyPVURI, token, "Authorization")
	return &Client{
		client: httpClient,
		device: device,
	}
}

func (c *Client) SetToken(token string) {
	c.client.SetToken(token)
}

func (c *Client) SetDevice(device string) {
	c.device = device
}

func (c *Client) LiveData() (*api.LiveData, error) {
	resp, err := c.client.Get(c.client.Host + fmt.Sprintf(LiveDataURIFormatString, c.device))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var liveData *api.LiveData
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &liveData); err != nil {
		return nil, err
	}
	return liveData, nil
}

func (c *Client) SetPowerWithDuration(powerWatts int, duration time.Duration) error {
	data := map[string]interface{}{
		"power":                 powerWatts,
		"validForMinutes":       int(duration.Minutes()),
		"pidPower":              nil,
		"timeBoostOverride":     0,
		"timeBoostValue":        0,
		"legionellaBoostBlock":  0,
		"batteryDischargeBlock": 0,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(c.client.Host+fmt.Sprintf("/api/v1/device/%s/power", c.device), bytes.NewReader(jsonData))
	if err != nil {
		log.Println("Error posting data:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Error setting power:", resp.Status)
	}
	return nil
}
