package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/snisid/platform/backend/internal/domain/identity/entity"
	"github.com/snisid/platform/backend/internal/platform/events"
)

type mockRepo struct {
	createFn  func(ctx context.Context, ident *entity.Identity, reason, changedBy string) error
	updateFn  func(ctx context.Context, ident *entity.Identity, reason, changedBy string) error
	getByIDFn func(ctx context.Context, id string) (*entity.Identity, error)
	historyFn func(ctx context.Context, id string) ([]entity.IdentityHistory, error)
	deleteFn  func(ctx context.Context, id, reason, changedBy string) error
}

func (m *mockRepo) Create(ctx context.Context, ident *entity.Identity, reason, changedBy string) error {
	if m.createFn != nil {
		return m.createFn(ctx, ident, reason, changedBy)
	}
	return nil
}

func (m *mockRepo) Update(ctx context.Context, ident *entity.Identity, reason, changedBy string) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, ident, reason, changedBy)
	}
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*entity.Identity, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) GetHistory(ctx context.Context, id string) ([]entity.IdentityHistory, error) {
	if m.historyFn != nil {
		return m.historyFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRepo) Delete(ctx context.Context, id, reason, changedBy string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id, reason, changedBy)
	}
	return nil
}

func newService(repo *mockRepo) IdentityService {
	return NewIdentityService(repo, nil)
}

func TestCreateIdentity_Success(t *testing.T) {
	var savedIdent *entity.Identity

	repo := &mockRepo{
		createFn: func(ctx context.Context, ident *entity.Identity, reason, changedBy string) error {
			savedIdent = ident
			return nil
		},
	}

	svc := newService(repo)
	ident := &entity.Identity{
		FirstName: "Marie",
		LastName:  "Pierre",
		DOB:       "1995-06-15",
		Gender:    "F",
		Agency:    "ONI",
	}

	result, err := svc.CreateIdentity(context.Background(), ident, "admin")
	if err != nil {
		t.Fatalf("CreateIdentity failed: %v", err)
	}

	if result.ID == "" {
		t.Error("Expected non-empty ID")
	}
	if result.Status != entity.StatePending {
		t.Errorf("Status = %s, want pending", result.Status)
	}
	if result.Version != 1 {
		t.Errorf("Version = %d, want 1", result.Version)
	}
	if savedIdent == nil {
		t.Fatal("Expected repo.Create to be called")
	}
}

func TestCreateIdentity_RepoError(t *testing.T) {
	repo := &mockRepo{
		createFn: func(ctx context.Context, ident *entity.Identity, reason, changedBy string) error {
			return errors.New("database error")
		},
	}

	svc := newService(repo)
	_, err := svc.CreateIdentity(context.Background(), &entity.Identity{}, "admin")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetIdentity_Success(t *testing.T) {
	expected := &entity.Identity{
		ID:        "ID-123",
		FirstName: "Jean",
		LastName:  "Dupont",
		Status:    entity.StateActive,
	}

	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, id string) (*entity.Identity, error) {
			if id != "ID-123" {
				return nil, errors.New("not found")
			}
			return expected, nil
		},
	}

	svc := newService(repo)
	result, err := svc.GetIdentity(context.Background(), "ID-123")
	if err != nil {
		t.Fatalf("GetIdentity failed: %v", err)
	}

	if result.FirstName != "Jean" {
		t.Errorf("FirstName = %s, want Jean", result.FirstName)
	}
}

func TestGetIdentity_NotFound(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, id string) (*entity.Identity, error) {
			return nil, errors.New("not found")
		},
	}

	svc := newService(repo)
	_, err := svc.GetIdentity(context.Background(), "ID-INVALID")
	if err == nil {
		t.Fatal("Expected error for non-existent identity")
	}
}

func TestUpdateIdentity_Success(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, id string) (*entity.Identity, error) {
			return &entity.Identity{
				ID:        id,
				FirstName: "Jean",
				LastName:  "Dupont",
				Status:    entity.StateActive,
				Version:   1,
			}, nil
		},
		updateFn: func(ctx context.Context, ident *entity.Identity, reason, changedBy string) error {
			return nil
		},
	}

	svc := newService(repo)
	result, err := svc.UpdateIdentity(context.Background(), "ID-123",
		func(i *entity.Identity) { i.FirstName = "Jean-Marc" },
		"Correction prénom", "admin")

	if err != nil {
		t.Fatalf("UpdateIdentity failed: %v", err)
	}
	if result.FirstName != "Jean-Marc" {
		t.Errorf("FirstName = %s, want Jean-Marc", result.FirstName)
	}
}

func TestFlagIdentity_SetsSuspended(t *testing.T) {
	var updatedIdent *entity.Identity

	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, id string) (*entity.Identity, error) {
			return &entity.Identity{
				ID:     id,
				Status: entity.StateActive,
			}, nil
		},
		updateFn: func(ctx context.Context, ident *entity.Identity, reason, changedBy string) error {
			updatedIdent = ident
			return nil
		},
	}

	svc := newService(repo)
	err := svc.FlagIdentity(context.Background(), "ID-123", "Fraude suspectée", "security-admin")
	if err != nil {
		t.Fatalf("FlagIdentity failed: %v", err)
	}

	if updatedIdent == nil {
		t.Fatal("Expected update to be called")
	}
	if updatedIdent.Status != entity.StateSuspended {
		t.Errorf("Status = %s, want suspended", updatedIdent.Status)
	}
}

func TestGetHistory_Success(t *testing.T) {
	expected := []entity.IdentityHistory{
		{HistoryID: "H-1", IdentityID: "ID-123", Version: 1},
		{HistoryID: "H-2", IdentityID: "ID-123", Version: 2},
	}

	repo := &mockRepo{
		historyFn: func(ctx context.Context, id string) ([]entity.IdentityHistory, error) {
			return expected, nil
		},
	}

	svc := newService(repo)
	result, err := svc.GetHistory(context.Background(), "ID-123")
	if err != nil {
		t.Fatalf("GetHistory failed: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("len(history) = %d, want 2", len(result))
	}
}

func TestNewIdentityService(t *testing.T) {
	svc := NewIdentityService(nil, nil)
	if svc == nil {
		t.Fatal("NewIdentityService returned nil")
	}
}
