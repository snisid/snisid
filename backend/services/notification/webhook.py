from __future__ import annotations

import hashlib
import hmac
import json
import time
from datetime import datetime, timezone
from typing import Any

import httpx
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from services.notification.models import (
    WebhookDeliveryLog,
    WebhookEvent,
    WebhookStatus,
    WebhookSubscription,
)
from shared.logging import get_logger

logger = get_logger(__name__)


class WebhookService:
    def __init__(self, session: AsyncSession) -> None:
        self._session = session

    async def dispatch(self, event_type: str, aggregate_id: str, payload: dict[str, Any]) -> None:
        result = await self._session.execute(
            select(WebhookSubscription).where(
                WebhookSubscription.status == WebhookStatus.ACTIVE,
                WebhookSubscription.events.contains(event_type),
                WebhookSubscription.is_deleted == False,
            )
        )
        subscriptions = result.scalars().all()
        if not subscriptions:
            return
        for sub in subscriptions:
            await self._deliver(sub, event_type, aggregate_id, payload)

    async def _deliver(
        self,
        sub: WebhookSubscription,
        event_type: str,
        aggregate_id: str,
        payload: dict[str, Any],
    ) -> None:
        body = json.dumps({
            "event_type": event_type,
            "aggregate_id": aggregate_id,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "data": payload,
        })
        signature = self._sign(body, sub.secret or "")
        start = time.monotonic()
        try:
            async with httpx.AsyncClient(timeout=sub.timeout_seconds) as client:
                resp = await client.post(
                    sub.url,
                    content=body,
                    headers={
                        "Content-Type": "application/json",
                        "X-SNISID-Signature": signature,
                        "X-SNISID-Event": event_type,
                    },
                )
            duration = int((time.monotonic() - start) * 1000)
            log = WebhookDeliveryLog(
                subscription_id=sub.id,
                event_type=event_type,
                payload=body,
                status_code=resp.status_code,
                response_body=resp.text[:1000],
                success=resp.is_success,
                duration_ms=duration,
            )
            if resp.is_success:
                sub.consecutive_failures = 0
            else:
                sub.consecutive_failures += 1
        except Exception as e:
            duration = int((time.monotonic() - start) * 1000)
            log = WebhookDeliveryLog(
                subscription_id=sub.id,
                event_type=event_type,
                payload=body,
                success=False,
                duration_ms=duration,
            )
            sub.consecutive_failures += 1
            logger.error("webhook_delivery_failed", url=sub.url, error=str(e))
        if sub.consecutive_failures >= sub.max_consecutive_failures:
            sub.status = WebhookStatus.FAILED
            sub.last_failure_reason = f"{sub.consecutive_failures} consecutive failures"
            logger.warn("webhook_auto_paused", url=sub.url, failures=sub.consecutive_failures)
        sub.last_triggered_at = datetime.now(timezone.utc)
        self._session.add(log)
        await self._session.flush()

    @staticmethod
    def _sign(payload: str, secret: str) -> str:
        return hmac.new(secret.encode(), payload.encode(), hashlib.sha256).hexdigest()
