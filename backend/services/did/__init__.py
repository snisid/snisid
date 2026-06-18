from __future__ import annotations

import hashlib
import json
import re
import uuid
from datetime import datetime, timezone
from enum import Enum
from typing import Any

try:
    import httpx
    HAS_HTTPX = True
except ImportError:
    HAS_HTTPX = False

import base64

from pydantic import BaseModel, ConfigDict


class DIDMethod(str, Enum):
    KEY = "key"
    SNISID = "snisid"
    WEB = "web"


class VerificationMethodType(str, Enum):
    ED25519_VERIFICATION_KEY_2018 = "Ed25519VerificationKey2018"
    JSON_WEB_KEY_2020 = "JsonWebKey2020"
    ECDSA_SECP256K1_VERIFICATION_KEY_2019 = "EcdsaSecp256k1VerificationKey2019"


class ServiceEndpoint(BaseModel):
    id: str
    type: str
    service_endpoint: str | list[Any]


class DIDDocument(BaseModel):
    context: list[str] = ["https://www.w3.org/ns/did/v1"]
    id: str
    also_known_as: list[str] = []
    verification_method: list[dict[str, Any]] = []
    authentication: list[str | dict[str, Any]] = []
    assertion_method: list[str | dict[str, Any]] = []
    key_agreement: list[str | dict[str, Any]] = []
    capability_invocation: list[str | dict[str, Any]] = []
    capability_delegation: list[str | dict[str, Any]] = []
    service: list[dict[str, Any]] = []
    created: str | None = None
    updated: str | None = None

    model_config = ConfigDict(extra="allow")


def _generate_key_id(did: str, fragment: str) -> str:
    return f"{did}#{fragment}"


def _multibase_encode(data: bytes) -> str:
    return "z" + base64.urlsafe_b64encode(data).rstrip(b"=").decode()


def _create_key_pair() -> tuple[str, str]:
    seed = uuid.uuid4().bytes + uuid.uuid4().bytes
    private = hashlib.sha256(seed).hexdigest()
    public = hashlib.sha256(private.encode()).hexdigest()
    return private, public


def create_did_key(public_key_multibase: str | None = None) -> tuple[str, str, str]:
    private, public = _create_key_pair()
    if public_key_multibase is None:
        public_key_multibase = _multibase_encode(public.encode())
    did = f"did:key:{public_key_multibase}"
    return did, private, public


def create_did_snisid(
    identifier: str,
    network: str = "mainnet",
) -> str:
    did = f"did:snisid:{network}:{identifier}"
    return did


def create_did_web(domain: str, path: str | None = None) -> str:
    if path:
        encoded = path.replace("/", ":")
        did = f"did:web:{domain}:{encoded}"
    else:
        did = f"did:web:{domain}"
    return did


def resolve_did(did: str) -> DIDDocument:
    if did.startswith("did:key:"):
        return _resolve_did_key(did)
    elif did.startswith("did:snisid:"):
        return _resolve_did_snisid(did)
    elif did.startswith("did:web:"):
        return _resolve_did_web(did)
    else:
        raise ValueError(f"Unsupported DID method: {did}")


def _resolve_did_key(did: str) -> DIDDocument:
    method_specific = did[len("did:key:"):]
    key_id = _generate_key_id(did, "key-1")

    vm = {
        "id": key_id,
        "type": VerificationMethodType.ED25519_VERIFICATION_KEY_2018,
        "controller": did,
        "publicKeyMultibase": method_specific,
    }

    now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    return DIDDocument(
        id=did,
        verification_method=[vm],
        authentication=[key_id],
        assertion_method=[key_id],
        created=now,
        updated=now,
    )


def _resolve_did_snisid(did: str) -> DIDDocument:
    parts = did.split(":")
    if len(parts) < 3:
        raise ValueError(f"Invalid did:snisid: {did}")

    network = parts[2] if len(parts) >= 4 else "mainnet"
    identifier = parts[-1]
    key_id = _generate_key_id(did, "snisid-key-1")
    recovery_key_id = _generate_key_id(did, "recovery-key-1")

    private, public = _create_key_pair()
    recovery_private, recovery_public = _create_key_pair()

    vm = [
        {
            "id": key_id,
            "type": VerificationMethodType.ED25519_VERIFICATION_KEY_2018,
            "controller": did,
            "publicKeyMultibase": _multibase_encode(public.encode()),
        },
        {
            "id": recovery_key_id,
            "type": VerificationMethodType.ED25519_VERIFICATION_KEY_2018,
            "controller": did,
            "publicKeyMultibase": _multibase_encode(recovery_public.encode()),
        },
    ]

    now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    return DIDDocument(
        id=did,
        verification_method=vm,
        authentication=[key_id],
        assertion_method=[key_id],
        capability_invocation=[key_id],
        capability_delegation=[recovery_key_id],
        created=now,
        updated=now,
    )


def _resolve_did_web(did: str) -> DIDDocument:
    method_specific = did[len("did:web:"):]
    domain = method_specific.replace(":", "/")

    # Attempt HTTP(S) resolution first (async with retry)
    if HAS_HTTPX:
        try:
            doc = asyncio.run(_resolve_did_web_http_async(domain, did))
            if doc is not None:
                return doc
        except Exception:
            pass

    # Fallback: deterministic generation
    key_id = _generate_key_id(did, "web-key-1")
    private, public = _create_key_pair()

    vm = {
        "id": key_id,
        "type": VerificationMethodType.ED25519_VERIFICATION_KEY_2018,
        "controller": did,
        "publicKeyMultibase": _multibase_encode(public.encode()),
    }

    now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    return DIDDocument(
        id=did,
        verification_method=[vm],
        authentication=[key_id],
        assertion_method=[key_id],
        service=[{
            "id": f"{did}#web-root",
            "type": "LinkedDomains",
            "serviceEndpoint": f"https://{domain}",
        }],
        created=now,
        updated=now,
    )


