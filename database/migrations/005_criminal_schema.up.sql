-- Migration 005 : Schéma criminel

CREATE SCHEMA IF NOT EXISTS snisid_criminal;

CREATE TYPE snisid_criminal.case_status AS ENUM (
    'OPEN',
    'UNDER_INVESTIGATION',
    'SUSPENDED',
    'CLOSED_SOLVED',
    'CLOSED_UNSOLVED'
);

CREATE TYPE snisid_criminal.evidence_type AS ENUM (
    'DNA',
    'FINGERPRINT',
    'WEAPON',
    'DOCUMENT',
    'DIGITAL',
    'PHOTO',
    'VIDEO',
    'WITNESS_TESTIMONY',
    'OTHER'
);

CREATE TABLE snisid_criminal.cases (
    case_id             UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    case_number         VARCHAR(50) UNIQUE NOT NULL,
    title               VARCHAR(300) NOT NULL,
    description         TEXT,
    status              snisid_criminal.case_status NOT NULL DEFAULT 'OPEN',
    severity            SMALLINT CHECK (severity BETWEEN 1 AND 5),
    primary_dept        VARCHAR(10),
    lead_investigator   CHAR(10),
    victim_count        INTEGER DEFAULT 0,
    suspect_count       INTEGER DEFAULT 0,
    opened_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at           TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cases_status ON snisid_criminal.cases (status);
CREATE INDEX idx_cases_dept ON snisid_criminal.cases (primary_dept);
CREATE INDEX idx_cases_number ON snisid_criminal.cases (case_number);

CREATE TABLE snisid_criminal.case_persons (
    id                  UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    case_id             UUID NOT NULL REFERENCES snisid_criminal.cases(case_id) ON DELETE CASCADE,
    niu                 CHAR(10) NOT NULL,
    role                VARCHAR(50) NOT NULL CHECK (role IN ('VICTIM', 'SUSPECT', 'WITNESS', 'INFORMANT')),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_case_persons_niu ON snisid_criminal.case_persons (niu);
CREATE INDEX idx_case_persons_case ON snisid_criminal.case_persons (case_id);

CREATE TABLE snisid_criminal.evidence (
    evidence_id         UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    case_id             UUID NOT NULL REFERENCES snisid_criminal.cases(case_id) ON DELETE CASCADE,
    evidence_type       snisid_criminal.evidence_type NOT NULL,
    description         TEXT,
    collection_date     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    collector_niu       CHAR(10),
    storage_location    VARCHAR(300),
    chain_of_custody    JSONB DEFAULT '[]',
    dna_profile_id      UUID,
    fingerprint_id      UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_evidence_case ON snisid_criminal.evidence (case_id);
CREATE INDEX idx_evidence_type ON snisid_criminal.evidence (evidence_type);

CREATE TABLE snisid_criminal.warrants (
    warrant_id          UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    warrant_number      VARCHAR(50) UNIQUE NOT NULL,
    case_id             UUID REFERENCES snisid_criminal.cases(case_id),
    target_niu          CHAR(10) NOT NULL,
    warrant_type        VARCHAR(50) NOT NULL,
    issued_by           VARCHAR(200) NOT NULL,
    issued_date         DATE NOT NULL,
    expiry_date         DATE,
    status              VARCHAR(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','EXECUTED','EXPIRED','REVOKED')),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_warrants_niu ON snisid_criminal.warrants (target_niu);
CREATE INDEX idx_warrants_status ON snisid_criminal.warrants (status);
