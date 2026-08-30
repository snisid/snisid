package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/sisal-svc/internal/domain"
)

type alertRepo struct {
	pool *pgxpool.Pool
}

func NewAlertRepo(pool *pgxpool.Pool) *alertRepo {
	return &alertRepo{pool: pool}
}

func (r *alertRepo) Create(alert *domain.SISALAlert) (*domain.SISALAlert, error) {
	ctx := context.Background()
	alert.ID = uuid.New()
	alert.NationalSisalID = "SISAL-HT-" + time.Now().Format("2006") + "-" + alert.ID.String()[:6]
	alert.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sisal_alerts
		 (alert_id, national_sisal_id, hazard_type, severity, title, message_fr, message_ht,
		  affected_depts, affected_pop_est, issued_at, source_agency, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		alert.ID, alert.NationalSisalID, alert.HazardType, alert.Severity,
		alert.Title, alert.MessageFR, alert.MessageHT, alert.AffectedDepts,
		alert.AffectedPopEst, alert.IssuedAt, alert.SourceAgency,
		alert.CreatedBy, alert.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return alert, nil
}

func (r *alertRepo) FindActive() ([]domain.SISALAlert, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT alert_id, national_sisal_id, hazard_type, severity, title, message_fr, message_ht,
		        affected_depts, issued_at, source_agency, created_at
		 FROM sisal_alerts
		 WHERE is_cancelled = FALSE AND (valid_until IS NULL OR valid_until > NOW())
		 ORDER BY issued_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.SISALAlert
	for rows.Next() {
		var a domain.SISALAlert
		if err := rows.Scan(&a.ID, &a.NationalSisalID, &a.HazardType, &a.Severity,
			&a.Title, &a.MessageFR, &a.MessageHT, &a.AffectedDepts,
			&a.IssuedAt, &a.SourceAgency, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r *alertRepo) FindHistory() ([]domain.SISALAlert, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT alert_id, national_sisal_id, hazard_type, severity, title, message_fr, message_ht,
		        affected_depts, issued_at, source_agency, created_at
		 FROM sisal_alerts ORDER BY issued_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.SISALAlert
	for rows.Next() {
		var a domain.SISALAlert
		if err := rows.Scan(&a.ID, &a.NationalSisalID, &a.HazardType, &a.Severity,
			&a.Title, &a.MessageFR, &a.MessageHT, &a.AffectedDepts,
			&a.IssuedAt, &a.SourceAgency, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r *alertRepo) Cancel(id uuid.UUID, reason string) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE sisal_alerts SET is_cancelled = TRUE, cancelled_at = NOW(), cancel_reason = $1 WHERE alert_id = $2`,
		reason, id)
	return err
}

func (r *alertRepo) CreateSubscription(sub *domain.Subscription) (*domain.Subscription, error) {
	ctx := context.Background()
	sub.ID = uuid.New()
	sub.CreatedAt = time.Now()
	isActive := true
	sub.IsActive = &isActive

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sisal_subscriptions
		 (sub_id, phone_number, email, dept_code, commune, min_severity, is_active, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		sub.ID, sub.PhoneNumber, sub.Email, sub.DeptCode, sub.Commune,
		sub.MinSeverity, sub.IsActive, sub.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return sub, nil
}
