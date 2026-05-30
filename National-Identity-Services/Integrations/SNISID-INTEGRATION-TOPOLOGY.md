# 🔗 SNISID — TOPOLOGIE D'INTÉGRATION NATIONALE
## Government Integration Bus & Agency API Contracts

**Document ID :** SNISID-INT-001  
**Version :** 1.0.0  
**Date :** Mai 2026  
**Classification :** SOUVERAIN / INTEROPÉRABILITÉ

---

## 1. ARCHITECTURE D'INTÉGRATION X-ROAD

SNISID utilise **X-Road** comme bus national de services gouvernementaux — une infrastructure d'interopérabilité souveraine à l'origine développée en Estonie et adoptée par de nombreux pays comme standard d'eGovernment.

```mermaid
graph TB
    subgraph SNISID_CORE["🏛️ SNISID Core"]
        SS_SNISID[X-Road Security Server\nSNISID Provider]
        ID_SVC[Identity Service]
        CIV_SVC[Civil Registry]
        BIO_SVC[Biometric Service]
        VERIFY_SVC[Verification Service]
    end

    subgraph XROAD["🔒 X-Road National Bus"]
        CS[Central Server\n(Certification + Registry)]
        TSA[Timestamp Authority\nRFC 3161]
        OCSP[OCSP Responder\nPKI SNISID]
    end

    subgraph AGENCIES["🏢 Agences Connectées"]
        SS_ONI[SS: ONI\nConsumer + Provider]
        SS_ANH[SS: ANH\nConsumer]
        SS_DGI[SS: DGI\nConsumer]
        SS_CEP[SS: CEP\nConsumer]
        SS_MSPP[SS: MSPP\nConsumer + Provider]
        SS_PNH[SS: PNH\nConsumer]
        SS_BCN[SS: BCN\nConsumer]
        SS_OMRH[SS: OMRH\nConsumer]
        SS_MENFP[SS: MENFP\nConsumer]
        SS_OFATMA[SS: OFATMA\nConsumer]
        SS_JUSTICE[SS: Justice/Tribunaux\nConsumer + Provider]
    end

    SS_SNISID <--> XROAD
    XROAD <--> SS_ONI & SS_ANH & SS_DGI & SS_CEP
    XROAD <--> SS_MSPP & SS_PNH & SS_BCN & SS_OMRH
    XROAD <--> SS_MENFP & SS_OFATMA & SS_JUSTICE

    note1["🔐 Toutes communications:\n• TLS 1.3 obligatoire\n• Signature messages\n• Logging centralisé\n• OCSP validation certificats"]
```

---

## 2. INTÉGRATIONS PAR AGENCE

### 2.1 ONI — Office National d'Identification

**Rôle :** Partenaire principal — registre civil historique + émission CIN/NNI physique.

