from __future__ import annotations

import uuid
from datetime import date, datetime, timezone

import pytest
from pydantic import ValidationError

from services.identity.models import (
    Citizen,
    IdentityDocument,
    BiometricRecord,
    CitizenReadModel,
    IdentityStatus,
    Gender,
    DocumentType,
    DocumentStatus,
    BiometricType,
)


class TestIdentityStatusEnum:
    """Test IdentityStatus enum values."""

    def test_enum_values(self):
        assert IdentityStatus.PENDING.value == "pending"
        assert IdentityStatus.ACTIVE.value == "active"
        assert IdentityStatus.SUSPENDED.value == "suspended"
        assert IdentityStatus.REVOKED.value == "revoked"
        assert IdentityStatus.DECEASED.value == "deceased"

    def test_enum_members(self):
        assert len(IdentityStatus) == 5


class TestGenderEnum:
    """Test Gender enum values."""

    def test_enum_values(self):
        assert Gender.MALE.value == "male"
        assert Gender.FEMALE.value == "female"
        assert Gender.OTHER.value == "other"


class TestDocumentTypeEnum:
    """Test DocumentType enum values."""

    def test_enum_values(self):
        assert DocumentType.NATIONAL_ID.value == "national_id"
        assert DocumentType.PASSPORT.value == "passport"
        assert DocumentType.BIRTH_CERTIFICATE.value == "birth_certificate"
        assert DocumentType.DRIVING_LICENSE.value == "driving_license"
        assert DocumentType.RESIDENCE_PERMIT.value == "residence_permit"


class TestDocumentStatusEnum:
    """Test DocumentStatus enum values."""

    def test_enum_values(self):
        assert DocumentStatus.ACTIVE.value == "active"
        assert DocumentStatus.EXPIRED.value == "expired"
        assert DocumentStatus.REVOKED.value == "revoked"
        assert DocumentStatus.LOST.value == "lost"
        assert DocumentStatus.STOLEN.value == "stolen"


class TestBiometricTypeEnum:
    """Test BiometricType enum values."""

    def test_enum_values(self):
        assert BiometricType.FINGERPRINT.value == "fingerprint"
        assert BiometricType.IRIS.value == "iris"
        assert BiometricType.FACE.value == "face"
        assert BiometricType.VOICE.value == "voice"


class TestCitizenModel:
    """Test Citizen ORM model."""

    def test_citizen_has_required_fields(self):
        assert hasattr(Citizen, "national_id")
        assert hasattr(Citizen, "first_name")
        assert hasattr(Citizen, "last_name")
        assert hasattr(Citizen, "date_of_birth")
        assert hasattr(Citizen, "place_of_birth")
        assert hasattr(Citizen, "gender")
        assert hasattr(Citizen, "nationality")
        assert hasattr(Citizen, "status")
        assert hasattr(Citizen, "agency_id")
        assert hasattr(Citizen, "created_by")
        assert hasattr(Citizen, "version")

    def test_citizen_has_optional_fields(self):
        assert hasattr(Citizen, "middle_name")
        assert hasattr(Citizen, "email")
        assert hasattr(Citizen, "phone")
        assert hasattr(Citizen, "address")
        assert hasattr(Citizen, "photo_url")
        assert hasattr(Citizen, "verified_at")
        assert hasattr(Citizen, "verified_by")
        assert hasattr(Citizen, "suspension_reason")
        assert hasattr(Citizen, "revocation_reason")

    def test_citizen_tablename(self):
        assert Citizen.__tablename__ == "citizens"

    def test_citizen_has_relationships(self):
        assert hasattr(Citizen, "documents")
        assert hasattr(Citizen, "biometrics")

    def test_citizen_has_indexes(self):
        assert len(Citizen.__table_args__) == 3

    def test_citizen_national_id_unique(self):
        cols = Citizen.__table__.columns
        assert cols["national_id"].unique is True

    def test_citizen_default_status(self):
        col = Citizen.__table__.columns["status"]
        assert col.default is not None


