package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/dipe-svc/internal/domain"
)

type missingRepo struct {
	pool *pgxpool.Pool
}

func NewMissingRepo(pool *pgxpool.Pool) *missingRepo {
	return &missingRepo{pool: pool}
}

func (r *missingRepo) CreateCase(c *domain.MissingPerson) error {
	c.CaseID = uuid.New()
	_, err := r.pool.Exec(context.Background(),
		`INSERT INTO dipe_missing_persons
		 (case_id, national_dipe_id, case_type, status, snisid_person_id, full_name, aliases,
		  dob, gender, nationality, occupation, photo_refs, height_cm, weight_kg, skin_tone,
		  eye_color, hair_color, distinguishing_marks, clothing_last_seen, last_seen_date,
		  last_seen_location, last_seen_dept_code, last_seen_commune, last_seen_lat, last_seen_lng,
		  circumstances, sivc_alert_id, gang_id, extors_case_id, reported_by_name, reported_by_phone,
		  reported_by_snisid, report_date, reporting_unit, afis_subject_id, dna_sample_ref,
		  dna_profile_id, interpol_notice_ref, ncmec_ref, rvin_case_id)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
		         $21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40)`,
		c.CaseID, c.NationalDipeID, c.CaseType, c.Status, c.SnisidPersonID, c.FullName, c.Aliases,
		c.DOB, c.Gender, c.Nationality, c.Occupation, c.PhotoRefs, c.HeightCM, c.WeightKG,
		c.SkinTone, c.EyeColor, c.HairColor, c.DistinguishingMarks, c.ClothingLastSeen,
		c.LastSeenDate, c.LastSeenLocation, c.LastSeenDeptCode, c.LastSeenCommune,
		c.LastSeenLat, c.LastSeenLng, c.Circumstances, c.SivcAlertID, c.GangID,
		c.ExtorsCaseID, c.ReportedByName, c.ReportedByPhone, c.ReportedBySnisid,
		c.ReportDate, c.ReportingUnit, c.AfisSubjectID, c.DnaSampleRef, c.DnaProfileID,
		c.InterpolNoticeRef, c.NcmecRef, c.RvinCaseID,
	)
	return err
}

func (r *missingRepo) FindByID(id uuid.UUID) (*domain.MissingPerson, error) {
	row := r.pool.QueryRow(context.Background(),
		`SELECT case_id, national_dipe_id, case_type, status, snisid_person_id, full_name, aliases,
		        dob, gender, nationality, occupation, photo_refs, height_cm, weight_kg, skin_tone,
		        eye_color, hair_color, distinguishing_marks, clothing_last_seen, last_seen_date,
		        last_seen_location, last_seen_dept_code, last_seen_commune, last_seen_lat, last_seen_lng,
		        circumstances, sivc_alert_id, gang_id, extors_case_id, reported_by_name, reported_by_phone,
		        reported_by_snisid, report_date, reporting_unit, afis_subject_id, dna_sample_ref,
		        dna_profile_id, interpol_notice_ref, ncmec_ref, resolution_date, resolution_notes,
		        rvin_case_id, created_at, updated_at
		 FROM dipe_missing_persons WHERE case_id=$1`, id)
	var c domain.MissingPerson
	err := row.Scan(
		&c.CaseID, &c.NationalDipeID, &c.CaseType, &c.Status, &c.SnisidPersonID, &c.FullName,
		&c.Aliases, &c.DOB, &c.Gender, &c.Nationality, &c.Occupation, &c.PhotoRefs, &c.HeightCM,
		&c.WeightKG, &c.SkinTone, &c.EyeColor, &c.HairColor, &c.DistinguishingMarks,
		&c.ClothingLastSeen, &c.LastSeenDate, &c.LastSeenLocation, &c.LastSeenDeptCode,
		&c.LastSeenCommune, &c.LastSeenLat, &c.LastSeenLng, &c.Circumstances, &c.SivcAlertID,
		&c.GangID, &c.ExtorsCaseID, &c.ReportedByName, &c.ReportedByPhone, &c.ReportedBySnisid,
		&c.ReportDate, &c.ReportingUnit, &c.AfisSubjectID, &c.DnaSampleRef, &c.DnaProfileID,
		&c.InterpolNoticeRef, &c.NcmecRef, &c.ResolutionDate, &c.ResolutionNotes,
		&c.RvinCaseID, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("case not found: %w", err)
	}
	return &c, nil
}

