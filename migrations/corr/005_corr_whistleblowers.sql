BEGIN;

CREATE TABLE corr_whistleblower_reports (
    report_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_token        VARCHAR(64) UNIQUE NOT NULL,
    allegation_type     corr_allegation_type NOT NULL,
    severity_estimate   corr_severity,
    officer_unit_hint   VARCHAR(50),
    officer_rank_hint   VARCHAR(50),
    description         TEXT NOT NULL,
    evidence_description TEXT,
    submission_date     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_hash             VARCHAR(64),
    processed           BOOLEAN DEFAULT FALSE,
    processed_by        UUID,
    integrity_case_id   UUID REFERENCES corr_integrity_cases(case_id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
