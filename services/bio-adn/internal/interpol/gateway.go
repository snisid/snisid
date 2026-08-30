package interpol

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"
)

type I247Format string

const (
	FormatDNAProfile I247Format = "DNA_PROFILE"
	FormatHitReport  I247Format = "HIT_REPORT"
	FormatFugitive   I247Format = "FUGITIVE_NOTICE"
)

type Submission struct {
	ID        string
	SampleIDs []string
	Reason    string
	CaseRef   string
	CreatedAt time.Time
}

type Gateway struct {
	endpoint    string
	tlsCertPath string
	tlsKeyPath  string
	caCertPath  string
	client      *tls.Conn
}

type InterpolClient interface {
	Connect(ctx context.Context) error
	ToI247XML(s *Submission) ([]byte, error)
	Submit(ctx context.Context, data []byte) error
	Close() error
}

var _ InterpolClient = (*Gateway)(nil)

func NewGateway(endpoint, certPath, keyPath, caPath string) *Gateway {
	return &Gateway{
		endpoint:    endpoint,
		tlsCertPath: certPath,
		tlsKeyPath:  keyPath,
		caCertPath:  caPath,
	}
}

func (g *Gateway) Connect(ctx context.Context) error {
	if g.tlsCertPath == "" || g.tlsKeyPath == "" {
		return fmt.Errorf("INTERPOL: certificats mTLS non configurés — chemin SNISID-PKI requis")
	}

	cert, err := tls.LoadX509KeyPair(g.tlsCertPath, g.tlsKeyPath)
	if err != nil {
		return fmt.Errorf("INTERPOL: chargement certificat: %w", err)
	}

	caCert, err := os.ReadFile(g.caCertPath)
	if err != nil {
		return fmt.Errorf("INTERPOL: chargement CA: %w", err)
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		ServerName:   "gateway.interpol.int",
		MinVersion:   tls.VersionTLS13,
	}

	conn, err := tls.Dial("tcp", g.endpoint, config)
	if err != nil {
		return fmt.Errorf("INTERPOL: connexion mTLS: %w", err)
	}
	g.client = conn
	return nil
}

func (g *Gateway) ToI247XML(s *Submission) ([]byte, error) {
	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<I247Message>
  <MessageID>%s</MessageID>
  <Timestamp>%s</Timestamp>
  <Origin>HTI-BCN-DCPJ</Origin>
  <Format>%s</Format>
  <Samples>`, s.ID, s.CreatedAt.Format(time.RFC3339), FormatDNAProfile)
	for _, sid := range s.SampleIDs {
		xml += fmt.Sprintf("<SampleID>%s</SampleID>", sid)
	}
	xml += fmt.Sprintf(`</Samples>
  <Reason>%s</Reason>
  <CaseRef>%s</CaseRef>
</I247Message>`, s.Reason, s.CaseRef)
	return []byte(xml), nil
}

func (g *Gateway) Submit(ctx context.Context, data []byte) error {
	if g.client == nil {
		return fmt.Errorf("INTERPOL: non connecté — appeler Connect() d'abord")
	}
	g.client.SetDeadline(time.Now().Add(30 * time.Second))
	n, err := g.client.Write(data)
	if err != nil {
		return fmt.Errorf("INTERPOL: envoi: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("INTERPOL: envoi partiel: %d/%d bytes", n, len(data))
	}
	return nil
}

func (g *Gateway) Close() error {
	if g.client != nil {
		return g.client.Close()
	}
	return nil
}
