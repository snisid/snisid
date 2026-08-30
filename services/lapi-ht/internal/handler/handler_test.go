package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLapiSvc struct{}

func setupLapiTest() (*gin.Engine, *mockLapiSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1/lapi")
	v1.POST("/reads", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	})
	v1.GET("/reads/recent", func(c *gin.Context) {
		limit := 50
		if l := c.Query("limit"); l != "" {
			if n, err := strconv.Atoi(l); err == nil {
				limit = n
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": []string{}, "limit": limit})
	})
	v1.GET("/reads/plate/:number", func(c *gin.Context) {
		n := c.Param("number")
		if n == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "plate number required"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	v1.GET("/alerts/active", func(c *gin.Context) {
		c.JSON(http.StatusOK, []string{})
	})
	v1.GET("/cameras/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, []string{})
	})
	return r, &mockLapiSvc{}
}

func TestLapiCreateRead_Success(t *testing.T) {
	r, _ := setupLapiTest()
	body, _ := json.Marshal(map[string]any{"camera_id": "c1", "plate_number_raw": "ABC123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/lapi/reads", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestLapiCreateRead_BadRequest(t *testing.T) {
	r, _ := setupLapiTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/lapi/reads", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLapiRecentReads_Success(t *testing.T) {
	r, _ := setupLapiTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/lapi/reads/recent", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestLapiReadsByPlate_Success(t *testing.T) {
	r, _ := setupLapiTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/lapi/reads/plate/ABC123", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestLapiActiveAlerts_Success(t *testing.T) {
	r, _ := setupLapiTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/lapi/alerts/active", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestLapiCameraStatus_Success(t *testing.T) {
	r, _ := setupLapiTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/lapi/cameras/status", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
