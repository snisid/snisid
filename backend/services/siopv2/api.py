from __future__ import annotations

import uuid
from typing import Any

from fastapi import APIRouter, HTTPException, Query

from services.siopv2 import (
    InputDescriptor,
    OIDC4VPRequest,
    OIDC4VPResponse,
    OIDC4VPVerifier,
    PresentationDefinition,
    SIOPWallet,
)

router = APIRouter(prefix="/v1/siopv2", tags=["siopv2"])

_wallet: SIOPWallet | None = None
_requests: dict[str, OIDC4VPRequest] = {}


def get_wallet() -> SIOPWallet:
    global _wallet
    if _wallet is None:
        _wallet = SIOPWallet()
    return _wallet


@router.post("/wallet/store")
def store_credential(credential: dict[str, Any]):
    wallet = get_wallet()
    wallet.store_credential(credential)
    return {"status": "stored", "count": len(wallet._credentials)}


@router.get("/wallet/did")
def wallet_did():
    return {"did": get_wallet().did}


@router.get("/wallet/credentials")
def wallet_credentials():
    return {"credentials": get_wallet()._credentials}


@router.post("/verifier/request")
def create_request(
    client_id: str = "did:key:verifier",
    scope: str = "openid",
    redirect_uri: str = "openid://",
):
    req = OIDC4VPRequest(
        client_id=client_id,
        scope=scope,
        redirect_uri=redirect_uri,
    )
    _requests[req.nonce] = req
    return {"request": req.to_dict(), "_internal_nonce": req.nonce}


@router.post("/verifier/request-with-pd")
def create_request_with_presentation_definition(
    client_id: str = "did:key:verifier",
    input_descriptor_ids: str = Query("identity-credential"),
    pd_name: str = "Identity Check",
):
    desc = InputDescriptor(
        id=input_descriptor_ids,
        name=pd_name,
        purpose="Verify identity",
    )
    pd = PresentationDefinition(
        id=str(uuid.uuid4()),
        name=pd_name,
        input_descriptors=[desc],
    )
    req = OIDC4VPRequest(
        client_id=client_id,
        presentation_definition=pd,
    )
    _requests[req.nonce] = req
    return {"request": req.to_dict(), "_internal_nonce": req.nonce}


@router.post("/wallet/respond")
def wallet_respond(request_data: dict[str, Any]):
    wallet = get_wallet()
    req = _deserialize_request(request_data)
    response = wallet.respond_to_request(req)
    return response.to_dict()


@router.post("/verifier/verify")
def verify_response(
    response_data: dict[str, Any],
    expected_nonce: str = Query(...),
    expected_state: str | None = Query(None),
):
    verifier = OIDC4VPVerifier()
    response = _deserialize_response(response_data)
    result = verifier.verify_response(
        response,
        expected_nonce=expected_nonce,
        expected_state=expected_state,
    )
    return result


@router.post("/wallet/respond/direct-post")
def wallet_direct_post(request_data: dict[str, Any]):
    """Wallet responds to a request via direct_post (POSTs response to verifier)."""
    wallet = get_wallet()
    req = _deserialize_request(request_data)
    response = wallet.post_response(req)
    return response.to_dict()


@router.post("/verifier/direct-post-callback")
def verifier_direct_post_callback(response_data: dict[str, Any]):
    """Verifier endpoint that receives direct_post responses from wallet."""
    verifier = OIDC4VPVerifier()
    response = _deserialize_response(response_data)
    nonce = response_data.get("_internal_nonce", "")
    result = verifier.verify_response(
        response,
        expected_nonce=nonce,
        expected_state=response.state,
    )
    return result


def _deserialize_request(data: dict[str, Any]) -> OIDC4VPRequest:
    pd_data = data.get("presentation_definition")
    pd = None
    if pd_data:
        descriptors = []
        for d in pd_data.get("input_descriptors", []):
            descriptors.append(
                InputDescriptor(
                    id=d.get("id", ""),
                    name=d.get("name", ""),
                    purpose=d.get("purpose", ""),
                    schema_uris=[s["uri"] for s in d.get("schema", [])] if d.get("schema") else None,
                    constraints=d.get("constraints"),
                )
            )
        pd = PresentationDefinition(
            id=pd_data.get("id", str(uuid.uuid4())),
            name=pd_data.get("name", ""),
            input_descriptors=descriptors,
        )
    return OIDC4VPRequest(
        client_id=data.get("client_id", ""),
        response_type=data.get("response_type", "id_token"),
        scope=data.get("scope", "openid"),
        nonce=data.get("nonce"),
        state=data.get("state"),
        presentation_definition=pd,
        redirect_uri=data.get("redirect_uri", "openid://"),
        response_mode=data.get("response_mode", "direct_post"),
    )


def _deserialize_response(data: dict[str, Any]) -> OIDC4VPResponse:
    vp_token = data.get("vp_token")
    ps_data = data.get("presentation_submission")
    ps = None
    if ps_data:
        from services.siopv2 import PresentationSubmission
        ps = PresentationSubmission(
            definition_id=ps_data.get("definition_id", ""),
            descriptor_map=ps_data.get("descriptor_map", []),
        )
    return OIDC4VPResponse(
        id_token=data.get("id_token", ""),
        vp_token=vp_token,
        presentation_submission=ps,
        state=data.get("state"),
    )
