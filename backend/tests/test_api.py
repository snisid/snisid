from __future__ import annotations

import uuid
from unittest.mock import AsyncMock, MagicMock

import pytest
from fastapi.testclient import TestClient


@pytest.fixture
def app(request, mock_db_session):
    import main
    main.app.dependency_overrides.clear()
    marker = request.node.get_closest_marker("db_healthy")
    db_healthy = marker.args[0] if marker else True
    identity_detail = {
        "id": "test-id-1", "national_id": "SN123456789",
        "first_name": "Jane", "last_name": "Doe", "full_name": "Jane Doe",
        "middle_name": None, "date_of_birth": "1990-06-15",
        "place_of_birth": "Dakar", "gender": "female",
        "nationality": "SEN", "status": "active",
        "agency_id": str(uuid.uuid4()), "verified": False,
        "verified_at": None, "verified_by": None, "photo_url": None,
        "email": None, "phone": None, "marital_status": None,
        "address": None, "document_count": 0, "has_biometrics": False,
        "documents": [], "biometrics": [], "version": 3,
        "created_at": "2025-01-01T00:00:00", "updated_at": "2025-01-01T00:00:00",
    }
    mock_query_handler = MagicMock()
    mock_query_handler.get_by_id = AsyncMock(return_value=identity_detail)
    mock_query_handler.get_by_national_id = AsyncMock(return_value=identity_detail)
    mock_query_handler.search = AsyncMock(return_value={"items": [identity_detail], "total": 1, "page": 1, "page_size": 20})
    mock_query_handler.get_history = AsyncMock(return_value={"items": [], "total": 0})
    mock_query_handler.get_stats = AsyncMock(return_value={"total": 100, "active": 80, "pending": 10, "suspended": 5, "revoked": 5})
    mock_command_handler = MagicMock()
    mock_command_handler.handle_create = AsyncMock(return_value={"identity_id": str(uuid.uuid4()), "status": "pending", "version": 1})
    mock_command_handler.handle_update = AsyncMock(return_value={"identity_id": "test-id-1", "status": "active", "version": 2})
    mock_command_handler.handle_verify = AsyncMock(return_value={"identity_id": "test-id-1", "status": "verified", "version": 3})
    mock_command_handler.handle_suspend = AsyncMock(return_value={"identity_id": "test-id-1", "status": "suspended", "version": 4})
    mock_command_handler.handle_revoke = AsyncMock(return_value={"identity_id": "test-id-1", "status": "revoked", "version": 5})
    mock_command_handler.handle_enroll_biometric = AsyncMock(return_value={"identity_id": "test-id-1", "biometric_type": "fingerprint", "version": 6})
    mock_command_handler.handle_issue_document = AsyncMock(return_value={"identity_id": "test-id-1", "document_type": "passport", "version": 7})
    agency_detail = {
        "agency_id": "test-agency-id", "name": "Test Agency",
        "code": "TST", "agency_type": "government",
        "department": "Interior", "region": "Dakar",
        "status": "active",
    }
    mock_agency_query_handler = MagicMock()
    mock_agency_query_handler.get_by_id = AsyncMock(return_value=agency_detail)
    mock_agency_query_handler.list_all = AsyncMock(return_value={"items": [agency_detail], "total": 1, "page": 1, "page_size": 20})
    mock_agency_command_handler = MagicMock()
    mock_agency_command_handler.handle_create = AsyncMock(return_value={"agency_id": str(uuid.uuid4())})
    mock_agency_command_handler.handle_update = AsyncMock(return_value={"agency_id": "test-agency-id", "status": "active"})
    mock_agency_command_handler.handle_deactivate = AsyncMock(return_value={"agency_id": "test-agency-id", "status": "inactive"})
    main.check_database_health = AsyncMock(return_value=db_healthy)
    main.app.dependency_overrides[main.get_db_session] = lambda: mock_db_session
    main.app.dependency_overrides[main.get_query_handler] = lambda: mock_query_handler
    main.app.dependency_overrides[main.get_command_handler] = lambda: mock_command_handler
    main.app.dependency_overrides[main.get_agency_query_handler] = lambda: mock_agency_query_handler
    main.app.dependency_overrides[main.get_agency_command_handler] = lambda: mock_agency_command_handler
    yield main.app
    main.app.dependency_overrides.clear()


