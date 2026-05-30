package forensics

import (
	"context"
	"fmt"

	"github.com/snisid/platform/backend/internal/platform/logger"
	"go.uber.org/zap"
)

type DeepfakeDetector struct {
	engine ForensicEngine
}

func NewDeepfakeDetector(engine ForensicEngine) *DeepfakeDetector {
	return &DeepfakeDetector{engine: engine}
}

func (d *DeepfakeDetector) Detect(ctx context.Context, identityID string, mediaData []byte) (float32, []string, error) {
	logger.Info(ctx, "Starting deepfake analysis", zap.String("identity_id", identityID))

	// 1. Preprocessing (Mock: in prod use OpenCV/ffmpeg to extract faces)
	
	// 2. Inference
	prob, anomalies, err := d.engine.Analyze(ctx, mediaData)
	if err != nil {
		return 0, nil, fmt.Errorf("forensic inference failed: %w", err)
	}

	logger.Info(ctx, "Deepfake analysis complete", 
		zap.String("identity_id", identityID), 
		zap.Float32("probability", prob),
	)

	return prob, anomalies, nil
}