| Aspect | Détail |
|--------|--------|
| **Type** | Synchrone REST + Kafka Async |
| **Niveau accès** | Provider (partage données) + Consumer (lit SNISID) |
| **Scope OAuth** | `identity:read, identity:write, biometric:enroll` |
| **Events consommés** | `identity.activated, identity.updated, identity.deceased` |
| **APIs exposées** | Citizen registration, NIN lookup, card issuance status |
| **Volume** | ~14,000 opérations/jour (période d'enrôlement) |

```yaml
# X-Road service descriptor — ONI
xroad_member:
  member_class: GOV
  member_code: ONI-001
  subsystem: NationalIdentity
  services:
    - name: getCitizenByNIU
      url: https://oni.gov.ht/api/v1/citizens/{niu}
      method: GET
      timeout_ms: 3000
    - name: updateCardStatus
      url: https://oni.gov.ht/api/v1/cards/status
      method: POST
      timeout_ms: 5000
```

### 2.2 ANH — Archives Nationales d'Haïti

**Rôle :** Dépositaire des archives civiles historiques + source de décès.

| Aspect | Détail |
|--------|--------|
| **Type** | Kafka Async (push décès) + REST Query |
| **Events produits** | `civil.death.anh.registered` → déclenche cascade |
| **APIs consommées** | Birth certificate archive lookup, Historical record migration |
| **Scope OAuth** | `civil:read, identity:deceased_update` |
| **SLA** | Notification décès → SNISID en < 2h |

### 2.3 Justice / Tribunaux (MJSP)

**Rôle :** Source de vérité judiciaire (jugements supplétifs, divorces, adoptions, ordonnances de suspension).

| Aspect | Détail |
|--------|--------|
| **Type** | REST sync (vérification jugements) + Kafka (décisions) |
| **Events produits** | `judicial.order.suspension, judicial.decree.birth` |
| **APIs clés** | `GET /judgments/{id}/verify` — vérification jugement supplétif |
| **Scope OAuth** | `civil:judicial_verify, identity:suspend` |
| **Cas d'usage** | Suspension NIU sur ordonnance, enrôlement EC-N03/N05 |

```yaml
# Contrat API Justice
paths:
  /v1/judgments/{judgment_id}/verify:
    get:
      summary: Vérifier l'authenticité d'un jugement de TPI
      parameters:
        - name: judgment_id
          in: path
          required: true
          schema: { type: string }
      responses:
        '200':
          content:
            application/json:
              schema:
                properties:
                  valid: { type: boolean }
                  tribunal: { type: string }
                  judge_name: { type: string }
                  date_judgment: { type: string, format: date }
                  seal_verified: { type: boolean }
                  subject_niu: { type: string }
```

### 2.4 PNH — Police Nationale d'Haïti

**Rôle :** Vérification d'identité aux points de contrôle, investigation fraude (DCPJ).

| Aspect | Détail |
|--------|--------|
| **Type** | REST sync + Event consumer |
| **Niveau accès** | Consumer (vérification identité) + DCPJ (investigation) |
| **Scope OAuth** | `identity:verify, biometric:verify, fraud:investigate` |
| **Cas d'usage** | Checkpoint identité, mandats d'arrêt, dossiers fraude |
| **Mode offline** | Cache edge node département (5000 templates) |
| **Accès DCPJ (Niveau 3)** | Full audit trail, forensic data, suspended/revoked NIUs |

### 2.5 CEP — Conseil Électoral Permanent

| Événements reçus | Action |
|-----------------|--------|
| `identity.activated` (âge ≥ 18) | Pré-inscrire sur listes électorales |
| `identity.deceased` | Retirer des listes électorales |
| `identity.suspended` | Suspendre droits de vote |
| `civil.address.updated` | Mettre à jour bureau de vote |

### 2.6 DGI — Direction Générale des Impôts

| Événements reçus | Action |
|-----------------|--------|
| `identity.activated` | Créer dossier fiscal |
| `civil.marriage.registered` | Déclaration fiscale commune |
| `identity.deceased` | Clôturer dossier, notifier héritiers |
| `identity.suspended` | Bloquer remboursements TVA |

### 2.7 MSPP — Ministère Santé Publique et Population

| Type | Détail |
|------|--------|
| **Intégration FHIR R4** | Vérification attestations médicales |
| **Events reçus** | `identity.activated, identity.deceased` |
| **APIs exposées** | `POST /fhir/Patient` — création patient sur activation NIU |
| **Scope** | `identity:read, civil:health_events` |

---

## 3. CONTRATS D'API INTER-AGENCES

### 3.1 Modèle OAuth 2.1 par Scope

```yaml
# iam/keycloak/realm-snisid/client-scopes.yaml
client_scopes:
  identity:read:
    description: "Lire le statut et profil minimal d'une identité"
    agencies: [ONI, DGI, CEP, MSPP, MENFP, BCN, OFATMA, PNH, OMRH]

  identity:write:
    description: "Créer ou modifier des identités (enrôlement)"
    agencies: [ONI]

  identity:suspend:
    description: "Suspendre/révoquer une identité (mesure judiciaire)"
    agencies: [JUSTICE, DCPJ]

  biometric:verify:
    description: "Vérification biométrique 1:1"
    agencies: [PNH, ONI, BANQUES_PARTENAIRES]

  civil:read:
    description: "Lire les actes d'état civil"
    agencies: [ONI, DGI, CEP, MSPP, OMRH, ANH]

  civil:write:
    description: "Créer des actes d'état civil"
    agencies: [ONI, ANH, JUSTICE]

  fraud:investigate:
    description: "Accès aux dossiers fraude (niveau DCPJ)"
    agencies: [DCPJ]

  audit:read:
    description: "Lire les journaux d'audit"
    agencies: [AND, NDPA, CNN_AUDIT]
```

### 3.2 Data Sharing Agreement Template (DSA)

```markdown
# Convention de Partage de Données — SNISID ↔ [AGENCE]

**Référence:** DSA-SNISID-[AGENCE]-2026-[N]
**Date d'entrée en vigueur:** [DATE]
**Parties:**
- AND (Autorité Nationale Numérique) — Gestionnaire SNISID
- [NOM AGENCE] — Consommateur

## Données Partagées
| Champ | Finalité | Base Légale | Durée Conservation |
|-------|---------|-------------|-------------------|
| NIU | Identification | Art. X Loi SNISID | Permanente |
| Statut identité | Vérification | Art. Y | Temps réel |
| [Autres champs] | [Finalité] | [Base] | [Durée] |

## Obligations
- Minimisation: utiliser uniquement les données strictement nécessaires
- Sécurité: TLS 1.3 + mTLS obligatoire
- Audit: tracer chaque accès
- Notification: signaler toute violation dans les 72h
- Sous-traitance: interdite sans accord écrit AND

## Sanctions
- Violation grave: suspension immédiate de l'accès
- Amende: jusqu'à 2% du budget annuel de l'agence
- Récidive: suspension définitive + référé NDPA

**Signataires:** DG-AND | DG-[Agence] | NDPA
```

---

## 4. PLATEFORME D'ENRÔLEMENT — ARCHITECTURE

### 4.1 Capacité Cible

| Canal | Kits | Capacité/Jour | Capacité/An |
|-------|------|--------------|-------------|
| Fixed-Site (10 centres) | — | 2,000 | 730,000 |
| Mobile Kit (MEK) | 150+ kits | 12,000 | 4,380,000 |
| **Total** | — | **14,000** | **5,110,000** |
| **Objectif Année 1** | | | **3,000,000** |

### 4.2 Matériel MEK (Mobile Enrollment Kit)

```yaml
# Manifest Matériel MEK v1.0
mek_hardware:
  tablet:
    model: "Panasonic TOUGHBOOK FZ-T1 (10 pouces ruggedisé)"
    specs: { os: Android 12, ram: 4GB, storage: 64GB, ip: IP68, mil: MIL-STD-810G }
    weight: 650g

  edge_compute:
    model: "Intel NUC 12 Pro"
    specs: { cpu: i5-1240P, ram: 16GB, storage: 512GB NVMe, tpm: TPM2.0 }
    weight: 1200g

  fingerprint_scanner:
    model: "Aware BioSled USB (FBI certifié 500dpi)"
    standard: "FBI Appendix F / STQC"
    weight: 400g

  iris_camera:
    model: "IrisAccess iCAM TD100"
    standard: "ISO 19794-6"
    weight: 350g

  face_camera:
    model: "Cognitec FaceVACS-Access 3D"
    standard: "ISO 19794-5 ICAO"
    weight: 200g

  document_scanner:
    model: "IRIScan Book 5 + MRZ reader"
    standard: "ICAO 9303"
    weight: 500g

  power:
    solar_panel: { model: "ETFE 100W foldable", weight: 2500g }
    battery: { model: "EcoFlow RIVER 2 288Wh LiFePO4", runtime: "8h full", weight: 3000g }

  connectivity:
    router: { model: "Sierra Wireless RV55 Dual-SIM", sims: ["Digicel HT", "Natcom HT"], weight: 300g }

  transport:
    case: { model: "Pelican 1535 AirCarryOn", ip: IP67, weight: 4000g }

  total_weight: "~13.1 kg"
  cost_per_kit_usd: 8500
  total_fleet_150_kits: "1,275,000 USD"
```

### 4.3 Programme de Formation Agents (5 jours)

| Jour | Contenu | Durée |
|------|---------|-------|
| **J1** | SNISID mission, éthique, protection données, droits citoyens | 8h |
| **J2** | Matériel MEK: montage, démontage, diagnostic pannes, énergie solaire | 8h |
| **J3** | Logiciel: app tablet, saisie données, capture biométrique, qualité | 8h |
| **J4** | Cas spéciaux: accommodations, minorités, enfants, personnes âgées | 8h |
| **J5** | Mode offline: synchronisation, gestion conflits, escalade incidents | 8h |
| **Certification** | Examen pratique + théorique (score min 80%) | 4h |

---

## 5. TABLEAU DE BORD OPÉRATIONNEL

### 5.1 Dashboard Exécutif — Métriques Prometheus

```yaml
# observability/prometheus/identity-dashboards.yaml
metrics:
  - name: snisid_citizens_total
    type: gauge
    labels: [statut, departement]
    description: "Total des identités par statut et département"

  - name: snisid_enrollments_daily
    type: counter
    labels: [canal, departement]
    description: "Enrôlements quotidiens par canal"

  - name: snisid_fraud_alerts_total
    type: counter
    labels: [type, severity]
    description: "Alertes fraude cumulées"

  - name: snisid_verification_duration_ms
    type: histogram
    buckets: [50, 100, 200, 500, 1000, 2000]
    description: "Latence des vérifications d'identité"

  - name: snisid_kit_battery_pct
    type: gauge
    labels: [kit_id, departement]
    description: "Niveau batterie des kits terrain"

  - name: snisid_offline_queue_depth
    type: gauge
    labels: [device_id]
    description: "Enrôlements en attente de sync"
```

### 5.2 Alertes KPI Phase 2

| KPI | Objectif | Alerte Warning | Alerte Critical |
|-----|---------|----------------|-----------------|
| Enrôlements/jour | ≥ 10,000 | < 8,000 | < 5,000 |
| Taux fraude détectée | ≤ 2% | > 3% | > 5% |
| Disponibilité API Identity | 99.95% | < 99.9% | < 99.5% |
| Délai cascade décès | < 5 min | > 10 min | > 30 min |
| Sync queue depth | < 500 | > 1000 | > 5000 |
| OEC signature validité | 100% | < 99% | < 95% |

---

## 6. RAPPORT DE VALIDATION FINALE — PHASE 2

### 6.1 Checklist Production-Readiness (80+ items)

**Identity Registry (20 items)**
- [x] Schéma CockroachDB complet (citizens, events, snapshots, documents)
- [x] Event Sourcing immutable (triggers + rules anti-UPDATE/DELETE)
- [x] CQRS: OpenSearch synchronisé via Debezium + Kafka
- [x] NIU généré cryptographiquement (crypto/rand, no demographic metadata)
- [x] Lifecycle state machine 7 états (PRE_REGISTERED → ARCHIVED)
- [x] API OpenAPI 3.1 complète (CRUD + history + search)
- [x] Kafka events catalog (Avro schemas v1)
- [x] Kubernetes manifests (Deployment, HPA, PDB, NetworkPolicy, KEDA)
- [x] ABAC policies OPA Rego (commune, département, rôle)
- [x] RLS PostgreSQL par agence
- [x] Partitionnement par département
- [x] HA: active-active PaP + Cap-Haïtien
- [x] Cascade notifications décès (8 agences, < 5 min)
- [x] Optimistic locking (version field)
- [x] GDPR anonymization function
- [x] Audit triggers automatiques
- [x] SLOs définis (99.95% availability, P99 < 2s)
- [x] Alerting Prometheus (10+ règles)
- [x] DR: RTO 2 min, RPO 0 (CockroachDB consensus)
- [x] Row-level security activée

**Biometric Platform (15 items)**
- [x] ABIS architecture multi-GPU documentée
- [x] 3 modalités: empreintes ISO 19794-2, iris ISO 19794-6, visage ISO 19794-5
- [x] PAD anti-spoofing (9 vecteurs d'attaque couverts)
- [x] Taux détection PAD: ≥ 97.8% sur tous les vecteurs
- [x] 1:N deduplication workflow (Mermaid + SLA 30s)
- [x] 1:1 verification workflow (SLA < 100ms)
- [x] Cryptographic template protection (AES-256-GCM + HSM)
- [x] API OpenAPI 3.1 (enroll, verify, identify, status, revoke)
- [x] Offline cache (5K templates/commune, 30j autonomie)
- [x] Exception accommodations documentées (6 cas)
- [x] Gouvernance données biométriques (RGPD + Convention 108+)
- [x] DCPJ adjudication workflow
- [x] HSM key hierarchy
- [x] Conflict resolution (doublon) + ConflictCase management
- [x] Audit log par vérification

**Civil Registry (15 items)**
- [x] 5 types de naissance (EC-N01 à EC-N05) — BPMN complets
- [x] Mariage civil + concordataire
- [x] Décès + cascade 8 agences (< 5 min SLA)
- [x] Divorce contentieux + mutuel
- [x] Adoption simple + plénière
- [x] Corrections administratives + annotations marginales
- [x] OEC authentication (FIDO2 + PKI)
- [x] PDF/A-3 + XAdES-LTA generation
- [x] QR code JWT (validité 5 ans, vérifiable offline)
- [x] Schéma SQL complet (civil_acts, officiers, corrections)
- [x] Kafka events pour chaque type d'acte
- [x] Mode offline (30j autonomie OEC)
- [x] SLA définis par type d'acte
- [x] DMN rules tables
- [x] API REST complète

**Fraud Detection (10 items)**
- [x] 14-dimension feature vector
- [x] Random Forest + DNN Ensemble
- [x] TensorFlow Serving deployment
- [x] 30+ velocity rules (YAML)
- [x] 5 modèles spécialisés (velocity, document, biometric, agent, geo)
- [x] Ghost worker detection (cross-OMRH)
- [x] Insider threat analytics
- [x] DCPJ referral workflow
- [x] Fraud KPIs + alerting
- [x] WORM audit trail (Merkle chain)

**Integrations (10 items)**
- [x] X-Road architecture (11 agences connectées)
- [x] DSA template (Data Sharing Agreement)
- [x] OAuth 2.1 scopes par agence
- [x] API contracts (OpenAPI) pour chaque agence
- [x] Rate limiting par agence
- [x] Kafka event contracts (AsyncAPI)
- [x] ONI: bidirectionnel (provider + consumer)
- [x] ANH: décès integration
- [x] Justice: jugements supplétifs
- [x] PNH: vérification + investigation DCPJ

**National Resilience (Phase 2 ZIP — intégré)**
- [x] 28 fichiers National-Resilience migrés et intégrés
- [x] Framework continuité (RTO/RPO définis)
- [x] DR multi-région (PaP + Cap-Haïtien)
- [x] Runbooks complets (5 types)
- [x] Crisis coordination platform
- [x] Offline survival (30j autonomie état)
- [x] Cyber resilience model
- [x] Power resilience (solaire + diesel)
- [x] KPI model résilience
- [x] Recovery automation (Terraform + ArgoCD + Velero)

### 6.2 Critères GO/NO-GO pour Phase 3

**GO si :**
- ✅ Tous les 80+ items checklist validés
- ✅ Test d'enrôlement: 1,000 citoyens en simulation sans erreur
- ✅ Test ABIS: déduplication fonctionnelle (< 30s pour 15M templates)
- ✅ Test cascade décès: 8 agences notifiées en < 5 min
- ✅ Test offline: 100 enrôlements + sync réussie
- ✅ Audit de sécurité externe (pentest + CIS K8s Benchmark)
- ✅ Formation: 200 agents ONI certifiés
- ✅ Approbation légale: Ministère Justice + NDPA

**NO-GO si :**
- ❌ ABIS 1:N > 60s sur 15M templates
- ❌ False Acceptance Rate (FAR) > 0.1%
- ❌ Cascade décès > 30 min
- ❌ Disponibilité API < 99.9% sur test de 72h
- ❌ Vulnérabilité critique non corrigée (CVSS ≥ 9.0)

---

*Document ID : SNISID-INT-001 + SNISID-PH2-VALID-001 v1.0.0 — Mai 2026*  
*Approuvé par : DG-AND | CISO | DG-ONI | Ministère Justice | NDPA*  
*Prochaine revue : Phase 3 Launch (Q3 2026)*
