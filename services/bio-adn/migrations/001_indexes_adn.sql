CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE SCHEMA IF NOT EXISTS bio_adn;
SET search_path TO bio_adn, public;

CREATE TABLE bio_laboratories (
    lab_id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    lab_code        VARCHAR(20) UNIQUE NOT NULL,
    lab_name        VARCHAR(200) NOT NULL,
    lab_level       VARCHAR(10) NOT NULL CHECK (lab_level IN ('LDIS','SDIS','NDIS')),
    department      VARCHAR(50),
    institution     VARCHAR(100),
    accreditation   VARCHAR(100),
    contact_email   VARCHAR(200),
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE bio_str_profiles (
    sample_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    specimen_number VARCHAR(100) UNIQUE NOT NULL,
    index_type      VARCHAR(10) NOT NULL CHECK (index_type IN ('BIO-CON','BIO-ARR','BIO-FSC','BIO-DIS','BIO-RNI')),
    loci_encrypted  BYTEA NOT NULL,
    loci_hash       VARCHAR(64) NOT NULL,
    amelogenin      CHAR(2),
    quality_score   DECIMAL(4,3) CHECK (quality_score BETWEEN 0 AND 1),
    loci_count      SMALLINT DEFAULT 20,
    lab_id          UUID REFERENCES bio_laboratories(lab_id),
    case_number     VARCHAR(100),
    collected_date  DATE NOT NULL,
    analysis_date   DATE,
    uploaded_ldis   BOOLEAN DEFAULT FALSE,
    uploaded_sdis   BOOLEAN DEFAULT FALSE,
    uploaded_ndis   BOOLEAN DEFAULT FALSE,
    ndis_upload_date TIMESTAMPTZ,
    is_expunged     BOOLEAN DEFAULT FALSE,
    expunge_date    TIMESTAMPTZ,
    expunge_order   VARCHAR(200),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_bio_str_hash ON bio_str_profiles(loci_hash);
CREATE INDEX idx_bio_str_index_type ON bio_str_profiles(index_type);
CREATE INDEX idx_bio_str_case ON bio_str_profiles(case_number);

CREATE TABLE bio_identity_links (
    link_id         UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    sample_id       UUID UNIQUE REFERENCES bio_str_profiles(sample_id),
    niu             VARCHAR(20),
    linked_by_agent VARCHAR(100) NOT NULL,
    linked_at       TIMESTAMPTZ DEFAULT NOW(),
    court_order_ref VARCHAR(200) NOT NULL,
    purpose         VARCHAR(100) NOT NULL,
    reviewed_by     VARCHAR(100),
    reviewed_at     TIMESTAMPTZ,
    review_outcome  VARCHAR(20)
);

CREATE TABLE bio_hits (
    hit_id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    query_sample_id UUID REFERENCES bio_str_profiles(sample_id),
    match_sample_id UUID REFERENCES bio_str_profiles(sample_id),
    match_type      VARCHAR(20) NOT NULL CHECK (match_type IN ('FULL_MATCH','PARTIAL','FAMILIAL')),
    confidence      DECIMAL(5,4) NOT NULL,
    matched_loci    SMALLINT NOT NULL,
    total_loci      SMALLINT NOT NULL,
    hit_level       VARCHAR(10) NOT NULL CHECK (hit_level IN ('LDIS','SDIS','NDIS')),
    alert_sent      BOOLEAN DEFAULT FALSE,
    alert_sent_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
