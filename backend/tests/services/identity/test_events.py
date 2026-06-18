"""
SNISID Identity Domain Events Tests
=====================================
Tests for identity domain events and event deserialization.
"""
from __future__ import annotations

import uuid
from datetime import datetime, timezone
from typing import Any

import pytest
from pydantic import ValidationError

from services.identity.events import (
    IdentityEvent,
    IdentityCreated,
    IdentityUpdated,
    IdentityVerified,
    IdentitySuspended,
    IdentityRevoked,
    IdentityReactivated,
    BiometricEnrolled,
    DocumentIssued,
    AddressChanged,
    EVENT_TYPE_MAP,
    deserialize_event,
)


class TestIdentityEventBase:
    """Test base IdentityEvent class."""

    def test_default_values(self):
        event = IdentityEvent()
        assert event.event_type == ""
        assert event.aggregate_type == "Identity"
        assert event.version == 0
        assert event.data == {}
        assert event.actor_id is None
        assert event.correlation_id is None
        assert isinstance(event.timestamp, datetime)
        assert uuid.UUID(event.event_id)

    def test_event_id_unique(self):
        e1 = IdentityEvent()
        e2 = IdentityEvent()
        assert e1.event_id != e2.event_id

    def test_timestamp_is_utc(self):
        event = IdentityEvent()
        assert event.timestamp.tzinfo is not None
        assert event.timestamp.tzinfo == timezone.utc

    def test_custom_data(self):
        event = IdentityEvent(data={"key": "value", "number": 42})
        assert event.data["key"] == "value"
        assert event.data["number"] == 42

    def test_actor_and_correlation_ids(self):
        event = IdentityEvent(
            actor_id="user-001",
            correlation_id="corr-abc",
        )
        assert event.actor_id == "user-001"
        assert event.correlation_id == "corr-abc"

    def test_invalid_data_type_raises(self):
        with pytest.raises(ValidationError):
            IdentityEvent(data="not-a-dict")


class TestIdentityCreated:
    """Test IdentityCreated event."""

    def test_event_type(self):
        assert IdentityCreated.model_fields["event_type"].default == "identity.created"

    def test_create_minimal(self):
        event = IdentityCreated.create(
            aggregate_id="agg-001",
            national_id="SN123456789",
            first_name="Jean",
            last_name="Dupont",
            date_of_birth="1990-01-15",
            place_of_birth="Port-au-Prince",
            gender="M",
            nationality="HTI",
            agency_id="ONI-AGENCY",
        )
        assert event.event_type == "identity.created"
        assert event.aggregate_id == "agg-001"
        assert event.data["national_id"] == "SN123456789"
        assert event.data["first_name"] == "Jean"
        assert event.data["last_name"] == "Dupont"
        assert event.data["date_of_birth"] == "1990-01-15"
        assert event.data["place_of_birth"] == "Port-au-Prince"
        assert event.data["gender"] == "M"
        assert event.data["nationality"] == "HTI"
        assert event.data["agency_id"] == "ONI-AGENCY"

    def test_create_with_actor(self):
        event = IdentityCreated.create(
            aggregate_id="agg-002",
            national_id="SN987654321",
            first_name="Marie",
            last_name="Pierre",
            date_of_birth="1985-06-20",
            place_of_birth="Cap-Haïtien",
            gender="F",
            nationality="HTI",
            agency_id="ONI-AGENCY",
            actor_id="agent-003",
            correlation_id="corr-456",
        )
        assert event.actor_id == "agent-003"
        assert event.correlation_id == "corr-456"
        assert event.aggregate_id == "agg-002"

    def test_data_immutable(self):
        event = IdentityCreated.create(
            aggregate_id="agg-003",
            national_id="SN112233445",
            first_name="Test",
            last_name="User",
            date_of_birth="2000-01-01",
            place_of_birth="Jacmel",
            gender="M",
            nationality="HTI",
            agency_id="ONI",
        )
        assert isinstance(event.data, dict)
        assert event.data["first_name"] == "Test"


