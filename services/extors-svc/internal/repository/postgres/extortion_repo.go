package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/extors-svc/internal/domain"
)

type extortionRepo struct {
	pool *pgxpool.Pool
}

func NewExtortionRepo(pool *pgxpool.Pool) *extortionRepo {
	return &extortionRepo{pool: pool}
}

func (r *extortionRepo) CreateCase(c *domain.ExtorsCase) (*domain.ExtorsCase, error) {
	ctx := context.Background()
	c.CaseID = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	if c.Status == "" {
		c.Status = domain.ACTIVE
	}

	err := r.pool.QueryRow(ctx,
		`INSERT INTO extors_cases
		 (case_id, national_extors_id, extors_type, status, gang_id, gang_name,
		  perpetrator_ids, chef_member_ids, victim_count, victim_snisid_ids, victim_types,
		  victim_nationality, is_foreigner_victim, incident_location, dept_code, commune,
		  lat, lng, route_number, demanded_amount, demanded_currency, paid_amount, paid_currency,
		  payment_channel, payment_ref, payment_date, first_contact_date, resolution_date,
		  case_reference, investigating_unit, ucref_str_id, blan_case_id, notes, created_by,
		  created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36)
		 RETURNING created_at, updated_at`,
		c.CaseID, c.NationalExtorsID, c.ExtorsType, c.Status, c.GangID, c.GangName,
		c.PerpetratorIDs, c.ChefMemberIDs, c.VictimCount, c.VictimSNISIDs, c.VictimTypes,
		c.VictimNationality, c.IsForeignerVictim, c.IncidentLocation, c.DeptCode, c.Commune,
		c.Lat, c.Lng, c.RouteNumber, c.DemandedAmount, c.DemandedCurrency, c.PaidAmount, c.PaidCurrency,
		c.PaymentChannel, c.PaymentRef, c.PaymentDate, c.FirstContactDate, c.ResolutionDate,
		c.CaseReference, c.InvestigatingUnit, c.UcrefStrID, c.BlanCaseID, c.Notes, c.CreatedBy,
		c.CreatedAt, c.UpdatedAt,
	).Scan(&c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create case: %w", err)
	}

	return c, nil
}

func (r *extortionRepo) FindByID(id uuid.UUID) (*domain.ExtorsCase, error) {
	ctx := context.Background()
	c := &domain.ExtorsCase{}
	err := r.pool.QueryRow(ctx,
		`SELECT case_id, national_extors_id, extors_type, status, gang_id, gang_name,
		        perpetrator_ids, chef_member_ids, victim_count, victim_snisid_ids, victim_types,
		        victim_nationality, is_foreigner_victim, incident_location, dept_code, commune,
		        lat, lng, route_number, demanded_amount, demanded_currency, paid_amount, paid_currency,
		        payment_channel, payment_ref, payment_date, first_contact_date, resolution_date,
		        case_reference, investigating_unit, ucref_str_id, blan_case_id, notes, created_by,
		        created_at, updated_at
		 FROM extors_cases WHERE case_id = $1`, id).Scan(
		&c.CaseID, &c.NationalExtorsID, &c.ExtorsType, &c.Status, &c.GangID, &c.GangName,
		&c.PerpetratorIDs, &c.ChefMemberIDs, &c.VictimCount, &c.VictimSNISIDs, &c.VictimTypes,
		&c.VictimNationality, &c.IsForeignerVictim, &c.IncidentLocation, &c.DeptCode, &c.Commune,
		&c.Lat, &c.Lng, &c.RouteNumber, &c.DemandedAmount, &c.DemandedCurrency, &c.PaidAmount, &c.PaidCurrency,
		&c.PaymentChannel, &c.PaymentRef, &c.PaymentDate, &c.FirstContactDate, &c.ResolutionDate,
		&c.CaseReference, &c.InvestigatingUnit, &c.UcrefStrID, &c.BlanCaseID, &c.Notes, &c.CreatedBy,
		&c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("find case by id: %w", err)
	}

	return c, nil
}

func (r *extortionRepo) AddNegotiation(n *domain.Negotiation) (*domain.Negotiation, error) {
	ctx := context.Background()
	n.NegID = uuid.New()
	n.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO extors_negotiations
		 (neg_id, case_id, negotiation_date, contact_method, contact_number,
		  demand_updated, demand_currency, position_update, recorded_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		 RETURNING created_at`,
		n.NegID, n.CaseID, n.NegotiationDate, n.ContactMethod, n.ContactNumber,
		n.DemandUpdated, n.DemandCurrency, n.PositionUpdate, n.RecordedBy, n.CreatedAt,
	).Scan(&n.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("add negotiation: %w", err)
	}

	return n, nil
}

func (r *extortionRepo) CreateTollPoint(t *domain.RoadTollPoint) (*domain.RoadTollPoint, error) {
	ctx := context.Background()
	t.TollID = uuid.New()
	t.IsActive = true
	t.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO extors_road_toll_points
		 (toll_id, gang_id, location_desc, route_number, dept_code, commune,
		  lat, lng, daily_revenue_usd, vehicle_types_taxed, toll_rates,
		  active_since, is_active, source_intel, last_confirmed_at, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		 RETURNING created_at`,
		t.TollID, t.GangID, t.LocationDesc, t.RouteNumber, t.DeptCode, t.Commune,
		t.Lat, t.Lng, t.DailyRevenueUSD, t.VehicleTypesTaxed, t.TollRates,
		t.ActiveSince, t.IsActive, t.SourceIntel, t.LastConfirmedAt, t.CreatedBy, t.CreatedAt,
	).Scan(&t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create toll point: %w", err)
	}

	return t, nil
}

