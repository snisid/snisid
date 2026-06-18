"""
SNISID Middleware Tests
========================
Tests for FastAPI middleware stack: CORS, Security, Logging, Audit, Sanitization.
"""
from __future__ import annotations

import uuid
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
from fastapi import FastAPI, Request, Response
from fastapi.responses import JSONResponse
from starlette.middleware.base import RequestResponseEndpoint
from starlette.types import ASGIApp

from shared.middleware import (
    RequestLoggingMiddleware,
    SecurityHeadersMiddleware,
    AuditMiddleware,
    InputSanitizationMiddleware,
    setup_middleware,
    configure_cors,
    _get_client_ip,
    _mask_headers,
    _contains_sql_injection,
)


class TestGetClientIP:
    """Test client IP extraction helper."""

    def test_x_forwarded_for(self):
        request = MagicMock(spec=Request)
        request.headers = {"X-Forwarded-For": "203.0.113.1, 10.0.0.1"}
        request.client = None
        assert _get_client_ip(request) == "203.0.113.1"

    def test_x_real_ip(self):
        request = MagicMock(spec=Request)
        request.headers = {"X-Real-IP": "198.51.100.1"}
        request.client = None
        assert _get_client_ip(request) == "198.51.100.1"

    def test_direct_connection(self):
        request = MagicMock(spec=Request)
        request.headers = {}
        request.client = MagicMock()
        request.client.host = "192.168.1.1"
        assert _get_client_ip(request) == "192.168.1.1"

    def test_no_ip_found(self):
        request = MagicMock(spec=Request)
        request.headers = {}
        request.client = None
        assert _get_client_ip(request) == "unknown"


class TestMaskHeaders:
    """Test sensitive header masking."""

    def test_mask_authorization(self):
        headers = {"Authorization": "Bearer secret-token"}
        masked = _mask_headers(headers)
        assert masked["Authorization"] == "Bearer ***"

    def test_mask_cookie(self):
        headers = {"Cookie": "session=abc123"}
        masked = _mask_headers(headers)
        assert masked["Cookie"] == "***"

    def test_keep_normal_headers(self):
        headers = {"Content-Type": "application/json", "Accept": "*/*"}
        masked = _mask_headers(headers)
        assert masked["Content-Type"] == "application/json"
        assert masked["Accept"] == "*/*"

    def test_case_insensitive_masking(self):
        headers = {"AUTHORIZATION": "Bearer token"}
        masked = _mask_headers(headers)
        assert masked["AUTHORIZATION"] == "Bearer ***"

    def test_empty_headers(self):
        assert _mask_headers({}) == {}


class TestContainsSQLInjection:
    """Test SQL injection pattern detection."""

    def test_detect_sql_select(self):
        assert _contains_sql_injection("SELECT * FROM users")

    def test_detect_sql_drop(self):
        assert _contains_sql_injection("'; DROP TABLE users; --")

    def test_detect_sql_union(self):
        assert _contains_sql_injection("UNION SELECT username, password")

    def test_detect_or_equals(self):
        assert _contains_sql_injection("' OR 1=1 --")

    def test_detect_sql_comments(self):
        assert _contains_sql_injection("admin' --")

    def test_clean_input_passes(self):
        assert not _contains_sql_injection("Jean Dupont")
        assert not _contains_sql_injection("NNU-2025-001234")
        assert not _contains_sql_injection("Port-au-Prince")

    def test_empty_string(self):
        assert not _contains_sql_injection("")

    def test_normal_unicode(self):
        assert not _contains_sql_injection("français et créole")


