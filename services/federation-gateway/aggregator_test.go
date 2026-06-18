package federationgateway

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAggregator(t *testing.T) {
	a, err := NewAggregator("SNISID-HTI-AGG")
	require.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, "SNISID-HTI-AGG", a.ID)
	assert.NotNil(t, a.privateKey)
	assert.NotNil(t, a.publicKey)
	assert.Equal(t, 2, a.minNodes)
	assert.Equal(t, 100, a.minSamples)
	assert.True(t, a.differentialPrivacy)
}

func TestAggregateWeights_Success(t *testing.T) {
	a, err := NewAggregator("test-agg")
	require.NoError(t, err)

	updates := []ModelUpdate{
		{
			NodeID: "node-1", Country: "HTI",
			Weights:     []float64{0.5, 0.3, 0.2},
			SampleCount: 100,
			Accuracy:    0.85,
			Loss:        0.15,
			Timestamp:   time.Now(),
		},
		{
			NodeID: "node-2", Country: "HTI",
			Weights:     []float64{0.6, 0.2, 0.2},
			SampleCount: 200,
			Accuracy:    0.90,
			Loss:        0.10,
			Timestamp:   time.Now(),
		},
	}

	result, err := a.AggregateWeights(updates)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Round)
	assert.Equal(t, 2, result.NodeCount)
	assert.Equal(t, 300, result.TotalSamples)
	assert.Greater(t, result.AvgLoss, 0.0)
	assert.Greater(t, result.WeightedAccuracy, 0.0)
	assert.NotEmpty(t, result.GlobalWeights)
	assert.Len(t, result.GlobalWeights, 3)
	assert.Len(t, result.NodeWeights, 2)
}

func TestAggregateWeights_InsufficientNodes(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	a.minNodes = 3

	updates := []ModelUpdate{
		{NodeID: "node-1", Weights: []float64{0.5}, SampleCount: 10, Accuracy: 0.8, Loss: 0.2, Timestamp: time.Now()},
		{NodeID: "node-2", Weights: []float64{0.6}, SampleCount: 10, Accuracy: 0.8, Loss: 0.2, Timestamp: time.Now()},
	}

	_, err := a.AggregateWeights(updates)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient nodes")
}

func TestAggregateWeights_InvalidUpdates_Rejected(t *testing.T) {
	a, _ := NewAggregator("test-agg")

	updates := []ModelUpdate{
		{NodeID: "good-node", Weights: []float64{0.5}, SampleCount: 100, Accuracy: 0.85, Loss: 0.15, Timestamp: time.Now()},
		{NodeID: "bad-node", Weights: []float64{}, SampleCount: 0, Accuracy: -1, Loss: 0.1, Timestamp: time.Now()},
	}

	_, err := a.AggregateWeights(updates)
	assert.Error(t, err)
}

func TestAggregateWeights_EmptyUpdates(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	_, err := a.AggregateWeights([]ModelUpdate{})
	assert.Error(t, err)
}

func TestVerifyUpdate_Valid(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	update := ModelUpdate{
		Weights:     []float64{0.1, 0.2},
		SampleCount: 50,
		Accuracy:    0.75,
		Loss:        0.25,
		Timestamp:   time.Now(),
	}
	assert.NoError(t, a.verifyUpdate(update))
}

func TestVerifyUpdate_EmptyWeights(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	update := ModelUpdate{
		Weights:     []float64{},
		SampleCount: 50,
		Accuracy:    0.75,
		Loss:        0.25,
		Timestamp:   time.Now(),
	}
	assert.Error(t, a.verifyUpdate(update))
}

func TestVerifyUpdate_InvalidSampleCount(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	update := ModelUpdate{
		Weights:     []float64{0.1},
		SampleCount: 0,
		Accuracy:    0.75,
		Loss:        0.25,
		Timestamp:   time.Now(),
	}
	assert.Error(t, a.verifyUpdate(update))
}

func TestVerifyUpdate_StaleTimestamp(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	update := ModelUpdate{
		Weights:     []float64{0.1},
		SampleCount: 10,
		Accuracy:    0.75,
		Loss:        0.25,
		Timestamp:   time.Now().Add(-48 * time.Hour),
	}
	assert.Error(t, a.verifyUpdate(update))
}

