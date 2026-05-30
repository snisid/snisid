package authz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PolicyInput struct {
	Subject      string                 `json:"subject"`
	Role         string                 `json:"role"`
	Resource     string                 `json:"resource"`
	Action       string                 `json:"action"`
	MTLSVerified bool                   `json:"mtls_verified"`
	Context      map[string]interface{} `json:"context"`
}

type OPAClient struct {
	Endpoint string
}

func (c *OPAClient) Authorize(input PolicyInput) (bool, string) {
	payload := map[string]interface{}{"input": input}
	data, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Post(c.Endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return false, fmt.Sprintf("error: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Result bool `json:"result"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Result, ""
}
