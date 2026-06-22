package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/card-svc/internal/domain"
)

type Repository interface {
	CreateProfile(ctx context.Context, profile *domain.CardProfile) error
	FindProfileByID(ctx context.Context, profileID uuid.UUID) (*domain.CardProfile, error)
	CreatePersonalizationOrder(ctx context.Context, order *domain.PersonalizationRequest) error
	FindOrderBySerial(ctx context.Context, cardSerial string) (*domain.PersonalizationRequest, error)
	UpdateOrderStatus(ctx context.Context, cardSerial string, status domain.CardStatus) error
	UpdateOrderActivated(ctx context.Context, cardSerial string) error
	UpdateOrderBlocked(ctx context.Context, cardSerial string, reason string) error
	FindInventoryByProfileID(ctx context.Context, profileID uuid.UUID) (*domain.CardInventory, error)
	FindAllInventory(ctx context.Context) ([]domain.CardInventory, error)
	CreateStock(ctx context.Context, stock *domain.CardStock) error
	FindStockByProfileID(ctx context.Context, profileID uuid.UUID) ([]domain.CardStock, error)
	CreateShipment(ctx context.Context, shipment *domain.Shipment) error
	FindShipments(ctx context.Context) ([]domain.Shipment, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateProfile(ctx context.Context, profile *domain.CardProfile) error {
	query := `INSERT INTO card_profiles (profile_id, card_type, name, description, form_factor, material, has_chip, has_mrz, valid_days, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query,
		profile.ProfileID, profile.CardType, profile.Name, profile.Description,
		profile.FormFactor, profile.Material, profile.HasChip, profile.HasMRZ,
		profile.ValidDays, profile.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert card_profile: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindProfileByID(ctx context.Context, profileID uuid.UUID) (*domain.CardProfile, error) {
	query := `SELECT profile_id, card_type, name, description, form_factor, material, has_chip, has_mrz, valid_days, created_at
		FROM card_profiles WHERE profile_id = $1`
	p := &domain.CardProfile{}
	err := r.db.QueryRowContext(ctx, query, profileID).Scan(
		&p.ProfileID, &p.CardType, &p.Name, &p.Description, &p.FormFactor,
		&p.Material, &p.HasChip, &p.HasMRZ, &p.ValidDays, &p.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile not found: %s", profileID)
		}
		return nil, fmt.Errorf("query card_profile: %w", err)
	}
	return p, nil
}

func (r *postgresRepo) CreatePersonalizationOrder(ctx context.Context, order *domain.PersonalizationRequest) error {
	query := `INSERT INTO card_personalization (order_id, profile_id, card_serial, citizen_id, full_name, date_of_birth, nationality, photo_data, signature_data, status, ordered_at, personalized_at, issued_at, activated_at, blocked_at, block_reason, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`
	_, err := r.db.ExecContext(ctx, query,
		order.OrderID, order.ProfileID, order.CardSerial, order.CitizenID,
		order.FullName, order.DateOfBirth, order.Nationality, order.PhotoData,
		order.SignatureData, order.Status, order.OrderedAt, order.PersonalizedAt,
		order.IssuedAt, order.ActivatedAt, order.BlockedAt, order.BlockReason,
		order.ExpiresAt, order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert card_personalization: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindOrderBySerial(ctx context.Context, cardSerial string) (*domain.PersonalizationRequest, error) {
	query := `SELECT order_id, profile_id, card_serial, citizen_id, full_name, date_of_birth, nationality, photo_data, signature_data, status, ordered_at, personalized_at, issued_at, activated_at, blocked_at, block_reason, expires_at, created_at, updated_at
		FROM card_personalization WHERE card_serial = $1`
	o := &domain.PersonalizationRequest{}
	err := r.db.QueryRowContext(ctx, query, cardSerial).Scan(
		&o.OrderID, &o.ProfileID, &o.CardSerial, &o.CitizenID, &o.FullName,
		&o.DateOfBirth, &o.Nationality, &o.PhotoData, &o.SignatureData, &o.Status,
		&o.OrderedAt, &o.PersonalizedAt, &o.IssuedAt, &o.ActivatedAt, &o.BlockedAt,
		&o.BlockReason, &o.ExpiresAt, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("card not found: %s", cardSerial)
		}
		return nil, fmt.Errorf("query card_personalization: %w", err)
	}
	return o, nil
}

func (r *postgresRepo) UpdateOrderStatus(ctx context.Context, cardSerial string, status domain.CardStatus) error {
	query := `UPDATE card_personalization SET status = $1, updated_at = $2 WHERE card_serial = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now().UTC(), cardSerial)
	if err != nil {
		return fmt.Errorf("update card status: %w", err)
	}
	return nil
}

func (r *postgresRepo) UpdateOrderActivated(ctx context.Context, cardSerial string) error {
	now := time.Now().UTC()
	query := `UPDATE card_personalization SET status = 'ACTIVE', activated_at = $1, updated_at = $2 WHERE card_serial = $3`
	_, err := r.db.ExecContext(ctx, query, now, now, cardSerial)
	if err != nil {
		return fmt.Errorf("activate card: %w", err)
	}
	return nil
}

func (r *postgresRepo) UpdateOrderBlocked(ctx context.Context, cardSerial string, reason string) error {
	now := time.Now().UTC()
	query := `UPDATE card_personalization SET status = 'BLOCKED', blocked_at = $1, block_reason = $2, updated_at = $3 WHERE card_serial = $4`
	_, err := r.db.ExecContext(ctx, query, now, reason, now, cardSerial)
	if err != nil {
		return fmt.Errorf("block card: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindInventoryByProfileID(ctx context.Context, profileID uuid.UUID) (*domain.CardInventory, error) {
	query := `SELECT
		cs.profile_id, cp.name AS profile_name, cp.card_type,
		COALESCE(SUM(cs.quantity), 0) AS total_stock,
		COALESCE(SUM(cs.available_qty), 0) AS available,
		COALESCE((SELECT COUNT(*) FROM card_personalization cp2 WHERE cp2.profile_id = $1 AND cp2.status = 'PERSONALIZED'), 0) AS personalized,
		COALESCE((SELECT COUNT(*) FROM card_personalization cp2 WHERE cp2.profile_id = $1 AND cp2.status = 'ISSUED'), 0) AS issued,
		COALESCE((SELECT COUNT(*) FROM card_personalization cp2 WHERE cp2.profile_id = $1 AND cp2.status = 'BLOCKED'), 0) AS blocked,
		0 AS defective
		FROM card_stock cs
		JOIN card_profiles cp ON cp.profile_id = cs.profile_id
		WHERE cs.profile_id = $1
		GROUP BY cs.profile_id, cp.name, cp.card_type`
	inv := &domain.CardInventory{}
	err := r.db.QueryRowContext(ctx, query, profileID).Scan(
		&inv.ProfileID, &inv.ProfileName, &inv.CardType,
		&inv.TotalStock, &inv.Available, &inv.Personalized,
		&inv.Issued, &inv.Blocked, &inv.Defective,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("inventory not found for profile: %s", profileID)
		}
		return nil, fmt.Errorf("query card_inventory: %w", err)
	}
	return inv, nil
}

func (r *postgresRepo) FindAllInventory(ctx context.Context) ([]domain.CardInventory, error) {
	query := `SELECT
		cp.profile_id, cp.name AS profile_name, cp.card_type,
		COALESCE((SELECT SUM(quantity) FROM card_stock cs WHERE cs.profile_id = cp.profile_id), 0) AS total_stock,
		COALESCE((SELECT SUM(available_qty) FROM card_stock cs WHERE cs.profile_id = cp.profile_id), 0) AS available,
		COALESCE((SELECT COUNT(*) FROM card_personalization cp2 WHERE cp2.profile_id = cp.profile_id AND cp2.status = 'PERSONALIZED'), 0) AS personalized,
		COALESCE((SELECT COUNT(*) FROM card_personalization cp2 WHERE cp2.profile_id = cp.profile_id AND cp2.status = 'ISSUED'), 0) AS issued,
		COALESCE((SELECT COUNT(*) FROM card_personalization cp2 WHERE cp2.profile_id = cp.profile_id AND cp2.status = 'BLOCKED'), 0) AS blocked,
		0 AS defective
		FROM card_profiles cp
		ORDER BY cp.name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all inventory: %w", err)
	}
	defer rows.Close()

	var invs []domain.CardInventory
	for rows.Next() {
		var inv domain.CardInventory
		if err := rows.Scan(
			&inv.ProfileID, &inv.ProfileName, &inv.CardType,
			&inv.TotalStock, &inv.Available, &inv.Personalized,
			&inv.Issued, &inv.Blocked, &inv.Defective,
		); err != nil {
			return nil, err
		}
		invs = append(invs, inv)
	}
	return invs, rows.Err()
}

func (r *postgresRepo) CreateStock(ctx context.Context, stock *domain.CardStock) error {
	query := `INSERT INTO card_stock (stock_id, profile_id, serial_from, serial_to, quantity, available_qty, location, received_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		stock.StockID, stock.ProfileID, stock.SerialFrom, stock.SerialTo,
		stock.Quantity, stock.AvailableQty, stock.Location, stock.ReceivedAt,
	)
	if err != nil {
		return fmt.Errorf("insert card_stock: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindStockByProfileID(ctx context.Context, profileID uuid.UUID) ([]domain.CardStock, error) {
	query := `SELECT stock_id, profile_id, serial_from, serial_to, quantity, available_qty, location, received_at
		FROM card_stock WHERE profile_id = $1 ORDER BY received_at DESC`
	rows, err := r.db.QueryContext(ctx, query, profileID)
	if err != nil {
		return nil, fmt.Errorf("query card_stock: %w", err)
	}
	defer rows.Close()

	var stocks []domain.CardStock
	for rows.Next() {
		var s domain.CardStock
		if err := rows.Scan(
			&s.StockID, &s.ProfileID, &s.SerialFrom, &s.SerialTo,
			&s.Quantity, &s.AvailableQty, &s.Location, &s.ReceivedAt,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}
	return stocks, rows.Err()
}

func (r *postgresRepo) CreateShipment(ctx context.Context, shipment *domain.Shipment) error {
	query := `INSERT INTO card_shipments (shipment_id, profile_id, serial_from, serial_to, quantity, tracking_ref, vendor, received_by, received_at, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		shipment.ShipmentID, shipment.ProfileID, shipment.SerialFrom, shipment.SerialTo,
		shipment.Quantity, shipment.TrackingRef, shipment.Vendor, shipment.ReceivedBy,
		shipment.ReceivedAt, shipment.Notes, shipment.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert card_shipment: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindShipments(ctx context.Context) ([]domain.Shipment, error) {
	query := `SELECT shipment_id, profile_id, serial_from, serial_to, quantity, tracking_ref, vendor, received_by, received_at, notes, created_at
		FROM card_shipments ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query card_shipments: %w", err)
	}
	defer rows.Close()

	var shipments []domain.Shipment
	for rows.Next() {
		var s domain.Shipment
		if err := rows.Scan(
			&s.ShipmentID, &s.ProfileID, &s.SerialFrom, &s.SerialTo,
			&s.Quantity, &s.TrackingRef, &s.Vendor, &s.ReceivedBy,
			&s.ReceivedAt, &s.Notes, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		shipments = append(shipments, s)
	}
	return shipments, rows.Err()
}
