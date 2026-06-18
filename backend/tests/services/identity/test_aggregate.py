from __future__ import annotations

import uuid
from datetime import date, datetime, timezone

import pytest

from services.identity.aggregate import IdentityAggregate, InvalidStateError
from services.identity.events import (
    IdentityCreated,
    IdentityUpdated,
    IdentityVerified,
    IdentitySuspended,
    IdentityRevoked,
    IdentityReactivated,
    BiometricEnrolled,
    DocumentIssued,
    AddressChanged,
)


class TestIdentityAggregateCreate:
    """Test identity creation command."""

    def test_create_valid_identity(self, identity_create_data):
        aggregate = IdentityAggregate.create(**identity_create_data)
        assert aggregate.id is not None
        assert len(aggregate.id) > 0
        assert aggregate.status == "pending"
        assert aggregate.version == 1
        assert aggregate.national_id == identity_create_data["national_id"]
        assert len(aggregate.pending_events) == 1
        assert isinstance(aggregate.pending_events[0], IdentityCreated)

    def test_create_sets_fields_correctly(self, identity_create_data):
        aggregate = IdentityAggregate.create(**identity_create_data)
        assert aggregate._first_name == "John"
        assert aggregate._last_name == "Doe"
        assert aggregate._gender == "male"
        assert aggregate._nationality == "SEN"
        assert aggregate._date_of_birth == "1990-01-15"

    def test_create_rejects_future_dob(self, identity_create_data):
        data = {**identity_create_data, "date_of_birth": "2099-01-01"}
        with pytest.raises(ValueError, match="Date of birth cannot be in the future"):
            IdentityAggregate.create(**data)

    def test_create_rejects_invalid_nationality(self, identity_create_data):
        data = {**identity_create_data, "nationality": "SENEGAL"}
        with pytest.raises(ValueError, match="Nationality must be a 3-letter ISO country code"):
            IdentityAggregate.create(**data)

    def test_create_rejects_empty_nationality(self, identity_create_data):
        data = {**identity_create_data, "nationality": ""}
        with pytest.raises(ValueError):
            IdentityAggregate.create(**data)

    def test_create_trims_whitespace(self, identity_create_data):
        data = {**identity_create_data, "first_name": "  John  ", "last_name": "  Doe  "}
        aggregate = IdentityAggregate.create(**data)
        assert aggregate._first_name == "John"
        assert aggregate._last_name == "Doe"


class TestIdentityAggregateUpdate:
    """Test identity update command."""

    @pytest.fixture
    def aggregate(self, identity_create_data):
        return IdentityAggregate.create(**identity_create_data)

    def test_update_valid_fields(self, aggregate):
        aggregate.update(changes={"first_name": "Jane", "phone": "+221123456"}, actor_id="user-2")
        assert aggregate._first_name == "Jane"
        assert len(aggregate.pending_events) == 2
        assert aggregate.version == 2

    def test_update_rejects_invalid_status(self, identity_create_data):
        aggregate = IdentityAggregate.create(**identity_create_data)
        aggregate._status = "revoked"
        with pytest.raises(InvalidStateError, match="Cannot update"):
            aggregate.update(changes={"first_name": "Jane"}, actor_id="user-2")

    def test_update_filters_unknown_fields(self, aggregate):
        aggregate.update(
            changes={"first_name": "Jane", "unknown_field": "value"},
            actor_id="user-2",
        )
        assert aggregate._first_name == "Jane"

    def test_update_requires_valid_fields(self, aggregate):
        with pytest.raises(ValueError, match="No valid fields to update"):
            aggregate.update(changes={"unknown": "value"}, actor_id="user-2")

    def test_update_suspended_identity(self, identity_create_data):
        aggregate = IdentityAggregate.create(**identity_create_data)
        aggregate._status = "suspended"
        with pytest.raises(InvalidStateError):
            aggregate.update(changes={"first_name": "Jane"}, actor_id="user-2")


