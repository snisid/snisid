package forensics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ForensicEngine interface {
	Analyze(ctx context.Context, mediaData []byte) (*ForensicResult, error)
}

type ForensicResult struct {
	DeepfakeProbability float64  `json:"deepfake_probability"`
	Anomalies           []string `json:"anomalies"`
	ModelVersion        string   `json:"model_version"`
	ProcessingTimeMs    int64    `json:"processing_time_ms"`
}

type MesoNetForensicEngine struct {
	endpoint string
	timeout  time.Duration
	client   *http.Client
}

func NewMesoNetForensicEngine(endpoint string, timeout time.Duration) *MesoNetForensicEngine {
	if endpoint == "" {
		endpoint = os.Getenv("MESONET_SERVICE_ENDPOINT")
		if endpoint == "" {
			endpoint = "http://forensics-service.snisid.svc.cluster.local:8080"
		}
	}
	return &MesoNetForensicEngine{
		endpoint: endpoint,
		timeout:  timeout,
		client:   &http.Client{Timeout: timeout},
	}
}

func (e *MesoNetForensicEngine) Analyze(ctx context.Context, mediaData []byte) (*ForensicResult, error) {
	if len(mediaData) == 0 {
		return nil, fmt.Errorf("empty media data")
	}

	reqBody := map[string]interface{}{
		"data": mediaData,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint+"/v1/analyze", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mesonet api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mesonet api error: %s", string(respBody))
	}

	var result ForensicResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}
