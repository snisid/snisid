from __future__ import annotations

from unittest.mock import AsyncMock, MagicMock

import pytest

from services.agency.commands import (
    AgencyCommandHandler,
    AgencyNotFoundError,
    CreateAgencyCommand,
    DeactivateAgencyCommand,
    UpdateAgencyCommand,
)
from services.agency.models import Agency, AgencyType


@pytest.fixture
def mock_session():
    session = AsyncMock()
    session.execute = AsyncMock()
    session.flush = AsyncMock()
    session.add = MagicMock()
    return session


class TestAgencyCommandHandler:
    async def test_create_agency(self, mock_session):
        handler = AgencyCommandHandler(mock_session)
        cmd = CreateAgencyCommand(
            name="Test Agency",
            code="TEST01",
            agency_type=AgencyType.LOCAL,
            city="Port-au-Prince",
            department="Ouest",
        )
        result = await handler.handle_create(cmd)
        assert result["status"] == "active"
        assert result["code"] == "TEST01"
        mock_session.add.assert_called_once()

    async def test_update_agency(self, mock_session):
        agency = MagicMock(spec=Agency)
        agency.id = "agency-1"
        agency.is_deleted = False
        mock_session.execute.return_value.scalar_one_or_none = MagicMock(return_value=agency)
        handler = AgencyCommandHandler(mock_session)
        cmd = UpdateAgencyCommand(agency_id="agency-1", changes={"name": "Updated Name"})
        result = await handler.handle_update(cmd)
        assert result["agency_id"] == "agency-1"
        assert agency.name == "Updated Name"

    async def test_update_agency_not_found(self, mock_session):
        mock_session.execute.return_value.scalar_one_or_none = MagicMock(return_value=None)
        handler = AgencyCommandHandler(mock_session)
        cmd = UpdateAgencyCommand(agency_id="nonexistent", changes={"name": "X"})
        with pytest.raises(AgencyNotFoundError):
            await handler.handle_update(cmd)

    async def test_deactivate_agency(self, mock_session):
        agency = MagicMock(spec=Agency)
        agency.id = "agency-1"
        agency.is_deleted = False
        mock_session.execute.return_value.scalar_one_or_none = MagicMock(return_value=agency)
        handler = AgencyCommandHandler(mock_session)
        cmd = DeactivateAgencyCommand(agency_id="agency-1", reason="Permanent closure")
        result = await handler.handle_deactivate(cmd)
        assert result["status"] == "inactive"
