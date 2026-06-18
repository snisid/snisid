from __future__ import annotations

from typing import Any

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from services.pki import CertificateAuthorityInfo, KeyAlgorithm
from services.pki.ca import InternalCA


class IssueCertRequest(BaseModel):
    subject_cn: str
    validity_days: int = 365
    key_algorithm: str = "ECDSA-P256"
    subject_alt_names: list[str] | None = None


class RevokeCertRequest(BaseModel):
    serial_number: str
    reason: str = "unspecified"


def create_pki_router(ca: InternalCA | None = None) -> APIRouter:
    router = APIRouter(prefix="/v1/pki", tags=["pki"])
    _ca = ca or InternalCA()

    @router.get("/ca", response_model=CertificateAuthorityInfo)
    async def get_ca_info():
        return _ca.get_ca_info()

    @router.post("/issue", response_model=dict[str, Any])
    async def issue_certificate(req: IssueCertRequest):
        try:
            alg = KeyAlgorithm(req.key_algorithm)
        except ValueError:
            raise HTTPException(400, f"Unsupported algorithm: {req.key_algorithm}")
        cert = _ca.issue_certificate(
            subject_cn=req.subject_cn,
            validity_days=req.validity_days,
            key_algorithm=alg,
            subject_alt_names=req.subject_alt_names,
        )
        return {
            "serial_number": cert.serial_number,
            "subject": cert.subject,
            "certificate_pem": cert.certificate_pem,
            "not_after": cert.not_after.isoformat(),
            "fingerprint_sha256": cert.fingerprint_sha256,
        }

    @router.post("/revoke", response_model=dict[str, Any])
    async def revoke_certificate(req: RevokeCertRequest):
        ok = _ca.revoke_certificate(req.serial_number, req.reason)
        if not ok:
            raise HTTPException(404, "Certificate not found or already revoked")
        return {"serial_number": req.serial_number, "status": "revoked"}

    @router.get("/status/{serial}", response_model=dict[str, Any])
    async def check_status(serial: str):
        cert = _ca.get_certificate(serial)
        if not cert:
            raise HTTPException(404, "Certificate not found")
        return {
            "serial_number": serial,
            "status": _ca.check_status(serial).value,
            "subject": cert.subject,
            "not_before": cert.not_before.isoformat(),
            "not_after": cert.not_after.isoformat(),
        }

    @router.get("/revoked", response_model=list[dict[str, Any]])
    async def list_revoked():
        return _ca.list_revoked()

    return router
