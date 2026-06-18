package unit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/fir/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRecord_Success(t *testing.T) {
	svc := service.NewRecordService()

	personID := uuid.New()
	record, err := svc.Create(context.Background(), personID, true, nil)

	require.NoError(t, err)
	require.NotNil(t, record)
	assert.NotNil(t, record.RecordID)
	assert.Equal(t, personID, record.SNISIDPersonID)
	assert.True(t, record.IsHaitianNational)
	assert.True(t, record.IsActive)
	assert.False(t, record.IsExpunged)
	assert.Contains(t, record.NationalFIRID, "FIR-HT-")
}

func TestCreateRecord_DuplicatePerson_ReturnsError(t *testing.T) {
	svc := service.NewRecordService()

	personID := uuid.New()
	_, err := svc.Create(context.Background(), personID, true, nil)
	require.NoError(t, err)

	_, err = svc.Create(context.Background(), personID, true, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "casier existe déjà")
}

func TestRecordGetByID_Found(t *testing.T) {
	svc := service.NewRecordService()

	personID := uuid.New()
	created, _ := svc.Create(context.Background(), personID, false, nil)

	found, err := svc.GetByID(context.Background(), created.RecordID)
	require.NoError(t, err)
	assert.Equal(t, created.RecordID, found.RecordID)
}

func TestRecordGetByID_NotFound(t *testing.T) {
	svc := service.NewRecordService()

	_, err := svc.GetByID(context.Background(), uuid.New())
	assert.ErrorIs(t, err, service.ErrRecordNotFound)
}

func TestGetByPersonID_Found(t *testing.T) {
	svc := service.NewRecordService()

	personID := uuid.New()
	created, _ := svc.Create(context.Background(), personID, true, nil)

	found, err := svc.GetByPersonID(context.Background(), personID)
	require.NoError(t, err)
	assert.Equal(t, created.RecordID, found.RecordID)
}

func TestGetByPersonID_NotFound(t *testing.T) {
	svc := service.NewRecordService()

	_, err := svc.GetByPersonID(context.Background(), uuid.New())
	assert.ErrorIs(t, err, service.ErrRecordNotFound)
}

func TestExpungeRecord(t *testing.T) {
	svc := service.NewRecordService()

	personID := uuid.New()
	created, _ := svc.Create(context.Background(), personID, true, nil)

	expunged, err := svc.Expunge(context.Background(), created.RecordID)
	require.NoError(t, err)
	assert.True(t, expunged.IsExpunged)
	assert.False(t, expunged.IsActive)
}

func TestReactivateRecord(t *testing.T) {
	svc := service.NewRecordService()

	personID := uuid.New()
	created, _ := svc.Create(context.Background(), personID, true, nil)

	svc.Expunge(context.Background(), created.RecordID)
	reactivated, err := svc.Reactivate(context.Background(), created.RecordID)
	require.NoError(t, err)
	assert.False(t, reactivated.IsExpunged)
	assert.True(t, reactivated.IsActive)
}

func TestUpdateRecord(t *testing.T) {
	svc := service.NewRecordService()

	personID := uuid.New()
	created, _ := svc.Create(context.Background(), personID, true, nil)

	haitian := false
	afisID := uuid.New()
	updated, err := svc.Update(context.Background(), created.RecordID, &haitian, &afisID)
	require.NoError(t, err)
	assert.False(t, updated.IsHaitianNational)
	assert.Equal(t, afisID, *updated.AFISSubjectID)
}

func TestListRecords(t *testing.T) {
	svc := service.NewRecordService()

	svc.Create(context.Background(), uuid.New(), true, nil)
	svc.Create(context.Background(), uuid.New(), false, nil)

	records, err := svc.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, records, 2)
}

func TestSearchByFIRID(t *testing.T) {
	svc := service.NewRecordService()

	created, _ := svc.Create(context.Background(), uuid.New(), true, nil)

	results, err := svc.Search(context.Background(), created.NationalFIRID)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, created.RecordID, results[0].RecordID)
}
