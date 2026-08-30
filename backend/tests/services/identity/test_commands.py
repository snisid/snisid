from __future__ import annotations

import uuid
from datetime import date, datetime, timezone
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
import pytest_asyncio

from services.identity.commands import (
    IdentityCommandHandler,
    IdentityNotFoundError,
    CreateIdentityCommand,
    UpdateIdentityCommand,
    VerifyIdentityCommand,
    SuspendIdentityCommand,
    RevokeIdentityCommand,
    EnrollBiometricCommand,
    IssueDocumentCommand,
)
from services.identity.aggregate import IdentityAggregate
from services.identity.events import IdentityEvent


class TestCreateIdentityCommand:
    """Test CreateIdentityCommand handler."""

    @pytest.mark.asyncio
    async def test_handle_create(self, mock_db_session, mock_event_store, mock_kafka, identity_create_data):
        handler = IdentityCommandHandler(mock_db_session, mock_kafka)
        cmd = CreateIdentityCommand(**identity_create_data)

        result = await handler.handle_create(cmd)

        assert "identity_id" in result
        assert result["status"] == "pending"
        assert "version" in result
        assert result["version"] == 1

    @pytest.mark.asyncio
    async def test_handle_create_rejects_future_dob(self, mock_db_session, identity_create_data):
        handler = IdentityCommandHandler(mock_db_session)
        data = {**identity_create_data, "date_of_birth": "2099-12-31"}
        cmd = CreateIdentityCommand(**data)

        with pytest.raises(ValueError, match="Date of birth cannot be in the future"):
            await handler.handle_create(cmd)

    @pytest.mark.asyncio
    async def test_handle_create_rejects_invalid_nationality(self, mock_db_session, identity_create_data):
        handler = IdentityCommandHandler(mock_db_session)
        data = {**identity_create_data, "nationality": "123"}
        cmd = CreateIdentityCommand(**data)

        with pytest.raises(ValueError):
            await handler.handle_create(cmd)

    @pytest.mark.asyncio
    async def test_handle_create_publishes_to_kafka(self, mock_db_session, mock_event_store, mock_kafka, identity_create_data):
        handler = IdentityCommandHandler(mock_db_session, mock_kafka)
        cmd = CreateIdentityCommand(**identity_create_data)

        await handler.handle_create(cmd)
        assert mock_kafka.publish.await_count >= 0


class TestUpdateIdentityCommand:
    """Test UpdateIdentityCommand handler."""

    @pytest.mark.asyncio
    async def test_handle_update(self, mock_db_session, mock_event_store, mock_kafka, sample_identity_id):
        mock_event_store.get_events = AsyncMock(return_value=[])
        mock_event_store.get_snapshot = AsyncMock(return_value=None)

        handler = IdentityCommandHandler(mock_db_session, mock_kafka)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN123456789",
                first_name="John",
                last_name="Doe",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            mock_load.return_value = aggregate

            cmd = UpdateIdentityCommand(
                identity_id=aggregate.id,
                changes={"first_name": "Jane"},
                actor_id="user-2",
            )
            result = await handler.handle_update(cmd)
            assert result["identity_id"] == aggregate.id

    @pytest.mark.asyncio
    async def test_handle_update_not_found(self, mock_db_session, mock_event_store):
        mock_event_store.get_events = AsyncMock(return_value=[])
        mock_event_store.get_snapshot = AsyncMock(return_value=None)

        handler = IdentityCommandHandler(mock_db_session)

        with pytest.raises(IdentityNotFoundError):
            cmd = UpdateIdentityCommand(
                identity_id="nonexistent",
                changes={"first_name": "Jane"},
                actor_id="user-2",
            )
            await handler.handle_update(cmd)


class TestVerifyIdentityCommand:
    """Test VerifyIdentityCommand handler."""

    @pytest.mark.asyncio
    async def test_handle_verify(self, mock_db_session, mock_event_store, mock_kafka):
        handler = IdentityCommandHandler(mock_db_session, mock_kafka)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN123",
                first_name="J",
                last_name="D",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            mock_load.return_value = aggregate

            cmd = VerifyIdentityCommand(
                identity_id=aggregate.id,
                verification_method="biometric",
                verifier_id="verifier-1",
                actor_id="verifier-1",
            )
            result = await handler.handle_verify(cmd)
            assert result["status"] == "active"

    @pytest.mark.asyncio
    async def test_handle_verify_revoked_fails(self, mock_db_session, mock_event_store):
        handler = IdentityCommandHandler(mock_db_session)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN123",
                first_name="J",
                last_name="D",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            aggregate._status = "revoked"
            mock_load.return_value = aggregate

            cmd = VerifyIdentityCommand(
                identity_id=aggregate.id,
                verification_method="biometric",
                verifier_id="v-1",
                actor_id="v-1",
            )
            with pytest.raises(Exception):
                await handler.handle_verify(cmd)


class TestSuspendIdentityCommand:
    """Test SuspendIdentityCommand handler."""

    @pytest.mark.asyncio
    async def test_handle_suspend(self, mock_db_session, mock_event_store, mock_kafka):
        handler = IdentityCommandHandler(mock_db_session, mock_kafka)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN999",
                first_name="A", last_name="B",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            mock_load.return_value = aggregate

            cmd = SuspendIdentityCommand(
                identity_id=aggregate.id,
                reason="Fraud investigation ongoing",
                actor_id="user-2",
            )
            result = await handler.handle_suspend(cmd)
            assert result["status"] == "suspended"

    @pytest.mark.asyncio
    async def test_handle_suspend_valid_reason(self, mock_db_session, mock_event_store):
        handler = IdentityCommandHandler(mock_db_session)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN111",
                first_name="X", last_name="Y",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            mock_load.return_value = aggregate

            cmd = SuspendIdentityCommand(
                identity_id=aggregate.id,
                reason="Temporary hold short",
                actor_id="user-2",
            )
            result = await handler.handle_suspend(cmd)
            assert result["status"] == "suspended"


