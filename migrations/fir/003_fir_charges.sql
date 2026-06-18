BEGIN;

CREATE TABLE fir_charges (
    charge_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id              UUID NOT NULL REFERENCES fir_criminal_records(record_id),
    is_arrest              BOOLEAN NOT NULL DEFAULT TRUE,
    arrest_date            TIMESTAMPTZ,
    arresting_unit         VARCHAR(50),
    arresting_officer      UUID,
    arrest_location        VARCHAR(300),
    dept_code              CHAR(2),
    charges_text           TEXT,
    offense_class          fir_offense_class NOT NULL,
    case_reference         VARCHAR(100),
    release_date           TIMESTAMPTZ,
    release_reason         TEXT,
    court_name             VARCHAR(150),
    court_dept             CHAR(2),
    offense_description    TEXT,
    ipc_code               VARCHAR(30),
    verdict_date           TIMESTAMPTZ,
    case_status            fir_case_status NOT NULL DEFAULT 'OPEN',
    sentence_type          fir_sentence_type,
    sentence_duration_days INTEGER,
    fine_amount_gdes       DECIMAL(12,2),
    sentence_start         TIMESTAMPTZ,
    sentence_end           TIMESTAMPTZ,
    is_foreign_record      BOOLEAN DEFAULT FALSE,
    foreign_country        CHAR(3),
    interpol_ccc_ref       VARCHAR(50),
    judge_name             VARCHAR(150),
    notes                  TEXT,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
