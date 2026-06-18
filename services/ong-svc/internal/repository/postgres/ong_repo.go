package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/ong-svc/internal/domain"
)

type ongRepo struct {
	pool *pgxpool.Pool
}

func NewONGRepo(pool *pgxpool.Pool) *ongRepo {
	return &ongRepo{pool: pool}
}

func (r *ongRepo) Create(org *domain.Organization) (*domain.Organization, error) {
	ctx := context.Background()
	org.ID = uuid.New()
	org.NationalONGID = "ONG-HT-" + time.Now().Format("2006") + "-" + org.ID.String()[:6]
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO ong_organizations
		 (org_id, national_ong_id, org_name, org_type, registration_status,
		  headquarter_country, operating_depts, sectors, director_name,
		  contact_email, contact_phone, risk_flag, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		org.ID, org.NationalONGID, org.OrgName, org.OrgType, org.RegistrationStatus,
		org.HeadquarterCountry, org.OperatingDepts, org.Sectors, org.DirectorName,
		org.ContactEmail, org.ContactPhone, org.RiskFlag, org.CreatedBy,
		org.CreatedAt, org.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (r *ongRepo) FindByID(id uuid.UUID) (*domain.Organization, error) {
	ctx := context.Background()
	o := &domain.Organization{}
	err := r.pool.QueryRow(ctx,
		`SELECT org_id, national_ong_id, org_name, org_type, registration_status,
		        headquarter_country, operating_depts, sectors, risk_flag, created_at
		 FROM ong_organizations WHERE org_id = $1`, id).Scan(
		&o.ID, &o.NationalONGID, &o.OrgName, &o.OrgType, &o.RegistrationStatus,
		&o.HeadquarterCountry, &o.OperatingDepts, &o.Sectors, &o.RiskFlag, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *ongRepo) FindAll() ([]domain.Organization, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT org_id, national_ong_id, org_name, org_type, registration_status,
		        headquarter_country, risk_flag, created_at
		 FROM ong_organizations ORDER BY org_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []domain.Organization
	for rows.Next() {
		var o domain.Organization
		if err := rows.Scan(&o.ID, &o.NationalONGID, &o.OrgName, &o.OrgType,
			&o.RegistrationStatus, &o.HeadquarterCountry, &o.RiskFlag, &o.CreatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, nil
}

func (r *ongRepo) FindFlagged() ([]domain.Organization, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT org_id, national_ong_id, org_name, org_type, registration_status,
		        headquarter_country, risk_flag, created_at
		 FROM ong_organizations WHERE risk_flag != 'NONE' ORDER BY risk_flag`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []domain.Organization
	for rows.Next() {
		var o domain.Organization
		if err := rows.Scan(&o.ID, &o.NationalONGID, &o.OrgName, &o.OrgType,
			&o.RegistrationStatus, &o.HeadquarterCountry, &o.RiskFlag, &o.CreatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, nil
}

func (r *ongRepo) FindUnregistered() ([]domain.Organization, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT org_id, national_ong_id, org_name, org_type, registration_status, headquarter_country
		 FROM ong_organizations WHERE registration_status = 'OPERATING_WITHOUT_REGISTRATION'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []domain.Organization
	for rows.Next() {
		var o domain.Organization
		if err := rows.Scan(&o.ID, &o.NationalONGID, &o.OrgName, &o.OrgType,
			&o.RegistrationStatus, &o.HeadquarterCountry); err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, nil
}

func (r *ongRepo) UpdateRiskFlag(id uuid.UUID, flag domain.RiskFlag) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE ong_organizations SET risk_flag = $1, updated_at = NOW() WHERE org_id = $2`,
		flag, id)
	return err
}

func (r *ongRepo) CreateStaff(staff *domain.Staff) (*domain.Staff, error) {
	ctx := context.Background()
	staff.ID = uuid.New()
	staff.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO ong_staff_registry
		 (staff_id, org_id, full_name, nationality, role, is_expatriate, passport_number, is_active, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		staff.ID, staff.OrgID, staff.FullName, staff.Nationality, staff.Role,
		staff.IsExpatriate, staff.PassportNumber, staff.IsActive, staff.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (r *ongRepo) CreateAccessRequest(ar *domain.AccessRequest) (*domain.AccessRequest, error) {
	ctx := context.Background()
	ar.ID = uuid.New()
	ar.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO ong_field_access_requests
		 (request_id, org_id, access_type, requested_zones, access_date, purpose, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		ar.ID, ar.OrgID, ar.AccessType, ar.RequestedZones, ar.AccessDate,
		ar.Purpose, ar.Status, ar.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (r *ongRepo) UpdateAccessStatus(id uuid.UUID, status string, notes string) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE ong_field_access_requests SET status = $1, approval_notes = $2 WHERE request_id = $3`,
		status, notes, id)
	return err
}
