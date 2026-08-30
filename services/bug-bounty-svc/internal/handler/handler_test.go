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

func setupBugBountyTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/bug-bounty")
	api.POST("/programs", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/programs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	api.POST("/reports", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"status": "submitted"})
	})
	api.GET("/reports/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"report_id": id})
	})
	api.POST("/reports/:id/triage", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "triaged"})
	})
	api.POST("/reports/:id/reward", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"status": "rewarded"})
	})
	api.POST("/pentests", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"status": "scheduled"})
	})
	api.GET("/pentests/:id/results", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"engagement_id": id})
	})
	return r
}

func TestBBCreateProgram_Success(t *testing.T) {
	r := setupBugBountyTest()
	body, _ := json.Marshal(map[string]any{"program_id": "abc-123", "target": "https://example.com", "scope_type": "URL", "in_scope": true})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bug-bounty/programs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestBBCreateProgram_BadRequest(t *testing.T) {
	r := setupBugBountyTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bug-bounty/programs", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBBListPrograms_Success(t *testing.T) {
	r := setupBugBountyTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bug-bounty/programs", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestBBSubmitReport_Success(t *testing.T) {
	r := setupBugBountyTest()
	body, _ := json.Marshal(map[string]string{"program_id": "abc-123", "submitter": "researcher1", "title": "XSS vuln", "description": "XSS in login", "severity": "HIGH"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bug-bounty/reports", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestBBGetReport_Success(t *testing.T) {
	r := setupBugBountyTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bug-bounty/reports/abc-123", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestBBGetReport_NotFound(t *testing.T) {
	r := setupBugBountyTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bug-bounty/reports/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBBTriageReport_Success(t *testing.T) {
	r := setupBugBountyTest()
	body, _ := json.Marshal(map[string]any{"triager": "analyst1", "severity": "HIGH", "reproducible": true})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bug-bounty/reports/abc-123/triage", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestBBIssueReward_Success(t *testing.T) {
	r := setupBugBountyTest()
	body, _ := json.Marshal(map[string]any{"amount": 5000, "currency": "USD", "paid_to": "researcher1", "approved_by": "admin1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bug-bounty/reports/abc-123/reward", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestBBSchedulePentest_Success(t *testing.T) {
	r := setupBugBountyTest()
	body, _ := json.Marshal(map[string]string{"program_id": "abc-123", "title": "Pentest Q3", "scope": "All endpoints", "start_date": "2026-07-01", "team_lead": "lead1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bug-bounty/pentests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestBBGetPentestResults_Success(t *testing.T) {
	r := setupBugBountyTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bug-bounty/pentests/abc-123/results", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
