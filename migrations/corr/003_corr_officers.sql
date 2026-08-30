BEGIN;

CREATE TABLE corr_officers (
    officer_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snisid_id           UUID NOT NULL UNIQUE,
    badge               VARCHAR(30),
    full_name           VARCHAR(150) NOT NULL,
    unit                VARCHAR(50),
    rank                VARCHAR(50),

    under_investigation BOOLEAN DEFAULT FALSE,
    active_case_id      UUID REFERENCES corr_integrity_cases(case_id),
    investigation_count INTEGER DEFAULT 0,
    total_risk_score    SMALLINT DEFAULT 0,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_corr_officers_snisid ON corr_officers(snisid_id);
CREATE INDEX idx_corr_officers_unit   ON corr_officers(unit);
CREATE INDEX idx_corr_officers_active ON corr_officers(under_investigation) WHERE under_investigation = TRUE;

COMMIT;
