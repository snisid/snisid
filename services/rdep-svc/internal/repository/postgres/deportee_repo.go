package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/rdep-svc/internal/domain"
)

type DeporteeRepo struct {
	pool *pgxpool.Pool
}

func NewDeporteeRepo(pool *pgxpool.Pool) *DeporteeRepo {
	return &DeporteeRepo{pool: pool}
}

func (r *DeporteeRepo) Create(ctx context.Context, deportee *domain.Deportee) error {
	query := `
		INSERT INTO rdep_deportees 
			(deportee_id, national_rdep_id, snisid_person_id, fir_record_id, afis_subject_id,
			 deportation_country, deportation_date, arrival_port, arrival_dept_code,
			 deporting_agency, deportation_reason, flight_number, foreign_name,
			 foreign_aliases, foreign_id_number, foreign_country_id, has_foreign_record,
			 criminal_risk_level, convicted_offenses, gang_affiliated, gang_name,
			 monitoring_required, monitoring_status, monitoring_unit, monitoring_officer,
			 monitoring_end_date, current_address, current_commune, current_dept_code,
			 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
				$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)
	`
	_, err := r.pool.Exec(ctx, query,
		deportee.DeporteeID, deportee.NationalRdepID, deportee.SNISIDPersonID,
		deportee.FIRRecordID, deportee.AFISSubjectID, deportee.DeportationCountry,
		deportee.DeportationDate, deportee.ArrivalPort, deportee.ArrivalDeptCode,
		deportee.DeportingAgency, deportee.DeportationReason, deportee.FlightNumber,
		deportee.ForeignName, deportee.ForeignAliases, deportee.ForeignIDNumber,
		deportee.ForeignCountryID, deportee.HasForeignRecord, deportee.CriminalRiskLevel,
		deportee.ConvictedOffenses, deportee.GangAffiliated, deportee.GangName,
		deportee.MonitoringRequired, deportee.MonitoringStatus, deportee.MonitoringUnit,
		deportee.MonitoringOfficer, deportee.MonitoringEndDate, deportee.CurrentAddress,
		deportee.CurrentCommune, deportee.CurrentDeptCode, deportee.CreatedAt, deportee.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create deportee: %w", err)
	}
	return nil
}

func (r *DeporteeRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Deportee, error) {
	query := `
		SELECT deportee_id, national_rdep_id, snisid_person_id, fir_record_id, afis_subject_id,
			   deportation_country, deportation_date, arrival_port, arrival_dept_code,
			   deporting_agency, deportation_reason, flight_number, foreign_name,
			   foreign_aliases, foreign_id_number, foreign_country_id, has_foreign_record,
			   criminal_risk_level, convicted_offenses, gang_affiliated, gang_name,
			   monitoring_required, monitoring_status, monitoring_unit, monitoring_officer,
			   monitoring_end_date, current_address, current_commune, current_dept_code,
			   created_at, updated_at
		FROM rdep_deportees
		WHERE deportee_id = $1
	`
	deportee := &domain.Deportee{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&deportee.DeporteeID, &deportee.NationalRdepID, &deportee.SNISIDPersonID,
		&deportee.FIRRecordID, &deportee.AFISSubjectID, &deportee.DeportationCountry,
		&deportee.DeportationDate, &deportee.ArrivalPort, &deportee.ArrivalDeptCode,
		&deportee.DeportingAgency, &deportee.DeportationReason, &deportee.FlightNumber,
		&deportee.ForeignName, &deportee.ForeignAliases, &deportee.ForeignIDNumber,
		&deportee.ForeignCountryID, &deportee.HasForeignRecord, &deportee.CriminalRiskLevel,
		&deportee.ConvictedOffenses, &deportee.GangAffiliated, &deportee.GangName,
		&deportee.MonitoringRequired, &deportee.MonitoringStatus, &deportee.MonitoringUnit,
		&deportee.MonitoringOfficer, &deportee.MonitoringEndDate, &deportee.CurrentAddress,
		&deportee.CurrentCommune, &deportee.CurrentDeptCode, &deportee.CreatedAt, &deportee.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find deportee: %w", err)
	}
	return deportee, nil
}

