from __future__ import annotations

import enum
from datetime import datetime, timezone

from sqlalchemy import Boolean, DateTime, Enum, Float, String, Text
from sqlalchemy.orm import Mapped, mapped_column

from shared.database import Base


class AgencyStatus(str, enum.Enum):
    ACTIVE = "active"
    INACTIVE = "inactive"
    SUSPENDED = "suspended"


class AgencyType(str, enum.Enum):
    CENTRAL = "central"
    REGIONAL = "regional"
    LOCAL = "local"
    MOBILE = "mobile"


class Agency(Base):
    __tablename__ = "agencies"

    name: Mapped[str] = mapped_column(String(200), nullable=False)
    code: Mapped[str] = mapped_column(String(20), unique=True, nullable=False, index=True)
    agency_type: Mapped[AgencyType] = mapped_column(
        Enum(AgencyType, name="agency_type_enum"), nullable=False
    )
    status: Mapped[AgencyStatus] = mapped_column(
        Enum(AgencyStatus, name="agency_status_enum"),
        nullable=False,
        default=AgencyStatus.ACTIVE,
        index=True,
    )
    address: Mapped[str | None] = mapped_column(Text, nullable=True)
    city: Mapped[str | None] = mapped_column(String(100), nullable=True)
    department: Mapped[str | None] = mapped_column(String(100), nullable=True)
    phone: Mapped[str | None] = mapped_column(String(20), nullable=True)
    email: Mapped[str | None] = mapped_column(String(255), nullable=True)
    latitude: Mapped[float | None] = mapped_column(Float, nullable=True)
    longitude: Mapped[float | None] = mapped_column(Float, nullable=True)
    max_daily_enrollments: Mapped[int] = mapped_column(default=500)
    is_headquarters: Mapped[bool] = mapped_column(Boolean, default=False)
    parent_agency_id: Mapped[str | None] = mapped_column(String(36), nullable=True)
    opened_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), default=lambda: datetime.now(timezone.utc)
    )
    closed_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)
