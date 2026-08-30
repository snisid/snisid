CREATE TYPE vict_crime_type AS ENUM (
    'HOMICIDE', 'MASS_KILLING', 'RAPE', 'GANG_RAPE',
    'TORTURE', 'FORCED_DISAPPEARANCE', 'EXTRAJUDICIAL_KILLING',
    'KIDNAPPING_VICTIM', 'MUTILATION', 'OTHER_GRAVE'
);

CREATE TYPE vict_victim_status AS ENUM (
    'ALIVE_SURVIVOR', 'DECEASED_IDENTIFIED',
    'DECEASED_UNIDENTIFIED', 'MISSING_PRESUMED_DEAD'
);

CREATE TABLE vict_victims (
    victim_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_vict_id    VARCHAR(25) UNIQUE NOT NULL,
    snisid_person_id    UUID,
    crime_type          vict_crime_type NOT NULL,
    victim_status       vict_victim_status NOT NULL,
    full_name           VARCHAR(200),
    dob                 DATE,
    gender              VARCHAR(10),
    nationality         CHAR(3) DEFAULT 'HTI',
    occupation          VARCHAR(100),
    incident_date       TIMESTAMPTZ NOT NULL,
    incident_location   VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    perpetrator_ids     UUID[] DEFAULT '{}',
    gang_id             UUID,
    case_reference      VARCHAR(100),
    parquet_ref         VARCHAR(100),
    medical_report_ref  VARCHAR(200),
    autopsy_ref         VARCHAR(200),
    dna_sample_ref      VARCHAR(100),
    afis_subject_id     UUID,
    rvin_case_id        UUID,
    iachr_ref           VARCHAR(50),
    un_special_rap_ref  VARCHAR(50),
    needs_reparation    BOOLEAN DEFAULT FALSE,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE vict_mass_incidents (
    mass_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_name       VARCHAR(200) NOT NULL,
    crime_type          vict_crime_type NOT NULL,
    incident_date       TIMESTAMPTZ NOT NULL,
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    victim_count        INTEGER NOT NULL,
    survivor_count      INTEGER DEFAULT 0,
    perpetrator_gang_id UUID,
    description         TEXT,
    documented_by       TEXT[] DEFAULT '{}',
    iachr_case_ref      VARCHAR(50),
    linked_victim_ids   UUID[] DEFAULT '{}',
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vict_crime_type ON vict_victims(crime_type, victim_status);
CREATE INDEX idx_vict_dept ON vict_victims(dept_code, incident_date DESC);
CREATE INDEX idx_vict_gang ON vict_victims(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_vict_mass ON vict_mass_incidents(dept_code, incident_date DESC);
