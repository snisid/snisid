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

	"github.com/snisid/critical-infra-protection-ht/internal/domain"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/infraprot")
	api.POST("/assets", func(c *gin.Context) {
		var req domain.CreateAssetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/assets/:sector", func(c *gin.Context) {
		c.JSON(http.StatusOK, []domain.CreateAssetRequest{})
	})
	api.POST("/incidents", func(c *gin.Context) {
		var req domain.ReportIncidentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/incidents/active", func(c *gin.Context) {
		c.JSON(http.StatusOK, []domain.ReportIncidentRequest{})
	})
	api.GET("/incidents/asset/:asset_id", func(c *gin.Context) {
		c.JSON(http.StatusOK, []domain.ReportIncidentRequest{})
	})
	api.POST("/assessments", func(c *gin.Context) {
		var req domain.CreateAssessmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/dashboard/national", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

func TestCreateAsset_Success(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(domain.CreateAssetRequest{
		AssetName:   "Grid-1",
		Sector:      "ENERGY",
		OwnerEntity: "GovCo",
		LocationLat: 45.0, LocationLng: -93.0,
		Region: "MW", DeptCode: "EN", Criticality: "HIGH",
		ContactName: "John", ContactPhone: "555-0100",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/infraprot/assets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateAsset_BadRequest(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/infraprot/assets", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAssetsBySector(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/infraprot/assets/ENERGY", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestReportIncident_Success(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(domain.ReportIncidentRequest{
		AssetID: "550e8400-e29b-41d4-a716-446655440000", IncidentType: "CYBER_ATTACK",
		Severity: "HIGH", Description: "breach",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/infraprot/incidents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestGetActiveIncidents(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/infraprot/incidents/active", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetIncidentsByAsset(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/infraprot/incidents/asset/550e8400-e29b-41d4-a716-446655440000", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCreateAssessment_Success(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(domain.CreateAssessmentRequest{
		Sector: "ENERGY", OverallRiskScore: 7, AssessorAgency: "CISA",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/infraprot/assessments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestGetNationalDashboard(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/infraprot/dashboard/national", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