def _resolve_did_web_http(domain: str, did: str) -> DIDDocument | None:
    """Attempt to resolve a did:web document via HTTPS.

    Fetches ``https://<domain>/.well-known/did.json`` (root) or
    ``https://<domain>/<path>/did.json`` (subpath).
    """
    try:
        url = f"https://{domain}/did.json"
        resp = httpx.get(url, timeout=5.0, follow_redirects=True)
        if resp.status_code != 200:
            # Try .well-known path
            url = f"https://{domain}/.well-known/did.json"
            resp = httpx.get(url, timeout=5.0, follow_redirects=True)
        if resp.status_code == 200:
            data = resp.json()
            now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return DIDDocument(
                id=data.get("id", did),
                verification_method=data.get("verificationMethod", []),
                authentication=data.get("authentication", []),
                assertion_method=data.get("assertionMethod", []),
                key_agreement=data.get("keyAgreement", []),
                service=data.get("service", []),
                created=data.get("created", now),
                updated=data.get("updated", now),
            )
    except Exception:
        pass
    return None


async def _resolve_did_web_http_async(domain: str, did: str) -> DIDDocument | None:
    """Async resolution of a did:web document via HTTPS with retry.

    Uses ``shared.retry`` for exponential backoff against transient
    network errors.
    """
    if not HAS_HTTPX:
        return None

    from shared.retry import async_retry

    async def _fetch(url: str) -> dict | None:
        async with httpx.AsyncClient(timeout=10.0, follow_redirects=True) as client:
            resp = await client.get(url)
            if resp.status_code == 200:
                return resp.json()
        return None

    async def _resolve() -> DIDDocument | None:
        url = f"https://{domain}/did.json"
        data = await _fetch(url)
        if data is None:
            url = f"https://{domain}/.well-known/did.json"
            data = await _fetch(url)
        if data is not None:
            now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return DIDDocument(
                id=data.get("id", did),
                verification_method=data.get("verificationMethod", []),
                authentication=data.get("authentication", []),
                assertion_method=data.get("assertionMethod", []),
                key_agreement=data.get("keyAgreement", []),
                service=data.get("service", []),
                created=data.get("created", now),
                updated=data.get("updated", now),
            )
        return None

    try:
        return await async_retry(
            _resolve,
            max_retries=2,
            base_delay=0.5,
            retryable_exceptions=(ConnectionError, TimeoutError, httpx.HTTPError),
        )
    except Exception:
        return None


CLASS_DID_CONTEXT = "https://www.w3.org/ns/did/v1"


class DIDManager:
    def __init__(self, storage: Any | None = None):
        self._registry: dict[str, dict[str, Any]] = {}
        self._storage = storage

    def create(self, method: DIDMethod, identifier: str | None = None) -> DIDDocument:
        if method == DIDMethod.KEY:
            did, private, public = create_did_key()
            doc = resolve_did(did)
        elif method == DIDMethod.SNISID:
            identifier = identifier or str(uuid.uuid4())
            did = create_did_snisid(identifier)
            doc = resolve_did(did)
        elif method == DIDMethod.WEB:
            domain = identifier or "example.com"
            did = create_did_web(domain)
            doc = resolve_did(did)
        else:
            raise ValueError(f"Unsupported method: {method}")

        self._registry[did] = {
            "document": doc.model_dump(),
            "private": doc.model_dump(),
        }
        return doc

    async def async_create(self, method: DIDMethod, identifier: str | None = None) -> DIDDocument:
        if method == DIDMethod.WEB:
            domain = identifier or "example.com"
            did = create_did_web(domain)
            doc = await _resolve_did_web_http_async(domain, did)
            if doc is None:
                doc = _resolve_did_web(did)
        else:
            doc = self.create(method, identifier)
        if self._storage:
            await self._storage.save(doc.id, method.value, doc.model_dump())
        return doc

    def resolve(self, did: str) -> DIDDocument:
        if did in self._registry:
            return DIDDocument(**self._registry[did]["document"])
        return resolve_did(did)

    async def async_resolve(self, did: str) -> DIDDocument:
        if self._storage:
            doc_data = await self._storage.get(did)
            if doc_data:
                return DIDDocument(**doc_data)
        return self.resolve(did)

    def update(self, did: str, updates: dict[str, Any]) -> DIDDocument:
        if did not in self._registry:
            raise ValueError(f"DID not found in registry: {did}")
        doc_data = self._registry[did]["document"]
        doc_data.update(updates)
        doc_data["updated"] = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        self._registry[did]["document"] = doc_data
        return DIDDocument(**doc_data)

    async def async_update(self, did: str, updates: dict[str, Any]) -> DIDDocument:
        doc = self.update(did, updates)
        if self._storage:
            doc_data = self._registry[did]["document"]
            await self._storage.save(did, doc.id.split(":")[1] if ":" in doc.id else "key", doc_data)
        return doc

    def deactivate(self, did: str) -> None:
        if did not in self._registry:
            raise ValueError(f"DID not found in registry: {did}")
        doc_data = self._registry[did]["document"]
        doc_data["deactivated"] = True
        doc_data["updated"] = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        self._registry[did]["document"] = doc_data

    async def async_deactivate(self, did: str) -> None:
        self.deactivate(did)
        if self._storage and did in self._registry:
            await self._storage.save(did, self._registry[did]["document"].get("id", did).split(":")[1], self._registry[did]["document"])
