"""
SNISID OpenTelemetry Tracing & Metrics
=======================================
Initializes OpenTelemetry SDK for distributed tracing and metrics
export using the ``OTEL_*`` environment variables from ``ObservabilityConfig``.

Usage in ``main.py``::

    from shared.telemetry import init_telemetry, close_telemetry

    @asynccontextmanager
    async def lifespan(app):
        init_telemetry()
        yield
        close_telemetry()
"""

from __future__ import annotations

import logging
from typing import Any

import structlog

from shared.config import get_settings

logger = structlog.get_logger(__name__)

_tracer_provider: Any = None
_meter_provider: Any = None


def init_telemetry() -> None:
    """Initialize OpenTelemetry SDK if observability is enabled."""
    settings = get_settings()
    otel = settings.observability

    if not otel.enabled:
        logger.info("telemetry_disabled")
        return

    try:
        from opentelemetry import trace
        from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import (
            OTLPSpanExporter,
        )
        from opentelemetry.sdk.resources import Resource
        from opentelemetry.sdk.trace import TracerProvider
        from opentelemetry.sdk.trace.export import BatchSpanProcessor

        global _tracer_provider

        resource = Resource.create({
            "service.name": otel.service_name,
            "service.version": settings.service_version,
        })

        _tracer_provider = TracerProvider(resource=resource)
        span_processor = BatchSpanProcessor(
            OTLPSpanExporter(endpoint=otel.exporter_endpoint, insecure=True),
        )
        _tracer_provider.add_span_processor(span_processor)
        trace.set_tracer_provider(_tracer_provider)

        logger.info(
            "telemetry_initialized",
            endpoint=otel.exporter_endpoint,
            service=otel.service_name,
        )
    except ImportError:
        logger.warning(
            "telemetry_import_error",
            detail="Install opentelemetry packages to enable tracing",
        )
    except Exception as exc:
        logger.error("telemetry_init_failed", error=str(exc))


def close_telemetry() -> None:
    """Shut down the OpenTelemetry SDK, flushing spans."""
    if _tracer_provider is not None:
        try:
            _tracer_provider.shutdown()
            logger.info("telemetry_shutdown")
        except Exception as exc:
            logger.error("telemetry_shutdown_failed", error=str(exc))


def get_tracer(name: str = "snisid") -> Any:
    """Get an OpenTelemetry tracer instance."""
    try:
        from opentelemetry import trace

        return trace.get_tracer(name)
    except ImportError:
        return None
