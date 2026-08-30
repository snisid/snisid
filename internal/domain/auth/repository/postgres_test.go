package repository

import (
	"context"
	"testing"
	"time"

	"github.com/snisid/platform/internal/domain/auth/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&entity.UserCredentials{}, &entity.WebAuthnCredential{})
	require.NoError(t, err)
	return db
}

func TestPostgresAuthRepository_CreateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAuthRepository(db)

	user := &entity.UserCredentials{
		UserID:       "u1",
		Username:     "jdoe",
		PasswordHash: "$argon2id$hash",
		Roles:        "user",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	var saved entity.UserCredentials
	err = db.First(&saved, "user_id = ?", "u1").Error
	require.NoError(t, err)
	assert.Equal(t, "jdoe", saved.Username)
	assert.Equal(t, "user", saved.Roles)
}

func TestPostgresAuthRepository_CreateUser_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAuthRepository(db)

	user := &entity.UserCredentials{
		UserID:   "u1",
		Username: "duplicate",
		Roles:    "user",
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	err = repo.CreateUser(context.Background(), user)
	assert.Error(t, err)
}

func TestPostgresAuthRepository_GetUserByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAuthRepository(db)

	user := &entity.UserCredentials{
		UserID:   "u1",
		Username: "jdoe",
		Roles:    "admin",
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	found, err := repo.GetUserByUsername(context.Background(), "jdoe")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "u1", found.UserID)
	assert.Equal(t, "admin", found.Roles)
}

func TestPostgresAuthRepository_GetUserByUsername_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAuthRepository(db)

	_, err := repo.GetUserByUsername(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestPostgresAuthRepository_UpdateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAuthRepository(db)

	user := &entity.UserCredentials{
		UserID:   "u1",
		Username: "jdoe",
		Roles:    "user",
	}
	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	user.Roles = "admin,auditor"
	err = repo.UpdateUser(context.Background(), user)
	require.NoError(t, err)

	var saved entity.UserCredentials
	err = db.First(&saved, "user_id = ?", "u1").Error
	require.NoError(t, err)
	assert.Equal(t, "admin,auditor", saved.Roles)
}

func TestPostgresAuthRepository_RegisterWebAuthn(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAuthRepository(db)

	cred := &entity.WebAuthnCredential{
		ID:              []byte("credential-id"),
		UserID:          "u1",
		PublicKey:       []byte("public-key-bytes"),
		AttestationType: "none",
		SignCount:       1,
		CreatedAt:       time.Now().UTC(),
	}

	err := repo.RegisterWebAuthn(context.Background(), cred)
	require.NoError(t, err)

	var saved entity.WebAuthnCredential
	err = db.First(&saved, "id = ?", []byte("credential-id")).Error
	require.NoError(t, err)
	assert.Equal(t, "u1", saved.UserID)
	assert.Equal(t, uint32(1), saved.SignCount)
}

func TestPostgresAuthRepository_GetWebAuthnCredentials(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAuthRepository(db)

	for i := 0; i < 3; i++ {
		cred := &entity.WebAuthnCredential{
			ID:        []byte{byte(i)},
			UserID:    "u1",
			PublicKey: []byte("key"),
		}
		err := repo.RegisterWebAuthn(context.Background(), cred)
		require.NoError(t, err)
	}

	creds, err := repo.GetWebAuthnCredentials(context.Background(), "u1")
	require.NoError(t, err)
	assert.Len(t, creds, 3)
}

func TestPostgresAuthRepository_GetWebAuthnCredentials_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAuthRepository(db)

	creds, err := repo.GetWebAuthnCredentials(context.Background(), "nonexistent")
	require.NoError(t, err)
	assert.Empty(t, creds)
}
