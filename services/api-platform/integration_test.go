package apiplatform

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestManager() *TenantManager {
	return &TenantManager{
		Tenants: map[string]Tenant{
			"key-valid-oni": {
				Country:     "HTI",
				APIKey:      "key-valid-oni",
				Permissions: []Permission{PermVerify, PermScore},
			},
			"key-valid-biometric": {
				Country:     "DOM",
				APIKey:      "key-valid-biometric",
				Permissions: []Permission{PermVerify, PermScore, PermBiocheck},
			},
			"key-readonly": {
				Country:     "HTI",
				APIKey:      "key-readonly",
				Permissions: []Permission{PermVerify},
			},
		},
	}
}

func TestValidateAccess_ValidKey(t *testing.T) {
	m := setupTestManager()
	ok, reason := m.ValidateAccess("key-valid-oni", PermVerify)
	assert.True(t, ok)
	assert.Empty(t, reason)
}

func TestValidateAccess_InvalidKey(t *testing.T) {
	m := setupTestManager()
	ok, reason := m.ValidateAccess("nonexistent-key", PermVerify)
	assert.False(t, ok)
	assert.Equal(t, "INVALID_API_KEY", reason)
}

func TestValidateAccess_InsufficientPermissions(t *testing.T) {
	m := setupTestManager()
	ok, reason := m.ValidateAccess("key-readonly", PermBiocheck)
	assert.False(t, ok)
	assert.Equal(t, "INSUFFICIENT_PERMISSIONS", reason)
}

func TestValidateAccess_MultiplePermissions(t *testing.T) {
	m := setupTestManager()
	ok, reason := m.ValidateAccess("key-valid-biometric", PermBiocheck)
	assert.True(t, ok)
	assert.Empty(t, reason)
	ok, reason = m.ValidateAccess("key-valid-biometric", PermScore)
	assert.True(t, ok)
	assert.Empty(t, reason)
}

func TestTenantManager_EmptyTenants(t *testing.T) {
	m := &TenantManager{Tenants: make(map[string]Tenant)}
	ok, reason := m.ValidateAccess("any-key", PermVerify)
	assert.False(t, ok)
	assert.Equal(t, "INVALID_API_KEY", reason)
}

func TestValidateHTTP_Success(t *testing.T) {
	m := setupTestManager()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/validate", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			APIKey   string     `json:"api_key"`
			Required Permission `json:"required"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok, reason := m.ValidateAccess(req.APIKey, req.Required)
		json.NewEncoder(w).Encode(map[string]interface{}{"allowed": ok, "reason": reason})
	})

	body := `{"api_key":"key-valid-oni","required":"CITIZEN_VERIFY"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/validate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["allowed"].(bool))
	assert.Empty(t, resp["reason"])
}

func TestValidateHTTP_InvalidKey(t *testing.T) {
	m := setupTestManager()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/validate", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			APIKey   string     `json:"api_key"`
			Required Permission `json:"required"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok, reason := m.ValidateAccess(req.APIKey, req.Required)
		json.NewEncoder(w).Encode(map[string]interface{}{"allowed": ok, "reason": reason})
	})

	body := `{"api_key":"bad-key","required":"FRAUD_SCORE"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/validate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp["allowed"].(bool))
	assert.Equal(t, "INVALID_API_KEY", resp["reason"])
}

func TestValidateHTTP_BadRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/validate", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			APIKey   string     `json:"api_key"`
			Required Permission `json:"required"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	})

	body := `{invalid json`
	req := httptest.NewRequest(http.MethodPost, "/v1/validate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHealthEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "ok", resp["status"])
}

func TestConcurrentAccessValidation(t *testing.T) {
	m := setupTestManager()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, reason := m.ValidateAccess("key-valid-oni", PermVerify)
			assert.True(t, ok)
			assert.Empty(t, reason)
		}()
	}
	wg.Wait()
}

func TestPermConstants(t *testing.T) {
	assert.Equal(t, Permission("CITIZEN_VERIFY"), PermVerify)
	assert.Equal(t, Permission("FRAUD_SCORE"), PermScore)
	assert.Equal(t, Permission("BIOMETRIC_CHECK"), PermBiocheck)
}
