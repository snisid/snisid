package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/snisid/mil-c2-ht/internal/domain"
)

type MilC2Repo interface {
	CreateUnit(u domain.MilitaryUnit) error
	GetDeployedUnits() ([]domain.MilitaryUnit, error)
	CreateOperation(o domain.Operation) error
	GetActiveOperations() ([]domain.Operation, error)
	CreateTacticalReport(r domain.TacticalReport) error
	GetReportsByOperation(opID uuid.UUID) ([]domain.TacticalReport, error)
	GetAllUnits() ([]domain.MilitaryUnit, error)
	GetAllOperations() ([]domain.Operation, error)
	GetAllReports() ([]domain.TacticalReport, error)
}

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) CreateUnit(u domain.MilitaryUnit) error {
	query := `INSERT INTO milc2_units 
		(unit_id, unit_name, branch, parent_unit_id, commander_name, personnel_count, 
		 location_lat, location_lng, operational_status, equipment_summary, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.db.Exec(query,
		u.UnitID, u.UnitName, u.Branch, u.ParentUnitID, u.CommanderName, u.PersonnelCount,
		u.LocationLat, u.LocationLng, u.OperationalStatus, u.EquipmentSummary, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *PostgresRepo) GetDeployedUnits() ([]domain.MilitaryUnit, error) {
	query := `SELECT unit_id, unit_name, branch, parent_unit_id, commander_name, personnel_count, 
		location_lat, location_lng, operational_status, equipment_summary, created_at, updated_at 
		FROM milc2_units WHERE operational_status = $1 ORDER BY unit_name`
	rows, err := r.db.Query(query, domain.OpStatusDeployed)
	if err != nil {
		return nil, fmt.Errorf("query deployed units: %w", err)
	}
	defer rows.Close()

	var units []domain.MilitaryUnit
	for rows.Next() {
		var u domain.MilitaryUnit
		if err := rows.Scan(
			&u.UnitID, &u.UnitName, &u.Branch, &u.ParentUnitID, &u.CommanderName,
			&u.PersonnelCount, &u.LocationLat, &u.LocationLng, &u.OperationalStatus,
			&u.EquipmentSummary, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan unit: %w", err)
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

func (r *PostgresRepo) CreateOperation(o domain.Operation) error {
	query := `INSERT INTO milc2_operations 
		(operation_id, operation_name, operation_type, status, commander_id, start_date, 
		 expected_end_date, operational_area, rules_of_engagement, mission_objective, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.db.Exec(query,
		o.OperationID, o.OperationName, o.OperationType, o.Status, o.CommanderID,
		o.StartDate, o.ExpectedEndDate, o.OperationalArea, o.RulesOfEngagement,
		o.MissionObjective, o.CreatedAt, o.UpdatedAt)
	return err
}

func (r *PostgresRepo) GetActiveOperations() ([]domain.Operation, error) {
	query := `SELECT operation_id, operation_name, operation_type, status, commander_id, start_date, 
		expected_end_date, operational_area, rules_of_engagement, mission_objective, created_at, updated_at 
		FROM milc2_operations WHERE status = $1 ORDER BY start_date DESC`
	rows, err := r.db.Query(query, domain.OpStatusActiveOp)
	if err != nil {
		return nil, fmt.Errorf("query active operations: %w", err)
	}
	defer rows.Close()

	var ops []domain.Operation
	for rows.Next() {
		var o domain.Operation
		if err := rows.Scan(
			&o.OperationID, &o.OperationName, &o.OperationType, &o.Status, &o.CommanderID,
			&o.StartDate, &o.ExpectedEndDate, &o.OperationalArea, &o.RulesOfEngagement,
			&o.MissionObjective, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan operation: %w", err)
		}
		ops = append(ops, o)
	}
	return ops, rows.Err()
}

func (r *PostgresRepo) CreateTacticalReport(t domain.TacticalReport) error {
	query := `INSERT INTO milc2_tactical_reports 
		(report_id, operation_id, reporting_unit_id, report_type, position_lat, position_lng, 
		 enemy_activity, civilian_interactions, casualties, detainees, equipment_status, submitted_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.db.Exec(query,
		t.ReportID, t.OperationID, t.ReportingUnitID, t.ReportType,
		t.PositionLat, t.PositionLng, t.EnemyActivity, t.CivilianInteractions,
		t.Casualties, t.Detainees, t.EquipmentStatus, t.SubmittedAt)
	return err
}

func (r *PostgresRepo) GetReportsByOperation(opID uuid.UUID) ([]domain.TacticalReport, error) {
	query := `SELECT report_id, operation_id, reporting_unit_id, report_type, position_lat, position_lng, 
		enemy_activity, civilian_interactions, casualties, detainees, equipment_status, submitted_at 
		FROM milc2_tactical_reports WHERE operation_id = $1 ORDER BY submitted_at DESC`
	rows, err := r.db.Query(query, opID)
	if err != nil {
		return nil, fmt.Errorf("query reports by operation: %w", err)
	}
	defer rows.Close()

	var reports []domain.TacticalReport
	for rows.Next() {
		var t domain.TacticalReport
		if err := rows.Scan(
			&t.ReportID, &t.OperationID, &t.ReportingUnitID, &t.ReportType,
			&t.PositionLat, &t.PositionLng, &t.EnemyActivity, &t.CivilianInteractions,
			&t.Casualties, &t.Detainees, &t.EquipmentStatus, &t.SubmittedAt); err != nil {
			return nil, fmt.Errorf("scan report: %w", err)
		}
		reports = append(reports, t)
	}
	return reports, rows.Err()
}

func (r *PostgresRepo) GetAllUnits() ([]domain.MilitaryUnit, error) {
	query := `SELECT unit_id, unit_name, branch, parent_unit_id, commander_name, personnel_count, 
		location_lat, location_lng, operational_status, equipment_summary, created_at, updated_at 
		FROM milc2_units ORDER BY unit_name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query all units: %w", err)
	}
	defer rows.Close()

	var units []domain.MilitaryUnit
	for rows.Next() {
		var u domain.MilitaryUnit
		if err := rows.Scan(
			&u.UnitID, &u.UnitName, &u.Branch, &u.ParentUnitID, &u.CommanderName,
			&u.PersonnelCount, &u.LocationLat, &u.LocationLng, &u.OperationalStatus,
			&u.EquipmentSummary, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan unit: %w", err)
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

func (r *PostgresRepo) GetAllOperations() ([]domain.Operation, error) {
	query := `SELECT operation_id, operation_name, operation_type, status, commander_id, start_date, 
		expected_end_date, operational_area, rules_of_engagement, mission_objective, created_at, updated_at 
		FROM milc2_operations ORDER BY start_date DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query all operations: %w", err)
	}
	defer rows.Close()

	var ops []domain.Operation
	for rows.Next() {
		var o domain.Operation
		if err := rows.Scan(
			&o.OperationID, &o.OperationName, &o.OperationType, &o.Status, &o.CommanderID,
			&o.StartDate, &o.ExpectedEndDate, &o.OperationalArea, &o.RulesOfEngagement,
			&o.MissionObjective, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan operation: %w", err)
		}
		ops = append(ops, o)
	}
	return ops, rows.Err()
}

func (r *PostgresRepo) GetAllReports() ([]domain.TacticalReport, error) {
	query := `SELECT report_id, operation_id, reporting_unit_id, report_type, position_lat, position_lng, 
		enemy_activity, civilian_interactions, casualties, detainees, equipment_status, submitted_at 
		FROM milc2_tactical_reports ORDER BY submitted_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query all reports: %w", err)
	}
	defer rows.Close()

	var reports []domain.TacticalReport
	for rows.Next() {
		var t domain.TacticalReport
		if err := rows.Scan(
			&t.ReportID, &t.OperationID, &t.ReportingUnitID, &t.ReportType,
			&t.PositionLat, &t.PositionLng, &t.EnemyActivity, &t.CivilianInteractions,
			&t.Casualties, &t.Detainees, &t.EquipmentStatus, &t.SubmittedAt); err != nil {
			return nil, fmt.Errorf("scan report: %w", err)
		}
		reports = append(reports, t)
	}
	return reports, rows.Err()
}
