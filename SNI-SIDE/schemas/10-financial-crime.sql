-- ============================================================
-- SNI-SIDE: Financial Crime Database
-- PostgreSQL 16
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_financial;
SET search_path TO snisid_financial;

-- ============ SUSPICIOUS TRANSACTIONS ============
CREATE TABLE suspicious_transactions (
    transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_ref VARCHAR(100) UNIQUE NOT NULL,
    transaction_date TIMESTAMPTZ NOT NULL,
    transaction_type VARCHAR(50) CHECK (transaction_type IN (
        'WIRE_TRANSFER','CASH_DEPOSIT','CASH_WITHDRAWAL','CHECK','CRYPTO_TRANSFER',
        'TRADE','LOAN','INSURANCE','REAL_ESTATE','TRUST','SHELL_COMPANY'
    )),
    amount DECIMAL(20,2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    amount_usd_equivalent DECIMAL(20,2),
    source_account VARCHAR(100),
    source_bank VARCHAR(255),
    source_country VARCHAR(100),
    destination_account VARCHAR(100),
    destination_bank VARCHAR(255),
    destination_country VARCHAR(100),
    sender_name VARCHAR(255),
    sender_niu VARCHAR(10),
    sender_id_type VARCHAR(50),
    sender_id_number VARCHAR(100),
    sender_address TEXT,
    sender_phone VARCHAR(50),
    beneficiary_name VARCHAR(255),
    beneficiary_niu VARCHAR(10),
    beneficiary_id_type VARCHAR(50),
    beneficiary_id_number VARCHAR(100),
    beneficiary_address TEXT,
    beneficiary_phone VARCHAR(50),
    relationship VARCHAR(100),
    purpose TEXT,
    source_of_funds VARCHAR(255),
    risk_score DECIMAL(5,2),
    mlro_filed BOOLEAN DEFAULT FALSE,
    mlro_ref VARCHAR(100),
    status VARCHAR(20) CHECK (status IN ('PENDING','REVIEWING','ESCALATED','CLOSED','FILED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (transaction_date);

CREATE INDEX idx_txn_ref ON suspicious_transactions(transaction_ref);
CREATE INDEX idx_txn_date ON suspicious_transactions(transaction_date DESC);
CREATE INDEX idx_txn_amount ON suspicious_transactions(amount DESC);
CREATE INDEX idx_txn_sender ON suspicious_transactions(sender_niu);
CREATE INDEX idx_txn_beneficiary ON suspicious_transactions(beneficiary_niu);
CREATE INDEX idx_txn_risk ON suspicious_transactions(risk_score DESC);
CREATE INDEX idx_txn_status ON suspicious_transactions(status);
CREATE INDEX idx_txn_type ON suspicious_transactions(transaction_type);
CREATE INDEX idx_txn_source_country ON suspicious_transactions(source_country);
CREATE INDEX idx_txn_dest_country ON suspicious_transactions(destination_country);

-- ============ AML ALERTS ============
CREATE TABLE aml_alerts (
    alert_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES suspicious_transactions(transaction_id),
    alert_type VARCHAR(50) CHECK (alert_type IN (
        'STRUCTURING','HIGH_RISK_JURISDICTION','UNUSUAL_PATTERN','LARGE_CASH',
        'PEP_RELATED','SANCTIONS_MATCH','FALSE_POSITIVE','RAPID_MOVEMENT',
        'TRADE_BASED_ML','SHELL_COMPANY','CRYPTO_ANOMALY'
    )),
    alert_score DECIMAL(5,2),
    alert_description TEXT,
    detected_by VARCHAR(100),
    detection_rule VARCHAR(100),
    assigned_to VARCHAR(100),
    status VARCHAR(20) CHECK (status IN ('NEW','INVESTIGATING','CONFIRMED','FALSE_POSITIVE','ESCALATED_FIU')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_aml_alert_txn ON aml_alerts(transaction_id);
CREATE INDEX idx_aml_alert_type ON aml_alerts(alert_type);
CREATE INDEX idx_aml_alert_score ON aml_alerts(alert_score DESC);
CREATE INDEX idx_aml_alert_status ON aml_alerts(status);

-- ============ BENEFICIAL OWNERS ============
CREATE TABLE beneficial_owners (
    owner_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    niu VARCHAR(10) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    entity_type VARCHAR(30) CHECK (entity_type IN (
        'COMPANY','TRUST','FOUNDATION','PARTNERSHIP','SHELL_COMPANY','OFFSHORE'
    )),
    entity_name VARCHAR(255) NOT NULL,
    entity_registration_number VARCHAR(100),
    entity_country VARCHAR(100),
    ownership_percentage DECIMAL(5,2),
    role_in_entity VARCHAR(100),
    politically_exposed BOOLEAN DEFAULT FALSE,
    pep_position VARCHAR(255),
    pep_country VARCHAR(100),
    source_of_wealth TEXT,
    risk_score DECIMAL(5,2),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_bo_niu ON beneficial_owners(niu);
CREATE INDEX idx_bo_entity ON beneficial_owners(entity_name);
CREATE INDEX idx_bo_pep ON beneficial_owners(politically_exposed);
CREATE INDEX idx_bo_risk ON beneficial_owners(risk_score DESC);

-- ============ POLITICALLY EXPOSED PERSONS ============
CREATE TABLE politically_exposed_persons (
    pep_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    niu VARCHAR(10) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    position VARCHAR(255) NOT NULL,
    institution VARCHAR(255),
    country VARCHAR(100),
    start_date DATE,
    end_date DATE,
    pep_level VARCHAR(20) CHECK (pep_level IN ('HEAD_OF_STATE','GOVERNMENT','JUDICIARY','MILITARY','LEGISLATIVE','DIPLOMATIC','STATE_ENTERPRISE')),
    family_members JSONB DEFAULT '[]',
    close_associates JSONB DEFAULT '[]',
    risk_score DECIMAL(5,2),
    status VARCHAR(20) CHECK (status IN ('ACTIVE','FORMER','DECEASED')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_pep_niu UNIQUE (niu)
);

CREATE INDEX idx_pep_niu ON politically_exposed_persons(niu);
CREATE INDEX idx_pep_position ON politically_exposed_persons(position);
CREATE INDEX idx_pep_risk ON politically_exposed_persons(risk_score DESC);

-- ============ CORRUPTION CASES ============
CREATE TABLE corruption_cases (
    case_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_number VARCHAR(50) UNIQUE NOT NULL,
    case_type VARCHAR(50) CHECK (case_type IN (
        'BRIBERY','EMBEZZLEMENT','KICKBACK','FRAUD','ABUSE_OF_OFFICE','MONEY_LAUNDERING','ILLICIT_ENRICHMENT'
    )),
    allegation_date DATE,
    investigation_started DATE,
    investigating_agency VARCHAR(100),
    prosecutor VARCHAR(255),
    accused_persons TEXT[],
    accused_niu_array TEXT[],
    estimated_amount DECIMAL(20,2),
    currency VARCHAR(10),
    Sector VARCHAR(100),
    modus_operandi TEXT,
    status VARCHAR(20) CHECK (status IN ('PENDING','UNDER_INVESTIGATION','PROSECUTED','CONVICTED','ACQUITTED','DISMISSED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============ FINANCIAL NETWORK ANALYSIS ============
CREATE TABLE financial_networks (
    network_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    network_name VARCHAR(255),
    network_type VARCHAR(50) CHECK (network_type IN (
        'ML_SCHEME','CORRUPTION_RING','TAX_EVASION','TRADE_BASED_ML','CRYPTO_ML','HAWALA'
    )),
    total_value DECIMAL(20,2),
    total_transactions INT,
    time_period_start DATE,
    time_period_end DATE,
    entities JSONB DEFAULT '[]',
    transactions JSONB DEFAULT '[]',
    risk_score DECIMAL(5,2),
    status VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE suspicious_transactions ENABLE ROW LEVEL SECURITY;
CREATE POLICY financial_fiu_select ON suspicious_transactions FOR SELECT USING (
    current_setting('snisid.agency') IN ('FIU','PNH','DCPJ','ANTI_CORRUPTION','SNISID_ADMIN')
);