class TestIdentityAggregateVerify:
    """Test identity verification."""

    @pytest.fixture
    def aggregate(self, identity_create_data):
        return IdentityAggregate.create(**identity_create_data)

    def test_verify_identity(self, aggregate):
        aggregate.verify(
            verification_method="biometric",
            verifier_id="verifier-1",
            actor_id="user-2",
        )
        assert aggregate.status == "active"
        assert aggregate._verified_by == "verifier-1"
        assert aggregate._verified_at is not None
        assert len(aggregate.pending_events) == 2
        assert isinstance(aggregate.pending_events[1], IdentityVerified)

    def test_verify_revoked_identity(self, aggregate):
        aggregate._status = "revoked"
        with pytest.raises(InvalidStateError, match="Cannot verify"):
            aggregate.verify(
                verification_method="biometric",
                verifier_id="v-1",
                actor_id="user-2",
            )

    def test_verify_deceased_identity(self, aggregate):
        aggregate._status = "deceased"
        with pytest.raises(InvalidStateError):
            aggregate.verify(
                verification_method="biometric",
                verifier_id="v-1",
                actor_id="user-2",
            )


class TestIdentityAggregateSuspend:
    """Test identity suspension."""

    @pytest.fixture
    def aggregate(self, identity_create_data):
        return IdentityAggregate.create(**identity_create_data)

    def test_suspend_identity(self, aggregate):
        aggregate.suspend(reason="Fraud investigation", actor_id="user-2")
        assert aggregate.status == "suspended"
        assert aggregate._suspension_reason == "Fraud investigation"
        assert isinstance(aggregate.pending_events[1], IdentitySuspended)

    def test_suspend_requires_reason(self, aggregate):
        with pytest.raises(ValueError, match="Suspension reason is required"):
            aggregate.suspend(reason="", actor_id="user-2")

    def test_suspend_already_suspended(self, aggregate):
        aggregate._status = "suspended"
        with pytest.raises(InvalidStateError, match="Cannot suspend"):
            aggregate.suspend(reason="Another reason", actor_id="user-2")

    def test_suspend_revoked_identity(self, aggregate):
        aggregate._status = "revoked"
        with pytest.raises(InvalidStateError):
            aggregate.suspend(reason="test", actor_id="user-2")

    def test_suspend_deceased(self, aggregate):
        aggregate._status = "deceased"
        with pytest.raises(InvalidStateError):
            aggregate.suspend(reason="test", actor_id="user-2")


class TestIdentityAggregateRevoke:
    """Test identity revocation."""

    @pytest.fixture
    def aggregate(self, identity_create_data):
        return IdentityAggregate.create(**identity_create_data)

    def test_revoke_identity(self, aggregate):
        aggregate.revoke(reason="Permanent revocation", actor_id="user-2")
        assert aggregate.status == "revoked"
        assert aggregate._revocation_reason == "Permanent revocation"
        assert isinstance(aggregate.pending_events[1], IdentityRevoked)

    def test_revoke_requires_reason(self, aggregate):
        with pytest.raises(ValueError, match="Revocation reason is required"):
            aggregate.revoke(reason="", actor_id="user-2")

    def test_revoke_already_revoked(self, aggregate):
        aggregate._status = "revoked"
        with pytest.raises(InvalidStateError, match="already revoked"):
            aggregate.revoke(reason="Again", actor_id="user-2")


class TestIdentityAggregateReactivate:
    """Test identity reactivation."""

    @pytest.fixture
    def aggregate(self, identity_create_data):
        return IdentityAggregate.create(**identity_create_data)

    def test_reactivate_suspended(self, aggregate):
        aggregate._status = "suspended"
        aggregate.reactivate(reason="Investigation cleared", actor_id="user-2")
        assert aggregate.status == "active"
        assert aggregate._suspension_reason is None
        assert isinstance(aggregate.pending_events[-1], IdentityReactivated)

    def test_reactivate_non_suspended_fails(self, aggregate):
        with pytest.raises(InvalidStateError, match="Can only reactivate suspended"):
            aggregate.reactivate(reason="test", actor_id="user-2")

    def test_reactivate_revoked_fails(self, aggregate):
        aggregate._status = "revoked"
        with pytest.raises(InvalidStateError):
            aggregate.reactivate(reason="test", actor_id="user-2")


