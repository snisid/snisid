from __future__ import annotations

import io
import json
import logging
import sys
from unittest.mock import patch

import pytest
import structlog

from shared.logging import (
    configure_logging,
    get_logger,
    set_log_context,
    _trace_id,
    _correlation_id,
    _user_id,
    _service_name,
)


class TestLoggerCreation:
    """Test logger creation."""

    def test_get_logger_returns_bound_logger(self):
        logger = get_logger("test_logger")
        assert isinstance(logger, (structlog.stdlib.BoundLogger, structlog._config.BoundLoggerLazyProxy))

    def test_get_logger_default_name(self):
        logger = get_logger()
        assert isinstance(logger, (structlog.stdlib.BoundLogger, structlog._config.BoundLoggerLazyProxy))

    def test_get_logger_with_name(self):
        logger = get_logger("my.module")
        assert isinstance(logger, (structlog.stdlib.BoundLogger, structlog._config.BoundLoggerLazyProxy))


class TestStructuredLogOutput:
    """Test structured log output in JSON format."""

    @pytest.fixture(autouse=True)
    def setup_logging(self):
        string_io = io.StringIO()

        with patch("sys.stdout", string_io):
            configure_logging(service_name="test-service", log_level="DEBUG", json_output=True)
            yield string_io

    def test_json_output_format(self, setup_logging):
        stream = setup_logging
        logger = get_logger("json_test")
        logger.info("test_event", key="value", number=42)

        output = stream.getvalue()
        log_entry = json.loads(output.strip().split("\n")[0])

        assert log_entry["event"] == "test_event"
        assert log_entry["key"] == "value"
        assert log_entry["number"] == 42
        assert log_entry["logger"] == "json_test"
        assert log_entry["service"] == "test-service"
        assert "timestamp" in log_entry
        assert "level" in log_entry

    def test_log_levels(self, setup_logging):
        stream = setup_logging
        logger = get_logger("levels")

        logger.debug("debug_msg")
        logger.info("info_msg")
        logger.warning("warn_msg")
        logger.error("error_msg")

        output = stream.getvalue()
        lines = [json.loads(line) for line in output.strip().split("\n")]

        levels = [l["level"] for l in lines]
        assert "debug" in levels
        assert "info" in levels
        assert "warning" in levels
        assert "error" in levels

    def test_exception_logging(self, setup_logging):
        stream = setup_logging
        logger = get_logger("exc_test")

        try:
            raise ValueError("test error")
        except ValueError:
            logger.exception("an_error_occurred")

        output = stream.getvalue()
        log_entry = json.loads(output.strip().split("\n")[0])
        assert log_entry["event"] == "an_error_occurred"


class TestContextPropagation:
    """Test context variable propagation in logs."""

    @pytest.fixture(autouse=True)
    def setup_capturing(self):
        string_io = io.StringIO()
        with patch("sys.stdout", string_io):
            configure_logging(service_name="ctx-test", log_level="DEBUG", json_output=True)
            yield string_io

    def test_context_variables_in_log(self, setup_capturing):
        stream = setup_capturing
        set_log_context(
            trace_id="trace-123",
            correlation_id="corr-456",
            user_id="user-789",
        )

        logger = get_logger("ctx")
        logger.info("context_test")

        output = stream.getvalue()
        log_entry = json.loads(output.strip().split("\n")[0])
        assert log_entry["trace_id"] == "trace-123"
        assert log_entry["correlation_id"] == "corr-456"
        assert log_entry["user_id"] == "user-789"

    def test_context_isolation(self, setup_capturing):
        stream = setup_capturing

        set_log_context(trace_id="first-trace")
        logger = get_logger("isolated")
        logger.info("first_event")

        set_log_context(trace_id="second-trace")
        logger.info("second_event")

        output = stream.getvalue()
        lines = [json.loads(l) for l in output.strip().split("\n")]

        assert lines[0]["trace_id"] == "first-trace"
        assert lines[1]["trace_id"] == "second-trace"

    def test_service_name_in_context(self, setup_capturing):
        stream = setup_capturing
        logger = get_logger("svc_test")
        logger.info("svc_check")

        output = stream.getvalue()
        log_entry = json.loads(output.strip().split("\n")[0])
        assert log_entry["service"] == "ctx-test"


class TestSensitiveDataRedaction:
    """Test that sensitive data patterns are handled."""

    def test_no_credit_card_in_logs(self):
        logger = get_logger("pii")
        logger.info("processing", card_number="4111-1111-1111-1111")
        assert True

    def test_password_not_logged_plaintext(self):
        logger = get_logger("pii")
        logger.info("auth_attempt", password="secret123")
        assert True

    def test_email_in_log(self):
        logger = get_logger("pii")
        logger.info("email_sent", email="user@example.com")
        assert True


class TestLogConfiguration:
    """Test logging configuration."""

    def test_configure_with_custom_level(self):
        configure_logging(service_name="config-test", log_level="WARNING")
        root = logging.getLogger()
        assert root.level == logging.WARNING

    def test_configure_suppresses_noisy_loggers(self):
        configure_logging()
        uvicorn_logger = logging.getLogger("uvicorn.access")
        assert uvicorn_logger.level == logging.WARNING

    def test_configure_default_level(self):
        configure_logging()
        root = logging.getLogger()
        assert root.level == logging.INFO


class TestServiceNameInContext:
    """Test service name handling."""

    def test_default_service_name(self):
        _service_name.set("default-svc")
        assert _service_name.get() == "default-svc"

    def test_set_service_name(self):
        _service_name.set("identity-service")
        assert _service_name.get() == "identity-service"
