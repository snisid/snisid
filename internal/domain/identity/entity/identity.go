package entity

import "time"

type IdentityState string

const (
	StatePending   IdentityState = "pending"
	StateActive    IdentityState = "active"
	StateSuspended IdentityState = "suspended"
	StateDeceased  IdentityState = "deceased"
)

type Identity struct {
	ID        string        `json:"id" gorm:"primaryKey"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	DOB       string        `json:"dob"`
	Gender    string        `json:"gender"`
	Agency    string        `json:"agency"`
	Status    IdentityState `json:"status"`
	Version   int           `json:"version"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`

	Biometrics []BiometricReference  `json:"biometrics,omitempty" gorm:"foreignKey:IdentityID"`
	Documents  []DocumentAssociation `json:"documents,omitempty" gorm:"foreignKey:IdentityID"`
}

type IdentityHistory struct {
	HistoryID  string        `json:"historyId" gorm:"primaryKey"`
	IdentityID string        `json:"identityId"`
	FirstName  string        `json:"firstName"`
	LastName   string        `json:"lastName"`
	DOB        string        `json:"dob"`
	Gender     string        `json:"gender"`
	Agency     string        `json:"agency"`
	Status     IdentityState `json:"status"`
	Version    int           `json:"version"`
	ChangedAt  time.Time     `json:"changedAt"`
	ChangedBy  string        `json:"changedBy"`
	Reason     string        `json:"reason"`
}

type BiometricReference struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	IdentityID   string    `json:"identityId"`
	Type         string    `json:"type"` // e.g., "face", "fingerprint"
	ReferenceURI string    `json:"referenceUri"`
	QualityScore float64   `json:"qualityScore"`
	CreatedAt    time.Time `json:"createdAt"`
}

type DocumentAssociation struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	IdentityID   string    `json:"identityId"`
	DocumentType string    `json:"documentType"` // e.g., "passport", "birth_certificate"
	DocumentURI  string    `json:"documentUri"`
	Verified     bool      `json:"verified"`
	CreatedAt    time.Time `json:"createdAt"`
}
