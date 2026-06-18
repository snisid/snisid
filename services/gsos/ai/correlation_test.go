package ai

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCorrelationLayer(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	assert.Equal(t, "v1.0", c.ModelVersion)
	assert.Empty(t, c.eventHistory)
	assert.Equal(t, 15*time.Minute, c.windowSize)
}

func TestIngestEvent(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	c.IngestEvent(SecurityEvent{
		ID: "EVT-001", Country: "HTI", EventType: "LOGIN_FAILURE",
		Severity: 0.5, Timestamp: time.Now().Unix(),
	})
	assert.Len(t, c.eventHistory, 1)
}

func TestIngestEvent_PrunesOldEvents(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	oldTimestamp := time.Now().Add(-2 * time.Hour).Unix()
	for i := 0; i < 10; i++ {
		c.IngestEvent(SecurityEvent{
			ID:        "EVT-OLD",
			Country:   "HTI",
			Timestamp: oldTimestamp,
		})
	}
	// Fresh event triggers pruning
	c.IngestEvent(SecurityEvent{
		ID: "EVT-FRESH", Country: "HTI", Timestamp: time.Now().Unix(),
	})
	assert.Len(t, c.eventHistory, 1)
	assert.Equal(t, "EVT-FRESH", c.eventHistory[0].ID)
}

func TestAnalyzeGlobalThreats_InsufficientEvents(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	results := c.AnalyzeGlobalThreats([]SecurityEvent{
		{ID: "EVT-1", Country: "HTI", EventType: "LOGIN_FAILURE", Severity: 0.5},
	})
	assert.Nil(t, results)
}

func TestAnalyzeGlobalThreats_SingleTypeCorrelation(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	now := time.Now().Unix()

	events := []SecurityEvent{
		{ID: "E1", Country: "HTI", EventType: "LOGIN_FAILURE", Severity: 0.8, Timestamp: now},
		{ID: "E2", Country: "DOM", EventType: "LOGIN_FAILURE", Severity: 0.7, Timestamp: now + 60},
		{ID: "E3", Country: "CUB", EventType: "LOGIN_FAILURE", Severity: 0.9, Timestamp: now + 120},
	}

	results := c.AnalyzeGlobalThreats(events)
	require.NotEmpty(t, results)
	assert.Equal(t, "LOGIN_FAILURE", results[0].EventTypes[0])
	assert.True(t, results[0].Score > 0)
}

func TestAnalyzeGlobalThreats_CoordinatedAttack(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	now := time.Now().Unix()

	events := make([]SecurityEvent, 6)
	for i := range events {
		events[i] = SecurityEvent{
			ID:        "E",
			Country:   []string{"HTI", "DOM", "CUB", "JAM", "BRB", "BHS"}[i],
			EventType: "DATA_EXFILTRATION",
			Severity:  0.9,
			Timestamp: now + int64(i*30),
		}
	}

	results := c.AnalyzeGlobalThreats(events)
	require.NotEmpty(t, results)
	assert.True(t, results[0].IsCoordinated)
	assert.Equal(t, "DATA_EXFIL", results[0].Technique)
}

func TestAnalyzeGlobalThreats_HighScoreWarning(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	now := time.Now().Unix()

	events := []SecurityEvent{
		{ID: "E1", Country: "HTI", EventType: "FRAUD_TRANSACTION", Severity: 0.95, Timestamp: now},
		{ID: "E2", Country: "DOM", EventType: "FRAUD_TRANSACTION", Severity: 0.95, Timestamp: now + 30},
		{ID: "E3", Country: "CUB", EventType: "FRAUD_TRANSACTION", Severity: 0.95, Timestamp: now + 60},
	}

	results := c.AnalyzeGlobalThreats(events)
	require.NotEmpty(t, results)
	assert.Equal(t, "SYNTHETIC_FRAUD", results[0].Technique)
}

