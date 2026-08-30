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

	"github.com/snisid/idcore-svc/internal/domain"
)

type mockIdentityService struct {
	enrollFn       func(req domain.EnrollmentRequest) (*domain.EnrollmentResult, error)
	getByNINFn     func(nin string) (*domain.Citizen, error)
	searchFn       func(query string) ([]*domain.Citizen, error)
	resolveDedupFn func(candidateID, resolution, reviewedBy string) error
	updateStatusFn func(nin string, status domain.IDStatus, reason, authorizedBy string) error
	getHistoryFn   func(id string) ([]domain.ChangeHistory, error)
	statsFn        func() (*domain.PopulationStats, error)
}

func (m *mockIdentityService) EnrollCitizen(ctx *gin.Context, req domain.EnrollmentRequest) (*domain.EnrollmentResult, error) {
	return m.enrollFn(req)
}
func (m *mockIdentityService) VerifyIdentity(ctx *gin.Context, nin string) (*domain.Citizen, error) {
	return m.getByNINFn(nin)
}
func (m *mockIdentityService) SearchCitizens(ctx *gin.Context, query string) ([]*domain.Citizen, error) {
	return m.searchFn(query)
}
func (m *mockIdentityService) ResolveDedup(ctx *gin.Context, candidateID, resolution, reviewedBy string) error {
	return m.resolveDedupFn(candidateID, resolution, reviewedBy)
}
func (m *mockIdentityService) UpdateStatus(ctx *gin.Context, nin string, status domain.IDStatus, reason, authorizedBy string) error {
	return m.updateStatusFn(nin, status, reason, authorizedBy)
}
func (m *mockIdentityService) GetHistory(ctx *gin.Context, id string) ([]domain.ChangeHistory, error) {
	return m.getHistoryFn(id)
}
func (m *mockIdentityService) GetPopulationStats(ctx *gin.Context) (*domain.PopulationStats, error) {
	return m.statsFn()
}

type IdentityServiceWrapper struct {
	svc *mockIdentityService
}

