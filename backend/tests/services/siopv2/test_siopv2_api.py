from __future__ import annotations

import pytest
from fastapi.testclient import TestClient


@pytest.fixture
def client():
    import main
    main.app.dependency_overrides.clear()
    import services.siopv2.api as api
    api._wallet = None
    api._requests = {}
    yield TestClient(main.app)
    main.app.dependency_overrides.clear()


class TestSIOPv2API:
    def test_wallet_did(self, client):
        resp = client.get("/v1/siopv2/wallet/did")
        assert resp.status_code == 200
        data = resp.json()
        assert data["did"].startswith("did:")

    def test_wallet_empty_credentials(self, client):
        resp = client.get("/v1/siopv2/wallet/credentials")
        assert resp.status_code == 200
        assert resp.json()["credentials"] == []

    def test_store_credential(self, client):
        cred = {
            "type": ["VerifiableCredential"],
            "credentialSubject": {"name": "Alice"},
        }
        resp = client.post("/v1/siopv2/wallet/store", json=cred)
        assert resp.status_code == 200
        assert resp.json()["count"] == 1

    def test_create_verifier_request(self, client):
        resp = client.post(
            "/v1/siopv2/verifier/request",
            params={"client_id": "did:key:test-verifier"},
        )
        assert resp.status_code == 200
        data = resp.json()
        assert "request" in data
        assert data["request"]["client_id"] == "did:key:test-verifier"
        assert data["request"]["scope"] == "openid"

    def test_create_request_with_pd(self, client):
        resp = client.post(
            "/v1/siopv2/verifier/request-with-pd",
            params={"input_descriptor_ids": "my-cred", "pd_name": "Check"},
        )
        assert resp.status_code == 200
        data = resp.json()
        pd = data["request"]["presentation_definition"]
        assert pd["input_descriptors"][0]["id"] == "my-cred"

    def test_wallet_respond(self, client):
        client.post("/v1/siopv2/wallet/store", json={
            "type": ["VerifiableCredential"],
            "credentialSubject": {"national_id": "SN-001"},
        })
        req_resp = client.post("/v1/siopv2/verifier/request")
        req_data = req_resp.json()["request"]

        resp = client.post("/v1/siopv2/wallet/respond", json=req_data)
        assert resp.status_code == 200
        data = resp.json()
        assert "id_token" in data

    def test_full_verifier_verify(self, client):
        client.post("/v1/siopv2/wallet/store", json={
            "type": ["VerifiableCredential", "SNISIDIdentityCredential"],
            "credentialSubject": {"national_id": "SN-002"},
        })

        req_resp = client.post("/v1/siopv2/verifier/request-with-pd")
        assert req_resp.status_code == 200
        req_data = req_resp.json()
        nonce = req_data["_internal_nonce"]
        request = req_data["request"]

        wallet_resp = client.post("/v1/siopv2/wallet/respond", json=request)
        assert wallet_resp.status_code == 200
        response_data = wallet_resp.json()

        verify_resp = client.post(
            "/v1/siopv2/verifier/verify",
            params={"expected_nonce": nonce},
            json=response_data,
        )
        assert verify_resp.status_code == 200
        result = verify_resp.json()
        assert result["valid"] is True

    def test_verify_fails_with_wrong_nonce(self, client):
        client.post("/v1/siopv2/wallet/store", json={
            "type": ["VerifiableCredential"],
        })

        req_resp = client.post("/v1/siopv2/verifier/request")
        request = req_resp.json()["request"]

        wallet_resp = client.post("/v1/siopv2/wallet/respond", json=request)
        response_data = wallet_resp.json()

        verify_resp = client.post(
            "/v1/siopv2/verifier/verify",
            params={"expected_nonce": "wrong-nonce"},
            json=response_data,
        )
        result = verify_resp.json()
        assert result["valid"] is False

    def test_respond_with_pd(self, client):
        client.post("/v1/siopv2/wallet/store", json={
            "type": ["VerifiableCredential"],
            "credentialSubject": {"name": "Bob"},
        })

        req_resp = client.post(
            "/v1/siopv2/verifier/request-with-pd",
            params={"input_descriptor_ids": "id-cred", "pd_name": "Identity"},
        )
        request = req_resp.json()["request"]

        resp = client.post("/v1/siopv2/wallet/respond", json=request)
        data = resp.json()
        assert "vp_token" in data
        assert "presentation_submission" in data
