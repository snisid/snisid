-- Migration 004 : Tables SSI (Self-Sovereign Identity)

CREATE SCHEMA IF NOT EXISTS snisid_ssi;

CREATE TABLE snisid_ssi.did_records (
    did             VARCHAR(500) NOT NULL PRIMARY KEY,
    method          VARCHAR(50) NOT NULL,
    document        JSONB NOT NULL,
    controller_niu  CHAR(10),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE snisid_ssi.verifiable_credentials (
    credential_id       UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    issuer_did          VARCHAR(500) NOT NULL,
    subject_did         VARCHAR(500) NOT NULL,
    credential_type     VARCHAR(100) NOT NULL,
    document            JSONB NOT NULL,
    status_list_id      UUID,
    is_revoked          BOOLEAN NOT NULL DEFAULT FALSE,
    issued_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at          TIMESTAMPTZ
);

CREATE INDEX idx_vc_subject ON snisid_ssi.verifiable_credentials (subject_did, credential_type);
CREATE INDEX idx_vc_issuer ON snisid_ssi.verifiable_credentials (issuer_did);
CREATE INDEX idx_vc_active ON snisid_ssi.verifiable_credentials (is_revoked, expires_at) WHERE is_revoked = FALSE;

CREATE TABLE snisid_ssi.status_lists (
    list_id         UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    purpose         VARCHAR(50) NOT NULL DEFAULT 'revocation',
    bitstring       TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE snisid_ssi.credential_flows (
    flow_id         UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    flow_type       VARCHAR(50) NOT NULL,
    subject_did     VARCHAR(500),
    issuer_did      VARCHAR(500),
    state           VARCHAR(50) NOT NULL DEFAULT 'INITIATED',
    payload         JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ
);

CREATE TABLE snisid_ssi.revocation_events (
    event_id        UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    credential_id   UUID NOT NULL REFERENCES snisid_ssi.verifiable_credentials(credential_id),
    reason          TEXT,
    revoked_by_did  VARCHAR(500),
    revoked_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
