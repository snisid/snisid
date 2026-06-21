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

type mockInfraSvc struct{}

func setupInfraTest() (*gin.Engine, *mockInfraSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/infra")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	api.GET("/clusters", func(c *gin.Context) {
		c.JSON(http.StatusOK, []string{})
	})
	api.POST("/dr/drill", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	})
	return r, &mockInfraSvc{}
}

func TestInfraHealth_Success(t *testing.T) {
	r, _ := setupInfraTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/infra/health", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestInfraClusters_Success(t *testing.T) {
	r, _ := setupInfraTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/infra/clusters", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestInfraDRDrill_Success(t *testing.T) {
	r, _ := setupInfraTest()
	body, _ := json.Marshal(map[string]any{"scenario": "earthquake", "success": true})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/infra/dr/drill", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestInfraDRDrill_BadRequest(t *testing.T) {
	r, _ := setupInfraTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/infra/dr/drill", nil)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
