package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidKey = errors.New("invalid key format")
)

func LoadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	privPEM, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, ErrInvalidKey
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func LoadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	pubPEM, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pubPEM)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, ErrInvalidKey
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("not RSA public key")
	}
}

func GenerateRS256Token(subject string, roles []string, ttl time.Duration, key *rsa.PrivateKey) (string, error) {
	if key == nil {
		return "", errors.New("private key is nil")
	}
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   subject,
		"roles": roles,
		"exp":   now.Add(ttl).Unix(),
		"iat":   now.Unix(),
		"nbf":   now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}
