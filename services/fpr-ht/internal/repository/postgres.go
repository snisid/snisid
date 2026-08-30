package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/snisid/fpr-ht/internal/domain"
)

type Repository interface {
	SaveWarrant(w *domain.Warrant) error
	FindWarrantsByName(name string) ([]domain.Warrant, error)
	SaveCheckLog(cl *domain.CheckLog) error
	SaveSighting(s *domain.Sighting) error
	UpdateWarrantExecuted(id uuid.UUID, executedAt time.Time) error
	GetArmedDangerousWarrants() ([]domain.Warrant, error)
	GetDashboardStats() (*domain.DashboardStats, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveWarrant(w *domain.Warrant) error {
	query := `INSERT INTO fpr_warrants (id, full_name, aliases, afis_subject_id, warrant_type, charges, issuing_court, danger_level, photo_refs, vehicle_plates_known, interpol_notice_ref, is_executed, issued_at, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := r.db.Exec(query, w.ID, w.FullName, pq.Array(w.Aliases), w.AfisSubjectID, w.WarrantType, pq.Array(w.Charges), w.IssuingCourt, w.DangerLevel, pq.Array(w.PhotoRefs), pq.Array(w.VehiclePlatesKnown), w.InterpolNoticeRef, w.IsExecuted, w.IssuedAt, w.CreatedAt, w.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetWarrantByID(id uuid.UUID) (*domain.Warrant, error) {
	query := `SELECT id, full_name, aliases, afis_subject_id, warrant_type, charges, issuing_court, danger_level, photo_refs, vehicle_plates_known, interpol_notice_ref, is_executed, issued_at, executed_at, created_at, updated_at FROM fpr_warrants WHERE id = $1`
	var w domain.Warrant
	var afisID, interpolRef, dangerLevel sql.NullString
	var executedAt sql.NullTime
	var aliases, charges, photoRefs, plates []string
	err := r.db.QueryRow(query, id).Scan(&w.ID, &w.FullName, pq.Array(&aliases), &afisID, &w.WarrantType, pq.Array(&charges), &w.IssuingCourt, &dangerLevel, pq.Array(&photoRefs), pq.Array(&plates), &interpolRef, &w.IsExecuted, &w.IssuedAt, &executedAt, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	w.Aliases = aliases
	w.Charges = charges
	w.PhotoRefs = photoRefs
	w.VehiclePlatesKnown = plates
	if afisID.Valid {
		w.AfisSubjectID = &afisID.String
	}
	if interpolRef.Valid {
		w.InterpolNoticeRef = &interpolRef.String
	}
	if dangerLevel.Valid {
		d := domain.DangerLevel(dangerLevel.String)
		w.DangerLevel = &d
	}
	if executedAt.Valid {
		w.ExecutedAt = &executedAt.Time
	}
	return &w, nil
}

func (r *PostgresRepository) FindWarrantsByName(name string) ([]domain.Warrant, error) {
	query := `SELECT id, full_name, aliases, afis_subject_id, warrant_type, charges, issuing_court, danger_level, photo_refs, vehicle_plates_known, interpol_notice_ref, is_executed, issued_at, executed_at, created_at, updated_at FROM fpr_warrants WHERE full_name ILIKE $1 OR $1 = ANY(aliases) ORDER BY issued_at DESC`
	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanWarrants(rows)
}

func (r *PostgresRepository) GetArmedDangerousWarrants() ([]domain.Warrant, error) {
	query := `SELECT id, full_name, aliases, afis_subject_id, warrant_type, charges, issuing_court, danger_level, photo_refs, vehicle_plates_known, interpol_notice_ref, is_executed, issued_at, executed_at, created_at, updated_at FROM fpr_warrants WHERE danger_level = 'ARMED_AND_DANGEROUS' AND is_executed = false ORDER BY issued_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanWarrants(rows)
}

func (r *PostgresRepository) UpdateWarrantExecuted(id uuid.UUID, executedAt time.Time) error {
	query := `UPDATE fpr_warrants SET is_executed = true, executed_at = $2, updated_at = $2 WHERE id = $1`
	_, err := r.db.Exec(query, id, executedAt)
	return err
}

func (r *PostgresRepository) SaveSighting(s *domain.Sighting) error {
	query := `INSERT INTO fpr_sighting_log (id, warrant_id, citizen_id, latitude, longitude, description, reported_by, sighted_at, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.db.Exec(query, s.ID, s.WarrantID, s.CitizenID, s.Latitude, s.Longitude, s.Description, s.ReportedBy, s.SightedAt, s.CreatedAt)
	return err
}

func (r *PostgresRepository) SaveCheckLog(cl *domain.CheckLog) error {
	query := `INSERT INTO fpr_check_log (id, citizen_id, warrant_id, result, checked_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.Exec(query, cl.ID, cl.CitizenID, cl.WarrantID, cl.Result, cl.CheckedAt)
	return err
}

func (r *PostgresRepository) GetDashboardStats() (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}
	err := r.db.QueryRow(`SELECT COUNT(*) FROM fpr_warrants`).Scan(&stats.TotalWarrants)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(`SELECT COUNT(*) FROM fpr_warrants WHERE is_executed = false`).Scan(&stats.ActiveWarrants)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(`SELECT COUNT(*) FROM fpr_warrants WHERE is_executed = true`).Scan(&stats.ExecutedWarrants)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(`SELECT COUNT(*) FROM fpr_warrants WHERE danger_level = 'ARMED_AND_DANGEROUS' AND is_executed = false`).Scan(&stats.ArmedDangerous)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(`SELECT COUNT(*) FROM fpr_sighting_log`).Scan(&stats.TotalSightings)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(`SELECT COUNT(*) FROM fpr_check_log`).Scan(&stats.TotalChecks)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *PostgresRepository) scanWarrants(rows *sql.Rows) ([]domain.Warrant, error) {
	var warrants []domain.Warrant
	for rows.Next() {
		var w domain.Warrant
		var afisID, interpolRef, dangerLevel sql.NullString
		var executedAt sql.NullTime
		var aliases, charges, photoRefs, plates []string
		err := rows.Scan(&w.ID, &w.FullName, pq.Array(&aliases), &afisID, &w.WarrantType, pq.Array(&charges), &w.IssuingCourt, &dangerLevel, pq.Array(&photoRefs), pq.Array(&plates), &interpolRef, &w.IsExecuted, &w.IssuedAt, &executedAt, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}
		w.Aliases = aliases
		w.Charges = charges
		w.PhotoRefs = photoRefs
		w.VehiclePlatesKnown = plates
		if afisID.Valid {
			w.AfisSubjectID = &afisID.String
		}
		if interpolRef.Valid {
			w.InterpolNoticeRef = &interpolRef.String
		}
		if dangerLevel.Valid {
			d := domain.DangerLevel(dangerLevel.String)
			w.DangerLevel = &d
		}
		if executedAt.Valid {
			w.ExecutedAt = &executedAt.Time
		}
		warrants = append(warrants, w)
	}
	return warrants, rows.Err()
}
