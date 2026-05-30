"""
SNISID Event Infrastructure — Kafka Producer, Consumer & In-Process EventBus
==============================================================================
Async Kafka messaging layer built on ``aiokafka`` with:

- **KafkaProducer** — publish domain events with automatic JSON
  serialisation, retry, batching, and transactional-outbox support.
- **KafkaConsumer** — consume events with handler dispatch, idempotency
  via Redis, dead-letter queue routing, and semaphore-based backpressure.
- **EventBus** — lightweight in-process pub/sub for local event dispatch
  within a single service.

All classes rely on ``shared.logging.get_logger`` for structured logging
and ``shared.config.get_settings`` for runtime configuration.
"""
from __future__ import annotations

import asyncio
import time
import uuid
from collections import defaultdict
from datetime import datetime, timezone
from typing import Any, Awaitable, Callable

import orjson
from pydantic import BaseModel

from shared.config import Settings, get_settings
from shared.cqrs import DomainEvent
from shared.logging import get_logger

logger = get_logger(__name__)

# Type alias for event handler callables
EventHandler = Callable[[DomainEvent], Awaitable[None]]


# ── Kafka Producer ────────────────────────────────────────────────────


class KafkaProducer:
    """
    Async Kafka producer that publishes :class:`DomainEvent` instances
    to topics with automatic JSON serialisation, retries, and optional
    transactional-outbox persistence.

    Usage::

        producer = KafkaProducer(
            bootstrap_servers="kafka:9092",
            topic_prefix="snisid",
        )
        await producer.start()
        await producer.publish("citizens", key=citizen_id, event=event)
        await producer.stop()

    Args:
        bootstrap_servers: Comma-separated Kafka broker addresses.
        topic_prefix:      Prefix prepended to topic names
                           (e.g. ``snisid.citizens``).
        settings:          Optional :class:`Settings` override. When
                           *None*, ``get_settings()`` is called lazily.
    """

    def __init__(
        self,
        bootstrap_servers: str | None = None,
        topic_prefix: str | None = None,
        settings: Settings | None = None,
    ) -> None:
        self._settings = settings or get_settings()
        kafka_cfg = self._settings.kafka

        self._bootstrap_servers = (
            bootstrap_servers or kafka_cfg.bootstrap_servers
        )
        self._topic_prefix = topic_prefix or kafka_cfg.topic_prefix
        self._producer: Any | None = None  # aiokafka.AIOKafkaProducer

        # Producer tuning from config
        self._acks = kafka_cfg.producer_acks
        self._retries = kafka_cfg.producer_retries
        self._batch_size = kafka_cfg.producer_batch_size
        self._linger_ms = kafka_cfg.producer_linger_ms

    async def start(self) -> None:
        """
        Create and start the underlying ``AIOKafkaProducer``.

        Raises:
            RuntimeError: If the producer is already running.
        """
        if self._producer is not None:
            raise RuntimeError("KafkaProducer is already started")

        from aiokafka import AIOKafkaProducer

        self._producer = AIOKafkaProducer(
            bootstrap_servers=self._bootstrap_servers,
            value_serializer=self._serialize_value,
            key_serializer=self._serialize_key,
            acks=self._acks,
            retry_backoff_ms=200,
            max_batch_size=self._batch_size,
            linger_ms=self._linger_ms,
            enable_idempotence=True,
        )
        await self._producer.start()
        logger.info(
            "kafka_producer_started",
            bootstrap_servers=self._bootstrap_servers,
            topic_prefix=self._topic_prefix,
        )

    async def stop(self) -> None:
        """Flush pending messages and close the producer."""
        if self._producer is not None:
            await self._producer.stop()
            self._producer = None
            logger.info("kafka_producer_stopped")

    async def publish(
        self,
        topic: str,
        key: str,
        event: DomainEvent,
        headers: dict[str, str] | None = None,
    ) -> None:
        """
        Publish a single domain event to a Kafka topic.

        The event is serialised to JSON via ``orjson``.  The method
        retries up to ``producer_retries`` times on transient failures.

        Args:
            topic:   Topic name (will be prefixed automatically).
            key:     Partition key (typically the aggregate ID).
            event:   The domain event to publish.
            headers: Optional Kafka headers as string key-value pairs.

        Raises:
            RuntimeError: If the producer has not been started.
            KafkaPublishError: After exhausting retries.
        """
        if self._producer is None:
            raise RuntimeError(
                "KafkaProducer not started. Call start() first."
            )

        full_topic = self._full_topic(topic)
        kafka_headers = self._encode_headers(headers, event)

        last_exc: Exception | None = None
        for attempt in range(1, self._retries + 1):
            try:
                await self._producer.send_and_wait(
                    topic=full_topic,
                    key=key,
                    value=event.model_dump(mode="json"),
                    headers=kafka_headers,
                )
                logger.info(
                    "event_published",
                    topic=full_topic,
                    event_type=event.event_type,
                    event_id=event.event_id,
                    key=key,
                )
                return
            except Exception as exc:
                last_exc = exc
                wait = min(0.2 * (2 ** (attempt - 1)), 5.0)
                logger.warning(
                    "event_publish_retry",
                    topic=full_topic,
                    event_id=event.event_id,
                    attempt=attempt,
                    max_retries=self._retries,
                    wait_seconds=wait,
                    error=str(exc),
                )
                await asyncio.sleep(wait)

        raise KafkaPublishError(
            f"Failed to publish event {event.event_id} to {full_topic} "
            f"after {self._retries} attempts: {last_exc}"
        )

    async def publish_batch(
        self,
        topic: str,
        events: list[tuple[str, DomainEvent]],
    ) -> int:
        """
        Publish multiple events in a single batch.

        Each item in *events* is a ``(key, event)`` tuple.

        Args:
            topic:  Topic name (will be prefixed).
            events: List of ``(partition_key, DomainEvent)`` tuples.

        Returns:
            Number of events successfully published.

        Raises:
            RuntimeError: If the producer has not been started.
        """
        if self._producer is None:
            raise RuntimeError(
                "KafkaProducer not started. Call start() first."
            )

        full_topic = self._full_topic(topic)
        published = 0
        batch = self._producer.create_batch()

        for key, event in events:
            payload = self._serialize_value(event.model_dump(mode="json"))
            kafka_headers = self._encode_headers(None, event)
            metadata = batch.append(
                key=self._serialize_key(key),
                value=payload,
                timestamp=None,
                headers=kafka_headers,
            )
            if metadata is None:
                # Batch full — send current batch and create a new one
                await self._producer.send_batch(
                    batch, full_topic, partition=None
                )
                logger.debug(
                    "event_batch_flushed",
                    topic=full_topic,
                    count=published,
                )
                batch = self._producer.create_batch()
                batch.append(
                    key=self._serialize_key(key),
                    value=payload,
                    timestamp=None,
                    headers=kafka_headers,
                )
            published += 1

        # Send remaining
        if published > 0:
            await self._producer.send_batch(
                batch, full_topic, partition=None
            )

        logger.info(
            "event_batch_published",
            topic=full_topic,
            count=published,
        )
        return published

    async def publish_to_outbox(
        self,
        session: Any,
        topic: str,
        key: str,
        event: DomainEvent,
    ) -> None:
        """
        Persist an event to the transactional outbox table within the
        current database transaction.

        A separate relay process polls the outbox and publishes pending
        messages to Kafka, ensuring exactly-once delivery semantics.

        Args:
            session: An active ``AsyncSession`` within a transaction.
            topic:   Destination topic name.
            key:     Partition key.
            event:   The domain event to persist.
        """
        from sqlalchemy import text as sa_text

        full_topic = self._full_topic(topic)
        payload = orjson.dumps(
            event.model_dump(mode="json"),
            option=orjson.OPT_NAIVE_UTC | orjson.OPT_SERIALIZE_UUID,
        ).decode("utf-8")

        await session.execute(
            sa_text(
                """
                INSERT INTO outbox_events
                    (id, topic, partition_key, payload, event_type,
                     event_id, created_at, published_at)
                VALUES
                    (:id, :topic, :partition_key, :payload, :event_type,
                     :event_id, :created_at, NULL)
                """
            ),
            {
                "id": str(uuid.uuid4()),
                "topic": full_topic,
                "partition_key": key,
                "payload": payload,
                "event_type": event.event_type,
                "event_id": event.event_id,
                "created_at": datetime.now(timezone.utc).isoformat(),
            },
        )
        logger.debug(
            "event_persisted_to_outbox",
            topic=full_topic,
            event_id=event.event_id,
            event_type=event.event_type,
        )

    # ── Internals ─────────────────────────────────────────────────────

    def _full_topic(self, topic: str) -> str:
        """Prepend the configured prefix to the topic name."""
        return f"{self._topic_prefix}.{topic}"

    @staticmethod
    def _serialize_value(value: Any) -> bytes:
        """Serialise a value to JSON bytes using orjson."""
        if isinstance(value, bytes):
            return value
        return orjson.dumps(
            value, option=orjson.OPT_NAIVE_UTC | orjson.OPT_SERIALIZE_UUID
        )

    @staticmethod
    def _serialize_key(key: str) -> bytes:
        """Encode partition key to UTF-8 bytes."""
        if isinstance(key, bytes):
            return key
        return key.encode("utf-8")

    @staticmethod
    def _encode_headers(
        headers: dict[str, str] | None,
        event: DomainEvent,
    ) -> list[tuple[str, bytes]]:
        """Build Kafka record headers including standard SNISID fields."""
        kafka_headers: list[tuple[str, bytes]] = [
            ("event_type", event.event_type.encode("utf-8")),
            ("event_id", event.event_id.encode("utf-8")),
            ("aggregate_id", event.aggregate_id.encode("utf-8")),
            ("aggregate_type", event.aggregate_type.encode("utf-8")),
            ("timestamp", event.timestamp.isoformat().encode("utf-8")),
        ]
        if event.correlation_id:
            kafka_headers.append(
                ("correlation_id", event.correlation_id.encode("utf-8"))
            )
        if headers:
            kafka_headers.extend(
                (k, v.encode("utf-8")) for k, v in headers.items()
            )
        return kafka_headers


