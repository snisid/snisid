BEGIN;

CREATE TABLE rdep_extraditions (
    extradition_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_rdep_id      VARCHAR(25) UNIQUE NOT NULL,   -- RDEP-EXT-AAAA-NNNNNN
    snisid_person_id      UUID NOT NULL,
    fir_record_id         UUID,

    requesting_country    rdep_deportation_country NOT NULL,
    extradition_status    rdep_extradition_status NOT NULL DEFAULT 'REQUESTED',
    request_date          TIMESTAMPTZ NOT NULL,
    approval_date         TIMESTAMPTZ,
    execution_date        TIMESTAMPTZ,

    charges_summary       TEXT NOT NULL,
    legal_reference       VARCHAR(200),
    treaty_article        VARCHAR(100),

    departure_port        VARCHAR(100),
    departure_dept_code   CHAR(2),
    escorting_agency      VARCHAR(100),
    extradition_officer   UUID,

    notes                 TEXT,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
