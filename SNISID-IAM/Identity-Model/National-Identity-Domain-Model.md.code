# National Identity Domain Model

> **Document officiel du domaine d'identité nationale SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. VISION

L'identité nationale est la **source unique de vérité** pour tous les systèmes SNISID. Chaque entité (citoyen, agent, système) possède une identité :

- **Unique** — Un seul enregistrement par individu
- **Immutable** — Le noyau identitaire ne change jamais
- **Auditée** — Chaque modification est tracée
- **Versionnée** — Historique complet des changements

---

## 2. DOMAINES DU MODÈLE

### 2.1 Citizen — Identité Individuelle

```
┌─────────────────────────────────────────────────────────┐
│                    CITIZEN IDENTITY                      │
├─────────────────────────────────────────────────────────┤
│ nnu:                 UUID souverain national             │
│ given_name:          String (UTF-8)                      │
│ family_name:         String (UTF-8)                      │
│ date_of_birth:       Date (ISO 8601)                     │
│ place_of_birth:      GeoPoint (ISO 6709)                 │
│ sex:                 M / F / Non-binary                  │
│ nationality:         HT (Haïti)                          │
│ mother_maiden_name:  String (pour vérification)          │
│ status:              active | suspended | revoked         │
│ created_at:          Timestamp (UTC)                     │
│ updated_at:          Timestamp (UTC)                     │
│ version:             Integer (auto-incrément)            │
│ audit_trail_id:      UUID → Audit Log                    │
│ biometric_refs:      Array[UUID] → Biometrics Domain     │
│ credential_refs:     Array[UUID] → Credentials Domain    │
│ role_refs:           Array[UUID] → Roles Domain          │
│ consent_records:     Array[UUID] → Consent Domain        │
│ judicial_status_ref: UUID → Judicial Status Domain       │
└─────────────────────────────────────────────────────────┘
```

#### Contraintes

| Champ | Règle |
|-------|-------|
| nnu | Unique, généré à l'enrôlement, jamais recyclé |
| given_name | 2-50 caractères, lettres + espaces |
| family_name | 2-50 caractères, lettres + espaces |
| date_of_birth | Date valide ≤ aujourd'hui |
| place_of_birth | Référentiel géographique officiel |
| status | Mutation via workflow signé uniquement |
| version | Incrémenté à chaque modification |

#### Cycle de vie

```
[Enrôlement] → [Créé] → [Vérifié] → [Actif]
                        ↓
                   [Suspendu] → [Révoqué]
                        ↓
                   [Corrigé] → [Actif]
                        ↓
                   [Contesté] → [Audité] → [Décision]
```

---

### 2.2 Legal Identity — État Légal

| Champ | Type | Description |
|-------|------|-------------|
| identity_id | UUID → Citizen | Référence identité citoyenne |
| birth_certificate | String | Numéro acte de naissance |
| birth_reg_date | Date | Date enregistrement |
| birth_reg_office | String | Bureau d'état civil |
| national_id_card | String | Numéro carte d'identité |
| passport_number | String | Passeport (optionnel) |
| tax_id | String | NIF |
| voter_id | String | Carte électorale (optionnel) |
| legal_status | Enum | valid \| expired \| revoked \| pending |
| provenance | Enum | original \| derived \| attested |
| attestation_level | Integer | 1-5 (5 = souverain) |
| legal_restrictions | Array[String] | Restrictions légales |
| verified_by | UUID | Enrollment Officer certifié |

#### Niveaux d'Attestation

| Niveau | Description | Usage |
|--------|-------------|-------|
| 1 | Déclaration unique | Accès basique |
| 2 | Document identité national | Services citoyens |
| 3 | Acte naissance officiel | Services financiers |
| 4 | Vérification biométrique | Services gouvernementaux |
| 5 | Double vérification + notaire | Accès souverain/critique |

---

### 2.3 Biometrics — Biométrie

| Champ | Type | Description |
|-------|------|-------------|
| biometric_id | UUID | Identifiant unique |
| citizen_ref | UUID → Citizen | Référence citoyen |
| template_type | Enum | fingerprint \| face \| iris \| voice |
| template_hash | String | SHA-256 du gabarit |
| template_encrypted | Binary | AES-256-GCM |
| encryption_key_ref | UUID → PKI | Référence clé HSM |
| quality_score | Float | 0.0 - 1.0 |
| capture_device | String | ID terminal certifié |
| capture_location | GeoPoint | Lieu de capture |
| capture_timestamp | Timestamp | UTC |
| capture_officer | UUID | Enrollment Officer |
| liveness_check | Boolean | Détection vivant |
| liveness_score | Float | 0.0 - 1.0 |
| status | Enum | active \| revoked \| expired |
| version | Integer | Versionnement |
| audit_trail_id | UUID | Trace d'audit |
| consent_ref | UUID | Référence consentement |

