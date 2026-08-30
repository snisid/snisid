from __future__ import annotations

import pytest
from fastapi.testclient import TestClient


@pytest.fixture
def client():
    import main
    main.app.dependency_overrides.clear()
    import services.wallet.api as api
    api._wallet = None
    yield TestClient(main.app)
    main.app.dependency_overrides.clear()


class TestWalletAPI:
    def test_did(self, client):
        resp = client.get("/v1/wallet/did")
        assert resp.status_code == 200
        assert resp.json()["did"].startswith("did:")

    def test_did_document(self, client):
        resp = client.get("/v1/wallet/did-document")
        assert resp.status_code == 200
        assert "verification_method" in resp.json()

    def test_store_and_list(self, client):
        cred = {
            "type": ["VerifiableCredential", "SNISIDIdentityCredential"],
            "issuer": "did:key:issuer",
            "credentialSubject": {"national_id": "SN-001"},
        }
        store_resp = client.post("/v1/wallet/credentials", json=cred, params={"label": "My ID"})
        assert store_resp.status_code == 200
        stored = store_resp.json()
        assert stored["label"] == "My ID"

        list_resp = client.get("/v1/wallet/credentials")
        assert list_resp.status_code == 200
        assert list_resp.json()["total"] == 1

    def test_get_credential(self, client):
        cred = {"type": ["VerifiableCredential"], "credentialSubject": {"name": "Alice"}}
        store_resp = client.post("/v1/wallet/credentials", json=cred)
        cid = store_resp.json()["id"]

        get_resp = client.get(f"/v1/wallet/credentials/{cid}")
        assert get_resp.status_code == 200
        assert get_resp.json()["id"] == cid

    def test_get_credential_not_found(self, client):
        resp = client.get("/v1/wallet/credentials/nonexistent")
        assert resp.status_code == 404

    def test_count(self, client):
        resp = client.get("/v1/wallet/credentials/count")
        assert resp.status_code == 200
        assert resp.json()["count"] == 0

    def test_create_presentation(self, client):
        cred = {"type": ["VerifiableCredential"], "credentialSubject": {"name": "Bob"}}
        client.post("/v1/wallet/credentials", json=cred)

        resp = client.post("/v1/wallet/presentations")
        assert resp.status_code == 200
        data = resp.json()
        assert "holder" in data
        assert "verifiableCredential" in data

    def test_verify_presentation(self, client):
        cred = {"type": ["VerifiableCredential"], "credentialSubject": {"name": "Bob"}}
        client.post("/v1/wallet/credentials", json=cred)

        vp_resp = client.post("/v1/wallet/presentations")
        verify_resp = client.post("/v1/wallet/presentations/verify", json=vp_resp.json())
        assert verify_resp.status_code == 200
        assert verify_resp.json()["valid"] is True

    def test_delete_credential(self, client):
        cred = {"type": ["VerifiableCredential"], "credentialSubject": {"name": "Temp"}}
        store_resp = client.post("/v1/wallet/credentials", json=cred)
        cid = store_resp.json()["id"]

        del_resp = client.delete(f"/v1/wallet/credentials/{cid}")
        assert del_resp.status_code == 200

        get_resp = client.get(f"/v1/wallet/credentials/{cid}")
        assert get_resp.status_code == 404

    def test_search(self, client):
        client.post("/v1/wallet/credentials", json={
            "type": ["VerifiableCredential"],
            "credentialSubject": {"name": "Alice Wonderland"},
        })
        resp = client.get("/v1/wallet/credentials/search", params={"q": "alice"})
        assert resp.status_code == 200
        assert resp.json()["total"] == 1

    def test_export_import(self, client):
        client.post("/v1/wallet/credentials", json={
            "type": ["VerifiableCredential"],
            "credentialSubject": {"name": "Export Me"},
        })

        export_resp = client.post("/v1/wallet/export")
        assert export_resp.status_code == 200
        data = export_resp.json()["credentials"]
        assert len(data) == 1

        import services.wallet.api as api
        api._wallet = None

        import_resp = client.post("/v1/wallet/import", json=data)
        assert import_resp.status_code == 200
        assert import_resp.json()["count"] == 1

    def test_list_by_type(self, client):
        client.post("/v1/wallet/credentials", json={
            "type": ["VerifiableCredential", "SNISIDIdentityCredential"],
            "credentialSubject": {"national_id": "SN-001"},
        })
        client.post("/v1/wallet/credentials", json={
            "type": ["VerifiableCredential", "DiplomaCredential"],
            "credentialSubject": {"degree": "BSc"},
        })
        resp = client.get("/v1/wallet/credentials", params={"credential_type": "SNISIDIdentityCredential"})
        assert resp.status_code == 200
        assert resp.json()["total"] == 1

    def test_by_issuer(self, client):
        client.post("/v1/wallet/credentials", json={
            "type": ["VerifiableCredential"],
            "issuer": "did:key:issuer-a",
            "credentialSubject": {"name": "A"},
        })
        resp = client.get("/v1/wallet/credentials/by-issuer/did:key:issuer-a")
        assert resp.status_code == 200
        assert resp.json()["total"] == 1


class TestWalletDIDCommAPI:
    def test_didcomm_send_not_found(self, client):
        resp = client.post(
            "/v1/wallet/didcomm/send/nonexistent",
            params={"to_did": "did:key:bob"},
        )
        assert resp.status_code == 404

    def test_didcomm_send_and_receive(self, client):
        store_resp = client.post("/v1/wallet/credentials", json={
            "type": ["VerifiableCredential"],
            "credentialSubject": {"name": "DIDComm Test"},
        })
        cid = store_resp.json()["id"]

        send_resp = client.post(
            f"/v1/wallet/didcomm/send/{cid}",
            params={"to_did": "did:key:bob"},
        )
        assert send_resp.status_code == 200
        packed = send_resp.json()["packed"]

        import services.wallet.api as api
        api._wallet = None

        recv_resp = client.post(
            "/v1/wallet/didcomm/receive",
            json=packed,
            params={"label": "Via DIDComm"},
        )
        assert recv_resp.status_code == 200
        assert recv_resp.json()["label"] == "Via DIDComm"

    def test_didcomm_receive_empty(self, client):
        recv_resp = client.post(
            "/v1/wallet/didcomm/receive",
            json={"ciphertext": "bad-data"},
            params={"label": "Bad"},
        )
        assert recv_resp.status_code == 400

    def test_didcomm_message(self, client):
        resp = client.post(
            "/v1/wallet/didcomm/message",
            params={
                "message_type": "https://didcomm.org/trust-ping/2.0/ping",
                "to_did": "did:key:friend",
            },
            json={"comment": "hello"},
        )
        assert resp.status_code == 200
        assert resp.json()["status"] == "sent"
