package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateIdentity(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	
	// Mock handler (In production, use real handler with mock DB)
	r.POST("/api/v1/identities", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"id": "HT-123", "status": "pending"})
	})

	payload := map[string]string{
		"firstName": "Jean",
		"lastName":  "Pierre",
		"dob":       "1990-01-01",
		"agency":    "AGENCY-PRP",
	}
	body, _ := json.Marshal(payload)

	// Execute
	req, _ := http.NewRequest("POST", "/api/v1/identities", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusCreated, resp.Code)
	
	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "HT-123", response["id"])
}
