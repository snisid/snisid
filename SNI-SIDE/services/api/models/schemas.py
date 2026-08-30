"""SNI-SIDE API — Pydantic models pour toutes les 15 bases"""

from datetime import datetime, date
from typing import Optional, List, Dict, Any
from pydantic import BaseModel, Field, ConfigDict


# ============ PAGINATION ============
class PaginationParams(BaseModel):
    page: int = Field(default=1, ge=1)
    limit: int = Field(default=20, ge=1, le=100)


class PaginatedResponse(BaseModel):
    total: int
    page: int
    limit: int
    data: List[Any]


# ============ NCID — CRIMINAL INTELLIGENCE ============
class WantedPersonCreate(BaseModel):
    niu: str = Field(pattern=r'^[A-Z0-9]{10}$')
    full_name: str
    alias: Optional[str] = None
    date_of_birth: Optional[date] = None
    place_of_birth: Optional[str] = None
    gender: Optional[str] = Field(None, pattern=r'^[MFO]$')
    nationality: Optional[str] = None
    height_cm: Optional[float] = None
    weight_kg: Optional[float] = None
    eye_color: Optional[str] = None
    hair_color: Optional[str] = None
    skin_tone: Optional[str] = None
    scars_marks: Optional[str] = None
    last_known_address: Optional[str] = None
    occupation: Optional[str] = None
    risk_level: str = Field(default="MEDIUM", pattern=r'^(CRITICAL|HIGH|MEDIUM|LOW)$')
    photos: List[str] = []


class WantedPerson(WantedPersonCreate):
    status: str = "ACTIVE"
    created_at: datetime
    updated_at: datetime

    model_config = ConfigDict(from_attributes=True)


class Warrant(BaseModel):
    warrant_id: str
    warrant_type: str
    issuing_authority: str
    person_niu: str
    case_id: Optional[str] = None
    charges: List[str]
    status: str
    issued_date: date
    expiry_date: Optional[date] = None
    risk_level: Optional[str] = None

    model_config = ConfigDict(from_attributes=True)


class CriminalCase(BaseModel):
    case_id: str
    case_number: str
    case_type: str
    case_category: Optional[str] = None
    status: str
    lead_agency: Optional[str] = None
    description: Optional[str] = None
    incident_date: Optional[date] = None
    opened_date: date
    latitude: Optional[float] = None
    longitude: Optional[float] = None

    model_config = ConfigDict(from_attributes=True)


class Gang(BaseModel):
    gang_id: str
    name: str
    alias: Optional[str] = None
    territory: Optional[str] = None
    criminal_activities: List[str] = []
    risk_level: Optional[str] = None
    member_count: Optional[int] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


class CriminalOrganization(BaseModel):
    org_id: str
    name: str
    type: str
    structure: Optional[str] = None
    geographic_reach: Optional[str] = None
    primary_activities: List[str] = []
    risk_level: Optional[str] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


class InterpolNotice(BaseModel):
    notice_id: str
    notice_type: str
    person_name: str
    nationality: Optional[str] = None
    issuing_country: str
    charges: List[str] = []
    status: str
    issued_date: date
    expiry_date: Optional[date] = None

    model_config = ConfigDict(from_attributes=True)


# ============ HN-NGI — BIOMETRICS ============
class BiometricVerificationRequest(BaseModel):
    niu: str
    biometric_type: str = Field(pattern=r'^(FACE|FINGERPRINT|IRIS|VOICE|PALM)$')
    template_data: str


class BiometricVerificationResponse(BaseModel):
    verified: bool
    match_score: float
    threshold: float
    search_duration_ms: int


class IdentifyBiometricRequest(BaseModel):
    biometric_type: str
    template_data: str
    gallery_id: Optional[str] = None
    max_candidates: int = Field(default=10, ge=1, le=100)
    threshold: float = Field(default=0.6, ge=0.0, le=1.0)


class BiometricCandidate(BaseModel):
    niu: str
    full_name: Optional[str] = None
    score: float
    rank: int


class IdentifyBiometricResponse(BaseModel):
    candidates: List[BiometricCandidate]
    gallery_size: int
    search_duration_ms: int


class DuplicateDetectionResult(BaseModel):
    niu_1: str
    niu_2: str
    score: float
    biometric_type: str
    status: str = "PENDING_REVIEW"


