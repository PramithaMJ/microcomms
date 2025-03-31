package httpclient

import (
	"errors"
	"log"
	"net/http"
	"time"
	"math/rand"
)

// Config holds configuration options for the HTTP client
type Config struct {
	Timeout       time.Duration
	RetryAttempts int
}

// Client struct
type Client struct {
	config Config
}

// NewClient initializes a new HTTP client
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}
	if config.RetryAttempts == 0 {
		config.RetryAttempts = 3
	}
	return &Client{config: config}
}

// Get makes an HTTP GET request with retries
func (c *Client) Get(url string) (*http.Response, error) {
	var lastErr error
	for i := 0; i < c.config.RetryAttempts; i++ {
		resp, err := c.doRequest(url)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		log.Printf("Request failed: %v, retrying... (%d/%d)", err, i+1, c.config.RetryAttempts)
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // Exponential backoff can be added here
	}
	return nil, lastErr
}

// doRequest handles the actual HTTP request
func (c *Client) doRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: c.config.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 500 {
		return nil, errors.New("server error")
	}
	return resp, nil
}
