import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from services.status_list import (
    StatusList,
    StatusList2021Credential,
    StatusListEntry,
    StatusListManager,
)
from services.status_list.api import router as status_list_router


class TestStatusList:
    def test_create_default_size(self):
        sl = StatusList()
        assert sl.size == 131072

    def test_set_and_get_status(self):
        sl = StatusList(size=64)
        assert sl.get_status(0) is False
        sl.set_status(0, revoked=True)
        assert sl.get_status(0) is True
        sl.set_status(0, revoked=False)
        assert sl.get_status(0) is False

    def test_multiple_bits(self):
        sl = StatusList(size=16)
        sl.set_status(0, True)
        sl.set_status(7, True)
        sl.set_status(8, True)
        sl.set_status(15, True)
        assert sl.get_status(0) is True
        assert sl.get_status(7) is True
        assert sl.get_status(8) is True
        assert sl.get_status(15) is True
        assert sl.get_status(1) is False

    def test_encode_decode_roundtrip(self):
        sl = StatusList(size=256)
        sl.set_status(42, True)
        sl.set_status(100, True)
        sl.set_status(200, True)
        encoded = sl.encode()
        restored = StatusList.decode(encoded, size=256)
        assert restored.size == 256
        assert restored.get_status(42) is True
        assert restored.get_status(100) is True
        assert restored.get_status(200) is True
        assert restored.get_status(0) is False

    def test_encode_decode_deterministic(self):
        sl = StatusList(size=128)
        sl.set_status(10, True)
        encoded1 = sl.encode()
        encoded2 = sl.encode()
        assert encoded1 == encoded2

    def test_index_out_of_range(self):
        sl = StatusList(size=8)
        with pytest.raises(ValueError, match="out of range"):
            sl.set_status(8, True)
        with pytest.raises(ValueError, match="out of range"):
            sl.get_status(8)

    def test_to_dict(self):
        sl = StatusList(size=64)
        d = sl.to_dict()
        assert d["size"] == 64
        assert isinstance(d["encodedList"], str)


class TestStatusList2021Credential:
    def test_to_vc_dict_structure(self):
        sl = StatusList(size=64)
        credential = StatusList2021Credential(
            issuer_id="https://issuer.example",
            status_list=sl,
            purpose="revocation",
        )
        vc = credential.to_vc_dict()
        assert vc["issuer"] == "https://issuer.example"
        assert "StatusList2021Credential" in vc["type"]
        assert vc["credentialSubject"]["statusPurpose"] == "revocation"
        assert "encodedList" in vc["credentialSubject"]

    def test_to_vc_dict_with_proof(self):
        sl = StatusList(size=64)
        credential = StatusList2021Credential(
            issuer_id="https://issuer.example",
            status_list=sl,
        )
        credential.set_proof({"type": "TestProof", "proofValue": "abc"})
        vc = credential.to_vc_dict()
        assert vc["proof"]["proofValue"] == "abc"

    def test_context_includes_status(self):
        sl = StatusList(size=64)
        credential = StatusList2021Credential(
            issuer_id="https://issuer.example",
            status_list=sl,
        )
        vc = credential.to_vc_dict()
        assert "https://www.w3.org/ns/credentials/status/v2" in vc["@context"]


class TestStatusListEntry:
    def test_to_dict(self):
        entry = StatusListEntry(
            status_list_credential_id="https://issuer.example/credentials/status/1",
            index=42,
        )
        d = entry.to_dict()
        assert d["type"] == "StatusList2021Entry"
        assert d["statusListIndex"] == "42"
        assert "statusListCredential" in d
        assert d["id"] == "https://issuer.example/credentials/status/1#42"

    def test_custom_purpose(self):
        entry = StatusListEntry(
            status_list_credential_id="https://issuer.example/credentials/status/1",
            index=0,
            purpose="suspension",
        )
        assert entry.statusPurpose == "suspension"


class TestStatusListManager:
    @pytest.fixture
    def manager(self):
        return StatusListManager(issuer_id="https://snisid.ht/issuer")

    def test_create_entry(self, manager):
        entry = manager.create_entry()
        assert entry.type == "StatusList2021Entry"
        assert entry.statusListIndex == "0"
        assert "statusListCredential" in entry.to_dict()

    def test_create_multiple_entries(self, manager):
        e1 = manager.create_entry()
        e2 = manager.create_entry()
        assert e1.statusListIndex == "0"
        assert e2.statusListIndex == "1"

    def test_revoke_and_check(self, manager):
        entry = manager.create_entry()
        assert manager.is_revoked(entry.id) is False
        assert manager.revoke(entry.id) is True
        assert manager.is_revoked(entry.id) is True

    def test_unrevoke(self, manager):
        entry = manager.create_entry()
        manager.revoke(entry.id)
        assert manager.is_revoked(entry.id) is True
        manager.unrevoke(entry.id)
        assert manager.is_revoked(entry.id) is False

    def test_revoke_nonexistent(self, manager):
        assert manager.revoke("nonexistent") is False

    def test_check_nonexistent(self, manager):
        assert manager.is_revoked("nonexistent") is None

    def test_get_status_list_credential(self, manager):
        manager.create_entry()
        list_id = f"{manager._issuer_id}/credentials/status/1"
        credential = manager.get_status_list_credential(list_id)
        assert credential is not None
        assert credential["issuer"] == manager._issuer_id

    def test_get_status_list_credential_not_found(self, manager):
        assert manager.get_status_list_credential("nonexistent") is None

    def test_auto_create_new_list_when_full(self, manager):
        sl, list_id = manager.get_or_create_list(size=8)
        for i in range(8):
            manager.create_entry()
        entry_9 = manager.create_entry()
        assert entry_9.statusListIndex == "0"
        assert "status/2" in entry_9.statusListCredential


class TestStatusListAPI:
    @pytest.fixture
    def client(self):
        app = FastAPI()
        app.include_router(status_list_router)
        return TestClient(app)

    def test_create_entry(self, client):
        resp = client.post("/status-list/entries")
        assert resp.status_code == 200
        data = resp.json()
        assert data["type"] == "StatusList2021Entry"
        assert data["statusListIndex"] == "0"

    def test_revoke_and_check(self, client):
        create_resp = client.post("/status-list/entries")
        entry_id = create_resp.json()["id"]
        revoke_resp = client.post("/status-list/revoke", params={"entry_id": entry_id})
        assert revoke_resp.status_code == 200
        check_resp = client.get("/status-list/check", params={"entry_id": entry_id})
        assert check_resp.json()["revoked"] is True

    def test_unrevoke(self, client):
        create_resp = client.post("/status-list/entries")
        entry_id = create_resp.json()["id"]
        client.post("/status-list/revoke", params={"entry_id": entry_id})
        resp = client.post("/status-list/unrevoke", params={"entry_id": entry_id})
        assert resp.status_code == 200
        assert client.get("/status-list/check", params={"entry_id": entry_id}).json()["revoked"] is False

    def test_revoke_not_found(self, client):
        resp = client.post("/status-list/revoke", params={"entry_id": "nonexistent"})
        assert resp.status_code == 404

    def test_check_not_found(self, client):
        resp = client.get("/status-list/check", params={"entry_id": "nonexistent"})
        assert resp.status_code == 404

    def test_get_status_list_credential(self, client):
        client.post("/status-list/entries")
        resp = client.get("/status-list/credential/https://snisid.ht/issuer/credentials/status/1")
        assert resp.status_code == 200
        data = resp.json()
        assert "StatusList2021Credential" in data["type"]
