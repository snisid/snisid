package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/enfl-svc/internal/domain"
)

type childRepo struct {
	pool *pgxpool.Pool
}

func NewChildRepo(pool *pgxpool.Pool) *childRepo {
	return &childRepo{pool: pool}
}

func (r *childRepo) Create(child *domain.Child) (*domain.Child, error) {
	ctx := context.Background()
	child.ID = uuid.New()
	child.NationalEnflID = "ENFL-HT-" + time.Now().Format("2006") + "-" + child.ID.String()[:6]
	child.CreatedAt = time.Now()
	child.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO enfl_children
		 (child_id, national_enfl_id, risk_category, status, full_name, dob, gender, nationality,
		  photo_refs, distinguishing_marks, height_cm, skin_tone, guardian_name, guardian_phone,
		  dept_code, commune, disappearance_date, gang_id, assistance_type, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)`,
		child.ID, child.NationalEnflID, child.RiskCategory, child.Status, child.FullName,
		child.DOB, child.Gender, child.Nationality, child.PhotoRefs, child.DistinguishingMarks,
		child.HeightCm, child.SkinTone, child.GuardianName, child.GuardianPhone,
		child.DeptCode, child.Commune, child.DisappearanceDate, child.GangID,
		child.AssistanceType, child.CreatedBy, child.CreatedAt, child.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return child, nil
}

func (r *childRepo) FindByID(id uuid.UUID) (*domain.Child, error) {
	ctx := context.Background()
	child := &domain.Child{}
	err := r.pool.QueryRow(ctx,
		`SELECT child_id, national_enfl_id, risk_category, status, full_name, dob, gender, nationality,
		        photo_refs, distinguishing_marks, height_cm, skin_tone, guardian_name, guardian_phone,
		        dept_code, commune, disappearance_date, gang_id, assistance_type, created_at, updated_at
		 FROM enfl_children WHERE child_id = $1`, id).Scan(
		&child.ID, &child.NationalEnflID, &child.RiskCategory, &child.Status, &child.FullName,
		&child.DOB, &child.Gender, &child.Nationality, &child.PhotoRefs, &child.DistinguishingMarks,
		&child.HeightCm, &child.SkinTone, &child.GuardianName, &child.GuardianPhone,
		&child.DeptCode, &child.Commune, &child.DisappearanceDate, &child.GangID,
		&child.AssistanceType, &child.CreatedAt, &child.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return child, nil
}

func (r *childRepo) FindMissing() ([]domain.Child, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT child_id, national_enfl_id, risk_category, status, full_name, dob, gender, nationality,
		        photo_refs, dept_code, commune, disappearance_date, created_at
		 FROM enfl_children WHERE status = 'MISSING' ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []domain.Child
	for rows.Next() {
		var c domain.Child
		if err := rows.Scan(&c.ID, &c.NationalEnflID, &c.RiskCategory, &c.Status, &c.FullName,
			&c.DOB, &c.Gender, &c.Nationality, &c.PhotoRefs, &c.DeptCode, &c.Commune,
			&c.DisappearanceDate, &c.CreatedAt); err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, nil
}

func (r *childRepo) FindRestaveks() ([]domain.Restavek, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT r.restavek_id, r.child_id, r.employing_household, r.household_dept, r.household_commune,
		        r.reported_conditions, r.school_attendance, r.ibesr_inspection, r.last_inspection_date, r.created_at
		 FROM enfl_restaveks r
		 JOIN enfl_children c ON r.child_id = c.child_id
		 WHERE c.risk_category = 'DOMESTIC_SERVITUDE_RESTAVEK'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaveks []domain.Restavek
	for rows.Next() {
		var r domain.Restavek
		if err := rows.Scan(&r.RestavekID, &r.ChildID, &r.EmployingHousehold, &r.HouseholdDept,
			&r.HouseholdCommune, &r.ReportedConditions, &r.SchoolAttendance, &r.IbesrInspection,
			&r.LastInspectionDate, &r.CreatedAt); err != nil {
			return nil, err
		}
		restaveks = append(restaveks, r)
	}
	return restaveks, nil
}

func (r *childRepo) UpdateStatus(id uuid.UUID, status domain.ChildStatus, location string) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE enfl_children SET status = $1, last_known_location = $2, updated_at = NOW() WHERE child_id = $3`,
		status, location, id)
	return err
}

func (r *childRepo) FindGangRecruited() ([]domain.Child, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT child_id, national_enfl_id, risk_category, status, full_name, dob, gender,
		        dept_code, commune, gang_id, created_at
		 FROM enfl_children WHERE risk_category = 'GANG_RECRUITMENT' AND gang_id IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []domain.Child
	for rows.Next() {
		var c domain.Child
		if err := rows.Scan(&c.ID, &c.NationalEnflID, &c.RiskCategory, &c.Status, &c.FullName,
			&c.DOB, &c.Gender, &c.DeptCode, &c.Commune, &c.GangID, &c.CreatedAt); err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, nil
}
