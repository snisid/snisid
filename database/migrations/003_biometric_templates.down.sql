-- Rollback migration 003
DROP TABLE IF EXISTS snisid_identity.identity_documents;
DROP TYPE IF EXISTS snisid_identity.document_type;
DROP TABLE IF EXISTS snisid_biometric.biometric_references;
DROP TYPE IF EXISTS snisid_biometric.biometric_modality;
