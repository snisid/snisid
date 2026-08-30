package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/internal/domain/audit/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuditDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&entity.AuditEvent{})
	require.NoError(t, err)
	return db
}

func makeEvent(seq int64, prevHash, payload string) *entity.AuditEvent {
	return &entity.AuditEvent{
		EventID:       uuid.NewString(),
		CorrelationID: "corr-" + uuid.NewString(),
		EventType:     "test.event",
		Actor:         "system",
		Action:        "test",
		Resource:      "test:resource",
		Status:        "success",
		Payload:       payload,
		PreviousHash:  prevHash,
		Hash:          "hash-" + uuid.NewString(),
		SequenceID:    seq,
		Timestamp:     time.Now().UTC(),
	}
}

func TestPostgresAuditRepository_Append(t *testing.T) {
	db := setupAuditDB(t)
	repo := NewPostgresAuditRepository(db)

	evt := makeEvent(1, "genesis-hash-snisid", `{"action":"create"}`)
	err := repo.Append(context.Background(), evt)
	require.NoError(t, err)

	var saved entity.AuditEvent
	err = db.First(&saved, "event_id = ?", evt.EventID).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), saved.SequenceID)
	assert.Equal(t, "test.event", saved.EventType)
}

func TestPostgresAuditRepository_GetLastEvent_Empty(t *testing.T) {
	db := setupAuditDB(t)
	repo := NewPostgresAuditRepository(db)

	evt, err := repo.GetLastEvent(context.Background())
	require.NoError(t, err)
	assert.Nil(t, evt)
}

func TestPostgresAuditRepository_GetLastEvent_WithData(t *testing.T) {
	db := setupAuditDB(t)
	repo := NewPostgresAuditRepository(db)

	evt1 := makeEvent(1, "genesis", `{"action":"create"}`)
	evt2 := makeEvent(2, "hash-1", `{"action":"update"}`)
	err := repo.Append(context.Background(), evt1)
	require.NoError(t, err)
	err = repo.Append(context.Background(), evt2)
	require.NoError(t, err)

	last, err := repo.GetLastEvent(context.Background())
	require.NoError(t, err)
	require.NotNil(t, last)
	assert.Equal(t, int64(2), last.SequenceID)
	assert.Equal(t, evt2.EventID, last.EventID)
}

func TestPostgresAuditRepository_GetEventsByCorrelationID(t *testing.T) {
	db := setupAuditDB(t)
	repo := NewPostgresAuditRepository(db)

	corrID := "corr-test-123"
	for i := 0; i < 3; i++ {
		evt := makeEvent(int64(i+1), "prev", `{"seq":`+string(rune('0'+i))+`}`)
		evt.CorrelationID = corrID
		err := repo.Append(context.Background(), evt)
		require.NoError(t, err)
	}

	events, err := repo.GetEventsByCorrelationID(context.Background(), corrID)
	require.NoError(t, err)
	assert.Len(t, events, 3)

	events, err = repo.GetEventsByCorrelationID(context.Background(), "nonexistent")
	require.NoError(t, err)
	assert.Empty(t, events)
}

func TestPostgresAuditRepository_GetEventsBySequenceRange(t *testing.T) {
	db := setupAuditDB(t)
	repo := NewPostgresAuditRepository(db)

	for i := 1; i <= 5; i++ {
		evt := makeEvent(int64(i), "prev", `{"seq":`+string(rune('0'+i))+`}`)
		err := repo.Append(context.Background(), evt)
		require.NoError(t, err)
	}

	events, err := repo.GetEventsBySequenceRange(context.Background(), 2, 4)
	require.NoError(t, err)
	assert.Len(t, events, 3)
	assert.Equal(t, int64(2), events[0].SequenceID)
	assert.Equal(t, int64(4), events[2].SequenceID)

	events, err = repo.GetEventsBySequenceRange(context.Background(), 10, 20)
	require.NoError(t, err)
	assert.Empty(t, events)
}
