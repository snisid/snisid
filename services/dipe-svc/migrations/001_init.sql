BEGIN;
CREATE TYPE dipe_case_type AS ENUM ('KIDNAPPING_SUSPECTED','VOLUNTARY_DISAPPEARANCE','DISASTER_RELATED','GANG_VIOLENCE','MIGRATION_RELATED','CHILD_ABDUCTION','TRAFFICKING_SUSPECTED','UNKNOWN');
CREATE TYPE dipe_case_status AS ENUM ('OPEN','LOCATED_ALIVE','BODY_IDENTIFIED','BODY_UNIDENTIFIED','CANCELLED','COLD_CASE');
CREATE TABLE dipe_missing_persons (
    case_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), national_dipe_id VARCHAR(25) UNIQUE NOT NULL,
    case_type dipe_case_type NOT NULL, status dipe_case_status NOT NULL DEFAULT 'OPEN',
    snisid_person_id UUID, full_name VARCHAR(200) NOT NULL, aliases TEXT[] DEFAULT '{}',
    dob DATE, gender VARCHAR(10), nationality CHAR(3) DEFAULT 'HTI',
    occupation VARCHAR(100), photo_refs TEXT[] DEFAULT '{}',
    height_cm SMALLINT, weight_kg SMALLINT, skin_tone VARCHAR(30),
    eye_color VARCHAR(30), hair_color VARCHAR(30), distinguishing_marks TEXT,
    clothing_last_seen TEXT, last_seen_date TIMESTAMPTZ NOT NULL,
    last_seen_location VARCHAR(300), last_seen_dept_code CHAR(2),
    last_seen_commune VARCHAR(100), last_seen_lat DECIMAL(10,7),
    last_seen_lng DECIMAL(10,7), circumstances TEXT,
    sivc_alert_id UUID, gang_id UUID, extors_case_id UUID,
    reported_by_name VARCHAR(200), reported_by_phone VARCHAR(30),
    reported_by_snisid UUID, report_date TIMESTAMPTZ NOT NULL,
    reporting_unit VARCHAR(50), afis_subject_id UUID,
    dna_sample_ref VARCHAR(100), dna_profile_id UUID,
    interpol_notice_ref VARCHAR(50), ncmec_ref VARCHAR(50),
    resolution_date TIMESTAMPTZ, resolution_notes TEXT,
    rvin_case_id UUID, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE dipe_sightings (
    sighting_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id UUID NOT NULL REFERENCES dipe_missing_persons(case_id),
    sighting_date TIMESTAMPTZ NOT NULL, location_desc VARCHAR(300),
    dept_code CHAR(2), lat DECIMAL(10,7), lng DECIMAL(10,7),
    reported_by UUID, report_method VARCHAR(30), confidence SMALLINT,
    photo_ref VARCHAR(500), verified BOOLEAN DEFAULT FALSE,
    verified_by UUID, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE dipe_disaster_missing (
    disaster_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id UUID NOT NULL REFERENCES dipe_missing_persons(case_id),
    disaster_type VARCHAR(30) NOT NULL, disaster_name VARCHAR(100),
    disaster_date DATE NOT NULL, last_known_address TEXT,
    shelter_checked TEXT[] DEFAULT '{}', hospital_checked TEXT[] DEFAULT '{}',
    morgue_checked TEXT[] DEFAULT '{}', rc_haiti_ref VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_dipe_status ON dipe_missing_persons(status, last_seen_date DESC);
CREATE INDEX idx_dipe_type ON dipe_missing_persons(case_type) WHERE status = 'OPEN';
CREATE INDEX idx_dipe_dept ON dipe_missing_persons(last_seen_dept_code) WHERE status = 'OPEN';
CREATE INDEX idx_dipe_person ON dipe_missing_persons(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_dipe_sightings ON dipe_sightings(case_id, sighting_date DESC);
CREATE INDEX idx_dipe_name_fts ON dipe_missing_persons USING gin(to_tsvector('simple', full_name));
COMMIT;
