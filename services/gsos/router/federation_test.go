package router

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFederationRouter(t *testing.T) {
	nodes := []string{"HTI", "DOM", "CUB"}
	r := NewFederationRouter(nodes)
	require.NotNil(t, r)
	assert.Equal(t, nodes, r.ActiveNodes)
	assert.Empty(t, r.alertHistory)
	assert.True(t, r.isNodeOnline("HTI"))
}

func TestNewFederationRouter_EmptyNodes(t *testing.T) {
	r := NewFederationRouter([]string{})
	assert.Empty(t, r.ActiveNodes)
}

func TestComputeAlertLevel(t *testing.T) {
	r := NewFederationRouter(nil)

	tests := []struct {
		risk  float64
		level AlertLevel
	}{
		{0.96, AlertGlobal},
		{0.95, AlertGlobal},
		{0.90, AlertCritical},
		{0.85, AlertCritical},
		{0.80, AlertWarning},
		{0.70, AlertWarning},
		{0.50, AlertInfo},
		{0.0, AlertInfo},
	}

	for _, tt := range tests {
		event := GSPEvent{Risk: tt.risk}
		assert.Equal(t, tt.level, r.computeAlertLevel(event))
	}
}

func TestRouteEvent_LowRisk(t *testing.T) {
	r := NewFederationRouter([]string{"HTI", "DOM"})
	event := GSPEvent{
		ID: "EVT-001", Country: "HTI", Category: "login",
		Risk: 0.3, Timestamp: time.Now().Unix(),
	}
	r.RouteEvent(event)

	assert.Len(t, r.alertHistory, 1)
	assert.Equal(t, AlertInfo, r.alertHistory[0].Severity)
}

func TestRouteEvent_CriticalTriggersPropagation(t *testing.T) {
	r := NewFederationRouter([]string{"HTI", "DOM", "CUB"})
	propagated := false
	r.SetPropagationFn(func(alert GlobalAlert, target string) error {
		propagated = true
		return nil
	})

	event := GSPEvent{
		ID: "EVT-CRIT", Country: "HTI", Category: "attack",
		Risk: 0.9, Timestamp: time.Now().Unix(),
	}
	r.RouteEvent(event)

	assert.True(t, propagated)
	assert.Len(t, r.alertHistory, 1)
}

func TestRouteEvent_GlobalAlert(t *testing.T) {
	r := NewFederationRouter([]string{"HTI", "DOM", "CUB", "JAM"})
	var receivedAlert GlobalAlert
	r.SetPropagationFn(func(alert GlobalAlert, target string) error {
		receivedAlert = alert
		return nil
	})

	event := GSPEvent{
		ID: "EVT-GLOBAL", Country: "HTI", Risk: 0.96, Timestamp: time.Now().Unix(),
	}
	r.RouteEvent(event)

	assert.Equal(t, AlertGlobal, receivedAlert.Severity)
	assert.Len(t, receivedAlert.AffectedCountries, 3)
}

func TestRouteEvent_AlertHistoryCap(t *testing.T) {
	r := NewFederationRouter([]string{"HTI"})
	for i := 0; i < 1500; i++ {
		r.RouteEvent(GSPEvent{ID: "EVT", Risk: 0.1})
	}
	r.mu.RLock()
	assert.LessOrEqual(t, len(r.alertHistory), 1000)
	r.mu.RUnlock()
}

func TestPropagateGlobalAlert_AcksReceived(t *testing.T) {
	r := NewFederationRouter([]string{"HTI", "DOM", "CUB"})
	ackCount := 0
	r.SetPropagationFn(func(alert GlobalAlert, target string) error {
		ackCount++
		return nil
	})

	alert := GlobalAlert{
		ID:                "GA-TEST",
		Severity:          AlertCritical,
		AffectedCountries: []string{"DOM", "CUB"},
	}
	r.PropagateGlobalAlert(alert)

	assert.Equal(t, 2, ackCount)
}

func TestPropagateGlobalAlert_PropagationFailure(t *testing.T) {
	r := NewFederationRouter([]string{"HTI", "DOM"})
	r.SetPropagationFn(func(alert GlobalAlert, target string) error {
		return errors.New("network error")
	})

	alert := GlobalAlert{
		ID:                "GA-FAIL",
		AffectedCountries: []string{"DOM"},
	}
	r.PropagateGlobalAlert(alert)

	// Should have 0 acks
	assert.Equal(t, 0, alert.AcksReceived)
}

