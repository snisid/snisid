"""
SNISID Health Check Framework
===============================
Kubernetes-ready health check system with liveness, readiness, and
Prometheus metrics endpoints.

Endpoints (no authentication required):
- ``GET /health``  — Liveness probe: always returns 200
- ``GET /ready``   — Readiness probe: runs all registered component checks
- ``GET /metrics`` — Prometheus exposition format metrics

Built-in checks:
- PostgreSQL (``SELECT 1``)
- Redis (``PING``)
- Kafka (metadata request)
"""
from __future__ import annotations

import time
from datetime import datetime, timezone
from enum import Enum
from typing import Any, Awaitable, Callable

import structlog
from fastapi import APIRouter, Response
from pydantic import BaseModel, Field

from shared.config import get_settings
from shared.logging import get_logger

logger = get_logger(__name__)


# ── Models ────────────────────────────────────────────────────────────


class HealthStatus(str, Enum):
    """Possible states for a health-checked component."""

    UP = "UP"
    DOWN = "DOWN"
    DEGRADED = "DEGRADED"


class ComponentHealth(BaseModel):
    """Health status of a single system component."""

    name: str = Field(..., description="Component identifier (e.g. 'database')")
    status: HealthStatus = Field(..., description="Current health status")
    details: dict[str, Any] = Field(
        default_factory=dict,
        description="Free-form details about the component state",
    )
    latency_ms: float = Field(
        0.0,
        description="Time taken to execute the health check (ms)",
    )


class ReadinessResponse(BaseModel):
    """Aggregate readiness response returned by ``GET /ready``."""

    status: HealthStatus = Field(..., description="Overall system health")
    timestamp: str = Field(..., description="ISO-8601 check timestamp (UTC)")
    service: str = Field(..., description="Service name")
    version: str = Field(..., description="Service version")
    checks: list[ComponentHealth] = Field(
        default_factory=list,
        description="Per-component health results",
    )
    uptime_seconds: float = Field(0.0, description="Process uptime in seconds")


class LivenessResponse(BaseModel):
    """Minimal liveness probe response for ``GET /health``."""

    status: str = Field("alive", description="Always 'alive' if process is running")
    timestamp: str = Field(..., description="ISO-8601 timestamp (UTC)")
    service: str = Field(..., description="Service name")


# ── Health Check Registry ─────────────────────────────────────────────

# Type alias for check functions: async () -> ComponentHealth
CheckFn = Callable[[], Awaitable[ComponentHealth]]


class HealthCheck:
    """
    Registry of named health-check functions.

    Register asynchronous check callables that probe individual
    infrastructure components.  The readiness endpoint executes all
    registered checks concurrently and aggregates the results.

    Example::

        health = HealthCheck()
        health.register("database", check_database)
        health.register("redis", check_redis)

        router = create_health_router(health)
        app.include_router(router)
    """

    def __init__(self) -> None:
        self._checks: dict[str, CheckFn] = {}
        self._start_time: float = time.monotonic()

    def register(self, name: str, check_fn: CheckFn) -> None:
        """
        Register a health-check function under the given *name*.

        Args:
            name: Unique identifier for the component (e.g. ``"database"``).
            check_fn: An async callable that returns a ``ComponentHealth``.

        Raises:
            ValueError: If a check with the same name is already registered.
        """
        if name in self._checks:
            raise ValueError(f"Health check '{name}' is already registered")
        self._checks[name] = check_fn
        logger.info("health_check_registered", component=name)

    def unregister(self, name: str) -> None:
        """
        Remove a previously registered health check.

        Args:
            name: The component name to remove.
        """
        self._checks.pop(name, None)
        logger.info("health_check_unregistered", component=name)

    async def run_all(self) -> list[ComponentHealth]:
        """
        Execute every registered check and return the results.

        Each check is executed sequentially to avoid thundering-herd
        issues against infrastructure components.  Individual failures
        are caught so one failing check does not prevent the others
        from reporting.

        Returns:
            A list of ``ComponentHealth`` results, one per registered check.
        """
        results: list[ComponentHealth] = []
        for name, check_fn in self._checks.items():
            start = time.perf_counter()
            try:
                result = await check_fn()
                result.latency_ms = round((time.perf_counter() - start) * 1000, 2)
                results.append(result)
            except Exception as exc:
                elapsed = round((time.perf_counter() - start) * 1000, 2)
                logger.error(
                    "health_check_exception",
                    component=name,
                    error=str(exc),
                    latency_ms=elapsed,
                )
                results.append(
                    ComponentHealth(
                        name=name,
                        status=HealthStatus.DOWN,
                        details={"error": str(exc)},
                        latency_ms=elapsed,
                    )
                )
        return results

    @property
    def uptime_seconds(self) -> float:
        """Seconds since the HealthCheck instance was created."""
        return round(time.monotonic() - self._start_time, 2)

    @property
    def registered_checks(self) -> list[str]:
        """Names of all registered health checks."""
        return list(self._checks.keys())


