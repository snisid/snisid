"""create bio_identity_links table (dissociation ADN/identité)

Revision ID: 0011
Revises: 0010
Create Date: 2026-06-11

"""
from __future__ import annotations

import sqlalchemy as sa
from alembic import op

revision: str = "0011"
down_revision: str | None = "0010"
branch_labels: str | None = None
depends_on: str | None = None


def upgrade() -> None:
    op.create_table(
        "bio_identity_links",
        sa.Column("sample_id", sa.String(36), nullable=False),
        sa.Column("niu", sa.String(20), nullable=False),
        sa.Column("linked_by", sa.String(100), nullable=False),
        sa.Column("linked_at", sa.DateTime(timezone=True), server_default=sa.text("NOW()"), nullable=False),
        sa.Column("court_order", sa.String(200), nullable=True),
        sa.ForeignKeyConstraint(["sample_id"], ["bio_str_profiles.sample_id"], ondelete="CASCADE"),
        sa.PrimaryKeyConstraint("sample_id"),
    )
    op.create_index("idx_bio_identity_links_niu", "bio_identity_links", ["niu"])


def downgrade() -> None:
    op.drop_table("bio_identity_links")
