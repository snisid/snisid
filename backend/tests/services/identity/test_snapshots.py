from __future__ import annotations

import uuid

from services.identity.aggregate import IdentityAggregate
from services.identity.events import IdentityCreated


class TestSnapshotCreation:
    def test_snapshot_contains_all_fields(self, identity_aggregate):
        snapshot = identity_aggregate.take_snapshot()
        expected_keys = {
            "national_id", "first_name", "last_name", "middle_name",
            "date_of_birth", "place_of_birth", "gender", "nationality",
            "status", "agency_id", "created_by", "address", "photo_url",
            "biometrics", "documents", "verified_at", "verified_by",
            "suspension_reason", "revocation_reason",
        }
        assert set(snapshot.keys()) == expected_keys

    def test_snapshot_reflects_initial_state(self, identity_aggregate, identity_create_data):
        snapshot = identity_aggregate.take_snapshot()
        assert snapshot["national_id"] == identity_create_data["national_id"]
        assert snapshot["first_name"] == identity_create_data["first_name"]
        assert snapshot["last_name"] == identity_create_data["last_name"]
        assert snapshot["status"] == "pending"

    def test_snapshot_version_tracking(self, identity_aggregate):
        assert identity_aggregate.version == 1

    def test_empty_snapshot_fields(self, identity_aggregate):
        snapshot = identity_aggregate.take_snapshot()
        assert snapshot["middle_name"] is None
        assert snapshot["address"] is None
        assert snapshot["photo_url"] is None
        assert snapshot["biometrics"] == []
        assert snapshot["documents"] == []


class TestSnapshotRestoration:
    def test_restore_from_snapshot(self, identity_aggregate, identity_create_data):
        snapshot = identity_aggregate.take_snapshot()
        restored = IdentityAggregate.from_snapshot(
            aggregate_id=identity_aggregate.id,
            snapshot=snapshot,
            version=identity_aggregate.version,
        )
        assert restored.id == identity_aggregate.id
        assert restored.national_id == identity_create_data["national_id"]
        assert restored.status == "pending"
        assert restored.version == 1

    def test_restored_aggregate_can_accept_commands(self, identity_aggregate):
        snapshot = identity_aggregate.take_snapshot()
        restored = IdentityAggregate.from_snapshot(
            aggregate_id=identity_aggregate.id,
            snapshot=snapshot,
            version=identity_aggregate.version,
        )
        restored.verify(
            verification_method="biometric",
            verifier_id="verifier-1",
            actor_id="actor-1",
        )
        assert restored.status == "active"
        assert restored.version == 2

    def test_snapshot_with_biometrics(self, identity_aggregate):
        identity_aggregate.enroll_biometric(
            biometric_type="fingerprint",
            template_hash="a" * 64,
            quality_score=0.95,
            actor_id="actor-1",
        )
        snapshot = identity_aggregate.take_snapshot()
        assert len(snapshot["biometrics"]) == 1
        assert snapshot["biometrics"][0]["biometric_type"] == "fingerprint"

    def test_restore_with_biometrics(self, identity_aggregate):
        identity_aggregate.enroll_biometric(
            biometric_type="iris",
            template_hash="b" * 64,
            quality_score=0.98,
            actor_id="actor-1",
        )
        snapshot = identity_aggregate.take_snapshot()
        restored = IdentityAggregate.from_snapshot(
            aggregate_id=identity_aggregate.id,
            snapshot=snapshot,
            version=identity_aggregate.version,
        )
        restored.enroll_biometric(
            biometric_type="face",
            template_hash="c" * 64,
            quality_score=0.90,
            actor_id="actor-1",
        )
        assert len(restored.pending_events) == 1


class TestEventReconstitution:
    def test_reconstitute_from_events(self, identity_create_data):
        aggregate = IdentityAggregate.create(**identity_create_data)
        events = aggregate.pending_events
        reconstituted = IdentityAggregate.from_events(aggregate.id, events)
        assert reconstituted.national_id == identity_create_data["national_id"]
        assert reconstituted.status == "pending"
        assert reconstituted.version == 1

    def test_reconstitution_preserves_event_order(self, identity_create_data):
        aggregate = IdentityAggregate.create(**identity_create_data)
        aggregate.verify(
            verification_method="biometric",
            verifier_id="verifier-1",
            actor_id="actor-1",
        )
        events = aggregate.pending_events
        reconstituted = IdentityAggregate.from_events(aggregate.id, events)
        assert reconstituted.status == "active"
        assert reconstituted.version == 2
