package security

import (
	"context"
	"fmt"

	"github.com/snisid/platform/backend/internal/platform/logger"
	"go.uber.org/zap"
)

type WorkloadIdentity struct {
	SPIFFEID string
	Agency   string
	TrustDomain string
}

type SPIFFEAdapter struct {
	trustDomain string
}

func NewSPIFFEAdapter(trustDomain string) *SPIFFEAdapter {
	return &SPIFFEAdapter{trustDomain: trustDomain}
}

func (a *SPIFFEAdapter) FetchIdentity(ctx context.Context) (*WorkloadIdentity, error) {
	logger.Debug(ctx, "Fetching workload identity from SPIRE agent")

	// Mock: Simulating SPIRE SVID fetch
	id := &WorkloadIdentity{
		SPIFFEID:    fmt.Sprintf("spiffe://%s/ns/snisid/sa/orchestrator", a.trustDomain),
		Agency:      "SNISID_CENTRAL",
		TrustDomain: a.trustDomain,
	}

	logger.Info(ctx, "Workload identity verified", zap.String("spiffe_id", id.SPIFFEID))
	return id, nil
}

func (a *SPIFFEAdapter) ValidatePeer(ctx context.Context, peerSVID string) error {
	logger.Debug(ctx, "Validating peer workload identity", zap.String("peer_svid", peerSVID))
	
	// Mock: Verify the peer belongs to the same trust domain
	expectedPrefix := fmt.Sprintf("spiffe://%s/", a.trustDomain)
	if !containsPrefix(peerSVID, expectedPrefix) {
		return fmt.Errorf("peer identity %s is outside of sovereign trust domain", peerSVID)
	}

	return nil
}

func containsPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}
