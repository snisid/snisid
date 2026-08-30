package oidc

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	once       sync.Once
	privateKey *rsa.PrivateKey
	keyID      string
)

func ensureKey() error {
	var err error
	once.Do(func() {
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return
		}
		kid := make([]byte, 16)
		rand.Read(kid)
		keyID = base64.RawURLEncoding.EncodeToString(kid)
	})
	return err
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	if err := ensureKey(); err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}
	return privateKey, nil
}

// JWKS serves the JSON Web Key Set for token signature verification.
// GET /.well-known/jwks
func (h *Handler) JWKS(c *gin.Context) {
	if err := ensureKey(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "key generation failed"})
		return
	}
	n := base64.RawURLEncoding.EncodeToString(privateKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(privateKey.E)).Bytes())
	c.JSON(http.StatusOK, gin.H{
		"keys": []gin.H{
			{
				"kty": "RSA",
				"use": "sig",
				"alg": "RS256",
				"kid": keyID,
				"n":   n,
				"e":   e,
			},
		},
	})
}