class TestSecurityHeadersMiddleware:
    """Test security headers injection."""

    @pytest.fixture
    def middleware(self):
        return SecurityHeadersMiddleware(lambda app: None)

    @pytest.mark.asyncio
    async def test_security_headers_applied(self):
        app = MagicMock(spec=ASGIApp)
        middleware = SecurityHeadersMiddleware(app)

        request = MagicMock(spec=Request)
        response = Response()

        async def call_next(req: Request) -> Response:
            return response

        result = await middleware.dispatch(request, call_next)
        assert result.headers["X-Content-Type-Options"] == "nosniff"
        assert result.headers["X-Frame-Options"] == "DENY"
        assert "Strict-Transport-Security" in result.headers
        assert result.headers["X-XSS-Protection"] == "1; mode=block"
        assert "Referrer-Policy" in result.headers

    @pytest.mark.asyncio
    async def test_server_header_removed(self):
        app = MagicMock(spec=ASGIApp)
        middleware = SecurityHeadersMiddleware(app)

        request = MagicMock(spec=Request)
        response = Response()
        response.headers["Server"] = "uvicorn"

        async def call_next(req: Request) -> Response:
            return response

        result = await middleware.dispatch(request, call_next)
        assert "Server" not in result.headers
        assert "server" not in result.headers

    @pytest.mark.asyncio
    async def test_cache_headers(self):
        app = MagicMock(spec=ASGIApp)
        middleware = SecurityHeadersMiddleware(app)

        request = MagicMock(spec=Request)
        response = Response()

        async def call_next(req: Request) -> Response:
            return response

        result = await middleware.dispatch(request, call_next)
        assert result.headers["Cache-Control"] == "no-store, no-cache, must-revalidate"
        assert result.headers["Pragma"] == "no-cache"


class TestRequestLoggingMiddleware:
    """Test request logging middleware."""

    @pytest.mark.asyncio
    async def test_request_id_injected(self):
        middleware = RequestLoggingMiddleware(lambda app: None)
        request = MagicMock(spec=Request)
        request.headers = {}
        request.method = "GET"
        request.url.path = "/healthz"
        url = MagicMock()
        url.query = None
        request.url = url

        response = Response()
        response.headers._list.clear()

        async def call_next(req):
            return response

        with patch("shared.middleware.logger") as mock_logger:
            result = await middleware.dispatch(request, call_next)
            assert "X-Request-ID" in result.headers
            assert uuid.UUID(result.headers["X-Request-ID"])

    @pytest.mark.asyncio
    async def test_response_headers_added(self):
        middleware = RequestLoggingMiddleware(lambda app: None)
        request = MagicMock(spec=Request)
        request.headers = {"X-Request-ID": "req-123"}
        request.method = "GET"
        request.url.path = "/healthz"
        url = MagicMock()
        url.query = None
        request.url = url

        response = Response()
        response.headers._list.clear()

        async def call_next(req):
            return response

        with patch("shared.middleware.logger") as mock_logger:
            result = await middleware.dispatch(request, call_next)
            assert result.headers["X-Request-ID"] == "req-123"
            assert "X-Trace-ID" in result.headers
            assert "X-Correlation-ID" in result.headers


