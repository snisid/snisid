from fastapi import APIRouter

from services.did import DIDDocument, DIDManager, DIDMethod, resolve_did

router = APIRouter(prefix="/did", tags=["did"])

_manager = DIDManager()


@router.post("/create")
async def create_did(method: str = "key", identifier: str | None = None):
    doc = await _manager.async_create(DIDMethod(method), identifier)
    return {"did": doc.id, "document": doc.model_dump()}


@router.get("/resolve/{did:path}")
async def resolve_did_endpoint(did: str):
    doc = await _manager.async_resolve(did)
    return doc.model_dump()


@router.post("/update/{did:path}")
async def update_did(did: str, updates: dict):
    doc = _manager.update(did, updates)
    return doc.model_dump()


@router.post("/deactivate/{did:path}")
async def deactivate_did(did: str):
    await _manager.async_deactivate(did)
    return {"status": "deactivated", "did": did}
