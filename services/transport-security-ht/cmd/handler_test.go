package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/snisid/transport-security-ht/internal/domain"
	"github.com/snisid/transport-security-ht/internal/handler"
	"github.com/snisid/transport-security-ht/internal/kafka"
	"github.com/snisid/transport-security-ht/internal/repository"
	"github.com/snisid/transport-security-ht/internal/service"
)

type mockTransportRepo struct {
	repository.TransportRepository
}

func (m *mockTransportRepo) GetRecentScreenings(ctx context.Context, limit int) ([]domain.PassengerScreening, error) {
	return []domain.PassengerScreening{}, nil
}

func (m *mockTransportRepo) CheckNoFly(ctx context.Context, identityRef string) (*domain.NoFlyPassenger, error) {
	return nil, nil
}

func (m *mockTransportRepo) GetZonesByAirport(ctx context.Context, airportCode string) ([]domain.AirportSecurityZone, error) {
	return []domain.AirportSecurityZone{}, nil
}

type mockKafkaProducer struct {
	kafka.Producer
}

func (m *mockKafkaProducer) Publish(ctx context.Context, key string, msg interface{}) error {
	return nil
}

func setupTransportRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	repo := &mockTransportRepo{}
	prod := &mockKafkaProducer{}
	svc := service.NewTransportService(repo, prod)
	h := handler.NewTransportHandler(svc)

	r := gin.New()
	h.RegisterRoutes(r)
	return r
}

func TestTransport_GetRecentScreenings(t *testing.T) {
	r := setupTransportRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/transport/screenings/recent", nil)
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

func TestTransport_CheckNoFly_NoMatch(t *testing.T) {
	r := setupTransportRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/transport/no-fly/check?identity=test123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"match":false`) {
		t.Errorf("expected no match, got %s", w.Body.String())
	}
}

func TestTransport_CheckNoFly_MissingIdentity(t *testing.T) {
	r := setupTransportRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/transport/no-fly/check", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTransport_GetZones(t *testing.T) {
	r := setupTransportRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/transport/zones/JFK", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTransport_LogScreening_BadJSON(t *testing.T) {
	r := setupTransportRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transport/screenings",
		strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTransport_ReportBreach_InvalidID(t *testing.T) {
	r := setupTransportRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transport/zones/bad-id/breach", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
