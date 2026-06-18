package unit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/fir/internal/domain"
	"github.com/snisid/platform/services/fir/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddAlias_Success(t *testing.T) {
	as := service.NewAliasService()
	recordID := uuid.New()

	alias := domain.Alias{
		FirstName:  strPtr("Jean"),
		LastName:   strPtr("Dupont"),
		IDDocument: strPtr("ID-12345"),
	}

	created, err := as.Add(context.Background(), recordID, alias)
	require.NoError(t, err)
	assert.NotNil(t, created.AliasID)
	assert.Equal(t, recordID, created.RecordID)
	assert.Equal(t, "Jean", *created.FirstName)
	assert.Equal(t, "Dupont", *created.LastName)
}

func TestListAliases_ByRecord(t *testing.T) {
	as := service.NewAliasService()
	recordID := uuid.New()

	as.Add(context.Background(), recordID, domain.Alias{FirstName: strPtr("Alias1")})
	as.Add(context.Background(), recordID, domain.Alias{FirstName: strPtr("Alias2")})

	aliases, err := as.ListByRecord(context.Background(), recordID)
	require.NoError(t, err)
	assert.Len(t, aliases, 2)
}

func TestRemoveAlias_Success(t *testing.T) {
	as := service.NewAliasService()
	recordID := uuid.New()

	created, _ := as.Add(context.Background(), recordID, domain.Alias{FirstName: strPtr("Test")})

	err := as.Remove(context.Background(), created.AliasID)
	require.NoError(t, err)

	aliases, err := as.ListByRecord(context.Background(), recordID)
	require.NoError(t, err)
	assert.Empty(t, aliases)
}

func TestRemoveAlias_NotFound(t *testing.T) {
	as := service.NewAliasService()

	err := as.Remove(context.Background(), uuid.New())
	assert.ErrorIs(t, err, service.ErrAliasNotFound)
}

func TestAliasGetByID_Found(t *testing.T) {
	as := service.NewAliasService()
	recordID := uuid.New()

	created, _ := as.Add(context.Background(), recordID, domain.Alias{
		FirstName: strPtr("Marie"),
		LastName:  strPtr("Charles"),
	})

	found, err := as.GetByID(context.Background(), created.AliasID)
	require.NoError(t, err)
	assert.Equal(t, "Marie", *found.FirstName)
	assert.Equal(t, "Charles", *found.LastName)
}

func TestAliasGetByID_NotFound(t *testing.T) {
	as := service.NewAliasService()

	_, err := as.GetByID(context.Background(), uuid.New())
	assert.ErrorIs(t, err, service.ErrAliasNotFound)
}
