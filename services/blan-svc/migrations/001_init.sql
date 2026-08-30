BEGIN;
CREATE TYPE blan_typology AS ENUM ('SMURFING','TRADE_BASED_ML','REAL_ESTATE','SHELL_COMPANY','CASH_INTENSIVE_BUSINESS','CRYPTO_MIXING','DIASPORA_TRANSFER','RANSOM_LAUNDERING','CORRUPTION_PROCEEDS');
CREATE TYPE blan_asset_type AS ENUM ('REAL_ESTATE','VEHICLE','BUSINESS','BANK_ACCOUNT','CRYPTO_WALLET','CASH','JEWELRY','LIVESTOCK','OTHER');
CREATE TABLE blan_cases (
    case_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), national_blan_id VARCHAR(25) UNIQUE NOT NULL,
    case_title TEXT NOT NULL, typology blan_typology NOT NULL,
    status VARCHAR(20) DEFAULT 'OPEN', total_amount_usd DECIMAL(18,2),
    predicate_crime VARCHAR(100), subject_ids UUID[] DEFAULT '{}', gang_id UUID,
    str_ids UUID[] DEFAULT '{}', opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    analyst_id UUID, parquet_ref VARCHAR(100), notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE blan_suspicious_assets (
    asset_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), case_id UUID NOT NULL REFERENCES blan_cases(case_id),
    asset_type blan_asset_type NOT NULL, description TEXT NOT NULL,
    address TEXT, dept_code CHAR(2), estimated_value_usd DECIMAL(15,2),
    acquisition_date DATE, owner_snisid_id UUID, owner_name VARCHAR(200),
    registered_in CHAR(3), is_frozen BOOLEAN DEFAULT FALSE,
    freeze_order_ref VARCHAR(100), is_seized BOOLEAN DEFAULT FALSE,
    seizure_date TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE blan_transaction_chains (
    chain_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), case_id UUID NOT NULL REFERENCES blan_cases(case_id),
    step_number SMALLINT NOT NULL, transaction_type VARCHAR(50),
    from_account VARCHAR(200), from_institution VARCHAR(150),
    to_account VARCHAR(200), to_institution VARCHAR(150),
    amount DECIMAL(18,2), currency CHAR(3), amount_usd DECIMAL(18,2),
    transaction_date TIMESTAMPTZ, is_suspicious_step BOOLEAN DEFAULT TRUE,
    notes TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE blan_real_estate_flagged (
    property_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), case_id UUID REFERENCES blan_cases(case_id),
    address TEXT NOT NULL, dept_code CHAR(2), commune VARCHAR(100),
    lat DECIMAL(10,7), lng DECIMAL(10,7), property_type VARCHAR(50),
    purchase_price_usd DECIMAL(15,2), purchase_date DATE,
    declared_owner VARCHAR(200), beneficial_owner_id UUID,
    suspicious_reasons TEXT[] DEFAULT '{}', is_frozen BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_blan_cases_status ON blan_cases(status, typology);
CREATE INDEX idx_blan_cases_gang ON blan_cases(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_blan_cases_persons ON blan_cases USING gin(subject_ids);
CREATE INDEX idx_blan_assets_type ON blan_suspicious_assets(asset_type, is_frozen);
CREATE INDEX idx_blan_real_estate ON blan_real_estate_flagged(dept_code);
COMMIT;
