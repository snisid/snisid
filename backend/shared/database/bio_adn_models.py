from __future__ import annotations

import uuid
from datetime import datetime, timezone

from sqlalchemy import Boolean, Date, DateTime, Float, ForeignKey, Integer, Numeric, SmallInteger, String, Text, UniqueConstraint
from sqlalchemy.dialects.postgresql import JSONB, UUID
from sqlalchemy.orm import Mapped, mapped_column

from shared.database import Base


class BioLaboratory(Base):
    __tablename__ = "bio_laboratories"

    lab_code: Mapped[str] = mapped_column(String(20), unique=True, nullable=False)
    lab_name: Mapped[str] = mapped_column(String(200), nullable=False)
    lab_level: Mapped[str] = mapped_column(String(10), nullable=False)
    department: Mapped[str | None] = mapped_column(String(50), nullable=True)
    institution: Mapped[str | None] = mapped_column(String(100), nullable=True)
    accreditation: Mapped[str | None] = mapped_column(String(100), nullable=True)
    contact_email: Mapped[str | None] = mapped_column(String(200), nullable=True)
    is_active: Mapped[bool] = mapped_column(Boolean, default=True)
    accreditation_body: Mapped[str | None] = mapped_column(String(100), nullable=True)
    accreditation_expiry: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    last_external_audit: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    external_quality_check_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)


class BioLabEquipment(Base):
    __tablename__ = "bio_lab_equipment"

    lab_code: Mapped[str] = mapped_column(String(20), ForeignKey("bio_laboratories.id"), nullable=False)
    equipment_name: Mapped[str] = mapped_column(String(200), nullable=False)
    model: Mapped[str | None] = mapped_column(String(200), nullable=True)
    serial_number: Mapped[str] = mapped_column(String(100), nullable=False)
    role: Mapped[str] = mapped_column(String(100), nullable=False)
    calibration_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    calibration_due: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")


class BioStaffTraining(Base):
    __tablename__ = "bio_staff_training"

    staff_niu: Mapped[str] = mapped_column(String(20), nullable=False)
    training_name: Mapped[str] = mapped_column(String(200), nullable=False)
    training_code: Mapped[str] = mapped_column(String(50), nullable=False)
    duration_hours: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    completed_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    valid_until: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    issued_by: Mapped[str] = mapped_column(String(100), nullable=False)
    frequency: Mapped[str | None] = mapped_column(String(20), nullable=True)


class BioSTRProfile(Base):
    __tablename__ = "bio_str_profiles"

    specimen_number: Mapped[str] = mapped_column(String(100), unique=True, nullable=False)
    index_type: Mapped[str] = mapped_column(String(10), nullable=False)
    loci_encrypted: Mapped[bytes] = mapped_column(nullable=False)
    loci_hash: Mapped[str] = mapped_column(String(64), nullable=False)
    amelogenin: Mapped[str | None] = mapped_column(String(2), nullable=True)
    quality_score: Mapped[float | None] = mapped_column(Numeric(4, 3), nullable=True)
    loci_count: Mapped[int | None] = mapped_column(SmallInteger, default=20)
    lab_id: Mapped[str | None] = mapped_column(String(36), ForeignKey("bio_laboratories.id"), nullable=True)
    case_number: Mapped[str | None] = mapped_column(String(100), nullable=True)
    collected_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    analysis_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    uploaded_ldis: Mapped[bool] = mapped_column(Boolean, default=False)
    uploaded_sdis: Mapped[bool] = mapped_column(Boolean, default=False)
    uploaded_ndis: Mapped[bool] = mapped_column(Boolean, default=False)
    ndis_upload_date: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)
    is_expunged: Mapped[bool] = mapped_column(Boolean, default=False)
    expunge_date: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)
    expunge_order: Mapped[str | None] = mapped_column(String(200), nullable=True)


