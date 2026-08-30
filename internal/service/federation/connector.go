package federation

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AgencyResult struct {
	Source string                 `json:"source"`
	Data   map[string]interface{} `json:"data"`
	Error  error                  `json:"-"`
}

type AgencyConnector interface {
	Name() string
	Fetch(ctx context.Context, query string) (AgencyResult, error)
}

type HTTPConnector struct {
	AgencyName string
	BaseURL    string
	APIKey     string
	Client     *http.Client
}

func NewHTTPConnector(name, baseURL, apiKey string, timeout time.Duration) *HTTPConnector {
	return &HTTPConnector{
		AgencyName: name,
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HTTPConnector) Name() string { return c.AgencyName }

func (c *HTTPConnector) Fetch(ctx context.Context, query string) (AgencyResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/api/v1/search?q=%s", c.BaseURL, query), nil)
	if err != nil {
		return AgencyResult{Source: c.AgencyName}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return AgencyResult{Source: c.AgencyName}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AgencyResult{Source: c.AgencyName}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return AgencyResult{Source: c.AgencyName},
			fmt.Errorf("agency %s returned status %d: %s", c.AgencyName, resp.StatusCode, string(body))
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return AgencyResult{Source: c.AgencyName}, fmt.Errorf("failed to parse response: %w", err)
	}

	return AgencyResult{
		Source: c.AgencyName,
		Data:   data,
	}, nil
}
