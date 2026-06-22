package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/hsm-svc/internal/domain"
)

type Repository interface {
	CreateKey(ctx context.Context, key *domain.HSMKey) error
	FindByKeyID(ctx context.Context, keyID uuid.UUID) (*domain.HSMKey, error)
	FindByAlgorithmAndState(ctx context.Context, algorithm domain.KeyAlgorithm, state domain.KeyState) ([]domain.HSMKey, error)
	FindAll(ctx context.Context) ([]domain.HSMKey, error)
	UpdateState(ctx context.Context, keyID uuid.UUID, state domain.KeyState) error
	UpdateRotatedAt(ctx context.Context, keyID uuid.UUID, rotatedAt time.Time) error
	FindSlotByID(ctx context.Context, slotID int) (*domain.HSMSlot, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateKey(ctx context.Context, key *domain.HSMKey) error {
	query := `INSERT INTO hsm_keys (key_id, key_label, algorithm, key_size, state, usages, slot_id, is_extractable, public_key_pem, key_hash, rotated_at, expires_at, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := r.db.ExecContext(ctx, query,
		key.KeyID, key.KeyLabel, key.Algorithm, key.KeySize, key.State,
		key.Usages, key.SlotID, key.IsExtractable, key.PublicKeyPEM, key.KeyHash,
		key.RotatedAt, key.ExpiresAt, key.CreatedAt, key.UpdatedAt, key.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("insert hsm_key: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindByKeyID(ctx context.Context, keyID uuid.UUID) (*domain.HSMKey, error) {
	query := `SELECT key_id, key_label, algorithm, key_size, state, usages, slot_id, is_extractable, public_key_pem, key_hash, rotated_at, expires_at, created_at, updated_at, created_by
		FROM hsm_keys WHERE key_id = $1`
	key := &domain.HSMKey{}
	err := r.db.QueryRowContext(ctx, query, keyID).Scan(
		&key.KeyID, &key.KeyLabel, &key.Algorithm, &key.KeySize, &key.State,
		&key.Usages, &key.SlotID, &key.IsExtractable, &key.PublicKeyPEM, &key.KeyHash,
		&key.RotatedAt, &key.ExpiresAt, &key.CreatedAt, &key.UpdatedAt, &key.CreatedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("key not found: %s", keyID)
		}
		return nil, fmt.Errorf("query hsm_key: %w", err)
	}
	return key, nil
}

func (r *postgresRepo) FindByAlgorithmAndState(ctx context.Context, algorithm domain.KeyAlgorithm, state domain.KeyState) ([]domain.HSMKey, error) {
	query := `SELECT key_id, key_label, algorithm, key_size, state, usages, slot_id, is_extractable, public_key_pem, key_hash, rotated_at, expires_at, created_at, updated_at, created_by
		FROM hsm_keys WHERE algorithm = $1 AND state = $2 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, algorithm, state)
	if err != nil {
		return nil, fmt.Errorf("query keys by algorithm and state: %w", err)
	}
	defer rows.Close()

	var keys []domain.HSMKey
	for rows.Next() {
		var key domain.HSMKey
		if err := rows.Scan(
			&key.KeyID, &key.KeyLabel, &key.Algorithm, &key.KeySize, &key.State,
			&key.Usages, &key.SlotID, &key.IsExtractable, &key.PublicKeyPEM, &key.KeyHash,
			&key.RotatedAt, &key.ExpiresAt, &key.CreatedAt, &key.UpdatedAt, &key.CreatedBy,
		); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func (r *postgresRepo) FindAll(ctx context.Context) ([]domain.HSMKey, error) {
	query := `SELECT key_id, key_label, algorithm, key_size, state, usages, slot_id, is_extractable, public_key_pem, key_hash, rotated_at, expires_at, created_at, updated_at, created_by
		FROM hsm_keys ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all keys: %w", err)
	}
	defer rows.Close()

	var keys []domain.HSMKey
	for rows.Next() {
		var key domain.HSMKey
		if err := rows.Scan(
			&key.KeyID, &key.KeyLabel, &key.Algorithm, &key.KeySize, &key.State,
			&key.Usages, &key.SlotID, &key.IsExtractable, &key.PublicKeyPEM, &key.KeyHash,
			&key.RotatedAt, &key.ExpiresAt, &key.CreatedAt, &key.UpdatedAt, &key.CreatedBy,
		); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func (r *postgresRepo) UpdateState(ctx context.Context, keyID uuid.UUID, state domain.KeyState) error {
	query := `UPDATE hsm_keys SET state = $1, updated_at = $2 WHERE key_id = $3`
	_, err := r.db.ExecContext(ctx, query, state, time.Now().UTC(), keyID)
	if err != nil {
		return fmt.Errorf("update key state: %w", err)
	}
	return nil
}

func (r *postgresRepo) UpdateRotatedAt(ctx context.Context, keyID uuid.UUID, rotatedAt time.Time) error {
	query := `UPDATE hsm_keys SET rotated_at = $1, state = 'ACTIVE', updated_at = $2 WHERE key_id = $3`
	_, err := r.db.ExecContext(ctx, query, rotatedAt, time.Now().UTC(), keyID)
	if err != nil {
		return fmt.Errorf("update key rotated_at: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindSlotByID(ctx context.Context, slotID int) (*domain.HSMSlot, error) {
	query := `SELECT slot_id, label, manufacturer, model, serial_number, firmware_version, is_logged_in, token_present, hardware_model
		FROM hsm_slots WHERE slot_id = $1`
	slot := &domain.HSMSlot{}
	err := r.db.QueryRowContext(ctx, query, slotID).Scan(
		&slot.SlotID, &slot.Label, &slot.Manufacturer, &slot.Model, &slot.SerialNumber,
		&slot.FirmwareVer, &slot.IsLoggedIn, &slot.TokenPresent, &slot.HardwareModel,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("slot not found: %d", slotID)
		}
		return nil, fmt.Errorf("query hsm_slot: %w", err)
	}
	return slot, nil
}
