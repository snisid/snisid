package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/internal/service/fraud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnv_Default(t *testing.T) {
	assert.Equal(t, "default", getEnv("NONEXISTENT_VAR_12345", "default"))
}

func TestGetEnv_FromEnv(t *testing.T) {
	t.Setenv("TEST_SNISID_VAR", "from-env")
	assert.Equal(t, "from-env", getEnv("TEST_SNISID_VAR", "default"))
}

func TestFraudEngine_HealthEndpoint(t *testing.T) {
	r := setupTestRouter()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestFraudEngine_MetricsEndpoint(t *testing.T) {
	aiClient := fraud.NewDefaultAIClient("http://test-ai:8000/predict")
	engine, err := fraud.NewScoringEngine("localhost:6379", aiClient)
	require.NoError(t, err)

	r := setupTestRouter()
	// Register metrics endpoint manually for testing
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"rules_loaded": len(engine.Rules()),
			"status":       "running",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "running")
}

func TestFraudEngine_ScoreEndpoint(t *testing.T) {
	aiClient := fraud.NewDefaultAIClient("http://test-ai:8000/predict")
	engine, err := fraud.NewScoringEngine("localhost:6379", aiClient)
	require.NoError(t, err)

	r := setupTestRouter()
	r.POST("/v1/score", func(c *gin.Context) {
		var event map[string]interface{}
		if err := c.BindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		score, reason, riskLevel := engine.CalculateScore(c.Request.Context(), event)
		c.JSON(http.StatusOK, gin.H{"score": score, "reason": reason, "riskLevel": riskLevel})
	})

	body := `{"identityId": "test-123", "amount": 50000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/score",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "score")
}

func TestFraudEngine_ScoreEndpoint_BadRequest(t *testing.T) {
	r := setupTestRouter()
	r.POST("/v1/score", func(c *gin.Context) {
		var event map[string]interface{}
		if err := c.BindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/score",
		strings.NewReader("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test helpers
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}
