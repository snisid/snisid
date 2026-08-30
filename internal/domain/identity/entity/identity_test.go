package entity

import (
	"testing"
	"time"
)

func TestIdentityState_Values(t *testing.T) {
	tests := []struct {
		state IdentityState
		want  string
	}{
		{StatePending, "pending"},
		{StateActive, "active"},
		{StateSuspended, "suspended"},
		{StateDeceased, "deceased"},
	}

	for _, tt := range tests {
		if string(tt.state) != tt.want {
			t.Errorf("IdentityState(%s) = %s, want %s", tt.want, string(tt.state), tt.want)
		}
	}
}

func TestIdentity_DefaultValues(t *testing.T) {
	now := time.Now()
	id := Identity{
		ID:        "ID-123",
		FirstName: "Jean",
		LastName:  "Dupont",
		DOB:       "1990-01-01",
		Gender:    "M",
		Agency:    "ONI",
		Status:    StateActive,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if id.ID != "ID-123" {
		t.Errorf("ID = %s, want ID-123", id.ID)
	}
	if id.Status != StateActive {
		t.Errorf("Status = %s, want active", id.Status)
	}
	if id.Version != 1 {
		t.Errorf("Version = %d, want 1", id.Version)
	}
}

func TestIdentity_Relations(t *testing.T) {
	id := Identity{
		ID: "ID-456",
		Biometrics: []BiometricReference{
			{ID: "BIO-1", IdentityID: "ID-456", Type: "face", QualityScore: 0.95},
		},
		Documents: []DocumentAssociation{
			{ID: "DOC-1", IdentityID: "ID-456", DocumentType: "passport", Verified: true},
		},
	}

	if len(id.Biometrics) != 1 {
		t.Errorf("len(Biometrics) = %d, want 1", len(id.Biometrics))
	}
	if len(id.Documents) != 1 {
		t.Errorf("len(Documents) = %d, want 1", len(id.Documents))
	}
	if id.Biometrics[0].Type != "face" {
		t.Errorf("Biometric type = %s, want face", id.Biometrics[0].Type)
	}
	if !id.Documents[0].Verified {
		t.Error("Document should be verified")
	}
}

func TestIdentityHistory_Fields(t *testing.T) {
	now := time.Now()
	h := IdentityHistory{
		HistoryID:  "H-1",
		IdentityID: "ID-123",
		FirstName:  "Jean",
		Version:    2,
		ChangedAt:  now,
		ChangedBy:  "admin",
		Reason:     "Correction du prénom",
	}

	if h.ChangedBy != "admin" {
		t.Errorf("ChangedBy = %s, want admin", h.ChangedBy)
	}
	if h.Version != 2 {
		t.Errorf("Version = %d, want 2", h.Version)
	}
}

func TestBiometricReference_Quality(t *testing.T) {
	ref := BiometricReference{
		ID:           "BIO-1",
		Type:         "fingerprint",
		QualityScore: 0.88,
	}

	if ref.QualityScore < 0.8 {
		t.Errorf("QualityScore = %f, want >= 0.8", ref.QualityScore)
	}
	if ref.Type != "fingerprint" {
		t.Errorf("Type = %s, want fingerprint", ref.Type)
	}
}

func TestDocumentAssociation_Verification(t *testing.T) {
	doc := DocumentAssociation{
		ID:           "DOC-1",
		DocumentType: "birth_certificate",
		Verified:     true,
	}

	if !doc.Verified {
		t.Error("Birth certificate should be verified")
	}
	if doc.DocumentType != "birth_certificate" {
		t.Errorf("DocumentType = %s, want birth_certificate", doc.DocumentType)
	}
}
