package service

import (
	"context"

	"github.com/snisid/platform/services/biar/internal/domain"
	"github.com/snisid/platform/services/biar/internal/repository"
)

type StatsService struct {
	weaponRepo repository.WeaponRepository
}

func NewStatsService(weaponRepo repository.WeaponRepository) *StatsService {
	return &StatsService{weaponRepo: weaponRepo}
}

func (s *StatsService) ByGang(ctx context.Context) ([]*domain.WeaponsByGang, error) {
	return s.weaponRepo.ByGang(ctx)
}

func (s *StatsService) ByOrigin(ctx context.Context) ([]*domain.WeaponsByOrigin, error) {
	return s.weaponRepo.ByOrigin(ctx)
}

func (s *StatsService) Routes(ctx context.Context) ([]*domain.TraffickingRoute, error) {
	return s.weaponRepo.Routes(ctx)
}
