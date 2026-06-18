from __future__ import annotations

from datetime import date
from unittest.mock import AsyncMock, MagicMock

import pytest

from services.identity.projections.citizen_projector import CitizenProjector


@pytest.fixture
def mock_session():
    session = AsyncMock()
    session.execute = AsyncMock()
    session.flush = AsyncMock()
    session.add = MagicMock()
    return session


class TestCitizenProjector:
    async def test_project_created(self, mock_session):
        projector = CitizenProjector(mock_session)
        await projector.project("IdentityCreated", "agg-1", {
            "national_id": "SN123456789",
            "first_name": "John",
            "last_name": "Doe",
            "date_of_birth": "1990-01-15",
            "gender": "male",
            "nationality": "SEN",
            "agency_id": "agency-1",
        })
        mock_session.add.assert_called_once()
        mock_session.flush.assert_called_once()

    async def test_project_verified(self, mock_session):
        read_model = MagicMock()
        read_model.status = "pending"
        mock_session.execute.return_value.scalar_one_or_none = MagicMock(return_value=read_model)
        projector = CitizenProjector(mock_session)
        await projector.project("IdentityVerified", "agg-1", {
            "verification_method": "biometric",
        })
        assert read_model.status == "active"
        assert read_model.verified is True

    async def test_project_suspended(self, mock_session):
        read_model = MagicMock()
        mock_session.execute.return_value.scalar_one_or_none = MagicMock(return_value=read_model)
        projector = CitizenProjector(mock_session)
        await projector.project("IdentitySuspended", "agg-1", {"reason": "Fraud investigation"})
        assert read_model.status == "suspended"

    async def test_project_revoked(self, mock_session):
        read_model = MagicMock()
        mock_session.execute.return_value.scalar_one_or_none = MagicMock(return_value=read_model)
        projector = CitizenProjector(mock_session)
        await projector.project("IdentityRevoked", "agg-1", {"reason": "Permanent"})
        assert read_model.status == "revoked"
        mock_session.delete.assert_called_once_with(read_model)

    async def test_unknown_event_ignored(self, mock_session):
        projector = CitizenProjector(mock_session)
        await projector.project("UnknownEvent", "agg-1", {})
        mock_session.add.assert_not_called()
        mock_session.flush.assert_not_called()