func (r *extortionRepo) FindActiveTollsByGang(gangID uuid.UUID) ([]domain.RoadTollPoint, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT toll_id, gang_id, location_desc, route_number, dept_code, commune,
		        lat, lng, daily_revenue_usd, vehicle_types_taxed, toll_rates,
		        active_since, is_active, source_intel, last_confirmed_at, created_by, created_at
		 FROM extors_road_toll_points
		 WHERE gang_id = $1 AND is_active = TRUE`, gangID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTollPoints(rows)
}

func (r *extortionRepo) FindPaidRansomsByGang(gangID uuid.UUID) ([]domain.ExtorsCase, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT case_id, national_extors_id, extors_type, status, gang_id, gang_name,
		        perpetrator_ids, chef_member_ids, victim_count, victim_snisid_ids, victim_types,
		        victim_nationality, is_foreigner_victim, incident_location, dept_code, commune,
		        lat, lng, route_number, demanded_amount, demanded_currency, paid_amount, paid_currency,
		        payment_channel, payment_ref, payment_date, first_contact_date, resolution_date,
		        case_reference, investigating_unit, ucref_str_id, blan_case_id, notes, created_by,
		        created_at, updated_at
		 FROM extors_cases
		 WHERE gang_id = $1 AND paid_amount IS NOT NULL`, gangID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCases(rows)
}

func (r *extortionRepo) FindActiveRacketsByGang(gangID uuid.UUID) ([]domain.ExtorsCase, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT case_id, national_extors_id, extors_type, status, gang_id, gang_name,
		        perpetrator_ids, chef_member_ids, victim_count, victim_snisid_ids, victim_types,
		        victim_nationality, is_foreigner_victim, incident_location, dept_code, commune,
		        lat, lng, route_number, demanded_amount, demanded_currency, paid_amount, paid_currency,
		        payment_channel, payment_ref, payment_date, first_contact_date, resolution_date,
		        case_reference, investigating_unit, ucref_str_id, blan_case_id, notes, created_by,
		        created_at, updated_at
		 FROM extors_cases
		 WHERE gang_id = $1 AND extors_type IN ('BUSINESS_PROTECTION_RACKET','REAL_ESTATE_EXTORTION','PUBLIC_SERVANT_EXTORTION','NGO_EXTORTION')
		   AND status = 'ACTIVE'`, gangID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCases(rows)
}

func (r *extortionRepo) GetTollsMap() ([]domain.RoadTollPoint, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT toll_id, gang_id, location_desc, route_number, dept_code, commune,
		        lat, lng, daily_revenue_usd, vehicle_types_taxed, toll_rates,
		        active_since, is_active, source_intel, last_confirmed_at, created_by, created_at
		 FROM extors_road_toll_points
		 WHERE is_active = TRUE AND lat IS NOT NULL AND lng IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTollPoints(rows)
}

func (r *extortionRepo) GetStatsByType() ([]domain.TypeStats, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT extors_type, COUNT(*) as count,
		        COALESCE(SUM(paid_amount), 0) as paid_total
		 FROM extors_cases
		 GROUP BY extors_type
		 ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []domain.TypeStats
	for rows.Next() {
		var s domain.TypeStats
		if err := rows.Scan(&s.ExtorsType, &s.Count, &s.PaidTotal); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func scanTollPoints(rows pgx.Rows) ([]domain.RoadTollPoint, error) {
	var points []domain.RoadTollPoint
	for rows.Next() {
		var t domain.RoadTollPoint
		if err := rows.Scan(
			&t.TollID, &t.GangID, &t.LocationDesc, &t.RouteNumber, &t.DeptCode, &t.Commune,
			&t.Lat, &t.Lng, &t.DailyRevenueUSD, &t.VehicleTypesTaxed, &t.TollRates,
			&t.ActiveSince, &t.IsActive, &t.SourceIntel, &t.LastConfirmedAt, &t.CreatedBy, &t.CreatedAt); err != nil {
			return nil, err
		}
		points = append(points, t)
	}
	return points, nil
}

func scanCases(rows pgx.Rows) ([]domain.ExtorsCase, error) {
	var cases []domain.ExtorsCase
	for rows.Next() {
		var c domain.ExtorsCase
		if err := rows.Scan(
			&c.CaseID, &c.NationalExtorsID, &c.ExtorsType, &c.Status, &c.GangID, &c.GangName,
			&c.PerpetratorIDs, &c.ChefMemberIDs, &c.VictimCount, &c.VictimSNISIDs, &c.VictimTypes,
			&c.VictimNationality, &c.IsForeignerVictim, &c.IncidentLocation, &c.DeptCode, &c.Commune,
			&c.Lat, &c.Lng, &c.RouteNumber, &c.DemandedAmount, &c.DemandedCurrency, &c.PaidAmount, &c.PaidCurrency,
			&c.PaymentChannel, &c.PaymentRef, &c.PaymentDate, &c.FirstContactDate, &c.ResolutionDate,
			&c.CaseReference, &c.InvestigatingUnit, &c.UcrefStrID, &c.BlanCaseID, &c.Notes, &c.CreatedBy,
			&c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, nil
}
