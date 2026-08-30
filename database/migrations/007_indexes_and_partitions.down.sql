-- Rollback migration 007
DROP INDEX IF EXISTS idx_citizens_dept_niu;
DROP INDEX IF EXISTS idx_citizens_status_created;
DROP INDEX IF EXISTS idx_audit_entity_agent;
DROP INDEX IF EXISTS idx_dna_active_profiles;
DROP INDEX IF EXISTS idx_person_active_records;
DROP INDEX IF EXISTS idx_property_stolen;
DROP INDEX IF EXISTS idx_person_subject_name_fts;
DROP INDEX IF EXISTS idx_feature_audit_updated;
