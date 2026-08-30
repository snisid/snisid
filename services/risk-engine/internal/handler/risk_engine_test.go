package handler

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func resetEngine() {
	engine = nil
	once = sync.Once{}
}

func TestRiskEngine_Evaluate_Approve(t *testing.T) {
	resetEngine()
	eng := getEngine()

	req := RiskRequest{
		TransactionID: "tx-001",
		UserID:        "user-01",
		Amount:        500,
		SourceIP:      "192.168.1.1",
		DeviceID:      "device-01",
		Timestamp:     time.Now().Unix(),
	}
	resp := eng.evaluate(req)
	assert.Equal(t, "tx-001", resp.TransactionID)
	assert.Equal(t, "APPROVE", resp.Decision)
	assert.Greater(t, resp.ComputedAt, int64(0))
	assert.NotEmpty(t, resp.Factors)
}

func TestRiskEngine_Evaluate_Review(t *testing.T) {
	resetEngine()
	eng := getEngine()

	req := RiskRequest{
		TransactionID: "tx-002",
		UserID:        "user-02",
		Amount:        20000,
		SourceIP:      "10.0.0.1",
		DeviceID:      "device-02",
		Timestamp:     time.Now().Unix(),
	}
	resp := eng.evaluate(req)
	assert.Equal(t, "REVIEW", resp.Decision)
}

func TestRiskEngine_Evaluate_Block(t *testing.T) {
	resetEngine()
	eng := getEngine()
	eng.thresholds.BlockThreshold = 0.3

	req := RiskRequest{
		TransactionID: "tx-003",
		UserID:        "user-03",
		Amount:        100000,
		SourceIP:      "10.0.0.1",
		DeviceID:      "device-03",
		Timestamp:     time.Now().Unix(),
	}
	resp := eng.evaluate(req)
	assert.Equal(t, "BLOCK", resp.Decision)
}

func TestRiskEngine_Evaluate_BlockViaBlacklist(t *testing.T) {
	resetEngine()
	eng := getEngine()
	eng.AddToBlacklist("1.2.3.4", "")

	req := RiskRequest{
		TransactionID: "tx-004",
		UserID:        "user-04",
		Amount:        100,
		SourceIP:      "1.2.3.4",
		DeviceID:      "device-04",
		Timestamp:     time.Now().Unix(),
	}
	resp := eng.evaluate(req)
	assert.Equal(t, "BLOCK", resp.Decision)
}

func TestRiskEngine_Evaluate_BlockViaDeviceBlacklist(t *testing.T) {
	resetEngine()
	eng := getEngine()
	eng.AddToBlacklist("", "stolen-device")

	req := RiskRequest{
		TransactionID: "tx-005",
		UserID:        "user-05",
		Amount:        100,
		SourceIP:      "10.0.0.1",
		DeviceID:      "stolen-device",
		Timestamp:     time.Now().Unix(),
	}
	resp := eng.evaluate(req)
	assert.Equal(t, "BLOCK", resp.Decision)
}

func TestRiskEngine_EvaluateAmount(t *testing.T) {
	resetEngine()
	eng := getEngine()

	tests := []struct {
		name   string
		amount float64
		expSc  float64
		expWt  float64
	}{
		{"normal amount", 500, 0.05, 0.35},
		{"over 1000", 1500, 0.2, 0.35},
		{"over 5000", 6000, 0.4, 0.35},
		{"over 10000", 20000, 0.7, 0.35},
		{"over 50000", 100000, 1.0, 0.35},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			factor := eng.evaluateAmount(tc.amount)
			assert.Equal(t, "amount", factor.Name)
			assert.InDelta(t, tc.expSc, factor.Score, 0.001)
			assert.InDelta(t, tc.expWt, factor.Weight, 0.001)
		})
	}
}

func TestRiskEngine_EvaluateVelocity_InsufficientHistory(t *testing.T) {
	resetEngine()
	eng := getEngine()

	factor := eng.evaluateVelocity("new-user")
	assert.Equal(t, "velocity", factor.Name)
	assert.InDelta(t, 0.1, factor.Score, 0.001)
	assert.Equal(t, "insufficient history", factor.Reason)
}

func TestRiskEngine_EvaluateVelocity_WithHistory(t *testing.T) {
	resetEngine()
	eng := getEngine()

	now := time.Now().Unix()
	eng.mu.Lock()
	eng.userHistory["frequent-user"] = []RiskRequest{
		{Timestamp: now - 10},
		{Timestamp: now - 5},
	}
	eng.mu.Unlock()

	factor := eng.evaluateVelocity("frequent-user")
	assert.Equal(t, "velocity", factor.Name)
	assert.Greater(t, factor.Score, 0.0)
}