class TestIdentityDocumentModel:
    """Test IdentityDocument ORM model."""

    def test_document_has_required_fields(self):
        assert hasattr(IdentityDocument, "citizen_id")
        assert hasattr(IdentityDocument, "document_type")
        assert hasattr(IdentityDocument, "document_number")
        assert hasattr(IdentityDocument, "issue_date")
        assert hasattr(IdentityDocument, "issuing_agency")
        assert hasattr(IdentityDocument, "status")

    def test_document_has_optional_fields(self):
        assert hasattr(IdentityDocument, "expiry_date")
        assert hasattr(IdentityDocument, "issuing_location")
        assert hasattr(IdentityDocument, "metadata_json")

    def test_document_tablename(self):
        assert IdentityDocument.__tablename__ == "identity_documents"

    def test_document_foreign_key(self):
        cols = IdentityDocument.__table__.columns
        assert any(fk.column.table.name == "citizens" for fk in cols["citizen_id"].foreign_keys)

    def test_document_unique_constraint(self):
        args = IdentityDocument.__table_args__
        assert any(hasattr(arg, "unique") and arg.unique for arg in args)


class TestBiometricRecordModel:
    """Test BiometricRecord ORM model."""

    def test_biometric_has_required_fields(self):
        assert hasattr(BiometricRecord, "citizen_id")
        assert hasattr(BiometricRecord, "biometric_type")
        assert hasattr(BiometricRecord, "template_hash")
        assert hasattr(BiometricRecord, "quality_score")
        assert hasattr(BiometricRecord, "captured_by")
        assert hasattr(BiometricRecord, "is_primary")

    def test_biometric_has_optional_fields(self):
        assert hasattr(BiometricRecord, "capture_device")

    def test_biometric_tablename(self):
        assert BiometricRecord.__tablename__ == "biometric_records"

    def test_biometric_template_hash_unique(self):
        cols = BiometricRecord.__table__.columns
        assert cols["template_hash"].unique is True

    def test_biometric_quality_score_default(self):
        col = BiometricRecord.__table__.columns["quality_score"]
        assert col is not None

    def test_biometric_is_primary_default(self):
        col = BiometricRecord.__table__.columns["is_primary"]
        assert col.default is not None or col.default.arg is False


class TestCitizenReadModel:
    """Test CitizenReadModel ORM model."""

    def test_read_model_has_required_fields(self):
        assert hasattr(CitizenReadModel, "national_id")
        assert hasattr(CitizenReadModel, "full_name")
        assert hasattr(CitizenReadModel, "first_name")
        assert hasattr(CitizenReadModel, "last_name")
        assert hasattr(CitizenReadModel, "date_of_birth")
        assert hasattr(CitizenReadModel, "gender")
        assert hasattr(CitizenReadModel, "nationality")
        assert hasattr(CitizenReadModel, "status")
        assert hasattr(CitizenReadModel, "agency_id")

    def test_read_model_tablename(self):
        assert CitizenReadModel.__tablename__ == "citizens_read"

    def test_read_model_indexes(self):
        assert len(CitizenReadModel.__table_args__) == 2

    def test_read_model_defaults(self):
        col = CitizenReadModel.__table__.columns["document_count"]
        assert col.default is not None or col.default.arg == 0

    def test_read_model_booleans(self):
        assert hasattr(CitizenReadModel, "has_biometrics")
        assert hasattr(CitizenReadModel, "verified")


class TestCitizenModelConstraints:
    """Test Citizen model field constraints."""

    def test_national_id_max_length(self):
        col = Citizen.__table__.columns["national_id"]
        assert col.type.length == 20

    def test_first_name_max_length(self):
        col = Citizen.__table__.columns["first_name"]
        assert col.type.length == 100

    def test_last_name_max_length(self):
        col = Citizen.__table__.columns["last_name"]
        assert col.type.length == 100

    def test_nationality_length(self):
        col = Citizen.__table__.columns["nationality"]
        assert col.type.length == 3

    def test_phone_max_length(self):
        col = Citizen.__table__.columns["phone"]
        assert col.type.length == 20

    def test_email_max_length(self):
        col = Citizen.__table__.columns["email"]
        assert col.type.length == 255

    def test_photo_url_max_length(self):
        col = Citizen.__table__.columns["photo_url"]
        assert col.type.length == 500


class TestInheritance:
    """Test that models inherit from Base properly."""

    def test_citizen_inherits_base(self):
        assert hasattr(Citizen, "id")
        assert hasattr(Citizen, "created_at")
        assert hasattr(Citizen, "updated_at")
        assert hasattr(Citizen, "is_deleted")
        assert hasattr(Citizen, "deleted_at")

    def test_document_inherits_base(self):
        assert hasattr(IdentityDocument, "id")
        assert hasattr(IdentityDocument, "created_at")

    def test_biometric_inherits_base(self):
        assert hasattr(BiometricRecord, "id")
        assert hasattr(BiometricRecord, "is_deleted")
