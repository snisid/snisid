package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/air-defense-ht/internal/domain"
)

type AirDefenseRepo interface {
	CreateRadarContact(c domain.RadarContact) error
	GetActiveContacts() ([]domain.RadarContact, error)
	GetContactByID(id uuid.UUID) (*domain.RadarContact, error)
	CreateIncident(i domain.AirDefenseIncident) error
	ResolveIncident(id uuid.UUID) error
	CreateNoFlyEntry(e domain.NoFlyListEntry) error
	GetNoFlyEntry(identity string) (*domain.NoFlyListEntry, error)
}

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) CreateRadarContact(c domain.RadarContact) error {
	query := `INSERT INTO airdef_radar_contacts 
		(contact_id, track_number, contact_type, latitude, longitude, altitude_m, speed_kmh, 
		 heading_deg, source_radar, identified, squawk_code, flight_plan_ref, threat_assessment, 
		 operator_notes, first_detected_at, last_updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
	_, err := r.db.Exec(query,
		c.ContactID, c.TrackNumber, c.ContactType, c.Latitude, c.Longitude,
		c.AltitudeM, c.SpeedKmh, c.HeadingDeg, c.SourceRadar, c.Identified,
		c.SquawkCode, c.FlightPlanRef, c.ThreatAssessment, c.OperatorNotes,
		c.FirstDetectedAt, c.LastUpdatedAt)
	return err
}

func (r *PostgresRepo) GetActiveContacts() ([]domain.RadarContact, error) {
	query := `SELECT contact_id, track_number, contact_type, latitude, longitude, altitude_m, 
		speed_kmh, heading_deg, source_radar, identified, squawk_code, flight_plan_ref, 
		threat_assessment, operator_notes, first_detected_at, last_updated_at 
		FROM airdef_radar_contacts WHERE last_updated_at > $1 ORDER BY last_updated_at DESC`
	rows, err := r.db.Query(query, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("query active contacts: %w", err)
	}
	defer rows.Close()

	var contacts []domain.RadarContact
	for rows.Next() {
		var c domain.RadarContact
		if err := rows.Scan(
			&c.ContactID, &c.TrackNumber, &c.ContactType, &c.Latitude, &c.Longitude,
			&c.AltitudeM, &c.SpeedKmh, &c.HeadingDeg, &c.SourceRadar, &c.Identified,
			&c.SquawkCode, &c.FlightPlanRef, &c.ThreatAssessment, &c.OperatorNotes,
			&c.FirstDetectedAt, &c.LastUpdatedAt); err != nil {
			return nil, fmt.Errorf("scan contact: %w", err)
		}
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}

func (r *PostgresRepo) GetContactByID(id uuid.UUID) (*domain.RadarContact, error) {
	query := `SELECT contact_id, track_number, contact_type, latitude, longitude, altitude_m, 
		speed_kmh, heading_deg, source_radar, identified, squawk_code, flight_plan_ref, 
		threat_assessment, operator_notes, first_detected_at, last_updated_at 
		FROM airdef_radar_contacts WHERE contact_id = $1`
	var c domain.RadarContact
	err := r.db.QueryRow(query, id).Scan(
		&c.ContactID, &c.TrackNumber, &c.ContactType, &c.Latitude, &c.Longitude,
		&c.AltitudeM, &c.SpeedKmh, &c.HeadingDeg, &c.SourceRadar, &c.Identified,
		&c.SquawkCode, &c.FlightPlanRef, &c.ThreatAssessment, &c.OperatorNotes,
		&c.FirstDetectedAt, &c.LastUpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get contact by id: %w", err)
	}
	return &c, nil
}

func (r *PostgresRepo) CreateIncident(i domain.AirDefenseIncident) error {
	query := `INSERT INTO airdef_incidents 
		(incident_id, severity, status, aircraft_id, interception_asset, pilot_response, 
		 engagement_rules_applied, duration_minutes, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.Exec(query,
		i.IncidentID, i.Severity, i.Status, i.AircraftID, i.InterceptionAsset,
		i.PilotResponse, i.EngagementRulesApplied, i.DurationMinutes, i.CreatedAt, i.UpdatedAt)
	return err
}

func (r *PostgresRepo) ResolveIncident(id uuid.UUID) error {
	query := `UPDATE airdef_incidents SET status = $1, updated_at = $2 WHERE incident_id = $3`
	_, err := r.db.Exec(query, domain.StatusClosed, time.Now(), id)
	return err
}

func (r *PostgresRepo) CreateNoFlyEntry(e domain.NoFlyListEntry) error {
	query := `INSERT INTO airdef_no_fly_list 
		(entry_id, identity_ref, full_name, document_number, reason, added_by, expires_at, interpol_notice_ref, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.db.Exec(query,
		e.EntryID, e.IdentityRef, e.FullName, e.DocumentNumber, e.Reason,
		e.AddedBy, e.ExpiresAt, e.InterpolNoticeRef, e.CreatedAt)
	return err
}

func (r *PostgresRepo) GetNoFlyEntry(identity string) (*domain.NoFlyListEntry, error) {
	query := `SELECT entry_id, identity_ref, full_name, document_number, reason, added_by, 
		expires_at, interpol_notice_ref, created_at 
		FROM airdef_no_fly_list WHERE identity_ref = $1 AND expires_at > $2`
	var e domain.NoFlyListEntry
	err := r.db.QueryRow(query, identity, time.Now()).Scan(
		&e.EntryID, &e.IdentityRef, &e.FullName, &e.DocumentNumber, &e.Reason,
		&e.AddedBy, &e.ExpiresAt, &e.InterpolNoticeRef, &e.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get no-fly entry: %w", err)
	}
	return &e, nil
}
