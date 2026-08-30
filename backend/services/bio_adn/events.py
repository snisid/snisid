from __future__ import annotations

from typing import Any, ClassVar

from shared.cqrs import DomainEvent


class STRProfileSubmitted(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.profile.submitted"

    def __init__(
        self,
        specimen_number: str,
        index_type: str,
        loci_data: dict,
        quality_score: float,
        case_number: str,
        correlation_id: str,
    ) -> None:
        super().__init__(
            event_type="STR_PROFILE_SUBMITTED",
            aggregate_id=specimen_number,
            aggregate_type="str_profile",
            data={
                "aggregate_id": specimen_number,
                "aggregate_type": "str_profile",
                "specimen_number": specimen_number,
                "index_type": index_type,
                "loci_data": loci_data,
                "quality_score": quality_score,
                "case_number": case_number,
                "correlation_id": correlation_id,
            },
        )


class HitFound(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.hits"

    def __init__(
        self,
        sample_id: str,
        hit_type: str,
        matched_profile: str,
        confidence: float,
        case_number: str,
    ) -> None:
        super().__init__(
            event_type="HIT_FOUND",
            aggregate_id=sample_id,
            aggregate_type="dna_hit",
            data={
                "aggregate_id": sample_id,
                "aggregate_type": "dna_hit",
                "sample_id": sample_id,
                "hit_type": hit_type,
                "matched_profile": matched_profile,
                "confidence": confidence,
                "case_number": case_number,
            },
        )


class SyncRequested(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.sync.requested"

    def __init__(self, source_index: str, target_index: str) -> None:
        super().__init__(
            event_type="SYNC_REQUESTED",
            aggregate_id=source_index,
            aggregate_type="sync",
            data={
                "aggregate_id": source_index,
                "aggregate_type": "sync",
                "source_index": source_index,
                "target_index": target_index,
            },
        )


class AuditLogCreated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.audit"

    def __init__(self, **kwargs: Any) -> None:
        agg_id = kwargs.get("audit_id", "")
        super().__init__(
            event_type="AUDIT_LOG_CREATED",
            aggregate_id=agg_id,
            aggregate_type="audit_log",
            data={
                "aggregate_id": agg_id,
                "aggregate_type": "audit_log",
                **kwargs,
            },
        )


class WantedPersonCreated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.wanted.events"

    def __init__(
        self,
        record_id: str,
        record_number: str,
        warrant_type: str,
        warrant_number: str,
        charges: list[str],
        danger_level: str,
        mco_contact: str,
        entering_officer: str,
        entering_agency: str,
    ) -> None:
        super().__init__(
            event_type="WANTED_PERSON_CREATED",
            aggregate_id=record_id,
            aggregate_type="wanted_person",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "wanted_person",
                "record_id": record_id,
                "record_number": record_number,
                "warrant_type": warrant_type,
                "warrant_number": warrant_number,
                "charges": charges,
                "danger_level": danger_level,
                "mco_contact": mco_contact,
                "entering_officer": entering_officer,
                "entering_agency": entering_agency,
            },
        )


class MissingPersonReported(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.missing.events"

    def __init__(
        self,
        record_id: str,
        record_number: str,
        category: str,
        entering_agency: str,
        citizen_portal_submission: bool,
        dna_sample_available: bool,
    ) -> None:
        super().__init__(
            event_type="MISSING_PERSON_REPORTED",
            aggregate_id=record_id,
            aggregate_type="missing_person",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "missing_person",
                "record_id": record_id,
                "record_number": record_number,
                "category": category,
                "entering_agency": entering_agency,
                "citizen_portal_submission": citizen_portal_submission,
                "dna_sample_available": dna_sample_available,
            },
        )


class VehicleStolen(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.vehicle.stolen"

    def __init__(self, record_id: str, record_number: str, plate_number: str, vin: str) -> None:
        super().__init__(
            event_type="VEHICLE_STOLEN",
            aggregate_id=record_id,
            aggregate_type="stolen_vehicle",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "stolen_vehicle",
                "record_id": record_id,
                "record_number": record_number,
                "plate_number": plate_number,
                "vin": vin,
            },
        )


class FirearmStolen(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.arm.stolen"

    def __init__(self, record_id: str, record_number: str, serial_number: str) -> None:
        super().__init__(
            event_type="FIREARM_STOLEN",
            aggregate_id=record_id,
            aggregate_type="stolen_firearm",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "stolen_firearm",
                "record_id": record_id,
                "record_number": record_number,
                "serial_number": serial_number,
            },
        )


class DocumentStolen(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.document.stolen"

    def __init__(self, record_id: str, record_number: str, document_type: str, document_number: str) -> None:
        super().__init__(
            event_type="DOCUMENT_STOLEN",
            aggregate_id=record_id,
            aggregate_type="stolen_document",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "stolen_document",
                "record_id": record_id,
                "record_number": record_number,
                "document_type": document_type,
                "document_number": document_number,
            },
        )


class VesselStolen(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.vessel.stolen"

    def __init__(self, record_id: str, record_number: str, vessel_name: str, registration_number: str) -> None:
        super().__init__(
            event_type="VESSEL_STOLEN",
            aggregate_id=record_id,
            aggregate_type="stolen_vessel",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "stolen_vessel",
                "record_id": record_id,
                "record_number": record_number,
                "vessel_name": vessel_name,
                "registration_number": registration_number,
            },
        )


class VehicleRecovered(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.vehicle.recovered"

    def __init__(self, record_id: str, record_number: str, recovered_location: str) -> None:
        super().__init__(
            event_type="VEHICLE_RECOVERED",
            aggregate_id=record_id,
            aggregate_type="vehicle_recovery",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "vehicle_recovery",
                "record_id": record_id,
                "record_number": record_number,
                "recovered_location": recovered_location,
            },
        )


class ArmCrimeSceneHit(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.arm.hit"

    def __init__(self, record_id: str, record_number: str, case_number: str) -> None:
        super().__init__(
            event_type="ARM_CRIME_SCENE_HIT",
            aggregate_id=record_id,
            aggregate_type="arm_crime_scene_hit",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "arm_crime_scene_hit",
                "record_id": record_id,
                "record_number": record_number,
                "case_number": case_number,
            },
        )


class StolenArticleCreated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.article.stolen"

    def __init__(self, record_id: str, record_number: str, category: str) -> None:
        super().__init__(
            event_type="ARTICLE_STOLEN",
            aggregate_id=record_id,
            aggregate_type="stolen_article",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "stolen_article",
                "record_id": record_id,
                "record_number": record_number,
                "category": category,
            },
        )


class StolenSecurityCreated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.security.stolen"

    def __init__(self, record_id: str, record_number: str, security_type: str) -> None:
        super().__init__(
            event_type="SECURITY_STOLEN",
            aggregate_id=record_id,
            aggregate_type="stolen_security",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "stolen_security",
                "record_id": record_id,
                "record_number": record_number,
                "security_type": security_type,
            },
        )


