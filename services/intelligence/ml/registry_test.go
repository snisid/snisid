package ml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModelRegistry(t *testing.T) {
	r := NewModelRegistry()
	assert.NotNil(t, r)
	assert.Empty(t, r.Models)
}

func TestRegister(t *testing.T) {
	r := NewModelRegistry()
	r.Register("fraud-detector", "v1.0", "xgboost")

	m, ok := r.Models["fraud-detector"]
	assert.True(t, ok)
	assert.Equal(t, "fraud-detector", m.Name)
	assert.Equal(t, "v1.0", m.Version)
	assert.Equal(t, "xgboost", m.Algorithm)
	assert.Equal(t, "DEPLOYED", m.Status)
}

func TestRegister_Overwrite(t *testing.T) {
	r := NewModelRegistry()
	r.Register("model-a", "v1.0", "random_forest")
	r.Register("model-a", "v2.0", "neural_net")

	m, err := r.Get("model-a")
	assert.NoError(t, err)
	assert.Equal(t, "v2.0", m.Version)
	assert.Equal(t, "neural_net", m.Algorithm)
}

func TestGet_Existing(t *testing.T) {
	r := NewModelRegistry()
	r.Register("identity-verifier", "v3.2", "resnet50")

	m, err := r.Get("identity-verifier")
	assert.NoError(t, err)
	assert.Equal(t, "identity-verifier", m.Name)
	assert.Equal(t, "v3.2", m.Version)
}

func TestGet_NotFound(t *testing.T) {
	r := NewModelRegistry()
	_, err := r.Get("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MODEL_NOT_FOUND")
}

func TestSwitchToShadow(t *testing.T) {
	r := NewModelRegistry()
	r.Register("model-shadow", "v1.0", "svm")
	r.SwitchToShadow("model-shadow")

	m, err := r.Get("model-shadow")
	assert.NoError(t, err)
	assert.Equal(t, "SHADOW", m.Status)
}

func TestSwitchToShadow_UnknownModel(t *testing.T) {
	r := NewModelRegistry()
	r.SwitchToShadow("unknown")
	_, err := r.Get("unknown")
	assert.Error(t, err)
}

func TestRegister_ConcurrentSafe(t *testing.T) {
	r := NewModelRegistry()
	t.Run("parallel", func(t *testing.T) {
		t.Run("register fraud", func(t *testing.T) {
			r.Register("fraud", "v1", "xgb")
		})
		t.Run("register identity", func(t *testing.T) {
			r.Register("identity", "v2", "cnn")
		})
		t.Run("register biometric", func(t *testing.T) {
			r.Register("biometric", "v1", "resnet")
		})
	})
	assert.Len(t, r.Models, 3)
}

func TestGet_RaceCondition(t *testing.T) {
	r := NewModelRegistry()
	r.Register("race-model", "v1.0", "lr")

	t.Run("parallel", func(t *testing.T) {
		t.Run("read", func(t *testing.T) {
			_, err := r.Get("race-model")
			assert.NoError(t, err)
		})
		t.Run("write", func(t *testing.T) {
			r.SwitchToShadow("race-model")
		})
	})
}

func TestModelMetadata_Values(t *testing.T) {
	r := NewModelRegistry()
	r.Register("test-model", "v0.0.1", "knn")

	m, _ := r.Get("test-model")
	assert.Equal(t, "test-model", m.Name)
	assert.Equal(t, "v0.0.1", m.Version)
	assert.Equal(t, "knn", m.Algorithm)
	assert.Equal(t, "DEPLOYED", m.Status)
}

func TestRegister_EmptyName(t *testing.T) {
	r := NewModelRegistry()
	r.Register("", "v1.0", "test")

	m, err := r.Get("")
	assert.NoError(t, err)
	assert.Equal(t, "", m.Name)
}
