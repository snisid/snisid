package security

import (
	"context"
	"testing"
)

func TestNewAdaptiveTrustScorer(t *testing.T) {
	s := NewAdaptiveTrustScorer()
	if s == nil {
		t.Fatal("NewAdaptiveTrustScorer returned nil")
	}
}

func TestCalculateScore_EmptySignals(t *testing.T) {
	s := NewAdaptiveTrustScorer()
	result := s.CalculateScore(context.Background(), "usr-001", []TrustSignal{})
	if result.Score != 0 {
		t.Errorf("Score = %d, want 0", result.Score)
	}
	if result.Level != "LOW" {
		t.Errorf("Level = %s, want LOW", result.Level)
	}
	if len(result.Signals) != 0 {
		t.Errorf("Signals count = %d, want 0", len(result.Signals))
	}
}

func TestCalculateScore_SingleSignal(t *testing.T) {
	s := NewAdaptiveTrustScorer()
	result := s.CalculateScore(context.Background(), "usr-001", []TrustSignal{
		{Type: "device_trust", Value: 90, Weight: 1.0},
	})
	if result.Score != 90 {
		t.Errorf("Score = %d, want 90", result.Score)
	}
	if result.Level != "HIGH" {
		t.Errorf("Level = %s, want HIGH", result.Level)
	}
	if len(result.Signals) != 1 {
		t.Errorf("Signals count = %d, want 1", len(result.Signals))
	}
}

func TestCalculateScore_MultipleSignals(t *testing.T) {
	s := NewAdaptiveTrustScorer()
	result := s.CalculateScore(context.Background(), "usr-002", []TrustSignal{
		{Type: "biometric_match", Value: 95, Weight: 0.5},
		{Type: "behavioral_consistency", Value: 70, Weight: 0.3},
		{Type: "location_trust", Value: 50, Weight: 0.2},
	})
	// Expected: 95*0.5 + 70*0.3 + 50*0.2 = 47.5 + 21 + 10 = 78.5 -> 78
	if result.Score != 78 {
		t.Errorf("Score = %d, want 78", result.Score)
	}
	if result.Level != "HIGH" {
		t.Errorf("Level = %s, want HIGH", result.Level)
	}
	if len(result.Signals) != 3 {
		t.Errorf("Signals count = %d, want 3", len(result.Signals))
	}
}

func TestCalculateScore_MediumLevel(t *testing.T) {
	s := NewAdaptiveTrustScorer()
	result := s.CalculateScore(context.Background(), "usr-003", []TrustSignal{
		{Type: "basic_auth", Value: 60, Weight: 1.0},
	})
	// Score 60 -> Level MEDIUM
	if result.Score != 60 {
		t.Errorf("Score = %d, want 60", result.Score)
	}
	if result.Level != "MEDIUM" {
		t.Errorf("Level = %s, want MEDIUM", result.Level)
	}
}

func TestCalculateScore_BoundaryLow(t *testing.T) {
	s := NewAdaptiveTrustScorer()
	result := s.CalculateScore(context.Background(), "usr-004", []TrustSignal{
		{Type: "anonymous", Value: 20, Weight: 1.0},
	})
	if result.Level != "LOW" {
		t.Errorf("Level = %s, want LOW (score = %d)", result.Level, result.Score)
	}
}

func TestCalculateScore_ZeroWeight(t *testing.T) {
	s := NewAdaptiveTrustScorer()
	result := s.CalculateScore(context.Background(), "usr-005", []TrustSignal{
		{Type: "irrelevant", Value: 100, Weight: 0.0},
	})
	if result.Score != 0 {
		t.Errorf("Score = %d, want 0 (weight = 0)", result.Score)
	}
}
