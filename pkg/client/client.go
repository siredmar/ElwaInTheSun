package client

import (
	"io"
	"net/http"
)

type Client struct {
	Host          string
	Token         string
	Client        *http.Client
	authHeaderKey string
}

// NewClient returns a new sonnen reader client
func NewClient(host string, token string, authHeaderKey string) *Client {
	return &Client{
		Host:          host,
		Token:         token,
		Client:        &http.Client{},
		authHeaderKey: authHeaderKey,
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set(c.authHeaderKey, c.Token)
	}
	return c.Client.Do(req)
}

func (c *Client) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set(c.authHeaderKey, c.Token)
		req.Header.Set("accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
	}

	return c.Client.Do(req)
}