func TestAnalyzeGlobalThreats_SpreadScore(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	now := time.Now().Unix()

	events := []SecurityEvent{
		{ID: "E1", Country: "HTI", EventType: "SUSPICIOUS_ACCESS", Severity: 0.6, Timestamp: now},
		{ID: "E2", Country: "DOM", EventType: "SUSPICIOUS_ACCESS", Severity: 0.6, Timestamp: now - 3600},
		{ID: "E3", Country: "CUB", EventType: "SUSPICIOUS_ACCESS", Severity: 0.6, Timestamp: now + 7200},
	}

	results := c.AnalyzeGlobalThreats(events)
	require.NotEmpty(t, results)
	assert.Equal(t, "PRIVILEGE_ESCALATION", results[0].Technique)
	assert.False(t, results[0].IsCoordinated)
}

func TestDetectCoordinatedAttack_NotEnoughEvents(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	assert.False(t, c.DetectCoordinatedAttack())
}

func TestDetectCoordinatedAttack_Detected(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	now := time.Now().Unix()
	for i := 0; i < 12; i++ {
		country := []string{"HTI", "DOM", "CUB", "JAM"}[i%4]
		eventType := []string{"LOGIN_FAILURE", "DATA_EXFILTRATION"}[i%2]
		c.IngestEvent(SecurityEvent{
			ID: "", Country: country, EventType: eventType, Timestamp: now + int64(i*10),
		})
	}
	assert.True(t, c.DetectCoordinatedAttack())
}

func TestDetectCoordinatedAttack_SingleCountryDominates(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	now := time.Now().Unix()
	for i := 0; i < 10; i++ {
		c.IngestEvent(SecurityEvent{
			Country: "HTI", EventType: "LOGIN_FAILURE", Timestamp: now + int64(i*10),
		})
	}
	assert.False(t, c.DetectCoordinatedAttack())
}

func TestClassifyTechnique_KnownPattern(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	assert.Equal(t, "CREDENTIAL_BRUTE_FORCE", c.classifyTechnique("LOGIN_FAILURE", nil))
	assert.Equal(t, "VEHICLE_THEFT", c.classifyTechnique("LAPI_HIT", nil))
	assert.Equal(t, "IDENTITY_SPOOFING", c.classifyTechnique("BIOMETRIC_MISMATCH", nil))
}

func TestClassifyTechnique_UnknownPattern_HighSeverity(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	events := []SecurityEvent{
		{Severity: 0.9}, {Severity: 0.95}, {Severity: 0.4},
	}
	assert.Equal(t, "TARGETED_ATTACK", c.classifyTechnique("UNKNOWN_TYPE", events))
}

func TestClassifyTechnique_UnknownPattern_LowSeverity(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	events := []SecurityEvent{
		{Severity: 0.3}, {Severity: 0.4}, {Severity: 0.5},
	}
	assert.Equal(t, "UNCLASSIFIED", c.classifyTechnique("UNKNOWN_TYPE", events))
}

func TestGetStats(t *testing.T) {
	c := NewCorrelationLayer("v2.0")
	c.IngestEvent(SecurityEvent{Country: "HTI", EventType: "LOGIN_FAILURE"})
	c.IngestEvent(SecurityEvent{Country: "DOM", EventType: "DATA_EXFILTRATION"})

	stats := c.GetStats()
	assert.Equal(t, 2, stats["total_events"])
	assert.Equal(t, "v2.0", stats["model_version"])
	assert.Equal(t, 15.0, stats["window_size_minutes"])
}

func TestCorrelationSortOrder(t *testing.T) {
	c := NewCorrelationLayer("v1.0")
	now := time.Now().Unix()

	results := c.AnalyzeGlobalThreats([]SecurityEvent{
		{Country: "HTI", EventType: "HIGH_IMPACT", Severity: 0.9, Timestamp: now},
		{Country: "DOM", EventType: "HIGH_IMPACT", Severity: 0.9, Timestamp: now + 30},
		{Country: "CUB", EventType: "HIGH_IMPACT", Severity: 0.9, Timestamp: now + 60},
		{Country: "HTI", EventType: "LOW_IMPACT", Severity: 0.2, Timestamp: now},
		{Country: "DOM", EventType: "LOW_IMPACT", Severity: 0.2, Timestamp: now + 30},
	})
	require.Len(t, results, 2)
	assert.GreaterOrEqual(t, results[0].Score, results[1].Score)
}
