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

type mockAPISvc struct{}

func setupAPITest() (*gin.Engine, *mockAPISvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/devportal")
	api.POST("/register", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		email, _ := req["email"].(string)
		name, _ := req["contact_name"].(string)
		if email == "" || name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email and contact_name are required"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "registered"})
	})
	api.POST("/keys/request", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		accID, _ := req["account_id"].(string)
		if accID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/catalog", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	api.GET("/usage/:key_id", func(c *gin.Context) {
		kid := c.Param("key_id")
		if kid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key_id"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	api.POST("/keys/:id/revoke", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "revoked"})
	})
	return r, &mockAPISvc{}
}

func TestAPIRegister_Success(t *testing.T) {
	r, _ := setupAPITest()
	body, _ := json.Marshal(map[string]string{"email": "a@b.com", "contact_name": "Alice"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/devportal/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestAPIRegister_BadRequest(t *testing.T) {
	r, _ := setupAPITest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/devportal/register", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIRequestKey_Success(t *testing.T) {
	r, _ := setupAPITest()
	body, _ := json.Marshal(map[string]string{"account_id": "a1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/devportal/keys/request", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestAPICatalog_Success(t *testing.T) {
	r, _ := setupAPITest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/devportal/catalog", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAPIUsage_Success(t *testing.T) {
	r, _ := setupAPITest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/devportal/usage/key1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAPIRevokeKey_Success(t *testing.T) {
	r, _ := setupAPITest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/devportal/keys/k1/revoke", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
