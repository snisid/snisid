package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/crypt-svc/internal/domain"
)

type walletRepo struct {
	pool *pgxpool.Pool
}

func NewWalletRepo(pool *pgxpool.Pool) *walletRepo {
	return &walletRepo{pool: pool}
}

func (r *walletRepo) CreateWallet(wallet *domain.FlaggedWallet) (*domain.FlaggedWallet, error) {
	ctx := context.Background()
	wallet.ID = uuid.New()
	wallet.NationalCryptID = "CRYPT-HT-" + wallet.ID.String()[:8]
	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO crypt_flagged_wallets
		 (wallet_id, national_crypt_id, wallet_address, asset_type, blockchain_network,
		  suspicion_type, snisid_person_id, gang_id, estimated_balance_usd,
		  total_received_usd, total_sent_usd, first_tx_date, last_tx_date,
		  is_sanctioned, ofac_sdn_ref, chainalysis_ref, elliptic_ref, source_intel,
		  linked_cases, is_frozen, freeze_jurisdiction, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24)
		 RETURNING wallet_id, created_at, updated_at`,
		wallet.ID, wallet.NationalCryptID, wallet.WalletAddress, wallet.AssetType,
		wallet.BlockchainNetwork, wallet.SuspicionType, wallet.SNISIDPersonID,
		wallet.GangID, wallet.EstimatedBalanceUSD, wallet.TotalReceivedUSD,
		wallet.TotalSentUSD, wallet.FirstTxDate, wallet.LastTxDate,
		wallet.IsSanctioned, wallet.OfacSDNRef, wallet.ChainalysisRef,
		wallet.EllipticRef, wallet.SourceIntel, wallet.LinkedCases,
		wallet.IsFrozen, wallet.FreezeJurisdiction, wallet.CreatedBy,
		wallet.CreatedAt, wallet.UpdatedAt,
	).Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r *walletRepo) FindByAddress(address string) (*domain.FlaggedWallet, error) {
	ctx := context.Background()
	wallet := &domain.FlaggedWallet{}
	err := r.pool.QueryRow(ctx,
		`SELECT wallet_id, national_crypt_id, wallet_address, asset_type, blockchain_network,
		        suspicion_type, snisid_person_id, gang_id, estimated_balance_usd,
		        total_received_usd, total_sent_usd, first_tx_date, last_tx_date,
		        is_sanctioned, ofac_sdn_ref, chainalysis_ref, elliptic_ref, source_intel,
		        linked_cases, is_frozen, freeze_jurisdiction, created_by, created_at, updated_at
		 FROM crypt_flagged_wallets WHERE wallet_address = $1`, address).Scan(
		&wallet.ID, &wallet.NationalCryptID, &wallet.WalletAddress, &wallet.AssetType,
		&wallet.BlockchainNetwork, &wallet.SuspicionType, &wallet.SNISIDPersonID,
		&wallet.GangID, &wallet.EstimatedBalanceUSD, &wallet.TotalReceivedUSD,
		&wallet.TotalSentUSD, &wallet.FirstTxDate, &wallet.LastTxDate,
		&wallet.IsSanctioned, &wallet.OfacSDNRef, &wallet.ChainalysisRef,
		&wallet.EllipticRef, &wallet.SourceIntel, &wallet.LinkedCases,
		&wallet.IsFrozen, &wallet.FreezeJurisdiction, &wallet.CreatedBy,
		&wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *walletRepo) FindByID(id uuid.UUID) (*domain.FlaggedWallet, error) {
	ctx := context.Background()
	wallet := &domain.FlaggedWallet{}
	err := r.pool.QueryRow(ctx,
		`SELECT wallet_id, national_crypt_id, wallet_address, asset_type, blockchain_network,
		        suspicion_type, snisid_person_id, gang_id, estimated_balance_usd,
		        total_received_usd, total_sent_usd, first_tx_date, last_tx_date,
		        is_sanctioned, ofac_sdn_ref, chainalysis_ref, elliptic_ref, source_intel,
		        linked_cases, is_frozen, freeze_jurisdiction, created_by, created_at, updated_at
		 FROM crypt_flagged_wallets WHERE wallet_id = $1`, id).Scan(
		&wallet.ID, &wallet.NationalCryptID, &wallet.WalletAddress, &wallet.AssetType,
		&wallet.BlockchainNetwork, &wallet.SuspicionType, &wallet.SNISIDPersonID,
		&wallet.GangID, &wallet.EstimatedBalanceUSD, &wallet.TotalReceivedUSD,
		&wallet.TotalSentUSD, &wallet.FirstTxDate, &wallet.LastTxDate,
		&wallet.IsSanctioned, &wallet.OfacSDNRef, &wallet.ChainalysisRef,
		&wallet.EllipticRef, &wallet.SourceIntel, &wallet.LinkedCases,
		&wallet.IsFrozen, &wallet.FreezeJurisdiction, &wallet.CreatedBy,
		&wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *walletRepo) GetSanctionedWallets() ([]domain.FlaggedWallet, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT wallet_id, national_crypt_id, wallet_address, asset_type, blockchain_network,
		        suspicion_type, snisid_person_id, gang_id, estimated_balance_usd,
		        total_received_usd, total_sent_usd, first_tx_date, last_tx_date,
		        is_sanctioned, ofac_sdn_ref, chainalysis_ref, elliptic_ref, source_intel,
		        linked_cases, is_frozen, freeze_jurisdiction, created_by, created_at, updated_at
		 FROM crypt_flagged_wallets WHERE is_sanctioned = true
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanWallets(rows)
}

func (r *walletRepo) GetWalletsByGang(gangID uuid.UUID) ([]domain.FlaggedWallet, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT wallet_id, national_crypt_id, wallet_address, asset_type, blockchain_network,
		        suspicion_type, snisid_person_id, gang_id, estimated_balance_usd,
		        total_received_usd, total_sent_usd, first_tx_date, last_tx_date,
		        is_sanctioned, ofac_sdn_ref, chainalysis_ref, elliptic_ref, source_intel,
		        linked_cases, is_frozen, freeze_jurisdiction, created_by, created_at, updated_at
		 FROM crypt_flagged_wallets WHERE gang_id = $1
		 ORDER BY created_at DESC`, gangID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanWallets(rows)
}

