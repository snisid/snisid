from __future__ import annotations

from unittest.mock import AsyncMock, MagicMock, patch

import pytest

from services.notification.models import WebhookStatus, WebhookSubscription
from services.notification.webhook import WebhookService


@pytest.fixture
def mock_session():
    session = AsyncMock()
    return session


class TestWebhookService:
    async def test_dispatch_no_subscribers(self, mock_session):
        result = MagicMock()
        result.scalars.return_value.all = MagicMock(return_value=[])
        mock_session.execute = AsyncMock(return_value=result)
        service = WebhookService(mock_session)
        await service.dispatch("identity.created", "agg-1", {"key": "value"})

    async def test_dispatch_success(self, mock_session):
        sub = MagicMock(spec=WebhookSubscription)
        sub.id = "sub-1"
        sub.url = "https://example.com/webhook"
        sub.secret = "test-secret"
        sub.events = "identity.created"
        sub.status = WebhookStatus.ACTIVE
        sub.is_deleted = False
        sub.timeout_seconds = 10
        sub.consecutive_failures = 0
        sub.max_consecutive_failures = 10

        result = MagicMock()
        result.scalars.return_value.all = MagicMock(return_value=[sub])
        mock_session.execute = AsyncMock(return_value=result)

        with patch("httpx.AsyncClient") as mock_client:
            mock_client.return_value.__aenter__.return_value.post = AsyncMock()
            mock_client.return_value.__aenter__.return_value.post.return_value.is_success = True
            mock_client.return_value.__aenter__.return_value.post.return_value.status_code = 200
            mock_client.return_value.__aenter__.return_value.post.return_value.text = "OK"

            service = WebhookService(mock_session)
            await service.dispatch("identity.created", "agg-1", {"key": "value"})

    async def test_dispatch_failure_pauses_after_max(self, mock_session):
        sub = MagicMock(spec=WebhookSubscription)
        sub.id = "sub-1"
        sub.url = "https://example.com/webhook"
        sub.secret = "test-secret"
        sub.events = "identity.created"
        sub.status = WebhookStatus.ACTIVE
        sub.is_deleted = False
        sub.timeout_seconds = 10
        sub.consecutive_failures = 9
        sub.max_consecutive_failures = 10

        result = MagicMock()
        result.scalars.return_value.all = MagicMock(return_value=[sub])
        mock_session.execute = AsyncMock(return_value=result)

        with patch("httpx.AsyncClient") as mock_client:
            mock_client.return_value.__aenter__.return_value.post = AsyncMock(
                side_effect=Exception("Connection refused")
            )
            service = WebhookService(mock_session)
            await service.dispatch("identity.created", "agg-1", {"key": "value"})

            assert sub.consecutive_failures == 10
            assert sub.status == WebhookStatus.FAILED