#### Règles de Protection

| Règle | Implémentation |
|-------|----------------|
| Chiffrement | AES-256-GCM au repos |
| Transmission | TLS 1.3 uniquement |
| Gabarit | Jamais stocké en clair |
| Accès | MFA + rôle autorisé |
| Rétention | Politique légale nationale |
| Suppression | Destruction sécurisée (DoD 5220.22-M) |

---

### 2.4 Credentials — Documents & Certificats

| Champ | Type | Description |
|-------|------|-------------|
| credential_id | UUID | Identifiant unique |
| citizen_ref | UUID → Citizen | Référence citoyen |
| credential_type | Enum | national_id \| passport \| driver \| voter_card \| digital_cert |
| issuer | String | Autorité émettrice |
| issue_date | Date | Date d'émission |
| expiry_date | Date | Date d'expiration |
| serial_number | String | Unique |
| digital_signature | String | Signature PKI |
| status | Enum | valid \| expired \| revoked \| stolen |
| verification_hash | String | SHA-256 |
| wallet_ref | UUID | Référence wallet |
| version | Integer | Versionnement |
| audit_trail_id | UUID | Trace d'audit |

---

### 2.5 Roles — Fonctions & Permissions

| Champ | Type | Description |
|-------|------|-------------|
| role_id | UUID | Identifiant unique |
| role_name | String | Nom unique du rôle |
| role_level | Enum | citizen \| medium \| high \| critical |
| description | Text | Description du rôle |
| permissions | Array[Permission] | Droits associés |
| inherits_from | Array[UUID] | Rôles parents |
| requires_mfa | Boolean | MFA obligatoire |
| requires_abac | Boolean | ABAC obligatoire |
| max_session_duration | Integer | Minutes |
| temporal_constraint | CRON | Contrainte horaire |
| status | Enum | active \| deprecated \| revoked |
| created_by | UUID → Admin | Créateur |
| created_at | Timestamp | Date création |
| version | Integer | Versionnement |
| audit_trail_id | UUID | Trace d'audit |

#### Permission

| Champ | Type | Description |
|-------|------|-------------|
| permission_id | UUID | Identifiant unique |
| resource | String | URI de la ressource |
| action | Enum | read \| write \| delete \| admin |
| condition | JSON | Règle ABAC |
| scope | Enum | national \| regional \| local |
| granted_by | UUID → Role | Rôle grant |
| effective_from | Timestamp | Début validité |
| effective_to | Timestamp | Fin validité (optionnel) |
| version | Integer | Versionnement |
| audit_trail_id | UUID | Trace d'audit |

---

### 2.6 Access — Journal d'Accès

| Champ | Type | Description |
|-------|------|-------------|
| access_id | UUID | Identifiant unique |
| citizen_ref | UUID | Citoyen concerné |
| actor_ref | UUID | Qui accède |
| resource_ref | String | Quoi |
| action | Enum | read \| write \| delete \| update |
| access_method | Enum | web \| api \| wallet \| kiosk |
| ip_address | String | Adresse IP |
| device_fingerprint | String | Empreinte appareil |
| location | GeoPoint | Géolocalisation |
| timestamp | Timestamp | UTC |
| mfa_used | Boolean | MFA utilisé |
| abac_decision | Enum | permit \| deny \| conditional |
| risk_score | Float | 0.0 - 1.0 |
| consent_ref | UUID | Référence consentement |
| status | Enum | success \| failed \| blocked |
| session_id | UUID | Session |
| audit_trail_id | UUID | Trace d'audit |

---

### 2.7 Consent — Autorisations Citoyen

| Champ | Type | Description |
|-------|------|-------------|
| consent_id | UUID | Identifiant unique |
| citizen_ref | UUID | Citoyen concerné |
| requester_ref | UUID | Entité demanderesse |
| purpose | String | Finalité explicite |
| data_scope | Array[String] | Données concernées |
| granted | Boolean | Accordé |
| granted_at | Timestamp | Date accord |
| granted_method | Enum | wallet \| in_person \| digital |
| valid_from | Timestamp | Début validité |
| valid_until | Timestamp | Fin validité |
| revoked | Boolean | Révoqué |
| revoked_at | Timestamp | Date révocation |
| revoked_reason | String | Raison révocation |
| version | Integer | Versionnement |
| audit_trail_id | UUID | Trace d'audit |

---

### 2.8 Judicial Status — Restrictions Judiciaires

