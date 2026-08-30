package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/internal/domain/identity/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupIdentityDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&entity.Identity{}, &entity.IdentityHistory{}, &entity.BiometricReference{}, &entity.DocumentAssociation{})
	require.NoError(t, err)
	return db
}

func TestPostgresRepository_Create(t *testing.T) {
	db := setupIdentityDB(t)
	repo := NewPostgresRepository(db)

	ident := &entity.Identity{
		ID:        uuid.NewString(),
		FirstName: "John",
		LastName:  "Doe",
		DOB:       "1990-01-01",
		Gender:    "M",
		Agency:    "ONI",
		Status:    entity.StatePending,
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := repo.Create(context.Background(), ident, "Initial creation", "system")
	require.NoError(t, err)

	var saved entity.Identity
	err = db.First(&saved, "id = ?", ident.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "John", saved.FirstName)
	assert.Equal(t, entity.StatePending, saved.Status)

	var history []entity.IdentityHistory
	err = db.Where("identity_id = ?", ident.ID).Find(&history).Error
	require.NoError(t, err)
	assert.Len(t, history, 1)
	assert.Equal(t, "Initial creation", history[0].Reason)
}

func TestPostgresRepository_GetByID(t *testing.T) {
	db := setupIdentityDB(t)
	repo := NewPostgresRepository(db)

	ident := &entity.Identity{
		ID:        "ID-1",
		FirstName: "Jane",
		LastName:  "Smith",
		Status:    entity.StateActive,
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err := repo.Create(context.Background(), ident, "test", "system")
	require.NoError(t, err)

	found, err := repo.GetByID(context.Background(), "ID-1")
	require.NoError(t, err)
	assert.Equal(t, "Jane", found.FirstName)
	assert.Equal(t, "Smith", found.LastName)
}

func TestPostgresRepository_GetByID_NotFound(t *testing.T) {
	db := setupIdentityDB(t)
	repo := NewPostgresRepository(db)

	_, err := repo.GetByID(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestPostgresRepository_Update(t *testing.T) {
	db := setupIdentityDB(t)
	repo := NewPostgresRepository(db)

	ident := &entity.Identity{
		ID:        "ID-1",
		FirstName: "John",
		Status:    entity.StatePending,
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err := repo.Create(context.Background(), ident, "create", "system")
	require.NoError(t, err)

	ident.FirstName = "Johnny"
	ident.Status = entity.StateActive
	err = repo.Update(context.Background(), ident, "name correction", "admin")
	require.NoError(t, err)

	var saved entity.Identity
	err = db.First(&saved, "id = ?", "ID-1").Error
	require.NoError(t, err)
	assert.Equal(t, "Johnny", saved.FirstName)
	assert.Equal(t, entity.StateActive, saved.Status)
	assert.Equal(t, 2, saved.Version)

	var history []entity.IdentityHistory
	err = db.Where("identity_id = ?", "ID-1").Order("version asc").Find(&history).Error
	require.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, "name correction", history[1].Reason)
}

func TestPostgresRepository_Delete(t *testing.T) {
	db := setupIdentityDB(t)
	repo := NewPostgresRepository(db)

	ident := &entity.Identity{
		ID:        "ID-1",
		FirstName: "John",
		Status:    entity.StateActive,
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err := repo.Create(context.Background(), ident, "create", "system")
	require.NoError(t, err)

	err = repo.Delete(context.Background(), "ID-1", "deceased", "registrar")
	require.NoError(t, err)

	var saved entity.Identity
	err = db.First(&saved, "id = ?", "ID-1").Error
	require.NoError(t, err)
	assert.Equal(t, entity.StateDeceased, saved.Status)
	assert.Equal(t, 2, saved.Version)

	var history []entity.IdentityHistory
	err = db.Where("identity_id = ?", "ID-1").Find(&history).Error
	require.NoError(t, err)
	assert.Len(t, history, 2)
}

func TestPostgresRepository_GetHistory(t *testing.T) {
	db := setupIdentityDB(t)
	repo := NewPostgresRepository(db)

	ident := &entity.Identity{
		ID:        "ID-1",
		FirstName: "John",
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err := repo.Create(context.Background(), ident, "v1", "system")
	require.NoError(t, err)

	ident.FirstName = "Johnny"
	err = repo.Update(context.Background(), ident, "v2", "admin")
	require.NoError(t, err)

	history, err := repo.GetHistory(context.Background(), "ID-1")
	require.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, "v2", history[0].Reason) // ordered by version desc
	assert.Equal(t, "v1", history[1].Reason)
}

func TestPostgresRepository_GetHistory_Empty(t *testing.T) {
	db := setupIdentityDB(t)
	repo := NewPostgresRepository(db)

	history, err := repo.GetHistory(context.Background(), "nonexistent")
	require.NoError(t, err)
	assert.Empty(t, history)
}
