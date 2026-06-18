from __future__ import annotations

import time
from datetime import datetime, timezone
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
import pytest_asyncio
from fastapi import FastAPI
from fastapi.testclient import TestClient

from shared.health import (
    HealthCheck,
    HealthStatus,
    ComponentHealth,
    ReadinessResponse,
    LivenessResponse,
    create_health_router,
    check_database,
    check_redis,
    check_kafka,
)


class TestHealthCheckRegistry:
    """Test health check registration and execution."""

    @pytest.mark.asyncio
    async def test_register_check(self):
        health = HealthCheck()
        check_fn = AsyncMock(return_value=ComponentHealth(
            name="test", status=HealthStatus.UP
        ))
        health.register("test", check_fn)
        assert "test" in health.registered_checks

    @pytest.mark.asyncio
    async def test_register_duplicate_raises(self):
        health = HealthCheck()
        health.register("dup", AsyncMock())
        with pytest.raises(ValueError, match="already registered"):
            health.register("dup", AsyncMock())

    @pytest.mark.asyncio
    async def test_unregister_check(self):
        health = HealthCheck()
        health.register("temp", AsyncMock())
        health.unregister("temp")
        assert "temp" not in health.registered_checks

    @pytest.mark.asyncio
    async def test_unregister_nonexistent(self):
        health = HealthCheck()
        health.unregister("nonexistent")
        assert True

    @pytest.mark.asyncio
    async def test_run_all_success(self):
        health = HealthCheck()
        health.register("db", AsyncMock(return_value=ComponentHealth(
            name="db", status=HealthStatus.UP, details={"version": "15"}
        )))
        health.register("redis", AsyncMock(return_value=ComponentHealth(
            name="redis", status=HealthStatus.UP, details={"ping": "PONG"}
        )))

        results = await health.run_all()
        assert len(results) == 2
        assert all(r.status == HealthStatus.UP for r in results)

    @pytest.mark.asyncio
    async def test_run_all_with_failures(self):
        health = HealthCheck()
        health.register("good", AsyncMock(return_value=ComponentHealth(
            name="good", status=HealthStatus.UP
        )))
        health.register("bad", AsyncMock(return_value=ComponentHealth(
            name="bad", status=HealthStatus.DOWN, details={"error": "timeout"}
        )))

        results = await health.run_all()
        assert len(results) == 2
        statuses = {r.name: r.status for r in results}
        assert statuses["good"] == HealthStatus.UP
        assert statuses["bad"] == HealthStatus.DOWN

    @pytest.mark.asyncio
    async def test_run_all_handles_exceptions(self):
        health = HealthCheck()
        health.register("crash", AsyncMock(side_effect=RuntimeError("unexpected")))

        results = await health.run_all()
        assert len(results) == 1
        assert results[0].status == HealthStatus.DOWN
        assert "unexpected" in results[0].details["error"]

    @pytest.mark.asyncio
    async def test_uptime(self):
        health = HealthCheck()
        start = health.uptime_seconds
        time.sleep(0.01)
        assert health.uptime_seconds > start


class TestBuiltinChecks:
    """Test built-in health check functions."""

    @pytest.mark.asyncio
    async def test_check_database_not_initialized(self):
        with patch("shared.database._engine", None):
            result = await check_database()
            assert result.status == HealthStatus.DOWN
            assert "not initialised" in result.details.get("error", "")

    @pytest.mark.asyncio
    async def test_check_database_error(self):
        mock_engine = MagicMock()
        mock_engine.connect = MagicMock(side_effect=Exception("connection refused"))

        with patch("shared.database._engine", mock_engine):
            result = await check_database()
            assert result.status == HealthStatus.DOWN
            assert "connection refused" in result.details.get("error", "")

    @pytest.mark.asyncio
    async def test_check_redis_error(self):
        with patch("redis.asyncio.from_url", side_effect=Exception("redis down")):
            result = await check_redis()
            assert result.status == HealthStatus.DOWN
            assert "redis down" in result.details.get("error", "")

    @pytest.mark.asyncio
    async def test_check_kafka_error(self):
        with patch("aiokafka.AIOKafkaProducer") as mock_cls:
            mock_instance = AsyncMock()
            mock_instance.start = AsyncMock(side_effect=Exception("kafka down"))
            mock_cls.return_value = mock_instance

            result = await check_kafka()
            assert result.status == HealthStatus.DOWN
            assert "kafka" in result.details.get("error", "")


class TestComponentHealthModel:
    """Test ComponentHealth model."""

    def test_create_up_component(self):
        ch = ComponentHealth(
            name="database",
            status=HealthStatus.UP,
            details={"version": "15"},
        )
        assert ch.name == "database"
        assert ch.status == HealthStatus.UP
        assert ch.latency_ms == 0.0

    def test_create_down_component(self):
        ch = ComponentHealth(
            name="redis",
            status=HealthStatus.DOWN,
            details={"error": "timeout"},
            latency_ms=150.5,
        )
        assert ch.status == HealthStatus.DOWN
        assert ch.latency_ms == 150.5


class TestLivenessResponse:
    """Test liveness response model."""

    def test_liveness_response(self):
        resp = LivenessResponse(
            status="alive",
            timestamp=datetime.now(timezone.utc).isoformat(),
            service="snisid",
        )
        assert resp.status == "alive"
        assert resp.service == "snisid"


class TestReadinessResponse:
    """Test readiness response model."""

    def test_readiness_response_up(self):
        resp = ReadinessResponse(
            status=HealthStatus.UP,
            timestamp=datetime.now(timezone.utc).isoformat(),
            service="snisid",
            version="1.0.0",
            checks=[
                ComponentHealth(name="db", status=HealthStatus.UP),
            ],
        )
        assert resp.status == HealthStatus.UP
        assert resp.uptime_seconds == 0.0


class TestHealthRouter:
    """Test health check HTTP endpoints."""

    @pytest.fixture
    def app(self):
        app = FastAPI()
        health = HealthCheck()
        health.register("mock", AsyncMock(return_value=ComponentHealth(
            name="mock", status=HealthStatus.UP
        )))
        app.include_router(create_health_router(health))
        return app

    @pytest.fixture
    def client(self, app):
        return TestClient(app)

    def test_liveness_endpoint(self, client):
        response = client.get("/health")
        assert response.status_code == 200
        data = response.json()
        assert data["status"] == "alive"
        assert "timestamp" in data
        assert "service" in data

    def test_readiness_endpoint_all_up(self, client):
        response = client.get("/ready")
        assert response.status_code == 200
        data = response.json()
        assert data["status"] == "UP"
        assert len(data["checks"]) > 0

    def test_readiness_endpoint_with_down(self):
        app = FastAPI()
        health = HealthCheck()
        health.register("failing", AsyncMock(return_value=ComponentHealth(
            name="failing", status=HealthStatus.DOWN, details={"error": "fail"}
        )))
        app.include_router(create_health_router(health))
        client = TestClient(app)

        response = client.get("/ready")
        assert response.status_code == 503
        data = response.json()
        assert data["status"] == "DOWN"

    def test_readiness_endpoint_degraded(self):
        app = FastAPI()
        health = HealthCheck()
        health.register("degraded", AsyncMock(return_value=ComponentHealth(
            name="degraded", status=HealthStatus.DEGRADED, details={"slow": True}
        )))
        app.include_router(create_health_router(health))
        client = TestClient(app)

        response = client.get("/ready")
        assert response.status_code == 200
        data = response.json()
        assert data["status"] == "DEGRADED"

    def test_metrics_endpoint(self, app):
        client = TestClient(app)
        response = client.get("/metrics")
        assert response.status_code == 200
