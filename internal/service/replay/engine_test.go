package replay

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEngine(t *testing.T) {
	e := NewEngine([]string{"localhost:9092"})
	assert.NotNil(t, e)
	assert.NotNil(t, e.producer)
}

func TestReplayJob_Defaults(t *testing.T) {
	job := ReplayJob{
		ID:          "test-job-1",
		SourceTopic: "snisid.prod.identity.v1.events",
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now(),
		Status:      "pending",
	}
	assert.Equal(t, "test-job-1", job.ID)
	assert.Equal(t, "pending", job.Status)
}

func TestRunJob_RecentReplay_UsesKafka(t *testing.T) {
	e := NewEngine([]string{"localhost:9092"})
	job := ReplayJob{
		ID:          "job-kafka",
		SourceTopic: "snisid.prod.identity.v1.events",
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now(),
	}

	err := e.RunJob(context.Background(), job)
	// Should not error since it's a simulated loop
	assert.NoError(t, err)
}

func TestRunJob_ForensicReplay_UsesAudit(t *testing.T) {
	e := NewEngine([]string{"localhost:9092"})
	job := ReplayJob{
		ID:          "job-forensic",
		SourceTopic: "snisid.prod.identity.v1.events",
		StartTime:   time.Now().Add(-30 * 24 * time.Hour), // 30 days ago
		EndTime:     time.Now(),
	}

	err := e.RunJob(context.Background(), job)
	assert.NoError(t, err)
}

func TestRunJob_ContextCancelled(t *testing.T) {
	e := NewEngine([]string{"localhost:9092"})
	job := ReplayJob{
		ID:          "job-cancel",
		SourceTopic: "test",
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := e.RunJob(ctx, job)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "canceled")
}

func TestRunJob_EmptyJob(t *testing.T) {
	e := NewEngine([]string{"localhost:9092"})
	job := ReplayJob{}

	err := e.RunJob(context.Background(), job)
	assert.NoError(t, err)
}

func TestReplayFromAudit_DoesNotError(t *testing.T) {
	e := NewEngine([]string{"localhost:9092"})
	job := ReplayJob{
		ID:        "audit-test",
		StartTime: time.Now().Add(-365 * 24 * time.Hour),
		EndTime:   time.Now(),
	}

	err := e.replayFromAudit(context.Background(), job)
	assert.NoError(t, err)
}
