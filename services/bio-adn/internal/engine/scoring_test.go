package engine

import (
	"math"
	"testing"
)

var profileFullMatch = STRLoci{
	"CSF1PO":   {Value1: "10", Value2: "12"},
	"D3S1358":  {Value1: "15", Value2: "17"},
	"D5S818":   {Value1: "11", Value2: "13"},
	"D7S820":   {Value1: "8", Value2: "10"},
	"D8S1179":  {Value1: "13", Value2: "14"},
	"D13S317":  {Value1: "9", Value2: "11"},
	"D16S539":  {Value1: "11", Value2: "12"},
	"D18S51":   {Value1: "14", Value2: "16"},
	"D21S11":   {Value1: "29", Value2: "30"},
	"FGA":      {Value1: "21", Value2: "23"},
	"TH01":     {Value1: "7", Value2: "9"},
	"TPOX":     {Value1: "8", Value2: "11"},
	"vWA":      {Value1: "16", Value2: "17"},
	"D1S1656":  {Value1: "13", Value2: "15"},
	"D2S441":   {Value1: "10", Value2: "11"},
	"D2S1338":  {Value1: "19", Value2: "23"},
	"D10S1248": {Value1: "13", Value2: "15"},
	"D12S391":  {Value1: "18", Value2: "20"},
	"D19S433":  {Value1: "13", Value2: "14"},
	"D22S1045": {Value1: "15", Value2: "16"},
}

var profilePartialMatch = STRLoci{
	"CSF1PO":   {Value1: "10", Value2: "12"},
	"D3S1358":  {Value1: "14", Value2: "16"},
	"D5S818":   {Value1: "11", Value2: "13"},
	"D7S820":   {Value1: "8", Value2: "10"},
	"D8S1179":  {Value1: "13", Value2: "14"},
	"D13S317":  {Value1: "9", Value2: "11"},
	"D16S539":  {Value1: "11", Value2: "12"},
	"D18S51":   {Value1: "14", Value2: "16"},
	"D21S11":   {Value1: "29", Value2: "30"},
	"FGA":      {Value1: "20", Value2: "24"},
	"TH01":     {Value1: "7", Value2: "9"},
	"TPOX":     {Value1: "8", Value2: "11"},
	"vWA":      {Value1: "16", Value2: "17"},
	"D1S1656":  {Value1: "13", Value2: "15"},
	"D2S441":   {Value1: "10", Value2: "11"},
	"D2S1338":  {Value1: "19", Value2: "23"},
	"D10S1248": {Value1: "13", Value2: "15"},
	"D12S391":  {Value1: "18", Value2: "20"},
	"D19S433":  {Value1: "13", Value2: "14"},
	"D22S1045": {Value1: "15", Value2: "16"},
}

var profileNoMatch = STRLoci{}
var profileIncomplete = STRLoci{
	"CSF1PO":  {Value1: "10", Value2: "12"},
	"D3S1358": {Value1: "15", Value2: "17"},
	"D5S818":  {Value1: "11", Value2: "13"},
	"D7S820":  {Value1: "8", Value2: "10"},
}

func init() {
	for k, v := range profileFullMatch {
		profileNoMatch[k] = Locus{
			Value1: addAllele(v.Value1, 3),
			Value2: addAllele(v.Value2, 2),
		}
	}
}

