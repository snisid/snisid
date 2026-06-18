-- Rollback migration 006
DROP TABLE IF EXISTS snisid_ml.model_registry;
DROP TABLE IF EXISTS snisid_ml.feature_cache_audit_2026_q3;
DROP TABLE IF EXISTS snisid_ml.feature_cache_audit;
DROP SCHEMA IF EXISTS snisid_ml CASCADE;
