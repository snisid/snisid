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

func TestCreateOrganization(t *testing.T) {
	repo := repository.NewInMemoryGangRepo()
	svc := service.NewGangService(repo)

	req := domain.CreateOrganizationRequest{
		Name:            "Viv Ansanm",
		PrimaryActivity: domain.ActivityMixed,
		PrimaryDeptCode: "OU",
		StructureType:   ptr(domain.GangStructureCoalition),
		OFACDesignation: true,
	}

	org, err := svc.CreateOrganization(context.Background(), req, uuid.Nil)
	require.NoError(t, err)
	assert.NotEmpty(t, org.GangID)
	assert.Equal(t, "Viv Ansanm", org.Name)
	assert.Equal(t, domain.GangActivityHigh, org.ActivityLevel)
	assert.True(t, org.IsActive)
	assert.Contains(t, org.NationalGangID, "GANG-HT-")
}

func TestGetOrganization(t *testing.T) {
	repo := repository.NewInMemoryGangRepo()
	svc := service.NewGangService(repo)

	req := domain.CreateOrganizationRequest{
		Name:            "G9 an Fanm",
		PrimaryActivity: domain.ActivityExtortion,
		PrimaryDeptCode: "OU",
	}
	org, err := svc.CreateOrganization(context.Background(), req, uuid.Nil)
	require.NoError(t, err)

	fetched, err := svc.GetOrganization(context.Background(), org.GangID)
	require.NoError(t, err)
	assert.Equal(t, org.GangID, fetched.GangID)
}

func TestListOrganizations(t *testing.T) {
	repo := repository.NewInMemoryGangRepo()
	svc := service.NewGangService(repo)

	svc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "Gang A", PrimaryActivity: domain.ActivityKidnapping, PrimaryDeptCode: "OU",
	}, uuid.Nil)
	svc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "Gang B", PrimaryActivity: domain.ActivityDrugTrafficking, PrimaryDeptCode: "AR",
	}, uuid.Nil)

	orgs, err := svc.ListOrganizations(context.Background())
	require.NoError(t, err)
	assert.Len(t, orgs, 2)
}

func TestByDeptCode(t *testing.T) {
	repo := repository.NewInMemoryGangRepo()
	svc := service.NewGangService(repo)

	svc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "400 Mawozo", PrimaryActivity: domain.ActivityKidnapping, PrimaryDeptCode: "OU",
	}, uuid.Nil)
	svc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "Gran Grif", PrimaryActivity: domain.ActivityExtortion, PrimaryDeptCode: "AR",
	}, uuid.Nil)

	orgs, err := svc.ByDeptCode(context.Background(), "OU")
	require.NoError(t, err)
	assert.Len(t, orgs, 1)
	assert.Equal(t, "400 Mawozo", orgs[0].Name)
}

func TestSanctioned(t *testing.T) {
	repo := repository.NewInMemoryGangRepo()
	svc := service.NewGangService(repo)

	svc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "400 Mawozo", PrimaryActivity: domain.ActivityKidnapping, PrimaryDeptCode: "OU",
		OFACDesignation: true,
	}, uuid.Nil)
	svc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "Gran Grif", PrimaryActivity: domain.ActivityExtortion, PrimaryDeptCode: "AR",
	}, uuid.Nil)

	orgs, err := svc.Sanctioned(context.Background())
	require.NoError(t, err)
	assert.Len(t, orgs, 1)
	assert.Equal(t, "400 Mawozo", orgs[0].Name)
}

func ptr[T any](v T) *T {
	return &v
}
