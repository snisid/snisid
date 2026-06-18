package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/sipci-svc/internal/domain"
)

type assetRepo struct {
	pool *pgxpool.Pool
}

func NewAssetRepo(pool *pgxpool.Pool) *assetRepo {
	return &assetRepo{pool: pool}
}

func (r *assetRepo) Create(asset *domain.Asset) (*domain.Asset, error) {
	ctx := context.Background()
	asset.ID = uuid.New()
	asset.NationalSipciID = "SIPCI-HT-" + asset.ID.String()[:6]
	asset.CreatedAt = time.Now()
	asset.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sipci_assets
		 (asset_id, national_sipci_id, asset_name, asset_category, owner_entity, operating_org,
		  dept_code, commune, lat, lng, criticality_score, population_served, single_point_failure,
		  current_threat_level, site_manager_phone, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
		asset.ID, asset.NationalSipciID, asset.AssetName, asset.AssetCategory,
		asset.OwnerEntity, asset.OperatingOrg, asset.DeptCode, asset.Commune,
		asset.Lat, asset.Lng, asset.CriticalityScore, asset.PopulationServed,
		asset.SinglePointFailure, asset.CurrentThreatLevel, asset.SiteManagerPhone,
		asset.CreatedBy, asset.CreatedAt, asset.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (r *assetRepo) FindByID(id uuid.UUID) (*domain.Asset, error) {
	ctx := context.Background()
	a := &domain.Asset{}
	err := r.pool.QueryRow(ctx,
		`SELECT asset_id, national_sipci_id, asset_name, asset_category, dept_code, commune,
		        lat, lng, criticality_score, population_served, single_point_failure,
		        current_threat_level, is_in_gang_zone, under_extortion, incident_count_12m, created_at
		 FROM sipci_assets WHERE asset_id = $1`, id).Scan(
		&a.ID, &a.NationalSipciID, &a.AssetName, &a.AssetCategory, &a.DeptCode, &a.Commune,
		&a.Lat, &a.Lng, &a.CriticalityScore, &a.PopulationServed, &a.SinglePointFailure,
		&a.CurrentThreatLevel, &a.IsInGangZone, &a.UnderExtortion, &a.IncidentCount12m, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *assetRepo) FindAll() ([]domain.Asset, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT asset_id, national_sipci_id, asset_name, asset_category, dept_code, lat, lng,
		        criticality_score, current_threat_level, created_at
		 FROM sipci_assets ORDER BY criticality_score DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []domain.Asset
	for rows.Next() {
		var a domain.Asset
		if err := rows.Scan(&a.ID, &a.NationalSipciID, &a.AssetName, &a.AssetCategory,
			&a.DeptCode, &a.Lat, &a.Lng, &a.CriticalityScore, &a.CurrentThreatLevel, &a.CreatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (r *assetRepo) FindCritical() ([]domain.Asset, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT asset_id, national_sipci_id, asset_name, asset_category, dept_code, lat, lng,
		        criticality_score, current_threat_level, created_at
		 FROM sipci_assets WHERE criticality_score >= 8 ORDER BY criticality_score DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []domain.Asset
	for rows.Next() {
		var a domain.Asset
		if err := rows.Scan(&a.ID, &a.NationalSipciID, &a.AssetName, &a.AssetCategory,
			&a.DeptCode, &a.Lat, &a.Lng, &a.CriticalityScore, &a.CurrentThreatLevel, &a.CreatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (r *assetRepo) FindUnderThreat() ([]domain.Asset, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT asset_id, national_sipci_id, asset_name, asset_category, dept_code, lat, lng,
		        criticality_score, current_threat_level, created_at
		 FROM sipci_assets WHERE current_threat_level IN ('HIGH','SEVERE','CRITICAL')
		 ORDER BY criticality_score DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []domain.Asset
	for rows.Next() {
		var a domain.Asset
		if err := rows.Scan(&a.ID, &a.NationalSipciID, &a.AssetName, &a.AssetCategory,
			&a.DeptCode, &a.Lat, &a.Lng, &a.CriticalityScore, &a.CurrentThreatLevel, &a.CreatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (r *assetRepo) CreateIncident(incident *domain.AssetIncident) (*domain.AssetIncident, error) {
	ctx := context.Background()
	incident.IncidentID = uuid.New()
	incident.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sipci_incidents
		 (incident_id, asset_id, incident_type, incident_date, perpetrator_type, gang_id,
		  description, impact_severity, economic_loss_usd, case_reference, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		incident.IncidentID, incident.AssetID, incident.IncidentType, incident.IncidentDate,
		incident.PerpetratorType, incident.GangID, incident.Description,
		incident.ImpactSeverity, incident.EconomicLossUSD, incident.CaseReference,
		incident.CreatedBy, incident.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return incident, nil
}

func (r *assetRepo) FindRecentIncidents() ([]domain.AssetIncident, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT incident_id, asset_id, incident_type, incident_date, description, impact_severity, created_at
		 FROM sipci_incidents
		 WHERE incident_date >= NOW() - INTERVAL '30 days'
		 ORDER BY incident_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []domain.AssetIncident
	for rows.Next() {
		var inc domain.AssetIncident
		if err := rows.Scan(&inc.IncidentID, &inc.AssetID, &inc.IncidentType,
			&inc.IncidentDate, &inc.Description, &inc.ImpactSeverity, &inc.CreatedAt); err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}
	return incidents, nil
}

func (r *assetRepo) UpdateThreatLevel(id uuid.UUID, level domain.ThreatLevel) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE sipci_assets SET current_threat_level = $1, updated_at = NOW() WHERE asset_id = $2`,
		level, id)
	return err
}

func (r *assetRepo) CountRecentIncidents(assetID uuid.UUID, months int) (int, error) {
	ctx := context.Background()
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM sipci_incidents
		 WHERE asset_id = $1 AND incident_date >= NOW() - ($2 || ' months')::interval`,
		assetID, months).Scan(&count)
	return count, err
}
