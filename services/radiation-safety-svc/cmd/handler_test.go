package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/snisid/radiation-safety-svc/internal/domain"
	"github.com/snisid/radiation-safety-svc/internal/handler"
	"github.com/snisid/radiation-safety-svc/internal/kafka"
	"github.com/snisid/radiation-safety-svc/internal/repository"
	"github.com/snisid/radiation-safety-svc/internal/service"
)

type mockRadiationRepo struct {
	repository.RadiationRepository
}

func (m *mockRadiationRepo) GetUnrespondedAlerts(ctx context.Context) ([]domain.RadiationAlert, error) {
	return []domain.RadiationAlert{}, nil
}

func (m *mockRadiationRepo) GetSuspiciousChemicals(ctx context.Context) ([]domain.ChemicalPrecursor, error) {
	return []domain.ChemicalPrecursor{}, nil
}

func (m *mockRadiationRepo) GetDashboardStats(ctx context.Context) (*repository.DashboardStats, error) {
	return &repository.DashboardStats{}, nil
}

type mockRadKafkaProducer struct {
	kafka.Producer
}

func (m *mockRadKafkaProducer) Publish(ctx context.Context, key string, msg interface{}) error {
	return nil
}

func setupRadiationRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	repo := &mockRadiationRepo{}
	prod := &mockRadKafkaProducer{}
	svc := service.NewRadiationService(repo, prod)
	h := handler.NewRadiationHandler(svc)

	r := gin.New()
	h.RegisterRoutes(r)
	return r
}

func TestRadiation_GetUnrespondedAlerts(t *testing.T) {
	r := setupRadiationRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiation/alerts/unresponded", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var result []interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
}

func TestRadiation_GetSuspiciousChemicals(t *testing.T) {
	r := setupRadiationRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiation/chemicals/suspicious", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRadiation_GetDashboard(t *testing.T) {
	r := setupRadiationRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiation/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRadiation_RegisterSource_BadJSON(t *testing.T) {
	r := setupRadiationRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/radiation/sources",
		strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRadiation_UpdateStatus_InvalidID(t *testing.T) {
	r := setupRadiationRouter()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/radiation/sources/bad/status",
		strings.NewReader(`{"status":"LOST"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRadiation_RegisterChemical_BadJSON(t *testing.T) {
	r := setupRadiationRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/radiation/chemicals",
		strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
