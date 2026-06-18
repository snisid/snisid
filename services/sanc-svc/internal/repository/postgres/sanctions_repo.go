package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/sanc-svc/internal/domain"
)

type sanctionsRepo struct {
	pool *pgxpool.Pool
}

func NewSanctionsRepo(pool *pgxpool.Pool) *sanctionsRepo {
	return &sanctionsRepo{pool: pool}
}

func (r *sanctionsRepo) UpsertEntry(entry *domain.SanctionEntry) error {
	ctx := context.Background()
	if entry.SancID == uuid.Nil {
		entry.SancID = uuid.New()
	}
	entry.CreatedAt = time.Now().UTC()
	entry.UpdatedAt = entry.CreatedAt

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sanc_entries (
			sanc_id, source, source_ref_id, entity_type, entity_name, aliases, nationality,
			date_of_birth, place_of_birth, passport_numbers, national_id_numbers, measure_types,
			listing_date, end_date, is_active, listing_reason, committee_notes,
			snisid_person_id, gang_id, chef_member_id, match_confidence, source_updated_at,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24)
		ON CONFLICT (source, source_ref_id) DO UPDATE SET
			entity_type = EXCLUDED.entity_type,
			entity_name = EXCLUDED.entity_name,
			aliases = EXCLUDED.aliases,
			nationality = EXCLUDED.nationality,
			date_of_birth = EXCLUDED.date_of_birth,
			place_of_birth = EXCLUDED.place_of_birth,
			passport_numbers = EXCLUDED.passport_numbers,
			national_id_numbers = EXCLUDED.national_id_numbers,
			measure_types = EXCLUDED.measure_types,
			listing_date = EXCLUDED.listing_date,
			end_date = EXCLUDED.end_date,
			is_active = EXCLUDED.is_active,
			listing_reason = EXCLUDED.listing_reason,
			committee_notes = EXCLUDED.committee_notes,
			snisid_person_id = EXCLUDED.snisid_person_id,
			gang_id = EXCLUDED.gang_id,
			chef_member_id = EXCLUDED.chef_member_id,
			match_confidence = EXCLUDED.match_confidence,
			source_updated_at = EXCLUDED.source_updated_at,
			updated_at = NOW()`,
		entry.SancID, entry.Source, entry.SourceRefID, entry.EntityType, entry.EntityName,
		entry.Aliases, entry.Nationality, entry.DateOfBirth, entry.PlaceOfBirth,
		entry.PassportNumbers, entry.NationalIDNumbers, entry.MeasureTypes,
		entry.ListingDate, entry.EndDate, entry.IsActive, entry.ListingReason,
		entry.CommitteeNotes, entry.SNISIDPersonID, entry.GangID, entry.ChefMemberID,
		entry.MatchConfidence, entry.SourceUpdatedAt, entry.CreatedAt, entry.UpdatedAt,
	)
	return err
}

func (r *sanctionsRepo) SearchByNameAndDOB(name string, dob *time.Time) ([]domain.SanctionEntry, error) {
	ctx := context.Background()
	var rows pgx.Rows
	var err error

	if dob != nil {
		rows, err = r.pool.Query(ctx,
			`SELECT sanc_id, source, source_ref_id, entity_type, entity_name, aliases, nationality,
			        date_of_birth, place_of_birth, passport_numbers, national_id_numbers, measure_types,
			        listing_date, end_date, is_active, listing_reason, committee_notes,
			        snisid_person_id, gang_id, chef_member_id, match_confidence, source_updated_at,
			        created_at, updated_at
			 FROM sanc_entries
			 WHERE is_active = true
			   AND (to_tsvector('simple', entity_name) @@ plainto_tsquery('simple', $1)
			        OR $1 = ANY(aliases))
			   AND date_of_birth = $2
			 ORDER BY created_at DESC`, name, dob)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT sanc_id, source, source_ref_id, entity_type, entity_name, aliases, nationality,
			        date_of_birth, place_of_birth, passport_numbers, national_id_numbers, measure_types,
			        listing_date, end_date, is_active, listing_reason, committee_notes,
			        snisid_person_id, gang_id, chef_member_id, match_confidence, source_updated_at,
			        created_at, updated_at
			 FROM sanc_entries
			 WHERE is_active = true
			   AND (to_tsvector('simple', entity_name) @@ plainto_tsquery('simple', $1)
			        OR $1 = ANY(aliases))
			 ORDER BY created_at DESC`, name)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEntries(rows)
}