class BioIdentityLink(Base):
    __tablename__ = "bio_identity_links"

    sample_id: Mapped[str] = mapped_column(String(36), ForeignKey("bio_str_profiles.id"), unique=True, nullable=False)
    niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    linked_by_agent: Mapped[str] = mapped_column(String(100), nullable=False)
    linked_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), default=lambda: datetime.now(timezone.utc))
    court_order_ref: Mapped[str] = mapped_column(String(200), nullable=False)
    purpose: Mapped[str] = mapped_column(String(100), nullable=False)
    reviewed_by: Mapped[str | None] = mapped_column(String(100), nullable=True)
    reviewed_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)
    review_outcome: Mapped[str | None] = mapped_column(String(20), nullable=True)


class BioHit(Base):
    __tablename__ = "bio_hits"

    query_sample_id: Mapped[str] = mapped_column(String(36), ForeignKey("bio_str_profiles.id"), nullable=False)
    match_sample_id: Mapped[str] = mapped_column(String(36), ForeignKey("bio_str_profiles.id"), nullable=False)
    match_type: Mapped[str] = mapped_column(String(20), nullable=False)
    confidence: Mapped[float] = mapped_column(Numeric(5, 4), nullable=False)
    matched_loci: Mapped[int] = mapped_column(SmallInteger, nullable=False)
    total_loci: Mapped[int] = mapped_column(SmallInteger, nullable=False)
    hit_level: Mapped[str] = mapped_column(String(10), nullable=False)
    alert_sent: Mapped[bool] = mapped_column(Boolean, default=False)
    alert_sent_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)


class PerWantedPerson(Base):
    __tablename__ = "per_wanted_persons"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    last_name: Mapped[str | None] = mapped_column(String(100), nullable=True)
    first_name: Mapped[str | None] = mapped_column(String(100), nullable=True)
    aliases: Mapped[dict | None] = mapped_column(JSONB, default=list)
    date_of_birth: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    gender: Mapped[str | None] = mapped_column(String(1), nullable=True)
    nationality: Mapped[str | None] = mapped_column(String(3), nullable=True)
    warrant_type: Mapped[str] = mapped_column(String(50), nullable=False)
    warrant_number: Mapped[str | None] = mapped_column(String(100), nullable=True)
    issuing_court: Mapped[str | None] = mapped_column(String(200), nullable=True)
    issuing_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    charges: Mapped[dict] = mapped_column(JSONB, nullable=False)
    danger_level: Mapped[str] = mapped_column(String(10), default="MEDIUM")
    armed_dangerous: Mapped[bool] = mapped_column(Boolean, default=False)
    height_cm: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    weight_kg: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    eye_color: Mapped[str | None] = mapped_column(String(30), nullable=True)
    hair_color: Mapped[str | None] = mapped_column(String(30), nullable=True)
    distinguishing_marks: Mapped[str | None] = mapped_column(Text, nullable=True)
    entering_agency: Mapped[str] = mapped_column(String(100), nullable=False)
    entering_officer: Mapped[str] = mapped_column(String(100), nullable=False)
    mco_contact: Mapped[str | None] = mapped_column(String(200), nullable=True)
    last_known_location: Mapped[str | None] = mapped_column(String(200), nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")
    expiry_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    fingerprint_ref: Mapped[str | None] = mapped_column(String(36), nullable=True)
    photo_refs: Mapped[dict | None] = mapped_column(JSONB, default=list)
    bio_sample_ref: Mapped[str | None] = mapped_column(String(36), ForeignKey("bio_str_profiles.id"), nullable=True)
    interpol_notice: Mapped[str | None] = mapped_column(String(50), nullable=True)
    last_hit_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)


