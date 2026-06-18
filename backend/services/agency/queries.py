from __future__ import annotations

from typing import Any

from pydantic import BaseModel, Field
from sqlalchemy import func, or_, select
from sqlalchemy.ext.asyncio import AsyncSession

from services.agency.models import Agency, AgencyStatus
from shared.logging import get_logger

logger = get_logger(__name__)


class GetAgencyByIdQuery(BaseModel):
    agency_id: str


class ListAgenciesQuery(BaseModel):
    status: str | None = None
    agency_type: str | None = None
    department: str | None = None
    search_term: str | None = None
    page: int = Field(default=1, ge=1)
    page_size: int = Field(default=20, ge=1, le=100)


class AgencyResponse(BaseModel):
    id: str
    name: str
    code: str
    agency_type: str
    status: str
    address: str | None = None
    city: str | None = None
    department: str | None = None
    phone: str | None = None
    email: str | None = None
    is_headquarters: bool = False
    parent_agency_id: str | None = None
    max_daily_enrollments: int = 500
    opened_at: str | None = None
    closed_at: str | None = None


class PaginatedAgencyResponse(BaseModel):
    items: list[AgencyResponse]
    total: int
    page: int
    page_size: int
    total_pages: int


class AgencyQueryHandler:
    def __init__(self, session: AsyncSession) -> None:
        self._session = session

    async def get_by_id(self, query: GetAgencyByIdQuery) -> AgencyResponse:
        result = await self._session.execute(
            select(Agency).where(Agency.id == query.agency_id, Agency.is_deleted == False)
        )
        agency = result.scalar_one_or_none()
        if not agency:
            raise AgencyNotFoundError(f"Agency {query.agency_id} not found")
        return self._to_response(agency)

    async def list_all(self, query: ListAgenciesQuery) -> PaginatedAgencyResponse:
        stmt = select(Agency).where(Agency.is_deleted == False)
        if query.status:
            stmt = stmt.where(Agency.status == query.status)
        if query.agency_type:
            stmt = stmt.where(Agency.agency_type == query.agency_type)
        if query.department:
            stmt = stmt.where(Agency.department == query.department)
        if query.search_term:
            search = f"%{query.search_term}%"
            stmt = stmt.where(
                or_(Agency.name.ilike(search), Agency.code.ilike(search))
            )
        count_stmt = select(func.count()).select_from(stmt.subquery())
        total = (await self._session.execute(count_stmt)).scalar_one()
        offset = (query.page - 1) * query.page_size
        stmt = stmt.offset(offset).limit(query.page_size)
        result = await self._session.execute(stmt)
        items = result.scalars().all()
        total_pages = (total + query.page_size - 1) // query.page_size
        return PaginatedAgencyResponse(
            items=[self._to_response(a) for a in items],
            total=total,
            page=query.page,
            page_size=query.page_size,
            total_pages=total_pages,
        )

    @staticmethod
    def _to_response(agency: Agency) -> AgencyResponse:
        return AgencyResponse(
            id=agency.id,
            name=agency.name,
            code=agency.code,
            agency_type=agency.agency_type.value if agency.agency_type else "",
            status=agency.status.value if agency.status else "",
            address=agency.address,
            city=agency.city,
            department=agency.department,
            phone=agency.phone,
            email=agency.email,
            is_headquarters=agency.is_headquarters,
            parent_agency_id=agency.parent_agency_id,
            max_daily_enrollments=agency.max_daily_enrollments,
            opened_at=agency.opened_at.isoformat() if agency.opened_at else None,
            closed_at=agency.closed_at.isoformat() if agency.closed_at else None,
        )


class AgencyNotFoundError(Exception):
    pass