func (r *sanctionsRepo) GetActiveEntries(limit, offset int) ([]domain.SanctionEntry, int, error) {
	ctx := context.Background()
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM sanc_entries WHERE is_active = true`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT sanc_id, source, source_ref_id, entity_type, entity_name, aliases, nationality,
		        date_of_birth, place_of_birth, passport_numbers, national_id_numbers, measure_types,
		        listing_date, end_date, is_active, listing_reason, committee_notes,
		        snisid_person_id, gang_id, chef_member_id, match_confidence, source_updated_at,
		        created_at, updated_at
		 FROM sanc_entries WHERE is_active = true
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	entries, err := scanEntries(rows)
	return entries, total, err
}

func (r *sanctionsRepo) GetEntriesBySource(source domain.Source, limit, offset int) ([]domain.SanctionEntry, int, error) {
	ctx := context.Background()
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM sanc_entries WHERE source = $1 AND is_active = true`, source).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT sanc_id, source, source_ref_id, entity_type, entity_name, aliases, nationality,
		        date_of_birth, place_of_birth, passport_numbers, national_id_numbers, measure_types,
		        listing_date, end_date, is_active, listing_reason, committee_notes,
		        snisid_person_id, gang_id, chef_member_id, match_confidence, source_updated_at,
		        created_at, updated_at
		 FROM sanc_entries WHERE source = $1 AND is_active = true
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`, source, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	entries, err := scanEntries(rows)
	return entries, total, err
}

func (r *sanctionsRepo) GetUnconfirmedMatches() ([]domain.IdentityMatch, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT match_id, sanc_id, snisid_person_id, match_score, match_fields,
		        confirmed_by, is_confirmed, is_false_positive, reviewed_at, created_at
		 FROM sanc_identity_matches
		 WHERE is_confirmed = false AND is_false_positive = false
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMatches(rows)
}

func (r *sanctionsRepo) SaveMatch(match *domain.IdentityMatch) error {
	ctx := context.Background()
	if match.MatchID == uuid.Nil {
		match.MatchID = uuid.New()
	}
	match.CreatedAt = time.Now().UTC()

	return r.pool.QueryRow(ctx,
		`INSERT INTO sanc_identity_matches (
			match_id, sanc_id, snisid_person_id, match_score, match_fields,
			confirmed_by, is_confirmed, is_false_positive, reviewed_at, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING match_id`,
		match.MatchID, match.SancID, match.SNISIDPersonID, match.MatchScore,
		match.MatchFields, match.ConfirmedBy, match.IsConfirmed,
		match.IsFalsePositive, match.ReviewedAt, match.CreatedAt,
	).Scan(&match.MatchID)
}

func (r *sanctionsRepo) ConfirmMatch(matchID uuid.UUID, confirmedBy uuid.UUID) error {
	ctx := context.Background()
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE sanc_identity_matches
		 SET is_confirmed = true, confirmed_by = $1, reviewed_at = $2
		 WHERE match_id = $3`, confirmedBy, now, matchID)
	return err
}

func (r *sanctionsRepo) SaveSyncLog(log *domain.SyncLog) error {
	ctx := context.Background()
	if log.SyncID == uuid.Nil {
		log.SyncID = uuid.New()
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sanc_sync_log (
			sync_id, source, started_at, status
		) VALUES ($1,$2,$3,$4)`,
		log.SyncID, log.Source, log.StartedAt, log.Status,
	)
	return err
}

func (r *sanctionsRepo) UpdateSyncLog(log *domain.SyncLog) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE sanc_sync_log SET
			completed_at = $1, entries_processed = $2, entries_added = $3,
			entries_updated = $4, entries_removed = $5, errors = $6,
			status = $7, error_details = $8
		WHERE sync_id = $9`,
		log.CompletedAt, log.EntriesProcessed, log.EntriesAdded,
		log.EntriesUpdated, log.EntriesRemoved, log.Errors,
		log.Status, log.ErrorDetails, log.SyncID,
	)
	return err
}

func (r *sanctionsRepo) GetSyncStatus(limit int) ([]domain.SyncLog, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT sync_id, source, started_at, completed_at, entries_processed,
		        entries_added, entries_updated, entries_removed, errors, status, error_details
		 FROM sanc_sync_log
		 ORDER BY started_at DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.SyncLog
	for rows.Next() {
		var l domain.SyncLog
		if err := rows.Scan(
			&l.SyncID, &l.Source, &l.StartedAt, &l.CompletedAt, &l.EntriesProcessed,
			&l.EntriesAdded, &l.EntriesUpdated, &l.EntriesRemoved, &l.Errors,
			&l.Status, &l.ErrorDetails,
		); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func scanEntries(rows pgx.Rows) ([]domain.SanctionEntry, error) {
	var entries []domain.SanctionEntry
	for rows.Next() {
		var e domain.SanctionEntry
		if err := rows.Scan(
			&e.SancID, &e.Source, &e.SourceRefID, &e.EntityType, &e.EntityName,
			&e.Aliases, &e.Nationality, &e.DateOfBirth, &e.PlaceOfBirth,
			&e.PassportNumbers, &e.NationalIDNumbers, &e.MeasureTypes,
			&e.ListingDate, &e.EndDate, &e.IsActive, &e.ListingReason,
			&e.CommitteeNotes, &e.SNISIDPersonID, &e.GangID, &e.ChefMemberID,
			&e.MatchConfidence, &e.SourceUpdatedAt, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func scanMatches(rows pgx.Rows) ([]domain.IdentityMatch, error) {
	var matches []domain.IdentityMatch
	for rows.Next() {
		var m domain.IdentityMatch
		if err := rows.Scan(
			&m.MatchID, &m.SancID, &m.SNISIDPersonID, &m.MatchScore,
			&m.MatchFields, &m.ConfirmedBy, &m.IsConfirmed,
			&m.IsFalsePositive, &m.ReviewedAt, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, nil
}
