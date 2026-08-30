package soc

import (
	"testing"
)

func TestNewSeverityEngine(t *testing.T) {
	e := NewSeverityEngine()
	if e == nil {
		t.Fatal("NewSeverityEngine returned nil")
	}
}

func TestClassify_Critical(t *testing.T) {
	e := NewSeverityEngine()
	alert := map[string]interface{}{
		"fraudScore": 95,
	}
	sev := e.Classify(alert)
	if sev != SeverityCritical {
		t.Errorf("Severity = %s, want CRITICAL", sev)
	}
}

func TestClassify_High(t *testing.T) {
	e := NewSeverityEngine()
	alert := map[string]interface{}{
		"fraudScore": 75,
	}
	sev := e.Classify(alert)
	if sev != SeverityHigh {
		t.Errorf("Severity = %s, want HIGH", sev)
	}
}

func TestClassify_Medium(t *testing.T) {
	e := NewSeverityEngine()
	alert := map[string]interface{}{
		"fraudScore": 50,
	}
	sev := e.Classify(alert)
	if sev != SeverityMedium {
		t.Errorf("Severity = %s, want MEDIUM", sev)
	}
}

func TestClassify_Low(t *testing.T) {
	e := NewSeverityEngine()
	alert := map[string]interface{}{
		"fraudScore": 10,
	}
	sev := e.Classify(alert)
	if sev != SeverityLow {
		t.Errorf("Severity = %s, want LOW", sev)
	}
}

func TestClassify_ScoreFromMetadata(t *testing.T) {
	e := NewSeverityEngine()
	alert := map[string]interface{}{
		"score": float64(85),
	}
	sev := e.Classify(alert)
	if sev != SeverityHigh {
		t.Errorf("Severity = %s, want HIGH", sev)
	}
}

func TestClassify_StatusCritical(t *testing.T) {
	e := NewSeverityEngine()
	alert := map[string]interface{}{
		"fraudScore": 0,
		"status":     "CRITICAL",
	}
	sev := e.Classify(alert)
	if sev != SeverityCritical {
		t.Errorf("Severity = %s, want CRITICAL", sev)
	}
}

func TestClassify_EmptyAlert(t *testing.T) {
	e := NewSeverityEngine()
	sev := e.Classify(map[string]interface{}{})
	if sev != SeverityLow {
		t.Errorf("Severity = %s, want LOW", sev)
	}
}

func TestGetEscalationPath_Critical(t *testing.T) {
	e := NewSeverityEngine()
	path := e.GetEscalationPath(SeverityCritical)
	if len(path) != 3 {
		t.Errorf("Path length = %d, want 3", len(path))
	}
	if path[0] != "pagerduty" {
		t.Errorf("path[0] = %s, want pagerduty", path[0])
	}
}

func TestGetEscalationPath_High(t *testing.T) {
	e := NewSeverityEngine()
	path := e.GetEscalationPath(SeverityHigh)
	if len(path) != 2 {
		t.Errorf("Path length = %d, want 2", len(path))
	}
}

func TestGetEscalationPath_Low(t *testing.T) {
	e := NewSeverityEngine()
	path := e.GetEscalationPath(SeverityLow)
	if len(path) != 1 {
		t.Errorf("Path length = %d, want 1", len(path))
	}
	if path[0] != "audit" {
		t.Errorf("path[0] = %s, want audit", path[0])
	}
}

func TestSeverityConstants(t *testing.T) {
	if SeverityLow != "LOW" {
		t.Errorf("SeverityLow = %s", SeverityLow)
	}
	if SeverityCritical != "CRITICAL" {
		t.Errorf("SeverityCritical = %s", SeverityCritical)
	}
}
