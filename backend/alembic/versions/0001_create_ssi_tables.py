"""create ssi tables

Revision ID: 0001
Revises:
Create Date: 2026-06-09
"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import JSONB

revision: str = "0001"
down_revision: Union[str, None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "ssi_did",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("did", sa.String(256), unique=True, nullable=False, index=True),
        sa.Column("method", sa.String(32), nullable=False),
        sa.Column("document", JSONB(), nullable=False),
    )
    op.create_table(
        "ssi_vc",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("credential_id", sa.String(128), unique=True, nullable=False, index=True),
        sa.Column("issuer_id", sa.String(256), nullable=False, index=True),
        sa.Column("subject_id", sa.String(256), nullable=False, index=True),
        sa.Column("credential_type", sa.String(128), nullable=False),
        sa.Column("document", JSONB(), nullable=False),
        sa.Column("status_list_id", sa.String(128), nullable=True),
        sa.Column("revoked", sa.Boolean(), default=False),
        sa.Column("issued_at", sa.DateTime(timezone=True), nullable=False),
    )
    op.create_table(
        "ssi_status_list",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("list_id", sa.String(128), unique=True, nullable=False, index=True),
        sa.Column("purpose", sa.String(64), nullable=False),
        sa.Column("bitstring", sa.Text(), nullable=False),
    )
    op.create_table(
        "ssi_wallet_credential",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("wallet_did", sa.String(256), nullable=False, index=True),
        sa.Column("credential_id", sa.String(128), nullable=False),
        sa.Column("document", JSONB(), nullable=False),
        sa.Column("issuer_id", sa.String(256), nullable=False),
        sa.Column("credential_type", sa.String(128), nullable=False),
        sa.Column("issued_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("stored_at", sa.DateTime(timezone=True), nullable=False),
        sa.UniqueConstraint("wallet_did", "credential_id", name="uq_wallet_credential"),
    )
    op.create_table(
        "ssi_didcomm_message",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("message_id", sa.String(128), unique=True, nullable=False, index=True),
        sa.Column("sender_did", sa.String(256), nullable=False),
        sa.Column("receiver_did", sa.String(256), nullable=False, index=True),
        sa.Column("message_type", sa.String(128), nullable=False),
        sa.Column("message_body", JSONB(), nullable=False),
        sa.Column("thread_id", sa.String(128), nullable=True),
        sa.Column("is_read", sa.Boolean(), default=False),
    )
    op.create_table(
        "ssi_credential_flow",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("flow_id", sa.String(128), unique=True, nullable=False, index=True),
        sa.Column("issuer_id", sa.String(256), nullable=False),
        sa.Column("offer_data", JSONB(), nullable=False),
        sa.Column("request_data", JSONB(), nullable=True),
        sa.Column("credential_id", sa.String(128), nullable=True),
        sa.Column("status", sa.String(32), nullable=False),
    )
    op.create_table(
        "ssi_chapi",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("holder_did", sa.String(256), nullable=False, index=True),
        sa.Column("credential_id", sa.String(128), nullable=False),
        sa.Column("document", JSONB(), nullable=False),
        sa.Column("query_frame", JSONB(), nullable=True),
        sa.Column("stored_at", sa.DateTime(timezone=True), nullable=False),
        sa.UniqueConstraint("holder_did", "credential_id", name="uq_chapi_credential"),
    )
    op.create_table(
        "ssi_credential_manifest",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("manifest_id", sa.String(128), unique=True, nullable=False, index=True),
        sa.Column("issuer_id", sa.String(256), nullable=False),
        sa.Column("document", JSONB(), nullable=False),
    )
    op.create_table(
        "ssi_revocation_event",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("event_id", sa.String(128), unique=True, nullable=False, index=True),
        sa.Column("credential_id", sa.String(128), nullable=False, index=True),
        sa.Column("subject_id", sa.String(256), nullable=False, index=True),
        sa.Column("event_type", sa.String(64), nullable=False),
        sa.Column("reason", sa.Text(), nullable=True),
    )


def downgrade() -> None:
    op.drop_table("ssi_revocation_event")
    op.drop_table("ssi_credential_manifest")
    op.drop_table("ssi_chapi")
    op.drop_table("ssi_credential_flow")
    op.drop_table("ssi_didcomm_message")
    op.drop_table("ssi_wallet_credential")
    op.drop_table("ssi_status_list")
    op.drop_table("ssi_vc")
    op.drop_table("ssi_did")
