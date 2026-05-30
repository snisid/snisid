package ml

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

type FeatureVector struct {
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Velocity  float64 `json:"velocity"`
	GraphRisk float64 `json:"graph_risk"`
	Timestamp int64   `json:"timestamp"`
}

type FeatureExtractor struct{}

func (e *FeatureExtractor) Extract(payload map[string]interface{}) FeatureVector {
	userID, _ := payload["user_id"].(string)
	amount, _ := payload["amount"].(float64)
	
	vector := FeatureVector{
		UserID:    userID,
		Amount:    amount,
		Velocity:  0.88, // Mock: computed from window
		GraphRisk: 0.42, // Mock: fetched from Neo4j cache
		Timestamp: time.Now().Unix(),
	}

	logger.Info(fmt.Sprintf("NEXUS-ML: Extracted features for user %s", userID))
	return vector
}

func (e *FeatureExtractor) SaveOnline(v FeatureVector) {
	data, _ := json.Marshal(v)
	fmt.Printf("📦 NEXUS-ML: Saving online feature vector to Redis for %s\n", v.UserID)
	// redis.Set(ctx, "fv:"+v.UserID, data, 0)
	_ = data
}