func (r *DeporteeRepo) FindByPersonID(ctx context.Context, personID uuid.UUID) (*domain.Deportee, error) {
	query := `
		SELECT deportee_id, national_rdep_id, snisid_person_id, fir_record_id, afis_subject_id,
			   deportation_country, deportation_date, arrival_port, arrival_dept_code,
			   deporting_agency, deportation_reason, flight_number, foreign_name,
			   foreign_aliases, foreign_id_number, foreign_country_id, has_foreign_record,
			   criminal_risk_level, convicted_offenses, gang_affiliated, gang_name,
			   monitoring_required, monitoring_status, monitoring_unit, monitoring_officer,
			   monitoring_end_date, current_address, current_commune, current_dept_code,
			   created_at, updated_at
		FROM rdep_deportees
		WHERE snisid_person_id = $1
	`
	deportee := &domain.Deportee{}
	err := r.pool.QueryRow(ctx, query, personID).Scan(
		&deportee.DeporteeID, &deportee.NationalRdepID, &deportee.SNISIDPersonID,
		&deportee.FIRRecordID, &deportee.AFISSubjectID, &deportee.DeportationCountry,
		&deportee.DeportationDate, &deportee.ArrivalPort, &deportee.ArrivalDeptCode,
		&deportee.DeportingAgency, &deportee.DeportationReason, &deportee.FlightNumber,
		&deportee.ForeignName, &deportee.ForeignAliases, &deportee.ForeignIDNumber,
		&deportee.ForeignCountryID, &deportee.HasForeignRecord, &deportee.CriminalRiskLevel,
		&deportee.ConvictedOffenses, &deportee.GangAffiliated, &deportee.GangName,
		&deportee.MonitoringRequired, &deportee.MonitoringStatus, &deportee.MonitoringUnit,
		&deportee.MonitoringOfficer, &deportee.MonitoringEndDate, &deportee.CurrentAddress,
		&deportee.CurrentCommune, &deportee.CurrentDeptCode, &deportee.CreatedAt, &deportee.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find deportee by person: %w", err)
	}
	return deportee, nil
}

func (r *DeporteeRepo) Update(ctx context.Context, deportee *domain.Deportee) error {
	query := `
		UPDATE rdep_deportees
		SET deportation_country = $3, deportation_date = $4, arrival_port = $5,
			arrival_dept_code = $6, deporting_agency = $7, deportation_reason = $8,
			flight_number = $9, foreign_name = $10, foreign_aliases = $11,
			foreign_id_number = $12, foreign_country_id = $13, has_foreign_record = $14,
			criminal_risk_level = $15, convicted_offenses = $16, gang_affiliated = $17,
			gang_name = $18, monitoring_required = $19, monitoring_status = $20,
			monitoring_unit = $21, monitoring_officer = $22, monitoring_end_date = $23,
			current_address = $24, current_commune = $25, current_dept_code = $26,
			updated_at = $27
		WHERE deportee_id = $1 AND snisid_person_id = $2
	`
	_, err := r.pool.Exec(ctx, query,
		deportee.DeporteeID, deportee.SNISIDPersonID, deportee.DeportationCountry,
		deportee.DeportationDate, deportee.ArrivalPort, deportee.ArrivalDeptCode,
		deportee.DeportingAgency, deportee.DeportationReason, deportee.FlightNumber,
		deportee.ForeignName, deportee.ForeignAliases, deportee.ForeignIDNumber,
		deportee.ForeignCountryID, deportee.HasForeignRecord, deportee.CriminalRiskLevel,
		deportee.ConvictedOffenses, deportee.GangAffiliated, deportee.GangName,
		deportee.MonitoringRequired, deportee.MonitoringStatus, deportee.MonitoringUnit,
		deportee.MonitoringOfficer, deportee.MonitoringEndDate, deportee.CurrentAddress,
		deportee.CurrentCommune, deportee.CurrentDeptCode, deportee.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update deportee: %w", err)
	}
	return nil
}

