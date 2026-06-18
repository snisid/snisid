from __future__ import annotations

import base64
import gzip
import json
import uuid
from datetime import datetime, timezone
from typing import Any


BITS_PER_CREDENTIAL = 1
BITS_PER_BYTE = 8


class StatusList:
    """W3C StatusList2021 — bitstring-based credential status registry."""

    def __init__(self, size: int = 131072):
        self._size = size
        self._bits = bytearray((size + 7) // 8)

    @property
    def size(self) -> int:
        return self._size

    def set_status(self, index: int, revoked: bool = True) -> None:
        if index < 0 or index >= self._size:
            raise ValueError(f"Index {index} out of range [0, {self._size})")
        byte_idx = index // 8
        bit_idx = index % 8
        if revoked:
            self._bits[byte_idx] |= 1 << bit_idx
        else:
            self._bits[byte_idx] &= ~(1 << bit_idx)

    def get_status(self, index: int) -> bool:
        if index < 0 or index >= self._size:
            raise ValueError(f"Index {index} out of range [0, {self._size})")
        byte_idx = index // 8
        bit_idx = index % 8
        return bool(self._bits[byte_idx] & (1 << bit_idx))

    def encode(self) -> str:
        compressed = gzip.compress(bytes(self._bits))
        return base64.urlsafe_b64encode(compressed).rstrip(b"=").decode()

    @classmethod
    def decode(cls, encoded: str, size: int | None = None) -> StatusList:
        padding = 4 - len(encoded) % 4
        if padding != 4:
            encoded += "=" * padding
        compressed = base64.urlsafe_b64decode(encoded)
        bits = bytearray(gzip.decompress(compressed))
        if size is None:
            size = len(bits) * 8
        sl = cls.__new__(cls)
        sl._size = size
        sl._bits = bits
        return sl

    def to_dict(self) -> dict[str, Any]:
        return {"size": self._size, "encodedList": self.encode()}


class StatusList2021Credential:
    """A StatusList2021 Credential wrapping an encoded bitstring."""

    def __init__(
        self,
        issuer_id: str,
        status_list: StatusList,
        purpose: str = "revocation",
        list_id: str | None = None,
    ):
        self.id = list_id or f"{issuer_id}/credentials/status/1"
        self.issuer_id = issuer_id
        self.purpose = purpose
        self.status_list = status_list
        self._proof: dict[str, Any] | None = None

    def set_proof(self, proof: dict[str, Any]) -> None:
        self._proof = proof

    def to_vc_dict(self) -> dict[str, Any]:
        now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        vc = {
            "@context": [
                "https://www.w3.org/ns/credentials/v2",
                "https://www.w3.org/ns/credentials/status/v2",
            ],
            "id": self.id,
            "type": ["VerifiableCredential", "StatusList2021Credential"],
            "issuer": self.issuer_id,
            "issuanceDate": now,
            "credentialSubject": {
                "id": f"{self.id}#list",
                "type": "StatusList2021",
                "statusPurpose": self.purpose,
                "encodedList": self.status_list.encode(),
            },
        }
        if self._proof:
            vc["proof"] = self._proof
        return vc


class StatusListEntry:
    """A credentialStatus entry pointing to a StatusList2021 position."""

    def __init__(
        self,
        status_list_credential_id: str,
        index: int,
        purpose: str = "revocation",
    ):
        self.id = f"{status_list_credential_id}#{index}"
        self.type = "StatusList2021Entry"
        self.statusPurpose = purpose
        self.statusListIndex = str(index)
        self.statusListCredential = status_list_credential_id

    def to_dict(self) -> dict[str, Any]:
        return {
            "id": self.id,
            "type": self.type,
            "statusPurpose": self.statusPurpose,
            "statusListIndex": self.statusListIndex,
            "statusListCredential": self.statusListCredential,
        }


class StatusListManager:
    """Manages StatusList2021 credentials and entries."""

    def __init__(self, issuer_id: str, storage: Any | None = None):
        self._issuer_id = issuer_id
        self._lists: list[tuple[str, StatusList]] = []
        self._entries: dict[str, int] = {}
        self._entry_count: int = 0
        self._storage = storage

    def get_or_create_list(
        self, purpose: str = "revocation", size: int = 131072
    ) -> tuple[StatusList, str]:
        for list_id, sl in self._lists:
            return sl, list_id
        list_id = f"{self._issuer_id}/credentials/status/1"
        sl = StatusList(size)
        self._lists.append((list_id, sl))
        return sl, list_id

    def create_entry(self, purpose: str = "revocation") -> StatusListEntry:
        if not self._lists:
            self.get_or_create_list(purpose)
        list_id, sl = self._lists[-1]
        index = self._entry_count
        if index >= sl.size:
            new_list_id = f"{self._issuer_id}/credentials/status/{len(self._lists) + 1}"
            sl = StatusList(sl.size)
            self._lists.append((new_list_id, sl))
            list_id = new_list_id
            index = 0
        entry = StatusListEntry(list_id, index, purpose)
        self._entries[entry.id] = index
        self._entry_count += 1
        return entry

    async def async_create_entry(self, purpose: str = "revocation") -> StatusListEntry:
        entry = self.create_entry(purpose)
        if self._storage:
            sl_found = next((sl for lid, sl in self._lists if lid == entry.statusListCredential), None)
            encoded = sl_found.encode() if sl_found else ""
            await self._storage.save(
                list_id=entry.statusListCredential,
                purpose=purpose,
                bitstring=encoded,
            )
        return entry

    def revoke(self, entry_id: str) -> bool:
        index = self._entries.get(entry_id)
        if index is None:
            return False
        for list_id, sl in self._lists:
            if index < sl.size:
                sl.set_status(index, revoked=True)
                return True
        return False

    async def async_revoke(self, entry_id: str) -> bool:
        ok = self.revoke(entry_id)
        if ok and self._storage:
            list_id = entry_id.rsplit("#", 1)[0] if "#" in entry_id else entry_id
            sl_vc = self.get_status_list_credential(list_id)
            if sl_vc:
                encoded = sl_vc.get("credentialSubject", {}).get("encodedList", "")
                await self._storage.update_bitstring(list_id, encoded)
        return ok

    def unrevoke(self, entry_id: str) -> bool:
        index = self._entries.get(entry_id)
        if index is None:
            return False
        for list_id, sl in self._lists:
            if index < sl.size:
                sl.set_status(index, revoked=False)
                return True
        return False

    def is_revoked(self, entry_id: str) -> bool | None:
        index = self._entries.get(entry_id)
        if index is None:
            return None
        for list_id, sl in self._lists:
            if index < sl.size:
                return sl.get_status(index)
        return None

    async def async_is_revoked(self, entry_id: str) -> bool | None:
        return self.is_revoked(entry_id)

    def get_status_list_credential(self, list_id: str) -> dict[str, Any] | None:
        for lid, sl in self._lists:
            if lid == list_id:
                credential = StatusList2021Credential(self._issuer_id, sl)
                return credential.to_vc_dict()
        return None

    @property
    def status_list(self) -> StatusList | None:
        if self._lists:
            return self._lists[-1][1]
        return None
