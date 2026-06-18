from __future__ import annotations

import uuid
from datetime import datetime, timezone
from typing import Any

from services.did import create_did_key, resolve_did
from services.didcomm import DIDCommMessage, DIDCommMessenger
from services.vp import VPIssuer, VerifiablePresentation


class CredentialRecord:
    def __init__(
        self,
        id: str,
        credential: dict[str, Any],
        issuer_did: str,
        credential_type: list[str],
        issued_at: str,
        label: str = "",
    ):
        self.id = id
        self.credential = credential
        self.issuer_did = issuer_did
        self.credential_type = credential_type
        self.issued_at = issued_at
        self.label = label

    def to_dict(self) -> dict[str, Any]:
        return {
            "id": self.id,
            "credential": self.credential,
            "issuer_did": self.issuer_did,
            "credential_type": self.credential_type,
            "issued_at": self.issued_at,
            "label": self.label,
        }


class Wallet:
    def __init__(self, wallet_key: str = "dev-wallet-key", storage: Any | None = None):
        self._key = wallet_key
        self._did, self._private_key, self._public_key = create_did_key()
        self._vp_issuer = VPIssuer(proof_key=wallet_key)
        self._credentials: dict[str, CredentialRecord] = {}
        self._storage = storage

    @property
    def did(self) -> str:
        return self._did

    @property
    def did_document(self) -> dict[str, Any]:
        doc = resolve_did(self._did)
        return doc.model_dump()

    def store(self, credential: dict[str, Any], label: str = "") -> CredentialRecord:
        record = CredentialRecord(
            id=str(uuid.uuid4()),
            credential=credential,
            issuer_did=credential.get("issuer", ""),
            credential_type=credential.get("type", ["VerifiableCredential"]),
            issued_at=credential.get("issuanceDate", datetime.now(timezone.utc).isoformat()),
            label=label,
        )
        self._credentials[record.id] = record
        return record

    async def async_store(self, credential: dict[str, Any], label: str = "") -> CredentialRecord:
        record = self.store(credential, label)
        if self._storage:
            await self._storage.save(
                wallet_did=self._did,
                credential_id=record.id,
                document=record.credential,
                issuer_id=record.issuer_did,
                credential_type=record.credential_type[0] if record.credential_type else "VerifiableCredential",
            )
        return record

    def get(self, credential_id: str) -> CredentialRecord | None:
        return self._credentials.get(credential_id)

    async def async_get(self, credential_id: str) -> CredentialRecord | None:
        if self._storage:
            doc = await self._storage.get(self._did, credential_id)
            if doc:
                return CredentialRecord(
                    id=credential_id,
                    credential=doc,
                    issuer_did=doc.get("issuer", ""),
                    credential_type=doc.get("type", ["VerifiableCredential"]),
                    issued_at=doc.get("issuanceDate", ""),
                )
        return self.get(credential_id)

    def list(self, credential_type: str | None = None) -> list[CredentialRecord]:
        records = list(self._credentials.values())
        if credential_type:
            records = [r for r in records if credential_type in r.credential_type]
        return records

    async def async_list(self, credential_type: str | None = None) -> list[CredentialRecord]:
        if self._storage:
            docs = await self._storage.list_by_wallet(self._did)
            records = [
                CredentialRecord(
                    id=doc.get("id", str(uuid.uuid4())),
                    credential=doc,
                    issuer_did=doc.get("issuer", ""),
                    credential_type=doc.get("type", ["VerifiableCredential"]),
                    issued_at=doc.get("issuanceDate", ""),
                )
                for doc in docs
            ]
            if credential_type:
                records = [r for r in records if credential_type in r.credential_type]
            return records
        return self.list(credential_type)

    def delete(self, credential_id: str) -> bool:
        if credential_id in self._credentials:
            del self._credentials[credential_id]
            return True
        return False

    async def async_delete(self, credential_id: str) -> bool:
        if self._storage:
            ok = await self._storage.delete(self._did, credential_id)
            self._credentials.pop(credential_id, None)
            return ok
        return self.delete(credential_id)

    def count(self) -> int:
        return len(self._credentials)

    def create_presentation(
        self,
        credential_ids: list[str] | None = None,
        credential_type: str | None = None,
    ) -> VerifiablePresentation:
        if credential_ids:
            records = [self._credentials[cid] for cid in credential_ids if cid in self._credentials]
        elif credential_type:
            records = self.list(credential_type)
        else:
            records = self.list()

        vcs = [r.credential for r in records]
        return self._vp_issuer.create_presentation(holder_did=self._did, verifiable_credentials=vcs)

    def verify_presentation(self, vp: VerifiablePresentation) -> bool:
        return self._vp_issuer.verify_presentation(vp)

    def send_via_didcomm(
        self,
        credential_id: str,
        to_did: str,
        message_type: str = "https://didcomm.org/issue-credential/2.0/issue-credential",
    ) -> dict[str, Any]:
        record = self._credentials.get(credential_id)
        if not record:
            raise ValueError(f"Credential not found: {credential_id}")
        messenger = DIDCommMessenger()
        msg = DIDCommMessage(
            id=str(uuid.uuid4()),
            type=message_type,
            body={"credentials": [record.credential]},
            from_did=self._did,
            to_did=to_did,
        )
        return messenger.send(msg, self._did, to_did)

    def receive_via_didcomm(self, packed: dict[str, Any], label: str = "") -> CredentialRecord | None:
        messenger = DIDCommMessenger()
        msg = messenger.receive(packed)
        credentials = msg.body.get("credentials", [])
        if not credentials:
            return None
        stored = None
        for cred in credentials:
            stored = self.store(cred, label=label)
        return stored

    def send_message(
        self,
        message_type: str,
        body: dict[str, Any],
        to_did: str,
    ) -> dict[str, Any]:
        messenger = DIDCommMessenger()
        msg = DIDCommMessage(
            id=str(uuid.uuid4()),
            type=message_type,
            body=body,
            from_did=self._did,
            to_did=to_did,
        )
        return messenger.send(msg, self._did, to_did)

    def search(self, query: str) -> list[CredentialRecord]:
        q = query.lower()
        results = []
        for record in self._credentials.values():
            if q in record.label.lower():
                results.append(record)
                continue
            sub = record.credential.get("credentialSubject", {})
            if isinstance(sub, dict):
                if any(q in str(v).lower() for v in sub.values()):
                    results.append(record)
                    continue
            if q in record.issuer_did.lower():
                results.append(record)
                continue
        return results

    async def async_search(self, query: str) -> list[CredentialRecord]:
        if self._storage:
            docs = await self._storage.search(self._did, query)
            return [
                CredentialRecord(
                    id=doc.get("id", str(uuid.uuid4())),
                    credential=doc,
                    issuer_did=doc.get("issuer", ""),
                    credential_type=doc.get("type", ["VerifiableCredential"]),
                    issued_at=doc.get("issuanceDate", ""),
                )
                for doc in docs
            ]
        return self.search(query)

    def get_credentials_by_issuer(self, issuer_did: str) -> list[CredentialRecord]:
        return [r for r in self._credentials.values() if r.issuer_did == issuer_did]

    def export_all(self) -> list[dict[str, Any]]:
        return [r.to_dict() for r in self._credentials.values()]

    def import_credentials(self, records: list[dict[str, Any]]):
        for r in records:
            record = CredentialRecord(
                id=r.get("id", str(uuid.uuid4())),
                credential=r.get("credential", {}),
                issuer_did=r.get("issuer_did", ""),
                credential_type=r.get("credential_type", ["VerifiableCredential"]),
                issued_at=r.get("issued_at", ""),
                label=r.get("label", ""),
            )
            self._credentials[record.id] = record
