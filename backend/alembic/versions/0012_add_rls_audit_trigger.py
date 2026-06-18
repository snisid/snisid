"""add row-level security policies and audit trigger

Revision ID: 0012
Revises: 0011
Create Date: 2026-06-11

"""
from __future__ import annotations

from alembic import op

revision: str = "0012"
down_revision: str | None = "0011"
branch_labels: str | None = None
depends_on: str | None = None


def upgrade() -> None:
    # ── Audit log table ──────────────────────────────────────────────────────
    op.execute("""
        CREATE TABLE IF NOT EXISTS bio_audit_log (
            id          BIGSERIAL PRIMARY KEY,
            event_type  VARCHAR(50) NOT NULL,
            table_name  VARCHAR(50) NOT NULL,
            record_id   VARCHAR(36),
            sample_id   VARCHAR(36),
            officer_niu VARCHAR(20),
            agency_code VARCHAR(20),
            purpose     VARCHAR(50),
            case_number VARCHAR(50),
            action      VARCHAR(20) NOT NULL,
            details     JSONB,
            ip_address_hash VARCHAR(64),
            signature   VARCHAR(256),
            created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )
    """)

    op.execute("CREATE INDEX IF NOT EXISTS idx_audit_log_created ON bio_audit_log (created_at DESC)")
    op.execute("CREATE INDEX IF NOT EXISTS idx_audit_log_officer ON bio_audit_log (officer_niu)")
    op.execute("CREATE INDEX IF NOT EXISTS idx_audit_log_table_rec ON bio_audit_log (table_name, record_id)")

    # ── Audit trigger function ───────────────────────────────────────────────
    op.execute("""
        CREATE OR REPLACE FUNCTION bio_audit_trigger_func()
        RETURNS TRIGGER AS $$
        BEGIN
            INSERT INTO bio_audit_log (
                event_type, table_name, record_id, sample_id,
                officer_niu, agency_code, purpose, case_number,
                action, details, ip_address_hash
            ) VALUES (
                'data.' || TG_OP,
                TG_TABLE_NAME,
                COALESCE(NEW.record_id::text, OLD.record_id::text, NEW.id::text, OLD.id::text),
                COALESCE(NEW.sample_id::text, OLD.sample_id::text),
                current_setting('snisid.officer_niu', TRUE),
                current_setting('snisid.agency_code', TRUE),
                current_setting('snisid.purpose', TRUE),
                current_setting('snisid.case_number', TRUE),
                TG_OP,
                row_to_json(COALESCE(NEW, OLD)),
                current_setting('snisid.client_ip', TRUE)
            );
            RETURN COALESCE(NEW, OLD);
        END;
        $$ LANGUAGE plpgsql SECURITY DEFINER
    """)

    # ── Attach trigger to forensic tables ────────────────────────────────────
    for tbl in [
        "bio_str_profiles", "bio_identity_links",
        "per_wanted_persons", "per_missing_persons",
        "per_foreign_fugitives", "per_unidentified_persons",
        "per_terrorism_watches", "per_protection_orders",
        "per_supervised_releases", "per_sex_offenders",
        "per_gang_members", "per_violence_records",
        "per_identity_thefts",
    ]:
        op.execute(f"""
            DROP TRIGGER IF EXISTS trg_bio_audit_{tbl.replace('per_', '').replace('bio_', '')} ON {tbl}
        """)
        op.execute(f"""
            CREATE TRIGGER trg_bio_audit_{tbl.replace('per_', '').replace('bio_', '')}
            AFTER INSERT OR UPDATE OR DELETE ON {tbl}
            FOR EACH ROW EXECUTE FUNCTION bio_audit_trigger_func()
        """)

    # ── Row-Level Security on bio_str_profiles ───────────────────────────────
    op.execute("ALTER TABLE bio_str_profiles ENABLE ROW LEVEL SECURITY")
    op.execute("""
        DROP POLICY IF EXISTS bio_str_profiles_lab_isolation ON bio_str_profiles
    """)
    op.execute("""
        CREATE POLICY bio_str_profiles_lab_isolation ON bio_str_profiles
        FOR ALL
        USING (
            current_setting('snisid.lab_id', TRUE) = '' OR
            lab_id::text = current_setting('snisid.lab_id', TRUE) OR
            current_setting('snisid.role', TRUE) IN ('bio.sdis.operator', 'bio.ndis.analyst', 'bio.admin')
        )
    """)

    # ── Row-Level Security on bio_identity_links ─────────────────────────────
    op.execute("ALTER TABLE bio_identity_links ENABLE ROW LEVEL SECURITY")
    op.execute("""
        DROP POLICY IF EXISTS bio_identity_links_dcpj_only ON bio_identity_links
    """)
    op.execute("""
        CREATE POLICY bio_identity_links_dcpj_only ON bio_identity_links
        FOR ALL
        USING (
            current_setting('snisid.role', TRUE) = 'bio.dcpj.director'
        )
    """)

    # ── RLS on Per tables ────────────────────────────────────────────────────
    for tbl in [
        "per_wanted_persons", "per_missing_persons",
        "per_foreign_fugitives", "per_unidentified_persons",
        "per_terrorism_watches", "per_protection_orders",
        "per_supervised_releases", "per_sex_offenders",
        "per_gang_members", "per_violence_records",
        "per_identity_thefts",
    ]:
        op.execute(f"ALTER TABLE {tbl} ENABLE ROW LEVEL SECURITY")
        op.execute(f"""
            DROP POLICY IF EXISTS {tbl}_dcpj_access ON {tbl}
        """)
        op.execute(f"""
            CREATE POLICY {tbl}_dcpj_access ON {tbl}
            FOR ALL
            USING (
                current_setting('snisid.role', TRUE) IN (
                    'bio.dcpj.investigator', 'bio.dcpj.director',
                    'bio.sdis.operator', 'bio.ndis.analyst',
                    'bio.admin'
                )
            )
        """)


def downgrade() -> None:
    for tbl in [
        "per_wanted_persons", "per_missing_persons",
        "per_foreign_fugitives", "per_unidentified_persons",
        "per_terrorism_watches", "per_protection_orders",
        "per_supervised_releases", "per_sex_offenders",
        "per_gang_members", "per_violence_records",
        "per_identity_thefts",
    ]:
        op.execute(f"DROP POLICY IF EXISTS {tbl}_dcpj_access ON {tbl}")
        op.execute(f"ALTER TABLE {tbl} DISABLE ROW LEVEL SECURITY")

    op.execute("DROP POLICY IF EXISTS bio_identity_links_dcpj_only ON bio_identity_links")
    op.execute("ALTER TABLE bio_identity_links DISABLE ROW LEVEL SECURITY")

    op.execute("DROP POLICY IF EXISTS bio_str_profiles_lab_isolation ON bio_str_profiles")
    op.execute("ALTER TABLE bio_str_profiles DISABLE ROW LEVEL SECURITY")

    for tbl in [
        "bio_str_profiles", "bio_identity_links",
        "per_wanted_persons", "per_missing_persons",
        "per_foreign_fugitives", "per_unidentified_persons",
        "per_terrorism_watches", "per_protection_orders",
        "per_supervised_releases", "per_sex_offenders",
        "per_gang_members", "per_violence_records",
        "per_identity_thefts",
    ]:
        suf = tbl.replace("per_", "").replace("bio_", "")
        op.execute(f"DROP TRIGGER IF EXISTS trg_bio_audit_{suf} ON {tbl}")

    op.execute("DROP FUNCTION IF EXISTS bio_audit_trigger_func")
    op.execute("DROP TABLE IF EXISTS bio_audit_log")
