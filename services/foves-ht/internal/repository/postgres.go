package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/foves-ht/internal/domain"
)

type Repository interface {
	CreateVehicle(ctx context.Context, v *domain.Vehicle) error
	FindByPlate(ctx context.Context, plate string) (*domain.Vehicle, error)
	FindByVIN(ctx context.Context, vin string) (*domain.Vehicle, error)
	FindByOwner(ctx context.Context, citizenID uuid.UUID) ([]domain.Vehicle, error)
	CreateTransfer(ctx context.Context, t *domain.OwnershipTransfer) error
	UpdateVehicleOwner(ctx context.Context, vehicleID, newOwnerID uuid.UUID) error
	CreateLicense(ctx context.Context, l *domain.DriverLicense) error
	FindLicenseByCitizen(ctx context.Context, citizenID uuid.UUID) (*domain.DriverLicense, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateVehicle(ctx context.Context, v *domain.Vehicle) error {
	q := `INSERT INTO foves_vehicles (id, plate_number, vin, make, model, year, color, category, owner_citizen_id, is_stolen, is_active, registered_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.ExecContext(ctx, q,
		v.ID, v.PlateNumber, v.VIN, v.Make, v.Model, v.Year, v.Color, v.Category,
		v.OwnerCitizenID, v.IsStolen, v.IsActive, v.RegisteredAt, v.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) FindByPlate(ctx context.Context, plate string) (*domain.Vehicle, error) {
	q := `SELECT id, plate_number, vin, make, model, year, color, category, owner_citizen_id, is_stolen, is_active, registered_at, updated_at
		FROM foves_vehicles WHERE plate_number = $1`
	v := &domain.Vehicle{}
	err := r.db.QueryRowContext(ctx, q, plate).Scan(
		&v.ID, &v.PlateNumber, &v.VIN, &v.Make, &v.Model, &v.Year, &v.Color, &v.Category,
		&v.OwnerCitizenID, &v.IsStolen, &v.IsActive, &v.RegisteredAt, &v.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle not found: %s", plate)
		}
		return nil, fmt.Errorf("query vehicle: %w", err)
	}
	return v, nil
}

func (r *postgresRepo) FindByVIN(ctx context.Context, vin string) (*domain.Vehicle, error) {
	q := `SELECT id, plate_number, vin, make, model, year, color, category, owner_citizen_id, is_stolen, is_active, registered_at, updated_at
		FROM foves_vehicles WHERE vin = $1`
	v := &domain.Vehicle{}
	err := r.db.QueryRowContext(ctx, q, vin).Scan(
		&v.ID, &v.PlateNumber, &v.VIN, &v.Make, &v.Model, &v.Year, &v.Color, &v.Category,
		&v.OwnerCitizenID, &v.IsStolen, &v.IsActive, &v.RegisteredAt, &v.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle not found by vin: %s", vin)
		}
		return nil, fmt.Errorf("query vehicle: %w", err)
	}
	return v, nil
}

func (r *postgresRepo) FindByOwner(ctx context.Context, citizenID uuid.UUID) ([]domain.Vehicle, error) {
	q := `SELECT id, plate_number, vin, make, model, year, color, category, owner_citizen_id, is_stolen, is_active, registered_at, updated_at
		FROM foves_vehicles WHERE owner_citizen_id = $1 ORDER BY registered_at DESC`
	rows, err := r.db.QueryContext(ctx, q, citizenID)
	if err != nil {
		return nil, fmt.Errorf("query vehicles by owner: %w", err)
	}
	defer rows.Close()

	var vehicles []domain.Vehicle
	for rows.Next() {
		var v domain.Vehicle
		if err := rows.Scan(&v.ID, &v.PlateNumber, &v.VIN, &v.Make, &v.Model, &v.Year, &v.Color, &v.Category,
			&v.OwnerCitizenID, &v.IsStolen, &v.IsActive, &v.RegisteredAt, &v.UpdatedAt,
		); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, rows.Err()
}

func (r *postgresRepo) CreateTransfer(ctx context.Context, t *domain.OwnershipTransfer) error {
	q := `INSERT INTO foves_ownership_transfers (id, vehicle_id, from_citizen_id, to_citizen_id, transfer_date, contract_ref, approved_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, q,
		t.ID, t.VehicleID, t.FromCitizenID, t.ToCitizenID, t.TransferDate,
		t.ContractRef, t.ApprovedBy, t.CreatedAt,
	)
	return err
}

func (r *postgresRepo) UpdateVehicleOwner(ctx context.Context, vehicleID, newOwnerID uuid.UUID) error {
	now := time.Now().UTC()
	q := `UPDATE foves_vehicles SET owner_citizen_id = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, q, newOwnerID, now, vehicleID)
	return err
}

func (r *postgresRepo) CreateLicense(ctx context.Context, l *domain.DriverLicense) error {
	q := `INSERT INTO foves_driver_licenses (id, citizen_id, license_number, category_a, category_b, category_c, category_d, category_e, category_f, issued_date, expiry_date, points_balance, is_suspended, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := r.db.ExecContext(ctx, q,
		l.ID, l.CitizenID, l.LicenseNumber, l.CategoryA, l.CategoryB, l.CategoryC,
		l.CategoryD, l.CategoryE, l.CategoryF, l.IssuedDate, l.ExpiryDate,
		l.PointsBalance, l.IsSuspended, l.CreatedAt, l.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) FindLicenseByCitizen(ctx context.Context, citizenID uuid.UUID) (*domain.DriverLicense, error) {
	q := `SELECT id, citizen_id, license_number, category_a, category_b, category_c, category_d, category_e, category_f, issued_date, expiry_date, points_balance, is_suspended, created_at, updated_at
		FROM foves_driver_licenses WHERE citizen_id = $1`
	l := &domain.DriverLicense{}
	err := r.db.QueryRowContext(ctx, q, citizenID).Scan(
		&l.ID, &l.CitizenID, &l.LicenseNumber, &l.CategoryA, &l.CategoryB, &l.CategoryC,
		&l.CategoryD, &l.CategoryE, &l.CategoryF, &l.IssuedDate, &l.ExpiryDate,
		&l.PointsBalance, &l.IsSuspended, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("license not found for citizen: %s", citizenID)
		}
		return nil, fmt.Errorf("query license: %w", err)
	}
	return l, nil
}
