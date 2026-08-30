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

type mockInteropSvc struct{}

func setupInteropTest() (*gin.Engine, *mockInteropSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/interop")
	api.POST("/exchange", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "exchanged"})
	})
	api.POST("/agreements", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	})
	api.GET("/logs/:agreement_id", func(c *gin.Context) {
		aid := c.Param("agreement_id")
		if aid == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "logs not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	return r, &mockInteropSvc{}
}

func TestInteropExchange_Success(t *testing.T) {
	r, _ := setupInteropTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/interop/exchange", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestInteropCreateAgreement_Success(t *testing.T) {
	r, _ := setupInteropTest()
	body, _ := json.Marshal(map[string]string{"provider_agency_id": "a", "consumer_agency_id": "b", "service_name": "s"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/interop/agreements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestInteropCreateAgreement_BadRequest(t *testing.T) {
	r, _ := setupInteropTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/interop/agreements", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInteropGetLogs_Success(t *testing.T) {
	r, _ := setupInteropTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/interop/logs/agreement-1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestInteropGetLogs_NotFound(t *testing.T) {
	r, _ := setupInteropTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/interop/logs/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
