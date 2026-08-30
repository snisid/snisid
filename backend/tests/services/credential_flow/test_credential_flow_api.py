from __future__ import annotations

import pytest
from fastapi.testclient import TestClient


@pytest.fixture(autouse=True)
def reset_flow():
    import services.credential_flow.api as api
    api._flow = None


@pytest.fixture
def client():
    import main
    main.app.dependency_overrides.clear()
    yield TestClient(main.app)
    main.app.dependency_overrides.clear()


class TestCredentialFlowAPI:
    def test_create_offer(self, client):
        resp = client.post(
            "/v1/credential-flow/offer",
            params={
                "issuer_did": "did:key:issuer-001",
                "holder_did": "did:snisid:mainnet:citizen-001",
            },
        )
        assert resp.status_code == 200
        data = resp.json()
        assert data["issuer_did"] == "did:key:issuer-001"
        assert data["holder_did"] == "did:snisid:mainnet:citizen-001"
        assert data["status"] == "pending"

    def test_create_and_send_offer(self, client):
        resp = client.post(
            "/v1/credential-flow/offer",
            params={
                "issuer_did": "did:key:issuer-002",
                "holder_did": "did:snisid:mainnet:citizen-002",
            },
        )
        assert resp.status_code == 200
        offer = resp.json()
        offer_id = offer["offer_id"]

        send_resp = client.post(f"/v1/credential-flow/offer/{offer_id}/send")
        assert send_resp.status_code == 200
        send_data = send_resp.json()
        assert send_data["status"] == "sent"

    def test_send_offer_not_found(self, client):
        resp = client.post("/v1/credential-flow/offer/nonexistent/send")
        assert resp.status_code == 404

    def test_request_requires_issuer_did(self, client):
        resp = client.post(
            "/v1/credential-flow/request",
            params={
                "offer_id": "offer-999",
                "holder_did": "did:snisid:mainnet:citizen",
                "issuer_did": "did:key:issuer",
            },
        )
        assert resp.status_code == 200
        data = resp.json()
        assert data["status"] == "requested"
        assert "packed" in data

    def test_receive_request(self, client):
        resp = client.post(
            "/v1/credential-flow/request",
            params={
                "offer_id": "offer-recv",
                "holder_did": "did:snisid:mainnet:recv",
                "issuer_did": "did:key:issuer",
            },
        )
        packed = resp.json()["packed"]

        recv_resp = client.post("/v1/credential-flow/receive-request", json=packed)
        assert recv_resp.status_code == 200
        req = recv_resp.json()
        assert req["offer_id"] == "offer-recv"

    def test_full_api_flow(self, client):
        iss_did = "did:key:issuer-api"
        hol_did = "did:snisid:mainnet:citizen-api"

        resp = client.post(
            "/v1/credential-flow/offer",
            params={"issuer_did": iss_did, "holder_did": hol_did},
        )
        assert resp.status_code == 200
        offer = resp.json()

        resp = client.post(
            "/v1/credential-flow/request",
            params={
                "offer_id": offer["offer_id"],
                "holder_did": hol_did,
                "issuer_did": iss_did,
            },
        )
        assert resp.status_code == 200
        req_data = resp.json()

        resp = client.post(
            "/v1/credential-flow/receive-request",
            json=req_data["packed"],
        )
        assert resp.status_code == 200
        received_req = resp.json()

        resp = client.post(
            "/v1/credential-flow/issue",
            json=received_req,
        )
        assert resp.status_code == 200
        vc_data = resp.json()
        assert vc_data["credentialSubject"]["national_id"] != ""
