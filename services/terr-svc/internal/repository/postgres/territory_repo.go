package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/terr-svc/internal/domain"
)

type territoryRepo struct {
	pool *pgxpool.Pool
}

func NewTerritoryRepo(pool *pgxpool.Pool) *territoryRepo {
	return &territoryRepo{pool: pool}
}

func (r *territoryRepo) FindAllZones() ([]domain.TerritoryZone, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT zone_id, gang_id, zone_name, dept_code, commune, section_communale,
		        ST_AsGeoJSON(geom) as geom, area_km2, centroid_lat, centroid_lng,
		        control_level, estimated_population, strategic_importance,
		        controls_national_road, road_numbers, controls_port, controls_airport,
		        controls_market, valid_from, valid_to, is_current, intelligence_source,
		        confidence_level, analyst_notes, created_by, created_at
		 FROM terr_zones
		 WHERE is_current = TRUE
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanZones(rows)
}

func (r *territoryRepo) FindZonesContainingPoint(lat, lng float64) ([]domain.TerritoryZone, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT zone_id, gang_id, zone_name, dept_code, commune, section_communale,
		        ST_AsGeoJSON(geom) as geom, area_km2, centroid_lat, centroid_lng,
		        control_level, estimated_population, strategic_importance,
		        controls_national_road, road_numbers, controls_port, controls_airport,
		        controls_market, valid_from, valid_to, is_current, intelligence_source,
		        confidence_level, analyst_notes, created_by, created_at
		 FROM terr_zones
		 WHERE is_current = TRUE AND ST_Contains(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))`,
		lng, lat)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanZones(rows)
}

func (r *territoryRepo) FindNearbyCheckpoints(lat, lng float64, radiusMeters float64) ([]domain.Checkpoint, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT checkpoint_id, gang_id, ST_AsGeoJSON(location) as location, location_desc,
		        dept_code, road_number, is_armed, extortion_type, reported_at, is_active, created_at
		 FROM terr_checkpoints
		 WHERE is_active = TRUE
		   AND ST_DWithin(location::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3)`,
		lng, lat, radiusMeters)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCheckpoints(rows)
}

func (r *territoryRepo) FindZonesByDept(deptCode string) ([]domain.TerritoryZone, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT zone_id, gang_id, zone_name, dept_code, commune, section_communale,
		        ST_AsGeoJSON(geom) as geom, area_km2, centroid_lat, centroid_lng,
		        control_level, estimated_population, strategic_importance,
		        controls_national_road, road_numbers, controls_port, controls_airport,
		        controls_market, valid_from, valid_to, is_current, intelligence_source,
		        confidence_level, analyst_notes, created_by, created_at
		 FROM terr_zones
		 WHERE is_current = TRUE AND dept_code = $1`,
		deptCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanZones(rows)
}

func (r *territoryRepo) FindZonesByGang(gangID uuid.UUID) ([]domain.TerritoryZone, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT zone_id, gang_id, zone_name, dept_code, commune, section_communale,
		        ST_AsGeoJSON(geom) as geom, area_km2, centroid_lat, centroid_lng,
		        control_level, estimated_population, strategic_importance,
		        controls_national_road, road_numbers, controls_port, controls_airport,
		        controls_market, valid_from, valid_to, is_current, intelligence_source,
		        confidence_level, analyst_notes, created_by, created_at
		 FROM terr_zones
		 WHERE is_current = TRUE AND gang_id = $1`,
		gangID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanZones(rows)
}

