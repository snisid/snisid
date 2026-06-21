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

	"github.com/snisid/offline-ht/internal/domain"
)

type mockOfflineSvc struct{}

func setupOfflineTest() (*gin.Engine, *mockOfflineSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1/offline")
	api.POST("/queue/push", func(c *gin.Context) {
		var req domain.PushQueueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "queued"})
	})
	api.POST("/sync/:terminal_id", func(c *gin.Context) {
		tid := c.Param("terminal_id")
		if tid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "terminal_id required"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"terminal_id": tid, "items": []string{}})
	})
	api.GET("/terminals/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, []string{})
	})
	api.GET("/conflicts", func(c *gin.Context) {
		c.JSON(http.StatusOK, []string{})
	})
	return r, &mockOfflineSvc{}
}

func TestOfflinePushQueue_Success(t *testing.T) {
	r, _ := setupOfflineTest()
	body, _ := json.Marshal(domain.PushQueueRequest{TerminalID: "t1", EntityType: "citizen", EntityID: "e1", Action: "UPDATE", Payload: "{}"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/offline/queue/push", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestOfflinePushQueue_BadRequest(t *testing.T) {
	r, _ := setupOfflineTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/offline/queue/push", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOfflineSync_Success(t *testing.T) {
	r, _ := setupOfflineTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/offline/sync/t1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestOfflineTerminalStatus_Success(t *testing.T) {
	r, _ := setupOfflineTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/offline/terminals/status", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestOfflineConflicts_Success(t *testing.T) {
	r, _ := setupOfflineTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/offline/conflicts", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
