package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/snisid/counterintel-ht/internal/domain"
)

type Repository interface {
	CreateInvestigation(ctx context.Context, inv *domain.BackgroundInvestigation) error
	GetInvestigation(ctx context.Context, id uuid.UUID) (*domain.BackgroundInvestigation, error)
	GetPendingInvestigations(ctx context.Context) ([]domain.BackgroundInvestigation, error)
	UpdateInvestigation(ctx context.Context, inv *domain.BackgroundInvestigation) error
	CreateThreatAlert(ctx context.Context, alert *domain.InsiderThreatAlert) error
	GetActiveThreats(ctx context.Context) ([]domain.InsiderThreatAlert, error)
	CreateForeignContact(ctx context.Context, fc *domain.ForeignContact) error
	GetContactsBySubject(ctx context.Context, subjectID string) ([]domain.ForeignContact, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateInvestigation(ctx context.Context, inv *domain.BackgroundInvestigation) error {
	query := `INSERT INTO counterintel_background_investigations
		(id, subject_identity_ref, investigation_type, status, criminal_record_check, financial_check,
		 foreign_contacts_check, social_media_check, drug_test, psych_eval, adjudicator, adjudication_notes,
		 completed_at, clearance_level_granted, expires_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.db.ExecContext(ctx, query,
		inv.ID, inv.SubjectIdentityRef, inv.InvestigationType, inv.Status,
		inv.CriminalRecordCheck, inv.FinancialCheck, inv.ForeignContactsCheck,
		inv.SocialMediaCheck, inv.DrugTest, inv.PsychEval,
		inv.Adjudicator, inv.AdjudicationNotes, inv.CompletedAt,
		inv.ClearanceLevelGranted, inv.ExpiresAt, inv.CreatedAt, inv.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) GetInvestigation(ctx context.Context, id uuid.UUID) (*domain.BackgroundInvestigation, error) {
	query := `SELECT id, subject_identity_ref, investigation_type, status, criminal_record_check, financial_check,
		foreign_contacts_check, social_media_check, drug_test, psych_eval, adjudicator, adjudication_notes,
		completed_at, clearance_level_granted, expires_at, created_at, updated_at
		FROM counterintel_background_investigations WHERE id = $1`
	inv := &domain.BackgroundInvestigation{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inv.ID, &inv.SubjectIdentityRef, &inv.InvestigationType, &inv.Status,
		&inv.CriminalRecordCheck, &inv.FinancialCheck, &inv.ForeignContactsCheck,
		&inv.SocialMediaCheck, &inv.DrugTest, &inv.PsychEval,
		&inv.Adjudicator, &inv.AdjudicationNotes, &inv.CompletedAt,
		&inv.ClearanceLevelGranted, &inv.ExpiresAt, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("investigation not found")
		}
		return nil, err
	}
	return inv, nil
}

func (r *postgresRepo) GetPendingInvestigations(ctx context.Context) ([]domain.BackgroundInvestigation, error) {
	query := `SELECT id, subject_identity_ref, investigation_type, status, criminal_record_check, financial_check,
		foreign_contacts_check, social_media_check, drug_test, psych_eval, adjudicator, adjudication_notes,
		completed_at, clearance_level_granted, expires_at, created_at, updated_at
		FROM counterintel_background_investigations WHERE status = 'PENDING' ORDER BY created_at DESC`
	return r.scanInvestigations(ctx, query)
}

func (r *postgresRepo) UpdateInvestigation(ctx context.Context, inv *domain.BackgroundInvestigation) error {
	query := `UPDATE counterintel_background_investigations SET
		status = $2, adjudicator = $3, adjudication_notes = $4, completed_at = $5,
		clearance_level_granted = $6, expires_at = $7, updated_at = $8
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		inv.ID, inv.Status, inv.Adjudicator, inv.AdjudicationNotes,
		inv.CompletedAt, inv.ClearanceLevelGranted, inv.ExpiresAt, inv.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) CreateThreatAlert(ctx context.Context, alert *domain.InsiderThreatAlert) error {
	query := `INSERT INTO counterintel_insider_threats
		(id, subject_id, alert_type, severity, description, evidence_refs, detected_by, status, investigation_ref, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.ExecContext(ctx, query,
		alert.ID, alert.SubjectID, alert.AlertType, alert.Severity, alert.Description,
		pq.StringArray(alert.EvidenceRefs), alert.DetectedBy, alert.Status,
		alert.InvestigationRef, alert.CreatedAt,
	)
	return err
}

func (r *postgresRepo) GetActiveThreats(ctx context.Context) ([]domain.InsiderThreatAlert, error) {
	query := `SELECT id, subject_id, alert_type, severity, description, evidence_refs, detected_by, status, investigation_ref, created_at
		FROM counterintel_insider_threats WHERE status IN ('OPEN','INVESTIGATING','CONFIRMED') ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.InsiderThreatAlert
	for rows.Next() {
		var a domain.InsiderThreatAlert
		var evidenceRefs pq.StringArray
		if err := rows.Scan(
			&a.ID, &a.SubjectID, &a.AlertType, &a.Severity, &a.Description,
			&evidenceRefs, &a.DetectedBy, &a.Status, &a.InvestigationRef, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		a.EvidenceRefs = []string(evidenceRefs)
		alerts = append(alerts, a)
	}
	return alerts, rows.Err()
}

func (r *postgresRepo) CreateForeignContact(ctx context.Context, fc *domain.ForeignContact) error {
	query := `INSERT INTO counterintel_foreign_contacts
		(id, subject_id, contact_name, foreign_government, relationship_type, last_contact_at, frequency, approved_by, notes, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.ExecContext(ctx, query,
		fc.ID, fc.SubjectID, fc.ContactName, fc.ForeignGovernment, fc.RelationshipType,
		fc.LastContactAt, fc.Frequency, fc.ApprovedBy, fc.Notes, fc.CreatedAt,
	)
	return err
}

func (r *postgresRepo) GetContactsBySubject(ctx context.Context, subjectID string) ([]domain.ForeignContact, error) {
	query := `SELECT id, subject_id, contact_name, foreign_government, relationship_type, last_contact_at, frequency, approved_by, notes, created_at
		FROM counterintel_foreign_contacts WHERE subject_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []domain.ForeignContact
	for rows.Next() {
		var c domain.ForeignContact
		if err := rows.Scan(
			&c.ID, &c.SubjectID, &c.ContactName, &c.ForeignGovernment, &c.RelationshipType,
			&c.LastContactAt, &c.Frequency, &c.ApprovedBy, &c.Notes, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}

func (r *postgresRepo) scanInvestigations(ctx context.Context, query string, args ...any) ([]domain.BackgroundInvestigation, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invs []domain.BackgroundInvestigation
	for rows.Next() {
		var inv domain.BackgroundInvestigation
		if err := rows.Scan(
			&inv.ID, &inv.SubjectIdentityRef, &inv.InvestigationType, &inv.Status,
			&inv.CriminalRecordCheck, &inv.FinancialCheck, &inv.ForeignContactsCheck,
			&inv.SocialMediaCheck, &inv.DrugTest, &inv.PsychEval,
			&inv.Adjudicator, &inv.AdjudicationNotes, &inv.CompletedAt,
			&inv.ClearanceLevelGranted, &inv.ExpiresAt, &inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		invs = append(invs, inv)
	}
	return invs, rows.Err()
}
