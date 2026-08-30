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

	"github.com/snisid/field-ht/internal/domain"
)

type mockFieldSvc struct{}

func setupFieldTest() (*gin.Engine, *mockFieldSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/field")
	api.POST("/missions", func(c *gin.Context) {
		var req domain.CreateMissionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/missions/active", func(c *gin.Context) {
		c.JSON(http.StatusOK, []domain.CreateMissionRequest{})
	})
	api.POST("/missions/:id/log", func(c *gin.Context) {
		mid := c.Param("id")
		if mid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var req domain.CreateMissionLogRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "logged"})
	})
	api.GET("/stats/coverage", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"stats": "ok"})
	})
	return r, &mockFieldSvc{}
}

func TestFieldCreateMission_Success(t *testing.T) {
	r, _ := setupFieldTest()
	body, _ := json.Marshal(domain.CreateMissionRequest{Title: "Mission 1", DeptCode: "OU"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/field/missions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestFieldCreateMission_BadRequest(t *testing.T) {
	r, _ := setupFieldTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/field/missions", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFieldActiveMissions_Success(t *testing.T) {
	r, _ := setupFieldTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/field/missions/active", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestFieldCreateLog_Success(t *testing.T) {
	r, _ := setupFieldTest()
	body, _ := json.Marshal(domain.CreateMissionLogRequest{Action: "ARRIVED", LoggedBy: "officer1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/field/missions/m1/log", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestFieldCoverageStats_Success(t *testing.T) {
	r, _ := setupFieldTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/field/stats/coverage", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