func (r *missingRepo) GetOpenCases(limit, offset int) ([]*domain.MissingPerson, int, error) {
	var total int
	err := r.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM dipe_missing_persons WHERE status='OPEN'`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(context.Background(),
		`SELECT case_id, national_dipe_id, case_type, status, snisid_person_id, full_name, aliases,
		        dob, gender, nationality, occupation, photo_refs, height_cm, weight_kg, skin_tone,
		        eye_color, hair_color, distinguishing_marks, clothing_last_seen, last_seen_date,
		        last_seen_location, last_seen_dept_code, last_seen_commune, last_seen_lat, last_seen_lng,
		        circumstances, sivc_alert_id, gang_id, extors_case_id, reported_by_name, reported_by_phone,
		        reported_by_snisid, report_date, reporting_unit, afis_subject_id, dna_sample_ref,
		        dna_profile_id, interpol_notice_ref, ncmec_ref, resolution_date, resolution_notes,
		        rvin_case_id, created_at, updated_at
		 FROM dipe_missing_persons WHERE status='OPEN'
		 ORDER BY last_seen_date DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cases []*domain.MissingPerson
	for rows.Next() {
		var c domain.MissingPerson
		if err := rows.Scan(
			&c.CaseID, &c.NationalDipeID, &c.CaseType, &c.Status, &c.SnisidPersonID, &c.FullName,
			&c.Aliases, &c.DOB, &c.Gender, &c.Nationality, &c.Occupation, &c.PhotoRefs, &c.HeightCM,
			&c.WeightKG, &c.SkinTone, &c.EyeColor, &c.HairColor, &c.DistinguishingMarks,
			&c.ClothingLastSeen, &c.LastSeenDate, &c.LastSeenLocation, &c.LastSeenDeptCode,
			&c.LastSeenCommune, &c.LastSeenLat, &c.LastSeenLng, &c.Circumstances, &c.SivcAlertID,
			&c.GangID, &c.ExtorsCaseID, &c.ReportedByName, &c.ReportedByPhone, &c.ReportedBySnisid,
			&c.ReportDate, &c.ReportingUnit, &c.AfisSubjectID, &c.DnaSampleRef, &c.DnaProfileID,
			&c.InterpolNoticeRef, &c.NcmecRef, &c.ResolutionDate, &c.ResolutionNotes,
			&c.RvinCaseID, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		cases = append(cases, &c)
	}
	return cases, total, nil
}

func (r *missingRepo) AddSighting(s *domain.Sighting) error {
	s.SightingID = uuid.New()
	_, err := r.pool.Exec(context.Background(),
		`INSERT INTO dipe_sightings
		 (sighting_id, case_id, sighting_date, location_desc, dept_code, lat, lng,
		  reported_by, report_method, confidence, photo_ref)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		s.SightingID, s.CaseID, s.SightingDate, s.LocationDesc, s.DeptCode,
		s.Lat, s.Lng, s.ReportedBy, s.ReportMethod, s.Confidence, s.PhotoRef,
	)
	return err
}

func (r *missingRepo) GetSightingsByCase(caseID uuid.UUID) ([]*domain.Sighting, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT sighting_id, case_id, sighting_date, location_desc, dept_code, lat, lng,
		        reported_by, report_method, confidence, photo_ref, verified, verified_by, created_at
		 FROM dipe_sightings WHERE case_id=$1 ORDER BY sighting_date DESC`, caseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sightings []*domain.Sighting
	for rows.Next() {
		var s domain.Sighting
		if err := rows.Scan(
			&s.SightingID, &s.CaseID, &s.SightingDate, &s.LocationDesc, &s.DeptCode,
			&s.Lat, &s.Lng, &s.ReportedBy, &s.ReportMethod, &s.Confidence,
			&s.PhotoRef, &s.Verified, &s.VerifiedBy, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		sightings = append(sightings, &s)
	}
	return sightings, nil
}

func (r *missingRepo) ResolveCase(id uuid.UUID, status domain.CaseStatus, notes *string) error {
	_, err := r.pool.Exec(context.Background(),
		`UPDATE dipe_missing_persons SET status=$1, resolution_date=NOW(),
		        resolution_notes=$2, updated_at=NOW()
		 WHERE case_id=$3`,
		status, notes, id,
	)
	return err
}

func (r *missingRepo) GetStatsByType() (map[domain.CaseType]int, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT case_type, COUNT(*) FROM dipe_missing_persons GROUP BY case_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[domain.CaseType]int)
	for rows.Next() {
		var ct domain.CaseType
		var count int
		if err := rows.Scan(&ct, &count); err != nil {
			return nil, err
		}
		stats[ct] = count
	}
	return stats, nil
}
