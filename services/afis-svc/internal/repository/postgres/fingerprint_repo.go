package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/snisid/platform/services/afis-svc/internal/domain"
)

type FingerprintRepo struct {
	pool *pgxpool.Pool
}

func NewFingerprintRepo(pool *pgxpool.Pool) *FingerprintRepo {
	return &FingerprintRepo{pool: pool}
}

func (r *FingerprintRepo) Create(ctx context.Context, fp *domain.Fingerprint) error {
	query := `INSERT INTO afis_fingerprints (print_id, subject_id, finger_position, capture_method, nfiq2_score, image_ref, minutiae_count, milvus_vector_id, is_primary, captured_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(ctx, query,
		fp.PrintID, fp.SubjectID, fp.FingerPosition, fp.CaptureMethod, fp.NFIQ2Score,
		fp.ImageRef, fp.MinutiaeCount, fp.MilvusVectorID, false, fp.CapturedAt, fp.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("insert fingerprint: %w", err)
	}
	return nil
}

func (r *FingerprintRepo) GetBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]domain.Fingerprint, error) {
	query := `SELECT print_id, subject_id, finger_position, capture_method, nfiq2_score, quality_accepted, image_ref, minutiae_count, milvus_vector_id, captured_at, created_by
		FROM afis_fingerprints WHERE subject_id = $1 ORDER BY finger_position`
	rows, err := r.pool.Query(ctx, query, subjectID)
	if err != nil {
		return nil, fmt.Errorf("get fingerprints: %w", err)
	}
	defer rows.Close()

	var results []domain.Fingerprint
	for rows.Next() {
		var fp domain.Fingerprint
		if err := rows.Scan(&fp.PrintID, &fp.SubjectID, &fp.FingerPosition, &fp.CaptureMethod,
			&fp.NFIQ2Score, &fp.QualityAccepted, &fp.ImageRef, &fp.MinutiaeCount,
			&fp.MilvusVectorID, &fp.CapturedAt, &fp.CreatedBy); err != nil {
			return nil, fmt.Errorf("scan fingerprint: %w", err)
		}
		results = append(results, fp)
	}
	return results, nil
}

func (r *FingerprintRepo) UpdateMilvusVectorID(ctx context.Context, printID uuid.UUID, vectorID string) error {
	query := `UPDATE afis_fingerprints SET milvus_vector_id = $1 WHERE print_id = $2`
	_, err := r.pool.Exec(ctx, query, vectorID, printID)
	return err
}
