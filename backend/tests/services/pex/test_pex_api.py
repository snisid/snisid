import pytest
from fastapi.testclient import TestClient
from services.wallet import Wallet


@pytest.fixture
def client():
    import main
    main.app.dependency_overrides.clear()
    import services.pex.api as api
    api._wallet_cache.clear()
    yield TestClient(main.app)
    main.app.dependency_overrides.clear()


@pytest.fixture
def wallet_with_creds():
    wallet = Wallet()
    wallet.store(
        credential={
            "issuer": "did:key:issuer1",
            "issuanceDate": "2026-01-01T00:00:00Z",
            "type": ["VerifiableCredential", "IdentityCredential"],
            "credentialSubject": {
                "id": "did:key:alice",
                "givenName": "Alice",
                "familyName": "Smith",
            },
        },
        label="Identity Credential",
    )
    wallet.store(
        credential={
            "issuer": "did:key:dmv",
            "issuanceDate": "2026-03-01T00:00:00Z",
            "type": ["VerifiableCredential", "DriversLicenseCredential"],
            "credentialSubject": {
                "id": "did:key:alice",
                "licenseNumber": "DL-12345",
            },
        },
        label="Drivers License",
    )
    import services.pex.api as api
    api._wallet_cache["did:key:alice"] = wallet
    return wallet


_descriptor = {
    "id": "pd-1",
    "input_descriptors": [
        {
            "id": "id_1",
            "name": "Identity Check",
            "purpose": "Verify name",
            "schema": [{"uri": "https://example.com/IdentityCredential"}],
        }
    ],
}


class TestPEXAPI:
    def test_evaluate_valid(self, client, wallet_with_creds):
        resp = client.post(
            "/v1/pex/evaluate",
            json={"definition": _descriptor, "holder_did": "did:key:alice"},
        )
        assert resp.status_code == 200
        assert resp.json()["valid"] is True

    def test_evaluate_no_match(self, client, wallet_with_creds):
        desc = {
            "id": "pd-2",
            "input_descriptors": [
                {
                    "id": "id_2",
                    "schema": [{"uri": "https://example.com/PassportCredential"}],
                }
            ],
        }
        resp = client.post(
            "/v1/pex/evaluate",
            json={"definition": desc, "holder_did": "did:key:alice"},
        )
        assert resp.status_code == 200
        assert resp.json()["valid"] is False

    def test_filter_credentials(self, client, wallet_with_creds):
        resp = client.post(
            "/v1/pex/filter",
            json={"definition": _descriptor, "holder_did": "did:key:alice"},
        )
        assert resp.status_code == 200
        data = resp.json()
        assert data["total_matched"] == 1
        assert data["total_input"] == 2

    def test_filter_empty_wallet(self, client):
        desc = {
            "id": "pd-3",
            "input_descriptors": [
                {"id": "id_1", "schema": [{"uri": "https://example.com/IdentityCredential"}]}
            ],
        }
        resp = client.post(
            "/v1/pex/filter",
            json={"definition": desc, "holder_did": "did:key:empty"},
        )
        assert resp.status_code == 200
        assert resp.json()["total_matched"] == 0

    def test_create_presentation(self, client, wallet_with_creds):
        resp = client.post(
            "/v1/pex/present",
            json={
                "definition": _descriptor,
                "holder_did": "did:key:alice",
                "issuer_did": "did:key:issuer1",
            },
        )
        assert resp.status_code == 200
        data = resp.json()
        assert "holder" in data
        assert data["holder"] == "did:key:alice"

    def test_create_presentation_no_match(self, client):
        desc = {
            "id": "pd-4",
            "input_descriptors": [
                {"id": "id_1", "schema": [{"uri": "https://example.com/NonexistentCredential"}]}
            ],
        }
        resp = client.post(
            "/v1/pex/present",
            json={"definition": desc, "holder_did": "did:key:alice"},
        )
        assert resp.status_code == 400
