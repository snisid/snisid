package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

type IntelReportRepo struct {
	db *sqlx.DB
}

func NewIntelReportRepo(db *sqlx.DB) *IntelReportRepo {
	return &IntelReportRepo{db: db}
}

func (r *IntelReportRepo) Create(ctx context.Context, report *domain.IntelligenceReport) error {
	query := `
		INSERT INTO sivc_intelligence_reports (
			report_id, report_number, title, report_type, classification,
			summary, full_report, alert_ids, plate_ids, person_ids,
			originating_unit, author_id, recipient_units, published_at,
			expiry_date, attachments, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		report.ReportID, report.ReportNumber, report.Title, report.ReportType,
		report.Classification, report.Summary, report.FullReport, report.AlertIDs,
		report.PlateIDs, report.PersonIDs, report.OriginatingUnit, report.AuthorID,
		report.RecipientUnits, report.PublishedAt, report.ExpiryDate, report.Attachments,
		report.CreatedAt, report.UpdatedAt,
	)
	return err
}

func (r *IntelReportRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.IntelligenceReport, error) {
	var report domain.IntelligenceReport
	query := `SELECT * FROM sivc_intelligence_reports WHERE report_id = $1`
	if err := r.db.GetContext(ctx, &report, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &report, nil
}

func (r *IntelReportRepo) FindByUnit(ctx context.Context, unit string) ([]*domain.IntelligenceReport, error) {
	var reports []*domain.IntelligenceReport
	query := `SELECT * FROM sivc_intelligence_reports WHERE originating_unit = $1 ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &reports, query, unit); err != nil {
		return nil, err
	}
	return reports, nil
}
