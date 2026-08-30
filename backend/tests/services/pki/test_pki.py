from __future__ import annotations

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from cryptography import x509

from services.pki import CertificateStatus, KeyAlgorithm
from services.pki.api import create_pki_router
from services.pki.ca import InternalCA


class TestInternalCA:
    def test_ca_generation(self):
        ca = InternalCA()
        info = ca.get_ca_info()
        assert info.ca_subject == "CN=SNISID Internal CA, O=SNISID, C=HT"
        assert info.ca_cert_pem.startswith("-----BEGIN CERTIFICATE-----")
        assert len(info.ca_fingerprint) == 64

    def test_issue_certificate(self):
        ca = InternalCA()
        cert = ca.issue_certificate(subject_cn="test.snisid.ht")
        assert cert.subject == "test.snisid.ht"
        assert cert.issuer == "CN=SNISID Internal CA, O=SNISID, C=HT"
        assert cert.certificate_pem.startswith("-----BEGIN CERTIFICATE-----")
        assert cert.status == CertificateStatus.ACTIVE
        assert cert.serial_number is not None

    def test_issue_certificate_with_sans(self):
        ca = InternalCA()
        cert = ca.issue_certificate(
            subject_cn="api.snisid.ht",
            subject_alt_names=["api.snisid.ht", "localhost"],
        )
        assert "api.snisid.ht" in cert.subject_alt_names
        assert "localhost" in cert.subject_alt_names

    def test_issue_rsa_certificate(self):
        ca = InternalCA(key_algorithm=KeyAlgorithm.RSA_2048)
        cert = ca.issue_certificate(
            subject_cn="rsa-test.snisid.ht",
            key_algorithm=KeyAlgorithm.RSA_2048,
        )
        assert cert.key_algorithm == KeyAlgorithm.RSA_2048
        assert cert.certificate_pem is not None

    def test_revoke_certificate(self):
        ca = InternalCA()
        cert = ca.issue_certificate(subject_cn="revoke-test.snisid.ht")
        assert ca.revoke_certificate(cert.serial_number, "key compromise") is True
        assert ca.check_status(cert.serial_number) == CertificateStatus.REVOKED
        assert cert.revocation_reason == "key compromise"

    def test_revoke_nonexistent(self):
        ca = InternalCA()
        assert ca.revoke_certificate("99999") is False

    def test_double_revoke(self):
        ca = InternalCA()
        cert = ca.issue_certificate(subject_cn="double-revoke.snisid.ht")
        assert ca.revoke_certificate(cert.serial_number) is True
        assert ca.revoke_certificate(cert.serial_number) is False

    def test_check_status_active(self):
        ca = InternalCA()
        cert = ca.issue_certificate(subject_cn="active-test.snisid.ht")
        assert ca.check_status(cert.serial_number) == CertificateStatus.ACTIVE

    def test_check_status_revoked(self):
        ca = InternalCA()
        cert = ca.issue_certificate(subject_cn="revoked-test.snisid.ht")
        ca.revoke_certificate(cert.serial_number)
        assert ca.check_status(cert.serial_number) == CertificateStatus.REVOKED

    def test_get_certificate(self):
        ca = InternalCA()
        cert = ca.issue_certificate(subject_cn="get-test.snisid.ht")
        found = ca.get_certificate(cert.serial_number)
        assert found is not None
        assert found.subject == "get-test.snisid.ht"

    def test_get_nonexistent_certificate(self):
        ca = InternalCA()
        assert ca.get_certificate("nonexistent") is None

    def test_list_revoked(self):
        ca = InternalCA()
        cert1 = ca.issue_certificate(subject_cn="rev-1.snisid.ht")
        cert2 = ca.issue_certificate(subject_cn="rev-2.snisid.ht")
        ca.revoke_certificate(cert1.serial_number, "reason A")
        ca.revoke_certificate(cert2.serial_number, "reason B")
        revoked = ca.list_revoked()
        assert len(revoked) == 2
        reasons = {r["reason"] for r in revoked}
        assert "reason A" in reasons
        assert "reason B" in reasons

    def test_jwt_signing_key(self):
        ca = InternalCA()
        key, cert = ca.get_jwt_signing_key()
        cn = cert.subject.get_attributes_for_oid(x509.oid.NameOID.COMMON_NAME)[0].value
        assert cn == "JWT Signing Key"
        issuer_cn = cert.issuer.get_attributes_for_oid(x509.oid.NameOID.COMMON_NAME)[0].value
        assert issuer_cn == "CN=SNISID Internal CA, O=SNISID, C=HT"

    def test_ecc_p384_ca(self):
        ca = InternalCA(key_algorithm=KeyAlgorithm.ECDSA_P384)
        info = ca.get_ca_info()
        assert info.key_algorithm == KeyAlgorithm.ECDSA_P384

    def test_issue_multiple_certificates(self):
        ca = InternalCA()
        certs = [ca.issue_certificate(subject_cn=f"test-{i}.snisid.ht") for i in range(5)]
        assert len(certs) == 5
        serials = {c.serial_number for c in certs}
        assert len(serials) == 5


