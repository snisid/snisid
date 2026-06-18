package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/sifr-svc/internal/domain"
)

type borderRepo struct {
	pool *pgxpool.Pool
}

func NewBorderRepo(pool *pgxpool.Pool) *borderRepo {
	return &borderRepo{pool: pool}
}

func (r *borderRepo) CreateCrossing(crossing *domain.Crossing) (*domain.Crossing, error) {
	ctx := context.Background()
	crossing.CrossingID = uuid.New()
	if crossing.CrossingDatetime.IsZero() {
		crossing.CrossingDatetime = time.Now().UTC()
	}
	crossing.CreatedAt = time.Now().UTC()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO sifr_crossings
		 (crossing_id, post_id, direction, crossing_datetime, snisid_person_id,
		  document_type, document_number, document_country, document_expiry,
		  traveler_name, traveler_dob, traveler_nationality, vehicle_plate,
		  lane_number, processing_officer, alert_triggered, alert_type,
		  alert_action_taken, processing_time_sec, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)
		 RETURNING crossing_id, created_at`,
		crossing.CrossingID, crossing.PostID, crossing.Direction, crossing.CrossingDatetime,
		crossing.SNISIDPersonID, crossing.DocumentType, crossing.DocumentNumber,
		crossing.DocumentCountry, crossing.DocumentExpiry, crossing.TravelerName,
		crossing.TravelerDob, crossing.TravelerNationality, crossing.VehiclePlate,
		crossing.LaneNumber, crossing.ProcessingOfficer, crossing.AlertTriggered,
		crossing.AlertType, crossing.AlertActionTaken, crossing.ProcessingTimeSec,
		crossing.CreatedAt,
	).Scan(&crossing.CrossingID, &crossing.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert crossing: %w", err)
	}
	return crossing, nil
}

func (r *borderRepo) FindCrossingsByPerson(personID uuid.UUID) ([]domain.Crossing, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT crossing_id, post_id, direction, crossing_datetime, snisid_person_id,
		        document_type, document_number, document_country, document_expiry,
		        traveler_name, traveler_dob, traveler_nationality, vehicle_plate,
		        lane_number, processing_officer, alert_triggered, alert_type,
		        alert_action_taken, processing_time_sec, created_at
		 FROM sifr_crossings WHERE snisid_person_id = $1 ORDER BY crossing_datetime DESC`, personID)
	if err != nil {
		return nil, fmt.Errorf("find crossings by person: %w", err)
	}
	defer rows.Close()
	return scanCrossings(rows)
}

func (r *borderRepo) FindCrossingsByPost(postID *uuid.UUID, limit, offset int) ([]domain.Crossing, int, error) {
	ctx := context.Background()

	var total int
	if postID != nil {
		err := r.pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM sifr_crossings WHERE post_id = $1`, *postID).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("count crossings by post: %w", err)
		}
	} else {
		err := r.pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM sifr_crossings`).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("count crossings: %w", err)
		}
	}

	query := `SELECT crossing_id, post_id, direction, crossing_datetime, snisid_person_id,
		        document_type, document_number, document_country, document_expiry,
		        traveler_name, traveler_dob, traveler_nationality, vehicle_plate,
		        lane_number, processing_officer, alert_triggered, alert_type,
		        alert_action_taken, processing_time_sec, created_at
		 FROM sifr_crossings`
	args := []interface{}{}
	argIdx := 1

	if postID != nil {
		query += fmt.Sprintf(" WHERE post_id = $%d", argIdx)
		args = append(args, *postID)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY crossing_datetime DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("find crossings by post: %w", err)
	}
	defer rows.Close()

	crossings, err := scanCrossings(rows)
	if err != nil {
		return nil, 0, err
	}
	return crossings, total, nil
}

