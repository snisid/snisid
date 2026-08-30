from __future__ import annotations

import uuid
from datetime import datetime, timezone
from typing import Any

from services.wallet import Wallet
from services.vp import VerifiablePresentation


class CHAPIStoreRequest:
    def __init__(
        self,
        credential: dict[str, Any],
        credential_id: str | None = None,
        protocol: str = "vc",
    ):
        self.credential = credential
        self.credential_id = credential_id or str(uuid.uuid4())
        self.protocol = protocol
        self.timestamp = datetime.now(timezone.utc).isoformat()

    def to_dict(self) -> dict[str, Any]:
        return {
            "credential_id": self.credential_id,
            "credential": self.credential,
            "protocol": self.protocol,
            "timestamp": self.timestamp,
        }


class CHAPIGetRequest:
    def __init__(
        self,
        query: list[dict[str, Any]],
        credential_id: str | None = None,
        protocol: str = "vc",
    ):
        self.query = query
        self.credential_id = credential_id or str(uuid.uuid4())
        self.protocol = protocol
        self.timestamp = datetime.now(timezone.utc).isoformat()

    def to_dict(self) -> dict[str, Any]:
        return {
            "credential_id": self.credential_id,
            "query": self.query,
            "protocol": self.protocol,
            "timestamp": self.timestamp,
        }


class CHAPIResponse:
    def __init__(
        self,
        data: list[dict[str, Any]] | dict[str, Any] | None = None,
        error: str | None = None,
    ):
        self.data = data
        self.error = error
        self.timestamp = datetime.now(timezone.utc).isoformat()

    def to_dict(self) -> dict[str, Any]:
        d: dict[str, Any] = {"timestamp": self.timestamp}
        if self.data is not None:
            d["data"] = self.data
        if self.error is not None:
            d["error"] = self.error
        return d


class CHAPIMediator:
    def __init__(self, wallet: Wallet | None = None):
        self._wallet = wallet or Wallet()

    @property
    def wallet(self) -> Wallet:
        return self._wallet

    def handle_store(self, request: CHAPIStoreRequest) -> CHAPIResponse:
        try:
            record = self._wallet.store(request.credential)
            return CHAPIResponse(
                data={
                    "status": "stored",
                    "credential_id": record.id,
                }
            )
        except Exception as e:
            return CHAPIResponse(error=str(e))

    def handle_get(self, request: CHAPIGetRequest) -> CHAPIResponse:
        try:
            all_records = {r.id: r for r in self._wallet.list()}
            matched_ids: set[str] = set()

            for q in request.query:
                qtype = q.get("type", "VerifiableCredential")
                qframe = q.get("credentialFrame", {})

                for rid, record in all_records.items():
                    if qtype in record.credential_type and self._matches_frame(record.credential, qframe):
                        matched_ids.add(rid)

            if matched_ids:
                vp = self._wallet.create_presentation(credential_ids=list(matched_ids))
                return CHAPIResponse(data=vp.to_dict())
            return CHAPIResponse(data={"verifiableCredential": []})

        except Exception as e:
            return CHAPIResponse(error=str(e))

    def _matches_frame(self, credential: dict[str, Any], frame: dict[str, Any]) -> bool:
        if not frame:
            return True
        for key, expected in frame.items():
            actual = credential.get(key)
            if isinstance(expected, dict) and isinstance(actual, dict):
                if not self._matches_frame(actual, expected):
                    return False
            elif isinstance(expected, list) and isinstance(actual, list):
                if not any(a in expected for a in actual if isinstance(a, str)):
                    return False
            elif actual != expected:
                return False
        return True

    def handler_registration(self, origin: str) -> CHAPIResponse:
        return CHAPIResponse(
            data={
                "handler": "snisid-wallet",
                "origin": origin,
                "capabilities": ["store", "get"],
                "version": "1.0",
            }
        )
