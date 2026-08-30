from __future__ import annotations

from fastapi import APIRouter, HTTPException

from services.didcomm_mediator import DIDCommMediator

router = APIRouter(prefix="/v1/didcomm/mediator", tags=["didcomm-mediator"])

_mediator: DIDCommMediator | None = None


def get_mediator() -> DIDCommMediator:
    global _mediator
    if _mediator is None:
        _mediator = DIDCommMediator()
    return _mediator


@router.post("/forward")
def forward_message(payload: dict):
    """Forward a DIDComm message through the mediator."""
    mediator = get_mediator()
    msg = mediator.forward_request(payload)
    if msg is None:
        raise HTTPException(status_code=400, detail="Invalid forward request")
    return {
        "message_id": msg.id,
        "recipient_did": msg.recipient_did,
        "sender_did": msg.sender_did,
        "created_at": msg.created_at.isoformat(),
    }


@router.post("/forward/packed")
def forward_packed(recipient_did: str, packed_message: dict, sender_did: str = ""):
    """Forward a packed DIDComm message for a recipient."""
    mediator = get_mediator()
    msg = mediator.forward(recipient_did, packed_message, sender_did)
    return {
        "message_id": msg.id,
        "recipient_did": msg.recipient_did,
        "created_at": msg.created_at.isoformat(),
    }


@router.get("/inbox/{recipient_did}")
def get_inbox(recipient_did: str):
    """Get all messages for a recipient."""
    mediator = get_mediator()
    msgs = mediator.get_inbox(recipient_did)
    return {
        "messages": [
            {
                "id": m.id,
                "sender_did": m.sender_did,
                "created_at": m.created_at.isoformat(),
                "is_delivered": m.is_delivered,
            }
            for m in msgs
        ],
        "total": len(msgs),
    }


@router.get("/inbox/{recipient_did}/pending")
def get_pending(recipient_did: str):
    """Get undelivered messages count and list."""
    mediator = get_mediator()
    msgs = mediator.fetch_messages(recipient_did)
    return {
        "messages": [
            {
                "id": m.id,
                "sender_did": m.sender_did,
                "created_at": m.created_at.isoformat(),
            }
            for m in msgs
        ],
        "total": len(msgs),
    }


@router.post("/deliver/{message_id}")
def deliver_message(message_id: str):
    """Mark a message as delivered and return its content."""
    mediator = get_mediator()
    msg = mediator.deliver(message_id)
    if msg is None:
        raise HTTPException(status_code=404, detail="Message not found")
    return {
        "id": msg.id,
        "sender_did": msg.sender_did,
        "recipient_did": msg.recipient_did,
        "packed_message": msg.packed_message,
        "delivered_at": msg.delivered_at.isoformat() if msg.delivered_at else None,
    }


@router.get("/pending-count/{recipient_did}")
def pending_count(recipient_did: str):
    """Get count of undelivered messages."""
    mediator = get_mediator()
    return {"recipient_did": recipient_did, "pending_count": mediator.get_pending_count(recipient_did)}


@router.delete("/messages/{message_id}")
def delete_message(message_id: str):
    """Delete a message from the mediator."""
    mediator = get_mediator()
    if not mediator.delete_message(message_id):
        raise HTTPException(status_code=404, detail="Message not found")
    return {"status": "deleted"}
