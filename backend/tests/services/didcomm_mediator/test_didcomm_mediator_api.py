import pytest
from fastapi.testclient import TestClient


@pytest.fixture
def client():
    import main
    main.app.dependency_overrides.clear()
    import services.didcomm_mediator.api as api
    api._mediator = None
    yield TestClient(main.app)
    main.app.dependency_overrides.clear()


class TestDIDCommMediatorAPI:
    def test_forward_packed(self, client):
        resp = client.post(
            "/v1/didcomm/mediator/forward/packed",
            params={
                "recipient_did": "did:key:bob",
                "sender_did": "did:key:alice",
            },
            json={"ciphertext": "encrypted"},
        )
        assert resp.status_code == 200
        data = resp.json()
        assert data["recipient_did"] == "did:key:bob"
        assert "message_id" in data

    def test_forward_request(self, client):
        payload = {
            "type": "https://didcomm.org/routing/2.0/forward",
            "id": "fwd-1",
            "to": ["did:key:bob"],
            "from": "did:key:alice",
            "body": {"next": "did:key:bob"},
            "attachments": [{"data": {"json": {"ciphertext": "secret"}}}],
        }
        resp = client.post("/v1/didcomm/mediator/forward", json=payload)
        assert resp.status_code == 200
        assert resp.json()["recipient_did"] == "did:key:bob"

    def test_forward_invalid(self, client):
        resp = client.post("/v1/didcomm/mediator/forward", json={})
        assert resp.status_code == 400

    def test_get_inbox(self, client):
        client.post(
            "/v1/didcomm/mediator/forward/packed",
            params={"recipient_did": "did:key:bob", "sender_did": "did:key:alice"},
            json={"data": "msg1"},
        )
        resp = client.get("/v1/didcomm/mediator/inbox/did:key:bob")
        assert resp.status_code == 200
        assert resp.json()["total"] == 1

    def test_get_inbox_empty(self, client):
        resp = client.get("/v1/didcomm/mediator/inbox/did:key:nobody")
        assert resp.status_code == 200
        assert resp.json()["total"] == 0

    def test_pending_messages(self, client):
        client.post(
            "/v1/didcomm/mediator/forward/packed",
            params={"recipient_did": "did:key:bob"},
            json={"data": "pending"},
        )
        resp = client.get("/v1/didcomm/mediator/inbox/did:key:bob/pending")
        assert resp.status_code == 200
        assert resp.json()["total"] == 1

    def test_deliver_message(self, client):
        forward = client.post(
            "/v1/didcomm/mediator/forward/packed",
            params={"recipient_did": "did:key:bob"},
            json={"data": "to-deliver"},
        )
        msg_id = forward.json()["message_id"]
        resp = client.post(f"/v1/didcomm/mediator/deliver/{msg_id}")
        assert resp.status_code == 200
        assert resp.json()["delivered_at"] is not None

    def test_deliver_nonexistent(self, client):
        resp = client.post("/v1/didcomm/mediator/deliver/nonexistent")
        assert resp.status_code == 404

    def test_pending_count(self, client):
        client.post(
            "/v1/didcomm/mediator/forward/packed",
            params={"recipient_did": "did:key:bob"},
            json={"data": "count-me"},
        )
        resp = client.get("/v1/didcomm/mediator/pending-count/did:key:bob")
        assert resp.status_code == 200
        assert resp.json()["pending_count"] == 1

    def test_delete_message(self, client):
        forward = client.post(
            "/v1/didcomm/mediator/forward/packed",
            params={"recipient_did": "did:key:bob"},
            json={"data": "delete-me"},
        )
        msg_id = forward.json()["message_id"]
        resp = client.delete(f"/v1/didcomm/mediator/messages/{msg_id}")
        assert resp.status_code == 200

    def test_delete_nonexistent(self, client):
        resp = client.delete("/v1/didcomm/mediator/messages/nonexistent")
        assert resp.status_code == 404