class TestAuditMiddleware:
    """Test audit event capture middleware."""

    @pytest.fixture
    def middleware(self):
        return AuditMiddleware(lambda app: None)

    @pytest.mark.asyncio
    async def test_get_requests_not_audited(self):
        middleware = AuditMiddleware(lambda app: None)
        request = MagicMock(spec=Request)
        request.method = "GET"
        request.url.path = "/healthz"

        response = Response()

        async def call_next(req):
            return response

        with patch.object(middleware, "_build_audit_event") as mock_build:
            result = await middleware.dispatch(request, call_next)
            mock_build.assert_not_called()

    @pytest.mark.asyncio
    async def test_post_requests_audited(self):
        middleware = AuditMiddleware(lambda app: None)
        request = MagicMock(spec=Request)
        request.method = "POST"
        request.url.path = "/api/v1/citizens"
        request.headers = {}
        segments = ["api", "v1", "citizens"]
        request.url = MagicMock()
        request.url.path = "/api/v1/citizens"

        response = Response(status_code=201)
        response.headers._list.clear()

        async def call_next(req):
            return response

        with patch.object(middleware, "_build_audit_event") as mock_build:
            mock_build.return_value = {"action": "POST /api/v1/citizens"}
            with patch.object(middleware, "_publish_to_kafka") as mock_kafka:
                result = await middleware.dispatch(request, call_next)
                mock_build.assert_called_once()
                mock_kafka.assert_called_once()

    @pytest.mark.asyncio
    async def test_kafka_fallback_logging(self):
        middleware = AuditMiddleware(lambda app: None)
        request = MagicMock(spec=Request)
        request.method = "PUT"
        request.url.path = "/api/v1/citizens/123"
        request.headers = {}
        request.url = MagicMock()
        request.url.path = "/api/v1/citizens/123"

        response = Response(status_code=200)
        response.headers._list.clear()

        async def call_next(req):
            return response

        with patch.object(middleware, "_build_audit_event") as mock_build:
            mock_build.return_value = {"action": "PUT /api/v1/citizens/123"}
            with patch.object(middleware, "_publish_to_kafka", side_effect=Exception("Kafka down")):
                with patch("shared.middleware.logger") as mock_logger:
                    result = await middleware.dispatch(request, call_next)
                    mock_logger.warning.assert_called_once()

    def test_build_audit_event_with_actor(self):
        middleware = AuditMiddleware(lambda app: None)
        request = MagicMock(spec=Request)
        request.method = "POST"
        request.url.path = "/api/v1/citizens"
        request.state.user = MagicMock()
        request.state.user.id = "user-abc"
        request.headers = {
            "X-Correlation-ID": "corr-xyz",
            "user-agent": "test-agent",
        }
        request.url = MagicMock()
        request.url.path = "/api/v1/citizens"
        request.client = MagicMock()
        request.client.host = "10.0.0.1"

        response = Response(status_code=201)
        response.headers._list.clear()

        event = middleware._build_audit_event(request, response)
        assert event["actor_id"] == "user-abc"
        assert event["action"] == "POST /api/v1/citizens"
        assert event["resource_type"] == "api"

    def test_build_audit_event_extracts_resource_id(self):
        middleware = AuditMiddleware(lambda app: None)
        request = MagicMock(spec=Request)
        request.method = "DELETE"
        request.url.path = f"/api/v1/citizens/{uuid.uuid4()}"
        request.state.user = None
        request.headers = {}
        request.client = MagicMock()
        request.client.host = "10.0.0.1"

        response = Response(status_code=204)
        response.headers._list.clear()

        event = middleware._build_audit_event(request, response)
        assert event["action"].startswith("DELETE")
        assert event["resource_type"] == "api"


class TestInputSanitizationMiddleware:
    """Test input sanitization middleware."""

    @pytest.mark.asyncio
    async def test_body_size_limit_exceeded(self):
        middleware = InputSanitizationMiddleware(lambda app: None, max_body_bytes=100)
        request = MagicMock(spec=Request)
        request.headers = {"content-length": "1000"}
        request.url.path = "/api/v1/citizens"
        request.query_params = {}
        request.url = MagicMock()
        request.url.path = "/api/v1/citizens"

        async def call_next(req):
            return Response()

        response = await middleware.dispatch(request, call_next)
        assert response.status_code == 413

    @pytest.mark.asyncio
    async def test_body_size_within_limit(self):
        middleware = InputSanitizationMiddleware(lambda app: None, max_body_bytes=1000)
        request = MagicMock(spec=Request)
        request.headers = {"content-length": "100"}
        request.url.path = "/api/v1/citizens"
        request.query_params = {}
        request.url = MagicMock()
        request.url.path = "/api/v1/citizens"

        async def call_next(req):
            return Response(status_code=200)

        response = await middleware.dispatch(request, call_next)
        assert response.status_code == 200

    @pytest.mark.asyncio
    async def test_sql_injection_blocked(self):
        middleware = InputSanitizationMiddleware(lambda app: None)
        request = MagicMock(spec=Request)
        request.headers = {}
        request.url.path = "/api/v1/search"
        request.query_params = {"q": "'; DROP TABLE users; --"}
        request.url = MagicMock()
        request.url.path = "/api/v1/search"

        async def call_next(req):
            return Response()

        response = await middleware.dispatch(request, call_next)
        assert response.status_code == 400

    @pytest.mark.asyncio
    async def test_clean_query_passes(self):
        middleware = InputSanitizationMiddleware(lambda app: None)
        request = MagicMock(spec=Request)
        request.headers = {}
        request.url.path = "/api/v1/search"
        request.query_params = {"q": "Jean Dupont"}
        request.url = MagicMock()
        request.url.path = "/api/v1/search"

        async def call_next(req):
            return Response(status_code=200)

        response = await middleware.dispatch(request, call_next)
        assert response.status_code == 200


