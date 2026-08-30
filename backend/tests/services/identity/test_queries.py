from __future__ import annotations

import uuid
from datetime import date, datetime, timezone
from unittest.mock import AsyncMock, MagicMock

import pytest
import pytest_asyncio

from services.identity.queries import (
    IdentityQueryHandler,
    IdentityQueryNotFoundError,
    GetIdentityByIdQuery,
    GetIdentityByNationalIdQuery,
    SearchIdentitiesQuery,
    GetIdentityHistoryQuery,
    GetIdentityStatsQuery,
    IdentityResponse,
    IdentityDetailResponse,
    PaginatedResponse,
    IdentityHistoryEntry,
    IdentityStatsResponse,
)


class TestGetIdentityQuery:
    """Test get identity by ID query."""

    @pytest.mark.asyncio
    async def test_get_by_id_found(self, mock_db_session):
        mock_citizen = MagicMock()
        mock_citizen.id = str(uuid.uuid4())
        mock_citizen.national_id = "SN123456789"
        mock_citizen.first_name = "John"
        mock_citizen.last_name = "Doe"
        mock_citizen.date_of_birth = date(1990, 1, 15)
        mock_citizen.place_of_birth = "Dakar"
        mock_citizen.gender = MagicMock(value="male")
        mock_citizen.nationality = "SEN"
        mock_citizen.status = MagicMock(value="active")
        mock_citizen.agency_id = str(uuid.uuid4())
        mock_citizen.verified_at = datetime.now(timezone.utc)
        mock_citizen.verified_by = "verifier-1"
        mock_citizen.photo_url = None
        mock_citizen.email = None
        mock_citizen.phone = None
        mock_citizen.marital_status = None
        mock_citizen.address = None
        mock_citizen.middle_name = None
        mock_citizen.version = 2
        mock_citizen.created_at = datetime.now(timezone.utc)
        mock_citizen.updated_at = datetime.now(timezone.utc)
        mock_citizen.documents = []
        mock_citizen.biometrics = []

        mock_result = MagicMock()
        mock_result.scalar_one_or_none = MagicMock(return_value=mock_citizen)
        mock_db_session.execute = AsyncMock(return_value=mock_result)

        handler = IdentityQueryHandler(mock_db_session)
        query = GetIdentityByIdQuery(identity_id=mock_citizen.id)
        result = await handler.get_by_id(query)

        assert isinstance(result, IdentityDetailResponse)
        assert result.id == mock_citizen.id
        assert result.national_id == "SN123456789"
        assert result.full_name == "John Doe"
        assert result.verified is True

    @pytest.mark.asyncio
    async def test_get_by_id_not_found(self, mock_db_session):
        mock_result = MagicMock()
        mock_result.scalar_one_or_none = MagicMock(return_value=None)
        mock_db_session.execute = AsyncMock(return_value=mock_result)

        handler = IdentityQueryHandler(mock_db_session)
        query = GetIdentityByIdQuery(identity_id="nonexistent")
        with pytest.raises(IdentityQueryNotFoundError, match="nonexistent"):
            await handler.get_by_id(query)

    @pytest.mark.asyncio
    async def test_get_by_id_with_documents(self, mock_db_session):
        mock_doc = MagicMock()
        mock_doc.id = str(uuid.uuid4())
        mock_doc.document_type = MagicMock(value="national_id")
        mock_doc.document_number = "NID123"
        mock_doc.issue_date = date(2024, 1, 1)
        mock_doc.expiry_date = date(2034, 1, 1)
        mock_doc.status = MagicMock(value="active")
        mock_doc.issuing_agency = "AN-AGENCY"

        mock_bio = MagicMock()
        mock_bio.id = str(uuid.uuid4())
        mock_bio.biometric_type = MagicMock(value="fingerprint")
        mock_bio.quality_score = 0.95
        mock_bio.captured_at = datetime.now(timezone.utc)
        mock_bio.is_primary = True

        mock_citizen = MagicMock()
        mock_citizen.id = str(uuid.uuid4())
        mock_citizen.national_id = "SN999"
        mock_citizen.first_name = "Jane"
        mock_citizen.last_name = "Smith"
        mock_citizen.date_of_birth = date(1995, 5, 20)
        mock_citizen.place_of_birth = "Thies"
        mock_citizen.gender = MagicMock(value="female")
        mock_citizen.nationality = "SEN"
        mock_citizen.status = MagicMock(value="active")
        mock_citizen.agency_id = "agency-1"
        mock_citizen.verified_at = None
        mock_citizen.verified_by = None
        mock_citizen.photo_url = "https://example.com/photo.jpg"
        mock_citizen.email = "jane@example.com"
        mock_citizen.phone = "+221123456"
        mock_citizen.marital_status = "single"
        mock_citizen.address = {"city": "Dakar"}
        mock_citizen.middle_name = "Marie"
        mock_citizen.version = 1
        mock_citizen.created_at = datetime.now(timezone.utc)
        mock_citizen.updated_at = datetime.now(timezone.utc)
        mock_citizen.documents = [mock_doc]
        mock_citizen.biometrics = [mock_bio]

        mock_result = MagicMock()
        mock_result.scalar_one_or_none = MagicMock(return_value=mock_citizen)
        mock_db_session.execute = AsyncMock(return_value=mock_result)

        handler = IdentityQueryHandler(mock_db_session)
        query = GetIdentityByIdQuery(identity_id=mock_citizen.id)
        result = await handler.get_by_id(query)

        assert result.document_count == 1
        assert result.has_biometrics is True
        assert result.photo_url == "https://example.com/photo.jpg"
        assert result.email == "jane@example.com"
        assert len(result.documents) == 1
        assert len(result.biometrics) == 1


