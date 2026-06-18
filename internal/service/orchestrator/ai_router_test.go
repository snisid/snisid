package orchestrator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAIRouter_DispatchAnalysis_Success(t *testing.T) {
	router := &AIRouter{}
	mediaData := []byte("fake-image-data-1024-bytes-minimum")

	verdict, err := router.DispatchAnalysis(context.Background(), mediaData)
	require.NoError(t, err)
	require.NotNil(t, verdict)

	assert.GreaterOrEqual(t, verdict.BiometricScore, float32(0))
	assert.LessOrEqual(t, verdict.BiometricScore, float32(100))
	assert.GreaterOrEqual(t, verdict.DeepfakeProb, float32(0))
	assert.GreaterOrEqual(t, verdict.FraudScore, 0)
	assert.NotEmpty(t, verdict.RiskLevel)
}

func TestAIRouter_DispatchAnalysis_EmptyData(t *testing.T) {
	router := &AIRouter{}
	verdict, err := router.DispatchAnalysis(context.Background(), []byte{})
	require.NoError(t, err)
	require.NotNil(t, verdict)
	// Should still return default mock values
	assert.Equal(t, float32(98.5), verdict.BiometricScore)
}