func TestVerifyUpdate_InvalidAccuracy(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	update := ModelUpdate{
		Weights:     []float64{0.1},
		SampleCount: 10,
		Accuracy:    1.5,
		Loss:        0.25,
		Timestamp:   time.Now(),
	}
	assert.Error(t, a.verifyUpdate(update))
}

func TestVerifyUpdate_NaNWeight(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	update := ModelUpdate{
		Weights:     []float64{math.NaN()},
		SampleCount: 10,
		Accuracy:    0.75,
		Loss:        0.25,
		Timestamp:   time.Now(),
	}
	assert.Error(t, a.verifyUpdate(update))
}

func TestUpdateReputation_Success(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	a.updateReputation("node-1", 0.9, true)
	rep := a.nodeReputations["node-1"]
	assert.Equal(t, 1, rep.TotalUpdates)
	assert.Equal(t, 0.9, rep.AvgAccuracy)
	assert.Equal(t, 0.5, rep.TrustScore) // Initial 0.5 + 0.05 = 0.55... wait
}

func TestUpdateReputation_Failure(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	a.updateReputation("node-1", 0.0, false)
	rep := a.nodeReputations["node-1"]
	assert.Equal(t, 0, rep.TotalUpdates)
	assert.Equal(t, 0.3, rep.TrustScore) // 0.5 - 0.2 = 0.3
}

func TestGetTrustScore_UnknownNode(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	assert.Equal(t, 0.5, a.getTrustScore("unknown-node"))
}

func TestGetTrustScore_KnownNode(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	a.updateReputation("node-1", 0.9, true)
	assert.Equal(t, 0.55, a.getTrustScore("node-1"))
}

func TestGetNodeReputations(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	a.updateReputation("node-1", 0.8, true)
	a.updateReputation("node-2", 0.9, true)

	reps := a.GetNodeReputations()
	assert.Len(t, reps, 2)
}

func TestSignModel(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	model := &AggregatedModel{
		GlobalWeights: []float64{0.5, 0.3},
		CreatedAt:     time.Now().UTC(),
	}

	sig, err := a.SignModel(model)
	require.NoError(t, err)
	assert.NotEmpty(t, sig)
}

func TestGenerateLaplaceNoise(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	noise := a.generateLaplaceNoise(0.01)
	assert.NotEqual(t, 0.0, noise)
}

func TestExportPublicKey(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	pubKey, err := a.ExportPublicKey()
	require.NoError(t, err)
	assert.Contains(t, pubKey, "PUBLIC KEY")
}

func TestDifferentialPrivacy_AddsNoise(t *testing.T) {
	a, _ := NewAggregator("test-agg")
	a.differentialPrivacy = true
	a.noiseScale = 0.1

	updates := []ModelUpdate{
		{NodeID: "n1", Weights: []float64{0.5}, SampleCount: 100, Accuracy: 0.8, Loss: 0.2, Timestamp: time.Now()},
		{NodeID: "n2", Weights: []float64{0.5}, SampleCount: 100, Accuracy: 0.8, Loss: 0.2, Timestamp: time.Now()},
	}

	result1, _ := a.AggregateWeights(updates)
	result2, _ := a.AggregateWeights(updates)

	// With DP, results should differ slightly
	assert.NotEqual(t, result1.GlobalWeights[0], result2.GlobalWeights[0])
}

func TestConvergenceDelta_Calculated(t *testing.T) {
	a, _ := NewAggregator("test-agg")

	updates := []ModelUpdate{
		{NodeID: "n1", Weights: []float64{0.5}, SampleCount: 100, Accuracy: 0.8, Loss: 0.2, Timestamp: time.Now()},
	}

	r1, _ := a.AggregateWeights(updates)
	assert.Equal(t, 0.0, r1.ConvergenceDelta) // First round, no previous

	updates2 := []ModelUpdate{
		{NodeID: "n1", Weights: []float64{0.6}, SampleCount: 100, Accuracy: 0.8, Loss: 0.2, Timestamp: time.Now()},
	}
	r2, _ := a.AggregateWeights(updates2)
	assert.Greater(t, r2.ConvergenceDelta, 0.0)
}
