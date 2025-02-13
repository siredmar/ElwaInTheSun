package sonnen

import (
	"encoding/json"
	"io/ioutil"

	api "github.com/siredmar/ElwaInTheSun/pkg/api/sonnenv2"
	"github.com/siredmar/ElwaInTheSun/pkg/client"
)

const (
	// StatusURI is the URI for the status endpoint
	StatusURI = "/api/v2/status"
)

type Client struct {
	client *client.Client
}

func New(host, token string) *Client {
	httpClient := client.NewClient(host, token, "Auth-Token")
	return &Client{
		client: httpClient,
	}
}

func (c *Client) SetToken(token string) {
	c.client.SetToken(token)
}

func (c *Client) SetHost(host string) {
	c.client.SetHost(host)
}

// Status returns the status of the sonnen reader
func (c *Client) Status() (*api.Status, error) {
	resp, err := c.client.Get(c.client.Host + StatusURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var status *api.Status
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &status); err != nil {
		return nil, err
	}
	return status, nil
}
