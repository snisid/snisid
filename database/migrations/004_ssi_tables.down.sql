-- Rollback migration 004
DROP TABLE IF EXISTS snisid_ssi.revocation_events;
DROP TABLE IF EXISTS snisid_ssi.credential_flows;
DROP TABLE IF EXISTS snisid_ssi.status_lists;
DROP TABLE IF EXISTS snisid_ssi.verifiable_credentials;
DROP TABLE IF EXISTS snisid_ssi.did_records;
DROP SCHEMA IF EXISTS snisid_ssi CASCADE;
