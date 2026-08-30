package detection

import (
	"testing"
)

func TestNewFuzzyMatcher(t *testing.T) {
	m := NewFuzzyMatcher()
	if m == nil {
		t.Fatal("NewFuzzyMatcher returned nil")
	}
	if m.Name() != "fuzzy_jaro_winkler" {
		t.Errorf("Name = %s, want fuzzy_jaro_winkler", m.Name())
	}
}

func TestFuzzyMatcher_ExactMatch(t *testing.T) {
	m := NewFuzzyMatcher()
	score := m.Score("Jean Dupont", "Jean Dupont")
	if score != 100 {
		t.Errorf("Score = %d, want 100", score)
	}
}

func TestFuzzyMatcher_Similar(t *testing.T) {
	m := NewFuzzyMatcher()
	score := m.Score("Jean Dupont", "Jean Dupond")
	if score < 80 {
		t.Errorf("Score = %d, want >= 80", score)
	}
	t.Logf("Jaro-Winkler similarity: %d", score)
}

func TestFuzzyMatcher_Different(t *testing.T) {
	m := NewFuzzyMatcher()
	score := m.Score("Alice", "Bob")
	if score > 50 {
		t.Errorf("Score = %d, want <= 50", score)
	}
}

func TestNewPhoneticMatcher(t *testing.T) {
	m := NewPhoneticMatcher()
	if m == nil {
		t.Fatal("NewPhoneticMatcher returned nil")
	}
	if m.Name() != "phonetic_soundex" {
		t.Errorf("Name = %s, want phonetic_soundex", m.Name())
	}
}

func TestPhoneticMatcher_SamePhonetic(t *testing.T) {
	m := NewPhoneticMatcher()
	// "Smith" and "Smyth" sound similar
	score := m.Score("Smith", "Smyth")
	if score != 100 {
		t.Errorf("Score = %d, want 100 (same phonetic encoding)", score)
	}
}

func TestPhoneticMatcher_DifferentPhonetic(t *testing.T) {
	m := NewPhoneticMatcher()
	score := m.Score("Jean", "Marie")
	if score != 0 {
		t.Errorf("Score = %d, want 0 (different phonetic encoding)", score)
	}
}

func TestFuzzyMatcher_EmptyStrings(t *testing.T) {
	m := NewFuzzyMatcher()
	score := m.Score("", "")
	if score != 100 {
		t.Errorf("Score = %d, want 100 (both empty)", score)
	}
}

func TestPhoneticMatcher_EmptyString(t *testing.T) {
	m := NewPhoneticMatcher()
	score := m.Score("", "test")
	if score != 0 {
		t.Errorf("Score = %d, want 0", score)
	}
}

func TestMatchStrategyInterface(t *testing.T) {
	var f MatchStrategy = NewFuzzyMatcher()
	var p MatchStrategy = NewPhoneticMatcher()

	if f.Name() == p.Name() {
		t.Error("Strategies should have different names")
	}
}
