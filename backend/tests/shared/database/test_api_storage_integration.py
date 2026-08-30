"""API-level integration tests: HTTP endpoints backed by SQLAlchemy storage."""
from __future__ import annotations

import pytest
from fastapi import FastAPI
from httpx import AsyncClient, ASGITransport
from sqlalchemy import select

from shared.database.ssi_models import (
    DIDRecord, VerifiableCredentialRecord, WalletCredentialRecord,
    DIDCommMessageRecord, StatusListRecord, CredentialFlowRecord,
    CredentialManifestRecord, RevocationEventRecord,
)
from shared.ssi_storage import (
    DIDStorage, VCStorage, StatusListStorage, WalletCredentialStorage,
    DIDCommMessageStorage, CredentialFlowStorage, CredentialManifestStorage,
    RevocationEventStorage,
)
from services.did import DIDManager, DIDMethod
from services.vc.issuer import VCIssuer
from services.vc.verifier import VCVerifier
from services.vc.api import create_vc_router
from services.wallet import Wallet
from services.didcomm import DIDCommMessage, DIDCommMessenger
from services.status_list import StatusListManager
from services.credential_flow import CredentialFlow
from services.credential_manifest import ManifestManager, OutputDescriptor
from services.revocation import RevocationNotifier, RevocationEventType


class TestVCAPIWithStorage:
    """Test VC API endpoints backed by SQLAlchemy storage."""

    @pytest.mark.asyncio
    async def test_issue_and_stored_in_db(self, db_session):
        vc_storage = VCStorage(db_session)
        status_storage = StatusListStorage(db_session)
        issuer = VCIssuer(
            issuer_id="did:key:vc-api-test",
            vc_storage=vc_storage,
            status_storage=status_storage,
        )
        verifier = VCVerifier()
        router = create_vc_router(issuer=issuer, verifier=verifier)
        app = FastAPI()
        app.include_router(router)

        async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
            body = {
                "id": "did:key:subject",
                "national_id": "HT-API-001",
                "first_name": "Alice",
                "last_name": "Citizen",
                "date_of_birth": "1990-01-15",
                "gender": "female",
                "nationality": "HTI",
                "status": "active",
            }
            resp = await ac.post("/v1/vc/issue/identity", json=body)
            assert resp.status_code == 200
            data = resp.json()
            assert "id" in data
            assert data["issuer"] == "did:key:vc-api-test"

        row = await db_session.scalar(
            select(VerifiableCredentialRecord).where(
                VerifiableCredentialRecord.credential_id == data["id"]
            )
        )
        assert row is not None
        assert row.subject_id == "did:key:subject"

    @pytest.mark.asyncio
    async def test_revoke_and_verify_status(self, db_session):
        vc_storage = VCStorage(db_session)
        status_storage = StatusListStorage(db_session)
        issuer = VCIssuer(
            issuer_id="did:key:vc-revoke-test",
            vc_storage=vc_storage,
            status_storage=status_storage,
        )
        verifier = VCVerifier()
        router = create_vc_router(issuer=issuer, verifier=verifier)
        app = FastAPI()
        app.include_router(router)

        async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
            body = {
                "id": "did:key:revsubject",
                "national_id": "HT-REV-001",
                "first_name": "Bob",
                "last_name": "Test",
                "date_of_birth": "1985-03-10",
                "gender": "male",
                "nationality": "HTI",
                "status": "active",
            }
            resp = await ac.post("/v1/vc/issue/identity", json=body)
            vc_id = resp.json()["id"]

            resp = await ac.post(f"/v1/vc/revoke?vc_id={vc_id}")
            assert resp.status_code == 200
            assert resp.json()["status"] == "revoked"

            resp = await ac.get(f"/v1/vc/status?vc_id={vc_id}")
            assert resp.status_code == 200
            assert resp.json()["status"] == "revoked"

    @pytest.mark.asyncio
    async def test_verify_endpoint(self, db_session):
        vc_storage = VCStorage(db_session)
        status_storage = StatusListStorage(db_session)
        issuer = VCIssuer(
            issuer_id="did:key:verify-test",
            vc_storage=vc_storage,
            status_storage=status_storage,
        )
        verifier = VCVerifier(trusted_issuers=["did:key:verify-test"])
        router = create_vc_router(issuer=issuer, verifier=verifier)
        app = FastAPI()
        app.include_router(router)

        async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
            body = {
                "id": "did:key:verifysub",
                "national_id": "HT-VFY-001",
                "first_name": "Charlie",
                "last_name": "Verify",
                "date_of_birth": "1992-07-22",
                "gender": "male",
                "nationality": "HTI",
                "status": "active",
            }
            resp = await ac.post("/v1/vc/issue/identity", json=body)
            vc_data = resp.json()

            resp = await ac.post("/v1/vc/verify", json=vc_data)
            assert resp.status_code == 200
            assert resp.json()["verified"] is True

        row = await db_session.scalar(
            select(VerifiableCredentialRecord).where(
                VerifiableCredentialRecord.credential_id == vc_data["id"]
            )
        )
        assert row is not None