class PerMissingPerson(Base):
    __tablename__ = "per_missing_persons"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    last_name: Mapped[str] = mapped_column(String(100), nullable=False)
    first_name: Mapped[str] = mapped_column(String(100), nullable=False)
    date_of_birth: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    age_at_missing: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    gender: Mapped[str | None] = mapped_column(String(1), nullable=True)
    nationality: Mapped[str | None] = mapped_column(String(3), nullable=True)
    category: Mapped[str] = mapped_column(String(20), nullable=False)
    missing_date: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    missing_location: Mapped[str] = mapped_column(String(200), nullable=False)
    circumstances: Mapped[str | None] = mapped_column(Text, nullable=True)
    last_seen_clothing: Mapped[str | None] = mapped_column(String(500), nullable=True)
    height_cm: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    weight_kg: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    distinctive_features: Mapped[str | None] = mapped_column(Text, nullable=True)
    family_contact: Mapped[str | None] = mapped_column(String(200), nullable=True)
    family_phone: Mapped[str | None] = mapped_column(String(50), nullable=True)
    photo_refs: Mapped[dict | None] = mapped_column(JSONB, default=list)
    bio_sample_ref: Mapped[str | None] = mapped_column(String(36), ForeignKey("bio_str_profiles.id"), nullable=True)
    family_bio_refs: Mapped[dict | None] = mapped_column(JSONB, default=list)
    medical_conditions: Mapped[str | None] = mapped_column(Text, nullable=True)
    medications: Mapped[str | None] = mapped_column(Text, nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")
    located_date: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)
    entering_agency: Mapped[str] = mapped_column(String(100), nullable=False)
    ncmec_notified: Mapped[bool] = mapped_column(Boolean, default=False)
    citizen_portal_submission: Mapped[bool] = mapped_column(Boolean, default=False)
    bpm_notified: Mapped[bool] = mapped_column(Boolean, default=False)
    auto_cross_bio_dis: Mapped[bool] = mapped_column(Boolean, default=False)


class PerSexOffender(Base):
    __tablename__ = "per_sex_offenders"

    niu: Mapped[str] = mapped_column(String(20), nullable=False)
    conviction_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    conviction_court: Mapped[str] = mapped_column(String(200), nullable=False)
    offenses: Mapped[dict] = mapped_column(JSONB, nullable=False)
    risk_level: Mapped[str | None] = mapped_column(String(10), nullable=True)
    registration_expiry: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    current_address: Mapped[str | None] = mapped_column(Text, nullable=True)
    employer: Mapped[str | None] = mapped_column(String(200), nullable=True)
    restrictions: Mapped[str | None] = mapped_column(Text, nullable=True)
    address_declared: Mapped[str | None] = mapped_column(Text, nullable=True)
    geographic_restrictions: Mapped[str | None] = mapped_column(Text, nullable=True)
    last_verified: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")


class PerGangMember(Base):
    __tablename__ = "per_gang_members"

    niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    last_name: Mapped[str | None] = mapped_column(String(100), nullable=True)
    first_name: Mapped[str | None] = mapped_column(String(100), nullable=True)
    aliases: Mapped[dict | None] = mapped_column(JSONB, default=list)
    gang_name: Mapped[str] = mapped_column(String(200), nullable=False)
    gang_code: Mapped[str | None] = mapped_column(String(50), nullable=True)
    membership_type: Mapped[str | None] = mapped_column(String(30), nullable=True)
    territory: Mapped[str | None] = mapped_column(String(200), nullable=True)
    known_weapons: Mapped[dict | None] = mapped_column(JSONB, default=list)
    criminal_activities: Mapped[dict | None] = mapped_column(JSONB, default=list)
    threat_level: Mapped[str] = mapped_column(String(10), default="HIGH")
    intelligence_notes: Mapped[str | None] = mapped_column(Text, nullable=True)
    source_reliability: Mapped[str | None] = mapped_column(String(10), nullable=True)
    last_review_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    auto_removal_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")


class PerForeignFugitive(Base):
    __tablename__ = "per_foreign_fugitives"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    interpol_notice_number: Mapped[str] = mapped_column(String(50), nullable=False)
    notice_type: Mapped[str] = mapped_column(String(20), nullable=False)
    last_name: Mapped[str] = mapped_column(String(100), nullable=False)
    first_name: Mapped[str | None] = mapped_column(String(100), nullable=True)
    aliases: Mapped[dict | None] = mapped_column(JSONB, default=list)
    date_of_birth: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    gender: Mapped[str | None] = mapped_column(String(1), nullable=True)
    nationality: Mapped[str | None] = mapped_column(String(3), nullable=True)
    charges: Mapped[dict] = mapped_column(JSONB, nullable=False)
    issuing_country: Mapped[str] = mapped_column(String(100), nullable=False)
    entering_agency: Mapped[str] = mapped_column(String(100), nullable=False)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")


