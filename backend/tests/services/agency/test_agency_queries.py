from __future__ import annotations

from unittest.mock import AsyncMock, MagicMock

import pytest

from services.agency.models import Agency, AgencyStatus, AgencyType
from services.agency.queries import (
    AgencyNotFoundError,
    AgencyQueryHandler,
    GetAgencyByIdQuery,
    ListAgenciesQuery,
)


@pytest.fixture
def mock_session():
    session = AsyncMock()
    return session


class TestAgencyQueryHandler:
    async def _make_mock_agency(self):
        agency = MagicMock(spec=Agency)
        agency.id = "agency-1"
        agency.name = "Test Agency"
        agency.code = "TEST01"
        agency.agency_type = AgencyType.LOCAL
        agency.status = AgencyStatus.ACTIVE
        agency.is_deleted = False
        agency.address = None
        agency.city = "Port-au-Prince"
        agency.department = "Ouest"
        agency.phone = None
        agency.email = None
        agency.is_headquarters = False
        agency.parent_agency_id = None
        agency.max_daily_enrollments = 500
        agency.opened_at = None
        agency.closed_at = None
        return agency

    async def test_get_by_id_found(self, mock_session):
        agency = await self._make_mock_agency()
        mock_session.execute.return_value.scalar_one_or_none = MagicMock(return_value=agency)
        handler = AgencyQueryHandler(mock_session)
        result = await handler.get_by_id(GetAgencyByIdQuery(agency_id="agency-1"))
        assert result.id == "agency-1"
        assert result.name == "Test Agency"
        assert result.code == "TEST01"

    async def test_get_by_id_not_found(self, mock_session):
        mock_session.execute.return_value.scalar_one_or_none = MagicMock(return_value=None)
        handler = AgencyQueryHandler(mock_session)
        with pytest.raises(AgencyNotFoundError):
            await handler.get_by_id(GetAgencyByIdQuery(agency_id="nonexistent"))

    async def test_list_all(self, mock_session):
        agency = await self._make_mock_agency()
        count_result = MagicMock()
        count_result.scalar_one = MagicMock(return_value=1)
        data_result = MagicMock()
        data_result.scalars.return_value.all = MagicMock(return_value=[agency])
        mock_session.execute = AsyncMock(side_effect=[count_result, data_result])
        handler = AgencyQueryHandler(mock_session)
        result = await handler.list_all(ListAgenciesQuery())
        assert len(result.items) == 1
