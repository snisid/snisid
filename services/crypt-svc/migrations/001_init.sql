BEGIN;
CREATE TYPE crypt_asset_type AS ENUM ('BITCOIN','ETHEREUM','USDT','USDC','MONERO','ZCASH','LITECOIN','OTHER_ERC20','UNKNOWN');
CREATE TYPE crypt_suspicion_type AS ENUM ('RANSOM_RECEIPT','SANCTIONS_EVASION','DARKWEB_PAYMENT','MIXER_SERVICE','PEER_TO_PEER_UNREGULATED','EXCHANGE_HIGH_RISK','GANG_PAYMENT','UNKNOWN');
CREATE TABLE crypt_flagged_wallets (
    wallet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), national_crypt_id VARCHAR(25) UNIQUE NOT NULL,
    wallet_address VARCHAR(200) NOT NULL, asset_type crypt_asset_type NOT NULL,
    blockchain_network VARCHAR(50), suspicion_type crypt_suspicion_type NOT NULL,
    snisid_person_id UUID, gang_id UUID, estimated_balance_usd DECIMAL(18,2),
    total_received_usd DECIMAL(18,2), total_sent_usd DECIMAL(18,2),
    first_tx_date TIMESTAMPTZ, last_tx_date TIMESTAMPTZ,
    is_sanctioned BOOLEAN DEFAULT FALSE, ofac_sdn_ref VARCHAR(50),
    chainalysis_ref VARCHAR(100), elliptic_ref VARCHAR(100), source_intel TEXT,
    linked_cases UUID[] DEFAULT '{}', is_frozen BOOLEAN DEFAULT FALSE,
    freeze_jurisdiction VARCHAR(50), created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE crypt_transactions (
    tx_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), wallet_id UUID REFERENCES crypt_flagged_wallets(wallet_id),
    tx_hash VARCHAR(100) NOT NULL, asset_type crypt_asset_type NOT NULL,
    direction VARCHAR(10) NOT NULL, from_address VARCHAR(200),
    to_address VARCHAR(200), amount_crypto DECIMAL(30,18),
    amount_usd_at_tx DECIMAL(18,2), tx_timestamp TIMESTAMPTZ NOT NULL,
    block_number BIGINT, is_mixer_involved BOOLEAN DEFAULT FALSE,
    mixer_service VARCHAR(100), risk_score SMALLINT,
    suspicion_flags TEXT[] DEFAULT '{}', extors_case_id UUID,
    ucref_str_id UUID, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE crypt_exchange_accounts (
    exchange_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), snisid_person_id UUID,
    exchange_name VARCHAR(100) NOT NULL, exchange_country CHAR(3),
    account_ref VARCHAR(200), kyc_level VARCHAR(20),
    total_volume_usd DECIMAL(18,2), is_flagged BOOLEAN DEFAULT FALSE,
    flagging_reason TEXT, legal_hold_request BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_crypt_wallets_address ON crypt_flagged_wallets(wallet_address);
CREATE INDEX idx_crypt_wallets_gang ON crypt_flagged_wallets(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_crypt_wallets_sanctioned ON crypt_flagged_wallets(is_sanctioned) WHERE is_sanctioned = TRUE;
CREATE INDEX idx_crypt_tx_wallet ON crypt_transactions(wallet_id, tx_timestamp DESC);
CREATE INDEX idx_crypt_tx_hash ON crypt_transactions(tx_hash);
CREATE INDEX idx_crypt_tx_mixer ON crypt_transactions(is_mixer_involved) WHERE is_mixer_involved = TRUE;
COMMIT;
