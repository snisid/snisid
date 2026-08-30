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

	"github.com/snisid/bio-ht/internal/domain"
)

type mockBioSvc struct {
	enrollFn   func(domain.EnrollRequest) (*domain.BioTemplate, error)
	verifyFn   func(domain.VerifyRequest) (*domain.VerifyResult, error)
	identifyFn func(domain.IdentifyRequest) (*domain.IdentifyResult, error)
	qualityFn  func(string) (float64, error)
}

func setupBio() (*gin.Engine, *mockBioSvc) {
	gin.SetMode(gin.TestMode)
	m := &mockBioSvc{}
	r := gin.New()
	api := r.Group("/api/v1/bio")
	api.POST("/enroll", func(c *gin.Context) {
		var req domain.EnrollRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		result, err := m.enrollFn(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, result)
	})
	api.POST("/verify", func(c *gin.Context) {
		var req domain.VerifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		result, err := m.verifyFn(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})
	api.POST("/identify", func(c *gin.Context) {
		var req domain.IdentifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		result, err := m.identifyFn(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})
	api.GET("/quality/:template_id", func(c *gin.Context) {
		tid := c.Param("template_id")
		q, err := m.qualityFn(tid)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"template_id": tid, "quality_score": q})
	})
	return r, m
}

func TestBioEnroll_Success(t *testing.T) {
	r, m := setupBio()
	m.enrollFn = func(req domain.EnrollRequest) (*domain.BioTemplate, error) {
		return &domain.BioTemplate{}, nil
	}
	body, _ := json.Marshal(domain.EnrollRequest{CitizenID: "a", Modality: "FINGERPRINT", CapturedBy: "u"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bio/enroll", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestBioEnroll_BadRequest(t *testing.T) {
	r, _ := setupBio()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bio/enroll", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBioEnroll_InternalError(t *testing.T) {
	r, m := setupBio()
	m.enrollFn = func(req domain.EnrollRequest) (*domain.BioTemplate, error) {
		return nil, assert.AnError
	}
	body, _ := json.Marshal(domain.EnrollRequest{CitizenID: "a", Modality: "FINGERPRINT", CapturedBy: "u"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bio/enroll", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBioVerify_Success(t *testing.T) {
	r, m := setupBio()
	m.verifyFn = func(req domain.VerifyRequest) (*domain.VerifyResult, error) {
		return &domain.VerifyResult{IsMatch: true}, nil
	}
	body, _ := json.Marshal(domain.VerifyRequest{CitizenID: "a", Modality: "FINGERPRINT"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bio/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestBioIdentify_Success(t *testing.T) {
	r, m := setupBio()
	m.identifyFn = func(req domain.IdentifyRequest) (*domain.IdentifyResult, error) {
		return &domain.IdentifyResult{}, nil
	}
	body, _ := json.Marshal(domain.IdentifyRequest{Modality: "FACE"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bio/identify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestBioQuality_Success(t *testing.T) {
	r, m := setupBio()
	m.qualityFn = func(id string) (float64, error) { return 0.85, nil }
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bio/quality/t1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestBioQuality_NotFound(t *testing.T) {
	r, m := setupBio()
	m.qualityFn = func(id string) (float64, error) { return 0, assert.AnError }
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bio/quality/t1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