# ============ HN-CODIS — DNA ============
class DNAProfile(BaseModel):
    profile_id: str
    profile_type: str
    niu: Optional[str] = None
    sample_id: str
    laboratory_id: str
    profile_quality: str
    collection_date: date
    status: str

    model_config = ConfigDict(from_attributes=True)


class DNAMatch(BaseModel):
    match_id: str
    profile_id_1: str
    profile_id_2: str
    match_type: str
    match_probability: float
    relationship: Optional[str] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


# ============ MISSING PERSONS ============
class MissingPerson(BaseModel):
    missing_id: str
    full_name: str
    case_type: str
    date_of_birth: Optional[date] = None
    age_at_disappearance: Optional[int] = None
    gender: Optional[str] = None
    nationality: Optional[str] = None
    last_seen_date: datetime
    last_seen_location: str
    last_seen_lat: Optional[float] = None
    last_seen_lng: Optional[float] = None
    status: str
    risk_level: Optional[str] = None
    photos: List[str] = []

    model_config = ConfigDict(from_attributes=True)


class Sighting(BaseModel):
    sighting_id: str
    missing_id: str
    sighting_date: datetime
    sighting_location: str
    latitude: Optional[float] = None
    longitude: Optional[float] = None
    source_type: str
    confidence_score: Optional[float] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


class TriggerAlertRequest(BaseModel):
    missing_id: str
    alert_type: str = Field(pattern=r'^(AMBER|SILVER|CLEAR|NATIONAL)$')
    channels: List[str] = Field(min_length=1)


# ============ VEHICLE INTELLIGENCE ============
class Vehicle(BaseModel):
    vehicle_id: str
    vin: str
    plate_number: Optional[str] = None
    make: str
    model: str
    year: int
    color: Optional[str] = None
    body_type: Optional[str] = None
    owner_niu: Optional[str] = None
    owner_name: Optional[str] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


class OwnershipRecord(BaseModel):
    ownership_id: str
    vehicle_id: str
    owner_niu: str
    owner_name: str
    ownership_start: date
    ownership_end: Optional[date] = None

    model_config = ConfigDict(from_attributes=True)


class VehicleNetwork(BaseModel):
    vin: str
    persons: List[dict] = []
    vehicles: List[dict] = []
    phones: List[dict] = []
    connections: List[dict] = []


# ============ ALPR ============
class ALPRRead(BaseModel):
    read_id: str
    plate_text: str
    plate_country: str = "HT"
    camera_code: str
    latitude: Optional[float] = None
    longitude: Optional[float] = None
    read_timestamp: datetime
    speed_kmh: Optional[int] = None
    vehicle_make: Optional[str] = None
    vehicle_model: Optional[str] = None
    vehicle_color: Optional[str] = None
    ocr_confidence: float
    alert_triggered: bool = False

    model_config = ConfigDict(from_attributes=True)


class RoutePoint(BaseModel):
    camera_code: str
    timestamp: datetime
    latitude: float
    longitude: float


class RouteAnalysis(BaseModel):
    plate: str
    points: List[RoutePoint]
    total_distance_km: float
    first_seen: datetime
    last_seen: datetime
    frequent_areas: List[str]


# ============ FIREARMS ============
class Firearm(BaseModel):
    firearm_id: str
    serial_number: str
    make: str
    model: str
    caliber: str
    type: str
    owner_niu: Optional[str] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


class BallisticEvidence(BaseModel):
    ballistic_id: str
    evidence_number: str
    caliber: str
    case_id: Optional[str] = None

    model_config = ConfigDict(from_attributes=True)


# ============ BORDER ============
class BorderCrossing(BaseModel):
    crossing_id: str
    niu: Optional[str] = None
    passport_number: Optional[str] = None
    full_name: str
    nationality: str
    direction: str
    border_point: str
    crossing_date: datetime
    crossing_method: str
    visa_type: Optional[str] = None
    risk_score: float = 0.0
    alert_triggered: bool = False

    model_config = ConfigDict(from_attributes=True)


class Visa(BaseModel):
    visa_id: str
    visa_number: str
    visa_type: str
    issuing_post: str
    issue_date: date
    expiry_date: date
    status: str

    model_config = ConfigDict(from_attributes=True)