func TestPropagateGlobalAlert_OfflineNode(t *testing.T) {
	r := NewFederationRouter([]string{"HTI", "DOM", "CUB"})
	r.SetNodeStatus("DOM", false)

	ackCount := 0
	r.SetPropagationFn(func(alert GlobalAlert, target string) error {
		ackCount++
		return nil
	})

	alert := GlobalAlert{
		ID:                "GA-OFF",
		AffectedCountries: []string{"HTI", "DOM", "CUB"},
	}
	r.PropagateGlobalAlert(alert)

	assert.Equal(t, 2, ackCount)
}

func TestSetNodeStatus(t *testing.T) {
	r := NewFederationRouter([]string{"HTI"})
	assert.True(t, r.isNodeOnline("HTI"))

	r.SetNodeStatus("HTI", false)
	assert.False(t, r.isNodeOnline("HTI"))

	r.SetNodeStatus("HTI", true)
	assert.True(t, r.isNodeOnline("HTI"))
}

func TestSyncLocalPolicies_Cached(t *testing.T) {
	r := NewFederationRouter([]string{"HTI"})
	r.policyCache["HTI"] = []string{"CUSTOM_POLICY"}

	policies := r.SyncLocalPolicies("HTI")
	assert.Equal(t, []string{"CUSTOM_POLICY"}, policies)
}

func TestSyncLocalPolicies_Defaults(t *testing.T) {
	r := NewFederationRouter([]string{"HTI"})
	policies := r.SyncLocalPolicies("DOM")
	require.Len(t, policies, 5)
	assert.Contains(t, policies, "EVENT_SHARING_AGREEMENT")
	assert.Contains(t, policies, "FRAUD_ALERT_PROPAGATION")
}

func TestGetAlertHistory_FilterBySince(t *testing.T) {
	r := NewFederationRouter([]string{"HTI"})
	now := time.Now().Unix()

	oldEvent := GSPEvent{ID: "old", Risk: 0.3, Timestamp: now - 3600}
	newEvent := GSPEvent{ID: "new", Risk: 0.3, Timestamp: now}

	r.RouteEvent(oldEvent)
	r.RouteEvent(newEvent)

	all := r.GetAlertHistory(0)
	assert.Len(t, all, 2)

	recent := r.GetAlertHistory(now - 1800)
	assert.Len(t, recent, 1)
	assert.Equal(t, "new", recent[0].ID)
}

func TestGetAlertHistory_SortedByPropagatedAt(t *testing.T) {
	r := NewFederationRouter([]string{"HTI"})
	r.RouteEvent(GSPEvent{ID: "evt1", Risk: 0.3, Timestamp: 100})
	r.RouteEvent(GSPEvent{ID: "evt2", Risk: 0.3, Timestamp: 200})

	history := r.GetAlertHistory(0)
	require.Len(t, history, 2)
	assert.True(t, history[0].PropagatedAt >= history[1].PropagatedAt)
}

func TestConcurrentRouteEvent(t *testing.T) {
	r := NewFederationRouter([]string{"HTI", "DOM", "CUB"})
	t.Run("parallel", func(t *testing.T) {
		t.Run("event1", func(t *testing.T) {
			r.RouteEvent(GSPEvent{ID: "e1", Risk: 0.5, Timestamp: 1})
		})
		t.Run("event2", func(t *testing.T) {
			r.RouteEvent(GSPEvent{ID: "e2", Risk: 0.8, Timestamp: 2})
		})
		t.Run("event3", func(t *testing.T) {
			r.RouteEvent(GSPEvent{ID: "e3", Risk: 0.9, Timestamp: 3})
		})
	})
	r.mu.RLock()
	assert.Len(t, r.alertHistory, 3)
	r.mu.RUnlock()
}

func TestComputeAffectedCountries(t *testing.T) {
	r := NewFederationRouter([]string{"HTI", "DOM", "CUB"})
	affected := r.computeAffectedCountries(GSPEvent{Country: "HTI"})
	assert.Len(t, affected, 2)
	assert.NotContains(t, affected, "HTI")
	assert.Contains(t, affected, "DOM")
}
