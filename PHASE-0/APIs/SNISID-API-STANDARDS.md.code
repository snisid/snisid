# 🔌 SNISID — API STANDARDS & CONTRACTS
## Référentiel National des Contrats API

**Document ID :** SNISID-API-001  
**Version :** 1.0.0  
**Date :** Mai 2026  
**Propriétaire :** Bureau de Gouvernance de l'Interopérabilité (BGI)  
**Classification :** Usage Gouvernemental  

---

## 1. PRINCIPES FONDAMENTAUX DES APIs SNISID

### 1.1 Règle d'Or API
> Toute fonctionnalité du SNISID est exposée via API documentée en OpenAPI 3.1. Il n'existe pas de backdoor, de connexion directe à la base de données, ni d'interface non documentée.

### 1.2 Principes API-First
1. **API avant UI** — L'interface est construite sur les APIs, pas l'inverse
2. **Documentation-driven** — L'OpenAPI Spec est le contrat, pas le code
3. **Versioning strict** — Toute API versionnée, pas de breaking changes sans préavis 6 mois
4. **Security by design** — OAuth 2.1 + mTLS mandatory sur toutes les APIs
5. **Idempotence** — Toutes les opérations write sont idempotentes (header Idempotency-Key)
6. **Auditabilité** — Toute requête API loguée de façon immuable

---

## 2. STANDARDS DE CONCEPTION API

### 2.1 Structure URL
```
https://api.snisid.gouv.ht/v{n}/{domain}/{resource}

Exemples:
  GET    /v1/identity/citizens/{nin}
  POST   /v1/civil/births
  PUT    /v1/identity/citizens/{nin}/status
  DELETE (interdit sur données identité — soft delete uniquement)
  GET    /v1/biometric/matches/{match-id}
  POST   /v1/auth/tokens
  GET    /v1/agencies/oni/citizens/{nin}  ← inter-agency via bus
```

### 2.2 Conventions REST
```yaml
méthodes_autorisées: [GET, POST, PUT, PATCH, DELETE (restricted)]
codes_réponse:
  200: OK (GET, PUT, PATCH successful)
  201: Created (POST successful)
  202: Accepted (async operation)
  204: No Content (DELETE)
  400: Bad Request (validation error)
  401: Unauthorized (missing/invalid token)
  403: Forbidden (insufficient permissions - OPA denied)
  404: Not Found
  409: Conflict (duplicate, optimistic lock)
  422: Unprocessable Entity (business rule violation)
  429: Too Many Requests (rate limit)
  500: Internal Server Error
  503: Service Unavailable (circuit breaker open)

pagination: cursor-based
  params: cursor, limit (max 100)
  response: { data: [], next_cursor: "...", total_count: n }

filtering: query params standardisés
  ?status=active&created_after=2026-01-01&sort=created_at:desc

content_type: application/json; charset=UTF-8
accept: application/json (primary), application/cbor (edge)
```

### 2.3 Format Standard de Réponse
```json
{
  "data": { ... },
  "meta": {
    "request_id": "uuid-v4",
    "timestamp": "ISO-8601",
    "api_version": "1.0.0",
    "processing_time_ms": 45
  },
  "pagination": {
    "next_cursor": "...",
    "has_more": true,
    "total_count": 1234
  }
}
```

### 2.4 Format Standard d'Erreur (RFC 7807)
```json
{
  "type": "https://api.snisid.gouv.ht/errors/citizen-not-found",
  "title": "Citizen Not Found",
  "status": 404,
  "detail": "No citizen found with NIN: 1234567890123",
  "instance": "/v1/identity/citizens/1234567890123",
  "trace_id": "uuid-v4",
  "timestamp": "2026-05-25T14:32:00Z"
}
```

---

## 3. SÉCURITÉ DES APIs

### 3.1 Authentification API (OAuth 2.1 + mTLS)
```
Flux Autorisés:
  Client Credentials (machine-to-machine agences)
  Authorization Code + PKCE (citoyens via portail)
  Device Flow (offline kits terrain)

Token:
  Type: JWT (signed RS256 ou ES256)
  TTL: 900 secondes (15 minutes) — non négociable
  Claims obligatoires:
    sub: {nin ou service-account-id}
    iss: https://auth.snisid.gouv.ht
    aud: https://api.snisid.gouv.ht
    exp, iat, jti (unique ID anti-replay)
    scope: {liste des scopes autorisés}
    agency_code: {code agence}
    commune_code: {code commune pour restriction géo}

mTLS: Obligatoire pour toutes les agences (Machine-to-Machine)
  Certificate: Émis par AN-PKI (Issuing CA Agences)
  Validation: via service mesh Istio
```

### 3.2 Scopes OAuth Définis
```yaml
scopes:
  identity:read: "Lire les données d'identité civile basiques"
  identity:verify: "Vérifier l'identité (1:1 match)"
  identity:enroll: "Créer/mettre à jour des enrôlements"
  identity:admin: "Administration identités (révocation, correction)"
  civil:read: "Lire les événements état civil"
  civil:write: "Créer des actes état civil"
  biometric:match: "Soumettre des requêtes de matching biométrique"
  biometric:admin: "Administration biométrique (audit uniquement)"
  audit:read: "Accès aux logs d'audit (inspecteurs)"
  kyc:verify: "KYC tiers (banques, télécom) - consentement requis"
  agency:admin: "Administration des agences"
  soc:monitor: "Monitoring SOC temps réel"
```

