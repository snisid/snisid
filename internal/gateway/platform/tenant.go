package platform

import (
	"context"
	"errors"
)

type Tenant struct {
	ID          string
	Country     string
	APIKey      string
	Permissions []string
	Namespace   string // Kubernetes namespace for isolation
}

type PlatformGateway struct {
	Tenants map[string]*Tenant
}

func (g *PlatformGateway) Authenticate(ctx context.Context, apiKey string) (*Tenant, error) {
	for _, tenant := range g.Tenants {
		if tenant.APIKey == apiKey {
			return tenant, nil
		}
	}
	return nil, errors.New("invalid_api_key")
}

func (g *PlatformGateway) Authorize(tenant *Tenant, action string) bool {
	for _, p := range tenant.Permissions {
		if p == action {
			return true
		}
	}
	return false
}
