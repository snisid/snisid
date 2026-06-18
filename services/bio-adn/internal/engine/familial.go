package engine

import (
	"math"
	"sort"
)

type FamilialMatch struct {
	SampleID       string  `json:"sample_id"`
	Score          float64 `json:"score"`
	SharedAlleles  int     `json:"shared_alleles"`
	TotalAlleles   int     `json:"total_alleles"`
	LikelyRelation string  `json:"likely_relation"` // PARENT, SIBLING, HALF_SIBLING, COUSIN, UNCLE, GRANDPARENT
}

type AlleleCount struct {
	Count      int
	ObservedIn int
}

func (m *Matcher) FamilialSearch(query STRLoci, candidates []STRLoci, candidateIDs []string) []FamilialMatch {
	var results []FamilialMatch

	for i, candidate := range candidates {
		shared, total := countSharedAlleles(query, candidate)
		if total == 0 {
			continue
		}
		score := float64(shared) / float64(total)
		if score < 0.25 {
			continue
		}
		relation := classifyFamilial(score, shared, total)
		results = append(results, FamilialMatch{
			SampleID:       candidateIDs[i],
			Score:          math.Round(score*10000) / 10000,
			SharedAlleles:  shared,
			TotalAlleles:   total,
			LikelyRelation: relation,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) > 50 {
		results = results[:50]
	}
	return results
}

func countSharedAlleles(query, candidate STRLoci) (shared, total int) {
	for locus, qv := range query {
		cv, ok := candidate[locus]
		if !ok {
			continue
		}
		total += 2
		if qv.Value1 == cv.Value1 || qv.Value1 == cv.Value2 {
			shared++
		}
		if qv.Value2 == cv.Value1 || qv.Value2 == cv.Value2 {
			shared++
		}
	}
	return
}

func classifyFamilial(score float64, shared, total int) string {
	if total == 0 {
		return "UNKNOWN"
	}
	switch {
	case score >= 0.90 && shared >= 35:
		return "PARENT"
	case score >= 0.75 && shared >= 25:
		return "SIBLING"
	case score >= 0.50 && shared >= 17:
		return "HALF_SIBLING"
	case score >= 0.35 && shared >= 12:
		return "GRANDPARENT"
	case score >= 0.35:
		return "UNCLE"
	case score >= 0.25:
		return "COUSIN"
	default:
		return "DISTANT"
	}
}
