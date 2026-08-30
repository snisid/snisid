from __future__ import annotations

import uuid
from datetime import datetime, timezone

from sqlalchemy import BigInteger, Boolean, DateTime, Float, ForeignKey, Integer, String, Text, UniqueConstraint
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.orm import Mapped, mapped_column, relationship

from shared.database import Base

JSON = dict | list | str | int | float | bool | None


class DIDRecord(Base):
    __tablename__ = "ssi_did"

    did: Mapped[str] = mapped_column(String(256), unique=True, nullable=False, index=True)
    method: Mapped[str] = mapped_column(String(32), nullable=False)
    document: Mapped[JSON] = mapped_column(JSONB, nullable=False)


class VerifiableCredentialRecord(Base):
    __tablename__ = "ssi_vc"

    credential_id: Mapped[str] = mapped_column(String(128), unique=True, nullable=False, index=True)
    issuer_id: Mapped[str] = mapped_column(String(256), nullable=False, index=True)
    subject_id: Mapped[str] = mapped_column(String(256), nullable=False, index=True)
    credential_type: Mapped[str] = mapped_column(String(128), nullable=False)
    document: Mapped[JSON] = mapped_column(JSONB, nullable=False)
    status_list_id: Mapped[str | None] = mapped_column(String(128), nullable=True)
    revoked: Mapped[bool] = mapped_column(Boolean, default=False)
    issued_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)


class StatusListRecord(Base):
    __tablename__ = "ssi_status_list"

    list_id: Mapped[str] = mapped_column(String(128), unique=True, nullable=False, index=True)
    purpose: Mapped[str] = mapped_column(String(64), nullable=False)
    bitstring: Mapped[str] = mapped_column(Text, nullable=False)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)


class WalletCredentialRecord(Base):
    __tablename__ = "ssi_wallet_credential"

    wallet_did: Mapped[str] = mapped_column(String(256), nullable=False, index=True)
    credential_id: Mapped[str] = mapped_column(String(128), nullable=False)
    document: Mapped[JSON] = mapped_column(JSONB, nullable=False)
    issuer_id: Mapped[str] = mapped_column(String(256), nullable=False)
    credential_type: Mapped[str] = mapped_column(String(128), nullable=False)
    issued_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    stored_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)

    __table_args__ = (
        UniqueConstraint("wallet_did", "credential_id", name="uq_wallet_credential"),
    )


class DIDCommMessageRecord(Base):
    __tablename__ = "ssi_didcomm_message"

    message_id: Mapped[str] = mapped_column(String(128), unique=True, nullable=False, index=True)
    sender_did: Mapped[str] = mapped_column(String(256), nullable=False)
    receiver_did: Mapped[str] = mapped_column(String(256), nullable=False, index=True)
    message_type: Mapped[str] = mapped_column(String(128), nullable=False)
    message_body: Mapped[JSON] = mapped_column(JSONB, nullable=False)
    thread_id: Mapped[str | None] = mapped_column(String(128), nullable=True)
    is_read: Mapped[bool] = mapped_column(Boolean, default=False)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)


class CredentialFlowRecord(Base):
    __tablename__ = "ssi_credential_flow"

    flow_id: Mapped[str] = mapped_column(String(128), unique=True, nullable=False, index=True)
    issuer_id: Mapped[str] = mapped_column(String(256), nullable=False)
    offer_data: Mapped[JSON] = mapped_column(JSONB, nullable=False)
    request_data: Mapped[JSON | None] = mapped_column(JSONB, nullable=True)
    credential_id: Mapped[str | None] = mapped_column(String(128), nullable=True)
    status: Mapped[str] = mapped_column(String(32), nullable=False)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)


class CHAPIRecord(Base):
    __tablename__ = "ssi_chapi"

    holder_did: Mapped[str] = mapped_column(String(256), nullable=False, index=True)
    credential_id: Mapped[str] = mapped_column(String(128), nullable=False)
    document: Mapped[JSON] = mapped_column(JSONB, nullable=False)
    query_frame: Mapped[JSON | None] = mapped_column(JSONB, nullable=True)
    stored_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)

    __table_args__ = (
        UniqueConstraint("holder_did", "credential_id", name="uq_chapi_credential"),
    )


class CredentialManifestRecord(Base):
    __tablename__ = "ssi_credential_manifest"

    manifest_id: Mapped[str] = mapped_column(String(128), unique=True, nullable=False, index=True)
    issuer_id: Mapped[str] = mapped_column(String(256), nullable=False)
    document: Mapped[JSON] = mapped_column(JSONB, nullable=False)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)


class RevocationEventRecord(Base):
    __tablename__ = "ssi_revocation_event"

    event_id: Mapped[str] = mapped_column(String(128), unique=True, nullable=False, index=True)
    credential_id: Mapped[str] = mapped_column(String(128), nullable=False, index=True)
    subject_id: Mapped[str] = mapped_column(String(256), nullable=False, index=True)
    event_type: Mapped[str] = mapped_column(String(64), nullable=False)
    reason: Mapped[str | None] = mapped_column(Text, nullable=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
