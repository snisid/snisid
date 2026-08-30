from __future__ import annotations

import json
from datetime import datetime, timezone
from typing import Any

from services.did import create_did_key, resolve_did


def _create_ephemeral_key() -> tuple[str, str, str]:
    return create_did_key()


class DIDCommMessage:
    def __init__(
        self,
        id: str,
        type: str,
        body: dict[str, Any],
        from_did: str | None = None,
        to_did: str | None = None,
        thid: str | None = None,
        pthid: str | None = None,
        created_time: int | None = None,
        attachments: list[dict] | None = None,
    ):
        self.id = id
        self.type = type
        self.body = body
        self.from_did = from_did
        self.to_did = to_did
        self.thid = thid
        self.pthid = pthid
        self.created_time = created_time or int(datetime.now(timezone.utc).timestamp())
        self.attachments = attachments or []

    def to_dict(self) -> dict:
        msg = {
            "id": self.id,
            "type": self.type,
            "body": self.body,
            "created_time": self.created_time,
        }
        if self.from_did:
            msg["from"] = self.from_did
        if self.to_did:
            msg["to"] = [self.to_did] if isinstance(self.to_did, str) else self.to_did
        if self.thid:
            msg["thid"] = self.thid
        if self.pthid:
            msg["pthid"] = self.pthid
        if self.attachments:
            msg["attachments"] = self.attachments
        return msg

    @classmethod
    def from_dict(cls, data: dict) -> DIDCommMessage:
        to = data.get("to", [])
        return cls(
            id=data["id"],
            type=data["type"],
            body=data.get("body", {}),
            from_did=data.get("from"),
            to_did=to[0] if isinstance(to, list) and to else to,
            thid=data.get("thid"),
            pthid=data.get("pthid"),
            created_time=data.get("created_time"),
            attachments=data.get("attachments"),
        )


class DIDCommPacker:
    def __init__(self, signing_key: str = "dev-didcomm-key"):
        self._key = signing_key

    def pack(
        self,
        message: DIDCommMessage,
        sender_did: str | None = None,
        recipient_did: str | None = None,
    ) -> dict:
        signed = self._sign(message, sender_did)
        if recipient_did:
            return self._encrypt(signed, recipient_did)
        return signed

    def unpack(self, packed: dict) -> DIDCommMessage:
        data = packed

        if "ciphertext" in data:
            data = self._decrypt(data)

        if "signatures" in data:
            data = self._verify(data)

        return DIDCommMessage.from_dict(data)

    def _sign(self, message: DIDCommMessage, sender_did: str | None) -> dict:
        msg_dict = message.to_dict()
        payload = json.dumps(msg_dict, separators=(",", ":"), sort_keys=True)
        signature = __import__("hashlib").sha256(
            f"{payload}{self._key}".encode()
        ).hexdigest()

        signed = {
            "payload": __import__("base64").urlsafe_b64encode(
                payload.encode()
            ).rstrip(b"=").decode(),
            "signatures": [
                {
                    "header": {"kid": f"{sender_did}#key-1" if sender_did else "did:key:default#key-1"},
                    "signature": signature,
                    "protected": json.dumps({"alg": "SNISID-HMAC-SHA256-2026"}, separators=(",", ":")),
                }
            ],
        }
        return signed

    def _encrypt(self, signed: dict, recipient_did: str) -> dict:
        payload = json.dumps(signed, separators=(",", ":"), sort_keys=True)
        _, _, public = _create_ephemeral_key()

        encrypted = {
            "ciphertext": __import__("base64").urlsafe_b64encode(
                payload.encode()
            ).rstrip(b"=").decode(),
            "protected": json.dumps(
                {"typ": "didcomm-enveloped", "alg": "A256GCM"},
                separators=(",", ":"),
            ),
            "recipients": [
                {
                    "header": {"kid": f"{recipient_did}#key-1"},
                }
            ],
        }
        return encrypted

    def _decrypt(self, data: dict) -> dict:
        raw = __import__("base64").urlsafe_b64decode(
            data["ciphertext"] + "=="
        )
        return json.loads(raw)

    def _verify(self, data: dict) -> dict:
        raw = __import__("base64").urlsafe_b64decode(
            data["payload"] + "=="
        )
        return json.loads(raw)


class DIDCommMessenger:
    def __init__(self, packer: DIDCommPacker | None = None, storage: Any | None = None):
        self._packer = packer or DIDCommPacker()
        self._inbox: list[dict] = []
        self._sent: list[dict] = []
        self._storage = storage

    def send(self, message: DIDCommMessage, sender_did: str, recipient_did: str) -> dict:
        packed = self._packer.pack(
            message=message,
            sender_did=sender_did,
            recipient_did=recipient_did,
        )
        self._sent.append(packed)
        return packed

    async def async_send(self, message: DIDCommMessage, sender_did: str, recipient_did: str) -> dict:
        packed = self.send(message, sender_did, recipient_did)
        if self._storage:
            await self._storage.save(
                message_id=message.id,
                sender_did=sender_did,
                receiver_did=recipient_did,
                message_type=message.type,
                message_body=message.body,
                thread_id=message.thid,
            )
        return packed

    def receive(self, packed: dict) -> DIDCommMessage:
        message = self._packer.unpack(packed)
        self._inbox.append(packed)
        return message

    async def async_receive(self, packed: dict) -> DIDCommMessage:
        message = self.receive(packed)
        if self._storage:
            existing = await self._storage.get_inbox(message.to_did or "")
            if not any(r.message_id == message.id for r in existing):
                await self._storage.save(
                    message_id=message.id,
                    sender_did=message.from_did or "",
                    receiver_did=message.to_did or "",
                    message_type=message.type,
                    message_body=message.body,
                    thread_id=message.thid,
                )
        return message

    async def async_get_inbox(self, receiver_did: str) -> list[DIDCommMessage]:
        if self._storage:
            records = await self._storage.get_inbox(receiver_did)
            return [DIDCommMessage.from_dict({
                "id": r.message_id,
                "type": r.message_type,
                "body": r.message_body,
                "from": r.sender_did,
                "to": [r.receiver_did],
                "thid": r.thread_id,
                "created_time": int(r.created_at.timestamp()) if r.created_at else 0,
            }) for r in records]
        return [DIDCommMessage.from_dict(p) for p in self._inbox]

    async def async_get_sent(self, sender_did: str) -> list[DIDCommMessage]:
        if self._storage:
            records = await self._storage.get_sent(sender_did)
            return [DIDCommMessage.from_dict({
                "id": r.message_id,
                "type": r.message_type,
                "body": r.message_body,
                "from": r.sender_did,
                "to": [r.receiver_did],
                "thid": r.thread_id,
                "created_time": int(r.created_at.timestamp()) if r.created_at else 0,
            }) for r in records]
        return [DIDCommMessage.from_dict(p) for p in self._sent]

    def create_trust_ping(
        self, from_did: str, to_did: str, response_requested: bool = True
    ) -> DIDCommMessage:
        import uuid
        return DIDCommMessage(
            id=str(uuid.uuid4()),
            type="https://didcomm.org/trust-ping/2.0/ping",
            body={"response_requested": response_requested},
            from_did=from_did,
            to_did=to_did,
        )

    def create_trust_ping_response(self, ping: DIDCommMessage) -> DIDCommMessage:
        import uuid
        return DIDCommMessage(
            id=str(uuid.uuid4()),
            type="https://didcomm.org/trust-ping/2.0/ping_response",
            body={},
            from_did=ping.to_did,
            to_did=ping.from_did,
            thid=ping.id,
        )
