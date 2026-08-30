from __future__ import annotations

import uuid
from datetime import datetime, timezone
from typing import Any

from services.didcomm import DIDCommMessage, DIDCommMessenger
from services.vc.issuer import VCIssuer
from services.vc.verifier import VCVerifier


class CredentialOffer:
    def __init__(
        self,
        offer_id: str,
        issuer_did: str,
        holder_did: str,
        credential_type: str,
        claims_requested: dict[str, str],
    ):
        self.offer_id = offer_id
        self.issuer_did = issuer_did
        self.holder_did = holder_did
        self.credential_type = credential_type
        self.claims_requested = claims_requested
        self.created_at = datetime.now(timezone.utc).isoformat()
        self.status = "pending"

    def to_dict(self) -> dict[str, Any]:
        return {
            "offer_id": self.offer_id,
            "issuer_did": self.issuer_did,
            "holder_did": self.holder_did,
            "credential_type": self.credential_type,
            "claims_requested": self.claims_requested,
            "created_at": self.created_at,
            "status": self.status,
        }


class CredentialRequest:
    def __init__(
        self,
        request_id: str,
        offer_id: str,
        holder_did: str,
        claims: dict[str, Any],
    ):
        self.request_id = request_id
        self.offer_id = offer_id
        self.holder_did = holder_did
        self.claims = claims
        self.created_at = datetime.now(timezone.utc).isoformat()

    def to_dict(self) -> dict[str, Any]:
        return {
            "request_id": self.request_id,
            "offer_id": self.offer_id,
            "holder_did": self.holder_did,
            "claims": self.claims,
            "created_at": self.created_at,
        }


class CredentialFlow:
    def __init__(
        self,
        issuer: VCIssuer,
        verifier: VCVerifier | None = None,
        messenger: DIDCommMessenger | None = None,
        storage: Any | None = None,
    ):
        self._issuer = issuer
        self._verifier = verifier or VCVerifier(trusted_issuers=[issuer._issuer_id])
        self._messenger = messenger or DIDCommMessenger()
        self._offers: dict[str, CredentialOffer] = {}
        self._storage = storage

    def create_offer(
        self,
        issuer_did: str,
        holder_did: str,
        credential_type: str = "SNISIDIdentityCredential",
        claims_requested: dict[str, str] | None = None,
    ) -> CredentialOffer:
        offer = CredentialOffer(
            offer_id=str(uuid.uuid4()),
            issuer_did=issuer_did,
            holder_did=holder_did,
            credential_type=credential_type,
            claims_requested=claims_requested or {
                "national_id": "National ID number",
                "first_name": "First name",
                "last_name": "Last name",
                "date_of_birth": "Date of birth",
                "gender": "Gender",
                "nationality": "Nationality",
            },
        )
        self._offers[offer.offer_id] = offer
        return offer

    async def async_create_offer(
        self,
        issuer_did: str,
        holder_did: str,
        credential_type: str = "SNISIDIdentityCredential",
        claims_requested: dict[str, str] | None = None,
    ) -> CredentialOffer:
        offer = self.create_offer(issuer_did, holder_did, credential_type, claims_requested)
        if self._storage:
            await self._storage.save(
                flow_id=offer.offer_id,
                issuer_id=issuer_did,
                offer_data=offer.to_dict(),
                status="pending",
            )
        return offer

    def send_offer(self, offer: CredentialOffer) -> dict[str, Any]:
        msg = DIDCommMessage(
            id=str(uuid.uuid4()),
            type="https://didcomm.org/issue-credential/2.0/offer-credential",
            body=offer.to_dict(),
            from_did=offer.issuer_did,
            to_did=offer.holder_did,
        )
        return self._messenger.send(msg, offer.issuer_did, offer.holder_did)

    def build_request(
        self,
        offer_id: str,
        holder_did: str,
        claims: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        return {
            "request_id": str(uuid.uuid4()),
            "offer_id": offer_id,
            "holder_did": holder_did,
            "claims": claims or {
                "national_id": "SN-CITIZEN-001",
                "first_name": "Alice",
                "last_name": "Citizen",
                "date_of_birth": "1990-01-15",
                "gender": "female",
                "nationality": "HTI",
            },
        }

    def pack_and_send_request(
        self, request_body: dict[str, Any], holder_did: str, issuer_did: str
    ) -> dict[str, Any]:
        msg = DIDCommMessage(
            id=str(uuid.uuid4()),
            type="https://didcomm.org/issue-credential/2.0/request-credential",
            body=request_body,
            from_did=holder_did,
            to_did=issuer_did,
        )
        return self._messenger.send(msg, holder_did, issuer_did)

    def receive_request(self, packed: dict[str, Any]) -> CredentialRequest:
        msg = self._messenger.receive(packed)
        body = msg.body
        return CredentialRequest(
            request_id=body.get("request_id", str(uuid.uuid4())),
            offer_id=body.get("offer_id", ""),
            holder_did=body.get("holder_did", msg.from_did or ""),
            claims=body.get("claims", {}),
        )

    async def async_receive_request(self, packed: dict[str, Any],
                                     offer_sender_did: str | None = None) -> CredentialRequest:
        req = self.receive_request(packed)
        if self._storage:
            await self._storage.update_request(req.offer_id, req.to_dict())
        return req

    def issue_from_request(
        self, credential_request: CredentialRequest
    ) -> dict[str, Any] | None:
        offer = self._offers.get(credential_request.offer_id)
        if not offer:
            return None
        offer.status = "fulfilled"

        claims = credential_request.claims
        vc = self._issuer.issue_identity_credential(
            subject_id=credential_request.holder_did,
            national_id=claims.get("national_id", ""),
            first_name=claims.get("first_name", ""),
            last_name=claims.get("last_name", ""),
            date_of_birth=claims.get("date_of_birth", ""),
            gender=claims.get("gender", ""),
            nationality=claims.get("nationality", ""),
        )
        return vc.model_dump()

    async def async_issue_from_request(
        self, credential_request: CredentialRequest
    ) -> dict[str, Any] | None:
        offer = self._offers.get(credential_request.offer_id)
        if not offer:
            return None
        offer.status = "fulfilled"

        claims = credential_request.claims
        vc = self._issuer.issue_identity_credential(
            subject_id=credential_request.holder_did,
            national_id=claims.get("national_id", ""),
            first_name=claims.get("first_name", ""),
            last_name=claims.get("last_name", ""),
            date_of_birth=claims.get("date_of_birth", ""),
            gender=claims.get("gender", ""),
            nationality=claims.get("nationality", ""),
        )
        vc_data = vc.model_dump()

        if self._storage:
            await self._storage.update_issued(credential_request.offer_id, vc.id)

        return vc_data

    def send_credential(
        self, vc_data: dict[str, Any], issuer_did: str, holder_did: str
    ) -> dict[str, Any]:
        msg = DIDCommMessage(
            id=str(uuid.uuid4()),
            type="https://didcomm.org/issue-credential/2.0/issue-credential",
            body={"credentials": [vc_data]},
            from_did=issuer_did,
            to_did=holder_did,
        )
        return self._messenger.send(msg, issuer_did, holder_did)

    def receive_credential(self, packed: dict[str, Any]) -> dict[str, Any] | None:
        msg = self._messenger.receive(packed)
        credentials = msg.body.get("credentials", [])
        return credentials[0] if credentials else None
