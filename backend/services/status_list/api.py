from fastapi import APIRouter, HTTPException

from services.status_list import StatusListManager

router = APIRouter(prefix="/status-list", tags=["status-list"])

_manager = StatusListManager(issuer_id="https://snisid.ht/issuer")


@router.post("/entries")
async def create_entry(purpose: str = "revocation"):
    entry = await _manager.async_create_entry(purpose)
    return entry.to_dict()


@router.post("/revoke")
async def revoke_entry(entry_id: str):
    if not await _manager.async_revoke(entry_id):
        raise HTTPException(status_code=404, detail="Entry not found")
    return {"status": "revoked", "entry_id": entry_id}


@router.post("/unrevoke")
async def unrevoke_entry(entry_id: str):
    if not _manager.unrevoke(entry_id):
        raise HTTPException(status_code=404, detail="Entry not found")
    return {"status": "unrevoked", "entry_id": entry_id}


@router.get("/check")
async def check_entry(entry_id: str):
    revoked = await _manager.async_is_revoked(entry_id)
    if revoked is None:
        raise HTTPException(status_code=404, detail="Entry not found")
    return {"entry_id": entry_id, "revoked": revoked}


@router.get("/credential/{list_id:path}")
async def get_status_list_credential(list_id: str):
    credential = _manager.get_status_list_credential(list_id)
    if credential is None:
        raise HTTPException(status_code=404, detail="Status list not found")
    return credential
