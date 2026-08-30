package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

type KidnappingIncidentRepo struct {
	db *sqlx.DB
}

func NewKidnappingIncidentRepo(db *sqlx.DB) *KidnappingIncidentRepo {
	return &KidnappingIncidentRepo{db: db}
}

func (r *KidnappingIncidentRepo) Create(ctx context.Context, incident *domain.KidnappingIncident) error {
	query := `
		INSERT INTO sivc_kidnapping_incidents (
			incident_id, alert_id, victim_count, victim_snisid_ids,
			victims_nationality, victims_description, abduction_date,
			abduction_location, abduction_dept_code, abduction_commune,
			abduction_context, ransom_demanded, ransom_amount, ransom_currency,
			ransom_channel, incident_status, cae_case_number, dcpj_case_number,
			created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		incident.IncidentID, incident.AlertID, incident.VictimCount, incident.VictimSnisidIDs,
		incident.VictimsNationality, incident.VictimsDescription, incident.AbductionDate,
		incident.AbductionLocation, incident.AbductionDeptCode, incident.AbductionCommune,
		incident.AbductionContext, incident.RansomDemanded, incident.RansomAmount,
		incident.RansomCurrency, incident.RansomChannel, incident.IncidentStatus,
		incident.CaeCaseNumber, incident.DcpjCaseNumber, incident.CreatedBy,
		incident.CreatedAt, incident.UpdatedAt,
	)
	return err
}

func (r *KidnappingIncidentRepo) FindByAlertID(ctx context.Context, alertID uuid.UUID) (*domain.KidnappingIncident, error) {
	var incident domain.KidnappingIncident
	query := `SELECT * FROM sivc_kidnapping_incidents WHERE alert_id = $1 LIMIT 1`
	if err := r.db.GetContext(ctx, &incident, query, alertID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &incident, nil
}

func (r *KidnappingIncidentRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.KidnappingStatus) error {
	query := `UPDATE sivc_kidnapping_incidents SET incident_status = $1, updated_at = NOW() WHERE incident_id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
