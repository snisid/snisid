-- Migration 001 : Schéma identité core

CREATE SCHEMA IF NOT EXISTS snisid_identity;
CREATE SCHEMA IF NOT EXISTS snisid_civil;
CREATE SCHEMA IF NOT EXISTS snisid_audit;
CREATE SCHEMA IF NOT EXISTS snisid_biometric;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TYPE snisid_identity.identity_status AS ENUM (
    'PRE_REGISTERED',
    'PENDING_VERIFICATION',
    'ACTIVE',
    'SUSPENDED',
    'DECEASED',
    'ARCHIVED'
);

CREATE TYPE snisid_identity.department_code AS ENUM (
    'OUEST', 'NORD', 'NORD_EST', 'NORD_OUEST',
    'ARTIBONITE', 'CENTRE', 'SUD', 'SUD_EST',
    'GRAND_ANSE', 'NIPPES'
);

CREATE TABLE snisid_identity.citizens (
    niu                    CHAR(10) NOT NULL,
    niu_version            INTEGER NOT NULL DEFAULT 1,
    statut_identite        snisid_identity.identity_status NOT NULL DEFAULT 'PRE_REGISTERED',
    prenom                 VARCHAR(100) NOT NULL,
    nom                    VARCHAR(100) NOT NULL,
    prenom_creole          VARCHAR(100),
    nom_creole             VARCHAR(100),
    date_naissance         DATE NOT NULL,
    lieu_naissance         VARCHAR(200),
    genre                  CHAR(1) CHECK (genre IN ('M', 'F', 'X')),
    nationalite            CHAR(3) NOT NULL DEFAULT 'HTI',
    pere_niu               CHAR(10),
    mere_niu               CHAR(10),
    departement_residence  snisid_identity.department_code NOT NULL,
    commune_residence      VARCHAR(100),
    adresse_complete       TEXT,
    agence_enregistrement  VARCHAR(50),
    agent_niu              CHAR(10),
    date_enregistrement    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    cert_thumbprint        CHAR(64),
    fraud_score_initial    DECIMAL(5,4) CHECK (fraud_score_initial BETWEEN 0 AND 1),
    fraud_flags            JSONB DEFAULT '{}',
    risk_level             VARCHAR(10) DEFAULT 'LOW' CHECK (risk_level IN ('LOW','MEDIUM','HIGH','CRITICAL')),
    version                INTEGER NOT NULL DEFAULT 1,
    data_hash              CHAR(64),
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_citizens PRIMARY KEY (niu, departement_residence)
) PARTITION BY LIST (departement_residence);

CREATE TABLE snisid_identity.citizens_ouest     PARTITION OF snisid_identity.citizens FOR VALUES IN ('OUEST');
CREATE TABLE snisid_identity.citizens_nord      PARTITION OF snisid_identity.citizens FOR VALUES IN ('NORD');
CREATE TABLE snisid_identity.citizens_nord_est  PARTITION OF snisid_identity.citizens FOR VALUES IN ('NORD_EST');
CREATE TABLE snisid_identity.citizens_nord_ouest PARTITION OF snisid_identity.citizens FOR VALUES IN ('NORD_OUEST');
CREATE TABLE snisid_identity.citizens_artibonite PARTITION OF snisid_identity.citizens FOR VALUES IN ('ARTIBONITE');
CREATE TABLE snisid_identity.citizens_centre    PARTITION OF snisid_identity.citizens FOR VALUES IN ('CENTRE');
CREATE TABLE snisid_identity.citizens_sud       PARTITION OF snisid_identity.citizens FOR VALUES IN ('SUD');
CREATE TABLE snisid_identity.citizens_sud_est   PARTITION OF snisid_identity.citizens FOR VALUES IN ('SUD_EST');
CREATE TABLE snisid_identity.citizens_grand_anse PARTITION OF snisid_identity.citizens FOR VALUES IN ('GRAND_ANSE');
CREATE TABLE snisid_identity.citizens_nippes    PARTITION OF snisid_identity.citizens FOR VALUES IN ('NIPPES');

CREATE INDEX idx_citizens_niu ON snisid_identity.citizens (niu);
CREATE INDEX idx_citizens_nom_prenom ON snisid_identity.citizens USING gin(
    (nom || ' ' || prenom) gin_trgm_ops
);
CREATE INDEX idx_citizens_date_naissance ON snisid_identity.citizens (date_naissance);
CREATE INDEX idx_citizens_fraud_score ON snisid_identity.citizens (fraud_score_initial) WHERE fraud_score_initial > 0.5;

CREATE OR REPLACE FUNCTION snisid_identity.update_updated_at()
RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = NOW(); RETURN NEW; END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_citizens_updated_at
    BEFORE UPDATE ON snisid_identity.citizens
    FOR EACH ROW EXECUTE FUNCTION snisid_identity.update_updated_at();

CREATE OR REPLACE FUNCTION snisid_identity.increment_version()
RETURNS TRIGGER AS $$
BEGIN NEW.version = OLD.version + 1; RETURN NEW; END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_citizens_version
    BEFORE UPDATE ON snisid_identity.citizens
    FOR EACH ROW EXECUTE FUNCTION snisid_identity.increment_version();

ALTER TABLE snisid_identity.citizens ENABLE ROW LEVEL SECURITY;

CREATE POLICY policy_oni_full ON snisid_identity.citizens
    FOR ALL TO oni_role
    USING (true);

CREATE POLICY policy_pnh_read ON snisid_identity.citizens
    FOR SELECT TO pnh_role
    USING (true);

CREATE TABLE snisid_identity.identity_events (
    event_id        UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    niu             CHAR(10) NOT NULL,
    event_type      VARCHAR(100) NOT NULL,
    event_version   INTEGER NOT NULL,
    payload         JSONB NOT NULL,
    agent_niu       CHAR(10),
    agence          VARCHAR(50),
    trace_id        VARCHAR(64),
    span_id         VARCHAR(32),
    previous_hash   CHAR(64),
    event_hash      CHAR(64),
    signature       TEXT,
    kafka_topic     VARCHAR(200),
    kafka_offset    BIGINT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (niu, event_version)
);

CREATE RULE no_update_events AS ON UPDATE TO snisid_identity.identity_events DO INSTEAD NOTHING;
CREATE RULE no_delete_events AS ON DELETE TO snisid_identity.identity_events DO INSTEAD NOTHING;

CREATE INDEX idx_events_niu ON snisid_identity.identity_events (niu, event_version);
CREATE INDEX idx_events_type ON snisid_identity.identity_events (event_type, created_at);

CREATE TABLE snisid_identity.identity_snapshots (
    niu             CHAR(10) NOT NULL PRIMARY KEY,
    latest_version  INTEGER NOT NULL,
    state           JSONB NOT NULL,
    state_hash      CHAR(64),
    snapshot_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
