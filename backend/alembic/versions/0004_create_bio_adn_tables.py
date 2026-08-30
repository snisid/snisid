"""create bio adn tables

Revision ID: 0004
Revises: 0003
Create Date: 2026-06-11
"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import JSONB

revision: str = "0004"
down_revision: Union[str, None] = "0003"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "bio_laboratories",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("lab_code", sa.String(20), unique=True, nullable=False),
        sa.Column("lab_name", sa.String(200), nullable=False),
        sa.Column("lab_level", sa.String(10), nullable=False),
        sa.Column("department", sa.String(50), nullable=True),
        sa.Column("institution", sa.String(100), nullable=True),
        sa.Column("accreditation", sa.String(100), nullable=True),
        sa.Column("contact_email", sa.String(200), nullable=True),
        sa.Column("is_active", sa.Boolean(), default=True),
        sa.CheckConstraint("lab_level IN ('LDIS','SDIS','NDIS')", name="ck_lab_level"),
    )
    op.create_table(
        "bio_str_profiles",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("specimen_number", sa.String(100), unique=True, nullable=False),
        sa.Column("index_type", sa.String(10), nullable=False),
        sa.Column("loci_encrypted", sa.LargeBinary, nullable=False),
        sa.Column("loci_hash", sa.String(64), nullable=False),
        sa.Column("amelogenin", sa.String(2), nullable=True),
        sa.Column("quality_score", sa.Numeric(4, 3), nullable=True),
        sa.Column("loci_count", sa.SmallInteger, default=20),
        sa.Column("lab_id", sa.String(36), sa.ForeignKey("bio_laboratories.id"), nullable=True),
        sa.Column("case_number", sa.String(100), nullable=True),
        sa.Column("collected_date", sa.Date, nullable=False),
        sa.Column("analysis_date", sa.Date, nullable=True),
        sa.Column("uploaded_ldis", sa.Boolean(), default=False),
        sa.Column("uploaded_sdis", sa.Boolean(), default=False),
        sa.Column("uploaded_ndis", sa.Boolean(), default=False),
        sa.Column("ndis_upload_date", sa.DateTime(timezone=True), nullable=True),
        sa.Column("is_expunged", sa.Boolean(), default=False),
        sa.Column("expunge_date", sa.DateTime(timezone=True), nullable=True),
        sa.Column("expunge_order", sa.String(200), nullable=True),
        sa.CheckConstraint("index_type IN ('BIO-CON','BIO-ARR','BIO-FSC','BIO-DIS','BIO-RNI')", name="ck_index_type"),
        sa.CheckConstraint("quality_score BETWEEN 0 AND 1", name="ck_quality_score"),
    )
    op.create_index("idx_bio_str_hash", "bio_str_profiles", ["loci_hash"])
    op.create_index("idx_bio_str_index_type", "bio_str_profiles", ["index_type"])
    op.create_index("idx_bio_str_case", "bio_str_profiles", ["case_number"])

    op.create_table(
        "bio_identity_links",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("sample_id", sa.String(36), sa.ForeignKey("bio_str_profiles.id"), unique=True, nullable=False),
        sa.Column("niu", sa.String(20), nullable=True),
        sa.Column("linked_by_agent", sa.String(100), nullable=False),
        sa.Column("linked_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("court_order_ref", sa.String(200), nullable=False),
        sa.Column("purpose", sa.String(100), nullable=False),
        sa.Column("reviewed_by", sa.String(100), nullable=True),
        sa.Column("reviewed_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("review_outcome", sa.String(20), nullable=True),
    )
    op.create_table(
        "bio_hits",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("query_sample_id", sa.String(36), sa.ForeignKey("bio_str_profiles.id"), nullable=False),
        sa.Column("match_sample_id", sa.String(36), sa.ForeignKey("bio_str_profiles.id"), nullable=False),
        sa.Column("match_type", sa.String(20), nullable=False),
        sa.Column("confidence", sa.Numeric(5, 4), nullable=False),
        sa.Column("matched_loci", sa.SmallInteger, nullable=False),
        sa.Column("total_loci", sa.SmallInteger, nullable=False),
        sa.Column("hit_level", sa.String(10), nullable=False),
        sa.Column("alert_sent", sa.Boolean(), default=False),
        sa.Column("alert_sent_at", sa.DateTime(timezone=True), nullable=True),
        sa.CheckConstraint("match_type IN ('FULL_MATCH','PARTIAL','FAMILIAL')", name="ck_match_type"),
        sa.CheckConstraint("hit_level IN ('LDIS','SDIS','NDIS')", name="ck_hit_level"),
    )
    op.create_table(
        "per_wanted_persons",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("niu", sa.String(20), nullable=True),
        sa.Column("last_name", sa.String(100), nullable=True),
        sa.Column("first_name", sa.String(100), nullable=True),
        sa.Column("aliases", JSONB, default=list),
        sa.Column("date_of_birth", sa.Date, nullable=True),
        sa.Column("gender", sa.String(1), nullable=True),
        sa.Column("nationality", sa.String(3), nullable=True),
        sa.Column("warrant_type", sa.String(50), nullable=False),
        sa.Column("warrant_number", sa.String(100), nullable=True),
        sa.Column("issuing_court", sa.String(200), nullable=True),
        sa.Column("issuing_date", sa.Date, nullable=False),
        sa.Column("charges", JSONB, nullable=False),
        sa.Column("danger_level", sa.String(10), default="MEDIUM"),
        sa.Column("armed_dangerous", sa.Boolean(), default=False),
        sa.Column("height_cm", sa.SmallInteger, nullable=True),
        sa.Column("weight_kg", sa.SmallInteger, nullable=True),
        sa.Column("eye_color", sa.String(30), nullable=True),
        sa.Column("hair_color", sa.String(30), nullable=True),
        sa.Column("distinguishing_marks", sa.Text, nullable=True),
        sa.Column("entering_agency", sa.String(100), nullable=False),
        sa.Column("entering_officer", sa.String(100), nullable=False),
        sa.Column("mco_contact", sa.String(200), nullable=True),
        sa.Column("last_known_location", sa.String(200), nullable=True),
        sa.Column("status", sa.String(20), default="ACTIVE"),
        sa.Column("expiry_date", sa.Date, nullable=True),
        sa.Column("fingerprint_ref", sa.String(36), nullable=True),
        sa.Column("photo_refs", JSONB, default=list),
        sa.Column("bio_sample_ref", sa.String(36), sa.ForeignKey("bio_str_profiles.id"), nullable=True),
        sa.Column("interpol_notice", sa.String(50), nullable=True),
        sa.Column("last_hit_at", sa.DateTime(timezone=True), nullable=True),
        sa.CheckConstraint("danger_level IN ('LOW','MEDIUM','HIGH','CRITICAL')", name="ck_danger_level"),
        sa.CheckConstraint("status IN ('ACTIVE','CLEARED','EXPIRED','SUSPENDED')", name="ck_wanted_status"),
    )
    op.create_index("idx_per_wanted_niu", "per_wanted_persons", ["niu"])
    op.create_index("idx_per_wanted_status", "per_wanted_persons", ["status"])
    op.create_index("idx_per_wanted_name", "per_wanted_persons",
                    [sa.text("to_tsvector('french', COALESCE(last_name,'') || ' ' || COALESCE(first_name,''))")],
                    postgresql_using="gin")

    op.create_table(
        "per_missing_persons",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("niu", sa.String(20), nullable=True),
        sa.Column("last_name", sa.String(100), nullable=False),
        sa.Column("first_name", sa.String(100), nullable=False),
        sa.Column("date_of_birth", sa.Date, nullable=True),
        sa.Column("age_at_missing", sa.SmallInteger, nullable=True),
        sa.Column("gender", sa.String(1), nullable=True),
        sa.Column("nationality", sa.String(3), nullable=True),
        sa.Column("category", sa.String(20), nullable=False),
        sa.Column("missing_date", sa.DateTime(timezone=True), nullable=False),
        sa.Column("missing_location", sa.String(200), nullable=False),
        sa.Column("circumstances", sa.Text, nullable=True),
        sa.Column("last_seen_clothing", sa.String(500), nullable=True),
        sa.Column("height_cm", sa.SmallInteger, nullable=True),
        sa.Column("weight_kg", sa.SmallInteger, nullable=True),
        sa.Column("distinctive_features", sa.Text, nullable=True),
        sa.Column("family_contact", sa.String(200), nullable=True),
        sa.Column("family_phone", sa.String(50), nullable=True),
        sa.Column("photo_refs", JSONB, default=list),
        sa.Column("bio_sample_ref", sa.String(36), sa.ForeignKey("bio_str_profiles.id"), nullable=True),
        sa.Column("family_bio_refs", JSONB, default=list),
        sa.Column("medical_conditions", sa.Text, nullable=True),
        sa.Column("medications", sa.Text, nullable=True),
        sa.Column("status", sa.String(20), default="ACTIVE"),
        sa.Column("located_date", sa.DateTime(timezone=True), nullable=True),
        sa.Column("entering_agency", sa.String(100), nullable=False),
        sa.Column("ncmec_notified", sa.Boolean(), default=False),
        sa.CheckConstraint("category IN ('CHILD','ENDANGERED','INVOLUNTARY','CATASTROPHE','UNEMANCIPATED','OTHER')", name="ck_missing_category"),
        sa.CheckConstraint("status IN ('ACTIVE','LOCATED','DECEASED','CANCELLED')", name="ck_missing_status"),
    )
    op.create_table(
        "per_sex_offenders",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("niu", sa.String(20), nullable=False),
        sa.Column("conviction_date", sa.Date, nullable=False),
        sa.Column("conviction_court", sa.String(200), nullable=False),
        sa.Column("offenses", JSONB, nullable=False),
        sa.Column("risk_level", sa.String(10), nullable=True),
        sa.Column("registration_expiry", sa.Date, nullable=True),
        sa.Column("current_address", sa.Text, nullable=True),
        sa.Column("employer", sa.String(200), nullable=True),
        sa.Column("restrictions", sa.Text, nullable=True),
        sa.Column("last_verified", sa.Date, nullable=True),
        sa.Column("status", sa.String(20), default="ACTIVE"),
        sa.CheckConstraint("risk_level IN ('LOW','MEDIUM','HIGH')", name="ck_risk_level"),
    )
    op.create_table(
        "per_gang_members",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("niu", sa.String(20), nullable=True),
        sa.Column("last_name", sa.String(100), nullable=True),
        sa.Column("first_name", sa.String(100), nullable=True),
        sa.Column("aliases", JSONB, default=list),
        sa.Column("gang_name", sa.String(200), nullable=False),
        sa.Column("gang_code", sa.String(50), nullable=True),
        sa.Column("membership_type", sa.String(30), nullable=True),
        sa.Column("territory", sa.String(200), nullable=True),
        sa.Column("known_weapons", JSONB, default=list),
        sa.Column("criminal_activities", JSONB, default=list),
        sa.Column("threat_level", sa.String(10), default="HIGH"),
        sa.Column("intelligence_notes", sa.Text, nullable=True),
        sa.Column("source_reliability", sa.String(10), nullable=True),
        sa.Column("status", sa.String(20), default="ACTIVE"),
        sa.CheckConstraint("membership_type IN ('LEADER','MEMBER','ASSOCIATE','PROSPECT')", name="ck_membership_type"),
        sa.CheckConstraint("threat_level IN ('LOW','MEDIUM','HIGH','CRITICAL')", name="ck_threat_level"),
    )
    op.create_table(
        "bie_stolen_vehicles",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("vin", sa.String(17), unique=True, nullable=True),
        sa.Column("plate_number", sa.String(20), nullable=True),
        sa.Column("plate_dept", sa.String(50), nullable=True),
        sa.Column("vehicle_make", sa.String(100), nullable=True),
        sa.Column("vehicle_model", sa.String(100), nullable=True),
        sa.Column("vehicle_year", sa.SmallInteger, nullable=True),
        sa.Column("vehicle_color", sa.String(50), nullable=True),
        sa.Column("vehicle_type", sa.String(50), nullable=True),
        sa.Column("theft_date", sa.Date, nullable=False),
        sa.Column("theft_location", sa.String(200), nullable=False),
        sa.Column("theft_department", sa.String(50), nullable=True),
        sa.Column("owner_niu", sa.String(20), nullable=True),
        sa.Column("owner_name", sa.String(200), nullable=True),
        sa.Column("owner_phone", sa.String(50), nullable=True),
        sa.Column("foves_record_id", sa.String(36), nullable=True),
        sa.Column("status", sa.String(20), default="STOLEN"),
        sa.Column("recovered_date", sa.Date, nullable=True),
        sa.Column("recovered_location", sa.String(200), nullable=True),
        sa.Column("entering_agency", sa.String(100), nullable=False),
        sa.CheckConstraint("status IN ('STOLEN','RECOVERED','CANCELLED')", name="ck_vehicle_status"),
    )
    op.create_table(
        "bie_stolen_firearms",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("serial_number", sa.String(100), unique=True, nullable=False),
        sa.Column("make", sa.String(100), nullable=True),
        sa.Column("model", sa.String(100), nullable=True),
        sa.Column("caliber", sa.String(50), nullable=True),
        sa.Column("firearm_type", sa.String(50), nullable=True),
        sa.Column("barrel_length", sa.Float, nullable=True),
        sa.Column("theft_date", sa.Date, nullable=False),
        sa.Column("theft_location", sa.String(200), nullable=True),
        sa.Column("owner_niu", sa.String(20), nullable=True),
        sa.Column("status", sa.String(20), default="STOLEN"),
        sa.Column("recovered_date", sa.Date, nullable=True),
        sa.Column("entering_agency", sa.String(100), nullable=False),
        sa.CheckConstraint("status IN ('STOLEN','RECOVERED','CANCELLED')", name="ck_firearm_status"),
    )
    op.create_table(
        "bie_stolen_documents",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("document_type", sa.String(50), nullable=False),
        sa.Column("document_number", sa.String(100), nullable=True),
        sa.Column("issuing_agency", sa.String(100), nullable=True),
        sa.Column("issue_date", sa.Date, nullable=True),
        sa.Column("expiry_date", sa.Date, nullable=True),
        sa.Column("owner_niu", sa.String(20), nullable=True),
        sa.Column("owner_name", sa.String(200), nullable=True),
        sa.Column("report_date", sa.Date, nullable=False),
        sa.Column("report_location", sa.String(200), nullable=True),
        sa.Column("theft_type", sa.String(20), default="STOLEN"),
        sa.Column("status", sa.String(20), default="ACTIVE"),
        sa.CheckConstraint("document_type IN ('PASSPORT','CIN','ACTE_NAISSANCE','PERMIS_CONDUIRE','TITRE_FONCIER','AUTRE')", name="ck_document_type"),
        sa.CheckConstraint("theft_type IN ('STOLEN','LOST','FORGED')", name="ck_theft_type"),
        sa.CheckConstraint("status IN ('ACTIVE','RECOVERED','CANCELLED')", name="ck_document_status"),
    )
    op.create_table(
        "bie_stolen_vessels",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("vessel_name", sa.String(200), nullable=True),
        sa.Column("registration_number", sa.String(100), nullable=True),
        sa.Column("hull_id_number", sa.String(50), nullable=True),
        sa.Column("vessel_type", sa.String(50), nullable=True),
        sa.Column("vessel_make", sa.String(100), nullable=True),
        sa.Column("vessel_length_m", sa.Float, nullable=True),
        sa.Column("hull_color", sa.String(50), nullable=True),
        sa.Column("home_port", sa.String(200), nullable=True),
        sa.Column("theft_location", sa.String(200), nullable=False),
        sa.Column("theft_date", sa.Date, nullable=False),
        sa.Column("owner_niu", sa.String(20), nullable=True),
        sa.Column("owner_name", sa.String(200), nullable=True),
        sa.Column("status", sa.String(20), default="STOLEN"),
    )
    op.create_table(
        "bio_audit_log",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("event_type", sa.String(100), nullable=False),
        sa.Column("table_name", sa.String(100), nullable=True),
        sa.Column("record_id", sa.String(36), nullable=True),
        sa.Column("officer_niu", sa.String(20), nullable=False),
        sa.Column("agency_code", sa.String(50), nullable=False),
        sa.Column("purpose", sa.String(200), nullable=False),
        sa.Column("case_number", sa.String(100), nullable=True),
        sa.Column("ip_hash", sa.String(64), nullable=True),
        sa.Column("action", sa.String(20), nullable=False),
        sa.Column("details", JSONB, nullable=True),
        sa.Column("signature", sa.Text, nullable=True),
        sa.CheckConstraint("action IN ('CREATE','READ','UPDATE','DELETE','SEARCH','HIT')", name="ck_action"),
    )

    op.execute("ALTER TABLE bio_str_profiles ENABLE ROW LEVEL SECURITY")
    op.execute("ALTER TABLE per_wanted_persons ENABLE ROW LEVEL SECURITY")
    op.execute("ALTER TABLE per_gang_members ENABLE ROW LEVEL SECURITY")
    op.execute("ALTER TABLE bio_identity_links ENABLE ROW LEVEL SECURITY")
    op.execute("""
        CREATE POLICY bio_identity_links_policy ON bio_identity_links
            USING (current_user = 'snisid_dcpj_director' OR current_user = 'snisid_admin')
    """)


def downgrade() -> None:
    op.execute("DROP POLICY IF EXISTS bio_identity_links_policy ON bio_identity_links")
    op.execute("ALTER TABLE bio_identity_links DISABLE ROW LEVEL SECURITY")
    op.execute("ALTER TABLE per_gang_members DISABLE ROW LEVEL SECURITY")
    op.execute("ALTER TABLE per_wanted_persons DISABLE ROW LEVEL SECURITY")
    op.execute("ALTER TABLE bio_str_profiles DISABLE ROW LEVEL SECURITY")
    op.drop_table("bio_audit_log")
    op.drop_table("bie_stolen_vessels")
    op.drop_table("bie_stolen_documents")
    op.drop_table("bie_stolen_firearms")
    op.drop_table("bie_stolen_vehicles")
    op.drop_table("per_gang_members")
    op.drop_table("per_sex_offenders")
    op.drop_table("per_missing_persons")
    op.drop_table("per_wanted_persons")
    op.drop_table("bio_hits")
    op.drop_table("bio_identity_links")
    op.drop_table("bio_str_profiles")
    op.drop_table("bio_laboratories")
