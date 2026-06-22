package phase13

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type Signer struct {
	key *rsa.PrivateKey
}

func NewSigner() (*Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &Signer{key: key}, nil
}

func (s *Signer) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, s.key, crypto.SHA256, hash[:])
}

func (s *Signer) Verify(data []byte, sig []byte) error {
	hash := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(&s.key.PublicKey, crypto.SHA256, hash[:], sig)
}
