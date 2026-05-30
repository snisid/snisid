-- Revoke default public schema access globally
REVOKE ALL ON SCHEMA public FROM PUBLIC;

-- 1. Create Schema Manager (Used by Alembic only, restricted after deployment)
CREATE ROLE schema_manager WITH LOGIN PASSWORD 'vault_injected_password' CREATEROLE CREATEDB;
GRANT ALL ON SCHEMA public TO schema_manager;

-- 2. Create Microservice Application Roles (Principle of Least Privilege)
CREATE ROLE identity_svc WITH LOGIN PASSWORD 'vault_injected_password';
GRANT USAGE ON SCHEMA public TO identity_svc;
GRANT SELECT, INSERT, UPDATE ON TABLE citizens, events_store TO identity_svc;
-- identity_svc cannot delete citizens (immutable append-only logs)

CREATE ROLE biometric_svc WITH LOGIN PASSWORD 'vault_injected_password';
GRANT USAGE ON SCHEMA public TO biometric_svc;
GRANT SELECT, INSERT, DELETE ON TABLE biometric_templates TO biometric_svc;
GRANT SELECT ON TABLE citizens TO biometric_svc; -- Read-only access to citizens

-- 3. Create Auditor Role (Read-Only across specific schemas)
CREATE ROLE auditor_role WITH LOGIN PASSWORD 'vault_injected_password';
GRANT USAGE ON SCHEMA public TO auditor_role;
GRANT SELECT ON TABLE audit_logs TO auditor_role;

-- 4. Create Backup Role for pgBackRest
CREATE ROLE backup_user WITH LOGIN REPLICATION PASSWORD 'vault_injected_password';
GRANT EXECUTE ON FUNCTION pg_start_backup(text, boolean, boolean) TO backup_user;
GRANT EXECUTE ON FUNCTION pg_stop_backup(boolean, boolean) TO backup_user;

-- 5. Hardening: Drop the default 'postgres' superuser privileges if required by strict CIS benchmark
-- NOTE: In Patroni environments, 'postgres' is often needed internally for cluster management.
-- Instead, we restrict its remote login via pg_hba.conf.
