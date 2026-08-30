from __future__ import annotations

from typing import Any

from fastapi import APIRouter, HTTPException, Query

from services.wallet import Wallet
from services.vp import VerifiablePresentation

router = APIRouter(prefix="/v1/wallet", tags=["wallet"])

_wallet: Wallet | None = None


def get_wallet() -> Wallet:
    global _wallet
    if _wallet is None:
        _wallet = Wallet()
    return _wallet


@router.get("/did")
def wallet_did():
    return {"did": get_wallet().did}


@router.get("/did-document")
def wallet_did_document():
    return get_wallet().did_document


@router.post("/credentials")
async def store_credential(credential: dict[str, Any], label: str = ""):
    record = await get_wallet().async_store(credential, label=label)
    return record.to_dict()


@router.get("/credentials")
async def list_credentials(credential_type: str | None = Query(None)):
    records = await get_wallet().async_list(credential_type=credential_type)
    return {"credentials": [r.to_dict() for r in records], "total": len(records)}


@router.get("/credentials/count")
def wallet_count():
    return {"count": get_wallet().count()}


@router.get("/credentials/search")
async def search_credentials(q: str = Query(...)):
    records = await get_wallet().async_search(q)
    return {"credentials": [r.to_dict() for r in records], "total": len(records)}


@router.get("/credentials/by-issuer/{issuer_did}")
def get_by_issuer(issuer_did: str):
    records = get_wallet().get_credentials_by_issuer(issuer_did)
    return {"credentials": [r.to_dict() for r in records], "total": len(records)}


@router.get("/credentials/{credential_id}")
async def get_credential(credential_id: str):
    record = await get_wallet().async_get(credential_id)
    if not record:
        raise HTTPException(404, "Credential not found")
    return record.to_dict()


@router.delete("/credentials/{credential_id}")
async def delete_credential(credential_id: str):
    if not await get_wallet().async_delete(credential_id):
        raise HTTPException(404, "Credential not found")
    return {"status": "deleted"}


@router.post("/presentations")
def create_presentation(
    credential_ids: list[str] | None = None,
    credential_type: str | None = None,
):
    vp = get_wallet().create_presentation(
        credential_ids=credential_ids,
        credential_type=credential_type,
    )
    return vp.to_dict()


@router.post("/presentations/verify")
def verify_presentation(presentation: dict[str, Any]):
    try:
        vp = VerifiablePresentation.from_dict(presentation)
    except Exception as e:
        raise HTTPException(400, f"Invalid presentation: {e}")
    valid = get_wallet().verify_presentation(vp)
    return {"valid": valid}


@router.post("/didcomm/send/{credential_id}")
def didcomm_send(credential_id: str, to_did: str):
    try:
        packed = get_wallet().send_via_didcomm(credential_id, to_did)
        return {"status": "sent", "packed": packed}
    except ValueError as e:
        raise HTTPException(404, str(e))


@router.post("/didcomm/receive")
def didcomm_receive(packed: dict[str, Any], label: str = ""):
    try:
        stored = get_wallet().receive_via_didcomm(packed, label=label)
    except Exception:
        raise HTTPException(400, "Invalid DIDComm message")
    if not stored:
        raise HTTPException(400, "No credential in message")
    return stored.to_dict()


@router.post("/didcomm/message")
def didcomm_message(message_type: str, body: dict[str, Any], to_did: str):
    packed = get_wallet().send_message(message_type, body, to_did)
    return {"status": "sent", "packed": packed}


@router.post("/export")
def export_wallet():
    return {"credentials": get_wallet().export_all()}


@router.post("/import")
def import_wallet(data: list[dict[str, Any]]):
    get_wallet().import_credentials(data)
    return {"status": "imported", "count": get_wallet().count()}
