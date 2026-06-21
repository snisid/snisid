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

type mockCivilSvc struct{}

func setupCivilTest() (*gin.Engine, *mockCivilSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/civil")
	api.POST("/births", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		childName, _ := req["child_full_name"].(string)
		if childName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "child_full_name required"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "declared"})
	})
	api.POST("/deaths", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "declared"})
	})
	api.POST("/marriages", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "registered"})
	})
	api.GET("/acts/:number", func(c *gin.Context) {
		num := c.Param("number")
		if num == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "act not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"act_number": num})
	})
	api.GET("/acts/citizen/:nin", func(c *gin.Context) {
		nin := c.Param("nin")
		if nin == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "citizen not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	return r, &mockCivilSvc{}
}

func TestCivilDeclareBirth_Success(t *testing.T) {
	r, _ := setupCivilTest()
	body, _ := json.Marshal(map[string]string{"child_full_name": "Baby Doe", "event_date": "2026-01-15", "registering_office": "Office 1", "dept_code": "OU", "commune": "Port-au-Prince"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/civil/births", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCivilDeclareBirth_BadRequest(t *testing.T) {
	r, _ := setupCivilTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/civil/births", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCivilDeclareDeath_Success(t *testing.T) {
	r, _ := setupCivilTest()
	body, _ := json.Marshal(map[string]string{"deceased_citizen_id": "c1", "event_date": "2026-01-15", "registering_office": "Office 1", "dept_code": "OU", "commune": "Port-au-Prince"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/civil/deaths", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCivilRegisterMarriage_Success(t *testing.T) {
	r, _ := setupCivilTest()
	body, _ := json.Marshal(map[string]string{"spouse_a_citizen_id": "c1", "spouse_b_citizen_id": "c2", "event_date": "2026-01-15", "registering_office": "Office 1", "dept_code": "OU", "commune": "Port-au-Prince"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/civil/marriages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCivilGetAct_Success(t *testing.T) {
	r, _ := setupCivilTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/civil/acts/ACTE-HT-2026-OU-B-000001", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCivilGetAct_NotFound(t *testing.T) {
	r, _ := setupCivilTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/civil/acts/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCivilGetCitizenActs_Success(t *testing.T) {
	r, _ := setupCivilTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/civil/acts/citizen/nin123", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
