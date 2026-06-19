package detection

import (
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
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

func (m *PhoneticMatcher) Name() string { return "phonetic_soundex" }
func (m *PhoneticMatcher) Score(s1, s2 string) int {
	if soundex(s1) == soundex(s2) {
		return 100
	}
	return 0
}

// soundex is a simple replacement for the missing phonetic library
func soundex(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ToUpper(s)
	result := []byte{s[0]}
	lastCode := charToSoundex(s[0])
	for i := 1; i < len(s) && len(result) < 4; i++ {
		code := charToSoundex(s[i])
		if code != 0 && code != lastCode {
			result = append(result, code)
			lastCode = code
		}
	}
	for len(result) < 4 {
		result = append(result, '0')
	}
	return string(result)
}

func charToSoundex(c byte) byte {
	switch c {
	case 'B', 'F', 'P', 'V':
		return '1'
	case 'C', 'G', 'J', 'K', 'Q', 'S', 'X', 'Z':
		return '2'
	case 'D', 'T':
		return '3'
	case 'L':
		return '4'
	case 'M', 'N':
		return '5'
	case 'R':
		return '6'
	default:
		return 0
	}
}