class ONIDocumentRevoked(DomainEvent):
    topic: ClassVar[str] = "snisid.oni.document.revoked"

    def __init__(self, document_number: str, revocation_reason: str) -> None:
        super().__init__(
            event_type="ONI_DOCUMENT_REVOKED",
            aggregate_id=document_number,
            aggregate_type="oni_document",
            data={
                "aggregate_id": document_number,
                "aggregate_type": "oni_document",
                "document_number": document_number,
                "revocation_reason": revocation_reason,
            },
        )


class ForeignFugitiveCreated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.fugitive.events"

    def __init__(
        self,
        record_id: str,
        record_number: str,
        interpol_notice_number: str,
        notice_type: str,
        issuing_country: str,
        entering_agency: str,
    ) -> None:
        super().__init__(
            event_type="FOREIGN_FUGITIVE_CREATED",
            aggregate_id=record_id,
            aggregate_type="foreign_fugitive",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "foreign_fugitive",
                "record_id": record_id,
                "record_number": record_number,
                "interpol_notice_number": interpol_notice_number,
                "notice_type": notice_type,
                "issuing_country": issuing_country,
                "entering_agency": entering_agency,
            },
        )


class UnidentifiedPersonCreated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.unidentified.events"

    def __init__(
        self,
        record_id: str,
        record_number: str,
        discovery_date: str,
        discovery_location: str,
        entering_agency: str,
        dna_sample_ref: str,
    ) -> None:
        super().__init__(
            event_type="UNIDENTIFIED_PERSON_CREATED",
            aggregate_id=record_id,
            aggregate_type="unidentified_person",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "unidentified_person",
                "record_id": record_id,
                "record_number": record_number,
                "discovery_date": discovery_date,
                "discovery_location": discovery_location,
                "entering_agency": entering_agency,
                "dna_sample_ref": dna_sample_ref,
            },
        )


