from __future__ import annotations

from typing import Any

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from services.sd_jwt import SDJWTBuilder, SDJWTIssuer, SDJWTVerifier

router = APIRouter(prefix="/sd-jwt", tags=["sd-jwt"])

_issuer = SDJWTIssuer(issuer_id="https://snisid.ht/issuer")
_verifier = SDJWTVerifier(trusted_issuers={"https://snisid.ht/issuer"})


class IssueRequest(BaseModel):
    subject: str
    disclosed_claims: dict[str, Any] = {}
    sd_claims: dict[str, Any] = {}
    expiration_seconds: int = 3600


class IssueResponse(BaseModel):
    sd_jwt: str
    disclosures: list[str]


class VerifyRequest(BaseModel):
    sd_jwt: str
    disclosures: list[str] = []
    required_claims: list[str] | None = None


class VerifyResponse(BaseModel):
    claims: dict[str, Any]


class PresentRequest(BaseModel):
    sd_jwt: str
    all_disclosures: list[str]
    disclose: list[str]


class PresentResponse(BaseModel):
    sd_jwt: str
    disclosures: list[str]


@router.post("/issue", response_model=IssueResponse)
async def issue_sd_jwt(req: IssueRequest):
    sd_jwt, disclosures = _issuer.issue(
        subject=req.subject,
        disclosed_claims=req.disclosed_claims,
        sd_claims=req.sd_claims,
        expiration_seconds=req.expiration_seconds,
    )
    return IssueResponse(sd_jwt=sd_jwt, disclosures=disclosures)


@router.post("/verify", response_model=VerifyResponse)
async def verify_sd_jwt(req: VerifyRequest):
    try:
        claims = _verifier.verify(
            req.sd_jwt,
            disclosures=req.disclosures,
            required_claims=req.required_claims,
        )
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    return VerifyResponse(claims=claims)


@router.post("/present", response_model=PresentResponse)
async def present_sd_jwt(req: PresentRequest):
    sd_jwt, disclosures = SDJWTBuilder.create_presentation(
        req.sd_jwt,
        all_disclosures=req.all_disclosures,
        disclose=req.disclose,
    )
    return PresentResponse(sd_jwt=sd_jwt, disclosures=disclosures)
