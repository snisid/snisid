package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/sigdc-svc/internal/domain"
)

type disasterRepo struct {
	pool *pgxpool.Pool
}

func NewDisasterRepo(pool *pgxpool.Pool) *disasterRepo {
	return &disasterRepo{pool: pool}
}

func (r *disasterRepo) CreateDisaster(d *domain.Disaster) (*domain.Disaster, error) {
	ctx := context.Background()
	d.ID = uuid.New()
	d.NationalSigdcID = "SIGDC-HT-" + time.Now().Format("2006") + "-" + d.ID.String()[:6]
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	status := "ACTIVE"
	d.Status = &status

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sigdc_disasters
		 (disaster_id, national_sigdc_id, disaster_type, disaster_name, alert_level, status,
		  onset_date, affected_depts, epicenter_lat, epicenter_lng, magnitude, response_agencies,
		  created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		d.ID, d.NationalSigdcID, d.DisasterType, d.DisasterName, d.AlertLevel, d.Status,
		d.OnsetDate, d.AffectedDepts, d.EpicenterLat, d.EpicenterLng, d.Magnitude,
		d.ResponseAgencies, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (r *disasterRepo) FindActiveDisasters() ([]domain.Disaster, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT disaster_id, national_sigdc_id, disaster_type, disaster_name, alert_level, status,
		        onset_date, affected_depts, epicenter_lat, epicenter_lng, magnitude,
		        confirmed_dead, confirmed_injured, confirmed_missing, created_at
		 FROM sigdc_disasters WHERE status = 'ACTIVE' ORDER BY onset_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var disasters []domain.Disaster
	for rows.Next() {
		var d domain.Disaster
		if err := rows.Scan(&d.ID, &d.NationalSigdcID, &d.DisasterType, &d.DisasterName,
			&d.AlertLevel, &d.Status, &d.OnsetDate, &d.AffectedDepts,
			&d.EpicenterLat, &d.EpicenterLng, &d.Magnitude,
			&d.ConfirmedDead, &d.ConfirmedInjured, &d.ConfirmedMissing, &d.CreatedAt); err != nil {
			return nil, err
		}
		disasters = append(disasters, d)
	}
	return disasters, nil
}

func (r *disasterRepo) SaveWarning(w *domain.EarlyWarning) error {
	ctx := context.Background()
	w.WarningID = uuid.New()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sigdc_early_warnings
		 (warning_id, disaster_type, alert_level, source_agency, message_text, affected_depts,
		  issued_at, channels_sent)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		w.WarningID, w.DisasterType, w.AlertLevel, w.SourceAgency, w.MessageText,
		w.AffectedDepts, w.IssuedAt, w.ChannelsSent,
	)
	return err
}

func (r *disasterRepo) CreateVictimRegistration(vr *domain.VictimRegistration) (*domain.VictimRegistration, error) {
	ctx := context.Background()
	vr.RegistrationID = uuid.New()
	vr.RegistrationDate = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sigdc_victim_registrations
		 (registration_id, disaster_id, full_name, status, location_found, dept_code,
		  registration_date, registered_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		vr.RegistrationID, vr.DisasterID, vr.FullName, vr.Status,
		vr.LocationFound, vr.DeptCode, vr.RegistrationDate, vr.RegisteredBy,
	)
	if err != nil {
		return nil, err
	}

	return vr, nil
}

func (r *disasterRepo) FindResources(disasterID uuid.UUID) ([]domain.Resource, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT resource_id, disaster_id, resource_type, provider_org, quantity, dept_code, status, created_at
		 FROM sigdc_resources WHERE disaster_id = $1 AND status = 'AVAILABLE'`, disasterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []domain.Resource
	for rows.Next() {
		var res domain.Resource
		if err := rows.Scan(&res.ResourceID, &res.DisasterID, &res.ResourceType, &res.ProviderOrg,
			&res.Quantity, &res.DeptCode, &res.Status, &res.CreatedAt); err != nil {
			return nil, err
		}
		resources = append(resources, res)
	}
	return resources, nil
}
