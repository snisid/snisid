package behavioralprofiling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzeBehavior_NormalTransaction(t *testing.T) {
	profile := BehaviorProfile{
		SubjectID:      "SUBJ-001",
		AvgTransaction: 100.0,
		FrequentLocs:   []string{"loc-1", "loc-2"},
		ActiveHours:    []int{8, 9, 10, 14, 15, 16},
	}
	anomaly := AnalyzeBehavior(150.0, profile)
	assert.False(t, anomaly)
}

func TestAnalyzeBehavior_AnomalousTransaction(t *testing.T) {
	profile := BehaviorProfile{
		SubjectID:      "SUBJ-002",
		AvgTransaction: 100.0,
	}
	anomaly := AnalyzeBehavior(400.0, profile)
	assert.True(t, anomaly)
}

func TestAnalyzeBehavior_ExactlyAtThreshold(t *testing.T) {
	profile := BehaviorProfile{
		SubjectID:      "SUBJ-003",
		AvgTransaction: 100.0,
	}
	anomaly := AnalyzeBehavior(300.0, profile)
	assert.False(t, anomaly)
}

func TestAnalyzeBehavior_ZeroAverage(t *testing.T) {
	profile := BehaviorProfile{
		SubjectID:      "SUBJ-004",
		AvgTransaction: 0,
	}
	anomaly := AnalyzeBehavior(1.0, profile)
	assert.True(t, anomaly)
}

func TestUpdateProfile_InitialUpdate(t *testing.T) {
	profile := &BehaviorProfile{
		SubjectID:      "SUBJ-005",
		AvgTransaction: 100.0,
	}
	UpdateProfile(profile, 200.0)
	expected := 100.0*0.9 + 200.0*0.1
	assert.InDelta(t, expected, profile.AvgTransaction, 0.01)
}

func TestUpdateProfile_MultipleUpdates(t *testing.T) {
	profile := &BehaviorProfile{
		SubjectID:      "SUBJ-006",
		AvgTransaction: 100.0,
	}
	UpdateProfile(profile, 200.0)
	UpdateProfile(profile, 50.0)
	expected := (100.0*0.9+200.0*0.1)*0.9 + 50.0*0.1
	assert.InDelta(t, expected, profile.AvgTransaction, 0.01)
}

func TestUpdateProfile_ZeroToNew(t *testing.T) {
	profile := &BehaviorProfile{
		SubjectID:      "SUBJ-007",
		AvgTransaction: 0,
	}
	UpdateProfile(profile, 500.0)
	assert.InDelta(t, 50.0, profile.AvgTransaction, 0.01)
}

func TestUpdateProfile_DoesNotChangeOtherFields(t *testing.T) {
	profile := &BehaviorProfile{
		SubjectID:      "SUBJ-008",
		AvgTransaction: 100.0,
		FrequentLocs:   []string{"loc-a"},
		ActiveHours:    []int{9, 17},
	}
	UpdateProfile(profile, 150.0)
	assert.Equal(t, "SUBJ-008", profile.SubjectID)
	assert.Equal(t, []string{"loc-a"}, profile.FrequentLocs)
	assert.Equal(t, []int{9, 17}, profile.ActiveHours)
}

func TestAnalyzeBehavior_ConsecutiveCalls(t *testing.T) {
	profile := BehaviorProfile{
		SubjectID:      "SUBJ-009",
		AvgTransaction: 100.0,
	}
	assert.False(t, AnalyzeBehavior(250.0, profile))
	assert.False(t, AnalyzeBehavior(299.0, profile))
	assert.True(t, AnalyzeBehavior(301.0, profile))
}
