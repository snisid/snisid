from __future__ import annotations

import enum
from dataclasses import dataclass, field
from datetime import datetime
from typing import Any


class KeyAlgorithm(str, enum.Enum):
    RSA_2048 = "RSA-2048"
    RSA_4096 = "RSA-4096"
    ECDSA_P256 = "ECDSA-P256"
    ECDSA_P384 = "ECDSA-P384"


class CertificateStatus(str, enum.Enum):
    ACTIVE = "active"
    EXPIRED = "expired"
    REVOKED = "revoked"


@dataclass
class Certificate:
    serial_number: str
    subject: str
    issuer: str
    not_before: datetime
    not_after: datetime
    public_key_pem: str
    certificate_pem: str
    fingerprint_sha256: str
    status: CertificateStatus = CertificateStatus.ACTIVE
    revocation_date: datetime | None = None
    revocation_reason: str | None = None
    key_algorithm: KeyAlgorithm = KeyAlgorithm.RSA_2048
    subject_alt_names: list[str] = field(default_factory=list)


@dataclass
class CertificateSigningRequest:
    csr_pem: str
    subject: str
    key_algorithm: KeyAlgorithm


@dataclass
class CertificateAuthorityInfo:
    ca_cert_pem: str
    ca_subject: str
    ca_serial: str
    ca_fingerprint: str
    not_after: datetime
    key_algorithm: KeyAlgorithm
