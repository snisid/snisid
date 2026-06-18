from __future__ import annotations

from fastapi import APIRouter, HTTPException, Query

from services.revocation import RevocationNotifier, RevocationEventType, WalletRevocationHook

router = APIRouter(prefix="/v1/revocation", tags=["revocation"])

_notifier: RevocationNotifier | None = None
_hooks: dict[str, WalletRevocationHook] = {}


def get_notifier() -> RevocationNotifier:
    global _notifier
    if _notifier is None:
        _notifier = RevocationNotifier()
    return _notifier


def get_hook(wallet_did: str) -> WalletRevocationHook:
    if wallet_did not in _hooks:
        _hooks[wallet_did] = WalletRevocationHook(get_notifier(), wallet_did)
    return _hooks[wallet_did]


@router.post("/notify/revoke")
async def notify_revoke(credential_id: str, subject_id: str, reason: str = ""):
    event = await get_notifier().async_notify(
        credential_id=credential_id,
        event_type=RevocationEventType.CREDENTIAL_REVOKED,
        subject_id=subject_id,
        reason=reason,
    )
    return event.to_dict()


@router.post("/notify/suspend")
async def notify_suspend(credential_id: str, subject_id: str, reason: str = ""):
    event = await get_notifier().async_notify(
        credential_id=credential_id,
        event_type=RevocationEventType.CREDENTIAL_SUSPENDED,
        subject_id=subject_id,
        reason=reason,
    )
    return event.to_dict()


@router.post("/notify/reinstate")
async def notify_reinstate(credential_id: str, subject_id: str, reason: str = ""):
    event = await get_notifier().async_notify(
        credential_id=credential_id,
        event_type=RevocationEventType.CREDENTIAL_REINSTATED,
        subject_id=subject_id,
        reason=reason,
    )
    return event.to_dict()


@router.get("/history")
async def get_history(
    credential_id: str | None = Query(None),
    subject_id: str | None = Query(None),
    limit: int = 50,
):
    events = await get_notifier().async_get_history(
        credential_id=credential_id,
        subject_id=subject_id,
        limit=limit,
    )
    return {"events": [e.to_dict() for e in events], "total": len(events)}


@router.post("/wallet/{wallet_did}/track/{credential_id}")
def track_credential(wallet_did: str, credential_id: str):
    hook = get_hook(wallet_did)
    hook.track_credential(credential_id)
    return {"status": "tracking"}


@router.post("/wallet/{wallet_did}/untrack/{credential_id}")
def untrack_credential(wallet_did: str, credential_id: str):
    hook = get_hook(wallet_did)
    hook.untrack_credential(credential_id)
    return {"status": "untracked"}


@router.get("/wallet/{wallet_did}/notifications")
def wallet_notifications(wallet_did: str):
    hook = get_hook(wallet_did)
    return {
        "notifications": [e.to_dict() for e in hook.get_notifications()],
        "total": len(hook.get_notifications()),
    }


@router.post("/wallet/{wallet_did}/notifications/clear")
def clear_notifications(wallet_did: str):
    hook = get_hook(wallet_did)
    hook.clear_notifications()
    return {"status": "cleared"}


@router.get("/wallet/{wallet_did}/status/{credential_id}")
def check_status(wallet_did: str, credential_id: str):
    hook = get_hook(wallet_did)
    status = hook.check_status(credential_id)
    if status is None:
        return {"status": "unknown"}
    return {"status": status.value}
