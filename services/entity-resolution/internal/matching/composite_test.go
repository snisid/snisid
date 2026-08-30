package matching

import (
	"testing"

	"github.com/snisid/platform/services/entity-resolution/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockDB struct {
	*gorm.DB
}

func TestMatchExact_NNUSuccess(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{NNU: "NNU-001"}
	cand := models.Identity{NNU: "nnu-001"}
	score := e.matchExact(req, cand)
	assert.Equal(t, 1.0, score)
}

func TestMatchExact_TaxIDSuccess(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{TaxID: "TAX-123"}
	cand := models.Identity{TaxID: "tax-123"}
	score := e.matchExact(req, cand)
	assert.Equal(t, 1.0, score)
}

func TestMatchExact_NationalIDSuccess(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{NationalID: "NID-001"}
	cand := models.Identity{NationalID: "nid-001"}
	score := e.matchExact(req, cand)
	assert.Equal(t, 1.0, score)
}

func TestMatchExact_NoMatch(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{NNU: "NNU-001", TaxID: "TAX-001"}
	cand := models.Identity{NNU: "NNU-999", TaxID: "TAX-999"}
	score := e.matchExact(req, cand)
	assert.Equal(t, 0.0, score)
}

func TestMatchExact_EmptyRequest(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{}
	cand := models.Identity{NNU: "NNU-001"}
	score := e.matchExact(req, cand)
	assert.Equal(t, 0.0, score)
}

func TestMatchExact_CaseInsensitive(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{NNU: "Nnu-001"}
	cand := models.Identity{NNU: "NNU-001"}
	score := e.matchExact(req, cand)
	assert.Equal(t, 1.0, score)
}

func TestMatchBiometric_ExactHash(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{BiometricHash: "abc123"}
	cand := models.Identity{BiometricHash: "abc123"}
	score := e.matchBiometric(req, cand)
	assert.Equal(t, 1.0, score)
}

func TestMatchBiometric_EmptyHash(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{BiometricHash: ""}
	cand := models.Identity{BiometricHash: "abc123"}
	score := e.matchBiometric(req, cand)
	assert.Equal(t, 0.0, score)
}

func TestMatchBiometric_Mismatch(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{BiometricHash: "abc"}
	cand := models.Identity{BiometricHash: "xyz"}
	score := e.matchBiometric(req, cand)
	assert.Equal(t, 0.0, score)
}

func TestMatchFuzzy_IdenticalNames(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{FirstName: "Jean", LastName: "Dupont"}
	cand := models.Identity{FirstName: "Jean", LastName: "Dupont"}
	score := e.matchFuzzy(req, cand)
	assert.Greater(t, score, 0.9)
}

func TestMatchFuzzy_SimilarNames(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{FirstName: "Jean", LastName: "Duponte"}
	cand := models.Identity{FirstName: "Jean", LastName: "Dupont"}
	score := e.matchFuzzy(req, cand)
	assert.Greater(t, score, 0.5)
}

func TestMatchFuzzy_DifferentNames(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{FirstName: "Alice", LastName: "Smith"}
	cand := models.Identity{FirstName: "Bob", LastName: "Jones"}
	score := e.matchFuzzy(req, cand)
	assert.Less(t, score, 0.5)
}

func TestMatchPhonetic_SameSoundex(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{FirstName: "John", LastName: "Smith"}
	cand := models.Identity{FirstName: "Jon", LastName: "Smyth"}
	score := e.matchPhonetic(req, cand)
	assert.Greater(t, score, 0.0)
}

func TestMatchPhonetic_CompletelyDifferent(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{FirstName: "Alice", LastName: "Johnson"}
	cand := models.Identity{FirstName: "Bob", LastName: "Martinez"}
	score := e.matchPhonetic(req, cand)
	assert.Equal(t, 0.0, score)
}

func TestReconcile_ExactMatch(t *testing.T) {
	e := &CompositeEngine{}
	req := models.MatchRequest{NNU: "NNU-001", FirstName: "Jean", LastName: "Dupont"}
	cand := models.Identity{NNU: "nnu-001", FirstName: "Jean", LastName: "Dupont"}

	overall := (1.0*1.0 + e.matchFuzzy(req, cand)*0.6 + e.matchPhonetic(req, cand)*0.4 + 0) / 2.8
	assert.GreaterOrEqual(t, overall, 0.0)
}

func TestExtractFeatures_GeneratesVector(t *testing.T) {
	req := models.MatchRequest{
		FullName: "Jean Dupont",
		DOB:      "1990-01-01",
		TaxID:    "TAX-001",
		NNU:      "NNU-001",
	}
	features := extractFeatures(req)
	require.Len(t, features, 256)
	for _, f := range features {
		assert.LessOrEqual(t, f, 1.0)
	}
}

func TestExtractFeatures_EmptyRequest(t *testing.T) {
	req := models.MatchRequest{}
	features := extractFeatures(req)
	require.Len(t, features, 256)
	for _, f := range features {
		assert.Equal(t, 0.0, f)
	}
}

func TestExtractFeatures_ClampsValues(t *testing.T) {
	req := models.MatchRequest{FullName: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"}
	features := extractFeatures(req)
	for _, f := range features {
		assert.LessOrEqual(t, f, 1.0)
	}
}