# ── Kafka Consumer ────────────────────────────────────────────────────


class KafkaConsumer:
    """
    Async Kafka consumer with handler dispatch, idempotency checks,
    dead-letter queue routing, and semaphore-based backpressure.

    Usage::

        consumer = KafkaConsumer(
            bootstrap_servers="kafka:9092",
            group_id="citizens-service",
            topics=["snisid.citizens"],
            handler_map={
                "CitizenRegistered": handle_citizen_registered,
                "CitizenUpdated": handle_citizen_updated,
            },
        )
        await consumer.start()

    Args:
        bootstrap_servers: Comma-separated Kafka broker addresses.
        group_id:          Consumer group identifier.
        topics:            List of topics to subscribe to.
        handler_map:       Mapping of ``event_type`` → async handler.
        redis_client:      Optional ``redis.asyncio.Redis`` for idempotency.
        max_concurrency:   Maximum concurrent handler invocations.
        max_retries:       Handler retries before DLQ routing.
        dlq_suffix:        Suffix for dead-letter topic names.
        settings:          Optional ``Settings`` override.
    """

    DEFAULT_MAX_CONCURRENCY: int = 10
    DEFAULT_MAX_RETRIES: int = 3
    DEFAULT_DLQ_SUFFIX: str = ".dlq"

    def __init__(
        self,
        bootstrap_servers: str | None = None,
        group_id: str | None = None,
        topics: list[str] | None = None,
        handler_map: dict[str, EventHandler] | None = None,
        redis_client: Any | None = None,
        max_concurrency: int = DEFAULT_MAX_CONCURRENCY,
        max_retries: int = DEFAULT_MAX_RETRIES,
        dlq_suffix: str = DEFAULT_DLQ_SUFFIX,
        settings: Settings | None = None,
    ) -> None:
        self._settings = settings or get_settings()
        kafka_cfg = self._settings.kafka

        self._bootstrap_servers = (
            bootstrap_servers or kafka_cfg.bootstrap_servers
        )
        self._group_id = group_id or f"{kafka_cfg.consumer_group_prefix}-default"
        self._topics = topics or []
        self._handler_map: dict[str, EventHandler] = handler_map or {}
        self._redis = redis_client
        self._max_retries = max_retries
        self._dlq_suffix = dlq_suffix

        self._consumer: Any | None = None  # aiokafka.AIOKafkaConsumer
        self._dlq_producer: KafkaProducer | None = None
        self._semaphore = asyncio.Semaphore(max_concurrency)
        self._running = False
        self._consume_task: asyncio.Task[None] | None = None

        # Redis key prefix for idempotency
        self._idempotency_prefix = "snisid:event_idempotency"
        # How long to remember processed event IDs (seconds)
        self._idempotency_ttl: int = 7 * 24 * 3600  # 7 days

    async def start(self) -> None:
        """
        Create the underlying consumer, subscribe to topics, and start
        the background consumption loop.

        Raises:
            RuntimeError: If the consumer is already running.
        """
        if self._running:
            raise RuntimeError("KafkaConsumer is already running")

        from aiokafka import AIOKafkaConsumer

        self._consumer = AIOKafkaConsumer(
            *self._topics,
            bootstrap_servers=self._bootstrap_servers,
            group_id=self._group_id,
            auto_offset_reset=self._settings.kafka.consumer_auto_offset_reset,
            enable_auto_commit=False,
            max_poll_records=self._settings.kafka.consumer_max_poll_records,
            session_timeout_ms=self._settings.kafka.consumer_session_timeout_ms,
            value_deserializer=self._deserialize_value,
            key_deserializer=self._deserialize_key,
        )

        # Start a DLQ producer for routing failed messages
        self._dlq_producer = KafkaProducer(
            bootstrap_servers=self._bootstrap_servers,
            topic_prefix="",  # DLQ topics use the original full name + suffix
            settings=self._settings,
        )
        await self._dlq_producer.start()

        await self._consumer.start()
        self._running = True
        self._consume_task = asyncio.create_task(
            self._consume_loop(), name="kafka-consumer-loop"
        )

        logger.info(
            "kafka_consumer_started",
            bootstrap_servers=self._bootstrap_servers,
            group_id=self._group_id,
            topics=self._topics,
        )

    async def stop(self) -> None:
        """Stop the consumer, cancel the background loop, and clean up."""
        self._running = False
        if self._consume_task is not None:
            self._consume_task.cancel()
            try:
                await self._consume_task
            except asyncio.CancelledError:
                pass
            self._consume_task = None

        if self._consumer is not None:
            await self._consumer.stop()
            self._consumer = None

        if self._dlq_producer is not None:
            await self._dlq_producer.stop()
            self._dlq_producer = None

        logger.info("kafka_consumer_stopped", group_id=self._group_id)

    def register_handler(
        self, event_type: str, handler: EventHandler
    ) -> None:
        """
        Register a handler for a specific event type at runtime.

        Args:
            event_type: The ``event_type`` string to match.
            handler:    Async callable accepting a ``DomainEvent``.
        """
        self._handler_map[event_type] = handler
        logger.debug(
            "consumer_handler_registered",
            event_type=event_type,
            group_id=self._group_id,
        )

    # ── Consumption Loop ──────────────────────────────────────────────

    async def _consume_loop(self) -> None:
        """Background task that polls Kafka and dispatches messages."""
        assert self._consumer is not None

        while self._running:
            try:
                records = await self._consumer.getmany(
                    timeout_ms=1000,
                    max_records=self._settings.kafka.consumer_max_poll_records,
                )
                for tp, messages in records.items():
                    for msg in messages:
                        await self._semaphore.acquire()
                        asyncio.create_task(
                            self._handle_message(msg, tp.topic)
                        )
            except asyncio.CancelledError:
                break
            except Exception:
                logger.exception(
                    "kafka_consume_error", group_id=self._group_id
                )
                await asyncio.sleep(1.0)

    async def _handle_message(self, msg: Any, topic: str) -> None:
        """
        Process a single Kafka message with idempotency, retry, and
        DLQ routing.
        """
        try:
            event = self._build_event(msg)
            if event is None:
                logger.warning(
                    "event_deserialization_failed",
                    topic=topic,
                    offset=msg.offset,
                    partition=msg.partition,
                )
                return

            # Idempotency check
            if await self._is_duplicate(event.event_id):
                logger.debug(
                    "duplicate_event_skipped",
                    event_id=event.event_id,
                    event_type=event.event_type,
                )
                await self._commit_offset(msg)
                return

            handler = self._handler_map.get(event.event_type)
            if handler is None:
                logger.debug(
                    "no_handler_for_event",
                    event_type=event.event_type,
                    topic=topic,
                )
                await self._commit_offset(msg)
                return

            # Retry loop
            last_exc: Exception | None = None
            for attempt in range(1, self._max_retries + 1):
                try:
                    await handler(event)
                    await self._mark_processed(event.event_id)
                    await self._commit_offset(msg)
                    logger.info(
                        "event_handled",
                        event_type=event.event_type,
                        event_id=event.event_id,
                        topic=topic,
                    )
                    return
                except Exception as exc:
                    last_exc = exc
                    logger.warning(
                        "event_handler_retry",
                        event_type=event.event_type,
                        event_id=event.event_id,
                        attempt=attempt,
                        max_retries=self._max_retries,
                        error=str(exc),
                    )
                    await asyncio.sleep(
                        min(0.5 * (2 ** (attempt - 1)), 10.0)
                    )

            # All retries exhausted — route to DLQ
            await self._route_to_dlq(topic, event, last_exc)
            await self._commit_offset(msg)

        except Exception:
            logger.exception(
                "event_processing_error",
                topic=topic,
            )
        finally:
            self._semaphore.release()

    # ── Helpers ───────────────────────────────────────────────────────

    def _build_event(self, msg: Any) -> DomainEvent | None:
        """Deserialise a Kafka message into a ``DomainEvent``."""
        try:
            data: dict[str, Any] = msg.value
            if data is None:
                return None
            return DomainEvent.model_validate(data)
        except Exception:
            logger.warning(
                "event_parse_failed",
                raw=str(msg.value)[:500],
                exc_info=True,
            )
            return None

    async def _is_duplicate(self, event_id: str) -> bool:
        """Check Redis for a previously processed event ID."""
        if self._redis is None:
            return False
        try:
            key = f"{self._idempotency_prefix}:{event_id}"
            return bool(await self._redis.exists(key))
        except Exception:
            logger.warning(
                "idempotency_check_failed",
                event_id=event_id,
                exc_info=True,
            )
            return False

    async def _mark_processed(self, event_id: str) -> None:
        """Record a processed event ID in Redis for idempotency."""
        if self._redis is None:
            return
        try:
            key = f"{self._idempotency_prefix}:{event_id}"
            await self._redis.set(key, "1", ex=self._idempotency_ttl)
        except Exception:
            logger.warning(
                "idempotency_mark_failed",
                event_id=event_id,
                exc_info=True,
            )

    async def _commit_offset(self, msg: Any) -> None:
        """Manually commit the consumer offset for the processed message."""
        if self._consumer is None:
            return
        try:
            await self._consumer.commit()
        except Exception:
            logger.warning(
                "offset_commit_failed",
                offset=msg.offset,
                exc_info=True,
            )

    async def _route_to_dlq(
        self,
        source_topic: str,
        event: DomainEvent,
        error: Exception | None,
    ) -> None:
        """
        Publish a failed event to the dead-letter queue topic.

        DLQ topic name is ``<source_topic>.dlq``.
        """
        dlq_topic = f"{source_topic}{self._dlq_suffix}"

        # Enrich with failure metadata
        dlq_event = event.model_copy(
            update={
                "data": {
                    **event.data,
                    "_dlq_source_topic": source_topic,
                    "_dlq_error": str(error) if error else None,
                    "_dlq_timestamp": datetime.now(timezone.utc).isoformat(),
                },
            }
        )

        try:
            if self._dlq_producer is not None:
                await self._dlq_producer.publish(
                    topic=dlq_topic,
                    key=event.aggregate_id,
                    event=dlq_event,
                    headers={"dlq_source": source_topic},
                )
                logger.error(
                    "event_routed_to_dlq",
                    event_id=event.event_id,
                    event_type=event.event_type,
                    source_topic=source_topic,
                    dlq_topic=dlq_topic,
                    error=str(error),
                )
        except Exception:
            logger.exception(
                "dlq_routing_failed",
                event_id=event.event_id,
                dlq_topic=dlq_topic,
            )

    @staticmethod
    def _deserialize_value(raw: bytes) -> dict[str, Any] | None:
        """Deserialise Kafka message value from JSON bytes."""
        try:
            return orjson.loads(raw)
        except Exception:
            return None

    @staticmethod
    def _deserialize_key(raw: bytes | None) -> str:
        """Decode Kafka message key from bytes."""
        if raw is None:
            return ""
        return raw.decode("utf-8", errors="replace")


