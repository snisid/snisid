package federation

import (
	"context"
	"fmt"
	"time"

	"github.com/snisid/platform/backend/internal/platform/logger"
	"go.uber.org/zap"
)

type AgencyResult struct {
	Source string
	Data   map[string]interface{}
	Error  error
}

type AgencyConnector interface {
	Name() string
	Fetch(ctx context.Context, query string) (AgencyResult, error)
}

type MockConnector struct {
	AgencyName string
	Delay      time.Duration
}

func (c *MockConnector) Name() string { return c.AgencyName }

func (c *MockConnector) Fetch(ctx context.Context, query string) (AgencyResult, error) {
	// Simulation of external API call with context awareness
	select {
	case <-time.After(c.Delay):
		return AgencyResult{
			Source: c.AgencyName,
			Data:   map[string]interface{}{"status": "active", "record_id": "MOCK-123"},
		}, nil
	case <-ctx.Done():
		return AgencyResult{Source: c.AgencyName}, ctx.Err()
	}
}