class TestIdentityUpdated:
    """Test IdentityUpdated event."""

    def test_event_type(self):
        assert IdentityUpdated.model_fields["event_type"].default == "identity.updated"

    def test_create_with_changes(self):
        event = IdentityUpdated.create(
            aggregate_id="agg-001",
            changes={"last_name": "Dupont-Michel", "address": "Rue 5, Pétion-Ville"},
        )
        assert event.data["changes"]["last_name"] == "Dupont-Michel"
        assert "address" in event.data["changes"]

    def test_create_with_actor(self):
        event = IdentityUpdated.create(
            aggregate_id="agg-001",
            changes={"phone": "+50912345678"},
            actor_id="agent-007",
        )
        assert event.actor_id == "agent-007"

    def test_empty_changes(self):
        event = IdentityUpdated.create(
            aggregate_id="agg-001",
            changes={},
        )
        assert event.data["changes"] == {}


class TestIdentityVerified:
    """Test IdentityVerified event."""

    def test_event_type(self):
        assert IdentityVerified.model_fields["event_type"].default == "identity.verified"

    def test_create_with_biometric(self):
        event = IdentityVerified.create(
            aggregate_id="agg-001",
            verification_method="biometric_face",
            verifier_id="biometric-engine-v2",
        )
        assert event.data["verification_method"] == "biometric_face"
        assert event.data["verifier_id"] == "biometric-engine-v2"

    def test_create_with_document(self):
        event = IdentityVerified.create(
            aggregate_id="agg-002",
            verification_method="document_check",
            verifier_id="agent-005",
            actor_id="agent-005",
        )
        assert event.actor_id == "agent-005"


class TestIdentitySuspended:
    """Test IdentitySuspended event."""

    def test_event_type(self):
        assert IdentitySuspended.model_fields["event_type"].default == "identity.suspended"

    def test_create_with_reason(self):
        event = IdentitySuspended.create(
            aggregate_id="agg-001",
            reason="Suspicious activity detected",
            actor_id="soc-analyst-01",
        )
        assert event.data["reason"] == "Suspicious activity detected"
        assert event.actor_id == "soc-analyst-01"


class TestIdentityRevoked:
    """Test IdentityRevoked event."""

    def test_event_type(self):
        assert IdentityRevoked.model_fields["event_type"].default == "identity.revoked"

    def test_create_with_reason(self):
        event = IdentityRevoked.create(
            aggregate_id="agg-001",
            reason="Fraudulent enrollment",
            actor_id="admin-001",
            correlation_id="corr-fraud-001",
        )
        assert event.data["reason"] == "Fraudulent enrollment"
        assert event.correlation_id == "corr-fraud-001"


class TestIdentityReactivated:
    """Test IdentityReactivated event."""

    def test_event_type(self):
        assert IdentityReactivated.model_fields["event_type"].default == "identity.reactivated"

    def test_create_with_reason(self):
        event = IdentityReactivated.create(
            aggregate_id="agg-001",
            reason="Appeal approved",
            actor_id="admin-002",
        )
        assert event.data["reason"] == "Appeal approved"


class TestBiometricEnrolled:
    """Test BiometricEnrolled event."""

    def test_event_type(self):
        assert BiometricEnrolled.model_fields["event_type"].default == "identity.biometric_enrolled"

    def test_create_with_biometric_data(self):
        event = BiometricEnrolled.create(
            aggregate_id="agg-001",
            biometric_type="face",
            template_hash="sha256:a" * 8,
            quality_score=0.95,
        )
        assert event.data["biometric_type"] == "face"
        assert event.data["quality_score"] == 0.95

    def test_create_with_low_quality(self):
        event = BiometricEnrolled.create(
            aggregate_id="agg-002",
            biometric_type="fingerprint",
            template_hash="sha256:b" * 8,
            quality_score=0.45,
        )
        assert event.data["quality_score"] == 0.45


class TestDocumentIssued:
    """Test DocumentIssued event."""

    def test_event_type(self):
        assert DocumentIssued.model_fields["event_type"].default == "identity.document_issued"

    def test_create_with_document_data(self):
        event = DocumentIssued.create(
            aggregate_id="agg-001",
            document_type="national_id_card",
            document_number="SNID-2025-001234",
            issue_date="2025-01-15",
            expiry_date="2035-01-15",
            issuing_agency="ONI",
        )
        assert event.data["document_type"] == "national_id_card"
        assert event.data["document_number"] == "SNID-2025-001234"
        assert event.data["expiry_date"] == "2035-01-15"

    def test_create_with_actor(self):
        event = DocumentIssued.create(
            aggregate_id="agg-002",
            document_type="passport",
            document_number="HT-PP-987654",
            issue_date="2025-03-01",
            expiry_date="2030-03-01",
            issuing_agency="DGI",
            actor_id="agent-010",
        )
        assert event.actor_id == "agent-010"


