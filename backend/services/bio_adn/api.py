from __future__ import annotations

import uuid
from datetime import datetime, timezone
from typing import Optional

from fastapi import APIRouter, Body, Depends, HTTPException, Query, status

from shared.events import KafkaProducer
from services.bio_adn.events import (
    STRProfileSubmitted,
    HitFound,
    WantedPersonCreated,
    MissingPersonReported,
    VehicleStolen,
    VehicleRecovered,
    DocumentStolen,
    VesselStolen,
    FirearmStolen,
    ArmCrimeSceneHit,
    StolenArticleCreated,
    StolenSecurityCreated,
    ONIDocumentRevoked,
    ForeignFugitiveCreated,
    UnidentifiedPersonCreated,
    TerrorismWatchCreated,
    ProtectionOrderCreated,
    SupervisedReleaseCreated,
    ViolenceRecordCreated,
    IdentityTheftRecorded,
    DuplicateSpecimenDetected,
    ProfileExpunged,
    EquipmentRegistered,
    TrainingRecorded,
    LDISUploadCompleted,
    CrossDeptHitDetected,
    InterpolSubmissionRequested,
    NDISReportGenerated,
    BioSecurityAlert,
    SDISHeartbeat,
    SDISIntraDeptMatch,
    BioIdentityLinked,
)
from services.bio_adn.quality import validate_profile
from services.bio_adn.models import (
    STRProfileSubmit,
    STRProfileResponse,
    DNASearchRequest,
    DNASearchResponse,
    CreateWantedPersonRequest,
    WantedPersonResponse,
    CreateMissingPersonRequest,
    MissingPersonResponse,
    StolenVehicleRequest,
    VehicleRecoverRequest,
    StolenFirearmRequest,
    StolenDocumentRequest,
    StolenVesselRequest,
    StolenArticleRequest,
    StolenArticleResponse,
    StolenSecurityRequest,
    StolenSecurityResponse,
    LAPIPlateResponse,
    CreateForeignFugitiveRequest,
    ForeignFugitiveResponse,
    CreateUnidentifiedPersonRequest,
    UnidentifiedPersonResponse,
    CreateTerrorismWatchRequest,
    TerrorismWatchResponse,
    CreateProtectionOrderRequest,
    ProtectionOrderResponse,
    CreateSupervisedReleaseRequest,
    SupervisedReleaseResponse,
    CreateViolenceRecordRequest,
    ViolenceRecordResponse,
    ViolenceIncidentType,
    CreateIdentityTheftRequest,
    IdentityTheftResponse,
    FraudType,
    SexOffenderUpdateRequest,
    GangMemberReviewRequest,
    LabEquipmentCreate,
    LabEquipmentResponse,
    LabEquipmentUpdate,
    StaffTrainingCreate,
    StaffTrainingResponse,
    UploadRequest,
    UploadResponse,
    NDISStats,
    NDISReportResponse,
    NDISCrossDeptHit,
    InterpolSubmissionCreate,
    InterpolSubmissionResponse,
    SdisNodeResponse,
    SdisStats,
    SdisMatchResponse,
    AlertLevel,
    WarrantType,
    WarrantStatus,
    MissingStatus,
    StolenStatus,
    InterpolNoticeType,
)

router = APIRouter(prefix="/v1/bio-adn", tags=["BIO-ADN"])

MAX_LAPI_RESPONSE_MS = 200

_producer: KafkaProducer | None = None


def get_producer() -> KafkaProducer | None:
    return _producer


async def _publish(event, key: str = ""):
    producer = get_producer()
    if producer is None:
        return
    topic = getattr(type(event), "topic", "")
    if not topic:
        return
    try:
        await producer.publish(topic=topic, key=key or str(uuid.uuid4()), event=event)
    except Exception:
        pass


def _error(code: str, message: str, status_code: int = 400):
    return HTTPException(status_code=status_code, detail={"code": code, "message": message})


# ── DNA Profiles ───────────────────────────────────────────────────────────


_duplicate_specimens: set[str] = set()


