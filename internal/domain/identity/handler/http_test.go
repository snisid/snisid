package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/internal/domain/identity/entity"
	"github.com/snisid/platform/internal/domain/identity/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockIdentityService struct {
	createIdentityFn func(ctx context.Context, ident *entity.Identity, changedBy string) (*entity.Identity, error)
	updateIdentityFn func(ctx context.Context, id string, updateFn func(*entity.Identity), reason, changedBy string) (*entity.Identity, error)
	getIdentityFn    func(ctx context.Context, id string) (*entity.Identity, error)
	flagIdentityFn   func(ctx context.Context, id, reason, changedBy string) error
	getHistoryFn     func(ctx context.Context, id string) ([]entity.IdentityHistory, error)
}

func (m *mockIdentityService) CreateIdentity(ctx context.Context, ident *entity.Identity, changedBy string) (*entity.Identity, error) {
	if m.createIdentityFn != nil {
		return m.createIdentityFn(ctx, ident, changedBy)
	}
	ident.ID = "ID-123"
	ident.Status = entity.StatePending
	return ident, nil
}

func (m *mockIdentityService) UpdateIdentity(ctx context.Context, id string, updateFn func(*entity.Identity), reason, changedBy string) (*entity.Identity, error) {
	if m.updateIdentityFn != nil {
		return m.updateIdentityFn(ctx, id, updateFn, reason, changedBy)
	}
	return &entity.Identity{ID: id, FirstName: "Updated", Status: entity.StateActive}, nil
}

func (m *mockIdentityService) GetIdentity(ctx context.Context, id string) (*entity.Identity, error) {
	if m.getIdentityFn != nil {
		return m.getIdentityFn(ctx, id)
	}
	return &entity.Identity{ID: id, FirstName: "John", LastName: "Doe", Status: entity.StateActive}, nil
}

func (m *mockIdentityService) FlagIdentity(ctx context.Context, id, reason, changedBy string) error {
	if m.flagIdentityFn != nil {
		return m.flagIdentityFn(ctx, id, reason, changedBy)
	}
	return nil
}

func (m *mockIdentityService) GetHistory(ctx context.Context, id string) ([]entity.IdentityHistory, error) {
	if m.getHistoryFn != nil {
		return m.getHistoryFn(ctx, id)
	}
	return []entity.IdentityHistory{}, nil
}

func setupIdentityRouter(svc usecase.IdentityService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHttpHandler(svc)
	h.RegisterRoutes(r.Group("/v1"))
	return r
}

func TestCreateIdentity_Success(t *testing.T) {
	svc := &mockIdentityService{}
	router := setupIdentityRouter(svc)

	body := `{"firstName":"John","lastName":"Doe","dob":"1990-01-01","gender":"M"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/identities", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp entity.Identity
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "John", resp.FirstName)
	assert.Equal(t, "Doe", resp.LastName)
}

func TestCreateIdentity_InvalidBody(t *testing.T) {
	svc := &mockIdentityService{}
	router := setupIdentityRouter(svc)

	body := `{invalid}`
	req := httptest.NewRequest(http.MethodPost, "/v1/identities", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateIdentity_ServiceError(t *testing.T) {
	svc := &mockIdentityService{
		createIdentityFn: func(ctx context.Context, ident *entity.Identity, changedBy string) (*entity.Identity, error) {
			return nil, assert.AnError
		},
	}
	router := setupIdentityRouter(svc)

	body := `{"firstName":"John"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/identities", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetIdentity_Success(t *testing.T) {
	svc := &mockIdentityService{}
	router := setupIdentityRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/identities/ID-123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp entity.Identity
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ID-123", resp.ID)
	assert.Equal(t, "John", resp.FirstName)
}

func TestGetIdentity_NotFound(t *testing.T) {
	svc := &mockIdentityService{
		getIdentityFn: func(ctx context.Context, id string) (*entity.Identity, error) {
			return nil, assert.AnError
		},
	}
	router := setupIdentityRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/identities/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateIdentity_Success(t *testing.T) {
	svc := &mockIdentityService{}
	router := setupIdentityRouter(svc)

	body := `{"firstName":"Jane","lastName":"Smith"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/identities/ID-123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp entity.Identity
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ID-123", resp.ID)
}

func TestUpdateIdentity_InvalidBody(t *testing.T) {
	svc := &mockIdentityService{}
	router := setupIdentityRouter(svc)

	req := httptest.NewRequest(http.MethodPut, "/v1/identities/ID-123", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFlagIdentity_Success(t *testing.T) {
	svc := &mockIdentityService{}
	router := setupIdentityRouter(svc)

	body := `{"reason":"Suspicious activity detected"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/identities/ID-123/flag", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "flagged", resp["status"])
}

func TestFlagIdentity_Error(t *testing.T) {
	svc := &mockIdentityService{
		flagIdentityFn: func(ctx context.Context, id, reason, changedBy string) error {
			return assert.AnError
		},
	}
	router := setupIdentityRouter(svc)

	body := `{"reason":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/identities/ID-123/flag", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetHistory_Success(t *testing.T) {
	svc := &mockIdentityService{}
	router := setupIdentityRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/identities/ID-123/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetHistory_Error(t *testing.T) {
	svc := &mockIdentityService{
		getHistoryFn: func(ctx context.Context, id string) ([]entity.IdentityHistory, error) {
			return nil, assert.AnError
		},
	}
	router := setupIdentityRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/identities/ID-123/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
