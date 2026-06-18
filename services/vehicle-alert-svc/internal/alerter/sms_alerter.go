package alerter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SMSAlerter struct {
	endpoint   string
	httpClient *http.Client
}

func NewSMSAlerter(endpoint string) *SMSAlerter {
	return &SMSAlerter{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type SMSRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func (a *SMSAlerter) SendSMS(ctx context.Context, phoneNumbers []string, message string) error {
	for _, phone := range phoneNumbers {
		req := SMSRequest{
			To:      phone,
			Message: message,
		}

		data, _ := json.Marshal(req)
		httpReq, err := http.NewRequestWithContext(ctx, "POST", a.endpoint, bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("failed to create SMS request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := a.httpClient.Do(httpReq)
		if err != nil {
			return fmt.Errorf("failed to send SMS to %s: %w", phone, err)
		}
		resp.Body.Close()
	}
	return nil
}
