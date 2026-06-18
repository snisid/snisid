package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/cybre-svc/internal/domain"
)

type cybreRepo struct {
	pool *pgxpool.Pool
}

func NewCybreRepo(pool *pgxpool.Pool) *cybreRepo {
	return &cybreRepo{pool: pool}
}

func (r *cybreRepo) CreateIncident(incident *domain.CyberIncident) (*domain.CyberIncident, error) {
	ctx := context.Background()
	incident.ID = uuid.New()
	incident.NationalCybreID = "CYBRE-HT-" + time.Now().Format("2006") + "-" + incident.ID.String()[:6]
	incident.CreatedAt = time.Now()
	incident.UpdatedAt = time.Now()
	status := "OPEN"
	incident.Status = &status

	_, err := r.pool.Exec(ctx,
		`INSERT INTO cybre_incidents
		 (incident_id, national_cybre_id, crime_type, severity, status, victim_count,
		  total_financial_loss_usd, incident_date, reported_date, attack_vector, targeted_platform,
		  suspect_phone, suspect_email, case_reference, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
		incident.ID, incident.NationalCybreID, incident.CrimeType, incident.Severity,
		incident.Status, incident.VictimCount, incident.TotalFinancialLossUSD,
		incident.IncidentDate, incident.ReportedDate, incident.AttackVector,
		incident.TargetedPlatform, incident.SuspectPhone, incident.SuspectEmail,
		incident.CaseReference, incident.CreatedBy, incident.CreatedAt, incident.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return incident, nil
}

func (r *cybreRepo) FindByID(id uuid.UUID) (*domain.CyberIncident, error) {
	ctx := context.Background()
	i := &domain.CyberIncident{}
	err := r.pool.QueryRow(ctx,
		`SELECT incident_id, national_cybre_id, crime_type, severity, status, victim_count,
		        total_financial_loss_usd, incident_date, reported_date, attack_vector,
		        targeted_platform, case_reference, created_at
		 FROM cybre_incidents WHERE incident_id = $1`, id).Scan(
		&i.ID, &i.NationalCybreID, &i.CrimeType, &i.Severity, &i.Status,
		&i.VictimCount, &i.TotalFinancialLossUSD, &i.IncidentDate, &i.ReportedDate,
		&i.AttackVector, &i.TargetedPlatform, &i.CaseReference, &i.CreatedAt)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (r *cybreRepo) FindRecentIntrusions() ([]domain.IntrusionAttempt, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT attempt_id, target_system, attack_timestamp, attack_type, source_country, was_successful, created_at
		 FROM cybre_intrusion_attempts
		 WHERE attack_timestamp >= NOW() - INTERVAL '30 days'
		 ORDER BY attack_timestamp DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []domain.IntrusionAttempt
	for rows.Next() {
		var a domain.IntrusionAttempt
		if err := rows.Scan(&a.ID, &a.TargetSystem, &a.AttackTimestamp, &a.AttackType,
			&a.SourceCountry, &a.WasSuccessful, &a.CreatedAt); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, nil
}

func (r *cybreRepo) CreateThreatIndicator(ti *domain.ThreatIndicator) (*domain.ThreatIndicator, error) {
	ctx := context.Background()
	ti.ID = uuid.New()
	ti.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO cybre_threat_intelligence
		 (threat_id, indicator_type, indicator_value, threat_category, confidence_score,
		  source, is_active, first_seen, last_seen, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		ti.ID, ti.IndicatorType, ti.IndicatorValue, ti.ThreatCategory, ti.ConfidenceScore,
		ti.Source, ti.IsActive, ti.FirstSeen, ti.LastSeen, ti.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return ti, nil
}

func (r *cybreRepo) FindActiveIndicator(indicatorType string, value string) (*domain.ThreatIndicator, error) {
	ctx := context.Background()
	ti := &domain.ThreatIndicator{}
	err := r.pool.QueryRow(ctx,
		`SELECT threat_id, indicator_type, indicator_value, threat_category, confidence_score,
		        source, is_active, first_seen, last_seen, created_at
		 FROM cybre_threat_intelligence
		 WHERE indicator_type = $1 AND indicator_value = $2 AND is_active = TRUE`,
		indicatorType, value).Scan(
		&ti.ID, &ti.IndicatorType, &ti.IndicatorValue, &ti.ThreatCategory,
		&ti.ConfidenceScore, &ti.Source, &ti.IsActive, &ti.FirstSeen, &ti.LastSeen, &ti.CreatedAt)
	if err != nil {
		return nil, err
	}
	return ti, nil
}

func (r *cybreRepo) GetStatsByType() ([]domain.CyberStats, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT crime_type, COUNT(*) as count FROM cybre_incidents GROUP BY crime_type ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []domain.CyberStats
	for rows.Next() {
		var s domain.CyberStats
		if err := rows.Scan(&s.CrimeType, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
