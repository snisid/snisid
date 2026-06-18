package unit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewStolenPlate(t *testing.T) {
	req := domain.DeclareStolenPlateRequest{
		PlateNumber:   "PP-5678",
		PlateCategory: domain.PlateCategoryPrivate,
		TheftDate:     time.Now(),
		TheftDeptCode: "OU",
		ReportingUnit: "BLVV",
	}

	userID := uuid.New()
	plate := domain.NewStolenPlate(req, userID)

	assert.NotEmpty(t, plate.PlateID)
	assert.Equal(t, "PP-5678", plate.PlateNumber)
	assert.Equal(t, domain.PlateCategoryPrivate, plate.PlateCategory)
	assert.Equal(t, domain.StolenPlateStatusStolen, plate.Status)
	assert.Equal(t, "BLVV", plate.ReportingUnit)
}

func TestDeclareStolenPlateRequest_StateClone(t *testing.T) {
	req := domain.DeclareStolenPlateRequest{
		PlateNumber:        "SE-00871",
		PlateCategory:      domain.PlateCategoryState,
		TheftDate:          time.Now(),
		TheftDeptCode:      "OU",
		ReportingUnit:      "DCPJ",
		IsStatePlateClone:  true,
		ImpersonatedAgency: strPtr("PNH"),
	}

	userID := uuid.New()
	plate := domain.NewStolenPlate(req, userID)

	assert.True(t, plate.IsStatePlateClone)
	assert.NotNil(t, plate.ImpersonatedAgency)
	assert.Equal(t, "PNH", *plate.ImpersonatedAgency)
}

func TestValidatePlateNumber_VariousFormats(t *testing.T) {
	validPlates := []string{
		"PP-1234",
		"SE-00871",
		"ABC-123456",
		"M-1234",
		"TC-5678",
		"PL-901234",
	}

	for _, plate := range validPlates {
		t.Run(plate, func(t *testing.T) {
			err := domain.ValidatePlateNumber(plate)
			assert.NoError(t, err, "plate %s should be valid", plate)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
