package interpol

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

type SMVClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	ncbCode    string
}

func NewSMVClient(baseURL, apiKey, ncbCode string) *SMVClient {
	return &SMVClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		ncbCode: ncbCode,
	}
}

func (c *SMVClient) SubmitSMV(
	ctx context.Context,
	alert *domain.CriminalAlert,
) (string, error) {
	record := domain.SMVVehicleRecord{
		NCBReference:  alert.AlertID.String(),
		OriginCountry: c.ncbCode,
		PlateNumber:   alert.PlateNumber,
		Make:          alert.Make,
		Model:         alert.Model,
		ColorPrimary:  alert.ColorPrimary,
		StolenDate:    alert.IncidentDate,
		CrimeType:     string(alert.CrimeCategory),
	}
	if alert.VIN != nil {
		record.VIN = *alert.VIN
	}
	if alert.Year != nil {
		record.Year = int(*alert.Year)
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("failed to marshal SMV record: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/smv/vehicles", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("INTERPOL SMV request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("INTERPOL SMV returned status %d", resp.StatusCode)
	}

	var result struct {
		SMVID string `json:"smv_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode SMV response: %w", err)
	}

	return result.SMVID, nil
}

func (c *SMVClient) SubmitSMVAsync(ctx context.Context, alert *domain.CriminalAlert) {
	_, _ = c.SubmitSMV(ctx, alert)
}
