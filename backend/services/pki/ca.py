from __future__ import annotations

import datetime
import hashlib
import uuid
from typing import Any

from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import ec, padding, rsa
from cryptography.x509.oid import NameOID

from services.pki import (
    Certificate,
    CertificateAuthorityInfo,
    CertificateStatus,
    KeyAlgorithm,
)


class InternalCA:
    """Simple internal Certificate Authority for development/testing."""

    def __init__(
        self,
        ca_subject: str = "CN=SNISID Internal CA, O=SNISID, C=HT",
        key_algorithm: KeyAlgorithm = KeyAlgorithm.ECDSA_P256,
        validity_years: int = 10,
    ):
        self._ca_subject_str = ca_subject
        self._key_algorithm = key_algorithm
        self._validity_years = validity_years
        self._private_key, self._certificate = self._generate_ca()
        self._revoked: dict[str, tuple[datetime.datetime, str]] = {}
        self._issued: dict[str, Certificate] = {}

    def _generate_ca(self) -> tuple[Any, x509.Certificate]:
        key = self._generate_key()
        subject = issuer = x509.Name([x509.NameAttribute(NameOID.COMMON_NAME, self._ca_subject_str)])
        now = datetime.datetime.now(datetime.timezone.utc)
        cert = (
            x509.CertificateBuilder()
            .subject_name(subject)
            .issuer_name(issuer)
            .public_key(key.public_key())
            .serial_number(int(uuid.uuid4().int % (2**63)))
            .not_valid_before(now)
            .not_valid_after(now + datetime.timedelta(days=365 * self._validity_years))
            .add_extension(
                x509.BasicConstraints(ca=True, path_length=None), critical=True
            )
            .add_extension(
                x509.KeyUsage(
                    digital_signature=True,
                    key_cert_sign=True,
                    key_encipherment=False,
                    content_commitment=False,
                    data_encipherment=False,
                    crl_sign=True,
                    encipher_only=False,
                    decipher_only=False,
                    key_agreement=False,
                ),
                critical=True,
            )
            .sign(key, self._hash_algorithm())
        )
        return key, cert

    def _generate_key(self) -> Any:
        if self._key_algorithm in (KeyAlgorithm.ECDSA_P256, KeyAlgorithm.ECDSA_P384):
            curve = ec.SECP256R1() if self._key_algorithm == KeyAlgorithm.ECDSA_P256 else ec.SECP384R1()
            return ec.generate_private_key(curve)
        size = 2048 if self._key_algorithm == KeyAlgorithm.RSA_2048 else 4096
        return rsa.generate_private_key(public_exponent=65537, key_size=size)

    def _hash_algorithm(self) -> hashes.HashAlgorithm:
        return hashes.SHA256() if "SHA256" in str(self._key_algorithm) else hashes.SHA384()

    def issue_certificate(
        self,
        subject_cn: str,
        validity_days: int = 365,
        key_algorithm: KeyAlgorithm = KeyAlgorithm.ECDSA_P256,
        subject_alt_names: list[str] | None = None,
    ) -> Certificate:
        subject_key = self._generate_key()
        subject = x509.Name([x509.NameAttribute(NameOID.COMMON_NAME, subject_cn)])
        now = datetime.datetime.now(datetime.timezone.utc)
        serial = uuid.uuid4().int % (2**63)

        builder = (
            x509.CertificateBuilder()
            .subject_name(subject)
            .issuer_name(self._certificate.subject)
            .public_key(subject_key.public_key())
            .serial_number(serial)
            .not_valid_before(now)
            .not_valid_after(now + datetime.timedelta(days=validity_days))
            .add_extension(
                x509.BasicConstraints(ca=False, path_length=None), critical=True
            )
        )

        if subject_alt_names:
            builder = builder.add_extension(
                x509.SubjectAlternativeName(
                    [x509.DNSName(san) for san in subject_alt_names]
                ),
                critical=False,
            )

        cert = builder.sign(self._private_key, self._hash_algorithm())
        fingerprint = cert.fingerprint(hashes.SHA256()).hex()

        result = Certificate(
            serial_number=str(serial),
            subject=subject_cn,
            issuer=self._ca_subject_str,
            not_before=now,
            not_after=now + datetime.timedelta(days=validity_days),
            public_key_pem=subject_key.public_key().public_bytes(
                encoding=serialization.Encoding.PEM,
                format=serialization.PublicFormat.SubjectPublicKeyInfo,
            ).decode(),
            certificate_pem=cert.public_bytes(serialization.Encoding.PEM).decode(),
            fingerprint_sha256=fingerprint,
            key_algorithm=key_algorithm,
            subject_alt_names=subject_alt_names or [],
        )
        self._issued[result.serial_number] = result
        return result

    def revoke_certificate(self, serial_number: str, reason: str = "unspecified") -> bool:
        cert = self._issued.get(serial_number)
        if not cert or cert.status == CertificateStatus.REVOKED:
            return False
        cert.status = CertificateStatus.REVOKED
        cert.revocation_date = datetime.datetime.now(datetime.timezone.utc)
        cert.revocation_reason = reason
        self._revoked[serial_number] = (cert.revocation_date, reason)
        return True

    def get_certificate(self, serial_number: str) -> Certificate | None:
        return self._issued.get(serial_number)

    def check_status(self, serial_number: str) -> CertificateStatus:
        cert = self._issued.get(serial_number)
        if not cert:
            return CertificateStatus.REVOKED
        if cert.status == CertificateStatus.REVOKED:
            return CertificateStatus.REVOKED
        if datetime.datetime.now(datetime.timezone.utc) > cert.not_after:
            return CertificateStatus.EXPIRED
        return CertificateStatus.ACTIVE

    def get_ca_info(self) -> CertificateAuthorityInfo:
        fp = self._certificate.fingerprint(hashes.SHA256()).hex()
        return CertificateAuthorityInfo(
            ca_cert_pem=self._certificate.public_bytes(
                serialization.Encoding.PEM
            ).decode(),
            ca_subject=self._ca_subject_str,
            ca_serial=str(self._certificate.serial_number),
            ca_fingerprint=fp,
            not_after=self._certificate.not_valid_after_utc,
            key_algorithm=self._key_algorithm,
        )

    def list_revoked(self) -> list[dict[str, Any]]:
        return [
            {"serial": sn, "date": d.isoformat(), "reason": r}
            for sn, (d, r) in self._revoked.items()
        ]

    def get_jwt_signing_key(self) -> tuple[Any, x509.Certificate]:
        key = self._generate_key()
        subject = x509.Name([x509.NameAttribute(NameOID.COMMON_NAME, "JWT Signing Key")])
        now = datetime.datetime.now(datetime.timezone.utc)
        cert = (
            x509.CertificateBuilder()
            .subject_name(subject)
            .issuer_name(self._certificate.subject)
            .public_key(key.public_key())
            .serial_number(int(uuid.uuid4().int % (2**63)))
            .not_valid_before(now)
            .not_valid_after(now + datetime.timedelta(days=365))
            .sign(self._private_key, self._hash_algorithm())
        )
        return key, cert

    def sign_data(self, data: bytes) -> tuple[bytes, str, str]:
        """Sign arbitrary data with the CA key.
        Returns (signature_bytes, key_algorithm_name, thumbprint)."""
        algo = hashes.SHA256()
        if hasattr(self._private_key, 'curve'):
            signature = self._private_key.sign(data, ec.ECDSA(algo))
        else:
            signature = self._private_key.sign(data, padding.PKCS1v15(), algo)
        pub_bytes = self._private_key.public_key().public_bytes(
            encoding=serialization.Encoding.DER,
            format=serialization.PublicFormat.SubjectPublicKeyInfo,
        )
        thumbprint = hashlib.sha256(pub_bytes).hexdigest()[:16]
        alg_name = "ES256" if hasattr(self._private_key, 'curve') else "RS256"
        return signature, alg_name, thumbprint

    @staticmethod
    def verify_data_signature(data: bytes, signature: bytes, public_key_pem: str) -> bool:
        """Verify a signature against a PEM-encoded public key."""
        try:
            pub_key = serialization.load_pem_public_key(public_key_pem.encode())
            algo = hashes.SHA256()
            if hasattr(pub_key, 'curve'):
                pub_key.verify(signature, data, ec.ECDSA(algo))
            else:
                pub_key.verify(signature, data, padding.PKCS1v15(), algo)
            return True
        except Exception:
            return False
