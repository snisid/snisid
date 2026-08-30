package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

type InterpolSyncRepo struct {
	db *sqlx.DB
}

func NewInterpolSyncRepo(db *sqlx.DB) *InterpolSyncRepo {
	return &InterpolSyncRepo{db: db}
}

func (r *InterpolSyncRepo) Create(ctx context.Context, log *domain.InterpolSyncLog) error {
	query := `
		INSERT INTO sivc_interpol_sync_log (
			sync_id, alert_id, stolen_plate_id, interpol_smv_id, interpol_sad_id,
			sync_direction, sync_status, sync_timestamp, retry_count,
			request_payload, response_payload, error_code, error_message,
			processed_by, processed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		log.SyncID, log.AlertID, log.StolenPlateID, log.InterpolSMVID,
		log.InterpolSADID, log.SyncDirection, log.SyncStatus, log.SyncTimestamp,
		log.RetryCount, log.RequestPayload, log.ResponsePayload, log.ErrorCode,
		log.ErrorMessage, log.ProcessedBy, log.ProcessedAt,
	)
	return err
}

func (r *InterpolSyncRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.SyncStatus, response interface{}) error {
	query := `UPDATE sivc_interpol_sync_log SET sync_status = $1, response_payload = $2, processed_at = NOW() WHERE sync_id = $3`
	_, err := r.db.ExecContext(ctx, query, status, response, id)
	return err
}

func (r *InterpolSyncRepo) FindPending(ctx context.Context) ([]*domain.InterpolSyncLog, error) {
	var logs []*domain.InterpolSyncLog
	query := `SELECT * FROM sivc_interpol_sync_log WHERE sync_status = 'PENDING' ORDER BY sync_timestamp ASC LIMIT 50`
	if err := r.db.SelectContext(ctx, &logs, query); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return logs, nil
}
