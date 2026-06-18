package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/trait-svc/internal/domain"
)

type traffickingRepo struct {
	pool *pgxpool.Pool
}

func NewTraffickingRepo(pool *pgxpool.Pool) *traffickingRepo {
	return &traffickingRepo{pool: pool}
}

func (r *traffickingRepo) CreateCase(c *domain.TraiffickingCase) (*domain.TraiffickingCase, error) {
	ctx := context.Background()
	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO trait_cases
		 (case_id, national_trait_id, trait_type, status, victim_count, minor_count,
		  origin_country, transit_countries, destination_country, route_description,
		  transport_mode, mar_incident_id, sifr_crossing_ids, gang_id, recruiter_ids,
		  total_amount_paid, amount_per_person, currency, investigating_unit,
		  case_reference, iom_case_ref, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24)
		 RETURNING case_id, created_at, updated_at`,
		c.ID, c.NationalTraitID, c.TrafficType, c.Status, c.VictimCount, c.MinorCount,
		c.OriginCountry, c.TransitCountries, c.DestinationCountry, c.RouteDescription,
		c.TransportMode, c.MarIncidentID, c.SifrcrossingIDs, c.GangID, c.RecruiterIDs,
		c.TotalAmountPaid, c.AmountPerPerson, c.Currency, c.InvestigatingUnit,
		c.CaseReference, c.IomCaseRef, c.CreatedBy, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *traffickingRepo) FindByID(id uuid.UUID) (*domain.TraiffickingCase, error) {
	ctx := context.Background()
	c := &domain.TraiffickingCase{}
	err := r.pool.QueryRow(ctx,
		`SELECT case_id, national_trait_id, trait_type, status, victim_count, minor_count,
		        origin_country, transit_countries, destination_country, route_description,
		        transport_mode, mar_incident_id, sifr_crossing_ids, gang_id, recruiter_ids,
		        total_amount_paid, amount_per_person, currency, investigating_unit,
		        case_reference, iom_case_ref, created_by, created_at, updated_at
		 FROM trait_cases WHERE case_id = $1`, id).Scan(
		&c.ID, &c.NationalTraitID, &c.TrafficType, &c.Status, &c.VictimCount, &c.MinorCount,
		&c.OriginCountry, &c.TransitCountries, &c.DestinationCountry, &c.RouteDescription,
		&c.TransportMode, &c.MarIncidentID, &c.SifrcrossingIDs, &c.GangID, &c.RecruiterIDs,
		&c.TotalAmountPaid, &c.AmountPerPerson, &c.Currency, &c.InvestigatingUnit,
		&c.CaseReference, &c.IomCaseRef, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *traffickingRepo) FindCaseByNationalID(nationalID string) (*domain.TraiffickingCase, error) {
	ctx := context.Background()
	c := &domain.TraiffickingCase{}
	err := r.pool.QueryRow(ctx,
		`SELECT case_id, national_trait_id, trait_type, status, victim_count, minor_count,
		        origin_country, transit_countries, destination_country, route_description,
		        transport_mode, mar_incident_id, sifr_crossing_ids, gang_id, recruiter_ids,
		        total_amount_paid, amount_per_person, currency, investigating_unit,
		        case_reference, iom_case_ref, created_by, created_at, updated_at
		 FROM trait_cases WHERE national_trait_id = $1`, nationalID).Scan(
		&c.ID, &c.NationalTraitID, &c.TrafficType, &c.Status, &c.VictimCount, &c.MinorCount,
		&c.OriginCountry, &c.TransitCountries, &c.DestinationCountry, &c.RouteDescription,
		&c.TransportMode, &c.MarIncidentID, &c.SifrcrossingIDs, &c.GangID, &c.RecruiterIDs,
		&c.TotalAmountPaid, &c.AmountPerPerson, &c.Currency, &c.InvestigatingUnit,
		&c.CaseReference, &c.IomCaseRef, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *traffickingRepo) AddVictim(v *domain.TraiffickingVictim) (*domain.TraiffickingVictim, error) {
	ctx := context.Background()
	v.ID = uuid.New()
	v.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO trait_victims
		 (victim_id, case_id, snisid_person_id, victim_status, full_name, nationality,
		  dob, gender, is_minor, exploitation_type, rescue_date, rescue_location,
		  current_location, assistance_provided, dipe_case_id, afis_subject_id, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		 RETURNING victim_id, created_at`,
		v.ID, v.CaseID, v.SnisidPersonID, v.VictimStatus, v.FullName, v.Nationality,
		v.Dob, v.Gender, v.IsMinor, v.ExploitationType, v.RescueDate, v.RescueLocation,
		v.CurrentLocation, v.AssistanceProvided, v.DipeCaseID, v.AfisSubjectID, v.CreatedAt,
	).Scan(&v.ID, &v.CreatedAt)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (r *traffickingRepo) GetVictimsByCase(caseID uuid.UUID) ([]domain.TraiffickingVictim, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT victim_id, case_id, snisid_person_id, victim_status, full_name, nationality,
		        dob, gender, is_minor, exploitation_type, rescue_date, rescue_location,
		        current_location, assistance_provided, dipe_case_id, afis_subject_id, created_at
		 FROM trait_victims
		 WHERE case_id = $1
		 ORDER BY created_at DESC`, caseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanVictims(rows)
}

func (r *traffickingRepo) GetMinorVictims() ([]domain.TraiffickingVictim, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT v.victim_id, v.case_id, v.snisid_person_id, v.victim_status, v.full_name, v.nationality,
		        v.dob, v.gender, v.is_minor, v.exploitation_type, v.rescue_date, v.rescue_location,
		        v.current_location, v.assistance_provided, v.dipe_case_id, v.afis_subject_id, v.created_at
		 FROM trait_victims v
		 WHERE v.is_minor = TRUE
		 ORDER BY v.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanVictims(rows)
}

func (r *traffickingRepo) CreateNetwork(n *domain.TraiffickingNetwork) (*domain.TraiffickingNetwork, error) {
	ctx := context.Background()
	n.ID = uuid.New()
	n.IsActive = true
	n.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO trait_networks
		 (network_id, network_name, primary_route, origin_dept, known_members,
		  gang_affiliations, monthly_volume_est, fee_per_person_usd, is_active,
		  intel_confidence, linked_cases, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		 RETURNING network_id, created_at`,
		n.ID, n.NetworkName, n.PrimaryRoute, n.OriginDept, n.KnownMembers,
		n.GangAffiliations, n.MonthlyVolumeEst, n.FeePerPersonUsd, n.IsActive,
		n.IntelConfidence, n.LinkedCases, n.CreatedBy, n.CreatedAt,
	).Scan(&n.ID, &n.CreatedAt)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (r *traffickingRepo) GetActiveNetworks() ([]domain.TraiffickingNetwork, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT network_id, network_name, primary_route, origin_dept, known_members,
		        gang_affiliations, monthly_volume_est, fee_per_person_usd, is_active,
		        intel_confidence, linked_cases, created_by, created_at
		 FROM trait_networks
		 WHERE is_active = TRUE
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNetworks(rows)
}

