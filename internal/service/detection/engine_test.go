package detection

import (
	"context"
	"testing"
)

func TestNewDetectionEngine(t *testing.T) {
	e := NewDetectionEngine()
	if e == nil {
		t.Fatal("NewDetectionEngine returned nil")
	}
	if len(e.strategies) != 2 {
		t.Errorf("Strategy count = %d, want 2", len(e.strategies))
	}
}

func TestDetect_ExactMatch(t *testing.T) {
	e := NewDetectionEngine()
	newID := map[string]interface{}{
		"fullName": "Jean Dupont",
		"taxId":    "123-45-6789",
	}
	candidates := []map[string]interface{}{
		{"fullName": "Jean Dupont", "identityId": "ID-001", "taxId": "123-45-6789"},
	}

	score, matchID, evidence := e.Detect(context.Background(), newID, candidates)
	if score != 100 {
		t.Errorf("Score = %d, want 100 (exact taxId match)", score)
	}
	if matchID != "ID-001" {
		t.Errorf("MatchID = %s, want ID-001", matchID)
	}
	if len(evidence) == 0 {
		t.Error("Evidence should not be empty")
	}
}

func TestDetect_PartialMatch(t *testing.T) {
	e := NewDetectionEngine()
	newID := map[string]interface{}{
		"fullName": "John Smith",
	}
	candidates := []map[string]interface{}{
		{"fullName": "John Smyth", "identityId": "ID-002"},
		{"fullName": "Jane Doe", "identityId": "ID-003"},
	}

	score, matchID, _ := e.Detect(context.Background(), newID, candidates)
	if score <= 0 {
		t.Error("Should have a positive match score for similar names")
	}
	if matchID != "ID-002" {
		t.Errorf("MatchID = %s, want ID-002 (closest match)", matchID)
	}
}

func TestDetect_NoMatch(t *testing.T) {
	e := NewDetectionEngine()
	newID := map[string]interface{}{
		"fullName": "Unknown Person",
	}
	candidates := []map[string]interface{}{
		{"fullName": "Alice Wonderland", "identityId": "ID-004"},
		{"fullName": "Bob Builder", "identityId": "ID-005"},
	}

	score, matchID, _ := e.Detect(context.Background(), newID, candidates)
	if matchID == "" {
		t.Log("No close match found, best score =", score)
	}
}

func TestDetect_EmptyCandidates(t *testing.T) {
	e := NewDetectionEngine()
	newID := map[string]interface{}{
		"fullName": "Test Person",
	}

	score, matchID, _ := e.Detect(context.Background(), newID, []map[string]interface{}{})
	if score != 0 {
		t.Errorf("Score = %d, want 0", score)
	}
	if matchID != "" {
		t.Errorf("MatchID = %s, want empty", matchID)
	}
}

func TestDetect_MultipleStrategies(t *testing.T) {
	e := NewDetectionEngine()
	if len(e.strategies) < 2 {
		t.Error("Should have at least 2 detection strategies")
	}

	names := []string{
		"fuzzy_jaro_winkler",
		"phonetic_soundex",
	}
	for i, s := range e.strategies {
		if s.Name() != names[i] {
			t.Errorf("Strategy %d name = %s, want %s", i, s.Name(), names[i])
		}
	}
}
