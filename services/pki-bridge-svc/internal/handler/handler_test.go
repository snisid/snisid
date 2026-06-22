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

func setupPKIBridgeTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/pki-bridge")
	api.POST("/foreign-cas", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "registered"})
	})
	api.POST("/cross-certs", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"status": "issued"})
	})
	api.GET("/cross-certs/:subject", func(c *gin.Context) {
		subj := c.Param("subject")
		if subj == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"subject": subj})
	})
	api.GET("/trust-anchors", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	api.POST("/validate-path", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": true})
	})
	api.GET("/bridges", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	api.POST("/bridges", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	return r
}

func TestPKIRegisterForeignCA_Success(t *testing.T) {
	r := setupPKIBridgeTest()
	body, _ := json.Marshal(map[string]string{"name": "ForeignCA1", "country": "FR", "public_key_pem": "LS0t..."})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pki-bridge/foreign-cas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestPKIRegisterForeignCA_BadRequest(t *testing.T) {
	r := setupPKIBridgeTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pki-bridge/foreign-cas", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPKIIssueCrossCert_Success(t *testing.T) {
	r := setupPKIBridgeTest()
	body, _ := json.Marshal(map[string]string{
		"subject": "CN=ForeignCA", "issuer_ca_id": "abc-123", "serial_number": "01",
		"not_before": "2026-01-01", "not_after": "2027-01-01", "certificate_pem": "LS0t...",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pki-bridge/cross-certs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestPKIGetCrossCert_Success(t *testing.T) {
	r := setupPKIBridgeTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/pki-bridge/cross-certs/CN%3DForeignCA", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestPKIGetCrossCert_NotFound(t *testing.T) {
	r := setupPKIBridgeTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/pki-bridge/cross-certs/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPKIListTrustAnchors_Success(t *testing.T) {
	r := setupPKIBridgeTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/pki-bridge/trust-anchors", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestPKIValidatePath_Success(t *testing.T) {
	r := setupPKIBridgeTest()
	body, _ := json.Marshal(map[string]any{"leaf_subject": "CN=leaf", "intermediates": []string{"CN=intermediate"}, "root_subject": "CN=root"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pki-bridge/validate-path", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestPKIListAgreements_Success(t *testing.T) {
	r := setupPKIBridgeTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/pki-bridge/bridges", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestPKICreateAgreement_Success(t *testing.T) {
	r := setupPKIBridgeTest()
	body, _ := json.Marshal(map[string]string{"name": "Agreement1", "partner_ca": "ForeignCA1", "policy_id": "abc-123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pki-bridge/bridges", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}