func (r *territoryRepo) CreateZone(zone *domain.TerritoryZone) (*domain.TerritoryZone, error) {
	ctx := context.Background()
	zone.ZoneID = uuid.New()
	zone.ValidFrom = time.Now()
	zone.IsCurrent = true
	zone.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO terr_zones
		 (zone_id, gang_id, zone_name, dept_code, commune, section_communale, geom,
		  control_level, estimated_population, strategic_importance, controls_national_road,
		  road_numbers, controls_port, controls_airport, controls_market,
		  valid_from, is_current, intelligence_source, confidence_level, analyst_notes, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6, ST_SetSRID(ST_GeomFromGeoJSON($7), 4326),$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)
		 RETURNING zone_id, created_at`,
		zone.ZoneID, zone.GangID, zone.ZoneName, zone.DeptCode, zone.Commune, zone.SectionCommunale,
		zone.Geom, zone.ControlLevel, zone.EstimatedPopulation, zone.StrategicImportance,
		zone.ControlsNationalRoad, zone.RoadNumbers, zone.ControlsPort, zone.ControlsAirport,
		zone.ControlsMarket, zone.ValidFrom, zone.IsCurrent, zone.IntelligenceSource,
		zone.ConfidenceLevel, zone.AnalystNotes, zone.CreatedBy, zone.CreatedAt,
	).Scan(&zone.ZoneID, &zone.CreatedAt)
	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (r *territoryRepo) UpdateZone(zone *domain.TerritoryZone) (*domain.TerritoryZone, error) {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE terr_zones SET gang_id=$1, zone_name=$2, dept_code=$3, commune=$4, section_communale=$5,
		        geom=ST_SetSRID(ST_GeomFromGeoJSON($6), 4326), control_level=$7, estimated_population=$8,
		        strategic_importance=$9, controls_national_road=$10, road_numbers=$11, controls_port=$12,
		        controls_airport=$13, controls_market=$14, intelligence_source=$15, confidence_level=$16,
		        analyst_notes=$17, valid_to=$18, is_current=$19
		 WHERE zone_id = $20`,
		zone.GangID, zone.ZoneName, zone.DeptCode, zone.Commune, zone.SectionCommunale,
		zone.Geom, zone.ControlLevel, zone.EstimatedPopulation, zone.StrategicImportance,
		zone.ControlsNationalRoad, zone.RoadNumbers, zone.ControlsPort, zone.ControlsAirport,
		zone.ControlsMarket, zone.IntelligenceSource, zone.ConfidenceLevel, zone.AnalystNotes,
		zone.ValidTo, zone.IsCurrent, zone.ZoneID)
	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (r *territoryRepo) GetZoneHistory(zoneID uuid.UUID) ([]domain.ZoneHistory, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT history_id, zone_id, change_type, previous_control, new_control, change_date, trigger_event, created_at
		 FROM terr_zone_history
		 WHERE zone_id = $1
		 ORDER BY change_date DESC`,
		zoneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.ZoneHistory
	for rows.Next() {
		var h domain.ZoneHistory
		if err := rows.Scan(&h.HistoryID, &h.ZoneID, &h.ChangeType, &h.PreviousControl,
			&h.NewControl, &h.ChangeDate, &h.TriggerEvent, &h.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func (r *territoryRepo) CreateCheckpoint(cp *domain.Checkpoint) (*domain.Checkpoint, error) {
	ctx := context.Background()
	cp.CheckpointID = uuid.New()
	cp.IsActive = true
	cp.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO terr_checkpoints
		 (checkpoint_id, gang_id, location, location_desc, dept_code, road_number, is_armed,
		  extortion_type, reported_at, is_active, created_at)
		 VALUES ($1,$2, ST_SetSRID(ST_GeomFromGeoJSON($3), 4326),$4,$5,$6,$7,$8,$9,$10,$11)
		 RETURNING checkpoint_id, created_at`,
		cp.CheckpointID, cp.GangID, cp.Location, cp.LocationDesc, cp.DeptCode, cp.RoadNumber,
		cp.IsArmed, cp.ExtortionType, cp.ReportedAt, cp.IsActive, cp.CreatedAt,
	).Scan(&cp.CheckpointID, &cp.CreatedAt)
	if err != nil {
		return nil, err
	}

	return cp, nil
}

func scanZones(rows pgx.Rows) ([]domain.TerritoryZone, error) {
	var zones []domain.TerritoryZone
	for rows.Next() {
		var z domain.TerritoryZone
		if err := rows.Scan(&z.ZoneID, &z.GangID, &z.ZoneName, &z.DeptCode, &z.Commune,
			&z.SectionCommunale, &z.Geom, &z.AreaKm2, &z.CentroidLat, &z.CentroidLng,
			&z.ControlLevel, &z.EstimatedPopulation, &z.StrategicImportance,
			&z.ControlsNationalRoad, &z.RoadNumbers, &z.ControlsPort, &z.ControlsAirport,
			&z.ControlsMarket, &z.ValidFrom, &z.ValidTo, &z.IsCurrent, &z.IntelligenceSource,
			&z.ConfidenceLevel, &z.AnalystNotes, &z.CreatedBy, &z.CreatedAt); err != nil {
			return nil, err
		}
		zones = append(zones, z)
	}
	return zones, nil
}

func scanCheckpoints(rows pgx.Rows) ([]domain.Checkpoint, error) {
	var cps []domain.Checkpoint
	for rows.Next() {
		var cp domain.Checkpoint
		if err := rows.Scan(&cp.CheckpointID, &cp.GangID, &cp.Location, &cp.LocationDesc,
			&cp.DeptCode, &cp.RoadNumber, &cp.IsArmed, &cp.ExtortionType, &cp.ReportedAt,
			&cp.IsActive, &cp.CreatedAt); err != nil {
			return nil, err
		}
		cps = append(cps, cp)
	}
	return cps, nil
}
