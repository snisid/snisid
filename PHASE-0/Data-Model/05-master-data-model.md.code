# 🗂️ SNISID — National Master Data Model

**Document N° :** SNISID-DAT-005
**Étape Phase 0 :** 5/16
**Principe :** *Les données doivent être unifiées nationalement.*

---

## 1. Objectif

Définir le **modèle de données canonique** partagé par toutes les agences haïtiennes utilisatrices de SNISID. Une seule vérité — la **golden record** — pour chaque citoyen, événement civil, document, transaction.

---

## 2. Domaines Maîtres

| Domaine | Entité racine | Référentiel responsable |
|---------|---------------|--------------------------|
| Citoyens | `Person` | ONI |
| Biométrie | `BiometricRecord` | ONI / Police (AFIS) |
| État civil | `CivilEvent` (Naissance, Mariage, Décès, Divorce, Adoption) | Officiers État Civil / Ministère Justice |
| Justice | `JudicialRecord` | Ministère de la Justice |
| Police | `PoliceRecord` | DCPJ / PNH |
| Immigration | `TravelRecord`, `Visa`, `Passport` | DIE + MAE |
| Géographie | `AdministrativeUnit` (Dépt, Commune, Section) | CNIGS |
| Documents | `OfficialDocument` | Tous |

---

## 3. Identifiants Nationaux

| ID | Format | Émetteur | Usage |
|----|--------|----------|-------|
| **NIN** (National Identity Number) | 13 chiffres + checksum Luhn | ONI | Identifiant unique citoyen à vie |
| **NIF** (Num Identification Fiscale) | DGI | Fiscal |
| **NIS** (Num Identification Sécurité) | OFATMA/ONA | Social |
| **DocumentID** | UUID v4 | SNISID | Tout document officiel |
| **EventID** | UUID v4 | SNISID | Tout événement civil |

Règles :
- NIN attribué à la naissance (ou à l'enrôlement pour les non-déclarés)
- NIN immuable, ne contient **aucune information sensible** (pas de date de naissance, pas de sexe)
- Mapping NIN ↔ autres identifiants via service d'identité fédérée

---

## 4. Entité `Person` (Schéma logique)

```yaml
Person:
  nin: string (13 + checksum)            # PK national
  status: enum [Active, Deceased, Suspended, Revoked]
  names:
    surname: string
    given_names: array<string>
    other_names: array<string>           # surnoms, alias légaux
  sex: enum [M, F, X]                    # X : non binaire / indéterminé
  birth:
    date: date
    place: AdministrativeUnit
    register_ref: CivilEvent             # acte de naissance
  nationality: array<string>             # ISO 3166-1 alpha-3 (HTI par défaut)
  parents:
    father: Person?
    mother: Person?
  civil_status: enum [Single, Married, Divorced, Widowed]
  spouse: Person?
  biometrics: BiometricRecord
  contacts:
    phones: array<Phone>
    emails: array<Email>
    addresses: array<Address>
  consent_flags: ConsentSet
  audit:
    created_at, created_by, updated_at, updated_by
    version (event-sourcing)
```

---

## 5. Entité `BiometricRecord`

```yaml
BiometricRecord:
  id: UUID
  person_nin: ref Person
  fingerprints:
    template_iso_19794_2: bytes (chiffré)
    quality_nfiq2: array<int>             # 10 doigts
  face:
    image_iso_19794_5: bytes (chiffré)
    template: bytes
  iris:
    left_template: bytes
    right_template: bytes
  capture_metadata:
    device_id, operator_id, location, timestamp
  consent_id: ref Consent
  retention_policy: ref Policy
```

> Standards : **ISO/IEC 19794** (templates), **NFIQ 2.0** (qualité empreintes), **ICAO 9303** (photo faciale passeport).

---

## 6. Entité `CivilEvent`

```yaml
CivilEvent:
  id: UUID
  type: enum [Birth, Marriage, Divorce, Death, Adoption, Recognition]
  subtype: enum [Simple, Recognition, LateDeclaration, Decree, JudgmentMinutes]
  date_event: date
  date_registration: date
  location: AdministrativeUnit
  officer: PublicOfficer
  parties: array<Person>
  witnesses: array<Person>
  document_id: ref OfficialDocument
  legal_basis: string                     # référence loi/article
  status: enum [Draft, Validated, Sealed, Annulled]
```

---

## 7. Entité `OfficialDocument`

```yaml
OfficialDocument:
  id: UUID
  type: enum [BirthCertificate, NationalIDCard, Passport, Judgment, ...]
  related_person: ref Person
  related_event: ref CivilEvent?
  issued_by: Agency
  issued_at: timestamp
  valid_until: date?
  signature: DigitalSignature             # XAdES-LTA
  qr_code: string                          # vérification publique
  pdf_a3_hash: sha-256
```

---

## 8. Référentiel Géographique (CNIGS)

```yaml
AdministrativeUnit:
  code: string                            # code officiel CNIGS
  level: enum [Country, Department, Arrondissement, Commune, CommunalSection, Locality]
  name_fr: string
  name_ht: string                         # créole
  parent: ref AdministrativeUnit
  geometry: GeoJSON
```

> Source de vérité unique : **CNIGS**, publiée comme **bien public** au format CSV + GeoJSON + API REST.

---

## 9. Gouvernance des Données Maîtres (MDM)

- **Data Steward** désigné par domaine (ex. ONI pour Person, MJ pour CivilEvent)
- **Comité MDM** mensuel — résolution doublons, arbitrages
- **Golden Record** maintenue par moteur de matching (probabilistic + déterministe)
- **Lineage** tracé bout-en-bout (Apache Atlas ou équivalent)

---

## 10. Qualité des Données — KPI

| KPI | Cible 2027 | Cible 2030 |
|-----|------------|------------|
| Complétude (champs obligatoires renseignés) | ≥ 95 % | ≥ 99 % |
| Doublons biométriques résiduels | ≤ 1 % | ≤ 0,1 % |
| Cohérence inter-agences | ≥ 90 % | ≥ 98 % |
| Délai résolution incohérence | ≤ 30 j | ≤ 7 j |

---

## 11. Politique de Rétention

| Donnée | Rétention | Justification |
|--------|-----------|---------------|
| Person (vivant) | Permanente | Identité légale |
| Person (décédé) | 100 ans | Historique civil |
| Biométriques (vivant) | Permanente, chiffrée | Authentification |
| Biométriques (décédé) | 10 ans puis destruction | Privacy |
| Logs accès | 5 ans | Audit, forensics |
| Documents officiels | Permanente | Valeur légale |
| Données KYC tiers | 5 ans après fin relation | Loi anti-blanchiment |

---

## 12. Interopérabilité Internationale

Compatibilité maintenue avec :
- **ICAO 9303** (passeport biométrique)
- **OACI MRTD**
- **ISO/IEC 24727** (cartes à puce)
- **FHIR** (santé, pour interconnexion MSPP)
- **OpenCRVS** (modèle de données état civil — référence)
- **MOSIP** (modular open-source identity platform — référence d'inspiration)

---
*Fin du document — Étape 5/16*
