from __future__ import annotations

from shared.events import KafkaProducer

_producer: KafkaProducer | None = None


async def init_bio_adn_kafka() -> KafkaProducer | None:
    global _producer
    try:
        producer = KafkaProducer()
        await producer.start()
        _producer = producer
        import services.bio_adn.api as api
        api._producer = producer
        return producer
    except Exception:
        return None


async def shutdown_bio_adn_kafka() -> None:
    global _producer
    if _producer is not None:
        try:
            await _producer.stop()
        except Exception:
            pass
        _producer = None
