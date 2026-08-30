package alerter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RadioAlerter struct {
	endpoint   string
	httpClient *http.Client
}

func NewRadioAlerter(endpoint string) *RadioAlerter {
	return &RadioAlerter{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type RadioAlert struct {
	Unit    string `json:"unit"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

func (a *RadioAlerter) BroadcastToUnit(ctx context.Context, unit string, alert interface{}) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create radio request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Target-Unit", unit)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to broadcast to unit %s: %w", unit, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("radio broadcast returned status %d", resp.StatusCode)
	}

	return nil
}
