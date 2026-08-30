package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/snisid/platform/services/bio-adn/pkg/models"
)

const (
	MatchTypeFull     = "FULL_MATCH"
	MatchTypePartial  = "PARTIAL"
	MatchTypeFamilial = "FAMILIAL"
	MatchTypeNoMatch  = "NO_MATCH"
)

type NDISMatcher struct {
	db      models.Database
	matcher *Matcher
}

func NewNDISMatcher(db models.Database) *NDISMatcher {
	return &NDISMatcher{
		db:      db,
		matcher: NewMatcher(),
	}
}

type CrossDeptMatchResult struct {
	HitID         string
	QuerySampleID string
	MatchSampleID string
	MatchType     string
	Confidence    float64
	QuerySDIS     string
	MatchSDIS     string
	AlertLevel    string
}

func (m *NDISMatcher) MatchCrossDept(ctx context.Context, sampleID string, lociHash string, indexType string, querySDIS string) (*CrossDeptMatchResult, error) {
	indexTypes := []string{"BIO-CON", "BIO-ARR", "BIO-FSC"}
	if indexType == "BIO-FSC" {
		indexTypes = []string{"BIO-CON", "BIO-ARR"}
	}

	for _, targetIndex := range indexTypes {
		profiles, total, err := m.db.SearchDNAProfiles(ctx, targetIndex, 100, 0)
		if err != nil || total == 0 {
			continue
		}

		for _, profile := range profiles {
			if profile.LociHash != lociHash {
				continue
			}
			sdisCode := extractSDIS(profile.LabID)
			return &CrossDeptMatchResult{
				HitID:         fmt.Sprintf("NDIS-%d", time.Now().UnixNano()),
				QuerySampleID: sampleID,
				MatchSampleID: profile.SampleID,
				MatchType:     MatchTypeFull,
				Confidence:    1.0,
				QuerySDIS:     querySDIS,
				MatchSDIS:     sdisCode,
				AlertLevel:    "CRITICAL",
			}, nil
		}
	}

	return nil, nil
}

func (m *NDISMatcher) WeeklyStats(ctx context.Context) (*models.NdisStats, error) {
	return m.db.GetNdisStats(ctx)
}

func extractSDIS(labID string) string {
	if len(labID) >= 4 {
		return "SDIS-" + labID[:3]
	}
	return "SDIS-UNK"
}
