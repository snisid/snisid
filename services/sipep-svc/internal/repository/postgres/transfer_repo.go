package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/sipep-svc/internal/domain"
)

type TransferRepo struct {
	pool *pgxpool.Pool
}

func NewTransferRepo(pool *pgxpool.Pool) *TransferRepo {
	return &TransferRepo{pool: pool}
}

func (r *TransferRepo) Create(ctx context.Context, transfer *domain.Transfer) error {
	query := `
		INSERT INTO sipep_transfers 
			(transfer_id, inmate_id, from_facility, to_facility, transfer_date,
			 transfer_reason, authorized_by, transport_unit, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query,
		transfer.TransferID, transfer.InmateID, transfer.FromFacility,
		transfer.ToFacility, transfer.TransferDate, transfer.TransferReason,
		transfer.AuthorizedBy, transfer.TransportUnit, transfer.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create transfer: %w", err)
	}
	return nil
}

func (r *TransferRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Transfer, error) {
	query := `
		SELECT transfer_id, inmate_id, from_facility, to_facility, transfer_date,
			   transfer_reason, authorized_by, transport_unit, created_at
		FROM sipep_transfers
		WHERE transfer_id = $1
	`
	transfer := &domain.Transfer{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&transfer.TransferID, &transfer.InmateID, &transfer.FromFacility,
		&transfer.ToFacility, &transfer.TransferDate, &transfer.TransferReason,
		&transfer.AuthorizedBy, &transfer.TransportUnit, &transfer.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find transfer: %w", err)
	}
	return transfer, nil
}

func (r *TransferRepo) FindByInmateID(ctx context.Context, inmateID uuid.UUID) ([]*domain.Transfer, error) {
	query := `
		SELECT transfer_id, inmate_id, from_facility, to_facility, transfer_date,
			   transfer_reason, authorized_by, transport_unit, created_at
		FROM sipep_transfers
		WHERE inmate_id = $1
		ORDER BY transfer_date DESC
	`
	rows, err := r.pool.Query(ctx, query, inmateID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transfers: %w", err)
	}
	defer rows.Close()

	var transfers []*domain.Transfer
	for rows.Next() {
		transfer := &domain.Transfer{}
		err := rows.Scan(
			&transfer.TransferID, &transfer.InmateID, &transfer.FromFacility,
			&transfer.ToFacility, &transfer.TransferDate, &transfer.TransferReason,
			&transfer.AuthorizedBy, &transfer.TransportUnit, &transfer.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transfer: %w", err)
		}
		transfers = append(transfers, transfer)
	}
	return transfers, nil
}
