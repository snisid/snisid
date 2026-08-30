package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/all-source-fusion-ht/internal/domain"
	"github.com/snisid/all-source-fusion-ht/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFusionSvc struct {
	service.FusionService
	createProductFn        func(ctx context.Context, p *domain.IntelProduct) error
	getRecentProductsFn    func(ctx context.Context) ([]domain.IntelProduct, error)
	createThreatActorFn    func(ctx context.Context, a *domain.ThreatActor) error
	getHighRiskActorsFn    func(ctx context.Context) ([]domain.ThreatActor, error)
	createCorrelationFn    func(ctx context.Context, c *domain.CrossDisciplineCorrelation) error
	getSourceMapFn         func(ctx context.Context, productID uuid.UUID) (*domain.IntelProduct, error)
	getNationalEstimatesFn func(ctx context.Context) ([]domain.IntelProduct, error)
}

func (m *mockFusionSvc) CreateProduct(ctx context.Context, p *domain.IntelProduct) error {
	return m.createProductFn(ctx, p)
}
func (m *mockFusionSvc) GetRecentProducts(ctx context.Context) ([]domain.IntelProduct, error) {
	return m.getRecentProductsFn(ctx)
}
func (m *mockFusionSvc) CreateThreatActor(ctx context.Context, a *domain.ThreatActor) error {
	return m.createThreatActorFn(ctx, a)
}
func (m *mockFusionSvc) GetHighRiskActors(ctx context.Context) ([]domain.ThreatActor, error) {
	return m.getHighRiskActorsFn(ctx)
}
func (m *mockFusionSvc) CreateCorrelation(ctx context.Context, c *domain.CrossDisciplineCorrelation) error {
	return m.createCorrelationFn(ctx, c)
}
func (m *mockFusionSvc) GetSourceMap(ctx context.Context, productID uuid.UUID) (*domain.IntelProduct, error) {
	return m.getSourceMapFn(ctx, productID)
}
func (m *mockFusionSvc) GetNationalEstimates(ctx context.Context) ([]domain.IntelProduct, error) {
	return m.getNationalEstimatesFn(ctx)
}

func setupFusionHandler(svc service.FusionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewFusionHandler(svc)
	r := gin.New()
	h.RegisterRoutes(r)
	return r
}

func TestCreateProductHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcErr     error
		wantStatus int
	}{
		{name: "success", body: `{"title":"Intel Report","classification":"SECRET","source_disciplines":["SIGINT"],"analyst_assessment":"Assessment","confidence_level":"HIGH","related_threat_actors":[],"related_regions":[],"created_by":"00000000-0000-0000-0000-000000000001"}`, wantStatus: http.StatusCreated},
		{name: "bad json", body: `invalid`, wantStatus: http.StatusBadRequest},
		{name: "service error", body: `{"title":"Intel Report","classification":"SECRET","source_disciplines":["SIGINT"],"analyst_assessment":"Assessment","confidence_level":"HIGH","related_threat_actors":[],"related_regions":[],"created_by":"00000000-0000-0000-0000-000000000001"}`, svcErr: errors.New("svc error"), wantStatus: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockFusionSvc{
				createProductFn: func(_ context.Context, _ *domain.IntelProduct) error {
					return tt.svcErr
				},
			}
			r := setupFusionHandler(svc)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/fusion/products", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetRecentProductsHandler(t *testing.T) {
	tests := []struct {
		name       string
		result     []domain.IntelProduct
		svcErr     error
		wantStatus int
	}{
		{name: "success", result: []domain.IntelProduct{}, wantStatus: http.StatusOK},
		{name: "service error", svcErr: errors.New("svc error"), wantStatus: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockFusionSvc{
				getRecentProductsFn: func(_ context.Context) ([]domain.IntelProduct, error) {
					return tt.result, tt.svcErr
				},
			}
			r := setupFusionHandler(svc)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/fusion/products/recent", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetSourceMapHandler(t *testing.T) {
	validID := uuid.New().String()
	product := &domain.IntelProduct{Title: "Report 1"}
	tests := []struct {
		name       string
		id         string
		result     *domain.IntelProduct
		svcErr     error
		wantStatus int
	}{
		{name: "found", id: validID, result: product, wantStatus: http.StatusOK},
		{name: "not found", id: validID, result: nil, wantStatus: http.StatusNotFound},
		{name: "invalid id", id: "bad-id", wantStatus: http.StatusBadRequest},
		{name: "service error", id: validID, svcErr: errors.New("svc error"), wantStatus: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockFusionSvc{
				getSourceMapFn: func(_ context.Context, _ uuid.UUID) (*domain.IntelProduct, error) {
					return tt.result, tt.svcErr
				},
			}
			r := setupFusionHandler(svc)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/fusion/products/"+tt.id+"/source-map", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCreateThreatActorHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcErr     error
		wantStatus int
	}{
		{name: "success", body: `{"name":"APT-42","type":"STATE","cap_level":5,"intent_level":4,"opportunity_level":3,"overall_risk":4,"primary_region":"EU","associated_groups":[],"ofac_designated":true,"notes":"State actor"}`, wantStatus: http.StatusCreated},
		{name: "bad json", body: `invalid`, wantStatus: http.StatusBadRequest},
		{name: "service error", body: `{"name":"APT-42","type":"STATE","cap_level":5,"intent_level":4,"opportunity_level":3,"overall_risk":4,"primary_region":"EU","associated_groups":[],"ofac_designated":true,"notes":"State actor"}`, svcErr: errors.New("svc error"), wantStatus: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockFusionSvc{
				createThreatActorFn: func(_ context.Context, _ *domain.ThreatActor) error {
					return tt.svcErr
				},
			}
			r := setupFusionHandler(svc)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/fusion/threat-actors", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetHighRiskActorsHandler(t *testing.T) {
	tests := []struct {
		name       string
		result     []domain.ThreatActor
		svcErr     error
		wantStatus int
	}{
		{name: "success", result: []domain.ThreatActor{}, wantStatus: http.StatusOK},
		{name: "service error", svcErr: errors.New("svc error"), wantStatus: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockFusionSvc{
				getHighRiskActorsFn: func(_ context.Context) ([]domain.ThreatActor, error) {
					return tt.result, tt.svcErr
				},
			}
			r := setupFusionHandler(svc)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/fusion/threat-actors/high-risk", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCreateCorrelationHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcErr     error
		wantStatus int
	}{
		{name: "success", body: `{"discipline_a":"SIGINT","reference_a":"REF-A","discipline_b":"HUMINT","reference_b":"REF-B","correlation_type":"SUPPORTS","analyst_notes":"Matches","score":0.95}`, wantStatus: http.StatusCreated},
		{name: "bad json", body: `invalid`, wantStatus: http.StatusBadRequest},
		{name: "service error", body: `{"discipline_a":"SIGINT","reference_a":"REF-A","discipline_b":"HUMINT","reference_b":"REF-B","correlation_type":"SUPPORTS","analyst_notes":"Matches","score":0.95}`, svcErr: errors.New("svc error"), wantStatus: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockFusionSvc{
				createCorrelationFn: func(_ context.Context, _ *domain.CrossDisciplineCorrelation) error {
					return tt.svcErr
				},
			}
			r := setupFusionHandler(svc)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/fusion/correlations", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetNationalEstimatesHandler(t *testing.T) {
	tests := []struct {
		name       string
		result     []domain.IntelProduct
		svcErr     error
		wantStatus int
	}{
		{name: "success", result: []domain.IntelProduct{}, wantStatus: http.StatusOK},
		{name: "service error", svcErr: errors.New("svc error"), wantStatus: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockFusionSvc{
				getNationalEstimatesFn: func(_ context.Context) ([]domain.IntelProduct, error) {
					return tt.result, tt.svcErr
				},
			}
			r := setupFusionHandler(svc)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/fusion/estimates/national", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestNewFusionHandler(t *testing.T) {
	h := NewFusionHandler(&mockFusionSvc{})
	require.NotNil(t, h)
}
