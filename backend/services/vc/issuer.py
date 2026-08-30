from __future__ import annotations

import base64
import hashlib
import hmac
import json
import uuid
from datetime import datetime, timedelta, timezone
from typing import Any

from cryptography import x509
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec, rsa

from services.pki.ca import InternalCA
from services.status_list import StatusListManager
from services.vc import (
    IdentityCredential,
    IdentityCredentialSubject,
    VCStatus,
    VerifiableCredential,
)


class VCIssuer:
    """Issues and manages Verifiable Credentials with StatusList2021 + PKI."""

    def __init__(
        self,
        issuer_id: str,
        signing_key: str | None = None,
        ca: InternalCA | None = None,
        vc_storage: Any | None = None,
        status_storage: Any | None = None,
    ):
        self._issuer_id = issuer_id
        self._signing_key = signing_key or "dev-signing-key"
        self._ca = ca
        self._ca_pub_key_pem: str | None = None
        self._status_manager = StatusListManager(issuer_id)
        self._vc_entry_map: dict[str, str] = {}
        self._vc_storage = vc_storage
        self._status_storage = status_storage

    def issue_identity_credential(
        self,
        subject_id: str,
        national_id: str,
        first_name: str,
        last_name: str,
        date_of_birth: str,
        gender: str,
        nationality: str,
        status: str = "active",
        expiry_days: int = 365,
    ) -> IdentityCredential:
        vc_id = f"urn:uuid:{uuid.uuid4()}"
        now = datetime.now(timezone.utc)
        subject = IdentityCredentialSubject(
            id=subject_id,
            national_id=national_id,
            first_name=first_name,
            last_name=last_name,
            date_of_birth=date_of_birth,
            gender=gender,
            nationality=nationality,
            status=status,
        )

        sl_entry = self._status_manager.create_entry()
        vc = IdentityCredential(
            id=vc_id,
            issuer=self._issuer_id,
            issuanceDate=now.isoformat(),
            expirationDate=(now + timedelta(days=expiry_days)).isoformat(),
            credentialSubject=subject,
            credentialStatus=sl_entry.to_dict(),
        )

        vc.proof = self._sign(vc)
        self._vc_entry_map[vc_id] = sl_entry.id
        return vc

    async def async_issue_identity_credential(
        self,
        subject_id: str,
        national_id: str,
        first_name: str,
        last_name: str,
        date_of_birth: str,
        gender: str,
        nationality: str,
        status: str = "active",
        expiry_days: int = 365,
    ) -> IdentityCredential:
        vc_id = f"urn:uuid:{uuid.uuid4()}"
        now = datetime.now(timezone.utc)
        subject = IdentityCredentialSubject(
            id=subject_id,
            national_id=national_id,
            first_name=first_name,
            last_name=last_name,
            date_of_birth=date_of_birth,
            gender=gender,
            nationality=nationality,
            status=status,
        )

        sl_entry = self._status_manager.create_entry()
        vc = IdentityCredential(
            id=vc_id,
            issuer=self._issuer_id,
            issuanceDate=now.isoformat(),
            expirationDate=(now + timedelta(days=expiry_days)).isoformat(),
            credentialSubject=subject,
            credentialStatus=sl_entry.to_dict(),
        )

        vc.proof = self._sign(vc)
        self._vc_entry_map[vc_id] = sl_entry.id

        if self._vc_storage:
            await self._vc_storage.save(
                credential_id=vc.id,
                issuer_id=self._issuer_id,
                subject_id=subject_id,
                credential_type="SNISIDIdentityCredential",
                document=vc.model_dump(),
                status_list_id=sl_entry.id,
            )
        if self._status_storage:
            sl_vc = self._status_manager.get_status_list_credential(
                sl_entry.statusListCredential
            )
            if sl_vc:
                await self._status_storage.save(
                    list_id=sl_entry.statusListCredential,
                    purpose="revocation",
                    bitstring=sl_vc.get("credentialSubject", {}).get("encodedList", ""),
                )

        return vc

    def get_credential_status(self, vc_id: str) -> VCStatus | None:
        entry_id = self._vc_entry_map.get(vc_id)
        if entry_id is None:
            return None
        revoked = self._status_manager.is_revoked(entry_id)
        if revoked is True:
            return VCStatus.REVOKED
        return VCStatus.ACTIVE

    async def async_get_credential_status(self, vc_id: str) -> VCStatus | None:
        return self.get_credential_status(vc_id)

    def revoke_credential(self, vc_id: str) -> bool:
        entry_id = self._vc_entry_map.get(vc_id)
        if entry_id is None:
            return False
        return self._status_manager.revoke(entry_id)

    async def async_revoke_credential(self, vc_id: str) -> bool:
        entry_id = self._vc_entry_map.get(vc_id)
        if entry_id is None:
            return False
        ok = self._status_manager.revoke(entry_id)
        if ok and self._status_storage:
            sl = self._status_manager.get_status_list_credential(
                self._vc_entry_map[vc_id].rsplit("#", 1)[0]
            )
            if sl:
                encoded = sl.get("credentialSubject", {}).get("encodedList", "")
                await self._status_storage.update_bitstring(
                    self._vc_entry_map[vc_id].rsplit("#", 1)[0], encoded
                )
        return ok

    def suspend_credential(self, vc_id: str) -> bool:
        entry_id = self._vc_entry_map.get(vc_id)
        if entry_id is None:
            return False
        return self._status_manager.revoke(entry_id)

    def _sign(self, vc: VerifiableCredential) -> dict[str, Any]:
        payload = vc.model_dump_json(exclude={"proof"})
        now = datetime.now(timezone.utc).isoformat()

        if self._ca:
            sig_bytes, alg_name, thumbprint = self._ca.sign_data(payload.encode())
            proof_value = base64.urlsafe_b64encode(sig_bytes).rstrip(b"=").decode()
            return {
                "type": f"SNISID-{alg_name}-2026",
                "created": now,
                "verificationMethod": f"{self._issuer_id}#{thumbprint}",
                "proofPurpose": "assertionMethod",
                "proofValue": proof_value,
            }

        signature = hmac.new(
            self._signing_key.encode(),
            payload.encode(),
            hashlib.sha256,
        ).hexdigest()
        return {
            "type": "SNISID-HMAC-SHA256-2026",
            "created": now,
            "verificationMethod": f"{self._issuer_id}#key-1",
            "proofPurpose": "assertionMethod",
            "proofValue": signature,
        }

    @staticmethod
    def _is_expired(entry: dict[str, Any]) -> bool:
        issuance = entry.get("issuanceDate", "")
        if not issuance:
            return False
        try:
            issued = datetime.fromisoformat(issuance)
            return datetime.now(timezone.utc) > issued + timedelta(days=365)
        except (ValueError, TypeError):
            return False

    def verify_signature(self, vc: VerifiableCredential) -> bool:
        if not vc.proof:
            return False
        if self._ca:
            if self._ca_pub_key_pem is None:
                ca_info = self._ca.get_ca_info()
                ca_cert = x509.load_pem_x509_certificate(ca_info.ca_cert_pem.encode())
                self._ca_pub_key_pem = ca_cert.public_key().public_bytes(
                    encoding=serialization.Encoding.PEM,
                    format=serialization.PublicFormat.SubjectPublicKeyInfo,
                ).decode()
            proof_value = vc.proof.get("proofValue", "")
            try:
                sig_bytes = base64.urlsafe_b64decode(proof_value + "==")
            except Exception:
                return False
            payload = vc.model_dump_json(exclude={"proof"})
            return InternalCA.verify_data_signature(
                payload.encode(), sig_bytes, self._ca_pub_key_pem
            )
        expected = self._sign(vc)
        return vc.proof.get("proofValue") == expected["proofValue"]
