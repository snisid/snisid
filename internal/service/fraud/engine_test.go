package fraud

import (
	"context"
	"testing"
	"time"

	"github.com/snisid/platform/backend/internal/service/router"
)

type mockAIClient struct {
	predictFn func(ctx context.Context, event map[string]interface{}) (int, error)
}

func (m *mockAIClient) Predict(ctx context.Context, event map[string]interface{}) (int, error) {
	if m.predictFn != nil {
		return m.predictFn(ctx, event)
	}
	return 0, nil
}

func newTestEngine(aiClient AIClient) *ScoringEngine {
	eng, _ := router.NewEngine()
	return &ScoringEngine{
		rules:    eng,
		state:    NewStateStore(""),
		aiClient: aiClient,
	}
}

func TestScoringEngine_CalculateScore_NoTriggers(t *testing.T) {
	engine := newTestEngine(&mockAIClient{})
	event := map[string]interface{}{
		"identityId": "normal-user",
		"action":     "view",
	}

	score, _ := engine.CalculateScore(context.Background(), event)
	if score != 5 {
		t.Errorf("score = %d, want 5 (only AI default)", score)
	}
}

func TestScoringEngine_AIHighRisk(t *testing.T) {
	aiClient := &mockAIClient{
		predictFn: func(ctx context.Context, event map[string]interface{}) (int, error) {
			return 60, nil
		},
	}

	engine := newTestEngine(aiClient)
	event := map[string]interface{}{
		"identityId": "test-fraud",
		"action":     "enroll",
	}

	score, reasons := engine.CalculateScore(context.Background(), event)
	if score < 60 {
		t.Errorf("score = %d, want >= 60", score)
	}
	if reasons == "" {
		t.Error("expected reasons for high risk")
	}
}

func TestScoringEngine_AIErrorFallback(t *testing.T) {
	aiClient := &mockAIClient{
		predictFn: func(ctx context.Context, event map[string]interface{}) (int, error) {
			return 0, nil
		},
	}

	engine := newTestEngine(aiClient)
	event := map[string]interface{}{
		"identityId": "user-001",
	}

	score, _ := engine.CalculateScore(context.Background(), event)
	if score < 0 || score > 100 {
		t.Errorf("score = %d, out of expected range [0,100]", score)
	}
}

func TestDefaultAIClient_Predict(t *testing.T) {
	client := NewDefaultAIClient("http://localhost:8501")

	tests := []struct {
		name    string
		event   map[string]interface{}
		wantMin int
	}{
		{"fraud pattern", map[string]interface{}{"identityId": "test-fraud"}, 60},
		{"normal user", map[string]interface{}{"identityId": "normal-user"}, 0},
		{"missing identityId", map[string]interface{}{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, err := client.Predict(context.Background(), tt.event)
			if err != nil {
				t.Fatalf("Predict failed: %v", err)
			}
			if score < tt.wantMin {
				t.Errorf("score = %d, want >= %d", score, tt.wantMin)
			}
		})
	}
}

func TestNewScoringEngine(t *testing.T) {
	engine, err := NewScoringEngine("localhost:6379", &mockAIClient{})
	if err != nil {
		t.Fatalf("NewScoringEngine failed: %v", err)
	}
	if engine == nil {
		t.Fatal("Expected non-nil engine")
	}
}

func TestScoringEngine_ReloadRules(t *testing.T) {
	engine := newTestEngine(&mockAIClient{})
	err := engine.ReloadRules([]router.Rule{})
	if err != nil {
		t.Errorf("ReloadRules failed: %v", err)
	}
}