class PerUnidentifiedPerson(Base):
    __tablename__ = "per_unidentified_persons"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    discovery_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    discovery_location: Mapped[str] = mapped_column(String(200), nullable=False)
    discovery_department: Mapped[str | None] = mapped_column(String(50), nullable=True)
    estimated_age_min: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    estimated_age_max: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    gender: Mapped[str | None] = mapped_column(String(1), nullable=True)
    estimated_height_cm: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    estimated_weight_kg: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    hair_color: Mapped[str | None] = mapped_column(String(50), nullable=True)
    eye_color: Mapped[str | None] = mapped_column(String(50), nullable=True)
    distinctive_features: Mapped[str | None] = mapped_column(Text, nullable=True)
    clothing_description: Mapped[str | None] = mapped_column(Text, nullable=True)
    dna_sample_ref: Mapped[str | None] = mapped_column(String(36), nullable=True)
    fingerprint_ref: Mapped[str | None] = mapped_column(String(36), nullable=True)
    dental_records_ref: Mapped[str | None] = mapped_column(String(36), nullable=True)
    photo_refs: Mapped[dict | None] = mapped_column(JSONB, default=list)
    entering_agency: Mapped[str] = mapped_column(String(100), nullable=False)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")
    matched_to_niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    matched_date: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)


class PerTerrorismWatch(Base):
    __tablename__ = "per_terrorism_watch"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    last_name: Mapped[str] = mapped_column(String(100), nullable=False)
    first_name: Mapped[str | None] = mapped_column(String(100), nullable=True)
    aliases: Mapped[dict | None] = mapped_column(JSONB, default=list)
    date_of_birth: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    nationality: Mapped[str | None] = mapped_column(String(3), nullable=True)
    risk_level: Mapped[str] = mapped_column(String(10), default="HIGH")
    threat_type: Mapped[str] = mapped_column(String(100), nullable=False)
    groups: Mapped[dict | None] = mapped_column(JSONB, default=list)
    known_associates: Mapped[dict | None] = mapped_column(JSONB, default=list)
    last_known_location: Mapped[str | None] = mapped_column(String(200), nullable=True)
    entering_agency: Mapped[str] = mapped_column(String(100), nullable=False)
    approved_by_director: Mapped[str] = mapped_column(String(100), nullable=False)
    approved_by_pg: Mapped[str] = mapped_column(String(100), nullable=False)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")


class PerProtectionOrder(Base):
    __tablename__ = "per_protection_orders"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    order_type: Mapped[str] = mapped_column(String(30), nullable=False)
    issuing_court: Mapped[str] = mapped_column(String(200), nullable=False)
    issuing_judge: Mapped[str] = mapped_column(String(100), nullable=False)
    beneficiary_niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    beneficiary_name: Mapped[str] = mapped_column(String(100), nullable=False)
    protected_person: Mapped[str] = mapped_column(String(100), nullable=False)
    restrained_person: Mapped[str] = mapped_column(String(100), nullable=False)
    restrictions: Mapped[dict] = mapped_column(JSONB, nullable=False)
    issue_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    expiration_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    emergency_contact: Mapped[str | None] = mapped_column(String(100), nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")


class PerSupervisedRelease(Base):
    __tablename__ = "per_supervised_releases"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    niu: Mapped[str] = mapped_column(String(20), nullable=False)
    last_name: Mapped[str] = mapped_column(String(100), nullable=False)
    first_name: Mapped[str | None] = mapped_column(String(100), nullable=True)
    supervision_type: Mapped[str] = mapped_column(String(30), nullable=False)
    start_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    end_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    conditions: Mapped[dict] = mapped_column(JSONB, nullable=False)
    supervising_officer: Mapped[str] = mapped_column(String(100), nullable=False)
    supervising_agency: Mapped[str] = mapped_column(String(100), nullable=False)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")


class NdisCrossDeptHit(Base):
    __tablename__ = "ndis_cross_dept_hits"

    query_sample_id: Mapped[str] = mapped_column(String(36), nullable=False)
    match_sample_id: Mapped[str] = mapped_column(String(36), nullable=False)
    match_type: Mapped[str] = mapped_column(String(20), nullable=False)
    confidence: Mapped[float] = mapped_column(Numeric(5, 4), nullable=False)
    query_sdis: Mapped[str] = mapped_column(String(20), nullable=False)
    match_sdis: Mapped[str] = mapped_column(String(20), nullable=False)
    alert_level: Mapped[str] = mapped_column(String(10), default="HIGH")
    notified_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)


