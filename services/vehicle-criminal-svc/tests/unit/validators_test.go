package unit

import (
	"testing"

	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNormalizePlate(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"PP-1234", "PP1234"},
		{"SE 00871", "SE00871"},
		{"ABC-123456", "ABC123456"},
		{"PP1234", "PP1234"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := domain.NormalizePlate(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateVIN(t *testing.T) {
	tests := []struct {
		name    string
		vin     string
		wantErr bool
	}{
		{"valid 17 char VIN", "1HGBH41JXMN109186", false},
		{"too short", "1234567890123456", true},
		{"too long", "1HGBH41JXMN1091860", true},
		{"empty", "", true},
		{"invalid chars I", "1HGBH41JXMN109I86", true},
		{"invalid chars O", "1HGBH41JXMN109O86", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidateVIN(tt.vin)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsStatePlate_EdgeCases(t *testing.T) {
	tests := []struct {
		plate string
		want  bool
	}{
		{"SE-00871", true},
		{"SE00871", true},
		{"SE 00871", true},
		{"SE-123", true},
		{"SE-123456789", false},
		{"SE-", false},
		{"SE", false},
		{"PP-00871", false},
		{"se-00871", false},
	}

	for _, tt := range tests {
		t.Run(tt.plate, func(t *testing.T) {
			got := domain.IsStatePlate(tt.plate)
			assert.Equal(t, tt.want, got, "plate: %s", tt.plate)
		})
	}
}
