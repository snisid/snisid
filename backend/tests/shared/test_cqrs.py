from __future__ import annotations

import pytest
import pytest_asyncio
from unittest.mock import AsyncMock, MagicMock

from shared.cqrs import (
    Command,
    CommandBus,
    CommandHandlerNotFoundError,
    CommandMiddleware,
    DomainEvent,
    Query,
    QueryBus,
    QueryHandlerNotFoundError,
)

from tests.conftest import (
    SampleCommand,
    SampleCommandHandler,
    SampleQuery,
    SampleQueryHandler,
    AnotherCommand,
    AnotherCommandHandler,
    AnotherQuery,
    AnotherQueryHandler,
)


class TestDomainEvent:
    """Test DomainEvent creation and serialization."""

    def test_create_domain_event(self):
        event = DomainEvent(
            event_type="citizen.registered",
            aggregate_id="agg-1",
            aggregate_type="Citizen",
            data={"name": "John"},
        )
        assert event.event_type == "citizen.registered"
        assert event.aggregate_id == "agg-1"
        assert event.aggregate_type == "Citizen"
        assert event.data == {"name": "John"}
        assert event.event_id is not None
        assert event.timestamp is not None
        assert event.version == 1

    def test_domain_event_defaults(self):
        event = DomainEvent(
            event_type="test.event",
            aggregate_id="agg-1",
            aggregate_type="Test",
        )
        assert event.actor_id is None
        assert event.correlation_id is None
        assert event.data == {}
        assert event.version == 1

    def test_domain_event_frozen(self):
        event = DomainEvent(
            event_type="test.event",
            aggregate_id="agg-1",
            aggregate_type="Test",
        )
        with pytest.raises((TypeError, ValueError)):
            event.event_type = "changed"

    def test_domain_event_serialization(self):
        event = DomainEvent(
            event_type="test.event",
            aggregate_id="agg-1",
            aggregate_type="Test",
            data={"key": "value"},
        )
        data = event.model_dump(mode="json")
        assert data["event_type"] == "test.event"
        assert data["data"]["key"] == "value"
        assert "event_id" in data

        restored = DomainEvent.model_validate(data)
        assert restored.event_type == event.event_type
        assert restored.data == event.data
        assert restored.event_id == event.event_id

    def test_domain_event_with_actor_and_correlation(self):
        event = DomainEvent(
            event_type="test.event",
            aggregate_id="agg-1",
            aggregate_type="Test",
            actor_id="user-1",
            correlation_id="corr-1",
        )
        assert event.actor_id == "user-1"
        assert event.correlation_id == "corr-1"

    def test_version_increments_independently(self):
        event1 = DomainEvent(event_type="e1", aggregate_id="a", aggregate_type="T")
        event2 = DomainEvent(event_type="e2", aggregate_id="a", aggregate_type="T")
        assert event1.version == 1
        assert event2.version == 1


class TestCommandBus:
    """Test command registration and dispatch."""

    @pytest.mark.asyncio
    async def test_register_and_dispatch(self, command_bus):
        handler = SampleCommandHandler()
        command_bus.register(SampleCommand, handler)
        cmd = SampleCommand(value="hello")
        result = await command_bus.dispatch(cmd)
        assert result["handled"] == "hello"

    @pytest.mark.asyncio
    async def test_handler_not_found(self, command_bus):
        cmd = SampleCommand(value="test")
        with pytest.raises(CommandHandlerNotFoundError):
            await command_bus.dispatch(cmd)

    @pytest.mark.asyncio
    async def test_duplicate_registration(self, command_bus):
        handler = SampleCommandHandler()
        command_bus.register(SampleCommand, handler)
        with pytest.raises(ValueError, match="Handler already registered"):
            command_bus.register(SampleCommand, handler)

    @pytest.mark.asyncio
    async def test_multiple_handlers(self, command_bus):
        command_bus.register(SampleCommand, SampleCommandHandler())
        command_bus.register(AnotherCommand, AnotherCommandHandler())

        result1 = await command_bus.dispatch(SampleCommand(value="a"))
        assert result1["handled"] == "a"

        result2 = await command_bus.dispatch(AnotherCommand(number=99))
        assert result2["number"] == 99


