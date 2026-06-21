package milvus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/snisid/idcore-svc/internal/domain"
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

func (c *Client) CheckDuplicate(ctx context.Context, sample []byte) (*domain.BiometricCheckResult, error) {
	if sample == nil {
		return &domain.BiometricCheckResult{HasMatch: false}, nil
	}

	return &domain.BiometricCheckResult{
		HasMatch:         false,
		MatchedCitizenID: uuid.Nil,
		Confidence:       0.0,
	}, nil
}

func (c *Client) StoreTemplate(ctx context.Context, citizenID uuid.UUID, sample []byte) (uuid.UUID, error) {
	if sample == nil {
		return uuid.Nil, fmt.Errorf("biometric sample is nil")
	}

	templateID := uuid.New()
	return templateID, nil
}
