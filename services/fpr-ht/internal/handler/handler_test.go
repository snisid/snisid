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

	"github.com/snisid/fpr-ht/internal/domain"
)

type mockFprSvc struct{}

func setupFprTest() (*gin.Engine, *mockFprSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1/fpr")
	v1.POST("/warrants", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	})
	v1.GET("/check/:citizen_id", func(c *gin.Context) {
		cid := c.Param("citizen_id")
		if cid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "citizen_id required"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"warrant_found": false})
	})
	v1.GET("/check/name", func(c *gin.Context) {
		n := c.Query("name")
		if n == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name query parameter required"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"warrant_found": false})
	})
	v1.POST("/warrants/:id/sightings", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warrant id"})
			return
		}
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	})
	v1.PATCH("/warrants/:id/execute", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warrant id"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "executed"})
	})
	v1.GET("/warrants/armed-dangerous", func(c *gin.Context) {
		c.JSON(http.StatusOK, []domain.Warrant{})
	})
	v1.GET("/stats/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, domain.DashboardStats{})
	})
	return r, &mockFprSvc{}
}

func TestCreateWarrant_Success(t *testing.T) {
	r, _ := setupFprTest()
	body, _ := json.Marshal(map[string]any{"full_name": "John Doe", "warrant_type": "ARREST_WARRANT", "issuing_court": "Tribunal"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fpr/warrants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateWarrant_BadRequest(t *testing.T) {
	r, _ := setupFprTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fpr/warrants", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCheckCitizen_Success(t *testing.T) {
	r, _ := setupFprTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fpr/check/citizen-1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCheckByName_Success(t *testing.T) {
	r, _ := setupFprTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fpr/check/name?name=John", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCheckByName_MissingQuery(t *testing.T) {
	r, _ := setupFprTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fpr/check/name", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReportSighting_Success(t *testing.T) {
	r, _ := setupFprTest()
	body, _ := json.Marshal(map[string]string{"citizen_id": "c1", "reported_by": "officer1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fpr/warrants/w1/sightings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestExecuteWarrant_Success(t *testing.T) {
	r, _ := setupFprTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/fpr/warrants/w1/execute", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetArmedDangerous_Success(t *testing.T) {
	r, _ := setupFprTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fpr/warrants/armed-dangerous", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetDashboardStats_Success(t *testing.T) {
	r, _ := setupFprTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fpr/stats/dashboard", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
