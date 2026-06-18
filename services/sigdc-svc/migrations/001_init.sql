CREATE TYPE sigdc_disaster_type AS ENUM (
    'EARTHQUAKE', 'HURRICANE', 'TSUNAMI', 'FLOOD',
    'LANDSLIDE', 'FIRE_MASS', 'INDUSTRIAL_ACCIDENT',
    'EPIDEMIC', 'SECURITY_MASS_CASUALTY'
);

CREATE TYPE sigdc_alert_level AS ENUM (
    'WATCH', 'WARNING', 'EMERGENCY', 'CATASTROPHE'
);

CREATE TABLE sigdc_disasters (
    disaster_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_sigdc_id   VARCHAR(25) UNIQUE NOT NULL,
    disaster_type       sigdc_disaster_type NOT NULL,
    disaster_name       VARCHAR(200),
    alert_level         sigdc_alert_level NOT NULL,
    status              VARCHAR(20) DEFAULT 'ACTIVE',
    onset_date          TIMESTAMPTZ NOT NULL,
    affected_depts      CHAR(2)[] DEFAULT '{}',
    affected_communes   TEXT[] DEFAULT '{}',
    epicenter_lat       DECIMAL(10,7),
    epicenter_lng       DECIMAL(10,7),
    magnitude           DECIMAL(4,2),
    wind_speed_kmh      INTEGER,
    estimated_affected  INTEGER,
    confirmed_dead      INTEGER DEFAULT 0,
    confirmed_injured   INTEGER DEFAULT 0,
    confirmed_missing   INTEGER DEFAULT 0,
    confirmed_displaced INTEGER DEFAULT 0,
    response_agencies   TEXT[] DEFAULT '{}',
    coordination_center VARCHAR(200),
    ocha_flash_ref      VARCHAR(100),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sigdc_victim_registrations (
    registration_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disaster_id         UUID NOT NULL REFERENCES sigdc_disasters(disaster_id),
    snisid_person_id    UUID,
    full_name           VARCHAR(200),
    dob                 DATE,
    gender              VARCHAR(10),
    status              VARCHAR(30) NOT NULL,
    injury_description  TEXT,
    location_found      VARCHAR(300),
    dept_code           CHAR(2),
    hospital_sent_to    VARCHAR(150),
    morgue_location     VARCHAR(150),
    afis_subject_id     UUID,
    dna_sample_taken    BOOLEAN DEFAULT FALSE,
    dna_sample_ref      VARCHAR(100),
    rvin_case_id        UUID,
    dpide_idp_id        UUID,
    registration_date   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    registered_by       UUID NOT NULL,
    org_registering     VARCHAR(100)
);

CREATE TABLE sigdc_resources (
    resource_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disaster_id         UUID NOT NULL REFERENCES sigdc_disasters(disaster_id),
    resource_type       VARCHAR(50) NOT NULL,
    provider_org        VARCHAR(150),
    quantity            INTEGER,
    unit                VARCHAR(30),
    location_lat        DECIMAL(10,7),
    location_lng        DECIMAL(10,7),
    dept_code           CHAR(2),
    available_from      TIMESTAMPTZ,
    available_until     TIMESTAMPTZ,
    status              VARCHAR(20) DEFAULT 'AVAILABLE',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sigdc_early_warnings (
    warning_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disaster_type       sigdc_disaster_type NOT NULL,
    alert_level         sigdc_alert_level NOT NULL,
    source_agency       VARCHAR(100),
    message_text        TEXT NOT NULL,
    affected_depts      CHAR(2)[] DEFAULT '{}',
    issued_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at          TIMESTAMPTZ,
    channels_sent       TEXT[] DEFAULT '{}',
    population_reached  INTEGER
);

CREATE INDEX idx_sigdc_disasters_status ON sigdc_disasters(status, disaster_type);
CREATE INDEX idx_sigdc_disasters_dept ON sigdc_disasters USING gin(affected_depts);
CREATE INDEX idx_sigdc_victims_disaster ON sigdc_victim_registrations(disaster_id, status);
CREATE INDEX idx_sigdc_resources ON sigdc_resources(disaster_id, resource_type, status);
CREATE INDEX idx_sigdc_warnings_date ON sigdc_early_warnings(issued_at DESC, alert_level);
