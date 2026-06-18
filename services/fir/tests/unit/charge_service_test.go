package unit

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/fir/internal/domain"
	"github.com/snisid/platform/services/fir/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestData(t *testing.T) (*service.ChargeService, uuid.UUID) {
	t.Helper()
	cs := service.NewChargeService()
	recordID := uuid.New()
	return cs, recordID
}

func TestCreateArrest_Success(t *testing.T) {
	cs, recordID := setupTestData(t)

	now := time.Now()
	charge := domain.Charge{
		ArrestDate:     &now,
		ArrestingUnit:  strPtr("DCPJ-PAP"),
		ArrestLocation: strPtr("Port-au-Prince"),
		DeptCode:       strPtr("OU"),
		ChargesText:    strPtr("Vol à main armée"),
		OffenseClass:   domain.OffenseCrime,
		CaseReference:  strPtr("PARQUET-2024-00123"),
	}

	created, err := cs.CreateArrest(context.Background(), recordID, charge)
	require.NoError(t, err)
	assert.NotNil(t, created.ChargeID)
	assert.True(t, created.IsArrest)
	assert.Equal(t, domain.CaseStatusOpen, created.CaseStatus)
	assert.Equal(t, domain.OffenseCrime, created.OffenseClass)
}

func TestCreateConviction_Success(t *testing.T) {
	cs, recordID := setupTestData(t)

	now := time.Now()
	sentence := domain.SentencePrison
	days := 365
	charge := domain.Charge{
		CourtName:            strPtr("Tribunal Correctionnel Port-au-Prince"),
		CourtDept:            strPtr("OU"),
		OffenseDescription:   strPtr("Vol qualifié"),
		OffenseClass:         domain.OffenseDelit,
		CaseReference:        strPtr("TRIB-2024-00456"),
		VerdictDate:          &now,
		CaseStatus:           domain.CaseStatusConvicted,
		SentenceType:         &sentence,
		SentenceDurationDays: &days,
		JudgeName:            strPtr("M. Joseph"),
	}

	created, err := cs.CreateConviction(context.Background(), recordID, charge)
	require.NoError(t, err)
	assert.NotNil(t, created.ChargeID)
	assert.False(t, created.IsArrest)
	assert.Equal(t, domain.CaseStatusConvicted, created.CaseStatus)
}

func TestListByRecord(t *testing.T) {
	cs, recordID := setupTestData(t)

	charge := domain.Charge{
		OffenseClass: domain.OffenseDelit,
		CaseStatus:   domain.CaseStatusOpen,
	}
	cs.CreateArrest(context.Background(), recordID, charge)
	cs.CreateArrest(context.Background(), recordID, charge)

	charges, err := cs.ListByRecord(context.Background(), recordID)
	require.NoError(t, err)
	assert.Len(t, charges, 2)
}

func TestChargeGetByID_NotFound(t *testing.T) {
	cs, _ := setupTestData(t)

	_, err := cs.GetByID(context.Background(), uuid.New())
	assert.ErrorIs(t, err, service.ErrChargeNotFound)
}

func TestUpdateStatus(t *testing.T) {
	cs, recordID := setupTestData(t)

	charge := domain.Charge{
		OffenseClass: domain.OffenseCrime,
		CaseStatus:   domain.CaseStatusOpen,
	}
	created, _ := cs.CreateArrest(context.Background(), recordID, charge)

	updated, err := cs.UpdateStatus(context.Background(), created.ChargeID, domain.CaseStatusConvicted)
	require.NoError(t, err)
	assert.Equal(t, domain.CaseStatusConvicted, updated.CaseStatus)
}

func TestDeleteCharge(t *testing.T) {
	cs, recordID := setupTestData(t)

	charge := domain.Charge{
		OffenseClass: domain.OffenseDelit,
		CaseStatus:   domain.CaseStatusOpen,
	}
	created, _ := cs.CreateArrest(context.Background(), recordID, charge)

	err := cs.Delete(context.Background(), created.ChargeID)
	require.NoError(t, err)

	_, err = cs.GetByID(context.Background(), created.ChargeID)
	assert.ErrorIs(t, err, service.ErrChargeNotFound)
}

func strPtr(s string) *string { return &s }
