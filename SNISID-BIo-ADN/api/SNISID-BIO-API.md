# SNISID-BIO-ADN — API REST (FastAPI)
**Document ID :** SNISID-BIO-API-001 | **Version :** 1.0.0

---

## 1. BASE URL ET AUTHENTIFICATION

```
Base URL    : https://bio-adn.snisid.gov.ht/api/v1
Auth        : Bearer JWT (Keycloak SNISID) + mTLS pour services internes
Rate Limit  : 100 req/min par agent (sauf LAPI : illimité)
```

---

## 2. ENDPOINTS ADN

### POST /dna/profiles
Soumettre un nouveau profil STR ADN.

**Accès :** `bio.lab.technician`, `bio.lab.supervisor`

```python
# Exemple FastAPI
from fastapi import APIRouter, Depends, HTTPException, status
from pydantic import BaseModel, Field, validator
from typing import Optional
from uuid import UUID

router = APIRouter(prefix="/dna", tags=["DNA Profiles"])

class STRProfileSubmit(BaseModel):
    specimen_number: str = Field(..., min_length=5, max_length=100)
    index_type: str = Field(..., pattern="^(BIO-CON|BIO-ARR|BIO-FSC|BIO-DIS|BIO-RNI)$")
    loci_data: dict = Field(..., description="20 loci STR {locus: {value1, value2}}")
    quality_score: float = Field(..., ge=0.0, le=1.0)
    case_number: Optional[str] = None
    collected_date: str = Field(..., description="YYYY-MM-DD")
    correlation_id: str

    @validator("loci_data")
    def validate_loci(cls, v):
        required_loci = {
            "CSF1PO","D3S1358","D5S818","D7S820","D8S1179",
            "D13S317","D16S539","D18S51","D21S11","FGA",
            "TH01","TPOX","vWA","D1S1656","D2S441",
            "D2S1338","D10S1248","D12S391","D19S433","D22S1045"
        }
        missing = required_loci - set(v.keys())
        if missing:
            raise ValueError(f"Loci manquants: {missing}")
        return v

class STRProfileResponse(BaseModel):
    sample_id: UUID
    accepted: bool
    rejection_reason: Optional[str] = None
    message: str

@router.post("/profiles", response_model=STRProfileResponse, status_code=201)
async def submit_dna_profile(
    profile: STRProfileSubmit,
    current_user = Depends(require_role(["bio.lab.technician", "bio.lab.supervisor"])),
    grpc_client = Depends(get_grpc_client)
):
    """Soumet un profil STR ADN au système SNISID-BIO-ADN."""

    # Valider que l'agent appartient au bon lab
    if current_user.lab_id is None:
        raise HTTPException(status_code=403, detail="Aucun laboratoire associé à cet agent")

    # Appel gRPC vers le moteur de matching
    response = await grpc_client.SubmitProfile(
        specimen_number=profile.specimen_number,
        index_type=profile.index_type,
        loci_data=serialize_loci(profile.loci_data),
        quality_score=profile.quality_score,
        lab_id=str(current_user.lab_id),
        case_number=profile.case_number,
        collected_date=profile.collected_date,
        correlation_id=profile.correlation_id
    )

    return STRProfileResponse(
        sample_id=response.sample_id,
        accepted=response.accepted,
        rejection_reason=response.rejection_reason or None,
        message="Profil soumis avec succès" if response.accepted else "Profil rejeté"
    )
```

### POST /dna/search
Recherche de correspondances dans la base.

**Accès :** `bio.ndis.analyst`, `bio.dcpj.investigator` (avec `case_number` obligatoire)

```python
class DNASearchRequest(BaseModel):
    loci_data: dict
    index_type: str = "BIO-FSC"
    case_number: str = Field(..., min_length=5)
    purpose: str = Field(..., pattern="^(criminal_investigation|missing_person|identification|mass_disaster)$")
    min_confidence: float = Field(0.85, ge=0.60, le=1.0)
    include_familial: bool = False

class DNASearchResponse(BaseModel):
    hits: list
    total_hits: int
    search_duration_ms: int
    case_number: str

@router.post("/search", response_model=DNASearchResponse)
async def search_dna_profile(
    request: DNASearchRequest,
    current_user = Depends(require_role(["bio.ndis.analyst","bio.dcpj.investigator"]))
):
    """Recherche des correspondances ADN dans tous les index autorisés."""
    # ... implémentation
```

### GET /dna/hits/{hit_id}
Détails d'un hit ADN.

---

## 3. ENDPOINTS PERSONNES

### POST /persons/wanted
Créer un mandat de recherche.

**Accès :** `bio.dcpj.investigator`

