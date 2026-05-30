package detection

import (
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/dotcypress/phonetic"
)

type MatchStrategy interface {
	Name() string
	Score(s1, s2 string) int
}

type FuzzyMatcher struct {
	metric strutil.StringMetric
}

func NewFuzzyMatcher() *FuzzyMatcher {
	return &FuzzyMatcher{
		metric: metrics.NewJaroWinkler(),
	}
}

func (m *FuzzyMatcher) Name() string { return "fuzzy_jaro_winkler" }
func (m *FuzzyMatcher) Score(s1, s2 string) int {
	score := m.metric.Compare(s1, s2)
	return int(score * 100)
}

type PhoneticMatcher struct{}

func NewPhoneticMatcher() *PhoneticMatcher {
	return &PhoneticMatcher{}
}

func (m *PhoneticMatcher) Name() string { return "phonetic_metaphone" }
func (m *PhoneticMatcher) Score(s1, s2 string) int {
	p := phonetic.NewMetaphone()
	if p.Encode(s1) == p.Encode(s2) {
		return 100
	}
	return 0
}
