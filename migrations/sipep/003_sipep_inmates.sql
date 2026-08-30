BEGIN;

CREATE TABLE sipep_inmates (
    inmate_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_inmate_id    VARCHAR(25) UNIQUE NOT NULL,
    snisid_person_id      UUID NOT NULL,
    fir_record_id         UUID,
    afis_subject_id       UUID,
    current_facility      VARCHAR(100) NOT NULL,
    current_dept_code     CHAR(2),
    cell_block            VARCHAR(20),
    is_currently_detained BOOLEAN DEFAULT TRUE,
    is_minor              BOOLEAN DEFAULT FALSE,
    is_female             BOOLEAN DEFAULT FALSE,
    has_special_needs     BOOLEAN DEFAULT FALSE,
    special_needs_notes   TEXT,
    intake_date           TIMESTAMPTZ NOT NULL,
    expected_release_date TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipep_detentions (
    detention_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inmate_id              UUID NOT NULL REFERENCES sipep_inmates(inmate_id),
    facility               VARCHAR(100) NOT NULL,
    detention_basis        detention_type NOT NULL,
    legal_status           inmate_status NOT NULL DEFAULT 'AWAITING_TRIAL',
    case_reference         VARCHAR(100),
    court_name             VARCHAR(150),
    arresting_authority    VARCHAR(100),
    warrant_number         VARCHAR(100),
    intake_date            TIMESTAMPTZ NOT NULL,
    intake_officer         UUID NOT NULL,
    sentence_duration_days INTEGER,
    release_date           TIMESTAMPTZ,
    release_type           VARCHAR(30),
    releasing_authority    VARCHAR(100),
    notes                  TEXT,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