func addAllele(a string, delta int) string {
	var val int
	for _, c := range a {
		if c >= '0' && c <= '9' {
			val = val*10 + int(c-'0')
		}
	}
	return itoa(val + delta)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [8]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

func TestCalculateLR_FullMatch(t *testing.T) {
	m := NewMatcher()
	r := m.CalculateLR(profileFullMatch, profileFullMatch, nil)
	if r == nil {
		t.Fatal("expected result")
	}
	if r.MatchType != "FULL_MATCH" {
		t.Fatalf("expected FULL_MATCH, got %s", r.MatchType)
	}
	if r.TotalLoci != 40 {
		t.Fatalf("expected 40 total alleles, got %d", r.TotalLoci)
	}
	if r.MatchedLoci != 40 {
		t.Fatalf("expected 40 matched alleles, got %d", r.MatchedLoci)
	}
	if r.LikelihoodRatio < 1e10 {
		t.Fatalf("expected LR >= 1e10 for full match, got %.2f", r.LikelihoodRatio)
	}
}

func TestCalculateLR_PartialMatch(t *testing.T) {
	m := NewMatcher()
	r := m.CalculateLR(profileFullMatch, profilePartialMatch, nil)
	if r == nil {
		t.Fatal("expected result")
	}
	// D3S1358 and FGA differ → 36 matched / 40 total
	if r.MatchedLoci != 36 {
		t.Fatalf("expected 36 matched alleles (2 loci differ), got %d", r.MatchedLoci)
	}
	logLR := math.Log10(r.LikelihoodRatio)
	if logLR < 3 || logLR >= 10 {
		t.Fatalf("expected partial match logLR between 3 and 10, got %.2f", logLR)
	}
}

func TestCalculateLR_NoMatch(t *testing.T) {
	m := NewMatcher()
	r := m.CalculateLR(profileFullMatch, profileNoMatch, nil)
	if r == nil {
		t.Fatal("expected result")
	}
	if r.MatchType != "NO_MATCH" {
		t.Fatalf("expected NO_MATCH, got %s", r.MatchType)
	}
	if r.LikelihoodRatio != 1.0 {
		t.Fatalf("expected LR clamped to 1.0 for no match, got %.2f", r.LikelihoodRatio)
	}
}

func TestCalculateLR_IncompleteProfile(t *testing.T) {
	m := NewMatcher()
	// query is full, candidate has only 4 loci (8 alleles)
	r := m.CalculateLR(profileFullMatch, profileIncomplete, nil)
	if r == nil {
		t.Fatal("expected result")
	}
	if r.TotalLoci != 8 {
		t.Fatalf("expected 8 total alleles (4 loci * 2), got %d", r.TotalLoci)
	}
	if r.MatchedLoci < 7 {
		t.Fatalf("expected at least 7 matched, got %d", r.MatchedLoci)
	}
}

func TestCalculateLR_NilOnEmptyQuery(t *testing.T) {
	m := NewMatcher()
	r := m.CalculateLR(STRLoci{}, profileFullMatch, nil)
	if r != nil {
		t.Fatal("expected nil for empty query")
	}
}

func TestCalculateLR_KnownAlleleFrequencies(t *testing.T) {
	m := NewMatcher()
	// CSF1PO 10/12 → use known freq: 10=0.220, 12=0.350
	query := STRLoci{"CSF1PO": {Value1: "10", Value2: "12"}}
	r := m.CalculateLR(query, query, nil)
	if r == nil {
		t.Fatal("expected result")
	}
	expectedLR := 1.0 / (2 * 0.220 * 0.350)
	if r.LikelihoodRatio < expectedLR*0.9 || r.LikelihoodRatio > expectedLR*1.1 {
		t.Fatalf("expected LR ~%.2f, got %.2f", expectedLR, r.LikelihoodRatio)
	}
}

func TestRankByLR(t *testing.T) {
	m := NewMatcher()
	candidates := []STRLoci{profileFullMatch, profilePartialMatch, profileNoMatch}
	ids := []string{"full-001", "partial-001", "no-match-001"}
	results := m.RankByLR(profileFullMatch, candidates, ids, nil)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	// Sorted descending by LR
	if results[0].SampleID != "full-001" {
		t.Fatalf("expected full-001 ranked first, got %s", results[0].SampleID)
	}
	if results[2].SampleID != "no-match-001" {
		t.Fatalf("expected no-match-001 ranked last, got %s", results[2].SampleID)
	}
}

func TestRankByLR_TruncateAt100(t *testing.T) {
	m := NewMatcher()
	candidates := make([]STRLoci, 150)
	ids := make([]string, 150)
	for i := range candidates {
		candidates[i] = profileFullMatch
		ids[i] = "id"
	}
	results := m.RankByLR(profileFullMatch, candidates, ids, nil)
	if len(results) > 100 {
		t.Fatalf("expected max 100 results, got %d", len(results))
	}
}

func TestCalculateLRFromProfiles(t *testing.T) {
	m := NewMatcher()
	lr := m.CalculateLRFromProfiles(profileFullMatch, profileFullMatch, nil)
	if lr < 1e10 {
		t.Fatalf("expected high LR for full match, got %.2f", lr)
	}
	lrNoMatch := m.CalculateLRFromProfiles(profileFullMatch, profileNoMatch, nil)
	if lrNoMatch != 1.0 {
		t.Fatalf("expected LR 1.0 for no match, got %.2f", lrNoMatch)
	}
}

func TestDefaultAlleleFrequenciesHaveAllLoci(t *testing.T) {
	expectedLoci := []string{
		"CSF1PO", "D3S1358", "D5S818", "D7S820", "D8S1179",
		"D13S317", "D16S539", "D18S51", "D21S11", "FGA",
		"TH01", "TPOX", "vWA", "D1S1656", "D2S441",
		"D2S1338", "D10S1248", "D12S391", "D19S433", "D22S1045",
	}
	for _, locus := range expectedLoci {
		if _, ok := DefaultAlleleFrequencies[locus]; !ok {
			t.Errorf("missing allele frequency for %s", locus)
		}
	}
}
