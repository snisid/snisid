package security

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGenerateHashChain_Consistency(t *testing.T) {
	h1 := GenerateHashChain("", "first-event")
	h2 := GenerateHashChain(h1, "second-event")
	h3 := GenerateHashChain(h2, "third-event")

	if h1 == "" {
		t.Error("hash should not be empty")
	}
	if h2 == "" {
		t.Error("hash should not be empty")
	}
	if h3 == "" {
		t.Error("hash should not be empty")
	}
	if h1 == h2 {
		t.Error("consecutive hashes should be different")
	}
}

func TestGenerateHashChain_Deterministic(t *testing.T) {
	h1 := GenerateHashChain("prev-hash", "payload-data")
	h2 := GenerateHashChain("prev-hash", "payload-data")
	if h1 != h2 {
		t.Error("hash chain should be deterministic")
	}
}

func TestGenerateHashChain_DifferentPayload(t *testing.T) {
	h1 := GenerateHashChain("prev", "payload-a")
	h2 := GenerateHashChain("prev", "payload-b")
	if h1 == h2 {
		t.Error("different payloads should produce different hashes")
	}
}

func TestVerifyHashChain_Valid(t *testing.T) {
	prevHash := "initial-hash"
	payload := "identity:created:user-001"
	currentHash := GenerateHashChain(prevHash, payload)

	if !VerifyHashChain(currentHash, prevHash, payload) {
		t.Error("VerifyHashChain should return true for valid chain")
	}
}

func TestVerifyHashChain_Invalid(t *testing.T) {
	prevHash := "initial-hash"
	payload := "identity:created:user-001"
	currentHash := GenerateHashChain(prevHash, payload)

	if VerifyHashChain(currentHash, "wrong-prev", payload) {
		t.Error("VerifyHashChain should return false with wrong previous hash")
	}
	if VerifyHashChain(currentHash, prevHash, "wrong-payload") {
		t.Error("VerifyHashChain should return false with wrong payload")
	}
}

func TestGenerateHashChain_Format(t *testing.T) {
	hash := GenerateHashChain("prev", "data")
	// SHA-256 produces 64 hex characters
	if len(hash) != 64 {
		t.Errorf("hash length = %d, want 64", len(hash))
	}
	// Verify it's valid hex
	_, err := hex.DecodeString(hash)
	if err != nil {
		t.Errorf("hash is not valid hex: %v", err)
	}
}

func TestGenerateHashChain_MatchExpected(t *testing.T) {
	prev := "previous-hash-value"
	payload := "event-data"
	// Manually compute expected hash
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s|%s", prev, payload)))
	expected := hex.EncodeToString(hasher.Sum(nil))

	got := GenerateHashChain(prev, payload)
	if got != expected {
		t.Errorf("hash = %s, want %s", got, expected)
	}
}

func TestVerifyHashChain_EmptyPrevious(t *testing.T) {
	hash := GenerateHashChain("", "first-event")
	if !VerifyHashChain(hash, "", "first-event") {
		t.Error("VerifyHashChain should work with empty previous hash")
	}
}

func TestVerifyHashChain_EmptyPayload(t *testing.T) {
	hash := GenerateHashChain("prev", "")
	if !VerifyHashChain(hash, "prev", "") {
		t.Error("VerifyHashChain should work with empty payload")
	}
}
