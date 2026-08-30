-- ============================================================
-- SNI-SIDE: HN-CODIS (Combined DNA Index System)
-- PostgreSQL 16
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_codis;
SET search_path TO snisid_codis;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============ DNA PROFILES ============
CREATE TABLE dna_profiles (
    profile_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_type VARCHAR(30) CHECK (profile_type IN (
        'CRIMINAL','CRIME_SCENE','MISSING_PERSON','HUMAN_REMAINS','FAMILIAL','VOLUNTARY'
    )),
    niu VARCHAR(10),
    person_name VARCHAR(255),
    sample_id VARCHAR(100) UNIQUE NOT NULL,
    laboratory_id VARCHAR(100) NOT NULL,
    analyst_id VARCHAR(100) NOT NULL,
    collection_date DATE NOT NULL,
    collection_location VARCHAR(500),
    collection_method VARCHAR(100),
    sample_type VARCHAR(100) CHECK (sample_type IN (
        'BLOOD','SALIVA','HAIR','TISSUE','BONE','SEMEN','SKIN','OTHER'
    )),
    dna_profile_hash VARCHAR(64) UNIQUE NOT NULL,
    locus_set JSONB NOT NULL,
    alleles TEXT[] NOT NULL,
    profile_quality VARCHAR(20) CHECK (profile_quality IN ('FULL','PARTIAL','DEGRADED','MIXED')),
    rar_score DECIMAL(10,8),
    status VARCHAR(20) CHECK (status IN ('ACTIVE','ARCHIVED','EXPUNGED','PENDING')),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_citizen FOREIGN KEY (niu) REFERENCES snisid_identity.citizens(niu)
) PARTITION BY LIST (profile_type);

CREATE INDEX idx_dna_type ON dna_profiles(profile_type);
CREATE INDEX idx_dna_niu ON dna_profiles(niu);
CREATE INDEX idx_dna_hash ON dna_profiles(dna_profile_hash);
CREATE INDEX idx_dna_status ON dna_profiles(status);
CREATE INDEX idx_dna_rar ON dna_profiles(rar_score DESC);
CREATE INDEX idx_dna_created ON dna_profiles(created_at DESC);

-- ============ DNA MATCHES ============
CREATE TABLE dna_matches (
    match_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id_1 UUID NOT NULL REFERENCES dna_profiles(profile_id),
    profile_id_2 UUID NOT NULL REFERENCES dna_profiles(profile_id),
    match_type VARCHAR(30) CHECK (match_type IN (
        'EXACT','FAMILIAL','PARTIAL','MITOCHONDRIAL','Y_STR'
    )),
    match_probability DECIMAL(15,12) NOT NULL,
    random_match_probability DECIMAL(15,12),
    shared_alleles INT,
    total_loci INT,
    statistical_weight DECIMAL(20,10),
    relationship VARCHAR(50),
    case_id UUID,
    status VARCHAR(20) CHECK (status IN ('PENDING_REVIEW','CONFIRMED','FALSE_POSITIVE','RESOLVED')),
    reviewed_by VARCHAR(100),
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_dna_match UNIQUE (profile_id_1, profile_id_2),
    CONSTRAINT chk_different_profiles CHECK (profile_id_1 <> profile_id_2)
);

CREATE INDEX idx_dna_match_p1 ON dna_matches(profile_id_1);
CREATE INDEX idx_dna_match_p2 ON dna_matches(profile_id_2);
CREATE INDEX idx_dna_match_type ON dna_matches(match_type);
CREATE INDEX idx_dna_match_prob ON dna_matches(match_probability DESC);
CREATE INDEX idx_dna_match_status ON dna_matches(status);

-- ============ CODIS LABORATORIES ============
CREATE TABLE codis_laboratories (
    lab_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(20) UNIQUE NOT NULL,
    type VARCHAR(30) CHECK (type IN ('NATIONAL','REGIONAL','MOBILE','CONTRACTOR')),
    address TEXT,
    accreditation VARCHAR(100),
    contact_person VARCHAR(255),
    contact_phone VARCHAR(50),
    contact_email VARCHAR(255),
    status VARCHAR(20) CHECK (status IN ('ACTIVE','SUSPENDED','DECOMMISSIONED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_lab_code ON codis_laboratories(code);
CREATE INDEX idx_lab_type ON codis_laboratories(type);

-- ============ DNA CHAIN OF CUSTODY ============
CREATE TABLE dna_chain_of_custody (
    custody_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sample_id VARCHAR(100) NOT NULL,
    profile_id UUID REFERENCES dna_profiles(profile_id),
    action VARCHAR(50) NOT NULL,
    actor_id VARCHAR(100) NOT NULL,
    actor_name VARCHAR(255) NOT NULL,
    actor_agency VARCHAR(100) NOT NULL,
    previous_location VARCHAR(500),
    new_location VARCHAR(500),
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    signature_hash VARCHAR(64),
    notes TEXT
);

CREATE INDEX idx_custody_sample ON dna_chain_of_custody(sample_id);
CREATE INDEX idx_custody_profile ON dna_chain_of_custody(profile_id);
CREATE INDEX idx_custody_time ON dna_chain_of_custody(timestamp);

-- ============ FAMILIAL DNA SEARCH LOG ============
CREATE TABLE familial_dna_searches (
    search_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query_profile_id UUID NOT NULL REFERENCES dna_profiles(profile_id),
    search_type VARCHAR(30) CHECK (search_type IN ('PARENT','CHILD','SIBLING','COUSIN','GRANDPARENT')),
    min_threshold DECIMAL(10,8),
    results_count INT,
    top_match_id UUID REFERENCES dna_profiles(profile_id),
    top_match_score DECIMAL(10,8),
    case_id UUID,
    authorized_by VARCHAR(100),
    authorization_ref VARCHAR(100),
    search_duration_ms INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_familial_query ON familial_dna_searches(query_profile_id);
CREATE INDEX idx_familial_time ON familial_dna_searches(created_at DESC);

-- ============ RLS ============
ALTER TABLE dna_profiles ENABLE ROW LEVEL LEVEL SECURITY;

CREATE POLICY codis_forensic_select ON dna_profiles FOR SELECT USING (
    current_setting('snisid.agency') IN ('PNH','DCPJ','LABORATORY','SNISID_ADMIN')
);

CREATE POLICY codis_forensic_insert ON dna_profiles FOR INSERT WITH CHECK (
    current_setting('snisid.agency') IN ('LABORATORY','DCPJ')
);