func (r *DeporteeRepo) FindHighRisk(ctx context.Context) ([]*domain.Deportee, error) {
	query := `
		SELECT deportee_id, national_rdep_id, snisid_person_id, fir_record_id, afis_subject_id,
			   deportation_country, deportation_date, arrival_port, arrival_dept_code,
			   deporting_agency, deportation_reason, flight_number, foreign_name,
			   foreign_aliases, foreign_id_number, foreign_country_id, has_foreign_record,
			   criminal_risk_level, convicted_offenses, gang_affiliated, gang_name,
			   monitoring_required, monitoring_status, monitoring_unit, monitoring_officer,
			   monitoring_end_date, current_address, current_commune, current_dept_code,
			   created_at, updated_at
		FROM rdep_deportees
		WHERE criminal_risk_level IN ('HIGH', 'VERY_HIGH')
		ORDER BY created_at DESC
	`
	return r.queryDeportees(ctx, query)
}

func (r *DeporteeRepo) FindGangAffiliated(ctx context.Context) ([]*domain.Deportee, error) {
	query := `
		SELECT deportee_id, national_rdep_id, snisid_person_id, fir_record_id, afis_subject_id,
			   deportation_country, deportation_date, arrival_port, arrival_dept_code,
			   deporting_agency, deportation_reason, flight_number, foreign_name,
			   foreign_aliases, foreign_id_number, foreign_country_id, has_foreign_record,
			   criminal_risk_level, convicted_offenses, gang_affiliated, gang_name,
			   monitoring_required, monitoring_status, monitoring_unit, monitoring_officer,
			   monitoring_end_date, current_address, current_commune, current_dept_code,
			   created_at, updated_at
		FROM rdep_deportees
		WHERE gang_affiliated = TRUE
		ORDER BY created_at DESC
	`
	return r.queryDeportees(ctx, query)
}

func (r *DeporteeRepo) GetStatsByCountry(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT deportation_country, COUNT(*) as count
		FROM rdep_deportees
		GROUP BY deportation_country
		ORDER BY count DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var country string
		var count int
		if err := rows.Scan(&country, &count); err != nil {
			return nil, fmt.Errorf("failed to scan stat: %w", err)
		}
		stats[country] = count
	}
	return stats, nil
}

func (r *DeporteeRepo) queryDeportees(ctx context.Context, query string) ([]*domain.Deportee, error) {
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query deportees: %w", err)
	}
	defer rows.Close()

	var deportees []*domain.Deportee
	for rows.Next() {
		deportee := &domain.Deportee{}
		err := rows.Scan(
			&deportee.DeporteeID, &deportee.NationalRdepID, &deportee.SNISIDPersonID,
			&deportee.FIRRecordID, &deportee.AFISSubjectID, &deportee.DeportationCountry,
			&deportee.DeportationDate, &deportee.ArrivalPort, &deportee.ArrivalDeptCode,
			&deportee.DeportingAgency, &deportee.DeportationReason, &deportee.FlightNumber,
			&deportee.ForeignName, &deportee.ForeignAliases, &deportee.ForeignIDNumber,
			&deportee.ForeignCountryID, &deportee.HasForeignRecord, &deportee.CriminalRiskLevel,
			&deportee.ConvictedOffenses, &deportee.GangAffiliated, &deportee.GangName,
			&deportee.MonitoringRequired, &deportee.MonitoringStatus, &deportee.MonitoringUnit,
			&deportee.MonitoringOfficer, &deportee.MonitoringEndDate, &deportee.CurrentAddress,
			&deportee.CurrentCommune, &deportee.CurrentDeptCode, &deportee.CreatedAt, &deportee.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deportee: %w", err)
		}
		deportees = append(deportees, deportee)
	}
	return deportees, nil
}
