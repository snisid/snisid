from __future__ import annotations

from enum import Enum
from typing import Any, Optional
from uuid import UUID

from pydantic import BaseModel, Field, field_validator


class IndexType(str, Enum):
    BIO_CON = "BIO-CON"
    BIO_ARR = "BIO-ARR"
    BIO_FSC = "BIO-FSC"
    BIO_DIS = "BIO-DIS"
    BIO_RNI = "BIO-RNI"


class MatchType(str, Enum):
    FULL_MATCH = "FULL_MATCH"
    PARTIAL = "PARTIAL"
    FAMILIAL = "FAMILIAL"


class AlertLevel(str, Enum):
    LOW = "LOW"
    MEDIUM = "MEDIUM"
    HIGH = "HIGH"
    CRITICAL = "CRITICAL"


class WarrantType(str, Enum):
    ARREST = "MAN-ARR"
    EXTRADITION = "MAN-EXT"
    SEARCH = "MAN-REC"
    NOTICE = "AVIS-REC"


class DangerLevel(str, Enum):
    LOW = "LOW"
    MEDIUM = "MEDIUM"
    HIGH = "HIGH"
    CRITICAL = "CRITICAL"


class MissingCategory(str, Enum):
    CHILD = "CHILD"
    ENDANGERED = "ENDANGERED"
    INVOLUNTARY = "INVOLUNTARY"
    CATASTROPHE = "CATASTROPHE"
    UNEMANCIPATED = "UNEMANCIPATED"
    OTHER = "OTHER"


class WarrantStatus(str, Enum):
    ACTIVE = "ACTIVE"
    CLEARED = "CLEARED"
    EXPIRED = "EXPIRED"
    SUSPENDED = "SUSPENDED"


class MissingStatus(str, Enum):
    ACTIVE = "ACTIVE"
    LOCATED = "LOCATED"
    DECEASED = "DECEASED"
    CANCELLED = "CANCELLED"


class StolenStatus(str, Enum):
    STOLEN = "STOLEN"
    RECOVERED = "RECOVERED"
    CANCELLED = "CANCELLED"


class FirearmType(str, Enum):
    PISTOL = "PISTOL"
    REVOLVER = "REVOLVER"
    RIFLE = "RIFLE"
    SHOTGUN = "SHOTGUN"
    MACHINEGUN = "MACHINEGUN"
    EXPLOSIVE = "EXPLOSIVE"
    OTHER = "OTHER"


class ArticleCategory(str, Enum):
    JEWELRY = "JEWELRY"
    ART = "ART"
    ELECTRONICS = "ELECTRONICS"
    CURRENCY = "CURRENCY"
    CATTLE = "CATTLE"
    MACHINERY = "MACHINERY"
    OTHER = "OTHER"


class SecurityType(str, Enum):
    CHEQUE = "CHEQUE"
    BOND = "BOND"
    PROPERTY_TITLE = "PROPERTY_TITLE"
    LETTER_CREDIT = "LETTER_CREDIT"
    OTHER = "OTHER"


class VesselType(str, Enum):
    FISHING_CANOE = "FISHING_CANOE"
    MOTORBOAT = "MOTORBOAT"
    SAILBOAT = "SAILBOAT"
    FERRY = "FERRY"
    CARGO_SMALL = "CARGO_SMALL"
    PATROL_BOAT = "PATROL_BOAT"
    OTHER = "OTHER"


class InterpolNoticeType(str, Enum):
    RED = "RED"
    BLUE = "BLUE"
    YELLOW = "YELLOW"
    BLACK = "BLACK"
    ORANGE = "ORANGE"
    PURPLE = "PURPLE"
    INTERPOL_UNKNOWN = "UNKNOWN"


class SexOffenderRiskLevel(str, Enum):
    LOW = "LOW"
    MEDIUM = "MEDIUM"
    HIGH = "HIGH"


class ProtectionOrderType(str, Enum):
    RESTRAINING = "RESTRAINING"
    HARASSMENT = "HARASSMENT"
    CHILD_PROTECTION = "CHILD_PROTECTION"
    DOMESTIC_VIOLENCE = "DOMESTIC_VIOLENCE"
    EMERGENCY = "EMERGENCY"
    OTHER = "OTHER"


