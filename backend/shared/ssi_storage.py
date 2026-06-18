"""
Async SQLAlchemy storage backends for SSI modules.

Each backend wraps CRUD operations against the ``ssi_models`` tables.
Designed to be injected into SSI services as a drop-in for in-memory stores.
"""
from __future__ import annotations

from datetime import datetime, timezone
from typing import Any, Sequence

from sqlalchemy import select, func, delete, or_, and_
from sqlalchemy.ext.asyncio import AsyncSession

from shared.database.ssi_models import (
    CHAPIRecord,
    CredentialFlowRecord,
    CredentialManifestRecord,
    DIDCommMessageRecord,
    DIDRecord,
    RevocationEventRecord,
    StatusListRecord,
    VerifiableCredentialRecord,
    WalletCredentialRecord,
)


class DIDStorage:
    def __init__(self, db: AsyncSession) -> None:
        self._db = db

    async def save(self, did: str, method: str, document: dict) -> None:
        self._db.add(DIDRecord(did=did, method=method, document=document))
        await self._db.flush()

    async def get(self, did: str) -> dict | None:
        row = await self._db.scalar(select(DIDRecord).where(DIDRecord.did == did))
        return row.document if row else None

    async def list_by_method(self, method: str) -> list[dict]:
        rows = await self._db.scalars(select(DIDRecord).where(DIDRecord.method == method))
        return [r.document for r in rows.all()]

    async def delete(self, did: str) -> bool:
        row = await self._db.scalar(select(DIDRecord).where(DIDRecord.did == did))
        if row is None:
            return False
        await self._db.delete(row)
        return True


class VCStorage:
    def __init__(self, db: AsyncSession) -> None:
        self._db = db

    async def save(self, credential_id: str, issuer_id: str, subject_id: str,
                    credential_type: str, document: dict, status_list_id: str | None = None) -> None:
        self._db.add(VerifiableCredentialRecord(
            credential_id=credential_id, issuer_id=issuer_id, subject_id=subject_id,
            credential_type=credential_type, document=document,
            status_list_id=status_list_id, revoked=False, issued_at=datetime.now(timezone.utc),
        ))
        await self._db.flush()

    async def get(self, credential_id: str) -> dict | None:
        row = await self._db.scalar(
            select(VerifiableCredentialRecord).where(VerifiableCredentialRecord.credential_id == credential_id)
        )
        return row.document if row else None

    async def list_by_subject(self, subject_id: str) -> list[dict]:
        rows = await self._db.scalars(
            select(VerifiableCredentialRecord).where(VerifiableCredentialRecord.subject_id == subject_id)
        )
        return [r.document for r in rows.all()]


class StatusListStorage:
    def __init__(self, db: AsyncSession) -> None:
        self._db = db

    async def save(self, list_id: str, purpose: str, bitstring: str) -> None:
        now = datetime.now(timezone.utc)
        self._db.add(StatusListRecord(list_id=list_id, purpose=purpose, bitstring=bitstring, created_at=now, updated_at=now))
        await self._db.flush()

    async def get(self, list_id: str) -> StatusListRecord | None:
        return await self._db.scalar(select(StatusListRecord).where(StatusListRecord.list_id == list_id))

    async def update_bitstring(self, list_id: str, bitstring: str) -> bool:
        row = await self._db.scalar(select(StatusListRecord).where(StatusListRecord.list_id == list_id))
        if row is None:
            return False
        row.bitstring = bitstring
        row.updated_at = datetime.now(timezone.utc)
        return True

    async def delete(self, list_id: str) -> bool:
        row = await self._db.scalar(select(StatusListRecord).where(StatusListRecord.list_id == list_id))
        if row is None:
            return False
        await self._db.delete(row)
        return True


