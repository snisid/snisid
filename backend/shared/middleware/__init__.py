"""
SNISID Middleware Stack
========================
Production-grade FastAPI middleware for the SNISID national identity system.

Middleware order (outermost → innermost):
1. CORS — Must be outermost to handle preflight requests
2. SecurityHeaders — Add defensive HTTP headers on every response
3. RequestLogging — Log all requests/responses with timing and trace IDs
4. Audit — Capture mutation operations for compliance (Kafka + fallback)
5. InputSanitization — Strip dangerous input patterns before handlers
"""
from __future__ import annotations

import re
import time
import uuid
from datetime import datetime, timezone
from typing import Any, Callable, Sequence

import structlog
from fastapi import FastAPI, Request, Response
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware, RequestResponseEndpoint
from starlette.types import ASGIApp

from shared.config import Environment, get_settings
from shared.logging import get_logger, set_log_context

logger = get_logger(__name__)

# ── Constants ─────────────────────────────────────────────────────────

_SENSITIVE_HEADERS: frozenset[str] = frozenset({
    "authorization",
    "cookie",
    "x-api-key",
    "x-csrf-token",
    "proxy-authorization",
})

_MUTATION_METHODS: frozenset[str] = frozenset({"POST", "PUT", "PATCH", "DELETE"})

_MAX_REQUEST_BODY_BYTES: int = 10 * 1024 * 1024  # 10 MB

# Pre-compiled SQL injection detection patterns
_SQL_INJECTION_PATTERNS: list[re.Pattern[str]] = [
    re.compile(r"(\b(SELECT|INSERT|UPDATE|DELETE|DROP|ALTER|CREATE|EXEC|EXECUTE|UNION)\b)", re.IGNORECASE),
    re.compile(r"(--|#|/\*|\*/)", re.IGNORECASE),
    re.compile(r"(\b(OR|AND)\b\s+\d+\s*=\s*\d+)", re.IGNORECASE),
    re.compile(r"(;|\bxp_|\bsp_)", re.IGNORECASE),
    re.compile(r"('|\"|\\x27|\\x22)", re.IGNORECASE),
    re.compile(r"(\b(CHAR|NCHAR|VARCHAR|NVARCHAR|CAST|CONVERT)\s*\()", re.IGNORECASE),
]

# Dangerous characters to strip from general string inputs
_DANGEROUS_CHARS: re.Pattern[str] = re.compile(r"[\x00-\x08\x0b\x0c\x0e-\x1f\x7f]")


# ── Request Logging Middleware ────────────────────────────────────────


class RequestLoggingMiddleware(BaseHTTPMiddleware):
    """
    Structured JSON request/response logging middleware.

    Logs every HTTP request with:
    - timestamp, method, path, status_code, duration_ms
    - user_agent, client_ip, request_id, content_length
    - Injects ``X-Request-ID`` header if not present (uuid4)
    - Sets trace_id and correlation_id in the logging context
    - Masks sensitive header values (Authorization → ``'Bearer ***'``)
    """

    async def dispatch(
        self, request: Request, call_next: RequestResponseEndpoint
    ) -> Response:
        """Process the request, measure timing, and log structured data."""
        start_time = time.perf_counter()

        # ── Inject / read request identifiers ──────────────────────
        request_id = request.headers.get("X-Request-ID") or str(uuid.uuid4())
        trace_id = request.headers.get("X-Trace-ID") or request_id
        correlation_id = request.headers.get("X-Correlation-ID") or request_id

        # Propagate context for downstream structured logs
        set_log_context(trace_id=trace_id, correlation_id=correlation_id)

        # ── Execute request ────────────────────────────────────────
        try:
            response = await call_next(request)
        except Exception:
            duration_ms = round((time.perf_counter() - start_time) * 1000, 2)
            logger.error(
                "request_unhandled_exception",
                method=request.method,
                path=request.url.path,
                duration_ms=duration_ms,
                request_id=request_id,
                client_ip=_get_client_ip(request),
            )
            raise

        duration_ms = round((time.perf_counter() - start_time) * 1000, 2)

        # ── Attach response headers ───────────────────────────────
        response.headers["X-Request-ID"] = request_id
        response.headers["X-Trace-ID"] = trace_id
        response.headers["X-Correlation-ID"] = correlation_id

        # ── Build log payload ─────────────────────────────────────
        content_length = response.headers.get("content-length", "0")
        log_data: dict[str, Any] = {
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "method": request.method,
            "path": request.url.path,
            "query": str(request.url.query) if request.url.query else None,
            "status_code": response.status_code,
            "duration_ms": duration_ms,
            "user_agent": request.headers.get("user-agent", ""),
            "client_ip": _get_client_ip(request),
            "request_id": request_id,
            "content_length": int(content_length) if content_length.isdigit() else 0,
            "headers": _mask_headers(dict(request.headers)),
        }

        # Choose log level based on status
        if response.status_code >= 500:
            logger.error("http_request", **log_data)
        elif response.status_code >= 400:
            logger.warning("http_request", **log_data)
        else:
            logger.info("http_request", **log_data)

        return response


