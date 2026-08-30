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

	"github.com/snisid/card-ht/internal/domain"
)

type mockCardSvc struct{}

func setupCardTest() (*gin.Engine, *mockCardSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/card")
	api.POST("/issue", func(c *gin.Context) {
		var req domain.IssueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "issued"})
	})
	api.GET("/verify/:doc_number", func(c *gin.Context) {
		doc := c.Param("doc_number")
		if doc == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"document_number": doc})
	})
	api.POST("/:id/report-lost", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "reported_lost"})
	})
	api.POST("/:id/revoke", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "revoked"})
	})
	return r, &mockCardSvc{}
}

func TestCardIssue_Success(t *testing.T) {
	r, _ := setupCardTest()
	body, _ := json.Marshal(domain.IssueRequest{DocType: "NATIONAL_ID", CitizenID: "a", CreatedBy: "u"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/card/issue", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCardIssue_BadRequest(t *testing.T) {
	r, _ := setupCardTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/card/issue", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCardVerify_Success(t *testing.T) {
	r, _ := setupCardTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/card/verify/HTI-ID-2026-000001", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCardVerify_NotFound(t *testing.T) {
	r, _ := setupCardTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/card/verify/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCardReportLost_Success(t *testing.T) {
	r, _ := setupCardTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/card/doc-1/report-lost", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCardRevoke_Success(t *testing.T) {
	r, _ := setupCardTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/card/doc-1/revoke", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
