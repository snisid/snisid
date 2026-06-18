from __future__ import annotations

import uuid
from datetime import datetime, timezone
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
import pytest_asyncio

from shared.cqrs import DomainEvent
from shared.events import EventBus, KafkaProducer, KafkaConsumer, KafkaPublishError


class TestEventBus:
    """Test in-process event bus."""

    @pytest.mark.asyncio
    async def test_publish_with_subscriber(self, sample_event):
        bus = EventBus()
        handler = AsyncMock()
        bus.subscribe("test.event", handler)
        await bus.publish(sample_event)
        handler.assert_awaited_once_with(sample_event)

    @pytest.mark.asyncio
    async def test_publish_no_subscribers(self, sample_event):
        bus = EventBus()
        await bus.publish(sample_event)
        assert True

    @pytest.mark.asyncio
    async def test_multiple_subscribers(self, sample_event):
        bus = EventBus()
        handler1 = AsyncMock()
        handler2 = AsyncMock()
        bus.subscribe("test.event", handler1)
        bus.subscribe("test.event", handler2)
        await bus.publish(sample_event)
        handler1.assert_awaited_once_with(sample_event)
        handler2.assert_awaited_once_with(sample_event)

    @pytest.mark.asyncio
    async def test_subscriber_error_does_not_affect_others(self, sample_event):
        bus = EventBus()
        failing_handler = AsyncMock(side_effect=ValueError("handler failed"))
        working_handler = AsyncMock()

        bus.subscribe("test.event", failing_handler)
        bus.subscribe("test.event", working_handler)

        await bus.publish(sample_event)
        working_handler.assert_awaited_once_with(sample_event)

    @pytest.mark.asyncio
    async def test_unsubscribe(self, sample_event):
        bus = EventBus()
        handler = AsyncMock()
        bus.subscribe("test.event", handler)
        bus.unsubscribe("test.event", handler)
        await bus.publish(sample_event)
        handler.assert_not_awaited()

    @pytest.mark.asyncio
    async def test_unsubscribe_nonexistent_handler(self):
        bus = EventBus()
        bus.unsubscribe("nonexistent", lambda e: None)
        assert True

    @pytest.mark.asyncio
    async def test_clear(self, sample_event):
        bus = EventBus()
        handler = AsyncMock()
        bus.subscribe("test.event", handler)
        bus.clear()
        await bus.publish(sample_event)
        handler.assert_not_awaited()

    @pytest.mark.asyncio
    async def test_publish_many(self, sample_events):
        bus = EventBus()
        handler = AsyncMock()
        bus.subscribe("test.event", handler)
        bus.subscribe("test.event2", handler)
        await bus.publish_many(sample_events)
        assert handler.await_count == 2


class TestEventSerialization:
    """Test event serialization and deserialization."""

    @pytest.mark.asyncio
    async def test_event_to_dict(self, sample_event):
        data = sample_event.model_dump(mode="json")
        assert data["event_type"] == "test.event"
        assert data["aggregate_id"] == sample_event.aggregate_id
        assert data["data"]["key"] == "value"

    @pytest.mark.asyncio
    async def test_event_from_dict(self, sample_event):
        data = sample_event.model_dump(mode="json")
        restored = DomainEvent.model_validate(data)
        assert restored.event_id == sample_event.event_id
        assert restored.event_type == sample_event.event_type

    @pytest.mark.asyncio
    async def test_event_with_complex_data(self):
        event = DomainEvent(
            event_type="complex.event",
            aggregate_id=str(uuid.uuid4()),
            aggregate_type="Test",
            data={"list": [1, 2, 3], "nested": {"a": 1}, "flag": True},
        )
        data = event.model_dump(mode="json")
        restored = DomainEvent.model_validate(data)
        assert restored.data["list"] == [1, 2, 3]
        assert restored.data["nested"]["a"] == 1


class TestKafkaProducer:
    """Test Kafka producer (mocked)."""

    @pytest.mark.asyncio
    async def test_publish(self):
        with patch("aiokafka.AIOKafkaProducer") as mock_producer_cls:
            mock_instance = AsyncMock()
            mock_producer_cls.return_value = mock_instance

            producer = KafkaProducer(
                bootstrap_servers="localhost:9092",
                topic_prefix="snisid",
            )
            await producer.start()
            assert producer._producer is not None

            event = DomainEvent(
                event_type="test.event",
                aggregate_id="agg-1",
                aggregate_type="Test",
            )
            await producer.publish(topic="test", key="agg-1", event=event)
            mock_instance.send_and_wait.assert_awaited_once()

            await producer.stop()
            mock_instance.stop.assert_awaited_once()

    @pytest.mark.asyncio
    async def test_publish_without_start(self):
        producer = KafkaProducer(
            bootstrap_servers="localhost:9092",
        )
        event = DomainEvent(
            event_type="test.event",
            aggregate_id="agg-1",
            aggregate_type="Test",
        )
        with pytest.raises(RuntimeError, match="not started"):
            await producer.publish(topic="test", key="key", event=event)

    @pytest.mark.asyncio
    async def test_double_start(self):
        with patch("aiokafka.AIOKafkaProducer") as mock_cls:
            mock_instance = AsyncMock()
            mock_cls.return_value = mock_instance

            producer = KafkaProducer(bootstrap_servers="localhost:9092")
            await producer.start()
            with pytest.raises(RuntimeError, match="already started"):
                await producer.start()

    @pytest.mark.asyncio
    async def test_publish_headers(self):
        with patch("aiokafka.AIOKafkaProducer") as mock_cls:
            mock_instance = AsyncMock()
            mock_cls.return_value = mock_instance

            producer = KafkaProducer(bootstrap_servers="localhost:9092")
            await producer.start()

            event = DomainEvent(
                event_type="test.event",
                aggregate_id="agg-1",
                aggregate_type="Test",
                correlation_id="corr-1",
            )
            await producer.publish(
                topic="test",
                key="agg-1",
                event=event,
                headers={"source": "test-suite"},
            )

            call_args = mock_instance.send_and_wait.call_args
            headers = call_args[1]["headers"]
            header_dict = {k: v.decode() for k, v in headers}
            assert header_dict["event_type"] == "test.event"
            assert header_dict["correlation_id"] == "corr-1"
            assert header_dict["source"] == "test-suite"

            await producer.stop()

    @pytest.mark.asyncio
    async def test_publish_batch(self):
        with patch("aiokafka.AIOKafkaProducer") as mock_cls:
            mock_instance = AsyncMock()
            batch = MagicMock()
            batch.append.return_value = MagicMock()
            mock_instance.create_batch = MagicMock(return_value=batch)
            mock_cls.return_value = mock_instance

            producer = KafkaProducer(
                bootstrap_servers="localhost:9092",
                topic_prefix="snisid",
            )
            await producer.start()

            events = [
                ("key1", DomainEvent(event_type="e1", aggregate_id="a1", aggregate_type="T")),
                ("key2", DomainEvent(event_type="e2", aggregate_id="a2", aggregate_type="T")),
            ]

            count = await producer.publish_batch(topic="batch", events=events)
            assert count == 2


