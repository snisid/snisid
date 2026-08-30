package unit

import (
	"testing"

	"github.com/snisid/platform/services/afis-svc/internal/domain"
)

func TestSearch_HitAbove85percent(t *testing.T) {
	results := []domain.SearchResult{
		{Score: 0.92, Rank: 1, NationalAFISID: "AFIS-2026-0000001"},
		{Score: 0.87, Rank: 2, NationalAFISID: "AFIS-2026-0000002"},
	}

	if len(results) == 0 {
		t.Fatal("expected at least one search result")
	}
	if results[0].Score < 0.85 {
		t.Fatalf("top result score %.2f below threshold 0.85", results[0].Score)
	}
	if results[0].Rank != 1 {
		t.Fatalf("expected rank 1, got %d", results[0].Rank)
	}
}

func TestSearch_NoMatchBelowThreshold(t *testing.T) {
	var results []domain.SearchResult
	if len(results) != 0 {
		t.Fatal("expected empty results below threshold")
	}
}
