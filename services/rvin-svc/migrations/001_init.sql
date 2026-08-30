CREATE TYPE rvin_source AS ENUM (
    'CRIME_SCENE', 'DISASTER_SITE', 'MASS_GRAVE',
    'RIVER', 'STREET', 'HOSPITAL_DOA', 'OTHER'
);

CREATE TYPE rvin_status AS ENUM (
    'UNIDENTIFIED', 'TENTATIVE_MATCH', 'CONFIRMED_IDENTIFIED', 'CLAIMED'
);

CREATE TABLE rvin_unidentified_remains (
    remains_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_rvin_id    VARCHAR(25) UNIQUE NOT NULL,
    discovery_date      TIMESTAMPTZ NOT NULL,
    discovery_location  VARCHAR(300) NOT NULL,
    dept_code           CHAR(2) NOT NULL,
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    discovery_source    rvin_source NOT NULL,
    status              rvin_status NOT NULL DEFAULT 'UNIDENTIFIED',
    estimated_sex       VARCHAR(10),
    estimated_age_min   SMALLINT,
    estimated_age_max   SMALLINT,
    estimated_height_cm SMALLINT,
    skin_tone           VARCHAR(30),
    hair_type           VARCHAR(50),
    clothing_description TEXT,
    distinguishing_marks TEXT,
    decomposition_level SMALLINT CHECK (decomposition_level BETWEEN 1 AND 5),
    afis_latent_id      UUID,
    dna_sample_taken    BOOLEAN DEFAULT FALSE,
    dna_sample_ref      VARCHAR(100),
    dna_profile_id      UUID,
    dental_chart_ref    VARCHAR(200),
    photo_refs          TEXT[] DEFAULT '{}',
    xray_refs           TEXT[] DEFAULT '{}',
    morgue_location     VARCHAR(200),
    morgue_ref          VARCHAR(50),
    storage_date        TIMESTAMPTZ,
    estimated_death_date TIMESTAMPTZ,
    disaster_id         UUID,
    mass_incident_id    UUID,
    gang_id             UUID,
    case_reference      VARCHAR(100),
    interpol_dvi_ref    VARCHAR(50),
    matched_dipe_case_id UUID,
    matched_snisid_id    UUID,
    identification_method VARCHAR(50),
    identification_date  TIMESTAMPTZ,
    identified_by        UUID,
    examiner_id         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rvin_dna_comparisons (
    comparison_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    remains_id          UUID NOT NULL REFERENCES rvin_unidentified_remains(remains_id),
    dipe_case_id        UUID,
    reference_dna_ref   VARCHAR(100),
    comparison_date     TIMESTAMPTZ NOT NULL,
    match_probability   DECIMAL(10,8),
    is_match            BOOLEAN DEFAULT FALSE,
    lab_reference       VARCHAR(100),
    examiner_id         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rvin_status ON rvin_unidentified_remains(status);
CREATE INDEX idx_rvin_dept ON rvin_unidentified_remains(dept_code, discovery_date DESC);
CREATE INDEX idx_rvin_disaster ON rvin_unidentified_remains(disaster_id) WHERE disaster_id IS NOT NULL;
