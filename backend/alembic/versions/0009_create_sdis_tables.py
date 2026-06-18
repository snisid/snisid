"""create sdis tables: nodes, matches, sync errors, quality reviews

Revision ID: 0009
Revises: 0008
Create Date: 2026-06-11
"""
from __future__ import annotations

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import JSONB

revision: str = "0009"
down_revision: Union[str, None] = "0008"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # ── SDIS Nodes ────────────────────────────────────────────────────────────
    op.create_table(
        "sdis_nodes",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("code", sa.String(20), unique=True, nullable=False),
        sa.Column("department", sa.String(50), nullable=False),
        sa.Column("dc_location", sa.String(100), nullable=False),
        sa.Column("dc_type", sa.String(20), nullable=False),
        sa.Column("lab_codes", JSONB, default=list),
        sa.Column("status", sa.String(20), default="ACTIVE"),
        sa.Column("last_heartbeat", sa.DateTime(timezone=True), nullable=True),
        sa.Column("is_primary", sa.Boolean(), default=True),
    )

    # ── Intra-departmental matches ──────────────────────────────────────────
    op.create_table(
        "sdis_matches",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("sdis_code", sa.String(20), nullable=False),
        sa.Column("query_sample_id", sa.String(36), nullable=False),
        sa.Column("match_sample_id", sa.String(36), nullable=False),
        sa.Column("match_type", sa.String(20), nullable=False),
        sa.Column("confidence", sa.Numeric(5, 4), nullable=False),
        sa.Column("alerted", sa.Boolean(), default=False),
        sa.Column("alerted_at", sa.DateTime(timezone=True), nullable=True),
    )

    # ── Sync errors log ─────────────────────────────────────────────────────
    op.create_table(
        "sdis_sync_errors",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("sdis_code", sa.String(20), nullable=False),
        sa.Column("error_type", sa.String(30), nullable=False),
        sa.Column("details", sa.Text, nullable=True),
        sa.Column("retry_count", sa.Integer(), default=0),
        sa.Column("resolved", sa.Boolean(), default=False),
        sa.Column("resolved_at", sa.DateTime(timezone=True), nullable=True),
    )

    # ── Quality review queue ────────────────────────────────────────────────
    op.create_table(
        "sdis_quality_reviews",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("sample_id", sa.String(36), nullable=False),
        sa.Column("sdis_code", sa.String(20), nullable=False),
        sa.Column("quality_score", sa.Numeric(4, 3), nullable=False),
        sa.Column("reason", sa.String(100), nullable=False),
        sa.Column("reviewed", sa.Boolean(), default=False),
        sa.Column("reviewed_by", sa.String(100), nullable=True),
        sa.Column("reviewed_at", sa.DateTime(timezone=True), nullable=True),
    )

    # ── Seed 10 SDIS nodes ──────────────────────────────────────────────────
    nodes_table = sa.table(
        "sdis_nodes",
        sa.column("id", sa.String),
        sa.column("code", sa.String),
        sa.column("department", sa.String),
        sa.column("dc_location", sa.String),
        sa.column("dc_type", sa.String),
        sa.column("lab_codes", JSONB),
        sa.column("status", sa.String),
        sa.column("is_primary", sa.Boolean),
    )
    op.bulk_insert(nodes_table, [
        {"id": "b001", "code": "SDIS-OUEST", "department": "Ouest", "dc_location": "SNISID PAP", "dc_type": "DC principal", "lab_codes": ["LDIS-PAP-001"], "status": "ACTIVE", "is_primary": True},
        {"id": "b002", "code": "SDIS-NORD", "department": "Nord", "dc_location": "SNISID CAP", "dc_type": "DC secondaire", "lab_codes": ["LDIS-CAP-001"], "status": "ACTIVE", "is_primary": True},
        {"id": "b003", "code": "SDIS-ARTIBONITE", "department": "Artibonite", "dc_location": "Gonaïves", "dc_type": "Noeud SDIS", "lab_codes": ["LDIS-GON-001"], "status": "ACTIVE", "is_primary": True},
        {"id": "b004", "code": "SDIS-SUD", "department": "Sud", "dc_location": "Les Cayes", "dc_type": "Noeud SDIS", "lab_codes": ["LDIS-LES-001"], "status": "ACTIVE", "is_primary": True},
        {"id": "b005", "code": "SDIS-SUDEST", "department": "Sud-Est", "dc_location": "Jacmel", "dc_type": "Noeud SDIS", "lab_codes": ["LDIS-JAC-001"], "status": "ACTIVE", "is_primary": True},
        {"id": "b006", "code": "SDIS-CENTRE", "department": "Centre", "dc_location": "Hinche", "dc_type": "Noeud SDIS", "lab_codes": ["LDIS-HIN-001"], "status": "ACTIVE", "is_primary": True},
        {"id": "b007", "code": "SDIS-NORDEST", "department": "Nord-Est", "dc_location": "Fort-Liberté", "dc_type": "Noeud SDIS", "lab_codes": [], "status": "INACTIVE", "is_primary": False},
        {"id": "b008", "code": "SDIS-NORDOUEST", "department": "Nord-Ouest", "dc_location": "Port-de-Paix", "dc_type": "Noeud SDIS", "lab_codes": [], "status": "INACTIVE", "is_primary": False},
        {"id": "b009", "code": "SDIS-GRANDANSE", "department": "Grand-Anse", "dc_location": "Jérémie", "dc_type": "Noeud SDIS", "lab_codes": [], "status": "INACTIVE", "is_primary": False},
        {"id": "b010", "code": "SDIS-NIPPES", "department": "Nippes", "dc_location": "Miragoâne", "dc_type": "Noeud SDIS", "lab_codes": [], "status": "INACTIVE", "is_primary": False},
    ])


def downgrade() -> None:
    op.drop_table("sdis_quality_reviews")
    op.drop_table("sdis_sync_errors")
    op.drop_table("sdis_matches")
    op.drop_table("sdis_nodes")
