from __future__ import annotations

import json
import uuid
from datetime import datetime, timezone
from hashlib import sha256
from typing import Any

from services.did import resolve_did, create_did_key, DIDDocument
from services.vp import VPIssuer, VerifiablePresentation

# ─── Presentation Definition (DIF PE) ────────────────────────────────────────


class InputDescriptor:
    def __init__(
        self,
        id: str,
        name: str,
        purpose: str,
        schema_uris: list[str] | None = None,
        constraints: dict[str, Any] | None = None,
    ):
        self.id = id
        self.name = name
        self.purpose = purpose
        self.schema_uris = schema_uris or [
            "https://www.w3.org/ns/credentials/v2"
        ]
        self.constraints = constraints or {
            "fields": [
                {
                    "path": ["$.type"],
                    "filter": {"type": "array", "contains": {"const": "VerifiableCredential"}},
                }
            ]
        }

    def to_dict(self) -> dict[str, Any]:
        return {
            "id": self.id,
            "name": self.name,
            "purpose": self.purpose,
            "schema": [{"uri": u} for u in self.schema_uris],
            "constraints": self.constraints,
        }


class PresentationDefinition:
    def __init__(self, id: str, input_descriptors: list[InputDescriptor], name: str = ""):
        self.id = id
        self.name = name
        self.input_descriptors = input_descriptors

    def to_dict(self) -> dict[str, Any]:
        return {
            "id": self.id,
            "name": self.name,
            "input_descriptors": [d.to_dict() for d in self.input_descriptors],
        }


class PresentationSubmission:
    def __init__(
        self,
        definition_id: str,
        descriptor_map: list[dict[str, Any]],
    ):
        self.definition_id = definition_id
        self.descriptor_map = descriptor_map

    def to_dict(self) -> dict[str, Any]:
        return {
            "id": str(uuid.uuid4()),
            "definition_id": self.definition_id,
            "descriptor_map": self.descriptor_map,
        }


# ─── Self-Issued ID Token ────────────────────────────────────────────────────


class SelfIssuedIDToken:
    def __init__(
        self,
        sub: str,
        sub_jwk: dict[str, Any],
        iss: str,
        nonce: str,
        exp: int | None = None,
        iat: int | None = None,
        did: str | None = None,
    ):
        self.sub = sub
        self.sub_jwk = sub_jwk
        self.iss = iss
        self.nonce = nonce
        self.did = did
        now = int(datetime.now(timezone.utc).timestamp())
        self.iat = iat or now
        self.exp = exp or (now + 3600)

    def to_dict(self) -> dict[str, Any]:
        payload: dict[str, Any] = {
            "sub": self.sub,
            "sub_jwk": self.sub_jwk,
            "iss": self.iss,
            "nonce": self.nonce,
            "iat": self.iat,
            "exp": self.exp,
        }
        if self.did:
            payload["did"] = self.did
        return payload

    def sign(self, wallet_key: str = "dev-siop-key") -> str:
        payload = self.to_dict()
        header = {"alg": "HS256", "typ": "JWT"}
        body = _base64url(json.dumps(header).encode()) + "." + _base64url(json.dumps(payload).encode())
        signature = _base64url(sha256((body + "." + wallet_key).encode()).digest())
        return body + "." + signature

    @classmethod
    def verify(cls, token: str, wallet_key: str = "dev-siop-key") -> SelfIssuedIDToken | None:
        parts = token.split(".")
        if len(parts) != 3:
            return None
        body = parts[0] + "." + parts[1]
        expected = _base64url(sha256((body + "." + wallet_key).encode()).digest())
        if parts[2] != expected:
            return None
        payload = json.loads(_base64url_decode(parts[1]))
        return cls(
            sub=payload.get("sub", ""),
            sub_jwk=payload.get("sub_jwk", {}),
            iss=payload.get("iss", ""),
            nonce=payload.get("nonce", ""),
            exp=payload.get("exp"),
            iat=payload.get("iat"),
            did=payload.get("did"),
        )


# ─── OIDC4VP Request / Response ──────────────────────────────────────────────


class OIDC4VPRequest:
    def __init__(
        self,
        client_id: str,
        response_type: str = "id_token",
        scope: str = "openid",
        nonce: str | None = None,
        state: str | None = None,
        presentation_definition: PresentationDefinition | None = None,
        redirect_uri: str = "openid://",
        response_mode: str = "direct_post",
        response_uri: str | None = None,
    ):
        self.client_id = client_id
        self.response_type = response_type
        self.scope = scope
        self.nonce = nonce or str(uuid.uuid4())
        self.state = state or str(uuid.uuid4())
        self.presentation_definition = presentation_definition
        self.redirect_uri = redirect_uri
        self.response_mode = response_mode
        self.response_uri = response_uri

    def to_dict(self) -> dict[str, Any]:
        d: dict[str, Any] = {
            "client_id": self.client_id,
            "response_type": self.response_type,
            "scope": self.scope,
            "nonce": self.nonce,
            "state": self.state,
            "redirect_uri": self.redirect_uri,
            "response_mode": self.response_mode,
        }
        if self.response_uri:
            d["response_uri"] = self.response_uri
        if self.presentation_definition:
            d["presentation_definition"] = self.presentation_definition.to_dict()
        return d