class TestGetIdentityByNationalIdQuery:
    """Test get identity by national ID query."""

    @pytest.mark.asyncio
    async def test_get_by_national_id_found(self, mock_db_session):
        mock_citizen = MagicMock()
        mock_citizen.id = str(uuid.uuid4())
        mock_citizen.national_id = "SN123456789"
        mock_citizen.first_name = "John"
        mock_citizen.last_name = "Doe"
        mock_citizen.date_of_birth = date(1990, 1, 15)
        mock_citizen.place_of_birth = "Dakar"
        mock_citizen.gender = MagicMock(value="male")
        mock_citizen.nationality = "SEN"
        mock_citizen.status = MagicMock(value="active")
        mock_citizen.agency_id = str(uuid.uuid4())
        mock_citizen.verified_at = None
        mock_citizen.verified_by = None
        mock_citizen.photo_url = None
        mock_citizen.email = None
        mock_citizen.phone = None
        mock_citizen.marital_status = None
        mock_citizen.address = None
        mock_citizen.middle_name = None
        mock_citizen.version = 1
        mock_citizen.created_at = datetime.now(timezone.utc)
        mock_citizen.updated_at = datetime.now(timezone.utc)
        mock_citizen.documents = []
        mock_citizen.biometrics = []

        mock_result = MagicMock()
        mock_result.scalar_one_or_none = MagicMock(return_value=mock_citizen)
        mock_db_session.execute = AsyncMock(return_value=mock_result)

        handler = IdentityQueryHandler(mock_db_session)
        query = GetIdentityByNationalIdQuery(national_id="SN123456789")
        result = await handler.get_by_national_id(query)
        assert result.national_id == "SN123456789"

    @pytest.mark.asyncio
    async def test_get_by_national_id_not_found(self, mock_db_session):
        mock_result = MagicMock()
        mock_result.scalar_one_or_none = MagicMock(return_value=None)
        mock_db_session.execute = AsyncMock(return_value=mock_result)

        handler = IdentityQueryHandler(mock_db_session)
        query = GetIdentityByNationalIdQuery(national_id="NONEXISTENT")
        with pytest.raises(IdentityQueryNotFoundError):
            await handler.get_by_national_id(query)


