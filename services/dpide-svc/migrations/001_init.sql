CREATE TYPE dpide_displacement_cause AS ENUM (
    'GANG_VIOLENCE', 'EARTHQUAKE', 'HURRICANE',
    'FLOOD', 'FIRE', 'POLITICAL_VIOLENCE', 'OTHER'
);

CREATE TYPE dpide_idp_status AS ENUM (
    'DISPLACED', 'IN_CAMP', 'WITH_HOST_FAMILY',
    'RELOCATED', 'RETURNED_HOME', 'EMIGRATED', 'DECEASED'
);

CREATE TABLE dpide_idps (
    idp_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_dpide_id   VARCHAR(25) UNIQUE NOT NULL,
    snisid_person_id    UUID,
    full_name           VARCHAR(200) NOT NULL,
    dob                 DATE,
    gender              VARCHAR(10),
    household_size      SMALLINT DEFAULT 1,
    minors_count        SMALLINT DEFAULT 0,
    displacement_cause  dpide_displacement_cause NOT NULL,
    displacement_date   TIMESTAMPTZ NOT NULL,
    origin_address      TEXT,
    origin_dept_code    CHAR(2) NOT NULL,
    origin_commune      VARCHAR(100),
    status              dpide_idp_status NOT NULL DEFAULT 'DISPLACED',
    current_location    TEXT,
    current_dept_code   CHAR(2),
    current_commune     VARCHAR(100),
    current_lat         DECIMAL(10,7),
    current_lng         DECIMAL(10,7),
    camp_id             UUID,
    shelter_type        VARCHAR(50),
    has_nfi             BOOLEAN DEFAULT FALSE,
    receives_food_aid   BOOLEAN DEFAULT FALSE,
    has_latrines        BOOLEAN DEFAULT FALSE,
    has_water_access    BOOLEAN DEFAULT FALSE,
    medical_needs       TEXT[] DEFAULT '{}',
    iom_dtm_ref         VARCHAR(50),
    ocha_ref            VARCHAR(50),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE dpide_camps (
    camp_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    camp_name           VARCHAR(150) NOT NULL,
    dept_code           CHAR(2) NOT NULL,
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    displacement_cause  dpide_displacement_cause,
    managing_org        VARCHAR(150),
    capacity            INTEGER,
    current_population  INTEGER DEFAULT 0,
    is_active           BOOLEAN DEFAULT TRUE,
    has_medical_post    BOOLEAN DEFAULT FALSE,
    has_school          BOOLEAN DEFAULT FALSE,
    water_source        TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dpide_status ON dpide_idps(status, displacement_cause);
CREATE INDEX idx_dpide_dept ON dpide_idps(current_dept_code) WHERE status IN ('DISPLACED','IN_CAMP');
CREATE INDEX idx_dpide_cause ON dpide_idps(displacement_cause);
CREATE INDEX idx_dpide_camp ON dpide_idps(camp_id) WHERE camp_id IS NOT NULL;
CREATE INDEX idx_dpide_snisid ON dpide_idps(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
