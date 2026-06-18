package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/biar/internal/domain"
	"github.com/snisid/platform/services/biar/internal/repository"
)

type SyncService struct {
	iarmsClient *IARMSClient
	weaponRepo  repository.WeaponRepository
	syncRepo    repository.SyncRepository
}

func NewSyncService(iarmsClient *IARMSClient, weaponRepo repository.WeaponRepository, syncRepo repository.SyncRepository) *SyncService {
	return &SyncService{
		iarmsClient: iarmsClient,
		weaponRepo:  weaponRepo,
		syncRepo:    syncRepo,
	}
}

func (s *SyncService) SubmitToIARMS(ctx context.Context, weaponID uuid.UUID) (string, error) {
	w, err := s.weaponRepo.GetByID(ctx, weaponID)
	if err != nil {
		return "", fmt.Errorf("arme introuvable: %w", err)
	}

	iarmsRef, err := s.iarmsClient.SubmitIllicitWeapon(ctx, w)
	if err != nil {
		log := &domain.IARMSyncLog{
			SyncID:       uuid.New(),
			WeaponID:     &weaponID,
			Direction:    string(domain.SyncOutbound),
			SyncStatus:   string(domain.SyncFailed),
			ErrorMessage: strPtr(err.Error()),
			CreatedAt:    time.Now(),
		}
		_ = s.syncRepo.Create(ctx, log)
		return "", err
	}

	log := &domain.IARMSyncLog{
		SyncID:     uuid.New(),
		WeaponID:   &weaponID,
		Direction:  string(domain.SyncOutbound),
		IARMSRef:   &iarmsRef,
		SyncStatus: string(domain.SyncSuccess),
		SyncedAt:   timeNowPtr(),
		CreatedAt:  time.Now(),
	}
	_ = s.syncRepo.Create(ctx, log)

	return iarmsRef, nil
}

func (s *SyncService) SyncFromIARMS(ctx context.Context) (*domain.SyncResult, error) {
	result := &domain.SyncResult{StartedAt: time.Now()}

	entries, err := s.iarmsClient.FetchRecentEntries(ctx, "HTI")
	if err != nil {
		result.Error = strPtr(err.Error())
		result.CompletedAt = time.Now()
		return result, err
	}

	for _, e := range entries {
		w := convertIARMSEntry(e)
		_ = s.weaponRepo.UpsertFromIARMS(ctx, w)
		result.EntriesProcessed++

		log := &domain.IARMSyncLog{
			SyncID:     uuid.New(),
			WeaponID:   &w.WeaponID,
			Direction:  string(domain.SyncInbound),
			IARMSRef:   &e.IARMSRef,
			SyncStatus: string(domain.SyncSuccess),
			SyncedAt:   timeNowPtr(),
			CreatedAt:  time.Now(),
		}
		_ = s.syncRepo.Create(ctx, log)
	}

	result.CompletedAt = time.Now()
	return result, nil
}

func convertIARMSEntry(e *domain.IARMSEntry) *domain.IllicitWeapon {
	return &domain.IllicitWeapon{
		WeaponID:       uuid.New(),
		NationalBIARID: fmt.Sprintf("BIAR-HT-%06d", time.Now().UnixMilli()%1000000),
		WeaponType:     e.WeaponType,
		SerialNumber:   strPtr(e.SerialNumber),
		Make:           strPtr(e.Make),
		Model:          strPtr(e.Model),
		Caliber:        strPtr(e.Caliber),
		OriginCountry:  strPtr(e.OriginCountry),
		IARMSRef:       strPtr(e.IARMSRef),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func strPtr(s string) *string {
	return &s
}

func timeNowPtr() *time.Time {
	t := time.Now()
	return &t
}
