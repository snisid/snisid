package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/sla-svc/internal/domain"
)

type Repository interface {
	CreateSLA(ctx context.Context, sla *domain.SLA) error
	ListSLAs(ctx context.Context) ([]domain.SLA, error)
	GetSLA(ctx context.Context, slaID uuid.UUID) (*domain.SLA, error)
	CreateSLO(ctx context.Context, slo *domain.SLO) error
	GetSLOs(ctx context.Context, slaID uuid.UUID) ([]domain.SLO, error)
	RecordSLI(ctx context.Context, sli *domain.ServiceLevelIndicator) error
	GetSLIs(ctx context.Context, sloID uuid.UUID, from, to time.Time) ([]domain.ServiceLevelIndicator, error)
	CreateBreach(ctx context.Context, breach *domain.BreachRecord) error
	GetBreaches(ctx context.Context, slaID uuid.UUID) ([]domain.BreachRecord, error)
	CreateEscalationPolicy(ctx context.Context, p *domain.EscalationPolicy) error
	GetUptime(ctx context.Context, slaID uuid.UUID) ([]domain.UptimeWindow, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateSLA(ctx context.Context, sla *domain.SLA) error {
	query := `INSERT INTO slas (sla_id, name, description, owner, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		sla.SLAID, sla.Name, sla.Description, sla.Owner, sla.IsActive, sla.CreatedAt, sla.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert sla: %w", err)
	}
	return nil
}

func (r *postgresRepo) ListSLAs(ctx context.Context) ([]domain.SLA, error) {
	query := `SELECT sla_id, name, description, owner, is_active, created_at, updated_at FROM slas ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query slas: %w", err)
	}
	defer rows.Close()

	var slas []domain.SLA
	for rows.Next() {
		var sla domain.SLA
		if err := rows.Scan(&sla.SLAID, &sla.Name, &sla.Description, &sla.Owner,
			&sla.IsActive, &sla.CreatedAt, &sla.UpdatedAt); err != nil {
			return nil, err
		}
		slas = append(slas, sla)
	}
	return slas, rows.Err()
}

func (r *postgresRepo) GetSLA(ctx context.Context, slaID uuid.UUID) (*domain.SLA, error) {
	query := `SELECT sla_id, name, description, owner, is_active, created_at, updated_at FROM slas WHERE sla_id = $1`
	sla := &domain.SLA{}
	err := r.db.QueryRowContext(ctx, query, slaID).Scan(
		&sla.SLAID, &sla.Name, &sla.Description, &sla.Owner,
		&sla.IsActive, &sla.CreatedAt, &sla.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sla not found: %s", slaID)
		}
		return nil, fmt.Errorf("query sla: %w", err)
	}
	return sla, nil
}

func (r *postgresRepo) CreateSLO(ctx context.Context, slo *domain.SLO) error {
	query := `INSERT INTO slos (slo_id, sla_id, name, target_value, threshold, time_window_days, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		slo.SLOID, slo.SLAID, slo.Name, slo.TargetValue, slo.Threshold, slo.TimeWindowDays, slo.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert slo: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetSLOs(ctx context.Context, slaID uuid.UUID) ([]domain.SLO, error) {
	query := `SELECT slo_id, sla_id, name, target_value, threshold, time_window_days, created_at
		FROM slos WHERE sla_id = $1 ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query, slaID)
	if err != nil {
		return nil, fmt.Errorf("query slos: %w", err)
	}
	defer rows.Close()

	var slos []domain.SLO
	for rows.Next() {
		var slo domain.SLO
		if err := rows.Scan(&slo.SLOID, &slo.SLAID, &slo.Name, &slo.TargetValue,
			&slo.Threshold, &slo.TimeWindowDays, &slo.CreatedAt); err != nil {
			return nil, err
		}
		slos = append(slos, slo)
	}
	return slos, rows.Err()
}