class OIDC4VPResponse:
    def __init__(
        self,
        id_token: str,
        vp_token: list[dict[str, Any]] | None = None,
        presentation_submission: PresentationSubmission | None = None,
        state: str | None = None,
    ):
        self.id_token = id_token
        self.vp_token = vp_token
        self.presentation_submission = presentation_submission
        self.state = state

    def to_dict(self) -> dict[str, Any]:
        d: dict[str, Any] = {
            "id_token": self.id_token,
        }
        if self.vp_token:
            d["vp_token"] = self.vp_token
        if self.presentation_submission:
            d["presentation_submission"] = self.presentation_submission.to_dict()
        if self.state:
            d["state"] = self.state
        return d


# ─── Wallet ──────────────────────────────────────────────────────────────────


class SIOPWallet:
    def __init__(self, wallet_key: str = "dev-siop-key"):
        self._key = wallet_key
        self._did, self._priv, self._pub = create_did_key()
        self._vp_issuer = VPIssuer(proof_key=wallet_key)
        self._credentials: list[dict[str, Any]] = []

    @property
    def did(self) -> str:
        return self._did

    def store_credential(self, credential: dict[str, Any]):
        self._credentials.append(credential)

    def store_credentials(self, credentials: list[dict[str, Any]]):
        self._credentials.extend(credentials)

    def respond_to_request(
        self, request: OIDC4VPRequest
    ) -> OIDC4VPResponse:
        id_token = SelfIssuedIDToken(
            sub=self._did,
            sub_jwk={"kty": "oct", "k": self._key},
            iss=self._did,
            nonce=request.nonce,
            did=self._did,
        ).sign(self._key)

        vp_token = None
        presentation_submission = None

        if request.presentation_definition:
            vp_token = []
            descriptors = []

            for i, desc in enumerate(request.presentation_definition.input_descriptors):
                matching = self._find_matching_credentials(desc)
                if matching:
                    vp = self._vp_issuer.create_presentation(
                        holder_did=self._did,
                        verifiable_credentials=matching,
                    )
                    vp_token.append(vp.to_dict())
                    descriptors.append({
                        "id": desc.id,
                        "format": "vp",
                        "path": f"$.vp_token[{i}]",
                    })

            if descriptors:
                presentation_submission = PresentationSubmission(
                    definition_id=request.presentation_definition.id,
                    descriptor_map=descriptors,
                )

        return OIDC4VPResponse(
            id_token=id_token,
            vp_token=vp_token,
            presentation_submission=presentation_submission,
            state=request.state,
        )

    def post_response(self, request: OIDC4VPRequest) -> OIDC4VPResponse:
        """Send a direct_post response to the verifier's ``response_uri``."""
        response = self.respond_to_request(request)
        if request.response_uri:
            try:
                import httpx
                httpx.post(
                    request.response_uri,
                    json=response.to_dict(),
                    timeout=10.0,
                )
            except Exception:
                pass
        return response

    def _find_matching_credentials(self, descriptor: InputDescriptor) -> list[dict[str, Any]]:
        matching = []
        for cred in self._credentials:
            types = cred.get("type", [])
            if "VerifiableCredential" in types:
                matching.append(cred)
        return matching


# ─── Verifier ────────────────────────────────────────────────────────────────


class OIDC4VPVerifier:
    def __init__(self, wallet_key: str = "dev-siop-key"):
        self._wallet_key = wallet_key

    def verify_response(
        self, response: OIDC4VPResponse, expected_nonce: str, expected_state: str | None = None
    ) -> dict[str, Any]:
        result: dict[str, Any] = {
            "valid": True,
            "errors": [],
        }

        id_token = SelfIssuedIDToken.verify(response.id_token, self._wallet_key)
        if id_token is None:
            result["valid"] = False
            result["errors"].append("Invalid id_token signature")
        else:
            if id_token.nonce != expected_nonce:
                result["valid"] = False
                result["errors"].append("nonce mismatch")
            result["subject"] = id_token.sub
            result["did"] = id_token.did

        if response.state and expected_state and response.state != expected_state:
            result["valid"] = False
            result["errors"].append("state mismatch")

        if response.vp_token:
            result["vp_count"] = len(response.vp_token)
            for i, vp_data in enumerate(response.vp_token):
                try:
                    vp = VerifiablePresentation.from_dict(vp_data)
                    vp_valid = self._vp_issuer.verify_presentation(vp)
                    if not vp_valid:
                        result["valid"] = False
                        result["errors"].append(f"vp_token[{i}] invalid signature")
                except Exception as e:
                    result["valid"] = False
                    result["errors"].append(f"vp_token[{i}] parse error: {e}")

        return result

    @property
    def _vp_issuer(self) -> VPIssuer:
        return VPIssuer(proof_key=self._wallet_key)


# ─── Helper ──────────────────────────────────────────────────────────────────


def _base64url(data: bytes) -> str:
    import base64
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode()


def _base64url_decode(s: str) -> bytes:
    import base64
    padding = 4 - len(s) % 4
    if padding != 4:
        s += "=" * padding
    return base64.urlsafe_b64decode(s)
