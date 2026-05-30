-- ============================================================
-- SNISID-Security — Criminal Case Schema (PostgreSQL/CockroachDB)
-- Schema Type: Event Sourcing & CQRS Projections
-- ============================================================

CREATE SCHEMA IF NOT EXISTS justice;
SET search_path TO justice;

-- ---------------------------------------------------------
-- 1. EVENT STORE (Immutable Append-Only Log)
-- ---------------------------------------------------------
CREATE TABLE criminal_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id VARCHAR(50) NOT NULL,
    event_type VARCHAR(100) NOT NULL, -- ex: 'CaseOpened', 'SuspectAdded', 'WarrantIssued'
    aggregate_version INT NOT NULL,
    payload JSONB NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    agent_id VARCHAR(100) NOT NULL, -- PKI de l'agent
    agency_id VARCHAR(50) NOT NULL, -- PNH, DCPJ, PARQUET_PAP
    signature TEXT NOT NULL, -- Signature cryptographique XAdES/ECDSA du payload
    UNIQUE(case_id, aggregate_version)
);

CREATE INDEX idx_criminal_events_case ON criminal_events(case_id);
CREATE INDEX idx_criminal_events_type ON criminal_events(event_type);

-- ---------------------------------------------------------
-- 2. CQRS READ PROJECTIONS (Materialized Views for querying)
-- ---------------------------------------------------------

CREATE TABLE cases (
    case_id VARCHAR(50) PRIMARY KEY, -- DOS-{YEAR}-{JURIDICTION}-{SEQ}
    status VARCHAR(50) NOT NULL,     -- OPENED, INVESTIGATION, TRIAL, CLOSED...
    jurisdiction_code VARCHAR(10) NOT NULL, -- PAP, CAP, CAY, etc.
    crime_category VARCHAR(100) NOT NULL, -- HOMICIDE, FRAUD, KIDNAPPING...
    incident_date TIMESTAMPTZ,
    location_commune VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    assigned_agency VARCHAR(50),
    assigned_investigator VARCHAR(100)
);

CREATE TABLE case_entities (
    entity_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id VARCHAR(50) REFERENCES cases(case_id),
    niu VARCHAR(10), -- Optionnel si l'entité n'a pas (encore) de NIU
    entity_role VARCHAR(50) NOT NULL, -- SUSPECT, VICTIM, WITNESS, ACCOMPLICE
    temporary_alias VARCHAR(200), -- Si NIU inconnu (ex: 'John Doe 1')
    added_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_case_entities_niu ON case_entities(niu);

CREATE TABLE warrants (
    warrant_id VARCHAR(50) PRIMARY KEY, -- MAN-{YEAR}-{SEQ}
    case_id VARCHAR(50) REFERENCES cases(case_id),
    target_niu VARCHAR(10) NOT NULL, -- Lié formellement à un NIU civil
    warrant_type VARCHAR(50) NOT NULL, -- AMENER, ARRET, PERQUISITION
    issuing_authority VARCHAR(100) NOT NULL, -- Juge d'instruction ou parquet
    issued_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL, -- ACTIVE, EXECUTED, CANCELLED, EXPIRED
    executed_at TIMESTAMPTZ,
    executed_by VARCHAR(100)
);

CREATE INDEX idx_warrants_target ON warrants(target_niu) WHERE status = 'ACTIVE';

CREATE TABLE evidence (
    evidence_id VARCHAR(50) PRIMARY KEY, -- PRV-{YEAR}-{SEQ}
    case_id VARCHAR(50) REFERENCES cases(case_id),
    evidence_type VARCHAR(50) NOT NULL, -- PHYSICAL, DIGITAL_VIDEO, DIGITAL_AUDIO, DOCUMENT
    description TEXT,
    sha256_hash VARCHAR(64), -- Pour preuves numériques
    storage_location VARCHAR(200) NOT NULL, -- MinIO bucket URI ou salle des scellés
    collected_by VARCHAR(100) NOT NULL,
    collected_at TIMESTAMPTZ NOT NULL,
    current_custody VARCHAR(100) NOT NULL -- PNH_EVIDENCE_ROOM, TRIBUNAL, LABO
);

CREATE TABLE chain_of_custody (
    custody_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evidence_id VARCHAR(50) REFERENCES evidence(evidence_id),
    transferred_from VARCHAR(100) NOT NULL,
    transferred_to VARCHAR(100) NOT NULL,
    transfer_date TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    transfer_reason TEXT,
    receiver_signature TEXT NOT NULL -- Validation PKI de la réception
);

-- ---------------------------------------------------------
-- 3. RLS (Row Level Security) - ABAC Policies
-- ---------------------------------------------------------

ALTER TABLE cases ENABLE ROW LEVEL SECURITY;
ALTER TABLE warrants ENABLE ROW LEVEL SECURITY;

-- PNH ne voit que les dossiers ouverts dans sa juridiction, sauf DCPJ qui voit tout
CREATE POLICY pnh_jurisdiction_policy ON cases
    FOR SELECT
    USING (
        current_setting('snisid.agency') = 'DCPJ' 
        OR (current_setting('snisid.agency') = 'PNH' AND jurisdiction_code = current_setting('snisid.jurisdiction'))
    );

-- Les mandats ACTIFS sont publics pour toutes les agences de sécurité (police, frontières)
CREATE POLICY active_warrants_visible_to_all_security ON warrants
    FOR SELECT
    USING (
        status = 'ACTIVE' 
        AND current_setting('snisid.agency') IN ('PNH', 'DCPJ', 'BORDER', 'PRISON', 'JUSTICE')
    );
