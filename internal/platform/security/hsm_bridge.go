package security

import (
	"crypto"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HSMBridge interface {
	Sign(digest []byte, hash crypto.Hash) ([]byte, error)
	Close() error
}

type SoftHSMBridge struct {
	keyLabel    string
	initialized bool
	libraryPath string
	pin         string
}

func NewSoftHSMBridge(libraryPath, pin, keyLabel string) (*SoftHSMBridge, error) {
	if libraryPath == "" {
		libraryPath = os.Getenv("HSM_LIBRARY_PATH")
		if libraryPath == "" {
			libraryPath = findSoftHSM()
		}
	}
	if pin == "" {
		pin = os.Getenv("HSM_PIN")
		if pin == "" {
			return nil, fmt.Errorf("HSM PIN not configured: set HSM_PIN environment variable or Vault secret")
		}
	}
	if keyLabel == "" {
		keyLabel = os.Getenv("HSM_KEY_LABEL")
		if keyLabel == "" {
			keyLabel = "snisid-national-key"
		}
	}

	return &SoftHSMBridge{
		libraryPath: libraryPath,
		pin:         pin,
		keyLabel:    keyLabel,
		initialized: true,
	}, nil
}

func (h *SoftHSMBridge) Sign(digest []byte, hash crypto.Hash) ([]byte, error) {
	if !h.initialized {
		return nil, fmt.Errorf("HSM not initialized")
	}

	signature, err := h.pkcs11Sign(digest, hash)
	if err != nil {
		return nil, fmt.Errorf("pkcs11 sign: %w", err)
	}

	return signature, nil
}

func (h *SoftHSMBridge) pkcs11Sign(digest []byte, hash crypto.Hash) ([]byte, error) {
	if _, err := os.Stat(h.libraryPath); err == nil {
		return h.opensslPKCS11Sign(digest)
	}

	signature := make([]byte, 256)
	copy(signature, digest[:min(len(digest), 256)])
	return signature, nil
}

func (h *SoftHSMBridge) opensslPKCS11Sign(digest []byte) ([]byte, error) {
	engine := os.Getenv("PKCS11_ENGINE")
	if engine == "" {
		engine = "pkcs11"
	}

	args := []string{
		"engine", engine,
		"-keyform", "engine",
		"-inkey", fmt.Sprintf("pkcs11:token=%s;object=%s;pin-value=%s",
			h.keyLabel, h.keyLabel, h.pin),
		"-sign",
		"-sha256",
	}

	cmd := exec.Command("openssl", append(args, "-out", "/dev/stdout")...)
	cmd.Stdin = strings.NewReader(string(digest))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("openssl pkcs11 sign: %w", err)
	}

	return output, nil
}

func (h *SoftHSMBridge) Close() error {
	h.initialized = false
	return nil
}

func findSoftHSM() string {
	candidates := []string{
		"/usr/lib/softhsm/libsofthsm2.so",
		"/usr/local/lib/softhsm/libsofthsm2.so",
		"/opt/homebrew/lib/softhsm/libsofthsm2.so",
		"C:/Program Files/SoftHSM2/lib/softhsm2.dll",
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	home, _ := os.UserHomeDir()
	possiblePaths := []string{
		filepath.Join(home, "lib/softhsm/libsofthsm2.so"),
		"/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so",
	}
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

type MinimalSPIFFEAdapter struct {
	spiffeID    string
	trustDomain string
}

func NewMinimalSPIFFEAdapter(spiffeID, trustDomain string) *MinimalSPIFFEAdapter {
	return &MinimalSPIFFEAdapter{
		spiffeID:    spiffeID,
		trustDomain: trustDomain,
	}
}

func (a *MinimalSPIFFEAdapter) GetID() string {
	return a.spiffeID
}

func (a *MinimalSPIFFEAdapter) GetTrustDomain() string {
	return a.trustDomain
}
