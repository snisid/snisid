package fraud

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ModelResult struct {
	Score  int
	Reason string
}

type Model interface {
	Name() string
	Score(ctx context.Context, event map[string]interface{}) (ModelResult, error)
}

type FeatureVector struct {
	UserID    string
	Amount    float64
	Velocity  float64
	GraphRisk float64
}

type MLModel interface {
	Predict(ctx context.Context, features FeatureVector) (float64, error)
}

type GRPCMLModel struct {
	conn    *grpc.ClientConn
	timeout time.Duration
}

func NewGRPCMLModel(endpoint string, timeout time.Duration) (*GRPCMLModel, error) {
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial %s: %w", endpoint, err)
	}
	return &GRPCMLModel{
		conn:    conn,
		timeout: timeout,
	}, nil
}

func (m *GRPCMLModel) Predict(ctx context.Context, features FeatureVector) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	_ = ctx

	score := features.Velocity*0.4 + features.Amount*0.001*0.3 + features.GraphRisk*0.3
	if score > 1.0 {
		score = 1.0
	}
	if score < 0.0 {
		score = 0.0
	}

	return score, nil
}

func (m *GRPCMLModel) Close() error {
	return m.conn.Close()
}

type DefaultAIClient struct {
	model MLModel
}

func NewDefaultAIClient(model MLModel) *DefaultAIClient {
	return &DefaultAIClient{model: model}
}

func (c *DefaultAIClient) Predict(ctx context.Context, features FeatureVector) (float64, error) {
	return c.model.Predict(ctx, features)
}
