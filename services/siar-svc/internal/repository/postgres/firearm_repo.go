package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/siar-svc/internal/domain"
)

type firearmRepo struct {
	pool *pgxpool.Pool
}

func NewFirearmRepo(pool *pgxpool.Pool) *firearmRepo {
	return &firearmRepo{pool: pool}
}

func (r *firearmRepo) CreateFirearm(f *domain.Firearm) (*domain.Firearm, error) {
	ctx := context.Background()
	f.ID = uuid.New()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO siar_firearms
		 (id, national_siar_id, serial_number, make, model, caliber, weapon_type, manufacture_year,
		  manufacture_country, status, reg_type, owner_snisid_id, owner_entity_name, license_number,
		  license_expiry, import_date, import_country, import_permit_ref, importer_name, customs_entry_ref,
		  current_dept_code, storage_location, fir_record_id, gang_id, case_references, iarms_ref,
		  atf_etrace_ref, notes, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31)
		 RETURNING created_at, updated_at`,
		f.ID, f.NationalSiarID, f.SerialNumber, f.Make, f.Model, f.Caliber, f.WeaponType,
		f.ManufactureYear, f.ManufactureCountry, f.Status, f.RegType, f.OwnerSnisidID,
		f.OwnerEntityName, f.LicenseNumber, f.LicenseExpiry, f.ImportDate, f.ImportCountry,
		f.ImportPermitRef, f.ImporterName, f.CustomsEntryRef, f.CurrentDeptCode, f.StorageLocation,
		f.FirRecordID, f.GangID, f.CaseReferences, f.IarmsRef, f.AtfEtraceRef, f.Notes,
		f.CreatedBy, f.CreatedAt, f.UpdatedAt,
	).Scan(&f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *firearmRepo) FindBySerial(serialNumber string) (*domain.Firearm, error) {
	ctx := context.Background()
	f := &domain.Firearm{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, national_siar_id, serial_number, make, model, caliber, weapon_type, manufacture_year,
		        manufacture_country, status, reg_type, owner_snisid_id, owner_entity_name, license_number,
		        license_expiry, import_date, import_country, import_permit_ref, importer_name, customs_entry_ref,
		        current_dept_code, storage_location, fir_record_id, gang_id, case_references, iarms_ref,
		        atf_etrace_ref, notes, created_by, created_at, updated_at
		 FROM siar_firearms WHERE serial_number = $1`, serialNumber).Scan(
		&f.ID, &f.NationalSiarID, &f.SerialNumber, &f.Make, &f.Model, &f.Caliber, &f.WeaponType,
		&f.ManufactureYear, &f.ManufactureCountry, &f.Status, &f.RegType, &f.OwnerSnisidID,
		&f.OwnerEntityName, &f.LicenseNumber, &f.LicenseExpiry, &f.ImportDate, &f.ImportCountry,
		&f.ImportPermitRef, &f.ImporterName, &f.CustomsEntryRef, &f.CurrentDeptCode, &f.StorageLocation,
		&f.FirRecordID, &f.GangID, &f.CaseReferences, &f.IarmsRef, &f.AtfEtraceRef, &f.Notes,
		&f.CreatedBy, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *firearmRepo) FindByID(id uuid.UUID) (*domain.Firearm, error) {
	ctx := context.Background()
	f := &domain.Firearm{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, national_siar_id, serial_number, make, model, caliber, weapon_type, manufacture_year,
		        manufacture_country, status, reg_type, owner_snisid_id, owner_entity_name, license_number,
		        license_expiry, import_date, import_country, import_permit_ref, importer_name, customs_entry_ref,
		        current_dept_code, storage_location, fir_record_id, gang_id, case_references, iarms_ref,
		        atf_etrace_ref, notes, created_by, created_at, updated_at
		 FROM siar_firearms WHERE id = $1`, id).Scan(
		&f.ID, &f.NationalSiarID, &f.SerialNumber, &f.Make, &f.Model, &f.Caliber, &f.WeaponType,
		&f.ManufactureYear, &f.ManufactureCountry, &f.Status, &f.RegType, &f.OwnerSnisidID,
		&f.OwnerEntityName, &f.LicenseNumber, &f.LicenseExpiry, &f.ImportDate, &f.ImportCountry,
		&f.ImportPermitRef, &f.ImporterName, &f.CustomsEntryRef, &f.CurrentDeptCode, &f.StorageLocation,
		&f.FirRecordID, &f.GangID, &f.CaseReferences, &f.IarmsRef, &f.AtfEtraceRef, &f.Notes,
		&f.CreatedBy, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *firearmRepo) UpdateFirearm(f *domain.Firearm) error {
	ctx := context.Background()
	f.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx,
		`UPDATE siar_firearms SET
		 serial_number = $2, make = $3, model = $4, caliber = $5, weapon_type = $6,
		 manufacture_year = $7, manufacture_country = $8, status = $9, reg_type = $10,
		 owner_snisid_id = $11, owner_entity_name = $12, license_number = $13, license_expiry = $14,
		 import_date = $15, import_country = $16, import_permit_ref = $17, importer_name = $18,
		 customs_entry_ref = $19, current_dept_code = $20, storage_location = $21, fir_record_id = $22,
		 gang_id = $23, case_references = $24, iarms_ref = $25, atf_etrace_ref = $26, notes = $27,
		 updated_at = $28
		 WHERE id = $1`,
		f.ID, f.SerialNumber, f.Make, f.Model, f.Caliber, f.WeaponType,
		f.ManufactureYear, f.ManufactureCountry, f.Status, f.RegType, f.OwnerSnisidID,
		f.OwnerEntityName, f.LicenseNumber, f.LicenseExpiry, f.ImportDate, f.ImportCountry,
		f.ImportPermitRef, f.ImporterName, f.CustomsEntryRef, f.CurrentDeptCode, f.StorageLocation,
		f.FirRecordID, f.GangID, f.CaseReferences, f.IarmsRef, f.AtfEtraceRef, f.Notes,
		f.UpdatedAt)
	return err
}

func (r *firearmRepo) CreateSeizure(s *domain.Seizure) error {
	ctx := context.Background()
	s.ID = uuid.New()
	s.CreatedAt = time.Now()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO siar_seizures
		 (id, firearm_id, seizure_date, seizing_unit, seizing_officer, location_desc,
		  dept_code, context, from_person_id, gang_id, case_reference, disposed_of,
		  disposal_method, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		s.ID, s.FirearmID, s.SeizureDate, s.SeizingUnit, s.SeizingOfficer, s.LocationDesc,
		s.DeptCode, s.Context, s.FromPersonID, s.GangID, s.CaseReference, s.DisposedOf,
		s.DisposalMethod, s.CreatedAt)
	return err
}

func (r *firearmRepo) CreateLicense(l *domain.License) error {
	ctx := context.Background()
	l.ID = uuid.New()
	l.CreatedAt = time.Now()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO siar_licenses
		 (id, license_number, holder_snisid_id, license_type, firearms_authorized, issue_date,
		  expiry_date, issuing_authority, is_active, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		l.ID, l.LicenseNumber, l.HolderSnisidID, l.LicenseType, l.FirearmsAuthorized,
		l.IssueDate, l.ExpiryDate, l.IssuingAuthority, l.IsActive, l.CreatedAt)
	return err
}

func (r *firearmRepo) GetLicensesByPerson(personID uuid.UUID) ([]domain.License, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT id, license_number, holder_snisid_id, license_type, firearms_authorized,
		        issue_date, expiry_date, issuing_authority, is_active, revocation_reason,
		        revoked_at, created_at
		 FROM siar_licenses WHERE holder_snisid_id = $1 ORDER BY created_at DESC`, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var licenses []domain.License
	for rows.Next() {
		var l domain.License
		if err := rows.Scan(
			&l.ID, &l.LicenseNumber, &l.HolderSnisidID, &l.LicenseType, &l.FirearmsAuthorized,
			&l.IssueDate, &l.ExpiryDate, &l.IssuingAuthority, &l.IsActive, &l.RevocationReason,
			&l.RevokedAt, &l.CreatedAt); err != nil {
			return nil, err
		}
		licenses = append(licenses, l)
	}
	return licenses, nil
}

func (r *firearmRepo) CreateTransfer(t *domain.Transfer) error {
	ctx := context.Background()
	t.ID = uuid.New()
	t.CreatedAt = time.Now()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO siar_transfers
		 (id, firearm_id, from_owner_id, to_owner_id, transfer_type, transfer_date,
		  permit_ref, authorized_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		t.ID, t.FirearmID, t.FromOwnerID, t.ToOwnerID, t.TransferType, t.TransferDate,
		t.PermitRef, t.AuthorizedBy, t.CreatedAt)
	return err
}

func (r *firearmRepo) GetStatsByType() ([]domain.StatsByType, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT weapon_type, COUNT(*) as count FROM siar_firearms GROUP BY weapon_type ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []domain.StatsByType
	for rows.Next() {
		var s domain.StatsByType
		if err := rows.Scan(&s.WeaponType, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func scanFirearms(rows pgx.Rows) ([]domain.Firearm, error) {
	var firearms []domain.Firearm
	for rows.Next() {
		var f domain.Firearm
		if err := rows.Scan(
			&f.ID, &f.NationalSiarID, &f.SerialNumber, &f.Make, &f.Model, &f.Caliber, &f.WeaponType,
			&f.ManufactureYear, &f.ManufactureCountry, &f.Status, &f.RegType, &f.OwnerSnisidID,
			&f.OwnerEntityName, &f.LicenseNumber, &f.LicenseExpiry, &f.ImportDate, &f.ImportCountry,
			&f.ImportPermitRef, &f.ImporterName, &f.CustomsEntryRef, &f.CurrentDeptCode, &f.StorageLocation,
			&f.FirRecordID, &f.GangID, &f.CaseReferences, &f.IarmsRef, &f.AtfEtraceRef, &f.Notes,
			&f.CreatedBy, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		firearms = append(firearms, f)
	}
	return firearms, nil
}
