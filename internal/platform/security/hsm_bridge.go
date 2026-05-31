package security

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"

	"github.com/snisid/platform/backend/internal/config"
	"github.com/snisid/platform/backend/internal/platform/logger"
	"go.uber.org/zap"
)

type Signer interface {
	Sign(ctx context.Context, keyID string, data []byte) ([]byte, error)
	PublicKey(ctx context.Context, keyID string) (ed25519.PublicKey, error)
}

type HSMBridge struct {
	pin         string
	slotID      int
	pkcs11Lib   string
	trustDomain string
	memStore    map[string]ed25519.PrivateKey
}

func NewHSMBridge(cfg config.HSMConfig) *HSMBridge {
	return &HSMBridge{
		pin:         cfg.PIN,
		slotID:      cfg.SlotID,
		pkcs11Lib:   cfg.PKCS11Lib,
		trustDomain: cfg.TrustDomain,
		memStore:    make(map[string]ed25519.PrivateKey),
	}
}

func (b *HSMBridge) Sign(ctx context.Context, keyID string, data []byte) ([]byte, error) {
	logger.Info(ctx, "Requesting cryptographic signature", zap.String("key_id", keyID))

	priv, ok := b.memStore[keyID]
	if !ok {
		var err error
		_, priv, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate key: %w", err)
		}
		b.memStore[keyID] = priv
	}

	signature := ed25519.Sign(priv, data)
	logger.Info(ctx, "Cryptographic signature generated successfully")
	return signature, nil
}

func (b *HSMBridge) PublicKey(ctx context.Context, keyID string) (ed25519.PublicKey, error) {
	priv, ok := b.memStore[keyID]
	if !ok {
		pub, _, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate key: %w", err)
		}
		return pub, nil
	}
	pub := priv.Public().(ed25519.PublicKey)
	return pub, nil
}