class TestPKIApi:
    def setup_method(self):
        self.ca = InternalCA()
        app = FastAPI()
        router = create_pki_router(ca=self.ca)
        app.include_router(router)
        self.client = TestClient(app)

    def test_get_ca_info(self):
        resp = self.client.get("/v1/pki/ca")
        assert resp.status_code == 200
        data = resp.json()
        assert data["ca_subject"] == "CN=SNISID Internal CA, O=SNISID, C=HT"
        assert data["ca_cert_pem"].startswith("-----BEGIN CERTIFICATE-----")

    def test_issue_certificate_api(self):
        resp = self.client.post("/v1/pki/issue", json={
            "subject_cn": "api-test.snisid.ht",
            "validity_days": 30,
        })
        assert resp.status_code == 200
        data = resp.json()
        assert data["subject"] == "api-test.snisid.ht"
        assert data["certificate_pem"].startswith("-----BEGIN CERTIFICATE-----")

    def test_revoke_certificate_api(self):
        issue_resp = self.client.post("/v1/pki/issue", json={
            "subject_cn": "revoke-api.snisid.ht",
        })
        serial = issue_resp.json()["serial_number"]
        revoke_resp = self.client.post("/v1/pki/revoke", json={
            "serial_number": serial,
            "reason": "compromised",
        })
        assert revoke_resp.status_code == 200
        assert revoke_resp.json()["status"] == "revoked"

        status_resp = self.client.get(f"/v1/pki/status/{serial}")
        assert status_resp.status_code == 200
        assert status_resp.json()["status"] == "revoked"

    def test_status_active(self):
        issue_resp = self.client.post("/v1/pki/issue", json={
            "subject_cn": "active-api.snisid.ht",
        })
        serial = issue_resp.json()["serial_number"]
        resp = self.client.get(f"/v1/pki/status/{serial}")
        assert resp.json()["status"] == "active"

    def test_status_not_found(self):
        resp = self.client.get("/v1/pki/status/99999")
        assert resp.status_code == 404

    def test_list_revoked_api(self):
        r1 = self.client.post("/v1/pki/issue", json={"subject_cn": "lr-1.snisid.ht"}).json()
        r2 = self.client.post("/v1/pki/issue", json={"subject_cn": "lr-2.snisid.ht"}).json()
        self.client.post("/v1/pki/revoke", json={"serial_number": r1["serial_number"]})
        self.client.post("/v1/pki/revoke", json={"serial_number": r2["serial_number"]})
        resp = self.client.get("/v1/pki/revoked")
        assert len(resp.json()) == 2

    def test_issue_with_sans_api(self):
        resp = self.client.post("/v1/pki/issue", json={
            "subject_cn": "san-test.snisid.ht",
            "subject_alt_names": ["san-test.snisid.ht", "127.0.0.1"],
        })
        assert resp.status_code == 200
