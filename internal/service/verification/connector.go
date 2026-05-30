package verification

import (
	"context"
	"fmt"
)

type CheckStatus string

const (
	StatusSuccess CheckStatus = "SUCCESS"
	StatusFailed  CheckStatus = "FAILED"
	StatusError   CheckStatus = "ERROR"
)

type Result struct {
	Status CheckStatus
	Reason string
	Score  int
}

type Connector interface {
	Name() string
	Verify(ctx context.Context, data map[string]interface{}) (Result, error)
}

// MockBiometricConnector simulates face/fingerprint validation
type MockBiometricConnector struct{}

func (c *MockBiometricConnector) Name() string { return "biometric" }
func (c *MockBiometricConnector) Verify(ctx context.Context, data map[string]interface{}) (Result, error) {
	if _, ok := data["biometricData"]; !ok {
		return Result{Status: StatusFailed, Reason: "No biometric data provided", Score: 0}, nil
	}
	return Result{Status: StatusSuccess, Reason: "Face match confirmed (Mock)", Score: 98}, nil
}

// MockAgencyConnector simulates government database checks
type MockAgencyConnector struct {
	AgencyName string
}

func (c *MockAgencyConnector) Name() string { return c.AgencyName }
func (c *MockAgencyConnector) Verify(ctx context.Context, data map[string]interface{}) (Result, error) {
	id, _ := data["identityId"].(string)
	if id == "revoked-id" {
		return Result{Status: StatusFailed, Reason: "Identity revoked by agency", Score: 0}, nil
	}
	return Result{Status: StatusSuccess, Reason: "Agency record verified", Score: 100}, nil
}
