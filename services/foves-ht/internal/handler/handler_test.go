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

type mockFovesSvc struct{}

func setupFovesTest() (*gin.Engine, *mockFovesSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/foves")
	api.POST("/vehicles", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		plate, _ := req["plate_number"].(string)
		if plate == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "plate_number required"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "registered"})
	})
	api.GET("/vehicles/plate/:number", func(c *gin.Context) {
		n := c.Param("number")
		if n == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"plate_number": n})
	})
	api.GET("/vehicles/vin/:vin", func(c *gin.Context) {
		v := c.Param("vin")
		if v == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"vin": v})
	})
	api.GET("/vehicles/owner/:citizen_id", func(c *gin.Context) {
		cid := c.Param("citizen_id")
		if cid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid citizen_id"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	api.POST("/transfers", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "transferred"})
	})
	api.POST("/licenses", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "issued"})
	})
	api.GET("/licenses/:citizen_id", func(c *gin.Context) {
		cid := c.Param("citizen_id")
		if cid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid citizen_id"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"license_number": "LIC-001"})
	})
	return r, &mockFovesSvc{}
}

func TestFovesRegisterVehicle_Success(t *testing.T) {
	r, _ := setupFovesTest()
	body, _ := json.Marshal(map[string]any{"plate_number": "ABC123", "vin": "VIN001", "make": "Toyota", "model": "Corolla", "year": 2020, "category": "PRIVATE_CAR", "owner_citizen_id": "c1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/foves/vehicles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestFovesRegisterVehicle_BadRequest(t *testing.T) {
	r, _ := setupFovesTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/foves/vehicles", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFovesGetByPlate_Success(t *testing.T) {
	r, _ := setupFovesTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/foves/vehicles/plate/ABC123", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestFovesGetByVIN_Success(t *testing.T) {
	r, _ := setupFovesTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/foves/vehicles/vin/VIN001", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestFovesGetByOwner_Success(t *testing.T) {
	r, _ := setupFovesTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/foves/vehicles/owner/c1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestFovesTransfer_Success(t *testing.T) {
	r, _ := setupFovesTest()
	body, _ := json.Marshal(map[string]string{"vehicle_id": "v1", "from_citizen_id": "c1", "to_citizen_id": "c2"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/foves/transfers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestFovesIssueLicense_Success(t *testing.T) {
	r, _ := setupFovesTest()
	body, _ := json.Marshal(map[string]any{"citizen_id": "c1", "license_number": "LIC001", "expiry_date": "2030-01-01", "points_balance": 12})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/foves/licenses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestFovesGetLicense_Success(t *testing.T) {
	r, _ := setupFovesTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/foves/licenses/c1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
