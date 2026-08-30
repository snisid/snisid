-- ============================================================
-- SNI-SIDE: National Watchlist Database
-- PostgreSQL 16
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_watchlist;
SET search_path TO snisid_watchlist;

-- ============ WATCHLIST ENTRIES ============
CREATE TABLE watchlist_entries (
    entry_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    watchlist_type VARCHAR(30) CHECK (watchlist_type IN (
        'PERSON','VEHICLE','ORGANIZATION','DOCUMENT','PHONE','EMAIL','WEAPON','VESSEL','AIRCRAFT'
    )),
    entry_category VARCHAR(50) CHECK (entry_category IN (
        'TERRORISM','NARCOTICS','MONEY_LAUNDERING','WAR_CRIMES','HUMAN_TRAFFICKING',
        'ORGANIZED_CRIME','CYBER_CRIME','FRAUD','SANCTIONS','AMBER_ALERT','FUGITIVE',
        'HIGH_RISK_TRAVELER','PROLIFERATION','CORRUPTION'
    )),
    value_primary VARCHAR(255) NOT NULL,
    value_secondary TEXT,
    full_name VARCHAR(255),
    alias TEXT[],
    niu VARCHAR(10),
    document_number VARCHAR(100),
    phone_number VARCHAR(50),
    email_address VARCHAR(255),
    vehicle_plate VARCHAR(50),
    vehicle_vin VARCHAR(17),
    organization_name VARCHAR(255),
    bank_account VARCHAR(100),
    crypto_wallet VARCHAR(255),
    passport_number VARCHAR(100),
    listing_authority VARCHAR(255) NOT NULL,
    listing_country VARCHAR(100),
    listing_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expiry_date TIMESTAMPTZ,
    reason TEXT NOT NULL,
    legal_basis TEXT,
    case_reference VARCHAR(100),
    risk_level VARCHAR(20) CHECK (risk_level IN ('CRITICAL','HIGH','MEDIUM','LOW')),
    confidence VARCHAR(20) CHECK (confidence IN ('CONFIRMED','HIGH_CONFIDENCE','MEDIUM_CONFIDENCE','LOW_CONFIDENCE','UNVERIFIED')),
    source_intelligence TEXT,
    related_entries JSONB DEFAULT '[]',
    tlp_level VARCHAR(10) CHECK (tlp_level IN ('RED','AMBER','GREEN','WHITE')),
    status VARCHAR(20) CHECK (status IN ('ACTIVE','EXPIRED','REMOVED','SUSPENDED')),
    reviewed_by VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_watchlist_type ON watchlist_entries(watchlist_type);
CREATE INDEX idx_watchlist_category ON watchlist_entries(entry_category);
CREATE INDEX idx_watchlist_primary ON watchlist_entries(value_primary);
CREATE INDEX idx_watchlist_name ON watchlist_entries(full_name);
CREATE INDEX idx_watchlist_niu ON watchlist_entries(niu);
CREATE INDEX idx_watchlist_document ON watchlist_entries(document_number);
CREATE INDEX idx_watchlist_phone ON watchlist_entries(phone_number);
CREATE INDEX idx_watchlist_email ON watchlist_entries(email_address);
CREATE INDEX idx_watchlist_plate ON watchlist_entries(vehicle_plate);
CREATE INDEX idx_watchlist_vin ON watchlist_entries(vehicle_vin);
CREATE INDEX idx_watchlist_org ON watchlist_entries(organization_name);
CREATE INDEX idx_watchlist_bank ON watchlist_entries(bank_account);
CREATE INDEX idx_watchlist_crypto ON watchlist_entries(crypto_wallet);
CREATE INDEX idx_watchlist_passport ON watchlist_entries(passport_number);
CREATE INDEX idx_watchlist_risk ON watchlist_entries(risk_level);
CREATE INDEX idx_watchlist_authority ON watchlist_entries(listing_authority);
CREATE INDEX idx_watchlist_status ON watchlist_entries(status);
CREATE INDEX idx_watchlist_tlp ON watchlist_entries(tlp_level);
CREATE INDEX idx_watchlist_created ON watchlist_entries(created_at DESC);

-- ============ WATCHLIST MATCH LOG ============
CREATE TABLE watchlist_matches (
    match_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL REFERENCES watchlist_entries(entry_id),
    match_source VARCHAR(100) NOT NULL,
    match_source_id VARCHAR(255),
    match_value VARCHAR(255) NOT NULL,
    match_type VARCHAR(30) NOT NULL,
    match_confidence DECIMAL(5,2),
    match_context TEXT,
    detected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    detected_by VARCHAR(100),
    detection_location VARCHAR(500),
    detection_method VARCHAR(100),
    alert_generated BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) CHECK (status IN ('NEW','REVIEWED','ACTIONED','FALSE_POSITIVE','ESCALATED')),
    action_taken TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_wm_entry ON watchlist_matches(entry_id);
CREATE INDEX idx_wm_source ON watchlist_matches(match_source);
CREATE INDEX idx_wm_value ON watchlist_matches(match_value);
CREATE INDEX idx_wm_detected ON watchlist_matches(detected_at DESC);
CREATE INDEX idx_wm_status ON watchlist_matches(status);

-- ============ WATCHLIST ALERTS ============
CREATE TABLE watchlist_alerts (
    alert_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES watchlist_matches(match_id),
    alert_level VARCHAR(20) CHECK (alert_level IN ('CRITICAL','HIGH','MEDIUM','LOW','INFO')),
    notification_channels TEXT[],
    notified_at TIMESTAMPTZ,
    acknowledged_by VARCHAR(100),
    acknowledged_at TIMESTAMPTZ,
    status VARCHAR(20) CHECK (status IN ('SENT','ACKNOWLEDGED','ESCALATED','RESOLVED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE watchlist_entries ENABLE ROW LEVEL SECURITY;
CREATE POLICY watchlist_national_select ON watchlist_entries FOR SELECT USING (
    current_setting('snisid.agency') IN ('PNH','DCPJ','IMMIGRATION','SOC','SNISID_ADMIN')
);
