package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/blan-svc/internal/domain"
)

type caseRepo struct {
	pool *pgxpool.Pool
}

func NewCaseRepo(pool *pgxpool.Pool) *caseRepo {
	return &caseRepo{pool: pool}
}

func (r *caseRepo) CreateCase(c *domain.BLANCase) (*domain.BLANCase, error) {
	ctx := context.Background()
	c.CaseID = uuid.New()
	c.Status = "OPEN"
	c.OpenedAt = time.Now()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO blan_cases
		 (case_id, national_blan_id, case_title, typology, status, total_amount_usd,
		  predicate_crime, subject_ids, gang_id, str_ids, opened_at, analyst_id,
		  parquet_ref, notes, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		 RETURNING case_id, created_at, updated_at`,
		c.CaseID, c.NationalBlanID, c.CaseTitle, c.Typology, c.Status, c.TotalAmountUSD,
		c.PredicateCrime, c.SubjectIDs, c.GangID, c.StrIDs, c.OpenedAt, c.AnalystID,
		c.ParquetRef, c.Notes, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.CaseID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create case: %w", err)
	}

	return c, nil
}

func (r *caseRepo) FindByID(id uuid.UUID) (*domain.BLANCase, error) {
	ctx := context.Background()
	c := &domain.BLANCase{}
	err := r.pool.QueryRow(ctx,
		`SELECT case_id, national_blan_id, case_title, typology, status, total_amount_usd,
		        predicate_crime, subject_ids, gang_id, str_ids, opened_at, analyst_id,
		        parquet_ref, notes, created_at, updated_at
		 FROM blan_cases WHERE case_id = $1`, id).Scan(
		&c.CaseID, &c.NationalBlanID, &c.CaseTitle, &c.Typology, &c.Status, &c.TotalAmountUSD,
		&c.PredicateCrime, &c.SubjectIDs, &c.GangID, &c.StrIDs, &c.OpenedAt, &c.AnalystID,
		&c.ParquetRef, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("find case by id: %w", err)
	}
	return c, nil
}

func (r *caseRepo) AddAsset(a *domain.SuspiciousAsset) (*domain.SuspiciousAsset, error) {
	ctx := context.Background()
	a.AssetID = uuid.New()
	a.IsFrozen = false
	a.IsSeized = false
	a.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO blan_suspicious_assets
		 (asset_id, case_id, asset_type, description, address, dept_code, estimated_value_usd,
		  acquisition_date, owner_snisid_id, owner_name, registered_in, is_frozen,
		  freeze_order_ref, is_seized, seizure_date, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		 RETURNING asset_id, created_at`,
		a.AssetID, a.CaseID, a.AssetType, a.Description, a.Address, a.DeptCode, a.EstimatedValueUSD,
		a.AcquisitionDate, a.OwnerSnisidID, a.OwnerName, a.RegisteredIn, a.IsFrozen,
		a.FreezeOrderRef, a.IsSeized, a.SeizureDate, a.CreatedAt,
	).Scan(&a.AssetID, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("add asset: %w", err)
	}

	return a, nil
}

func (r *caseRepo) AddChainStep(step *domain.TransactionChain) (*domain.TransactionChain, error) {
	ctx := context.Background()
	step.ChainID = uuid.New()
	step.IsSuspiciousStep = true
	step.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO blan_transaction_chains
		 (chain_id, case_id, step_number, transaction_type, from_account, from_institution,
		  to_account, to_institution, amount, currency, amount_usd, transaction_date,
		  is_suspicious_step, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		 RETURNING chain_id, created_at`,
		step.ChainID, step.CaseID, step.StepNumber, step.TransactionType, step.FromAccount,
		step.FromInstitution, step.ToAccount, step.ToInstitution, step.Amount, step.Currency,
		step.AmountUSD, step.TransactionDate, step.IsSuspiciousStep, step.Notes, step.CreatedAt,
	).Scan(&step.ChainID, &step.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("add chain step: %w", err)
	}

	return step, nil
}

func (r *caseRepo) GetFlaggedRealEstate() ([]domain.RealEstateFlagged, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT property_id, case_id, address, dept_code, commune, lat, lng,
		        property_type, purchase_price_usd, purchase_date, declared_owner,
		        beneficial_owner_id, suspicious_reasons, is_frozen, created_at
		 FROM blan_real_estate_flagged
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("get flagged real estate: %w", err)
	}
	defer rows.Close()

	var estates []domain.RealEstateFlagged
	for rows.Next() {
		var e domain.RealEstateFlagged
		if err := rows.Scan(
			&e.PropertyID, &e.CaseID, &e.Address, &e.DeptCode, &e.Commune, &e.Lat, &e.Lng,
			&e.PropertyType, &e.PurchasePriceUSD, &e.PurchaseDate, &e.DeclaredOwner,
			&e.BeneficialOwnerID, &e.SuspiciousReasons, &e.IsFrozen, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan flagged real estate: %w", err)
		}
		estates = append(estates, e)
	}
	return estates, nil
}

func (r *caseRepo) GetFrozenAssets() ([]domain.SuspiciousAsset, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT asset_id, case_id, asset_type, description, address, dept_code,
		        estimated_value_usd, acquisition_date, owner_snisid_id, owner_name,
		        registered_in, is_frozen, freeze_order_ref, is_seized, seizure_date, created_at
		 FROM blan_suspicious_assets
		 WHERE is_frozen = TRUE
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("get frozen assets: %w", err)
	}
	defer rows.Close()

	var assets []domain.SuspiciousAsset
	for rows.Next() {
		var a domain.SuspiciousAsset
		if err := rows.Scan(
			&a.AssetID, &a.CaseID, &a.AssetType, &a.Description, &a.Address, &a.DeptCode,
			&a.EstimatedValueUSD, &a.AcquisitionDate, &a.OwnerSnisidID, &a.OwnerName,
			&a.RegisteredIn, &a.IsFrozen, &a.FreezeOrderRef, &a.IsSeized, &a.SeizureDate,
			&a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan frozen asset: %w", err)
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (r *caseRepo) GetStatsByTypology() ([]domain.TypologyStats, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT typology, COUNT(*) as case_count, COALESCE(SUM(total_amount_usd), 0) as total_usd
		 FROM blan_cases
		 GROUP BY typology
		 ORDER BY case_count DESC`)
	if err != nil {
		return nil, fmt.Errorf("get stats by typology: %w", err)
	}
	defer rows.Close()

	var stats []domain.TypologyStats
	for rows.Next() {
		var s domain.TypologyStats
		if err := rows.Scan(&s.Typology, &s.CaseCount, &s.TotalUSD); err != nil {
			return nil, fmt.Errorf("scan typology stats: %w", err)
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func (r *caseRepo) CountCasesByPrefix(prefix string) (int, error) {
	ctx := context.Background()
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM blan_cases WHERE national_blan_id LIKE $1`,
		fmt.Sprintf("%s%%", prefix)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count cases by prefix: %w", err)
	}
	return count, nil
}

func scanChains(rows pgx.Rows) ([]domain.TransactionChain, error) {
	var chains []domain.TransactionChain
	for rows.Next() {
		var c domain.TransactionChain
		if err := rows.Scan(
			&c.ChainID, &c.CaseID, &c.StepNumber, &c.TransactionType, &c.FromAccount,
			&c.FromInstitution, &c.ToAccount, &c.ToInstitution, &c.Amount, &c.Currency,
			&c.AmountUSD, &c.TransactionDate, &c.IsSuspiciousStep, &c.Notes, &c.CreatedAt); err != nil {
			return nil, err
		}
		chains = append(chains, c)
	}
	return chains, nil
}
