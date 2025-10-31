package uptimerobot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultAPIURL = "https://api.uptimerobot.com/v2/getMonitors"

type GetMetricsResponse struct {
	Stat       string     `json:"stat"`
	Error      *APIError  `json:"error,omitempty"`
	Pagination Pagination `json:"pagination"`
	Monitors   []Monitor  `json:"monitors"`
}

type APIError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type Monitor struct {
	ID           int    `json:"id"`
	FriendlyName string `json:"friendly_name"`
	URL          string `json:"url"`
	Type         int    `json:"type"`
	Status       int    `json:"status"`
}

// Client is a client for the UptimeRobot API.
type Client struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// NewClient creates a new client with a default HTTP client.
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: defaultAPIURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// NewClientWithHTTP creates a client using a provided HTTP client.
func NewClientWithHTTP(apiKey string, client *http.Client) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: defaultAPIURL,
		Client:  client,
	}
}

// GetMonitors retrieves the list of monitors from UptimeRobot.
func (c *Client) GetMonitors() ([]Monitor, error) {
	form := url.Values{}
	form.Set("api_key", c.APIKey)
	form.Set("format", "json")

	req, err := http.NewRequest(http.MethodPost, c.BaseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status %d", res.StatusCode)
	}

	var decoded GetMetricsResponse
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if decoded.Stat != "ok" {
		if decoded.Error != nil {
			return nil, fmt.Errorf("api error: %s - %s", decoded.Error.Type, decoded.Error.Message)
		}
		return nil, errors.New("uptimerobot api returned failure")
	}

	return decoded.Monitors, nil
}
