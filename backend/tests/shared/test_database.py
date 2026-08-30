from __future__ import annotations

import uuid
from datetime import datetime, timezone
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
import pytest_asyncio
from sqlalchemy.ext.asyncio import AsyncSession

from shared.database import Base, check_database_health, close_database, get_session
from shared.database.event_store import (
    AggregateSnapshot,
    ConcurrencyError,
    EventStore,
    StoredEvent,
)


class TestBaseModel:
    """Test Base model columns and soft delete."""

    def test_base_columns_exist(self):
        assert hasattr(Base, "id")
        assert hasattr(Base, "created_at")
        assert hasattr(Base, "updated_at")
        assert hasattr(Base, "is_deleted")
        assert hasattr(Base, "deleted_at")

    def test_base_defaults(self):
        assert hasattr(Base, "is_deleted")
        assert hasattr(Base, "deleted_at")

    def test_soft_delete(self):
        instance = Base()
        instance.soft_delete()
        assert instance.is_deleted is True
        assert instance.deleted_at is not None


class TestSessionManagement:
    """Test database session management."""

    @pytest.mark.asyncio
    async def test_get_session_raises_if_not_initialized(self):
        with pytest.raises(RuntimeError, match="Database not initialized"):
            async with get_session():
                pass

    @pytest.mark.asyncio
    async def test_check_health_false_when_not_initialized(self):
        result = await check_database_health()
        assert result is False

    @pytest.mark.asyncio
    async def test_close_database_when_not_initialized(self):
        await close_database()
        assert True


class TestEventStore:
    """Test event store operations."""

    @pytest.fixture
    def session(self):
        return AsyncMock(spec=AsyncSession)

    @pytest.fixture
    def store(self, session):
        return EventStore(session)

    @pytest.mark.asyncio
    async def test_append_events(self, store, session):
        session.execute = AsyncMock(return_value=MagicMock(scalar_one_or_none=MagicMock(return_value=0)))

        events = [
            {"event_type": "identity.created", "data": {"name": "John"}},
        ]
        result = await store.append_events(
            aggregate_id="agg-1",
            aggregate_type="Identity",
            events=events,
            expected_version=0,
            actor_id="user-1",
            correlation_id="corr-1",
        )
        assert len(result) == 1
        session.add.assert_called_once()
        session.flush.assert_awaited_once()

    @pytest.mark.asyncio
    async def test_concurrency_conflict(self, store, session):
        session.execute = AsyncMock(return_value=MagicMock(scalar_one_or_none=MagicMock(return_value=5)))

        events = [
            {"event_type": "identity.updated", "data": {"field": "name"}},
        ]
        with pytest.raises(ConcurrencyError, match="expected version 3, found 5"):
            await store.append_events(
                aggregate_id="agg-1",
                aggregate_type="Identity",
                events=events,
                expected_version=3,
            )

    @pytest.mark.asyncio
    async def test_get_events(self, store, session):
        mock_event = MagicMock(spec=StoredEvent)
        mock_event.event_id = str(uuid.uuid4())
        mock_event.event_type = "identity.created"
        mock_event.version = 1
        mock_event.event_data = {"name": "John"}
        mock_event.actor_id = "user-1"
        mock_event.timestamp = datetime.now(timezone.utc)
        mock_event.aggregate_id = "agg-1"
        mock_event.correlation_id = None
        mock_event.event_metadata = None

        mock_scalars = MagicMock()
        mock_scalars.all = MagicMock(return_value=[mock_event])
        session.execute = AsyncMock(return_value=MagicMock(scalars=MagicMock(return_value=mock_scalars)))

        events = await store.get_events(aggregate_id="agg-1")
        assert len(events) == 1
        assert events[0].event_type == "identity.created"

    @pytest.mark.asyncio
    async def test_get_events_after_version(self, store, session):
        mock_scalars = MagicMock()
        mock_scalars.all = MagicMock(return_value=[])
        session.execute = AsyncMock(return_value=MagicMock(scalars=MagicMock(return_value=mock_scalars)))

        events = await store.get_events(aggregate_id="agg-1", after_version=5)
        assert len(events) == 0

    @pytest.mark.asyncio
    async def test_get_events_with_limit(self, store, session):
        mock_scalars = MagicMock()
        mock_scalars.all = MagicMock(return_value=[])
        session.execute = AsyncMock(return_value=MagicMock(scalars=MagicMock(return_value=mock_scalars)))

        events = await store.get_events(aggregate_id="agg-1", limit=10)
        assert len(events) == 0

    @pytest.mark.asyncio
    async def test_get_events_empty(self, store, session):
        mock_scalars = MagicMock()
        mock_scalars.all = MagicMock(return_value=[])
        session.execute = AsyncMock(return_value=MagicMock(scalars=MagicMock(return_value=mock_scalars)))

        events = await store.get_events(aggregate_id="nonexistent")
        assert len(events) == 0