@router.post("/dna/profiles", response_model=STRProfileResponse, status_code=201)
async def submit_dna_profile(profile: STRProfileSubmit):
    loci_count = len(profile.loci_data)
    errors = validate_profile(profile.index_type.value, profile.quality_score, loci_count)
    if errors:
        raise _error("BIO-001", "; ".join(errors), 422)

    if profile.specimen_number in _duplicate_specimens:
        event = DuplicateSpecimenDetected(
            specimen_number=profile.specimen_number,
            existing_sample_id="",
            new_submission_id=str(uuid.uuid4()),
        )
        await _publish(event, key=profile.specimen_number)
        raise _error("BIO-409", f"Specimen_number {profile.specimen_number} déjà soumis", 409)

    _duplicate_specimens.add(profile.specimen_number)

    sample_id = str(uuid.uuid4())
    event = STRProfileSubmitted(
        specimen_number=profile.specimen_number,
        index_type=profile.index_type.value,
        loci_data=profile.loci_data,
        quality_score=profile.quality_score,
        case_number=profile.case_number or "",
        correlation_id=profile.correlation_id,
    )
    await _publish(event, key=sample_id)

    return STRProfileResponse(
        sample_id=uuid.UUID(sample_id),
        accepted=True,
        message="Profil soumis avec succès",
    )


@router.get("/dna/profiles")
async def list_dna_profiles(lab_id: Optional[str] = Query(None)):
    return {"profiles": [], "total": 0, "lab_id": lab_id or "all"}


@router.post("/dna/search", response_model=DNASearchResponse)
async def search_dna_profiles(request: DNASearchRequest):
    return DNASearchResponse(
        hits=[],
        total_hits=0,
        search_duration_ms=0,
        case_number=request.case_number,
    )


@router.get("/dna/hits/{hit_id}")
async def get_dna_hit(hit_id: str):
    if not hit_id:
        raise _error("BIO-001", "Identifiant de hit invalide", 404)

    return {"hit_id": hit_id, "status": "pending_review"}


@router.post("/dna/profiles/{sample_id}/expunge", status_code=202)
async def expunge_dna_profile(
    sample_id: str,
    court_order_ref: str = Query(...),
    reason: str = "court_order",
    officer_niu: str = "system",
):
    event = ProfileExpunged(
        sample_id=sample_id,
        court_order_ref=court_order_ref,
        reason=reason,
        officer_niu=officer_niu,
    )
    await _publish(event, key=sample_id)
    return {
        "success": True,
        "expunge_id": str(uuid.uuid4()),
        "timestamp": datetime.now(timezone.utc).isoformat(),
    }


# ── Identity Linkage (restreint DCPJ-DIR) ────────────────────────────────


@router.post("/dna/identity-links", status_code=201)
async def create_identity_link(
    sample_id: str = Query(...),
    niu: str = Query(..., min_length=5),
    linked_by: str = Query(...),
    court_order_ref: str = Query(default=""),
):
    event = BioIdentityLinked(
        sample_id=sample_id,
        niu=niu,
        linked_by=linked_by,
        court_order=court_order_ref or None,
    )
    await _publish(event, key=sample_id)

    return {
        "link_id": str(uuid.uuid4()),
        "sample_id": sample_id,
        "niu": niu,
        "linked_by": linked_by,
        "linked_at": datetime.now(timezone.utc).isoformat(),
        "status": "active",
    }


@router.get("/dna/identity-links/by-niu/{niu}")
async def get_identity_links_by_niu(niu: str):
    return {"niu": niu, "links": []}


@router.get("/dna/identity-links/{sample_id}")
async def get_identity_link(sample_id: str):
    return {"sample_id": sample_id, "linked": True}


# ── PER-REC: Wanted Persons ────────────────────────────────────────────────


@router.post("/persons/wanted", response_model=WantedPersonResponse, status_code=201)
async def create_wanted_person(request: CreateWantedPersonRequest):
    if request.warrant_type in (WarrantType.ARREST, WarrantType.EXTRADITION):
        if not request.warrant_number:
            raise _error("PER-001", "warrant_number obligatoire pour MAN-ARR et MAN-EXT", 422)

    record_id = str(uuid.uuid4())
    record_number = f"PRE-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = WantedPersonCreated(
        record_id=record_id,
        record_number=record_number,
        warrant_type=request.warrant_type.value,
        warrant_number=request.warrant_number or "",
        charges=request.charges,
        danger_level=request.danger_level.value,
        mco_contact=request.mco_contact,
        entering_officer=request.entering_officer or "",
        entering_agency=request.mco_contact,
    )
    await _publish(event, key=record_id)

    return WantedPersonResponse(
        record_id=uuid.UUID(record_id),
        record_number=record_number,
        status=WarrantStatus.ACTIVE.value,
        mco_contact=request.mco_contact,
    )


