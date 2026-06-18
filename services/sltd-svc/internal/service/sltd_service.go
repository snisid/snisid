package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sltd-svc/internal/domain"
)

type SLTDService struct {
	repo domain.Repository
	log  *zap.Logger
}

func NewSLTDService(repo domain.Repository, log *zap.Logger) *SLTDService {
	return &SLTDService{repo: repo, log: log}
}

func (s *SLTDService) CheckDocument(docNumber, issuingCountry, checkedBy, source string) (*domain.CheckResult, error) {
	doc, err := s.repo.FindByNumber(docNumber, issuingCountry)
	if err != nil {
		s.log.Info("document check - not found",
			zap.String("doc_number", docNumber),
			zap.String("issuing_country", issuingCountry))
		_ = s.repo.CreateCheckLog(&domain.SltdCheckLog{
			CheckID:        uuid.New(),
			DocumentNumber: docNumber,
			CheckedBy:      uuid.New(),
			Result:         "NOT_FOUND",
			Source:         source,
			CheckedAt:      time.Now(),
		})
		return &domain.CheckResult{
			IsBlacklisted: false,
			IsStolen:      false,
			IsLost:        false,
			Message:       "Document not found in SLTD database",
		}, nil
	}

	isStolen := doc.Status == domain.DocStatusStolen
	isLost := doc.Status == domain.DocStatusLost

	_ = s.repo.CreateCheckLog(&domain.SltdCheckLog{
		CheckID:        uuid.New(),
		DocumentNumber: docNumber,
		DocType:        doc.DocType,
		CheckedBy:      uuid.New(),
		Result:         string(doc.Status),
		Source:         source,
		SltdDocID:      &doc.DocID,
		CheckedAt:      time.Now(),
	})

	return &domain.CheckResult{
		IsBlacklisted: isStolen || isLost,
		IsStolen:      isStolen,
		IsLost:        isLost,
		Document:      doc,
		Message:       fmt.Sprintf("Document found with status: %s", doc.Status),
	}, nil
}

func (s *SLTDService) ReportLost(req *domain.ReportLostRequest) (*domain.SltdDocument, error) {
	doc := &domain.SltdDocument{
		DocID:             uuid.New(),
		NationalSltdID:    fmt.Sprintf("SLTD-%s", uuid.New().String()[:8]),
		DocType:           domain.DocType(req.DocType),
		DocumentNumber:    req.DocumentNumber,
		IssuingCountry:    req.IssuingCountry,
		HolderName:        req.HolderName,
		HolderNationality: req.HolderNationality,
		Status:            domain.DocStatusLost,
		ReportedDate:      time.Now(),
		ReportingDeptCode: req.ReportingDeptCode,
		TheftContext:      req.TheftContext,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if req.HolderSnisidID != "" {
		id, err := uuid.Parse(req.HolderSnisidID)
		if err == nil {
			doc.HolderSnisidID = &id
		}
	}
	if req.HolderDOB != "" {
		t, err := time.Parse("2006-01-02", req.HolderDOB)
		if err == nil {
			doc.HolderDOB = &t
		}
	}
	if req.IssueDate != "" {
		t, err := time.Parse("2006-01-02", req.IssueDate)
		if err == nil {
			doc.IssueDate = &t
		}
	}
	if req.ExpiryDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpiryDate)
		if err == nil {
			doc.ExpiryDate = &t
		}
	}
	if req.ReportedBy != "" {
		id, err := uuid.Parse(req.ReportedBy)
		if err == nil {
			doc.ReportedBy = id
		}
	}

	if err := s.repo.CreateDocument(doc); err != nil {
		return nil, fmt.Errorf("failed to create SLTD document: %w", err)
	}

	s.log.Info("document reported lost",
		zap.String("doc_number", req.DocumentNumber),
		zap.String("sltd_id", doc.NationalSltdID))
	return doc, nil
}

func (s *SLTDService) ReportStolen(req *domain.ReportStolenRequest) (*domain.SltdDocument, error) {
	doc := &domain.SltdDocument{
		DocID:             uuid.New(),
		NationalSltdID:    fmt.Sprintf("SLTD-%s", uuid.New().String()[:8]),
		DocType:           domain.DocType(req.DocType),
		DocumentNumber:    req.DocumentNumber,
		IssuingCountry:    req.IssuingCountry,
		HolderName:        req.HolderName,
		HolderNationality: req.HolderNationality,
		Status:            domain.DocStatusStolen,
		ReportedDate:      time.Now(),
		ReportingDeptCode: req.ReportingDeptCode,
		TheftContext:      req.TheftContext,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if req.HolderSnisidID != "" {
		id, err := uuid.Parse(req.HolderSnisidID)
		if err == nil {
			doc.HolderSnisidID = &id
		}
	}
	if req.HolderDOB != "" {
		t, err := time.Parse("2006-01-02", req.HolderDOB)
		if err == nil {
			doc.HolderDOB = &t
		}
	}
	if req.IssueDate != "" {
		t, err := time.Parse("2006-01-02", req.IssueDate)
		if err == nil {
			doc.IssueDate = &t
		}
	}
	if req.ExpiryDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpiryDate)
		if err == nil {
			doc.ExpiryDate = &t
		}
	}
	if req.ReportedBy != "" {
		id, err := uuid.Parse(req.ReportedBy)
		if err == nil {
			doc.ReportedBy = id
		}
	}

	if err := s.repo.CreateDocument(doc); err != nil {
		return nil, fmt.Errorf("failed to create SLTD document: %w", err)
	}

	s.log.Info("document reported stolen",
		zap.String("doc_number", req.DocumentNumber),
		zap.String("sltd_id", doc.NationalSltdID))
	return doc, nil
}

func (s *SLTDService) MarkFound(docID uuid.UUID, foundLocation, reportedBy string) (*domain.SltdDocument, error) {
	doc, err := s.repo.FindByID(docID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	now := time.Now()
	doc.Status = domain.DocStatusFound
	doc.FoundDate = &now
	doc.FoundLocation = foundLocation

	if reportedBy != "" {
		id, err := uuid.Parse(reportedBy)
		if err == nil {
			doc.ReportedBy = id
		}
	}

	if err := s.repo.UpdateDocument(doc); err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	s.log.Info("document marked found",
		zap.String("doc_id", docID.String()),
		zap.String("location", foundLocation))
	return doc, nil
}

func (s *SLTDService) GetStats() (*domain.SLTDStats, error) {
	return s.repo.GetStatsByType()
}
