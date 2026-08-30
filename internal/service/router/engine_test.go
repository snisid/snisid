package router

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEngine(t *testing.T) {
	engine, err := NewEngine()
	require.NoError(t, err)
	assert.NotNil(t, engine)
	assert.NotNil(t, engine.env)
}

func TestUpdateRules_ValidExpression(t *testing.T) {
	engine, err := NewEngine()
	require.NoError(t, err)

	rules := []Rule{
		{ID: "high-value", Expression: "event.amount > 10000", Targets: []string{"compliance", "soc"}},
		{ID: "suspicious-country", Expression: "event.country == 'untrusted'", Targets: []string{"soc"}},
	}

	err = engine.UpdateRules(rules)
	require.NoError(t, err)
	assert.Len(t, engine.rules, 2)
}

func TestUpdateRules_InvalidExpression(t *testing.T) {
	engine, err := NewEngine()
	require.NoError(t, err)

	rules := []Rule{
		{ID: "bad-rule", Expression: "event >>> invalid", Targets: []string{"soc"}},
	}

	err = engine.UpdateRules(rules)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to compile rule")
}

func TestEvaluate_MatchesRule(t *testing.T) {
	engine, err := NewEngine()
	require.NoError(t, err)

	err = engine.UpdateRules([]Rule{
		{ID: "high-amount", Expression: "event.amount > 5000", Targets: []string{"fraud"}},
	})
	require.NoError(t, err)

	event := map[string]interface{}{
		"amount": 15000,
		"type":   "transfer",
	}

	targets := engine.Evaluate(context.Background(), event)
	assert.Len(t, targets, 1)
	assert.Equal(t, "fraud", targets[0])
}

func TestEvaluate_NoMatch(t *testing.T) {
	engine, err := NewEngine()
	require.NoError(t, err)

	err = engine.UpdateRules([]Rule{
		{ID: "high-amount", Expression: "event.amount > 5000", Targets: []string{"fraud"}},
	})
	require.NoError(t, err)

	event := map[string]interface{}{
		"amount": 100,
		"type":   "transfer",
	}

	targets := engine.Evaluate(context.Background(), event)
	assert.Empty(t, targets)
}

func TestEvaluate_MultipleTargets(t *testing.T) {
	engine, err := NewEngine()
	require.NoError(t, err)

	err = engine.UpdateRules([]Rule{
		{ID: "critical", Expression: "event.severity == 'critical'", Targets: []string{"soc", "compliance", "admin"}},
	})
	require.NoError(t, err)

	event := map[string]interface{}{
		"severity": "critical",
	}

	targets := engine.Evaluate(context.Background(), event)
	assert.Len(t, targets, 3)
	assert.Contains(t, targets, "soc")
	assert.Contains(t, targets, "compliance")
	assert.Contains(t, targets, "admin")
}

func TestEvaluate_EmptyRules(t *testing.T) {
	engine, err := NewEngine()
	require.NoError(t, err)

	event := map[string]interface{}{"type": "test"}
	targets := engine.Evaluate(context.Background(), event)
	assert.Empty(t, targets)
}

func TestEvaluate_NoDuplicateTargets(t *testing.T) {
	engine, err := NewEngine()
	require.NoError(t, err)

	err = engine.UpdateRules([]Rule{
		{ID: "rule1", Expression: "event.type == 'alert'", Targets: []string{"soc"}},
		{ID: "rule2", Expression: "event.priority > 5", Targets: []string{"soc"}},
	})
	require.NoError(t, err)

	event := map[string]interface{}{
		"type":     "alert",
		"priority": 10,
	}

	targets := engine.Evaluate(context.Background(), event)
	assert.Len(t, targets, 1)
	assert.Equal(t, "soc", targets[0])
}
