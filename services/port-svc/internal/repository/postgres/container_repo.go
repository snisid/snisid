package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/port-svc/internal/domain"
)

type containerRepo struct {
	pool *pgxpool.Pool
}

func NewContainerRepo(pool *pgxpool.Pool) *containerRepo {
	return &containerRepo{pool: pool}
}

func (r *containerRepo) CreateArrival(arrival *domain.VesselArrival) (*domain.VesselArrival, error) {
	ctx := context.Background()
	arrival.ID = uuid.New()
	arrival.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO port_vessels_arrivals
		 (arrival_id, port_code, vessel_imo, vessel_name, flag_country, shipping_company,
		  arrival_date, origin_port, origin_country, container_count, manifest_ref, mar_vessel_id,
		  risk_score, risk_level, cbp_targeting_ref, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		 RETURNING arrival_id, created_at`,
		arrival.ID, arrival.PortCode, arrival.VesselIMO, arrival.VesselName, arrival.FlagCountry,
		arrival.ShippingCompany, arrival.ArrivalDate, arrival.OriginPort, arrival.OriginCountry,
		arrival.ContainerCount, arrival.ManifestRef, arrival.MARVesselID,
		arrival.RiskScore, arrival.RiskLevel, arrival.CBPTargetingRef, arrival.CreatedAt,
	).Scan(&arrival.ID, &arrival.CreatedAt)
	if err != nil {
		return nil, err
	}

	return arrival, nil
}

func (r *containerRepo) FindArrivalByID(id uuid.UUID) (*domain.VesselArrival, error) {
	ctx := context.Background()
	arrival := &domain.VesselArrival{}
	err := r.pool.QueryRow(ctx,
		`SELECT arrival_id, port_code, vessel_imo, vessel_name, flag_country, shipping_company,
		        arrival_date, origin_port, origin_country, container_count, manifest_ref, mar_vessel_id,
		        risk_score, risk_level, cbp_targeting_ref, created_at
		 FROM port_vessels_arrivals WHERE arrival_id = $1`, id).Scan(
		&arrival.ID, &arrival.PortCode, &arrival.VesselIMO, &arrival.VesselName, &arrival.FlagCountry,
		&arrival.ShippingCompany, &arrival.ArrivalDate, &arrival.OriginPort, &arrival.OriginCountry,
		&arrival.ContainerCount, &arrival.ManifestRef, &arrival.MARVesselID,
		&arrival.RiskScore, &arrival.RiskLevel, &arrival.CBPTargetingRef, &arrival.CreatedAt)
	if err != nil {
		return nil, err
	}
	return arrival, nil
}

func (r *containerRepo) GetHighRiskContainers() ([]domain.Container, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT container_id, arrival_id, container_number, container_type, declared_content,
		        declared_weight_kg, declared_value_usd, shipper_name, shipper_country,
		        consignee_name, consignee_snisid_id, status, risk_score, risk_level,
		        risk_flags, selected_for_scan, scan_date, scan_result, seized,
		        seizure_description, case_reference, cbp_targeting_match, created_at, updated_at
		 FROM port_containers
		 WHERE risk_level IN ('HIGH','CRITICAL')
		 ORDER BY risk_score DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containers []domain.Container
	for rows.Next() {
		var c domain.Container
		if err := rows.Scan(
			&c.ID, &c.ArrivalID, &c.ContainerNumber, &c.ContainerType, &c.DeclaredContent,
			&c.DeclaredWeightKg, &c.DeclaredValueUSD, &c.ShipperName, &c.ShipperCountry,
			&c.ConsigneeName, &c.ConsigneeSNISIDID, &c.Status, &c.RiskScore, &c.RiskLevel,
			&c.RiskFlags, &c.SelectedForScan, &c.ScanDate, &c.ScanResult, &c.Seized,
			&c.SeizureDescription, &c.CaseReference, &c.CBPTargetingMatch, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		containers = append(containers, c)
	}
	return containers, nil
}

func (r *containerRepo) CreateContainer(container *domain.Container) (*domain.Container, error) {
	ctx := context.Background()
	container.ID = uuid.New()
	container.CreatedAt = time.Now()
	container.UpdatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO port_containers
		 (container_id, arrival_id, container_number, container_type, declared_content,
		  declared_weight_kg, declared_value_usd, shipper_name, shipper_country,
		  consignee_name, consignee_snisid_id, status, risk_score, risk_level, risk_flags,
		  selected_for_scan, scan_date, scan_result, seized, seizure_description,
		  case_reference, cbp_targeting_match, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24)
		 RETURNING container_id, created_at, updated_at`,
		container.ID, container.ArrivalID, container.ContainerNumber, container.ContainerType,
		container.DeclaredContent, container.DeclaredWeightKg, container.DeclaredValueUSD,
		container.ShipperName, container.ShipperCountry, container.ConsigneeName,
		container.ConsigneeSNISIDID, container.Status, container.RiskScore, container.RiskLevel,
		container.RiskFlags, container.SelectedForScan, container.ScanDate, container.ScanResult,
		container.Seized, container.SeizureDescription, container.CaseReference,
		container.CBPTargetingMatch, container.CreatedAt, container.UpdatedAt,
	).Scan(&container.ID, &container.CreatedAt, &container.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (r *containerRepo) ScanContainer(id uuid.UUID, scanResult string) (*domain.Container, error) {
	ctx := context.Background()
	now := time.Now()
	container := &domain.Container{}
	err := r.pool.QueryRow(ctx,
		`UPDATE port_containers
		 SET scan_date = $1, scan_result = $2, selected_for_scan = true, updated_at = $3
		 WHERE container_id = $4
		 RETURNING container_id, arrival_id, container_number, container_type, declared_content,
		        declared_weight_kg, declared_value_usd, shipper_name, shipper_country,
		        consignee_name, consignee_snisid_id, status, risk_score, risk_level,
		        risk_flags, selected_for_scan, scan_date, scan_result, seized,
		        seizure_description, case_reference, cbp_targeting_match, created_at, updated_at`,
		now, scanResult, now, id,
	).Scan(
		&container.ID, &container.ArrivalID, &container.ContainerNumber, &container.ContainerType,
		&container.DeclaredContent, &container.DeclaredWeightKg, &container.DeclaredValueUSD,
		&container.ShipperName, &container.ShipperCountry, &container.ConsigneeName,
		&container.ConsigneeSNISIDID, &container.Status, &container.RiskScore, &container.RiskLevel,
		&container.RiskFlags, &container.SelectedForScan, &container.ScanDate, &container.ScanResult,
		&container.Seized, &container.SeizureDescription, &container.CaseReference,
		&container.CBPTargetingMatch, &container.CreatedAt, &container.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func (r *containerRepo) SeizeContainer(id uuid.UUID, description string, caseRef string) (*domain.Container, error) {
	ctx := context.Background()
	now := time.Now()
	container := &domain.Container{}
	err := r.pool.QueryRow(ctx,
		`UPDATE port_containers
		 SET seized = true, status = 'SEIZED', seizure_description = $1, case_reference = $2, updated_at = $3
		 WHERE container_id = $4
		 RETURNING container_id, arrival_id, container_number, container_type, declared_content,
		        declared_weight_kg, declared_value_usd, shipper_name, shipper_country,
		        consignee_name, consignee_snisid_id, status, risk_score, risk_level,
		        risk_flags, selected_for_scan, scan_date, scan_result, seized,
		        seizure_description, case_reference, cbp_targeting_match, created_at, updated_at`,
		description, caseRef, now, id,
	).Scan(
		&container.ID, &container.ArrivalID, &container.ContainerNumber, &container.ContainerType,
		&container.DeclaredContent, &container.DeclaredWeightKg, &container.DeclaredValueUSD,
		&container.ShipperName, &container.ShipperCountry, &container.ConsigneeName,
		&container.ConsigneeSNISIDID, &container.Status, &container.RiskScore, &container.RiskLevel,
		&container.RiskFlags, &container.SelectedForScan, &container.ScanDate, &container.ScanResult,
		&container.Seized, &container.SeizureDescription, &container.CaseReference,
		&container.CBPTargetingMatch, &container.CreatedAt, &container.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func (r *containerRepo) GetSeizureStats() (*domain.SeizureStats, error) {
	ctx := context.Background()
	stats := &domain.SeizureStats{}
	err := r.pool.QueryRow(ctx,
		`SELECT
			COALESCE(SUM(CASE WHEN seized = true THEN 1 ELSE 0 END), 0) AS total_seized,
			COALESCE(SUM(CASE WHEN scan_date IS NOT NULL THEN 1 ELSE 0 END), 0) AS total_scanned,
			COALESCE(SUM(CASE WHEN risk_level = 'HIGH' THEN 1 ELSE 0 END), 0) AS high_risk_count,
			COALESCE(SUM(CASE WHEN risk_level = 'CRITICAL' THEN 1 ELSE 0 END), 0) AS critical_risk_count
		 FROM port_containers`).Scan(&stats.TotalSeized, &stats.TotalScanned, &stats.HighRiskCount, &stats.CriticalRiskCount)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
