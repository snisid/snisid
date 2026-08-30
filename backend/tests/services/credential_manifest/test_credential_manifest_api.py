from __future__ import annotations

import pytest
from fastapi.testclient import TestClient


@pytest.fixture
def client():
    import main
    main.app.dependency_overrides.clear()
    import services.credential_manifest.api as api
    api._manager = None
    yield TestClient(main.app)
    main.app.dependency_overrides.clear()


class TestCredentialManifestAPI:
    def test_create_manifest(self, client):
        resp = client.post(
            "/v1/credential-manifest/manifests",
            params={
                "issuer_did": "did:key:issuer",
                "name": "National ID",
            },
        )
        assert resp.status_code == 200
        data = resp.json()
        assert data["issuer"] == "did:key:issuer"
        assert data["name"] == "National ID"

    def test_list_manifests(self, client):
        client.post("/v1/credential-manifest/manifests", params={"issuer_did": "did:key:a"})
        client.post("/v1/credential-manifest/manifests", params={"issuer_did": "did:key:b"})
        resp = client.get("/v1/credential-manifest/manifests")
        assert resp.status_code == 200
        assert resp.json()["total"] == 2

    def test_get_manifest_not_found(self, client):
        resp = client.get("/v1/credential-manifest/manifests/nonexistent")
        assert resp.status_code == 404

    def test_get_manifest(self, client):
        create_resp = client.post(
            "/v1/credential-manifest/manifests",
            params={"issuer_did": "did:key:test", "name": "Test"},
        )
        mid = create_resp.json()["id"]
        get_resp = client.get(f"/v1/credential-manifest/manifests/{mid}")
        assert get_resp.status_code == 200
        assert get_resp.json()["name"] == "Test"

    def test_apply_success(self, client):
        create_resp = client.post(
            "/v1/credential-manifest/manifests",
            params={"issuer_did": "did:key:issuer"},
        )
        mid = create_resp.json()["id"]

        apply_resp = client.post("/v1/credential-manifest/apply", json={
            "manifest_id": mid,
            "applicant": "did:key:alice",
            "claims": {
                "national_id": "SN-API-001",
                "first_name": "Alice",
                "last_name": "Smith",
                "date_of_birth": "1990-01-01",
                "gender": "female",
                "nationality": "HTI",
            },
        })
        assert apply_resp.status_code == 200
        data = apply_resp.json()
        assert len(data["credentials"]) == 1
        assert data["credentials"][0]["credentialSubject"]["national_id"] == "SN-API-001"

    def test_apply_missing_fields(self, client):
        resp = client.post("/v1/credential-manifest/apply", json={})
        assert resp.status_code == 400

    def test_apply_unknown_manifest(self, client):
        resp = client.post("/v1/credential-manifest/apply", json={
            "manifest_id": "nonexistent",
            "applicant": "did:key:alice",
        })
        assert resp.status_code == 400

    def test_list_applications(self, client):
        create_resp = client.post(
            "/v1/credential-manifest/manifests",
            params={"issuer_did": "did:key:issuer"},
        )
        mid = create_resp.json()["id"]
        client.post("/v1/credential-manifest/apply", json={
            "manifest_id": mid,
            "applicant": "did:key:bob",
            "claims": {"national_id": "SN-002", "first_name": "Bob", "last_name": "Test",
                       "date_of_birth": "2000-01-01", "gender": "male", "nationality": "HTI"},
        })
        resp = client.get("/v1/credential-manifest/applications")
        assert resp.status_code == 200
        assert resp.json()["total"] == 1

    def test_get_response(self, client):
        create_resp = client.post(
            "/v1/credential-manifest/manifests",
            params={"issuer_did": "did:key:issuer"},
        )
        mid = create_resp.json()["id"]

        apply_resp = client.post("/v1/credential-manifest/apply", json={
            "manifest_id": mid,
            "applicant": "did:key:charlie",
            "claims": {"national_id": "SN-003", "first_name": "Charlie", "last_name": "Test",
                       "date_of_birth": "1995-05-05", "gender": "female", "nationality": "HTI"},
        })
        rid = apply_resp.json()["id"]

        get_resp = client.get(f"/v1/credential-manifest/responses/{rid}")
        assert get_resp.status_code == 200
        assert len(get_resp.json()["credentials"]) == 1

    def test_get_response_not_found(self, client):
        resp = client.get("/v1/credential-manifest/responses/nonexistent")
        assert resp.status_code == 404
