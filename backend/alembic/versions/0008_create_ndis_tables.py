"""create ndis tables for cross-dept hits, reports, interpol submissions

Revision ID: 0008
Revises: 0007
Create Date: 2026-06-11
"""
from __future__ import annotations

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import JSONB

revision: str = "0008"
down_revision: Union[str, None] = "0007"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # ── NDIS Cross-Department Hits ──────────────────────────────────────────
    op.create_table(
        "ndis_cross_dept_hits",
        sa.Column("id", sa.String(50), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("query_sample_id", sa.String(36), nullable=False),
        sa.Column("match_sample_id", sa.String(36), nullable=False),
        sa.Column("match_type", sa.String(20), nullable=False),
        sa.Column("confidence", sa.Numeric(5, 4), nullable=False),
        sa.Column("query_sdis", sa.String(20), nullable=False),
        sa.Column("match_sdis", sa.String(20), nullable=False),
        sa.Column("alert_level", sa.String(10), default="HIGH"),
        sa.Column("notified_at", sa.DateTime(timezone=True), nullable=True),
    )

    # ── NDIS Reports ────────────────────────────────────────────────────────
    op.create_table(
        "ndis_reports",
        sa.Column("id", sa.String(50), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("report_type", sa.String(20), nullable=False),
        sa.Column("status", sa.String(20), default="GENERATED"),
        sa.Column("file_path", sa.String(500), nullable=True),
        sa.CheckConstraint(
            "report_type IN ('STATS','HITS','UNMATCHED','QUALITY','INTERPOL')",
            name="ck_report_type",
        ),
    )

    # ── INTERPOL Submissions ────────────────────────────────────────────────
    op.create_table(
        "ndis_interpol_submissions",
        sa.Column("id", sa.String(50), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("sample_ids", JSONB, nullable=False),
        sa.Column("reason", sa.String(50), nullable=False),
        sa.Column("case_number", sa.String(100), nullable=True),
        sa.Column("status", sa.String(20), default="PENDING"),
        sa.CheckConstraint(
            "reason IN ('disaster_victim','international_fugitive','trafficking_victim','unidentified')",
            name="ck_interpol_reason",
        ),
    )


def downgrade() -> None:
    op.drop_table("ndis_interpol_submissions")
    op.drop_table("ndis_reports")
    op.drop_table("ndis_cross_dept_hits")