class SupervisionType(str, Enum):
    CONDITIONAL_RELEASE = "CONDITIONAL_RELEASE"
    PAROLE = "PAROLE"
    PROBATION = "PROBATION"
    JUDICIAL_CONTROL = "JUDICIAL_CONTROL"
    OTHER_SUPERVISION = "OTHER"


# ── CODIS STR Locus Models ──────────────────────────────────────────────────


class STRLocus(BaseModel):
    locus: str
    value1: float
    value2: float | None = None


class STRLociData(BaseModel):
    CSF1PO: STRLocus
    D3S1358: STRLocus
    D5S818: STRLocus
    D7S820: STRLocus
    D8S1179: STRLocus
    D13S317: STRLocus
    D16S539: STRLocus
    D18S51: STRLocus
    D21S11: STRLocus
    FGA: STRLocus
    TH01: STRLocus
    TPOX: STRLocus
    vWA: STRLocus
    D1S1656: STRLocus
    D2S441: STRLocus
    D2S1338: STRLocus
    D10S1248: STRLocus
    D12S391: STRLocus
    D19S433: STRLocus
    D22S1045: STRLocus


class STRProfileCreate(BaseModel):
    specimen_number: str = Field(..., min_length=5)
    index_type: IndexType
    loci_data: STRLociData
    amelogenin: str | None = Field(None, pattern=r"^(XX|XY)$")
    quality_score: float = Field(..., ge=0, le=1)
    lab_id: UUID
    case_number: str | None = None
    collected_date: str
    correlation_id: str


class STRProfileDB(STRProfileCreate):
    sample_id: UUID
    loci_encrypted: bytes
    loci_hash: str
    uploaded_ldis: bool = False
    uploaded_sdis: bool = False
    uploaded_ndis: bool = False
    is_expunged: bool = False
    created_at: str
    updated_at: str


# ── Wanted Person (arch ref) ────────────────────────────────────────────────


class WantedPersonCreate(BaseModel):
    niu: str | None = None
    last_name: str | None = None
    first_name: str | None = None
    aliases: list[str] = []
    date_of_birth: str | None = None
    gender: str | None = Field(None, pattern=r"^[MFU]$")
    nationality: str | None = Field(None, min_length=3, max_length=3)
    warrant_type: WarrantType
    warrant_number: str | None = None
    issuing_court: str | None = None
    issuing_date: str
    charges: list[str] = Field(..., min_length=1)
    danger_level: DangerLevel = DangerLevel.MEDIUM
    armed_dangerous: bool = False
    height_cm: int | None = Field(None, ge=50, le=250)
    weight_kg: int | None = Field(None, ge=10, le=300)
    mco_contact: str
    expiry_date: str | None = None
    interpol_notice: str | None = None


# ── Missing Person (arch ref) ────────────────────────────────────────────────


class MissingPersonCreate(BaseModel):
    niu: str | None = None
    last_name: str
    first_name: str
    date_of_birth: str | None = None
    age_at_missing: int | None = None
    gender: str | None = None
    category: str = Field(..., pattern=r"^(CHILD|ENDANGERED|INVOLUNTARY|CATASTROPHE|OTHER)$")
    missing_date: str
    missing_location: str
    circumstances: str | None = None
    height_cm: int | None = None
    weight_kg: int | None = None
    family_contact: str | None = None
    family_phone: str | None = None
    medical_conditions: str | None = None
    entering_agency: str


# ── Stolen Vehicle (arch ref) ────────────────────────────────────────────────


class StolenVehicleCreate(BaseModel):
    vin: str | None = Field(None, pattern=r"^[A-HJ-NPR-Z0-9]{17}$")
    plate_number: str
    plate_dept: str | None = None
    vehicle_make: str
    vehicle_model: str
    vehicle_year: int = Field(..., ge=1950, le=2030)
    vehicle_color: str
    vehicle_type: str | None = None
    theft_date: str
    theft_location: str
    theft_department: str | None = None
    owner_niu: str | None = None
    owner_name: str | None = None
    owner_phone: str | None = None
    entering_agency: str


# ── Stolen Vessel (arch ref) ────────────────────────────────────────────────


