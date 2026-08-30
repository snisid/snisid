from __future__ import annotations

import pytest
from fastapi.testclient import TestClient


@pytest.fixture
def client():
    import main
    main.app.dependency_overrides.clear()
    import services.revocation.api as api
    api._notifier = None
    api._hooks = {}
    yield TestClient(main.app)
    main.app.dependency_overrides.clear()


class TestRevocationAPI:
    def test_notify_revoke(self, client):
        resp = client.post(
            "/v1/revocation/notify/revoke",
            params={
                "credential_id": "vc-001",
                "subject_id": "did:key:alice",
                "reason": "Lost ID",
            },
        )
        assert resp.status_code == 200
        data = resp.json()
        assert data["event_type"] == "credential.revoked"
        assert data["credential_id"] == "vc-001"

    def test_notify_suspend(self, client):
        resp = client.post(
            "/v1/revocation/notify/suspend",
            params={
                "credential_id": "vc-002",
                "subject_id": "did:key:bob",
            },
        )
        assert resp.status_code == 200
        assert resp.json()["event_type"] == "credential.suspended"

    def test_notify_reinstate(self, client):
        resp = client.post(
            "/v1/revocation/notify/reinstate",
            params={
                "credential_id": "vc-003",
                "subject_id": "did:key:charlie",
                "reason": "Error corrected",
            },
        )
        assert resp.status_code == 200
        assert resp.json()["event_type"] == "credential.reinstated"

    def test_history(self, client):
        client.post("/v1/revocation/notify/revoke", params={"credential_id": "vc-1", "subject_id": "did:key:a"})
        client.post("/v1/revocation/notify/revoke", params={"credential_id": "vc-2", "subject_id": "did:key:b"})
        resp = client.get("/v1/revocation/history")
        assert resp.status_code == 200
        assert resp.json()["total"] == 2

    def test_history_filtered(self, client):
        client.post("/v1/revocation/notify/revoke", params={"credential_id": "vc-a", "subject_id": "did:key:a"})
        client.post("/v1/revocation/notify/suspend", params={"credential_id": "vc-a", "subject_id": "did:key:a"})
        resp = client.get("/v1/revocation/history", params={"credential_id": "vc-a"})
        assert resp.status_code == 200
        assert resp.json()["total"] == 2

    def test_track_and_notify(self, client):
        client.post("/v1/revocation/wallet/did:key:alice/track/vc-tracked")
        client.post(
            "/v1/revocation/notify/revoke",
            params={"credential_id": "vc-tracked", "subject_id": "did:key:alice"},
        )
        resp = client.get("/v1/revocation/wallet/did:key:alice/notifications")
        assert resp.status_code == 200
        assert resp.json()["total"] == 1

    def test_untrack(self, client):
        client.post("/v1/revocation/wallet/did:key:bob/track/vc-untrack")
        client.post("/v1/revocation/wallet/did:key:bob/untrack/vc-untrack")
        client.post(
            "/v1/revocation/notify/revoke",
            params={"credential_id": "vc-untrack", "subject_id": "did:key:bob"},
        )
        resp = client.get("/v1/revocation/wallet/did:key:bob/notifications")
        assert resp.status_code == 200
        assert resp.json()["total"] == 0

    def test_check_status(self, client):
        client.post("/v1/revocation/wallet/did:key:charlie/track/vc-status")
        client.post(
            "/v1/revocation/notify/suspend",
            params={"credential_id": "vc-status", "subject_id": "did:key:charlie"},
        )
        resp = client.get("/v1/revocation/wallet/did:key:charlie/status/vc-status")
        assert resp.status_code == 200
        assert resp.json()["status"] == "credential.suspended"

    def test_clear_notifications(self, client):
        client.post("/v1/revocation/wallet/did:key:dave/track/vc-clear")
        client.post(
            "/v1/revocation/notify/revoke",
            params={"credential_id": "vc-clear", "subject_id": "did:key:dave"},
        )
        client.post("/v1/revocation/wallet/did:key:dave/notifications/clear")
        resp = client.get("/v1/revocation/wallet/did:key:dave/notifications")
        assert resp.json()["total"] == 0
