"""
SNISID Event Store — Append-Only Event Persistence
====================================================
Provides the immutable event store for Event Sourcing.
Events are append-only with optimistic concurrency control via version checks.
Supports snapshot storage for aggregate rebuild optimization.
"""
from __future__ import annotations

import uuid
from datetime import datetime, timezone
from typing import Any, Sequence

from sqlalchemy import (
    BigInteger,
    DateTime,
    Index,
    Integer,
    String,
    Text,
    select,
    func,
)
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import Mapped, mapped_column

from shared.database import Base
from shared.logging import get_logger

logger = get_logger(__name__)


class StoredEvent(Base):
    """
    Immutable event record in the event store.
    Each event represents a state change to an aggregate.
    """

    __tablename__ = "event_store"

    # Override base id with BigInt auto-increment for ordering
    id: Mapped[int] = mapped_column(  # type: ignore[assignment]
        BigInteger, primary_key=True, autoincrement=True
    )
    event_id: Mapped[str] = mapped_column(
        String(36), unique=True, nullable=False, default=lambda: str(uuid.uuid4())
    )

    # Aggregate identity
    aggregate_id: Mapped[str] = mapped_column(String(36), nullable=False, index=True)
    aggregate_type: Mapped[str] = mapped_column(String(100), nullable=False, index=True)

    # Event metadata
    event_type: Mapped[str] = mapped_column(String(200), nullable=False, index=True)
    event_data: Mapped[dict[str, Any]] = mapped_column(JSONB, nullable=False)
    event_metadata: Mapped[dict[str, Any] | None] = mapped_column(JSONB, nullable=True)

    # Versioning for optimistic concurrency
    version: Mapped[int] = mapped_column(Integer, nullable=False)

    # Audit fields
    actor_id: Mapped[str | None] = mapped_column(String(36), nullable=True)
    correlation_id: Mapped[str | None] = mapped_column(String(36), nullable=True, index=True)
    causation_id: Mapped[str | None] = mapped_column(String(36), nullable=True)
    timestamp: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        nullable=False,
        default=lambda: datetime.now(timezone.utc),
    )

    __table_args__ = (
        # Unique constraint: each aggregate has unique version numbers
        Index("uq_aggregate_version", "aggregate_id", "version", unique=True),
        # Efficient time-range queries
        Index("ix_event_store_timestamp", "timestamp"),
        # Aggregate type + id for type-scoped queries
        Index("ix_event_store_type_id", "aggregate_type", "aggregate_id"),
    )


class AggregateSnapshot(Base):
    """
    Snapshot of an aggregate state for fast reconstitution.
    Periodically created to avoid replaying full event history.
    """

    __tablename__ = "aggregate_snapshots"

    aggregate_id: Mapped[str] = mapped_column(String(36), nullable=False, index=True)
    aggregate_type: Mapped[str] = mapped_column(String(100), nullable=False)
    version: Mapped[int] = mapped_column(Integer, nullable=False)
    state: Mapped[dict[str, Any]] = mapped_column(JSONB, nullable=False)
    snapshot_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        nullable=False,
        default=lambda: datetime.now(timezone.utc),
    )

    __table_args__ = (
        Index("uq_snapshot_aggregate", "aggregate_id", "aggregate_type", unique=True),
    )