class TestEventStoreSnapshots:
    """Test snapshot operations."""

    @pytest.fixture
    def session(self):
        return AsyncMock(spec=AsyncSession)

    @pytest.fixture
    def store(self, session):
        return EventStore(session)

    @pytest.mark.asyncio
    async def test_save_new_snapshot(self, store, session):
        session.execute = AsyncMock(return_value=MagicMock(scalar_one_or_none=MagicMock(return_value=None)))

        snapshot = await store.save_snapshot(
            aggregate_id="agg-1",
            aggregate_type="Identity",
            version=50,
            state={"status": "active"},
        )
        assert snapshot.aggregate_id == "agg-1"
        assert snapshot.version == 50
        session.add.assert_called_once()
        session.flush.assert_awaited_once()

    @pytest.mark.asyncio
    async def test_get_snapshot(self, store, session):
        mock_snapshot = MagicMock(spec=AggregateSnapshot)
        mock_snapshot.aggregate_id = "agg-1"
        mock_snapshot.aggregate_type = "Identity"
        mock_snapshot.version = 50
        mock_snapshot.state = {"status": "active"}

        session.execute = AsyncMock(
            return_value=MagicMock(scalar_one_or_none=MagicMock(return_value=mock_snapshot))
        )

        snapshot = await store.get_snapshot("agg-1", "Identity")
        assert snapshot is not None
        assert snapshot.version == 50
        assert snapshot.state["status"] == "active"

    @pytest.mark.asyncio
    async def test_get_snapshot_nonexistent(self, store, session):
        session.execute = AsyncMock(
            return_value=MagicMock(scalar_one_or_none=MagicMock(return_value=None))
        )

        snapshot = await store.get_snapshot("nonexistent", "Identity")
        assert snapshot is None

    @pytest.mark.asyncio
    async def test_save_snapshot_updates_existing(self, store, session):
        existing = MagicMock(spec=AggregateSnapshot)
        existing.version = 40
        existing.state = {"status": "pending"}

        session.execute = AsyncMock(
            return_value=MagicMock(scalar_one_or_none=MagicMock(return_value=existing))
        )

        snapshot = await store.save_snapshot(
            aggregate_id="agg-1",
            aggregate_type="Identity",
            version=60,
            state={"status": "active"},
        )
        assert snapshot is existing
        assert existing.version == 60
        assert existing.state == {"status": "active"}


class TestEventStoreByType:
    """Test retrieving events by type."""

    @pytest.fixture
    def session(self):
        return AsyncMock(spec=AsyncSession)

    @pytest.fixture
    def store(self, session):
        return EventStore(session)

    @pytest.mark.asyncio
    async def test_get_events_by_type(self, store, session):
        mock_scalars = MagicMock()
        mock_scalars.all = MagicMock(return_value=[])
        session.execute = AsyncMock(return_value=MagicMock(scalars=MagicMock(return_value=mock_scalars)))

        events = await store.get_events_by_type("Identity")
        assert len(events) == 0

    @pytest.mark.asyncio
    async def test_get_events_by_correlation(self, store, session):
        mock_scalars = MagicMock()
        mock_scalars.all = MagicMock(return_value=[])
        session.execute = AsyncMock(return_value=MagicMock(scalars=MagicMock(return_value=mock_scalars)))

        events = await store.get_events_by_correlation("corr-1")
        assert len(events) == 0

    @pytest.mark.asyncio
    async def test_get_aggregate_ids(self, store, session):
        mock_scalars = MagicMock()
        mock_scalars.all = MagicMock(return_value=["agg-1", "agg-2"])
        session.execute = AsyncMock(return_value=MagicMock(scalars=MagicMock(return_value=mock_scalars)))

        ids = await store.get_aggregate_ids("Identity")
        assert len(ids) == 2
