CREATE SCHEMA IF NOT EXISTS snisid_core;

CREATE TABLE IF NOT EXISTS snisid_core.citizens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    niu VARCHAR(20) UNIQUE NOT NULL,
    first_name VARCHAR(200) NOT NULL,
    last_name VARCHAR(200) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(10),
    nationality VARCHAR(100) DEFAULT 'HAITI',
    place_of_birth VARCHAR(200),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS snisid_core.audit_trail (
    id BIGSERIAL,
    event_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(36) NOT NULL,
    actor_id VARCHAR(36) NOT NULL,
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE TABLE IF NOT EXISTS snisid_core.bio_identity_links (
    sample_id VARCHAR(36) PRIMARY KEY,
    niu VARCHAR(20) NOT NULL,
    linked_by VARCHAR(100) NOT NULL,
    linked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    court_order_ref VARCHAR(200),
    purpose VARCHAR(100) NOT NULL
);
