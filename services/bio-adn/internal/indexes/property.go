package indexes

import (
	"context"

	"github.com/snisid/platform/services/bio-adn/pkg/models"
)

type PropertyIndex struct {
	db models.Database
}

func NewPropertyIndex(db models.Database) *PropertyIndex {
	return &PropertyIndex{db: db}
}

func (idx *PropertyIndex) CreateVehicle(ctx context.Context, v *models.StolenVehicle) error {
	return idx.db.CreateStolenVehicle(ctx, v)
}

func (idx *PropertyIndex) RecoverVehicle(ctx context.Context, id, location, agency string) error {
	return idx.db.UpdateVehicleStatus(ctx, id, "RECOVERED", location, agency)
}

func (idx *PropertyIndex) CreateFirearm(ctx context.Context, f *models.StolenFirearm) error {
	return idx.db.CreateStolenFirearm(ctx, f)
}

func (idx *PropertyIndex) CreateDocument(ctx context.Context, d *models.StolenDocument) error {
	return idx.db.CreateStolenDocument(ctx, d)
}

func (idx *PropertyIndex) CreateVessel(ctx context.Context, v *models.StolenVessel) error {
	return idx.db.CreateStolenVessel(ctx, v)
}

func (idx *PropertyIndex) CreateArticle(ctx context.Context, a *models.StolenArticle) error {
	return idx.db.CreateStolenArticle(ctx, a)
}

func (idx *PropertyIndex) CreateSecurity(ctx context.Context, s *models.StolenSecurity) error {
	return idx.db.CreateStolenSecurity(ctx, s)
}
