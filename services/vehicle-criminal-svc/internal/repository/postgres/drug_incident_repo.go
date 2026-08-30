package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

type DrugIncidentRepo struct {
	db *sqlx.DB
}

func NewDrugIncidentRepo(db *sqlx.DB) *DrugIncidentRepo {
	return &DrugIncidentRepo{db: db}
}

func (r *DrugIncidentRepo) Create(ctx context.Context, incident *domain.DrugIncident) error {
	query := `
		INSERT INTO sivc_drug_incidents (
			incident_id, alert_id, drug_types, seizure_weight_kg, estimated_value_usd,
			seizure_date, seizure_location, seizure_dept_code, seizure_commune,
			route_type, origin_country, transit_points, destination,
			suspected_cartel, blts_case_number, interpol_ref, concealment_method,
			notes, created_by, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		incident.IncidentID, incident.AlertID, incident.DrugTypes, incident.SeizureWeightKg,
		incident.EstimatedValueUSD, incident.SeizureDate, incident.SeizureLocation,
		incident.SeizureDeptCode, incident.SeizureCommune, incident.RouteType,
		incident.OriginCountry, incident.TransitPoints, incident.Destination,
		incident.SuspectedCartel, incident.BltsCaseNumber, incident.InterpolRef,
		incident.ConcealmentMethod, incident.Notes, incident.CreatedBy, incident.CreatedAt,
	)
	return err
}

func (r *DrugIncidentRepo) FindByAlertID(ctx context.Context, alertID uuid.UUID) (*domain.DrugIncident, error) {
	var incident domain.DrugIncident
	query := `SELECT * FROM sivc_drug_incidents WHERE alert_id = $1 LIMIT 1`
	if err := r.db.GetContext(ctx, &incident, query, alertID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &incident, nil
}