class NdisReport(Base):
    __tablename__ = "ndis_reports"

    report_type: Mapped[str] = mapped_column(String(20), nullable=False)
    status: Mapped[str] = mapped_column(String(20), default="GENERATED")
    file_path: Mapped[str | None] = mapped_column(String(500), nullable=True)


class NdisInterpolSubmission(Base):
    __tablename__ = "ndis_interpol_submissions"

    sample_ids: Mapped[dict] = mapped_column(JSONB, nullable=False)
    reason: Mapped[str] = mapped_column(String(50), nullable=False)
    case_number: Mapped[str | None] = mapped_column(String(100), nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="PENDING")


class BieStolenVehicle(Base):
    __tablename__ = "bie_stolen_vehicles"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    vin: Mapped[str | None] = mapped_column(String(17), unique=True, nullable=True)
    plate_number: Mapped[str | None] = mapped_column(String(20), nullable=True)
    plate_dept: Mapped[str | None] = mapped_column(String(50), nullable=True)
    vehicle_make: Mapped[str | None] = mapped_column(String(100), nullable=True)
    vehicle_model: Mapped[str | None] = mapped_column(String(100), nullable=True)
    vehicle_year: Mapped[int | None] = mapped_column(SmallInteger, nullable=True)
    vehicle_color: Mapped[str | None] = mapped_column(String(50), nullable=True)
    vehicle_type: Mapped[str | None] = mapped_column(String(50), nullable=True)
    theft_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    theft_location: Mapped[str] = mapped_column(String(200), nullable=False)
    theft_department: Mapped[str | None] = mapped_column(String(50), nullable=True)
    owner_niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    owner_name: Mapped[str | None] = mapped_column(String(200), nullable=True)
    owner_phone: Mapped[str | None] = mapped_column(String(50), nullable=True)
    foves_record_id: Mapped[str | None] = mapped_column(String(36), nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="STOLEN")
    recovered_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    recovered_location: Mapped[str | None] = mapped_column(String(200), nullable=True)
    entering_agency: Mapped[str] = mapped_column(String(100), nullable=False)


class BieStolenFirearm(Base):
    __tablename__ = "bie_stolen_firearms"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    serial_number: Mapped[str] = mapped_column(String(100), unique=True, nullable=False)
    make: Mapped[str | None] = mapped_column(String(100), nullable=True)
    model: Mapped[str | None] = mapped_column(String(100), nullable=True)
    caliber: Mapped[str | None] = mapped_column(String(50), nullable=True)
    firearm_type: Mapped[str | None] = mapped_column(String(50), nullable=True)
    barrel_length: Mapped[float | None] = mapped_column(Float, nullable=True)
    theft_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    theft_location: Mapped[str | None] = mapped_column(String(200), nullable=True)
    owner_niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="STOLEN")
    recovered_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    entering_agency: Mapped[str] = mapped_column(String(100), nullable=False)


class BieStolenDocument(Base):
    __tablename__ = "bie_stolen_documents"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    document_type: Mapped[str] = mapped_column(String(50), nullable=False)
    document_number: Mapped[str | None] = mapped_column(String(100), nullable=True)
    issuing_agency: Mapped[str | None] = mapped_column(String(100), nullable=True)
    issue_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    expiry_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    owner_niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    owner_name: Mapped[str | None] = mapped_column(String(200), nullable=True)
    report_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    report_location: Mapped[str | None] = mapped_column(String(200), nullable=True)
    theft_type: Mapped[str] = mapped_column(String(20), default="STOLEN")
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")