@router.get("/persons/wanted/query")
async def query_wanted_persons(
    last_name: Optional[str] = Query(None),
    first_name: Optional[str] = Query(None),
    niu: Optional[str] = Query(None),
    plate_number: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


@router.get("/persons/wanted/{record_id}")
async def get_wanted_person(record_id: str):
    return {"record_id": record_id, "status": "ACTIVE"}


@router.patch("/persons/wanted/{record_id}/status")
async def update_wanted_status(record_id: str, status: WarrantStatus = Query(...)):
    return {"record_id": record_id, "status": status.value, "updated_at": datetime.now(timezone.utc).isoformat()}


# ── PER-FUG: Foreign Fugitives ──────────────────────────────────────────────


@router.post("/persons/foreign-fugitives", response_model=ForeignFugitiveResponse, status_code=201)
async def create_foreign_fugitive(request: CreateForeignFugitiveRequest):
    record_id = str(uuid.uuid4())
    record_number = f"FUG-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = ForeignFugitiveCreated(
        record_id=record_id,
        record_number=record_number,
        interpol_notice_number=request.interpol_notice_number,
        notice_type=request.notice_type.value,
        issuing_country=request.issuing_country,
        entering_agency=request.entering_agency,
    )
    await _publish(event, key=record_id)

    return ForeignFugitiveResponse(
        record_id=uuid.UUID(record_id),
        record_number=record_number,
        interpol_notice_number=request.interpol_notice_number,
        status="ACTIVE",
    )


@router.get("/persons/foreign-fugitives/query")
async def query_foreign_fugitives(
    last_name: Optional[str] = Query(None),
    nationality: Optional[str] = Query(None),
    notice_type: Optional[InterpolNoticeType] = Query(None),
):
    return {"total": 0, "results": []}


@router.get("/persons/foreign-fugitives/{record_id}")
async def get_foreign_fugitive(record_id: str):
    return {"record_id": record_id, "status": "ACTIVE"}


# ── PER-DIS: Missing Persons ────────────────────────────────────────────────


@router.post("/persons/missing", response_model=MissingPersonResponse, status_code=201)
async def create_missing_person(request: CreateMissingPersonRequest):
    record_id = str(uuid.uuid4())
    record_number = f"DIS-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    bpm_notified = bool(request.category.value == "CHILD")

    event = MissingPersonReported(
        record_id=record_id,
        record_number=record_number,
        category=request.category.value,
        entering_agency=request.entering_agency,
        citizen_portal_submission=request.citizen_portal_submission,
        dna_sample_available=request.dna_sample_available,
    )
    await _publish(event, key=record_id)

    return MissingPersonResponse(
        record_id=uuid.UUID(record_id),
        record_number=record_number,
        status=MissingStatus.ACTIVE.value,
        bpm_notified=bpm_notified,
    )


@router.get("/persons/missing/query")
async def query_missing_persons(
    last_name: Optional[str] = Query(None),
    first_name: Optional[str] = Query(None),
    niu: Optional[str] = Query(None),
    category: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


@router.get("/persons/missing/{record_id}")
async def get_missing_person(record_id: str):
    return {"record_id": record_id, "status": "ACTIVE"}


# ── PER-NID: Unidentified Persons ────────────────────────────────────────────


@router.post("/persons/unidentified", response_model=UnidentifiedPersonResponse, status_code=201)
async def create_unidentified_person(request: CreateUnidentifiedPersonRequest):
    record_id = str(uuid.uuid4())
    record_number = f"NID-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = UnidentifiedPersonCreated(
        record_id=record_id,
        record_number=record_number,
        discovery_date=request.discovery_date,
        discovery_location=request.discovery_location,
        entering_agency=request.entering_agency,
        dna_sample_ref=request.dna_sample_ref or "",
    )
    await _publish(event, key=record_id)

    return UnidentifiedPersonResponse(
        record_id=uuid.UUID(record_id),
        record_number=record_number,
        status="ACTIVE",
    )


@router.get("/persons/unidentified/query")
async def query_unidentified_persons(
    discovery_department: Optional[str] = Query(None),
    gender: Optional[str] = Query(None),
    estimated_age_min: Optional[int] = Query(None),
    estimated_age_max: Optional[int] = Query(None),
):
    return {"total": 0, "results": []}


@router.get("/persons/unidentified/{record_id}")
async def get_unidentified_person(record_id: str):
    return {"record_id": record_id, "status": "ACTIVE"}


