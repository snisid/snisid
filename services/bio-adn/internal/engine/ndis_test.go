package engine

import (
	"context"
	"testing"

	"github.com/snisid/platform/services/bio-adn/pkg/models"
)

type mockNDISDB struct {
	models.Database
}

func (m *mockNDISDB) SearchDNAProfiles(ctx context.Context, indexType string, limit, offset int) ([]models.DNAProfile, int, error) {
	return []models.DNAProfile{
		{SampleID: "MATCH-001", LociHash: "hash_abcdef", IndexType: indexType, LabID: "PAP"},
		{SampleID: "MATCH-002", LociHash: "hash_123456", IndexType: indexType, LabID: "CAP"},
	}, 2, nil
}

func TestNDISMatcher_CrossDept(t *testing.T) {
	db := &mockNDISDB{}
	matcher := NewNDISMatcher(db)
	result, err := matcher.MatchCrossDept(context.Background(), "QUERY-001", "hash_abcdef", "BIO-FSC", "SDIS-ART")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected a match, got nil")
	}
	if result.MatchSDIS != "SDIS-PAP" {
		t.Fatalf("expected SDIS-PAP, got %s", result.MatchSDIS)
	}
}

func TestNDISMatcher_NoMatch(t *testing.T) {
	db := &mockNDISDB{}
	matcher := NewNDISMatcher(db)
	result, err := matcher.MatchCrossDept(context.Background(), "QUERY-002", "hash_nonexistent", "BIO-FSC", "SDIS-ART")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected no match")
	}
}
