package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

type StolenPlateRepo struct {
	db *sqlx.DB
}

func NewStolenPlateRepo(db *sqlx.DB) *StolenPlateRepo {
	return &StolenPlateRepo{db: db}
}

func (r *StolenPlateRepo) Create(ctx context.Context, plate *domain.StolenPlate) error {
	query := `
		INSERT INTO sivc_stolen_plates (
			plate_id, plate_number, plate_category, original_vehicle_id,
			original_make, original_model, original_vin, theft_date,
			theft_location, theft_dept_code, theft_commune, theft_context,
			reporting_unit, reporting_officer_id, blvv_case_number, status,
			is_state_plate_clone, impersonated_agency, notes, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		plate.PlateID, plate.PlateNumber, plate.PlateCategory, plate.OriginalVehicleID,
		plate.OriginalMake, plate.OriginalModel, plate.OriginalVIN, plate.TheftDate,
		plate.TheftLocation, plate.TheftDeptCode, plate.TheftCommune, plate.TheftContext,
		plate.ReportingUnit, plate.ReportingOfficerID, plate.BlvvCaseNumber, plate.Status,
		plate.IsStatePlateClone, plate.ImpersonatedAgency, plate.Notes, plate.CreatedBy, plate.CreatedAt, plate.UpdatedAt,
	)
	return err
}

func (r *StolenPlateRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.StolenPlate, error) {
	var plate domain.StolenPlate
	query := `SELECT * FROM sivc_stolen_plates WHERE plate_id = $1`
	if err := r.db.GetContext(ctx, &plate, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &plate, nil
}

func (r *StolenPlateRepo) FindByPlate(ctx context.Context, plateNumber string) (*domain.StolenPlate, error) {
	var plate domain.StolenPlate
	query := `SELECT * FROM sivc_stolen_plates WHERE plate_number = $1 ORDER BY created_at DESC LIMIT 1`
	if err := r.db.GetContext(ctx, &plate, query, plateNumber); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &plate, nil
}

func (r *StolenPlateRepo) FindStolenByPlate(ctx context.Context, plateNumber string) (*domain.StolenPlate, error) {
	var plate domain.StolenPlate
	query := `SELECT * FROM sivc_stolen_plates WHERE plate_number = $1 AND status = 'STOLEN' LIMIT 1`
	if err := r.db.GetContext(ctx, &plate, query, plateNumber); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &plate, nil
}

func (r *StolenPlateRepo) MarkRecovered(ctx context.Context, id uuid.UUID, location string, deptCode string) error {
	query := `
		UPDATE sivc_stolen_plates
		SET status = 'RECOVERED', recovered_date = NOW(), recovery_location = $1, recovery_dept_code = $2, updated_at = NOW()
		WHERE plate_id = $3
	`
	_, err := r.db.ExecContext(ctx, query, location, deptCode, id)
	return err
}
