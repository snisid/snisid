-- Rollback migration 005
DROP TABLE IF EXISTS snisid_criminal.warrants;
DROP TABLE IF EXISTS snisid_criminal.evidence;
DROP TABLE IF EXISTS snisid_criminal.case_persons;
DROP TABLE IF EXISTS snisid_criminal.cases;
DROP TYPE IF EXISTS snisid_criminal.evidence_type;
DROP TYPE IF EXISTS snisid_criminal.case_status;
DROP SCHEMA IF EXISTS snisid_criminal CASCADE;
