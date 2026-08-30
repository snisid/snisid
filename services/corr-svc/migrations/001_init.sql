CREATE TYPE corr_allegation_type AS ENUM (
    'DATA_LEAK_TO_GANG', 'RECORD_TAMPERING', 'UNAUTHORIZED_ACCESS',
    'BRIBERY', 'EXTORTION_OF_CIVILIANS', 'FACILITATED_PRISON_ESCAPE',
    'STOLEN_CREDENTIALS', 'FINANCIAL_CORRUPTION', 'GANG_AFFILIATION', 'OTHER'
);

CREATE TYPE corr_severity AS ENUM ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL');

CREATE TYPE corr_status AS ENUM (
    'REPORTED', 'UNDER_INVESTIGATION', 'SUBSTANTIATED',
    'UNSUBSTANTIATED', 'REFERRED_TO_PARQUET', 'CLOSED'
);

CREATE TABLE corr_integrity_cases (
    case_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_corr_id    VARCHAR(25) UNIQUE NOT NULL,
    officer_snisid_id   UUID NOT NULL,
    officer_badge       VARCHAR(30),
    officer_unit        VARCHAR(50),
    officer_rank        VARCHAR(50),
    allegation_type     corr_allegation_type NOT NULL,
    severity            corr_severity NOT NULL,
    status              corr_status NOT NULL DEFAULT 'REPORTED',
    allegation_summary  TEXT NOT NULL,
    incident_date_from  TIMESTAMPTZ,
    incident_date_to    TIMESTAMPTZ,
    evidence_refs       TEXT[] DEFAULT '{}',
    gang_id             UUID,
    gang_member_ids     UUID[] DEFAULT '{}',
    financial_gain_usd  DECIMAL(15,2),
    blan_case_id        UUID,
    reported_by_type    VARCHAR(30),
    reported_by_id      UUID,
    reporting_date      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_whistleblower    BOOLEAN DEFAULT FALSE,
    whistleblower_protected BOOLEAN DEFAULT FALSE,
    igpnh_investigator  UUID,
    investigation_start TIMESTAMPTZ,
    investigation_end   TIMESTAMPTZ,
    investigation_notes TEXT,
    sanctions_applied   TEXT,
    referred_to_parquet BOOLEAN DEFAULT FALSE,
    parquet_ref         VARCHAR(100),
    ulcc_ref            VARCHAR(100),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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
                            real_estate_usd + vehicles_usd + bank_accounts_usd + other_assets_usd
                        ) STORED,
    known_salary_annual_usd DECIMAL(12,2),
    unexplained_wealth_usd DECIMAL(15,2),
    is_flagged          BOOLEAN DEFAULT FALSE,
    verified_by         UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_corr_cases_officer ON corr_integrity_cases(officer_snisid_id);
CREATE INDEX idx_corr_cases_status ON corr_integrity_cases(status, severity);
CREATE INDEX idx_corr_alerts_officer ON corr_behavioral_alerts(officer_snisid_id, created_at DESC);
CREATE INDEX idx_corr_alerts_unrev ON corr_behavioral_alerts(reviewed) WHERE reviewed = FALSE;
