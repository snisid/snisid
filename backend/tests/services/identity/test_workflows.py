from __future__ import annotations

import uuid
from unittest.mock import patch

import pytest

from services.identity.aggregate import IdentityAggregate, InvalidStateError
from services.identity.commands import (
    CreateIdentityCommand,
    IdentityCommandHandler,
    SuspendIdentityCommand,
    VerifyIdentityCommand,
    RevokeIdentityCommand,
)


@pytest.fixture
def cmd_handler(mock_db_session, mock_kafka):
    return IdentityCommandHandler(session=mock_db_session, kafka_publisher=mock_kafka)


def _make_aggregate(identity_create_data):
    return IdentityAggregate.create(**identity_create_data)


class TestCreateThenVerifyWorkflow:
    async def test_full_lifecycle(self, cmd_handler, identity_create_data):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load, \
             patch.object(cmd_handler, "_persist_events") as mock_persist:
            agg = _make_aggregate(identity_create_data)
            mock_load.return_value = agg
            mock_persist.side_effect = lambda a: a.clear_pending_events()

            verify_cmd = VerifyIdentityCommand(
                identity_id=agg.id,
                verification_method="biometric",
                verifier_id="verifier-1",
                actor_id="actor-1",
            )
            verify_result = await cmd_handler.handle_verify(verify_cmd)
            assert verify_result["status"] == "active"

    async def test_create_then_suspend_then_reactivate(self, cmd_handler, identity_create_data):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load, \
             patch.object(cmd_handler, "_persist_events") as mock_persist:
            agg = _make_aggregate(identity_create_data)
            mock_load.return_value = agg
            mock_persist.side_effect = lambda a: a.clear_pending_events()

            suspend_cmd = SuspendIdentityCommand(
                identity_id=agg.id,
                reason="Temporary administrative hold",
                actor_id="actor-1",
            )
            suspend_result = await cmd_handler.handle_suspend(suspend_cmd)
            assert suspend_result["status"] == "suspended"

    async def test_cannot_verify_suspended_identity(self, cmd_handler, identity_create_data):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load, \
             patch.object(cmd_handler, "_persist_events") as mock_persist:
            agg = _make_aggregate(identity_create_data)
            agg._status = "suspended"
            mock_load.return_value = agg
            mock_persist.side_effect = lambda a: a.clear_pending_events()

            verify_cmd = VerifyIdentityCommand(
                identity_id=agg.id,
                verification_method="biometric",
                verifier_id="verifier-1",
                actor_id="actor-1",
            )
            with pytest.raises(InvalidStateError):
                await cmd_handler.handle_verify(verify_cmd)

    async def test_cannot_suspend_revoked_identity(self, cmd_handler, identity_create_data):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load, \
             patch.object(cmd_handler, "_persist_events") as mock_persist:
            agg = _make_aggregate(identity_create_data)
            agg.revoke(
                reason="Permanent revocation due to fraud",
                actor_id="actor-1",
            )
            agg.clear_pending_events()
            mock_load.return_value = agg
            mock_persist.side_effect = lambda a: a.clear_pending_events()

            suspend_cmd = SuspendIdentityCommand(
                identity_id=agg.id,
                reason="Cannot suspend revoked identity",
                actor_id="actor-1",
            )
            with pytest.raises(InvalidStateError):
                await cmd_handler.handle_suspend(suspend_cmd)


class TestVerifyAfterCreate:
    async def test_verify_after_create_changes_status(self, cmd_handler, identity_create_data):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load, \
             patch.object(cmd_handler, "_persist_events") as mock_persist:
            agg = _make_aggregate(identity_create_data)
            mock_load.return_value = agg
            mock_persist.side_effect = lambda a: a.clear_pending_events()

            verify_cmd = VerifyIdentityCommand(
                identity_id=agg.id,
                verification_method="document_review",
                verifier_id="verifier-1",
                actor_id="actor-1",
            )
            result = await cmd_handler.handle_verify(verify_cmd)
            assert result["status"] == "active"

    async def test_verify_updates_version(self, cmd_handler, identity_create_data):
        with patch.object(cmd_handler, "_load_aggregate") as mock_load, \
             patch.object(cmd_handler, "_persist_events") as mock_persist:
            agg = _make_aggregate(identity_create_data)
            mock_load.return_value = agg
            mock_persist.side_effect = lambda a: a.clear_pending_events()

            verify_cmd = VerifyIdentityCommand(
                identity_id=agg.id,
                verification_method="biometric",
                verifier_id="verifier-1",
                actor_id="actor-1",
            )
            result = await cmd_handler.handle_verify(verify_cmd)
            assert result["version"] > 1