func (r *borderRepo) FindActiveAlerts() ([]domain.AlertLog, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT alert_id, crossing_id, post_id, alert_type, snisid_person_id,
		        document_number, vehicle_plate, alert_source, source_record_id,
		        notified_units, action_taken, resolved, resolved_by, resolved_at, created_at
		 FROM sifr_alerts_log WHERE resolved = FALSE ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("find active alerts: %w", err)
	}
	defer rows.Close()

	var alerts []domain.AlertLog
	for rows.Next() {
		var a domain.AlertLog
		if err := rows.Scan(
			&a.AlertID, &a.CrossingID, &a.PostID, &a.AlertType, &a.SNISIDPersonID,
			&a.DocumentNumber, &a.VehiclePlate, &a.AlertSource, &a.SourceRecordID,
			&a.NotifiedUnits, &a.ActionTaken, &a.Resolved, &a.ResolvedBy,
			&a.ResolvedAt, &a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan alert: %w", err)
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r *borderRepo) GetBorderPosts() ([]domain.BorderPost, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT post_id, post_code, name, dept_code, border_country, post_lat, post_lng,
		        is_official, is_active, lanes_count, has_biometric_scanner, has_vehicle_scanner,
		        operating_hours, commanding_officer, created_at
		 FROM sifr_border_posts WHERE is_active = TRUE ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("get border posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.BorderPost
	for rows.Next() {
		var p domain.BorderPost
		if err := rows.Scan(
			&p.PostID, &p.PostCode, &p.Name, &p.DeptCode, &p.BorderCountry,
			&p.PostLat, &p.PostLng, &p.IsOfficial, &p.IsActive, &p.LanesCount,
			&p.HasBiometricScanner, &p.HasVehicleScanner, &p.OperatingHours,
			&p.CommandingOfficer, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan border post: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *borderRepo) CreateClandestineReport(report *domain.ClandestineCrossing) (*domain.ClandestineCrossing, error) {
	ctx := context.Background()
	report.ReportID = uuid.New()
	if report.ReportedDate.IsZero() {
		report.ReportedDate = time.Now().UTC()
	}
	report.CreatedAt = time.Now().UTC()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO sifr_clandestine_crossings
		 (report_id, location_desc, dept_code, lat, lng, reported_date,
		  crossing_type, estimated_persons, gang_related, gang_id,
		  trafficking_type, reported_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		 RETURNING report_id, created_at`,
		report.ReportID, report.LocationDesc, report.DeptCode, report.Lat, report.Lng,
		report.ReportedDate, report.CrossingType, report.EstimatedPersons,
		report.GangRelated, report.GangID, report.TraffickingType,
		report.ReportedBy, report.CreatedAt,
	).Scan(&report.ReportID, &report.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert clandestine report: %w", err)
	}
	return report, nil
}

func (r *borderRepo) GetDailyStats(postID *uuid.UUID) (map[string]interface{}, error) {
	ctx := context.Background()
	stats := make(map[string]interface{})

	query := `SELECT
		COUNT(*) AS total_crossings,
		COUNT(*) FILTER (WHERE direction = 'ENTRY') AS total_entries,
		COUNT(*) FILTER (WHERE direction = 'EXIT') AS total_exits,
		COUNT(*) FILTER (WHERE alert_triggered = TRUE) AS alerts_triggered,
		COALESCE(AVG(processing_time_sec), 0) AS avg_processing_time
	  FROM sifr_crossings
	  WHERE crossing_datetime >= CURRENT_DATE`
	args := []interface{}{}
	argIdx := 1

	if postID != nil {
		query += fmt.Sprintf(" AND post_id = $%d", argIdx)
		args = append(args, *postID)
		argIdx++
	}

	var total, entries, exits, alerts int
	var avgTime float64
	err := r.pool.QueryRow(ctx, query, args...).Scan(&total, &entries, &exits, &alerts, &avgTime)
	if err != nil {
		return nil, fmt.Errorf("get daily stats: %w", err)
	}

	stats["date"] = time.Now().Format("2006-01-02")
	stats["total_crossings"] = total
	stats["total_entries"] = entries
	stats["total_exits"] = exits
	stats["alerts_triggered"] = alerts
	stats["avg_processing_time_sec"] = avgTime

	return stats, nil
}

func scanCrossings(rows pgx.Rows) ([]domain.Crossing, error) {
	var crossings []domain.Crossing
	for rows.Next() {
		var c domain.Crossing
		if err := rows.Scan(
			&c.CrossingID, &c.PostID, &c.Direction, &c.CrossingDatetime,
			&c.SNISIDPersonID, &c.DocumentType, &c.DocumentNumber,
			&c.DocumentCountry, &c.DocumentExpiry, &c.TravelerName,
			&c.TravelerDob, &c.TravelerNationality, &c.VehiclePlate,
			&c.LaneNumber, &c.ProcessingOfficer, &c.AlertTriggered,
			&c.AlertType, &c.AlertActionTaken, &c.ProcessingTimeSec,
			&c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan crossing: %w", err)
		}
		crossings = append(crossings, c)
	}
	return crossings, nil
}
