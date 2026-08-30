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

func TestCreateMember(t *testing.T) {
	gangRepo := repository.NewInMemoryGangRepo()
	memberRepo := repository.NewInMemoryMemberRepo()
	gangSvc := service.NewGangService(gangRepo)
	memberSvc := service.NewMemberService(memberRepo, gangRepo)

	org, _ := gangSvc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "Test Gang", PrimaryActivity: domain.ActivityMixed, PrimaryDeptCode: "OU",
	}, uuid.Nil)

	member, err := memberSvc.CreateMember(context.Background(), domain.CreateMemberRequest{
		GangID:   org.GangID,
		FullName: "John Doe",
		Role:     ptr("Chef"),
		IsLeader: true,
	}, uuid.Nil)
	require.NoError(t, err)
	assert.Equal(t, "John Doe", member.FullName)
	assert.True(t, member.IsLeader)
	assert.Contains(t, member.NationalMemberID, "GANG-M-")
}

func TestGetMembers(t *testing.T) {
	gangRepo := repository.NewInMemoryGangRepo()
	memberRepo := repository.NewInMemoryMemberRepo()
	gangSvc := service.NewGangService(gangRepo)
	memberSvc := service.NewMemberService(memberRepo, gangRepo)

	org, _ := gangSvc.CreateOrganization(context.Background(), domain.CreateOrganizationRequest{
		Name: "Test Gang", PrimaryActivity: domain.ActivityMixed, PrimaryDeptCode: "OU",
	}, uuid.Nil)

	memberSvc.CreateMember(context.Background(), domain.CreateMemberRequest{
		GangID: org.GangID, FullName: "Member A",
	}, uuid.Nil)
	memberSvc.CreateMember(context.Background(), domain.CreateMemberRequest{
		GangID: org.GangID, FullName: "Member B",
	}, uuid.Nil)

	members, err := memberSvc.GetMembers(context.Background(), org.GangID)
	require.NoError(t, err)
	assert.Len(t, members, 2)
}
