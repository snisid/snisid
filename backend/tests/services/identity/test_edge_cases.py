from __future__ import annotations

import uuid

from unittest.mock import patch

import pytest

from services.identity.aggregate import IdentityAggregate, InvalidStateError
from services.identity.commands import (
    CreateIdentityCommand,
    IdentityCommandHandler,
    IdentityNotFoundError,
    UpdateIdentityCommand,
    RevokeIdentityCommand,
    EnrollBiometricCommand,
)
from services.identity.commands import VerifyIdentityCommand


@pytest.fixture
def cmd_handler(mock_db_session, mock_kafka):
    return IdentityCommandHandler(session=mock_db_session, kafka_publisher=mock_kafka)


class TestDuplicateOperations:
    async def test_double_verify(self, cmd_handler, identity_create_data):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load, \
             patch.object(cmd_handler, "_persist_events") as mock_persist:
            agg = IdentityAggregate.create(**identity_create_data)
            mock_load.return_value = agg
            mock_persist.side_effect = lambda a: a.clear_pending_events()

            verify_cmd = VerifyIdentityCommand(
                identity_id=agg.id,
                verification_method="biometric",
                verifier_id="verifier-1",
                actor_id="actor-1",
            )
            await cmd_handler.handle_verify(verify_cmd)
            result2 = await cmd_handler.handle_verify(verify_cmd)
            assert result2["status"] == "active"

    async def test_double_revoke(self, cmd_handler, identity_create_data):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load, \
             patch.object(cmd_handler, "_persist_events") as mock_persist:
            agg = IdentityAggregate.create(**identity_create_data)
            mock_load.return_value = agg
            mock_persist.side_effect = lambda a: a.clear_pending_events()

            revoke_cmd = RevokeIdentityCommand(
                identity_id=agg.id,
                reason="Permanent revocation due to fraud investigation",
                actor_id="actor-1",
            )
            await cmd_handler.handle_revoke(revoke_cmd)

            with pytest.raises(InvalidStateError, match="already revoked"):
                await cmd_handler.handle_revoke(revoke_cmd)


class TestMultipleBiometricEnrollments:
    def test_enroll_multiple_types(self, identity_aggregate):
        types = ["fingerprint", "iris", "face", "voice"]
        for bt in types:
            identity_aggregate.enroll_biometric(
                biometric_type=bt,
                template_hash=uuid.uuid4().hex,
                quality_score=0.95,
                actor_id="actor-1",
            )
        assert len(identity_aggregate.pending_events) == len(types) + 1

    def test_enroll_duplicate_type(self, identity_aggregate):
        identity_aggregate.enroll_biometric(
            biometric_type="fingerprint",
            template_hash="a" * 64,
            quality_score=0.95,
            actor_id="actor-1",
        )
        identity_aggregate.enroll_biometric(
            biometric_type="fingerprint",
            template_hash="b" * 64,
            quality_score=0.90,
            actor_id="actor-1",
        )
        assert len(identity_aggregate.pending_events) == 3

    def test_enroll_on_revoked_raises_error(self, identity_aggregate):
        identity_aggregate.revoke(
            reason="National security concern",
            actor_id="actor-1",
        )
        with pytest.raises(InvalidStateError, match="revoked"):
            identity_aggregate.enroll_biometric(
                biometric_type="fingerprint",
                template_hash="a" * 64,
                quality_score=0.95,
                actor_id="actor-1",
            )


class TestIdentityNotFound:
    async def test_update_nonexistent_identity(self, cmd_handler):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load:
            mock_load.side_effect = IdentityNotFoundError("not found")
            update_cmd = UpdateIdentityCommand(
                identity_id=str(uuid.uuid4()),
                changes={"first_name": "Jane"},
                actor_id="actor-1",
            )
            with pytest.raises(IdentityNotFoundError):
                await cmd_handler.handle_update(update_cmd)

    async def test_verify_nonexistent_identity(self, cmd_handler):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load:
            mock_load.side_effect = IdentityNotFoundError("not found")
            verify_cmd = VerifyIdentityCommand(
                identity_id=str(uuid.uuid4()),
                verification_method="biometric",
                verifier_id="verifier-1",
                actor_id="actor-1",
            )
            with pytest.raises(IdentityNotFoundError):
                await cmd_handler.handle_verify(verify_cmd)


class TestEmptyChanges:
    async def test_update_with_no_valid_fields(self, cmd_handler, identity_create_data):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load:
            agg = IdentityAggregate.create(**identity_create_data)
            mock_load.return_value = agg
            update_cmd = UpdateIdentityCommand(
                identity_id=agg.id,
                changes={"invalid_field": "value"},
                actor_id="actor-1",
            )
            with pytest.raises(ValueError, match="No valid fields"):
                await cmd_handler.handle_update(update_cmd)

    async def test_update_with_unchanged_fields(self, cmd_handler, identity_create_data):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load:
            agg = IdentityAggregate.create(**identity_create_data)
            mock_load.return_value = agg
            update_cmd = UpdateIdentityCommand(
                identity_id=agg.id,
                changes={"first_name": "John"},
                actor_id="actor-1",
            )
            update_result = await cmd_handler.handle_update(update_cmd)
            assert update_result["identity_id"] == agg.id
