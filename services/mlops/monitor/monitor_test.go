package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordPrediction_Correct(t *testing.T) {
	m := &ModelMetrics{}
	m.RecordPrediction(0.8, 1)

	assert.Equal(t, 1, m.Total)
	assert.Equal(t, 1, m.Correct)
	assert.Len(t, m.Predictions, 1)
	assert.Len(t, m.GroundTruths, 1)
}

func TestRecordPrediction_Incorrect(t *testing.T) {
	m := &ModelMetrics{}
	m.RecordPrediction(0.8, 0)

	assert.Equal(t, 1, m.Total)
	assert.Equal(t, 0, m.Correct)
}

func TestRecordPrediction_BoundaryThreshold(t *testing.T) {
	m := &ModelMetrics{}

	tests := []struct {
		pred    float64
		actual  int
		correct bool
	}{
		{0.5, 0, false},
		{0.501, 1, true},
		{0.5, 1, false},
		{0.49, 0, true},
	}

	for _, tt := range tests {
		m.RecordPrediction(tt.pred, tt.actual)
	}
	assert.Equal(t, 4, m.Total)
	assert.Equal(t, 2, m.Correct)
}

func TestAccuracy_ZeroTotal(t *testing.T) {
	m := &ModelMetrics{}
	assert.InDelta(t, 0.0, m.Accuracy(), 0.001)
}

func TestAccuracy_Perfect(t *testing.T) {
	m := &ModelMetrics{}
	m.RecordPrediction(0.9, 1)
	m.RecordPrediction(0.1, 0)
	assert.InDelta(t, 1.0, m.Accuracy(), 0.001)
}

func TestAccuracy_Partial(t *testing.T) {
	m := &ModelMetrics{}
	m.RecordPrediction(0.9, 1)
	m.RecordPrediction(0.8, 1)
	m.RecordPrediction(0.1, 1)
	assert.InDelta(t, 2.0/3.0, m.Accuracy(), 0.001)
}

func TestCalculateDrift_NotEnoughSamples(t *testing.T) {
	m := &ModelMetrics{}
	m.RecordPrediction(0.5, 1)
	m.RecordPrediction(0.5, 0)

	assert.InDelta(t, 0.0, m.DriftScore, 0.001)
}

func TestCalculateDrift_NoDrift(t *testing.T) {
	m := &ModelMetrics{}
	for i := 0; i < 10; i++ {
		m.RecordPrediction(0.5, 1)
	}
	assert.InDelta(t, 0.0, m.DriftScore, 0.001)
}

func TestCalculateDrift_SignificantDrift(t *testing.T) {
	m := &ModelMetrics{}
	for i := 0; i < 10; i++ {
		m.RecordPrediction(0.9, 1)
	}
	assert.Greater(t, m.DriftScore, 0.2)
}

func TestCalculateDrift_DriftValue(t *testing.T) {
	m := &ModelMetrics{}
	for i := 0; i < 10; i++ {
		m.RecordPrediction(0.8, 1)
	}
	expectedDrift := 0.8 - 0.5
	assert.InDelta(t, expectedDrift, m.DriftScore, 0.001)
}

func TestRecordPrediction_ConcurrentSafe(t *testing.T) {
	m := &ModelMetrics{}
	t.Run("parallel writes", func(t *testing.T) {
		t.Run("writer1", func(t *testing.T) {
			m.RecordPrediction(0.9, 1)
		})
		t.Run("writer2", func(t *testing.T) {
			m.RecordPrediction(0.1, 0)
		})
		t.Run("writer3", func(t *testing.T) {
			m.RecordPrediction(0.7, 1)
		})
	})
	assert.Equal(t, 3, m.Total)
}

func TestMultipleRecords(t *testing.T) {
	m := &ModelMetrics{}
	for i := 0; i < 100; i++ {
		m.RecordPrediction(0.6, 1)
	}
	assert.Equal(t, 100, m.Total)
	assert.Equal(t, 100, m.Correct)
}

func TestAccuracy_AllIncorrect(t *testing.T) {
	m := &ModelMetrics{}
	m.RecordPrediction(0.9, 0)
	m.RecordPrediction(0.8, 0)
	assert.InDelta(t, 0.0, m.Accuracy(), 0.001)
}

func TestUpdatePrometheusMetrics(t *testing.T) {
	m := &ModelMetrics{}
	m.RecordPrediction(0.9, 1)
	UpdatePrometheusMetrics(m)
}

func TestMetricsRange(t *testing.T) {
	m := &ModelMetrics{}
	m.RecordPrediction(1.0, 1)
	m.RecordPrediction(0.0, 0)
	m.RecordPrediction(0.5, 0)
	m.RecordPrediction(0.99, 1)

	assert.InDelta(t, 0.75, m.Accuracy(), 0.001)
}
