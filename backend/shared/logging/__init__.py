"""
SNISID Structured Logging
==========================
JSON-structured logging with contextual fields for tracing and audit.
Uses structlog for high-performance structured logging.
"""
from __future__ import annotations

import logging
import sys
from contextvars import ContextVar
from typing import Any

import structlog

# Context variables for request-scoped logging
_trace_id: ContextVar[str] = ContextVar("trace_id", default="")
_correlation_id: ContextVar[str] = ContextVar("correlation_id", default="")
_user_id: ContextVar[str] = ContextVar("user_id", default="")
_service_name: ContextVar[str] = ContextVar("service_name", default="snisid")


def set_log_context(
    *,
    trace_id: str | None = None,
    correlation_id: str | None = None,
    user_id: str | None = None,
    service_name: str | None = None,
) -> None:
    """Set contextual logging fields for the current async context."""
    if trace_id is not None:
        _trace_id.set(trace_id)
    if correlation_id is not None:
        _correlation_id.set(correlation_id)
    if user_id is not None:
        _user_id.set(user_id)
    if service_name is not None:
        _service_name.set(service_name)


def _add_context(
    logger: Any, method_name: str, event_dict: dict[str, Any]
) -> dict[str, Any]:
    """Processor that injects context vars into every log entry."""
    trace_id = _trace_id.get("")
    if trace_id:
        event_dict["trace_id"] = trace_id

    correlation_id = _correlation_id.get("")
    if correlation_id:
        event_dict["correlation_id"] = correlation_id

    user_id = _user_id.get("")
    if user_id:
        event_dict["user_id"] = user_id

    event_dict["service"] = _service_name.get("snisid")
    return event_dict


def _add_call_info(
    logger: Any, method_name: str, event_dict: dict[str, Any]
) -> dict[str, Any]:
    """Add caller information for debug-level logs."""
    if method_name == "debug":
        record = event_dict.get("_record")
        if record:
            event_dict["filename"] = record.filename
            event_dict["lineno"] = record.lineno
            event_dict["func_name"] = record.funcName
    return event_dict


def configure_logging(
    service_name: str = "snisid",
    log_level: str = "INFO",
    json_output: bool = True,
) -> None:
    """
    Configure structured logging for the application.

    Args:
        service_name: Name of the microservice for log identification.
        log_level: Minimum log level (DEBUG, INFO, WARNING, ERROR, CRITICAL).
        json_output: If True, output JSON; if False, output colored console.
    """
    _service_name.set(service_name)

    shared_processors: list[structlog.types.Processor] = [
        structlog.contextvars.merge_contextvars,
        _add_context,
        structlog.stdlib.add_log_level,
        structlog.stdlib.add_logger_name,
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.StackInfoRenderer(),
        structlog.processors.UnicodeDecoder(),
    ]

    if json_output:
        renderer = structlog.processors.JSONRenderer()
    else:
        renderer = structlog.dev.ConsoleRenderer(colors=True)

    structlog.configure(
        processors=[
            *shared_processors,
            structlog.stdlib.ProcessorFormatter.wrap_for_formatter,
        ],
        logger_factory=structlog.stdlib.LoggerFactory(),
        wrapper_class=structlog.stdlib.BoundLogger,
        cache_logger_on_first_use=True,
    )

    formatter = structlog.stdlib.ProcessorFormatter(
        processors=[
            structlog.stdlib.ProcessorFormatter.remove_processors_meta,
            _add_call_info,
            renderer,
        ],
    )

    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(formatter)

    root_logger = logging.getLogger()
    root_logger.handlers.clear()
    root_logger.addHandler(handler)
    root_logger.setLevel(getattr(logging, log_level.upper(), logging.INFO))

    # Suppress noisy third-party loggers
    for logger_name in ("uvicorn.access", "sqlalchemy.engine", "aiokafka"):
        logging.getLogger(logger_name).setLevel(logging.WARNING)


def get_logger(name: str | None = None) -> structlog.stdlib.BoundLogger:
    """
    Get a structured logger instance.

    Args:
        name: Logger name, typically __name__ of the calling module.

    Returns:
        A bound structured logger.
    """
    return structlog.get_logger(name or "snisid")
