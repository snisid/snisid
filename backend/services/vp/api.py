from typing import Any

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from services.vp import VPIssuer, VerifiablePresentation

router = APIRouter(prefix="/vp", tags=["vp"])

_issuer = VPIssuer()


class CreatePresentationRequest(BaseModel):
    holder_did: str
    verifiable_credentials: list[dict[str, Any]]
    verification_method: str | None = None


class VerifyPresentationRequest(BaseModel):
    data: dict[str, Any]


@router.post("/create")
async def create_presentation(req: CreatePresentationRequest):
    vp = _issuer.create_presentation(
        holder_did=req.holder_did,
        verifiable_credentials=req.verifiable_credentials,
        verification_method=req.verification_method,
    )
    return vp.to_dict()


@router.post("/verify")
async def verify_presentation(req: VerifyPresentationRequest):
    try:
        vp = VerifiablePresentation.from_dict(req.data)
        valid = _issuer.verify_presentation(vp)
        return {"verified": valid}
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))
