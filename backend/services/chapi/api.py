from __future__ import annotations

from typing import Any

from fastapi import APIRouter, HTTPException

from services.chapi import (
    CHAPIMediator,
    CHAPIStoreRequest,
    CHAPIGetRequest,
)
from services.wallet import Wallet

router = APIRouter(prefix="/v1/chapi", tags=["chapi"])

_mediator: CHAPIMediator | None = None


def get_mediator() -> CHAPIMediator:
    global _mediator
    if _mediator is None:
        _mediator = CHAPIMediator()
    return _mediator


@router.post("/store")
def chapi_store(payload: dict[str, Any]):
    mediator = get_mediator()
    credential = payload.get("credential", payload)
    req = CHAPIStoreRequest(
        credential=credential,
        protocol=payload.get("protocol", "vc"),
    )
    result = mediator.handle_store(req)
    if result.error:
        raise HTTPException(400, result.error)
    return result.to_dict()


@router.post("/get")
def chapi_get(payload: dict[str, Any]):
    mediator = get_mediator()
    query = payload.get("query", [{"type": "VerifiableCredential"}])
    req = CHAPIGetRequest(
        query=query,
        protocol=payload.get("protocol", "vc"),
    )
    result = mediator.handle_get(req)
    if result.error:
        raise HTTPException(400, result.error)
    return result.to_dict()


@router.get("/register")
def chapi_register(origin: str = "https://example.com"):
    mediator = get_mediator()
    result = mediator.handler_registration(origin)
    return result.to_dict()


@router.get("/wallet")
def chapi_wallet_info():
    mediator = get_mediator()
    wallet = mediator.wallet
    return {
        "did": wallet.did,
        "credential_count": wallet.count(),
        "capabilities": ["store", "get"],
    }