# ── Built-in Checks ──────────────────────────────────────────────────


async def check_database() -> ComponentHealth:
    """
    Check PostgreSQL connectivity by executing ``SELECT 1``.

    Uses the shared async engine from ``shared.database``.
    """
    from shared.database import _engine  # noqa: WPS436 — internal access is intentional

    name = "database"
    if _engine is None:
        return ComponentHealth(
            name=name,
            status=HealthStatus.DOWN,
            details={"error": "Database engine not initialised"},
        )

    try:
        from sqlalchemy import text

        async with _engine.connect() as conn:
            result = await conn.execute(text("SELECT 1"))
            row = result.scalar()
            return ComponentHealth(
                name=name,
                status=HealthStatus.UP,
                details={"result": row, "pool_size": _engine.pool.size()},
            )
    except Exception as exc:
        logger.error("health_check_database_failed", error=str(exc))
        return ComponentHealth(
            name=name,
            status=HealthStatus.DOWN,
            details={"error": str(exc)},
        )


async def check_redis() -> ComponentHealth:
    """
    Check Redis connectivity by issuing a ``PING`` command.

    Creates a short-lived connection using settings from ``shared.config``.
    """
    name = "redis"
    try:
        import redis.asyncio as aioredis

        settings = get_settings()
        client = aioredis.from_url(
            settings.redis.url,
            socket_timeout=settings.redis.socket_timeout,
            socket_connect_timeout=settings.redis.socket_connect_timeout,
            decode_responses=True,
        )
        try:
            pong: Any = await client.ping()
            info = await client.info(section="server")
            return ComponentHealth(
                name=name,
                status=HealthStatus.UP,
                details={
                    "ping": pong,
                    "redis_version": info.get("redis_version", "unknown"),
                    "connected_clients": info.get("connected_clients", -1),
                },
            )
        finally:
            await client.aclose()
    except Exception as exc:
        logger.error("health_check_redis_failed", error=str(exc))
        return ComponentHealth(
            name=name,
            status=HealthStatus.DOWN,
            details={"error": str(exc)},
        )


async def check_kafka() -> ComponentHealth:
    """
    Check Kafka connectivity by requesting cluster metadata.

    Creates a short-lived ``AIOKafkaProducer`` and calls
    ``client.check_version()`` to verify broker reachability.
    """
    name = "kafka"
    try:
        from aiokafka import AIOKafkaProducer

        settings = get_settings()
        producer = AIOKafkaProducer(
            bootstrap_servers=settings.kafka.bootstrap_list,
            security_protocol=settings.kafka.security_protocol,
            request_timeout_ms=5000,
        )
        try:
            await producer.start()
            # Partition metadata confirms full broker connectivity
            partitions = await producer.partitions_for("__consumer_offsets")
            return ComponentHealth(
                name=name,
                status=HealthStatus.UP,
                details={
                    "bootstrap_servers": settings.kafka.bootstrap_servers,
                    "partitions_available": len(partitions) if partitions else 0,
                },
            )
        finally:
            await producer.stop()
    except Exception as exc:
        logger.error("health_check_kafka_failed", error=str(exc))
        return ComponentHealth(
            name=name,
            status=HealthStatus.DOWN,
            details={"error": str(exc)},
        )


