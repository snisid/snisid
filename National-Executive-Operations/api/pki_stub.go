package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
)

// PKIProvider defines the hardware integration interface for Smartcard
type PKIProvider interface {
	VerifyPIN(signerID string, pinCode string) error
	SealDocument(content string, signerID string) (string, error)
}

// HardwarePKIStub simulates the real hardware integration
type HardwarePKIStub struct{}

func (h *HardwarePKIStub) VerifyPIN(signerID string, pinCode string) error {
	log.Printf("HardwarePKIStub: Verifying PIN for %s via Smartcard reader", signerID)
	if pinCode != "1234" {
		return errors.New("invalid smartcard PIN")
	}
	return nil
}

func (h *HardwarePKIStub) SealDocument(content string, signerID string) (string, error) {
	log.Printf("HardwarePKIStub: Sealing document for %s via HSM/Smartcard", signerID)
	// Emulate the Qualified Electronic Signature (QES) logic
	hash := sha256.New()
	hash.Write([]byte(content))
	docHash := hex.EncodeToString(hash.Sum(nil))
	return "QES_SEALED_" + signerID + "_" + docHash, nil
}
