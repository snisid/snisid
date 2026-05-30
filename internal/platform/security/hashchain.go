package security

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateHashChain computes SHA-256(previousHash + payload)
func GenerateHashChain(previousHash string, payload string) string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s|%s", previousHash, payload)))
	return hex.EncodeToString(hasher.Sum(nil))
}

// VerifyHashChain validates if the current hash mathematically matches previousHash + payload
func VerifyHashChain(currentHash, previousHash, payload string) bool {
	expected := GenerateHashChain(previousHash, payload)
	return currentHash == expected
}
