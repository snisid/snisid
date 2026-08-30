from __future__ import annotations

import uuid
from datetime import datetime, timezone
from typing import Any

from services.vc.issuer import VCIssuer


class OutputDescriptor:
    def __init__(
        self,
        id: str,
        schema: str,
        name: str = "",
        description: str = "",
        display: dict[str, Any] | None = None,
    ):
        self.id = id
        self.schema = schema
        self.name = name
        self.description = description
        self.display = display or {}

    def to_dict(self) -> dict[str, Any]:
        return {
            "id": self.id,
            "schema": self.schema,
            "name": self.name,
            "description": self.description,
            "display": self.display,
        }


class CredentialManifest:
    def __init__(
        self,
        id: str,
        issuer: str,
        output_descriptors: list[OutputDescriptor],
        name: str = "",
        description: str = "",
        presentation_definition: dict[str, Any] | None = None,
    ):
        self.id = id
        self.issuer = issuer
        self.name = name
        self.description = description
        self.output_descriptors = output_descriptors
        self.presentation_definition = presentation_definition

    def to_dict(self) -> dict[str, Any]:
        d: dict[str, Any] = {
            "id": self.id,
            "issuer": self.issuer,
            "name": self.name,
            "description": self.description,
            "output_descriptors": [od.to_dict() for od in self.output_descriptors],
        }
        if self.presentation_definition:
            d["presentation_definition"] = self.presentation_definition
        return d


class CredentialApplication:
    def __init__(
        self,
        id: str,
        manifest_id: str,
        applicant: str,
        presentation_submission: dict[str, Any] | None = None,
        claims: dict[str, Any] | None = None,
    ):
        self.id = id
        self.manifest_id = manifest_id
        self.applicant = applicant
        self.presentation_submission = presentation_submission
        self.claims = claims or {}
        self.created_at = datetime.now(timezone.utc).isoformat()

    def to_dict(self) -> dict[str, Any]:
        d: dict[str, Any] = {
            "id": self.id,
            "manifest_id": self.manifest_id,
            "applicant": self.applicant,
            "created_at": self.created_at,
        }
        if self.presentation_submission:
            d["presentation_submission"] = self.presentation_submission
        if self.claims:
            d["claims"] = self.claims
        return d


class CredentialResponse:
    def __init__(
        self,
        id: str,
        manifest_id: str,
        applicant: str,
        credentials: list[dict[str, Any]] | None = None,
        error: str | None = None,
    ):
        self.id = id
        self.manifest_id = manifest_id
        self.applicant = applicant
        self.credentials = credentials
        self.error = error
        self.created_at = datetime.now(timezone.utc).isoformat()

    def to_dict(self) -> dict[str, Any]:
        d: dict[str, Any] = {
            "id": self.id,
            "manifest_id": self.manifest_id,
            "applicant": self.applicant,
            "created_at": self.created_at,
        }
        if self.credentials:
            d["credentials"] = self.credentials
        if self.error:
            d["error"] = self.error
        return d


class ManifestManager:
    def __init__(self, issuer: VCIssuer, storage: Any | None = None):
        self._issuer = issuer
        self._manifests: dict[str, CredentialManifest] = {}
        self._applications: dict[str, CredentialApplication] = {}
        self._responses: dict[str, CredentialResponse] = {}
        self._storage = storage

    def register(self, manifest: CredentialManifest):
        self._manifests[manifest.id] = manifest

    def create_manifest(
        self,
        issuer_did: str,
        output_descriptors: list[OutputDescriptor] | None = None,
        name: str = "",
        description: str = "",
        presentation_definition: dict[str, Any] | None = None,
    ) -> CredentialManifest:
        manifest = CredentialManifest(
            id=str(uuid.uuid4()),
            issuer=issuer_did,
            name=name,
            description=description,
            output_descriptors=output_descriptors or [
                OutputDescriptor(
                    id="identity-credential",
                    schema="https://www.w3.org/ns/credentials/v2",
                    name="SNISID Identity Credential",
                    description="National identity credential",
                ),
            ],
            presentation_definition=presentation_definition,
        )
        self.register(manifest)
        return manifest

    async def async_create_manifest(
        self,
        issuer_did: str,
        output_descriptors: list[OutputDescriptor] | None = None,
        name: str = "",
        description: str = "",
        presentation_definition: dict[str, Any] | None = None,
    ) -> CredentialManifest:
        manifest = self.create_manifest(issuer_did, output_descriptors, name, description, presentation_definition)
        if self._storage:
            await self._storage.save(manifest.id, issuer_did, manifest.to_dict())
        return manifest

    def get_manifest(self, manifest_id: str) -> CredentialManifest | None:
        return self._manifests.get(manifest_id)

    async def async_get_manifest(self, manifest_id: str) -> CredentialManifest | None:
        if self._storage:
            doc = await self._storage.get(manifest_id)
            if doc:
                od_list = [OutputDescriptor(**od) for od in doc.get("output_descriptors", [])]
                return CredentialManifest(
                    id=doc["id"],
                    issuer=doc["issuer"],
                    output_descriptors=od_list,
                    name=doc.get("name", ""),
                    description=doc.get("description", ""),
                    presentation_definition=doc.get("presentation_definition"),
                )
        return self.get_manifest(manifest_id)

    def list_manifests(self, issuer_did: str | None = None) -> list[CredentialManifest]:
        if issuer_did:
            return [m for m in self._manifests.values() if m.issuer == issuer_did]
        return list(self._manifests.values())

    async def async_list_manifests(self, issuer_did: str | None = None) -> list[CredentialManifest]:
        if self._storage:
            docs = await self._storage.list_by_issuer(issuer_did or "")
            manifests = []
            for doc in docs:
                od_list = [OutputDescriptor(**od) for od in doc.get("output_descriptors", [])]
                manifests.append(CredentialManifest(
                    id=doc["id"],
                    issuer=doc["issuer"],
                    output_descriptors=od_list,
                    name=doc.get("name", ""),
                    description=doc.get("description", ""),
                    presentation_definition=doc.get("presentation_definition"),
                ))
            return manifests
        return self.list_manifests(issuer_did)

    def submit_application(self, application: CredentialApplication) -> CredentialResponse:
        self._applications[application.id] = application
        manifest = self._manifests.get(application.manifest_id)
        if not manifest:
            return CredentialResponse(
                id=str(uuid.uuid4()),
                manifest_id=application.manifest_id,
                applicant=application.applicant,
                error="Manifest not found",
            )

        claims = application.claims
        vc = self._issuer.issue_identity_credential(
            subject_id=application.applicant,
            national_id=claims.get("national_id", application.applicant),
            first_name=claims.get("first_name", ""),
            last_name=claims.get("last_name", ""),
            date_of_birth=claims.get("date_of_birth", ""),
            gender=claims.get("gender", ""),
            nationality=claims.get("nationality", ""),
        )
        response = CredentialResponse(
            id=str(uuid.uuid4()),
            manifest_id=application.manifest_id,
            applicant=application.applicant,
            credentials=[vc.model_dump()],
        )
        self._responses[response.id] = response
        return response

    def get_response(self, response_id: str) -> CredentialResponse | None:
        return self._responses.get(response_id)

    def get_applications(
        self, manifest_id: str | None = None
    ) -> list[CredentialApplication]:
        if manifest_id:
            return [a for a in self._applications.values() if a.manifest_id == manifest_id]
        return list(self._applications.values())
