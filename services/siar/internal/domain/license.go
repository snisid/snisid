package domain

import (
	"time"

	"github.com/google/uuid"
)

type License struct {
	LicenseID          uuid.UUID   `json:"license_id"`
	LicenseNumber      string      `json:"license_number"`
	HolderSnisidID     uuid.UUID   `json:"holder_snisid_id"`
	HolderName         string      `json:"holder_name"`
	LicenseType        LicenseType `json:"license_type"`
	FirearmsAuthorized int         `json:"firearms_authorized"`
	IssueDate          time.Time   `json:"issue_date"`
	ExpiryDate         time.Time   `json:"expiry_date"`
	IssuingAuthority   string      `json:"issuing_authority"`
	IsActive           bool        `json:"is_active"`
	RevocationReason   string      `json:"revocation_reason,omitempty"`
	RevokedAt          *time.Time  `json:"revoked_at,omitempty"`
	CreatedBy          uuid.UUID   `json:"created_by"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

type CreateLicenseRequest struct {
	HolderSnisidID     uuid.UUID   `json:"holder_snisid_id" binding:"required"`
	HolderName         string      `json:"holder_name" binding:"required"`
	LicenseType        LicenseType `json:"license_type" binding:"required"`
	FirearmsAuthorized int         `json:"firearms_authorized"`
	IssueDate          time.Time   `json:"issue_date" binding:"required"`
	ExpiryDate         time.Time   `json:"expiry_date" binding:"required"`
	IssuingAuthority   string      `json:"issuing_authority" binding:"required"`
}
