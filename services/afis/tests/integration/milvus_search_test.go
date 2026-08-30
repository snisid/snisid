package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/afis/internal/domain"
	"github.com/snisid/platform/services/afis/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMilvus_SearchLatency_Sub500ms(t *testing.T) {
	if testing.Short() {
		t.Skip("saut du test d'intégration Milvus en mode court")
	}

	quality := service.NewQualityService(60)
	enrollment := service.NewEnrollmentService(quality)
	searchSvc := service.NewSearchService()

	officerID := uuid.New()

	for i := 0; i < 100; i++ {
		req := domain.EnrollmentRequest{
			SubjectType:   domain.SubjectCriminal,
			EnrollingUnit: "DCPJ-MASS",
			Fingerprints: []domain.FingerprintCapture{
				{Position: domain.FingerRightThumb, Method: domain.CaptureLivescanner, ImageBase64: "dummy", NFIQ2Score: 80},
				{Position: domain.FingerRightIndex, Method: domain.CaptureLivescanner, ImageBase64: "dummy", NFIQ2Score: 80},
			},
		}
		subject, fps, err := enrollment.Enroll(context.Background(), req, officerID)
		require.NoError(t, err)
		searchSvc.IndexSubject(subject, fps)
	}

	probeReq := domain.EnrollmentRequest{
		SubjectType:   domain.SubjectSuspect,
		EnrollingUnit: "TEST-LATENCY",
		Fingerprints: []domain.FingerprintCapture{
			{Position: domain.FingerRightThumb, Method: domain.CaptureLivescanner, ImageBase64: "probe1", NFIQ2Score: 85},
			{Position: domain.FingerRightIndex, Method: domain.CaptureLivescanner, ImageBase64: "probe2", NFIQ2Score: 82},
		},
	}

	start := time.Now()
	results, err := searchSvc.SearchTenprint(context.Background(), probeReq)
	duration := time.Since(start)

	require.NoError(t, err)
	t.Logf("Recherche terminée en %v avec %d résultats", duration, len(results))
	assert.Less(t, duration.Milliseconds(), int64(500), "la recherche doit prendre moins de 500ms")
}

func TestMilvus_SearchNearest_TopKResults(t *testing.T) {
	if testing.Short() {
		t.Skip("saut du test d'intégration Milvus en mode court")
	}

	quality := service.NewQualityService(60)
	enrollment := service.NewEnrollmentService(quality)
	searchSvc := service.NewSearchService()

	officerID := uuid.New()
	for i := 0; i < 50; i++ {
		req := domain.EnrollmentRequest{
			SubjectType:   domain.SubjectCriminal,
			EnrollingUnit: "DCPJ-BULK",
			Fingerprints: []domain.FingerprintCapture{
				{Position: domain.FingerRightThumb, Method: domain.CaptureLivescanner, ImageBase64: "bulk", NFIQ2Score: 75},
				{Position: domain.FingerRightIndex, Method: domain.CaptureLivescanner, ImageBase64: "bulk", NFIQ2Score: 75},
			},
		}
		subject, fps, err := enrollment.Enroll(context.Background(), req, officerID)
		require.NoError(t, err)
		searchSvc.IndexSubject(subject, fps)
	}

	probeReq := domain.EnrollmentRequest{
		SubjectType:   domain.SubjectSuspect,
		EnrollingUnit: "TEST-TOPK",
		Fingerprints: []domain.FingerprintCapture{
			{Position: domain.FingerRightThumb, Method: domain.CaptureLivescanner, ImageBase64: "probe", NFIQ2Score: 88},
			{Position: domain.FingerRightIndex, Method: domain.CaptureLivescanner, ImageBase64: "probe", NFIQ2Score: 85},
		},
	}

	results, err := searchSvc.SearchTenprint(context.Background(), probeReq)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(results), service.CandidateListSize)
}
