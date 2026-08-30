package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/ucref-svc/internal/domain"
)

type UCREFService struct {
	repo   domain.Repository
	logger *zap.Logger
}

func NewUCREFService(repo domain.Repository, logger *zap.Logger) *UCREFService {
	return &UCREFService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UCREFService) SubmitSTR(req *domain.SubmitSTRRequest) (*domain.STRReport, error) {
	nationalID, err := s.generateNationalStrID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate national STR ID: %w", err)
	}

	reportDate := time.Now()
	if req.ReportDate != nil {
		reportDate = *req.ReportDate
	}

	report := &domain.STRReport{
		NationalStrID:        nationalID,
		ReportType:           req.ReportType,
		Status:               domain.STRStatusReceived,
		ReportingInstitution: req.ReportingInstitution,
		InstitutionType:      req.InstitutionType,
		ReportDate:           reportDate,
		TransactionDate:      req.TransactionDate,
		TransactionAmount:    req.TransactionAmount,
		TransactionCurrency:  "HTG",
		TransactionAmountUSD: req.TransactionAmountUSD,
		SubjectSnisidIDs:     req.SubjectSnisidIDs,
		SubjectNames:         req.SubjectNames,
		SubjectAccounts:      req.SubjectAccounts,
		SuspiciousActivity:   req.SuspiciousActivity,
		MLTypology:           req.MLTypology,
		PredicateCrime:       req.PredicateCrime,
		GangID:               req.GangID,
		FPRPersonIDs:         req.FPRPersonIDs,
		SancMatchIDs:         req.SancMatchIDs,
	}

	if report.SubjectSnisidIDs == nil {
		report.SubjectSnisidIDs = []string{}
	}
	if report.SubjectNames == nil {
		report.SubjectNames = []string{}
	}
	if report.SubjectAccounts == nil {
		report.SubjectAccounts = []string{}
	}
	if report.FPRPersonIDs == nil {
		report.FPRPersonIDs = []string{}
	}
	if report.SancMatchIDs == nil {
		report.SancMatchIDs = []string{}
	}
	if report.DisseminatedTo == nil {
		report.DisseminatedTo = []string{}
	}

	if err := s.repo.CreateSTR(report); err != nil {
		s.logger.Error("failed to create STR report",
			zap.String("national_str_id", nationalID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to submit STR: %w", err)
	}

	s.logger.Info("STR report submitted",
		zap.String("str_id", report.StrID.String()),
		zap.String("national_str_id", nationalID),
	)

	return report, nil
}

func (s *UCREFService) GetSTRDetail(id uuid.UUID) (*domain.STRReport, error) {
	report, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.Error("failed to get STR report",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get STR: %w", err)
	}
	return report, nil
}

func (s *UCREFService) GetFinancialProfile(personID uuid.UUID) (*domain.FinancialProfile, error) {
	profile, err := s.repo.GetFinancialProfile(personID)
	if err != nil {
		s.logger.Error("failed to get financial profile",
			zap.String("person_id", personID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	return profile, nil
}

func (s *UCREFService) RecordMonCashPattern(pattern *domain.MonCashPattern) error {
	if err := s.repo.CreateMonCashPattern(pattern); err != nil {
		s.logger.Error("failed to record MonCash pattern",
			zap.String("phone_number", pattern.PhoneNumber),
			zap.Error(err),
		)
		return fmt.Errorf("failed to record MonCash pattern: %w", err)
	}

	s.logger.Info("MonCash pattern recorded",
		zap.String("pattern_id", pattern.PatternID.String()),
		zap.String("phone_number", pattern.PhoneNumber),
	)

	return nil
}

func (s *UCREFService) GetUnanalyzedSTRs() ([]domain.STRReport, error) {
	reports, err := s.repo.GetUnanalyzedSTRs()
	if err != nil {
		s.logger.Error("failed to get unanalyzed STRs", zap.Error(err))
		return nil, fmt.Errorf("failed to get unanalyzed STRs: %w", err)
	}
	return reports, nil
}

func (s *UCREFService) DisseminateSTR(id uuid.UUID, req *domain.DisseminateRequest) error {
	report, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to get STR for dissemination: %w", err)
	}

	if report.Status != domain.STRStatusReceived && report.Status != domain.STRStatusUnderAnalysis {
		return fmt.Errorf("STR cannot be disseminated in current status: %s", report.Status)
	}

	if err := s.repo.DisseminateSTR(id, req.DisseminatedTo); err != nil {
		s.logger.Error("failed to disseminate STR",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to disseminate STR: %w", err)
	}

	s.logger.Info("STR disseminated",
		zap.String("id", id.String()),
		zap.Strings("disseminated_to", req.DisseminatedTo),
	)

	return nil
}

func (s *UCREFService) GetGangFinances(gangID uuid.UUID) ([]domain.FinancialProfile, error) {
	profiles, err := s.repo.GetGangFinances(gangID)
	if err != nil {
		s.logger.Error("failed to get gang finances",
			zap.String("gang_id", gangID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get gang finances: %w", err)
	}
	return profiles, nil
}

func (s *UCREFService) generateNationalStrID() (string, error) {
	year := time.Now().Format("2006")
	seq, err := s.repo.GetNextSequence(year)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("STR-HT-%s-%06d", year, seq), nil
}
