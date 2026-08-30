from __future__ import annotations

import uuid
from datetime import datetime, timezone
from enum import Enum
from typing import Any, Callable

from services.vc import VCStatus


class RevocationEventType(str, Enum):
    CREDENTIAL_REVOKED = "credential.revoked"
    CREDENTIAL_SUSPENDED = "credential.suspended"
    CREDENTIAL_UNSUSPENDED = "credential.unsuspended"
    CREDENTIAL_REINSTATED = "credential.reinstated"


class RevocationEvent:
    def __init__(
        self,
        credential_id: str,
        event_type: RevocationEventType,
        subject_id: str,
        reason: str = "",
        timestamp: str | None = None,
    ):
        self.id = str(uuid.uuid4())
        self.credential_id = credential_id
        self.event_type = event_type
        self.subject_id = subject_id
        self.reason = reason
        self.timestamp = timestamp or datetime.now(timezone.utc).isoformat()

    def to_dict(self) -> dict[str, Any]:
        return {
            "id": self.id,
            "credential_id": self.credential_id,
            "event_type": self.event_type.value,
            "subject_id": self.subject_id,
            "reason": self.reason,
            "timestamp": self.timestamp,
        }


RevocationCallback = Callable[[RevocationEvent], None]


class RevocationNotifier:
    def __init__(self, storage: Any | None = None):
        self._subscribers: dict[str, list[RevocationCallback]] = {}
        self._history: list[RevocationEvent] = []
        self._wallet_hooks: dict[str, list[str]] = {}  # wallet_did -> [credential_ids]
        self._storage = storage

    def subscribe(self, event_type: RevocationEventType | str, callback: RevocationCallback):
        key = event_type.value if isinstance(event_type, RevocationEventType) else event_type
        if key not in self._subscribers:
            self._subscribers[key] = []
        self._subscribers[key].append(callback)

    def register_wallet(self, wallet_did: str, credential_id: str):
        if wallet_did not in self._wallet_hooks:
            self._wallet_hooks[wallet_did] = []
        self._wallet_hooks[wallet_did].append(credential_id)

    def unregister_wallet(self, wallet_did: str, credential_id: str):
        if wallet_did in self._wallet_hooks:
            self._wallet_hooks[wallet_did] = [
                cid for cid in self._wallet_hooks[wallet_did] if cid != credential_id
            ]

    def get_wallet_credentials(self, wallet_did: str) -> list[str]:
        return self._wallet_hooks.get(wallet_did, [])

    def notify(
        self,
        credential_id: str,
        event_type: RevocationEventType,
        subject_id: str,
        reason: str = "",
    ) -> RevocationEvent:
        event = RevocationEvent(
            credential_id=credential_id,
            event_type=event_type,
            subject_id=subject_id,
            reason=reason,
        )
        self._history.append(event)

        key = event_type.value
        if key in self._subscribers:
            for cb in self._subscribers[key]:
                cb(event)

        return event

    async def async_notify(
        self,
        credential_id: str,
        event_type: RevocationEventType,
        subject_id: str,
        reason: str = "",
    ) -> RevocationEvent:
        event = self.notify(credential_id, event_type, subject_id, reason)
        if self._storage:
            await self._storage.save(
                event_id=event.id,
                credential_id=event.credential_id,
                subject_id=event.subject_id,
                event_type=event.event_type.value,
                reason=event.reason,
            )
        return event

    def get_history(
        self,
        credential_id: str | None = None,
        subject_id: str | None = None,
        limit: int = 50,
    ) -> list[RevocationEvent]:
        results = self._history
        if credential_id:
            results = [e for e in results if e.credential_id == credential_id]
        if subject_id:
            results = [e for e in results if e.subject_id == subject_id]
        return results[-limit:]

    async def async_get_history(
        self,
        credential_id: str | None = None,
        subject_id: str | None = None,
        limit: int = 50,
    ) -> list[RevocationEvent]:
        if self._storage:
            records = await self._storage.get_history(
                credential_id=credential_id,
                subject_id=subject_id,
                limit=limit,
            )
            events = []
            for r in records:
                evt = RevocationEvent(
                    credential_id=r.credential_id,
                    event_type=RevocationEventType(r.event_type),
                    subject_id=r.subject_id,
                    reason=r.reason,
                    timestamp=r.created_at.isoformat() if r.created_at else None,
                )
                evt.id = r.event_id
                events.append(evt)
            return events
        return self.get_history(credential_id, subject_id, limit)

    def notify_revocation(self, credential_id: str, subject_id: str, reason: str = ""):
        return self.notify(
            credential_id=credential_id,
            event_type=RevocationEventType.CREDENTIAL_REVOKED,
            subject_id=subject_id,
            reason=reason,
        )

    def notify_suspension(self, credential_id: str, subject_id: str, reason: str = ""):
        return self.notify(
            credential_id=credential_id,
            event_type=RevocationEventType.CREDENTIAL_SUSPENDED,
            subject_id=subject_id,
            reason=reason,
        )

    def notify_reinstatement(self, credential_id: str, subject_id: str, reason: str = ""):
        return self.notify(
            credential_id=credential_id,
            event_type=RevocationEventType.CREDENTIAL_REINSTATED,
            subject_id=subject_id,
            reason=reason,
        )


class WalletRevocationHook:
    def __init__(self, notifier: RevocationNotifier, wallet_did: str):
        self._notifier = notifier
        self._wallet_did = wallet_did
        self._notifications: list[RevocationEvent] = []

        self._notifier.subscribe(
            RevocationEventType.CREDENTIAL_REVOKED,
            self._on_event,
        )
        self._notifier.subscribe(
            RevocationEventType.CREDENTIAL_SUSPENDED,
            self._on_event,
        )
        self._notifier.subscribe(
            RevocationEventType.CREDENTIAL_REINSTATED,
            self._on_event,
        )

    def _on_event(self, event: RevocationEvent):
        tracked = self._notifier.get_wallet_credentials(self._wallet_did)
        if event.credential_id in tracked:
            self._notifications.append(event)

    def track_credential(self, credential_id: str):
        self._notifier.register_wallet(self._wallet_did, credential_id)

    def untrack_credential(self, credential_id: str):
        self._notifier.unregister_wallet(self._wallet_did, credential_id)

    def check_status(self, credential_id: str) -> RevocationEventType | None:
        for event in reversed(self._notifications):
            if event.credential_id == credential_id:
                return event.event_type
        return None

    def get_notifications(self) -> list[RevocationEvent]:
        return list(self._notifications)

    def clear_notifications(self):
        self._notifications.clear()
