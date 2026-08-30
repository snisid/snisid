package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

type GangAssociationRepo struct {
	db *sqlx.DB
}

func NewGangAssociationRepo(db *sqlx.DB) *GangAssociationRepo {
	return &GangAssociationRepo{db: db}
}

func (r *GangAssociationRepo) Create(ctx context.Context, assoc *domain.GangAssociation) error {
	query := `
		INSERT INTO sivc_gang_associations (
			assoc_id, alert_id, gang_identifier, gang_territory_dept,
			gang_territory_communes, gang_snisid_id, vehicle_role,
			association_confidence, intelligence_source, source_classification,
			first_seen_date, last_confirmed_date, notes, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		assoc.AssocID, assoc.AlertID, assoc.GangIdentifier, assoc.GangTerritoryDept,
		assoc.GangTerritoryCommunes, assoc.GangSnisidID, assoc.VehicleRole,
		assoc.AssociationConfidence, assoc.IntelligenceSource, assoc.SourceClassification,
		assoc.FirstSeenDate, assoc.LastConfirmedDate, assoc.Notes, assoc.CreatedBy,
		assoc.CreatedAt, assoc.UpdatedAt,
	)
	return err
}

func (r *GangAssociationRepo) FindByAlertID(ctx context.Context, alertID uuid.UUID) ([]*domain.GangAssociation, error) {
	var assocs []*domain.GangAssociation
	query := `SELECT * FROM sivc_gang_associations WHERE alert_id = $1`
	if err := r.db.SelectContext(ctx, &assocs, query, alertID); err != nil {
		return nil, err
	}
	return assocs, nil
}

func (r *GangAssociationRepo) FindByGang(ctx context.Context, gangIdentifier string) ([]*domain.GangAssociation, error) {
	var assocs []*domain.GangAssociation
	query := `SELECT * FROM sivc_gang_associations WHERE gang_identifier = $1`
	if err := r.db.SelectContext(ctx, &assocs, query, gangIdentifier); err != nil {
		return nil, err
	}
	return assocs, nil
}
