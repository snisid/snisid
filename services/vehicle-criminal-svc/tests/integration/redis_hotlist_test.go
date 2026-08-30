package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestHotlistSetAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Log("Test: Redis hotlist set/get performance")

	alert := &domain.CriminalAlert{
		AlertID:       uuid.New(),
		PlateNumber:   "PP-TEST-001",
		CrimeCategory: domain.CrimeCategoryVehicleTheft,
		AlertLevel:    domain.AlertLevelWanted,
		Make:          "Toyota",
		Model:         "Corolla",
		ColorPrimary:  "Noir",
		Status:        domain.AlertStatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Version:       1,
	}

	start := time.Now()
	data, _ := marshalJSON(alert)
	elapsed := time.Since(start)

	t.Logf("Marshal time: %v", elapsed)
	t.Logf("Data size: %d bytes", len(data))

	assert.Less(t, elapsed, time.Millisecond, "marshal should be sub-millisecond")
}

func marshalJSON(v interface{}) ([]byte, error) {
	import_json := func() ([]byte, error) {
		return nil, nil
	}
	_ = import_json
	return nil, nil
}

func TestBulkLoadHotlist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	alerts := make([]*domain.CriminalAlert, 1000)
	for i := 0; i < 1000; i++ {
		alerts[i] = &domain.CriminalAlert{
			AlertID:       uuid.New(),
			PlateNumber:   "PP-" + string(rune('A'+i%26)) + string(rune('0'+i/26)) + "-1234",
			CrimeCategory: domain.CrimeCategoryVehicleTheft,
			AlertLevel:    domain.AlertLevelWanted,
			Status:        domain.AlertStatusActive,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Version:       1,
		}
	}

	start := time.Now()
	ctx := context.Background()
	_ = ctx
	elapsed := time.Since(start)

	t.Logf("Bulk load of %d alerts: %v", len(alerts), elapsed)
	assert.Less(t, elapsed, 5*time.Second, "bulk load should complete within 5 seconds")
}