class EventStore:
    """
    Event Store providing append-only event persistence with
    optimistic concurrency and snapshot support.
    """

    SNAPSHOT_INTERVAL: int = 50  # Create snapshot every N events

    def __init__(self, session: AsyncSession) -> None:
        self._session = session

    async def append_events(
        self,
        aggregate_id: str,
        aggregate_type: str,
        events: list[dict[str, Any]],
        expected_version: int,
        actor_id: str | None = None,
        correlation_id: str | None = None,
    ) -> list[StoredEvent]:
        """
        Append events to the store with optimistic concurrency control.

        Args:
            aggregate_id: ID of the aggregate the events belong to.
            aggregate_type: Type name of the aggregate.
            events: List of event dicts with 'event_type' and 'event_data'.
            expected_version: Expected current version (for concurrency check).
            actor_id: ID of the user/system performing the action.
            correlation_id: Correlation ID for distributed tracing.

        Returns:
            List of persisted StoredEvent instances.

        Raises:
            ConcurrencyError: If expected_version doesn't match current version.
        """
        # Check current version
        current_version = await self._get_current_version(aggregate_id)
        if current_version != expected_version:
            raise ConcurrencyError(
                f"Concurrency conflict for aggregate {aggregate_id}: "
                f"expected version {expected_version}, found {current_version}"
            )

        stored_events: list[StoredEvent] = []
        version = expected_version

        for event_data in events:
            version += 1
            stored_event = StoredEvent(
                event_id=str(uuid.uuid4()),
                aggregate_id=aggregate_id,
                aggregate_type=aggregate_type,
                event_type=event_data["event_type"],
                event_data=event_data.get("data", {}),
                event_metadata=event_data.get("metadata"),
                version=version,
                actor_id=actor_id,
                correlation_id=correlation_id,
                causation_id=event_data.get("causation_id"),
                timestamp=datetime.now(timezone.utc),
            )
            self._session.add(stored_event)
            stored_events.append(stored_event)

        await self._session.flush()

        logger.info(
            "events_appended",
            aggregate_id=aggregate_id,
            aggregate_type=aggregate_type,
            event_count=len(events),
            new_version=version,
        )

        return stored_events

    async def get_events(
        self,
        aggregate_id: str,
        after_version: int = 0,
        limit: int | None = None,
    ) -> Sequence[StoredEvent]:
        """
        Retrieve events for an aggregate, optionally after a version.

        Args:
            aggregate_id: ID of the aggregate.
            after_version: Only return events with version > this.
            limit: Maximum number of events to return.

        Returns:
            Ordered sequence of StoredEvent instances.
        """
        stmt = (
            select(StoredEvent)
            .where(
                StoredEvent.aggregate_id == aggregate_id,
                StoredEvent.version > after_version,
            )
            .order_by(StoredEvent.version.asc())
        )
        if limit is not None:
            stmt = stmt.limit(limit)

        result = await self._session.execute(stmt)
        return result.scalars().all()

    async def get_events_by_type(
        self,
        aggregate_type: str,
        after_timestamp: datetime | None = None,
        event_types: list[str] | None = None,
        limit: int = 1000,
    ) -> Sequence[StoredEvent]:
        """Retrieve events by aggregate type, optionally filtered."""
        stmt = (
            select(StoredEvent)
            .where(StoredEvent.aggregate_type == aggregate_type)
            .order_by(StoredEvent.timestamp.asc())
            .limit(limit)
        )
        if after_timestamp:
            stmt = stmt.where(StoredEvent.timestamp > after_timestamp)
        if event_types:
            stmt = stmt.where(StoredEvent.event_type.in_(event_types))

        result = await self._session.execute(stmt)
        return result.scalars().all()

    async def get_events_by_correlation(
        self, correlation_id: str
    ) -> Sequence[StoredEvent]:
        """Retrieve all events with a given correlation ID."""
        stmt = (
            select(StoredEvent)
            .where(StoredEvent.correlation_id == correlation_id)
            .order_by(StoredEvent.timestamp.asc())
        )
        result = await self._session.execute(stmt)
        return result.scalars().all()

    async def save_snapshot(
        self,
        aggregate_id: str,
        aggregate_type: str,
        version: int,
        state: dict[str, Any],
    ) -> AggregateSnapshot:
        """Save or update a snapshot of an aggregate's state."""
        # Upsert: delete existing snapshot, insert new one
        existing = await self._session.execute(
            select(AggregateSnapshot).where(
                AggregateSnapshot.aggregate_id == aggregate_id,
                AggregateSnapshot.aggregate_type == aggregate_type,
            )
        )
        existing_snapshot = existing.scalar_one_or_none()
        if existing_snapshot:
            existing_snapshot.version = version
            existing_snapshot.state = state
            existing_snapshot.snapshot_at = datetime.now(timezone.utc)
            return existing_snapshot

        snapshot = AggregateSnapshot(
            aggregate_id=aggregate_id,
            aggregate_type=aggregate_type,
            version=version,
            state=state,
        )
        self._session.add(snapshot)
        await self._session.flush()
        return snapshot

    async def get_snapshot(
        self, aggregate_id: str, aggregate_type: str
    ) -> AggregateSnapshot | None:
        """Get the latest snapshot for an aggregate."""
        result = await self._session.execute(
            select(AggregateSnapshot).where(
                AggregateSnapshot.aggregate_id == aggregate_id,
                AggregateSnapshot.aggregate_type == aggregate_type,
            )
        )
        return result.scalar_one_or_none()

    async def _get_current_version(self, aggregate_id: str) -> int:
        """Get the current (latest) version number for an aggregate."""
        result = await self._session.execute(
            select(func.max(StoredEvent.version)).where(
                StoredEvent.aggregate_id == aggregate_id
            )
        )
        version = result.scalar_one_or_none()
        return version or 0

    async def get_aggregate_ids(
        self, aggregate_type: str, limit: int = 1000, offset: int = 0
    ) -> Sequence[str]:
        """Get all aggregate IDs of a given type."""
        result = await self._session.execute(
            select(StoredEvent.aggregate_id)
            .where(StoredEvent.aggregate_type == aggregate_type)
            .distinct()
            .limit(limit)
            .offset(offset)
        )
        return result.scalars().all()


class ConcurrencyError(Exception):
    """Raised when optimistic concurrency check fails."""
    pass
