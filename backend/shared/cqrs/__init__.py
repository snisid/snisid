"""
SNISID CQRS Framework
======================
Command-Query Responsibility Segregation primitives for the SNISID platform.

Provides:
- **Command** / **CommandHandler** / **CommandBus** — write-side pipeline
  with automatic database transaction wrapping.
- **Query** / **QueryHandler** / **QueryBus** — read-side pipeline with
  optional Redis result caching.
- **DomainEvent** — canonical domain event envelope for Event Sourcing and
  cross-service messaging.

All buses use a handler-per-type registry to dispatch incoming messages.
"""
from __future__ import annotations

import asyncio
import hashlib
import uuid
from datetime import datetime, timezone
from typing import Any, Protocol, runtime_checkable

import orjson
from pydantic import BaseModel, Field

from shared.config import get_settings
from shared.logging import get_logger

logger = get_logger(__name__)


# ── Domain Event ──────────────────────────────────────────────────────


class DomainEvent(BaseModel):
    """
    Canonical domain event envelope.

    Every state change in the system is represented as a ``DomainEvent``.
    Events are immutable facts that are appended to the event store and
    published via Kafka for downstream consumers.

    Attributes:
        event_id:       Globally unique event identifier.
        event_type:     Qualified event name (e.g. ``CitizenRegistered``).
        aggregate_id:   Identity of the aggregate that produced the event.
        aggregate_type: Type name of the producing aggregate.
        timestamp:      UTC timestamp when the event was raised.
        version:        Aggregate version at event creation time.
        actor_id:       Identity of the user or system that triggered it.
        correlation_id: Correlation token for distributed tracing.
        data:           Arbitrary payload carrying the state delta.
    """

    event_id: str = Field(
        default_factory=lambda: str(uuid.uuid4()),
        description="Globally unique event identifier",
    )
    event_type: str = Field(
        ...,
        description="Qualified event type name, e.g. CitizenRegistered",
    )
    aggregate_id: str = Field(
        ...,
        description="Identity of the owning aggregate",
    )
    aggregate_type: str = Field(
        ...,
        description="Type name of the owning aggregate",
    )
    timestamp: datetime = Field(
        default_factory=lambda: datetime.now(timezone.utc),
        description="UTC timestamp when the event was created",
    )
    version: int = Field(
        default=1,
        ge=1,
        description="Aggregate version at the time of event creation",
    )
    actor_id: str | None = Field(
        default=None,
        description="Identity of the actor who triggered the event",
    )
    correlation_id: str | None = Field(
        default=None,
        description="Correlation token for distributed tracing",
    )
    data: dict[str, Any] = Field(
        default_factory=dict,
        description="Event payload carrying the state delta",
    )

    model_config = {"frozen": True}


# ── Command Side ──────────────────────────────────────────────────────


class Command(BaseModel):
    """
    Base class for all CQRS commands.

    Every command carries an identity, a creation timestamp, and the
    ``actor_id`` of the authenticated principal that issued it.
    """

    command_id: str = Field(
        default_factory=lambda: str(uuid.uuid4()),
        description="Unique command identifier",
    )
    timestamp: datetime = Field(
        default_factory=lambda: datetime.now(timezone.utc),
        description="UTC timestamp when the command was created",
    )
    actor_id: str | None = Field(
        default=None,
        description="Identity of the actor issuing the command",
    )

    model_config = {"frozen": True}


@runtime_checkable
class CommandHandler(Protocol):
    """
    Protocol for command handlers.

    Each concrete handler processes exactly one command type within a
    database transaction managed by the :class:`CommandBus`.
    """

    async def handle(self, command: Command) -> Any:
        """Execute the command and return an optional result."""
        ...


