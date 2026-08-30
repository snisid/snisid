CREATE TYPE enfl_risk_category AS ENUM (
    'MISSING_ABDUCTION', 'GANG_RECRUITMENT', 'DOMESTIC_SERVITUDE_RESTAVEK',
    'SEXUAL_EXPLOITATION', 'TRAFFICKING', 'UNACCOMPANIED_MIGRANT',
    'SEPARATED_DISASTER', 'STREET_CHILD', 'OTHER'
);

CREATE TYPE enfl_status AS ENUM (
    'AT_RISK', 'MISSING', 'LOCATED_SAFE', 'LOCATED_AT_RISK',
    'IN_CARE', 'REPATRIATED', 'DECEASED'
);

CREATE TABLE enfl_children (
    child_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_enfl_id    VARCHAR(25) UNIQUE NOT NULL,
    snisid_person_id    UUID,
    dipe_case_id        UUID,
    trait_case_id       UUID,
    risk_category       enfl_risk_category NOT NULL,
    status              enfl_status NOT NULL DEFAULT 'MISSING',
    full_name           VARCHAR(200) NOT NULL,
    dob                 DATE NOT NULL,
    age_at_registration SMALLINT,
    gender              VARCHAR(10),
    nationality         CHAR(3) DEFAULT 'HTI',
    photo_refs          TEXT[] DEFAULT '{}',
    distinguishing_marks TEXT,
    height_cm           SMALLINT,
    skin_tone           VARCHAR(30),
    guardian_name       VARCHAR(200),
    guardian_phone      VARCHAR(30),
    guardian_snisid_id  UUID,
    last_known_location VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    disappearance_date  TIMESTAMPTZ,
    gang_id             UUID,
    recruiter_snisid_id UUID,
    afis_subject_id     UUID,
    dna_profile_id      UUID,
    interpol_icse_ref   VARCHAR(50),
    ncmec_ref           VARCHAR(50),
    ibesr_ref           VARCHAR(50),
    assistance_type     TEXT[] DEFAULT '{}',
    current_shelter     VARCHAR(200),
    assigned_caseworker UUID,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE enfl_restaveks (
    restavek_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id            UUID NOT NULL REFERENCES enfl_children(child_id),
    employing_household VARCHAR(300),
    household_dept      CHAR(2),
    household_commune   VARCHAR(100),
    employing_person_id UUID,
    reported_conditions TEXT,
    school_attendance   BOOLEAN DEFAULT FALSE,
    ibesr_inspection    BOOLEAN DEFAULT FALSE,
    last_inspection_date DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_enfl_status ON enfl_children(status, risk_category);
CREATE INDEX idx_enfl_dept ON enfl_children(dept_code) WHERE status IN ('MISSING','AT_RISK');
CREATE INDEX idx_enfl_gang ON enfl_children(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_enfl_name_fts ON enfl_children USING gin(to_tsvector('simple', full_name));
