package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/rvin-svc/internal/domain"
)

type remainsRepo struct {
	pool *pgxpool.Pool
}

func NewRemainsRepo(pool *pgxpool.Pool) *remainsRepo {
	return &remainsRepo{pool: pool}
}

func (r *remainsRepo) Create(remains *domain.UnidentifiedRemains) (*domain.UnidentifiedRemains, error) {
	ctx := context.Background()
	remains.ID = uuid.New()
	remains.NationalRvinID = "RVIN-HT-" + time.Now().Format("2006") + "-" + remains.ID.String()[:6]
	remains.CreatedAt = time.Now()
	remains.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO rvin_unidentified_remains
		 (remains_id, national_rvin_id, discovery_date, discovery_location, dept_code, commune,
		  lat, lng, discovery_source, status, estimated_sex, estimated_age_min, estimated_age_max,
		  estimated_height_cm, skin_tone, distinguishing_marks, decomposition_level, morgue_location,
		  examiner_id, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)`,
		remains.ID, remains.NationalRvinID, remains.DiscoveryDate, remains.DiscoveryLocation,
		remains.DeptCode, remains.Commune, remains.Lat, remains.Lng, remains.DiscoverySource,
		remains.Status, remains.EstimatedSex, remains.EstimatedAgeMin, remains.EstimatedAgeMax,
		remains.EstimatedHeightCm, remains.SkinTone, remains.DistinguishingMarks,
		remains.DecompositionLevel, remains.MorgueLocation, remains.ExaminerID,
		remains.CreatedAt, remains.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return remains, nil
}

func (r *remainsRepo) FindByID(id uuid.UUID) (*domain.UnidentifiedRemains, error) {
	ctx := context.Background()
	rem := &domain.UnidentifiedRemains{}
	err := r.pool.QueryRow(ctx,
		`SELECT remains_id, national_rvin_id, discovery_date, discovery_location, dept_code, commune,
		        lat, lng, discovery_source, status, estimated_sex, estimated_age_min, estimated_age_max,
		        estimated_height_cm, distinguishing_marks, morgue_location, created_at, updated_at
		 FROM rvin_unidentified_remains WHERE remains_id = $1`, id).Scan(
		&rem.ID, &rem.NationalRvinID, &rem.DiscoveryDate, &rem.DiscoveryLocation,
		&rem.DeptCode, &rem.Commune, &rem.Lat, &rem.Lng, &rem.DiscoverySource,
		&rem.Status, &rem.EstimatedSex, &rem.EstimatedAgeMin, &rem.EstimatedAgeMax,
		&rem.EstimatedHeightCm, &rem.DistinguishingMarks, &rem.MorgueLocation,
		&rem.CreatedAt, &rem.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return rem, nil
}

func (r *remainsRepo) SubmitDNA(remainsID uuid.UUID, dna *domain.DNAResult) error {
	ctx := context.Background()
	dna.ComparisonID = uuid.New()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO rvin_dna_comparisons
		 (comparison_id, remains_id, reference_dna_ref, comparison_date, match_probability, is_match, lab_reference, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		dna.ComparisonID, dna.RemainsID, dna.ReferenceDNARef, dna.CreatedAt,
		dna.MatchProbability, dna.IsMatch, dna.LabReference, dna.CreatedAt,
	)
	return err
}

func (r *remainsRepo) FindUnidentified() ([]domain.UnidentifiedRemains, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT remains_id, national_rvin_id, discovery_date, discovery_location, dept_code,
		        discovery_source, status, estimated_sex, estimated_age_min, estimated_age_max
		 FROM rvin_unidentified_remains WHERE status = 'UNIDENTIFIED' ORDER BY discovery_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var remains []domain.UnidentifiedRemains
	for rows.Next() {
		var rem domain.UnidentifiedRemains
		if err := rows.Scan(&rem.ID, &rem.NationalRvinID, &rem.DiscoveryDate, &rem.DiscoveryLocation,
			&rem.DeptCode, &rem.DiscoverySource, &rem.Status, &rem.EstimatedSex,
			&rem.EstimatedAgeMin, &rem.EstimatedAgeMax); err != nil {
			return nil, err
		}
		remains = append(remains, rem)
	}
	return remains, nil
}

func (r *remainsRepo) GetStatsBySource() ([]domain.SourceStats, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT discovery_source, COUNT(*) as count FROM rvin_unidentified_remains GROUP BY discovery_source ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []domain.SourceStats
	for rows.Next() {
		var s domain.SourceStats
		if err := rows.Scan(&s.Source, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
