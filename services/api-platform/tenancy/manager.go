package apiplatform

import (
	"fmt"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

type Permission string

const (
	PermVerify   Permission = "CITIZEN_VERIFY"
	PermScore    Permission = "FRAUD_SCORE"
	PermBiocheck Permission = "BIOMETRIC_CHECK"
)

type Tenant struct {
	Country     string
	APIKey      string
	Permissions []Permission
}

type TenantManager struct {
	Tenants map[string]Tenant
}

func (m *TenantManager) ValidateAccess(apiKey string, required Permission) (bool, string) {
	tenant, ok := m.Tenants[apiKey]
	if !ok {
		return false, "INVALID_API_KEY"
	}

	for _, p := range tenant.Permissions {
		if p == required {
			logger.Info(fmt.Sprintf("API-PLATFORM: Access granted for country %s to %s", tenant.Country, required))
			return true, ""
		}
	}

	return false, "INSUFFICIENT_PERMISSIONS"
}
