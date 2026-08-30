package repository

import (
	"database/sql"

	"github.com/snisid/lapi-ht/internal/domain"
)

type Repository interface {
	SavePlateRead(read *domain.PlateRead) error
	GetRecentReads(limit int) ([]domain.PlateRead, error)
	GetReadsByPlate(plateNumber string) ([]domain.PlateRead, error)
	GetActiveAlerts() ([]domain.AlertDispatch, error)
	GetCameras() ([]domain.Camera, error)
	SaveAlertDispatch(alert *domain.AlertDispatch) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SavePlateRead(read *domain.PlateRead) error {
	query := `INSERT INTO lapi_reads (id, camera_id, plate_number_raw, plate_number_normalized, ocr_confidence, latitude, longitude, speed_estimate_kmh, alert_triggered, captured_at, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.db.Exec(query, read.ID, read.CameraID, read.PlateNumberRaw, read.PlateNumberNormalized, read.OcrConfidence, read.Latitude, read.Longitude, read.SpeedEstimateKmh, read.AlertTriggered, read.CapturedAt, read.CreatedAt)
	return err
}

func (r *PostgresRepository) GetRecentReads(limit int) ([]domain.PlateRead, error) {
	query := `SELECT id, camera_id, plate_number_raw, plate_number_normalized, ocr_confidence, latitude, longitude, speed_estimate_kmh, alert_triggered, captured_at, created_at FROM lapi_reads ORDER BY captured_at DESC LIMIT $1`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reads []domain.PlateRead
	for rows.Next() {
		var pr domain.PlateRead
		var lat, lng, speed sql.NullFloat64
		err := rows.Scan(&pr.ID, &pr.CameraID, &pr.PlateNumberRaw, &pr.PlateNumberNormalized, &pr.OcrConfidence, &lat, &lng, &speed, &pr.AlertTriggered, &pr.CapturedAt, &pr.CreatedAt)
		if err != nil {
			return nil, err
		}
		if lat.Valid {
			pr.Latitude = &lat.Float64
		}
		if lng.Valid {
			pr.Longitude = &lng.Float64
		}
		if speed.Valid {
			pr.SpeedEstimateKmh = &speed.Float64
		}
		reads = append(reads, pr)
	}
	return reads, rows.Err()
}

func (r *PostgresRepository) GetReadsByPlate(plateNumber string) ([]domain.PlateRead, error) {
	query := `SELECT id, camera_id, plate_number_raw, plate_number_normalized, ocr_confidence, latitude, longitude, speed_estimate_kmh, alert_triggered, captured_at, created_at FROM lapi_reads WHERE plate_number_normalized = $1 ORDER BY captured_at DESC`
	rows, err := r.db.Query(query, plateNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reads []domain.PlateRead
	for rows.Next() {
		var pr domain.PlateRead
		var lat, lng, speed sql.NullFloat64
		err := rows.Scan(&pr.ID, &pr.CameraID, &pr.PlateNumberRaw, &pr.PlateNumberNormalized, &pr.OcrConfidence, &lat, &lng, &speed, &pr.AlertTriggered, &pr.CapturedAt, &pr.CreatedAt)
		if err != nil {
			return nil, err
		}
		if lat.Valid {
			pr.Latitude = &lat.Float64
		}
		if lng.Valid {
			pr.Longitude = &lng.Float64
		}
		if speed.Valid {
			pr.SpeedEstimateKmh = &speed.Float64
		}
		reads = append(reads, pr)
	}
	return reads, rows.Err()
}

func (r *PostgresRepository) GetActiveAlerts() ([]domain.AlertDispatch, error) {
	query := `SELECT id, read_id, plate_number, reason, dispatched_at, resolved_at, is_active FROM lapi_alert_dispatches WHERE is_active = true ORDER BY dispatched_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.AlertDispatch
	for rows.Next() {
		var a domain.AlertDispatch
		var resolvedAt sql.NullTime
		err := rows.Scan(&a.ID, &a.ReadID, &a.PlateNumber, &a.Reason, &a.DispatchedAt, &resolvedAt, &a.IsActive)
		if err != nil {
			return nil, err
		}
		if resolvedAt.Valid {
			a.ResolvedAt = &resolvedAt.Time
		}
		alerts = append(alerts, a)
	}
	return alerts, rows.Err()
}

func (r *PostgresRepository) GetCameras() ([]domain.Camera, error) {
	query := `SELECT id, label, type, latitude, longitude, is_active, created_at, updated_at FROM lapi_cameras ORDER BY label`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cameras []domain.Camera
	for rows.Next() {
		var c domain.Camera
		err := rows.Scan(&c.ID, &c.Label, &c.Type, &c.Latitude, &c.Longitude, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cameras = append(cameras, c)
	}
	return cameras, rows.Err()
}

func (r *PostgresRepository) SaveAlertDispatch(alert *domain.AlertDispatch) error {
	query := `INSERT INTO lapi_alert_dispatches (id, read_id, plate_number, reason, dispatched_at, resolved_at, is_active) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(query, alert.ID, alert.ReadID, alert.PlateNumber, alert.Reason, alert.DispatchedAt, alert.ResolvedAt, alert.IsActive)
	return err
}
