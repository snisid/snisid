# 🔌 SNISID — National Interoperability Framework

**Document N° :** SNISID-INT-008
**Étape Phase 0 :** 8/16
**Principe :** *Toutes les agences doivent parler le même langage numérique.*

---

## 1. Vision

SNISID établit un **bus national d'interopérabilité gouvernementale (Government Service Bus — GSB)** sur lequel chaque agence publie ses services et événements, et consomme ceux des autres, selon des **standards ouverts opposables**.

---

## 2. Référence Internationale

- **X-Road** (Estonie) — modèle d'inspiration
- **European Interoperability Framework (EIF)**
- **GovStack** (initiative ONU/GIZ/ITU)
- **GA4GH / FHIR** (santé)
- **OASIS / W3C** standards

---

## 3. Agences Cibles (Cartographie)

| # | Agence | Données apportées | Données consommées |
|---|--------|-------------------|---------------------|
| 1 | **ONI** | Identité, biométrie | Événements état civil |
| 2 | **MJSP / DGEC** | État civil | Identité |
| 3 | **CNIGS** | Référentiel géographique | — |
| 4 | **PNH / DCPJ** | Casier, recherches | Identité, état civil |
| 5 | **DAP** | Détentions | Identité |
| 6 | **DIE / MAE** | Passeports, visas, frontières | Identité, état civil |
| 7 | **DGI** | Fiscalité (NIF) | Identité |
| 8 | **AGD** | Douanes | Identité, immigration |
| 9 | **MSPP** | Santé (FHIR) | Identité, état civil |
| 10 | **MENFP** | Éducation | Identité |
| 11 | **MAST / OFATMA / ONA** | Protection sociale | Identité, état civil |
| 12 | **CEP** | Liste électorale | Identité, état civil |
| 13 | **BRH / banques** | KYC | Identité, état civil |
| 14 | **CONATEL / opérateurs télécom** | SIM registration | Identité |
| 15 | **INARA** | Foncier | Identité |

---

## 4. Couches d'Interopérabilité (EIF)

```
┌────────────────────────────────────────────┐
│  Juridique     → Accords, conventions     │
├────────────────────────────────────────────┤
│  Organisationnel → Processus, RACI         │
├────────────────────────────────────────────┤
│  Sémantique    → Ontologies, vocabulaires  │
├────────────────────────────────────────────┤
│  Technique     → APIs, formats, protocoles │
└────────────────────────────────────────────┘
```

---

## 5. Standards Techniques Imposés

| Domaine | Standard |
|---------|----------|
| API REST | OpenAPI 3.1 |
| API GraphQL (option) | GraphQL spec June 2018+ |
| Évènements | CloudEvents 1.0 + Apache Avro / Protobuf |
| Sécurité API | OAuth 2.1 + OIDC + mTLS |
| Messagerie asynchrone | Kafka avec Schema Registry |
| Identité fédérée | SAML 2.0 + OIDC |
| Santé | HL7 FHIR R4 |
| Documents | PDF/A-3 + XAdES-LTA |
| Géographique | GeoJSON, WMS, WFS (OGC) |
| Cartes à puce | ISO/IEC 7816 + 24727 |
| Voyage | ICAO 9303 |

---

## 6. Architecture du Bus National

```
   Agence A ───┐                          ┌─── Agence D
                │                          │
   Agence B ──┐ │   ┌─────────────────┐    │ ┌── Agence E
              │ │   │  GOV SERVICE BUS │    │ │
              ├─┼──▶│  - Routing       │◀───┼─┤
              │ │   │  - Mediation     │    │ │
              │ │   │  - Authz         │    │ │
   Agence C ──┘ │   │  - Audit         │    │ └── Agence F
                │   └─────────────────┘    │
                │           │              │
                │    ┌──────▼──────┐       │
                │    │ EVENT BUS   │       │
                └────│ (Kafka)     │───────┘
                     └─────────────┘
```

Composants :
- **API Gateway national** (Kong/APISIX) en front
- **Service Mesh** (Istio) pour mTLS interne
- **Kafka cluster** événementiel
- **Schema Registry** pour gouvernance des contrats
- **Catalog API** (Backstage ou DevPortal) pour discovery
- **Audit immuable** (write-only log + WORM storage)

---

## 7. Data Exchange Patterns

| Pattern | Usage |
|---------|-------|
| **Sync REST** | Vérification identité 1:1 (latence courte) |
| **Async Event** | Notification événement civil (naissance, décès) |
| **Bulk file** | Référentiels (CNIGS), liste électorale |
| **Streaming** | Logs SIEM, télémétrie |
| **Federated query** | Recherche multi-bases avec consentement |

---

## 8. Identity Federation

- **IdP central** : Keycloak SNISID
- Agences = **brokers** SAML/OIDC fédérés
- Citoyens utilisent leur identité unique pour TOUS services publics
- Niveaux d'assurance (LoA) selon eIDAS : faible / substantiel / élevé
- Consentement explicite pour chaque partage de donnée

---

## 9. Gouvernance des APIs

- Catalogue national publié sur portail développeurs
- Cycle de vie : Design → Review → Approve → Publish → Deprecate
- Versioning obligatoire (semver, /v1, /v2)
- SLA défini par API (latence, dispo, throughput)
- Quotas + throttling par consommateur
- Tests de contrat (Pact) inter-agences

---

## 10. Conventions Inter-Agences

Modèle de **Data Sharing Agreement (DSA)** standardisé :
- Finalité du partage
- Base légale
- Données échangées (champ par champ)
- Durée de conservation côté consommateur
- Mesures de sécurité
- Responsabilités en cas d'incident
- Auditabilité

Toutes les DSA enregistrées auprès de la **NDPA**.

---

## 11. KPI Interopérabilité

| KPI | Cible 2028 |
|-----|------------|
| Agences connectées au bus | ≥ 15 |
| APIs publiées au catalogue | ≥ 200 |
| Disponibilité bus national | ≥ 99,95 % |
| Latence p95 API publique | < 300 ms |
| Couverture événements clés (naissance, décès, mariage) | 100 % |

---

## 12. Étapes de Déploiement

| Phase | Cible |
|-------|-------|
| Onboarding 1 (2027) | ONI, MJSP, CNIGS, PNH |
| Onboarding 2 (2028) | DGI, AGD, MSPP, MENFP |
| Onboarding 3 (2029) | OFATMA, ONA, CEP, INARA |
| Onboarding 4 (2030) | Banques, télécom, partenaires régionaux |

---
*Fin du document — Étape 8/16*
