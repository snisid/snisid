-- Rollback migration 002
DROP TABLE IF EXISTS snisid_audit.verification_log;
DROP RULE IF EXISTS no_update_audit ON snisid_audit.audit_trail;
DROP RULE IF EXISTS no_delete_audit ON snisid_audit.audit_trail;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2027_05;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2027_04;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2027_03;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2027_02;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2027_01;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2026_12;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2026_11;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2026_10;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2026_09;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2026_08;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2026_07;
DROP TABLE IF EXISTS snisid_audit.audit_trail_2026_06;
DROP TABLE IF EXISTS snisid_audit.audit_trail;
