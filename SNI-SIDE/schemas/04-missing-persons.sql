-- ============================================================
-- SNI-SIDE: Missing Persons Database
-- PostgreSQL 16 + PostGIS
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_missing;
SET search_path TO snisid_missing;

CREATE EXTENSION IF NOT EXISTS postgis;

-- ============ MISSING PERSONS ============
CREATE TABLE missing_persons (
    missing_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_type VARCHAR(30) CHECK (case_type IN ('CHILD','ADULT','KIDNAPPING','TRAFFICKING','FOUND')),
    niu VARCHAR(10),
    full_name VARCHAR(255) NOT NULL,
    alias VARCHAR(255),
    date_of_birth DATE,
    age_at_disappearance INT,
    gender VARCHAR(1),
    nationality VARCHAR(100),
    height_cm DECIMAL(5,2),
    weight_kg DECIMAL(5,2),
    eye_color VARCHAR(50),
    hair_color VARCHAR(50),
    skin_tone VARCHAR(50),
    scars_marks TEXT,
    last_seen_date TIMESTAMPTZ NOT NULL,
    last_seen_location VARCHAR(500),
    last_seen_lat DECIMAL(10,7),
    last_seen_lng DECIMAL(10,7),
    last_seen_geom GEOMETRY(Point, 4326),
    last_seen_circumstances TEXT,
    last_known_phone VARCHAR(50),
    last_known_email VARCHAR(255),
    clothing_description TEXT,
    medical_conditions TEXT,
    medications TEXT,
    risk_factors TEXT[],
    photos JSONB DEFAULT '[]',
    biometric_references JSONB DEFAULT '{}',
    dna_profile_id UUID,
    reported_by VARCHAR(255),
    reporting_agency VARCHAR(100),
    case_officer VARCHAR(255),
    status VARCHAR(20) CHECK (status IN (
        'MISSING','FOUND_SAFE','FOUND_DECEASED','FOUND_UNIDENTIFIED','ACTIVE_SEARCH','CLOSED'
    )),
    risk_level VARCHAR(20) CHECK (risk_level IN ('CRITICAL','HIGH','MEDIUM','LOW')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_citizen FOREIGN KEY (niu) REFERENCES snisid_identity.citizens(niu),
    CONSTRAINT fk_dna FOREIGN KEY (dna_profile_id) REFERENCES snisid_codis.dna_profiles(profile_id)
) PARTITION BY LIST (case_type);

CREATE INDEX idx_missing_name ON missing_persons USING gin(to_tsvector('french', full_name));
CREATE INDEX idx_missing_status ON missing_persons(status);
CREATE INDEX idx_missing_risk ON missing_persons(risk_level);
CREATE INDEX idx_missing_date ON missing_persons(last_seen_date DESC);
CREATE INDEX idx_missing_location ON missing_persons USING gist(last_seen_geom);
CREATE INDEX idx_missing_agency ON missing_persons(reporting_agency);
CREATE INDEX idx_missing_dna ON missing_persons(dna_profile_id);
CREATE INDEX idx_missing_phone ON missing_persons(last_known_phone);

-- ============ SIGHTINGS ============
CREATE TABLE sightings (
    sighting_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    missing_id UUID NOT NULL REFERENCES missing_persons(missing_id),
    sighting_date TIMESTAMPTZ NOT NULL,
    sighting_location VARCHAR(500),
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    location_geom GEOMETRY(Point, 4326),
    description TEXT,
    witness_name VARCHAR(255),
    witness_contact VARCHAR(100),
    source_type VARCHAR(30) CHECK (source_type IN (
        'WITNESS','CAMERA','POLICE_PATROL','BORDER_CROSSING','SOCIAL_MEDIA','ANONYMOUS_TIP','ALPR'
    )),
    confidence_score DECIMAL(5,2),
    evidence_files JSONB DEFAULT '[]',
    status VARCHAR(20) CHECK (status IN ('PENDING_VERIFICATION','VERIFIED','FALSE_POSITIVE','ACTIONED')),
    verified_by VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sighting_missing ON sightings(missing_id);
CREATE INDEX idx_sighting_date ON sightings(sighting_date DESC);
CREATE INDEX idx_sighting_location ON sightings USING gist(location_geom);
CREATE INDEX idx_sighting_source ON sightings(source_type);
CREATE INDEX idx_sighting_status ON sightings(status);
CREATE INDEX idx_sighting_confidence ON sightings(confidence_score DESC);

-- ============ KIDNAPPING DETAILS ============
CREATE TABLE kidnapping_details (
    kidnapping_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    missing_id UUID NOT NULL REFERENCES missing_persons(missing_id),
    abduction_date TIMESTAMPTZ NOT NULL,
    abduction_location VARCHAR(500),
    abduction_method VARCHAR(100),
    suspects TEXT[],
    vehicles_involved JSONB DEFAULT '[]',
    ransom_demanded BOOLEAN DEFAULT FALSE,
    ransom_amount DECIMAL(20,2),
    ransom_currency VARCHAR(10),
    ransom_paid BOOLEAN DEFAULT FALSE,
    communication_method VARCHAR(100),
    negotiation_status VARCHAR(30),
    threat_level VARCHAR(20),
    tactical_team_deployed BOOLEAN DEFAULT FALSE,
    case_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============ TRAFFICKING CASES ============
CREATE TABLE trafficking_cases (
    trafficking_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    missing_id UUID NOT NULL REFERENCES missing_persons(missing_id),
    trafficking_type VARCHAR(30) CHECK (trafficking_type IN (
        'SEXUAL','LABOR','ORGAN','CHILD_SOLDIER','FORCED_MARRIAGE','OTHER'
    )),
    recruiter_name VARCHAR(255),
    recruiter_method VARCHAR(255),
    transit_routes TEXT[],
    destination_countries VARCHAR(100)[],
    intermediary_names TEXT[],
    network_id UUID,
    case_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============ FOUND PERSONS ============
CREATE TABLE found_persons (
    found_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    missing_id UUID REFERENCES missing_persons(missing_id),
    found_date TIMESTAMPTZ NOT NULL,
    found_location VARCHAR(500),
    found_by_agency VARCHAR(100),
    condition_found VARCHAR(50) CHECK (condition_found IN (
        'ALIVE_WELL','ALIVE_INJURED','ALIVE_TRAFFICKED','DECEASED','UNIDENTIFIED'
    )),
    identifying_features TEXT,
    biometric_verified BOOLEAN DEFAULT FALSE,
    dna_confirmed BOOLEAN DEFAULT FALSE,
    reunification_date TIMESTAMPTZ,
    reunified_with VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============ AMBER/SILVER ALERTS ============
CREATE TABLE alerts_missing (
    alert_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    missing_id UUID NOT NULL REFERENCES missing_persons(missing_id),
    alert_type VARCHAR(20) CHECK (alert_type IN ('AMBER','SILVER','CLEAR','NATIONAL')),
    activated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deactivated_at TIMESTAMPTZ,
    broadcast_regions VARCHAR(100)[],
    channel VARCHAR(50) CHECK (channel IN (
        'SMS','EMAIL','SOCIAL_MEDIA','TV','RADIO','ALPR','BORDER','ALL'
    )),
    status VARCHAR(20) CHECK (status IN ('ACTIVE','DEACTIVATED','EXPIRED')),
    resolution VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_alert_missing ON alerts_missing(missing_id);
CREATE INDEX idx_alert_type ON alerts_missing(alert_type);
CREATE INDEX idx_alert_status ON alerts_missing(status);

-- ============ RLS ============
ALTER TABLE missing_persons ENABLE ROW LEVEL SECURITY;
CREATE POLICY missing_pnh_select ON missing_persons FOR SELECT USING (
    current_setting('snisid.agency') IN ('PNH','DCPJ','IMMIGRATION','SNISID_ADMIN')
);
