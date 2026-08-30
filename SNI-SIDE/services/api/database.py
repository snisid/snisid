"""SNI-SIDE API Server — Connexions aux bases de données"""

import asyncio
from typing import Optional
from contextlib import asynccontextmanager

import asyncpg
import redis.asyncio as redis
from neo4j import AsyncGraphDatabase, AsyncDriver
from pymilvus import MilvusClient
from clickhouse_driver import Client as ClickHouseClient
from aiokafka import AIOKafkaProducer, AIOKafkaConsumer
from minio import Minio

from config import settings


class DatabasePool:
    """Pool de connexions pour toutes les bases de données SNI-SIDE"""

    def __init__(self):
        self.pg_pool: Optional[asyncpg.Pool] = None
        self.cockroach_pool: Optional[asyncpg.Pool] = None
        self.neo4j_driver: Optional[AsyncDriver] = None
        self.milvus_client: Optional[MilvusClient] = None
        self.clickhouse_client: Optional[ClickHouseClient] = None
        self.kafka_producer: Optional[AIOKafkaProducer] = None
        self.redis_client: Optional[redis.Redis] = None
        self.minio_client: Optional[Minio] = None
        self._initialized = False

    async def initialize(self):
        """Initialise toutes les connexions aux bases de données"""
        if self._initialized:
            return

        # PostgreSQL
        self.pg_pool = await asyncpg.create_pool(
            host=settings.postgres_host,
            port=settings.postgres_port,
            database=settings.postgres_db,
            user=settings.postgres_user,
            password=settings.postgres_password,
            min_size=settings.postgres_min_connections,
            max_size=settings.postgres_max_connections,
            command_timeout=30,
        )

        # CockroachDB
        self.cockroach_pool = await asyncpg.create_pool(
            host=settings.cockroach_host,
            port=settings.cockroach_port,
            database=settings.cockroach_db,
            user=settings.cockroach_user,
            password=settings.cockroach_password,
            min_size=5,
            max_size=settings.cockroach_max_connections,
            command_timeout=30,
        )

        # Neo4j
        self.neo4j_driver = AsyncGraphDatabase.driver(
            settings.neo4j_uri,
            auth=(settings.neo4j_user, settings.neo4j_password),
            max_connection_pool_size=settings.neo4j_max_connection_pool_size,
        )

        # Milvus
        self.milvus_client = MilvusClient(
            host=settings.milvus_host,
            port=settings.milvus_port,
        )

        # ClickHouse
        self.clickhouse_client = ClickHouseClient(
            host=settings.clickhouse_host,
            port=settings.clickhouse_port,
            user=settings.clickhouse_user,
            password=settings.clickhouse_password,
            database=settings.clickhouse_db,
        )

        # Kafka Producer
        self.kafka_producer = AIOKafkaProducer(
            bootstrap_servers=settings.kafka_bootstrap_servers,
            client_id=settings.kafka_client_id,
            acks='all',
            compression_type='zstd',
            max_batch_size=1048576,
        )
        await self.kafka_producer.start()

        # Redis
        self.redis_client = redis.Redis(
            host=settings.redis_host,
            port=settings.redis_port,
            db=settings.redis_db,
            password=settings.redis_password or None,
            decode_responses=True,
        )

        # MinIO
        self.minio_client = Minio(
            settings.minio_endpoint,
            access_key=settings.minio_access_key,
            secret_key=settings.minio_secret_key,
            secure=settings.minio_secure,
        )
        self._ensure_buckets()

        self._initialized = True

    def _ensure_buckets(self):
        """Crée les buckets MinIO si nécessaire"""
        for bucket in [settings.minio_evidence_bucket, settings.minio_document_bucket]:
            if not self.minio_client.bucket_exists(bucket):
                self.minio_client.make_bucket(bucket)

    async def close(self):
        """Ferme toutes les connexions"""
        if self.pg_pool:
            await self.pg_pool.close()
        if self.cockroach_pool:
            await self.cockroach_pool.close()
        if self.neo4j_driver:
            await self.neo4j_driver.close()
        if self.kafka_producer:
            await self.kafka_producer.stop()
        if self.redis_client:
            await self.redis_client.close()
        self._initialized = False

    # ============ HELPERS ============
    @asynccontextmanager
    async def pg_conn(self):
        """Connexion PostgreSQL (context manager)"""
        async with self.pg_pool.acquire() as conn:
            yield conn

    @asynccontextmanager
    async def cockroach_conn(self):
        """Connexion CockroachDB (context manager)"""
        async with self.cockroach_pool.acquire() as conn:
            yield conn

    async def neo4j_session(self, database: str = "neo4j"):
        """Session Neo4j"""
        return self.neo4j_driver.session(database=database)

    async def publish_event(self, topic: str, key: str, value: bytes):
        """Publie un événement Kafka"""
        await self.kafka_producer.send(
            topic=topic,
            key=key.encode(),
            value=value,
        )

    async def cache_get(self, key: str) -> Optional[str]:
        """Récupère une valeur du cache Redis"""
        return await self.redis_client.get(key)

    async def cache_set(self, key: str, value: str, ttl: int = 300):
        """Stocke une valeur dans Redis avec TTL"""
        await self.redis_client.setex(key, ttl, value)

    async def cache_delete(self, key: str):
        """Supprime une clé Redis"""
        await self.redis_client.delete(key)


# Instance singleton
db = DatabasePool()
