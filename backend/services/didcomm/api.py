import uuid
from typing import Any

from fastapi import APIRouter
from pydantic import BaseModel

from services.didcomm import DIDCommMessage, DIDCommMessenger

router = APIRouter(prefix="/didcomm", tags=["didcomm"])

_messenger = DIDCommMessenger()


class PackRequest(BaseModel):
    type: str
    body: dict[str, Any]
    from_did: str | None = None
    to_did: str | None = None


class UnpackRequest(BaseModel):
    packed: dict[str, Any]


@router.post("/pack")
async def pack_message(req: PackRequest):
    msg = DIDCommMessage(
        id=str(uuid.uuid4()),
        type=req.type,
        body=req.body,
        from_did=req.from_did,
        to_did=req.to_did,
    )
    packed = _messenger.send(msg, req.from_did or "", req.to_did or "")
    return packed


@router.post("/unpack")
async def unpack_message(req: UnpackRequest):
    msg = _messenger.receive(req.packed)
    return msg.to_dict()


@router.post("/trust-ping")
async def trust_ping(from_did: str, to_did: str):
    ping = _messenger.create_trust_ping(from_did, to_did)
    packed = _messenger.send(ping, from_did, to_did)
    return packed


@router.post("/trust-ping-response")
async def trust_ping_response(packed: dict):
    ping = _messenger.receive(packed)
    response = _messenger.create_trust_ping_response(ping)
    packed_response = _messenger.send(
        response, response.from_did or "", response.to_did or ""
    )
    return packed_response
