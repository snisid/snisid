from __future__ import annotations

import hashlib
import json
import uuid
from datetime import datetime, timezone
from typing import Any

import base64


def _b64(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode()


def _sha256(data: str) -> str:
    return _b64(hashlib.sha256(data.encode()).digest())


def _random_salt(length: int = 16) -> str:
    return _b64(uuid.uuid4().bytes + uuid.uuid4().bytes)[:length]


class SDJWTIssuer:
    """Issues SD-JWTs with selectively disclosable claims."""

    def __init__(self, issuer_id: str, signing_key: str = "dev-sd-key"):
        self._issuer_id = issuer_id
        self._key = signing_key

    def issue(
        self,
        subject: str,
        disclosed_claims: dict[str, Any],
        sd_claims: dict[str, Any] | None = None,
        expiration_seconds: int = 3600,
    ) -> tuple[str, list[dict[str, Any]]]:
        """
        Issue an SD-JWT.
        Returns (sd_jwt_string, disclosures_list).
        """
        disclosures: list[dict[str, Any]] = []
        sd_hashes: list[str] = []

        sd_claims = sd_claims or {}
        for claim_name, claim_value in sd_claims.items():
            salt = _random_salt()
            disclosure = [salt, claim_name, claim_value]
            disclosure_json = json.dumps(disclosure, separators=(",", ":"))
            disclosure_hash = _sha256(disclosure_json)
            sd_hashes.append(disclosure_hash)
            disclosures.append({
                "salt": salt,
                "name": claim_name,
                "value": claim_value,
                "_digest": disclosure_hash,
            })

        now = int(datetime.now(timezone.utc).timestamp())
        payload = {
            "iss": self._issuer_id,
            "sub": subject,
            "iat": now,
            "exp": now + expiration_seconds,
            "_sd": sd_hashes,
            **disclosed_claims,
        }

        header = {"alg": "HS256", "typ": "sd+jwt"}
        header_b64 = _b64(json.dumps(header, separators=(",", ":")).encode())
        payload_b64 = _b64(json.dumps(payload, separators=(",", ":")).encode())
        signature = _b64(
            hashlib.sha256(
                f"{header_b64}.{payload_b64}{self._key}".encode()
            ).digest()
        )
        sd_jwt = f"{header_b64}.{payload_b64}.{signature}"

        disclosure_list = []
        for d in disclosures:
            dj = json.dumps([d["salt"], d["name"], d["value"]], separators=(",", ":"))
            disclosure_list.append(_b64(dj.encode()))

        return sd_jwt, disclosure_list


class SDJWTVerifier:
    """Verifies SD-JWTs and validates selective disclosures."""

    def __init__(self, trusted_issuers: set[str] | None = None):
        self._trusted_issuers = trusted_issuers or set()

    def verify(
        self,
        sd_jwt: str,
        disclosures: list[str],
        required_claims: list[str] | None = None,
    ) -> dict[str, Any]:
        """
        Verify an SD-JWT with given disclosures.
        Returns the disclosed claims.
        """
        parts = sd_jwt.split(".")
        if len(parts) != 3:
            raise ValueError("Invalid SD-JWT format")

        header_b64, payload_b64, signature_b64 = parts

        try:
            header = json.loads(base64.urlsafe_b64decode(header_b64 + "=="))
        except Exception:
            raise ValueError("Invalid header encoding")

        if header.get("typ") != "sd+jwt":
            raise ValueError("Not an SD-JWT")

        try:
            payload = json.loads(base64.urlsafe_b64decode(payload_b64 + "=="))
        except Exception:
            raise ValueError("Invalid payload encoding")

        if not self._verify_signature(header_b64, payload_b64, signature_b64):
            raise ValueError("Invalid signature")

        issuer = payload.get("iss", "")
        if self._trusted_issuers and issuer not in self._trusted_issuers:
            raise ValueError(f"Untrusted issuer: {issuer}")

        exp = payload.get("exp", 0)
        if exp and exp < int(datetime.now(timezone.utc).timestamp()):
            raise ValueError("SD-JWT has expired")

        sd_hashes: list[str] = payload.get("_sd", [])

        disclosed: dict[str, Any] = {}
        for disclosure_b64 in disclosures:
            try:
                decoded = base64.urlsafe_b64decode(disclosure_b64 + "==").decode()
                parts_decoded = json.loads(decoded)
                if not isinstance(parts_decoded, list) or len(parts_decoded) != 3:
                    continue
                salt, claim_name, claim_value = parts_decoded
                digest = _sha256(decoded)
                if digest not in sd_hashes:
                    raise ValueError(f"Disclosure not in _sd hash array: {claim_name}")
                disclosed[claim_name] = claim_value
            except (json.JSONDecodeError, ValueError) as e:
                if "not in _sd" in str(e):
                    raise
                raise ValueError(f"Invalid disclosure: {e}")

        if required_claims:
            missing = [c for c in required_claims if c not in disclosed]
            if missing:
                raise ValueError(f"Required claims not disclosed: {missing}")

        for key in payload:
            if key not in ("_sd", "iss", "sub", "iat", "exp"):
                disclosed.setdefault(key, payload[key])

        return disclosed

    def _verify_signature(
        self, header_b64: str, payload_b64: str, signature_b64: str
    ) -> bool:
        expected = _b64(
            hashlib.sha256(
                f"{header_b64}.{payload_b64}dev-sd-key".encode()
            ).digest()
        )
        return signature_b64 == expected


class SDJWTBuilder:
    """Builds SD-JWT presentation payloads."""

    @staticmethod
    def create_presentation(
        sd_jwt: str,
        all_disclosures: list[str],
        disclose: list[str],
    ) -> tuple[str, list[str]]:
        matching = []
        for d in all_disclosures:
            try:
                decoded = base64.urlsafe_b64decode(d + "==").decode()
                parts = json.loads(decoded)
                if len(parts) == 3 and parts[1] in disclose:
                    matching.append(d)
            except Exception:
                continue
        return sd_jwt, matching