```python
from enum import Enum

class WarrantType(str, Enum):
    ARREST      = "MAN-ARR"
    EXTRADITION = "MAN-EXT"
    SEARCH      = "MAN-REC"
    NOTICE      = "AVIS-REC"

class DangerLevel(str, Enum):
    LOW      = "LOW"
    MEDIUM   = "MEDIUM"
    HIGH     = "HIGH"
    CRITICAL = "CRITICAL"

class CreateWantedPersonRequest(BaseModel):
    niu: Optional[str] = None
    last_name: Optional[str] = None
    first_name: Optional[str] = None
    aliases: list[str] = []
    date_of_birth: Optional[str] = None
    gender: Optional[str] = Field(None, pattern="^[MFU]$")
    warrant_type: WarrantType
    warrant_number: Optional[str] = None
    issuing_court: Optional[str] = None
    issuing_date: str
    charges: list[str] = Field(..., min_items=1)
    danger_level: DangerLevel = DangerLevel.MEDIUM
    armed_dangerous: bool = False
    height_cm: Optional[int] = None
    weight_kg: Optional[int] = None
    mco_contact: str = Field(..., description="Contact agence entrante — obligatoire")
    expiry_date: Optional[str] = None

    @validator("last_name", "first_name")
    def validate_identity(cls, v, values):
        # Au moins un des deux doit être présent
        return v

@router.post("/wanted", status_code=201)
async def create_wanted_person(
    request: CreateWantedPersonRequest,
    current_user = Depends(require_role(["bio.dcpj.investigator"]))
):
    """Crée un enregistrement de personne recherchée (mandat d'arrêt, etc.)"""
    # ...
```

### GET /persons/wanted/query
Interrogation des mandats (usage terrain PNH).

**Accès :** Tous agents authentifiés

```python
@router.get("/wanted/query")
async def query_wanted_persons(
    last_name: Optional[str] = None,
    first_name: Optional[str] = None,
    niu: Optional[str] = None,
    plate_number: Optional[str] = None,
    current_user = Depends(get_current_user)
):
    """Interroge l'index des personnes recherchées."""
    # Audit log OBLIGATOIRE
    await audit_log(
        action="SEARCH",
        resource="per_wanted_persons",
        officer_niu=current_user.niu,
        agency=current_user.agency,
        purpose="field_query"
    )
    # ...
```

### POST /persons/missing
Signaler une disparition.

**Accès :** `bio.dcpj.investigator`, portail citoyen (avec validation agent J+24h)

---

## 4. ENDPOINTS BIENS

### POST /property/vehicles
Déclarer un véhicule volé.

```python
class StolenVehicleRequest(BaseModel):
    vin: Optional[str] = Field(None, regex="^[A-HJ-NPR-Z0-9]{17}$")
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

@router.post("/vehicles", status_code=201)
async def report_stolen_vehicle(
    request: StolenVehicleRequest,
    current_user = Depends(require_role(["bio.dcpj.investigator","pnh.officer"]))
):
    """Déclare un véhicule volé et synchronise avec FOVeS/SIV."""
    # ...
```

### POST /property/vessels
Déclarer une embarcation volée.

### GET /lapi/plate/{plate_number}
**SLA < 200ms** — Interrogation plaque LAPI (temps réel).

---

## 5. OPENAPI SPEC (extrait)

```yaml
openapi: 3.1.0
info:
  title: SNISID-BIO-ADN API
  version: 1.0.0
  description: |
    API souveraine nationale de base de données biométrique et criminelle d'Haïti.
    Équivalent CODIS+NCIC adapté au contexte haïtien.
  contact:
    name: Direction SNISID
    email: api@snisid.gov.ht

servers:
  - url: https://bio-adn.snisid.gov.ht/api/v1
    description: Production

security:
  - BearerAuth: []
  - mTLS: []

tags:
  - name: DNA Profiles
    description: Index ADN CODIS-HT (BIO-CON, BIO-ARR, BIO-FSC, BIO-DIS, BIO-RNI)
  - name: Wanted Persons
    description: Index personnes recherchées (PER-REC, PER-FUG)
  - name: Missing Persons
    description: Index personnes disparues (PER-DIS, PER-NID)
  - name: Property
    description: Index biens volés (BIE-VEH, BIE-ARM, BIE-DOC, BIE-EMB, etc.)
  - name: LAPI
    description: Interrogation temps réel pour LAPI (MP-16) — SLA 200ms
```

---

## 6. CODES D'ERREUR SPÉCIFIQUES

| Code | Nom | Description |
|------|-----|-------------|
| `BIO-001` | QUALITY_INSUFFICIENT | Quality score < seuil requis pour cet index |
| `BIO-002` | LOCI_INCOMPLETE | Nombre de loci insuffisant |
| `BIO-003` | LAB_NOT_ACCREDITED | Laboratoire non accrédité pour cet index |
| `BIO-004` | COURT_ORDER_REQUIRED | Ordonnance judiciaire requise pour cette opération |
| `BIO-005` | PROFILE_EXPUNGED | Profil effacé — opération impossible |
| `PER-001` | WARRANT_EXPIRED | Mandat expiré |
| `PER-002` | MCO_CONTACT_MISSING | Contact agence entrante obligatoire |
| `PER-003` | VERIFY_WITH_ENTERING_AGENCY | Vérifier avec l'agence avant arrestation |
| `LAPI-001`| SLA_BREACH | Réponse > 200ms — dégradé |
