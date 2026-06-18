package unit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/afis/internal/domain"
	"github.com/snisid/platform/services/afis/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_HitAbove85Percent(t *testing.T) {
	quality := service.NewQualityService(60)
	enrollment := service.NewEnrollmentService(quality)
	searchSvc := service.NewSearchService()

	officerID := uuid.New()

	enrollReq := domain.EnrollmentRequest{
		SubjectType:   domain.SubjectCriminal,
		EnrollingUnit: "DCPJ-PAP",
		Fingerprints: []domain.FingerprintCapture{
			{Position: domain.FingerRightThumb, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,REF1", NFIQ2Score: 85},
			{Position: domain.FingerRightIndex, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,REF2", NFIQ2Score: 82},
			{Position: domain.FingerLeftThumb, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,REF3", NFIQ2Score: 88},
			{Position: domain.FingerLeftIndex, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,REF4", NFIQ2Score: 80},
		},
	}
	subject, fps, err := enrollment.Enroll(context.Background(), enrollReq, officerID)
	require.NoError(t, err)
	searchSvc.IndexSubject(subject, fps)

	probeReq := domain.EnrollmentRequest{
		SubjectType:   domain.SubjectSuspect,
		EnrollingUnit: "BRH",
		Fingerprints: []domain.FingerprintCapture{
			{Position: domain.FingerRightThumb, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,PROBE1", NFIQ2Score: 90},
			{Position: domain.FingerRightIndex, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,PROBE2", NFIQ2Score: 87},
			{Position: domain.FingerLeftThumb, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,PROBE3", NFIQ2Score: 91},
			{Position: domain.FingerLeftIndex, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,PROBE4", NFIQ2Score: 84},
		},
	}

	results, err := searchSvc.SearchTenprint(context.Background(), probeReq)
	require.NoError(t, err)

	assert.NotEmpty(t, results)
	for _, r := range results {
		assert.GreaterOrEqual(t, r.Score, 0.85)
		assert.NotEmpty(t, r.Rank)
		assert.NotNil(t, r.CandidateID)
		assert.NotNil(t, r.SubjectID)
	}
}

func TestSearch_NoMatch_ReturnsEmpty(t *testing.T) {
	searchSvc := service.NewSearchService()

	probeReq := domain.EnrollmentRequest{
		SubjectType:   domain.SubjectSuspect,
		EnrollingUnit: "TEST",
		Fingerprints: []domain.FingerprintCapture{
			{Position: domain.FingerRightThumb, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,NODATA", NFIQ2Score: 90},
			{Position: domain.FingerRightIndex, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,NODATA", NFIQ2Score: 87},
		},
	}

	results, err := searchSvc.SearchTenprint(context.Background(), probeReq)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestSearch_LatentToTen(t *testing.T) {
	quality := service.NewQualityService(60)
	enrollment := service.NewEnrollmentService(quality)
	searchSvc := service.NewSearchService()

	officerID := uuid.New()
	enrollReq := domain.EnrollmentRequest{
		SubjectType:   domain.SubjectCriminal,
		EnrollingUnit: "DCPJ",
		Fingerprints: []domain.FingerprintCapture{
			{Position: domain.FingerRightThumb, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,REF1", NFIQ2Score: 85},
			{Position: domain.FingerRightIndex, Method: domain.CaptureLivescanner, ImageBase64: "data:image/wsq;base64,REF2", NFIQ2Score: 82},
		},
	}
	subject, fps, err := enrollment.Enroll(context.Background(), enrollReq, officerID)
	require.NoError(t, err)
	searchSvc.IndexSubject(subject, fps)

	latent := domain.LatentPrint{
		CaseReference:  "CS-2024-00123",
		ImageRef:       "s3://crime-scene/latent001.wsq",
		FingerPosition: domain.FingerRightIndex,
	}

	results, err := searchSvc.SearchLatent(context.Background(), latent)
	require.NoError(t, err)
	assert.NotEmpty(t, results)
}
