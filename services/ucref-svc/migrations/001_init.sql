BEGIN;
CREATE TYPE ucref_str_status AS ENUM ('RECEIVED','UNDER_ANALYSIS','DISSEMINATED','ARCHIVED','NO_ACTION');
CREATE TYPE ucref_report_type AS ENUM ('STR','CTR','INTERNATIONAL_WIRE','REAL_ESTATE','MONCASH_PATTERN','CRYPTO_PATTERN');
CREATE TABLE ucref_str_reports (
    str_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), national_str_id VARCHAR(25) UNIQUE NOT NULL,
    report_type ucref_report_type NOT NULL, status ucref_str_status NOT NULL DEFAULT 'RECEIVED',
    reporting_institution VARCHAR(200) NOT NULL, institution_type VARCHAR(30),
    report_date TIMESTAMPTZ NOT NULL, transaction_date TIMESTAMPTZ,
    transaction_amount DECIMAL(18,2), transaction_currency CHAR(3) DEFAULT 'HTG',
    transaction_amount_usd DECIMAL(18,2), subject_snisid_ids UUID[] DEFAULT '{}',
    subject_names TEXT[] DEFAULT '{}', subject_accounts TEXT[] DEFAULT '{}',
    suspicious_activity TEXT NOT NULL, ml_typology VARCHAR(100),
    predicate_crime VARCHAR(100), gang_id UUID, fpr_person_ids UUID[] DEFAULT '{}',
    sanc_match_ids UUID[] DEFAULT '{}', analyst_id UUID, analysis_notes TEXT,
    disseminated_to TEXT[] DEFAULT '{}', disseminated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE ucref_financial_profiles (
    profile_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), snisid_person_id UUID NOT NULL UNIQUE,
    total_str_count INTEGER DEFAULT 0, total_ctr_count INTEGER DEFAULT 0,
    estimated_illegal_assets_usd DECIMAL(18,2), known_accounts JSONB,
    known_properties JSONB, known_businesses TEXT[] DEFAULT '{}',
    ml_risk_score SMALLINT CHECK (ml_risk_score BETWEEN 0 AND 100),
    is_pep BOOLEAN DEFAULT FALSE, last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE ucref_moncash_patterns (
    pattern_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), str_id UUID REFERENCES ucref_str_reports(str_id),
    phone_number VARCHAR(20) NOT NULL, snisid_person_id UUID,
    pattern_type VARCHAR(50), transaction_count INTEGER,
    total_amount_htg DECIMAL(18,2), period_start TIMESTAMPTZ,
    period_end TIMESTAMPTZ, notes TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ucref_str_status ON ucref_str_reports(status, report_date DESC);
CREATE INDEX idx_ucref_str_gang ON ucref_str_reports(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_ucref_str_subjects ON ucref_str_reports USING gin(subject_snisid_ids);
CREATE INDEX idx_ucref_profiles ON ucref_financial_profiles(snisid_person_id);
CREATE INDEX idx_ucref_moncash ON ucref_moncash_patterns(phone_number);
COMMIT;
