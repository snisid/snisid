package foves

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/snisid/vehicle-criminal-svc/internal/repository"
)

type FovesClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewFovesClient(baseURL, token string) *FovesClient {
	return &FovesClient{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *FovesClient) VerifyStatePlate(ctx context.Context, plate string) (*repository.FovesStatePlateResult, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/vehicles/verify-state-plate?plate="+plate, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("FOVeS request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &repository.FovesStatePlateResult{IsRegistered: false}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FOVeS returned status %d", resp.StatusCode)
	}

	var result repository.FovesStatePlateResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode FOVeS response: %w", err)
	}

	return &result, nil
}