# ── PER-SEX: Sex Offender Registry ──────────────────────────────────────────


@router.post("/persons/sex-offenders", status_code=201)
async def register_sex_offender(
    niu: str,
    conviction_date: str,
    conviction_court: str,
    offenses: list[str],
    risk_level: str = "MEDIUM",
):
    if risk_level not in ("LOW", "MEDIUM", "HIGH"):
        raise _error("PER-002", "risk_level doit être LOW, MEDIUM ou HIGH", 422)
    return {"record_id": str(uuid.uuid4()), "status": "ACTIVE", "risk_level": risk_level}


@router.patch("/persons/sex-offenders/{record_id}/risk", status_code=200)
async def update_sex_offender_risk(record_id: str, request: SexOffenderUpdateRequest):
    return {
        "record_id": record_id,
        "risk_level": request.risk_level.value,
        "updated_at": datetime.now(timezone.utc).isoformat(),
    }


@router.get("/persons/sex-offenders/query")
async def query_sex_offenders(
    risk_level: Optional[str] = Query(None),
    location: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


# ── PER-GNG: Gang Members ────────────────────────────────────────────────────


@router.post("/persons/gang-members", status_code=201)
async def register_gang_member(
    niu: Optional[str] = None,
    last_name: Optional[str] = None,
    first_name: Optional[str] = None,
    gang_name: str = ...,
    membership_type: str = "MEMBER",
):
    return {"record_id": str(uuid.uuid4())}


@router.post("/persons/gang-members/{record_id}/review", status_code=200)
async def review_gang_member(record_id: str, request: GangMemberReviewRequest):
    return {
        "record_id": record_id,
        "reviewed_at": datetime.now(timezone.utc).isoformat(),
    }


@router.get("/persons/gang-members/query")
async def query_gang_members(
    gang_name: Optional[str] = Query(None),
    membership_type: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


# ── PER-TER: Terrorism Watch ────────────────────────────────────────────────


@router.post("/persons/terrorism", response_model=TerrorismWatchResponse, status_code=201)
async def create_terrorism_watch(request: CreateTerrorismWatchRequest):
    record_id = str(uuid.uuid4())
    record_number = f"TER-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = TerrorismWatchCreated(
        record_id=record_id,
        record_number=record_number,
        threat_type=request.threat_type,
        risk_level=request.risk_level,
        entering_agency=request.entering_agency,
    )
    await _publish(event, key=record_id)

    return TerrorismWatchResponse(
        record_id=uuid.UUID(record_id),
        record_number=record_number,
        status="ACTIVE",
    )


@router.get("/persons/terrorism/query")
async def query_terrorism_watches(
    risk_level: Optional[str] = Query(None),
    threat_type: Optional[str] = Query(None),
    nationality: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


@router.get("/persons/terrorism/{record_id}")
async def get_terrorism_watch(record_id: str):
    return {"record_id": record_id, "status": "ACTIVE"}


# ── PER-OPR: Protection Orders ──────────────────────────────────────────────


@router.post("/persons/protection-orders", response_model=ProtectionOrderResponse, status_code=201)
async def create_protection_order(request: CreateProtectionOrderRequest):
    record_id = str(uuid.uuid4())
    record_number = f"OPR-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = ProtectionOrderCreated(
        record_id=record_id,
        record_number=record_number,
        order_type=request.order_type.value,
        beneficiary_name=request.beneficiary_name,
        restrained_person=request.restrained_person,
        issuing_court=request.issuing_court,
    )
    await _publish(event, key=record_id)

    return ProtectionOrderResponse(
        record_id=uuid.UUID(record_id),
        record_number=record_number,
        order_type=request.order_type.value,
        status="ACTIVE",
    )


@router.get("/persons/protection-orders/query")
async def query_protection_orders(
    beneficiary_name: Optional[str] = Query(None),
    restrained_person: Optional[str] = Query(None),
    order_type: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


@router.get("/persons/protection-orders/urgent/{beneficiary_niu}")
async def get_protection_order_urgent(beneficiary_niu: str):
    return {"beneficiary_niu": beneficiary_niu, "active_orders": 0, "response_ms": 0}


# ── PER-LIB: Supervised Release ─────────────────────────────────────────────


@router.post("/persons/supervised-release", response_model=SupervisedReleaseResponse, status_code=201)
async def create_supervised_release(request: CreateSupervisedReleaseRequest):
    record_id = str(uuid.uuid4())
    record_number = f"LIB-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = SupervisedReleaseCreated(
        record_id=record_id,
        record_number=record_number,
        niu=request.niu,
        supervision_type=request.supervision_type.value,
        supervising_officer=request.supervising_officer,
        supervising_agency=request.supervising_agency,
    )
    await _publish(event, key=record_id)

    return SupervisedReleaseResponse(
        record_id=uuid.UUID(record_id),
        record_number=record_number,
        status="ACTIVE",
    )


@router.get("/persons/supervised-release/query")
async def query_supervised_releases(
    niu: Optional[str] = Query(None),
    supervision_type: Optional[str] = Query(None),
    status: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


@router.get("/persons/supervised-release/{record_id}")
async def get_supervised_release(record_id: str):
    return {"record_id": record_id, "status": "ACTIVE"}


# ── PER-VIO: Known Violence ─────────────────────────────────────────────────


@router.post("/persons/violence", response_model=ViolenceRecordResponse, status_code=201)
async def create_violence_record(request: CreateViolenceRecordRequest):
    record_id = str(uuid.uuid4())
    record_number = f"VIO-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = ViolenceRecordCreated(
        record_id=record_id,
        record_number=record_number,
        incident_type=request.incident_type.value,
        niu=request.niu,
        arresting_agency=request.arresting_agency,
        risk_level=request.risk_level,
    )
    await _publish(event, key=record_id)

    return ViolenceRecordResponse(record_id=uuid.UUID(record_id), record_number=record_number, status="ACTIVE")


@router.get("/persons/violence/query")
async def query_violence_records(
    niu: Optional[str] = Query(None),
    incident_type: Optional[str] = Query(None),
    status: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


@router.get("/persons/violence/{record_id}")
async def get_violence_record(record_id: str):
    return {"record_id": record_id, "status": "ACTIVE"}


# ── PER-IDV: Identity Theft ─────────────────────────────────────────────────


@router.post("/persons/identity-theft", response_model=IdentityTheftResponse, status_code=201)
async def create_identity_theft(request: CreateIdentityTheftRequest):
    record_id = str(uuid.uuid4())
    record_number = f"IDV-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = IdentityTheftRecorded(
        record_id=record_id,
        record_number=record_number,
        victim_niu=request.victim_niu,
        fraud_type=request.fraud_type.value,
        reporting_agency=request.reporting_agency,
    )
    await _publish(event, key=record_id)

    return IdentityTheftResponse(record_id=uuid.UUID(record_id), record_number=record_number, status="ACTIVE")


@router.get("/persons/identity-theft/query")
async def query_identity_thefts(
    victim_niu: Optional[str] = Query(None),
    fraud_type: Optional[str] = Query(None),
    status: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


@router.get("/persons/identity-theft/{record_id}")
async def get_identity_theft(record_id: str):
    return {"record_id": record_id, "status": "ACTIVE"}


# ── Stolen Property ────────────────────────────────────────────────────────


@router.post("/property/vehicles", status_code=201)
async def report_stolen_vehicle(request: StolenVehicleRequest):
    record_id = str(uuid.uuid4())
    record_number = f"VEH-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = VehicleStolen(
        record_id=record_id,
        record_number=record_number,
        plate_number=request.plate_number,
        vin=request.vin or "",
    )
    await _publish(event, key=record_id)

    return {
        "record_id": record_id,
        "record_number": record_number,
        "status": StolenStatus.STOLEN.value,
    }


@router.patch("/property/vehicles/{record_id}/recover", status_code=200)
async def recover_stolen_vehicle(record_id: str, request: VehicleRecoverRequest):
    record_number = f"VEH-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"
    event = VehicleRecovered(
        record_id=record_id,
        record_number=record_number,
        recovered_location=request.recovered_location,
    )
    await _publish(event, key=record_id)
    return {
        "success": True,
        "record_id": record_id,
        "status": StolenStatus.RECOVERED.value,
        "recovered_at": datetime.now(timezone.utc).isoformat(),
    }


@router.post("/property/firearms", status_code=201)
async def report_stolen_firearm(request: StolenFirearmRequest):
    record_id = str(uuid.uuid4())
    record_number = f"ARM-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = FirearmStolen(
        record_id=record_id,
        record_number=record_number,
        serial_number=request.serial_number,
    )
    await _publish(event, key=record_id)

    return {
        "record_id": record_id,
        "record_number": record_number,
        "status": StolenStatus.STOLEN.value,
    }


