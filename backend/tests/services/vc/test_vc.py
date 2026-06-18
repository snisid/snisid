from __future__ import annotations

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from services.pki import KeyAlgorithm
from services.pki.ca import InternalCA
from services.vc import (
    IdentityCredential,
    IdentityCredentialSubject,
    VCStatus,
    VerifiableCredential,
    VerifiablePresentation,
)
from services.vc.api import create_vc_router
from services.vc.issuer import VCIssuer
from services.vc.verifier import VCVerifier


class TestVCModel:
    def test_identity_credential_subject_defaults(self):
        sub = IdentityCredentialSubject(
            id="did:example:123",
            national_id="SN000001",
            first_name="Jean",
            last_name="Dupont",
            date_of_birth="1990-01-15",
            gender="male",
            nationality="HTI",
            status="active",
        )
        assert sub.id == "did:example:123"
        assert sub.national_id == "SN000001"

    def test_verifiable_credential_defaults(self):
        vc = VerifiableCredential()
        assert "VerifiableCredential" in vc.type
        assert vc.issuanceDate is not None
        assert vc.id == ""

    def test_identity_credential_type(self):
        sub = IdentityCredentialSubject(
            id="did:example:1", national_id="SN001",
            first_name="A", last_name="B",
            date_of_birth="2000-01-01", gender="male",
            nationality="HTI", status="active",
        )
        vc = IdentityCredential(credentialSubject=sub)
        assert "SNISIDIdentityCredential" in vc.type
        assert vc.credentialSubject.national_id == "SN001"


class TestVCIssuer:
    def setup_method(self):
        self.issuer = VCIssuer(issuer_id="did:snisid:issuer:1")

    def test_issue_identity_credential(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:123",
            national_id="SN123456",
            first_name="Marie",
            last_name="Curie",
            date_of_birth="1990-01-15",
            gender="female",
            nationality="HTI",
        )
        assert vc.id.startswith("urn:uuid:")
        assert vc.issuer == "did:snisid:issuer:1"
        assert vc.credentialSubject.national_id == "SN123456"
        assert vc.proof is not None
        assert vc.proof["type"] == "SNISID-HMAC-SHA256-2026"
        assert vc.expirationDate is not None

    def test_credential_status_defaults_active(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:456",
            national_id="SN654321",
            first_name="Pierre",
            last_name="Martin",
            date_of_birth="1985-05-20",
            gender="male",
            nationality="HTI",
        )
        status = self.issuer.get_credential_status(vc.id)
        assert status == VCStatus.ACTIVE

    def test_revoke_credential(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:789",
            national_id="SN789012",
            first_name="Alice",
            last_name="Dupuis",
            date_of_birth="1995-08-12",
            gender="female",
            nationality="HTI",
        )
        assert self.issuer.revoke_credential(vc.id) is True
        assert self.issuer.get_credential_status(vc.id) == VCStatus.REVOKED

    def test_revoke_nonexistent_credential(self):
        assert self.issuer.revoke_credential("nonexistent") is False

    def test_suspend_credential(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:101",
            national_id="SN101112",
            first_name="Paul",
            last_name="Pierre",
            date_of_birth="2000-03-01",
            gender="male",
            nationality="HTI",
        )
        assert self.issuer.suspend_credential(vc.id) is True
        assert self.issuer.get_credential_status(vc.id) == VCStatus.REVOKED

    def test_verify_signature_valid(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:202",
            national_id="SN202122",
            first_name="Sophie",
            last_name="Lemaire",
            date_of_birth="1998-12-25",
            gender="female",
            nationality="HTI",
        )
        assert self.issuer.verify_signature(vc) is True

    def test_verify_signature_tampered(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:303",
            national_id="SN303132",
            first_name="Lucas",
            last_name="Petit",
            date_of_birth="2002-06-30",
            gender="male",
            nationality="HTI",
        )
        vc.credentialSubject.national_id = "SN999999"
        assert self.issuer.verify_signature(vc) is False


class TestVCWithPKI:
    def setup_method(self):
        self.ca = InternalCA()
        self.issuer = VCIssuer(
            issuer_id="did:snisid:pki-issuer:1",
            ca=self.ca,
        )

    def test_pki_signed_vc_has_proof(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:pki-1",
            national_id="SN-PKI-001",
            first_name="PKI", last_name="Sign",
            date_of_birth="1990-01-01", gender="male", nationality="HTI",
        )
        assert vc.proof is not None
        assert "ES256" in vc.proof["type"] or "RS256" in vc.proof["type"]

    def test_pki_signed_vc_verify_valid(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:pki-2",
            national_id="SN-PKI-002",
            first_name="Verify", last_name="PKI",
            date_of_birth="1990-01-01", gender="female", nationality="HTI",
        )
        assert self.issuer.verify_signature(vc) is True

    def test_pki_signed_vc_verify_tampered(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:pki-3",
            national_id="SN-PKI-003",
            first_name="Tamper", last_name="Test",
            date_of_birth="1990-01-01", gender="male", nationality="HTI",
        )
        vc.credentialSubject.national_id = "SN-FAKE-999"
        assert self.issuer.verify_signature(vc) is False

    def test_pki_signed_vc_different_algo(self):
        ca_rsa = InternalCA(key_algorithm=KeyAlgorithm.RSA_2048)
        issuer_rsa = VCIssuer(
            issuer_id="did:snisid:rsa-issuer:1",
            ca=ca_rsa,
        )
        vc = issuer_rsa.issue_identity_credential(
            subject_id="did:snisid:citizen:rsa-1",
            national_id="SN-RSA-001",
            first_name="RSA", last_name="Sign",
            date_of_birth="1990-01-01", gender="male", nationality="HTI",
        )
        assert "RS256" in vc.proof["type"]
        assert issuer_rsa.verify_signature(vc) is True