class BieStolenVessel(Base):
    __tablename__ = "bie_stolen_vessels"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    vessel_name: Mapped[str | None] = mapped_column(String(200), nullable=True)
    registration_number: Mapped[str | None] = mapped_column(String(100), nullable=True)
    hull_id_number: Mapped[str | None] = mapped_column(String(50), nullable=True)
    vessel_type: Mapped[str | None] = mapped_column(String(50), nullable=True)
    vessel_make: Mapped[str | None] = mapped_column(String(100), nullable=True)
    vessel_length_m: Mapped[float | None] = mapped_column(Float, nullable=True)
    hull_color: Mapped[str | None] = mapped_column(String(50), nullable=True)
    home_port: Mapped[str | None] = mapped_column(String(200), nullable=True)
    engine_serial: Mapped[str | None] = mapped_column(String(100), nullable=True)
    distinctive_marks: Mapped[str | None] = mapped_column(Text, nullable=True)
    theft_location: Mapped[str] = mapped_column(String(200), nullable=False)
    theft_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    owner_niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    owner_name: Mapped[str | None] = mapped_column(String(200), nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="STOLEN")


class BieStolenArticle(Base):
    __tablename__ = "bie_stolen_articles"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    category: Mapped[str] = mapped_column(String(30), nullable=False)
    description: Mapped[str] = mapped_column(Text, nullable=False)
    serial_number: Mapped[str | None] = mapped_column(String(100), nullable=True)
    estimated_value: Mapped[float | None] = mapped_column(Numeric(12, 2), nullable=True)
    currency_code: Mapped[str] = mapped_column(String(3), default="HTG")
    theft_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    theft_location: Mapped[str] = mapped_column(String(200), nullable=False)
    theft_department: Mapped[str | None] = mapped_column(String(50), nullable=True)
    owner_niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    owner_name: Mapped[str | None] = mapped_column(String(200), nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="STOLEN")
    recovered_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    entering_agency: Mapped[str] = mapped_column(String(100), nullable=False)


class BieStolenSecurity(Base):
    __tablename__ = "bie_stolen_securities"

    record_number: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    security_type: Mapped[str] = mapped_column(String(30), nullable=False)
    issuer: Mapped[str] = mapped_column(String(200), nullable=False)
    security_number: Mapped[str] = mapped_column(String(100), nullable=False)
    face_value: Mapped[float | None] = mapped_column(Numeric(12, 2), nullable=True)
    currency_code: Mapped[str] = mapped_column(String(3), default="HTG")
    issue_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    theft_date: Mapped[datetime] = mapped_column(Date, nullable=False)
    theft_location: Mapped[str] = mapped_column(String(200), nullable=False)
    owner_niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    owner_name: Mapped[str | None] = mapped_column(String(200), nullable=True)
    status: Mapped[str] = mapped_column(String(20), default="STOLEN")
    recovered_date: Mapped[datetime | None] = mapped_column(Date, nullable=True)
    entering_agency: Mapped[str] = mapped_column(String(100), nullable=False)


class BioAuditLog(Base):
    __tablename__ = "bio_audit_log"

    event_type: Mapped[str] = mapped_column(String(100), nullable=False)
    table_name: Mapped[str | None] = mapped_column(String(100), nullable=True)
    record_id: Mapped[str | None] = mapped_column(String(36), nullable=True)
    officer_niu: Mapped[str] = mapped_column(String(20), nullable=False)
    agency_code: Mapped[str] = mapped_column(String(50), nullable=False)
    purpose: Mapped[str] = mapped_column(String(200), nullable=False)
    case_number: Mapped[str | None] = mapped_column(String(100), nullable=True)
    ip_hash: Mapped[str | None] = mapped_column(String(64), nullable=True)
    action: Mapped[str] = mapped_column(String(20), nullable=False)
    details: Mapped[dict | None] = mapped_column(JSONB, nullable=True)
    signature: Mapped[str | None] = mapped_column(Text, nullable=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)


class SdisNode(Base):
    __tablename__ = "sdis_nodes"

    code: Mapped[str] = mapped_column(String(20), unique=True, nullable=False)
    department: Mapped[str] = mapped_column(String(50), nullable=False)
    dc_location: Mapped[str] = mapped_column(String(100), nullable=False)
    dc_type: Mapped[str] = mapped_column(String(20), nullable=False)
    lab_codes: Mapped[dict] = mapped_column(JSONB, default=list)
    status: Mapped[str] = mapped_column(String(20), default="ACTIVE")
    last_heartbeat: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)
    is_primary: Mapped[bool] = mapped_column(Boolean, default=True)


