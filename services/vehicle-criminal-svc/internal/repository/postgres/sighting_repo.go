package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

type SightingRepo struct {
	db *sqlx.DB
}

func NewSightingRepo(db *sqlx.DB) *SightingRepo {
	return &SightingRepo{db: db}
}

func (r *SightingRepo) Create(ctx context.Context, sighting *domain.VehicleSighting) error {
	query := `
		INSERT INTO sivc_vehicle_sightings (
			sighting_id, plate_number, source_type, lapi_unit_id, reporting_agent_id,
			sighting_timestamp, location_lat, location_lng, location_desc,
			dept_code, commune, checkpoint_name, matched_alert_id, matched_plate_id,
			match_confidence, alert_triggered, alert_level, alert_sent_at,
			alert_recipients, image_ref, video_clip_ref, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		sighting.SightingID, sighting.PlateNumber, sighting.SourceType, sighting.LAPIUnitID,
		sighting.ReportingAgentID, sighting.SightingTimestamp, sighting.LocationLat,
		sighting.LocationLng, sighting.LocationDesc, sighting.DeptCode, sighting.Commune,
		sighting.CheckpointName, sighting.MatchedAlertID, sighting.MatchedPlateID,
		sighting.MatchConfidence, sighting.AlertTriggered, sighting.AlertLevel,
		sighting.AlertSentAt, sighting.AlertRecipients, sighting.ImageRef,
		sighting.VideoClipRef, sighting.CreatedAt,
	)
	return err
}

func (r *SightingRepo) FindByAlertID(ctx context.Context, alertID uuid.UUID) ([]*domain.VehicleSighting, error) {
	var sightings []*domain.VehicleSighting
	query := `SELECT * FROM sivc_vehicle_sightings WHERE matched_alert_id = $1 ORDER BY sighting_timestamp DESC`
	if err := r.db.SelectContext(ctx, &sightings, query, alertID); err != nil {
		return nil, err
	}
	return sightings, nil
}

func (r *SightingRepo) FindByPlate(ctx context.Context, plateNumber string) ([]*domain.VehicleSighting, error) {
	var sightings []*domain.VehicleSighting
	query := `SELECT * FROM sivc_vehicle_sightings WHERE plate_number = $1 ORDER BY sighting_timestamp DESC`
	if err := r.db.SelectContext(ctx, &sightings, query, plateNumber); err != nil {
		return nil, err
	}
	return sightings, nil
}
