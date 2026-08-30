package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/aero-svc/internal/domain"
)

type aircraftRepo struct {
	pool *pgxpool.Pool
}

func NewAircraftRepo(pool *pgxpool.Pool) *aircraftRepo {
	return &aircraftRepo{pool: pool}
}

func (r *aircraftRepo) CreateAircraft(aircraft *domain.Aircraft) (*domain.Aircraft, error) {
	ctx := context.Background()
	aircraft.AircraftID = uuid.New()
	aircraft.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO aero_aircraft_registry
		 (aircraft_id, registration_mark, icao_hex_code, aircraft_type, make, model,
		  manufacture_year, flag_country, owner_name, owner_snisid_id, operator_name,
		  is_registered, is_suspected, is_stolen, gang_id, drug_trafficking,
		  interpol_ref, faa_registry_ref, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)
		 RETURNING aircraft_id, created_at`,
		aircraft.AircraftID, aircraft.RegistrationMark, aircraft.ICAOHexCode,
		aircraft.AircraftType, aircraft.Make, aircraft.Model,
		aircraft.ManufactureYear, aircraft.FlagCountry, aircraft.OwnerName,
		aircraft.OwnerSnisidID, aircraft.OperatorName, aircraft.IsRegistered,
		aircraft.IsSuspected, aircraft.IsStolen, aircraft.GangID,
		aircraft.DrugTrafficking, aircraft.InterpolRef, aircraft.FAARegistryRef,
		aircraft.Notes, aircraft.CreatedAt,
	).Scan(&aircraft.AircraftID, &aircraft.CreatedAt)
	if err != nil {
		return nil, err
	}

	return aircraft, nil
}

func (r *aircraftRepo) FindByRegistration(mark string) (*domain.Aircraft, error) {
	ctx := context.Background()
	aircraft := &domain.Aircraft{}

	err := r.pool.QueryRow(ctx,
		`SELECT aircraft_id, registration_mark, icao_hex_code, aircraft_type, make, model,
		        manufacture_year, flag_country, owner_name, owner_snisid_id, operator_name,
		        is_registered, is_suspected, is_stolen, gang_id, drug_trafficking,
		        interpol_ref, faa_registry_ref, notes, created_at
		 FROM aero_aircraft_registry
		 WHERE registration_mark = $1`, mark).Scan(
		&aircraft.AircraftID, &aircraft.RegistrationMark, &aircraft.ICAOHexCode,
		&aircraft.AircraftType, &aircraft.Make, &aircraft.Model,
		&aircraft.ManufactureYear, &aircraft.FlagCountry, &aircraft.OwnerName,
		&aircraft.OwnerSnisidID, &aircraft.OperatorName, &aircraft.IsRegistered,
		&aircraft.IsSuspected, &aircraft.IsStolen, &aircraft.GangID,
		&aircraft.DrugTrafficking, &aircraft.InterpolRef, &aircraft.FAARegistryRef,
		&aircraft.Notes, &aircraft.CreatedAt)
	if err != nil {
		return nil, err
	}

	return aircraft, nil
}

func (r *aircraftRepo) CreateStrip(strip *domain.ClandestineStrip) (*domain.ClandestineStrip, error) {
	ctx := context.Background()
	strip.StripID = uuid.New()
	strip.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO aero_clandestine_strips
		 (strip_id, strip_name, dept_code, commune, lat, lng, length_m, surface_type,
		  status, capable_aircraft, gang_id, first_detected, last_activity_date,
		  source_intel, satellite_image_ref, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		 RETURNING strip_id, created_at`,
		strip.StripID, strip.StripName, strip.DeptCode, strip.Commune,
		strip.Lat, strip.Lng, strip.LengthM, strip.SurfaceType,
		strip.Status, strip.CapableAircraft, strip.GangID,
		strip.FirstDetected, strip.LastActivityDate, strip.SourceIntel,
		strip.SatelliteImageRef, strip.CreatedBy, strip.CreatedAt,
	).Scan(&strip.StripID, &strip.CreatedAt)
	if err != nil {
		return nil, err
	}

	return strip, nil
}

func (r *aircraftRepo) FindActiveStrips() ([]domain.ClandestineStrip, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT strip_id, strip_name, dept_code, commune, lat, lng, length_m, surface_type,
		        status, capable_aircraft, gang_id, first_detected, last_activity_date,
		        source_intel, satellite_image_ref, created_by, created_at
		 FROM aero_clandestine_strips
		 WHERE status = 'ACTIVE'
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanStrips(rows)
}

