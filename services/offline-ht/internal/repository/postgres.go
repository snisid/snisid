package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/offline-ht/internal/domain"
)

type Repository interface {
	PushQueue(ctx context.Context, item *domain.SyncQueueItem) error
	SyncTerminal(ctx context.Context, terminalID uuid.UUID) ([]domain.SyncQueueItem, error)
	GetConflictItems(ctx context.Context) ([]domain.SyncQueueItem, error)
	UpsertTerminal(ctx context.Context, t *domain.OfflineTerminal) error
	GetTerminalsStatus(ctx context.Context) ([]domain.OfflineTerminal, error)
	UpdateQueueItemStatus(ctx context.Context, id uuid.UUID, status domain.SyncStatus, errMsg *string) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) PushQueue(ctx context.Context, item *domain.SyncQueueItem) error {
	query := `INSERT INTO offline_sync_queue (id, terminal_id, entity_type, entity_id, action, payload, status, retry_count, error_msg, created_at, updated_at, synced_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.TerminalID, item.EntityType, item.EntityID, item.Action,
		item.Payload, item.Status, item.RetryCount, item.ErrorMsg,
		item.CreatedAt, item.UpdatedAt, item.SyncedAt,
	)
	return err
}

func (r *postgresRepo) SyncTerminal(ctx context.Context, terminalID uuid.UUID) ([]domain.SyncQueueItem, error) {
	query := `SELECT id, terminal_id, entity_type, entity_id, action, payload, status, retry_count, error_msg, created_at, updated_at, synced_at
		FROM offline_sync_queue WHERE terminal_id = $1 AND status = 'PENDING' ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, terminalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.SyncQueueItem
	for rows.Next() {
		var item domain.SyncQueueItem
		if err := rows.Scan(
			&item.ID, &item.TerminalID, &item.EntityType, &item.EntityID, &item.Action,
			&item.Payload, &item.Status, &item.RetryCount, &item.ErrorMsg,
			&item.CreatedAt, &item.UpdatedAt, &item.SyncedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *postgresRepo) GetConflictItems(ctx context.Context) ([]domain.SyncQueueItem, error) {
	query := `SELECT id, terminal_id, entity_type, entity_id, action, payload, status, retry_count, error_msg, created_at, updated_at, synced_at
		FROM offline_sync_queue WHERE status = 'CONFLICT' ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.SyncQueueItem
	for rows.Next() {
		var item domain.SyncQueueItem
		if err := rows.Scan(
			&item.ID, &item.TerminalID, &item.EntityType, &item.EntityID, &item.Action,
			&item.Payload, &item.Status, &item.RetryCount, &item.ErrorMsg,
			&item.CreatedAt, &item.UpdatedAt, &item.SyncedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *postgresRepo) UpsertTerminal(ctx context.Context, t *domain.OfflineTerminal) error {
	query := `INSERT INTO offline_terminals (id, name, location, last_sync_at, firmware_ver, is_online, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET name = $2, location = $3, last_sync_at = $4, firmware_ver = $5, is_online = $6, updated_at = $8`
	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.Name, t.Location, t.LastSyncAt, t.FirmwareVer, t.IsOnline,
		t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) GetTerminalsStatus(ctx context.Context) ([]domain.OfflineTerminal, error) {
	query := `SELECT id, name, location, last_sync_at, firmware_ver, is_online, created_at, updated_at
		FROM offline_terminals ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terminals []domain.OfflineTerminal
	for rows.Next() {
		var t domain.OfflineTerminal
		if err := rows.Scan(
			&t.ID, &t.Name, &t.Location, &t.LastSyncAt, &t.FirmwareVer,
			&t.IsOnline, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		terminals = append(terminals, t)
	}
	return terminals, rows.Err()
}

func (r *postgresRepo) UpdateQueueItemStatus(ctx context.Context, id uuid.UUID, status domain.SyncStatus, errMsg *string) error {
	now := time.Now().UTC()
	query := `UPDATE offline_sync_queue SET status = $1, error_msg = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, status, errMsg, now, id)
	return err
}
