package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/api-ht/internal/domain"
)

type Repository interface {
	CreateAccount(ctx context.Context, acc *domain.DeveloperAccount) error
	FindAccountByEmail(ctx context.Context, email string) (*domain.DeveloperAccount, error)
	CreateAPIKey(ctx context.Context, key *domain.APIKey) error
	FindKeyByValue(ctx context.Context, keyValue string) (*domain.APIKey, error)
	FindKeysByAccount(ctx context.Context, accountID uuid.UUID) ([]domain.APIKey, error)
	RevokeKey(ctx context.Context, id uuid.UUID) error
	FindKeyByID(ctx context.Context, id uuid.UUID) (*domain.APIKey, error)
	ListCatalog(ctx context.Context) ([]domain.APIEndpoint, error)
	ListUsageByKey(ctx context.Context, keyID uuid.UUID) ([]domain.UsageLog, error)
	LogUsage(ctx context.Context, log *domain.UsageLog) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateAccount(ctx context.Context, acc *domain.DeveloperAccount) error {
	q := `INSERT INTO api_developer_accounts (id, email, org_name, contact_name, contact_phone, is_approved, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, q,
		acc.ID, acc.Email, acc.OrgName, acc.ContactName, acc.ContactPhone,
		acc.IsApproved, acc.CreatedAt, acc.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) FindAccountByEmail(ctx context.Context, email string) (*domain.DeveloperAccount, error) {
	q := `SELECT id, email, org_name, contact_name, contact_phone, is_approved, created_at, updated_at
		FROM api_developer_accounts WHERE email = $1`
	acc := &domain.DeveloperAccount{}
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&acc.ID, &acc.Email, &acc.OrgName, &acc.ContactName, &acc.ContactPhone,
		&acc.IsApproved, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (r *postgresRepo) CreateAPIKey(ctx context.Context, key *domain.APIKey) error {
	q := `INSERT INTO api_keys (id, account_id, key_value, description, is_active, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, q,
		key.ID, key.AccountID, key.KeyValue, key.Description, key.IsActive, key.ExpiresAt, key.CreatedAt,
	)
	return err
}

func (r *postgresRepo) FindKeyByValue(ctx context.Context, keyValue string) (*domain.APIKey, error) {
	q := `SELECT id, account_id, key_value, description, is_active, expires_at, created_at, revoked_at
		FROM api_keys WHERE key_value = $1`
	k := &domain.APIKey{}
	err := r.db.QueryRowContext(ctx, q, keyValue).Scan(
		&k.ID, &k.AccountID, &k.KeyValue, &k.Description, &k.IsActive, &k.ExpiresAt, &k.CreatedAt, &k.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (r *postgresRepo) FindKeysByAccount(ctx context.Context, accountID uuid.UUID) ([]domain.APIKey, error) {
	q := `SELECT id, account_id, key_value, description, is_active, expires_at, created_at, revoked_at
		FROM api_keys WHERE account_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []domain.APIKey
	for rows.Next() {
		var k domain.APIKey
		if err := rows.Scan(&k.ID, &k.AccountID, &k.KeyValue, &k.Description, &k.IsActive, &k.ExpiresAt, &k.CreatedAt, &k.RevokedAt); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, rows.Err()
}

func (r *postgresRepo) RevokeKey(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	q := `UPDATE api_keys SET is_active = false, revoked_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, now, id)
	return err
}

func (r *postgresRepo) FindKeyByID(ctx context.Context, id uuid.UUID) (*domain.APIKey, error) {
	q := `SELECT id, account_id, key_value, description, is_active, expires_at, created_at, revoked_at
		FROM api_keys WHERE id = $1`
	k := &domain.APIKey{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&k.ID, &k.AccountID, &k.KeyValue, &k.Description, &k.IsActive, &k.ExpiresAt, &k.CreatedAt, &k.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (r *postgresRepo) ListCatalog(ctx context.Context) ([]domain.APIEndpoint, error) {
	q := `SELECT id, path, method, description, sensitivity, module_source, base_path, is_active, version, created_at, updated_at
		FROM api_catalog WHERE is_active = true ORDER BY base_path, path`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eps []domain.APIEndpoint
	for rows.Next() {
		var ep domain.APIEndpoint
		if err := rows.Scan(&ep.ID, &ep.Path, &ep.Method, &ep.Description, &ep.Sensitivity, &ep.ModuleSource, &ep.BasePath, &ep.IsActive, &ep.Version, &ep.CreatedAt, &ep.UpdatedAt); err != nil {
			return nil, err
		}
		eps = append(eps, ep)
	}
	return eps, rows.Err()
}

func (r *postgresRepo) ListUsageByKey(ctx context.Context, keyID uuid.UUID) ([]domain.UsageLog, error) {
	q := `SELECT id, key_id, endpoint, method, status, latency_ms, ip_address, user_agent, created_at
		FROM api_usage_log WHERE key_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, keyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.UsageLog
	for rows.Next() {
		var l domain.UsageLog
		if err := rows.Scan(&l.ID, &l.KeyID, &l.Endpoint, &l.Method, &l.Status, &l.LatencyMs, &l.IPAddress, &l.UserAgent, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func (r *postgresRepo) LogUsage(ctx context.Context, logEntry *domain.UsageLog) error {
	q := `INSERT INTO api_usage_log (id, key_id, endpoint, method, status, latency_ms, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, q,
		logEntry.ID, logEntry.KeyID, logEntry.Endpoint, logEntry.Method, logEntry.Status,
		logEntry.LatencyMs, logEntry.IPAddress, logEntry.UserAgent, logEntry.CreatedAt,
	)
	return err
}