class TestSearchIdentitiesQuery:
    """Test search identities query."""

    @pytest.mark.asyncio
    async def test_search_without_filters(self, mock_db_session):
        call_count = 0

        async def execute_side_effect(*args, **kwargs):
            nonlocal call_count
            call_count += 1
            if call_count == 1:
                result = MagicMock()
                result.scalar_one = MagicMock(return_value=1)
                return result
            result = MagicMock()
            mock_scalars = MagicMock()
            mock_citizen = MagicMock(
                id="id-1", national_id="SN123456789",
                full_name="John Doe", first_name="John", last_name="Doe",
                middle_name=None,
                date_of_birth=date(1990, 1, 15), place_of_birth="Dakar",
                gender="male", nationality="SEN",
                status="active", agency_id="agency-1",
                verified_at=None, verified_by=None,
                photo_url=None, email=None, phone=None,
                marital_status=None, address=None,
                documents=[], biometrics=[],
                version=1, created_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
                updated_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
                is_deleted=False,
            )
            mock_scalars.all = MagicMock(return_value=[mock_citizen])
            result.scalars = MagicMock(return_value=mock_scalars)
            return result

        mock_db_session.execute = execute_side_effect

        handler = IdentityQueryHandler(mock_db_session)
        query = SearchIdentitiesQuery()
        result = await handler.search(query)

        assert isinstance(result, PaginatedResponse)
        assert result.total == 1
        assert len(result.items) == 1
        assert result.items[0].full_name == "John Doe"

    @pytest.mark.asyncio
    async def test_search_with_status_filter(self, mock_db_session):
        async def execute_side_effect(*args, **kwargs):
            return MagicMock(scalar_one=MagicMock(return_value=0))

        mock_db_session.execute = execute_side_effect

        handler = IdentityQueryHandler(mock_db_session)
        query = SearchIdentitiesQuery(status="active")
        result = await handler.search(query)
        assert result.total == 0

    @pytest.mark.asyncio
    async def test_search_empty_results(self, mock_db_session):
        mock_scalars = MagicMock()
        mock_scalars.all = MagicMock(return_value=[])
        mock_execute = AsyncMock(return_value=MagicMock(
            scalars=MagicMock(return_value=mock_scalars)
        ))

        async def execute_side_effect(*args, **kwargs):
            return MagicMock(scalar_one=MagicMock(return_value=0))

        mock_db_session.execute = execute_side_effect

        handler = IdentityQueryHandler(mock_db_session)
        query = SearchIdentitiesQuery(search_term="nonexistent")
        result = await handler.search(query)
        assert result.total == 0
        assert len(result.items) == 0

    @pytest.mark.asyncio
    async def test_search_pagination(self, mock_db_session):
        mock_scalars = MagicMock()
        mock_scalars.all = MagicMock(return_value=[])
        mock_execute = AsyncMock(return_value=MagicMock(
            scalars=MagicMock(return_value=mock_scalars)
        ))

        async def execute_side_effect(*args, **kwargs):
            return MagicMock(scalar_one=MagicMock(return_value=25))

        mock_db_session.execute = execute_side_effect

        handler = IdentityQueryHandler(mock_db_session)
        query = SearchIdentitiesQuery(page=2, page_size=10)
        result = await handler.search(query)
        assert result.page == 2
        assert result.page_size == 10
        assert result.total_pages == 3

    @pytest.mark.asyncio
    async def test_search_with_all_filters(self, mock_db_session):
        async def execute_side_effect(*args, **kwargs):
            return MagicMock(scalar_one=MagicMock(return_value=0))

        mock_db_session.execute = execute_side_effect

        handler = IdentityQueryHandler(mock_db_session)
        query = SearchIdentitiesQuery(
            search_term="John",
            status="active",
            agency_id="agency-1",
            nationality="SEN",
            sort_by="last_name",
            sort_order="asc",
        )
        result = await handler.search(query)
        assert result.total == 0