class CommandBus:
    """
    Dispatches commands to their registered handlers.

    The bus wraps every handler invocation in a database session
    transaction — the handler receives the session via its constructor
    or an injected dependency.  If the handler raises, the transaction
    is rolled back automatically by the session context manager.

    Usage::

        bus = CommandBus()
        bus.register(RegisterCitizenCommand, RegisterCitizenHandler())
        result = await bus.dispatch(command)
    """

    def __init__(self) -> None:
        self._handlers: dict[type[Command], CommandHandler] = {}
        self._middleware: list[CommandMiddleware] = []

    def register(
        self, command_type: type[Command], handler: CommandHandler
    ) -> None:
        """
        Register a handler for a specific command type.

        Args:
            command_type: The command class to handle.
            handler:      The handler instance.

        Raises:
            ValueError: If a handler is already registered for this type.
        """
        if command_type in self._handlers:
            raise ValueError(
                f"Handler already registered for {command_type.__name__}"
            )
        self._handlers[command_type] = handler
        logger.debug(
            "command_handler_registered",
            command_type=command_type.__name__,
            handler=type(handler).__name__,
        )

    def add_middleware(self, middleware: CommandMiddleware) -> None:
        """
        Append middleware to the command processing pipeline.

        Middleware is executed in registration order *before* the handler.

        Args:
            middleware: A callable implementing the :class:`CommandMiddleware`
                        protocol.
        """
        self._middleware.append(middleware)

    async def dispatch(self, command: Command) -> Any:
        """
        Dispatch a command to its registered handler.

        The command is first passed through all registered middleware,
        then forwarded to the handler.  The handler is expected to run
        inside a database transaction provided externally (e.g., via
        ``get_session()``).

        Args:
            command: The command to dispatch.

        Returns:
            The result produced by the handler, if any.

        Raises:
            CommandHandlerNotFoundError: If no handler is registered.
        """
        handler = self._handlers.get(type(command))
        if handler is None:
            raise CommandHandlerNotFoundError(
                f"No handler registered for command {type(command).__name__}"
            )

        logger.info(
            "command_dispatching",
            command_type=type(command).__name__,
            command_id=command.command_id,
            actor_id=command.actor_id,
        )

        # Run middleware chain
        for mw in self._middleware:
            command = await mw.before(command)

        from shared.database import get_session

        try:
            async with get_session() as session:
                # Attach session so handlers can access it if designed for DI
                result = await handler.handle(command)
            logger.info(
                "command_handled",
                command_type=type(command).__name__,
                command_id=command.command_id,
            )
        except Exception:
            logger.exception(
                "command_failed",
                command_type=type(command).__name__,
                command_id=command.command_id,
            )
            raise

        # Post-processing middleware (reverse order)
        for mw in reversed(self._middleware):
            await mw.after(command, result)

        return result


@runtime_checkable
class CommandMiddleware(Protocol):
    """Protocol for command bus middleware."""

    async def before(self, command: Command) -> Command:
        """Pre-process the command. Return (possibly modified) command."""
        ...

    async def after(self, command: Command, result: Any) -> None:
        """Post-process after successful handling."""
        ...


class CommandHandlerNotFoundError(Exception):
    """Raised when no handler is registered for a dispatched command."""


# ── Query Side ────────────────────────────────────────────────────────


class Query(BaseModel):
    """
    Base class for all CQRS queries.

    Queries are read-only operations that return projections or
    computed views of application state.
    """

    model_config = {"frozen": True}


@runtime_checkable
class QueryHandler(Protocol):
    """
    Protocol for query handlers.

    Each handler processes exactly one query type.
    """

    async def handle(self, query: Query) -> Any:
        """Execute the query and return its result."""
        ...


