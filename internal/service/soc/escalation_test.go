package soc

import (
	"context"
	"testing"
)

func TestNewEscalationManager(t *testing.T) {
	m := NewEscalationManager("https://hooks.slack.com/test", "pd-key-123")
	if m == nil {
		t.Fatal("NewEscalationManager returned nil")
	}
	if m.slackWebhook != "https://hooks.slack.com/test" {
		t.Errorf("SlackWebhook = %s", m.slackWebhook)
	}
	if m.pdServiceKey != "pd-key-123" {
		t.Errorf("PdServiceKey = %s", m.pdServiceKey)
	}
}

func TestEscalate_Critical(t *testing.T) {
	m := NewEscalationManager("slack-url", "pd-key")
	err := m.Escalate(context.Background(), "INC-001", SeverityCritical, "Critical security incident detected")
	if err != nil {
		t.Fatalf("Escalate failed: %v", err)
	}
}

func TestEscalate_High(t *testing.T) {
	m := NewEscalationManager("slack-url", "pd-key")
	err := m.Escalate(context.Background(), "INC-002", SeverityHigh, "High severity alert")
	if err != nil {
		t.Fatalf("Escalate failed: %v", err)
	}
}

func TestEscalate_Low(t *testing.T) {
	m := NewEscalationManager("slack-url", "pd-key")
	err := m.Escalate(context.Background(), "INC-003", SeverityLow, "Low severity")
	if err != nil {
		t.Fatalf("Escalate failed: %v", err)
	}
}

func TestEscalate_EmptyIncidentID(t *testing.T) {
	m := NewEscalationManager("slack-url", "pd-key")
	err := m.Escalate(context.Background(), "", SeverityHigh, "test")
	if err != nil {
		t.Fatalf("Escalate with empty ID failed: %v", err)
	}
}

func TestNewEscalationManager_EmptyWebhooks(t *testing.T) {
	m := NewEscalationManager("", "")
	if m == nil {
		t.Fatal("NewEscalationManager should handle empty webhooks")
	}
}

func TestEscalate_EmptyDescription(t *testing.T) {
	m := NewEscalationManager("slack-url", "pd-key")
	err := m.Escalate(context.Background(), "INC-004", SeverityCritical, "")
	if err != nil {
		t.Fatalf("Escalate with empty description failed: %v", err)
	}
}

func TestEscalate_MediumLevel(t *testing.T) {
	m := NewEscalationManager("slack-url", "pd-key")
	// Medium severity doesn't exist, defaults to High path via sev >= SeverityHigh
	err := m.Escalate(context.Background(), "INC-005", "MEDIUM", "Test")
	if err != nil {
		t.Fatalf("Escalate failed: %v", err)
	}
}