# ============ NARCOTICS ============
class NarcoticsRoute(BaseModel):
    route_id: str
    route_type: str
    origin_country: str
    destination_country: str
    primary_drugs: List[str] = []
    risk_level: Optional[str] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


class Seizure(BaseModel):
    seizure_id: str
    seizure_date: datetime
    location: str
    drug_types: List[str]
    estimated_value_usd: Optional[float] = None
    seizure_agency: str
    arrests_made: int = 0

    model_config = ConfigDict(from_attributes=True)


class Cartel(BaseModel):
    cartel_id: str
    name: str
    country_of_origin: Optional[str] = None
    primary_drugs: List[str] = []
    risk_level: Optional[str] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


# ============ FINANCIAL ============
class SuspiciousTransaction(BaseModel):
    transaction_id: str
    transaction_ref: str
    transaction_date: datetime
    transaction_type: str
    amount: float
    currency: str
    sender_name: str
    beneficiary_name: str
    source_country: Optional[str] = None
    destination_country: Optional[str] = None
    risk_score: float = 0.0
    status: str

    model_config = ConfigDict(from_attributes=True)


class BeneficialOwner(BaseModel):
    owner_id: str
    niu: str
    full_name: str
    entity_name: str
    ownership_percentage: Optional[float] = None
    politically_exposed: bool = False
    risk_score: float = 0.0

    model_config = ConfigDict(from_attributes=True)


class PEP(BaseModel):
    pep_id: str
    niu: str
    full_name: str
    position: str
    institution: Optional[str] = None
    country: Optional[str] = None
    pep_level: str
    risk_score: float = 0.0

    model_config = ConfigDict(from_attributes=True)


class AMLRiskFactor(BaseModel):
    factor: str
    score: float
    description: str


class AMLAnalysisResponse(BaseModel):
    overall_risk_score: float
    risk_level: str
    risk_factors: List[AMLRiskFactor]
    recommendations: List[str]


# ============ CYBERCRIME ============
class IOC(BaseModel):
    ioc_id: str
    ioc_value: str
    ioc_type: str
    confidence: int = 50
    severity: str = "MEDIUM"
    malware_family: Optional[str] = None
    threat_actor: Optional[str] = None
    tlp_level: str = "AMBER"
    last_seen: Optional[datetime] = None
    status: str = "ACTIVE"

    model_config = ConfigDict(from_attributes=True)


class ThreatActor(BaseModel):
    actor_id: str
    name: str
    alias: List[str] = []
    actor_type: str
    country_of_origin: Optional[str] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


class CryptoWallet(BaseModel):
    wallet_id: str
    wallet_address: str
    blockchain: str
    wallet_type: Optional[str] = None
    risk_score: float = 0.0
    total_received_usd: Optional[float] = None

    model_config = ConfigDict(from_attributes=True)


# ============ WATCHLIST ============
class WatchlistEntry(BaseModel):
    entry_id: str
    watchlist_type: str
    entry_category: str
    value_primary: str
    full_name: Optional[str] = None
    niu: Optional[str] = None
    document_number: Optional[str] = None
    phone_number: Optional[str] = None
    vehicle_plate: Optional[str] = None
    risk_level: str
    listing_authority: str
    reason: str
    status: str

    model_config = ConfigDict(from_attributes=True)


class WatchlistMatch(BaseModel):
    match_id: str
    entry_id: str
    match_value: str
    match_source: str
    match_confidence: float
    detected_at: datetime
    status: str

    model_config = ConfigDict(from_attributes=True)


# ============ DOCUMENT FRAUD ============
class FraudDocument(BaseModel):
    document_id: str
    document_number: str
    document_type: str
    holder_niu: Optional[str] = None
    holder_name: str
    issuing_country: Optional[str] = None
    status: str
    risk_score: float = 0.0

    model_config = ConfigDict(from_attributes=True)


class DocumentFraudReport(BaseModel):
    fraud_id: str
    document_id: str
    fraud_type: str
    detection_date: datetime
    detection_method: str
    confidence_score: Optional[float] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


class DocumentVerificationResponse(BaseModel):
    authentic: bool
    confidence: float
    anomalies: List[str] = []
    verification_method: str


