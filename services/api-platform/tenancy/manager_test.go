package tenancy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAccess_InvalidAPIKey(t *testing.T) {
	m := TenantManager{Tenants: map[string]Tenant{}}
	ok, reason := m.ValidateAccess("invalid-key", PermVerify)
	assert.False(t, ok)
	assert.Equal(t, "INVALID_API_KEY", reason)
}

func TestValidateAccess_ValidKeyNoPermission(t *testing.T) {
	m := TenantManager{
		Tenants: map[string]Tenant{
			"key-1": {Country: "HT", Permissions: []Permission{PermScore}},
		},
	}
	ok, reason := m.ValidateAccess("key-1", PermVerify)
	assert.False(t, ok)
	assert.Equal(t, "INSUFFICIENT_PERMISSIONS", reason)
}

func TestValidateAccess_ValidKeyHasPermission(t *testing.T) {
	m := TenantManager{
		Tenants: map[string]Tenant{
			"key-1": {Country: "HT", Permissions: []Permission{PermVerify, PermScore, PermBiocheck}},
		},
	}
	ok, reason := m.ValidateAccess("key-1", PermVerify)
	assert.True(t, ok)
	assert.Empty(t, reason)
}

func TestValidateAccess_AllPermissions(t *testing.T) {
	m := TenantManager{
		Tenants: map[string]Tenant{
			"key-1": {Country: "HT", Permissions: []Permission{PermVerify, PermScore, PermBiocheck}},
		},
	}
	for _, perm := range []Permission{PermVerify, PermScore, PermBiocheck} {
		ok, reason := m.ValidateAccess("key-1", perm)
		assert.True(t, ok, "permission %s should be granted", perm)
		assert.Empty(t, reason)
	}
}

func TestValidateAccess_MultipleTenants(t *testing.T) {
	m := TenantManager{
		Tenants: map[string]Tenant{
			"key-ht": {Country: "HT", Permissions: []Permission{PermVerify}},
			"key-do": {Country: "DO", Permissions: []Permission{PermScore, PermBiocheck}},
		},
	}
	ok1, _ := m.ValidateAccess("key-ht", PermVerify)
	assert.True(t, ok1)

	ok2, _ := m.ValidateAccess("key-do", PermScore)
	assert.True(t, ok2)

	ok3, _ := m.ValidateAccess("key-ht", PermBiocheck)
	assert.False(t, ok3)
}

func TestValidateAccess_EmptyTenants(t *testing.T) {
	m := TenantManager{Tenants: map[string]Tenant{}}
	ok, reason := m.ValidateAccess("any", PermVerify)
	assert.False(t, ok)
	assert.Equal(t, "INVALID_API_KEY", reason)
}

func TestValidateAccess_PartialPermissions(t *testing.T) {
	m := TenantManager{
		Tenants: map[string]Tenant{
			"key-1": {Country: "HT", Permissions: []Permission{PermScore}},
		},
	}
	tests := []struct {
		perm Permission
		want bool
	}{
		{PermScore, true},
		{PermVerify, false},
		{PermBiocheck, false},
	}
	for _, tt := range tests {
		ok, _ := m.ValidateAccess("key-1", tt.perm)
		assert.Equal(t, tt.want, ok, "permission %s", tt.perm)
	}
}
