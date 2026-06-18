# SNISID-BIO-ADN — Modèles de Données
**Document ID :** SNISID-BIO-MDL-001 | **Version :** 1.0.0

---

## 1. MODÈLES PYTHON (Pydantic v2)

```python
# bio_adn/models/dna.py
from pydantic import BaseModel, Field, validator, UUID4
from typing import Optional
from enum import Enum
from datetime import date, datetime

class IndexType(str, Enum):
    CONVICTED    = "BIO-CON"
    ARRESTEE     = "BIO-ARR"
    FORENSIC     = "BIO-FSC"
    MISSING      = "BIO-DIS"
    UNIDENTIFIED = "BIO-RNI"

class STRLocus(BaseModel):
    locus:  str
    value1: float
    value2: Optional[float] = None

class STRLociData(BaseModel):
    """20 loci CODIS Core"""
    CSF1PO:   STRLocus
    D3S1358:  STRLocus
    D5S818:   STRLocus
    D7S820:   STRLocus
    D8S1179:  STRLocus
    D13S317:  STRLocus
    D16S539:  STRLocus
    D18S51:   STRLocus
    D21S11:   STRLocus
    FGA:      STRLocus
    TH01:     STRLocus
    TPOX:     STRLocus
    vWA:      STRLocus
    D1S1656:  STRLocus
    D2S441:   STRLocus
    D2S1338:  STRLocus
    D10S1248: STRLocus
    D12S391:  STRLocus
    D19S433:  STRLocus
    D22S1045: STRLocus

class STRProfileCreate(BaseModel):
    specimen_number: str = Field(..., min_length=5)
    index_type:      IndexType
    loci_data:       STRLociData
    amelogenin:      Optional[str] = Field(None, pattern="^(XX|XY)$")
    quality_score:   float = Field(..., ge=0, le=1)
    lab_id:          UUID4
    case_number:     Optional[str] = None
    collected_date:  date
    correlation_id:  str

class STRProfileDB(STRProfileCreate):
    sample_id:       UUID4
    loci_encrypted:  bytes          # Stocké chiffré en base
    loci_hash:       str            # SHA-256
    uploaded_ldis:   bool = False
    uploaded_sdis:   bool = False
    uploaded_ndis:   bool = False
    is_expunged:     bool = False
    created_at:      datetime
    updated_at:      datetime
```

```python
# bio_adn/models/persons.py
from pydantic import BaseModel, Field
from typing import Optional, List
from enum import Enum
from datetime import date, datetime
from uuid import UUID

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

class RecordStatus(str, Enum):
    ACTIVE    = "ACTIVE"
    CLEARED   = "CLEARED"
    EXPIRED   = "EXPIRED"
    SUSPENDED = "SUSPENDED"

class WantedPersonCreate(BaseModel):
    niu:               Optional[str] = None
    last_name:         Optional[str] = None
    first_name:        Optional[str] = None
    aliases:           List[str] = []
    date_of_birth:     Optional[date] = None
    gender:            Optional[str] = Field(None, pattern="^[MFU]$")
    nationality:       Optional[str] = Field(None, min_length=3, max_length=3)
    warrant_type:      WarrantType
    warrant_number:    Optional[str] = None
    issuing_court:     Optional[str] = None
    issuing_date:      date
    charges:           List[str] = Field(..., min_items=1)
    danger_level:      DangerLevel = DangerLevel.MEDIUM
    armed_dangerous:   bool = False
    height_cm:         Optional[int] = Field(None, ge=50, le=250)
    weight_kg:         Optional[int] = Field(None, ge=10, le=300)
    mco_contact:       str  # Obligatoire
    expiry_date:       Optional[date] = None
    interpol_notice:   Optional[str] = None

    @validator("last_name", "first_name", always=True)
    def at_least_one_name(cls, v, values):
        if not v and not values.get("niu"):
            raise ValueError("last_name, first_name ou niu requis")
        return v

class MissingPersonCreate(BaseModel):
    niu:               Optional[str] = None
    last_name:         str
    first_name:        str
    date_of_birth:     Optional[date] = None
    age_at_missing:    Optional[int] = None
    gender:            Optional[str] = None
    category:          str = Field(..., pattern="^(CHILD|ENDANGERED|INVOLUNTARY|CATASTROPHE|OTHER)$")
    missing_date:      datetime
    missing_location:  str
    circumstances:     Optional[str] = None
    height_cm:         Optional[int] = None
    weight_kg:         Optional[int] = None
    family_contact:    Optional[str] = None
    family_phone:      Optional[str] = None
    medical_conditions: Optional[str] = None
    entering_agency:   str
```