class TestCommandBusMiddleware:
    """Test middleware chain in command bus."""

    @pytest.mark.asyncio
    async def test_middleware_before(self, command_bus):
        class LoggingMiddleware:
            def __init__(self):
                self.called = False

            async def before(self, command):
                self.called = True
                return command

            async def after(self, command, result):
                pass

        middleware = LoggingMiddleware()
        command_bus.add_middleware(middleware)
        command_bus.register(SampleCommand, SampleCommandHandler())

        await command_bus.dispatch(SampleCommand(value="test"))
        assert middleware.called

    @pytest.mark.asyncio
    async def test_middleware_after(self, command_bus):
        results = []

        class TrackMiddleware:
            async def before(self, command):
                return command

            async def after(self, command, result):
                results.append(("after", command.command_id, result))

        command_bus.add_middleware(TrackMiddleware())
        command_bus.register(SampleCommand, SampleCommandHandler())

        cmd = SampleCommand(value="tracked")
        await command_bus.dispatch(cmd)
        assert len(results) == 1
        assert results[0][0] == "after"
        assert results[0][1] == cmd.command_id

    @pytest.mark.asyncio
    async def test_middleware_modifies_command(self, command_bus):
        class ModifyMiddleware:
            async def before(self, command):
                return command.model_copy(update={"value": "modified"})

            async def after(self, command, result):
                pass

        command_bus.add_middleware(ModifyMiddleware())

        class CaptureHandler:
            async def handle(self, command):
                return {"value": command.value}

        command_bus.register(SampleCommand, CaptureHandler())
        result = await command_bus.dispatch(SampleCommand(value="original"))
        assert result["value"] == "modified"

    @pytest.mark.asyncio
    async def test_middleware_order(self, command_bus):
        order = []

        class FirstMiddleware:
            async def before(self, command):
                order.append("first_before")
                return command
            async def after(self, command, result):
                order.append("first_after")

        class SecondMiddleware:
            async def before(self, command):
                order.append("second_before")
                return command
            async def after(self, command, result):
                order.append("second_after")

        command_bus.add_middleware(FirstMiddleware())
        command_bus.add_middleware(SecondMiddleware())
        command_bus.register(SampleCommand, SampleCommandHandler())

        await command_bus.dispatch(SampleCommand(value="order"))
        assert order == ["first_before", "second_before", "second_after", "first_after"]


class TestQueryBus:
    """Test query registration and dispatch."""

    @pytest.mark.asyncio
    async def test_register_and_dispatch(self, query_bus):
        handler = SampleQueryHandler()
        query_bus.register(SampleQuery, handler)
        q = SampleQuery(query_id="q1")
        result = await query_bus.dispatch(q)
        assert result["result"] == "sample"
        assert result["query_id"] == "q1"

    @pytest.mark.asyncio
    async def test_handler_not_found(self, query_bus):
        q = SampleQuery()
        with pytest.raises(QueryHandlerNotFoundError):
            await query_bus.dispatch(q)

    @pytest.mark.asyncio
    async def test_duplicate_registration(self, query_bus):
        query_bus.register(SampleQuery, SampleQueryHandler())
        with pytest.raises(ValueError, match="Handler already registered"):
            query_bus.register(SampleQuery, SampleQueryHandler())

    @pytest.mark.asyncio
    async def test_multiple_query_types(self, query_bus):
        query_bus.register(SampleQuery, SampleQueryHandler())
        query_bus.register(AnotherQuery, AnotherQueryHandler())

        r1 = await query_bus.dispatch(SampleQuery(query_id="q1"))
        assert r1["query_id"] == "q1"

        r2 = await query_bus.dispatch(AnotherQuery(term="search"))
        assert r2["term"] == "search"

    @pytest.mark.asyncio
    async def test_dispatch_without_caching(self, query_bus):
        call_count = 0

        class CountingHandler:
            async def handle(self, query):
                nonlocal call_count
                call_count += 1
                return {"count": call_count}

        query_bus.register(SampleQuery, CountingHandler(), cache_ttl=0)
        q = SampleQuery(query_id="q1")
        r1 = await query_bus.dispatch(q)
        r2 = await query_bus.dispatch(q)
        assert r1["count"] == 1
        assert r2["count"] == 2


