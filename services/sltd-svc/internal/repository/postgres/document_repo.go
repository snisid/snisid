package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/sltd-svc/internal/domain"
)

type documentRepo struct {
	pool *pgxpool.Pool
}

func NewDocumentRepo(pool *pgxpool.Pool) *documentRepo {
	return &documentRepo{pool: pool}
}

func (r *documentRepo) CreateDocument(doc *domain.SltdDocument) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO sltd_documents
		 (doc_id, national_sltd_id, doc_type, document_number, issuing_country,
		  holder_name, holder_snisid_id, holder_dob, holder_nationality,
		  issue_date, expiry_date, status, reported_date, reported_by,
		  reporting_dept_code, theft_context, found_date, found_location,
		  interpol_sltd_ref, reported_to_interpol, interpol_reported_at,
		  created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)`,
		doc.DocID, doc.NationalSltdID, doc.DocType, doc.DocumentNumber,
		doc.IssuingCountry, doc.HolderName, doc.HolderSnisidID, doc.HolderDOB,
		doc.HolderNationality, doc.IssueDate, doc.ExpiryDate, doc.Status,
		doc.ReportedDate, doc.ReportedBy, doc.ReportingDeptCode, doc.TheftContext,
		doc.FoundDate, doc.FoundLocation, doc.InterpolSltdRef, doc.ReportedToInterpol,
		doc.InterpolReportedAt, doc.CreatedAt, doc.UpdatedAt)
	return err
}

func (r *documentRepo) FindByNumber(docNumber string, issuingCountry string) (*domain.SltdDocument, error) {
	ctx := context.Background()
	doc := &domain.SltdDocument{}
	err := r.pool.QueryRow(ctx,
		`SELECT doc_id, national_sltd_id, doc_type, document_number, issuing_country,
		        holder_name, holder_snisid_id, holder_dob, holder_nationality,
		        issue_date, expiry_date, status, reported_date, reported_by,
		        reporting_dept_code, theft_context, found_date, found_location,
		        interpol_sltd_ref, reported_to_interpol, interpol_reported_at,
		        created_at, updated_at
		 FROM sltd_documents
		 WHERE document_number = $1 AND issuing_country = $2
		   AND status IN ('LOST','STOLEN','REVOKED')
		 ORDER BY created_at DESC LIMIT 1`,
		docNumber, issuingCountry).Scan(
		&doc.DocID, &doc.NationalSltdID, &doc.DocType, &doc.DocumentNumber,
		&doc.IssuingCountry, &doc.HolderName, &doc.HolderSnisidID, &doc.HolderDOB,
		&doc.HolderNationality, &doc.IssueDate, &doc.ExpiryDate, &doc.Status,
		&doc.ReportedDate, &doc.ReportedBy, &doc.ReportingDeptCode, &doc.TheftContext,
		&doc.FoundDate, &doc.FoundLocation, &doc.InterpolSltdRef, &doc.ReportedToInterpol,
		&doc.InterpolReportedAt, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (r *documentRepo) FindByID(docID uuid.UUID) (*domain.SltdDocument, error) {
	ctx := context.Background()
	doc := &domain.SltdDocument{}
	err := r.pool.QueryRow(ctx,
		`SELECT doc_id, national_sltd_id, doc_type, document_number, issuing_country,
		        holder_name, holder_snisid_id, holder_dob, holder_nationality,
		        issue_date, expiry_date, status, reported_date, reported_by,
		        reporting_dept_code, theft_context, found_date, found_location,
		        interpol_sltd_ref, reported_to_interpol, interpol_reported_at,
		        created_at, updated_at
		 FROM sltd_documents WHERE doc_id = $1`, docID).Scan(
		&doc.DocID, &doc.NationalSltdID, &doc.DocType, &doc.DocumentNumber,
		&doc.IssuingCountry, &doc.HolderName, &doc.HolderSnisidID, &doc.HolderDOB,
		&doc.HolderNationality, &doc.IssueDate, &doc.ExpiryDate, &doc.Status,
		&doc.ReportedDate, &doc.ReportedBy, &doc.ReportingDeptCode, &doc.TheftContext,
		&doc.FoundDate, &doc.FoundLocation, &doc.InterpolSltdRef, &doc.ReportedToInterpol,
		&doc.InterpolReportedAt, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (r *documentRepo) UpdateDocument(doc *domain.SltdDocument) error {
	ctx := context.Background()
	doc.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE sltd_documents SET status=$1, reported_date=$2, reported_by=$3,
		 reporting_dept_code=$4, theft_context=$5, found_date=$6, found_location=$7,
		 interpol_sltd_ref=$8, reported_to_interpol=$9, interpol_reported_at=$10, updated_at=$11
		 WHERE doc_id=$12`,
		doc.Status, doc.ReportedDate, doc.ReportedBy, doc.ReportingDeptCode,
		doc.TheftContext, doc.FoundDate, doc.FoundLocation, doc.InterpolSltdRef,
		doc.ReportedToInterpol, doc.InterpolReportedAt, doc.UpdatedAt, doc.DocID)
	return err
}

func (r *documentRepo) CreateCheckLog(log *domain.SltdCheckLog) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO sltd_check_log
		 (check_id, document_number, doc_type, checked_by, check_location,
		  post_id, result, source, sltd_doc_id, checked_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		log.CheckID, log.DocumentNumber, log.DocType, log.CheckedBy,
		log.CheckLocation, log.PostID, log.Result, log.Source,
		log.SltdDocID, log.CheckedAt)
	return err
}

func (r *documentRepo) GetStatsByType() (*domain.SLTDStats, error) {
	ctx := context.Background()
	stats := &domain.SLTDStats{
		ByStatus:  make(map[string]int),
		ByDocType: make(map[string]int),
	}

	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM sltd_documents`).Scan(&stats.TotalDocuments)
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT status::text, COUNT(*) FROM sltd_documents GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.ByStatus[status] = count
	}

	rows2, err := r.pool.Query(ctx,
		`SELECT doc_type::text, COUNT(*) FROM sltd_documents GROUP BY doc_type`)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var docType string
		var count int
		if err := rows2.Scan(&docType, &count); err != nil {
			return nil, err
		}
		stats.ByDocType[docType] = count
	}

	return stats, nil
}
