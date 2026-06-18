"""create event store, identity, and agency tables

Revision ID: 0003
Revises: 0002
Create Date: 2026-06-10
"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import JSONB

revision: str = "0003"
down_revision: Union[str, None] = "0002"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # ── Enums ────────────────────────────────────────────────────────
    sa.Enum("pending", "active", "suspended", "revoked", "deceased",
            name="identity_status_enum").create(op.get_bind(), checkfirst=True)
    sa.Enum("male", "female", "other",
            name="gender_enum").create(op.get_bind(), checkfirst=True)
    sa.Enum("national_id", "passport", "birth_certificate",
            "driving_license", "residence_permit",
            name="document_type_enum").create(op.get_bind(), checkfirst=True)
    sa.Enum("active", "expired", "revoked", "lost", "stolen",
            name="document_status_enum").create(op.get_bind(), checkfirst=True)
    sa.Enum("fingerprint", "iris", "face", "voice",
            name="biometric_type_enum").create(op.get_bind(), checkfirst=True)
    sa.Enum("active", "inactive", "suspended",
            name="agency_status_enum").create(op.get_bind(), checkfirst=True)
    sa.Enum("central", "regional", "local", "mobile",
            name="agency_type_enum").create(op.get_bind(), checkfirst=True)

    # ── Event Store ──────────────────────────────────────────────────
    op.create_table(
        "event_store",
        sa.Column("id", sa.BigInteger(), primary_key=True, autoincrement=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now(), onupdate=sa.func.now()),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("event_id", sa.String(36), unique=True, nullable=False),
        sa.Column("aggregate_id", sa.String(36), nullable=False, index=True),
        sa.Column("aggregate_type", sa.String(100), nullable=False, index=True),
        sa.Column("event_type", sa.String(200), nullable=False, index=True),
        sa.Column("event_data", JSONB(), nullable=False),
        sa.Column("event_metadata", JSONB(), nullable=True),
        sa.Column("version", sa.Integer(), nullable=False),
        sa.Column("actor_id", sa.String(36), nullable=True),
        sa.Column("correlation_id", sa.String(36), nullable=True),
        sa.Column("timestamp", sa.DateTime(timezone=True), nullable=False),
        Index("ix_event_store_aggregate_version", "aggregate_id", "version", unique=True),
    )
    op.create_table(
        "aggregate_snapshots",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now(), onupdate=sa.func.now()),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("aggregate_id", sa.String(36), unique=True, nullable=False, index=True),
        sa.Column("aggregate_type", sa.String(100), nullable=False),
        sa.Column("snapshot_data", JSONB(), nullable=False),
        sa.Column("version", sa.Integer(), nullable=False),
    )

    # ── Citizens ─────────────────────────────────────────────────────
    op.create_table(
        "citizens",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now(), onupdate=sa.func.now()),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("national_id", sa.String(20), unique=True, nullable=False, index=True),
        sa.Column("first_name", sa.String(100), nullable=False),
        sa.Column("last_name", sa.String(100), nullable=False),
        sa.Column("middle_name", sa.String(100), nullable=True),
        sa.Column("date_of_birth", sa.Date(), nullable=False),
        sa.Column("place_of_birth", sa.String(200), nullable=False),
        sa.Column("gender", sa.Enum("male", "female", "other",
                  name="gender_enum"), nullable=False),
        sa.Column("nationality", sa.String(3), nullable=False),
        sa.Column("marital_status", sa.String(20), nullable=True),
        sa.Column("email", sa.String(255), nullable=True),
        sa.Column("phone", sa.String(20), nullable=True),
        sa.Column("address", JSONB(), nullable=True),
        sa.Column("photo_url", sa.String(500), nullable=True),
        sa.Column("fingerprint_hash", sa.String(128), nullable=True),
        sa.Column("iris_hash", sa.String(128), nullable=True),
        sa.Column("face_encoding_hash", sa.String(128), nullable=True),
        sa.Column("status", sa.Enum("pending", "active", "suspended",
                  "revoked", "deceased", name="identity_status_enum"),
                  nullable=False, default="pending", index=True),
        sa.Column("verified_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("verified_by", sa.String(36), nullable=True),
        sa.Column("suspension_reason", sa.Text(), nullable=True),
        sa.Column("revocation_reason", sa.Text(), nullable=True),
        sa.Column("agency_id", sa.String(36), nullable=False, index=True),
        sa.Column("created_by", sa.String(36), nullable=False),
        sa.Column("version", sa.Integer(), nullable=False, default=0),
        Index("ix_citizens_name", "last_name", "first_name"),
        Index("ix_citizens_dob", "date_of_birth"),
        Index("ix_citizens_agency_status", "agency_id", "status"),
    )
    op.create_table(
        "identity_documents",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now(), onupdate=sa.func.now()),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("citizen_id", sa.String(36),
                  sa.ForeignKey("citizens.id", ondelete="CASCADE"),
                  nullable=False, index=True),
        sa.Column("document_type", sa.Enum("national_id", "passport",
                  "birth_certificate", "driving_license", "residence_permit",
                  name="document_type_enum"), nullable=False),
        sa.Column("document_number", sa.String(50), nullable=False),
        sa.Column("issue_date", sa.Date(), nullable=False),
        sa.Column("expiry_date", sa.Date(), nullable=True),
        sa.Column("issuing_agency", sa.String(100), nullable=False),
        sa.Column("issuing_location", sa.String(200), nullable=True),
        sa.Column("status", sa.Enum("active", "expired", "revoked",
                  "lost", "stolen", name="document_status_enum"),
                  nullable=False, default="active"),
        sa.Column("metadata_json", JSONB(), nullable=True),
        Index("uq_document_type_number", "document_type", "document_number", unique=True),
    )
    op.create_table(
        "biometric_records",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now(), onupdate=sa.func.now()),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("citizen_id", sa.String(36),
                  sa.ForeignKey("citizens.id", ondelete="CASCADE"),
                  nullable=False, index=True),
        sa.Column("biometric_type", sa.Enum("fingerprint", "iris", "face",
                  "voice", name="biometric_type_enum"), nullable=False),
        sa.Column("template_hash", sa.String(256), unique=True, nullable=False),
        sa.Column("quality_score", sa.Float(), nullable=False),
        sa.Column("capture_device", sa.String(100), nullable=True),
        sa.Column("captured_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("captured_by", sa.String(36), nullable=False),
        sa.Column("is_primary", sa.Boolean(), default=False, nullable=False),
        Index("ix_biometric_type_citizen", "biometric_type", "citizen_id"),
    )
    op.create_table(
        "citizens_read",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now(), onupdate=sa.func.now()),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("national_id", sa.String(20), unique=True, nullable=False, index=True),
        sa.Column("full_name", sa.String(300), nullable=False),
        sa.Column("first_name", sa.String(100), nullable=False),
        sa.Column("last_name", sa.String(100), nullable=False),
        sa.Column("date_of_birth", sa.Date(), nullable=False),
        sa.Column("gender", sa.String(10), nullable=False),
        sa.Column("nationality", sa.String(3), nullable=False),
        sa.Column("status", sa.String(20), nullable=False, index=True),
        sa.Column("agency_id", sa.String(36), nullable=False, index=True),
        sa.Column("document_count", sa.Integer(), default=0),
        sa.Column("has_biometrics", sa.Boolean(), default=False),
        sa.Column("verified", sa.Boolean(), default=False),
        sa.Column("photo_url", sa.String(500), nullable=True),
        sa.Column("address_summary", sa.String(500), nullable=True),
        sa.Column("last_event_at", sa.DateTime(timezone=True), nullable=True),
        Index("ix_read_fullname", "full_name"),
    )

    # ── Agencies ─────────────────────────────────────────────────────
    op.create_table(
        "agencies",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now(), onupdate=sa.func.now()),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("name", sa.String(200), nullable=False),
        sa.Column("code", sa.String(20), unique=True, nullable=False, index=True),
        sa.Column("agency_type", sa.Enum("central", "regional", "local",
                  "mobile", name="agency_type_enum"), nullable=False),
        sa.Column("status", sa.Enum("active", "inactive", "suspended",
                  name="agency_status_enum"),
                  nullable=False, default="active", index=True),
        sa.Column("address", sa.Text(), nullable=True),
        sa.Column("city", sa.String(100), nullable=True),
        sa.Column("department", sa.String(100), nullable=True),
        sa.Column("phone", sa.String(20), nullable=True),
        sa.Column("email", sa.String(255), nullable=True),
        sa.Column("latitude", sa.Float(), nullable=True),
        sa.Column("longitude", sa.Float(), nullable=True),
        sa.Column("max_daily_enrollments", sa.Integer(), default=500),
        sa.Column("is_headquarters", sa.Boolean(), default=False),
        sa.Column("parent_agency_id", sa.String(36), nullable=True),
        sa.Column("opened_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("closed_at", sa.DateTime(timezone=True), nullable=True),
    )


def downgrade() -> None:
    op.drop_table("agencies")
    op.drop_table("citizens_read")
    op.drop_table("biometric_records")
    op.drop_table("identity_documents")
    op.drop_table("citizens")
    op.drop_table("aggregate_snapshots")
    op.drop_table("event_store")
    for name in ("identity_status_enum", "gender_enum", "document_type_enum",
                 "document_status_enum", "biometric_type_enum",
                 "agency_status_enum", "agency_type_enum"):
        sa.Enum(name=name).drop(op.get_bind(), checkfirst=True)
