"""Integration tests: SSI services backed by SQLAlchemy storage."""
from __future__ import annotations

import pytest

import uuid
from sqlalchemy import select

from shared.database.ssi_models import (
    DIDRecord, VerifiableCredentialRecord, WalletCredentialRecord,
    DIDCommMessageRecord, CredentialFlowRecord, CHAPIRecord,
    CredentialManifestRecord, RevocationEventRecord,
)
from shared.ssi_storage import (
    DIDStorage, VCStorage, StatusListStorage, WalletCredentialStorage,
    DIDCommMessageStorage, CredentialFlowStorage, CHAPIStorage,
    CredentialManifestStorage, RevocationEventStorage,
)
from services.did import DIDManager, DIDMethod
from services.vc import VCStatus
from services.vc.issuer import VCIssuer
from services.wallet import Wallet
from services.didcomm import DIDCommMessage, DIDCommMessenger
from services.status_list import StatusListManager
from services.credential_flow import CredentialFlow
from services.credential_manifest import ManifestManager, OutputDescriptor
from services.revocation import RevocationNotifier, RevocationEventType
from services.chapi import CHAPIMediator, CHAPIStoreRequest, CHAPIGetRequest


class TestDIDWithStorage:
    @pytest.mark.asyncio
    async def test_async_create_and_resolve(self, db_session):
        storage = DIDStorage(db_session)
        mgr = DIDManager(storage=storage)
        doc = await mgr.async_create(DIDMethod.KEY)
        assert doc.id.startswith("did:key:")
        resolved = await mgr.async_resolve(doc.id)
        assert resolved.id == doc.id
        row = await db_session.scalar(select(DIDRecord).where(DIDRecord.did == doc.id))
        assert row is not None

    @pytest.mark.asyncio
    async def test_sync_methods_still_work(self, db_session):
        mgr = DIDManager(storage=DIDStorage(db_session))
        doc = mgr.create(DIDMethod.SNISID)
        assert doc.id.startswith("did:snisid:")
        resolved = mgr.resolve(doc.id)
        assert resolved.id == doc.id


class TestVCWithStorage:
    @pytest.mark.asyncio
    async def test_async_issue_identity_credential(self, db_session):
        vc_storage = VCStorage(db_session)
        status_storage = StatusListStorage(db_session)
        issuer = VCIssuer(
            issuer_id="did:key:issuer",
            vc_storage=vc_storage,
            status_storage=status_storage,
        )
        vc = await issuer.async_issue_identity_credential(
            subject_id="did:key:alice",
            national_id="HT-001",
            first_name="Alice",
            last_name="Citizen",
            date_of_birth="1990-01-15",
            gender="female",
            nationality="HTI",
        )
        assert vc.id is not None
        row = await db_session.scalar(
            select(VerifiableCredentialRecord).where(VerifiableCredentialRecord.credential_id == vc.id)
        )
        assert row is not None
        assert row.subject_id == "did:key:alice"

    @pytest.mark.asyncio
    async def test_async_revoke(self, db_session):
        vc_storage = VCStorage(db_session)
        status_storage = StatusListStorage(db_session)
        issuer = VCIssuer(
            issuer_id="did:key:issuer2",
            vc_storage=vc_storage,
            status_storage=status_storage,
        )
        vc = await issuer.async_issue_identity_credential(
            subject_id="did:key:bob",
            national_id="HT-002",
            first_name="Bob",
            last_name="Test",
            date_of_birth="1985-05-20",
            gender="male",
            nationality="HTI",
        )
        status = await issuer.async_get_credential_status(vc.id)
        assert status == VCStatus.ACTIVE
        ok = await issuer.async_revoke_credential(vc.id)
        assert ok is True
        status = await issuer.async_get_credential_status(vc.id)
        assert status == VCStatus.REVOKED


class TestWalletWithStorage:
    @pytest.mark.asyncio
    async def test_async_store_and_get(self, db_session):
        storage = WalletCredentialStorage(db_session)
        wallet = Wallet(storage=storage)
        cred = {"id": "vc-1", "issuer": "did:key:issuer", "type": ["VerifiableCredential"]}
        record = await wallet.async_store(cred, label="test")
        assert record.id is not None
        row = await db_session.scalar(
            select(WalletCredentialRecord).where(
                WalletCredentialRecord.wallet_did == wallet.did,
                WalletCredentialRecord.credential_id == record.id,
            )
        )
        assert row is not None

    @pytest.mark.asyncio
    async def test_async_list(self, db_session):
        storage = WalletCredentialStorage(db_session)
        wallet = Wallet(storage=storage)
        await wallet.async_store({"id": "vc-a", "type": ["TypeA"]}, label="a")
        await wallet.async_store({"id": "vc-b", "type": ["TypeB"]}, label="b")
        all_records = await wallet.async_list()
        assert len(all_records) == 2
        filtered = await wallet.async_list(credential_type="TypeA")
        assert len(filtered) == 1

    @pytest.mark.asyncio
    async def test_async_delete(self, db_session):
        storage = WalletCredentialStorage(db_session)
        wallet = Wallet(storage=storage)
        record = await wallet.async_store({"id": "vc-del", "type": ["VC"]})
        deleted = await wallet.async_delete(record.id)
        assert deleted is True
        assert await wallet.async_get(record.id) is None

    @pytest.mark.asyncio
    async def test_async_search(self, db_session):
        storage = WalletCredentialStorage(db_session)
        wallet = Wallet(storage=storage)
        await wallet.async_store({"id": "id-cred", "type": ["IdentityCredential"]}, label="ID")
        await wallet.async_store({"id": "dl-cred", "type": ["DriversLicense"]}, label="DL")
        results = await wallet.async_search("Identity")
        assert len(results) == 1


