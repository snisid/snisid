package forensics

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockForensicEngine struct {
	prob      float64
	anomalies []string
	err       error
}

func (e *mockForensicEngine) Analyze(ctx context.Context, mediaData []byte) (*ForensicResult, error) {
	if e.err != nil {
		return nil, e.err
	}
	return &ForensicResult{
		DeepfakeProbability: e.prob,
		Anomalies:           e.anomalies,
	}, nil
}

func TestNewDeepfakeDetector(t *testing.T) {
	engine := &mockForensicEngine{}
	d := NewDeepfakeDetector(engine)
	assert.NotNil(t, d)
	assert.Equal(t, engine, d.engine)
}

func TestDeepfakeDetector_Detect_LowProbability(t *testing.T) {
	engine := &mockForensicEngine{
		prob:      0.03,
		anomalies: []string{},
	}
	d := NewDeepfakeDetector(engine)

	prob, anomalies, err := d.Detect(context.Background(), "CITIZEN-001", []byte("fake-media"))
	require.NoError(t, err)
	assert.Equal(t, float32(0.03), prob)
	assert.Empty(t, anomalies)
}

func TestDeepfakeDetector_Detect_HighProbability(t *testing.T) {
	engine := &mockForensicEngine{
		prob:      0.92,
		anomalies: []string{"Inconsistent facial landmarks", "Eye blinking irregularity"},
	}
	d := NewDeepfakeDetector(engine)

	prob, anomalies, err := d.Detect(context.Background(), "CITIZEN-002", []byte("suspicious-media"))
	require.NoError(t, err)
	assert.Equal(t, float32(0.92), prob)
	assert.Len(t, anomalies, 2)
	assert.Contains(t, anomalies[0], "facial")
}

func TestDeepfakeDetector_Detect_EngineError(t *testing.T) {
	engine := &mockForensicEngine{
		err: errors.New("model inference failed"),
	}
	d := NewDeepfakeDetector(engine)

	prob, anomalies, err := d.Detect(context.Background(), "CITIZEN-003", []byte("bad-media"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forensic inference failed")
	assert.Equal(t, float32(0), prob)
	assert.Nil(t, anomalies)
}

func TestDeepfakeDetector_Detect_EmptyMedia(t *testing.T) {
	engine := &mockForensicEngine{
		prob:      0.5,
		anomalies: []string{},
	}
	d := NewDeepfakeDetector(engine)

	prob, anomalies, err := d.Detect(context.Background(), "CITIZEN-004", []byte{})
	require.NoError(t, err)
	assert.Equal(t, float32(0.5), prob)
}

func TestDeepfakeDetector_Detect_NilMedia(t *testing.T) {
	engine := &mockForensicEngine{
		prob:      0.1,
		anomalies: nil,
	}
	d := NewDeepfakeDetector(engine)

	prob, anomalies, err := d.Detect(context.Background(), "CITIZEN-005", nil)
	require.NoError(t, err)
	assert.Equal(t, float32(0.1), prob)
	assert.Nil(t, anomalies)
}