class StolenVesselCreate(BaseModel):
    vessel_name: str | None = None
    registration_number: str | None = None
    hull_id_number: str | None = None
    vessel_type: str = Field(..., pattern=r"^(FISHING_CANOE|MOTORBOAT|SAILBOAT|FERRY|CARGO_SMALL|PATROL_BOAT|OTHER)$")
    vessel_make: str | None = None
    vessel_length_m: float | None = None
    hull_color: str | None = None
    home_port: str | None = None
    theft_location: str
    theft_date: str
    owner_niu: str | None = None
    owner_name: str | None = None
    entering_agency: str


REQUIRED_LOCIS = {
    "CSF1PO", "D3S1358", "D5S818", "D7S820", "D8S1179",
    "D13S317", "D16S539", "D18S51", "D21S11", "FGA",
    "TH01", "TPOX", "vWA", "D1S1656", "D2S441",
    "D2S1338", "D10S1248", "D12S391", "D19S433", "D22S1045",
}


class STRProfileSubmit(BaseModel):
    specimen_number: str = Field(..., min_length=5, max_length=100)
    index_type: IndexType
    loci_data: dict[str, Any] = Field(..., description="20 loci STR {locus: {value1, value2}}")
    quality_score: float = Field(..., ge=0.0, le=1.0)
    case_number: Optional[str] = None
    collected_date: str = Field(..., description="YYYY-MM-DD")
    correlation_id: str

    @field_validator("loci_data")
    @classmethod
    def validate_loci(cls, v):
        missing = REQUIRED_LOCIS - set(v.keys())
        if missing:
            raise ValueError(f"Loci manquants: {missing}")
        return v


class STRProfileResponse(BaseModel):
    sample_id: UUID
    accepted: bool
    rejection_reason: Optional[str] = None
    message: str


class DNASearchRequest(BaseModel):
    loci_data: dict[str, Any]
    index_type: str = "BIO-FSC"
    case_number: str = Field(..., min_length=5)
    purpose: str = Field(..., pattern=r"^(criminal_investigation|missing_person|identification|mass_disaster)$")
    min_confidence: float = Field(0.85, ge=0.60, le=1.0)
    include_familial: bool = False


class DNASearchResponse(BaseModel):
    hits: list
    total_hits: int
    search_duration_ms: int
    case_number: str


class CreateWantedPersonRequest(BaseModel):
    niu: Optional[str] = None
    last_name: Optional[str] = None
    first_name: Optional[str] = None
    aliases: list[str] = []
    date_of_birth: Optional[str] = None
    gender: Optional[str] = Field(None, pattern=r"^[MFU]$")
    nationality: Optional[str] = None
    warrant_type: WarrantType
    warrant_number: Optional[str] = None
    issuing_court: Optional[str] = None
    issuing_date: str
    charges: list[str] = Field(..., min_length=1)
    danger_level: DangerLevel = DangerLevel.MEDIUM
    armed_dangerous: bool = False
    height_cm: Optional[int] = None
    weight_kg: Optional[int] = None
    mco_contact: str = Field(..., description="Contact agence entrante — obligatoire")
    entering_officer: Optional[str] = None
    expiry_date: Optional[str] = None

    @field_validator("last_name", "first_name")
    @classmethod
    def validate_identity(cls, v, info):
        values = info.data
        if not v and not values.get("last_name") and not values.get("first_name"):
            raise ValueError("Au moins un champ last_name ou first_name doit être présent")
        return v


class WantedPersonResponse(BaseModel):
    record_id: UUID
    record_number: str
    status: str
    mco_contact: str


class CreateMissingPersonRequest(BaseModel):
    niu: Optional[str] = None
    last_name: str = Field(..., min_length=1)
    first_name: str = Field(..., min_length=1)
    date_of_birth: Optional[str] = None
    age_at_missing: Optional[int] = None
    gender: Optional[str] = Field(None, pattern=r"^[MFU]$")
    nationality: Optional[str] = None
    category: MissingCategory
    missing_date: str
    missing_location: str = Field(..., min_length=1)
    circumstances: Optional[str] = None
    height_cm: Optional[int] = None
    weight_kg: Optional[int] = None
    distinctive_features: Optional[str] = None
    family_contact: Optional[str] = None
    family_phone: Optional[str] = None
    entering_agency: str = Field(..., min_length=1)
    citizen_portal_submission: bool = False
    dna_sample_available: bool = False


class MissingPersonResponse(BaseModel):
    record_id: UUID
    record_number: str
    status: str
    bpm_notified: bool = False