class TestDIDCommWithStorage:
    @pytest.mark.asyncio
    async def test_async_send_and_inbox(self, db_session):
        storage = DIDCommMessageStorage(db_session)
        messenger = DIDCommMessenger(storage=storage)
        msg = DIDCommMessage(
            id="msg-001",
            type="https://didcomm.org/basicmessage/2.0/message",
            body={"text": "hello"},
            from_did="did:key:alice",
            to_did="did:key:bob",
        )
        packed = await messenger.async_send(msg, "did:key:alice", "did:key:bob")
        assert packed is not None
        inbox = await messenger.async_get_inbox("did:key:bob")
        assert len(inbox) == 1
        assert inbox[0].body["text"] == "hello"

    @pytest.mark.asyncio
    async def test_async_receive(self, db_session):
        storage = DIDCommMessageStorage(db_session)
        messenger = DIDCommMessenger(storage=storage)
        msg = DIDCommMessage(
            id="msg-002",
            type="https://didcomm.org/basicmessage/2.0/message",
            body={"text": "hi"},
            from_did="did:key:carol",
            to_did="did:key:dave",
        )
        packed = await messenger.async_send(msg, "did:key:carol", "did:key:dave")
        received = await messenger.async_receive(packed)
        assert received.body["text"] == "hi"


class TestStatusListWithStorage:
    @pytest.mark.asyncio
    async def test_async_create_entry_and_revoke(self, db_session):
        storage = StatusListStorage(db_session)
        mgr = StatusListManager("did:key:issuer", storage=storage)
        entry = await mgr.async_create_entry("revocation")
        assert entry.id is not None
        assert await mgr.async_is_revoked(entry.id) is False
        ok = await mgr.async_revoke(entry.id)
        assert ok is True
        assert await mgr.async_is_revoked(entry.id) is True


class TestCredentialFlowWithStorage:
    @pytest.mark.asyncio
    async def test_async_create_offer(self, db_session):
        from services.vc.verifier import VCVerifier
        storage = CredentialFlowStorage(db_session)
        issuer = VCIssuer("did:key:issuer")
        flow = CredentialFlow(issuer=issuer, storage=storage)
        offer = await flow.async_create_offer(
            issuer_did="did:key:issuer",
            holder_did="did:key:holder",
        )
        assert offer.offer_id is not None
        assert offer.status == "pending"

    @pytest.mark.asyncio
    async def test_async_issue_from_request(self, db_session):
        from services.credential_flow import CredentialRequest
        storage = CredentialFlowStorage(db_session)
        issuer = VCIssuer("did:key:issuer")
        flow = CredentialFlow(issuer=issuer, storage=storage)
        offer = await flow.async_create_offer(
            issuer_did="did:key:issuer",
            holder_did="did:key:holder",
        )
        req = CredentialRequest(
            request_id="req-1",
            offer_id=offer.offer_id,
            holder_did="did:key:holder",
            claims={"national_id": "HT-001", "first_name": "A", "last_name": "B",
                    "date_of_birth": "1990-01-01", "gender": "female", "nationality": "HTI"},
        )
        vc_data = await flow.async_issue_from_request(req)
        assert vc_data is not None


class TestManifestWithStorage:
    @pytest.mark.asyncio
    async def test_async_create_and_get_manifest(self, db_session):
        storage = CredentialManifestStorage(db_session)
        issuer = VCIssuer("did:key:issuer")
        mgr = ManifestManager(issuer=issuer, storage=storage)
        manifest = await mgr.async_create_manifest(
            issuer_did="did:key:issuer",
            name="Test Manifest",
        )
        assert manifest.id is not None
        fetched = await mgr.async_get_manifest(manifest.id)
        assert fetched is not None
        assert fetched.name == "Test Manifest"

    @pytest.mark.asyncio
    async def test_async_list_manifests(self, db_session):
        storage = CredentialManifestStorage(db_session)
        issuer = VCIssuer("did:key:issuer2")
        mgr = ManifestManager(issuer=issuer, storage=storage)
        await mgr.async_create_manifest(issuer_did="did:key:issuer2", name="M1")
        await mgr.async_create_manifest(issuer_did="did:key:issuer2", name="M2")
        manifests = await mgr.async_list_manifests(issuer_did="did:key:issuer2")
        assert len(manifests) == 2


class TestRevocationWithStorage:
    @pytest.mark.asyncio
    async def test_async_notify_and_history(self, db_session):
        storage = RevocationEventStorage(db_session)
        notifier = RevocationNotifier(storage=storage)
        event = await notifier.async_notify(
            credential_id="vc-1",
            event_type=RevocationEventType.CREDENTIAL_REVOKED,
            subject_id="did:key:alice",
            reason="Lost",
        )
        assert event.id is not None
        history = await notifier.async_get_history(credential_id="vc-1")
        assert len(history) >= 1


class TestCHAPIWithStorage:
    @pytest.mark.asyncio
    async def test_chapi_store_with_wallet_storage(self, db_session):
        wallet_storage = WalletCredentialStorage(db_session)
        wallet = Wallet(storage=wallet_storage)
        mediator = CHAPIMediator(wallet=wallet)
        req = CHAPIStoreRequest(
            credential={"id": "chapi-vc-1", "issuer": "did:key:issuer", "type": ["VerifiableCredential"]},
            protocol="vc",
        )
        response = mediator.handle_store(req)
        assert response.error is None
        assert response.data["status"] == "stored"
        records = wallet.list()
        assert len(records) == 1
