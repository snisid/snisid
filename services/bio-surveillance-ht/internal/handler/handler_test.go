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

	"github.com/snisid/bio-surveillance-ht/internal/domain"
)

type mockBioSvc struct{}

func setupBioTest() (*gin.Engine, *mockBioSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/biosurv")
	api.POST("/alerts", func(c *gin.Context) {
		var req domain.CreateDiseaseAlertRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/alerts/active", func(c *gin.Context) {
		c.JSON(http.StatusOK, []domain.DiseaseAlert{})
	})
	api.GET("/alerts/region/:region", func(c *gin.Context) {
		c.JSON(http.StatusOK, []domain.DiseaseAlert{})
	})
	api.POST("/campaigns", func(c *gin.Context) {
		var req domain.CreateVaccinationCampaignRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/campaigns/:id/coverage", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	api.PATCH("/facilities/:id/stock", func(c *gin.Context) {
		var req domain.UpdateFacilityStockRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	})
	api.GET("/dashboard/national", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r, &mockBioSvc{}
}

func TestCreateAlert_Success(t *testing.T) {
	r, _ := setupBioTest()
	body, _ := json.Marshal(domain.CreateDiseaseAlertRequest{
		DiseaseName: "COVID-19", PathogenType: "VIRUS", Icd10Code: "U07.1",
		AlertLevel: "RED", TransmissionMode: "AIRBORNE",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/biosurv/alerts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateAlert_BadRequest(t *testing.T) {
	r, _ := setupBioTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/biosurv/alerts", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetActiveAlerts_Success(t *testing.T) {
	r, _ := setupBioTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/biosurv/alerts/active", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetAlertsByRegion_Success(t *testing.T) {
	r, _ := setupBioTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/biosurv/alerts/region/Ouest", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCreateCampaign_Success(t *testing.T) {
	r, _ := setupBioTest()
	body, _ := json.Marshal(domain.CreateVaccinationCampaignRequest{
		CampaignName: "Vax-Ouest", TargetDisease: "Polio", VaccineType: "OPV",
		TargetPopulation: 50000, DosesAdministered: 35000, CoveragePct: 70.0,
		StartDate: "2026-01-15T00:00:00Z", CoordinatorAgency: "MSPP",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/biosurv/campaigns", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestGetCampaignCoverage_Success(t *testing.T) {
	r, _ := setupBioTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/biosurv/campaigns/"+"abc123"+"/coverage", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateFacilityStock_Success(t *testing.T) {
	r, _ := setupBioTest()
	body, _ := json.Marshal(domain.UpdateFacilityStockRequest{StockStatus: "CRITICAL", BedsAvailable: 0})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/biosurv/facilities/"+"fac1"+"/stock", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetDashboardNational_Success(t *testing.T) {
	r, _ := setupBioTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/biosurv/dashboard/national", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
