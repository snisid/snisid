package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/civil-ht/internal/domain"
)

type Repository interface {
	CreateBirth(ctx context.Context, act *domain.CivilAct, birth *domain.BirthDeclaration) error
	CreateDeath(ctx context.Context, act *domain.CivilAct, death *domain.DeathDeclaration) error
	CreateMarriage(ctx context.Context, act *domain.CivilAct, marriage *domain.MarriageDeclaration) error
	FindByActNumber(ctx context.Context, actNumber string) (*domain.CivilAct, error)
	FindByCitizenID(ctx context.Context, citizenID uuid.UUID) ([]domain.CivilAct, error)
	FindBirthDetails(ctx context.Context, actID uuid.UUID) (*domain.BirthDeclaration, error)
	FindDeathDetails(ctx context.Context, actID uuid.UUID) (*domain.DeathDeclaration, error)
	FindMarriageDetails(ctx context.Context, actID uuid.UUID) (*domain.MarriageDeclaration, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateBirth(ctx context.Context, act *domain.CivilAct, birth *domain.BirthDeclaration) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	actQuery := `INSERT INTO civil_acts (act_id, act_number, act_type, citizen_id, registering_office, dept_code, commune, event_date, declared_date, officer_name, officer_id, is_late_declaration, is_reconstructed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	if _, err := tx.ExecContext(ctx, actQuery,
		act.ActID, act.ActNumber, act.ActType, act.CitizenID, act.RegisteringOffice,
		act.DeptCode, act.Commune, act.EventDate, act.DeclaredDate, act.OfficerName,
		act.OfficerID, act.IsLateDeclaration, act.IsReconstructed, time.Now().UTC(),
	); err != nil {
		return fmt.Errorf("insert civil_act: %w", err)
	}

	birthQuery := `INSERT INTO civil_birth_details (act_id, child_full_name, child_gender, mother_citizen_id, father_citizen_id, birth_weight_g, birth_facility, attending_professional)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	if _, err := tx.ExecContext(ctx, birthQuery,
		birth.ActID, birth.ChildFullName, birth.ChildGender, birth.MotherCitizenID,
		birth.FatherCitizenID, birth.BirthWeightG, birth.BirthFacility, birth.AttendingProfessional,
	); err != nil {
		return fmt.Errorf("insert birth_details: %w", err)
	}

	return tx.Commit()
}

func (r *postgresRepo) CreateDeath(ctx context.Context, act *domain.CivilAct, death *domain.DeathDeclaration) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	actQuery := `INSERT INTO civil_acts (act_id, act_number, act_type, citizen_id, registering_office, dept_code, commune, event_date, declared_date, officer_name, officer_id, is_late_declaration, is_reconstructed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	if _, err := tx.ExecContext(ctx, actQuery,
		act.ActID, act.ActNumber, act.ActType, act.CitizenID, act.RegisteringOffice,
		act.DeptCode, act.Commune, act.EventDate, act.DeclaredDate, act.OfficerName,
		act.OfficerID, act.IsLateDeclaration, act.IsReconstructed, time.Now().UTC(),
	); err != nil {
		return fmt.Errorf("insert civil_act: %w", err)
	}

	deathQuery := `INSERT INTO civil_death_details (act_id, deceased_citizen_id, cause_of_death, death_location, medical_certifier, is_violent_death, fir_case_reference)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := tx.ExecContext(ctx, deathQuery,
		death.ActID, death.DeceasedCitizenID, death.CauseOfDeath, death.DeathLocation,
		death.MedicalCertifier, death.IsViolentDeath, death.FIRCaseReference,
	); err != nil {
		return fmt.Errorf("insert death_details: %w", err)
	}

	return tx.Commit()
}

