-- SNISID-BIO-ADN : Schémas de base de données
-- Index ADN (CODIS-inspired), Personnes (NCIC-inspired), Biens (NCIC-inspired)

CREATE SCHEMA IF NOT EXISTS snisid_bio_adn;

-- ============================
-- INDEX ADN (Catégorie A)
-- ============================

CREATE TYPE snisid_bio_adn.dna_profile_status AS ENUM (
    'PENDING_ANALYSIS',
    'ANALYZED',
    'MATCHED',
    'EXCLUDED',
    'ARCHIVED'
);

CREATE TYPE snisid_bio_adn.dna_match_type AS ENUM (
    'EXACT',
    'PARTIAL',
    'Familial',
    'NONE'
);

CREATE TABLE snisid_bio_adn.dna_profiles (
    profile_id          UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    niu                 CHAR(10),
    profile_type        VARCHAR(10) NOT NULL CHECK (profile_type IN ('BIO-CON','BIO-ARR','BIO-FSC','BIO-DIS','BIO-RNI')),
    case_reference      VARCHAR(100),
    lab_case_number     VARCHAR(100),
    locus_data          JSONB NOT NULL,
    profile_hash        CHAR(64) NOT NULL,
    status              snisid_bio_adn.dna_profile_status NOT NULL DEFAULT 'PENDING_ANALYSIS',
    submitting_agency   VARCHAR(100) NOT NULL,
    submitting_officer  CHAR(10),
    analysis_date       DATE,
    expiration_date     DATE,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dna_niu ON snisid_bio_adn.dna_profiles (niu) WHERE niu IS NOT NULL;
CREATE INDEX idx_dna_type ON snisid_bio_adn.dna_profiles (profile_type);
CREATE INDEX idx_dna_status ON snisid_bio_adn.dna_profiles (status);
CREATE INDEX idx_dna_case ON snisid_bio_adn.dna_profiles (case_reference) WHERE case_reference IS NOT NULL;

CREATE TABLE snisid_bio_adn.dna_matches (
    match_id            UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    query_profile_id    UUID NOT NULL REFERENCES snisid_bio_adn.dna_profiles(profile_id),
    candidate_profile_id UUID NOT NULL REFERENCES snisid_bio_adn.dna_profiles(profile_id),
    match_type          snisid_bio_adn.dna_match_type NOT NULL,
    match_score         DECIMAL(8,6) CHECK (match_score BETWEEN 0 AND 1),
    match_date          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_by         CHAR(10),
    review_status       VARCHAR(20) DEFAULT 'PENDING' CHECK (review_status IN ('PENDING','CONFIRMED','EXCLUDED')),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dna_match_query ON snisid_bio_adn.dna_matches (query_profile_id);
CREATE INDEX idx_dna_match_candidate ON snisid_bio_adn.dna_matches (candidate_profile_id);
CREATE INDEX idx_dna_match_type ON snisid_bio_adn.dna_matches (match_type, match_date);

-- ============================
-- INDEX PERSONNES (Catégorie B)
-- ============================

CREATE TYPE snisid_bio_adn.person_record_type AS ENUM (
    'PER-REC', 'PER-FUG', 'PER-DIS', 'PER-NID',
    'PER-SEX', 'PER-OPR', 'PER-GNG', 'PER-TER',
    'PER-VIO', 'PER-IDV', 'PER-LIB'
);

CREATE TYPE snisid_bio_adn.person_status AS ENUM (
    'ACTIVE',
    'INACTIVE',
    'LOCATED',
    'DECEASED',
    'APPREHENDED'
);

CREATE TABLE snisid_bio_adn.person_records (
    record_id           UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    niu                 CHAR(10) NOT NULL,
    record_type         snisid_bio_adn.person_record_type NOT NULL,
    status              snisid_bio_adn.person_status NOT NULL DEFAULT 'ACTIVE',
    priority_level      SMALLINT DEFAULT 3 CHECK (priority_level BETWEEN 1 AND 5),
    subject_name        VARCHAR(200),
    subject_alias       TEXT[],
    subject_description TEXT,
    photo_refs          TEXT[],
    fingerprint_id      UUID,
    last_known_location VARCHAR(300),
    last_known_dept     VARCHAR(10),
    risk_assessment     VARCHAR(20) DEFAULT 'UNKNOWN',
    warrant_id          UUID,
    case_reference      VARCHAR(100),
    reporting_agency    VARCHAR(100) NOT NULL,
    reporting_officer   CHAR(10),
    alert_expiry        TIMESTAMPTZ,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_person_niu ON snisid_bio_adn.person_records (niu);
CREATE INDEX idx_person_type ON snisid_bio_adn.person_records (record_type);
CREATE INDEX idx_person_status ON snisid_bio_adn.person_records (status, is_active);
CREATE INDEX idx_person_dept ON snisid_bio_adn.person_records (last_known_dept);
CREATE INDEX idx_person_name ON snisid_bio_adn.person_records USING gin(
    to_tsvector('french', COALESCE(subject_name, ''))
);

-- ============================
-- INDEX BIENS (Catégorie C)
-- ============================

CREATE TYPE snisid_bio_adn.property_record_type AS ENUM (
    'BIE-VEH', 'BIE-ARM', 'BIE-DOC', 'BIE-OBJ',
    'BIE-PLQ', 'BIE-TIT', 'BIE-EMB'
);

CREATE TYPE snisid_bio_adn.property_status AS ENUM (
    'STOLEN',
    'RECOVERED',
    'RESTITUTED',
    'DESTROYED'
);

CREATE TABLE snisid_bio_adn.property_records (
    record_id           UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_type         snisid_bio_adn.property_record_type NOT NULL,
    status              snisid_bio_adn.property_status NOT NULL DEFAULT 'STOLEN',
    item_description    TEXT NOT NULL,
    serial_number       VARCHAR(200),
    make                VARCHAR(100),
    model               VARCHAR(100),
    color               VARCHAR(50),
    year_of_manufacture SMALLINT,
    registration_number VARCHAR(100),
    vin                 VARCHAR(17),
    plate_number        VARCHAR(20),
    calibre             VARCHAR(50),
    document_type       VARCHAR(100),
    estimated_value     DECIMAL(12,2),
    theft_date          TIMESTAMPTZ NOT NULL,
    theft_location      VARCHAR(300),
    theft_dept          VARCHAR(10),
    recovery_date       TIMESTAMPTZ,
    recovery_location   VARCHAR(300),
    case_reference      VARCHAR(100),
    reporting_agency    VARCHAR(100) NOT NULL,
    linked_person_niu   CHAR(10),
    foves_vehicle_id    UUID,
    lapi_alert_id       UUID,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_property_type ON snisid_bio_adn.property_records (record_type);
CREATE INDEX idx_property_status ON snisid_bio_adn.property_records (status, is_active);
CREATE INDEX idx_property_serial ON snisid_bio_adn.property_records (serial_number) WHERE serial_number IS NOT NULL;
CREATE INDEX idx_property_vin ON snisid_bio_adn.property_records (vin) WHERE vin IS NOT NULL;
CREATE INDEX idx_property_plate ON snisid_bio_adn.property_records (plate_number) WHERE plate_number IS NOT NULL;
CREATE INDEX idx_property_dept ON snisid_bio_adn.property_records (theft_dept);
CREATE INDEX idx_property_person ON snisid_bio_adn.property_records (linked_person_niu) WHERE linked_person_niu IS NOT NULL;
