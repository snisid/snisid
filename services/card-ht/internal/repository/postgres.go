package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/card-ht/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, doc *domain.CardDocument) error
	FindByDocumentNumber(ctx context.Context, docNumber string) (*domain.CardDocument, error)
	FindByCitizenID(ctx context.Context, citizenID uuid.UUID) ([]domain.CardDocument, error)
	UpdateStatus(ctx context.Context, documentID uuid.UUID, status domain.CardStatus) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, doc *domain.CardDocument) error {
	query := `INSERT INTO card_documents (document_id, document_number, doc_type, citizen_id, status, chip_serial, mrz_line1, mrz_line2, public_key_cert_ref, issue_date, expiry_date, issuing_office, personalization_facility, photo_ref, signature_ref, sltd_reported, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`
	_, err := r.db.ExecContext(ctx, query,
		doc.DocumentID, doc.DocumentNumber, doc.DocType, doc.CitizenID, doc.Status,
		doc.ChipSerial, doc.MRZLine1, doc.MRZLine2, doc.PublicKeyCertRef,
		doc.IssueDate, doc.ExpiryDate, doc.IssuingOffice, doc.PersonalizationFacility,
		doc.PhotoRef, doc.SignatureRef, doc.SLTDReported, doc.CreatedBy, time.Now().UTC(),
	)
	return err
}

func (r *postgresRepo) FindByDocumentNumber(ctx context.Context, docNumber string) (*domain.CardDocument, error) {
	query := `SELECT document_id, document_number, doc_type, citizen_id, status, chip_serial, mrz_line1, mrz_line2, public_key_cert_ref, issue_date, expiry_date, issuing_office, personalization_facility, photo_ref, signature_ref, sltd_reported, created_by, created_at
		FROM card_documents WHERE document_number = $1`
	doc := &domain.CardDocument{}
	err := r.db.QueryRowContext(ctx, query, docNumber).Scan(
		&doc.DocumentID, &doc.DocumentNumber, &doc.DocType, &doc.CitizenID, &doc.Status,
		&doc.ChipSerial, &doc.MRZLine1, &doc.MRZLine2, &doc.PublicKeyCertRef,
		&doc.IssueDate, &doc.ExpiryDate, &doc.IssuingOffice, &doc.PersonalizationFacility,
		&doc.PhotoRef, &doc.SignatureRef, &doc.SLTDReported, &doc.CreatedBy, &doc.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document not found")
		}
		return nil, err
	}
	return doc, nil
}

func (r *postgresRepo) FindByCitizenID(ctx context.Context, citizenID uuid.UUID) ([]domain.CardDocument, error) {
	query := `SELECT document_id, document_number, doc_type, citizen_id, status, chip_serial, mrz_line1, mrz_line2, public_key_cert_ref, issue_date, expiry_date, issuing_office, personalization_facility, photo_ref, signature_ref, sltd_reported, created_by, created_at
		FROM card_documents WHERE citizen_id = $1 ORDER BY issue_date DESC`
	rows, err := r.db.QueryContext(ctx, query, citizenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []domain.CardDocument
	for rows.Next() {
		var doc domain.CardDocument
		if err := rows.Scan(&doc.DocumentID, &doc.DocumentNumber, &doc.DocType, &doc.CitizenID, &doc.Status,
			&doc.ChipSerial, &doc.MRZLine1, &doc.MRZLine2, &doc.PublicKeyCertRef,
			&doc.IssueDate, &doc.ExpiryDate, &doc.IssuingOffice, &doc.PersonalizationFacility,
			&doc.PhotoRef, &doc.SignatureRef, &doc.SLTDReported, &doc.CreatedBy, &doc.CreatedAt,
		); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, rows.Err()
}

func (r *postgresRepo) UpdateStatus(ctx context.Context, documentID uuid.UUID, status domain.CardStatus) error {
	_, err := r.db.ExecContext(ctx, `UPDATE card_documents SET status = $1 WHERE document_id = $2`, status, documentID)
	return err
}