class QueryBus:
    """
    Dispatches queries to their registered handlers with optional
    Redis caching.

    Cached results are keyed by a deterministic hash of the query
    model's JSON representation.

    Usage::

        bus = QueryBus()
        bus.register(GetCitizenQuery, GetCitizenHandler(), cache_ttl=120)
        result = await bus.dispatch(query)
    """

    def __init__(self, redis_client: Any | None = None) -> None:
        """
        Args:
            redis_client: An ``redis.asyncio.Redis`` instance used for
                          query result caching.  If ``None``, caching is
                          disabled.
        """
        self._handlers: dict[type[Query], QueryHandler] = {}
        self._cache_ttl: dict[type[Query], int] = {}
        self._redis = redis_client

    def register(
        self,
        query_type: type[Query],
        handler: QueryHandler,
        cache_ttl: int = 0,
    ) -> None:
        """
        Register a handler for a specific query type.

        Args:
            query_type: The query class to handle.
            handler:    The handler instance.
            cache_ttl:  Cache lifetime in seconds.  ``0`` disables caching
                        for this query type.

        Raises:
            ValueError: If a handler is already registered for this type.
        """
        if query_type in self._handlers:
            raise ValueError(
                f"Handler already registered for {query_type.__name__}"
            )
        self._handlers[query_type] = handler
        if cache_ttl > 0:
            self._cache_ttl[query_type] = cache_ttl
        logger.debug(
            "query_handler_registered",
            query_type=query_type.__name__,
            handler=type(handler).__name__,
            cache_ttl=cache_ttl,
        )

    async def dispatch(self, query: Query) -> Any:
        """
        Dispatch a query to its registered handler.

        If a TTL is configured and a Redis client is available, the
        result is looked up in cache first.  On cache miss, the handler
        is invoked and the result is stored.

        Args:
            query: The query to dispatch.

        Returns:
            The query result (from cache or handler).

        Raises:
            QueryHandlerNotFoundError: If no handler is registered.
        """
        handler = self._handlers.get(type(query))
        if handler is None:
            raise QueryHandlerNotFoundError(
                f"No handler registered for query {type(query).__name__}"
            )

        cache_key: str | None = None
        ttl = self._cache_ttl.get(type(query), 0)

        # Attempt cache read
        if ttl > 0 and self._redis is not None:
            cache_key = self._build_cache_key(query)
            try:
                cached = await self._redis.get(cache_key)
                if cached is not None:
                    logger.debug(
                        "query_cache_hit",
                        query_type=type(query).__name__,
                        cache_key=cache_key,
                    )
                    return orjson.loads(cached)
            except Exception:
                logger.warning(
                    "query_cache_read_failed",
                    query_type=type(query).__name__,
                    cache_key=cache_key,
                    exc_info=True,
                )

        logger.debug(
            "query_dispatching",
            query_type=type(query).__name__,
        )

        result = await handler.handle(query)

        # Populate cache
        if cache_key is not None and ttl > 0 and self._redis is not None:
            try:
                serialized = orjson.dumps(
                    result,
                    option=orjson.OPT_NAIVE_UTC | orjson.OPT_SERIALIZE_UUID,
                )
                await self._redis.set(cache_key, serialized, ex=ttl)
                logger.debug(
                    "query_cache_set",
                    query_type=type(query).__name__,
                    cache_key=cache_key,
                    ttl=ttl,
                )
            except Exception:
                logger.warning(
                    "query_cache_write_failed",
                    query_type=type(query).__name__,
                    cache_key=cache_key,
                    exc_info=True,
                )

        return result

    async def invalidate(self, query: Query) -> None:
        """
        Explicitly remove a cached query result.

        Args:
            query: A query instance whose cached result should be evicted.
        """
        if self._redis is None:
            return
        cache_key = self._build_cache_key(query)
        try:
            await self._redis.delete(cache_key)
            logger.debug(
                "query_cache_invalidated",
                query_type=type(query).__name__,
                cache_key=cache_key,
            )
        except Exception:
            logger.warning(
                "query_cache_invalidation_failed",
                cache_key=cache_key,
                exc_info=True,
            )

    # ── Internals ─────────────────────────────────────────────────────

    @staticmethod
    def _build_cache_key(query: Query) -> str:
        """
        Build a deterministic Redis key from the query payload.

        Uses a SHA-256 digest of the sorted JSON representation to
        guarantee consistent keys regardless of field ordering.
        """
        payload = orjson.dumps(
            query.model_dump(mode="json"),
            option=orjson.OPT_SORT_KEYS,
        )
        digest = hashlib.sha256(payload).hexdigest()
        return f"snisid:query_cache:{type(query).__name__}:{digest}"


class QueryHandlerNotFoundError(Exception):
    """Raised when no handler is registered for a dispatched query."""


# ── Public API ────────────────────────────────────────────────────────

__all__ = [
    "Command",
    "CommandBus",
    "CommandHandler",
    "CommandHandlerNotFoundError",
    "CommandMiddleware",
    "DomainEvent",
    "Query",
    "QueryBus",
    "QueryHandler",
    "QueryHandlerNotFoundError",
]
