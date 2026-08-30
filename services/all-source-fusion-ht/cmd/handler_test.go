package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/all-source-fusion-ht/internal/domain"
	"github.com/snisid/all-source-fusion-ht/internal/handler"
	"github.com/snisid/all-source-fusion-ht/internal/kafka"
	"github.com/snisid/all-source-fusion-ht/internal/repository"
	"github.com/snisid/all-source-fusion-ht/internal/service"
)

type mockFusionRepo struct {
	repository.FusionRepository
}

func (m *mockFusionRepo) GetRecentProducts(ctx context.Context, limit int) ([]domain.IntelProduct, error) {
	return []domain.IntelProduct{}, nil
}

func (m *mockFusionRepo) GetHighRiskActors(ctx context.Context) ([]domain.ThreatActor, error) {
	return []domain.ThreatActor{}, nil
}

func (m *mockFusionRepo) GetNationalEstimates(ctx context.Context) ([]domain.IntelProduct, error) {
	return []domain.IntelProduct{}, nil
}

func (m *mockFusionRepo) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.IntelProduct, error) {
	return nil, nil
}

type mockFusionKafkaProducer struct {
	kafka.Producer
}

func (m *mockFusionKafkaProducer) Publish(ctx context.Context, key string, msg interface{}) error {
	return nil
}

func setupFusionRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	repo := &mockFusionRepo{}
	prod := &mockFusionKafkaProducer{}
	svc := service.NewFusionService(repo, prod)
	h := handler.NewFusionHandler(svc)

	r := gin.New()
	h.RegisterRoutes(r)
	return r
}

func TestFusion_GetRecentProducts(t *testing.T) {
	r := setupFusionRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/fusion/products/recent", nil)
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

func TestFusion_GetHighRiskActors(t *testing.T) {
	r := setupFusionRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/fusion/threat-actors/high-risk", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestFusion_GetNationalEstimates(t *testing.T) {
	r := setupFusionRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/fusion/estimates/national", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestFusion_GetSourceMap_NotFound(t *testing.T) {
	r := setupFusionRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/fusion/products/"+uuid.NewString()+"/source-map", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestFusion_CreateProduct_BadJSON(t *testing.T) {
	r := setupFusionRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fusion/products",
		strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestFusion_GetSourceMap_InvalidID(t *testing.T) {
	r := setupFusionRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/fusion/products/bad-id/source-map", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