@router.post("/property/firearms/{record_id}/crime-scene-hit", status_code=201)
async def report_arm_crime_scene_hit(
    record_id: str,
    crime_scene_ref: str,
    case_number: str,
):
    hit_id = str(uuid.uuid4())
    event = ArmCrimeSceneHit(
        record_id=record_id,
        record_number="",
        case_number=case_number,
    )
    await _publish(event, key=hit_id)
    return {"hit_id": hit_id, "status": "dispatched"}


@router.post("/property/documents", status_code=201)
async def report_stolen_document(request: StolenDocumentRequest):
    record_id = str(uuid.uuid4())
    record_number = f"DOC-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = DocumentStolen(
        record_id=record_id,
        record_number=record_number,
        document_type=request.document_type,
        document_number=request.document_number or "",
    )
    await _publish(event, key=record_id)

    return {
        "record_id": record_id,
        "record_number": record_number,
        "status": "ACTIVE",
    }


@router.post("/property/documents/oni-revoke", status_code=200, response_model=None)
async def oni_document_revoke(
    document_number: str = Body(...),
    revocation_reason: str = Body(...),
):
    event = ONIDocumentRevoked(document_number=document_number, revocation_reason=revocation_reason)
    await _publish(event, key=document_number)
    return {"success": True, "message": "Document révoqué et ajouté à BIE-DOC"}