func (r *postgresRepo) CreateMarriage(ctx context.Context, act *domain.CivilAct, marriage *domain.MarriageDeclaration) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	actQuery := `INSERT INTO civil_acts (act_id, act_number, act_type, registering_office, dept_code, commune, event_date, declared_date, officer_name, officer_id, is_late_declaration, is_reconstructed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	now := time.Now().UTC()
	if _, err := tx.ExecContext(ctx, actQuery,
		act.ActID, act.ActNumber, act.ActType, act.RegisteringOffice,
		act.DeptCode, act.Commune, act.EventDate, act.DeclaredDate, act.OfficerName,
		act.OfficerID, act.IsLateDeclaration, act.IsReconstructed, now,
	); err != nil {
		return fmt.Errorf("insert civil_act: %w", err)
	}

	marriageQuery := `INSERT INTO civil_marriage_details (act_id, spouse_a_citizen_id, spouse_b_citizen_id, marriage_regime, prenuptial_agreement)
		VALUES ($1, $2, $3, $4, $5)`
	if _, err := tx.ExecContext(ctx, marriageQuery,
		marriage.ActID, marriage.SpouseACitizenID, marriage.SpouseBCitizenID,
		marriage.MarriageRegime, marriage.PrenuptialAgreement,
	); err != nil {
		return fmt.Errorf("insert marriage_details: %w", err)
	}

	return tx.Commit()
}

func (r *postgresRepo) FindByActNumber(ctx context.Context, actNumber string) (*domain.CivilAct, error) {
	query := `SELECT act_id, act_number, act_type, citizen_id, registering_office, dept_code, commune, event_date, declared_date, officer_name, officer_id, is_late_declaration, is_reconstructed, created_at
		FROM civil_acts WHERE act_number = $1`

	act := &domain.CivilAct{}
	err := r.db.QueryRowContext(ctx, query, actNumber).Scan(
		&act.ActID, &act.ActNumber, &act.ActType, &act.CitizenID, &act.RegisteringOffice,
		&act.DeptCode, &act.Commune, &act.EventDate, &act.DeclaredDate, &act.OfficerName,
		&act.OfficerID, &act.IsLateDeclaration, &act.IsReconstructed, &act.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("act not found: %s", actNumber)
		}
		return nil, fmt.Errorf("query act: %w", err)
	}
	return act, nil
}

func (r *postgresRepo) FindByCitizenID(ctx context.Context, citizenID uuid.UUID) ([]domain.CivilAct, error) {
	query := `SELECT act_id, act_number, act_type, citizen_id, registering_office, dept_code, commune, event_date, declared_date, officer_name, officer_id, is_late_declaration, is_reconstructed, created_at
		FROM civil_acts WHERE citizen_id = $1 ORDER BY event_date DESC`

	rows, err := r.db.QueryContext(ctx, query, citizenID)
	if err != nil {
		return nil, fmt.Errorf("query acts by citizen: %w", err)
	}
	defer rows.Close()

	var acts []domain.CivilAct
	for rows.Next() {
		var act domain.CivilAct
		if err := rows.Scan(&act.ActID, &act.ActNumber, &act.ActType, &act.CitizenID, &act.RegisteringOffice,
			&act.DeptCode, &act.Commune, &act.EventDate, &act.DeclaredDate, &act.OfficerName,
			&act.OfficerID, &act.IsLateDeclaration, &act.IsReconstructed, &act.CreatedAt,
		); err != nil {
			return nil, err
		}
		acts = append(acts, act)
	}
	return acts, rows.Err()
}

func (r *postgresRepo) FindBirthDetails(ctx context.Context, actID uuid.UUID) (*domain.BirthDeclaration, error) {
	b := &domain.BirthDeclaration{}
	err := r.db.QueryRowContext(ctx, `SELECT act_id, child_full_name, child_gender, mother_citizen_id, father_citizen_id, birth_weight_g, birth_facility, attending_professional FROM civil_birth_details WHERE act_id = $1`, actID).Scan(
		&b.ActID, &b.ChildFullName, &b.ChildGender, &b.MotherCitizenID, &b.FatherCitizenID,
		&b.BirthWeightG, &b.BirthFacility, &b.AttendingProfessional,
	)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *postgresRepo) FindDeathDetails(ctx context.Context, actID uuid.UUID) (*domain.DeathDeclaration, error) {
	d := &domain.DeathDeclaration{}
	err := r.db.QueryRowContext(ctx, `SELECT act_id, deceased_citizen_id, cause_of_death, death_location, medical_certifier, is_violent_death, fir_case_reference FROM civil_death_details WHERE act_id = $1`, actID).Scan(
		&d.ActID, &d.DeceasedCitizenID, &d.CauseOfDeath, &d.DeathLocation,
		&d.MedicalCertifier, &d.IsViolentDeath, &d.FIRCaseReference,
	)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *postgresRepo) FindMarriageDetails(ctx context.Context, actID uuid.UUID) (*domain.MarriageDeclaration, error) {
	m := &domain.MarriageDeclaration{}
	err := r.db.QueryRowContext(ctx, `SELECT act_id, spouse_a_citizen_id, spouse_b_citizen_id, marriage_regime, prenuptial_agreement FROM civil_marriage_details WHERE act_id = $1`, actID).Scan(
		&m.ActID, &m.SpouseACitizenID, &m.SpouseBCitizenID,
		&m.MarriageRegime, &m.PrenuptialAgreement,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}
