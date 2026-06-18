package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/vict-svc/internal/domain"
)

type victimRepo struct {
	pool *pgxpool.Pool
}

func NewVictimRepo(pool *pgxpool.Pool) *victimRepo {
	return &victimRepo{pool: pool}
}

func (r *victimRepo) Create(victim *domain.Victim) (*domain.Victim, error) {
	ctx := context.Background()
	victim.ID = uuid.New()
	victim.NationalVictID = "VICT-HT-" + time.Now().Format("2006") + "-" + victim.ID.String()[:6]
	victim.CreatedAt = time.Now()
	victim.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO vict_victims
		 (victim_id, national_vict_id, crime_type, victim_status, full_name, dob, gender,
		  incident_date, incident_location, dept_code, commune, gang_id, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		victim.ID, victim.NationalVictID, victim.CrimeType, victim.VictimStatus,
		victim.FullName, victim.DOB, victim.Gender, victim.IncidentDate,
		victim.IncidentLocation, victim.DeptCode, victim.Commune, victim.GangID,
		victim.CreatedBy, victim.CreatedAt, victim.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return victim, nil
}

func (r *victimRepo) FindByID(id uuid.UUID) (*domain.Victim, error) {
	ctx := context.Background()
	v := &domain.Victim{}
	err := r.pool.QueryRow(ctx,
		`SELECT victim_id, national_vict_id, crime_type, victim_status, full_name, dob, gender,
		        incident_date, incident_location, dept_code, commune, gang_id, created_at, updated_at
		 FROM vict_victims WHERE victim_id = $1`, id).Scan(
		&v.ID, &v.NationalVictID, &v.CrimeType, &v.VictimStatus, &v.FullName, &v.DOB,
		&v.Gender, &v.IncidentDate, &v.IncidentLocation, &v.DeptCode, &v.Commune,
		&v.GangID, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (r *victimRepo) CreateMassIncident(mi *domain.MassIncident) (*domain.MassIncident, error) {
	ctx := context.Background()
	mi.ID = uuid.New()
	mi.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO vict_mass_incidents
		 (mass_id, incident_name, crime_type, incident_date, dept_code, commune,
		  victim_count, description, documented_by, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		mi.ID, mi.IncidentName, mi.CrimeType, mi.IncidentDate, mi.DeptCode,
		mi.Commune, mi.VictimCount, mi.Description, mi.DocumentedBy,
		mi.CreatedBy, mi.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return mi, nil
}

func (r *victimRepo) FindMassIncidents() ([]domain.MassIncident, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT mass_id, incident_name, crime_type, incident_date, dept_code, commune,
		        victim_count, survivor_count, description, documented_by, created_at
		 FROM vict_mass_incidents ORDER BY incident_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []domain.MassIncident
	for rows.Next() {
		var mi domain.MassIncident
		if err := rows.Scan(&mi.ID, &mi.IncidentName, &mi.CrimeType, &mi.IncidentDate,
			&mi.DeptCode, &mi.Commune, &mi.VictimCount, &mi.SurvivorCount,
			&mi.Description, &mi.DocumentedBy, &mi.CreatedAt); err != nil {
			return nil, err
		}
		incidents = append(incidents, mi)
	}
	return incidents, nil
}

func (r *victimRepo) FindByGang(gangID uuid.UUID) ([]domain.Victim, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT victim_id, national_vict_id, crime_type, victim_status, full_name, incident_date, dept_code
		 FROM vict_victims WHERE gang_id = $1 ORDER BY incident_date DESC`, gangID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var victims []domain.Victim
	for rows.Next() {
		var v domain.Victim
		if err := rows.Scan(&v.ID, &v.NationalVictID, &v.CrimeType, &v.VictimStatus,
			&v.FullName, &v.IncidentDate, &v.DeptCode); err != nil {
			return nil, err
		}
		victims = append(victims, v)
	}
	return victims, nil
}

func (r *victimRepo) GetStatsByType() ([]domain.CrimeStats, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT crime_type, COUNT(*) as count FROM vict_victims GROUP BY crime_type ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []domain.CrimeStats
	for rows.Next() {
		var s domain.CrimeStats
		if err := rows.Scan(&s.CrimeType, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func (r *victimRepo) GetReparationList() ([]domain.Victim, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT victim_id, national_vict_id, crime_type, victim_status, full_name, incident_date, dept_code
		 FROM vict_victims WHERE needs_reparation = TRUE OR victim_status IN ('DECEASED_IDENTIFIED','ALIVE_SURVIVOR')
		 ORDER BY incident_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var victims []domain.Victim
	for rows.Next() {
		var v domain.Victim
		if err := rows.Scan(&v.ID, &v.NationalVictID, &v.CrimeType, &v.VictimStatus,
			&v.FullName, &v.IncidentDate, &v.DeptCode); err != nil {
			return nil, err
		}
		victims = append(victims, v)
	}
	return victims, nil
}
