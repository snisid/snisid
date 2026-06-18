package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/opr-svc/internal/domain"
)

type ViolationRepo struct {
	pool *pgxpool.Pool
}

func NewViolationRepo(pool *pgxpool.Pool) *ViolationRepo {
	return &ViolationRepo{pool: pool}
}

func (r *ViolationRepo) Create(ctx context.Context, violation *domain.Violation) error {
	query := `
		INSERT INTO opr_violations 
			(violation_id, order_id, violation_date, violation_type, location_desc,
			 dept_code, reported_by, arrest_made, arrest_case_ref, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		violation.ViolationID, violation.OrderID, violation.ViolationDate,
		violation.ViolationType, violation.LocationDesc, violation.DeptCode,
		violation.ReportedBy, violation.ArrestMade, violation.ArrestCaseRef,
		violation.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create violation: %w", err)
	}
	return nil
}

func (r *ViolationRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Violation, error) {
	query := `
		SELECT violation_id, order_id, violation_date, violation_type, location_desc,
			   dept_code, reported_by, arrest_made, arrest_case_ref, created_at
		FROM opr_violations
		WHERE violation_id = $1
	`
	violation := &domain.Violation{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&violation.ViolationID, &violation.OrderID, &violation.ViolationDate,
		&violation.ViolationType, &violation.LocationDesc, &violation.DeptCode,
		&violation.ReportedBy, &violation.ArrestMade, &violation.ArrestCaseRef,
		&violation.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find violation: %w", err)
	}
	return violation, nil
}

func (r *ViolationRepo) FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.Violation, error) {
	query := `
		SELECT violation_id, order_id, violation_date, violation_type, location_desc,
			   dept_code, reported_by, arrest_made, arrest_case_ref, created_at
		FROM opr_violations
		WHERE order_id = $1
		ORDER BY violation_date DESC
	`
	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query violations: %w", err)
	}
	defer rows.Close()

	var violations []*domain.Violation
	for rows.Next() {
		violation := &domain.Violation{}
		err := rows.Scan(
			&violation.ViolationID, &violation.OrderID, &violation.ViolationDate,
			&violation.ViolationType, &violation.LocationDesc, &violation.DeptCode,
			&violation.ReportedBy, &violation.ArrestMade, &violation.ArrestCaseRef,
			&violation.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan violation: %w", err)
		}
		violations = append(violations, violation)
	}
	return violations, nil
}