class TestIdentityAggregateBiometric:
    """Test biometric enrollment."""

    @pytest.fixture
    def aggregate(self, identity_create_data):
        return IdentityAggregate.create(**identity_create_data)

    def test_enroll_biometric(self, aggregate):
        aggregate.enroll_biometric(
            biometric_type="fingerprint",
            template_hash="hash123..." * 8,
            quality_score=0.95,
            actor_id="user-2",
        )
        assert len(aggregate._biometrics) == 1
        assert aggregate._biometrics[0]["biometric_type"] == "fingerprint"
        assert aggregate._biometrics[0]["is_primary"] is True

    def test_enroll_biometric_revoked_fails(self, aggregate):
        aggregate._status = "revoked"
        with pytest.raises(InvalidStateError, match="Cannot enroll biometrics"):
            aggregate.enroll_biometric(
                biometric_type="fingerprint",
                template_hash="hash",
                quality_score=0.9,
                actor_id="user-2",
            )

    def test_enroll_biometric_invalid_quality(self, aggregate):
        with pytest.raises(ValueError, match="Quality score must be between"):
            aggregate.enroll_biometric(
                biometric_type="iris",
                template_hash="hash",
                quality_score=1.5,
                actor_id="user-2",
            )

    @pytest.mark.parametrize("score", [-0.1, 1.1, 100])
    def test_enroll_biometric_out_of_range_scores(self, aggregate, score):
        with pytest.raises(ValueError):
            aggregate.enroll_biometric(
                biometric_type="face",
                template_hash="hash",
                quality_score=score,
                actor_id="user-2",
            )


class TestIdentityAggregateDocument:
    """Test document issuance."""

    @pytest.fixture
    def aggregate(self, identity_create_data):
        return IdentityAggregate.create(**identity_create_data)

    def test_issue_document(self, aggregate):
        aggregate.issue_document(
            document_type="national_id",
            document_number="NID123456",
            issue_date="2024-01-01",
            expiry_date="2034-01-01",
            issuing_agency="AN-AGENCY",
            actor_id="user-2",
        )
        assert len(aggregate._documents) == 1
        assert aggregate._documents[0]["document_type"] == "national_id"
        assert aggregate._documents[0]["document_number"] == "NID123456"

    def test_issue_document_revoked_fails(self, aggregate):
        aggregate._status = "revoked"
        with pytest.raises(InvalidStateError):
            aggregate.issue_document(
                document_type="passport",
                document_number="PP123",
                issue_date="2024-01-01",
                expiry_date="2034-01-01",
                issuing_agency="AGENCY",
                actor_id="user-2",
            )

    def test_issue_document_deceased_fails(self, aggregate):
        aggregate._status = "deceased"
        with pytest.raises(InvalidStateError):
            aggregate.issue_document(
                document_type="passport",
                document_number="PP123",
                issue_date="2024-01-01",
                expiry_date="2034-01-01",
                issuing_agency="AGENCY",
                actor_id="user-2",
            )


class TestIdentityAggregateAddress:
    """Test address change."""

    @pytest.fixture
    def aggregate(self, identity_create_data):
        return IdentityAggregate.create(**identity_create_data)

    def test_change_address(self, aggregate):
        new_address = {"city": "Dakar", "street": "123 Main St"}
        aggregate.change_address(new_address=new_address, actor_id="user-2")
        assert aggregate._address == new_address
        assert isinstance(aggregate.pending_events[-1], AddressChanged)

    def test_change_address_revoked_fails(self, aggregate):
        aggregate._status = "revoked"
        with pytest.raises(InvalidStateError):
            aggregate.change_address(
                new_address={"city": "Dakar"}, actor_id="user-2"
            )


