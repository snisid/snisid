import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from services.vp import VPIssuer, VerifiablePresentation
from services.vp.api import router as vp_router


@pytest.fixture
def issuer():
    return VPIssuer()


@pytest.fixture
def sample_vc():
    return {
        "@context": ["https://www.w3.org/ns/credentials/v2"],
        "id": "urn:uuid:vc-123",
        "type": ["VerifiableCredential"],
        "issuer": "did:key:zabc123",
        "issuanceDate": "2026-01-01T00:00:00Z",
        "credentialSubject": {"id": "did:key:zholder456", "name": "Alice"},
    }


class TestVerifiablePresentation:
    def test_to_dict_minimal(self):
        vp = VerifiablePresentation(
            holder_did="did:key:holder123",
            verifiable_credential=[],
        )
        d = vp.to_dict()
        assert d["holder"] == "did:key:holder123"
        assert d["verifiableCredential"] == []
        assert "https://www.w3.org/ns/credentials/v2" in d["@context"]
        assert "VerifiablePresentation" in d["type"]

    def test_to_dict_with_proof(self):
        vp = VerifiablePresentation(
            holder_did="did:key:holder",
            verifiable_credential=[],
            proof={"type": "TestProof", "proofValue": "abc"},
        )
        d = vp.to_dict()
        assert d["proof"]["proofValue"] == "abc"

    def test_from_dict_roundtrip(self, sample_vc):
        original = VerifiablePresentation(
            holder_did="did:key:holder",
            verifiable_credential=[sample_vc],
            proof={"type": "Test", "proofValue": "xyz"},
        )
        d = original.to_dict()
        restored = VerifiablePresentation.from_dict(d)
        assert restored.holder == original.holder
        assert len(restored.verifiable_credential) == 1
        assert restored.proof is not None

    def test_from_dict_no_proof(self):
        data = {
            "holder": "did:key:holder",
            "verifiableCredential": [],
        }
        vp = VerifiablePresentation.from_dict(data)
        assert vp.proof is None


class TestVPIssuer:
    def test_create_presentation_adds_proof(self, issuer, sample_vc):
        vp = issuer.create_presentation(
            holder_did="did:key:holder",
            verifiable_credentials=[sample_vc],
        )
        assert vp.proof is not None
        assert vp.proof["type"] == "SNISID-HMAC-SHA256-2025"
        assert vp.proof["proofPurpose"] == "authentication"
        assert vp.proof["proofValue"] is not None
        assert vp.proof["verificationMethod"] == "did:key:holder#key-1"

    def test_create_presentation_empty_credentials(self, issuer):
        vp = issuer.create_presentation(
            holder_did="did:key:holder",
            verifiable_credentials=[],
        )
        assert vp.proof is not None
        assert vp.verifiable_credential == []

    def test_create_presentation_custom_vm(self, issuer, sample_vc):
        vm = "did:key:holder#assertion-key"
        vp = issuer.create_presentation(
            holder_did="did:key:holder",
            verifiable_credentials=[sample_vc],
            verification_method=vm,
        )
        assert vp.proof["verificationMethod"] == vm

    def test_verify_valid_presentation(self, issuer, sample_vc):
        vp = issuer.create_presentation(
            holder_did="did:key:holder",
            verifiable_credentials=[sample_vc],
        )
        assert issuer.verify_presentation(vp) is True

    def test_verify_tampered_credential(self, issuer, sample_vc):
        vp = issuer.create_presentation(
            holder_did="did:key:holder",
            verifiable_credentials=[sample_vc],
        )
        vp.verifiable_credential[0]["credentialSubject"]["name"] = "Eve"
        assert issuer.verify_presentation(vp) is False

    def test_verify_no_proof(self, issuer, sample_vc):
        vp = VerifiablePresentation(
            holder_did="did:key:holder",
            verifiable_credential=[sample_vc],
        )
        assert issuer.verify_presentation(vp) is False

    def test_verify_tampered_holder(self, issuer, sample_vc):
        vp = issuer.create_presentation(
            holder_did="did:key:holder",
            verifiable_credentials=[sample_vc],
        )
        vp.holder = "did:key:attacker"
        assert issuer.verify_presentation(vp) is False

    def test_verify_tampered_proof(self, issuer, sample_vc):
        vp = issuer.create_presentation(
            holder_did="did:key:holder",
            verifiable_credentials=[sample_vc],
        )
        vp.proof["proofValue"] = "fake"
        assert issuer.verify_presentation(vp) is False

    def test_verify_multiple_credentials(self, issuer, sample_vc):
        vc2 = dict(sample_vc, id="urn:uuid:vc-456")
        vp = issuer.create_presentation(
            holder_did="did:key:holder",
            verifiable_credentials=[sample_vc, vc2],
        )
        assert issuer.verify_presentation(vp) is True


class TestVPApi:
    @pytest.fixture
    def client(self):
        app = FastAPI()
        app.include_router(vp_router)
        return TestClient(app)

    def test_create_endpoint(self, client):
        resp = client.post(
            "/vp/create",
            json={
                "holder_did": "did:key:holder",
                "verifiable_credentials": [{"id": "vc-1", "type": ["VerifiableCredential"]}],
            },
        )
        assert resp.status_code == 200
        data = resp.json()
        assert data["holder"] == "did:key:holder"
        assert "proof" in data

    def test_create_endpoint_no_vcs(self, client):
        resp = client.post(
            "/vp/create",
            json={
                "holder_did": "did:key:holder",
                "verifiable_credentials": [],
            },
        )
        assert resp.status_code == 200
        assert "proof" in resp.json()

    def test_verify_endpoint_valid(self, client):
        create_resp = client.post(
            "/vp/create",
            json={
                "holder_did": "did:key:holder",
                "verifiable_credentials": [{"id": "vc-1", "type": ["VerifiableCredential"]}],
            },
        )
        vp_data = create_resp.json()
        resp = client.post("/vp/verify", json={"data": vp_data})
        assert resp.status_code == 200
        assert resp.json()["verified"] is True

    def test_verify_endpoint_tampered(self, client):
        create_resp = client.post(
            "/vp/create",
            json={
                "holder_did": "did:key:holder",
                "verifiable_credentials": [{"id": "vc-1", "type": ["VerifiableCredential"]}],
            },
        )
        vp_data = create_resp.json()
        vp_data["holder"] = "did:key:attacker"
        resp = client.post("/vp/verify", json={"data": vp_data})
        assert resp.status_code == 200
        assert resp.json()["verified"] is False

    def test_verify_endpoint_no_proof(self, client):
        resp = client.post("/vp/verify", json={"data": {"holder": "did:key:h", "verifiableCredential": []}})
        assert resp.status_code == 200
        assert resp.json()["verified"] is False
