import logging
import os
from datetime import datetime

logger = logging.getLogger("sniside.etl.writers")

PG_CONNECTION = os.getenv(
    "SNISIDE_PG_CONNECTION",
    "postgresql://sniside:sniside-secret@postgres:5432/sniside",
)


class PostgresWriter:
    def __init__(self, connection_string: str = None):
        self.conn_string = connection_string or self._build_connection_string()

    def _build_connection_string(self) -> str:
        return os.getenv(
            "SNISIDE_PG_CONNECTION",
            PG_CONNECTION,
        )

    def write_batch(self, records: list, schema: str, table: str):
        import psycopg2
        conn = psycopg2.connect(self.conn_string)
        try:
            with conn.cursor() as cur:
                if records:
                    columns = list(records[0].keys())
                    placeholders = ", ".join(["%s"] * len(columns))
                    col_names = ", ".join(columns)
                    rows = [[r.get(c) for c in columns] for r in records]
                    cur.executemany(
                        f"INSERT INTO {schema}.{table} ({col_names}) VALUES ({placeholders}) ON CONFLICT DO NOTHING",
                        rows,
                    )
                conn.commit()
            logger.debug(f"PostgresWriter: {len(records)} rows \u2192 {schema}.{table}")
        except Exception as e:
            conn.rollback()
            raise
        finally:
            conn.close()


class Neo4jWriter:
    def __init__(self):
        self.uri = os.getenv("SNISIDE_NEO4J_URI", "bolt://neo4j:7687")
        self.user = os.getenv("SNISIDE_NEO4J_USER", "neo4j")
        self.password = os.getenv("SNISIDE_NEO4J_PASSWORD", "sniside-neo4j")

    def update_graph(self, label: str, stats: dict):
        from neo4j import GraphDatabase
        driver = GraphDatabase.driver(self.uri, auth=(self.user, self.password))
        with driver.session() as session:
            session.run(f"""
                MERGE (n:ETLMetadata {{label: $label}})
                SET n.rows_inserted = $rows, n.last_run = timestamp()
            """, label=label or "UNKNOWN", rows=stats.get("rows_inserted", 0))
        driver.close()
        logger.info(f"Neo4jWriter: ETL metadata updated for {label}")


class KafkaWriter:
    def __init__(self, bootstrap_servers: str = None):
        self.servers = bootstrap_servers or os.getenv("SNISIDE_KAFKA_BROKERS", "kafka:9092")

    def emit_topic(self, topic: str, stats: dict):
        try:
            import json, uuid
            from kafka import KafkaProducer
            producer = KafkaProducer(bootstrap_servers=self.servers, compression_type="zstd")
            msg = {
                "event_id": str(uuid.uuid4()),
                "etl_run": datetime.utcnow().isoformat(),
                "rows_inserted": stats.get("rows_inserted", 0),
                "rows_errors": stats.get("rows_errors", 0),
                "source": "sniside-etl-toolkit",
            }
            producer.send(topic, key=b"etl", value=json.dumps(msg).encode("utf-8"))
            producer.flush()
            producer.close()
            logger.info(f"KafkaWriter: event emitted to {topic}")
        except Exception as e:
            logger.warning(f"KafkaWriter failed: {e}")


class CockroachWriter(PostgresWriter):
    def __init__(self, connection_string: str = None):
        super().__init__(
            connection_string
            or os.getenv(
                "SNISIDE_COCKROACH_CONNECTION",
                "postgresql://sniside:sniside-secret@cockroach:26257/sniside?sslmode=disable",
            )
        )
