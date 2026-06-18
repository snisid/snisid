package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/blkl-svc/internal/domain"
)

type blacklistRepo struct {
	pool *pgxpool.Pool
}

func NewBlacklistRepo(pool *pgxpool.Pool) *blacklistRepo {
	return &blacklistRepo{pool: pool}
}

func (r *blacklistRepo) CheckPerson(personID uuid.UUID) (*domain.BlacklistCheckResult, error) {
	ctx := context.Background()
	result := &domain.BlacklistCheckResult{
		PersonID:     personID,
		Restrictions: []domain.RestrictionType{},
	}

	rows, err := r.pool.Query(ctx,
		`SELECT restriction_type, armed_dangerous
		 FROM blkl_blacklist
		 WHERE snisid_person_id = $1 AND is_active = true`,
		personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rt domain.RestrictionType
		var ad bool
		if err := rows.Scan(&rt, &ad); err != nil {
			return nil, err
		}
		result.IsBlacklisted = true
		result.Restrictions = append(result.Restrictions, rt)
		if ad {
			result.ArmedDangerous = true
		}
	}

	return result, nil
}

func (r *blacklistRepo) AddEntry(entry *domain.BlklBlacklist) (*domain.BlklBlacklist, error) {
	ctx := context.Background()
	entry.ID = uuid.New()
	entry.EntryID = "BLKL-" + entry.ID.String()[:8]
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	entry.IsActive = true

	err := r.pool.QueryRow(ctx,
		`INSERT INTO blkl_blacklist
		 (id, entry_id, national_blkl_id, snisid_person_id, restriction_type, source, source_record_id,
		  reason, court_order_ref, ordered_by, effective_date, expiry_date, is_permanent, is_active,
		  alert_level, armed_dangerous, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		 RETURNING id, created_at, updated_at`,
		entry.ID, entry.EntryID, entry.NationalBlklID, entry.SNISIDPersonID, entry.RestrictionType,
		entry.Source, entry.SourceRecordID, entry.Reason, entry.CourtOrderRef, entry.OrderedBy,
		entry.EffectiveDate, entry.ExpiryDate, entry.IsPermanent, entry.IsActive,
		entry.AlertLevel, entry.ArmedDangerous, entry.CreatedBy, entry.CreatedAt, entry.UpdatedAt,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *blacklistRepo) LiftEntry(id uuid.UUID, liftedBy string) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE blkl_blacklist SET is_active = false, updated_at = $1 WHERE id = $2`,
		time.Now(), id)
	return err
}

func (r *blacklistRepo) GetActiveEntries() ([]domain.BlklBlacklist, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT id, entry_id, national_blkl_id, snisid_person_id, restriction_type, source, source_record_id,
		        reason, court_order_ref, ordered_by, effective_date, expiry_date, is_permanent, is_active,
		        alert_level, armed_dangerous, created_by, created_at, updated_at
		 FROM blkl_blacklist
		 WHERE is_active = true
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEntries(rows)
}

func (r *blacklistRepo) GetExpiringSoon(days int) ([]domain.BlklBlacklist, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT id, entry_id, national_blkl_id, snisid_person_id, restriction_type, source, source_record_id,
		        reason, court_order_ref, ordered_by, effective_date, expiry_date, is_permanent, is_active,
		        alert_level, armed_dangerous, created_by, created_at, updated_at
		 FROM blkl_blacklist
		 WHERE is_active = true AND expiry_date IS NOT NULL AND expiry_date <= $1
		 ORDER BY expiry_date ASC`,
		time.Now().AddDate(0, 0, days))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEntries(rows)
}

func (r *blacklistRepo) GetByID(id uuid.UUID) (*domain.BlklBlacklist, error) {
	ctx := context.Background()
	entry := &domain.BlklBlacklist{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, entry_id, national_blkl_id, snisid_person_id, restriction_type, source, source_record_id,
		        reason, court_order_ref, ordered_by, effective_date, expiry_date, is_permanent, is_active,
		        alert_level, armed_dangerous, created_by, created_at, updated_at
		 FROM blkl_blacklist WHERE id = $1`, id).Scan(
		&entry.ID, &entry.EntryID, &entry.NationalBlklID, &entry.SNISIDPersonID, &entry.RestrictionType,
		&entry.Source, &entry.SourceRecordID, &entry.Reason, &entry.CourtOrderRef, &entry.OrderedBy,
		&entry.EffectiveDate, &entry.ExpiryDate, &entry.IsPermanent, &entry.IsActive,
		&entry.AlertLevel, &entry.ArmedDangerous, &entry.CreatedBy, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (r *blacklistRepo) LogAlert(log *domain.BlklAlertsLog) error {
	ctx := context.Background()
	log.ID = uuid.New()
	log.CreatedAt = time.Now()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO blkl_alerts_log (id, blacklist_id, person_id, alert_type, message, acknowledged, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		log.ID, log.BlacklistID, log.PersonID, log.AlertType, log.Message, log.Acknowledged, log.CreatedAt)
	return err
}

func scanEntries(rows pgx.Rows) ([]domain.BlklBlacklist, error) {
	var entries []domain.BlklBlacklist
	for rows.Next() {
		var e domain.BlklBlacklist
		if err := rows.Scan(
			&e.ID, &e.EntryID, &e.NationalBlklID, &e.SNISIDPersonID, &e.RestrictionType,
			&e.Source, &e.SourceRecordID, &e.Reason, &e.CourtOrderRef, &e.OrderedBy,
			&e.EffectiveDate, &e.ExpiryDate, &e.IsPermanent, &e.IsActive,
			&e.AlertLevel, &e.ArmedDangerous, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