### 3.3 Rate Limiting par Agence
```yaml
tiers:
  TIER_CRITIQUE (ONI, MJSP, DIE):
    requests_per_second: 1000
    burst: 5000
    concurrent: 200

  TIER_STANDARD (DGI, MSPP, MENFP, DCPJ):
    requests_per_second: 200
    burst: 1000
    concurrent: 50

  TIER_BASIQUE (ONA, OFATMA, CEP, INARA):
    requests_per_second: 50
    burst: 200
    concurrent: 20

  TIER_PARTENAIRE_PRIVE (Banques, Télécom):
    requests_per_second: 10
    burst: 50
    concurrent: 5
    requires_citizen_consent: true
```

---

## 4. CONTRATS API PRINCIPAUX

### 4.1 API Identity Core

```yaml
openapi: 3.1.0
info:
  title: SNISID Identity Core API
  version: 1.0.0
  description: API principale d'identité nationale

paths:
  /v1/identity/citizens/{nin}:
    get:
      summary: Récupérer l'identité d'un citoyen
      security: [OAuth2: [identity:read], mTLS: []]
      parameters:
        - nin: string (13 chiffres + checksum)
        - fields: array<string> (sélection champs — minimisation)
      responses:
        200:
          schema:
            nin: string
            status: enum [active, deceased, suspended, revoked]
            names:
              surname: string
              given_names: array<string>
            birth_date: date (YYYY-MM-DD)
            birth_place: { code: string, name: string }
            sex: enum [M, F, X]
            nationality: array<string> (ISO 3166-1 alpha-3)
            _links:
              self: "/v1/identity/citizens/{nin}"
              documents: "/v1/identity/citizens/{nin}/documents"
              civil_events: "/v1/identity/citizens/{nin}/events"

  /v1/identity/verify:
    post:
      summary: Vérification 1:1 d'identité (avec NIN + biométrie)
      security: [OAuth2: [identity:verify], mTLS: []]
      request:
        nin: string
        verification_mode: enum [card_signature, biometric_1to1, document]
        biometric_probe: bytes? (template ISO 19794)
      response:
        verified: boolean
        confidence_score: float (0.0-1.0)
        match_type: enum [exact, probabilistic, failed]
        verification_id: uuid
```

### 4.2 API État Civil

```yaml
paths:
  /v1/civil/births:
    post:
      summary: Déclarer une naissance (EC-N01 à EC-N05)
      security: [OAuth2: [civil:write], mTLS: []]
      request:
        procedure_type: enum [simple, recognition, late, decree, judgment]
        child:
          surname: string
          given_names: array<string>
          sex: enum [M, F, X]
          birth_date: date
          birth_place: AdministrativeUnitCode
          birth_time: time?
        mother: { nin: string }
        father: { nin: string }?
        witnesses: array<{ nin: string, role: string }>
        officer_id: string
        medical_certificate_id: string? (FHIR ID si dispo)
        offline_mode: boolean
        offline_signature: bytes? (signature HSM terrain)
      response:
        event_id: uuid
        nin_child: string? (attribué si online, null si offline)
        status: enum [registered, pending_sync, pending_validation]
        document_url: string? (PDF/A-3 signé)
        receipt_code: string (pour suivi offline)

  /v1/civil/births/{event_id}:
    get:
      summary: Récupérer un acte de naissance
      security: [OAuth2: [civil:read], mTLS: []]
      response:
        event: CivilEvent
        document: { url: string, hash_sha256: string, signature: XAdESLTA }
```

### 4.3 API KYC (Tiers privés — Consentement requis)

```yaml
paths:
  /v1/kyc/verify:
    post:
      summary: Vérification KYC légère (banques, télécoms)
      security: [OAuth2: [kyc:verify], mTLS: []]
      note: "Requiert consent_token valide du citoyen"
      request:
        nin: string
        consent_token: string (JWT signé par clé privée citoyen)
        verification_level: enum [basic, standard, enhanced]
        purpose: string (ex: "Ouverture compte bancaire BNC")
      response:
        verified: boolean
        # Données minimales selon niveau + consentement
        basic:
          name_match: boolean
          is_alive: boolean
          is_adult: boolean
        standard: + {
          full_name: string
          birth_date: date
          nationality: string
        }
        enhanced: + {
          address: Address (si consentement accordé)
        }
        # JAMAIS: biométrie, NIF, données judiciaires
        consent_used_id: uuid
        audit_reference: uuid
```

### 4.4 API Événements (Kafka — CloudEvents)

