package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/fir-svc/internal/domain"
)

type ConvictionRepo struct {
	pool *pgxpool.Pool
}

func NewConvictionRepo(pool *pgxpool.Pool) *ConvictionRepo {
	return &ConvictionRepo{pool: pool}
}

func (r *ConvictionRepo) Create(ctx context.Context, conviction *domain.Conviction) error {
	query := `
		INSERT INTO fir_convictions 
			(conviction_id, record_id, case_reference, court_name, court_dept,
			 offense_class, offense_description, ipc_code, verdict_date, case_status,
			 sentence_type, sentence_duration_days, fine_amount_gdes, sentence_start,
			 sentence_end, is_foreign_record, foreign_country, interpol_ccc_ref,
			 judge_name, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`
	_, err := r.pool.Exec(ctx, query,
		conviction.ConvictionID, conviction.RecordID, conviction.CaseReference,
		conviction.CourtName, conviction.CourtDept, conviction.OffenseClass,
		conviction.OffenseDescription, conviction.IPCCode, conviction.VerdictDate,
		conviction.CaseStatus, conviction.SentenceType, conviction.SentenceDurationDays,
		conviction.FineAmountGDES, conviction.SentenceStart, conviction.SentenceEnd,
		conviction.IsForeignRecord, conviction.ForeignCountry, conviction.InterpolCCCRef,
		conviction.JudgeName, conviction.Notes, conviction.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create conviction: %w", err)
	}
	return nil
}

func (r *ConvictionRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Conviction, error) {
	query := `
		SELECT conviction_id, record_id, case_reference, court_name, court_dept,
			   offense_class, offense_description, ipc_code, verdict_date, case_status,
			   sentence_type, sentence_duration_days, fine_amount_gdes, sentence_start,
			   sentence_end, is_foreign_record, foreign_country, interpol_ccc_ref,
			   judge_name, notes, created_at
		FROM fir_convictions
		WHERE conviction_id = $1
	`
	conviction := &domain.Conviction{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&conviction.ConvictionID, &conviction.RecordID, &conviction.CaseReference,
		&conviction.CourtName, &conviction.CourtDept, &conviction.OffenseClass,
		&conviction.OffenseDescription, &conviction.IPCCode, &conviction.VerdictDate,
		&conviction.CaseStatus, &conviction.SentenceType, &conviction.SentenceDurationDays,
		&conviction.FineAmountGDES, &conviction.SentenceStart, &conviction.SentenceEnd,
		&conviction.IsForeignRecord, &conviction.ForeignCountry, &conviction.InterpolCCCRef,
		&conviction.JudgeName, &conviction.Notes, &conviction.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find conviction: %w", err)
	}
	return conviction, nil
}

func (r *ConvictionRepo) FindByRecordID(ctx context.Context, recordID uuid.UUID) ([]*domain.Conviction, error) {
	query := `
		SELECT conviction_id, record_id, case_reference, court_name, court_dept,
			   offense_class, offense_description, ipc_code, verdict_date, case_status,
			   sentence_type, sentence_duration_days, fine_amount_gdes, sentence_start,
			   sentence_end, is_foreign_record, foreign_country, interpol_ccc_ref,
			   judge_name, notes, created_at
		FROM fir_convictions
		WHERE record_id = $1
		ORDER BY verdict_date DESC
	`
	rows, err := r.pool.Query(ctx, query, recordID)
	if err != nil {
		return nil, fmt.Errorf("failed to query convictions: %w", err)
	}
	defer rows.Close()

	var convictions []*domain.Conviction
	for rows.Next() {
		conviction := &domain.Conviction{}
		err := rows.Scan(
			&conviction.ConvictionID, &conviction.RecordID, &conviction.CaseReference,
			&conviction.CourtName, &conviction.CourtDept, &conviction.OffenseClass,
			&conviction.OffenseDescription, &conviction.IPCCode, &conviction.VerdictDate,
			&conviction.CaseStatus, &conviction.SentenceType, &conviction.SentenceDurationDays,
			&conviction.FineAmountGDES, &conviction.SentenceStart, &conviction.SentenceEnd,
			&conviction.IsForeignRecord, &conviction.ForeignCountry, &conviction.InterpolCCCRef,
			&conviction.JudgeName, &conviction.Notes, &conviction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conviction: %w", err)
		}
		convictions = append(convictions, conviction)
	}
	return convictions, nil
}

func (r *ConvictionRepo) Update(ctx context.Context, conviction *domain.Conviction) error {
	query := `
		UPDATE fir_convictions
		SET case_reference = $3, court_name = $4, court_dept = $5,
			offense_class = $6, offense_description = $7, ipc_code = $8,
			verdict_date = $9, case_status = $10, sentence_type = $11,
			sentence_duration_days = $12, fine_amount_gdes = $13,
			sentence_start = $14, sentence_end = $15, is_foreign_record = $16,
			foreign_country = $17, interpol_ccc_ref = $18, judge_name = $19,
			notes = $20
		WHERE conviction_id = $1 AND record_id = $2
	`
	_, err := r.pool.Exec(ctx, query,
		conviction.ConvictionID, conviction.RecordID, conviction.CaseReference,
		conviction.CourtName, conviction.CourtDept, conviction.OffenseClass,
		conviction.OffenseDescription, conviction.IPCCode, conviction.VerdictDate,
		conviction.CaseStatus, conviction.SentenceType, conviction.SentenceDurationDays,
		conviction.FineAmountGDES, conviction.SentenceStart, conviction.SentenceEnd,
		conviction.IsForeignRecord, conviction.ForeignCountry, conviction.InterpolCCCRef,
		conviction.JudgeName, conviction.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to update conviction: %w", err)
	}
	return nil
}
