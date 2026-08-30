from __future__ import annotations

import hashlib
import json
from datetime import datetime, timezone
from typing import Any

from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.asymmetric.utils import decode_dss_signature, encode_dss_signature


def generate_audit_keypair() -> tuple[bytes, bytes]:
    private_key = ec.generate_private_key(ec.SECP256R1())
    pub = private_key.public_key()
    priv_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.NoEncryption(),
    )
    pub_pem = pub.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )
    return priv_pem, pub_pem


class SignedAuditEntry:
    """Immutable audit entry signed with ECDSA-P256 (SHA-256 digest)."""

    def __init__(
        self,
        event_type: str,
        table_name: str,
        record_id: str | None = None,
        officer_niu: str | None = None,
        agency_code: str | None = None,
        purpose: str | None = None,
        case_number: str | None = None,
        action: str | None = None,
        details: dict[str, Any] | None = None,
        sample_id: str | None = None,
    ) -> None:
        self.event_type = event_type
        self.table_name = table_name
        self.record_id = record_id
        self.officer_niu = officer_niu
        self.agency_code = agency_code
        self.purpose = purpose
        self.case_number = case_number
        self.action = action
        self.details = details or {}
        self.sample_id = sample_id
        self.timestamp = datetime.now(timezone.utc).isoformat()
        self.signature: str | None = None

    def _canonical(self) -> bytes:
        parts = "|".join([
            self.event_type or "",
            self.table_name or "",
            self.record_id or "",
            self.sample_id or "",
            self.officer_niu or "",
            self.agency_code or "",
            self.purpose or "",
            self.case_number or "",
            self.action or "",
            self.timestamp,
        ])
        return parts.encode("utf-8")

    def sign(self, private_key_pem: bytes) -> str:
        private_key = serialization.load_pem_private_key(private_key_pem, password=None)
        if not isinstance(private_key, ec.EllipticCurvePrivateKey):
            raise TypeError("Expected EC private key")

        digest = hashlib.sha256(self._canonical()).digest()
        signature = private_key.sign(digest, ec.ECDSA(hashes.SHA256()))
        der_sig = decode_dss_signature(signature)
        sig_bytes = der_sig[0].to_bytes(32, "big") + der_sig[1].to_bytes(32, "big")
        self.signature = sig_bytes.hex()
        return self.signature

    def verify(self, public_key_pem: bytes) -> bool:
        if not self.signature:
            return False
        public_key = serialization.load_pem_public_key(public_key_pem)
        if not isinstance(public_key, ec.EllipticCurvePublicKey):
            raise TypeError("Expected EC public key")

        digest = hashlib.sha256(self._canonical()).digest()
        sig_bytes = bytes.fromhex(self.signature)
        r = int.from_bytes(sig_bytes[:32], "big")
        s = int.from_bytes(sig_bytes[32:], "big")
        der_signature = encode_dss_signature(r, s)

        try:
            public_key.verify(der_signature, digest, ec.ECDSA(hashes.SHA256()))
            return True
        except InvalidSignature:
            return False

    def to_dict(self) -> dict[str, Any]:
        return {
            "event_type": self.event_type,
            "table_name": self.table_name,
            "record_id": self.record_id,
            "sample_id": self.sample_id,
            "officer_niu": self.officer_niu,
            "agency_code": self.agency_code,
            "purpose": self.purpose,
            "case_number": self.case_number,
            "action": self.action,
            "details": self.details,
            "timestamp": self.timestamp,
            "signature": self.signature,
        }

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> SignedAuditEntry:
        entry = cls(
            event_type=data.get("event_type", ""),
            table_name=data.get("table_name", ""),
            record_id=data.get("record_id"),
            officer_niu=data.get("officer_niu"),
            agency_code=data.get("agency_code"),
            purpose=data.get("purpose"),
            case_number=data.get("case_number"),
            action=data.get("action"),
            details=data.get("details"),
            sample_id=data.get("sample_id"),
        )
        entry.timestamp = data.get("timestamp", entry.timestamp)
        entry.signature = data.get("signature")
        return entry
