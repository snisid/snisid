package milvus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Client struct {
	addr string
}

func NewClient(addr string) (*Client, error) {
	return &Client{addr: addr}, nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) StoreVector(ctx context.Context, templateID uuid.UUID, citizenID uuid.UUID, modality string, embedding []float32) (string, error) {
	vectorID := fmt.Sprintf("%s_%s", modality, templateID.String())
	return vectorID, nil
}

func (c *Client) Verify(ctx context.Context, modality string, sample []float32, citizenID uuid.UUID) (float64, error) {
	return 0.95, nil
}

func (c *Client) Identify(ctx context.Context, modality string, sample []float32, threshold float64) ([]CandidateResult, error) {
	return nil, nil
}

type CandidateResult struct {
	CitizenID  string
	Score      float64
	TemplateID string
}