# ── Router Factory ────────────────────────────────────────────────────


def create_health_router(health_check: HealthCheck) -> APIRouter:
    """
    Create a FastAPI router with liveness, readiness, and metrics endpoints.

    These endpoints are **exempt from authentication** because they are
    consumed by Kubernetes probes and monitoring systems that do not
    carry JWT tokens.

    Args:
        health_check: The ``HealthCheck`` registry containing registered
            component checks.

    Returns:
        A configured ``APIRouter`` ready to be included in the application.

    Example::

        from shared.health import HealthCheck, create_health_router, check_database

        health = HealthCheck()
        health.register("database", check_database)

        app = FastAPI()
        app.include_router(create_health_router(health))
    """
    router = APIRouter(tags=["Health"])
    settings = get_settings()

    # ── GET /health  — Liveness ────────────────────────────────────

    @router.get(
        "/health",
        response_model=LivenessResponse,
        summary="Liveness probe",
        description=(
            "Returns HTTP 200 as long as the process is running.  "
            "Used by Kubernetes ``livenessProbe``."
        ),
        status_code=200,
    )
    async def liveness() -> LivenessResponse:
        """Liveness check — always returns 200 if the process is alive."""
        return LivenessResponse(
            status="alive",
            timestamp=datetime.now(timezone.utc).isoformat(),
            service=settings.service_name,
        )

    # ── GET /ready  — Readiness ────────────────────────────────────

    @router.get(
        "/ready",
        response_model=ReadinessResponse,
        summary="Readiness probe",
        description=(
            "Executes all registered health checks.  Returns 200 if every "
            "component is UP, 503 if any component is DOWN."
        ),
        responses={
            200: {"description": "All components healthy"},
            503: {"description": "One or more components unhealthy"},
        },
    )
    async def readiness(response: Response) -> ReadinessResponse:
        """
        Readiness check — runs all component health checks.

        Returns:
            200 with status ``UP`` when all checks pass.
            503 with status ``DOWN`` when any check fails.
        """
        checks = await health_check.run_all()

        any_down = any(c.status == HealthStatus.DOWN for c in checks)
        any_degraded = any(c.status == HealthStatus.DEGRADED for c in checks)

        if any_down:
            overall_status = HealthStatus.DOWN
            response.status_code = 503
        elif any_degraded:
            overall_status = HealthStatus.DEGRADED
            response.status_code = 200
        else:
            overall_status = HealthStatus.UP

        result = ReadinessResponse(
            status=overall_status,
            timestamp=datetime.now(timezone.utc).isoformat(),
            service=settings.service_name,
            version=settings.service_version,
            checks=checks,
            uptime_seconds=health_check.uptime_seconds,
        )

        logger.info(
            "readiness_check",
            status=overall_status.value,
            components={c.name: c.status.value for c in checks},
            uptime=health_check.uptime_seconds,
        )

        return result

    # ── GET /metrics  — Prometheus ─────────────────────────────────

    @router.get(
        "/metrics",
        summary="Prometheus metrics",
        description=(
            "Returns application metrics in Prometheus exposition format "
            "for scraping by Prometheus / Grafana Agent."
        ),
        response_class=Response,
    )
    async def metrics() -> Response:
        """
        Expose Prometheus metrics.

        Uses ``prometheus_client`` to generate the latest metrics
        snapshot in the standard text exposition format.
        """
        from prometheus_client import (
            CONTENT_TYPE_LATEST,
            CollectorRegistry,
            generate_latest,
            multiprocess,
        )

        # In multi-process mode (gunicorn), we need a fresh registry
        # that merges data from all workers via shared mmap files.
        try:
            registry = CollectorRegistry()
            multiprocess.MultiProcessCollector(registry)
            body = generate_latest(registry)
        except ValueError:
            # Not in multiprocess mode — use the default registry
            body = generate_latest()

        return Response(
            content=body,
            media_type=CONTENT_TYPE_LATEST,
            status_code=200,
        )

    return router
