package biometrics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestL2Normalize_UnitVector(t *testing.T) {
	v := []float32{1, 0, 0}
	result := l2Normalize(v)
	assert.InDelta(t, 1.0, result[0], 0.001)
	assert.InDelta(t, 0.0, result[1], 0.001)
}

func TestL2Normalize_ZeroVector(t *testing.T) {
	v := []float32{0, 0, 0}
	result := l2Normalize(v)
	assert.Equal(t, v, result)
}

func TestL2Normalize_Norm(t *testing.T) {
	v := []float32{3, 4}
	result := l2Normalize(v)
	var norm float32
	for _, x := range result {
		norm += x * x
	}
	assert.InDelta(t, 1.0, float64(norm), 0.001)
}

func TestONNXInferenceEngine_GenerateEmbedding(t *testing.T) {
	t.Skip("integration test: requires external ONNX service")
	engine, err := NewONNXInferenceEngine("dummy.onnx")
	assert.NoError(t, err)

	imageData := make([]byte, 112*112*3)
	for i := range imageData {
		imageData[i] = byte(i % 256)
	}

	embedding, err := engine.GenerateEmbedding(context.Background(), imageData, "face")
	assert.NoError(t, err)
	assert.Len(t, embedding, 512)

	var norm float32
	for _, x := range embedding {
		norm += x * x
	}
	assert.InDelta(t, 1.0, float64(norm), 0.01)
}

func TestCosineSimilarity_SameVector(t *testing.T) {
	a := []float32{1, 0, 0}
	sim := CosineSimilarity(a, a)
	assert.InDelta(t, 1.0, sim, 0.001)
}

func TestCosineSimilarity_Orthogonal(t *testing.T) {
	a := []float32{1, 0}
	b := []float32{0, 1}
	sim := CosineSimilarity(a, b)
	assert.InDelta(t, 0.0, sim, 0.001)
}

func TestPreprocessImage_OutputDimensions(t *testing.T) {
	data := make([]byte, 100)
	result, err := preprocessImage(data, 112, 112)
	assert.NoError(t, err)
	assert.Len(t, result, 3*112*112)
}