class TestAddressChanged:
    """Test AddressChanged event."""

    def test_event_type(self):
        assert AddressChanged.model_fields["event_type"].default == "identity.address_changed"

    def test_create_with_old_and_new(self):
        event = AddressChanged.create(
            aggregate_id="agg-001",
            old_address={"city": "Port-au-Prince", "street": "Rue 1"},
            new_address={"city": "Pétion-Ville", "street": "Rue 5"},
        )
        assert event.data["old_address"]["city"] == "Port-au-Prince"
        assert event.data["new_address"]["city"] == "Pétion-Ville"

    def test_create_without_old_address(self):
        event = AddressChanged.create(
            aggregate_id="agg-002",
            old_address=None,
            new_address={"city": "Jacmel", "street": "Rue 10"},
        )
        assert event.data["old_address"] is None
        assert event.data["new_address"]["city"] == "Jacmel"


class TestEventTypeMap:
    """Test EVENT_TYPE_MAP registry."""

    def test_all_event_types_registered(self):
        assert "identity.created" in EVENT_TYPE_MAP
        assert "identity.updated" in EVENT_TYPE_MAP
        assert "identity.verified" in EVENT_TYPE_MAP
        assert "identity.suspended" in EVENT_TYPE_MAP
        assert "identity.revoked" in EVENT_TYPE_MAP
        assert "identity.reactivated" in EVENT_TYPE_MAP
        assert "identity.biometric_enrolled" in EVENT_TYPE_MAP
        assert "identity.document_issued" in EVENT_TYPE_MAP
        assert "identity.address_changed" in EVENT_TYPE_MAP

    def test_map_count(self):
        assert len(EVENT_TYPE_MAP) == 9

    def test_map_classes_correct(self):
        assert EVENT_TYPE_MAP["identity.created"] == IdentityCreated
        assert EVENT_TYPE_MAP["identity.updated"] == IdentityUpdated
        assert EVENT_TYPE_MAP["identity.verified"] == IdentityVerified
        assert EVENT_TYPE_MAP["identity.revoked"] == IdentityRevoked


class TestDeserializeEvent:
    """Test event deserialization."""

    def test_deserialize_identity_created(self):
        data = {
            "aggregate_id": "agg-001",
            "data": {"first_name": "Jean", "national_id": "SN123"},
        }
        event = deserialize_event("identity.created", data)
        assert isinstance(event, IdentityCreated)
        assert event.aggregate_id == "agg-001"
        assert event.data["first_name"] == "Jean"

    def test_deserialize_identity_verified(self):
        data = {
            "aggregate_id": "agg-002",
            "data": {"verification_method": "biometric", "verifier_id": "engine-v1"},
        }
        event = deserialize_event("identity.verified", data)
        assert isinstance(event, IdentityVerified)
        assert event.data["verification_method"] == "biometric"

    def test_deserialize_unknown_type(self):
        data = {
            "aggregate_id": "agg-003",
            "event_type": "identity.unknown",
            "data": {},
        }
        event = deserialize_event("identity.unknown", data)
        assert isinstance(event, IdentityEvent)
        assert not isinstance(event, IdentityCreated)

    def test_deserialize_round_trip(self):
        original = IdentityCreated.create(
            aggregate_id="agg-004",
            national_id="SN999",
            first_name="Marie",
            last_name="Jean",
            date_of_birth="1995-05-10",
            place_of_birth="Les Cayes",
            gender="F",
            nationality="HTI",
            agency_id="ONI",
            actor_id="agent-010",
            correlation_id="corr-999",
        )
        data = original.model_dump()
        restored = deserialize_event("identity.created", data)
        assert restored.aggregate_id == original.aggregate_id
        assert restored.actor_id == original.actor_id
        assert restored.data["first_name"] == original.data["first_name"]

    def test_deserialize_with_none_fields(self):
        data = {
            "aggregate_id": "agg-005",
            "data": {},
        }
        event = deserialize_event("identity.created", data)
        assert event.actor_id is None
        assert event.correlation_id is None
