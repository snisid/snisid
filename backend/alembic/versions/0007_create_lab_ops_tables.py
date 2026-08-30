"""create lab ops tables, seed labs, add accreditation/equipment/training

Revision ID: 0007
Revises: 0006
Create Date: 2026-06-11
"""
from __future__ import annotations

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import JSONB

revision: str = "0007"
down_revision: Union[str, None] = "0006"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # ── BioLaboratory: new accreditation columns ────────────────────────────
    op.add_column("bio_laboratories", sa.Column("accreditation_body", sa.String(100), nullable=True))
    op.add_column("bio_laboratories", sa.Column("accreditation_expiry", sa.Date, nullable=True))
    op.add_column("bio_laboratories", sa.Column("last_external_audit", sa.Date, nullable=True))
    op.add_column("bio_laboratories", sa.Column("external_quality_check_date", sa.Date, nullable=True))

    # ── BioLabEquipment ─────────────────────────────────────────────────────
    op.create_table(
        "bio_lab_equipment",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("lab_code", sa.String(20), sa.ForeignKey("bio_laboratories.id"), nullable=False),
        sa.Column("equipment_name", sa.String(200), nullable=False),
        sa.Column("model", sa.String(200), nullable=True),
        sa.Column("serial_number", sa.String(100), nullable=False),
        sa.Column("role", sa.String(100), nullable=False),
        sa.Column("calibration_date", sa.Date, nullable=True),
        sa.Column("calibration_due", sa.Date, nullable=True),
        sa.Column("status", sa.String(20), default="ACTIVE"),
    )

    # ── BioStaffTraining ────────────────────────────────────────────────────
    op.create_table(
        "bio_staff_training",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("staff_niu", sa.String(20), nullable=False),
        sa.Column("training_name", sa.String(200), nullable=False),
        sa.Column("training_code", sa.String(50), nullable=False),
        sa.Column("duration_hours", sa.SmallInteger, nullable=True),
        sa.Column("completed_date", sa.Date, nullable=False),
        sa.Column("valid_until", sa.Date, nullable=True),
        sa.Column("issued_by", sa.String(100), nullable=False),
        sa.Column("frequency", sa.String(20), nullable=True),
    )

    # ── Seed LDIS laboratories ──────────────────────────────────────────────
    labs_table = sa.table(
        "bio_laboratories",
        sa.column("id", sa.String),
        sa.column("lab_code", sa.String),
        sa.column("lab_name", sa.String),
        sa.column("lab_level", sa.String),
        sa.column("department", sa.String),
        sa.column("institution", sa.String),
        sa.column("is_active", sa.Boolean),
    )
    op.bulk_insert(labs_table, [
        {"id": "a001", "lab_code": "LDIS-PAP-001", "lab_name": "Labo Médico-Légal PAP", "lab_level": "LDIS",
         "department": "Ouest", "institution": "DCPJ / MSPP", "is_active": True},
        {"id": "a002", "lab_code": "LDIS-CAP-001", "lab_name": "Labo Médico-Légal Cap-Haïtien", "lab_level": "LDIS",
         "department": "Nord", "institution": "PNH / MSPP", "is_active": True},
        {"id": "a003", "lab_code": "LDIS-LES-001", "lab_name": "Labo Médico-Légal Les Cayes", "lab_level": "LDIS",
         "department": "Sud", "institution": "MSPP", "is_active": False},
        {"id": "a004", "lab_code": "LDIS-GON-001", "lab_name": "Labo Médico-Légal Gonaïves", "lab_level": "LDIS",
         "department": "Artibonite", "institution": "MSPP", "is_active": False},
        {"id": "a005", "lab_code": "LDIS-JAC-001", "lab_name": "Labo Médico-Légal Jacmel", "lab_level": "LDIS",
         "department": "Sud-Est", "institution": "MSPP", "is_active": False},
        {"id": "a006", "lab_code": "LDIS-HIN-001", "lab_name": "Labo Médico-Légal Hinche", "lab_level": "LDIS",
         "department": "Centre", "institution": "MSPP", "is_active": False},
    ])


def downgrade() -> None:
    op.drop_table("bio_staff_training")
    op.drop_table("bio_lab_equipment")
    op.drop_column("bio_laboratories", "external_quality_check_date")
    op.drop_column("bio_laboratories", "last_external_audit")
    op.drop_column("bio_laboratories", "accreditation_expiry")
    op.drop_column("bio_laboratories", "accreditation_body")
    op.execute("DELETE FROM bio_laboratories WHERE lab_code LIKE 'LDIS-%'")
