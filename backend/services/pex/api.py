from __future__ import annotations

from typing import Any

from fastapi import APIRouter, Body, HTTPException

from services.pex import (
    PresentationDefinition,
    evaluate_definition,
    filter_credentials_by_definition,
    create_presentation_from_definition,
)
from services.wallet import Wallet

router = APIRouter(prefix="/v1/pex", tags=["pex"])


_wallet_cache: dict[str, Wallet] = {}


def _get_wallet(did: str) -> Wallet:
    if did not in _wallet_cache:
        _wallet_cache[did] = Wallet(did)
    return _wallet_cache[did]


@router.post("/evaluate")
def evaluate_presentation_definition(body: dict[str, Any] = Body(...)):
    definition = body["definition"]
    holder_did = body["holder_did"]
    wallet = _get_wallet(holder_did)
    pd = PresentationDefinition.from_dict(definition)
    credentials = [cred.credential for cred in wallet.list()]
    result = evaluate_definition(pd, credentials)
    return {
        "definition_id": result.definition_id,
        "valid": result.valid,
        "matches": [
            {
                "descriptor_id": m.descriptor_id,
                "matched": m.matched,
                "errors": m.errors,
            }
            for m in result.matches
        ],
        "errors": result.errors,
    }


@router.post("/filter")
def filter_credentials(body: dict[str, Any] = Body(...)):
    definition = body["definition"]
    holder_did = body["holder_did"]
    wallet = _get_wallet(holder_did)
    pd = PresentationDefinition.from_dict(definition)
    credentials = [cred.credential for cred in wallet.list()]
    filtered = filter_credentials_by_definition(pd, credentials)
    return {
        "total_input": len(credentials),
        "total_matched": len(filtered),
        "credentials": filtered,
    }


@router.post("/present")
def create_presentation(body: dict[str, Any] = Body(...)):
    definition = body["definition"]
    holder_did = body["holder_did"]
    issuer_did = body.get("issuer_did", "")
    wallet = _get_wallet(holder_did)
    pd = PresentationDefinition.from_dict(definition)
    credentials = [cred.credential for cred in wallet.list()]
    vp = create_presentation_from_definition(pd, credentials, holder_did, issuer_did or None)
    if vp is None:
        raise HTTPException(status_code=400, detail="Could not satisfy PresentationDefinition")
    return vp.to_dict() if hasattr(vp, "to_dict") else vp
