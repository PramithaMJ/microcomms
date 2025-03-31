package httpclient

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func NewClient(baseURL string) *Client {
	c := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(10 * time.Second)

	return &Client{client: c}
}

func (c *Client) Get(endpoint string) (*resty.Response, error) {
	return c.client.R().Get(endpoint)
}
