package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/snisid/classification-mgmt-ht/internal/domain"
)

type Repository interface {
	CreateRule(ctx context.Context, r *domain.ClassificationRule) error
	GetRulesByDataType(ctx context.Context, dataType string) ([]domain.ClassificationRule, error)
	CreateTag(ctx context.Context, t *domain.DataTag) error
	GetTagByURI(ctx context.Context, uri string) (*domain.DataTag, error)
	CreateAuditLog(ctx context.Context, a *domain.ClassificationAudit) error
	GetRecentAuditLogs(ctx context.Context) ([]domain.ClassificationAudit, error)
	GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateRule(ctx context.Context, rule *domain.ClassificationRule) error {
	query := `INSERT INTO classification_rules
		(id, data_type, sensitivity_level, handling_caveats, dissemination_limit,
		 encryption_required, access_control_mfa, audit_logging, retention_days,
		 destruction_required, created_by, active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := r.db.ExecContext(ctx, query,
		rule.ID, rule.DataType, rule.SensitivityLevel,
		pq.StringArray(rule.HandlingCaveats), rule.DisseminationLimit,
		rule.EncryptionRequired, rule.AccessControlMFA, rule.AuditLogging,
		rule.RetentionDays, rule.DestructionRequired, rule.CreatedBy,
		rule.Active, rule.CreatedAt, rule.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) GetRulesByDataType(ctx context.Context, dataType string) ([]domain.ClassificationRule, error) {
	query := `SELECT id, data_type, sensitivity_level, handling_caveats, dissemination_limit,
		encryption_required, access_control_mfa, audit_logging, retention_days,
		destruction_required, created_by, active, created_at, updated_at
		FROM classification_rules WHERE data_type = $1 AND active = true ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, dataType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []domain.ClassificationRule
	for rows.Next() {
		var rule domain.ClassificationRule
		var caveats pq.StringArray
		if err := rows.Scan(
			&rule.ID, &rule.DataType, &rule.SensitivityLevel, &caveats,
			&rule.DisseminationLimit, &rule.EncryptionRequired, &rule.AccessControlMFA,
			&rule.AuditLogging, &rule.RetentionDays, &rule.DestructionRequired,
			&rule.CreatedBy, &rule.Active, &rule.CreatedAt, &rule.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rule.HandlingCaveats = []string(caveats)
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

func (r *postgresRepo) CreateTag(ctx context.Context, tag *domain.DataTag) error {
	query := `INSERT INTO classification_tags
		(id, resource_uri, classification_top_level, classification_atomic, handling_caveats,
		 owner_agency, tagged_by, tagged_at, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.db.ExecContext(ctx, query,
		tag.ID, tag.ResourceURI, tag.ClassificationTop, tag.ClassificationAtomic,
		pq.StringArray(tag.HandlingCaveats), tag.OwnerAgency,
		tag.TaggedBy, tag.TaggedAt, tag.ExpiresAt,
	)
	return err
}

func (r *postgresRepo) GetTagByURI(ctx context.Context, uri string) (*domain.DataTag, error) {
	query := `SELECT id, resource_uri, classification_top_level, classification_atomic, handling_caveats,
		owner_agency, tagged_by, tagged_at, expires_at
		FROM classification_tags WHERE resource_uri = $1 ORDER BY tagged_at DESC LIMIT 1`
	tag := &domain.DataTag{}
	var caveats pq.StringArray
	err := r.db.QueryRowContext(ctx, query, uri).Scan(
		&tag.ID, &tag.ResourceURI, &tag.ClassificationTop, &tag.ClassificationAtomic,
		&caveats, &tag.OwnerAgency, &tag.TaggedBy, &tag.TaggedAt, &tag.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found for URI")
		}
		return nil, err
	}
	tag.HandlingCaveats = []string(caveats)
	return tag, nil
}

func (r *postgresRepo) CreateAuditLog(ctx context.Context, a *domain.ClassificationAudit) error {
	query := `INSERT INTO classification_audit
		(id, resource_uri, action, from_level, to_level, rationale,
		 authorized_by, classification_authority, timestamp, ip_address)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.ResourceURI, a.Action, a.FromLevel, a.ToLevel, a.Rationale,
		a.AuthorizedBy, a.ClassificationAuthority, a.Timestamp, a.IPAddress,
	)
	return err
}

func (r *postgresRepo) GetRecentAuditLogs(ctx context.Context) ([]domain.ClassificationAudit, error) {
	query := `SELECT id, resource_uri, action, from_level, to_level, rationale,
		authorized_by, classification_authority, timestamp, ip_address
		FROM classification_audit ORDER BY timestamp DESC LIMIT 100`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.ClassificationAudit
	for rows.Next() {
		var l domain.ClassificationAudit
		if err := rows.Scan(
			&l.ID, &l.ResourceURI, &l.Action, &l.FromLevel, &l.ToLevel, &l.Rationale,
			&l.AuthorizedBy, &l.ClassificationAuthority, &l.Timestamp, &l.IPAddress,
		); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func (r *postgresRepo) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM classification_rules`).Scan(&stats.TotalRules)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM classification_rules WHERE active = true`).Scan(&stats.ActiveRules)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM classification_tags`).Scan(&stats.TotalTags)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM classification_audit`).Scan(&stats.TotalAuditLogs)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM classification_tags WHERE classification_top_level IN ('CONFIDENTIAL','SECRET','TOP_SECRET')`).Scan(&stats.ClassifiedCount)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM classification_tags WHERE classification_top_level = 'SECRET'`).Scan(&stats.SecretCount)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM classification_tags WHERE classification_top_level = 'TOP_SECRET'`).Scan(&stats.TopSecretCount)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
