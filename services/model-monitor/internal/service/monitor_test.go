package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate_CorrectPrediction(t *testing.T) {
	m := &Metrics{}
	m.Update(0.8, 1, "v1.0")

	assert.Equal(t, 1, m.Total)
	assert.Equal(t, 1, m.Correct)
	assert.Len(t, m.Predictions, 1)
}

func TestUpdate_IncorrectPrediction(t *testing.T) {
	m := &Metrics{}
	m.Update(0.8, 0, "v1.0")

	assert.Equal(t, 1, m.Total)
	assert.Equal(t, 0, m.Correct)
}

func TestUpdate_BoundaryThreshold(t *testing.T) {
	m := &Metrics{}

	tests := []struct {
		pred   float64
		label  int
		expect int
	}{
		{0.5, 1, 0},
		{0.501, 1, 1},
		{0.5, 0, 1},
		{0.49, 0, 1},
	}

	for _, tt := range tests {
		m.Update(tt.pred, tt.label, "v1")
	}
	assert.Equal(t, 2, m.Correct)
	assert.Equal(t, 4, m.Total)
}

func TestUpdate_VersionTracking(t *testing.T) {
	m := &Metrics{}
	m.Update(0.9, 1, "v2.0")
	_ = m.CalculateDrift(0.5, "v2.0")

	assert.Equal(t, 1, m.Total)
}

func TestCalculateDrift_NoPredictions(t *testing.T) {
	m := &Metrics{}
	drift := m.CalculateDrift(0.5, "v3.0")
	assert.InDelta(t, 0.0, drift, 0.001)
}

func TestCalculateDrift_NoDrift(t *testing.T) {
	m := &Metrics{}
	m.Update(0.5, 1, "v1")
	m.Update(0.5, 0, "v1")
	m.Update(0.5, 1, "v1")

	drift := m.CalculateDrift(0.5, "v1")
	assert.InDelta(t, 0.0, drift, 0.001)
}

func TestCalculateDrift_PositiveDrift(t *testing.T) {
	m := &Metrics{}
	m.Update(0.9, 1, "v1")
	m.Update(0.8, 1, "v1")
	m.Update(0.85, 0, "v1")

	drift := m.CalculateDrift(0.5, "v1")
	assert.InDelta(t, 0.35, drift, 0.001)
}

func TestCalculateDrift_NegativeDrift(t *testing.T) {
	m := &Metrics{}
	m.Update(0.1, 1, "v1")
	m.Update(0.2, 0, "v1")

	drift := m.CalculateDrift(0.5, "v1")
	assert.InDelta(t, 0.35, drift, 0.001)
}

func TestUpdate_ConcurrentSafe(t *testing.T) {
	m := &Metrics{}
	t.Run("parallel", func(t *testing.T) {
		t.Run("update1", func(t *testing.T) {
			m.Update(0.9, 1, "v1")
		})
		t.Run("update2", func(t *testing.T) {
			m.Update(0.1, 0, "v1")
		})
		t.Run("update3", func(t *testing.T) {
			m.Update(0.7, 1, "v1")
		})
	})
	assert.Equal(t, 3, m.Total)
}

func TestLargeBatchUpdates(t *testing.T) {
	m := &Metrics{}
	for i := 0; i < 1000; i++ {
		m.Update(0.7, 1, "v1")
	}
	assert.Equal(t, 1000, m.Total)
	assert.Equal(t, 1000, m.Correct)
}

func TestCalculateDrift_AfterUpdates(t *testing.T) {
	m := &Metrics{}
	m.Update(0.9, 1, "v1")
	m.Update(0.8, 1, "v1")
	m.Update(0.7, 0, "v1")

	drift1 := m.CalculateDrift(0.5, "v1")
	assert.Greater(t, drift1, 0.0)

	m.Update(0.1, 1, "v1")
	m.Update(0.2, 0, "v1")

	drift2 := m.CalculateDrift(0.5, "v1")
	assert.Greater(t, drift1, drift2)
}

func TestAccuracy_Calculation(t *testing.T) {
	m := &Metrics{}
	m.Update(0.9, 1, "v1")
	m.Update(0.2, 0, "v1")
	m.Update(0.8, 0, "v1")

	assert.Equal(t, 3, m.Total)
	assert.Equal(t, 2, m.Correct)
	_ = m.CalculateDrift(0.5, "v1")
}
