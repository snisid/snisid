CREATE TYPE cybre_crime_type AS ENUM (
    'MONCASH_FRAUD', 'SIM_SWAPPING', 'PHISHING', 'IDENTITY_THEFT_DIGITAL',
    'SYSTEM_INTRUSION', 'RANSOMWARE', 'SOCIAL_MEDIA_MANIPULATION',
    'ONLINE_SCAM', 'DIGITAL_EXTORTION', 'CRYPTO_FRAUD',
    'CHILD_EXPLOITATION_ONLINE', 'STATE_SYSTEM_ATTACK', 'OTHER'
);

CREATE TYPE cybre_severity AS ENUM ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL');

CREATE TABLE cybre_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_cybre_id   VARCHAR(25) UNIQUE NOT NULL,
    crime_type          cybre_crime_type NOT NULL,
    severity            cybre_severity NOT NULL DEFAULT 'MEDIUM',
    status              VARCHAR(20) DEFAULT 'OPEN',
    victim_count        INTEGER DEFAULT 1,
    victim_snisid_ids   UUID[] DEFAULT '{}',
    victim_types        TEXT[] DEFAULT '{}',
    total_financial_loss_usd DECIMAL(15,2),
    incident_date       TIMESTAMPTZ NOT NULL,
    reported_date       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    attack_vector       TEXT,
    attack_method       TEXT,
    targeted_platform   VARCHAR(100),
    targeted_system     VARCHAR(100),
    suspect_ids         UUID[] DEFAULT '{}',
    suspect_phone       TEXT[] DEFAULT '{}',
    suspect_email       TEXT[] DEFAULT '{}',
    suspect_ip_hashes   TEXT[] DEFAULT '{}',
    crypto_wallet_ids   UUID[] DEFAULT '{}',
    suspect_countries   CHAR(3)[] DEFAULT '{}',
    digital_evidence_refs TEXT[] DEFAULT '{}',
    hash_evidence       TEXT[] DEFAULT '{}',
    chain_of_custody_ref VARCHAR(100),
    investigating_unit  VARCHAR(50) DEFAULT 'DCPJ_CYBER',
    conatel_ref         VARCHAR(50),
    case_reference      VARCHAR(100),
    parquet_ref         VARCHAR(100),
    ucref_str_id        UUID,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cybre_moncash_fraud_patterns (
    pattern_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id         UUID REFERENCES cybre_incidents(incident_id),
    fraud_type          VARCHAR(50) NOT NULL,
    moncash_phone       VARCHAR(20),
    amount_stolen_htg   DECIMAL(14,2),
    victims_count       INTEGER DEFAULT 1,
    modus_operandi      TEXT,
    detected_by         VARCHAR(50),
    linked_phone_numbers TEXT[] DEFAULT '{}',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cybre_intrusion_attempts (
    attempt_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id         UUID REFERENCES cybre_incidents(incident_id),
    target_system       VARCHAR(100) NOT NULL,
    attack_timestamp    TIMESTAMPTZ NOT NULL,
    attack_type         VARCHAR(50),
    source_ip_hash      VARCHAR(64),
    source_country      CHAR(3),
    was_successful      BOOLEAN DEFAULT FALSE,
    data_potentially_accessed TEXT,
    snisid_module_targeted VARCHAR(30),
    detection_source    VARCHAR(50),
    mitigated_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cybre_threat_intelligence (
    threat_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    indicator_type      VARCHAR(30) NOT NULL,
    indicator_value     VARCHAR(500) NOT NULL,
    threat_category     cybre_crime_type,
    confidence_score    SMALLINT CHECK (confidence_score BETWEEN 0 AND 100),
    source              VARCHAR(100),
    is_active           BOOLEAN DEFAULT TRUE,
    first_seen          TIMESTAMPTZ,
    last_seen           TIMESTAMPTZ,
    linked_incidents    UUID[] DEFAULT '{}',
    misp_ref            VARCHAR(100),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cybre_incidents_type ON cybre_incidents(crime_type, severity, status);
CREATE INDEX idx_cybre_incidents_date ON cybre_incidents(incident_date DESC);
CREATE INDEX idx_cybre_moncash_phone ON cybre_moncash_fraud_patterns(moncash_phone);
CREATE INDEX idx_cybre_intrusions_target ON cybre_intrusion_attempts(target_system, attack_timestamp DESC);
CREATE INDEX idx_cybre_intel_indicator ON cybre_threat_intelligence(indicator_type, indicator_value)
    WHERE is_active = TRUE;
