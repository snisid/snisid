-- SIPEP-HT Migration: Système d'Information Pénitentiaire d'Haïti

BEGIN;

CREATE TYPE sipep_detention_basis AS ENUM (
    'PREVENTIVE',        -- Détention provisoire (avant jugement)
    'SENTENCED',         -- Condamné purgeant peine
    'ADMINISTRATIVE',    -- Détention administrative (immigration, etc.)
    'CONTEMPT'           -- Outrage au tribunal
);

CREATE TYPE sipep_legal_status AS ENUM (
    'AWAITING_TRIAL',
    'ON_TRIAL',
    'SENTENCED',
    'APPEAL_PENDING',
    'CONDEMNED'
);

CREATE TYPE sipep_release_type AS ENUM (
    'SENTENCE_SERVED',
    'CONDITIONAL_RELEASE',
    'BAIL',
    'JUDICIAL_ORDER',
    'DEATH',
    'ESCAPE',
    'TRANSFER_OUT'
);

CREATE TABLE sipep_inmates (
    inmate_id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_inmate_id      VARCHAR(25) UNIQUE NOT NULL, -- Format: SIPEP-HT-AAAA-NNNNNN
    snisid_person_id        UUID NOT NULL,
    fir_record_id           UUID,
    afis_subject_id         UUID,
    current_facility        VARCHAR(100) NOT NULL,
    current_dept_code       CHAR(2),
    cell_block              VARCHAR(20),
    is_currently_detained   BOOLEAN DEFAULT TRUE,
    is_minor                BOOLEAN DEFAULT FALSE,
    is_female               BOOLEAN DEFAULT FALSE,
    has_special_needs       BOOLEAN DEFAULT FALSE,
    special_needs_notes     TEXT,
    intake_date             TIMESTAMPTZ NOT NULL,
    expected_release_date   TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipep_detentions (
    detention_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inmate_id               UUID NOT NULL REFERENCES sipep_inmates(inmate_id),
    facility                VARCHAR(100) NOT NULL,
    detention_basis         sipep_detention_basis NOT NULL,
    legal_status            sipep_legal_status NOT NULL DEFAULT 'AWAITING_TRIAL',
    case_reference          VARCHAR(100),
    court_name              VARCHAR(150),
    arresting_authority     VARCHAR(100),
    warrant_number          VARCHAR(100),
    intake_date             TIMESTAMPTZ NOT NULL,
    intake_officer          UUID NOT NULL,
    sentence_duration_days  INTEGER,
    release_date            TIMESTAMPTZ,
    release_type            sipep_release_type,
    releasing_authority     VARCHAR(100),
    notes                   TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipep_transfers (
    transfer_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inmate_id               UUID NOT NULL REFERENCES sipep_inmates(inmate_id),
    from_facility           VARCHAR(100) NOT NULL,
    to_facility             VARCHAR(100) NOT NULL,
    transfer_date           TIMESTAMPTZ NOT NULL,
    transfer_reason         TEXT,
    authorized_by           UUID NOT NULL,
    transport_unit          VARCHAR(50),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipep_health_events (
    event_id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inmate_id               UUID NOT NULL REFERENCES sipep_inmates(inmate_id),
    event_type              VARCHAR(50) NOT NULL, -- INJURY, ILLNESS, DEATH, PSYCHIATRIC
    event_date              TIMESTAMPTZ NOT NULL,
    description             TEXT,
    treating_facility       VARCHAR(150),
    outcome                 TEXT,
    reported_by             UUID,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sipep_inmates_facility ON sipep_inmates(current_facility) WHERE is_currently_detained = TRUE;
CREATE INDEX idx_sipep_inmates_person   ON sipep_inmates(snisid_person_id);
CREATE INDEX idx_sipep_detentions_case  ON sipep_detentions(case_reference);
CREATE INDEX idx_sipep_detentions_basis ON sipep_detentions(detention_basis, legal_status);

COMMIT;
