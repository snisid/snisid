"""Integration tests: health endpoints, security headers, auth middleware."""
from __future__ import annotations

import os

import pytest
from fastapi import FastAPI
from fastapi.responses import JSONResponse
from httpx import AsyncClient, ASGITransport

from shared.health import HealthCheck, create_health_router
from shared.middleware import setup_middleware

API_KEY = "test-api-key-123"


@pytest.fixture
def app() -> FastAPI:
    os.environ["ENVIRONMENT"] = "production"
    _app = FastAPI()

    health = HealthCheck()
    _app.include_router(create_health_router(health))

    setup_middleware(_app)

    @_app.middleware("http")
    async def api_key_middleware(request, call_next):
        if request.url.path in ("/health", "/ready", "/metrics", "/docs", "/openapi.json", "/redoc"):
            return await call_next(request)
        api_key = request.headers.get("X-API-Key")
        if not api_key or api_key != API_KEY:
            return JSONResponse(status_code=401, content={"detail": "Invalid or missing API key"})
        return await call_next(request)

    @_app.get("/protected")
    async def protected():
        return {"message": "you have access"}

    return _app


class TestAuthMiddleware:
    @pytest.mark.asyncio
    async def test_health_bypasses_auth(self, app: FastAPI):
        async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
            resp = await ac.get("/health")
        assert resp.status_code == 200

    @pytest.mark.asyncio
    async def test_protected_blocked_without_key(self, app: FastAPI):
        async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
            resp = await ac.get("/protected")
        assert resp.status_code == 401

    @pytest.mark.asyncio
    async def test_protected_allowed_with_key(self, app: FastAPI):
        async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
            resp = await ac.get("/protected", headers={"X-API-Key": API_KEY})
        assert resp.status_code == 200

    @pytest.mark.asyncio
    async def test_wrong_key_rejected(self, app: FastAPI):
        async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
            resp = await ac.get("/protected", headers={"X-API-Key": "wrong-key"})
        assert resp.status_code == 401


@pytest.mark.asyncio
async def test_liveness_returns_200(app: FastAPI):
    async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
        resp = await ac.get("/health")
    assert resp.status_code == 200
    data = resp.json()
    assert data["status"] == "alive"
    assert "timestamp" in data
    assert "service" in data


@pytest.mark.asyncio
async def test_readiness_returns_200(app: FastAPI):
    async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
        resp = await ac.get("/ready")
    assert resp.status_code == 200
    data = resp.json()
    assert data["status"] in ("UP", "DEGRADED", "DOWN")
    assert "checks" in data
    assert "uptime_seconds" in data


@pytest.mark.asyncio
async def test_health_not_found(app: FastAPI):
    async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
        resp = await ac.get("/nonexistent", headers={"X-API-Key": API_KEY})
    assert resp.status_code == 404


class TestSecurityHeadersIntegration:
    """Security headers applied to all responses via middleware."""

    @pytest.mark.asyncio
    async def test_security_headers_present_on_health(self, app: FastAPI):
        async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
            resp = await ac.get("/health")
        assert resp.headers.get("X-Content-Type-Options") == "nosniff"
        assert resp.headers.get("X-Frame-Options") == "DENY"
        assert resp.headers.get("X-XSS-Protection") == "1; mode=block"
        assert "Strict-Transport-Security" in resp.headers
        assert "Referrer-Policy" in resp.headers

    @pytest.mark.asyncio
    async def test_security_headers_on_404(self, app: FastAPI):
        async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
            resp = await ac.get("/nonexistent", headers={"X-API-Key": API_KEY})
        assert resp.headers.get("X-Content-Type-Options") == "nosniff"
        assert resp.headers.get("X-Frame-Options") == "DENY"
        assert resp.headers.get("X-XSS-Protection") == "1; mode=block"