func (r *aircraftRepo) GetStripsMap() ([]domain.ClandestineStrip, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT strip_id, strip_name, dept_code, commune, lat, lng, length_m, surface_type,
		        status, capable_aircraft, gang_id, first_detected, last_activity_date,
		        source_intel, satellite_image_ref, created_by, created_at
		 FROM aero_clandestine_strips
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanStrips(rows)
}

func (r *aircraftRepo) CreateFlight(flight *domain.SuspiciousFlight) (*domain.SuspiciousFlight, error) {
	ctx := context.Background()
	flight.FlightID = uuid.New()
	flight.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO aero_suspicious_flights
		 (flight_id, aircraft_id, registration_mark, flight_date, origin_airport,
		  destination_airport, origin_country, destination_country, landing_strip_id,
		  landing_location, flight_type, cargo_suspected, source_radar, source_informant,
		  case_reference, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		 RETURNING flight_id, created_at`,
		flight.FlightID, flight.AircraftID, flight.RegistrationMark,
		flight.FlightDate, flight.OriginAirport, flight.DestinationAirport,
		flight.OriginCountry, flight.DestinationCountry, flight.LandingStripID,
		flight.LandingLocation, flight.FlightType, flight.CargoSuspected,
		flight.SourceRadar, flight.SourceInformant, flight.CaseReference,
		flight.CreatedBy, flight.CreatedAt,
	).Scan(&flight.FlightID, &flight.CreatedAt)
	if err != nil {
		return nil, err
	}

	return flight, nil
}

func (r *aircraftRepo) GetFlightsByDate(from, to time.Time) ([]domain.SuspiciousFlight, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT flight_id, aircraft_id, registration_mark, flight_date, origin_airport,
		        destination_airport, origin_country, destination_country, landing_strip_id,
		        landing_location, flight_type, cargo_suspected, source_radar, source_informant,
		        case_reference, created_by, created_at
		 FROM aero_suspicious_flights
		 WHERE flight_date BETWEEN $1 AND $2
		 ORDER BY flight_date DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanFlights(rows)
}

func (r *aircraftRepo) GetStripStats() (*domain.StripStats, error) {
	ctx := context.Background()
	stats := &domain.StripStats{
		ByDepartment: make(map[string]int),
		ByStatus:     make(map[string]int),
	}

	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM aero_clandestine_strips`).Scan(&stats.TotalStrips)
	if err != nil {
		return nil, err
	}

	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM aero_clandestine_strips WHERE status = 'ACTIVE'`).Scan(&stats.ActiveStrips)
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT dept_code, COUNT(*) FROM aero_clandestine_strips WHERE status = 'ACTIVE' GROUP BY dept_code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dept string
		var count int
		if err := rows.Scan(&dept, &count); err != nil {
			return nil, err
		}
		stats.ByDepartment[dept] = count
	}

	rows2, err := r.pool.Query(ctx,
		`SELECT status::text, COUNT(*) FROM aero_clandestine_strips GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var status string
		var count int
		if err := rows2.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.ByStatus[status] = count
	}

	return stats, nil
}

func scanStrips(rows pgx.Rows) ([]domain.ClandestineStrip, error) {
	var strips []domain.ClandestineStrip
	for rows.Next() {
		var s domain.ClandestineStrip
		if err := rows.Scan(
			&s.StripID, &s.StripName, &s.DeptCode, &s.Commune, &s.Lat, &s.Lng,
			&s.LengthM, &s.SurfaceType, &s.Status, &s.CapableAircraft, &s.GangID,
			&s.FirstDetected, &s.LastActivityDate, &s.SourceIntel,
			&s.SatelliteImageRef, &s.CreatedBy, &s.CreatedAt); err != nil {
			return nil, err
		}
		strips = append(strips, s)
	}
	return strips, nil
}

func scanFlights(rows pgx.Rows) ([]domain.SuspiciousFlight, error) {
	var flights []domain.SuspiciousFlight
	for rows.Next() {
		var f domain.SuspiciousFlight
		if err := rows.Scan(
			&f.FlightID, &f.AircraftID, &f.RegistrationMark, &f.FlightDate,
			&f.OriginAirport, &f.DestinationAirport, &f.OriginCountry,
			&f.DestinationCountry, &f.LandingStripID, &f.LandingLocation,
			&f.FlightType, &f.CargoSuspected, &f.SourceRadar, &f.SourceInformant,
			&f.CaseReference, &f.CreatedBy, &f.CreatedAt); err != nil {
			return nil, err
		}
		flights = append(flights, f)
	}
	return flights, nil
}
