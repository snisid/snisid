package biometrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBiometricsEngine(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	assert.NotNil(t, e)
	assert.Equal(t, "v2.1.0", e.ModelVersion)
	assert.Equal(t, 0.97, e.faceThreshold)
	assert.Equal(t, 0.95, e.fingerThreshold)
	assert.Equal(t, 0.99, e.irisThreshold)
}

func TestExtractFaceEmbedding_Success(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	image := make([]byte, 2048)
	for i := range image {
		image[i] = byte(i % 256)
	}

	vector, err := e.ExtractFaceEmbedding(image)
	require.NoError(t, err)
	assert.Len(t, vector, 512)

	// Check normalized (unit vector)
	var sumSquares float64
	for _, v := range vector {
		sumSquares += v * v
	}
	assert.InDelta(t, 1.0, sumSquares, 0.01)
}

func TestExtractFaceEmbedding_EmptyImage(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	_, err := e.ExtractFaceEmbedding([]byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty image")
}

func TestExtractFaceEmbedding_TooSmall(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	_, err := e.ExtractFaceEmbedding(make([]byte, 512))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "image too small")
}

func TestCompareFaceEmbeddings_Identical(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	image := make([]byte, 2048)
	v1, _ := e.ExtractFaceEmbedding(image)
	v2, _ := e.ExtractFaceEmbedding(image)

	similarity := e.CompareFaceEmbeddings(v1, v2)
	assert.GreaterOrEqual(t, similarity, 0.99)
}

func TestCompareFaceEmbeddings_Different(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	img1 := make([]byte, 2048)
	img2 := make([]byte, 2048)
	for i := range img2 {
		img2[i] = 0xFF
	}

	v1, _ := e.ExtractFaceEmbedding(img1)
	v2, _ := e.ExtractFaceEmbedding(img2)

	similarity := e.CompareFaceEmbeddings(v1, v2)
	assert.Less(t, similarity, 0.99)
}

func TestCompareFaceEmbeddings_EmptyVectors(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	assert.Equal(t, 0.0, e.CompareFaceEmbeddings(BiometricVector{}, BiometricVector{512: 0}))
	assert.Equal(t, 0.0, e.CompareFaceEmbeddings(nil, BiometricVector{}))
}

func TestMatchFace_AboveThreshold(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	image := make([]byte, 2048)
	vector, _ := e.ExtractFaceEmbedding(image)

	result := e.MatchFace(image, []BiometricVector{vector})
	assert.True(t, result.Match)
	assert.GreaterOrEqual(t, result.Confidence, result.Threshold)
}

func TestMatchFace_BelowThreshold(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	img1 := make([]byte, 2048)
	img2 := make([]byte, 2048)
	for i := range img2 {
		img2[i] = 0xFF
	}

	v2, _ := e.ExtractFaceEmbedding(img2)
	result := e.MatchFace(img1, []BiometricVector{v2})
	assert.False(t, result.Match)
}

func TestMatchFace_EmptyImage(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	result := e.MatchFace([]byte{}, []BiometricVector{{0.5, 0.5}})
	assert.False(t, result.Match)
	assert.Equal(t, 0.0, result.Confidence)
}

func TestVerifyLiveness_Alive(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	image := make([]byte, 4096)
	result := e.VerifyLiveness(image)
	assert.True(t, result.Alive)
	assert.GreaterOrEqual(t, result.Score, 0.85)
	assert.Equal(t, "multimodal_liveness_v2", result.Method)
	assert.Len(t, result.Details, 4)
}

func TestVerifyLiveness_Empty(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	result := e.VerifyLiveness([]byte{})
	assert.False(t, result.Alive)
	assert.Equal(t, 0.0, result.Score)
}

func TestMatchFingerprint_AboveThreshold(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	template := make([]byte, 64)
	result := e.MatchFingerprint(template, [][]byte{template})
	assert.True(t, result.Match)
}

func TestMatchFingerprint_Empty(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	result := e.MatchFingerprint([]byte{}, [][]byte{{1, 2, 3}})
	assert.False(t, result.Match)
	assert.Equal(t, 0.0, result.Confidence)
}

func TestMatchIris_Success(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	image := make([]byte, 1024)
	id, confidence := e.MatchIris(image)
	assert.Contains(t, id, "iris_")
	assert.Greater(t, confidence, 0.0)
}

func TestMatchIris_TooSmall(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	id, confidence := e.MatchIris(make([]byte, 256))
	assert.Empty(t, id)
	assert.Equal(t, 0.0, confidence)
}

func TestCalculateQualityScore_Face(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	score := e.CalculateQualityScore(make([]byte, 4096), "face")
	assert.Greater(t, score, 0.0)
	assert.LessOrEqual(t, score, 1.0)
}

func TestCalculateQualityScore_Empty(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	score := e.CalculateQualityScore([]byte{}, "face")
	assert.Equal(t, 0.0, score)
}

func TestEnrollFace_Success(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	image := make([]byte, 4096)
	enrollment, err := e.EnrollFace("CITIZEN-001", image)
	require.NoError(t, err)
	assert.Equal(t, "CITIZEN-001", enrollment.CitizenID)
	assert.Equal(t, "face", enrollment.Modality)
	assert.GreaterOrEqual(t, enrollment.QualityScore, 0.5)
	assert.Len(t, enrollment.Vector, 512)
}

func TestEnrollFace_LowQuality(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	_, err := e.EnrollFace("CITIZEN-002", make([]byte, 100))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "image quality too low")
}

func TestIsDuplicate_AboveThreshold(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	image := make([]byte, 2048)
	vector, _ := e.ExtractFaceEmbedding(image)

	duplicate, similarity := e.IsDuplicate(vector, []BiometricVector{vector}, 0)
	assert.True(t, duplicate)
	assert.GreaterOrEqual(t, similarity, 0.97)
}

func TestIsDuplicate_BelowThreshold(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	img1 := make([]byte, 2048)
	img2 := make([]byte, 2048)
	for i := range img2 {
		img2[i] = 0xFF
	}

	v1, _ := e.ExtractFaceEmbedding(img1)
	v2, _ := e.ExtractFaceEmbedding(img2)

	duplicate, _ := e.IsDuplicate(v1, []BiometricVector{v2}, 0.99)
	assert.False(t, duplicate)
}

func TestNormalizeVector(t *testing.T) {
	e := NewBiometricsEngine("v2.1.0")
	v := BiometricVector{3, 4}
	e.normalizeVector(v)
	assert.InDelta(t, 0.6, v[0], 0.001)
	assert.InDelta(t, 0.8, v[1], 0.001)
}
