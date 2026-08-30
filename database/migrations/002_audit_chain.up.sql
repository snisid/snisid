-- Migration 002 : Chaîne d'audit immutable (Merkle chain)

CREATE TABLE snisid_audit.audit_trail (
    audit_id        UUID NOT NULL DEFAULT uuid_generate_v4(),
    event_type      VARCHAR(100) NOT NULL,
    entity_type     VARCHAR(50) NOT NULL,
    entity_id       VARCHAR(100) NOT NULL,
    agent_niu       CHAR(10),
    agence          VARCHAR(50),
    action          VARCHAR(50) NOT NULL,
    before_state    JSONB,
    after_state     JSONB,
    trace_id        VARCHAR(64),
    client_ip       INET,
    geo_location    POINT,
    audit_hash      CHAR(64) NOT NULL,
    chain_hash      CHAR(64),
    sequence_id     BIGSERIAL UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE TABLE snisid_audit.audit_trail_2026_06 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE snisid_audit.audit_trail_2026_07 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE snisid_audit.audit_trail_2026_08 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE snisid_audit.audit_trail_2026_09 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE snisid_audit.audit_trail_2026_10 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE snisid_audit.audit_trail_2026_11 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE snisid_audit.audit_trail_2026_12 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');
CREATE TABLE snisid_audit.audit_trail_2027_01 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2027-01-01') TO ('2027-02-01');
CREATE TABLE snisid_audit.audit_trail_2027_02 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2027-02-01') TO ('2027-03-01');
CREATE TABLE snisid_audit.audit_trail_2027_03 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2027-03-01') TO ('2027-04-01');
CREATE TABLE snisid_audit.audit_trail_2027_04 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2027-04-01') TO ('2027-05-01');
CREATE TABLE snisid_audit.audit_trail_2027_05 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2027-05-01') TO ('2027-06-01');

CREATE RULE no_update_audit AS ON UPDATE TO snisid_audit.audit_trail DO INSTEAD NOTHING;
CREATE RULE no_delete_audit AS ON DELETE TO snisid_audit.audit_trail DO INSTEAD NOTHING;

CREATE INDEX idx_audit_entity ON snisid_audit.audit_trail (entity_type, entity_id, created_at);
CREATE INDEX idx_audit_agent ON snisid_audit.audit_trail (agent_niu, created_at);
CREATE INDEX idx_audit_trace ON snisid_audit.audit_trail (trace_id) WHERE trace_id IS NOT NULL;

CREATE TABLE snisid_audit.verification_log (
    verification_id     UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    niu                 CHAR(10) NOT NULL,
    verification_type   VARCHAR(50) NOT NULL,
    result              BOOLEAN NOT NULL,
    confidence_score    DECIMAL(5,4),
    verifier_niu        CHAR(10),
    agence              VARCHAR(50),
    reason              TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_verif_niu ON snisid_audit.verification_log (niu, created_at);