class WalletCredentialStorage:
    def __init__(self, db: AsyncSession) -> None:
        self._db = db

    async def save(self, wallet_did: str, credential_id: str, document: dict,
                    issuer_id: str, credential_type: str, issued_at: datetime | None = None) -> None:
        self._db.add(WalletCredentialRecord(
            wallet_did=wallet_did, credential_id=credential_id, document=document,
            issuer_id=issuer_id, credential_type=credential_type,
            issued_at=issued_at or datetime.now(timezone.utc),
            stored_at=datetime.now(timezone.utc),
        ))
        await self._db.flush()

    async def get(self, wallet_did: str, credential_id: str) -> dict | None:
        row = await self._db.scalar(
            select(WalletCredentialRecord).where(
                WalletCredentialRecord.wallet_did == wallet_did,
                WalletCredentialRecord.credential_id == credential_id,
            )
        )
        return row.document if row else None

    async def list_by_wallet(self, wallet_did: str) -> list[dict]:
        rows = await self._db.scalars(
            select(WalletCredentialRecord).where(WalletCredentialRecord.wallet_did == wallet_did)
        )
        return [r.document for r in rows.all()]

    async def delete(self, wallet_did: str, credential_id: str) -> bool:
        row = await self._db.scalar(
            select(WalletCredentialRecord).where(
                WalletCredentialRecord.wallet_did == wallet_did,
                WalletCredentialRecord.credential_id == credential_id,
            )
        )
        if row is None:
            return False
        await self._db.delete(row)
        return True

    async def search(self, wallet_did: str, query: str) -> list[dict]:
        stmt = select(WalletCredentialRecord).where(
            WalletCredentialRecord.wallet_did == wallet_did,
            or_(
                WalletCredentialRecord.credential_type.ilike(f"%{query}%"),
                WalletCredentialRecord.issuer_id.ilike(f"%{query}%"),
                WalletCredentialRecord.credential_id.ilike(f"%{query}%"),
            ),
        )
        rows = await self._db.scalars(stmt)
        return [r.document for r in rows.all()]


class DIDCommMessageStorage:
    def __init__(self, db: AsyncSession) -> None:
        self._db = db

    async def save(self, message_id: str, sender_did: str, receiver_did: str,
                    message_type: str, message_body: dict, thread_id: str | None = None) -> None:
        self._db.add(DIDCommMessageRecord(
            message_id=message_id, sender_did=sender_did, receiver_did=receiver_did,
            message_type=message_type, message_body=message_body,
            thread_id=thread_id, created_at=datetime.now(timezone.utc),
        ))
        await self._db.flush()

    async def get_inbox(self, receiver_did: str) -> list[DIDCommMessageRecord]:
        rows = await self._db.scalars(
            select(DIDCommMessageRecord)
            .where(DIDCommMessageRecord.receiver_did == receiver_did)
            .order_by(DIDCommMessageRecord.created_at.desc())
        )
        return rows.all()

    async def get_sent(self, sender_did: str) -> list[DIDCommMessageRecord]:
        rows = await self._db.scalars(
            select(DIDCommMessageRecord)
            .where(DIDCommMessageRecord.sender_did == sender_did)
            .order_by(DIDCommMessageRecord.created_at.desc())
        )
        return rows.all()

    async def mark_read(self, message_id: str) -> bool:
        row = await self._db.scalar(
            select(DIDCommMessageRecord).where(DIDCommMessageRecord.message_id == message_id)
        )
        if row is None:
            return False
        row.is_read = True
        return True

    async def delete(self, message_id: str) -> bool:
        row = await self._db.scalar(
            select(DIDCommMessageRecord).where(DIDCommMessageRecord.message_id == message_id)
        )
        if row is None:
            return False
        await self._db.delete(row)
        return True