class TestKafkaConsumer:
    """Test Kafka consumer (mocked)."""

    @pytest.mark.asyncio
    async def test_start_and_stop(self):
        with patch("aiokafka.AIOKafkaConsumer") as mock_consumer_cls, \
             patch("shared.events.KafkaProducer") as mock_producer_cls:

            mock_consumer = AsyncMock()
            mock_consumer_cls.return_value = mock_consumer
            mock_producer = AsyncMock()
            mock_producer_cls.return_value = mock_producer

            consumer = KafkaConsumer(
                bootstrap_servers="localhost:9092",
                group_id="test-group",
                topics=["snisid.test"],
            )
            await consumer.start()
            assert consumer._running is True
            mock_consumer.start.assert_awaited_once()

            await consumer.stop()
            assert consumer._running is False
            mock_consumer.stop.assert_awaited_once()

    @pytest.mark.asyncio
    async def test_double_start(self):
        with patch("aiokafka.AIOKafkaConsumer") as mock_cls, \
             patch("shared.events.KafkaProducer") as mock_prod:
            mock_cls.return_value = AsyncMock()
            mock_prod.return_value = AsyncMock()
            consumer = KafkaConsumer(
                bootstrap_servers="localhost:9092",
                group_id="test-group",
            )
            await consumer.start()
            with pytest.raises(RuntimeError, match="already running"):
                await consumer.start()

    @pytest.mark.asyncio
    async def test_register_handler(self):
        consumer = KafkaConsumer(bootstrap_servers="localhost:9092", group_id="test")
        handler = AsyncMock()
        consumer.register_handler("test.event", handler)
        assert "test.event" in consumer._handler_map


class TestKafkaErrorHandling:
    """Test error handling in Kafka operations."""

    @pytest.mark.asyncio
    async def test_publish_retry_then_fail(self):
        with patch("aiokafka.AIOKafkaProducer") as mock_cls:
            mock_instance = AsyncMock()
            mock_instance.send_and_wait = AsyncMock(side_effect=ConnectionError("broker down"))
            mock_cls.return_value = mock_instance

            producer = KafkaProducer(
                bootstrap_servers="localhost:9092",
                topic_prefix="test",
            )
            await producer.start()

            event = DomainEvent(
                event_type="test.event",
                aggregate_id="agg-1",
                aggregate_type="Test",
            )
            with pytest.raises(KafkaPublishError):
                await producer.publish(topic="test", key="key", event=event)

    @pytest.mark.asyncio
    async def test_publish_succeeds_on_retry(self):
        with patch("aiokafka.AIOKafkaProducer") as mock_cls:
            mock_instance = AsyncMock()
            attempts = 0

            async def send_and_wait(**kwargs):
                nonlocal attempts
                attempts += 1
                if attempts < 2:
                    raise ConnectionError("transient")
                return None

            mock_instance.send_and_wait = send_and_wait
            mock_cls.return_value = mock_instance

            producer = KafkaProducer(
                bootstrap_servers="localhost:9092",
                topic_prefix="test",
            )
            await producer.start()

            event = DomainEvent(
                event_type="test.event",
                aggregate_id="agg-1",
                aggregate_type="Test",
            )
            await producer.publish(topic="test", key="key", event=event)
            assert attempts == 2


class TestEventBusTypeSafety:
    """Test event bus with different event types."""

    @pytest.mark.asyncio
    async def test_only_relevant_handlers_called(self):
        bus = EventBus()
        handler_a = AsyncMock()
        handler_b = AsyncMock()

        bus.subscribe("event.a", handler_a)
        bus.subscribe("event.b", handler_b)

        event_a = DomainEvent(
            event_type="event.a",
            aggregate_id="agg-1",
            aggregate_type="Test",
        )
        await bus.publish(event_a)

        handler_a.assert_awaited_once()
        handler_b.assert_not_awaited()
