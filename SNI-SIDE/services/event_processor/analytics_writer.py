import logging, os, json, asyncio
from datetime import datetime
from typing import Optional

logger = logging.getLogger("sniside.analytics-writer")

CLICKHOUSE_HOST = os.getenv("CLICKHOUSE_HOST", "clickhouse")
CLICKHOUSE_PORT = int(os.getenv("CLICKHOUSE_PORT", "9000"))
CLICKHOUSE_USER = os.getenv("CLICKHOUSE_USER", "default")
CLICKHOUSE_PASSWORD = os.getenv("CLICKHOUSE_PASSWORD", "")


class ClickHouseAnalyticsWriter:
    def __init__(self):
        self.client = None

    def start(self):
        try:
            from clickhouse_driver import Client
            self.client = Client(
                host=CLICKHOUSE_HOST, port=CLICKHOUSE_PORT,
                user=CLICKHOUSE_USER, password=CLICKHOUSE_PASSWORD,
            )
            self.client.execute("SELECT 1")
            logger.info("ClickHouse connected")
        except Exception as e:
            logger.warning(f"ClickHouse unavailable — analytics writes skipped: {e}")

    def stop(self):
        if self.client:
            self.client.disconnect()

    async def write_event(self, domain: str, event) -> bool:
        if not self.client:
            return False
        try:
            v = event.value
            ts = datetime.fromtimestamp(event.timestamp / 1000) if event.timestamp else datetime.utcnow()
            table = f"sniside_analytics.{domain}_events"
            insert = f"INSERT INTO {table} (event_id, topic, key, value, timestamp) VALUES"
            values = [(str(event.event_id), event.topic, event.key or "",
                       json.dumps(v), ts)]
            await asyncio.to_thread(self.client.execute, insert, values)
            return True
        except Exception as e:
            logger.debug(f"Analytics write failed for {domain}: {e}")
            return False
