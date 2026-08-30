package router

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/snisid/platform/internal/platform/logger"
	"go.uber.org/zap"
)

// Rule represents a routing condition and its target topics
type Rule struct {
	ID         string
	Expression string
	Targets    []string
	Program    cel.Program
}

// Engine evaluates events against a set of CEL rules
type Engine struct {
	env   *cel.Env
	rules []Rule
}

func NewEngine() (*Engine, error) {
	// Setup CEL environment with 'event' variable
	env, err := cel.NewEnv(
		cel.Variable("event", cel.MapType(cel.StringType, cel.DynType)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create cel env: %w", err)
	}

	return &Engine{
		env: env,
	}, nil
}

// UpdateRules compiles and loads new rules into the engine
func (e *Engine) UpdateRules(newRules []Rule) error {
	var compiledRules []Rule
	for _, r := range newRules {
		ast, iss := e.env.Compile(r.Expression)
		if iss.Err() != nil {
			return fmt.Errorf("failed to compile rule %s: %v", r.ID, iss.Err())
		}

		prg, err := e.env.Program(ast)
		if err != nil {
			return fmt.Errorf("failed to create program for rule %s: %w", r.ID, err)
		}

		r.Program = prg
		compiledRules = append(compiledRules, r)
	}

	e.rules = compiledRules
	return nil
}

// Rules returns the currently loaded rules
func (e *Engine) Rules() []Rule {
	return e.rules
}

// Evaluate checks an event against all rules and returns the combined list of target topics
func (e *Engine) Evaluate(ctx context.Context, event map[string]interface{}) []string {
	targetSet := make(map[string]struct{})

	for _, r := range e.rules {
		out, _, err := r.Program.Eval(map[string]interface{}{
			"event": event,
		})
		if err != nil {
			logger.Error(ctx, "Rule evaluation error", err, zap.String("rule_id", r.ID))
			continue
		}

		if out == types.True {
			for _, t := range r.Targets {
				targetSet[t] = struct{}{}
			}
		}
	}

	var targets []string
	for t := range targetSet {
		targets = append(targets, t)
	}
	return targets
}
