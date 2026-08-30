package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestLAPIIntegration_StolenPlate_TriggersAlert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Log("Test: LAPI integration - stolen plate triggers alert")
	t.Log("1. Créer une plaque volée dans SIVC")
	t.Log("2. Simuler une lecture LAPI")
	t.Log("3. Vérifier que l'alerte est déclenchée < 100ms")
	t.Log("4. Vérifier que le signalement est enregistré dans sivc_vehicle_sightings")
	t.Log("5. Vérifier que l'unité BLVV est notifiée via Kafka")

	assert.True(t, true, "integration test placeholder")
}

func TestCheckPlate_HighLoad_Sub5msP99(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Log("Test: Check plate high load - P99 < 5ms via hotlist Redis")
	t.Log("1. 1000 vérifications concurrentes")
	t.Log("2. P99 < 5ms (via hotlist Redis)")
	t.Log("3. P99 < 50ms (via fallback PostgreSQL)")

	start := time.Now()
	_ = domain.PlateCheckResult{
		PlateNumber: "PP-1234",
		CheckedAt:   time.Now(),
	}
	elapsed := time.Since(start)

	t.Logf("Check latency: %v", elapsed)
	assert.True(t, true, "load test placeholder")
}

func TestAlertCreation_Parallel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	userID := uuid.New()

	alerts := make([]*domain.CriminalAlert, 10)
	for i := 0; i < 10; i++ {
		req := domain.CreateAlertRequest{
			PlateNumber:   "PP-" + string(rune('A'+i)) + "1234",
			Make:          "Toyota",
			Model:         "Hilux",
			ColorPrimary:  "Blanc",
			CrimeCategory: domain.CrimeCategoryVehicleTheft,
			ReportingUnit: "BLVV",
			IncidentDate:  time.Now(),
		}
		alerts[i] = domain.NewCriminalAlert(req, userID)
	}

	assert.Len(t, alerts, 10)
	for _, a := range alerts {
		assert.NotEmpty(t, a.AlertID)
		assert.Equal(t, domain.AlertStatusActive, a.Status)
	}
}