@router.post("/property/vessels", status_code=201)
async def report_stolen_vessel(request: StolenVesselRequest):
    record_id = str(uuid.uuid4())
    record_number = f"EMB-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"

    event = VesselStolen(
        record_id=record_id,
        record_number=record_number,
        vessel_name=request.vessel_name or "",
        registration_number=request.registration_number or "",
    )
    await _publish(event, key=record_id)

    return {
        "record_id": record_id,
        "record_number": record_number,
        "status": StolenStatus.STOLEN.value,
    }


@router.get("/property/query")
async def query_stolen_property(
    plate_number: Optional[str] = Query(None),
    vin: Optional[str] = Query(None),
    serial_number: Optional[str] = Query(None),
    document_number: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


# ── BIE-OBJ (Articles / Biens Précieux) ───────────────────────────────────


@router.post("/property/articles", status_code=201, response_model=StolenArticleResponse)
async def report_stolen_article(request: StolenArticleRequest):
    record_id = str(uuid.uuid4())
    record_number = f"OBJ-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"
    event = StolenArticleCreated(
        record_id=record_id,
        record_number=record_number,
        category=request.category.value,
    )
    await _publish(event, key=record_id)
    return StolenArticleResponse(
        record_id=uuid.UUID(record_id),
        record_number=record_number,
        category=request.category.value,
        status=StolenStatus.STOLEN.value,
    )


@router.get("/property/articles/query")
async def query_stolen_articles(
    category: Optional[str] = Query(None),
    serial_number: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


# ── BIE-TIT (Titres et Valeurs) ───────────────────────────────────────────


@router.post("/property/securities", status_code=201, response_model=StolenSecurityResponse)
async def report_stolen_security(request: StolenSecurityRequest):
    record_id = str(uuid.uuid4())
    record_number = f"TIT-{datetime.now().year}-{uuid.uuid4().hex[:6].upper()}"
    event = StolenSecurityCreated(
        record_id=record_id,
        record_number=record_number,
        security_type=request.security_type.value,
    )
    await _publish(event, key=record_id)
    return StolenSecurityResponse(
        record_id=uuid.UUID(record_id),
        record_number=record_number,
        security_type=request.security_type.value,
        status=StolenStatus.STOLEN.value,
    )


@router.get("/property/securities/query")
async def query_stolen_securities(
    security_type: Optional[str] = Query(None),
    security_number: Optional[str] = Query(None),
):
    return {"total": 0, "results": []}


# ── LAPI (Real-time Plate Query, SLA < 200ms) ─────────────────────────────


@router.get("/lapi/plate/{plate_number}", response_model=LAPIPlateResponse)
async def query_plate_lapi(
    plate_number: str,
    camera_id: Optional[str] = Query(None),
    location: Optional[str] = Query(None),
):
    query_id = str(uuid.uuid4())
    return LAPIPlateResponse(
        query_id=query_id,
        hit_found=False,
        response_ms=0,
    )


@router.get("/lapi/vin/{vin}", response_model=LAPIPlateResponse)
async def query_vin_lapi(vin: str):
    query_id = str(uuid.uuid4())
    return LAPIPlateResponse(
        query_id=query_id,
        hit_found=False,
        response_ms=0,
    )


# ── Lab Equipment ────────────────────────────────────────────────────────────


@router.post("/lab/equipment", response_model=LabEquipmentResponse, status_code=201)
async def register_equipment(request: LabEquipmentCreate):
    equipment_id = str(uuid.uuid4())
    event = EquipmentRegistered(
        equipment_id=equipment_id,
        lab_code=request.lab_code,
        equipment_name=request.equipment_name,
    )
    await _publish(event, key=equipment_id)
    return LabEquipmentResponse(
        id=uuid.UUID(equipment_id),
        equipment_name=request.equipment_name,
        serial_number=request.serial_number,
        status="ACTIVE",
    )


@router.get("/lab/equipment")
async def list_equipment(lab_code: Optional[str] = Query(None)):
    return {"total": 0, "results": []}


@router.get("/lab/equipment/{equipment_id}")
async def get_equipment(equipment_id: str):
    return {"id": equipment_id, "status": "ACTIVE"}


@router.patch("/lab/equipment/{equipment_id}/calibration")
async def update_calibration(equipment_id: str, request: LabEquipmentUpdate):
    return {"id": equipment_id, "calibration_updated": True}


# ── Staff Training ────────────────────────────────────────────────────────────


@router.post("/lab/training", response_model=StaffTrainingResponse, status_code=201)
async def record_training(request: StaffTrainingCreate):
    training_id = str(uuid.uuid4())
    event = TrainingRecorded(
        training_id=training_id,
        staff_niu=request.staff_niu,
        training_name=request.training_name,
    )
    await _publish(event, key=training_id)
    return StaffTrainingResponse(
        id=uuid.UUID(training_id),
        training_name=request.training_name,
        training_code=request.training_code,
    )


@router.get("/lab/training")
async def list_training(staff_niu: Optional[str] = Query(None)):
    return {"total": 0, "results": []}


@router.get("/lab/training/{training_id}")
async def get_training(training_id: str):
    return {"id": training_id}


# ── LDIS Upload (CLI-style) ─────────────────────────────────────────────────


@router.post("/lab/upload", response_model=UploadResponse, status_code=202)
async def ldis_upload(request: UploadRequest):
    count = 0
    event = LDISUploadCompleted(
        lab_code=request.lab_code,
        uploaded_count=count,
        operator_niu=request.operator_niu,
    )
    await _publish(event, key=request.lab_code)
    return UploadResponse(success=True, uploaded_count=count, message="Upload LDIS→SDIS déclenché")


# ── Laboratory Info ──────────────────────────────────────────────────────────


# ── SDIS ──────────────────────────────────────────────────────────────────────


@router.get("/sdis/nodes", response_model=list[SdisNodeResponse])
async def sdis_nodes():
    return [
        {"code": "SDIS-OUEST", "department": "Ouest", "dc_location": "SNISID PAP", "dc_type": "DC principal", "status": "ACTIVE"},
        {"code": "SDIS-NORD", "department": "Nord", "dc_location": "SNISID CAP", "dc_type": "DC secondaire", "status": "ACTIVE"},
        {"code": "SDIS-ARTIBONITE", "department": "Artibonite", "dc_location": "Gonaïves", "dc_type": "Nœud SDIS", "status": "ACTIVE"},
        {"code": "SDIS-SUD", "department": "Sud", "dc_location": "Les Cayes", "dc_type": "Nœud SDIS", "status": "ACTIVE"},
        {"code": "SDIS-SUDEST", "department": "Sud-Est", "dc_location": "Jacmel", "dc_type": "Nœud SDIS", "status": "ACTIVE"},
        {"code": "SDIS-CENTRE", "department": "Centre", "dc_location": "Hinche", "dc_type": "Nœud SDIS", "status": "ACTIVE"},
        {"code": "SDIS-NORDEST", "department": "Nord-Est", "dc_location": "Fort-Liberté", "dc_type": "Nœud SDIS", "status": "À_CRÉER"},
        {"code": "SDIS-NORDOUEST", "department": "Nord-Ouest", "dc_location": "Port-de-Paix", "dc_type": "Nœud SDIS", "status": "À_CRÉER"},
        {"code": "SDIS-GRANDANSE", "department": "Grand-Anse", "dc_location": "Jérémie", "dc_type": "Nœud SDIS", "status": "À_CRÉER"},
        {"code": "SDIS-NIPPES", "department": "Nippes", "dc_location": "Miragoâne", "dc_type": "Nœud SDIS", "status": "À_CRÉER"},
    ]


@router.get("/sdis/stats/{sdis_code}", response_model=SdisStats)
async def sdis_stats(sdis_code: str):
    return SdisStats(
        sdis_code=sdis_code,
        total_profiles=0,
        intra_dept_matches=0,
        pending_reviews=0,
        sync_errors=0,
    )


@router.get("/sdis/errors/{sdis_code}")
async def sdis_errors(sdis_code: str):
    return {"total": 0, "results": []}


@router.get("/sdis/matches/{sdis_code}")
async def sdis_matches(sdis_code: str):
    return {"total": 0, "results": []}


@router.post("/sdis/security-alert", status_code=202)
async def sdis_security_alert(alert_type: str, sdis_code: str, details: str, sample_id: str):
    event = BioSecurityAlert(alert_type=alert_type, sdis_code=sdis_code, details=details, sample_id=sample_id)
    await _publish(event, key=sample_id)
    return {"alert_id": str(uuid.uuid4()), "status": "alerted"}


@router.post("/sdis/heartbeat", status_code=200)
async def sdis_heartbeat(sdis_code: str):
    event = SDISHeartbeat(sdis_code=sdis_code, node_healthy=True, upstream_ok=True)
    await _publish(event, key=sdis_code)
    return {"sdis_code": sdis_code, "last_heartbeat": datetime.now(timezone.utc).isoformat()}


# ── NDIS Stats ────────────────────────────────────────────────────────────────


@router.get("/ndis/stats", response_model=NDISStats)
async def ndis_stats():
    return NDISStats(
        total_bio_con=5000,
        total_bio_arr=1500,
        total_bio_fsc=2000,
        total_bio_dis=800,
        total_bio_rni=200,
        cross_dept_hits_this_week=12,
        hit_rate_percent=15.0,
        interpol_submissions_this_week=3,
    )


@router.get("/ndis/hits")
async def ndis_hits(sdis: Optional[str] = Query(None), match_type: Optional[str] = Query(None)):
    return {"total": 0, "results": []}


@router.get("/ndis/reports", response_model=list[NDISReportResponse])
async def ndis_reports():
    return []


@router.post("/ndis/reports/generate", status_code=202)
async def ndis_generate_report(report_type: str = Query(...)):
    report_id = str(uuid.uuid4())
    event = NDISReportGenerated(
        report_id=report_id,
        report_type=report_type,
        generated_by="system",
    )
    await _publish(event, key=report_id)
    return {"report_id": report_id, "status": "generating"}


# ── INTERPOL Gateway ──────────────────────────────────────────────────────────


@router.post("/ndis/interpol/submit", response_model=InterpolSubmissionResponse, status_code=202)
async def interpol_submit(request: InterpolSubmissionCreate):
    submission_id = str(uuid.uuid4())
    event = InterpolSubmissionRequested(
        submission_id=submission_id,
        sample_ids=request.sample_ids,
        reason=request.reason,
    )
    await _publish(event, key=submission_id)
    return InterpolSubmissionResponse(
        submission_id=uuid.UUID(submission_id),
        submitted_samples=len(request.sample_ids),
        status="submitted",
    )


@router.get("/lab/labs")
async def list_laboratories():
    return {
        "total": 6,
        "results": [
            {"code": "LDIS-PAP-001", "name": "Labo Médico-Légal PAP", "department": "Ouest", "status": "PRIORITAIRE"},
            {"code": "LDIS-CAP-001", "name": "Labo Médico-Légal Cap-Haïtien", "department": "Nord", "status": "PRIORITAIRE"},
            {"code": "LDIS-LES-001", "name": "Labo Médico-Légal Les Cayes", "department": "Sud", "status": "À_CRÉER"},
            {"code": "LDIS-GON-001", "name": "Labo Médico-Légal Gonaïves", "department": "Artibonite", "status": "À_CRÉER"},
            {"code": "LDIS-JAC-001", "name": "Labo Médico-Légal Jacmel", "department": "Sud-Est", "status": "À_CRÉER"},
            {"code": "LDIS-HIN-001", "name": "Labo Médico-Légal Hinche", "department": "Centre", "status": "À_CRÉER"},
        ],
    }
