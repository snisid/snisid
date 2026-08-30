package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/pki-ht/internal/domain"
	"github.com/snisid/pki-ht/internal/kafka"
	"github.com/snisid/pki-ht/internal/repository"
)

type PKIService struct {
	repo     repository.Repository
	producer *kafka.Producer
	ca       *domain.CertificateAuthority
}

func NewPKIService(repo repository.Repository, producer *kafka.Producer) *PKIService {
	return &PKIService{repo: repo, producer: producer}
}

func (s *PKIService) InitCA(ctx context.Context) error {
	ca, err := s.generateRootCA()
	if err != nil {
		return fmt.Errorf("generate root CA: %w", err)
	}
	s.ca = ca
	return nil
}

func (s *PKIService) Issue(ctx context.Context, req domain.IssueRequest) (*domain.IssuedCertificate, error) {
	if s.ca == nil {
		return nil, fmt.Errorf("CA not initialized")
	}
	serial, err := generateSerial()
	if err != nil {
		return nil, err
	}

	subjectType := domain.SubjectType(req.SubjectType)
	var subjectRef *uuid.UUID
	if req.SubjectRef != "" {
		if id, err := uuid.Parse(req.SubjectRef); err == nil {
			subjectRef = &id
		}
	}

	cert := &domain.IssuedCertificate{
		CertID:       uuid.New(),
		SerialNumber: serial,
		IssuingCAID:  s.ca.CAID,
		SubjectType:  subjectType,
		SubjectRef:   subjectRef,
		CommonName:   &req.CommonName,
		Status:       domain.CertValid,
		ValidFrom:    time.Now().UTC(),
		ValidUntil:   time.Now().UTC().AddDate(1, 0, 0),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.repo.CreateCertificate(ctx, cert); err != nil {
		return nil, fmt.Errorf("save certificate: %w", err)
	}

	s.publishEvent(ctx, "pki.certificate.issued", cert)
	return cert, nil
}

func (s *PKIService) Revoke(ctx context.Context, serial string, reason string) error {
	if err := s.repo.RevokeCertificate(ctx, serial, reason); err != nil {
		return fmt.Errorf("revoke: %w", err)
	}

	cert, err := s.repo.FindBySerial(ctx, serial)
	if err != nil {
		return err
	}

	s.publishEvent(ctx, "pki.certificate.revoked", cert)
	return nil
}

func (s *PKIService) CheckOCSP(ctx context.Context, serial string) (*domain.IssuedCertificate, error) {
	return s.repo.FindBySerial(ctx, serial)
}

func (s *PKIService) GetCRL(ctx context.Context, caID string) (*domain.CRL, error) {
	cid, err := uuid.Parse(caID)
	if err != nil {
		return nil, fmt.Errorf("invalid ca_id: %w", err)
	}
	return s.repo.GetActiveCRL(ctx, cid)
}

func (s *PKIService) generateRootCA() (*domain.CertificateAuthority, error) {
	serial, _ := generateSerial()
	ca := &domain.CertificateAuthority{
		CAID:         uuid.New(),
		CAType:       domain.CARoot,
		CommonName:   "SNISID Root CA",
		SerialNumber: serial,
		PublicKeyPEM: "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\n-----END PUBLIC KEY-----",
		ValidFrom:    time.Now().UTC(),
		ValidUntil:   time.Now().UTC().AddDate(20, 0, 0),
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
	}
	return ca, nil
}

func (s *PKIService) publishEvent(ctx context.Context, eventType string, cert *domain.IssuedCertificate) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType:    eventType,
		SerialNumber: cert.SerialNumber,
		Timestamp:    time.Now().UTC(),
		Data:         cert,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}

func generateSerial() (string, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(serial.Bytes()), nil
}

var _ = x509.Certificate{}
var _ = rsa.PrivateKey{}
var _ = pkix.Name{}
var _ = sha256.New
