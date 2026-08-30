package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/audit-svc/internal/domain"
)

type Repository interface {
	InsertEvent(ctx context.Context, evt *domain.AuditEvent) error
	GetEvent(ctx context.Context, eventID uuid.UUID) (*domain.AuditEvent, error)
	SearchEvents(ctx context.Context, q domain.AuditQuery) ([]domain.AuditEvent, error)
	InsertImmutableEntry(ctx context.Context, entry *domain.ImmutableEntry) error
	GetLastHash(ctx context.Context) (string, error)
	GetStats(ctx context.Context) (map[string]int, error)
	ApplyRetention(ctx context.Context, policy domain.RetentionPolicy) (int, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) InsertEvent(ctx context.Context, evt *domain.AuditEvent) error {
	lastHash, _ := r.GetLastHash(ctx)
	evt.PrevHash = lastHash

	hashInput, _ := json.Marshal(evt)
	h := sha256.Sum256(hashInput)
	evt.Hash = hex.EncodeToString(h[:])

	query := `INSERT INTO audit_events (event_id, source, event_type, category, actor_id, resource_id, action, payload, hash, prev_hash, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		evt.EventID, evt.Source, evt.EventType, evt.Category, evt.ActorID,
		evt.ResourceID, evt.Action, evt.Payload, evt.Hash, evt.PrevHash, evt.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("insert audit event: %w", err)
	}

	entry := &domain.ImmutableEntry{
		EntryID:   uuid.New(),
		EventID:   evt.EventID,
		Hash:      evt.Hash,
		PrevHash:  evt.PrevHash,
		Data:      map[string]any{"event_id": evt.EventID, "hash": evt.Hash},
		CreatedAt: time.Now().UTC(),
	}
	return r.InsertImmutableEntry(ctx, entry)
}

func (r *postgresRepo) GetEvent(ctx context.Context, eventID uuid.UUID) (*domain.AuditEvent, error) {
	query := `SELECT event_id, source, event_type, category, actor_id, resource_id, action, payload, hash, prev_hash, timestamp
		FROM audit_events WHERE event_id = $1`
	e := &domain.AuditEvent{}
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(
		&e.EventID, &e.Source, &e.EventType, &e.Category, &e.ActorID,
		&e.ResourceID, &e.Action, &e.Payload, &e.Hash, &e.PrevHash, &e.Timestamp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found: %s", eventID)
		}
		return nil, fmt.Errorf("query event: %w", err)
	}
	return e, nil
}

func (r *postgresRepo) SearchEvents(ctx context.Context, q domain.AuditQuery) ([]domain.AuditEvent, error) {
	query := `SELECT event_id, source, event_type, category, actor_id, resource_id, action, payload, hash, prev_hash, timestamp
		FROM audit_events WHERE 1=1`
	args := []any{}
	argIdx := 1

	if q.Source != nil {
		query += fmt.Sprintf(" AND source = $%d", argIdx)
		args = append(args, *q.Source)
		argIdx++
	}
	if q.EventType != nil {
		query += fmt.Sprintf(" AND event_type = $%d", argIdx)
		args = append(args, *q.EventType)
		argIdx++
	}
	if q.Category != nil {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, *q.Category)
		argIdx++
	}
	if q.ActorID != nil {
		query += fmt.Sprintf(" AND actor_id = $%d", argIdx)
		args = append(args, *q.ActorID)
		argIdx++
	}
	if q.From != nil {
		query += fmt.Sprintf(" AND timestamp >= $%d", argIdx)
		args = append(args, *q.From)
		argIdx++
	}
	if q.To != nil {
		query += fmt.Sprintf(" AND timestamp <= $%d", argIdx)
		args = append(args, *q.To)
		argIdx++
	}
	query += " ORDER BY timestamp DESC"
	if q.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, q.Limit)
		argIdx++
	} else {
		query += " LIMIT 100"
	}
	if q.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, q.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("search events: %w", err)
	}
	defer rows.Close()

	var events []domain.AuditEvent
	for rows.Next() {
		var e domain.AuditEvent
		if err := rows.Scan(&e.EventID, &e.Source, &e.EventType, &e.Category, &e.ActorID,
			&e.ResourceID, &e.Action, &e.Payload, &e.Hash, &e.PrevHash, &e.Timestamp); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *postgresRepo) InsertImmutableEntry(ctx context.Context, entry *domain.ImmutableEntry) error {
	query := `INSERT INTO immutable_audit_chain (entry_id, event_id, hash, prev_hash, data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		entry.EntryID, entry.EventID, entry.Hash, entry.PrevHash, entry.Data, entry.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert immutable entry: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetLastHash(ctx context.Context) (string, error) {
	query := `SELECT hash FROM immutable_audit_chain ORDER BY created_at DESC LIMIT 1`
	var hash string
	err := r.db.QueryRowContext(ctx, query).Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("get last hash: %w", err)
	}
	return hash, nil
}

func (r *postgresRepo) GetStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)
	query := `SELECT source, COUNT(*) FROM audit_events GROUP BY source`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var source string
		var count int
		if err := rows.Scan(&source, &count); err != nil {
			return nil, err
		}
		stats[source] = count
	}
	return stats, rows.Err()
}

func (r *postgresRepo) ApplyRetention(ctx context.Context, policy domain.RetentionPolicy) (int, error) {
	cutoff := time.Now().AddDate(0, 0, -policy.RetentionDays)
	query := `DELETE FROM audit_events WHERE category = $1 AND timestamp < $2`
	result, err := r.db.ExecContext(ctx, query, policy.Category, cutoff)
	if err != nil {
		return 0, fmt.Errorf("apply retention: %w", err)
	}
	affected, _ := result.RowsAffected()
	return int(affected), nil
}
