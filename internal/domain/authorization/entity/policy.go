package entity

import "time"

type Policy struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex"`
	Module    string    `json:"module"` // The raw Rego string
	Enabled   bool      `json:"enabled"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RoleGrant struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Role      string    `json:"role"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	CreatedAt time.Time `json:"createdAt"`
}

type AuthorizationRequest struct {
	Subject    SubjectData            `json:"subject"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	Attributes map[string]interface{} `json:"attributes"`
}

type SubjectData struct {
	UserID    string   `json:"userId"`
	Roles     []string `json:"roles"`
	Agency    string   `json:"agency"`
	Clearance string   `json:"clearance"`
}

type AuthorizationDecision struct {
	Allowed    bool   `json:"allowed"`
	Reason     string `json:"reason,omitempty"`
	PolicyName string `json:"policyName,omitempty"`
}
