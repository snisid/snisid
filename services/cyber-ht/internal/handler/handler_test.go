package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/snisid/cyber-ht/internal/domain"
)

type mockCyberSvc struct{}

func setupCyberTest() (*gin.Engine, *mockCyberSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/cyber")
	api.POST("/incidents", func(c *gin.Context) {
		var req domain.CreateIncidentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/incidents/active", func(c *gin.Context) {
		c.JSON(http.StatusOK, []domain.CreateIncidentRequest{})
	})
	api.POST("/policies", func(c *gin.Context) {
		var req domain.CreatePolicyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/threat-intel/check", func(c *gin.Context) {
		indicator := c.Query("indicator")
		if indicator == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "indicator query parameter is required"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"indicator": indicator})
	})
	return r, &mockCyberSvc{}
}

func TestCyberIncident_Success(t *testing.T) {
	r, _ := setupCyberTest()
	body, _ := json.Marshal(domain.CreateIncidentRequest{Title: "Test", Severity: "HIGH", DetectedBy: "soc"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cyber/incidents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCyberIncident_BadRequest(t *testing.T) {
	r, _ := setupCyberTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cyber/incidents", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCyberActiveIncidents_Success(t *testing.T) {
	r, _ := setupCyberTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cyber/incidents/active", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCyberPolicy_Success(t *testing.T) {
	r, _ := setupCyberTest()
	body, _ := json.Marshal(domain.CreatePolicyRequest{Name: "ZT-1", Description: "policy desc", PolicyType: "NETWORK", Rules: []string{"r1"}, CreatedBy: "admin"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cyber/policies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCyberThreatIntel_Success(t *testing.T) {
	r, _ := setupCyberTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cyber/threat-intel/check?indicator=1.2.3.4", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCyberThreatIntel_MissingParam(t *testing.T) {
	r, _ := setupCyberTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cyber/threat-intel/check", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