func TestRiskEngine_EvaluateVelocity_ExceededLimit(t *testing.T) {
	resetEngine()
	eng := getEngine()

	now := time.Now().Unix()
	eng.mu.Lock()
	history := make([]RiskRequest, 20)
	for i := range history {
		history[i] = RiskRequest{Timestamp: now - int64(i)*10}
	}
	eng.userHistory["rapid-user"] = history
	eng.mu.Unlock()

	factor := eng.evaluateVelocity("rapid-user")
	assert.Equal(t, 1.0, factor.Score)
	assert.Equal(t, "velocity limit exceeded", factor.Reason)
}

func TestRiskEngine_EvaluateBlacklists_NotBlacklisted(t *testing.T) {
	resetEngine()
	eng := getEngine()

	req := RiskRequest{SourceIP: "192.168.1.1", DeviceID: "device-01"}
	factor := eng.evaluateBlacklists(req)
	assert.Equal(t, 0.0, factor.Score)
}

func TestRiskEngine_EvaluateBlacklists_IPBlacklisted(t *testing.T) {
	resetEngine()
	eng := getEngine()
	eng.AddToBlacklist("192.168.1.1", "")

	req := RiskRequest{SourceIP: "192.168.1.1", DeviceID: "device-01"}
	factor := eng.evaluateBlacklists(req)
	assert.Equal(t, 1.0, factor.Score)
}

func TestRiskEngine_EvaluateBlacklists_DeviceBlacklisted(t *testing.T) {
	resetEngine()
	eng := getEngine()
	eng.AddToBlacklist("", "bad-device")

	req := RiskRequest{SourceIP: "192.168.1.1", DeviceID: "bad-device"}
	factor := eng.evaluateBlacklists(req)
	assert.Equal(t, 1.0, factor.Score)
}

func TestRiskEngine_AddToBlacklist(t *testing.T) {
	resetEngine()
	eng := getEngine()

	eng.AddToBlacklist("10.0.0.1", "device-x")
	eng.mu.RLock()
	assert.True(t, eng.ipBlacklist["10.0.0.1"])
	assert.True(t, eng.deviceBlacklist["device-x"])
	eng.mu.RUnlock()
}

func TestFormatVelocityReason(t *testing.T) {
	assert.Equal(t, "5 transactions in window", formatVelocityReason(5, 10))
	assert.Equal(t, "velocity limit exceeded", formatVelocityReason(10, 10))
	assert.Equal(t, "velocity limit exceeded", formatVelocityReason(15, 10))
}

func TestRiskEngine_ConcurrentAccess(t *testing.T) {
	resetEngine()
	eng := getEngine()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			req := RiskRequest{
				TransactionID: "tx-concurrent",
				UserID:        "user-concurrent",
				Amount:        1000,
				SourceIP:      "10.0.0.1",
				DeviceID:      "device-01",
				Timestamp:     time.Now().Unix(),
			}
			resp := eng.evaluate(req)
			assert.NotEmpty(t, resp.Decision)
		}(i)
	}
	wg.Wait()
}

func TestGetEngine_Singleton(t *testing.T) {
	resetEngine()
	e1 := getEngine()
	e2 := getEngine()
	assert.Same(t, e1, e2)
}

func TestRiskEngine_Evaluate_MultipleFactors(t *testing.T) {
	resetEngine()
	eng := getEngine()

	req := RiskRequest{
		TransactionID: "tx-multi",
		UserID:        "user-multi",
		Amount:        25000,
		SourceIP:      "10.0.0.1",
		DeviceID:      "device-01",
		Timestamp:     time.Now().Unix(),
	}
	resp := eng.evaluate(req)
	assert.Len(t, resp.Factors, 4)

	factorNames := make([]string, len(resp.Factors))
	for i, f := range resp.Factors {
		factorNames[i] = f.Name
	}
	assert.Contains(t, factorNames, "amount")
	assert.Contains(t, factorNames, "velocity")
	assert.Contains(t, factorNames, "blacklist")
	assert.Contains(t, factorNames, "time")
}

func TestFormatInt(t *testing.T) {
	assert.Equal(t, "0", formatInt(0))
	assert.Equal(t, "42", formatInt(42))
	assert.Equal(t, "-5", formatInt(-5))
}
