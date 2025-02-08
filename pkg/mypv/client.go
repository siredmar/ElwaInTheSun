package mypv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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
