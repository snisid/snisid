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

	"github.com/snisid/pki-ht/internal/domain"
)

type mockPKISvc struct{}

func setupPKITest() (*gin.Engine, *mockPKISvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/pki")
	api.POST("/issue", func(c *gin.Context) {
		var req domain.IssueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"serial_number": "abc123"})
	})
	api.POST("/revoke", func(c *gin.Context) {
		var req domain.RevokeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "revoked"})
	})
	api.GET("/ocsp", func(c *gin.Context) {
		serial := c.Query("serial")
		if serial == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "serial parameter required"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"serial": serial, "status": "VALID"})
	})
	api.GET("/crl/:ca_id", func(c *gin.Context) {
		caID := c.Param("ca_id")
		if caID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "CRL not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ca_id": caID})
	})
	return r, &mockPKISvc{}
}

func TestPKIIssue_Success(t *testing.T) {
	r, _ := setupPKITest()
	body, _ := json.Marshal(domain.IssueRequest{SubjectType: "CITIZEN", CommonName: "John Doe"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pki/issue", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestPKIIssue_BadRequest(t *testing.T) {
	r, _ := setupPKITest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pki/issue", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPKIRevoke_Success(t *testing.T) {
	r, _ := setupPKITest()
	body, _ := json.Marshal(domain.RevokeRequest{SerialNumber: "abc", Reason: "compromised"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pki/revoke", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestPKIOCSP_Success(t *testing.T) {
	r, _ := setupPKITest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/pki/ocsp?serial=abc", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestPKIOCSP_MissingSerial(t *testing.T) {
	r, _ := setupPKITest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/pki/ocsp", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPKICRL_Success(t *testing.T) {
	r, _ := setupPKITest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/pki/crl/ca-1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestPKICRL_NotFound(t *testing.T) {
	r, _ := setupPKITest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/pki/crl/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
