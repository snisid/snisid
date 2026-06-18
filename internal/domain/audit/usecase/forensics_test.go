package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/snisid/platform/internal/domain/audit/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAuditRepo struct {
	getEventsBySeqRangeFn func(ctx context.Context, startSeq, endSeq int64) ([]entity.AuditEvent, error)
	getByCorrelationIDFn  func(ctx context.Context, correlationID string) ([]entity.AuditEvent, error)
}

func (m *mockAuditRepo) GetEventsBySequenceRange(ctx context.Context, startSeq, endSeq int64) ([]entity.AuditEvent, error) {
	if m.getEventsBySeqRangeFn != nil {
		return m.getEventsBySeqRangeFn(ctx, startSeq, endSeq)
	}
	return nil, nil
}

func (m *mockAuditRepo) GetEventsByCorrelationID(ctx context.Context, correlationID string) ([]entity.AuditEvent, error) {
	if m.getByCorrelationIDFn != nil {
		return m.getByCorrelationIDFn(ctx, correlationID)
	}
	return nil, nil
}

func TestVerifyIntegrity_EmptyRange(t *testing.T) {
	svc := NewForensicsService(&mockAuditRepo{
		getEventsBySeqRangeFn: func(ctx context.Context, startSeq, endSeq int64) ([]entity.AuditEvent, error) {
			return nil, nil
		},
	})
	valid, err := svc.VerifyIntegrity(context.Background(), 1, 100)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestVerifyIntegrity_SingleEvent(t *testing.T) {
	svc := NewForensicsService(&mockAuditRepo{
		getEventsBySeqRangeFn: func(ctx context.Context, startSeq, endSeq int64) ([]entity.AuditEvent, error) {
			return []entity.AuditEvent{
				{SequenceID: 1, Hash: "hash-1", PreviousHash: "", Payload: "{}"},
			}, nil
		},
	})
	valid, err := svc.VerifyIntegrity(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestVerifyIntegrity_ValidChain(t *testing.T) {
	svc := NewForensicsService(&mockAuditRepo{
		getEventsBySeqRangeFn: func(ctx context.Context, startSeq, endSeq int64) ([]entity.AuditEvent, error) {
			return []entity.AuditEvent{
				{SequenceID: 1, Hash: "h1", PreviousHash: "", Payload: "{}"},
				{SequenceID: 2, Hash: "h2", PreviousHash: "h1", Payload: `{"action":"login"}`},
				{SequenceID: 3, Hash: "h3", PreviousHash: "h2", Payload: `{"action":"update"}`},
			}, nil
		},
	})
	valid, err := svc.VerifyIntegrity(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestVerifyIntegrity_BrokenChain(t *testing.T) {
	svc := NewForensicsService(&mockAuditRepo{
		getEventsBySeqRangeFn: func(ctx context.Context, startSeq, endSeq int64) ([]entity.AuditEvent, error) {
			return []entity.AuditEvent{
				{SequenceID: 1, Hash: "h1", PreviousHash: "", Payload: "{}"},
				{SequenceID: 2, Hash: "h2", PreviousHash: "wrong-hash", Payload: `{}`},
			}, nil
		},
	})
	valid, err := svc.VerifyIntegrity(context.Background(), 1, 10)
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestVerifyIntegrity_RepoError(t *testing.T) {
	svc := NewForensicsService(&mockAuditRepo{
		getEventsBySeqRangeFn: func(ctx context.Context, startSeq, endSeq int64) ([]entity.AuditEvent, error) {
			return nil, errors.New("db error")
		},
	})
	valid, err := svc.VerifyIntegrity(context.Background(), 1, 10)
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestQueryByCorrelationID_Success(t *testing.T) {
	expected := []entity.AuditEvent{
		{EventID: "e1", CorrelationID: "corr-1", Action: "login"},
	}
	svc := NewForensicsService(&mockAuditRepo{
		getByCorrelationIDFn: func(ctx context.Context, correlationID string) ([]entity.AuditEvent, error) {
			return expected, nil
		},
	})
	events, err := svc.QueryByCorrelationID(context.Background(), "corr-1")
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "e1", events[0].EventID)
}

func TestQueryByCorrelationID_Empty(t *testing.T) {
	svc := NewForensicsService(&mockAuditRepo{
		getByCorrelationIDFn: func(ctx context.Context, correlationID string) ([]entity.AuditEvent, error) {
			return nil, nil
		},
	})
	events, err := svc.QueryByCorrelationID(context.Background(), "nonexistent")
	require.NoError(t, err)
	assert.Empty(t, events)
}

func TestQueryByCorrelationID_Error(t *testing.T) {
	svc := NewForensicsService(&mockAuditRepo{
		getByCorrelationIDFn: func(ctx context.Context, correlationID string) ([]entity.AuditEvent, error) {
			return nil, errors.New("query failed")
		},
	})
	events, err := svc.QueryByCorrelationID(context.Background(), "corr-1")
	assert.Error(t, err)
	assert.Nil(t, events)
}
