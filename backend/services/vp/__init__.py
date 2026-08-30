from __future__ import annotations

import hashlib
import json
from datetime import datetime, timezone
from typing import Any

import base64

from services.did import VerificationMethodType


def _b64(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode()


VP_CONTEXT = "https://www.w3.org/ns/credentials/v2"


class VerifiablePresentation:
    def __init__(
        self,
        holder_did: str,
        verifiable_credential: list[dict[str, Any]],
        proof: dict[str, Any] | None = None,
        id: str | None = None,
        context: list[str] | None = None,
        type: list[str] | None = None,
    ):
        self.context = context or [VP_CONTEXT]
        self.id = id or f"urn:uuid:{__import__('uuid').uuid4()}"
        self.type = type or ["VerifiablePresentation"]
        self.holder = holder_did
        self.verifiable_credential = verifiable_credential
        self.proof = proof

    def to_dict(self) -> dict[str, Any]:
        d = {
            "@context": self.context,
            "id": self.id,
            "type": self.type,
            "holder": self.holder,
            "verifiableCredential": self.verifiable_credential,
        }
        if self.proof:
            d["proof"] = self.proof
        return d

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> VerifiablePresentation:
        return cls(
            holder_did=data.get("holder", ""),
            verifiable_credential=data.get("verifiableCredential", []),
            proof=data.get("proof"),
            id=data.get("id"),
            context=data.get("@context"),
            type=data.get("type"),
        )


class VPIssuer:
    def __init__(self, proof_key: str = "dev-vp-key"):
        self._key = proof_key

    def create_presentation(
        self,
        holder_did: str,
        verifiable_credentials: list[dict[str, Any]],
        verification_method: str | None = None,
    ) -> VerifiablePresentation:
        vp = VerifiablePresentation(
            holder_did=holder_did,
            verifiable_credential=verifiable_credentials,
        )

        proof = self._create_proof(vp, holder_did, verification_method)
        vp.proof = proof
        return vp

    def _create_proof(
        self,
        vp: VerifiablePresentation,
        holder_did: str,
        verification_method: str | None = None,
    ) -> dict[str, Any]:
        vp_data = vp.to_dict()
        vp_data.pop("proof", None)
        canonical = json.dumps(vp_data, separators=(",", ":"), sort_keys=True)

        now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        vm = verification_method or f"{holder_did}#key-1"
        proof_config = {
            "type": "SNISID-HMAC-SHA256-2025",
            "created": now,
            "verificationMethod": vm,
            "proofPurpose": "authentication",
        }
        proof_payload = json.dumps(proof_config, separators=(",", ":"), sort_keys=True)
        signature_input = f"{proof_payload}{canonical}{self._key}"
        proof_value = _b64(hashlib.sha256(signature_input.encode()).digest())

        return {
            **proof_config,
            "proofValue": proof_value,
        }

    def verify_presentation(self, vp: VerifiablePresentation) -> bool:
        if not vp.proof:
            return False

        proof = vp.proof
        expected_proof_value = proof.get("proofValue", "")
        vp_data = vp.to_dict()
        vp_data.pop("proof", None)
        canonical = json.dumps(vp_data, separators=(",", ":"), sort_keys=True)

        proof_config = {k: v for k, v in proof.items() if k != "proofValue"}
        proof_payload = json.dumps(proof_config, separators=(",", ":"), sort_keys=True)
        signature_input = f"{proof_payload}{canonical}{self._key}"
        expected = _b64(hashlib.sha256(signature_input.encode()).digest())

        return expected == expected_proof_value
