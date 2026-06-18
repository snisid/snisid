from __future__ import annotations

import uuid
from typing import Any

from fastapi import APIRouter, HTTPException

from services.credential_flow import CredentialFlow
from services.vc.issuer import VCIssuer
from services.vc.verifier import VCVerifier
from services.vc import VerifiableCredential, VCStatus

router = APIRouter(prefix="/v1/credential-flow", tags=["credential-flow"])

_flow: CredentialFlow | None = None


def get_flow() -> CredentialFlow:
    global _flow
    if _flow is None:
        _flow = CredentialFlow(
            issuer=VCIssuer(issuer_id="did:snisid:mainnet:issuer-system"),
            verifier=VCVerifier(trusted_issuers=["did:snisid:mainnet:issuer-system"]),
        )
    return _flow


@router.post("/offer")
async def create_offer(
    issuer_did: str,
    holder_did: str,
    credential_type: str = "SNISIDIdentityCredential",
):
    flow = get_flow()
    offer = await flow.async_create_offer(
        issuer_did=issuer_did,
        holder_did=holder_did,
        credential_type=credential_type,
    )
    return offer.to_dict()


@router.post("/offer/{offer_id}/send")
def send_offer(offer_id: str):
    flow = get_flow()
    offer = flow._offers.get(offer_id)
    if not offer:
        raise HTTPException(404, "Offer not found")
    packed = flow.send_offer(offer)
    return {"status": "sent", "packed": packed}


@router.post("/request")
def request_credential(
    offer_id: str,
    holder_did: str,
    issuer_did: str,
    claims: dict[str, Any] | None = None,
):
    flow = get_flow()
    body = flow.build_request(offer_id, holder_did, claims)
    packed = flow.pack_and_send_request(body, holder_did, issuer_did)
    return {"status": "requested", "packed": packed}


@router.post("/receive-request")
def receive_request(packed: dict[str, Any]):
    flow = get_flow()
    req = flow.receive_request(packed)
    return req.to_dict()


@router.post("/issue")
async def issue_credential(request_data: dict[str, Any]):
    flow = get_flow()
    from services.credential_flow import CredentialRequest

    req = CredentialRequest(
        request_id=request_data.get("request_id", str(uuid.uuid4())),
        offer_id=request_data["offer_id"],
        holder_did=request_data["holder_did"],
        claims=request_data.get("claims", {}),
    )
    vc_data = await flow.async_issue_from_request(req)
    if vc_data is None:
        raise HTTPException(404, "Offer not found")
    return vc_data


@router.post("/send-credential")
def send_credential(credential: dict[str, Any], issuer_did: str, holder_did: str):
    flow = get_flow()
    packed = flow.send_credential(credential, issuer_did, holder_did)
    return {"status": "sent", "packed": packed}


@router.post("/receive-credential")
def receive_credential(packed: dict[str, Any]):
    flow = get_flow()
    vc = flow.receive_credential(packed)
    if vc is None:
        raise HTTPException(400, "No credential in message")
    return vc
