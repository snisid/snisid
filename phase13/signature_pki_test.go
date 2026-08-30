package phase13

import (
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	signer, err := NewSigner()
	if err != nil {
		t.Fatalf("NewSigner() failed: %v", err)
	}

	data := []byte("test data for SNISID Phase 13 PKI")
	sig, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	if err := signer.Verify(data, sig); err != nil {
		t.Fatalf("Verify() failed: %v", err)
	}
}

func TestVerifyBadSignature(t *testing.T) {
	signer, err := NewSigner()
	if err != nil {
		t.Fatalf("NewSigner() failed: %v", err)
	}

	err = signer.Verify([]byte("data"), []byte("bad"))
	if err == nil {
		t.Fatal("expected error for bad signature")
	}
}
