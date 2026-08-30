"""SNI-SIDE API Server — Configuration centralisée"""

from pydantic_settings import BaseSettings
from typing import List, Optional
from enum import Enum


class Environment(str, Enum):
    DEVELOPMENT = "development"
    STAGING = "staging"
    PRODUCTION = "production"


class Settings(BaseSettings):
    # App
    app_name: str = "SNI-SIDE API"
    environment: Environment = Environment.DEVELOPMENT
    debug: bool = False
    @property
    def cors_origins(self) -> List[str]:
        if self.environment == Environment.PRODUCTION:
            return ["https://sniside.intranet.pnh.ht"]
        return ["*"]

    # Server
    host: str = "0.0.0.0"
    port: int = 8080
    workers: int = 4
    reload: bool = False

    # Security
    jwt_secret_key: str = "change-me-in-production"
    jwt_algorithm: str = "HS256"
    jwt_expiration_minutes: int = 60
    mtls_enabled: bool = True
    spire_socket_path: str = "unix:///tmp/spire-agent/public/api.sock"

    # PostgreSQL — NCID, CODIS, Missing, Vehicle, Firearms, Border, Narcotics, Financial, Cyber, Watchlist, Docs, GEOINT
    postgres_host: str = "sniside-ncid-postgres"
    postgres_port: int = 5432
    postgres_db: str = "sniside"
    postgres_user: str = "sniside"
    postgres_password: str = ""
    postgres_max_connections: int = 50
    postgres_min_connections: int = 10

    # CockroachDB — ALPR, Evidence
    cockroach_host: str = "sniside-cockroachdb"
    cockroach_port: int = 26257
    cockroach_db: str = "sniside"
    cockroach_user: str = "sniside"
    cockroach_password: str = ""
    cockroach_max_connections: int = 30

    # Neo4j — Graph Intelligence
    neo4j_uri: str = "bolt://sniside-neo4j:7687"
    neo4j_user: str = "neo4j"
    neo4j_password: str = ""
    neo4j_max_connection_pool_size: int = 50

    # Milvus — Vector Database
    milvus_host: str = "sniside-milvus-proxy"
    milvus_port: int = 19530
    milvus_aliases_collection: str = "sniside_face_embeddings"
    milvus_fingerprint_collection: str = "sniside_fingerprint_templates"

    # ClickHouse — Analytics
    clickhouse_host: str = "sniside-clickhouse"
    clickhouse_port: int = 9000
    clickhouse_db: str = "sniside"
    clickhouse_user: str = "default"
    clickhouse_password: str = ""

    # Kafka
    kafka_bootstrap_servers: str = "sniside-kafka-kafka-bootstrap:9092"
    kafka_client_id: str = "sniside-api"
    kafka_schema_registry_url: str = "http://sniside-schema-registry:8081"

    # Redis — Cache
    redis_host: str = "sniside-redis"
    redis_port: int = 6379
    redis_db: int = 0
    redis_password: str = ""

    # MinIO — Object Store
    minio_endpoint: str = "sniside-minio:9000"
    minio_access_key: str = ""
    minio_secret_key: str = ""
    minio_secure: bool = False
    minio_evidence_bucket: str = "sniside-evidence"
    minio_document_bucket: str = "sniside-documents"

    # OpenTelemetry
    otel_service_name: str = "sniside-api"
    otel_exporter_otlp_endpoint: str = "http://sniside-otel-collector:4318"

    # AI Models
    arcface_model_path: str = "/models/arcface_v1.pth"
    fraud_gnn_model_path: str = "/models/fraud_gnn_v1.pth"
    deepfake_model_path: str = "/models/deepfake_v1.pth"

    # Search Engine
    search_default_limit: int = 20
    search_max_limit: int = 100
    search_timeout_ms: int = 5000

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
        case_sensitive = False


settings = Settings()