class TestRevokeIdentityCommand:
    """Test RevokeIdentityCommand handler."""

    @pytest.mark.asyncio
    async def test_handle_revoke(self, mock_db_session, mock_event_store, mock_kafka):
        handler = IdentityCommandHandler(mock_db_session, mock_kafka)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN555",
                first_name="C", last_name="D",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            mock_load.return_value = aggregate

            cmd = RevokeIdentityCommand(
                identity_id=aggregate.id,
                reason="Permanent revocation ordered by court",
                actor_id="user-2",
            )
            result = await handler.handle_revoke(cmd)
            assert result["status"] == "revoked"

    @pytest.mark.asyncio
    async def test_handle_revoke_twice_fails(self, mock_db_session, mock_event_store):
        handler = IdentityCommandHandler(mock_db_session)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN666",
                first_name="E", last_name="F",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            aggregate._status = "revoked"
            mock_load.return_value = aggregate

            cmd = RevokeIdentityCommand(
                identity_id=aggregate.id,
                reason="Another revocation attempt",
                actor_id="user-2",
            )
            with pytest.raises(Exception):
                await handler.handle_revoke(cmd)


class TestEnrollBiometricCommand:
    """Test EnrollBiometricCommand handler."""

    @pytest.mark.asyncio
    async def test_handle_enroll_biometric(self, mock_db_session, mock_event_store):
        handler = IdentityCommandHandler(mock_db_session)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN777",
                first_name="G", last_name="H",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            mock_load.return_value = aggregate

            cmd = EnrollBiometricCommand(
                identity_id=aggregate.id,
                biometric_type="fingerprint",
                template_hash="a" * 32,
                quality_score=0.95,
                actor_id="user-2",
            )
            result = await handler.handle_enroll_biometric(cmd)
            assert result["identity_id"] == aggregate.id

    @pytest.mark.asyncio
    async def test_handle_enroll_biometric_invalid_type(self, mock_db_session, mock_event_store):
        handler = IdentityCommandHandler(mock_db_session)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN888",
                first_name="I", last_name="J",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            mock_load.return_value = aggregate

            with pytest.raises(Exception):
                cmd = EnrollBiometricCommand(
                    identity_id=aggregate.id,
                    biometric_type="dna",
                    template_hash="a" * 32,
                    quality_score=0.9,
                    actor_id="user-2",
                )
                await handler.handle_enroll_biometric(cmd)


class TestIssueDocumentCommand:
    """Test IssueDocumentCommand handler."""

    @pytest.mark.asyncio
    async def test_handle_issue_document(self, mock_db_session, mock_event_store):
        handler = IdentityCommandHandler(mock_db_session)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN999",
                first_name="K", last_name="L",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            mock_load.return_value = aggregate

            cmd = IssueDocumentCommand(
                identity_id=aggregate.id,
                document_type="passport",
                document_number="PP001122",
                issue_date="2024-06-01",
                expiry_date="2034-06-01",
                issuing_agency="AN-AGENCY",
                actor_id="user-2",
            )
            result = await handler.handle_issue_document(cmd)
            assert result["identity_id"] == aggregate.id

    @pytest.mark.asyncio
    async def test_handle_issue_document_to_revoked(self, mock_db_session, mock_event_store):
        handler = IdentityCommandHandler(mock_db_session)

        with patch.object(handler, "_load_aggregate") as mock_load:
            aggregate = IdentityAggregate.create(
                national_id="SN000",
                first_name="M", last_name="N",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id="user-1",
            )
            aggregate._status = "revoked"
            mock_load.return_value = aggregate

            cmd = IssueDocumentCommand(
                identity_id=aggregate.id,
                document_type="passport",
                document_number="PP003344",
                issue_date="2024-06-01",
                expiry_date="2034-06-01",
                issuing_agency="AGENCY",
                actor_id="user-2",
            )
            with pytest.raises(Exception):
                await handler.handle_issue_document(cmd)


class TestCommandValidation:
    """Test command model validation."""

    def test_create_identity_command_validation(self):
        with pytest.raises(Exception):
            CreateIdentityCommand(
                national_id="AB",
                first_name="",
                last_name="",
                date_of_birth="invalid",
                place_of_birth="",
                gender="unknown",
                nationality="TOOLONG",
                agency_id="ag-1",
                actor_id="user-1",
            )

    def test_create_identity_command_valid(self, identity_create_data):
        cmd = CreateIdentityCommand(**identity_create_data)
        assert cmd.national_id == identity_create_data["national_id"]

    def test_suspend_command_min_length(self):
        cmd = SuspendIdentityCommand(
            identity_id=str(uuid.uuid4()),
            reason="Fraud investigation ongoing",
            actor_id="user-1",
        )
        assert len(cmd.reason) >= 10

    def test_enroll_biometric_command_pattern(self):
        with pytest.raises(Exception):
            EnrollBiometricCommand(
                identity_id=str(uuid.uuid4()),
                biometric_type="invalid_type",
                template_hash="a" * 32,
                quality_score=0.9,
                actor_id="user-1",
            )

    def test_enroll_biometric_quality_range(self):
        with pytest.raises(Exception):
            EnrollBiometricCommand(
                identity_id=str(uuid.uuid4()),
                biometric_type="fingerprint",
                template_hash="a" * 32,
                quality_score=999,
                actor_id="user-1",
            )
