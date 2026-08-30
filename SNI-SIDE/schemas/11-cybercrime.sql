-- ============================================================
-- SNI-SIDE: Cybercrime Intelligence Database
-- PostgreSQL 16 + TimescaleDB
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_cybercrime;
SET search_path TO snisid_cybercrime;

-- ============ INDICATORS OF COMPROMISE (IOC) ============
CREATE TABLE iocs (
    ioc_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ioc_value TEXT NOT NULL,
    ioc_type VARCHAR(30) CHECK (ioc_type IN (
        'IPV4','IPV6','DOMAIN','URL','MD5','SHA1','SHA256','EMAIL','USERNAME',
        'REGISTRY_KEY','FILE_PATH','MUTEX','SERVICE_NAME','CVE','YARA_RULE'
    )),
    ioc_hash VARCHAR(64) GENERATED ALWAYS AS (encode(sha256(ioc_value::bytea), 'hex')) STORED,
    confidence INT CHECK (confidence >= 0 AND confidence <= 100),
    severity VARCHAR(20) CHECK (severity IN ('CRITICAL','HIGH','MEDIUM','LOW','INFO')),
    threat_actor_id UUID REFERENCES threat_actors(actor_id),
    malware_family VARCHAR(255),
    first_seen TIMESTAMPTZ DEFAULT NOW(),
    last_seen TIMESTAMPTZ DEFAULT NOW(),
    source VARCHAR(255),
    tags TEXT[],
    description TEXT,
    tlp_level VARCHAR(10) CHECK (tlp_level IN ('RED','AMBER','GREEN','WHITE')),
    status VARCHAR(20) CHECK (status IN ('ACTIVE','OBSOLETE','WHITELISTED','REVIEWING')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_ioc_hash UNIQUE (ioc_hash)
);

CREATE INDEX idx_ioc_value ON iocs(ioc_value);
CREATE INDEX idx_ioc_type ON iocs(ioc_type);
CREATE INDEX idx_ioc_hash_idx ON iocs(ioc_hash);
CREATE INDEX idx_ioc_severity ON iocs(severity);
CREATE INDEX idx_ioc_actor ON iocs(threat_actor_id);
CREATE INDEX idx_ioc_family ON iocs(malware_family);
CREATE INDEX idx_ioc_last_seen ON iocs(last_seen DESC);
CREATE INDEX idx_ioc_tlp ON iocs(tlp_level);
CREATE INDEX idx_ioc_tags ON iocs USING gin(tags);

-- ============ THREAT ACTORS ============
CREATE TABLE threat_actors (
    actor_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    alias TEXT[],
    actor_type VARCHAR(50) CHECK (actor_type IN (
        'APT','CYBER_CRIMINAL','HACKTIVIST','INSIDER_THREAT','STATE_SPONSORED',
        'CYBER_TERRORIST','SCRIPT_KIDDIE','ORGANIZED_CRIME'
    )),
    country_of_origin VARCHAR(100),
    motivation VARCHAR(255),
    target_sectors TEXT[],
    known_tools TEXT[],
    known_techniques TEXT[],
    first_observed DATE,
    last_observed DATE,
    aliases_attributed TEXT[],
    associated_campaigns TEXT[],
    confidence_level VARCHAR(20),
    status VARCHAR(20) CHECK (status IN ('ACTIVE','INACTIVE','ATTRIBUTED','UNKNOWN')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_actor_name ON threat_actors(name);
CREATE INDEX idx_actor_type ON threat_actors(actor_type);
CREATE INDEX idx_actor_country ON threat_actors(country_of_origin);

-- ============ MALWARE SAMPLES ============
CREATE TABLE malware_samples (
    sample_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_name VARCHAR(255),
    file_size_bytes BIGINT,
    file_type VARCHAR(100),
    md5_hash VARCHAR(32),
    sha1_hash VARCHAR(40),
    sha256_hash VARCHAR(64) UNIQUE NOT NULL,
    imphash VARCHAR(64),
    ssdeep_hash VARCHAR(255),
    authentihash VARCHAR(64),
    compilation_timestamp TIMESTAMPTZ,
    architecture VARCHAR(50),
    packer VARCHAR(100),
    malware_family VARCHAR(255),
    malware_type VARCHAR(100) CHECK (malware_type IN (
        'RANSOMWARE','TROJAN','WORM','VIRUS','BOTNET','ROOTKIT','BACKDOOR',
        'KEYLOGGER','SPYWARE','LOADER','DROPPER','INFO_STEALER','BANKER','RAT','OTHER'
    )),
    mitre_attack_ids TEXT[],
    yara_matches TEXT[],
    vt_detection_ratio VARCHAR(20),
    vt_link VARCHAR(500),
    first_submitted TIMESTAMPTZ,
    analysis_report_path VARCHAR(500),
    tlp_level VARCHAR(10),
    status VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_malware_sha256 ON malware_samples(sha256_hash);
CREATE INDEX idx_malware_md5 ON malware_samples(md5_hash);
CREATE INDEX idx_malware_family ON malware_samples(malware_family);
CREATE INDEX idx_malware_type ON malware_samples(malware_type);

-- ============ CRYPTO WALLETS ============
CREATE TABLE crypto_wallets (
    wallet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_address VARCHAR(255) UNIQUE NOT NULL,
    blockchain VARCHAR(50) CHECK (blockchain IN (
        'BITCOIN','ETHEREUM','USDT_TRC20','USDT_ERC20','SOLANA','MONERO','LITECOIN','BSC','POLYGON','AVALANCHE','OTHER'
    )),
    wallet_type VARCHAR(30) CHECK (wallet_type IN ('EXCHANGE','MIXER','DARKNET_MARKET','PERSONAL','SCAM','RANSOMWARE')),
    first_seen TIMESTAMPTZ,
    last_active TIMESTAMPTZ,
    total_received_usd DECIMAL(20,2),
    total_sent_usd DECIMAL(20,2),
    transaction_count INT,
    known_entity VARCHAR(255),
    risk_score DECIMAL(5,2),
    tags TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_wallet_address ON crypto_wallets(wallet_address);
CREATE INDEX idx_wallet_blockchain ON crypto_wallets(blockchain);
CREATE INDEX idx_wallet_type ON crypto_wallets(wallet_type);
CREATE INDEX idx_wallet_risk ON crypto_wallets(risk_score DESC);

-- ============ CYBER CAMPAIGNS ============
CREATE TABLE cyber_campaigns (
    campaign_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    campaign_type VARCHAR(50),
    threat_actor_id UUID REFERENCES threat_actors(actor_id),
    start_date DATE,
    end_date DATE,
    targeted_entities TEXT[],
    targeted_regions TEXT[],
    targeted_sectors TEXT[],
    iocs_used TEXT[],
    techniques_observed TEXT[],
    malware_families TEXT[],
    impact_assessment TEXT,
    status VARCHAR(20) CHECK (status IN ('ONGOING','PAST','ATTRIBUTED','UNKNOWN')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============ CYBER INCIDENTS ============
CREATE TABLE cyber_incidents (
    incident_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_number VARCHAR(50) UNIQUE NOT NULL,
    incident_date TIMESTAMPTZ NOT NULL,
    incident_type VARCHAR(50) CHECK (incident_type IN (
        'PHISHING','MALWARE','RANSOMWARE','DDoS','DATA_BREACH','UNAUTHORIZED_ACCESS',
        'INSIDER_THREAT','SOCIAL_ENGINEERING','WEB_DEFACEMENT','FRAUD','OTHER'
    )),
    severity VARCHAR(20) CHECK (severity IN ('CRITICAL','HIGH','MEDIUM','LOW')),
    affected_entity VARCHAR(255),
    affected_sector VARCHAR(100),
    description TEXT,
    iocs_involved JSONB DEFAULT '[]',
    compromised_data TEXT[],
    financial_loss DECIMAL(20,2),
    remediation_status VARCHAR(30),
    case_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE iocs ENABLE ROW LEVEL SECURITY;
CREATE POLICY cyber_soc_select ON iocs FOR SELECT USING (
    current_setting('snisid.agency') IN ('SOC','CYBER_DEFENSE','PNH','DCPJ','SNISID_ADMIN')
);
