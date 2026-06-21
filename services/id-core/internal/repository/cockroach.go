package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/idcore-svc/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, citizen *domain.Citizen) error
	FindByNIN(ctx context.Context, nin string) (*domain.Citizen, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Citizen, error)
	Update(ctx context.Context, citizen *domain.Citizen) error
	UpdateStatus(ctx context.Context, nin string, status domain.IDStatus, reason string, authorizedBy uuid.UUID) error
	FindDemographicMatches(ctx context.Context, fullName string, dob time.Time) ([]domain.DemographicMatch, error)
	CreateDedupCandidate(ctx context.Context, candidate domain.DedupCandidate) error
	GetHistory(ctx context.Context, citizenID uuid.UUID) ([]domain.ChangeHistory, error)
	GetPopulationStats(ctx context.Context) (*domain.PopulationStats, error)
}

type cockroachRepo struct {
	db *sql.DB
}

func NewCockroachRepo(db *sql.DB) Repository {
	return &cockroachRepo{db: db}
}

func (r *cockroachRepo) Create(ctx context.Context, citizen *domain.Citizen) error {
	query := `INSERT INTO citizens (
		citizen_id, nin, status, enrollment_type, full_name_legal,
		first_name, middle_names, last_name, maiden_name,
		dob, pob_commune, pob_dept_code, gender, nationality,
		dept_code, current_address, current_commune,
		biometric_template_id, photo_ref,
		mother_nin, father_nin,
		created_by, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5,
		$6, $7, $8, $9,
		$10, $11, $12, $13, $14,
		$15, $16, $17,
		$18, $19,
		$20, $21,
		$22, $23, $24
	)`

	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query,
		citizen.CitizenID,
		citizen.NIN,
		citizen.Status,
		citizen.EnrollmentType,
		citizen.FullNameLegal,
		citizen.FirstName,
		citizen.MiddleNames,
		citizen.LastName,
		citizen.MaidenName,
		citizen.DOB,
		citizen.PobCommune,
		citizen.PobDeptCode,
		citizen.Gender,
		citizen.Nationality,
		citizen.DeptCode,
		citizen.CurrentAddress,
		citizen.CurrentCommune,
		citizen.BiometricTemplateID,
		citizen.PhotoRef,
		citizen.MotherNIN,
		citizen.FatherNIN,
		citizen.CreatedBy,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("create citizen: %w", err)
	}
	citizen.CreatedAt = now
	citizen.UpdatedAt = now
	return nil
}

func (r *cockroachRepo) FindByNIN(ctx context.Context, nin string) (*domain.Citizen, error) {
	query := `SELECT
		citizen_id, nin, status, enrollment_type, full_name_legal,
		first_name, middle_names, last_name, maiden_name,
		dob, pob_commune, pob_dept_code, gender, nationality,
		dept_code, current_address, current_commune,
		biometric_template_id, photo_ref,
		mother_nin, father_nin,
		date_of_death, death_certificate_ref,
		is_merged, merged_into_citizen_id,
		created_by, created_at, updated_at
	FROM citizens WHERE nin = $1`

	return r.scanCitizen(ctx, query, nin)
}

func (r *cockroachRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Citizen, error) {
	query := `SELECT
		citizen_id, nin, status, enrollment_type, full_name_legal,
		first_name, middle_names, last_name, maiden_name,
		dob, pob_commune, pob_dept_code, gender, nationality,
		dept_code, current_address, current_commune,
		biometric_template_id, photo_ref,
		mother_nin, father_nin,
		date_of_death, death_certificate_ref,
		is_merged, merged_into_citizen_id,
		created_by, created_at, updated_at
	FROM citizens WHERE citizen_id = $1`

	return r.scanCitizen(ctx, query, id)
}

func (r *cockroachRepo) Update(ctx context.Context, citizen *domain.Citizen) error {
	query := `UPDATE citizens SET
		full_name_legal = $1, first_name = $2, middle_names = $3, last_name = $4, maiden_name = $5,
		pob_commune = $6, pob_dept_code = $7, gender = $8,
		current_address = $9, current_commune = $10,
		biometric_template_id = $11, photo_ref = $12,
		mother_nin = $13, father_nin = $14,
		date_of_death = $15, death_certificate_ref = $16,
		is_merged = $17, merged_into_citizen_id = $18,
		updated_at = $19
	WHERE citizen_id = $20`

	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query,
		citizen.FullNameLegal, citizen.FirstName, citizen.MiddleNames, citizen.LastName, citizen.MaidenName,
		citizen.PobCommune, citizen.PobDeptCode, citizen.Gender,
		citizen.CurrentAddress, citizen.CurrentCommune,
		citizen.BiometricTemplateID, citizen.PhotoRef,
		citizen.MotherNIN, citizen.FatherNIN,
		citizen.DateOfDeath, citizen.DeathCertificateRef,
		citizen.IsMerged, citizen.MergedIntoCitizenID,
		now,
		citizen.CitizenID,
	)
	if err != nil {
		return fmt.Errorf("update citizen: %w", err)
	}
	citizen.UpdatedAt = now
	return nil
}