class UnidentifiedIdentified(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.unidentified.events"

    def __init__(self, record_id: str, record_number: str, matched_to_niu: str) -> None:
        super().__init__(
            event_type="UNIDENTIFIED_IDENTIFIED",
            aggregate_id=record_id,
            aggregate_type="unidentified_person",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "unidentified_person",
                "record_id": record_id,
                "record_number": record_number,
                "matched_to_niu": matched_to_niu,
            },
        )


class TerrorismWatchCreated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.terrorism.events"

    def __init__(
        self,
        record_id: str,
        record_number: str,
        threat_type: str,
        risk_level: str,
        entering_agency: str,
    ) -> None:
        super().__init__(
            event_type="TERRORISM_WATCH_CREATED",
            aggregate_id=record_id,
            aggregate_type="terrorism_watch",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "terrorism_watch",
                "record_id": record_id,
                "record_number": record_number,
                "threat_type": threat_type,
                "risk_level": risk_level,
                "entering_agency": entering_agency,
            },
        )


class ProtectionOrderCreated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.protection.events"

    def __init__(
        self,
        record_id: str,
        record_number: str,
        order_type: str,
        beneficiary_name: str,
        restrained_person: str,
        issuing_court: str,
    ) -> None:
        super().__init__(
            event_type="PROTECTION_ORDER_CREATED",
            aggregate_id=record_id,
            aggregate_type="protection_order",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "protection_order",
                "record_id": record_id,
                "record_number": record_number,
                "order_type": order_type,
                "beneficiary_name": beneficiary_name,
                "restrained_person": restrained_person,
                "issuing_court": issuing_court,
            },
        )


class SupervisedReleaseCreated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.supervised.events"

    def __init__(
        self,
        record_id: str,
        record_number: str,
        niu: str,
        supervision_type: str,
        supervising_officer: str,
        supervising_agency: str,
    ) -> None:
        super().__init__(
            event_type="SUPERVISED_RELEASE_CREATED",
            aggregate_id=record_id,
            aggregate_type="supervised_release",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "supervised_release",
                "record_id": record_id,
                "record_number": record_number,
                "niu": niu,
                "supervision_type": supervision_type,
                "supervising_officer": supervising_officer,
                "supervising_agency": supervising_agency,
            },
        )


