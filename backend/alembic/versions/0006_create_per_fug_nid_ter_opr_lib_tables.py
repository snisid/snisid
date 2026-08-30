"""create per fug nid ter opr lib tables + update existing per

Revision ID: 0006
Revises: 0005
Create Date: 2026-06-11
"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import JSONB

revision: str = "0006"
down_revision: Union[str, None] = "0005"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # ── PER-FUG: Foreign Fugitives (mandats rouges Interpol, etc.) ────────
    op.create_table(
        "per_foreign_fugitives",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("interpol_notice_number", sa.String(50), nullable=False),
        sa.Column("notice_type", sa.String(20), nullable=False),
        sa.Column("last_name", sa.String(100), nullable=False),
        sa.Column("first_name", sa.String(100), nullable=True),
        sa.Column("aliases", JSONB, default=list),
        sa.Column("date_of_birth", sa.Date, nullable=True),
        sa.Column("gender", sa.String(1), nullable=True),
        sa.Column("nationality", sa.String(3), nullable=True),
        sa.Column("charges", JSONB, nullable=False),
        sa.Column("issuing_country", sa.String(100), nullable=False),
        sa.Column("entering_agency", sa.String(100), nullable=False),
        sa.Column("status", sa.String(20), default="ACTIVE"),
        sa.CheckConstraint(
            "notice_type IN ('RED','BLUE','YELLOW','BLACK','ORANGE','PURPLE','UNKNOWN')",
            name="ck_fugitive_notice_type",
        ),
    )

    # ── PER-NID: Unidentified Persons ─────────────────────────────────────
    op.create_table(
        "per_unidentified_persons",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("discovery_date", sa.Date, nullable=False),
        sa.Column("discovery_location", sa.String(200), nullable=False),
        sa.Column("discovery_department", sa.String(50), nullable=True),
        sa.Column("estimated_age_min", sa.SmallInteger, nullable=True),
        sa.Column("estimated_age_max", sa.SmallInteger, nullable=True),
        sa.Column("gender", sa.String(1), nullable=True),
        sa.Column("estimated_height_cm", sa.SmallInteger, nullable=True),
        sa.Column("estimated_weight_kg", sa.SmallInteger, nullable=True),
        sa.Column("hair_color", sa.String(50), nullable=True),
        sa.Column("eye_color", sa.String(50), nullable=True),
        sa.Column("distinctive_features", sa.Text, nullable=True),
        sa.Column("clothing_description", sa.Text, nullable=True),
        sa.Column("dna_sample_ref", sa.String(36), nullable=True),
        sa.Column("fingerprint_ref", sa.String(36), nullable=True),
        sa.Column("dental_records_ref", sa.String(36), nullable=True),
        sa.Column("photo_refs", JSONB, default=list),
        sa.Column("entering_agency", sa.String(100), nullable=False),
        sa.Column("status", sa.String(20), default="ACTIVE"),
        sa.Column("matched_to_niu", sa.String(20), nullable=True),
        sa.Column("matched_date", sa.DateTime(timezone=True), nullable=True),
    )

    # ── PER-TER: Terrorism Watch ──────────────────────────────────────────
    op.create_table(
        "per_terrorism_watch",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("niu", sa.String(20), nullable=True),
        sa.Column("last_name", sa.String(100), nullable=False),
        sa.Column("first_name", sa.String(100), nullable=True),
        sa.Column("aliases", JSONB, default=list),
        sa.Column("date_of_birth", sa.Date, nullable=True),
        sa.Column("nationality", sa.String(3), nullable=True),
        sa.Column("risk_level", sa.String(10), default="HIGH"),
        sa.Column("threat_type", sa.String(100), nullable=False),
        sa.Column("groups", JSONB, default=list),
        sa.Column("known_associates", JSONB, default=list),
        sa.Column("last_known_location", sa.String(200), nullable=True),
        sa.Column("entering_agency", sa.String(100), nullable=False),
        sa.Column("approved_by_director", sa.String(100), nullable=False),
        sa.Column("approved_by_pg", sa.String(100), nullable=False),
        sa.Column("status", sa.String(20), default="ACTIVE"),
    )

    # ── PER-OPR: Protection Orders ────────────────────────────────────────
    op.create_table(
        "per_protection_orders",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("order_type", sa.String(30), nullable=False),
        sa.Column("issuing_court", sa.String(200), nullable=False),
        sa.Column("issuing_judge", sa.String(100), nullable=False),
        sa.Column("beneficiary_niu", sa.String(20), nullable=True),
        sa.Column("beneficiary_name", sa.String(100), nullable=False),
        sa.Column("protected_person", sa.String(100), nullable=False),
        sa.Column("restrained_person", sa.String(100), nullable=False),
        sa.Column("restrictions", JSONB, nullable=False),
        sa.Column("issue_date", sa.Date, nullable=False),
        sa.Column("expiration_date", sa.Date, nullable=True),
        sa.Column("emergency_contact", sa.String(100), nullable=True),
        sa.Column("status", sa.String(20), default="ACTIVE"),
        sa.CheckConstraint(
            "order_type IN ('RESTRAINING','HARASSMENT','CHILD_PROTECTION','DOMESTIC_VIOLENCE','EMERGENCY','OTHER')",
            name="ck_protection_order_type",
        ),
    )

    # ── PER-LIB: Supervised Release ────────────────────────────────────────
    op.create_table(
        "per_supervised_releases",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("niu", sa.String(20), nullable=False),
        sa.Column("last_name", sa.String(100), nullable=False),
        sa.Column("first_name", sa.String(100), nullable=True),
        sa.Column("supervision_type", sa.String(30), nullable=False),
        sa.Column("start_date", sa.Date, nullable=False),
        sa.Column("end_date", sa.Date, nullable=True),
        sa.Column("conditions", JSONB, nullable=False),
        sa.Column("supervising_officer", sa.String(100), nullable=False),
        sa.Column("supervising_agency", sa.String(100), nullable=False),
        sa.Column("status", sa.String(20), default="ACTIVE"),
        sa.CheckConstraint(
            "supervision_type IN ('CONDITIONAL_RELEASE','PAROLE','PROBATION','JUDICIAL_CONTROL','OTHER')",
            name="ck_supervision_type",
        ),
    )

    # ── PER-REC: Add entering_officer column ────────────────────────────────
    op.add_column("per_wanted_persons", sa.Column("entering_officer", sa.String(100), nullable=True))

    # ── PER-DIS: Add citizen portal, BPM, auto-cross columns ───────────────
    op.add_column("per_missing_persons", sa.Column("citizen_portal_submission", sa.Boolean(), default=False))
    op.add_column("per_missing_persons", sa.Column("bpm_notified", sa.Boolean(), default=False))
    op.add_column("per_missing_persons", sa.Column("auto_cross_bio_dis", sa.Boolean(), default=False))

    # ── PER-SEX: Add address_declared, geographic_restrictions ──────────────
    op.add_column("per_sex_offenders", sa.Column("address_declared", sa.Text, nullable=True))
    op.add_column("per_sex_offenders", sa.Column("geographic_restrictions", sa.Text, nullable=True))

    # ── PER-GNG: Add review/removal tracking columns ─────────────────────────
    op.add_column("per_gang_members", sa.Column("last_review_date", sa.Date, nullable=True))
    op.add_column("per_gang_members", sa.Column("auto_removal_date", sa.Date, nullable=True))


def downgrade() -> None:
    op.drop_table("per_supervised_releases")
    op.drop_table("per_protection_orders")
    op.drop_table("per_terrorism_watch")
    op.drop_table("per_unidentified_persons")
    op.drop_table("per_foreign_fugitives")
    op.drop_column("per_gang_members", "auto_removal_date")
    op.drop_column("per_gang_members", "last_review_date")
    op.drop_column("per_sex_offenders", "geographic_restrictions")
    op.drop_column("per_sex_offenders", "address_declared")
    op.drop_column("per_missing_persons", "auto_cross_bio_dis")
    op.drop_column("per_missing_persons", "bpm_notified")
    op.drop_column("per_missing_persons", "citizen_portal_submission")
    op.drop_column("per_wanted_persons", "entering_officer")
