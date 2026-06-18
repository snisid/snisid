package engine

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/offline-sync-engine/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.OfflineEvent{}))
	return db
}

func TestNewSyncEngine(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)
	assert.NotNil(t, e)
	assert.NotNil(t, e.db)
}

func TestQueueEvent_Success(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	event := &models.OfflineEvent{
		EventType:  "enrollment",
		Payload:    `{"name":"test"}`,
		TerminalID: "terminal-01",
		Priority:   2,
	}
	err := e.QueueEvent(event)
	require.NoError(t, err)
	assert.NotEmpty(t, event.ID)
	assert.Equal(t, "pending", event.Status)
	assert.Equal(t, 0, event.RetryCount)
	assert.Equal(t, 3, event.MaxRetries)
	assert.NotEmpty(t, event.VectorClock)
}

func TestQueueEvent_MissingEventType(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	err := e.QueueEvent(&models.OfflineEvent{
		Payload:    `{}`,
		TerminalID: "term-01",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "event_type is required")
}

func TestQueueEvent_MissingPayload(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	err := e.QueueEvent(&models.OfflineEvent{
		EventType:  "test",
		TerminalID: "term-01",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payload is required")
}

func TestQueueEvent_MissingTerminalID(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	err := e.QueueEvent(&models.OfflineEvent{
		EventType: "test",
		Payload:   `{}`,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "terminal_id is required")
}

func TestQueueEvent_CustomMaxRetries(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	event := &models.OfflineEvent{
		EventType:  "sync",
		Payload:    `{"key":"val"}`,
		TerminalID: "term-02",
		MaxRetries: 5,
	}
	err := e.QueueEvent(event)
	require.NoError(t, err)
	assert.Equal(t, 5, event.MaxRetries)
}

func TestQueueEvent_ExistingVectorClock(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	event := &models.OfflineEvent{
		EventType:   "update",
		Payload:     `{"data":"1"}`,
		TerminalID:  "term-03",
		AggregateID: "agg-001",
		VectorClock: `{"term-03":5}`,
	}
	err := e.QueueEvent(event)
	require.NoError(t, err)

	vc := DeserializeVectorClock(event.VectorClock)
	assert.Equal(t, 6, vc["term-03"])
}

func TestSync_NoEvents(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	result, err := e.Sync()
	require.NoError(t, err)
	assert.Equal(t, 0, result.Total)
	assert.Equal(t, 0, result.Synced)
	assert.Equal(t, 0, result.Conflicts)
}

func TestSync_WithCallback(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	callbackCalled := false
	e.SetSyncCallback(func(event *models.OfflineEvent) error {
		callbackCalled = true
		return nil
	})

	require.NoError(t, e.QueueEvent(&models.OfflineEvent{
		EventType:  "test",
		Payload:    `{"a":1}`,
		TerminalID: "term-01",
	}))

	result, err := e.Sync()
	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, 1, result.Synced)
	assert.True(t, callbackCalled)
}

func TestSync_CallbackError(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	e.SetSyncCallback(func(event *models.OfflineEvent) error {
		return errors.New("callback failure")
	})

	require.NoError(t, e.QueueEvent(&models.OfflineEvent{
		EventType:  "fail-test",
		Payload:    `{"x":1}`,
		TerminalID: "term-01",
	}))

	result, err := e.Sync()
	require.NoError(t, err)
	assert.Equal(t, 1, result.Failed)
	assert.Equal(t, 0, result.Synced)
}

func TestSync_WithoutCallback(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	require.NoError(t, e.QueueEvent(&models.OfflineEvent{
		EventType:  "auto-sync",
		Payload:    `{"ok":true}`,
		TerminalID: "term-01",
	}))

	result, err := e.Sync()
	require.NoError(t, err)
	assert.Equal(t, 1, result.Synced)
	assert.Equal(t, 0, result.Failed)
}

func TestSync_ConflictDetection(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	require.NoError(t, e.QueueEvent(&models.OfflineEvent{
		EventType:   "update",
		Payload:     `{"v":1}`,
		TerminalID:  "term-a",
		AggregateID: "agg-001",
	}))

	// First sync establishes baseline
	e.SetSyncCallback(func(event *models.OfflineEvent) error { return nil })
	_, err := e.Sync()
	require.NoError(t, err)

	// Queue another event with older clock for same aggregate
	require.NoError(t, e.QueueEvent(&models.OfflineEvent{
		EventType:   "update",
		Payload:     `{"v":2}`,
		TerminalID:  "term-b",
		AggregateID: "agg-001",
	}))

	// Second sync
	result, err := e.Sync()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, result.Synced, 1)
}

func TestRemoveEvent_Success(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	event := &models.OfflineEvent{
		ID:         uuid.New().String(),
		EventType:  "delete-test",
		Payload:    `{}`,
		TerminalID: "term-01",
	}
	require.NoError(t, e.QueueEvent(event))

	err := e.RemoveEvent(event.ID)
	require.NoError(t, err)
}

func TestRemoveEvent_NotFound(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	err := e.RemoveEvent("nonexistent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "event not found")
}

func TestGetPendingCount(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	assert.Equal(t, 0, e.GetPendingCount())

	for i := 0; i < 3; i++ {
		require.NoError(t, e.QueueEvent(&models.OfflineEvent{
			EventType:  "pending-test",
			Payload:    `{}`,
			TerminalID: "term-01",
		}))
	}
	assert.Equal(t, 3, e.GetPendingCount())
}

func TestListEvents_Pagination(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	for i := 0; i < 5; i++ {
		require.NoError(t, e.QueueEvent(&models.OfflineEvent{
			EventType:  "list-test",
			Payload:    `{"i":` + string(rune('0'+i)) + `}`,
			TerminalID: "term-01",
		}))
	}

	events, total, err := e.ListEvents("", 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, events, 2)
}

func TestListEvents_InvalidPagination(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	events, total, err := e.ListEvents("", 0, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, events)
}

func TestListEvents_ByStatus(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	require.NoError(t, e.QueueEvent(&models.OfflineEvent{
		EventType:  "filter-test",
		Payload:    `{}`,
		TerminalID: "term-01",
	}))

	events, total, err := e.ListEvents("pending", 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, events, 1)

	events, total, err = e.ListEvents("synced", 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
}

func TestResetStuckEvents(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	require.NoError(t, e.QueueEvent(&models.OfflineEvent{
		EventType:  "stuck-test",
		Payload:    `{}`,
		TerminalID: "term-01",
	}))

	err := e.ResetStuckEvents()
	require.NoError(t, err)
}

func TestGetQueueStatus(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	status, err := e.GetQueueStatus()
	require.NoError(t, err)
	assert.Equal(t, int64(0), status.Total)

	require.NoError(t, e.QueueEvent(&models.OfflineEvent{
		EventType:  "status-test",
		Payload:    `{}`,
		TerminalID: "term-01",
	}))

	status, err = e.GetQueueStatus()
	require.NoError(t, err)
	assert.Equal(t, int64(1), status.Pending)
	assert.Equal(t, int64(1), status.Total)
}

func TestLastSynced_NotFound(t *testing.T) {
	db := setupTestDB(t)
	e := NewSyncEngine(db)

	event := &models.OfflineEvent{}
	err := e.LastSynced(event)
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestVectorClockOperations(t *testing.T) {
	vc1 := DeserializeVectorClock(`{"node-a":3,"node-b":5}`)
	vc2 := DeserializeVectorClock(`{"node-a":3,"node-b":7}`)
	assert.Equal(t, BEFORE, vc1.Compare(vc2))
	assert.Equal(t, AFTER, vc2.Compare(vc1))

	vc3 := DeserializeVectorClock(`{"node-a":4,"node-b":5}`)
	assert.Equal(t, CONCURRENT, vc1.Compare(vc3))

	vc4 := DeserializeVectorClock(`{"node-a":3,"node-b":5}`)
	assert.Equal(t, EQUAL, vc1.Compare(vc4))

	vc5 := make(VectorClock)
	vc6 := make(VectorClock)
	assert.Equal(t, EQUAL, vc5.Compare(vc6))

	vc7 := VectorClock{"a": 1}
	empty := make(VectorClock)
	assert.Equal(t, AFTER, vc7.Compare(empty))
	assert.Equal(t, BEFORE, empty.Compare(vc7))
}

func TestVectorClockMerge(t *testing.T) {
	vc1 := VectorClock{"a": 1, "b": 2}
	vc2 := VectorClock{"b": 3, "c": 4}
	vc1.Merge(vc2)

	assert.Equal(t, 1, vc1["a"])
	assert.Equal(t, 3, vc1["b"])
	assert.Equal(t, 4, vc1["c"])
}

func TestVectorClockSerialize(t *testing.T) {
	vc := VectorClock{"a": 1, "b": 2}
	s := vc.Serialize()
	assert.Contains(t, s, `"a":1`)
	assert.Contains(t, s, `"b":2`)

	var nilVC VectorClock
	assert.Equal(t, "{}", nilVC.Serialize())
}

func TestDeserializeVectorClock_Invalid(t *testing.T) {
	vc := DeserializeVectorClock("invalid json")
	assert.NotNil(t, vc)
	assert.Empty(t, vc)

	vc = DeserializeVectorClock("")
	assert.NotNil(t, vc)
	assert.Empty(t, vc)

	vc = DeserializeVectorClock("{}")
	assert.NotNil(t, vc)
	assert.Empty(t, vc)
}