# ── In-Process Event Bus ──────────────────────────────────────────────


class EventBus:
    """
    Lightweight in-process event bus for local handler dispatch.

    This bus is intended for **within-service** communication where Kafka
    is unnecessary.  It uses asyncio tasks to invoke handlers concurrently
    and reports any handler failures via structured logging.

    Usage::

        bus = EventBus()
        bus.subscribe("CitizenRegistered", my_handler)
        await bus.publish(event)
    """

    def __init__(self) -> None:
        self._handlers: dict[str, list[EventHandler]] = defaultdict(list)

    def subscribe(self, event_type: str, handler: EventHandler) -> None:
        """
        Register a handler for an event type.

        Multiple handlers can subscribe to the same event type.

        Args:
            event_type: The ``event_type`` string to listen for.
            handler:    Async callable accepting a ``DomainEvent``.
        """
        self._handlers[event_type].append(handler)
        logger.debug(
            "event_bus_subscribed",
            event_type=event_type,
            handler=getattr(handler, "__name__", repr(handler)),
        )

    def unsubscribe(self, event_type: str, handler: EventHandler) -> None:
        """
        Remove a previously registered handler.

        Args:
            event_type: The event type to unsubscribe from.
            handler:    The handler to remove.
        """
        handlers = self._handlers.get(event_type, [])
        try:
            handlers.remove(handler)
            logger.debug(
                "event_bus_unsubscribed",
                event_type=event_type,
                handler=getattr(handler, "__name__", repr(handler)),
            )
        except ValueError:
            logger.warning(
                "event_bus_unsubscribe_not_found",
                event_type=event_type,
            )

    async def publish(self, event: DomainEvent) -> None:
        """
        Dispatch a domain event to all subscribed handlers.

        Handlers are invoked concurrently via ``asyncio.gather``.
        Individual handler failures are logged but do **not** prevent
        other handlers from executing.

        Args:
            event: The domain event to publish.
        """
        handlers = self._handlers.get(event.event_type, [])
        if not handlers:
            logger.debug(
                "event_bus_no_subscribers",
                event_type=event.event_type,
                event_id=event.event_id,
            )
            return

        logger.debug(
            "event_bus_publishing",
            event_type=event.event_type,
            event_id=event.event_id,
            handler_count=len(handlers),
        )

        tasks = [
            self._safe_invoke(handler, event) for handler in handlers
        ]
        await asyncio.gather(*tasks)

    async def publish_many(self, events: list[DomainEvent]) -> None:
        """
        Publish multiple events sequentially.

        Events are published in order to preserve causal consistency.

        Args:
            events: Ordered list of domain events.
        """
        for event in events:
            await self.publish(event)

    def clear(self) -> None:
        """Remove all subscriptions."""
        self._handlers.clear()
        logger.debug("event_bus_cleared")

    # ── Internals ─────────────────────────────────────────────────────

    @staticmethod
    async def _safe_invoke(
        handler: EventHandler, event: DomainEvent
    ) -> None:
        """
        Invoke a handler and catch exceptions so that sibling handlers
        are not affected.
        """
        try:
            await handler(event)
        except Exception:
            logger.exception(
                "event_bus_handler_error",
                event_type=event.event_type,
                event_id=event.event_id,
                handler=getattr(handler, "__name__", repr(handler)),
            )


# ── Exceptions ────────────────────────────────────────────────────────


class KafkaPublishError(Exception):
    """Raised when a Kafka publish operation fails after exhausting retries."""


# ── Public API ────────────────────────────────────────────────────────

__all__ = [
    "EventBus",
    "EventHandler",
    "KafkaConsumer",
    "KafkaProducer",
    "KafkaPublishError",
]
