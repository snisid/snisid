package security

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"
	"time"
)

func generateTestKey(t *testing.T) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}
	return key
}

func TestLoadRSAPrivateKey_FileNotFound(t *testing.T) {
	_, err := LoadRSAPrivateKey("/nonexistent/key.pem")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLoadRSAPublicKey_FileNotFound(t *testing.T) {
	_, err := LoadRSAPublicKey("/nonexistent/key.pub")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestGenerateRS256Token_Valid(t *testing.T) {
	key := generateTestKey(t)
	token, err := GenerateRS256Token("usr-001", []string{"admin", "officer"}, time.Hour, key)
	if err != nil {
		t.Fatalf("GenerateRS256Token failed: %v", err)
	}
	if token == "" {
		t.Error("token should not be empty")
	}
}

func TestGenerateRS256Token_Expired(t *testing.T) {
	key := generateTestKey(t)
	token, err := GenerateRS256Token("usr-002", []string{"viewer"}, -time.Hour, key)
	if err != nil {
		t.Fatalf("GenerateRS256Token failed: %v", err)
	}

	if token == "" {
		t.Error("token should not be empty (expiry is embedded in claims)")
	}
}

func TestGenerateRS256Token_DifferentKeys(t *testing.T) {
	key1 := generateTestKey(t)
	key2 := generateTestKey(t)

	token1, _ := GenerateRS256Token("usr-001", nil, time.Hour, key1)
	token2, _ := GenerateRS256Token("usr-001", nil, time.Hour, key2)

	if token1 == token2 {
		t.Error("tokens signed with different keys should be different")
	}
}

func TestGenerateRS256Token_NilKey(t *testing.T) {
	_, err := GenerateRS256Token("usr-001", nil, time.Hour, nil)
	if err == nil {
		t.Error("Expected error with nil key")
	}
}

func TestErrInvalidKey(t *testing.T) {
	if ErrInvalidKey.Error() != "invalid key format" {
		t.Errorf("ErrorMessage = %s, want 'invalid key format'", ErrInvalidKey.Error())
	}
}

func TestGenerateRS256Token_WithRoles(t *testing.T) {
	key := generateTestKey(t)
	roles := []string{"admin", "officer"}
	token, err := GenerateRS256Token("usr-003", roles, 30*time.Minute, key)
	if err != nil {
		t.Fatalf("GenerateRS256Token failed: %v", err)
	}
	if token == "" {
		t.Error("token should not be empty")
	}
}

func TestGenerateRS256Token_NoRoles(t *testing.T) {
	key := generateTestKey(t)
	token, err := GenerateRS256Token("usr-004", nil, time.Hour, key)
	if err != nil {
		t.Fatalf("GenerateRS256Token failed: %v", err)
	}
	if token == "" {
		t.Error("token should not be empty")
	}
}

// writePEMToTemp writes a PEM key to a temp file and returns the path
func writePEMToTemp(t *testing.T, content string) string {
	f, err := os.CreateTemp("", "test-key-*.pem")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}