func (r *traffickingRepo) GetStatsByType() ([]domain.TypeStats, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT trait_type, COUNT(*)::int AS count, COALESCE(SUM(minor_count), 0)::int AS minor_count
		 FROM trait_cases
		 GROUP BY trait_type
		 ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []domain.TypeStats
	for rows.Next() {
		var s domain.TypeStats
		if err := rows.Scan(&s.TrafficType, &s.Count, &s.MinorCount); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func (r *traffickingRepo) GetMaritimeCases() ([]domain.TraiffickingCase, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT case_id, national_trait_id, trait_type, status, victim_count, minor_count,
		        origin_country, transit_countries, destination_country, route_description,
		        transport_mode, mar_incident_id, sifr_crossing_ids, gang_id, recruiter_ids,
		        total_amount_paid, amount_per_person, currency, investigating_unit,
		        case_reference, iom_case_ref, created_by, created_at, updated_at
		 FROM trait_cases
		 WHERE mar_incident_id IS NOT NULL
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCases(rows)
}

func (r *traffickingRepo) CountCasesByTypeAndYear(traitType domain.TraiffickingType, year int) (int, error) {
	ctx := context.Background()
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*)::int FROM trait_cases
		 WHERE trait_type = $1 AND EXTRACT(YEAR FROM created_at) = $2`,
		traitType, year).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func scanCases(rows pgx.Rows) ([]domain.TraiffickingCase, error) {
	var cases []domain.TraiffickingCase
	for rows.Next() {
		var c domain.TraiffickingCase
		if err := rows.Scan(
			&c.ID, &c.NationalTraitID, &c.TrafficType, &c.Status, &c.VictimCount, &c.MinorCount,
			&c.OriginCountry, &c.TransitCountries, &c.DestinationCountry, &c.RouteDescription,
			&c.TransportMode, &c.MarIncidentID, &c.SifrcrossingIDs, &c.GangID, &c.RecruiterIDs,
			&c.TotalAmountPaid, &c.AmountPerPerson, &c.Currency, &c.InvestigatingUnit,
			&c.CaseReference, &c.IomCaseRef, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, nil
}

func scanVictims(rows pgx.Rows) ([]domain.TraiffickingVictim, error) {
	var victims []domain.TraiffickingVictim
	for rows.Next() {
		var v domain.TraiffickingVictim
		if err := rows.Scan(
			&v.ID, &v.CaseID, &v.SnisidPersonID, &v.VictimStatus, &v.FullName, &v.Nationality,
			&v.Dob, &v.Gender, &v.IsMinor, &v.ExploitationType, &v.RescueDate, &v.RescueLocation,
			&v.CurrentLocation, &v.AssistanceProvided, &v.DipeCaseID, &v.AfisSubjectID, &v.CreatedAt); err != nil {
			return nil, err
		}
		victims = append(victims, v)
	}
	return victims, nil
}

func scanNetworks(rows pgx.Rows) ([]domain.TraiffickingNetwork, error) {
	var networks []domain.TraiffickingNetwork
	for rows.Next() {
		var n domain.TraiffickingNetwork
		if err := rows.Scan(
			&n.ID, &n.NetworkName, &n.PrimaryRoute, &n.OriginDept, &n.KnownMembers,
			&n.GangAffiliations, &n.MonthlyVolumeEst, &n.FeePerPersonUsd, &n.IsActive,
			&n.IntelConfidence, &n.LinkedCases, &n.CreatedBy, &n.CreatedAt); err != nil {
			return nil, err
		}
		networks = append(networks, n)
	}
	return networks, nil
}
