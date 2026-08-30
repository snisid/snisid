package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/snisid/transport-security-ht/internal/domain"
)

type TransportRepository interface {
	CreateScreening(ctx context.Context, s *domain.PassengerScreening) error
	GetRecentScreenings(ctx context.Context, limit int) ([]domain.PassengerScreening, error)
	AddNoFly(ctx context.Context, p *domain.NoFlyPassenger) error
	CheckNoFly(ctx context.Context, identityRef string) (*domain.NoFlyPassenger, error)
	GetZonesByAirport(ctx context.Context, airportCode string) ([]domain.AirportSecurityZone, error)
	ReportZoneBreach(ctx context.Context, zoneID uuid.UUID) error
}

type transportRepo struct {
	db *sql.DB
}

func NewTransportRepository(db *sql.DB) TransportRepository {
	return &transportRepo{db: db}
}

func (r *transportRepo) CreateScreening(ctx context.Context, s *domain.PassengerScreening) error {
	query := `INSERT INTO transport_passenger_screenings
		(screening_id, traveler_identity_ref, document_type, document_number, nationality,
		 travel_mode, screening_point_type, screening_point_name, flight_number, vessel_name,
		 departure_at, arrival_at, watchlist_match, watchlist_ref, screening_result, screening_officer, screened_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.db.ExecContext(ctx, query,
		s.ScreeningID, s.TravelerIdentityRef, s.DocumentType, s.DocumentNumber, s.Nationality,
		s.TravelMode, s.ScreeningPointType, s.ScreeningPointName, s.FlightNumber, s.VesselName,
		s.DepartureAt, s.ArrivalAt, s.WatchlistMatch, s.WatchlistRef, s.ScreeningResult, s.ScreeningOfficer, s.ScreenedAt)
	return err
}

func (r *transportRepo) GetRecentScreenings(ctx context.Context, limit int) ([]domain.PassengerScreening, error) {
	query := `SELECT screening_id, traveler_identity_ref, document_type, document_number, nationality,
		travel_mode, screening_point_type, screening_point_name, flight_number, vessel_name,
		departure_at, arrival_at, watchlist_match, watchlist_ref, screening_result, screening_officer, screened_at
		FROM transport_passenger_screenings ORDER BY screened_at DESC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.PassengerScreening
	for rows.Next() {
		var s domain.PassengerScreening
		err := rows.Scan(&s.ScreeningID, &s.TravelerIdentityRef, &s.DocumentType, &s.DocumentNumber,
			&s.Nationality, &s.TravelMode, &s.ScreeningPointType, &s.ScreeningPointName,
			&s.FlightNumber, &s.VesselName, &s.DepartureAt, &s.ArrivalAt, &s.WatchlistMatch,
			&s.WatchlistRef, &s.ScreeningResult, &s.ScreeningOfficer, &s.ScreenedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if result == nil {
		result = []domain.PassengerScreening{}
	}
	return result, rows.Err()
}

func (r *transportRepo) AddNoFly(ctx context.Context, p *domain.NoFlyPassenger) error {
	query := `INSERT INTO transport_no_fly (identity_ref, list_type, added_by, reason, court_order_ref, expires_at, interpol_ref)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.ExecContext(ctx, query,
		p.IdentityRef, p.ListType, p.AddedBy, p.Reason, p.CourtOrderRef, p.ExpiresAt, p.InterpolRef)
	return err
}

func (r *transportRepo) CheckNoFly(ctx context.Context, identityRef string) (*domain.NoFlyPassenger, error) {
	query := `SELECT identity_ref, list_type, added_by, reason, court_order_ref, expires_at, interpol_ref
		FROM transport_no_fly WHERE identity_ref = $1 AND (expires_at IS NULL OR expires_at > $2) LIMIT 1`
	p := &domain.NoFlyPassenger{}
	err := r.db.QueryRowContext(ctx, query, identityRef, time.Now()).Scan(
		&p.IdentityRef, &p.ListType, &p.AddedBy, &p.Reason, &p.CourtOrderRef, &p.ExpiresAt, &p.InterpolRef)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *transportRepo) GetZonesByAirport(ctx context.Context, airportCode string) ([]domain.AirportSecurityZone, error) {
	query := `SELECT zone_id, airport_code, zone_name, zone_type, access_level, camera_count, last_inspected_at, status
		FROM transport_airport_zones WHERE airport_code = $1 ORDER BY zone_name`
	rows, err := r.db.QueryContext(ctx, query, airportCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.AirportSecurityZone
	for rows.Next() {
		var z domain.AirportSecurityZone
		if err := rows.Scan(&z.ZoneID, &z.AirportCode, &z.ZoneName, &z.ZoneType, &z.AccessLevel,
			&z.CameraCount, &z.LastInspectedAt, &z.Status); err != nil {
			return nil, err
		}
		result = append(result, z)
	}
	if result == nil {
		result = []domain.AirportSecurityZone{}
	}
	return result, rows.Err()
}

func (r *transportRepo) ReportZoneBreach(ctx context.Context, zoneID uuid.UUID) error {
	query := `UPDATE transport_airport_zones SET status = 'BREACHED' WHERE zone_id = $1`
	res, err := r.db.ExecContext(ctx, query, zoneID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

var _ TransportRepository = (*transportRepo)(nil)
var _ = pq.Array
