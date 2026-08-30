"""create bie obj and security tables

Revision ID: 0005
Revises: 0004
Create Date: 2026-06-11
"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import JSONB

revision: str = "0005"
down_revision: Union[str, None] = "0004"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # ── BIE-EMB: Add new columns to vessels ────────────────────────────────
    op.add_column("bie_stolen_vessels", sa.Column("vessel_length_m", sa.Float(), nullable=True))
    op.add_column("bie_stolen_vessels", sa.Column("hull_color", sa.String(50), nullable=True))
    op.add_column("bie_stolen_vessels", sa.Column("home_port", sa.String(200), nullable=True))
    op.add_column("bie_stolen_vessels", sa.Column("engine_serial", sa.String(100), nullable=True))
    op.add_column("bie_stolen_vessels", sa.Column("distinctive_marks", sa.Text, nullable=True))

    # ── BIE-OBJ: Stolen Articles (bijoux, art, bétail, etc.) ───────────────
    op.create_table(
        "bie_stolen_articles",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("category", sa.String(30), nullable=False),
        sa.Column("description", sa.Text, nullable=False),
        sa.Column("serial_number", sa.String(100), nullable=True),
        sa.Column("estimated_value", sa.Numeric(12, 2), nullable=True),
        sa.Column("currency_code", sa.String(3), default="HTG"),
        sa.Column("theft_date", sa.Date, nullable=False),
        sa.Column("theft_location", sa.String(200), nullable=False),
        sa.Column("theft_department", sa.String(50), nullable=True),
        sa.Column("owner_niu", sa.String(20), nullable=True),
        sa.Column("owner_name", sa.String(200), nullable=True),
        sa.Column("status", sa.String(20), default="STOLEN"),
        sa.Column("recovered_date", sa.Date, nullable=True),
        sa.Column("entering_agency", sa.String(100), nullable=False),
        sa.CheckConstraint(
            "category IN ('JEWELRY','ART','ELECTRONICS','CURRENCY','CATTLE','MACHINERY','OTHER')",
            name="ck_article_category",
        ),
    )

    # ── BIE-TIT: Stolen Securities (chèques, obligations, titres) ──────────
    op.create_table(
        "bie_stolen_securities",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("record_number", sa.String(50), unique=True, nullable=False),
        sa.Column("security_type", sa.String(30), nullable=False),
        sa.Column("issuer", sa.String(200), nullable=False),
        sa.Column("security_number", sa.String(100), nullable=False),
        sa.Column("face_value", sa.Numeric(12, 2), nullable=True),
        sa.Column("currency_code", sa.String(3), default="HTG"),
        sa.Column("issue_date", sa.Date, nullable=True),
        sa.Column("theft_date", sa.Date, nullable=False),
        sa.Column("theft_location", sa.String(200), nullable=False),
        sa.Column("owner_niu", sa.String(20), nullable=True),
        sa.Column("owner_name", sa.String(200), nullable=True),
        sa.Column("status", sa.String(20), default="STOLEN"),
        sa.Column("recovered_date", sa.Date, nullable=True),
        sa.Column("entering_agency", sa.String(100), nullable=False),
        sa.CheckConstraint(
            "security_type IN ('CHEQUE','BOND','PROPERTY_TITLE','LETTER_CREDIT','OTHER')",
            name="ck_security_type",
        ),
    )

    # ── BIE-VEH: Add recovered_location column ─────────────────────────────
    op.add_column("bie_stolen_vehicles", sa.Column("recovered_location", sa.String(200), nullable=True))


def downgrade() -> None:
    op.drop_table("bie_stolen_securities")
    op.drop_table("bie_stolen_articles")
    op.drop_column("bie_stolen_vessels", "distinctive_marks")
    op.drop_column("bie_stolen_vessels", "engine_serial")
    op.drop_column("bie_stolen_vessels", "home_port")
    op.drop_column("bie_stolen_vessels", "hull_color")
    op.drop_column("bie_stolen_vessels", "vessel_length_m")
    op.drop_column("bie_stolen_vehicles", "recovered_location")