class TestGetIdentityHistoryQuery:
    """Test get identity history query."""

    @pytest.mark.asyncio
    async def test_get_history(self, mock_db_session, mock_event_store):
        mock_event = MagicMock()
        mock_event.event_id = "evt-1"
        mock_event.event_type = "identity.created"
        mock_event.version = 1
        mock_event.timestamp = datetime.now(timezone.utc)
        mock_event.actor_id = "user-1"
        mock_event.event_data = {"national_id": "SN123"}

        mock_event_store.get_events = AsyncMock(return_value=[mock_event])

        handler = IdentityQueryHandler(mock_db_session)
        handler._event_store = mock_event_store
        query = GetIdentityHistoryQuery(identity_id="agg-1")
        result = await handler.get_history(query)

        assert len(result) == 1
        assert isinstance(result[0], IdentityHistoryEntry)
        assert result[0].event_type == "identity.created"
        assert result[0].version == 1

    @pytest.mark.asyncio
    async def test_get_history_empty(self, mock_db_session, mock_event_store):
        mock_event_store.get_events = AsyncMock(return_value=[])

        handler = IdentityQueryHandler(mock_db_session)
        handler._event_store = mock_event_store
        query = GetIdentityHistoryQuery(identity_id="nonexistent")
        result = await handler.get_history(query)
        assert len(result) == 0

    @pytest.mark.asyncio
    async def test_get_history_after_version(self, mock_db_session, mock_event_store):
        mock_event_store.get_events = AsyncMock(return_value=[])

        handler = IdentityQueryHandler(mock_db_session)
        handler._event_store = mock_event_store
        query = GetIdentityHistoryQuery(identity_id="agg-1", after_version=5)
        result = await handler.get_history(query)
        assert len(result) == 0


class TestGetIdentityStatsQuery:
    """Test get identity statistics query."""

    @pytest.mark.asyncio
    async def test_get_stats(self, mock_db_session):
        async def execute_side_effect(stmt):
            compiled = str(stmt.compile(compile_kwargs={"literal_binds": True}))
            if "count" in compiled.lower():
                if "group_by" in compiled.lower() or "status" in compiled:
                    return MagicMock(all=MagicMock(return_value=[("active", 5), ("pending", 3)]))
                return MagicMock(scalar_one=MagicMock(return_value=8))
            return MagicMock(scalar_one=MagicMock(return_value=0))

        mock_db_session.execute = execute_side_effect

        handler = IdentityQueryHandler(mock_db_session)
        query = GetIdentityStatsQuery(agency_id=None)
        result = await handler.get_stats(query)

        assert isinstance(result, IdentityStatsResponse)
        assert result.total >= 0

    @pytest.mark.asyncio
    async def test_get_stats_with_agency(self, mock_db_session):
        async def execute_side_effect(stmt):
            return MagicMock(scalar_one=MagicMock(return_value=0))

        mock_db_session.execute = execute_side_effect

        handler = IdentityQueryHandler(mock_db_session)
        query = GetIdentityStatsQuery(agency_id="agency-1")
        result = await handler.get_stats(query)
        assert result.total == 0


class TestQueryPaginationLinks:
    """Test HATEOAS pagination links."""

    def test_build_pagination_links_first_page(self):
        query = SearchIdentitiesQuery(page=1, page_size=20)
        links = IdentityQueryHandler._build_pagination_links(query, total_pages=5)
        assert "self" in links
        assert "first" in links
        assert "last" in links
        assert "next" in links
        assert "prev" not in links

    def test_build_pagination_links_middle_page(self):
        query = SearchIdentitiesQuery(page=3, page_size=20)
        links = IdentityQueryHandler._build_pagination_links(query, total_pages=5)
        assert "prev" in links
        assert "next" in links

    def test_build_pagination_links_last_page(self):
        query = SearchIdentitiesQuery(page=5, page_size=20)
        links = IdentityQueryHandler._build_pagination_links(query, total_pages=5)
        assert "prev" in links
        assert "next" not in links

    def test_build_pagination_links_single_page(self):
        query = SearchIdentitiesQuery(page=1, page_size=20)
        links = IdentityQueryHandler._build_pagination_links(query, total_pages=1)
        assert "next" not in links
        assert "prev" not in links


class TestIdentityLinks:
    """Test HATEOAS identity links."""

    def test_build_identity_links(self):
        identity_id = str(uuid.uuid4())
        links = IdentityQueryHandler._build_identity_links(identity_id)
        assert "self" in links
        assert "update" in links
        assert "verify" in links
        assert "suspend" in links
        assert "revoke" in links
        assert "history" in links
        assert "documents" in links
        assert "biometrics" in links
        assert links["self"]["href"] == f"/v1/identities/{identity_id}"
