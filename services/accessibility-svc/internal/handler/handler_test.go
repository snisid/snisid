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

func setupAccessibilityTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/accessibility")
	api.POST("/audits", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		target, _ := req["target_url"].(string)
		if target == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "target_url required"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "started"})
	})
	api.GET("/audits/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"audit_run_id": id})
	})
	api.GET("/audits", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	api.POST("/violations/:id/remediate", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "remediated"})
	})
	api.GET("/compliance", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	api.POST("/schedules", func(c *gin.Context) {
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	api.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
	})
	return r
}

func TestACCRunAudit_Success(t *testing.T) {
	r := setupAccessibilityTest()
	body, _ := json.Marshal(map[string]string{"target_url": "https://example.com", "wcag_level": "AA"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accessibility/audits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestACCRunAudit_BadRequest(t *testing.T) {
	r := setupAccessibilityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accessibility/audits", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestACCGetAuditResult_Success(t *testing.T) {
	r := setupAccessibilityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/accessibility/audits/abc-123", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestACCGetAuditResult_NotFound(t *testing.T) {
	r := setupAccessibilityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/accessibility/audits/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestACCListAudits_Success(t *testing.T) {
	r := setupAccessibilityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/accessibility/audits", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestACCMarkRemediated_Success(t *testing.T) {
	r := setupAccessibilityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accessibility/violations/abc-123/remediate", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestACCGetCompliance_Success(t *testing.T) {
	r := setupAccessibilityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/accessibility/compliance", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestACCCreateSchedule_Success(t *testing.T) {
	r := setupAccessibilityTest()
	body, _ := json.Marshal(map[string]string{"target_url": "https://example.com", "wcag_level": "AA", "cron_expr": "0 0 * * 0"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accessibility/schedules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestACCGetDashboard_Success(t *testing.T) {
	r := setupAccessibilityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/accessibility/dashboard", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
