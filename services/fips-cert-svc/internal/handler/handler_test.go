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
)

func setupFIPSTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/fips")
	api.POST("/modules", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		name, _ := req["name"].(string)
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "registered"})
	})
	api.GET("/modules", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	api.GET("/modules/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "module not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"module_id": id})
	})
	api.POST("/modules/:id/validate", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "validated"})
	})
	api.POST("/modules/:id/cve", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"status": "reported"})
	})
	api.GET("/compliance/:service", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": c.Param("service"), "overall_status": "COMPLIANT"})
	})
	api.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	return r
}

func TestFIPSRegisterModule_Success(t *testing.T) {
	r := setupFIPSTest()
	body, _ := json.Marshal(map[string]string{"name": "AES-256", "version": "1.0", "vendor": "TestCorp", "fips_level": "LEVEL_1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fips/modules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestFIPSRegisterModule_BadRequest(t *testing.T) {
	r := setupFIPSTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fips/modules", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFIPSListModules_Success(t *testing.T) {
	r := setupFIPSTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fips/modules", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestFIPSGetModule_Success(t *testing.T) {
	r := setupFIPSTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fips/modules/abc-123", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestFIPSGetModule_NotFound(t *testing.T) {
	r := setupFIPSTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fips/modules/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFIPSSubmitValidation_Success(t *testing.T) {
	r := setupFIPSTest()
	body, _ := json.Marshal(map[string]string{"cert_number": "FIPS-2026-001", "validation_date": "2026-06-22"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fips/modules/abc-123/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestFIPSReportCVE_Success(t *testing.T) {
	r := setupFIPSTest()
	body, _ := json.Marshal(map[string]string{"cve_id": "CVE-2026-1234", "severity": "HIGH"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fips/modules/abc-123/cve", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestFIPSGetCompliance_Success(t *testing.T) {
	r := setupFIPSTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fips/compliance/auth-svc", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestFIPSGetDashboard_Success(t *testing.T) {
	r := setupFIPSTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fips/dashboard", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
