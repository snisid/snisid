BEGIN;

CREATE TABLE corr_behavioral_alerts (
    alert_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    officer_snisid_id   UUID NOT NULL,
    alert_type          VARCHAR(50) NOT NULL,
    description         TEXT NOT NULL,
    module_source       VARCHAR(30),
    risk_score          SMALLINT CHECK (risk_score BETWEEN 0 AND 100),
    auto_generated      BOOLEAN DEFAULT TRUE,
    reviewed            BOOLEAN DEFAULT FALSE,
    reviewed_by         UUID,
    is_false_positive   BOOLEAN DEFAULT FALSE,
    corr_case_id        UUID REFERENCES corr_integrity_cases(case_id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE corr_asset_declarations (
    declaration_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    officer_snisid_id   UUID NOT NULL,
    declaration_year    SMALLINT NOT NULL,
    real_estate_usd     DECIMAL(15,2) DEFAULT 0,
    vehicles_usd        DECIMAL(15,2) DEFAULT 0,
    bank_accounts_usd   DECIMAL(15,2) DEFAULT 0,
    other_assets_usd    DECIMAL(15,2) DEFAULT 0,
    total_assets_usd    DECIMAL(15,2) GENERATED ALWAYS AS (
                            real_estate_usd + vehicles_usd +
                            bank_accounts_usd + other_assets_usd
                        ) STORED,
    known_salary_annual_usd DECIMAL(12,2),
    unexplained_wealth_usd DECIMAL(15,2),
    is_flagged          BOOLEAN DEFAULT FALSE,
    verified_by         UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_corr_cases_officer   ON corr_integrity_cases(officer_snisid_id);
CREATE INDEX idx_corr_cases_status    ON corr_integrity_cases(status, severity);
CREATE INDEX idx_corr_alerts_officer  ON corr_behavioral_alerts(officer_snisid_id, created_at DESC);
CREATE INDEX idx_corr_alerts_unrev    ON corr_behavioral_alerts(reviewed) WHERE reviewed = FALSE;
CREATE INDEX idx_corr_decl_officer    ON corr_asset_declarations(officer_snisid_id);
CREATE INDEX idx_corr_decl_flagged    ON corr_asset_declarations(is_flagged) WHERE is_flagged = TRUE;

-- Row-Level Security
ALTER TABLE corr_integrity_cases      ENABLE ROW LEVEL SECURITY;
ALTER TABLE corr_behavioral_alerts    ENABLE ROW LEVEL SECURITY;
ALTER TABLE corr_asset_declarations   ENABLE ROW LEVEL SECURITY;
ALTER TABLE corr_officers             ENABLE ROW LEVEL SECURITY;
ALTER TABLE corr_evidence             ENABLE ROW LEVEL SECURITY;
ALTER TABLE corr_whistleblower_reports ENABLE ROW LEVEL SECURITY;

CREATE POLICY corr_strict_access ON corr_integrity_cases
    FOR ALL USING (
        current_setting('app.user_role', TRUE) IN ('IGPNH','MJSP_MINISTER','SUPERADMIN')
    );

CREATE POLICY corr_alerts_access ON corr_behavioral_alerts
    FOR ALL USING (
        current_setting('app.user_role', TRUE) IN ('IGPNH','MJSP_MINISTER','SUPERADMIN')
    );

CREATE POLICY corr_decl_access ON corr_asset_declarations
    FOR ALL USING (
        current_setting('app.user_role', TRUE) IN ('IGPNH','ULCC','SUPERADMIN')
    );

CREATE POLICY corr_officers_access ON corr_officers
    FOR ALL USING (
        current_setting('app.user_role', TRUE) IN ('IGPNH','MJSP_MINISTER','SUPERADMIN')
    );

CREATE POLICY corr_evidence_access ON corr_evidence
    FOR ALL USING (
        current_setting('app.user_role', TRUE) IN ('IGPNH','SUPERADMIN')
    );

CREATE POLICY corr_whistleblower_select ON corr_whistleblower_reports
    FOR SELECT USING (
        current_setting('app.user_role', TRUE) IN ('IGPNH','SUPERADMIN')
    );

COMMIT;
