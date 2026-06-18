"""create webhook subscription and delivery log tables

Revision ID: 0002
Revises: 0001
Create Date: 2026-06-10
"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa

revision: str = "0002"
down_revision: Union[str, None] = "0001"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    webhook_status_enum = sa.Enum("active", "paused", "failed", name="webhook_status_enum")
    webhook_status_enum.create(op.get_bind(), checkfirst=True)

    webhook_event_enum = sa.Enum(
        "identity.created", "identity.updated", "identity.verified",
        "identity.suspended", "identity.revoked",
        "biometric.enrolled", "document.issued",
        name="webhook_event_enum",
    )
    webhook_event_enum.create(op.get_bind(), checkfirst=True)

    op.create_table(
        "webhook_subscriptions",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now(), onupdate=sa.func.now()),
        sa.Column("is_deleted", sa.Boolean(), default=False, index=True),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("url", sa.String(500), nullable=False),
        sa.Column("events", sa.Text(), nullable=False,
                  doc="Comma-separated list of event types"),
        sa.Column("status", webhook_status_enum, nullable=False, default="active"),
        sa.Column("secret", sa.String(128), nullable=True,
                  doc="HMAC secret for payload signing"),
        sa.Column("retry_count", sa.Integer(), default=3),
        sa.Column("timeout_seconds", sa.Integer(), default=10),
        sa.Column("last_triggered_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("last_failure_reason", sa.Text(), nullable=True),
        sa.Column("consecutive_failures", sa.Integer(), default=0),
        sa.Column("max_consecutive_failures", sa.Integer(), default=10),
    )
    op.create_table(
        "webhook_delivery_logs",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False,
                  server_default=sa.func.now(), onupdate=sa.func.now()),
        sa.Column("is_deleted", sa.Boolean(), default=False),
        sa.Column("deleted_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("subscription_id", sa.String(36), nullable=False, index=True),
        sa.Column("event_type", sa.String(100), nullable=False),
        sa.Column("payload", sa.Text(), nullable=False),
        sa.Column("status_code", sa.Integer(), nullable=True),
        sa.Column("response_body", sa.Text(), nullable=True),
        sa.Column("error_message", sa.Text(), nullable=True),
        sa.Column("duration_ms", sa.Integer(), nullable=True),
        sa.Column("delivered", sa.Boolean(), default=False),
    )


def downgrade() -> None:
    op.drop_table("webhook_delivery_logs")
    op.drop_table("webhook_subscriptions")
    sa.Enum(name="webhook_event_enum").drop(op.get_bind(), checkfirst=True)
    sa.Enum(name="webhook_status_enum").drop(op.get_bind(), checkfirst=True)
