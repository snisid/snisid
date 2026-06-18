BEGIN;

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

COMMIT;
