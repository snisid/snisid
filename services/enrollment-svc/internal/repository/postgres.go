package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/enrollment-svc/internal/domain"
)

type Repository interface {
	CreateRequest(ctx context.Context, req *domain.EnrollmentRequest) error
	FindRequestByID(ctx context.Context, requestID uuid.UUID) (*domain.EnrollmentRequest, error)
	FindPendingRequests(ctx context.Context) ([]domain.EnrollmentRequest, error)
	UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status domain.EnrollmentStatus) error
	UpdateRequest(ctx context.Context, req *domain.EnrollmentRequest) error
	CreateDocument(ctx context.Context, doc *domain.IdentityDocument) error
	FindDocumentsByRequestID(ctx context.Context, requestID uuid.UUID) ([]domain.IdentityDocument, error)
	CreateBiometricSample(ctx context.Context, sample *domain.BiometricSample) error
	FindBiometricSamplesByRequestID(ctx context.Context, requestID uuid.UUID) ([]domain.BiometricSample, error)
	CreateReview(ctx context.Context, review *domain.EnrollmentReview) error
	FindReviewByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.EnrollmentReview, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateRequest(ctx context.Context, req *domain.EnrollmentRequest) error {
	query := `INSERT INTO enrollment_requests (request_id, citizen_id, full_name, date_of_birth, nationality, email, phone, proofing_level, status, submitted_at, updated_at, assigned_officer, remarks)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.ExecContext(ctx, query,
		req.RequestID, req.CitizenID, req.FullName, req.DateOfBirth, req.Nationality,
		req.Email, req.Phone, req.ProofingLevel, req.Status, req.SubmittedAt,
		req.UpdatedAt, req.AssignedOfficer, req.Remarks,
	)
	if err != nil {
		return fmt.Errorf("insert enrollment_request: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindRequestByID(ctx context.Context, requestID uuid.UUID) (*domain.EnrollmentRequest, error) {
	query := `SELECT request_id, citizen_id, full_name, date_of_birth, nationality, email, phone, proofing_level, status, submitted_at, updated_at, assigned_officer, remarks
		FROM enrollment_requests WHERE request_id = $1`
	req := &domain.EnrollmentRequest{}
	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&req.RequestID, &req.CitizenID, &req.FullName, &req.DateOfBirth, &req.Nationality,
		&req.Email, &req.Phone, &req.ProofingLevel, &req.Status, &req.SubmittedAt,
		&req.UpdatedAt, &req.AssignedOfficer, &req.Remarks,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("enrollment request not found: %s", requestID)
		}
		return nil, fmt.Errorf("query enrollment_request: %w", err)
	}
	return req, nil
}

func (r *postgresRepo) FindPendingRequests(ctx context.Context) ([]domain.EnrollmentRequest, error) {
	query := `SELECT request_id, citizen_id, full_name, date_of_birth, nationality, email, phone, proofing_level, status, submitted_at, updated_at, assigned_officer, remarks
		FROM enrollment_requests WHERE status IN ('PENDING_DOCUMENTS','PENDING_BIOMETRICS','PENDING_REVIEW') ORDER BY submitted_at ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query pending requests: %w", err)
	}
	defer rows.Close()

	var reqs []domain.EnrollmentRequest
	for rows.Next() {
		var req domain.EnrollmentRequest
		if err := rows.Scan(
			&req.RequestID, &req.CitizenID, &req.FullName, &req.DateOfBirth, &req.Nationality,
			&req.Email, &req.Phone, &req.ProofingLevel, &req.Status, &req.SubmittedAt,
			&req.UpdatedAt, &req.AssignedOfficer, &req.Remarks,
		); err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, rows.Err()
}