class DuplicateSpecimenDetected(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.lab.duplicate"

    def __init__(self, specimen_number: str, existing_sample_id: str, new_submission_id: str) -> None:
        super().__init__(
            event_type="BIO_DUPLICATE_SPECIMEN",
            aggregate_id=specimen_number,
            aggregate_type="specimen",
            data={
                "aggregate_id": specimen_number,
                "aggregate_type": "specimen",
                "specimen_number": specimen_number,
                "existing_sample_id": existing_sample_id,
                "new_submission_id": new_submission_id,
            },
        )


class ProfileExpunged(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.profile.expunged"

    def __init__(self, sample_id: str, court_order_ref: str, reason: str, officer_niu: str) -> None:
        super().__init__(
            event_type="PROFILE_EXPUNGED",
            aggregate_id=sample_id,
            aggregate_type="dna_profile",
            data={
                "aggregate_id": sample_id,
                "aggregate_type": "dna_profile",
                "sample_id": sample_id,
                "court_order_ref": court_order_ref,
                "reason": reason,
                "officer_niu": officer_niu,
            },
        )


class EquipmentRegistered(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.lab.equipment"

    def __init__(self, equipment_id: str, lab_code: str, equipment_name: str) -> None:
        super().__init__(
            event_type="EQUIPMENT_REGISTERED",
            aggregate_id=equipment_id,
            aggregate_type="lab_equipment",
            data={
                "aggregate_id": equipment_id,
                "aggregate_type": "lab_equipment",
                "equipment_id": equipment_id,
                "lab_code": lab_code,
                "equipment_name": equipment_name,
            },
        )


class TrainingRecorded(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.lab.training"

    def __init__(self, training_id: str, staff_niu: str, training_name: str) -> None:
        super().__init__(
            event_type="TRAINING_RECORDED",
            aggregate_id=training_id,
            aggregate_type="staff_training",
            data={
                "aggregate_id": training_id,
                "aggregate_type": "staff_training",
                "training_id": training_id,
                "staff_niu": staff_niu,
                "training_name": training_name,
            },
        )


class LDISUploadCompleted(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.lab.upload"

    def __init__(self, lab_code: str, uploaded_count: int, operator_niu: str) -> None:
        super().__init__(
            event_type="LDIS_UPLOAD_COMPLETED",
            aggregate_id=lab_code,
            aggregate_type="ldis_upload",
            data={
                "aggregate_id": lab_code,
                "aggregate_type": "ldis_upload",
                "lab_code": lab_code,
                "uploaded_count": uploaded_count,
                "operator_niu": operator_niu,
            },
        )


class NDISProfileUploaded(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.ndis.profile.uploaded"

    def __init__(self, sample_id: str, index_type: str, source_sdis: str) -> None:
        super().__init__(
            event_type="NDIS_PROFILE_UPLOADED",
            aggregate_id=sample_id,
            aggregate_type="ndis_profile",
            data={
                "aggregate_id": sample_id,
                "aggregate_type": "ndis_profile",
                "sample_id": sample_id,
                "index_type": index_type,
                "source_sdis": source_sdis,
            },
        )


class CrossDeptHitDetected(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.ndis.crossdept.hit"

    def __init__(
        self,
        hit_id: str,
        query_sample_id: str,
        match_sample_id: str,
        match_type: str,
        confidence: float,
        query_sdis: str,
        match_sdis: str,
    ) -> None:
        super().__init__(
            event_type="CROSS_DEPT_HIT_DETECTED",
            aggregate_id=hit_id,
            aggregate_type="cross_dept_hit",
            data={
                "aggregate_id": hit_id,
                "aggregate_type": "cross_dept_hit",
                "hit_id": hit_id,
                "query_sample_id": query_sample_id,
                "match_sample_id": match_sample_id,
                "match_type": match_type,
                "confidence": confidence,
                "query_sdis": query_sdis,
                "match_sdis": match_sdis,
            },
        )


class InterpolSubmissionRequested(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.ndis.interpol"

    def __init__(self, submission_id: str, sample_ids: list[str], reason: str) -> None:
        super().__init__(
            event_type="INTERPOL_SUBMISSION_REQUESTED",
            aggregate_id=submission_id,
            aggregate_type="interpol_submission",
            data={
                "aggregate_id": submission_id,
                "aggregate_type": "interpol_submission",
                "submission_id": submission_id,
                "sample_ids": sample_ids,
                "reason": reason,
            },
        )


class NDISReportGenerated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.ndis.reports"

    def __init__(self, report_id: str, report_type: str, generated_by: str) -> None:
        super().__init__(
            event_type="NDIS_REPORT_GENERATED",
            aggregate_id=report_id,
            aggregate_type="ndis_report",
            data={
                "aggregate_id": report_id,
                "aggregate_type": "ndis_report",
                "report_id": report_id,
                "report_type": report_type,
                "generated_by": generated_by,
            },
        )


class ViolenceRecordCreated(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.violence.events"

    def __init__(
        self,
        record_id: str,
        record_number: str,
        incident_type: str,
        niu: str | None,
        arresting_agency: str,
        risk_level: str,
    ) -> None:
        super().__init__(
            event_type="VIOLENCE_RECORD_CREATED",
            aggregate_id=record_id,
            aggregate_type="violence_record",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "violence_record",
                "record_id": record_id,
                "record_number": record_number,
                "incident_type": incident_type,
                "niu": niu,
                "arresting_agency": arresting_agency,
                "risk_level": risk_level,
            },
        )


class IdentityTheftRecorded(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.identitytheft.events"

    def __init__(
        self,
        record_id: str,
        record_number: str,
        victim_niu: str,
        fraud_type: str,
        reporting_agency: str,
    ) -> None:
        super().__init__(
            event_type="IDENTITY_THEFT_RECORDED",
            aggregate_id=record_id,
            aggregate_type="identity_theft",
            data={
                "aggregate_id": record_id,
                "aggregate_type": "identity_theft",
                "record_id": record_id,
                "record_number": record_number,
                "victim_niu": victim_niu,
                "fraud_type": fraud_type,
                "reporting_agency": reporting_agency,
            },
        )


class BioIdentityLinked(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.identity.linked"

    def __init__(
        self,
        sample_id: str,
        niu: str,
        linked_by: str,
        court_order: str | None = None,
    ) -> None:
        super().__init__(
            event_type="BIO_IDENTITY_LINKED",
            aggregate_id=sample_id,
            aggregate_type="bio_identity_link",
            data={
                "aggregate_id": sample_id,
                "aggregate_type": "bio_identity_link",
                "sample_id": sample_id,
                "niu": niu,
                "linked_by": linked_by,
                "court_order": court_order,
            },
        )


class BioSecurityAlert(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.security.alert"

    def __init__(self, alert_type: str, sdis_code: str, details: str, sample_id: str) -> None:
        super().__init__(
            event_type="BIO_SECURITY_ALERT",
            aggregate_id=sample_id,
            aggregate_type="bio_security_alert",
            data={
                "aggregate_id": sample_id,
                "aggregate_type": "bio_security_alert",
                "alert_type": alert_type,
                "sdis_code": sdis_code,
                "details": details,
                "sample_id": sample_id,
            },
        )


class SDISHeartbeat(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.sdis.heartbeat"

    def __init__(self, sdis_code: str, node_healthy: bool, upstream_ok: bool) -> None:
        super().__init__(
            event_type="SDIS_HEARTBEAT",
            aggregate_id=sdis_code,
            aggregate_type="sdis_heartbeat",
            data={
                "aggregate_id": sdis_code,
                "aggregate_type": "sdis_heartbeat",
                "sdis_code": sdis_code,
                "node_healthy": node_healthy,
                "upstream_ok": upstream_ok,
            },
        )


class SDISIntraDeptMatch(DomainEvent):
    topic: ClassVar[str] = "snisid.bio.sdis.intradept.match"

    def __init__(self, sdis_code: str, match_id: str, match_type: str, confidence: float) -> None:
        super().__init__(
            event_type="SDIS_INTRADEPT_MATCH",
            aggregate_id=match_id,
            aggregate_type="sdis_intradept_match",
            data={
                "aggregate_id": match_id,
                "aggregate_type": "sdis_intradept_match",
                "sdis_code": sdis_code,
                "match_id": match_id,
                "match_type": match_type,
                "confidence": confidence,
            },
        )


TOPIC_MAP: dict[str, type[DomainEvent]] = {
    "snisid.bio.profile.submitted": STRProfileSubmitted,
    "snisid.bio.hits": HitFound,
    "snisid.bio.sync.requested": SyncRequested,
    "snisid.bio.audit": AuditLogCreated,
    "snisid.bio.wanted.events": WantedPersonCreated,
    "snisid.bio.missing.events": MissingPersonReported,
    "snisid.bio.vehicle.stolen": VehicleStolen,
    "snisid.bio.arm.stolen": FirearmStolen,
    "snisid.bio.document.stolen": DocumentStolen,
    "snisid.bio.vessel.stolen": VesselStolen,
    "snisid.bio.vehicle.recovered": VehicleRecovered,
    "snisid.bio.arm.hit": ArmCrimeSceneHit,
    "snisid.bio.article.stolen": StolenArticleCreated,
    "snisid.bio.security.stolen": StolenSecurityCreated,
    "snisid.oni.document.revoked": ONIDocumentRevoked,
    "snisid.bio.fugitive.events": ForeignFugitiveCreated,
    "snisid.bio.unidentified.events": UnidentifiedPersonCreated,
    "snisid.bio.terrorism.events": TerrorismWatchCreated,
    "snisid.bio.protection.events": ProtectionOrderCreated,
    "snisid.bio.supervised.events": SupervisedReleaseCreated,
    "snisid.bio.lab.duplicate": DuplicateSpecimenDetected,
    "snisid.bio.profile.expunged": ProfileExpunged,
    "snisid.bio.lab.equipment": EquipmentRegistered,
    "snisid.bio.lab.training": TrainingRecorded,
    "snisid.bio.lab.upload": LDISUploadCompleted,
    "snisid.bio.ndis.profile.uploaded": NDISProfileUploaded,
    "snisid.bio.ndis.crossdept.hit": CrossDeptHitDetected,
    "snisid.bio.ndis.interpol": InterpolSubmissionRequested,
    "snisid.bio.ndis.reports": NDISReportGenerated,
    "snisid.bio.security.alert": BioSecurityAlert,
    "snisid.bio.sdis.heartbeat": SDISHeartbeat,
    "snisid.bio.sdis.intradept.match": SDISIntraDeptMatch,
    "snisid.bio.violence.events": ViolenceRecordCreated,
    "snisid.bio.identitytheft.events": IdentityTheftRecorded,
    "snisid.bio.identity.linked": BioIdentityLinked,
}