class TestQueryBusCaching:
    """Test query result caching."""

    @pytest.mark.asyncio
    async def test_cache_hit(self, query_bus_with_cache):
        call_count = 0

        class CountingHandler:
            async def handle(self, query):
                nonlocal call_count
                call_count += 1
                return {"count": call_count, "data": "cached"}

        query_bus_with_cache.register(SampleQuery, CountingHandler(), cache_ttl=60)
        q = SampleQuery(query_id="q1")

        r1 = await query_bus_with_cache.dispatch(q)
        assert r1["count"] == 1

        r2 = await query_bus_with_cache.dispatch(q)
        assert r2["count"] == 1
        assert r2["data"] == "cached"
        assert call_count == 1

    @pytest.mark.asyncio
    async def test_cache_miss(self, query_bus_with_cache):
        call_count = 0

        class CountingHandler:
            async def handle(self, query):
                nonlocal call_count
                call_count += 1
                return {"count": call_count}

        query_bus_with_cache.register(SampleQuery, CountingHandler(), cache_ttl=60)

        q1 = SampleQuery(query_id="q1")
        q2 = SampleQuery(query_id="q2")

        r1 = await query_bus_with_cache.dispatch(q1)
        assert r1["count"] == 1

        r2 = await query_bus_with_cache.dispatch(q2)
        assert r2["count"] == 2

    @pytest.mark.asyncio
    async def test_invalidate_cache(self, query_bus_with_cache):
        call_count = 0

        class CountingHandler:
            async def handle(self, query):
                nonlocal call_count
                call_count += 1
                return {"count": call_count}

        query_bus_with_cache.register(SampleQuery, CountingHandler(), cache_ttl=60)
        q = SampleQuery(query_id="invalidate")

        r1 = await query_bus_with_cache.dispatch(q)
        assert r1["count"] == 1

        await query_bus_with_cache.invalidate(q)

        r2 = await query_bus_with_cache.dispatch(q)
        assert r2["count"] == 2

    @pytest.mark.asyncio
    async def test_no_cache_when_ttl_zero(self, query_bus_with_cache):
        call_count = 0

        class CountingHandler:
            async def handle(self, query):
                nonlocal call_count
                call_count += 1
                return {"count": call_count}

        query_bus_with_cache.register(SampleQuery, CountingHandler(), cache_ttl=0)
        q = SampleQuery(query_id="no-cache")

        r1 = await query_bus_with_cache.dispatch(q)
        assert r1["count"] == 1

        r2 = await query_bus_with_cache.dispatch(q)
        assert r2["count"] == 2

    @pytest.mark.asyncio
    async def test_cache_key_uniqueness(self, query_bus_with_cache):
        query_bus_with_cache.register(SampleQuery, SampleQueryHandler(), cache_ttl=60)
        q1 = SampleQuery(query_id="id1")
        q2 = SampleQuery(query_id="id2")

        key1 = query_bus_with_cache._build_cache_key(q1)
        key2 = query_bus_with_cache._build_cache_key(q2)
        assert key1 != key2


class TestQueryBusWithoutRedis:
    """Test query bus behavior without Redis client."""

    @pytest.mark.asyncio
    async def test_dispatch_without_redis(self):
        bus = QueryBus(redis_client=None)
        bus.register(SampleQuery, SampleQueryHandler(), cache_ttl=60)

        result = await bus.dispatch(SampleQuery(query_id="no-redis"))
        assert result["result"] == "sample"

    @pytest.mark.asyncio
    async def test_invalidate_without_redis(self):
        bus = QueryBus(redis_client=None)
        bus.register(SampleQuery, SampleQueryHandler())

        await bus.invalidate(SampleQuery(query_id="x"))