| Champ | Type | Description |
|-------|------|-------------|
| judicial_id | UUID | Identifiant unique |
| citizen_ref | UUID | Citoyen concerné |
| restriction_type | Enum | travel_ban \| asset_freeze \| identity_suspension \| monitoring |
| court_reference | String | Numéro dossier |
| court_name | String | Nom tribunal |
| judge_ref | UUID | Magistrat |
| ordered_date | Date | Date ordonnance |
| effective_from | Timestamp | Début effet |
| effective_until | Timestamp | Fin effet |
| reason | Text | Motivation |
| status | Enum | active \| lifted \| expired |
| authorized_by | UUID | Signature PKI |
| version | Integer | Versionnement |
| audit_trail_id | UUID | Trace d'audit |

---

## 3. RELATIONS ENTRE DOMAINES

```
                    ┌──────────────┐
                    │   CITIZEN    │
                    │   (Noyau)    │
                    └──────┬───────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
           ▼               ▼               ▼
    ┌────────────┐ ┌─────────────┐ ┌──────────────┐
    │  LEGAL     │ │ BIOMETRICS  │ │ CREDENTIALS  │
    │  IDENTITY  │ │ (Preuve)    │ │ (Documents)  │
    └─────┬──────┘ └──────┬──────┘ └──────┬───────┘
          │               │               │
          ▼               ▼               ▼
    ┌────────────┐ ┌─────────────┐ ┌──────────────┐
    │  CONSENT   │ │   ACCESS    │ │   JUDICIAL   │
    │  (Droits)  │ │   (Logs)    │ │   (Statut)   │
    └────────────┘ └─────────────┘ └──────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │    ROLES    │
                    │ (Accès)     │
                    └─────────────┘
```

---

## 4. RÈGLES DE GOUVERNANCE DU MODÈLE

### 4.1 Immutabilité

| Règle | Description |
|-------|-------------|
| Noyau immuable | NNU, nom, date/lieu de naissance ne changent jamais |
| Correction | Toute correction crée une nouvelle version, jamais de modification directe |
| Historique | Chaque version est préservée et signée |

### 4.2 Versionnage

```
Version 1.0.0 — Création initiale (enrôlement)
Version 1.1.0 — Correction mineure (nom corrigé)
Version 2.0.0 — Changement majeur (statut légal)
Version X.Y.Z — X=changement majeur, Y=mineur, Z=correctif
```

### 4.3 Audit Trail

Chaque enregistrement possède un `audit_trail_id` qui pointe vers :

| Champ | Type | Description |
|-------|------|-------------|
| audit_id | UUID | Identifiant unique |
| entity_ref | UUID | Entité modifiée |
| entity_type | Enum | citizen \| biometric \| consent \| ... |
| operation | Enum | create \| update \| delete \| read |
| old_value | JSON | Avant modification |
| new_value | JSON | Après modification |
| performed_by | UUID | Acteur |
| performed_at | Timestamp | UTC |
| signature | String | Signature numérique |
| justification | Text | Raison |
| ip_address | String | IP acteur |
| session_id | UUID | Session |

### 4.4 Règle d'Or

> **Dans SNISID : aucune opération sans identité vérifiée.**

Chaque action sur le système IAM nécessite :
1. Une identité NNU valide
2. Un rôle autorisé
3. Une décision ABAC favorable
4. Un audit trail signé
5. Un consentement citoyen (si données personnelles)

---

## 5. SCHÉMA DE VALIDATION

```yaml
citizen_validation:
  required_fields:
    - nnu
    - given_name
    - family_name
    - date_of_birth
    - place_of_birth
    - nationality
    - status
  
  biometric_required:
    - at_least_one_template
    - liveness_check_passed
    - quality_score > 0.85
  
  legal_required:
    - birth_certificate_or_equivalent
    - verification_by_certified_officer
    - attestation_level >= 3
  
  consent_required:
    - explicit_consent_recorded
    - purpose_declared
    - data_scope_defined
```

---

## 6. INDEX DES ENTITÉS

| Domaine | Clé primaire | Clé unique | Référence externe |
|---------|-------------|------------|-------------------|
| Citizen | citizen_id | nnu | — |
| Legal Identity | identity_id | (citizen_id, serial_number) | État civil |
| Biometrics | biometric_id | (citizen_id, template_type) | Capteurs |
| Credentials | credential_id | serial_number | Autorités émettrices |
| Roles | role_id | role_name | IAM |
| Permissions | permission_id | (role_id, resource, action) | — |
| Access | access_id | — | SIEM |
| Consent | consent_id | — | Wallet citoyen |
| Judicial Status | judicial_id | (citizen_id, court_reference) | Système judiciaire |
| Audit | audit_id | — | Archive nationale |

---

> **Ce document est la référence officielle du domaine d'identité nationale SNISID.**  
> **Toute modification doit être approuvée par le comité de gouvernance IAM.**
