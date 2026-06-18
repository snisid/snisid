from __future__ import annotations

import uuid
from datetime import datetime, timezone
from typing import Any

from pydantic import BaseModel, Field
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from services.agency.models import Agency, AgencyStatus, AgencyType
from shared.logging import get_logger

logger = get_logger(__name__)


class CreateAgencyCommand(BaseModel):
    command_id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    name: str = Field(..., min_length=1, max_length=200)
    code: str = Field(..., min_length=2, max_length=20)
    agency_type: AgencyType
    address: str | None = None
    city: str | None = None
    department: str | None = None
    phone: str | None = None
    email: str | None = None
    latitude: float | None = None
    longitude: float | None = None
    is_headquarters: bool = False
    parent_agency_id: str | None = None
    actor_id: str = "system"


class UpdateAgencyCommand(BaseModel):
    command_id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    agency_id: str
    changes: dict[str, Any]
    actor_id: str = "system"


class DeactivateAgencyCommand(BaseModel):
    command_id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    agency_id: str
    reason: str = Field(..., min_length=10)
    actor_id: str = "system"


class AgencyNotFoundError(Exception):
    pass


class AgencyCommandHandler:
    def __init__(self, session: AsyncSession) -> None:
        self._session = session

    async def handle_create(self, cmd: CreateAgencyCommand) -> dict[str, Any]:
        agency = Agency(
            name=cmd.name,
            code=cmd.code,
            agency_type=cmd.agency_type,
            address=cmd.address,
            city=cmd.city,
            department=cmd.department,
            phone=cmd.phone,
            email=cmd.email,
            latitude=cmd.latitude,
            longitude=cmd.longitude,
            is_headquarters=cmd.is_headquarters,
            parent_agency_id=cmd.parent_agency_id,
        )
        self._session.add(agency)
        await self._session.flush()
        logger.info("agency_created", agency_id=agency.id, code=agency.code)
        return {"agency_id": agency.id, "code": agency.code, "status": "active"}

    async def handle_update(self, cmd: UpdateAgencyCommand) -> dict[str, Any]:
        result = await self._session.execute(
            select(Agency).where(Agency.id == cmd.agency_id, Agency.is_deleted == False)
        )
        agency = result.scalar_one_or_none()
        if not agency:
            raise AgencyNotFoundError(f"Agency {cmd.agency_id} not found")
        for key, value in cmd.changes.items():
            if hasattr(agency, key):
                setattr(agency, key, value)
        await self._session.flush()
        logger.info("agency_updated", agency_id=agency.id)
        return {"agency_id": agency.id}

    async def handle_deactivate(self, cmd: DeactivateAgencyCommand) -> dict[str, Any]:
        result = await self._session.execute(
            select(Agency).where(Agency.id == cmd.agency_id, Agency.is_deleted == False)
        )
        agency = result.scalar_one_or_none()
        if not agency:
            raise AgencyNotFoundError(f"Agency {cmd.agency_id} not found")
        agency.status = AgencyStatus.INACTIVE
        agency.closed_at = datetime.now(timezone.utc)
        await self._session.flush()
        logger.info("agency_deactivated", agency_id=agency.id, reason=cmd.reason)
        return {"agency_id": agency.id, "status": "inactive"}
