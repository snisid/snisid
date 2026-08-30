package unit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/snisid/platform/services/gang/internal/domain"
	"github.com/snisid/platform/services/gang/internal/repository"
	"github.com/snisid/platform/services/gang/internal/service"
)

func TestCreateTerritory(t *testing.T) {
	gangRepo := repository.NewInMemoryGangRepo()
	territoryRepo := repository.NewInMemoryTerritoryRepo()
	gangSvc := service.NewGangService(gangRepo)
	territorySvc := service.NewTerritoryService(territoryRepo, gangRepo)

	org, _ := gangSvc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "Test Gang", PrimaryActivity: domain.ActivityMixed, PrimaryDeptCode: "OU",
	}, uuid.Nil)

	territory, err := territorySvc.CreateTerritory(context.Background(), domain.CreateTerritoryRequest{
		GangID:    org.GangID,
		DeptCode:  "OU",
		Commune:   "Port-au-Prince",
		IsClaimed: ptr(true),
	}, uuid.Nil)
	require.NoError(t, err)
	assert.Equal(t, "OU", territory.DeptCode)
	assert.Equal(t, "Port-au-Prince", territory.Commune)
	assert.True(t, territory.IsClaimed)
}

func TestGetTerritories(t *testing.T) {
	gangRepo := repository.NewInMemoryGangRepo()
	territoryRepo := repository.NewInMemoryTerritoryRepo()
	gangSvc := service.NewGangService(gangRepo)
	territorySvc := service.NewTerritoryService(territoryRepo, gangRepo)

	org, _ := gangSvc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "Test Gang", PrimaryActivity: domain.ActivityMixed, PrimaryDeptCode: "OU",
	}, uuid.Nil)

	territorySvc.CreateTerritory(context.Background(), domain.CreateTerritoryRequest{
		GangID: org.GangID, DeptCode: "OU", Commune: "Cite Soleil",
	}, uuid.Nil)

	territories, err := territorySvc.GetTerritories(context.Background(), org.GangID)
	require.NoError(t, err)
	assert.Len(t, territories, 1)
}
