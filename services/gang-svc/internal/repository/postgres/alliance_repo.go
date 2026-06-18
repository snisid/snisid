package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/gang-svc/internal/domain"
)

type AllianceRepo struct {
	pool *pgxpool.Pool
}

func NewAllianceRepo(pool *pgxpool.Pool) *AllianceRepo {
	return &AllianceRepo{pool: pool}
}

func (r *AllianceRepo) Create(ctx context.Context, alliance *domain.Alliance) error {
	query := `
		INSERT INTO gang_alliances 
			(alliance_id, gang_a_id, gang_b_id, alliance_type, start_date,
			 end_date, confidence_level, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query,
		alliance.AllianceID, alliance.GangAID, alliance.GangBID,
		alliance.AllianceType, alliance.StartDate, alliance.EndDate,
		alliance.ConfidenceLevel, alliance.Notes, alliance.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create alliance: %w", err)
	}
	return nil
}

func (r *AllianceRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Alliance, error) {
	query := `
		SELECT alliance_id, gang_a_id, gang_b_id, alliance_type, start_date,
			   end_date, confidence_level, notes, created_at
		FROM gang_alliances
		WHERE alliance_id = $1
	`
	alliance := &domain.Alliance{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&alliance.AllianceID, &alliance.GangAID, &alliance.GangBID,
		&alliance.AllianceType, &alliance.StartDate, &alliance.EndDate,
		&alliance.ConfidenceLevel, &alliance.Notes, &alliance.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find alliance: %w", err)
	}
	return alliance, nil
}

func (r *AllianceRepo) FindByGangID(ctx context.Context, gangID uuid.UUID) ([]*domain.Alliance, error) {
	query := `
		SELECT alliance_id, gang_a_id, gang_b_id, alliance_type, start_date,
			   end_date, confidence_level, notes, created_at
		FROM gang_alliances
		WHERE gang_a_id = $1 OR gang_b_id = $1
		ORDER BY created_at DESC
	`
	return r.queryAlliances(ctx, query, gangID)
}

func (r *AllianceRepo) GetAllianceMap(ctx context.Context) ([]*domain.Alliance, error) {
	query := `
		SELECT alliance_id, gang_a_id, gang_b_id, alliance_type, start_date,
			   end_date, confidence_level, notes, created_at
		FROM gang_alliances
		WHERE end_date IS NULL
		ORDER BY confidence_level DESC
	`
	return r.queryAlliances(ctx, query)
}

func (r *AllianceRepo) queryAlliances(ctx context.Context, query string, args ...interface{}) ([]*domain.Alliance, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query alliances: %w", err)
	}
	defer rows.Close()

	var alliances []*domain.Alliance
	for rows.Next() {
		alliance := &domain.Alliance{}
		err := rows.Scan(
			&alliance.AllianceID, &alliance.GangAID, &alliance.GangBID,
			&alliance.AllianceType, &alliance.StartDate, &alliance.EndDate,
			&alliance.ConfidenceLevel, &alliance.Notes, &alliance.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alliance: %w", err)
		}
		alliances = append(alliances, alliance)
	}
	return alliances, nil
}
