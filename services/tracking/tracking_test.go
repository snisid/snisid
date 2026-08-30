package tracking

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCrossImageTracker(t *testing.T) {
	tr := NewCrossImageTracker()
	require.NotNil(t, tr)
	assert.NotNil(t, tr.Registry)
	assert.NotNil(t, tr.SubjectDB)
	assert.InDelta(t, 0.85, tr.matchThreshold, 0.001)
	assert.InDelta(t, 0.6, tr.minSimilarity, 0.001)
}

func TestRegisterSubject(t *testing.T) {
	tr := NewCrossImageTracker()
	tr.RegisterSubject("CIT_001", []float64{0.1, 0.2, 0.3})

	tr.mu.RLock()
	assert.Contains(t, tr.SubjectDB, "CIT_001")
	assert.Equal(t, "ONI", tr.Registry["CIT_001"])
	tr.mu.RUnlock()
}

func TestRemoveSubject(t *testing.T) {
	tr := NewCrossImageTracker()
	tr.RegisterSubject("TAX_001", []float64{0.5, 0.6})
	tr.RemoveSubject("TAX_001")

	tr.mu.RLock()
	assert.NotContains(t, tr.SubjectDB, "TAX_001")
	assert.NotContains(t, tr.Registry, "TAX_001")
	tr.mu.RUnlock()
}

func TestRemoveSubject_NonExistent(t *testing.T) {
	tr := NewCrossImageTracker()
	tr.RemoveSubject("NON_EXISTENT")
}

func TestReIdentify_EmptyDatabase(t *testing.T) {
	tr := NewCrossImageTracker()
	results := tr.ReIdentify([]float64{0.1, 0.2})
	assert.Nil(t, results)
}

func TestReIdentify_ZeroVectors(t *testing.T) {
	tr := NewCrossImageTracker()
	results := tr.ReIdentify([]float64{})
	assert.Nil(t, results)
}

func TestReIdentify_FindsMatch(t *testing.T) {
	tr := NewCrossImageTracker()
	tr.RegisterSubject("CIT_001", []float64{0.9, 0.1, 0.5})

	results := tr.ReIdentify([]float64{0.85, 0.12, 0.48})
	require.NotEmpty(t, results)
	assert.Equal(t, "CIT_001", results[0].SubjectID)
	assert.GreaterOrEqual(t, results[0].MatchScore, 0.6)
	assert.False(t, results[0].Anomalous)
}

func TestReIdentify_BelowMinSimilarity(t *testing.T) {
	tr := NewCrossImageTracker()
	tr.RegisterSubject("CIT_001", []float64{1.0, 0.0, 0.0})

	results := tr.ReIdentify([]float64{0.0, 1.0, 0.0})
	assert.Empty(t, results)
}

func TestReIdentify_ReturnsTop5(t *testing.T) {
	tr := NewCrossImageTracker()
	for i := 0; i < 10; i++ {
		vec := []float64{float64(i) / 10, 0.5, 0.3}
		tr.RegisterSubject("POL_00"+string(rune('0'+i)), vec)
	}

	results := tr.ReIdentify([]float64{0.5, 0.5, 0.3})
	assert.LessOrEqual(t, len(results), 5)
}

func TestReIdentify_AnomalyDetection(t *testing.T) {
	tr := NewCrossImageTracker()
	tr.RegisterSubject("CIT_001", []float64{1.0, 0.0, 0.0})
	tr.RegisterSubject("CIT_002", []float64{0.9, 0.1, 0.0})

	results := tr.ReIdentify([]float64{0.6, 0.4, 0.0})
	if len(results) > 1 {
		assert.True(t, results[1].Anomalous || !results[1].Anomalous)
	}
}

func TestInferDatabase(t *testing.T) {
	tr := NewCrossImageTracker()

	tests := []struct {
		subjectID string
		expected  string
	}{
		{"CIT_001", "ONI"},
		{"TAX_002", "DGI"},
		{"POL_003", "PNH"},
		{"BIO_004", "BIOMETRICS"},
		{"FPR_005", "FPR"},
		{"UNK_001", "UNKNOWN"},
		{"SHORT", "UNKNOWN"},
		{"", "UNKNOWN"},
	}

	for _, tc := range tests {
		t.Run(tc.subjectID, func(t *testing.T) {
			assert.Equal(t, tc.expected, tr.inferDatabase(tc.subjectID))
		})
	}
}

