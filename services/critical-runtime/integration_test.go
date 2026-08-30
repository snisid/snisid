package criticalruntimetest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/snisid/platform/services/critical-runtime/healer"
	"github.com/snisid/platform/services/critical-runtime/monitor"
	"github.com/snisid/platform/services/critical-runtime/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckInvariant_Pass(t *testing.T) {
	checker := &monitor.RuntimeChecker{ID: "test-runtime"}
	state := monitor.SystemState{
		RiskVector: map[string]int{"node-a": 30, "node-b": 50},
		Threshold:  80,
		Policies:   map[string]string{"node-a": "ALLOW", "node-b": "ALLOW"},
	}

	valid, msg := checker.CheckInvariant(state)
	assert.True(t, valid)
	assert.Equal(t, "PASS", msg)
}

func TestCheckInvariant_Violation(t *testing.T) {
	checker := &monitor.RuntimeChecker{ID: "test-runtime"}
	state := monitor.SystemState{
		RiskVector: map[string]int{"node-a": 90, "node-b": 20},
		Threshold:  80,
		Policies:   map[string]string{"node-a": "ALLOW", "node-b": "ALLOW"},
	}

	valid, msg := checker.CheckInvariant(state)
	assert.False(t, valid)
	assert.Contains(t, msg, "INVARIANT_VIOLATION")
	assert.Contains(t, msg, "node-a")
}

func TestCheckInvariant_ViolationSuppressedByPolicy(t *testing.T) {
	checker := &monitor.RuntimeChecker{ID: "test-runtime"}
	state := monitor.SystemState{
		RiskVector: map[string]int{"node-a": 90},
		Threshold:  80,
		Policies:   map[string]string{"node-a": "DENY"},
	}

	valid, msg := checker.CheckInvariant(state)
	assert.True(t, valid)
	assert.Equal(t, "PASS", msg)
}

func TestHealer_CriticalPlan(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	h := healer.NewHealer("test-platform", store)

	violation := healer.Violation{
		ID: "vio-001", Type: "SECURITY", Description: "Critical breach",
		Severity: healer.SeverityCritical, Affected: []string{"domain-1", "domain-2"},
		DetectedAt: time.Now(), Source: "monitor",
	}

	plan := h.Heal(violation)
	require.NotNil(t, plan)
	assert.Equal(t, healer.SeverityCritical, plan.Violation.Severity)
	assert.Len(t, plan.Steps, 5)
	assert.Equal(t, "ISOLATE", plan.Steps[0].Action)
	assert.Equal(t, "affected_domains", plan.Steps[0].Target)
}

func TestHealer_HighPlan(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	h := healer.NewHealer("test-platform", store)

	violation := healer.Violation{
		ID: "vio-002", Type: "PERFORMANCE", Severity: healer.SeverityHigh,
		Affected: []string{"api-gateway"}, DetectedAt: time.Now(),
	}

	plan := h.Heal(violation)
	require.NotNil(t, plan)
	assert.Len(t, plan.Steps, 3)
	assert.Equal(t, "ISOLATE", plan.Steps[0].Action)
	assert.Equal(t, "api-gateway", plan.Steps[0].Target)
}

func TestHealer_MediumPlan(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	h := healer.NewHealer("test-platform", store)

	violation := healer.Violation{
		ID: "vio-003", Type: "AVAILABILITY", Severity: healer.SeverityMedium,
		Affected: []string{"db-replica"}, DetectedAt: time.Now(),
	}

	plan := h.Heal(violation)
	require.NotNil(t, plan)
	assert.Len(t, plan.Steps, 2)
	assert.Equal(t, "RESTART", plan.Steps[0].Action)
}

func TestHealer_LowPlan(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	h := healer.NewHealer("test-platform", store)

	violation := healer.Violation{
		ID: "vio-004", Type: "INTEGRITY", Severity: healer.SeverityLow,
		Affected: []string{"config-service"}, DetectedAt: time.Now(),
	}

	plan := h.Heal(violation)
	require.NotNil(t, plan)
	assert.Len(t, plan.Steps, 1)
	assert.Equal(t, "RECONFIGURE", plan.Steps[0].Action)
}

func TestHealer_RateLimiting(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	h := healer.NewHealer("test-platform", store)

	v := healer.Violation{ID: "vio-rate", Type: "SECURITY", Severity: healer.SeverityHigh,
		Affected: []string{"node"}, DetectedAt: time.Now()}

	for i := 0; i < 3; i++ {
		h.Heal(v)
	}

	plan := h.Heal(v)
	assert.Nil(t, plan, "expected rate-limited plan to be nil")
}

