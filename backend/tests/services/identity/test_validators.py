from __future__ import annotations

import uuid
from datetime import date, timedelta

import pytest
from pydantic import ValidationError

from services.identity.commands import (
    CreateIdentityCommand,
    EnrollBiometricCommand,
    IssueDocumentCommand,
    SuspendIdentityCommand,
    RevokeIdentityCommand,
)
from services.identity.aggregate import IdentityAggregate


class TestCreateIdentityValidation:
    def test_valid_creation(self):
        cmd = CreateIdentityCommand(
            national_id="SN123456789",
            first_name="John",
            last_name="Doe",
            date_of_birth="1990-01-15",
            place_of_birth="Dakar",
            gender="male",
            nationality="SEN",
            agency_id=str(uuid.uuid4()),
            actor_id=str(uuid.uuid4()),
        )
        assert cmd.national_id == "SN123456789"

    def test_invalid_date_format(self):
        with pytest.raises(ValidationError):
            CreateIdentityCommand(
                national_id="SN123456789",
                first_name="John",
                last_name="Doe",
                date_of_birth="15-01-1990",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id=str(uuid.uuid4()),
            )

    def test_invalid_gender(self):
        with pytest.raises(ValidationError):
            CreateIdentityCommand(
                national_id="SN123456789",
                first_name="John",
                last_name="Doe",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="unknown",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id=str(uuid.uuid4()),
            )

    def test_invalid_nationality_length(self):
        with pytest.raises(ValidationError):
            CreateIdentityCommand(
                national_id="SN123456789",
                first_name="John",
                last_name="Doe",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SENEGAL",
                agency_id=str(uuid.uuid4()),
                actor_id=str(uuid.uuid4()),
            )

    def test_empty_first_name(self):
        with pytest.raises(ValidationError):
            CreateIdentityCommand(
                national_id="SN123456789",
                first_name="",
                last_name="Doe",
                date_of_birth="1990-01-15",
                place_of_birth="Dakar",
                gender="male",
                nationality="SEN",
                agency_id=str(uuid.uuid4()),
                actor_id=str(uuid.uuid4()),
            )


class TestBiometricValidation:
    def test_valid_biometric_enroll(self):
        cmd = EnrollBiometricCommand(
            identity_id=str(uuid.uuid4()),
            biometric_type="fingerprint",
            template_hash="a" * 64,
            quality_score=0.95,
            actor_id=str(uuid.uuid4()),
        )
        assert cmd.biometric_type == "fingerprint"
        assert cmd.quality_score == 0.95

    def test_invalid_biometric_type(self):
        with pytest.raises(ValidationError):
            EnrollBiometricCommand(
                identity_id=str(uuid.uuid4()),
                biometric_type="dna",
                template_hash="a" * 64,
                quality_score=0.95,
                actor_id=str(uuid.uuid4()),
            )

    def test_quality_score_too_low(self):
        with pytest.raises(ValidationError):
            EnrollBiometricCommand(
                identity_id=str(uuid.uuid4()),
                biometric_type="iris",
                template_hash="b" * 64,
                quality_score=-0.1,
                actor_id=str(uuid.uuid4()),
            )

    def test_quality_score_too_high(self):
        with pytest.raises(ValidationError):
            EnrollBiometricCommand(
                identity_id=str(uuid.uuid4()),
                biometric_type="face",
                template_hash="c" * 64,
                quality_score=1.5,
                actor_id=str(uuid.uuid4()),
            )

    def test_template_hash_too_short(self):
        with pytest.raises(ValidationError):
            EnrollBiometricCommand(
                identity_id=str(uuid.uuid4()),
                biometric_type="voice",
                template_hash="short",
                quality_score=0.8,
                actor_id=str(uuid.uuid4()),
            )


class TestDocumentValidation:
    def test_valid_document_issue(self):
        cmd = IssueDocumentCommand(
            identity_id=str(uuid.uuid4()),
            document_type="national_id_card",
            document_number="NIC-2024-00001",
            issue_date="2024-01-15",
            expiry_date="2034-01-15",
            issuing_agency="ANTA",
            actor_id=str(uuid.uuid4()),
        )
        assert cmd.document_type == "national_id_card"

    def test_invalid_issue_date_format(self):
        with pytest.raises(ValidationError):
            IssueDocumentCommand(
                identity_id=str(uuid.uuid4()),
                document_type="passport",
                document_number="PP-2024-001",
                issue_date="01/15/2024",
                expiry_date="2034-01-15",
                issuing_agency="ANTA",
                actor_id=str(uuid.uuid4()),
            )

    def test_empty_document_number(self):
        with pytest.raises(ValidationError):
            IssueDocumentCommand(
                identity_id=str(uuid.uuid4()),
                document_type="national_id_card",
                document_number="",
                issue_date="2024-01-15",
                expiry_date="2034-01-15",
                issuing_agency="ANTA",
                actor_id=str(uuid.uuid4()),
            )


class TestSuspendRevokeValidation:
    def test_suspend_reason_too_short(self):
        with pytest.raises(ValidationError):
            SuspendIdentityCommand(
                identity_id=str(uuid.uuid4()),
                reason="Short",
                actor_id=str(uuid.uuid4()),
            )

    def test_revoke_reason_too_short(self):
        with pytest.raises(ValidationError):
            RevokeIdentityCommand(
                identity_id=str(uuid.uuid4()),
                reason="Fraud",
                actor_id=str(uuid.uuid4()),
            )


class TestAggregateBusinessRules:
    def test_future_dob_raises_error(self, identity_create_data):
        future_date = (date.today() + timedelta(days=365)).isoformat()
        with pytest.raises(ValueError, match="future"):
            IdentityAggregate.create(
                national_id=identity_create_data["national_id"],
                first_name=identity_create_data["first_name"],
                last_name=identity_create_data["last_name"],
                date_of_birth=future_date,
                place_of_birth=identity_create_data["place_of_birth"],
                gender=identity_create_data["gender"],
                nationality=identity_create_data["nationality"],
                agency_id=identity_create_data["agency_id"],
                actor_id=identity_create_data["actor_id"],
            )

    def test_invalid_nationality_alpha(self, identity_create_data):
        with pytest.raises(ValueError, match="Nationality"):
            IdentityAggregate.create(
                national_id=identity_create_data["national_id"],
                first_name=identity_create_data["first_name"],
                last_name=identity_create_data["last_name"],
                date_of_birth=identity_create_data["date_of_birth"],
                place_of_birth=identity_create_data["place_of_birth"],
                gender=identity_create_data["gender"],
                nationality="123",
                agency_id=identity_create_data["agency_id"],
                actor_id=identity_create_data["actor_id"],
            )
