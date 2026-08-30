package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/entity-resolution/internal/handlers"
	"github.com/snisid/platform/services/entity-resolution/internal/matching"
	"github.com/snisid/platform/services/entity-resolution/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	db.AutoMigrate(&models.Identity{}, &models.ResolvedIdentity{})
	return db
}

func seedIdentities(db *gorm.DB) {
	now := time.Now()
	db.Create(&models.Identity{
		ID: "id-001", NNU: "NNU-001", FirstName: "Jean", LastName: "Dupont",
		FullName: "Jean Dupont", DOB: "1990-01-15", TaxID: "TAX-001",
		Status: "active", CreatedAt: now, UpdatedAt: now,
	})
	db.Create(&models.Identity{
		ID: "id-002", NNU: "NNU-002", FirstName: "Marie", LastName: "Pierre",
		FullName: "Marie Pierre", DOB: "1985-06-20", TaxID: "TAX-002",
		Status: "active", CreatedAt: now, UpdatedAt: now,
	})
	db.Create(&models.Identity{
		ID: "id-003", NNU: "NNU-003", FirstName: "Jean", LastName: "Dupont",
		FullName: "Jean Dupont", DOB: "1990-01-15", TaxID: "TAX-003",
		BiometricHash: "hash-same", Status: "active",
		CreatedAt: now, UpdatedAt: now,
	})
}

func TestMatchHandler_ExactMatch(t *testing.T) {
	db := setupTestDB(t)
	seedIdentities(db)

	lsh := matching.NewLSHIndex(100, 20, 5)
	var identities []models.Identity
	db.Find(&identities)
	lsh.Build(matching.IdentitiesToIndexable(identities))

	engine := matching.NewCompositeEngine(db, lsh)
	h := handlers.NewResolutionHandler(engine, db)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/resolution/match", h.MatchHandler)

	body := `{"nnu":"NNU-001","first_name":"Jean","last_name":"Dupont"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/resolution/match", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.GreaterOrEqual(t, resp["count"].(float64), float64(1))
}

func TestMatchHandler_FuzzyMatch(t *testing.T) {
	db := setupTestDB(t)
	seedIdentities(db)

	lsh := matching.NewLSHIndex(100, 20, 5)
	var identities []models.Identity
	db.Find(&identities)
	lsh.Build(matching.IdentitiesToIndexable(identities))

	engine := matching.NewCompositeEngine(db, lsh)
	h := handlers.NewResolutionHandler(engine, db)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/resolution/match", h.MatchHandler)

	body := `{"first_name":"Jea","last_name":"Dupont","dob":"1990-01-15"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/resolution/match", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.GreaterOrEqual(t, resp["count"].(float64), float64(1))
}

func TestMatchHandler_NoMatch(t *testing.T) {
	db := setupTestDB(t)
	seedIdentities(db)

	lsh := matching.NewLSHIndex(100, 20, 5)
	var identities []models.Identity
	db.Find(&identities)
	lsh.Build(matching.IdentitiesToIndexable(identities))

	engine := matching.NewCompositeEngine(db, lsh)
	h := handlers.NewResolutionHandler(engine, db)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/resolution/match", h.MatchHandler)

	body := `{"first_name":"Unknown","last_name":"Person","dob":"2000-01-01"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/resolution/match", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["count"].(float64))
}

func TestReconcileHandler_Success(t *testing.T) {
	db := setupTestDB(t)
	seedIdentities(db)

	lsh := matching.NewLSHIndex(100, 20, 5)
	engine := matching.NewCompositeEngine(db, lsh)
	h := handlers.NewResolutionHandler(engine, db)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/resolution/reconcile", h.ReconcileHandler)

	body := `{"primary_id":"id-001","secondary_id":"id-003"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/resolution/reconcile", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result models.ReconciliationResult
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "id-001", result.PrimaryID)
	assert.Equal(t, "id-003", result.SecondaryID)
	assert.Greater(t, result.OverallScore, float64(0))
	assert.NotEmpty(t, result.Decision)
}

func TestReconcileHandler_NotFound(t *testing.T) {
	db := setupTestDB(t)

	lsh := matching.NewLSHIndex(100, 20, 5)
	engine := matching.NewCompositeEngine(db, lsh)
	h := handlers.NewResolutionHandler(engine, db)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/resolution/reconcile", h.ReconcileHandler)

	body := `{"primary_id":"nonexistent","secondary_id":"id-002"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/resolution/reconcile", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMergeHandler(t *testing.T) {
	db := setupTestDB(t)
	seedIdentities(db)

	lsh := matching.NewLSHIndex(100, 20, 5)
	engine := matching.NewCompositeEngine(db, lsh)
	h := handlers.NewResolutionHandler(engine, db)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/resolution/merge", h.MergeHandler)

	body := `{"primary_id":"id-001","secondary_id":"id-002","resolved_by":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/resolution/merge", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["message"].(string), "merged")

	var merged models.Identity
	db.First(&merged, "id = ?", "id-002")
	assert.Equal(t, "merged", merged.Status)
}

func TestCandidatesHandler(t *testing.T) {
	db := setupTestDB(t)
	seedIdentities(db)

	lsh := matching.NewLSHIndex(100, 20, 5)
	var identities []models.Identity
	db.Find(&identities)
	lsh.Build(matching.IdentitiesToIndexable(identities))

	engine := matching.NewCompositeEngine(db, lsh)
	h := handlers.NewResolutionHandler(engine, db)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/v1/resolution/candidates/:id", h.CandidatesHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/resolution/candidates/id-001", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "id-001", resp["identity_id"].(string))
}

func TestStatsHandler(t *testing.T) {
	db := setupTestDB(t)
	seedIdentities(db)

	lsh := matching.NewLSHIndex(100, 20, 5)
	engine := matching.NewCompositeEngine(db, lsh)
	h := handlers.NewResolutionHandler(engine, db)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/v1/resolution/stats", h.StatsHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/resolution/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var stats models.StatsResponse
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	require.NoError(t, err)
	assert.Equal(t, int64(3), stats.TotalIdentities)
}

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "ok", resp["status"])
}