func TestHealer_GetActiveHealings(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	h := healer.NewHealer("test-platform", store)

	v := healer.Violation{ID: "vio-active", Type: "AVAILABILITY", Severity: healer.SeverityLow,
		Affected: []string{"svc"}, DetectedAt: time.Now()}
	h.Heal(v)

	active := h.GetActiveHealings()
	assert.Contains(t, active, "vio-active")
	assert.Equal(t, "RECONFIGURE", active["vio-active"].Steps[0].Action)
}

func TestHealer_ResumeClearsActive(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	h := healer.NewHealer("test-platform", store)

	h.Heal(healer.Violation{ID: "vio-resume", Type: "INTEGRITY", Severity: healer.SeverityLow,
		Affected: []string{"svc"}, DetectedAt: time.Now()})

	h.Resume()
	assert.Empty(t, h.GetActiveHealings())
}

func TestSnapshotStore_SaveAndGetLastValid(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	state := snapshot.ValidState{
		Timestamp: 1000,
		RiskData:  map[string]int{"node-a": 30},
		PolicySet: map[string]string{"node-a": "ALLOW"},
	}
	store.Save(state)

	last := store.GetLastValid()
	assert.Equal(t, int64(1000), last.Timestamp)
	assert.Equal(t, 30, last.RiskData["node-a"])
}

func TestSnapshotStore_MultipleSnapshots(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	store.Save(snapshot.ValidState{Timestamp: 1})
	store.Save(snapshot.ValidState{Timestamp: 2})
	store.Save(snapshot.ValidState{Timestamp: 3})

	last := store.GetLastValid()
	assert.Equal(t, int64(3), last.Timestamp)
	assert.Len(t, store.History, 3)
}

func TestSnapshotStore_EmptyReturnsZero(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	last := store.GetLastValid()
	assert.Equal(t, int64(0), last.Timestamp)
	assert.Nil(t, last.RiskData)
}

func TestCheckHTTPEndpoint_ValidState(t *testing.T) {
	checker := &monitor.RuntimeChecker{ID: "test-runtime"}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/check", func(w http.ResponseWriter, r *http.Request) {
		var state monitor.SystemState
		if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok, msg := checker.CheckInvariant(state)
		json.NewEncoder(w).Encode(map[string]interface{}{"valid": ok, "message": msg})
	})

	body := `{"risk_vector":{"node-a":30},"threshold":80,"policies":{"node-a":"ALLOW"}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/check", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["valid"].(bool))
	assert.Equal(t, "PASS", resp["message"])
}

func TestCheckHTTPEndpoint_Violation(t *testing.T) {
	checker := &monitor.RuntimeChecker{ID: "test-runtime"}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/check", func(w http.ResponseWriter, r *http.Request) {
		var state monitor.SystemState
		if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok, msg := checker.CheckInvariant(state)
		json.NewEncoder(w).Encode(map[string]interface{}{"valid": ok, "message": msg})
	})

	body := `{"risk_vector":{"node-a":90},"threshold":80,"policies":{"node-a":"ALLOW"}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/check", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp["valid"].(bool))
	assert.Contains(t, resp["message"].(string), "INVARIANT_VIOLATION")
}

func TestCheckHTTPEndpoint_BadRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/check", func(w http.ResponseWriter, r *http.Request) {
		var state monitor.SystemState
		if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	})

	body := `{invalid`
	req := httptest.NewRequest(http.MethodPost, "/v1/check", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHealHTTPEndpoint(t *testing.T) {
	store := &snapshot.SnapshotStore{}
	h := healer.NewHealer("test-platform", store)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/heal", func(w http.ResponseWriter, r *http.Request) {
		var v healer.Violation
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		plan := h.Heal(v)
		json.NewEncoder(w).Encode(plan)
	})

	body := `{"id":"http-heal","type":"PERFORMANCE","severity":3,"affected":["web-svc"]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/heal", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var plan healer.HealingPlan
	err := json.Unmarshal(w.Body.Bytes(), &plan)
	assert.NoError(t, err)
	assert.Equal(t, "http-heal", plan.Violation.ID)
	assert.GreaterOrEqual(t, len(plan.Steps), 1)
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

func TestViolationSeverityConstants(t *testing.T) {
	assert.Equal(t, healer.ViolationSeverity(1), healer.SeverityLow)
	assert.Equal(t, healer.ViolationSeverity(2), healer.SeverityMedium)
	assert.Equal(t, healer.ViolationSeverity(3), healer.SeverityHigh)
	assert.Equal(t, healer.ViolationSeverity(4), healer.SeverityCritical)
}