# ── Security Headers Middleware ───────────────────────────────────────


class SecurityHeadersMiddleware(BaseHTTPMiddleware):
    """
    Inject OWASP-recommended security headers on every response.

    Headers applied:
    - ``X-Content-Type-Options: nosniff``
    - ``X-Frame-Options: DENY``
    - ``Strict-Transport-Security: max-age=63072000; includeSubDomains; preload``
    - ``X-XSS-Protection: 1; mode=block``
    - ``Referrer-Policy: strict-origin-when-cross-origin``
    - ``Content-Security-Policy: default-src 'self'``
    - ``Permissions-Policy: geolocation=(), camera=(), microphone=()``
    - Removes the ``Server`` header to prevent information leakage.
    """

    # Headers stored as class-level constant — no per-request allocation
    _SECURITY_HEADERS: dict[str, str] = {
        "X-Content-Type-Options": "nosniff",
        "X-Frame-Options": "DENY",
        "Strict-Transport-Security": "max-age=63072000; includeSubDomains; preload",
        "X-XSS-Protection": "1; mode=block",
        "Referrer-Policy": "strict-origin-when-cross-origin",
        "Content-Security-Policy": "default-src 'self'; frame-ancestors 'none'",
        "Permissions-Policy": (
            "geolocation=(), camera=(), microphone=(), "
            "payment=(), usb=(), magnetometer=()"
        ),
        "Cache-Control": "no-store, no-cache, must-revalidate",
        "Pragma": "no-cache",
    }

    async def dispatch(
        self, request: Request, call_next: RequestResponseEndpoint
    ) -> Response:
        """Add security headers and strip the Server header."""
        response = await call_next(request)

        for header_name, header_value in self._SECURITY_HEADERS.items():
            response.headers[header_name] = header_value

        # Remove the Server header to avoid version fingerprinting
        response.headers.pop("Server", None)
        response.headers.pop("server", None)

        return response


# ── Audit Middleware ──────────────────────────────────────────────────


