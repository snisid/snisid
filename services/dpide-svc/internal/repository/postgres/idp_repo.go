package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/dpide-svc/internal/domain"
)

type idpRepo struct {
	pool *pgxpool.Pool
}

func NewIDPRepo(pool *pgxpool.Pool) *idpRepo {
	return &idpRepo{pool: pool}
}

func (r *idpRepo) Create(idp *domain.IDP) (*domain.IDP, error) {
	ctx := context.Background()
	idp.ID = uuid.New()
	idp.NationalDpideID = "DPIDE-HT-" + time.Now().Format("2006") + "-" + idp.ID.String()[:6]
	idp.CreatedAt = time.Now()
	idp.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO dpide_idps
		 (idp_id, national_dpide_id, full_name, dob, gender, household_size, minors_count,
		  displacement_cause, displacement_date, origin_dept_code, origin_commune,
		  status, current_location, current_dept_code, current_commune, current_lat, current_lng,
		  shelter_type, medical_needs, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)`,
		idp.ID, idp.NationalDpideID, idp.FullName, idp.DOB, idp.Gender, idp.HouseholdSize,
		idp.MinorsCount, idp.DisplacementCause, idp.DisplacementDate, idp.OriginDeptCode,
		idp.OriginCommune, idp.Status, idp.CurrentLocation, idp.CurrentDeptCode,
		idp.CurrentCommune, idp.CurrentLat, idp.CurrentLng, idp.ShelterType,
		idp.MedicalNeeds, idp.CreatedAt, idp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return idp, nil
}

func (r *idpRepo) FindByID(id uuid.UUID) (*domain.IDP, error) {
	ctx := context.Background()
	idp := &domain.IDP{}
	err := r.pool.QueryRow(ctx,
		`SELECT idp_id, national_dpide_id, full_name, dob, gender, displacement_cause,
		        displacement_date, origin_dept_code, origin_commune, status,
		        current_location, current_dept_code, current_commune, current_lat, current_lng,
		        shelter_type, medical_needs, created_at, updated_at
		 FROM dpide_idps WHERE idp_id = $1`, id).Scan(
		&idp.ID, &idp.NationalDpideID, &idp.FullName, &idp.DOB, &idp.Gender,
		&idp.DisplacementCause, &idp.DisplacementDate, &idp.OriginDeptCode, &idp.OriginCommune,
		&idp.Status, &idp.CurrentLocation, &idp.CurrentDeptCode, &idp.CurrentCommune,
		&idp.CurrentLat, &idp.CurrentLng, &idp.ShelterType, &idp.MedicalNeeds,
		&idp.CreatedAt, &idp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return idp, nil
}

func (r *idpRepo) FindCamps() ([]domain.Camp, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT camp_id, camp_name, dept_code, commune, lat, lng, displacement_cause,
		        managing_org, capacity, current_population, is_active, has_medical_post, has_school, created_at
		 FROM dpide_camps WHERE is_active = TRUE ORDER BY camp_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var camps []domain.Camp
	for rows.Next() {
		var c domain.Camp
		if err := rows.Scan(&c.ID, &c.CampName, &c.DeptCode, &c.Commune, &c.Lat, &c.Lng,
			&c.DisplacementCause, &c.ManagingOrg, &c.Capacity, &c.CurrentPopulation,
			&c.IsActive, &c.HasMedicalPost, &c.HasSchool, &c.CreatedAt); err != nil {
			return nil, err
		}
		camps = append(camps, c)
	}
	return camps, nil
}

func (r *idpRepo) GetStats() (*domain.IDPStats, error) {
	ctx := context.Background()
	stats := &domain.IDPStats{}
	err := r.pool.QueryRow(ctx,
		`SELECT
			COUNT(*) AS total_idps,
			COUNT(*) FILTER (WHERE status = 'DISPLACED') AS displaced_count,
			COUNT(*) FILTER (WHERE status = 'IN_CAMP') AS in_camp_count,
			COUNT(*) FILTER (WHERE status = 'RETURNED_HOME') AS returned_count,
			(SELECT COUNT(*) FROM dpide_camps WHERE is_active = TRUE) AS camp_count
		 FROM dpide_idps`).Scan(&stats.TotalIDPs, &stats.DisplacedCount, &stats.InCampCount,
		&stats.ReturnedCount, &stats.CampCount)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *idpRepo) UpdateStatus(id uuid.UUID, status domain.IDPStatus) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE dpide_idps SET status = $1, updated_at = NOW() WHERE idp_id = $2`,
		status, id)
	return err
}
