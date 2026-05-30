"""
SNISID Aggregate Root — Event-Sourced Aggregate Base
=====================================================
Provides the ``AggregateRoot`` base class for domain aggregates that
derive their state exclusively from a stream of domain events.

An aggregate:

1. Receives a **command** via a public method.
2. Validates business invariants.
3. Calls :meth:`raise_event` to record one or more domain events.
4. The events are applied immediately (via ``_when_<EventType>`` mutators)
   and queued for persistence in the event store.

Supports:
- **Event replay** — :meth:`load_from_history` replays persisted events.
- **Snapshots** — :meth:`take_snapshot` / :meth:`load_from_snapshot` for
  large event streams.
"""
from __future__ import annotations

import copy
from typing import Any

from shared.cqrs import DomainEvent
from shared.logging import get_logger

logger = get_logger(__name__)


class AggregateRoot:
    """
    Base class for event-sourced aggregates.

    Subclasses must implement ``_when_<EventType>`` methods for every
    event type they produce.  For example, if the aggregate raises a
    ``CitizenRegistered`` event the subclass must define::

        def _when_CitizenRegistered(self, event: DomainEvent) -> None:
            self._name = event.data["name"]

    Attributes:
        _id:              Aggregate identity (set on first event).
        _version:         Current event-stream version (0 = new).
        _pending_events:  Events raised but not yet persisted.
    """

    def __init__(self, aggregate_id: str | None = None) -> None:
        """
        Initialise a new or empty aggregate.

        Args:
            aggregate_id: Optional explicit ID.  If *None*, the ID is
                          expected to be set by the first applied event.
        """
        self._id: str = aggregate_id or ""
        self._version: int = 0
        self._pending_events: list[DomainEvent] = []

    # ── Properties ────────────────────────────────────────────────────

    @property
    def id(self) -> str:
        """Aggregate identity."""
        return self._id

    @property
    def version(self) -> int:
        """Current aggregate version (number of applied events)."""
        return self._version

    # ── Event Sourcing Core ───────────────────────────────────────────

    def raise_event(self, event: DomainEvent) -> None:
        """
        Record a new domain event, apply it to mutate state, and queue
        it for persistence.

        The event version is overridden to ``_version + 1`` to keep the
        aggregate stream strictly sequential.

        Args:
            event: The domain event to record.
        """
        self._version += 1

        # Override version to match aggregate's sequential stream
        event = event.model_copy(update={"version": self._version})

        self._apply(event)
        self._pending_events.append(event)

        logger.debug(
            "domain_event_raised",
            aggregate_id=self._id,
            aggregate_type=type(self).__name__,
            event_type=event.event_type,
            version=self._version,
        )

    def _apply(self, event: DomainEvent) -> None:
        """
        Dispatch the event to the appropriate ``_when_<EventType>``
        handler on the concrete subclass.

        The method name is derived from :attr:`DomainEvent.event_type`.
        If no handler exists the event is silently ignored (forward
        compatibility).

        Args:
            event: The domain event to apply.

        Raises:
            AggregateEventHandlerError: If the handler raises an
                                        unexpected exception.
        """
        handler_name = f"_when_{event.event_type}"
        handler = getattr(self, handler_name, None)

        if handler is None:
            logger.warning(
                "aggregate_event_handler_missing",
                aggregate_type=type(self).__name__,
                event_type=event.event_type,
                handler_name=handler_name,
            )
            return

        try:
            handler(event)
        except Exception as exc:
            raise AggregateEventHandlerError(
                f"Error applying {event.event_type} on "
                f"{type(self).__name__}({self._id}): {exc}"
            ) from exc

    # ── History Replay ────────────────────────────────────────────────

    def load_from_history(self, events: list[DomainEvent]) -> None:
        """
        Replay persisted events to reconstitute aggregate state.

        Events are applied in order **without** being added to the
        pending list (they are already persisted).

        Args:
            events: Ordered list of historical events from the event
                    store.
        """
        for event in events:
            self._apply(event)
            self._version = event.version

        logger.debug(
            "aggregate_loaded_from_history",
            aggregate_id=self._id,
            aggregate_type=type(self).__name__,
            version=self._version,
            event_count=len(events),
        )

    # ── Snapshot Support ──────────────────────────────────────────────

    def load_from_snapshot(
        self,
        snapshot_data: dict[str, Any],
        version: int,
    ) -> None:
        """
        Restore aggregate state from a snapshot.

        After loading, subsequent events (after the snapshot version)
        should be replayed via :meth:`load_from_history`.

        Args:
            snapshot_data: Dictionary of serialised aggregate state
                          (produced by :meth:`take_snapshot`).
            version:       Aggregate version at snapshot creation time.
        """
        self._id = snapshot_data.get("id", self._id)
        self._version = version

        # Let subclasses restore domain-specific state
        self._restore_from_snapshot(snapshot_data)

        logger.debug(
            "aggregate_loaded_from_snapshot",
            aggregate_id=self._id,
            aggregate_type=type(self).__name__,
            version=version,
        )

    def take_snapshot(self) -> dict[str, Any]:
        """
        Serialise the current aggregate state into a dictionary.

        Subclasses should override :meth:`_snapshot_state` to include
        domain-specific fields.  The base implementation captures ``id``
        and ``version``.

        Returns:
            A JSON-serialisable dictionary representing the aggregate.
        """
        state = {
            "id": self._id,
            "version": self._version,
            "aggregate_type": type(self).__name__,
        }
        domain_state = self._snapshot_state()
        state.update(domain_state)

        logger.debug(
            "aggregate_snapshot_taken",
            aggregate_id=self._id,
            aggregate_type=type(self).__name__,
            version=self._version,
        )
        return state

    def _snapshot_state(self) -> dict[str, Any]:
        """
        Return domain-specific state for snapshot serialisation.

        Override in subclasses to persist fields beyond ``id`` and
        ``version``.  The default implementation captures all public
        attributes (non-underscore, non-callable) via shallow copy.

        Returns:
            A dictionary of serialisable domain state.
        """
        state: dict[str, Any] = {}
        for attr_name in vars(self):
            # Skip private/internal attributes already handled
            if attr_name.startswith("_"):
                continue
            value = getattr(self, attr_name)
            if not callable(value):
                state[attr_name] = copy.deepcopy(value)
        return state

    def _restore_from_snapshot(self, snapshot_data: dict[str, Any]) -> None:
        """
        Restore domain-specific state from a snapshot dictionary.

        Override in subclasses to restore fields captured by
        :meth:`_snapshot_state`.  The default implementation sets every
        key from *snapshot_data* (except ``id``, ``version``,
        ``aggregate_type``) as an attribute on the instance.

        Args:
            snapshot_data: The snapshot dictionary.
        """
        reserved = {"id", "version", "aggregate_type"}
        for key, value in snapshot_data.items():
            if key not in reserved:
                setattr(self, key, copy.deepcopy(value))

    # ── Pending Events ────────────────────────────────────────────────

    def get_pending_events(self) -> list[DomainEvent]:
        """
        Return events that have been raised but not yet persisted.

        Returns:
            A *copy* of the pending events list.
        """
        return list(self._pending_events)

    def clear_pending_events(self) -> None:
        """
        Clear the pending events list after successful persistence.

        This should be called by the repository / unit-of-work once
        events have been written to the event store and published.
        """
        count = len(self._pending_events)
        self._pending_events.clear()
        logger.debug(
            "pending_events_cleared",
            aggregate_id=self._id,
            aggregate_type=type(self).__name__,
            cleared_count=count,
        )

    # ── Equality ──────────────────────────────────────────────────────

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, AggregateRoot):
            return NotImplemented
        return self._id == other._id and type(self) is type(other)

    def __hash__(self) -> int:
        return hash((type(self).__name__, self._id))

    def __repr__(self) -> str:
        return (
            f"<{type(self).__name__}(id={self._id!r}, "
            f"version={self._version})>"
        )


class AggregateEventHandlerError(Exception):
    """Raised when an aggregate ``_when_*`` mutator fails."""


__all__ = [
    "AggregateEventHandlerError",
    "AggregateRoot",
]