class StolenVehicleRequest(BaseModel):
    vin: Optional[str] = Field(None, pattern=r"^[A-HJ-NPR-Z0-9]{17}$")
    plate_number: str
    vehicle_make: str
    vehicle_model: str
    vehicle_year: int = Field(..., ge=1950, le=2030)
    vehicle_color: str
    theft_date: str
    theft_location: str
    theft_department: str
    owner_niu: Optional[str] = None
    owner_name: Optional[str] = None
    owner_phone: Optional[str] = None


class VehicleRecoverRequest(BaseModel):
    recovered_location: str
    recovering_agency: str
    notes: str = ""


class StolenFirearmRequest(BaseModel):
    serial_number: str
    make: Optional[str] = None
    model: Optional[str] = None
    caliber: Optional[str] = None
    firearm_type: Optional[FirearmType] = None
    barrel_length: Optional[float] = None
    theft_date: str
    theft_location: Optional[str] = None
    owner_niu: Optional[str] = None
    entering_agency: str


class StolenDocumentRequest(BaseModel):
    document_type: str = Field(..., pattern=r"^(PASSPORT|CIN|ACTE_NAISSANCE|PERMIS_CONDUIRE|TITRE_FONCIER|AUTRE)$")
    document_number: Optional[str] = None
    issuing_agency: Optional[str] = None
    issue_date: Optional[str] = None
    owner_niu: Optional[str] = None
    owner_name: Optional[str] = None
    report_date: str
    theft_type: str = Field("STOLEN", pattern=r"^(STOLEN|LOST|FORGED)$")


class StolenVesselRequest(BaseModel):
    vessel_name: Optional[str] = None
    registration_number: Optional[str] = None
    hull_id_number: Optional[str] = None
    vessel_type: Optional[VesselType] = None
    vessel_make: Optional[str] = None
    vessel_length_m: Optional[float] = None
    hull_color: Optional[str] = None
    home_port: Optional[str] = None
    engine_serial: Optional[str] = None
    distinctive_marks: Optional[str] = None
    theft_location: str
    theft_date: str
    owner_niu: Optional[str] = None
    owner_name: Optional[str] = None


class StolenArticleRequest(BaseModel):
    category: ArticleCategory
    description: str = Field(..., min_length=3)
    serial_number: Optional[str] = None
    estimated_value: Optional[float] = None
    currency_code: str = "HTG"
    theft_date: str
    theft_location: str
    theft_department: Optional[str] = None
    owner_niu: Optional[str] = None
    owner_name: Optional[str] = None
    entering_agency: str


class StolenArticleResponse(BaseModel):
    record_id: UUID
    record_number: str
    category: str
    status: str


class StolenSecurityRequest(BaseModel):
    security_type: SecurityType
    issuer: str
    security_number: str = Field(..., min_length=2)
    face_value: Optional[float] = None
    currency_code: str = "HTG"
    issue_date: Optional[str] = None
    theft_date: str
    theft_location: str
    owner_niu: Optional[str] = None
    owner_name: Optional[str] = None
    entering_agency: str


class StolenSecurityResponse(BaseModel):
    record_id: UUID
    record_number: str
    security_type: str
    status: str


# ── PER-FUG: Foreign Fugitives ─────────────────────────────────────────────


class CreateForeignFugitiveRequest(BaseModel):
    interpol_notice_number: str
    notice_type: InterpolNoticeType
    last_name: str = Field(..., min_length=1)
    first_name: Optional[str] = None
    aliases: list[str] = []
    date_of_birth: Optional[str] = None
    gender: Optional[str] = Field(None, pattern=r"^[MFU]$")
    nationality: Optional[str] = None
    charges: list[str] = Field(..., min_length=1)
    issuing_country: str = Field(..., min_length=2)
    entering_agency: str = Field(..., min_length=1)


class ForeignFugitiveResponse(BaseModel):
    record_id: UUID
    record_number: str
    interpol_notice_number: str
    status: str


# ── PER-NID: Unidentified Persons ──────────────────────────────────────────


class CreateUnidentifiedPersonRequest(BaseModel):
    discovery_date: str
    discovery_location: str = Field(..., min_length=1)
    discovery_department: Optional[str] = None
    estimated_age_min: Optional[int] = None
    estimated_age_max: Optional[int] = None
    gender: Optional[str] = Field(None, pattern=r"^[MFU]$")
    estimated_height_cm: Optional[int] = None
    estimated_weight_kg: Optional[int] = None
    hair_color: Optional[str] = None
    eye_color: Optional[str] = None
    distinctive_features: Optional[str] = None
    clothing_description: Optional[str] = None
    dna_sample_ref: Optional[str] = None
    fingerprint_ref: Optional[str] = None
    dental_records_ref: Optional[str] = None
    entering_agency: str = Field(..., min_length=1)


