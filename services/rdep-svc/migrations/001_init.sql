-- RDEP-HT Migration: Registre des Déportés et Extradés d'Haïti

BEGIN;

CREATE TYPE rdep_deportation_country AS ENUM (
    'USA','CAN','DOM','BHS','CUB','JAM','TTO','MEX','BRA','FRA','OTHER'
);

CREATE TYPE rdep_criminal_risk AS ENUM (
    'NONE','LOW','MEDIUM','HIGH','VERY_HIGH'
);

CREATE TYPE rdep_monitoring_status AS ENUM (
    'ACTIVE','SUSPENDED','COMPLETED','FLED','DECEASED'
);

CREATE TABLE rdep_deportees (
    deportee_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_rdep_id    VARCHAR(25) UNIQUE NOT NULL,   -- Format: RDEP-HT-AAAA-NNNNNN
    snisid_person_id    UUID NOT NULL,
    fir_record_id       UUID,
    afis_subject_id     UUID,

    -- Informations déportation
    deportation_country rdep_deportation_country NOT NULL,
    deportation_date    TIMESTAMPTZ NOT NULL,
    arrival_port        VARCHAR(100) NOT NULL,          -- PAP, CAP, Frontière Malpasse, etc.
    arrival_dept_code   CHAR(2),
    deporting_agency    VARCHAR(100),                   -- ICE, CBSA, DNCD-DOM, etc.
    deportation_reason  TEXT,
    flight_number       VARCHAR(20),

    -- Identité étrangère
    foreign_name        VARCHAR(200),
    foreign_aliases     TEXT[] DEFAULT '{}',
    foreign_id_number   VARCHAR(100),                   -- SSN/SIN masqué, etc.
    foreign_country_id  VARCHAR(50),

    -- Antécédents criminels étrangers
    has_foreign_record  BOOLEAN DEFAULT FALSE,
    criminal_risk_level rdep_criminal_risk DEFAULT 'NONE',
    convicted_offenses  TEXT[] DEFAULT '{}',
    gang_affiliated     BOOLEAN DEFAULT FALSE,
    gang_name           VARCHAR(100),

    -- Surveillance
    monitoring_required BOOLEAN DEFAULT FALSE,
    monitoring_status   rdep_monitoring_status DEFAULT 'ACTIVE',
    monitoring_unit     VARCHAR(50),
    monitoring_officer  UUID,
    monitoring_end_date TIMESTAMPTZ,

    -- Localisation actuelle
    current_address     TEXT,
    current_commune     VARCHAR(100),
    current_dept_code   CHAR(2),

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rdep_foreign_records (
    foreign_record_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deportee_id         UUID NOT NULL REFERENCES rdep_deportees(deportee_id),
    country             rdep_deportation_country NOT NULL,
    court_name          VARCHAR(200),
    offense_description TEXT NOT NULL,
    offense_date        TIMESTAMPTZ,
    conviction_date     TIMESTAMPTZ,
    sentence            TEXT,
    prison_served       TEXT,
    fbi_number          VARCHAR(50),
    interpol_ref        VARCHAR(50),
    source_document     VARCHAR(500),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rdep_monitoring_events (
    event_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deportee_id         UUID NOT NULL REFERENCES rdep_deportees(deportee_id),
    event_type          VARCHAR(50) NOT NULL,  -- CHECK_IN, VIOLATION, ADDRESS_CHANGE
    event_date          TIMESTAMPTZ NOT NULL,
    location_lat        DECIMAL(10,7),
    location_lng        DECIMAL(10,7),
    notes               TEXT,
    reported_by         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rdep_deportees_person      ON rdep_deportees(snisid_person_id);
CREATE INDEX idx_rdep_deportees_country     ON rdep_deportees(deportation_country);
CREATE INDEX idx_rdep_deportees_risk        ON rdep_deportees(criminal_risk_level) WHERE criminal_risk_level IN ('HIGH','VERY_HIGH');
CREATE INDEX idx_rdep_deportees_gang        ON rdep_deportees(gang_affiliated) WHERE gang_affiliated = TRUE;
CREATE INDEX idx_rdep_deportees_monitoring  ON rdep_deportees(monitoring_status) WHERE monitoring_required = TRUE;

COMMIT;
