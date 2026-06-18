from __future__ import annotations

import json
from enum import Enum
from pathlib import Path
from typing import Any, Dict, List, Optional

import yaml
from pydantic import BaseModel, Field, field_validator


class SourceType(str, Enum):
    csv = "csv"
    sqlserver = "sqlserver"
    postgresql = "postgresql"
    rest_api = "rest_api"


class CleansingRules(BaseModel):
    normalize_names: bool = True
    normalize_dates: bool = True
    normalize_phones: bool = True
    normalize_addresses: bool = True
    strip_accents: bool = True
    strip_special_chars: bool = True
    normalize_case: bool = True
    expand_abbreviations: bool = True
    min_name_length: int = 1
    max_valid_year: int = 2026
    min_valid_year: int = 1900


class MatchingConfig(BaseModel):
    jaro_winkler_threshold: float = 0.88
    auto_merge_threshold: float = 0.95
    pending_review_threshold: float = 0.80
    phonetic_enabled: bool = True
    exact_cin_match_boost: float = 0.15
    blocking_dob_year_window: int = 1
    blocking_name_prefix_len: int = 2
    max_candidates_per_block: int = 50


class RetryConfig(BaseModel):
    max_retries: int = 3
    backoff_base_seconds: float = 1.0
    backoff_max_seconds: float = 60.0
    backoff_multiplier: float = 2.0
    retryable_statuses: List[int] = Field(default_factory=lambda: [429, 500, 502, 503, 504])


class RateLimitConfig(BaseModel):
    max_requests_per_second: float = 50.0
    max_concurrent: int = 5
    burst_size: int = 10


class SourceSchema(BaseModel):
    fields: Dict[str, str] = Field(default_factory=dict)


class SourceConfig(BaseModel):
    type: SourceType
    path: Optional[str] = None
    connection_string: Optional[str] = None
    host: Optional[str] = None
    port: Optional[int] = None
    database: Optional[str] = None
    username: Optional[str] = None
    password: Optional[str] = None
    query: Optional[str] = None
    table: Optional[str] = None
    api_url: Optional[str] = None
    api_key: Optional[str] = None
    encoding: str = "utf-8-sig"
    delimiter: str = ";"
    schema: Optional[SourceSchema] = None
    ssl_mode: str = "prefer"
    connect_timeout: int = 30
    fetch_size: int = 10000

    @field_validator("port")
    @classmethod
    def validate_port(cls, v):
        if v is not None and (v < 1 or v > 65535):
            raise ValueError(f"Port must be between 1 and 65535, got {v}")
        return v


class TargetConfig(BaseModel):
    api_url: str
    api_key: Optional[str] = None
    auth_token: Optional[str] = None
    batch_size: int = 100
    timeout_seconds: int = 60
    retry: RetryConfig = Field(default_factory=RetryConfig)
    rate_limit: RateLimitConfig = Field(default_factory=RateLimitConfig)
    idempotency_key_header: str = "X-Idempotency-Key"
    endpoint_identities: str = "/identities"
    endpoint_batch: str = "/identities/batch"
    endpoint_status: str = "/migration/status"


class LoggingConfig(BaseModel):
    level: str = "INFO"
    file: Optional[str] = None
    format: str = "%(asctime)s [%(levelname)s] %(name)s: %(message)s"
    max_bytes: int = 10485760
    backup_count: int = 5


class PipelineConfig(BaseModel):
    name: str = "snisid_migration"
    batch_size: int = 1000
    max_retries: int = 3
    checkpoint_dir: str = "./checkpoints"
    quarantine_dir: str = "./quarantine"
    report_dir: str = "./reports"
    source: SourceConfig
    target: TargetConfig
    cleansing: CleansingRules = Field(default_factory=CleansingRules)
    matching: MatchingConfig = Field(default_factory=MatchingConfig)
    logging: LoggingConfig = Field(default_factory=LoggingConfig)
    dry_run: bool = False
    strict_mode: bool = False
    max_errors_per_batch: int = 50

    @field_validator("batch_size")
    @classmethod
    def batch_size_positive(cls, v):
        if v < 1:
            raise ValueError(f"batch_size must be >= 1, got {v}")
        return v

    @classmethod
    def from_yaml(cls, path: str) -> PipelineConfig:
        path_obj = Path(path)
        if not path_obj.exists():
            raise FileNotFoundError(f"Config file not found: {path}")
        with open(path_obj, "r", encoding="utf-8") as f:
            data = yaml.safe_load(f)
        return cls(**data)

    @classmethod
    def from_json(cls, path: str) -> PipelineConfig:
        path_obj = Path(path)
        if not path_obj.exists():
            raise FileNotFoundError(f"Config file not found: {path}")
        with open(path_obj, "r", encoding="utf-8") as f:
            data = json.load(f)
        return cls(**data)

    def to_dict(self) -> Dict[str, Any]:
        return json.loads(self.model_dump_json())

    def to_yaml(self, path: str) -> None:
        with open(path, "w", encoding="utf-8") as f:
            yaml.dump(json.loads(self.model_dump_json()), f, default_flow_style=False, allow_unicode=True)