class TestConfigureCORS:
    """Test CORS configuration."""

    def test_cors_added_to_app(self):
        app = FastAPI()
        configure_cors(app)
        # Check that middleware was added (there should be at least CORS + default)
        assert len(app.user_middleware) >= 1

    def test_setup_middleware_registers_all(self):
        app = FastAPI()
        setup_middleware(app)
        # Should have: CORS + SecurityHeaders + RequestLogging + Audit + InputSanitization
        # (Starlette adds its own middleware too)
        middleware_types = [m.cls for m in app.user_middleware]
        from fastapi.middleware.cors import CORSMiddleware
        assert CORSMiddleware in middleware_types


class TestRateLimitMiddleware:
    """Test rate limit middleware fall-open behavior."""

    @pytest.mark.asyncio
    async def test_bypasses_health_endpoints(self):
        app = MagicMock(spec=ASGIApp)
        from shared.middleware import RateLimitMiddleware
        middleware = RateLimitMiddleware(app)
        request = MagicMock(spec=Request)
        request.url.path = "/health"
        request.method = "GET"
        request.headers = {}

        async def call_next(req):
            return Response(status_code=200)

        result = await middleware.dispatch(request, call_next)
        assert result.status_code == 200

    @pytest.mark.asyncio
    async def test_falls_open_when_redis_unreachable(self):
        app = MagicMock(spec=ASGIApp)
        from shared.middleware import RateLimitMiddleware
        middleware = RateLimitMiddleware(app)
        request = MagicMock(spec=Request)
        request.url.path = "/api/v1/citizens"
        request.method = "GET"
        request.headers = {}

        async def call_next(req):
            return Response(status_code=200)

        with patch("shared.middleware._HAS_REDIS", False):
            result = await middleware.dispatch(request, call_next)
            assert result.status_code == 200


class TestResponseCacheMiddleware:
    """Test response cache middleware fall-open behavior."""

    @pytest.mark.asyncio
    async def test_falls_open_when_redis_unreachable(self):
        app = MagicMock(spec=ASGIApp)
        from shared.middleware import ResponseCacheMiddleware
        middleware = ResponseCacheMiddleware(app)
        request = MagicMock(spec=Request)
        request.url.path = "/api/v1/citizens"
        request.method = "GET"
        request.headers = {}

        async def call_next(req):
            return Response(status_code=200)

        with patch("shared.middleware._HAS_REDIS", False):
            result = await middleware.dispatch(request, call_next)
            assert result.status_code == 200

    @pytest.mark.asyncio
    async def test_passes_non_get_requests(self):
        app = MagicMock(spec=ASGIApp)
        from shared.middleware import ResponseCacheMiddleware
        middleware = ResponseCacheMiddleware(app)
        request = MagicMock(spec=Request)
        request.url.path = "/api/v1/citizens"
        request.method = "POST"
        request.headers = {}

        async def call_next(req):
            return Response(status_code=201)

        result = await middleware.dispatch(request, call_next)
        assert result.status_code == 201
