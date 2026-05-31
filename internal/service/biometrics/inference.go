package biometrics

import (
	"context"
	"math/rand"
	"time"
)

type BiometricType string

const (
	TypeFace        BiometricType = "FACE"
	TypeFingerprint BiometricType = "FINGERPRINT"
)

type InferenceEngine interface {
	GenerateEmbedding(ctx context.Context, rawData []byte, bType BiometricType) ([]float32, error)
}

type DefaultInferenceEngine struct {
	Endpoint string
	Timeout  time.Duration
}

func NewInferenceEngine(endpoint string) *DefaultInferenceEngine {
	return &DefaultInferenceEngine{
		Endpoint: endpoint,
		Timeout:  5 * time.Second,
	}
}

func (e *DefaultInferenceEngine) GenerateEmbedding(ctx context.Context, rawData []byte, bType BiometricType) ([]float32, error) {
	time.Sleep(50 * time.Millisecond)

	dim := 128
	vec := make([]float32, dim)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < dim; i++ {
		vec[i] = r.Float32()
	}
	return vec, nil
}