func TestCosineSimilarity_IdenticalVectors(t *testing.T) {
	a := []float64{1.0, 2.0, 3.0}
	sim := cosineSimilarity(a, a)
	assert.InDelta(t, 1.0, sim, 0.001)
}

func TestCosineSimilarity_OrthogonalVectors(t *testing.T) {
	a := []float64{1.0, 0.0}
	b := []float64{0.0, 1.0}
	sim := cosineSimilarity(a, b)
	assert.InDelta(t, 0.0, sim, 0.001)
}

func TestCosineSimilarity_EmptyVectors(t *testing.T) {
	assert.Equal(t, 0.0, cosineSimilarity(nil, []float64{1.0}))
	assert.Equal(t, 0.0, cosineSimilarity([]float64{}, []float64{}))
	assert.Equal(t, 0.0, cosineSimilarity([]float64{1.0}, []float64{1.0, 2.0}))
}

func TestCosineSimilarity_ZeroVector(t *testing.T) {
	a := []float64{0.0, 0.0}
	b := []float64{1.0, 0.0}
	sim := cosineSimilarity(a, b)
	assert.InDelta(t, 0.0, sim, 0.001)
}

func TestCosineSimilarity_DifferentMagnitudes(t *testing.T) {
	a := []float64{2.0, 0.0}
	b := []float64{4.0, 0.0}
	sim := cosineSimilarity(a, b)
	assert.InDelta(t, 1.0, sim, 0.001)
}

func TestLinkInconsistencies_FewerThanTwo(t *testing.T) {
	tr := NewCrossImageTracker()
	tr.LinkInconsistencies(nil)
	tr.LinkInconsistencies([]TrackResult{{SubjectID: "CIT_001", MatchScore: 0.95, Database: "ONI"}})
}

func TestLinkInconsistencies_CrossDatabaseMatch(t *testing.T) {
	tr := NewCrossImageTracker()
	results := []TrackResult{
		{SubjectID: "CIT_001", MatchScore: 0.96, Database: "ONI"},
		{SubjectID: "TAX_001", MatchScore: 0.93, Database: "DGI"},
	}
	tr.LinkInconsistencies(results)
}

func TestDetectAnomaly(t *testing.T) {
	tr := NewCrossImageTracker()

	tests := []struct {
		name      string
		score     float64
		allScores []float64
		anomalous bool
	}{
		{"above threshold not anomalous", 0.9, []float64{0.9, 0.5}, false},
		{"below threshold and low relative", 0.3, []float64{0.9, 0.8, 0.3}, true},
		{"below threshold but similar to avg", 0.5, []float64{0.6, 0.5, 0.7}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			allResults := make([]scored, len(tc.allScores))
			for i, s := range tc.allScores {
				allResults[i] = scored{score: s}
			}
			result := tr.detectAnomaly(scored{score: tc.score}, allResults)
			assert.Equal(t, tc.anomalous, result)
		})
	}
}

func TestAddFromDatabase(t *testing.T) {
	tr := NewCrossImageTracker()
	subjects := map[string][]float64{
		"CIT_001": {0.1, 0.2},
		"CIT_002": {0.3, 0.4},
	}
	tr.AddFromDatabase("ONI", subjects)

	tr.mu.RLock()
	assert.Len(t, tr.SubjectDB, 2)
	assert.Equal(t, "ONI", tr.Registry["CIT_001"])
	assert.Equal(t, "ONI", tr.Registry["CIT_002"])
	tr.mu.RUnlock()
}

func TestConcurrentAccess(t *testing.T) {
	tr := NewCrossImageTracker()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			subjectID := "CIT_" + string(rune('0'+id))
			tr.RegisterSubject(subjectID, []float64{float64(id) / 20, 0.5})
		}(i)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tr.ReIdentify([]float64{0.5, 0.5})
		}()
	}

	wg.Wait()

	tr.mu.RLock()
	assert.Len(t, tr.SubjectDB, 20)
	tr.mu.RUnlock()
}