@pytest.fixture
def client(app):
    return TestClient(app)


class TestHealth:
    def test_healthz_ok(self, client):
        resp = client.get("/health")
        assert resp.status_code == 200
        data = resp.json()
        assert data["status"] == "alive"
        assert data["service"] == "snisid"

    @pytest.mark.db_healthy(False)
    def test_healthz_db_unhealthy(self, client):
        resp = client.get("/ready")
        assert resp.status_code == 503


class TestIdentityAPI:
    def test_create_identity(self, client, identity_create_data):
        resp = client.post("/v1/identities", json=identity_create_data)
        assert resp.status_code == 201
        data = resp.json()
        assert "identity_id" in data
        assert data["status"] == "pending"

    def test_create_identity_invalid_dob(self, client):
        payload = {
            "national_id": "SN999999999",
            "first_name": "Bad",
            "last_name": "Date",
            "date_of_birth": "not-a-date",
            "place_of_birth": "Dakar",
            "gender": "male",
            "nationality": "SEN",
            "agency_id": str(uuid.uuid4()),
        }
        resp = client.post("/v1/identities", json=payload)
        assert resp.status_code == 422

    def test_get_identity_by_id(self, client):
        resp = client.get("/v1/identities/test-id-1")
        assert resp.status_code == 200

    def test_get_identity_not_found(self, client):
        resp = client.get("/v1/identities/nonexistent")
        assert resp.status_code == 200

    def test_get_identity_by_national_id(self, client):
        resp = client.get("/v1/identities/by-national-id/SN123456789")
        assert resp.status_code == 200

    def test_search_identities(self, client):
        resp = client.get("/v1/identities?search_term=Jane")
        assert resp.status_code == 200
        data = resp.json()
        assert "items" in data

    def test_get_identity_stats(self, client):
        resp = client.get("/v1/identities/stats")
        assert resp.status_code == 200

    def test_update_identity(self, client):
        resp = client.put("/v1/identities/test-id-1", json={"name": "Updated"})
        assert resp.status_code == 200

    def test_verify_identity(self, client):
        resp = client.post("/v1/identities/test-id-1/verify?verifier_id=verifier-1")
        assert resp.status_code == 200

    def test_suspend_identity(self, client):
        resp = client.post("/v1/identities/test-id-1/suspend?reason=Administrative+suspension")
        assert resp.status_code == 200

    def test_revoke_identity(self, client):
        resp = client.post("/v1/identities/test-id-1/revoke?reason= Permanent+revocation")
        assert resp.status_code == 200

    def test_get_identity_history(self, client):
        resp = client.get("/v1/identities/test-id-1/history")
        assert resp.status_code == 200


class TestAgencyAPI:
    def test_create_agency(self, client):
        payload = {
            "name": "Test Agency",
            "code": "TST",
            "agency_type": "central",
            "department": "Interior",
            "city": "Dakar",
        }
        resp = client.post("/v1/agencies", json=payload)
        assert resp.status_code == 201
        data = resp.json()
        assert "agency_id" in data

    def test_get_agency(self, client):
        resp = client.get("/v1/agencies/some-agency-id")
        assert resp.status_code == 200

    def test_list_agencies(self, client):
        resp = client.get("/v1/agencies")
        assert resp.status_code == 200

    def test_update_agency(self, client):
        resp = client.put("/v1/agencies/test-agency-id", json={"name": "Updated"})
        assert resp.status_code == 200

    def test_deactivate_agency(self, client):
        resp = client.post("/v1/agencies/test-agency-id/deactivate?reason=Administrative+shutdown")
        assert resp.status_code == 200
