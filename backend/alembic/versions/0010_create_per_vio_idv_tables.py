"""create per_violence_records and per_identity_thefts tables

Revision ID: 0010
Revises: 0009
Create Date: 2026-06-11

"""
from __future__ import annotations

from typing import ClassVar

import sqlalchemy as sa
from alembic import op

revision: str = "0010"
down_revision: str | None = "0009"
branch_labels: str | None = None
depends_on: str | None = None


def upgrade() -> None:
    op.create_table(
        "per_violence_records",
        sa.Column("id", sa.UUID(), server_default=sa.text("gen_random_uuid()"), nullable=False),
        sa.Column("record_number", sa.String(30), nullable=False),
        sa.Column("niu", sa.String(20), nullable=True),
        sa.Column("last_name", sa.String(100), nullable=True),
        sa.Column("first_name", sa.String(100), nullable=True),
        sa.Column("incident_type", sa.String(30), nullable=False),
        sa.Column("incident_date", sa.String(20), nullable=False),
        sa.Column("location", sa.String(200), nullable=False),
        sa.Column("victim_niu", sa.String(20), nullable=True),
        sa.Column("victim_name", sa.String(200), nullable=True),
        sa.Column("arresting_agency", sa.String(50), nullable=False),
        sa.Column("court_case_ref", sa.String(50), nullable=True),
        sa.Column("risk_level", sa.String(10), nullable=False, server_default="MEDIUM"),
        sa.Column("status", sa.String(20), nullable=False, server_default="ACTIVE"),
        sa.Column("created_at", sa.DateTime(timezone=True), server_default=sa.text("NOW()"), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), server_default=sa.text("NOW()"), nullable=False),
        sa.PrimaryKeyConstraint("id"),
        sa.UniqueConstraint("record_number"),
    )
    op.create_index("idx_per_violence_niu", "per_violence_records", ["niu"])
    op.create_index("idx_per_violence_incident_type", "per_violence_records", ["incident_type"])
    op.create_index("idx_per_violence_status", "per_violence_records", ["status"])

    op.create_table(
        "per_identity_thefts",
        sa.Column("id", sa.UUID(), server_default=sa.text("gen_random_uuid()"), nullable=False),
        sa.Column("record_number", sa.String(30), nullable=False),
        sa.Column("victim_niu", sa.String(20), nullable=False),
        sa.Column("victim_name", sa.String(200), nullable=True),
        sa.Column("fraud_type", sa.String(30), nullable=False),
        sa.Column("document_type_used", sa.String(50), nullable=True),
        sa.Column("perpetrator_known", sa.Boolean(), nullable=False, server_default=sa.text("false")),
        sa.Column("perpetrator_name", sa.String(200), nullable=True),
        sa.Column("report_date", sa.String(20), nullable=False),
        sa.Column("reporting_agency", sa.String(50), nullable=False),
        sa.Column("status", sa.String(20), nullable=False, server_default="ACTIVE"),
        sa.Column("created_at", sa.DateTime(timezone=True), server_default=sa.text("NOW()"), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), server_default=sa.text("NOW()"), nullable=False),
        sa.PrimaryKeyConstraint("id"),
        sa.UniqueConstraint("record_number"),
    )
    op.create_index("idx_per_identitytheft_victim_niu", "per_identity_thefts", ["victim_niu"])
    op.create_index("idx_per_identitytheft_fraud_type", "per_identity_thefts", ["fraud_type"])
    op.create_index("idx_per_identitytheft_status", "per_identity_thefts", ["status"])


def downgrade() -> None:
    op.drop_table("per_identity_thefts")
    op.drop_table("per_violence_records")
