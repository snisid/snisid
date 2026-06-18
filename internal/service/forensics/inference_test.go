package forensics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMesoNetForensicEngine_Analyze(t *testing.T) {
	engine := NewMesoNetForensicEngine("localhost:50051", 10)

	mediaData := make([]byte, 1024)
	for i := range mediaData {
		mediaData[i] = byte(i % 256)
	}

	result, err := engine.Analyze(context.Background(), mediaData)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, result.DeepfakeProbability, 0.0)
	assert.LessOrEqual(t, result.DeepfakeProbability, 1.0)
	assert.Equal(t, "mesonet4", result.ModelVersion)
}

func TestMesoNetForensicEngine_Analyze_EmptyData(t *testing.T) {
	engine := NewMesoNetForensicEngine("localhost:50051", 10)

	_, err := engine.Analyze(context.Background(), []byte{})
	assert.Error(t, err)
}

func TestMockForensicEngine_Analyze(t *testing.T) {
	engine := &MockForensicEngine{}

	result, err := engine.Analyze(context.Background(), []byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, 0.1, result.DeepfakeProbability)
}