class TestDIDAPIWithStorage:
    @pytest.mark.asyncio
    async def test_create_and_resolve_via_storage(self, db_session):
        import services.did.api as did_api

        storage = DIDStorage(db_session)
        original_manager = did_api._manager
        did_api._manager = DIDManager(storage=storage)

        try:
            app = FastAPI()
            app.include_router(did_api.router)
            async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
                resp = await ac.post("/did/create", params={"method": "key"})
                assert resp.status_code == 200
                data = resp.json()
                did = data["did"]
                assert did.startswith("did:key:")

                resp = await ac.get(f"/did/resolve/{did}")
                assert resp.status_code == 200
                assert resp.json()["id"] == did

            row = await db_session.scalar(
                select(DIDRecord).where(DIDRecord.did == did)
            )
            assert row is not None
        finally:
            did_api._manager = original_manager


class TestWalletAPIWithStorage:
    @pytest.mark.asyncio
    async def test_store_and_list_via_storage(self, db_session):
        import services.wallet.api as wallet_api

        storage = WalletCredentialStorage(db_session)
        original_wallet = wallet_api._wallet
        wallet_api._wallet = Wallet(storage=storage)

        try:
            app = FastAPI()
            app.include_router(wallet_api.router)
            async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
                resp = await ac.post(
                    "/v1/wallet/credentials",
                    json={"type": ["VerifiableCredential"], "issuer": "did:key:test"},
                )
                assert resp.status_code == 200
                cred_id = resp.json()["id"]

                resp = await ac.get(f"/v1/wallet/credentials/{cred_id}")
                assert resp.status_code == 200

            wallet_did = wallet_api.get_wallet().did
            row = await db_session.scalar(
                select(WalletCredentialRecord).where(
                    WalletCredentialRecord.wallet_did == wallet_did,
                    WalletCredentialRecord.credential_id == cred_id,
                )
            )
            assert row is not None
        finally:
            wallet_api._wallet = original_wallet


class TestDIDCommAPIWithStorage:
    @pytest.mark.asyncio
    async def test_pack_and_unpack_with_storage(self, db_session):
        import services.didcomm.api as didcomm_api

        storage = DIDCommMessageStorage(db_session)
        original_messenger = didcomm_api._messenger
        didcomm_api._messenger = DIDCommMessenger(storage=storage)

        try:
            app = FastAPI()
            app.include_router(didcomm_api.router)
            async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
                resp = await ac.post(
                    "/didcomm/pack",
                    json={
                        "type": "https://didcomm.org/basicmessage/2.0/message",
                        "body": {"text": "storage test"},
                        "from_did": "did:key:alice",
                        "to_did": "did:key:bob",
                    },
                )
                assert resp.status_code == 200
                packed = resp.json()

                resp = await ac.post("/didcomm/unpack", json={"packed": packed})
                assert resp.status_code == 200
                assert resp.json()["body"]["text"] == "storage test"
        finally:
            didcomm_api._messenger = original_messenger


