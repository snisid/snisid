package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/rdep-svc/internal/domain"
)

type ForeignRecordRepo struct {
	pool *pgxpool.Pool
}

func NewForeignRecordRepo(pool *pgxpool.Pool) *ForeignRecordRepo {
	return &ForeignRecordRepo{pool: pool}
}

func (r *ForeignRecordRepo) Create(ctx context.Context, record *domain.ForeignRecord) error {
	query := `
		INSERT INTO rdep_foreign_records 
			(foreign_record_id, deportee_id, country, court_name, offense_description,
			 offense_date, conviction_date, sentence, prison_served, fbi_number,
			 interpol_ref, source_document, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.pool.Exec(ctx, query,
		record.ForeignRecordID, record.DeporteeID, record.Country,
		record.CourtName, record.OffenseDescription, record.OffenseDate,
		record.ConvictionDate, record.Sentence, record.PrisonServed,
		record.FBINumber, record.InterpolRef, record.SourceDocument, record.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create foreign record: %w", err)
	}
	return nil
}

func (r *ForeignRecordRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.ForeignRecord, error) {
	query := `
		SELECT foreign_record_id, deportee_id, country, court_name, offense_description,
			   offense_date, conviction_date, sentence, prison_served, fbi_number,
			   interpol_ref, source_document, created_at
		FROM rdep_foreign_records
		WHERE foreign_record_id = $1
	`
	record := &domain.ForeignRecord{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&record.ForeignRecordID, &record.DeporteeID, &record.Country,
		&record.CourtName, &record.OffenseDescription, &record.OffenseDate,
		&record.ConvictionDate, &record.Sentence, &record.PrisonServed,
		&record.FBINumber, &record.InterpolRef, &record.SourceDocument, &record.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find foreign record: %w", err)
	}
	return record, nil
}

func (r *ForeignRecordRepo) FindByDeporteeID(ctx context.Context, deporteeID uuid.UUID) ([]*domain.ForeignRecord, error) {
	query := `
		SELECT foreign_record_id, deportee_id, country, court_name, offense_description,
			   offense_date, conviction_date, sentence, prison_served, fbi_number,
			   interpol_ref, source_document, created_at
		FROM rdep_foreign_records
		WHERE deportee_id = $1
		ORDER BY offense_date DESC
	`
	rows, err := r.pool.Query(ctx, query, deporteeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query foreign records: %w", err)
	}
	defer rows.Close()

	var records []*domain.ForeignRecord
	for rows.Next() {
		record := &domain.ForeignRecord{}
		err := rows.Scan(
			&record.ForeignRecordID, &record.DeporteeID, &record.Country,
			&record.CourtName, &record.OffenseDescription, &record.OffenseDate,
			&record.ConvictionDate, &record.Sentence, &record.PrisonServed,
			&record.FBINumber, &record.InterpolRef, &record.SourceDocument, &record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan foreign record: %w", err)
		}
		records = append(records, record)
	}
	return records, nil
}
