package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/afis-svc/internal/domain"
)

type FingerprintRepo interface {
	Create(ctx context.Context, fp *domain.Fingerprint) error
	GetBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]domain.Fingerprint, error)
	UpdateMilvusVectorID(ctx context.Context, printID uuid.UUID, vectorID string) error
}

type SubjectRepo interface {
	Create(ctx context.Context, s *domain.SubjectProfile) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.SubjectProfile, error)
	GetBySNISIDPersonID(ctx context.Context, personID uuid.UUID) (*domain.SubjectProfile, error)
}

type ImageRepo interface {
	Upload(ctx context.Context, objectName string, data []byte, contentType string) error
	Download(ctx context.Context, objectName string) ([]byte, error)
	Delete(ctx context.Context, objectName string) error
}

type VectorRepo interface {
	InsertVectors(ctx context.Context, printIDs []string, vectors [][]float32) error
	SearchNearest(ctx context.Context, queryVectors [][]float32, topK int) ([]interface{}, error)
}

type LatentRepo interface {
	Create(ctx context.Context, lp *domain.LatentPrint) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.LatentPrint, error)
	ConfirmMatch(ctx context.Context, latentID, subjectID uuid.UUID, score float64, examiner uuid.UUID) error
	GetUnidentified(ctx context.Context) ([]domain.LatentPrint, error)
}
