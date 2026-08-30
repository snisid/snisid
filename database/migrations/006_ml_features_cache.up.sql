-- Migration 006 : Cache ML features + model registry

CREATE SCHEMA IF NOT EXISTS snisid_ml;

CREATE TABLE snisid_ml.feature_cache_audit (
    audit_id        UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    niu             CHAR(10) NOT NULL,
    feature_name    VARCHAR(100) NOT NULL,
    old_value       DOUBLE PRECISION,
    new_value       DOUBLE PRECISION NOT NULL,
    source          VARCHAR(50) NOT NULL,
    kafka_offset    BIGINT,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (updated_at);

CREATE TABLE snisid_ml.feature_cache_audit_2026_q3
    PARTITION OF snisid_ml.feature_cache_audit
    FOR VALUES FROM ('2026-07-01') TO ('2026-10-01');

CREATE INDEX idx_feature_audit_niu ON snisid_ml.feature_cache_audit (niu, feature_name, updated_at);

CREATE TABLE snisid_ml.model_registry (
    model_id        UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    model_name      VARCHAR(100) NOT NULL,
    model_version   VARCHAR(50) NOT NULL,
    model_type      VARCHAR(50) NOT NULL,
    artifact_path   VARCHAR(500),
    metrics         JSONB,
    status          VARCHAR(20) NOT NULL DEFAULT 'SHADOW'
                    CHECK (status IN ('SHADOW','DEPLOYED','RETIRED')),
    deployed_at     TIMESTAMPTZ,
    retired_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (model_name, model_version)
);

CREATE INDEX idx_model_status ON snisid_ml.model_registry (model_name, status);