class TestVCVerifier:
    def setup_method(self):
        self.issuer = VCIssuer(issuer_id="did:snisid:issuer:trusted")
        self.verifier = VCVerifier(trusted_issuers=["did:snisid:issuer:trusted"])

    def test_verify_valid_credential(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:1",
            national_id="SN001", first_name="A", last_name="B",
            date_of_birth="2000-01-01", gender="male", nationality="HTI",
        )
        result = self.verifier.verify_credential(vc)
        assert result.valid is True
        assert len(result.errors) == 0

    def test_verify_untrusted_issuer(self):
        issuer2 = VCIssuer(issuer_id="did:untrusted:1")
        vc = issuer2.issue_identity_credential(
            subject_id="did:snisid:citizen:2",
            national_id="SN002", first_name="B", last_name="C",
            date_of_birth="2000-01-01", gender="female", nationality="HTI",
        )
        result = self.verifier.verify_credential(vc)
        assert result.valid is False
        assert any("not trusted" in e for e in result.errors)

    def test_verify_revoked_credential(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:3",
            national_id="SN003", first_name="C", last_name="D",
            date_of_birth="2000-01-01", gender="male", nationality="HTI",
        )
        self.issuer.revoke_credential(vc.id)
        self.verifier.register_status(vc.id, VCStatus.REVOKED)
        result = self.verifier.verify_credential(vc)
        assert result.valid is False
        assert any("revoked" in e for e in result.errors)

    def test_verify_missing_proof(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:4",
            national_id="SN004", first_name="D", last_name="E",
            date_of_birth="2000-01-01", gender="female", nationality="HTI",
        )
        vc.proof = None
        result = self.verifier.verify_credential(vc)
        assert result.valid is False
        assert any("Missing proof" in e for e in result.errors)

    def test_verify_presentation(self):
        vc1 = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:5",
            national_id="SN005", first_name="E", last_name="F",
            date_of_birth="2000-01-01", gender="male", nationality="HTI",
        )
        vc2 = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:6",
            national_id="SN006", first_name="F", last_name="G",
            date_of_birth="2000-01-01", gender="female", nationality="HTI",
        )
        pres = VerifiablePresentation(
            holder="did:snisid:citizen:5",
            verifiableCredential=[vc1, vc2],
        )
        results = self.verifier.verify_presentation(pres)
        assert len(results) == 2
        assert all(r.valid for r in results)

    def test_extract_subject_data(self):
        vc = self.issuer.issue_identity_credential(
            subject_id="did:snisid:citizen:7",
            national_id="SN007", first_name="G", last_name="H",
            date_of_birth="2000-01-01", gender="male", nationality="HTI",
        )
        data = self.verifier.extract_subject_data(vc)
        assert data["sub"] == "did:snisid:citizen:7"
        assert data["vc_id"] == vc.id
        assert data["issuer"] == "did:snisid:issuer:trusted"


class TestVCApi:
    def setup_method(self):
        self.issuer = VCIssuer(issuer_id="http://localhost:8000")
        self.verifier = VCVerifier(
            trusted_issuers=["http://localhost:8000"]
        )
        app = FastAPI()
        router = create_vc_router(issuer=self.issuer, verifier=self.verifier)
        app.include_router(router)
        self.client = TestClient(app)

    def test_issue_credential(self):
        resp = self.client.post("/v1/vc/issue/identity", json={
            "id": "did:snisid:citizen:api-1",
            "national_id": "SN-API-001",
            "first_name": "API",
            "last_name": "Test",
            "date_of_birth": "1990-01-15",
            "gender": "male",
            "nationality": "HTI",
            "status": "active",
        })
        assert resp.status_code == 200
        data = resp.json()
        assert data["credentialSubject"]["national_id"] == "SN-API-001"
        assert data["proof"] is not None

    def test_verify_credential(self):
        issue_resp = self.client.post("/v1/vc/issue/identity", json={
            "id": "did:snisid:citizen:api-2",
            "national_id": "SN-API-002",
            "first_name": "Verify",
            "last_name": "Test",
            "date_of_birth": "1990-01-15",
            "gender": "female",
            "nationality": "HTI",
            "status": "active",
        })
        vc_data = issue_resp.json()
        verify_resp = self.client.post("/v1/vc/verify", json=vc_data)
        assert verify_resp.status_code == 200
        assert verify_resp.json()["verified"] is True

    def test_revoke_credential(self):
        issue_resp = self.client.post("/v1/vc/issue/identity", json={
            "id": "did:snisid:citizen:api-3",
            "national_id": "SN-API-003",
            "first_name": "Revoke",
            "last_name": "Test",
            "date_of_birth": "1990-01-15",
            "gender": "male",
            "nationality": "HTI",
            "status": "active",
        })
        vc_data = issue_resp.json()
        vc_id = vc_data["id"]

        revoke_resp = self.client.post(f"/v1/vc/revoke?vc_id={vc_id}")
        assert revoke_resp.status_code == 200
        assert revoke_resp.json()["status"] == "revoked"

        status_resp = self.client.get(f"/v1/vc/status?vc_id={vc_id}")
        assert status_resp.status_code == 200
        assert status_resp.json()["status"] == "revoked"

        verify_resp = self.client.post("/v1/vc/verify", json=vc_data)
        assert verify_resp.status_code == 200
        assert verify_resp.json()["verified"] is False
