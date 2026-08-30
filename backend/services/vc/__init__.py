from __future__ import annotations

import enum
from datetime import datetime, timezone
from typing import Any

from pydantic import BaseModel, Field


# ── VC Status ────────────────────────────────────────────────────────────


class VCStatus(str, enum.Enum):
    ACTIVE = "active"
    SUSPENDED = "suspended"
    REVOKED = "revoked"
    EXPIRED = "expired"


# ── W3C VC Data Model 2.0 ───────────────────────────────────────────────


class CredentialSubject(BaseModel):
    id: str = ""
    additional: dict[str, Any] = Field(default_factory=dict)


class VerifiableCredential(BaseModel):
    """W3C Verifiable Credential Data Model 2.0."""
    id: str = ""
    type: list[str] = Field(default_factory=lambda: ["VerifiableCredential"])
    issuer: str = ""
    issuanceDate: str = Field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    expirationDate: str | None = None
    credentialSubject: CredentialSubject = Field(default_factory=CredentialSubject)
    proof: dict[str, Any] | None = None
    credentialStatus: dict[str, Any] | None = None
    validFrom: str | None = None
    validUntil: str | None = None


class VerifiablePresentation(BaseModel):
    """W3C Verifiable Presentation."""
    id: str = ""
    type: list[str] = Field(default_factory=lambda: ["VerifiablePresentation"])
    holder: str = ""
    verifiableCredential: list[VerifiableCredential] = Field(default_factory=list)
    proof: dict[str, Any] | None = None


# ── SNISID-specific VC Schemas ────────────────────────────────────────────


class IdentityCredentialSubject(BaseModel):
    """Subject data for a SNISID Identity VC."""
    id: str
    national_id: str
    first_name: str
    last_name: str
    date_of_birth: str
    gender: str
    nationality: str
    status: str


class IdentityCredential(VerifiableCredential):
    type: list[str] = Field(
        default_factory=lambda: ["VerifiableCredential", "SNISIDIdentityCredential"]
    )
    credentialSubject: IdentityCredentialSubject