```yaml
# Schema événement naissance
event_schema: etat-civil.naissance.declaree.v1
type: object
properties:
  specversion: "1.0"
  type: "ht.gov.snisid.etat-civil.naissance.declaree.v1"
  source: "https://api.snisid.gouv.ht/v1/civil/births"
  id: uuid
  time: datetime
  datacontenttype: "application/json"
  data:
    event_id: uuid
    nin_child: string?
    procedure_type: string
    birth_date: date
    birth_place_code: string
    officer_id: string
    timestamp_registered: datetime
    # PAS de données personnelles sensibles dans l'événement bus
    # Les agences consommatrices requêtent l'API pour les détails
```

---

## 5. CATALOGUE D'APIs PAR AGENCE

### 5.1 ONI — Office National d'Identification
| API | Endpoint | Consommateurs |
|-----|----------|--------------|
| Identity CRUD | /v1/identity/citizens | MJSP, DIE, DGI, DCPJ |
| Biometric Verify | /v1/identity/verify | Toutes agences autorisées |
| Enrollment Status | /v1/identity/enrollments/{id} | Agents ONI |
| NIN Generation | /v1/identity/nin/generate | Interne ONI uniquement |

### 5.2 MJSP/DGEC — Ministère Justice / État Civil
| API | Endpoint | Consommateurs |
|-----|----------|--------------|
| Civil Births | /v1/civil/births | ONI, MSPP, CEP |
| Civil Deaths | /v1/civil/deaths | ONI, DGI, OFATMA, Banques |
| Civil Marriages | /v1/civil/marriages | ONI, Notaires |
| Criminal Record | /v1/justice/criminal-records/{nin} | DCPJ, DIE, DGI |

### 5.3 CNIGS — Centre National Information Géo-Spatiale
| API | Endpoint | Consommateurs |
|-----|----------|--------------|
| Administrative Units | /v1/geo/units | Toutes agences |
| GeoSearch | /v1/geo/search?query=... | Toutes agences |
| GeoJSON Export | /v1/geo/units/{code}/geometry | SIG, ONI |

### 5.4 DGI — Direction Générale Impôts
| API | Endpoint | Consommateurs |
|-----|----------|--------------|
| Tax Status | /v1/tax/citizens/{nin}/status | Banques (avec consent) |
| NIF Verification | /v1/tax/nif/{nif}/verify | Banques, AGD |

### 5.5 PNH/DCPJ — Police
| API | Endpoint | Consommateurs |
|-----|----------|--------------|
| Wanted Persons | /v1/police/wanted/{nin} | DIE (frontières) |
| Biometric Search | /v1/police/biometric-search | DCPJ interne |

### 5.6 DIE — Immigration
| API | Endpoint | Consommateurs |
|-----|----------|--------------|
| Border Crossing | /v1/immigration/crossings/{nin} | Douanes, Sécurité |
| Passport Status | /v1/immigration/passports/{nin} | DIE interne |
| Visa Query | /v1/immigration/visas/{visa_id} | Ambassades |

---

## 6. GOUVERNANCE DU CATALOGUE API

### 6.1 Cycle de Vie API

```
1. PROPOSE    → Agence soumet demande API (template BGI)
2. DESIGN     → Architecte BGI + équipe agence conçoivent spec OpenAPI
3. REVIEW     → Comité BGI + CISO + juriste (privacy) review
4. APPROVE    → CCB approuve (4-yeux)
5. PILOT      → Déploiement environnement pilote (2 agences max)
6. PUBLISH    → Publication catalogue + documentation développeurs
7. MONITOR    → SLA + usage monitoring continu
8. DEPRECATE  → Notification 6 mois minimum + migration plan
9. RETIRE     → Suppression après migration confirmée
```

### 6.2 Portail Développeurs
- URL : `https://dev.snisid.gouv.ht`
- Contenu : catalogue OpenAPI interactif (Swagger UI + Redoc)
- Accès : intranet État pour APIs sensibles, internet pour APIs publiques
- Sandbox : environnement de test avec données anonymisées
- SDK : Python, Java, Go (générés depuis OpenAPI)

### 6.3 SLA API Standard

| Classe API | Disponibilité | Latence p95 | Débit max |
|-----------|--------------|-------------|-----------|
| Core Identity | 99,99% | 200ms | 1000 TPS |
| Civil Registry | 99,95% | 500ms | 200 TPS |
| Biometric Match (1:1) | 99,5% | 500ms | 100 TPS |
| Biometric Search (1:N) | 99% | 3000ms | 10 TPS |
| Analytics | 99% | 2000ms | 50 TPS |
| KYC (tiers) | 99,9% | 300ms | 50 TPS/agence |

---

## 7. TESTS & VALIDATION

### 7.1 Contract Testing (Pact)
```bash
# Chaque agence consommatrice maintient des tests Pact
# Pipeline CI vérifie la compatibilité avant tout déploiement
# Broker Pact central: https://pact.snisid.gouv.ht (interne)
```

### 7.2 API Health Checks
```
GET /health         → 200 { status: "up", version: "1.0.0" }
GET /health/live    → 200/503 (Kubernetes liveness)
GET /health/ready   → 200/503 (Kubernetes readiness)
GET /metrics        → Prometheus format (interne uniquement)
```

---

*Document approuvé par le Bureau de Gouvernance de l'Interopérabilité (BGI)*  
*SNISID — République d'Haïti*