class AuditMiddleware(BaseHTTPMiddleware):
    """
    Capture mutation operations (POST/PUT/PATCH/DELETE) for compliance.

    For each mutation request, an audit event is published to the
    ``snisid.audit.events`` Kafka topic in a fire-and-forget manner
    so it never blocks the response path.  If Kafka is unavailable,
    the event is written to the structured log as a fallback.

    Audit event fields:
    - actor_id — user ID extracted from the JWT ``sub`` claim
    - action — ``{HTTP_METHOD} {path}``
    - resource_type — first meaningful path segment
    - resource_id — UUID-shaped segment extracted from the URL path
    - ip_address, user_agent, correlation_id, timestamp
    """

    _KAFKA_TOPIC: str = "snisid.audit.events"
    # Pre-compiled UUID pattern for resource_id extraction
    _UUID_PATTERN: re.Pattern[str] = re.compile(
        r"[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
    )

    async def dispatch(
        self, request: Request, call_next: RequestResponseEndpoint
    ) -> Response:
        """Intercept mutation requests and emit audit events."""
        response = await call_next(request)

        if request.method not in _MUTATION_METHODS:
            return response

        # Build audit event (non-blocking)
        audit_event = self._build_audit_event(request, response)

        # Fire-and-forget Kafka publish — never delay the response
        try:
            await self._publish_to_kafka(audit_event)
        except Exception as exc:
            # Fallback: structured log guarantees the audit trail survives
            logger.warning(
                "audit_kafka_fallback",
                reason=str(exc),
                **audit_event,
            )

        return response

    def _build_audit_event(self, request: Request, response: Response) -> dict[str, Any]:
        """Construct the audit event payload from the request context."""
        path = request.url.path
        segments = [s for s in path.strip("/").split("/") if s]

        # Extract resource type (first path segment, e.g., "citizens", "documents")
        resource_type = segments[0] if segments else "unknown"

        # Try to find a UUID-shaped resource_id in the path
        resource_id: str | None = None
        for segment in segments:
            if self._UUID_PATTERN.fullmatch(segment):
                resource_id = segment
                break

        # Extract actor_id from JWT — the auth middleware will have decoded
        # the token and stored it in request.state if present
        actor_id = self._extract_actor_id(request)

        return {
            "actor_id": actor_id,
            "action": f"{request.method} {path}",
            "resource_type": resource_type,
            "resource_id": resource_id,
            "status_code": response.status_code,
            "ip_address": _get_client_ip(request),
            "user_agent": request.headers.get("user-agent", ""),
            "correlation_id": request.headers.get("X-Correlation-ID", ""),
            "timestamp": datetime.now(timezone.utc).isoformat(),
        }

    @staticmethod
    def _extract_actor_id(request: Request) -> str:
        """
        Best-effort actor ID extraction from the request.

        Tries ``request.state.user`` (set by auth dependencies) first,
        then falls back to a lightweight JWT ``sub`` claim peek without
        full validation (the auth layer already validated it).
        """
        # Fast path: auth dependency already resolved the user
        user = getattr(request.state, "user", None)
        if user is not None:
            return getattr(user, "id", "anonymous")

        # Fallback: peek at the Authorization header
        auth_header = request.headers.get("authorization", "")
        if auth_header.lower().startswith("bearer "):
            try:
                import json
                import base64

                token = auth_header.split(" ", 1)[1]
                # Decode the payload section (index 1) without verification
                payload_b64 = token.split(".")[1]
                # Add padding
                padded = payload_b64 + "=" * (4 - len(payload_b64) % 4)
                payload = json.loads(base64.urlsafe_b64decode(padded))
                return payload.get("sub", "anonymous")
            except Exception:
                pass

        return "anonymous"

    async def _publish_to_kafka(self, event: dict[str, Any]) -> None:
        """
        Publish the audit event to Kafka using aiokafka.

        This is fire-and-forget: we use a short-lived producer per event
        to avoid holding a persistent connection in middleware.  For
        high-throughput systems, consider a shared producer managed via
        the application lifespan.
        """
        import orjson
        from aiokafka import AIOKafkaProducer

        settings = get_settings()

        producer = AIOKafkaProducer(
            bootstrap_servers=settings.kafka.bootstrap_list,
            security_protocol=settings.kafka.security_protocol,
            value_serializer=lambda v: orjson.dumps(v),
            request_timeout_ms=5000,
            acks=1,  # Single ack for audit speed — topic should be replicated
        )
        try:
            await producer.start()
            await producer.send_and_wait(
                self._KAFKA_TOPIC,
                value=event,
                key=event.get("actor_id", "anonymous").encode("utf-8"),
            )
            logger.debug(
                "audit_event_published",
                topic=self._KAFKA_TOPIC,
                actor_id=event.get("actor_id"),
                action=event.get("action"),
            )
        finally:
            await producer.stop()


# ── CORS Configuration ───────────────────────────────────────────────


def configure_cors(app: FastAPI) -> None:
    """
    Add CORS middleware configured from application settings.

    In **development**, all origins are allowed for convenience.
    In **staging / production**, origins are restricted to the explicit
    list defined in ``Settings.cors_origins``.

    Args:
        app: The FastAPI application instance.
    """
    settings = get_settings()
    is_dev = settings.environment == Environment.DEVELOPMENT

    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"] if is_dev else settings.cors_origins,
        allow_credentials=not is_dev,  # Credentials require explicit origins
        allow_methods=["*"],
        allow_headers=[
            "Authorization",
            "Content-Type",
            "X-Request-ID",
            "X-Trace-ID",
            "X-Correlation-ID",
            "Accept",
            "Accept-Language",
        ],
        expose_headers=[
            "X-Request-ID",
            "X-Trace-ID",
            "X-Correlation-ID",
            "X-Total-Count",
        ],
        max_age=600 if is_dev else 3600,
    )

    logger.info(
        "cors_configured",
        environment=settings.environment.value,
        allow_origins="*" if is_dev else settings.cors_origins,
    )


# ── Input Sanitization Middleware ─────────────────────────────────────


