package engine

import (
	"testing"
)

func makeParentProfile() STRLoci {
	p := STRLoci{}
	for locus, v := range profileFullMatch {
		p[locus] = v
	}
	return p
}

func makeChildProfile(parent STRLoci) STRLoci {
	// Child inherits one allele per locus from parent, gets one from other parent
	child := STRLoci{}
	for locus, v := range parent {
		alleles := []string{v.Value1, v.Value2}
		inherited := alleles[0]
		other := addAllele(alleles[1], 1)
		child[locus] = Locus{Value1: inherited, Value2: other}
	}
	return child
}

func makeSiblingProfile(parent STRLoci) STRLoci {
	// Sibling shares ~50% of alleles
	sib := STRLoci{}
	for locus, v := range parent {
		alleles := []string{v.Value1, v.Value2}
		shared := alleles[1]
		other := addAllele(alleles[0], 2)
		sib[locus] = Locus{Value1: shared, Value2: other}
	}
	return sib
}

func TestFamilialSearch_ParentChild(t *testing.T) {
	m := NewMatcher()
	parent := makeParentProfile()
	child := makeChildProfile(parent)

	results := m.FamilialSearch(parent, []STRLoci{child}, []string{"child-001"})
	if len(results) == 0 {
		t.Fatal("expected at least one familial match")
	}
	r := results[0]
	if r.LikelyRelation != "PARENT" {
		t.Fatalf("expected PARENT relation, got %s", r.LikelyRelation)
	}
	if r.Score < 0.90 {
		t.Fatalf("expected score >= 0.90 for parent-child, got %.4f", r.Score)
	}
}

func TestFamilialSearch_Sibling(t *testing.T) {
	m := NewMatcher()
	parent := makeParentProfile()
	sibling := makeSiblingProfile(parent)

	results := m.FamilialSearch(parent, []STRLoci{sibling}, []string{"sib-001"})
	if len(results) == 0 {
		t.Fatal("expected at least one familial match")
	}
	r := results[0]
	if r.LikelyRelation != "SIBLING" && r.LikelyRelation != "PARENT" {
		t.Fatalf("expected SIBLING or PARENT, got %s", r.LikelyRelation)
	}
	if r.Score < 0.75 {
		t.Fatalf("expected score >= 0.75 for sibling, got %.4f", r.Score)
	}
}

func TestFamilialSearch_Cousin(t *testing.T) {
	m := NewMatcher()
	// Cousin: ~12.5% shared → each locus has ~25% chance of sharing one allele
	// Build a profile with ~25% allele overlap
	query := profileFullMatch
	cousin := STRLoci{}
	i := 0
	for locus, v := range query {
		if i%4 == 0 {
			// share one allele
			cousin[locus] = Locus{Value1: v.Value1, Value2: addAllele(v.Value2, 5)}
		} else {
			// share none
			cousin[locus] = Locus{Value1: addAllele(v.Value1, 3), Value2: addAllele(v.Value2, 4)}
		}
		i++
	}

	results := m.FamilialSearch(query, []STRLoci{cousin}, []string{"cousin-001"})
	if len(results) == 0 {
		t.Fatal("expected at least one familial match for cousin-level sharing")
	}
	r := results[0]
	if r.LikelyRelation != "COUSIN" && r.LikelyRelation != "DISTANT" {
		t.Fatalf("expected COUSIN or DISTANT, got %s", r.LikelyRelation)
	}
}

func TestFamilialSearch_BelowThreshold(t *testing.T) {
	m := NewMatcher()
	query := profileFullMatch
	// 2/20 loci share (score = 0.05) → below 0.25 threshold
	unrelated := STRLoci{}
	i := 0
	for locus, v := range query {
		if i < 2 {
			unrelated[locus] = v
		} else {
			unrelated[locus] = Locus{Value1: addAllele(v.Value1, 5), Value2: addAllele(v.Value2, 5)}
		}
		i++
	}

	results := m.FamilialSearch(query, []STRLoci{unrelated}, []string{"unrelated-001"})
	if len(results) != 0 {
		t.Fatalf("expected 0 results below threshold, got %d", len(results))
	}
}

func TestFamilialSearch_EmptyCandidates(t *testing.T) {
	m := NewMatcher()
	results := m.FamilialSearch(profileFullMatch, []STRLoci{}, []string{})
	if len(results) != 0 {
		t.Fatalf("expected 0 results for empty candidates, got %d", len(results))
	}
}

func TestFamilialSearch_TruncateAt50(t *testing.T) {
	m := NewMatcher()
	child := makeChildProfile(profileFullMatch)
	candidates := make([]STRLoci, 100)
	ids := make([]string, 100)
	for i := range candidates {
		candidates[i] = child
		ids[i] = "repeat"
	}
	results := m.FamilialSearch(profileFullMatch, candidates, ids)
	if len(results) > 50 {
		t.Fatalf("expected max 50 results, got %d", len(results))
	}
}

func TestCountSharedAlleles(t *testing.T) {
	a := STRLoci{"CSF1PO": {Value1: "10", Value2: "12"}}
	b := STRLoci{"CSF1PO": {Value1: "10", Value2: "14"}}
	shared, total := countSharedAlleles(a, b)
	if total != 2 {
		t.Fatalf("expected 2 total alleles, got %d", total)
	}
	if shared != 1 {
		t.Fatalf("expected 1 shared allele (10), got %d", shared)
	}
}

func TestCountSharedAlleles_MissingLocus(t *testing.T) {
	a := STRLoci{"CSF1PO": {Value1: "10", Value2: "12"}}
	b := STRLoci{"D3S1358": {Value1: "15", Value2: "18"}}
	shared, total := countSharedAlleles(a, b)
	if total != 0 {
		t.Fatalf("expected 0 total (no shared loci), got %d", total)
	}
	if shared != 0 {
		t.Fatalf("expected 0 shared, got %d", shared)
	}
}

func TestClassifyFamilial(t *testing.T) {
	tests := []struct {
		score  float64
		shared int
		total  int
		want   string
	}{
		{0.95, 38, 40, "PARENT"},
		{0.80, 32, 40, "SIBLING"},
		{0.60, 24, 40, "HALF_SIBLING"},
		{0.40, 16, 40, "GRANDPARENT"},
		{0.36, 14, 40, "UNCLE"},
		{0.28, 11, 40, "COUSIN"},
		{0.10, 4, 40, "DISTANT"},
	}
	for _, tc := range tests {
		got := classifyFamilial(tc.score, tc.shared, tc.total)
		if got != tc.want {
			t.Errorf("classifyFamilial(%.2f, %d, %d) = %s, want %s",
				tc.score, tc.shared, tc.total, got, tc.want)
		}
	}
}

func TestClassifyFamilial_Unknown(t *testing.T) {
	got := classifyFamilial(0.5, 10, 0)
	if got != "UNKNOWN" {
		t.Fatalf("expected UNKNOWN for zero total, got %s", got)
	}
}
