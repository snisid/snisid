package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/snisid/executive-protection-ht/internal/domain"
)

type Repository interface {
	CreateProtectee(ctx context.Context, p *domain.Protectee) error
	GetActiveProtectees(ctx context.Context) ([]domain.Protectee, error)
	CreateMovementPlan(ctx context.Context, m *domain.MovementPlan) error
	GetUpcomingMovements(ctx context.Context) ([]domain.MovementPlan, error)
	CreateThreatAssessment(ctx context.Context, t *domain.ThreatAssessment) error
	GetActiveThreatsByProtectee(ctx context.Context, protecteeID uuid.UUID) ([]domain.ThreatAssessment, error)
	GetDashboard(ctx context.Context) (*domain.DashboardProtection, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateProtectee(ctx context.Context, p *domain.Protectee) error {
	query := `INSERT INTO execprot_protectees (id, full_name, official_title, protection_level, risk_assessment, active_threats, primary_agent_id, secondary_agents, secure_vehicle_plate, residence_location, workplace_location, daily_schedule_refs, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.FullName, p.OfficialTitle, p.ProtectionLevel, p.RiskAssessment,
		p.ActiveThreats, p.PrimaryAgentID, pq.Array(p.SecondaryAgents),
		p.SecureVehiclePlate, p.ResidenceLocation, p.WorkplaceLocation,
		pq.Array(p.DailyScheduleRefs), time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert protectee: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetActiveProtectees(ctx context.Context) ([]domain.Protectee, error) {
	query := `SELECT id, full_name, official_title, protection_level, risk_assessment, active_threats, primary_agent_id, secondary_agents, secure_vehicle_plate, residence_location, workplace_location, daily_schedule_refs, created_at
		FROM execprot_protectees WHERE risk_assessment IN ('HIGH', 'CRITICAL') ORDER BY active_threats DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query active protectees: %w", err)
	}
	defer rows.Close()

	var protectees []domain.Protectee
	for rows.Next() {
		var p domain.Protectee
		if err := rows.Scan(
			&p.ID, &p.FullName, &p.OfficialTitle, &p.ProtectionLevel, &p.RiskAssessment,
			&p.ActiveThreats, &p.PrimaryAgentID, pq.Array(&p.SecondaryAgents),
			&p.SecureVehiclePlate, &p.ResidenceLocation, &p.WorkplaceLocation,
			pq.Array(&p.DailyScheduleRefs), &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		protectees = append(protectees, p)
	}
	return protectees, rows.Err()
}

func (r *postgresRepo) CreateMovementPlan(ctx context.Context, m *domain.MovementPlan) error {
	query := `INSERT INTO execprot_movement_plans (id, protectee_id, event_name, date, departure_location, arrival_location, transport_mode, route_plan, advance_done, cleared_by, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.ProtecteeID, m.EventName, m.Date, m.DepartureLocation, m.ArrivalLocation,
		m.TransportMode, m.RoutePlan, m.AdvanceDone, m.ClearedBy, m.Status, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert movement plan: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetUpcomingMovements(ctx context.Context) ([]domain.MovementPlan, error) {
	query := `SELECT id, protectee_id, event_name, date, departure_location, arrival_location, transport_mode, route_plan, advance_done, cleared_by, status, created_at
		FROM execprot_movement_plans WHERE status IN ('DRAFT', 'APPROVED', 'ACTIVE') AND date >= NOW() ORDER BY date ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query upcoming movements: %w", err)
	}
	defer rows.Close()

	var plans []domain.MovementPlan
	for rows.Next() {
		var m domain.MovementPlan
		if err := rows.Scan(
			&m.ID, &m.ProtecteeID, &m.EventName, &m.Date, &m.DepartureLocation, &m.ArrivalLocation,
			&m.TransportMode, &m.RoutePlan, &m.AdvanceDone, &m.ClearedBy, &m.Status, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		plans = append(plans, m)
	}
	return plans, rows.Err()
}

func (r *postgresRepo) CreateThreatAssessment(ctx context.Context, t *domain.ThreatAssessment) error {
	query := `INSERT INTO execprot_threat_assessments (id, protectee_id, threat_type, threat_level, threat_detail, source_info, assessed_by, mitigation, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.ProtecteeID, t.ThreatType, t.ThreatLevel, t.ThreatDetail, t.SourceInfo,
		t.AssessedBy, t.Mitigation, t.Status, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert threat assessment: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetActiveThreatsByProtectee(ctx context.Context, protecteeID uuid.UUID) ([]domain.ThreatAssessment, error) {
	query := `SELECT id, protectee_id, threat_type, threat_level, threat_detail, source_info, assessed_by, mitigation, status, created_at
		FROM execprot_threat_assessments WHERE protectee_id = $1 AND status IN ('PENDING', 'ACTIVE') ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, protecteeID)
	if err != nil {
		return nil, fmt.Errorf("query active threats: %w", err)
	}
	defer rows.Close()

	var threats []domain.ThreatAssessment
	for rows.Next() {
		var t domain.ThreatAssessment
		if err := rows.Scan(
			&t.ID, &t.ProtecteeID, &t.ThreatType, &t.ThreatLevel, &t.ThreatDetail, &t.SourceInfo,
			&t.AssessedBy, &t.Mitigation, &t.Status, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		threats = append(threats, t)
	}
	return threats, rows.Err()
}

func (r *postgresRepo) GetDashboard(ctx context.Context) (*domain.DashboardProtection, error) {
	query := `SELECT
		COALESCE((SELECT COUNT(*) FROM execprot_protectees), 0),
		COALESCE((SELECT COUNT(*) FROM execprot_protectees WHERE risk_assessment IN ('HIGH','CRITICAL')), 0),
		COALESCE((SELECT COUNT(*) FROM execprot_movement_plans WHERE status IN ('DRAFT','APPROVED','ACTIVE') AND date >= NOW()), 0),
		COALESCE((SELECT COUNT(*) FROM execprot_threat_assessments WHERE status IN ('PENDING','ACTIVE')), 0),
		COALESCE((SELECT COUNT(*) FROM execprot_protectees WHERE risk_assessment = 'CRITICAL'), 0),
		COALESCE((SELECT COUNT(*) FROM execprot_protectees WHERE risk_assessment = 'HIGH'), 0),
		COALESCE((SELECT COUNT(*) FROM execprot_protectees WHERE risk_assessment = 'MEDIUM'), 0),
		COALESCE((SELECT AVG(active_threats::FLOAT) FROM execprot_protectees), 0.0)`

	d := &domain.DashboardProtection{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&d.TotalProtectees, &d.ActiveProtectees, &d.UpcomingMovements, &d.ActiveThreats,
		&d.CriticalRiskCount, &d.HighRiskCount, &d.MediumRiskCount, &d.AvgActiveThreats,
	)
	if err != nil {
		return nil, fmt.Errorf("query dashboard: %w", err)
	}
	return d, nil
}