class TestStatusListAPIWithStorage:
    @pytest.mark.asyncio
    async def test_create_and_revoke_entry(self, db_session):
        import services.status_list.api as sl_api

        storage = StatusListStorage(db_session)
        original_manager = sl_api._manager
        sl_api._manager = StatusListManager("did:key:sl-test", storage=storage)

        try:
            app = FastAPI()
            app.include_router(sl_api.router)
            async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
                resp = await ac.post("/status-list/entries", params={"purpose": "revocation"})
                assert resp.status_code == 200
                entry_id = resp.json()["id"]

                resp = await ac.post(f"/status-list/revoke", params={"entry_id": entry_id})
                assert resp.status_code == 200

                resp = await ac.get(f"/status-list/check", params={"entry_id": entry_id})
                assert resp.status_code == 200
                assert resp.json()["revoked"] is True

            list_id = entry_id.rsplit("#", 1)[0]
            row = await db_session.scalar(
                select(StatusListRecord).where(StatusListRecord.list_id == list_id)
            )
            assert row is not None
        finally:
            sl_api._manager = original_manager


class TestCredentialFlowAPIWithStorage:
    @pytest.mark.asyncio
    async def test_create_offer_and_issue(self, db_session):
        import services.credential_flow.api as cf_api

        storage = CredentialFlowStorage(db_session)
        issuer = VCIssuer("did:key:cf-issuer")
        original_flow = cf_api._flow
        cf_api._flow = CredentialFlow(issuer=issuer, storage=storage)

        try:
            app = FastAPI()
            app.include_router(cf_api.router)
            async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
                resp = await ac.post(
                    "/v1/credential-flow/offer",
                    params={
                        "issuer_did": "did:key:cf-issuer",
                        "holder_did": "did:key:cf-holder",
                    },
                )
                assert resp.status_code == 200
                offer_id = resp.json()["offer_id"]

            row = await db_session.scalar(
                select(CredentialFlowRecord).where(
                    CredentialFlowRecord.flow_id == offer_id
                )
            )
            assert row is not None
            assert row.status == "pending"
        finally:
            cf_api._flow = original_flow


class TestManifestAPIWithStorage:
    @pytest.mark.asyncio
    async def test_create_manifest(self, db_session):
        import services.credential_manifest.api as cm_api

        storage = CredentialManifestStorage(db_session)
        issuer = VCIssuer("did:key:cm-issuer")
        original_manager = cm_api._manager
        cm_api._manager = ManifestManager(issuer=issuer, storage=storage)

        try:
            app = FastAPI()
            app.include_router(cm_api.router)
            async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
                resp = await ac.post(
                    "/v1/credential-manifest/manifests",
                    params={"issuer_did": "did:key:cm-issuer", "name": "API Test Manifest"},
                )
                assert resp.status_code == 200
                manifest_id = resp.json()["id"]

                resp = await ac.get(f"/v1/credential-manifest/manifests/{manifest_id}")
                assert resp.status_code == 200
                assert resp.json()["name"] == "API Test Manifest"

            row = await db_session.scalar(
                select(CredentialManifestRecord).where(
                    CredentialManifestRecord.manifest_id == manifest_id
                )
            )
            assert row is not None
        finally:
            cm_api._manager = original_manager


class TestRevocationAPIWithStorage:
    @pytest.mark.asyncio
    async def test_notify_and_history(self, db_session):
        import services.revocation.api as rev_api

        storage = RevocationEventStorage(db_session)
        original_notifier = rev_api._notifier
        rev_api._notifier = RevocationNotifier(storage=storage)

        try:
            app = FastAPI()
            app.include_router(rev_api.router)
            async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
                resp = await ac.post(
                    "/v1/revocation/notify/revoke",
                    params={
                        "credential_id": "vc-revoke-1",
                        "subject_id": "did:key:rev-subject",
                        "reason": "Lost card",
                    },
                )
                assert resp.status_code == 200
                event_data = resp.json()
                assert event_data["credential_id"] == "vc-revoke-1"

                resp = await ac.get(
                    "/v1/revocation/history",
                    params={"credential_id": "vc-revoke-1", "limit": 10},
                )
                assert resp.status_code == 200
                assert resp.json()["total"] >= 1

            rows = await db_session.scalars(
                select(RevocationEventRecord).where(
                    RevocationEventRecord.credential_id == "vc-revoke-1"
                )
            )
            assert len(rows.all()) >= 1
        finally:
            rev_api._notifier = original_notifier