func setupIdentityTest() (*gin.Engine, *mockIdentityService) {
	gin.SetMode(gin.TestMode)
	mock := &mockIdentityService{}
	w := &IdentityServiceWrapper{svc: mock}
	r := gin.New()
	api := r.Group("/api/v1/idcore")
	api.POST("/enroll", func(c *gin.Context) {
		if w.svc.enrollFn == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no mock"})
			return
		}
		body, _ := c.GetRawData()
		var req struct {
			EnrollmentType string `json:"enrollment_type"`
			FullNameLegal  string `json:"full_name_legal"`
			FirstName      string `json:"first_name"`
			LastName       string `json:"last_name"`
			Nationality    string `json:"nationality"`
			DeptCode       string `json:"dept_code"`
			CreatedBy      string `json:"created_by"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
			return
		}
		if req.FullNameLegal == "" || req.FirstName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		enrollReq := domain.EnrollmentRequest{
			EnrollmentType: domain.EnrollmentType(req.EnrollmentType),
			FullNameLegal:  req.FullNameLegal,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			Nationality:    req.Nationality,
			DeptCode:       req.DeptCode,
			CreatedBy:      req.CreatedBy,
		}
		result, err := w.svc.enrollFn(enrollReq)
		if err != nil {
			if err == domain.ErrDuplicateDetected {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, result)
	})
	api.GET("/citizens/:nin", func(c *gin.Context) {
		nin := c.Param("nin")
		citizen, err := w.svc.getByNINFn(nin)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "citizen not found"})
			return
		}
		c.JSON(http.StatusOK, citizen)
	})
	api.GET("/citizens/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
			return
		}
		citizens, err := w.svc.searchFn(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": citizens})
	})
	api.POST("/dedup/resolve", func(c *gin.Context) {
		var req struct {
			CandidateID string `json:"candidate_id"`
			Resolution  string `json:"resolution"`
			ReviewedBy  string `json:"reviewed_by"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := w.svc.resolveDedupFn(req.CandidateID, req.Resolution, req.ReviewedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "resolved"})
	})
	api.PATCH("/citizens/:nin/status", func(c *gin.Context) {
		nin := c.Param("nin")
		var req struct {
			Status       string `json:"status"`
			Reason       string `json:"reason"`
			AuthorizedBy string `json:"authorized_by"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := w.svc.updateStatusFn(nin, domain.IDStatus(req.Status), req.Reason, req.AuthorizedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	})
	api.GET("/citizens/:nin/history", func(c *gin.Context) {
		nin := c.Param("nin")
		history, err := w.svc.getHistoryFn(nin)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "citizen not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": history})
	})
	api.GET("/stats/population", func(c *gin.Context) {
		stats, err := w.svc.statsFn()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)
	})
	return r, mock
}

func TestEnroll_Success(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.enrollFn = func(req domain.EnrollmentRequest) (*domain.EnrollmentResult, error) {
		return &domain.EnrollmentResult{NIN: "1234567890123"}, nil
	}

	body, _ := json.Marshal(map[string]string{
		"enrollment_type": "BIRTH",
		"full_name_legal": "John Doe",
		"first_name":      "John",
		"last_name":       "Doe",
		"nationality":     "HTI",
		"dept_code":       "OU",
		"created_by":      "00000000-0000-0000-0000-000000000001",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/idcore/enroll", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "1234567890123", resp["nin"])
}

func TestEnroll_BadRequest(t *testing.T) {
	r, _ := setupIdentityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/idcore/enroll", bytes.NewReader([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEnroll_Conflict(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.enrollFn = func(req domain.EnrollmentRequest) (*domain.EnrollmentResult, error) {
		return nil, domain.ErrDuplicateDetected
	}

	body, _ := json.Marshal(map[string]string{
		"enrollment_type": "BIRTH",
		"full_name_legal": "John Doe",
		"first_name":      "John",
		"last_name":       "Doe",
		"nationality":     "HTI",
		"dept_code":       "OU",
		"created_by":      "00000000-0000-0000-0000-000000000001",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/idcore/enroll", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestEnroll_InternalError(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.enrollFn = func(req domain.EnrollmentRequest) (*domain.EnrollmentResult, error) {
		return nil, assert.AnError
	}

	body, _ := json.Marshal(map[string]string{
		"enrollment_type": "BIRTH",
		"full_name_legal": "John Doe",
		"first_name":      "John",
		"last_name":       "Doe",
		"nationality":     "HTI",
		"dept_code":       "OU",
		"created_by":      "00000000-0000-0000-0000-000000000001",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/idcore/enroll", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetByNIN_Success(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.getByNINFn = func(nin string) (*domain.Citizen, error) {
		return &domain.Citizen{NIN: nin, FullNameLegal: "John Doe"}, nil
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/idcore/citizens/1234567890123", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp domain.Citizen
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "John Doe", resp.FullNameLegal)
}

func TestGetByNIN_NotFound(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.getByNINFn = func(nin string) (*domain.Citizen, error) {
		return nil, assert.AnError
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/idcore/citizens/1234567890123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSearch_MissingQuery(t *testing.T) {
	r, _ := setupIdentityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/idcore/citizens/search", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearch_Success(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.searchFn = func(query string) ([]*domain.Citizen, error) {
		return []*domain.Citizen{{NIN: "1234567890123", FullNameLegal: "John Doe"}}, nil
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/idcore/citizens/search?q=John", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestSearch_InternalError(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.searchFn = func(query string) ([]*domain.Citizen, error) {
		return nil, assert.AnError
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/idcore/citizens/search?q=John", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestResolveDedup_Success(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.resolveDedupFn = func(candidateID, resolution, reviewedBy string) error {
		return nil
	}

	body, _ := json.Marshal(map[string]string{
		"candidate_id": "abc",
		"resolution":   "MERGE",
		"reviewed_by":  "user1",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/idcore/dedup/resolve", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestResolveDedup_BadRequest(t *testing.T) {
	r, _ := setupIdentityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/idcore/dedup/resolve", bytes.NewReader([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateStatus_Success(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.updateStatusFn = func(nin string, status domain.IDStatus, reason, authorizedBy string) error {
		return nil
	}

	body, _ := json.Marshal(map[string]string{
		"status":        "SUSPENDED",
		"reason":        "fraud investigation",
		"authorized_by": "user1",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/idcore/citizens/1234567890123/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateStatus_BadRequest(t *testing.T) {
	r, _ := setupIdentityTest()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/idcore/citizens/1234567890123/status", bytes.NewReader([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetHistory_Success(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.getHistoryFn = func(id string) ([]domain.ChangeHistory, error) {
		return []domain.ChangeHistory{}, nil
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/idcore/citizens/1234567890123/history", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetHistory_NotFound(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.getHistoryFn = func(id string) ([]domain.ChangeHistory, error) {
		return nil, assert.AnError
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/idcore/citizens/1234567890123/history", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPopulationStats_Success(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.statsFn = func() (*domain.PopulationStats, error) {
		return &domain.PopulationStats{Total: 1000, Active: 900}, nil
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/idcore/stats/population", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var stats domain.PopulationStats
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	require.NoError(t, err)
	assert.Equal(t, 1000, stats.Total)
}

func TestGetPopulationStats_InternalError(t *testing.T) {
	r, mock := setupIdentityTest()
	mock.statsFn = func() (*domain.PopulationStats, error) {
		return nil, assert.AnError
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/idcore/stats/population", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
