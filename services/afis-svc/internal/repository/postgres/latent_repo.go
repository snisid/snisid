package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/snisid/platform/services/afis-svc/internal/domain"
)

type LatentRepo struct {
	pool *pgxpool.Pool
}

func NewLatentRepo(pool *pgxpool.Pool) *LatentRepo {
	return &LatentRepo{pool: pool}
}

func (r *LatentRepo) Create(ctx context.Context, lp *domain.LatentPrint) error {
	query := `INSERT INTO afis_latent_prints (latent_id, case_reference, crime_scene_id, location_desc, dept_code, found_at, image_ref, finger_position, examined_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query,
		lp.LatentID, lp.CaseReference, lp.CrimeSceneID, lp.LocationDesc, lp.DeptCode,
		lp.FoundAt, lp.ImageRef, lp.FingerPosition, lp.ExaminedBy,
	)
	if err != nil {
		return fmt.Errorf("insert latent: %w", err)
	}
	return nil
}

func (r *LatentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.LatentPrint, error) {
	query := `SELECT latent_id, case_reference, crime_scene_id, location_desc, dept_code, found_at, image_ref, nfiq2_score, finger_position, is_identified, matched_subject_id, match_score, examined_by, created_at
		FROM afis_latent_prints WHERE latent_id = $1`
	lp := &domain.LatentPrint{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&lp.LatentID, &lp.CaseReference, &lp.CrimeSceneID, &lp.LocationDesc, &lp.DeptCode,
		&lp.FoundAt, &lp.ImageRef, &lp.NFIQ2Score, &lp.FingerPosition, &lp.IsIdentified,
		&lp.MatchedSubjectID, &lp.MatchScore, &lp.ExaminedBy, &lp.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get latent: %w", err)
	}
	return lp, nil
}

func (r *LatentRepo) ConfirmMatch(ctx context.Context, latentID, subjectID uuid.UUID, score float64, examiner uuid.UUID) error {
	query := `UPDATE afis_latent_prints SET is_identified = TRUE, matched_subject_id = $1, match_score = $2, examined_by = $3 WHERE latent_id = $4`
	_, err := r.pool.Exec(ctx, query, subjectID, score, examiner, latentID)
	return err
}

func (r *LatentRepo) GetUnidentified(ctx context.Context) ([]domain.LatentPrint, error) {
	query := `SELECT latent_id, case_reference, crime_scene_id, location_desc, dept_code, found_at, image_ref, nfiq2_score, finger_position, is_identified, matched_subject_id, match_score, examined_by, created_at
		FROM afis_latent_prints WHERE is_identified = FALSE ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get unidentified: %w", err)
	}
	defer rows.Close()

	var results []domain.LatentPrint
	for rows.Next() {
		var lp domain.LatentPrint
		if err := rows.Scan(&lp.LatentID, &lp.CaseReference, &lp.CrimeSceneID, &lp.LocationDesc, &lp.DeptCode,
			&lp.FoundAt, &lp.ImageRef, &lp.NFIQ2Score, &lp.FingerPosition, &lp.IsIdentified,
			&lp.MatchedSubjectID, &lp.MatchScore, &lp.ExaminedBy, &lp.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan latent: %w", err)
		}
		results = append(results, lp)
	}
	return results, nil
}
