import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from pydantic import ValidationError

from services.did import (
    DIDDocument,
    DIDMethod,
    DIDManager,
    create_did_key,
    create_did_snisid,
    create_did_web,
    resolve_did,
)
from services.did.api import router as did_router


class TestDIDCore:
    def test_create_did_key(self):
        did, private, public = create_did_key()
        assert did.startswith("did:key:")
        assert len(private) == 64
        assert len(public) == 64

    def test_create_did_key_with_existing_key(self):
        did, _, public = create_did_key()
        did2, _, _ = create_did_key(public_key_multibase=public)
        assert did != did2

    def test_create_did_snisid(self):
        did = create_did_snisid("user-123")
        assert did == "did:snisid:mainnet:user-123"

    def test_create_did_snisid_custom_network(self):
        did = create_did_snisid("user-123", network="testnet")
        assert did == "did:snisid:testnet:user-123"

    def test_resolve_did_key(self):
        did, _, _ = create_did_key()
        doc = resolve_did(did)
        assert doc.id == did
        assert len(doc.verification_method) == 1
        assert doc.verification_method[0]["type"] == "Ed25519VerificationKey2018"
        assert doc.authentication == [f"{did}#key-1"]

    def test_resolve_did_snisid(self):
        did = create_did_snisid("user-456")
        doc = resolve_did(did)
        assert doc.id == did
        assert len(doc.verification_method) == 2
        assert len(doc.authentication) == 1
        assert len(doc.capability_delegation) == 1

    def test_resolve_did_key_has_context(self):
        did, _, _ = create_did_key()
        doc = resolve_did(did)
        assert "https://www.w3.org/ns/did/v1" in doc.context

    def test_resolve_unsupported_method(self):
        with pytest.raises(ValueError, match="Unsupported DID method"):
            resolve_did("did:unsupported:123")

    def test_create_did_web(self):
        did = create_did_web("example.com")
        assert did == "did:web:example.com"

    def test_create_did_web_with_path(self):
        did = create_did_web("example.com", "path/to/did")
        assert did == "did:web:example.com:path:to:did"

    def test_resolve_did_web(self):
        did = create_did_web("snisid.ht")
        doc = resolve_did(did)
        assert doc.id == did
        assert len(doc.verification_method) == 1
        assert doc.service[0]["serviceEndpoint"] == "https://snisid.ht"

    def test_did_web_has_linked_domains(self):
        did = create_did_web("gov.ht")
        doc = resolve_did(did)
        assert doc.service[0]["type"] == "LinkedDomains"

    def test_did_document_extra_fields(self):
        doc = DIDDocument(id="did:example:123", custom_field="hello")
        assert doc.custom_field == "hello"

    def test_did_document_validation(self):
        with pytest.raises(ValidationError):
            DIDDocument()


class TestDIDManager:
    @pytest.fixture
    def manager(self):
        return DIDManager()

    def test_create_did_key(self, manager):
        doc = manager.create(DIDMethod.KEY)
        assert doc.id.startswith("did:key:")

    def test_create_did_snisid(self, manager):
        doc = manager.create(DIDMethod.SNISID, identifier="alice")
        assert doc.id == "did:snisid:mainnet:alice"

    def test_create_did_web(self, manager):
        doc = manager.create(DIDMethod.WEB, identifier="example.gov")
        assert doc.id == "did:web:example.gov"
        assert len(doc.service) == 1
        assert doc.service[0]["serviceEndpoint"] == "https://example.gov"

    def test_resolve_created(self, manager):
        doc = manager.create(DIDMethod.KEY)
        resolved = manager.resolve(doc.id)
        assert resolved.id == doc.id

    def test_resolve_external(self, manager):
        did, _, _ = create_did_key()
        resolved = manager.resolve(did)
        assert resolved.id == did

    def test_update_did(self, manager):
        doc = manager.create(DIDMethod.KEY)
        updated = manager.update(doc.id, {"also_known_as": ["did:example:alias"]})
        assert "did:example:alias" in updated.also_known_as
        assert updated.updated is not None

    def test_update_not_found(self, manager):
        with pytest.raises(ValueError, match="DID not found"):
            manager.update("did:key:nonexistent", {})

    def test_deactivate_did(self, manager):
        doc = manager.create(DIDMethod.KEY)
        manager.deactivate(doc.id)
        resolved = manager.resolve(doc.id)
        assert resolved.model_dump().get("deactivated") is True

    def test_deactivate_not_found(self, manager):
        with pytest.raises(ValueError, match="DID not found"):
            manager.deactivate("did:key:nonexistent")


class TestDIDApi:
    @pytest.fixture
    def client(self):
        app = FastAPI()
        app.include_router(did_router)
        return TestClient(app)

    def test_create_did_key(self, client):
        resp = client.post("/did/create", params={"method": "key"})
        assert resp.status_code == 200
        data = resp.json()
        assert data["did"].startswith("did:key:")

    def test_create_did_snisid(self, client):
        resp = client.post("/did/create", params={"method": "snisid", "identifier": "bob"})
        assert resp.status_code == 200
        assert resp.json()["did"] == "did:snisid:mainnet:bob"

    def test_resolve(self, client):
        create_resp = client.post("/did/create", params={"method": "key"})
        did = create_resp.json()["did"]
        resp = client.get(f"/did/resolve/{did}")
        assert resp.status_code == 200
        assert resp.json()["id"] == did

    def test_update(self, client):
        create_resp = client.post("/did/create", params={"method": "key"})
        did = create_resp.json()["did"]
        resp = client.post(f"/did/update/{did}", json={"also_known_as": ["did:example:alias"]})
        assert resp.status_code == 200
        assert "did:example:alias" in resp.json()["also_known_as"]

    def test_deactivate(self, client):
        create_resp = client.post("/did/create", params={"method": "key"})
        did = create_resp.json()["did"]
        resp = client.post(f"/did/deactivate/{did}")
        assert resp.status_code == 200
        assert resp.json()["status"] == "deactivated"
