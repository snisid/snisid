"""
DIDComm Mediator - Routing and Forwarding Service.

Implements DIDComm v2 routing mediator pattern:
- Forward messages to/from recipients via mediator
- Store encrypted messages for pickup
- Provide inbox/sent/fetch semantics
- Supports forward message type per DIDComm v2 spec
"""
from __future__ import annotations

import uuid
from dataclasses import dataclass, field
from datetime import datetime, timezone
from typing import Any


@dataclass
class ForwardedMessage:
    """A forwarded DIDComm message stored by the mediator."""

    id: str = field(default_factory=lambda: str(uuid.uuid4()))
    recipient_did: str = ""
    sender_did: str = ""
    packed_message: dict[str, Any] | None = None
    forward_request: dict[str, Any] | None = None
    created_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))
    is_delivered: bool = False
    delivered_at: datetime | None = None


class DIDCommMediator:
    """
    DIDComm v2 message mediator.

    Receives forward requests and stores them for the recipient to pick up.
    Supports both packed and unpacked forwarding.
    """

    def __init__(self) -> None:
        self._messages: dict[str, ForwardedMessage] = {}
        self._locks: list[ForwardedMessage] = []

    def forward(self, recipient_did: str, packed_message: dict[str, Any], sender_did: str = "") -> ForwardedMessage:
        """
        Accept a Forward message and store it for the recipient.

        Args:
            recipient_did: DID of the intended recipient.
            packed_message: The encrypted DIDComm message (JWM).
            sender_did: DID of the sender (optional).

        Returns:
            The stored ForwardedMessage.
        """
        msg = ForwardedMessage(
            recipient_did=recipient_did,
            sender_did=sender_did,
            packed_message=packed_message,
            is_delivered=False,
        )
        self._messages[msg.id] = msg
        return msg

    def forward_request(self, forward_msg: dict[str, Any]) -> ForwardedMessage | None:
        """
        Process a DIDComm Forward message per the spec.

        Expected format (per DIDComm v2):
        ```json
        {
            "type": "https://didcomm.org/routing/2.0/forward",
            "id": "<uuid>",
            "to": ["<recipient_did>"],
            "body": {
                "next": "<recipient_did>"
            },
            "attachments": [{
                "data": { "json": <packed_message> }
            }]
        }
        ```
        """
        if not isinstance(forward_msg, dict):
            return None
        msg_type = forward_msg.get("type", "")
        if "forward" not in msg_type.lower():
            return None
        to = forward_msg.get("to", [])
        if not to:
            return None
        recipient_did = to[0] if isinstance(to, list) else to
        attachments = forward_msg.get("attachments", [])
        packed = None
        if attachments:
            data = attachments[0].get("data", {})
            packed = data.get("json") or data.get("base64") or data
        sender_did = forward_msg.get("from", "")
        return self.forward(recipient_did, packed or forward_msg, sender_did=sender_did)

    def fetch_messages(self, recipient_did: str) -> list[ForwardedMessage]:
        """Return all undelivered messages for a recipient."""
        pending = []
        for msg in self._messages.values():
            if msg.recipient_did == recipient_did and not msg.is_delivered:
                pending.append(msg)
        return pending

    def deliver(self, message_id: str) -> ForwardedMessage | None:
        """Mark a message as delivered and return it."""
        msg = self._messages.get(message_id)
        if msg is None:
            return None
        msg.is_delivered = True
        msg.delivered_at = datetime.now(timezone.utc)
        return msg

    def get_inbox(self, recipient_did: str) -> list[ForwardedMessage]:
        """Return all messages (delivered and pending) for a recipient, newest first."""
        msgs = [m for m in self._messages.values() if m.recipient_did == recipient_did]
        msgs.sort(key=lambda m: m.created_at, reverse=True)
        return msgs

    def get_pending_count(self, recipient_did: str) -> int:
        """Return count of undelivered messages."""
        return sum(1 for m in self._messages.values() if m.recipient_did == recipient_did and not m.is_delivered)

    def delete_message(self, message_id: str) -> bool:
        """Remove a message from the mediator."""
        if message_id in self._messages:
            del self._messages[message_id]
            return True
        return False

    def clear(self) -> None:
        """Remove all messages."""
        self._messages.clear()
