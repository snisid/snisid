package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/snisid/critical-infra-protection-ht/internal/domain"
)

type Repository interface {
	CreateAsset(ctx context.Context, a *domain.CriticalAsset) error
	GetAssetsBySector(ctx context.Context, sector string) ([]domain.CriticalAsset, error)
	CreateIncident(ctx context.Context, inc *domain.InfrastructureIncident) error
	GetActiveIncidents(ctx context.Context) ([]domain.InfrastructureIncident, error)
	GetIncidentsByAsset(ctx context.Context, assetID uuid.UUID) ([]domain.InfrastructureIncident, error)
	CreateAssessment(ctx context.Context, a *domain.SectorRiskAssessment) error
	GetNationalDashboard(ctx context.Context) ([]domain.SectorRiskAssessment, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateAsset(ctx context.Context, a *domain.CriticalAsset) error {
	query := `INSERT INTO infraprot_assets
		(id, asset_name, sector, owner_entity, location_lat, location_lng, region, dept_code,
		 criticality, cyber_maturity_score, physical_security_score, last_cisa_assessment_at,
		 contact_name, contact_phone, has_backup_generator, has_cyber_insurance, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.AssetName, a.Sector, a.OwnerEntity, a.LocationLat, a.LocationLng,
		a.Region, a.DeptCode, a.Criticality, a.CyberMaturityScore, a.PhysicalSecurityScore,
		a.LastCISAAssessmentAt, a.ContactName, a.ContactPhone,
		a.HasBackupGenerator, a.HasCyberInsurance, a.CreatedAt,
	)
	return err
}

func (r *postgresRepo) GetAssetsBySector(ctx context.Context, sector string) ([]domain.CriticalAsset, error) {
	query := `SELECT id, asset_name, sector, owner_entity, location_lat, location_lng, region, dept_code,
		criticality, cyber_maturity_score, physical_security_score, last_cisa_assessment_at,
		contact_name, contact_phone, has_backup_generator, has_cyber_insurance, created_at
		FROM infraprot_assets WHERE sector = $1 ORDER BY asset_name`
	rows, err := r.db.QueryContext(ctx, query, sector)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []domain.CriticalAsset
	for rows.Next() {
		var a domain.CriticalAsset
		if err := rows.Scan(
			&a.ID, &a.AssetName, &a.Sector, &a.OwnerEntity,
			&a.LocationLat, &a.LocationLng, &a.Region, &a.DeptCode,
			&a.Criticality, &a.CyberMaturityScore, &a.PhysicalSecurityScore,
			&a.LastCISAAssessmentAt, &a.ContactName, &a.ContactPhone,
			&a.HasBackupGenerator, &a.HasCyberInsurance, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

func (r *postgresRepo) CreateIncident(ctx context.Context, inc *domain.InfrastructureIncident) error {
	query := `INSERT INTO infraprot_incidents
		(id, asset_id, incident_type, severity, description, impact_assessment, downtime_hours,
		 estimated_loss_usd, responded_by, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.db.ExecContext(ctx, query,
		inc.ID, inc.AssetID, inc.IncidentType, inc.Severity, inc.Description,
		inc.ImpactAssessment, inc.DowntimeHours, inc.EstimatedLossUSD,
		inc.RespondedBy, inc.Status, inc.CreatedAt,
	)
	return err
}

func (r *postgresRepo) GetActiveIncidents(ctx context.Context) ([]domain.InfrastructureIncident, error) {
	query := `SELECT id, asset_id, incident_type, severity, description, impact_assessment, downtime_hours,
		estimated_loss_usd, responded_by, status, created_at
		FROM infraprot_incidents WHERE status IN ('REPORTED','RESPONDING','CONTAINED') ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []domain.InfrastructureIncident
	for rows.Next() {
		var inc domain.InfrastructureIncident
		if err := rows.Scan(
			&inc.ID, &inc.AssetID, &inc.IncidentType, &inc.Severity, &inc.Description,
			&inc.ImpactAssessment, &inc.DowntimeHours, &inc.EstimatedLossUSD,
			&inc.RespondedBy, &inc.Status, &inc.CreatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}
	return incidents, rows.Err()
}

func (r *postgresRepo) GetIncidentsByAsset(ctx context.Context, assetID uuid.UUID) ([]domain.InfrastructureIncident, error) {
	query := `SELECT id, asset_id, incident_type, severity, description, impact_assessment, downtime_hours,
		estimated_loss_usd, responded_by, status, created_at
		FROM infraprot_incidents WHERE asset_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []domain.InfrastructureIncident
	for rows.Next() {
		var inc domain.InfrastructureIncident
		if err := rows.Scan(
			&inc.ID, &inc.AssetID, &inc.IncidentType, &inc.Severity, &inc.Description,
			&inc.ImpactAssessment, &inc.DowntimeHours, &inc.EstimatedLossUSD,
			&inc.RespondedBy, &inc.Status, &inc.CreatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}
	return incidents, rows.Err()
}

func (r *postgresRepo) CreateAssessment(ctx context.Context, a *domain.SectorRiskAssessment) error {
	query := `INSERT INTO infraprot_sector_assessments
		(id, sector, assessment_date, overall_risk_score, top_threats, vulnerabilities, recommendations,
		 assessor_agency, next_assessment_due, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.Sector, a.AssessmentDate, a.OverallRiskScore,
		pq.StringArray(a.TopThreats), pq.StringArray(a.Vulnerabilities), pq.StringArray(a.Recommendations),
		a.AssessorAgency, a.NextAssessmentDue, a.CreatedAt,
	)
	return err
}

func (r *postgresRepo) GetNationalDashboard(ctx context.Context) ([]domain.SectorRiskAssessment, error) {
	query := `SELECT id, sector, assessment_date, overall_risk_score, top_threats, vulnerabilities, recommendations,
		assessor_agency, next_assessment_due, created_at
		FROM infraprot_sector_assessments ORDER BY assessment_date DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []domain.SectorRiskAssessment
	for rows.Next() {
		var a domain.SectorRiskAssessment
		var threats, vulns, recs pq.StringArray
		if err := rows.Scan(
			&a.ID, &a.Sector, &a.AssessmentDate, &a.OverallRiskScore,
			&threats, &vulns, &recs,
			&a.AssessorAgency, &a.NextAssessmentDue, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		a.TopThreats = []string(threats)
		a.Vulnerabilities = []string(vulns)
		a.Recommendations = []string(recs)
		assessments = append(assessments, a)
	}
	return assessments, rows.Err()
}