func (r *cockroachRepo) UpdateStatus(ctx context.Context, nin string, status domain.IDStatus, reason string, authorizedBy uuid.UUID) error {
	query := `UPDATE citizens SET status = $1, updated_at = $2 WHERE nin = $3`
	now := time.Now().UTC()
	res, err := r.db.ExecContext(ctx, query, status, now, nin)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("citizen not found: %s", nin)
	}

	historyQuery := `INSERT INTO id_change_history (history_id, citizen_id, field_changed, old_value, new_value, change_reason, authorized_by, changed_at)
		SELECT gen_random_uuid(), citizen_id, 'status', status, $1, $2, $3, $4 FROM citizens WHERE nin = $5`
	_, err = r.db.ExecContext(ctx, historyQuery, status, reason, authorizedBy, now, nin)
	return err
}

func (r *cockroachRepo) FindDemographicMatches(ctx context.Context, fullName string, dob time.Time) ([]domain.DemographicMatch, error) {
	query := `SELECT citizen_id, 0.9 AS score FROM citizens
		WHERE full_name_legal % $1 AND dob = $2
		LIMIT 10`

	rows, err := r.db.QueryContext(ctx, query, fullName, dob)
	if err != nil {
		return nil, fmt.Errorf("demographic search: %w", err)
	}
	defer rows.Close()

	var matches []domain.DemographicMatch
	for rows.Next() {
		var m domain.DemographicMatch
		if err := rows.Scan(&m.CitizenID, &m.Score); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

func (r *cockroachRepo) CreateDedupCandidate(ctx context.Context, candidate domain.DedupCandidate) error {
	query := `INSERT INTO id_dedup_candidates (candidate_id, citizen_id_a, citizen_id_b, biometric_score, demographic_score, composite_score)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		candidate.CitizenIDA, candidate.CitizenIDB,
		candidate.BiometricScore, candidate.DemographicScore, candidate.CompositeScore,
	)
	return err
}

func (r *cockroachRepo) GetHistory(ctx context.Context, citizenID uuid.UUID) ([]domain.ChangeHistory, error) {
	query := `SELECT history_id, citizen_id, field_changed, old_value, new_value, change_reason, authorized_by, changed_at
		FROM id_change_history WHERE citizen_id = $1 ORDER BY changed_at DESC`

	rows, err := r.db.QueryContext(ctx, query, citizenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.ChangeHistory
	for rows.Next() {
		var h domain.ChangeHistory
		if err := rows.Scan(&h.HistoryID, &h.CitizenID, &h.FieldChanged, &h.OldValue, &h.NewValue, &h.ChangeReason, &h.AuthorizedBy, &h.ChangedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, rows.Err()
}

func (r *cockroachRepo) GetPopulationStats(ctx context.Context) (*domain.PopulationStats, error) {
	query := `SELECT
		COUNT(*) AS total,
		COUNT(*) FILTER (WHERE status = 'ACTIVE') AS active,
		COUNT(*) FILTER (WHERE status = 'SUSPENDED') AS suspended,
		COUNT(*) FILTER (WHERE status = 'DECEASED') AS deceased,
		COUNT(*) FILTER (WHERE status = 'CANCELLED') AS cancelled,
		COUNT(*) FILTER (WHERE is_merged) AS merged
	FROM citizens`

	stats := &domain.PopulationStats{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.Total, &stats.Active, &stats.Suspended,
		&stats.Deceased, &stats.Cancelled, &stats.Merged,
	)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *cockroachRepo) scanCitizen(ctx context.Context, query string, args ...any) (*domain.Citizen, error) {
	c := &domain.Citizen{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&c.CitizenID, &c.NIN, &c.Status, &c.EnrollmentType, &c.FullNameLegal,
		&c.FirstName, &c.MiddleNames, &c.LastName, &c.MaidenName,
		&c.DOB, &c.PobCommune, &c.PobDeptCode, &c.Gender, &c.Nationality,
		&c.DeptCode, &c.CurrentAddress, &c.CurrentCommune,
		&c.BiometricTemplateID, &c.PhotoRef,
		&c.MotherNIN, &c.FatherNIN,
		&c.DateOfDeath, &c.DeathCertificateRef,
		&c.IsMerged, &c.MergedIntoCitizenID,
		&c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("citizen not found")
		}
		return nil, fmt.Errorf("query citizen: %w", err)
	}
	return c, nil
}
