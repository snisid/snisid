package unit

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/snisid/platform/services/gang/internal/domain"
	"github.com/snisid/platform/services/gang/internal/repository"
	"github.com/snisid/platform/services/gang/internal/service"
)

func TestCreateIncident(t *testing.T) {
	gangRepo := repository.NewInMemoryGangRepo()
	incidentRepo := repository.NewInMemoryIncidentRepo()
	gangSvc := service.NewGangService(gangRepo)
	incidentSvc := service.NewIncidentService(incidentRepo, gangRepo)

	org, _ := gangSvc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "Test Gang", PrimaryActivity: domain.ActivityMixed, PrimaryDeptCode: "OU",
	}, uuid.Nil)

	now := time.Now()
	inc, err := incidentSvc.CreateIncident(context.Background(), domain.CreateIncidentRequest{
		GangID:       org.GangID,
		IncidentType: "KIDNAPPING",
		IncidentDate: now,
		DeptCode:     ptr("OU"),
		Casualties:   ptr(int16(3)),
	}, uuid.Nil)
	require.NoError(t, err)
	assert.Equal(t, "KIDNAPPING", inc.IncidentType)
	assert.Equal(t, int16(3), inc.Casualties)
}

func TestGetIncidents(t *testing.T) {
	gangRepo := repository.NewInMemoryGangRepo()
	incidentRepo := repository.NewInMemoryIncidentRepo()
	gangSvc := service.NewGangService(gangRepo)
	incidentSvc := service.NewIncidentService(incidentRepo, gangRepo)

	org, _ := gangSvc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "Test Gang", PrimaryActivity: domain.ActivityMixed, PrimaryDeptCode: "OU",
	}, uuid.Nil)

	incidentSvc.CreateIncident(context.Background(), domain.CreateIncidentRequest{
		GangID: org.GangID, IncidentType: "KIDNAPPING", IncidentDate: time.Now(),
	}, uuid.Nil)
	incidentSvc.CreateIncident(context.Background(), domain.CreateIncidentRequest{
		GangID: org.GangID, IncidentType: "EXTORTION", IncidentDate: time.Now(),
	}, uuid.Nil)

	incidents, err := incidentSvc.GetIncidents(context.Background(), org.GangID)
	require.NoError(t, err)
	assert.Len(t, incidents, 2)
}
