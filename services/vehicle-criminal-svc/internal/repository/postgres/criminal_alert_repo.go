package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/snisid/vehicle-criminal-svc/internal/repository"
)

type CriminalAlertRepo struct {
	db *sqlx.DB
}

func NewCriminalAlertRepo(db *sqlx.DB) *CriminalAlertRepo {
	return &CriminalAlertRepo{db: db}
}

func (r *CriminalAlertRepo) Create(ctx context.Context, alert *domain.CriminalAlert) error {
	query := `
		INSERT INTO sivc_criminal_alerts (
			alert_id, plate_number, plate_category, vin, chassis_number,
			vehicle_type, make, model, year, color_primary, color_secondary,
			distinguishing_marks, foves_vehicle_id, crime_category, crime_subcategory,
			alert_level, status, armed_and_dangerous, do_not_stop_alone,
			officer_safety_notes, reporting_unit, reporting_officer_id,
			incident_reference, incident_date, expiry_date,
			associated_person_ids, associated_case_ids, associated_alert_ids,
			interpol_smv_id, interpol_reported, interpol_reported_at,
			last_seen_lat, last_seen_lng, last_seen_location,
			last_seen_dept_code, last_seen_commune, last_seen_at,
			photo_refs, document_refs, created_by, created_at, updated_at, version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		alert.AlertID, alert.PlateNumber, alert.PlateCategory, alert.VIN, alert.ChassisNumber,
		alert.VehicleType, alert.Make, alert.Model, alert.Year, alert.ColorPrimary, alert.ColorSecondary,
		alert.DistinguishingMarks, alert.FovesVehicleID, alert.CrimeCategory, alert.CrimeSubcategory,
		alert.AlertLevel, alert.Status, alert.ArmedAndDangerous, alert.DoNotStopAlone,
		alert.OfficerSafetyNotes, alert.ReportingUnit, alert.ReportingOfficerID,
		alert.IncidentReference, alert.IncidentDate, alert.ExpiryDate,
		alert.AssociatedPersonIDs, alert.AssociatedCaseIDs, alert.AssociatedAlertIDs,
		alert.InterpolSMVID, alert.InterpolReported, alert.InterpolReportedAt,
		alert.LastSeenLat, alert.LastSeenLng, alert.LastSeenLocation,
		alert.LastSeenDeptCode, alert.LastSeenCommune, alert.LastSeenAt,
		alert.PhotoRefs, alert.DocumentRefs, alert.CreatedBy, alert.CreatedAt, alert.UpdatedAt, alert.Version,
	)
	return err
}

func (r *CriminalAlertRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.CriminalAlert, error) {
	var alert domain.CriminalAlert
	query := `SELECT * FROM sivc_criminal_alerts WHERE alert_id = $1`
	if err := r.db.GetContext(ctx, &alert, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &alert, nil
}

func (r *CriminalAlertRepo) FindActiveByPlate(ctx context.Context, plateNumber string) (*domain.CriminalAlert, error) {
	var alert domain.CriminalAlert
	query := `SELECT * FROM sivc_criminal_alerts WHERE plate_number = $1 AND status = 'ACTIVE' LIMIT 1`
	if err := r.db.GetContext(ctx, &alert, query, plateNumber); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &alert, nil
}

func (r *CriminalAlertRepo) FindAll(ctx context.Context, filter repository.AlertFilter) ([]*domain.CriminalAlert, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if filter.DeptCode != "" {
		where = append(where, fmt.Sprintf("last_seen_dept_code = $%d", argIdx))
		args = append(args, filter.DeptCode)
		argIdx++
	}
	if filter.Category != "" {
		where = append(where, fmt.Sprintf("crime_category = $%d", argIdx))
		args = append(args, filter.Category)
		argIdx++
	}
	if filter.Level != "" {
		where = append(where, fmt.Sprintf("alert_level = $%d", argIdx))
		args = append(args, filter.Level)
		argIdx++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	} else {
		where = append(where, "status = 'ACTIVE'")
	}
	if filter.ReportingUnit != "" {
		where = append(where, fmt.Sprintf("reporting_unit = $%d", argIdx))
		args = append(args, filter.ReportingUnit)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM sivc_criminal_alerts WHERE %s", whereClause)
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}
	offset := (filter.Page - 1) * filter.Limit

	dataQuery := fmt.Sprintf(
		"SELECT * FROM sivc_criminal_alerts WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		whereClause, argIdx, argIdx+1,
	)
	args = append(args, filter.Limit, offset)

	var alerts []*domain.CriminalAlert
	if err := r.db.SelectContext(ctx, &alerts, dataQuery, args...); err != nil {
		return nil, 0, err
	}
	return alerts, total, nil
}

func (r *CriminalAlertRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AlertStatus, updatedBy uuid.UUID) error {
	query := `UPDATE sivc_criminal_alerts SET status = $1, updated_by = $2 WHERE alert_id = $3`
	_, err := r.db.ExecContext(ctx, query, status, updatedBy, id)
	return err
}

func (r *CriminalAlertRepo) UpdateLastSeen(ctx context.Context, id uuid.UUID, lat, lng float64, location, deptCode, commune string) error {
	query := `
		UPDATE sivc_criminal_alerts
		SET last_seen_lat = $1, last_seen_lng = $2, last_seen_location = $3,
		    last_seen_dept_code = $4, last_seen_commune = $5, last_seen_at = NOW()
		WHERE alert_id = $6
	`
	_, err := r.db.ExecContext(ctx, query, lat, lng, location, deptCode, commune, id)
	return err
}

func (r *CriminalAlertRepo) Search(ctx context.Context, query string, filters repository.AlertFilter) ([]*domain.CriminalAlert, int, error) {
	where := []string{
		"to_tsvector('french', COALESCE(plate_number, '') || ' ' || COALESCE(make, '') || ' ' || COALESCE(model, '') || ' ' || COALESCE(color_primary, '')) @@ plainto_tsquery('french', $1)",
	}
	args := []interface{}{query}
	argIdx := 2

	if filters.DeptCode != "" {
		where = append(where, fmt.Sprintf("last_seen_dept_code = $%d", argIdx))
		args = append(args, filters.DeptCode)
		argIdx++
	}
	if filters.Category != "" {
		where = append(where, fmt.Sprintf("crime_category = $%d", argIdx))
		args = append(args, filters.Category)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM sivc_criminal_alerts WHERE %s", whereClause)
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 20
	}
	offset := (filters.Page - 1) * filters.Limit

	dataQuery := fmt.Sprintf(
		"SELECT * FROM sivc_criminal_alerts WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		whereClause, argIdx, argIdx+1,
	)
	args = append(args, filters.Limit, offset)

	var alerts []*domain.CriminalAlert
	if err := r.db.SelectContext(ctx, &alerts, dataQuery, args...); err != nil {
		return nil, 0, err
	}
	return alerts, total, nil
}