# ============ GEOINT ============
class Hotspot(BaseModel):
    hotspot_id: str
    name: Optional[str] = None
    type: str
    latitude: float
    longitude: float
    risk_level: Optional[str] = None
    last_incident_date: Optional[datetime] = None

    model_config = ConfigDict(from_attributes=True)


class RiskZone(BaseModel):
    zone_id: str
    zone_name: str
    zone_type: str
    risk_level: Optional[str] = None
    risk_score: Optional[float] = None

    model_config = ConfigDict(from_attributes=True)


# ============ EVIDENCE ============
class EvidenceItem(BaseModel):
    evidence_id: str
    evidence_number: str
    evidence_type: str
    case_number: Optional[str] = None
    title: str
    mime_type: Optional[str] = None
    file_size: Optional[int] = None
    captured_date: Optional[datetime] = None
    captured_by_agency: Optional[str] = None
    status: str

    model_config = ConfigDict(from_attributes=True)


class EvidenceFaceMatch(BaseModel):
    evidence_id: str
    confidence: float
    person_niu: Optional[str] = None
    bbox: Optional[dict] = None


class MultimodalSearchResult(BaseModel):
    evidence_id: str
    modality: str
    score: float
    description: Optional[str] = None


# ============ SEARCH ============
class UnifiedSearchQuery(BaseModel):
    q: str
    type: str = "ALL"
    databases: List[str] = []
    page: int = Field(default=1, ge=1)
    limit: int = Field(default=20, ge=1, le=100)
    fuzzy: bool = True


class SearchResultItem(BaseModel):
    database: str
    result_type: str
    id: str
    title: str
    description: str = ""
    score: float = 0.0
    confidence: float = 0.0
    risk_score: float = 0.0
    metadata: dict = {}


class UnifiedSearchResponse(BaseModel):
    query: str
    total_results: int
    page: int
    limit: int
    search_duration_ms: float
    databases_searched: int
    results: Dict[str, List[SearchResultItem]]
    graph_context: Optional[dict] = None
    suggested_queries: List[str] = []


class GraphSearchQuery(BaseModel):
    niu: Optional[str] = None
    phone: Optional[str] = None
    plate: Optional[str] = None
    depth: int = Field(default=2, ge=1, le=5)
    relationship_types: List[str] = []


class GraphNode(BaseModel):
    id: str
    label: str
    type: str
    properties: dict = {}
    risk_score: float = 0.0


class GraphEdge(BaseModel):
    source: str
    target: str
    relationship: str
    weight: float = 1.0
    properties: dict = {}


class GraphSearchResponse(BaseModel):
    nodes: List[GraphNode]
    edges: List[GraphEdge]
    centrality_score: float = 0.0
    network_size: int = 0
    detected_patterns: List[str] = []
    analysis_summary: Optional[str] = None


# ============ ALERTS ============
class Alert(BaseModel):
    alert_id: str
    source: str
    alert_type: str
    severity: str
    title: str
    description: str
    entity_ids: List[str] = []
    status: str = "NEW"
    created_at: datetime
    acknowledged_at: Optional[datetime] = None
    acknowledged_by: Optional[str] = None
    resolution: Optional[str] = None

    model_config = ConfigDict(from_attributes=True)


class AcknowledgeAlertResponse(BaseModel):
    success: bool
    alert_id: str


class ResolveAlertRequest(BaseModel):
    resolution: str
    action_taken: Optional[str] = None


# ============ AI FUSION ============
class FraudAnalysisRequest(BaseModel):
    entity_niu: str
    graph_depth: int = Field(default=2, ge=1, le=5)


class FraudAnalysisResponse(BaseModel):
    risk_score: float
    risk_level: str
    fraud_indicators: List[str]
    feature_importance: Dict[str, float]
    model_version: str
    network_context: Optional[dict] = None


class GraphRAGQuery(BaseModel):
    query: str
    seed_entity_id: Optional[str] = None
    seed_entity_type: Optional[str] = None
    graph_depth: int = Field(default=2, ge=1, le=5)


class GraphRAGResponse(BaseModel):
    analysis: str
    context_nodes: List[GraphNode] = []
    context_edges: List[GraphEdge] = []
    overall_risk_score: float = 0.0
    key_findings: List[str] = []
    confidence: float = 0.0
