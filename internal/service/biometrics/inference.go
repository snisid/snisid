package biometrics

import (
	"context"
	"math/rand"
	"time"
)

type BiometricType string

const (
	TypeFace       BiometricType = "FACE"
	TypeFingerprint BiometricType = "FINGERPRINT"
)

type InferenceEngine interface {
	GenerateEmbedding(ctx context.Context, rawData []byte, bType BiometricType) ([]float32, error)
}

type MockInferenceEngine struct{}

func (e *MockInferenceEngine) GenerateEmbedding(ctx context.Context, rawData []byte, bType BiometricType) ([]float32, error) {
	// Simulation of GPU-accelerated embedding generation
	// In production, this would be a gRPC call to a Triton/TensorRT service
	time.Sleep(150 * time.Millisecond)
	
	dim := 128
	vec := make([]float32, dim)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < dim; i++ {
		vec[i] = r.Float32()
	}
	return vec, nil
}
