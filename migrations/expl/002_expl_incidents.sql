BEGIN;

CREATE TABLE expl_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_expl_id    VARCHAR(25) UNIQUE NOT NULL,
    incident_type       expl_incident_type NOT NULL,
    explosive_type      expl_type NOT NULL,
    status              expl_status NOT NULL DEFAULT 'RECOVERED',
    quantity            INTEGER DEFAULT 1,
    weight_kg           DECIMAL(10,3),
    manufacturer        VARCHAR(100),
    lot_number          VARCHAR(50),
    manufacture_country CHAR(3),
    estimated_date      DATE,
    incident_date       TIMESTAMPTZ NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    responding_unit     VARCHAR(50),
    eod_officer         UUID,
    casualties          SMALLINT DEFAULT 0,
    gang_id             UUID,
    from_person_id      UUID,
    case_reference      VARCHAR(100),
    dna_sample_taken    BOOLEAN DEFAULT FALSE,
    bio_sample_ref      VARCHAR(100),
    photo_refs          TEXT[] DEFAULT '{}',
    interpol_exploint_ref VARCHAR(50),
    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
