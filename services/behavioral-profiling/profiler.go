package behavioralprofiling

import (
	"context"
	"fmt"

	"github.com/snisid/platform/internal/platform/logger"
)

type BehaviorProfile struct {
	SubjectID      string
	AvgTransaction float64
	FrequentLocs   []string
	ActiveHours    []int
}

func AnalyzeBehavior(currentTransaction float64, profile BehaviorProfile) bool {
	logger.Info(context.Background(), fmt.Sprintf("PROFILER: Analyzing behavior for subject %s", profile.SubjectID))

	// Anomaly detection: pattern-of-life deviation
	if currentTransaction > profile.AvgTransaction*3 {
		fmt.Println("ANOMALY: High transaction velocity detected.")
		return true
	}

	return false
}

func UpdateProfile(profile *BehaviorProfile, newTx float64) {
	// Sliding window average update
	profile.AvgTransaction = (profile.AvgTransaction * 0.9) + (newTx * 0.1)
}
