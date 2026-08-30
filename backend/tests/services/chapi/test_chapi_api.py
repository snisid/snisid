from __future__ import annotations

import pytest
from fastapi.testclient import TestClient


@pytest.fixture
def client():
    import main
    main.app.dependency_overrides.clear()
    import services.chapi.api as api
    api._mediator = None
    yield TestClient(main.app)
    main.app.dependency_overrides.clear()


class TestCHAPIAPI:
    def test_wallet_info(self, client):
        resp = client.get("/v1/chapi/wallet")
        assert resp.status_code == 200
        data = resp.json()
        assert data["did"].startswith("did:")
        assert data["credential_count"] == 0

    def test_register(self, client):
        resp = client.get("/v1/chapi/register", params={"origin": "https://app.example.com"})
        assert resp.status_code == 200
        data = resp.json()["data"]
        assert data["origin"] == "https://app.example.com"
        assert data["handler"] == "snisid-wallet"

    def test_store_credential(self, client):
        cred = {
            "type": ["VerifiableCredential", "SNISIDIdentityCredential"],
            "issuer": "did:key:issuer",
            "credentialSubject": {"national_id": "SN-001"},
        }
        resp = client.post("/v1/chapi/store", json=cred)
        assert resp.status_code == 200
        data = resp.json()
        assert data["data"]["status"] == "stored"

    def test_get_credentials(self, client):
        client.post("/v1/chapi/store", json={
            "type": ["VerifiableCredential"],
            "credentialSubject": {"name": "Alice"},
        })
        resp = client.post("/v1/chapi/get", json={
            "query": [{"type": "VerifiableCredential"}],
        })
        assert resp.status_code == 200
        data = resp.json()
        assert "verifiableCredential" in data["data"]
        assert len(data["data"]["verifiableCredential"]) == 1

    def test_get_filtered_by_type(self, client):
        client.post("/v1/chapi/store", json={
            "type": ["VerifiableCredential", "PassportCredential"],
            "credentialSubject": {"passport": "PP-001"},
        })
        resp = client.post("/v1/chapi/get", json={
            "query": [{"type": "PassportCredential"}],
        })
        assert resp.status_code == 200
        vcs = resp.json()["data"]["verifiableCredential"]
        assert len(vcs) == 1

    def test_get_no_match(self, client):
        resp = client.post("/v1/chapi/get", json={
            "query": [{"type": "NonexistentCredential"}],
        })
        assert resp.status_code == 200
        assert len(resp.json()["data"]["verifiableCredential"]) == 0

    def test_store_wrapped_in_payload(self, client):
        payload = {
            "credential": {
                "type": ["VerifiableCredential"],
                "credentialSubject": {"name": "Wrapped"},
            },
            "protocol": "vc",
        }
        resp = client.post("/v1/chapi/store", json=payload)
        assert resp.status_code == 200
        assert resp.json()["data"]["status"] == "stored"

    def test_store_then_get_roundtrip(self, client):
        client.post("/v1/chapi/store", json={
            "type": ["VerifiableCredential", "NationalIDCredential"],
            "credentialSubject": {"national_id": "SN-ROUNDTRIP"},
        })
        resp = client.post("/v1/chapi/get", json={
            "query": [{"type": "NationalIDCredential"}],
        })
        vcs = resp.json()["data"]["verifiableCredential"]
        assert len(vcs) == 1
        assert vcs[0]["credentialSubject"]["national_id"] == "SN-ROUNDTRIP"