```python
# bio_adn/models/property.py
from pydantic import BaseModel, Field, validator
from typing import Optional
from datetime import date
from enum import Enum

class StolenVehicleCreate(BaseModel):
    vin:               Optional[str] = Field(None, regex="^[A-HJ-NPR-Z0-9]{17}$")
    plate_number:      str
    plate_dept:        Optional[str] = None
    vehicle_make:      str
    vehicle_model:     str
    vehicle_year:      int = Field(..., ge=1950, le=2030)
    vehicle_color:     str
    vehicle_type:      Optional[str] = None
    theft_date:        date
    theft_location:    str
    theft_department:  Optional[str] = None
    owner_niu:         Optional[str] = None
    owner_name:        Optional[str] = None
    owner_phone:       Optional[str] = None
    entering_agency:   str

class StolenVesselCreate(BaseModel):
    vessel_name:       Optional[str] = None
    registration_number: Optional[str] = None
    hull_id_number:    Optional[str] = None   # HIN
    vessel_type:       str = Field(..., pattern="^(FISHING_CANOE|MOTORBOAT|SAILBOAT|FERRY|CARGO_SMALL|PATROL_BOAT|OTHER)$")
    vessel_make:       Optional[str] = None
    vessel_length_m:   Optional[float] = None
    hull_color:        Optional[str] = None
    home_port:         Optional[str] = None
    theft_location:    str
    theft_date:        date
    owner_niu:         Optional[str] = None
    owner_name:        Optional[str] = None
    entering_agency:   str
```

---

## 2. RELATIONS ENTRE MODÈLES

```
STRProfileDB
    │
    ├── bio_identity_links (1:1, accès restreint DCPJ Director)
    │       └── niu → SNISID Core (NIU)
    │
    └── bio_hits (1:N)
            ├── query_sample_id
            └── match_sample_id

WantedPersonDB
    ├── niu → SNISID Core (optionnel)
    ├── bio_sample_ref → STRProfileDB (optionnel)
    └── fingerprint_ref → NGI-HT (optionnel)

MissingPersonDB
    ├── niu → SNISID Core (optionnel)
    └── bio_sample_ref → STRProfileDB BIO-DIS (optionnel)

StolenVehicleDB
    └── foves_record_id → FOVeS/SIV (MP-15, bidirectionnel)
```

---

## 3. CONSTANTES ET ENUMS PARTAGÉS

```python
# bio_adn/constants.py

# Départements haïtiens
HAITI_DEPARTMENTS = [
    "OUEST", "NORD", "NORD-EST", "NORD-OUEST",
    "ARTIBONITE", "CENTRE", "SUD", "SUD-EST",
    "GRAND-ANSE", "NIPPES"
]

# Ports d'attache marins (BIE-EMB)
HAITI_PORTS = [
    "PORT-AU-PRINCE", "CAP-HAÏTIEN", "JACMEL",
    "LES-CAYES", "JÉRÉMIE", "SAINT-MARC", "GONAÏVES",
    "PORT-DE-PAIX", "FORT-LIBERTÉ", "MIRAGOÂNE"
]

# Niveaux d'alerte
ALERT_LEVELS = {
    "FULL_MATCH": "CRITICAL",
    "PARTIAL": "HIGH",
    "FAMILIAL": "MEDIUM",
    "VEHICLE_STOLEN": "HIGH",
    "PERSON_WANTED": "HIGH",
    "PERSON_WANTED_ARMED": "CRITICAL",
}

# Seuils qualité par index
QUALITY_THRESHOLDS = {
    "BIO-CON": {"min_score": 0.95, "min_loci": 20},
    "BIO-ARR": {"min_score": 0.90, "min_loci": 18},
    "BIO-FSC": {"min_score": 0.60, "min_loci": 10},
    "BIO-DIS": {"min_score": 0.85, "min_loci": 15},
    "BIO-RNI": {"min_score": 0.50, "min_loci": 8},
}

# Durées de rétention (en jours)
RETENTION_DAYS = {
    "BIO-ARR": 3 * 365,          # 3 ans max
    "BIO-CON": None,             # Durée condamnation + 10 ans (variable)
    "PER-REC": None,             # Durée mandat + 2 ans (variable)
    "BIE-VEH": 5 * 365,         # 5 ans après récupération
}
```
