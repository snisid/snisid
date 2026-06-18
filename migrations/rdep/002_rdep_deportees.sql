BEGIN;

CREATE TABLE rdep_deportees (
    deportee_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_rdep_id      VARCHAR(25) UNIQUE NOT NULL,   -- RDEP-HT-AAAA-NNNNNN
    snisid_person_id      UUID NOT NULL,
    fir_record_id         UUID,
    afis_subject_id       UUID,

    deportation_country   rdep_deportation_country NOT NULL,
    deportation_date      TIMESTAMPTZ NOT NULL,
    arrival_port          VARCHAR(100) NOT NULL,
    arrival_dept_code     CHAR(2),
    deporting_agency      VARCHAR(100),
    deportation_reason    TEXT,
    flight_id             UUID,
    flight_number         VARCHAR(20),

    foreign_name          VARCHAR(200),
    foreign_aliases       TEXT[] DEFAULT '{}',
    foreign_id_number     VARCHAR(100),
    foreign_country_id    VARCHAR(50),

    has_foreign_record    BOOLEAN DEFAULT FALSE,
    criminal_risk_level   rdep_criminal_risk DEFAULT 'NONE',
    convicted_offenses    TEXT[] DEFAULT '{}',
    gang_affiliated       BOOLEAN DEFAULT FALSE,
    gang_name             VARCHAR(100),

    monitoring_required   BOOLEAN DEFAULT FALSE,
    monitoring_status     rdep_monitoring_status DEFAULT 'ACTIVE',
    monitoring_unit       VARCHAR(50),
    monitoring_officer    UUID,
    monitoring_end_date   TIMESTAMPTZ,

    current_address       TEXT,
    current_commune       VARCHAR(100),
    current_dept_code     CHAR(2),

    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
