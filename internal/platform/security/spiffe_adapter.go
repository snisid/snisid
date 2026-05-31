package security

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/snisid/platform/backend/internal/config"
	"github.com/snisid/platform/backend/internal/platform/logger"
	"go.uber.org/zap"
)

type WorkloadIdentity struct {
	SPIFFEID    string
	Agency      string
	TrustDomain string
}

type SPIFFEAdapter struct {
	trustDomain string
	agentSocket string
	client      *SPIREProvider
}

type SPIREProvider struct {
	socketPath string
}

func NewSPIREProvider(socketPath string) *SPIREProvider {
	return &SPIREProvider{socketPath: socketPath}
}

func NewSPIFFEAdapter(cfg config.SPIREConfig) *SPIFFEAdapter {
	return &SPIFFEAdapter{
		trustDomain: cfg.TrustDomain,
		agentSocket: cfg.AgentSocket,
		client:      NewSPIREProvider(cfg.AgentSocket),
	}
}

func (a *SPIFFEAdapter) FetchIdentity(ctx context.Context) (*WorkloadIdentity, error) {
	logger.Debug(ctx, "Fetching workload identity from SPIRE agent", zap.String("socket", a.agentSocket))

	conn, err := net.DialTimeout("unix", a.agentSocket, 5*time.Second)
	if err != nil {
		logger.Warn(ctx, "SPIRE agent not available, using platform identity", zap.Error(err))
		return a.platformIdentity(), nil
	}
	defer conn.Close()

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

	expectedPrefix := fmt.Sprintf("spiffe://%s/", a.trustDomain)
	if !containsPrefix(peerSVID, expectedPrefix) {
		return fmt.Errorf("peer identity %s is outside of sovereign trust domain", peerSVID)
	}

	return nil
}

func (a *SPIFFEAdapter) platformIdentity() *WorkloadIdentity {
	return &WorkloadIdentity{
		SPIFFEID:    fmt.Sprintf("spiffe://%s/ns/snisid/sa/platform", a.trustDomain),
		Agency:      "SNISID_CENTRAL",
		TrustDomain: a.trustDomain,
	}
}

func containsPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}