class TestEventSourcingReconstitution:
    """Test rebuilding aggregate state from events."""

    def test_replay_events(self, identity_create_data):
        original = IdentityAggregate.create(**identity_create_data)
        original.verify(
            verification_method="biometric",
            verifier_id="v-1",
            actor_id="user-2",
        )
        original.suspend(reason="Fraud investigation", actor_id="user-2")

        events = original.pending_events
        rebuilt = IdentityAggregate.from_events(original.id, events)
        assert rebuilt.id == original.id
        assert rebuilt.status == "suspended"
        assert rebuilt.version == 3
        assert rebuilt._first_name == "John"
        assert rebuilt._verified_by == "v-1"

    def test_replay_full_lifecycle(self, identity_create_data):
        original = IdentityAggregate.create(**identity_create_data)
        original.verify(verification_method="document", verifier_id="v-1", actor_id="user-2")
        original.suspend(reason="Investigation", actor_id="user-2")
        original.reactivate(reason="Cleared", actor_id="user-2")
        original.revoke(reason="Permanent", actor_id="user-2")

        events = original.pending_events
        rebuilt = IdentityAggregate.from_events(original.id, events)

        assert rebuilt.status == "revoked"
        assert rebuilt.version == 5
        assert rebuilt._revocation_reason == "Permanent"

    def test_replay_with_biometrics(self, identity_create_data):
        original = IdentityAggregate.create(**identity_create_data)
        original.enroll_biometric(
            biometric_type="fingerprint",
            template_hash="hash1",
            quality_score=0.9,
            actor_id="user-2",
        )
        original.enroll_biometric(
            biometric_type="iris",
            template_hash="hash2",
            quality_score=0.85,
            actor_id="user-2",
        )

        events = original.pending_events
        rebuilt = IdentityAggregate.from_events(original.id, events)
        assert len(rebuilt._biometrics) == 2
        assert rebuilt.version == 3

    def test_replay_empty_events(self):
        aggregate = IdentityAggregate.from_events(str(uuid.uuid4()), [])
        assert aggregate.version == 0
        assert aggregate.status == "pending"

    def test_from_snapshot(self, identity_create_data):
        original = IdentityAggregate.create(**identity_create_data)
        original.verify(verification_method="biometric", verifier_id="v-1", actor_id="user-2")
        snapshot = original.take_snapshot()

        rebuilt = IdentityAggregate.from_snapshot(
            aggregate_id=original.id,
            snapshot=snapshot,
            version=original.version,
        )
        assert rebuilt.id == original.id
        assert rebuilt.status == original.status
        assert rebuilt.version == original.version
        assert rebuilt._first_name == original._first_name

    def test_take_snapshot_includes_all_fields(self, identity_create_data):
        aggregate = IdentityAggregate.create(**identity_create_data)
        snapshot = aggregate.take_snapshot()
        assert snapshot["national_id"] == identity_create_data["national_id"]
        assert snapshot["first_name"] == "John"
        assert snapshot["status"] == "pending"
        assert snapshot["biometrics"] == []
        assert snapshot["documents"] == []


class TestVersionConcurrency:
    """Test version tracking for concurrency control."""

    def test_version_increments(self, identity_create_data):
        aggregate = IdentityAggregate.create(**identity_create_data)
        assert aggregate.version == 1
        aggregate.update(changes={"first_name": "Jane"}, actor_id="user-2")
        assert aggregate.version == 2

    def test_pending_events_have_correct_versions(self, identity_create_data):
        aggregate = IdentityAggregate.create(**identity_create_data)
        aggregate.suspend(reason="Fraud", actor_id="user-2")
        aggregate.revoke(reason="Final", actor_id="user-2")
        events = aggregate.pending_events
        assert events[0].version == 1
        assert events[1].version == 2
        assert events[2].version == 3

    def test_clear_pending_events(self, identity_create_data):
        aggregate = IdentityAggregate.create(**identity_create_data)
        assert len(aggregate.pending_events) == 1
        aggregate.clear_pending_events()
        assert len(aggregate.pending_events) == 0