func (r *walletRepo) AddTransaction(tx *domain.CryptoTransaction) (*domain.CryptoTransaction, error) {
	ctx := context.Background()
	tx.ID = uuid.New()
	tx.CreatedAt = time.Now()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO crypt_transactions
		 (tx_id, wallet_id, tx_hash, asset_type, direction, from_address, to_address,
		  amount_crypto, amount_usd_at_tx, tx_timestamp, block_number, is_mixer_involved,
		  mixer_service, risk_score, suspicion_flags, extors_case_id, ucref_str_id, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
		 RETURNING tx_id, created_at`,
		tx.ID, tx.WalletID, tx.TxHash, tx.AssetType, tx.Direction,
		tx.FromAddress, tx.ToAddress, tx.AmountCrypto, tx.AmountUSDAtTx,
		tx.TxTimestamp, tx.BlockNumber, tx.IsMixerInvolved, tx.MixerService,
		tx.RiskScore, tx.SuspicionFlags, tx.ExtorsCaseID, tx.UCRefStrID,
		tx.CreatedAt,
	).Scan(&tx.ID, &tx.CreatedAt)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *walletRepo) GetTransactionsByWallet(walletID uuid.UUID) ([]domain.CryptoTransaction, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT tx_id, wallet_id, tx_hash, asset_type, direction, from_address, to_address,
		        amount_crypto, amount_usd_at_tx, tx_timestamp, block_number, is_mixer_involved,
		        mixer_service, risk_score, suspicion_flags, extors_case_id, ucref_str_id, created_at
		 FROM crypt_transactions WHERE wallet_id = $1
		 ORDER BY tx_timestamp DESC`, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTransactions(rows)
}

func (r *walletRepo) GetStatsByAsset() ([]domain.AssetStats, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT asset_type,
		        COUNT(*) AS total_wallets,
		        COUNT(*) FILTER (WHERE is_sanctioned = true) AS sanctioned_count,
		        COUNT(*) FILTER (WHERE is_frozen = true) AS frozen_count,
		        COALESCE(SUM(estimated_balance_usd), 0) AS total_balance_usd,
		        0.0 AS avg_risk_score
		 FROM crypt_flagged_wallets
		 GROUP BY asset_type
		 ORDER BY asset_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []domain.AssetStats
	for rows.Next() {
		var s domain.AssetStats
		if err := rows.Scan(
			&s.AssetType, &s.TotalWallets, &s.SanctionedCount,
			&s.FrozenCount, &s.TotalBalanceUSD, &s.AvgRiskScore); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func scanWallets(rows pgx.Rows) ([]domain.FlaggedWallet, error) {
	var wallets []domain.FlaggedWallet
	for rows.Next() {
		var w domain.FlaggedWallet
		if err := rows.Scan(
			&w.ID, &w.NationalCryptID, &w.WalletAddress, &w.AssetType,
			&w.BlockchainNetwork, &w.SuspicionType, &w.SNISIDPersonID,
			&w.GangID, &w.EstimatedBalanceUSD, &w.TotalReceivedUSD,
			&w.TotalSentUSD, &w.FirstTxDate, &w.LastTxDate,
			&w.IsSanctioned, &w.OfacSDNRef, &w.ChainalysisRef,
			&w.EllipticRef, &w.SourceIntel, &w.LinkedCases,
			&w.IsFrozen, &w.FreezeJurisdiction, &w.CreatedBy,
			&w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}

func scanTransactions(rows pgx.Rows) ([]domain.CryptoTransaction, error) {
	var txs []domain.CryptoTransaction
	for rows.Next() {
		var t domain.CryptoTransaction
		if err := rows.Scan(
			&t.ID, &t.WalletID, &t.TxHash, &t.AssetType, &t.Direction,
			&t.FromAddress, &t.ToAddress, &t.AmountCrypto, &t.AmountUSDAtTx,
			&t.TxTimestamp, &t.BlockNumber, &t.IsMixerInvolved, &t.MixerService,
			&t.RiskScore, &t.SuspicionFlags, &t.ExtorsCaseID, &t.UCRefStrID,
			&t.CreatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}
