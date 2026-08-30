from __future__ import annotations

from datetime import datetime, timezone
from typing import Any

from services.vc import VerifiableCredential, VerifiablePresentation, VCStatus


class VerificationResult:
    """Result of a VC verification."""

    def __init__(
        self,
        valid: bool,
        vc_id: str = "",
        errors: list[str] | None = None,
    ):
        self.valid = valid
        self.vc_id = vc_id
        self.errors = errors or []

    def __bool__(self) -> bool:
        return self.valid


class VCVerifier:
    """Verifies Verifiable Credentials and Presentations."""

    def __init__(self, trusted_issuers: list[str] | None = None):
        self._trusted_issuers = set(trusted_issuers or [])
        self._status_resolver: dict[str, VCStatus] = {}

    def register_status(self, vc_id: str, status: VCStatus) -> None:
        self._status_resolver[vc_id] = status

    def register_trusted_issuer(self, issuer_id: str) -> None:
        self._trusted_issuers.add(issuer_id)

    def verify_credential(self, vc: VerifiableCredential) -> VerificationResult:
        errors: list[str] = []

        if not vc.id:
            errors.append("Missing credential ID")

        if not vc.issuer:
            errors.append("Missing issuer")
        elif self._trusted_issuers and vc.issuer not in self._trusted_issuers:
            errors.append(f"Issuer not trusted: {vc.issuer}")

        if not vc.credentialSubject or not vc.credentialSubject.id:
            errors.append("Missing or empty credential subject")

        if not vc.proof:
            errors.append("Missing proof")

        status = self._status_resolver.get(vc.id)
        if status == VCStatus.REVOKED:
            errors.append("Credential has been revoked")
        elif status == VCStatus.SUSPENDED:
            errors.append("Credential is suspended")

        if vc.expirationDate:
            try:
                exp = datetime.fromisoformat(vc.expirationDate)
                if exp < datetime.now(timezone.utc):
                    errors.append("Credential has expired")
            except (ValueError, TypeError):
                errors.append("Invalid expiration date")

        return VerificationResult(
            valid=len(errors) == 0,
            vc_id=vc.id,
            errors=errors,
        )

    def verify_presentation(
        self,
        presentation: VerifiablePresentation,
    ) -> list[VerificationResult]:
        if not presentation.verifiableCredential:
            return [
                VerificationResult(
                    valid=False,
                    errors=["Presentation contains no credentials"],
                )
            ]

        results = []
        for vc in presentation.verifiableCredential:
            result = self.verify_credential(vc)
            results.append(result)
        return results

    def extract_subject_data(
        self, vc: VerifiableCredential
    ) -> dict[str, Any]:
        """Extract claims from a verified credential."""
        subject = vc.credentialSubject
        data = {"sub": subject.id, "vc_id": vc.id, "issuer": vc.issuer}
        extra = getattr(subject, "additional", None) or {}
        if isinstance(extra, dict):
            data.update(extra)
        return data
