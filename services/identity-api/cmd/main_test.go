package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/snisid/platform/services/identity-api/internal/handlers"
	"github.com/snisid/platform/services/identity-api/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.Identity{}, &models.IdentityHistory{})
	require.NoError(t, err)
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := handlers.New(db, nil)
	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("actor_id", "test-user")
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

func TestCreateIdentity_Success(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	payload := map[string]string{
		"first_name": "Marie",
		"last_name":  "Pierre",
		"gender":     "F",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/identities", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var result models.Identity
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.NNU)
	assert.Equal(t, "Marie", result.FirstName)
	assert.Equal(t, "Pierre", result.LastName)
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, 1, result.Version)

	var count int64
	db.Model(&models.Identity{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestCreateIdentity_ValidationError(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	payload := map[string]string{"first_name": ""}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/identities", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var errResp map[string]string
	json.Unmarshal(resp.Body.Bytes(), &errResp)
	assert.Contains(t, errResp["error"], "first_name")
}

func TestGetIdentity_Found(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	ident := models.Identity{
		ID:        "test-id-1",
		NNU:       "TES250000001",
		FirstName: "Jean",
		LastName:  "Dupont",
		Status:    "active",
		Version:   1,
	}
	db.Create(&ident)

	req, _ := http.NewRequest("GET", "/api/v1/identities/test-id-1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result models.Identity
	json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, "Jean", result.FirstName)
	assert.Equal(t, "Dupont", result.LastName)
}

func TestGetIdentity_NotFound(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	req, _ := http.NewRequest("GET", "/api/v1/identities/non-existent", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestListIdentities_Empty(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	req, _ := http.NewRequest("GET", "/api/v1/identities", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]any
	json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, float64(0), result["total"])
	assert.Equal(t, float64(1), result["page"])
	assert.Equal(t, float64(20), result["limit"])

	data := result["data"].([]any)
	assert.Empty(t, data)
}

func TestListIdentities_WithPagination(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("id-%d", i)
		db.Create(&models.Identity{
			ID: id, NNU: fmt.Sprintf("TES25%05X", i),
			FirstName: "User", LastName: fmt.Sprintf("Num%d", i),
			Status: "pending", Version: 1,
		})
	}

	req, _ := http.NewRequest("GET", "/api/v1/identities?page=1&limit=2", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]any
	json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, float64(5), result["total"])
	assert.Equal(t, float64(1), result["page"])
	assert.Equal(t, float64(2), result["limit"])

	data := result["data"].([]any)
	assert.Len(t, data, 2)
}

func TestListIdentities_StatusFilter(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	db.Create(&models.Identity{ID: "id-1", NNU: "TES25A0001", FirstName: "A", LastName: "One", Status: "active", Version: 1})
	db.Create(&models.Identity{ID: "id-2", NNU: "TES25A0002", FirstName: "B", LastName: "Two", Status: "pending", Version: 1})
	db.Create(&models.Identity{ID: "id-3", NNU: "TES25A0003", FirstName: "C", LastName: "Three", Status: "suspended", Version: 1})

	req, _ := http.NewRequest("GET", "/api/v1/identities?status=active", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]any
	json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, float64(1), result["total"])

	data := result["data"].([]any)
	assert.Len(t, data, 1)
	assert.Equal(t, "active", data[0].(map[string]any)["status"])
}

func TestUpdateIdentity_Success(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	db.Create(&models.Identity{
		ID: "update-id", NNU: "TES25U0001",
		FirstName: "OldName", LastName: "User",
		Status: "active", Version: 1,
	})

	payload := map[string]string{"first_name": "NewName"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/api/v1/identities/update-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result models.Identity
	json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, "NewName", result.FirstName)
	assert.Equal(t, 2, result.Version)

	var historyCount int64
	db.Model(&models.IdentityHistory{}).Where("identity_id = ?", "update-id").Count(&historyCount)
	assert.Equal(t, int64(1), historyCount)
}

func TestUpdateIdentity_NotFound(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	payload := map[string]string{"first_name": "NewName"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/api/v1/identities/non-existent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestDeleteIdentity_Success(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	db.Create(&models.Identity{
		ID: "delete-id", NNU: "TES25D0001",
		FirstName: "Delete", LastName: "Me",
		Status: "active", Version: 1,
	})

	req, _ := http.NewRequest("DELETE", "/api/v1/identities/delete-id", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]string
	json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, "suspended", result["status"])

	var ident models.Identity
	db.First(&ident, "id = ?", "delete-id")
	assert.Equal(t, "suspended", ident.Status)
}

func TestDeleteIdentity_NotFound(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	req, _ := http.NewRequest("DELETE", "/api/v1/identities/non-existent", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestFlagIdentity_Success(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	db.Create(&models.Identity{
		ID: "flag-id", NNU: "TES25F0001",
		FirstName: "Flag", LastName: "User",
		Status: "active", Version: 1,
	})

	req, _ := http.NewRequest("POST", "/api/v1/identities/flag/flag-id", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]string
	json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, "suspended", result["status"])

	var ident models.Identity
	db.First(&ident, "id = ?", "flag-id")
	assert.Equal(t, "suspended", ident.Status)
	assert.Equal(t, 2, ident.Version)
}

func TestFlagIdentity_NotFound(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	req, _ := http.NewRequest("POST", "/api/v1/identities/flag/non-existent", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req, _ := http.NewRequest("GET", "/healthz", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]string
	json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, "ok", result["status"])
}

func TestGetEnv_Default(t *testing.T) {
	result := getEnv("NONEXISTENT_VAR_12345", "default-val")
	assert.Equal(t, "default-val", result)
}

func TestGetEnv_FromEnv(t *testing.T) {
	os.Setenv("TEST_GETENV_KEY", "env-value")
	defer os.Unsetenv("TEST_GETENV_KEY")
	result := getEnv("TEST_GETENV_KEY", "default-val")
	assert.Equal(t, "env-value", result)
}