class CredentialFlowStorage:
    def __init__(self, db: AsyncSession) -> None:
        self._db = db

    async def save(self, flow_id: str, issuer_id: str, offer_data: dict, status: str) -> None:
        now = datetime.now(timezone.utc)
        self._db.add(CredentialFlowRecord(
            flow_id=flow_id, issuer_id=issuer_id, offer_data=offer_data,
            status=status, created_at=now, updated_at=now,
        ))
        await self._db.flush()

    async def get(self, flow_id: str) -> CredentialFlowRecord | None:
        return await self._db.scalar(select(CredentialFlowRecord).where(CredentialFlowRecord.flow_id == flow_id))

    async def update_request(self, flow_id: str, request_data: dict) -> bool:
        row = await self._db.scalar(select(CredentialFlowRecord).where(CredentialFlowRecord.flow_id == flow_id))
        if row is None:
            return False
        row.request_data = request_data
        row.status = "requested"
        row.updated_at = datetime.now(timezone.utc)
        return True

    async def update_issued(self, flow_id: str, credential_id: str) -> bool:
        row = await self._db.scalar(select(CredentialFlowRecord).where(CredentialFlowRecord.flow_id == flow_id))
        if row is None:
            return False
        row.credential_id = credential_id
        row.status = "issued"
        row.updated_at = datetime.now(timezone.utc)
        return True


class CHAPIStorage:
    def __init__(self, db: AsyncSession) -> None:
        self._db = db

    async def save(self, holder_did: str, credential_id: str, document: dict, query_frame: dict | None = None) -> None:
        self._db.add(CHAPIRecord(
            holder_did=holder_did, credential_id=credential_id, document=document,
            query_frame=query_frame, stored_at=datetime.now(timezone.utc),
        ))
        await self._db.flush()

    async def get(self, holder_did: str, credential_id: str) -> dict | None:
        row = await self._db.scalar(
            select(CHAPIRecord).where(
                CHAPIRecord.holder_did == holder_did,
                CHAPIRecord.credential_id == credential_id,
            )
        )
        return row.document if row else None

    async def list_by_holder(self, holder_did: str) -> list[dict]:
        rows = await self._db.scalars(
            select(CHAPIRecord).where(CHAPIRecord.holder_did == holder_did)
        )
        return [r.document for r in rows.all()]

    async def delete(self, holder_did: str, credential_id: str) -> bool:
        row = await self._db.scalar(
            select(CHAPIRecord).where(
                CHAPIRecord.holder_did == holder_did,
                CHAPIRecord.credential_id == credential_id,
            )
        )
        if row is None:
            return False
        await self._db.delete(row)
        return True


class CredentialManifestStorage:
    def __init__(self, db: AsyncSession) -> None:
        self._db = db

    async def save(self, manifest_id: str, issuer_id: str, document: dict) -> None:
        self._db.add(CredentialManifestRecord(
            manifest_id=manifest_id, issuer_id=issuer_id, document=document,
            created_at=datetime.now(timezone.utc),
        ))
        await self._db.flush()

    async def get(self, manifest_id: str) -> dict | None:
        row = await self._db.scalar(
            select(CredentialManifestRecord).where(CredentialManifestRecord.manifest_id == manifest_id)
        )
        return row.document if row else None

    async def list_by_issuer(self, issuer_id: str) -> list[dict]:
        rows = await self._db.scalars(
            select(CredentialManifestRecord).where(CredentialManifestRecord.issuer_id == issuer_id)
        )
        return [r.document for r in rows.all()]


class RevocationEventStorage:
    def __init__(self, db: AsyncSession) -> None:
        self._db = db

    async def save(self, event_id: str, credential_id: str, subject_id: str,
                    event_type: str, reason: str | None = None) -> None:
        self._db.add(RevocationEventRecord(
            event_id=event_id, credential_id=credential_id, subject_id=subject_id,
            event_type=event_type, reason=reason, created_at=datetime.now(timezone.utc),
        ))
        await self._db.flush()

    async def get_history(self, credential_id: str | None = None,
                           subject_id: str | None = None, limit: int = 50) -> list[RevocationEventRecord]:
        stmt = select(RevocationEventRecord).order_by(RevocationEventRecord.created_at.desc())
        if credential_id:
            stmt = stmt.where(RevocationEventRecord.credential_id == credential_id)
        if subject_id:
            stmt = stmt.where(RevocationEventRecord.subject_id == subject_id)
        rows = await self._db.scalars(stmt.limit(limit))
        return rows.all()
