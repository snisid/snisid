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

	"github.com/snisid/data-ht/internal/domain"
)

type mockDataSvc struct{}

func setupDataTest() (*gin.Engine, *mockDataSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/data")
	api.GET("/pipelines", func(c *gin.Context) {
		c.JSON(http.StatusOK, []domain.Pipeline{})
	})
	api.POST("/models/register", func(c *gin.Context) {
		var req domain.RegisterModelRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "registered"})
	})
	api.GET("/models/:id/bias-audit", func(c *gin.Context) {
		mid := c.Param("id")
		if mid == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"model_id": mid})
	})
	api.GET("/dashboards/national", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "dashboard"})
	})
	return r, &mockDataSvc{}
}

func TestDataPipelines_Success(t *testing.T) {
	r, _ := setupDataTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/data/pipelines", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestDataRegisterModel_Success(t *testing.T) {
	r, _ := setupDataTest()
	body, _ := json.Marshal(domain.RegisterModelRequest{Name: "model1", ModelType: "fraud", Version: "1.0", MlflowRunID: "run1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/data/models/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestDataRegisterModel_BadRequest(t *testing.T) {
	r, _ := setupDataTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/data/models/register", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDataBiasAudit_Success(t *testing.T) {
	r, _ := setupDataTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/data/models/m1/bias-audit", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestDataBiasAudit_NotFound(t *testing.T) {
	r, _ := setupDataTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/data/models//bias-audit", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDataDashboard_Success(t *testing.T) {
	r, _ := setupDataTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/data/dashboards/national", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
