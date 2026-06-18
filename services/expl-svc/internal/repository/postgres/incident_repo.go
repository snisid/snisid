package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/expl-svc/internal/domain"
)

type incidentRepo struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewIncidentRepository(db *pgxpool.Pool, log *zap.Logger) *incidentRepo {
	return &incidentRepo{db: db, log: log}
}

func (r *incidentRepo) CreateIncident(incident *domain.ExplIncident) error {
	ctx := context.Background()
	_, err := r.db.Exec(ctx,
		`INSERT INTO expl_incidents (
			national_expl_id, incident_type, explosive_type, status,
			quantity, weight_kg, manufacturer, lot_number, manufacture_country,
			estimated_date, incident_date, location_desc, dept_code, commune,
			lat, lng, responding_unit, eod_officer, casualties,
			gang_id, from_person_id, case_reference, dna_sample_taken, bio_sample_ref,
			photo_refs, interpol_exploint_ref, notes, created_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28)`,
		incident.NationalExplID, incident.IncidentType, incident.ExplosiveType,
		incident.Status, incident.Quantity, incident.WeightKg,
		incident.Manufacturer, incident.LotNumber, incident.ManufactureCountry,
		incident.EstimatedDate, incident.IncidentDate, incident.LocationDesc,
		incident.DeptCode, incident.Commune, incident.Lat, incident.Lng,
		incident.RespondingUnit, incident.EODOfficer, incident.Casualties,
		incident.GangID, incident.FromPersonID, incident.CaseReference,
		incident.DNASampleTaken, incident.BioSampleRef, incident.PhotoRefs,
		incident.InterpolExplointRef, incident.Notes, incident.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("create incident: %w", err)
	}
	return nil
}

func (r *incidentRepo) FindByID(id uuid.UUID) (*domain.ExplIncident, error) {
	ctx := context.Background()
	incident := &domain.ExplIncident{}
	err := r.db.QueryRow(ctx,
		`SELECT incident_id, national_expl_id, incident_type, explosive_type, status,
			quantity, weight_kg, manufacturer, lot_number, manufacture_country,
			estimated_date, incident_date, location_desc, dept_code, commune,
			lat, lng, responding_unit, eod_officer, casualties,
			gang_id, from_person_id, case_reference, dna_sample_taken, bio_sample_ref,
			photo_refs, interpol_exploint_ref, notes, created_by, created_at
		 FROM expl_incidents WHERE incident_id = $1`, id).Scan(
		&incident.IncidentID, &incident.NationalExplID, &incident.IncidentType,
		&incident.ExplosiveType, &incident.Status, &incident.Quantity, &incident.WeightKg,
		&incident.Manufacturer, &incident.LotNumber, &incident.ManufactureCountry,
		&incident.EstimatedDate, &incident.IncidentDate, &incident.LocationDesc,
		&incident.DeptCode, &incident.Commune, &incident.Lat, &incident.Lng,
		&incident.RespondingUnit, &incident.EODOfficer, &incident.Casualties,
		&incident.GangID, &incident.FromPersonID, &incident.CaseReference,
		&incident.DNASampleTaken, &incident.BioSampleRef, &incident.PhotoRefs,
		&incident.InterpolExplointRef, &incident.Notes, &incident.CreatedBy, &incident.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find incident by id: %w", err)
	}
	return incident, nil
}

func (r *incidentRepo) FindByDept(deptCode string, limit, offset int) ([]domain.ExplIncident, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT incident_id, national_expl_id, incident_type, explosive_type, status,
			quantity, weight_kg, manufacturer, lot_number, manufacture_country,
			estimated_date, incident_date, location_desc, dept_code, commune,
			lat, lng, responding_unit, eod_officer, casualties,
			gang_id, from_person_id, case_reference, dna_sample_taken, bio_sample_ref,
			photo_refs, interpol_exploint_ref, notes, created_by, created_at
		 FROM expl_incidents WHERE dept_code = $1
		 ORDER BY incident_date DESC
		 LIMIT $2 OFFSET $3`, deptCode, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("find incidents by dept: %w", err)
	}
	defer rows.Close()

	var incidents []domain.ExplIncident
	for rows.Next() {
		var inc domain.ExplIncident
		if err := rows.Scan(
			&inc.IncidentID, &inc.NationalExplID, &inc.IncidentType,
			&inc.ExplosiveType, &inc.Status, &inc.Quantity, &inc.WeightKg,
			&inc.Manufacturer, &inc.LotNumber, &inc.ManufactureCountry,
			&inc.EstimatedDate, &inc.IncidentDate, &inc.LocationDesc,
			&inc.DeptCode, &inc.Commune, &inc.Lat, &inc.Lng,
			&inc.RespondingUnit, &inc.EODOfficer, &inc.Casualties,
			&inc.GangID, &inc.FromPersonID, &inc.CaseReference,
			&inc.DNASampleTaken, &inc.BioSampleRef, &inc.PhotoRefs,
			&inc.InterpolExplointRef, &inc.Notes, &inc.CreatedBy, &inc.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan incident row: %w", err)
		}
		incidents = append(incidents, inc)
	}
	return incidents, nil
}

func (r *incidentRepo) CreateLegalStock(stock *domain.LegalStock) error {
	ctx := context.Background()
	_, err := r.db.Exec(ctx,
		`INSERT INTO expl_legal_stocks (
			holder_entity, holder_type, explosive_type, quantity_kg,
			storage_location, dept_code, license_ref, last_audit_date,
			next_audit_date, is_secured
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		stock.HolderEntity, stock.HolderType, stock.ExplosiveType,
		stock.QuantityKg, stock.StorageLocation, stock.DeptCode,
		stock.LicenseRef, stock.LastAuditDate, stock.NextAuditDate, stock.IsSecured,
	)
	if err != nil {
		return fmt.Errorf("create legal stock: %w", err)
	}
	return nil
}

func (r *incidentRepo) GetLegalStocks(deptCode string, limit, offset int) ([]domain.LegalStock, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT stock_id, holder_entity, holder_type, explosive_type, quantity_kg,
			storage_location, dept_code, license_ref, last_audit_date,
			next_audit_date, is_secured, created_at
		 FROM expl_legal_stocks WHERE dept_code = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`, deptCode, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get legal stocks: %w", err)
	}
	defer rows.Close()

	var stocks []domain.LegalStock
	for rows.Next() {
		var s domain.LegalStock
		if err := rows.Scan(
			&s.StockID, &s.HolderEntity, &s.HolderType, &s.ExplosiveType,
			&s.QuantityKg, &s.StorageLocation, &s.DeptCode, &s.LicenseRef,
			&s.LastAuditDate, &s.NextAuditDate, &s.IsSecured, &s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan legal stock row: %w", err)
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}

func (r *incidentRepo) CountIncidents() (int, error) {
	ctx := context.Background()
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM expl_incidents`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count incidents: %w", err)
	}
	return count, nil
}