func (r *postgresRepo) RecordSLI(ctx context.Context, sli *domain.ServiceLevelIndicator) error {
	query := `INSERT INTO sli_data (sli_id, slo_id, sla_id, name, value, recorded_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		sli.SLIID, sli.SLOID, sli.SLAID, sli.Name, sli.Value, sli.RecordedAt,
	)
	if err != nil {
		return fmt.Errorf("insert sli: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetSLIs(ctx context.Context, sloID uuid.UUID, from, to time.Time) ([]domain.ServiceLevelIndicator, error) {
	query := `SELECT sli_id, slo_id, sla_id, name, value, recorded_at
		FROM sli_data WHERE slo_id = $1 AND recorded_at >= $2 AND recorded_at <= $3 ORDER BY recorded_at DESC`
	rows, err := r.db.QueryContext(ctx, query, sloID, from, to)
	if err != nil {
		return nil, fmt.Errorf("query slis: %w", err)
	}
	defer rows.Close()

	var slis []domain.ServiceLevelIndicator
	for rows.Next() {
		var sli domain.ServiceLevelIndicator
		if err := rows.Scan(&sli.SLIID, &sli.SLOID, &sli.SLAID, &sli.Name, &sli.Value, &sli.RecordedAt); err != nil {
			return nil, err
		}
		slis = append(slis, sli)
	}
	return slis, rows.Err()
}

func (r *postgresRepo) CreateBreach(ctx context.Context, breach *domain.BreachRecord) error {
	query := `INSERT INTO breach_records (breach_id, sla_id, slo_id, sli_value, threshold, detected_at, resolved_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		breach.BreachID, breach.SLAID, breach.SLOID, breach.SLIValue, breach.Threshold,
		breach.DetectedAt, breach.ResolvedAt, breach.IsActive,
	)
	if err != nil {
		return fmt.Errorf("insert breach: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetBreaches(ctx context.Context, slaID uuid.UUID) ([]domain.BreachRecord, error) {
	query := `SELECT breach_id, sla_id, slo_id, sli_value, threshold, detected_at, resolved_at, is_active
		FROM breach_records WHERE sla_id = $1 ORDER BY detected_at DESC`
	rows, err := r.db.QueryContext(ctx, query, slaID)
	if err != nil {
		return nil, fmt.Errorf("query breaches: %w", err)
	}
	defer rows.Close()

	var breaches []domain.BreachRecord
	for rows.Next() {
		var b domain.BreachRecord
		if err := rows.Scan(&b.BreachID, &b.SLAID, &b.SLOID, &b.SLIValue, &b.Threshold,
			&b.DetectedAt, &b.ResolvedAt, &b.IsActive); err != nil {
			return nil, err
		}
		breaches = append(breaches, b)
	}
	return breaches, rows.Err()
}

func (r *postgresRepo) CreateEscalationPolicy(ctx context.Context, p *domain.EscalationPolicy) error {
	query := `INSERT INTO escalation_policies (policy_id, sla_id, escalate_after, notify_channel, notify_target, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		p.PolicyID, p.SLAID, p.EscalateAfter, p.NotifyChannel, p.NotifyTarget, p.IsActive,
	)
	if err != nil {
		return fmt.Errorf("insert escalation policy: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetUptime(ctx context.Context, slaID uuid.UUID) ([]domain.UptimeWindow, error) {
	query := `SELECT window_id, sla_id, start_time, end_time, is_up, duration_ms
		FROM uptime_windows WHERE sla_id = $1 ORDER BY start_time DESC`
	rows, err := r.db.QueryContext(ctx, query, slaID)
	if err != nil {
		return nil, fmt.Errorf("query uptime: %w", err)
	}
	defer rows.Close()

	var windows []domain.UptimeWindow
	for rows.Next() {
		var w domain.UptimeWindow
		if err := rows.Scan(&w.WindowID, &w.SLAID, &w.StartTime, &w.EndTime, &w.IsUp, &w.DurationMs); err != nil {
			return nil, err
		}
		windows = append(windows, w)
	}
	return windows, rows.Err()
}
