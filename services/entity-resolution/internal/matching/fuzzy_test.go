package matching

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJaroWinkler_Identical(t *testing.T) {
	score := JaroWinkler("Jean", "Jean")
	assert.Equal(t, 1.0, score)
}

func TestJaroWinkler_Empty(t *testing.T) {
	score := JaroWinkler("", "Jean")
	assert.Equal(t, 0.0, score)

	score2 := JaroWinkler("Jean", "")
	assert.Equal(t, 0.0, score2)

	score3 := JaroWinkler("", "")
	assert.Equal(t, 0.0, score3)
}

func TestJaroWinkler_Similar(t *testing.T) {
	score := JaroWinkler("Jon", "John")
	assert.Greater(t, score, 0.7)
}

func TestJaroWinkler_Different(t *testing.T) {
	score := JaroWinkler("Alice", "Bob")
	assert.Less(t, score, 0.5)
}

func TestJaroWinkler_CaseInsensitive(t *testing.T) {
	score1 := JaroWinkler("jean", "JEAN")
	assert.Equal(t, 1.0, score1)
}

func TestJaroWinkler_WhitespaceTrimmed(t *testing.T) {
	score := JaroWinkler("  Jean  ", "JEAN")
	assert.Equal(t, 1.0, score)
}

func TestLevenshtein_Identical(t *testing.T) {
	d := Levenshtein("Jean", "Jean")
	assert.Equal(t, 0, d)
}

func TestLevenshtein_Empty(t *testing.T) {
	d := Levenshtein("", "Jean")
	assert.Equal(t, 4, d)

	d2 := Levenshtein("Jean", "")
	assert.Equal(t, 4, d2)
}

func TestLevenshtein_Different(t *testing.T) {
	d := Levenshtein("kitten", "sitting")
	assert.Equal(t, 3, d)
}

func TestLevenshtein_SingleChar(t *testing.T) {
	d := Levenshtein("a", "b")
	assert.Equal(t, 1, d)
}

func TestNormalizedLevenshtein_Identical(t *testing.T) {
	s := NormalizedLevenshtein("Jean", "Jean")
	assert.Equal(t, 1.0, s)
}

func TestNormalizedLevenshtein_Empty(t *testing.T) {
	s := NormalizedLevenshtein("", "Jean")
	assert.Equal(t, 0.0, s)
}

func TestNormalizedLevenshtein_Similar(t *testing.T) {
	s := NormalizedLevenshtein("kitten", "sitten")
	assert.Greater(t, s, 0.8)
}

func TestSoundex_Basic(t *testing.T) {
	assert.Equal(t, "J500", Soundex("Jean"))
	assert.Equal(t, "D153", Soundex("Dupont"))
}

func TestSoundex_Empty(t *testing.T) {
	assert.Equal(t, "", Soundex(""))
}

func TestSoundex_DifferentInputsSameCode(t *testing.T) {
	s1 := Soundex("Smith")
	s2 := Soundex("Smyth")
	assert.Equal(t, s1, s2)
}

func TestSoundex_ShortName(t *testing.T) {
	code := Soundex("Jo")
	assert.Equal(t, "J000", code)
}

func TestMetaphone_Basic(t *testing.T) {
	assert.Equal(t, "JN", Metaphone("Jean"))
	assert.Equal(t, "TPNT", Metaphone("Dupont"))
}

func TestMetaphone_Empty(t *testing.T) {
	assert.Equal(t, "", Metaphone(""))
}

func TestMetaphone_SilentPrefix(t *testing.T) {
	assert.Equal(t, "N", Metaphone("Knight"))
	assert.Equal(t, "N", Metaphone("Gnome"))
}

func TestMetaphone_Complex(t *testing.T) {
	assert.Equal(t, "KMPL", Metaphone("Complex"))
}

func TestNameSimilarity_Identical(t *testing.T) {
	s := NameSimilarity("Jean", "Dupont", "Jean", "Dupont")
	assert.Greater(t, s, 0.95)
}

func TestNameSimilarity_Different(t *testing.T) {
	s := NameSimilarity("Alice", "Smith", "Bob", "Jones")
	assert.Less(t, s, 0.5)
}

func TestNameSimilarity_Similar(t *testing.T) {
	s := NameSimilarity("Jon", "Smith", "John", "Smyth")
	assert.Greater(t, s, 0.5)
}

func TestMin(t *testing.T) {
	assert.Equal(t, 1, min(1, 5))
	assert.Equal(t, -3, min(-3, 0))
	assert.Equal(t, 5, min(10, 5))
}

func TestMax(t *testing.T) {
	assert.Equal(t, 5, max(1, 5))
	assert.Equal(t, 0, max(-3, 0))
	assert.Equal(t, 10, max(10, 5))
}
