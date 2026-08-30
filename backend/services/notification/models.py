from __future__ import annotations

import enum
from datetime import datetime, timezone

from sqlalchemy import Boolean, DateTime, Enum, Integer, String, Text
from sqlalchemy.orm import Mapped, mapped_column

from shared.database import Base


class WebhookStatus(str, enum.Enum):
    ACTIVE = "active"
    PAUSED = "paused"
    FAILED = "failed"


class WebhookEvent(str, enum.Enum):
    IDENTITY_CREATED = "identity.created"
    IDENTITY_UPDATED = "identity.updated"
    IDENTITY_VERIFIED = "identity.verified"
    IDENTITY_SUSPENDED = "identity.suspended"
    IDENTITY_REVOKED = "identity.revoked"
    BIOMETRIC_ENROLLED = "biometric.enrolled"
    DOCUMENT_ISSUED = "document.issued"


class WebhookSubscription(Base):
    __tablename__ = "webhook_subscriptions"

    url: Mapped[str] = mapped_column(String(500), nullable=False)
    events: Mapped[str] = mapped_column(
        Text, nullable=False, doc="Comma-separated list of event types"
    )
    status: Mapped[WebhookStatus] = mapped_column(
        Enum(WebhookStatus, name="webhook_status_enum"),
        nullable=False,
        default=WebhookStatus.ACTIVE,
    )
    secret: Mapped[str | None] = mapped_column(
        String(128), nullable=True, doc="HMAC secret for payload signing"
    )
    retry_count: Mapped[int] = mapped_column(Integer, default=3)
    timeout_seconds: Mapped[int] = mapped_column(Integer, default=10)
    last_triggered_at: Mapped[datetime | None] = mapped_column(
        DateTime(timezone=True), nullable=True
    )
    last_failure_reason: Mapped[str | None] = mapped_column(Text, nullable=True)
    consecutive_failures: Mapped[int] = mapped_column(Integer, default=0)
    max_consecutive_failures: Mapped[int] = mapped_column(Integer, default=10)


class WebhookDeliveryLog(Base):
    __tablename__ = "webhook_delivery_logs"

    subscription_id: Mapped[str] = mapped_column(String(36), nullable=False, index=True)
    event_type: Mapped[str] = mapped_column(String(100), nullable=False)
    payload: Mapped[str] = mapped_column(Text, nullable=False)
    status_code: Mapped[int | None] = mapped_column(nullable=True)
    response_body: Mapped[str | None] = mapped_column(Text, nullable=True)
    success: Mapped[bool] = mapped_column(default=False)
    duration_ms: Mapped[int | None] = mapped_column(nullable=True)
    delivered_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), default=lambda: datetime.now(timezone.utc)
    )