class SdisMatch(Base):
    __tablename__ = "sdis_matches"

    sdis_code: Mapped[str] = mapped_column(String(20), nullable=False)
    query_sample_id: Mapped[str] = mapped_column(String(36), nullable=False)
    match_sample_id: Mapped[str] = mapped_column(String(36), nullable=False)
    match_type: Mapped[str] = mapped_column(String(20), nullable=False)
    confidence: Mapped[float] = mapped_column(Numeric(5, 4), nullable=False)
    alerted: Mapped[bool] = mapped_column(Boolean, default=False)
    alerted_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)


class SdisSyncError(Base):
    __tablename__ = "sdis_sync_errors"

    sdis_code: Mapped[str] = mapped_column(String(20), nullable=False)
    error_type: Mapped[str] = mapped_column(String(30), nullable=False)
    details: Mapped[str | None] = mapped_column(Text, nullable=True)
    retry_count: Mapped[int] = mapped_column(default=0)
    resolved: Mapped[bool] = mapped_column(Boolean, default=False)
    resolved_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)


class SdisQualityReview(Base):
    __tablename__ = "sdis_quality_reviews"

    sample_id: Mapped[str] = mapped_column(String(36), nullable=False)
    sdis_code: Mapped[str] = mapped_column(String(20), nullable=False)
    quality_score: Mapped[float] = mapped_column(Numeric(4, 3), nullable=False)
    reason: Mapped[str] = mapped_column(String(100), nullable=False)
    reviewed: Mapped[bool] = mapped_column(Boolean, default=False)
    reviewed_by: Mapped[str | None] = mapped_column(String(100), nullable=True)
    reviewed_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)


class PerViolenceRecord(Base):
    __tablename__ = "per_violence_records"

    record_number: Mapped[str] = mapped_column(String(30), unique=True, nullable=False, index=True)
    niu: Mapped[str | None] = mapped_column(String(20), nullable=True, index=True)
    last_name: Mapped[str | None] = mapped_column(String(100), nullable=True)
    first_name: Mapped[str | None] = mapped_column(String(100), nullable=True)
    incident_type: Mapped[str] = mapped_column(String(30), nullable=False)
    incident_date: Mapped[str] = mapped_column(String(20), nullable=False)
    location: Mapped[str] = mapped_column(String(200), nullable=False)
    victim_niu: Mapped[str | None] = mapped_column(String(20), nullable=True)
    victim_name: Mapped[str | None] = mapped_column(String(200), nullable=True)
    arresting_agency: Mapped[str] = mapped_column(String(50), nullable=False)
    court_case_ref: Mapped[str | None] = mapped_column(String(50), nullable=True)
    risk_level: Mapped[str] = mapped_column(String(10), nullable=False, default="MEDIUM")
    status: Mapped[str] = mapped_column(String(20), nullable=False, default="ACTIVE")


class PerIdentityTheft(Base):
    __tablename__ = "per_identity_thefts"

    record_number: Mapped[str] = mapped_column(String(30), unique=True, nullable=False, index=True)
    victim_niu: Mapped[str] = mapped_column(String(20), nullable=False, index=True)
    court_order: Mapped[str | None] = mapped_column(String(200), nullable=True)
    victim_name: Mapped[str | None] = mapped_column(String(200), nullable=True)
    fraud_type: Mapped[str] = mapped_column(String(30), nullable=False)
    document_type_used: Mapped[str | None] = mapped_column(String(50), nullable=True)
    perpetrator_known: Mapped[bool] = mapped_column(Boolean, default=False)
    perpetrator_name: Mapped[str | None] = mapped_column(String(200), nullable=True)
    report_date: Mapped[str] = mapped_column(String(20), nullable=False)
    reporting_agency: Mapped[str] = mapped_column(String(50), nullable=False)
    status: Mapped[str] = mapped_column(String(20), nullable=False, default="ACTIVE")