class UnidentifiedPersonResponse(BaseModel):
    record_id: UUID
    record_number: str
    status: str


# ── PER-TER: Terrorism ──────────────────────────────────────────────────────


class CreateTerrorismWatchRequest(BaseModel):
    niu: Optional[str] = None
    last_name: str = Field(..., min_length=1)
    first_name: Optional[str] = None
    aliases: list[str] = []
    date_of_birth: Optional[str] = None
    nationality: Optional[str] = None
    risk_level: str = "HIGH"
    threat_type: str = Field(..., description="Ex: RADICALISATION, FINANCEMENT, RECRUTEMENT")
    groups: list[str] = []
    known_associates: list[str] = []
    last_known_location: Optional[str] = None
    entering_agency: str = Field(..., min_length=1)
    approved_by_director: str
    approved_by_pg: str


class TerrorismWatchResponse(BaseModel):
    record_id: UUID
    record_number: str
    status: str


# ── PER-OPR: Protection Orders ─────────────────────────────────────────────


class CreateProtectionOrderRequest(BaseModel):
    order_type: ProtectionOrderType
    issuing_court: str
    issuing_judge: str
    beneficiary_niu: Optional[str] = None
    beneficiary_name: str = Field(..., min_length=1)
    protected_person: str = Field(..., min_length=1)
    restrained_person: str = Field(..., min_length=1)
    restrictions: list[str] = Field(..., min_length=1)
    issue_date: str
    expiration_date: Optional[str] = None
    emergency_contact: Optional[str] = None


class ProtectionOrderResponse(BaseModel):
    record_id: UUID
    record_number: str
    order_type: str
    status: str


# ── PER-LIB: Supervised Release ────────────────────────────────────────────


class CreateSupervisedReleaseRequest(BaseModel):
    niu: str = Field(..., min_length=1)
    last_name: str = Field(..., min_length=1)
    first_name: Optional[str] = None
    supervision_type: SupervisionType
    start_date: str
    end_date: Optional[str] = None
    conditions: list[str] = Field(..., min_length=1)
    supervising_officer: str
    supervising_agency: str


class SupervisedReleaseResponse(BaseModel):
    record_id: UUID
    record_number: str
    status: str


# ── Sex Offender Risk Management ────────────────────────────────────────────


class SexOffenderUpdateRequest(BaseModel):
    risk_level: SexOffenderRiskLevel
    current_address: Optional[str] = None
    employer: Optional[str] = None
    restrictions: Optional[list[str]] = None


# ── Gang Member Review ──────────────────────────────────────────────────────


class GangMemberReviewRequest(BaseModel):
    review_notes: str = ""
    auto_removal_years: int = 5


# ── Shared ──────────────────────────────────────────────────────────────────


class LAPIPlateResponse(BaseModel):
    query_id: str
    hit_found: bool
    hit_type: Optional[str] = None
    record_number: Optional[str] = None
    alert_level: Optional[str] = None
    mco_contact: Optional[str] = None
    response_ms: int
    clone_warning: Optional[bool] = None


# ── Lab Equipment ──────────────────────────────────────────────────────────


class LabEquipmentCreate(BaseModel):
    lab_code: str
    equipment_name: str
    model: Optional[str] = None
    serial_number: str
    role: str
    calibration_date: Optional[str] = None
    calibration_due: Optional[str] = None


class LabEquipmentResponse(BaseModel):
    id: UUID
    equipment_name: str
    serial_number: str
    status: str


class LabEquipmentUpdate(BaseModel):
    calibration_date: Optional[str] = None
    calibration_due: Optional[str] = None
    status: Optional[str] = None


# ── Staff Training ─────────────────────────────────────────────────────────


class StaffTrainingCreate(BaseModel):
    staff_niu: str
    training_name: str
    training_code: str
    duration_hours: Optional[int] = None
    completed_date: str
    valid_until: Optional[str] = None
    issued_by: str
    frequency: Optional[str] = None


