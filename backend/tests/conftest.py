from __future__ import annotations

import uuid
from unittest.mock import AsyncMock, MagicMock

import pytest
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.ext.compiler import compiles

from shared.cqrs import Command, CommandBus, DomainEvent, Query, QueryBus
from services.identity.aggregate import IdentityAggregate


# ── JSONB Compatibility for SQLite ───────────────────────────────────────


@compiles(JSONB, "sqlite")
def compile_jsonb_sqlite(type_, compiler, **kw):
    return "JSON"


# ── CQRS Test Helpers ────────────────────────────────────────────────────


class SampleCommand(Command):
    value: str


class SampleCommandHandler:
    async def handle(self, command: SampleCommand) -> dict:
        return {"handled": command.value}


class AnotherCommand(Command):
    number: int


class AnotherCommandHandler:
    async def handle(self, command: AnotherCommand) -> dict:
        return {"number": command.number}


class SampleQuery(Query):
    query_id: str = "default"


class SampleQueryHandler:
    async def handle(self, query: SampleQuery) -> dict:
        return {"result": "sample", "query_id": query.query_id}


class AnotherQuery(Query):
    term: str = ""


class AnotherQueryHandler:
    async def handle(self, query: AnotherQuery) -> dict:
        return {"term": query.term}


# ── CQRS Fixtures ────────────────────────────────────────────────────────


@pytest.fixture
def command_bus():
    from unittest.mock import AsyncMock, patch
    from contextlib import asynccontextmanager

    mock_session = AsyncMock()

    @asynccontextmanager
    async def fake_get_session():
        yield mock_session

    with patch("shared.database.get_session", fake_get_session):
        yield CommandBus()


@pytest.fixture
def query_bus():
    return QueryBus()


@pytest.fixture
def query_bus_with_cache():
    from fakeredis.aioredis import FakeRedis
    return QueryBus(redis_client=FakeRedis())


# ── Event Bus Fixtures ────────────────────────────────────────────────────


@pytest.fixture
def sample_event():
    return DomainEvent(
        event_type="test.event",
        aggregate_id=str(uuid.uuid4()),
        aggregate_type="Test",
        data={"key": "value"},
    )


@pytest.fixture
def sample_events():
    return [
        DomainEvent(
            event_type="test.event",
            aggregate_id=str(uuid.uuid4()),
            aggregate_type="Test",
            data={"key": "value"},
        ),
        DomainEvent(
            event_type="test.event2",
            aggregate_id=str(uuid.uuid4()),
            aggregate_type="Test",
            data={"key2": "value2"},
        ),
    ]


# ── Redis Fixtures ────────────────────────────────────────────────────────


@pytest.fixture
def redis_client():
    from fakeredis.aioredis import FakeRedis
    return FakeRedis()


# ── Identity Test Fixtures ───────────────────────────────────────────────


@pytest.fixture
def identity_create_data():
    return {
        "national_id": "SN123456789",
        "first_name": "John",
        "last_name": "Doe",
        "date_of_birth": "1990-01-15",
        "place_of_birth": "Dakar",
        "gender": "male",
        "nationality": "SEN",
        "agency_id": str(uuid.uuid4()),
        "actor_id": "user-1",
        "correlation_id": str(uuid.uuid4()),
    }


@pytest.fixture
def identity_aggregate(identity_create_data):
    return IdentityAggregate.create(**identity_create_data)


@pytest.fixture
def sample_identity_id():
    return str(uuid.uuid4())


@pytest.fixture
def mock_db_session():
    session = AsyncMock(spec=AsyncSession)
    result = AsyncMock()
    result.scalar_one_or_none = MagicMock(return_value=None)
    result.scalar = MagicMock(return_value=0)
    result.scalars = MagicMock()
    result.scalars.return_value.all = MagicMock(return_value=[])
    session.execute.return_value = result
    return session


@pytest.fixture
def mock_event_store():
    store = AsyncMock()
    store.get_events = AsyncMock(return_value=[])
    store.get_snapshot = AsyncMock(return_value=None)
    store.append_events = AsyncMock()
    store.save_snapshot = AsyncMock()
    return store


@pytest.fixture
def mock_kafka():
    kafka = AsyncMock()
    kafka.publish = AsyncMock()
    return kafka