func (r *postgresRepo) UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status domain.EnrollmentStatus) error {
	query := `UPDATE enrollment_requests SET status = $1, updated_at = $2 WHERE request_id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now().UTC(), requestID)
	if err != nil {
		return fmt.Errorf("update request status: %w", err)
	}
	return nil
}

func (r *postgresRepo) UpdateRequest(ctx context.Context, req *domain.EnrollmentRequest) error {
	query := `UPDATE enrollment_requests SET citizen_id=$1, full_name=$2, date_of_birth=$3, nationality=$4, email=$5, phone=$6, proofing_level=$7, status=$8, updated_at=$9, assigned_officer=$10, remarks=$11 WHERE request_id=$12`
	_, err := r.db.ExecContext(ctx, query,
		req.CitizenID, req.FullName, req.DateOfBirth, req.Nationality,
		req.Email, req.Phone, req.ProofingLevel, req.Status, time.Now().UTC(),
		req.AssignedOfficer, req.Remarks, req.RequestID,
	)
	if err != nil {
		return fmt.Errorf("update enrollment_request: %w", err)
	}
	return nil
}

func (r *postgresRepo) CreateDocument(ctx context.Context, doc *domain.IdentityDocument) error {
	query := `INSERT INTO enrollment_documents (doc_id, request_id, doc_type, doc_number, issuing_auth, issue_date, expiry_date, front_image, back_image, is_verified, verified_at, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query,
		doc.DocID, doc.RequestID, doc.DocType, doc.DocNumber, doc.IssuingAuth,
		doc.IssueDate, doc.ExpiryDate, doc.FrontImage, doc.BackImage,
		doc.IsVerified, doc.VerifiedAt, doc.UploadedAt,
	)
	if err != nil {
		return fmt.Errorf("insert enrollment_document: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindDocumentsByRequestID(ctx context.Context, requestID uuid.UUID) ([]domain.IdentityDocument, error) {
	query := `SELECT doc_id, request_id, doc_type, doc_number, issuing_auth, issue_date, expiry_date, front_image, back_image, is_verified, verified_at, uploaded_at
		FROM enrollment_documents WHERE request_id = $1 ORDER BY uploaded_at DESC`
	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, fmt.Errorf("query enrollment_documents: %w", err)
	}
	defer rows.Close()

	var docs []domain.IdentityDocument
	for rows.Next() {
		var doc domain.IdentityDocument
		if err := rows.Scan(
			&doc.DocID, &doc.RequestID, &doc.DocType, &doc.DocNumber, &doc.IssuingAuth,
			&doc.IssueDate, &doc.ExpiryDate, &doc.FrontImage, &doc.BackImage,
			&doc.IsVerified, &doc.VerifiedAt, &doc.UploadedAt,
		); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, rows.Err()
}

func (r *postgresRepo) CreateBiometricSample(ctx context.Context, sample *domain.BiometricSample) error {
	query := `INSERT INTO enrollment_biometrics (sample_id, request_id, sample_type, format, data, quality, captured_at, device_id, operator_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query,
		sample.SampleID, sample.RequestID, sample.SampleType, sample.Format,
		sample.Data, sample.Quality, sample.CapturedAt, sample.DeviceID, sample.OperatorID,
	)
	if err != nil {
		return fmt.Errorf("insert enrollment_biometric: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindBiometricSamplesByRequestID(ctx context.Context, requestID uuid.UUID) ([]domain.BiometricSample, error) {
	query := `SELECT sample_id, request_id, sample_type, format, data, quality, captured_at, device_id, operator_id
		FROM enrollment_biometrics WHERE request_id = $1 ORDER BY captured_at DESC`
	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, fmt.Errorf("query enrollment_biometrics: %w", err)
	}
	defer rows.Close()

	var samples []domain.BiometricSample
	for rows.Next() {
		var s domain.BiometricSample
		if err := rows.Scan(
			&s.SampleID, &s.RequestID, &s.SampleType, &s.Format, &s.Data,
			&s.Quality, &s.CapturedAt, &s.DeviceID, &s.OperatorID,
		); err != nil {
			return nil, err
		}
		samples = append(samples, s)
	}
	return samples, rows.Err()
}

func (r *postgresRepo) CreateReview(ctx context.Context, review *domain.EnrollmentReview) error {
	query := `INSERT INTO enrollment_reviews (review_id, request_id, officer_id, officer_name, decision, reason, verified_level, reviewed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		review.ReviewID, review.RequestID, review.OfficerID, review.OfficerName,
		review.Decision, review.Reason, review.VerifiedLevel, review.ReviewedAt,
	)
	if err != nil {
		return fmt.Errorf("insert enrollment_review: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindReviewByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.EnrollmentReview, error) {
	query := `SELECT review_id, request_id, officer_id, officer_name, decision, reason, verified_level, reviewed_at
		FROM enrollment_reviews WHERE request_id = $1`
	review := &domain.EnrollmentReview{}
	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&review.ReviewID, &review.RequestID, &review.OfficerID, &review.OfficerName,
		&review.Decision, &review.Reason, &review.VerifiedLevel, &review.ReviewedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found for request: %s", requestID)
		}
		return nil, fmt.Errorf("query enrollment_review: %w", err)
	}
	return review, nil
}
