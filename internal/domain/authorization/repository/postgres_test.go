package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/snisid/platform/internal/domain/authorization/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}
	return gormDB, mock
}

func TestGetActivePolicies_Success(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewPostgresPolicyRepository(gormDB)

	rows := sqlmock.NewRows([]string{"id", "name", "module", "enabled", "version"}).
		AddRow("pol-001", "abac_rules", `package test`, true, 1).
		AddRow("pol-002", "soc_rules", `package test`, true, 2)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "policies" WHERE enabled = $1`)).
		WithArgs(true).
		WillReturnRows(rows)

	policies, err := repo.GetActivePolicies(context.Background())
	if err != nil {
		t.Fatalf("GetActivePolicies failed: %v", err)
	}
	if len(policies) != 2 {
		t.Errorf("policies count = %d, want 2", len(policies))
	}
	if policies[0].Name != "abac_rules" {
		t.Errorf("Name = %s, want abac_rules", policies[0].Name)
	}
}

func TestGetActivePolicies_Empty(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewPostgresPolicyRepository(gormDB)

	rows := sqlmock.NewRows([]string{"id", "name", "module", "enabled", "version"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "policies" WHERE enabled = $1`)).
		WithArgs(true).
		WillReturnRows(rows)

	policies, err := repo.GetActivePolicies(context.Background())
	if err != nil {
		t.Fatalf("GetActivePolicies failed: %v", err)
	}
	if len(policies) != 0 {
		t.Errorf("policies count = %d, want 0", len(policies))
	}
}

func TestGetActivePolicies_DBError(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewPostgresPolicyRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "policies" WHERE enabled = $1`)).
		WithArgs(true).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.GetActivePolicies(context.Background())
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetRoleGrants_Success(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewPostgresPolicyRepository(gormDB)

	rows := sqlmock.NewRows([]string{"id", "role", "action", "resource"}).
		AddRow("grant-001", "admin", "write", "identities:*").
		AddRow("grant-002", "officer", "read", "identities:NNU-*")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "role_grants"`)).
		WillReturnRows(rows)

	grants, err := repo.GetRoleGrants(context.Background())
	if err != nil {
		t.Fatalf("GetRoleGrants failed: %v", err)
	}
	if len(grants) != 2 {
		t.Errorf("grants count = %d, want 2", len(grants))
	}
}

func TestNewPostgresPolicyRepository(t *testing.T) {
	gormDB, _ := newMockDB(t)
	repo := NewPostgresPolicyRepository(gormDB)
	if repo == nil {
		t.Fatal("NewPostgresPolicyRepository returned nil")
	}
}
