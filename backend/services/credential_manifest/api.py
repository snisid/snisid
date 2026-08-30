from __future__ import annotations

from typing import Any

from fastapi import APIRouter, HTTPException, Query

from services.credential_manifest import (
    CredentialApplication,
    CredentialManifest,
    ManifestManager,
    OutputDescriptor,
)
from services.vc.issuer import VCIssuer

router = APIRouter(prefix="/v1/credential-manifest", tags=["credential-manifest"])

_manager: ManifestManager | None = None


def get_manager() -> ManifestManager:
    global _manager
    if _manager is None:
        _manager = ManifestManager(
            issuer=VCIssuer(issuer_id="did:snisid:mainnet:manifest-issuer")
        )
    return _manager


@router.get("/manifests")
def list_manifests(issuer_did: str | None = Query(None)):
    manifests = get_manager().list_manifests(issuer_did=issuer_did)
    return {"manifests": [m.to_dict() for m in manifests], "total": len(manifests)}


@router.post("/manifests")
async def create_manifest(
    issuer_did: str,
    name: str = "",
    description: str = "",
):
    manifest = await get_manager().async_create_manifest(
        issuer_did=issuer_did,
        name=name,
        description=description,
    )
    return manifest.to_dict()


@router.get("/manifests/{manifest_id}")
async def get_manifest(manifest_id: str):
    manifest = await get_manager().async_get_manifest(manifest_id)
    if not manifest:
        raise HTTPException(404, "Manifest not found")
    return manifest.to_dict()


@router.post("/apply")
def apply(payload: dict[str, Any]):
    manager = get_manager()
    manifest_id = payload.get("manifest_id", "")
    applicant = payload.get("applicant", "")
    if not manifest_id or not applicant:
        raise HTTPException(400, "manifest_id and applicant required")

    application = CredentialApplication(
        id=payload.get("id", ""),
        manifest_id=manifest_id,
        applicant=applicant,
        claims=payload.get("claims", {}),
        presentation_submission=payload.get("presentation_submission"),
    )
    response = manager.submit_application(application)
    if response.error:
        raise HTTPException(400, response.error)
    return response.to_dict()


@router.get("/applications")
def list_applications(manifest_id: str | None = Query(None)):
    apps = get_manager().get_applications(manifest_id=manifest_id)
    return {"applications": [a.to_dict() for a in apps], "total": len(apps)}


@router.get("/responses/{response_id}")
def get_response(response_id: str):
    response = get_manager().get_response(response_id)
    if not response:
        raise HTTPException(404, "Response not found")
    return response.to_dict()