class StaffTrainingResponse(BaseModel):
    id: UUID
    training_name: str
    training_code: str
    valid_until: Optional[str] = None


# ── Upload CLI ─────────────────────────────────────────────────────────────


class UploadRequest(BaseModel):
    level: str = Field(..., pattern=r"^(ldis-to-sdis)$")
    lab_code: str = Field(..., min_length=5)
    date_from: str
    date_to: str
    operator_niu: str = Field(..., min_length=5)


class UploadResponse(BaseModel):
    success: bool
    uploaded_count: int
    message: str


# ── NDIS ──────────────────────────────────────────────────────────────────────


class NDISCrossDeptHit(BaseModel):
    hit_id: str
    query_sample_id: str
    match_sample_id: str
    match_type: str
    confidence: float
    query_sdis: str
    match_sdis: str
    notified_at: Optional[str] = None


class NDISStats(BaseModel):
    total_bio_con: int
    total_bio_arr: int
    total_bio_fsc: int
    total_bio_dis: int
    total_bio_rni: int
    cross_dept_hits_this_week: int
    hit_rate_percent: float
    interpol_submissions_this_week: int


class NDISReportResponse(BaseModel):
    report_id: UUID
    report_type: str
    generated_at: str
    status: str


class InterpolSubmissionCreate(BaseModel):
    sample_ids: list[str] = Field(..., min_length=1)
    reason: str = Field(..., pattern=r"^(disaster_victim|international_fugitive|trafficking_victim|unidentified)$")
    case_number: Optional[str] = None


class InterpolSubmissionResponse(BaseModel):
    submission_id: UUID
    submitted_samples: int
    status: str


# ── SDIS ──────────────────────────────────────────────────────────────────────


class SdisNodeResponse(BaseModel):
    code: str
    department: str
    dc_location: str
    dc_type: str
    status: str
    last_heartbeat: Optional[str] = None


class SdisStats(BaseModel):
    sdis_code: str
    total_profiles: int
    intra_dept_matches: int
    pending_reviews: int
    sync_errors: int
    last_ndis_upload: Optional[str] = None


class SdisMatchResponse(BaseModel):
    id: UUID
    sdis_code: str
    match_type: str
    confidence: float
    alerted: bool


# ── PER-VIO: Known Violence ────────────────────────────────────────────────


class ViolenceIncidentType(str, Enum):
    DOMESTIC_VIOLENCE = "DOMESTIC_VIOLENCE"
    ASSAULT = "ASSAULT"
    BATTERY = "BATTERY"
    WEAPON_OFFENSE = "WEAPON_OFFENSE"
    HOMICIDE_ATTEMPT = "HOMICIDE_ATTEMPT"
    OTHER = "OTHER"


class CreateViolenceRecordRequest(BaseModel):
    niu: Optional[str] = None
    last_name: Optional[str] = None
    first_name: Optional[str] = None
    incident_type: ViolenceIncidentType
    incident_date: str
    location: str
    victim_niu: Optional[str] = None
    victim_name: Optional[str] = None
    arresting_agency: str
    court_case_ref: Optional[str] = None
    risk_level: str = "MEDIUM"


class ViolenceRecordResponse(BaseModel):
    record_id: UUID
    record_number: str
    status: str


# ── PER-IDV: Identity Theft ────────────────────────────────────────────────


class FraudType(str, Enum):
    CIN_FRAUD = "CIN_FRAUD"
    PASSPORT_FRAUD = "PASSPORT_FRAUD"
    FINANCIAL_FRAUD = "FINANCIAL_FRAUD"
    SOCIAL_MEDIA = "SOCIAL_MEDIA"
    OTHER = "OTHER"


class CreateIdentityTheftRequest(BaseModel):
    victim_niu: str = Field(..., min_length=5)
    victim_name: Optional[str] = None
    fraud_type: FraudType
    document_type_used: Optional[str] = None
    perpetrator_known: bool = False
    perpetrator_name: Optional[str] = None
    report_date: str
    reporting_agency: str


class IdentityTheftResponse(BaseModel):
    record_id: UUID
    record_number: str
    status: str


class AuditLogEntry(BaseModel):
    event_type: str
    table_name: str
    record_id: Optional[UUID] = None
    officer_niu: str
    agency_code: str
    purpose: str
    case_number: Optional[str] = None
    action: str
    details: Optional[dict] = None
