from __future__ import annotations

from typing import Any

from fastapi import APIRouter, HTTPException

from services.vc import IdentityCredentialSubject, VCStatus, VerifiablePresentation
from services.vc.issuer import VCIssuer
from services.vc.verifier import VCVerifier, VerificationResult


def create_vc_router(
    issuer: VCIssuer | None = None,
    verifier: VCVerifier | None = None,
) -> APIRouter:
    router = APIRouter(prefix="/v1/vc", tags=["verifiable-credentials"])
    _issuer = issuer or VCIssuer(issuer_id="http://localhost:8000")
    _verifier = verifier or VCVerifier()

    @router.post("/issue/identity", response_model=dict[str, Any])
    async def issue_identity_credential(body: IdentityCredentialSubject):
        vc = await _issuer.async_issue_identity_credential(
            subject_id=body.id,
            national_id=body.national_id,
            first_name=body.first_name,
            last_name=body.last_name,
            date_of_birth=body.date_of_birth,
            gender=body.gender,
            nationality=body.nationality,
            status=body.status,
        )
        _verifier.register_status(vc.id, VCStatus.ACTIVE)
        return vc.model_dump(mode="json")

    @router.post("/verify", response_model=dict[str, Any])
    async def verify_credential(body: dict[str, Any]):
        try:
            vc = VerifiablePresentation(**body) if "holder" in body else None
            if vc and vc.verifiableCredential:
                results = _verifier.verify_presentation(vc)
                return {
                    "verified": all(r.valid for r in results),
                    "results": [r.__dict__ for r in results],
                }
            from services.vc import VerifiableCredential
            vc_single = VerifiableCredential(**body)
            result = _verifier.verify_credential(vc_single)
            return {
                "verified": result.valid,
                "vc_id": result.vc_id,
                "errors": result.errors,
            }
        except Exception as e:
            raise HTTPException(status_code=400, detail=f"Invalid credential format: {e}")

    @router.get("/status", response_model=dict[str, Any])
    async def get_credential_status(vc_id: str):
        status = _issuer.get_credential_status(vc_id) if _issuer else None
        if status is None:
            raise HTTPException(status_code=404, detail="Credential not found")
        return {"vc_id": vc_id, "status": status.value}

    @router.post("/revoke", response_model=dict[str, Any])
    async def revoke_credential(vc_id: str):
        if not _issuer:
            raise HTTPException(status_code=503, detail="Issuer not available")
        ok = await _issuer.async_revoke_credential(vc_id)
        if not ok:
            raise HTTPException(status_code=404, detail="Credential not found")
        _verifier.register_status(vc_id, VCStatus.REVOKED)
        return {"vc_id": vc_id, "status": "revoked"}

    return router