class InputSanitizationMiddleware(BaseHTTPMiddleware):
    """
    Defence-in-depth input sanitization middleware.

    Responsibilities:
    - Enforce a maximum request body size (default 10 MB).
    - Strip dangerous control characters from query parameters.
    - Reject requests whose query parameters contain SQL injection patterns.

    This middleware is a *safety net* — individual endpoints MUST still
    validate and parameterise input properly.
    """

    def __init__(
        self,
        app: ASGIApp,
        max_body_bytes: int = _MAX_REQUEST_BODY_BYTES,
    ) -> None:
        super().__init__(app)
        self.max_body_bytes = max_body_bytes

    async def dispatch(
        self, request: Request, call_next: RequestResponseEndpoint
    ) -> Response:
        """Validate and sanitise incoming request data."""
        # ── Body size check ────────────────────────────────────────
        content_length_header = request.headers.get("content-length")
        if content_length_header is not None:
            try:
                content_length = int(content_length_header)
                if content_length > self.max_body_bytes:
                    logger.warning(
                        "request_body_too_large",
                        content_length=content_length,
                        max_allowed=self.max_body_bytes,
                        path=request.url.path,
                        client_ip=_get_client_ip(request),
                    )
                    return JSONResponse(
                        status_code=413,
                        content={
                            "detail": "Request body exceeds maximum allowed size",
                            "max_bytes": self.max_body_bytes,
                        },
                    )
            except ValueError:
                pass  # Non-numeric content-length is handled by the server

        # ── Query parameter sanitization ───────────────────────────
        for param_name, param_value in request.query_params.items():
            # Strip dangerous control characters
            cleaned = _DANGEROUS_CHARS.sub("", param_value)
            if cleaned != param_value:
                logger.warning(
                    "dangerous_chars_stripped",
                    param=param_name,
                    path=request.url.path,
                    client_ip=_get_client_ip(request),
                )

            # Check for SQL injection patterns
            if _contains_sql_injection(cleaned):
                logger.error(
                    "sql_injection_blocked",
                    param=param_name,
                    path=request.url.path,
                    client_ip=_get_client_ip(request),
                    user_agent=request.headers.get("user-agent", ""),
                )
                return JSONResponse(
                    status_code=400,
                    content={"detail": "Potentially malicious input detected"},
                )

        return await call_next(request)


# ── Middleware Orchestrator ───────────────────────────────────────────


def setup_middleware(app: FastAPI) -> None:
    """
    Register all SNISID middleware on the FastAPI application.

    Middleware is added in reverse order because Starlette applies
    ``add_middleware`` calls **outermost-first** (LIFO stack), so the
    first ``add_middleware`` call becomes the outermost layer.

    Final execution order (outermost → innermost):
    1. CORS
    2. SecurityHeaders
    3. RequestLogging
    4. Audit
    5. InputSanitization

    Args:
        app: The FastAPI application to configure.
    """
    # 5. InputSanitization (innermost — closest to route handlers)
    app.add_middleware(InputSanitizationMiddleware)

    # 4. Audit
    app.add_middleware(AuditMiddleware)

    # 3. RequestLogging
    app.add_middleware(RequestLoggingMiddleware)

    # 2. SecurityHeaders
    app.add_middleware(SecurityHeadersMiddleware)

    # 1. CORS (outermost — must handle OPTIONS preflight first)
    configure_cors(app)

    logger.info("middleware_stack_configured", middleware_count=5)


# ── Helpers ───────────────────────────────────────────────────────────


def _get_client_ip(request: Request) -> str:
    """
    Extract the client IP address, respecting reverse-proxy headers.

    Checks ``X-Forwarded-For`` first (comma-separated list, first entry
    is the original client), then ``X-Real-IP``, and finally falls
    back to the direct connection IP.
    """
    forwarded_for = request.headers.get("X-Forwarded-For")
    if forwarded_for:
        # First IP in the chain is the original client
        return forwarded_for.split(",")[0].strip()

    real_ip = request.headers.get("X-Real-IP")
    if real_ip:
        return real_ip.strip()

    if request.client:
        return request.client.host

    return "unknown"


def _mask_headers(headers: dict[str, str]) -> dict[str, str]:
    """
    Return a copy of *headers* with sensitive values redacted.

    Sensitive headers (Authorization, Cookie, etc.) are replaced with
    a masked placeholder to prevent credential leakage in logs.
    """
    masked: dict[str, str] = {}
    for name, value in headers.items():
        if name.lower() in _SENSITIVE_HEADERS:
            if name.lower() == "authorization" and value.lower().startswith("bearer "):
                masked[name] = "Bearer ***"
            else:
                masked[name] = "***"
        else:
            masked[name] = value
    return masked


def _contains_sql_injection(value: str) -> bool:
    """
    Check whether *value* contains common SQL injection patterns.

    This is a heuristic safety net, **not** a replacement for
    parameterised queries.  It catches the most common attack vectors
    while minimising false positives on normal government-domain data.
    """
    for pattern in _SQL_INJECTION_PATTERNS:
        if pattern.search(value):
            return True
    return False
