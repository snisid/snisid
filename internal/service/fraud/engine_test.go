package fraud

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/snisid/platform/internal/service/router"
)

type mockAIClient struct {
	predictFn func(ctx context.Context, event map[string]interface{}) (int, error)
}

func (m *mockAIClient) Predict(ctx context.Context, event map[string]interface{}) (int, error) {
	if m.predictFn != nil {
		return m.predictFn(ctx, event)
	}
	return 5, nil // match DefaultAIClient baseline
}

func newTestEngine(aiClient AIClient) (*ScoringEngine, *miniredis.Miniredis) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	eng, _ := router.NewEngine()
	return &ScoringEngine{
		rules:    eng,
		state:    &StateStore{client: client},
		aiClient: aiClient,
	}, s
}

func TestScoringEngine_CalculateScore_NoTriggers(t *testing.T) {
	engine, mini := newTestEngine(&mockAIClient{})
	t.Cleanup(func() { mini.Close() })
	event := map[string]interface{}{
		"identityId": "normal-user",
		"action":     "view",
	}

	score, _, riskLevel := engine.CalculateScore(context.Background(), event)
	if score != 5 {
		t.Errorf("score = %d, want 5 (only AI default)", score)
	}
	if riskLevel == "" {
		t.Error("expected non-empty risk level")
	}
}

func TestScoringEngine_AIHighRisk(t *testing.T) {
	aiClient := &mockAIClient{
		predictFn: func(ctx context.Context, event map[string]interface{}) (int, error) {
			return 60, nil
		},
	}

	engine, mini := newTestEngine(aiClient)
	t.Cleanup(func() { mini.Close() })
	event := map[string]interface{}{
		"identityId": "test-fraud",
		"action":     "enroll",
	}

	score, reasons, riskLevel := engine.CalculateScore(context.Background(), event)
	if score < 60 {
		t.Errorf("score = %d, want >= 60", score)
	}
	if reasons == "" {
		t.Error("expected reasons for high risk")
	}
	if riskLevel == "" {
		t.Error("expected non-empty risk level")
	}
}

func TestScoringEngine_AIErrorFallback(t *testing.T) {
	aiClient := &mockAIClient{
		predictFn: func(ctx context.Context, event map[string]interface{}) (int, error) {
			return 0, nil
		},
	}

	engine, mini := newTestEngine(aiClient)
	t.Cleanup(func() { mini.Close() })
	event := map[string]interface{}{
		"identityId": "user-001",
	}

	score, _, riskLevel := engine.CalculateScore(context.Background(), event)
	if score < 0 || score > 100 {
		t.Errorf("score = %d, out of expected range [0,100]", score)
	}
	if riskLevel == "" {
		t.Error("expected non-empty risk level")
	}
}

func TestScoringEngine_IntelligenceFusion(t *testing.T) {
	engine, mini := newTestEngine(&mockAIClient{})
	t.Cleanup(func() { mini.Close() })
	event := map[string]interface{}{
		"identityId": "high-graph-risk",
		"graph_risk": 0.95,
	}

	score, reasons, riskLevel := engine.CalculateScore(context.Background(), event)
	if score < 5 {
		t.Errorf("score = %d, want >= 5", score)
	}
	if reasons == "" {
		t.Error("expected reasons from intelligence fusion")
	}
	if riskLevel != "CRITICAL" {
		t.Errorf("riskLevel = %s, want CRITICAL (graph_risk=0.95)", riskLevel)
	}
}

func TestScoringEngine_VelocityScore(t *testing.T) {
	engine, mini := newTestEngine(&mockAIClient{})
	t.Cleanup(func() { mini.Close() })
	ctx := context.Background()

	event := map[string]interface{}{
		"identityId": "vel-user",
		"action":     "enroll",
	}

	for i := 0; i < 6; i++ {
		engine.CalculateScore(ctx, event)
	}

	score, reasons, _ := engine.CalculateScore(ctx, event)
	if score < 50 {
		t.Errorf("score = %d, want >= 50 (velocity penalty)", score)
	}
	if reasons == "" {
		t.Error("expected reasons including velocity")
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
	engine, mini := newTestEngine(&mockAIClient{})
	t.Cleanup(func() { mini.Close() })
	err := engine.ReloadRules([]router.Rule{})
	if err != nil {
		t.Errorf("ReloadRules failed: %v", err)
	}
}
