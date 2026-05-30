package biometrics

import (
	"context"
	"fmt"

	"github.com/snisid/platform/backend/internal/platform/logger"
	"go.uber.org/zap"
)

type BiometricService struct {
	milvus    *MilvusBridge
	inference InferenceEngine
}

func NewBiometricService(milvus *MilvusBridge, inference InferenceEngine) *BiometricService {
	return &BiometricService{
		milvus:    milvus,
		inference: inference,
	}
}

func (s *BiometricService) Enroll(ctx context.Context, identityID string, rawData []byte, bType BiometricType) error {
	logger.Info(ctx, "Enrolling biometric profile", zap.String("identity_id", identityID), zap.String("type", string(bType)))

	vector, err := s.inference.GenerateEmbedding(ctx, rawData, bType)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	collection := fmt.Sprintf("snisid_biometrics_%s", bType)
	return s.milvus.InsertBiometric(ctx, collection, identityID, vector)
}

func (s *BiometricService) Verify(ctx context.Context, rawData []byte, bType BiometricType) (string, float32, error) {
	logger.Info(ctx, "Starting biometric verification", zap.String("type", string(bType)))

	vector, err := s.inference.GenerateEmbedding(ctx, rawData, bType)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate verification embedding: %w", err)
	}

	collection := fmt.Sprintf("snisid_biometrics_%s", bType)
	matchID, distance, err := s.milvus.Search(ctx, collection, vector)
	if err != nil {
		return "", 0, fmt.Errorf("biometric search failed: %w", err)
	}

	// Score normalization: L2 distance to confidence percentage
	// Simplified: smaller distance = higher confidence
	confidence := 100.0 - (distance * 10.0)
	if confidence < 0 {
		confidence = 0
	}

	logger.Info(ctx, "Biometric match found", zap.String("match_id", matchID), zap.Float32("confidence", float32(confidence)))

	return matchID, float32(confidence), nil
}
