package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/opr-svc/internal/domain"
)

type ProtectionOrderRepo struct {
	pool *pgxpool.Pool
}

func NewProtectionOrderRepo(pool *pgxpool.Pool) *ProtectionOrderRepo {
	return &ProtectionOrderRepo{pool: pool}
}

func (r *ProtectionOrderRepo) Create(ctx context.Context, order *domain.ProtectionOrder) error {
	query := `
		INSERT INTO opr_protection_orders 
			(order_id, order_number, order_type, status, protected_person_id,
			 subject_person_id, subject_fir_id, exclusion_radius_m, exclusion_addresses,
			 no_contact_modes, geographic_ban_geojson, issuing_court, issuing_judge,
			 issue_date, expiry_date, is_renewable, violation_count, last_violation_at,
			 created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`
	_, err := r.pool.Exec(ctx, query,
		order.OrderID, order.OrderNumber, order.OrderType, order.Status,
		order.ProtectedPersonID, order.SubjectPersonID, order.SubjectFIRID,
		order.ExclusionRadiusM, order.ExclusionAddresses, order.NoContactModes,
		order.GeographicBanGeoJSON, order.IssuingCourt, order.IssuingJudge,
		order.IssueDate, order.ExpiryDate, order.IsRenewable, order.ViolationCount,
		order.LastViolationAt, order.CreatedBy, order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create protection order: %w", err)
	}
	return nil
}

func (r *ProtectionOrderRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProtectionOrder, error) {
	query := `
		SELECT order_id, order_number, order_type, status, protected_person_id,
			   subject_person_id, subject_fir_id, exclusion_radius_m, exclusion_addresses,
			   no_contact_modes, geographic_ban_geojson, issuing_court, issuing_judge,
			   issue_date, expiry_date, is_renewable, violation_count, last_violation_at,
			   created_by, created_at, updated_at
		FROM opr_protection_orders
		WHERE order_id = $1
	`
	order := &domain.ProtectionOrder{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.OrderID, &order.OrderNumber, &order.OrderType, &order.Status,
		&order.ProtectedPersonID, &order.SubjectPersonID, &order.SubjectFIRID,
		&order.ExclusionRadiusM, &order.ExclusionAddresses, &order.NoContactModes,
		&order.GeographicBanGeoJSON, &order.IssuingCourt, &order.IssuingJudge,
		&order.IssueDate, &order.ExpiryDate, &order.IsRenewable, &order.ViolationCount,
		&order.LastViolationAt, &order.CreatedBy, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find protection order: %w", err)
	}
	return order, nil
}

func (r *ProtectionOrderRepo) FindActiveBySubject(ctx context.Context, personID uuid.UUID) ([]*domain.ProtectionOrder, error) {
	query := `
		SELECT order_id, order_number, order_type, status, protected_person_id,
			   subject_person_id, subject_fir_id, exclusion_radius_m, exclusion_addresses,
			   no_contact_modes, geographic_ban_geojson, issuing_court, issuing_judge,
			   issue_date, expiry_date, is_renewable, violation_count, last_violation_at,
			   created_by, created_at, updated_at
		FROM opr_protection_orders
		WHERE subject_person_id = $1 AND status = 'ACTIVE'
		ORDER BY issue_date DESC
	`
	return r.queryOrders(ctx, query, personID)
}

func (r *ProtectionOrderRepo) FindExpiringSoon(ctx context.Context, days int) ([]*domain.ProtectionOrder, error) {
	query := `
		SELECT order_id, order_number, order_type, status, protected_person_id,
			   subject_person_id, subject_fir_id, exclusion_radius_m, exclusion_addresses,
			   no_contact_modes, geographic_ban_geojson, issuing_court, issuing_judge,
			   issue_date, expiry_date, is_renewable, violation_count, last_violation_at,
			   created_by, created_at, updated_at
		FROM opr_protection_orders
		WHERE status = 'ACTIVE' AND expiry_date <= NOW() + INTERVAL '1 day' * $1
		ORDER BY expiry_date ASC
	`
	rows, err := r.pool.Query(ctx, query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query expiring orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.ProtectionOrder
	for rows.Next() {
		order := &domain.ProtectionOrder{}
		err := rows.Scan(
			&order.OrderID, &order.OrderNumber, &order.OrderType, &order.Status,
			&order.ProtectedPersonID, &order.SubjectPersonID, &order.SubjectFIRID,
			&order.ExclusionRadiusM, &order.ExclusionAddresses, &order.NoContactModes,
			&order.GeographicBanGeoJSON, &order.IssuingCourt, &order.IssuingJudge,
			&order.IssueDate, &order.ExpiryDate, &order.IsRenewable, &order.ViolationCount,
			&order.LastViolationAt, &order.CreatedBy, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *ProtectionOrderRepo) FindByGangID(ctx context.Context, gangID uuid.UUID) ([]*domain.ProtectionOrder, error) {
	query := `
		SELECT order_id, order_number, order_type, status, protected_person_id,
			   subject_person_id, subject_fir_id, exclusion_radius_m, exclusion_addresses,
			   no_contact_modes, geographic_ban_geojson, issuing_court, issuing_judge,
			   issue_date, expiry_date, is_renewable, violation_count, last_violation_at,
			   created_by, created_at, updated_at
		FROM opr_protection_orders
		WHERE order_type = 'GANG_EXCLUSION_ZONE' AND status = 'ACTIVE'
		ORDER BY issue_date DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query gang orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.ProtectionOrder
	for rows.Next() {
		order := &domain.ProtectionOrder{}
		err := rows.Scan(
			&order.OrderID, &order.OrderNumber, &order.OrderType, &order.Status,
			&order.ProtectedPersonID, &order.SubjectPersonID, &order.SubjectFIRID,
			&order.ExclusionRadiusM, &order.ExclusionAddresses, &order.NoContactModes,
			&order.GeographicBanGeoJSON, &order.IssuingCourt, &order.IssuingJudge,
			&order.IssueDate, &order.ExpiryDate, &order.IsRenewable, &order.ViolationCount,
			&order.LastViolationAt, &order.CreatedBy, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *ProtectionOrderRepo) Update(ctx context.Context, order *domain.ProtectionOrder) error {
	query := `
		UPDATE opr_protection_orders
		SET status = $3, violation_count = $4, last_violation_at = $5, updated_at = $6
		WHERE order_id = $1 AND protected_person_id = $2
	`
	_, err := r.pool.Exec(ctx, query,
		order.OrderID, order.ProtectedPersonID, order.Status,
		order.ViolationCount, order.LastViolationAt, order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update protection order: %w", err)
	}
	return nil
}

func (r *ProtectionOrderRepo) queryOrders(ctx context.Context, query string, args ...interface{}) ([]*domain.ProtectionOrder, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.ProtectionOrder
	for rows.Next() {
		order := &domain.ProtectionOrder{}
		err := rows.Scan(
			&order.OrderID, &order.OrderNumber, &order.OrderType, &order.Status,
			&order.ProtectedPersonID, &order.SubjectPersonID, &order.SubjectFIRID,
			&order.ExclusionRadiusM, &order.ExclusionAddresses, &order.NoContactModes,
			&order.GeographicBanGeoJSON, &order.IssuingCourt, &order.IssuingJudge,
			&order.IssueDate, &order.ExpiryDate, &order.IsRenewable, &order.ViolationCount,
			&order.LastViolationAt, &order.CreatedBy, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}
