package compiler

import (
	"fmt"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

type WorldState struct {
	InfraHealth    float64
	AttackPressure float64
	Load           float64
}

type AIInfraCompiler struct{}

func NewAIInfraCompiler() *AIInfraCompiler {
	return &AIInfraCompiler{}
}

func (c *AIInfraCompiler) Compile(state WorldState) (string, error) {
	logger.Info("AI Compiler: analyzing world state for patch generation")

	if state.InfraHealth < 0.7 {
		return c.generateHealingPlan(), nil
	}

	if state.Load > 0.8 {
		return c.generateScalingPlan(), nil
	}

	return "", nil
}

func (c *AIInfraCompiler) generateHealingPlan() string {
	return `
# AI Generated Healing Plan
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: snisid-backend
spec:
  values:
    replicaCount: 5
    resources:
      limits:
        cpu: "2000m"
`
}

func (c *AIInfraCompiler) generateScalingPlan() string {
	return `
# AI Generated Scaling Plan
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: snisid-backend-hpa
spec:
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 60
`
}
