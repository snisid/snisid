BEGIN;

CREATE TABLE corr_evidence (
    evidence_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id             UUID NOT NULL REFERENCES corr_integrity_cases(case_id) ON DELETE CASCADE,
    evidence_type       VARCHAR(50) NOT NULL,
    description         TEXT NOT NULL,
    file_hash           VARCHAR(128),
    storage_ref         VARCHAR(255),
    collected_by        UUID,
    collected_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    chain_of_custody    TEXT[] DEFAULT '{}',
    is_verified         BOOLEAN DEFAULT FALSE,
    verified_by         UUID,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_corr_evidence_case ON corr_evidence(case_id);

COMMIT;
