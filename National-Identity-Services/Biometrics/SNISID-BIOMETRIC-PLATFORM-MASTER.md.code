# 🔬 SNISID — PLATEFORME BIOMÉTRIQUE NATIONALE
## National Biometric Platform — Architecture & Anti-Fraud Model

**Document ID :** SNISID-BIO-001  
**Version :** 1.0.0  
**Date :** Mai 2026  
**Classification :** SOUVERAIN / BIOMÉTRIQUE SENSIBLE  
**Principe absolu :** Un citoyen, une identité biométrique — déduplication ABIS obligatoire

---

## 1. VISION & MISSION BIOMÉTRIQUE

La plateforme biométrique SNISID est l'**ancre absolue de l'unicité identitaire** haïtienne. Elle garantit qu'aucun citoyen ne peut être enregistré deux fois sous des identités différentes, même avec des noms et dates de naissance altérés.

**La biométrie répond à 3 questions souveraines :**
1. **Qui êtes-vous ?** (1:N identification — recherche dans 15M+ templates)
2. **Êtes-vous bien cette personne ?** (1:1 vérification — comparaison ciblée)
3. **Êtes-vous vivant ?** (PAD — liveness detection anti-spoofing)

---

## 2. ARCHITECTURE MULTI-MODALE

### 2.1 Vue d'Ensemble

```mermaid
graph TB
    subgraph CAPTURE["📸 Couche Capture (Multi-Modal)"]
        FP_SCAN[Scanners Empreintes\nFBI certifiés ISO 19794-2]
        IRIS_CAM[Caméras Iris Dual NIR\nISO 19794-6]
        FACE_3D[Caméra Faciale 3D\nISO 19794-5 ICAO]
        PAD_ENGINE[PAD Engine\nLiveness Detection]
    end

    subgraph PROCESSING["⚙️ Couche Traitement"]
        QUALITY[Quality Assessment\nNFIQ2 + ISO checks]
        EXTRACT[Feature Extraction\nTemplate Generation]
        ENCRYPT[Template Encryption\nAES-256-GCM + HSM]
    end

    subgraph ABIS_CORE["🖥️ ABIS — Automated Biometric ID System"]
        ABIS_API[ABIS API Controller\nGolang]
        GPU_MATCH[GPU Matching Nodes\nCUDA — C++ Engine]
        GALLERY[(Encrypted Template Gallery\nIn-Memory — LUKS+TPM)]
        ADJUD[Adjudication Interface\nDCPJ Forensic Examiner]
    end

    subgraph SYNC["🔄 Offline Sync"]
        EDGE_CACHE[Edge Commune Cache\n5000 templates/commune]
        NATS[NATS JetStream\nOffline Queue]
    end

    CAPTURE --> PROCESSING
    PROCESSING --> ABIS_CORE
    ABIS_CORE --> SYNC
    GPU_MATCH --> ADJUD
```

### 2.2 Spécifications Techniques par Modalité

| Modalité | Standard | Format Template | Qualité Min | Retries Max | Fallback |
|---------|---------|-----------------|-------------|-------------|----------|
| **Empreintes (10-print)** | ISO 19794-2 | Minutiae WSQ | NFIQ2 ≥ 40 (par doigt) | 3 par doigt | Accepter 8/10 doigts |
| **Iris (Dual)** | ISO 19794-6 | IrisCode NIR | Usable area ≥ 70% | 3 par œil | Skip si cert. médical |
| **Visage (3D)** | ISO 19794-5 ICAO | Feature Vector + 3D | 23 checks ICAO | 5 | Reposition + lumière |
| **Liveness** | ISO 30107-3 | PAD Score (0-100) | ≥ 95% | 3 | Escalade superviseur |

---

## 3. PIPELINE ANTI-SPOOFING (PAD)

### 3.1 Architecture PAD Multi-Couches

```mermaid
graph TD
    subgraph FP_PAD["Empreintes — PAD"]
        FP_MULTI[Capteur Multispectral\nVisible + NIR + UV]
        FP_SWEAT[Détection pores sueur\nCNN]
        FP_PULSE[Oxymétrie pouls\nsignal vasculaire]
        FP_TEXTURE[Texture peau\nResNet-50]
    end

    subgraph IRIS_PAD["Iris — PAD"]
        IR_NIR[Fusion NIR + Visible]
        IR_PUPIL[Réponse pupillaire\nflash test]
        IR_REFLECT[Réflexion cornéenne]
        IR_TEXTURE[Texture iris\nCNN]
    end

    subgraph FACE_PAD["Visage — PAD"]
        FC_3D[Carte profondeur\n3D Structured Light]
        FC_BLINK[Clignement + micro-expr.\nOptical Flow]
        FC_MORPH[Détection morphing\nMEF-GAN-Detector]
        FC_DEEP[Artefacts deepfake\nTemporel CNN]
    end

    FP_SWEAT & FP_PULSE & FP_TEXTURE --> FP_SCORE[PAD Score\nEmpreintes]
    IR_PUPIL & IR_REFLECT & IR_TEXTURE --> IR_SCORE[PAD Score\nIris]
    FC_3D & FC_BLINK & FC_MORPH & FC_DEEP --> FC_SCORE[PAD Score\nVisage]

    FP_SCORE --> FUSION[Fusion Engine\nWeighted Average]
    IR_SCORE --> FUSION
    FC_SCORE --> FUSION

    FUSION --> DECISION{Score combiné}
    DECISION -->|"≥ 95% : VIVANT"| PROCEED[✅ Continuer]
    DECISION -->|"80-94% : INCERTAIN"| MANUAL[⚠️ Override agent]
    DECISION -->|"< 80% : ATTAQUE"| REJECT[❌ Rejet + Alerte SOC]
```

### 3.2 Matrice des Vecteurs d'Attaque

| Type d'Attaque | Modalité | Méthode Détection | Taux Détection | FRR |
|---------------|---------|------------------|----------------|-----|
| Photo 2D imprimée | Visage | Analyse carte profondeur | **99.9%** | < 0.01% |
| Vidéo sur écran | Visage | Détection motif Moiré + reflets | **99.7%** | < 0.05% |
| Masque silicone 3D | Visage | Texture peau + thermique | **99.2%** | < 0.1% |
| Deepfake vidéo | Visage | Artefacts temporels CNN | **98.5%** | < 0.2% |
| Photo morphée | Visage | MAD différentiel | **97.8%** | < 0.3% |
| Doigt gélatine/silicone | Empreinte | Multispectral + pouls | **99.5%** | < 0.05% |
| Empreinte latente levée | Empreinte | Absence pores sueur | **99.8%** | < 0.01% |
| Iris imprimé papier | Iris | Test réponse pupillaire | **99.9%** | < 0.01% |
| Lentille de contact | Iris | Irrégularité bord pupille | **98.0%** | < 0.5% |

---

## 4. WORKFLOWS BIOMÉTRIQUES ABIS

### 4.1 Déduplication 1:N (Enrôlement)

```mermaid
sequenceDiagram
    participant AGENT as 👤 Agent ONI
    participant GW as API Gateway
    participant BIO_SVC as Biometric Service
    participant PAD as PAD Engine
    participant ABIS as ABIS GPU Cluster
    participant WF as Temporal Workflow
    participant DCPJ as DCPJ Adjudicateur

    AGENT->>GW: POST /v1/biometrics/enroll\n{niu, templates, pad_scores}
    GW->>BIO_SVC: Forward (mTLS vérifié)
    BIO_SVC->>PAD: Vérifier PAD scores (threshold ≥ 95%)
    PAD-->>BIO_SVC: PAD: LIVE (score: 97.3%)

    BIO_SVC->>BIO_SVC: Quality check NFIQ2 + ISO 19794
    BIO_SVC->>WF: Start DeduplicationWorkflow(niu, templates)
    WF->>ABIS: 1:N Search (15M+ templates)

    Note over ABIS: GPU Search in progress...\nTypically 10-30 seconds on\nNVIDIA A100 cluster

    alt Doublon détecté (score ≥ 85%)
        ABIS-->>WF: Match: NIU-7891234560 (score: 93.2%)
        WF->>WF: Geler les 2 NIUs → SUSPENDED
        WF->>DCPJ: Créer ConflictCase #CC-2026-001\nRouter vers adjudication manuelle
        WF-->>AGENT: 202 CONFLICT — Référence #CC-2026-001
        Note over DCPJ: Examinateur DCPJ compare\nempreintes côte à côte (48h SLA)
    else Unique (aucun doublon)
        ABIS-->>WF: NO_MATCH — Citoyen unique
        WF->>BIO_SVC: Stocker templates chiffrés\ndans ABIS Gallery
        WF->>WF: Mettre statut NIU → ACTIVE
        WF-->>AGENT: 200 ENROLLED — NIU: 7392851046
    end
```

### 4.2 Vérification 1:1 (Authentication)

```mermaid
sequenceDiagram
    participant BANK as 🏦 Banque\n(Agence habilitée)
    participant GW as API Gateway\n(Kong + OAuth2.1)
    participant BIO_SVC as Biometric Service
    participant CACHE as Redis Cache\n(30s TTL)
    participant ABIS as ABIS API
    participant AUDIT as Audit Log\n(WORM Kafka)

    BANK->>GW: POST /v1/biometrics/verify\n{niu: "7392851046", modality: "FINGERPRINT", template_b64: "..."}
    Note over GW: Valide: OAuth scope "biometric:verify"\nVérifie: mTLS certificat banque\nApplique: rate limit 500K/jour

    GW->>BIO_SVC: Forward (scope vérifié)
    BIO_SVC->>CACHE: GET biometric_ref:7392851046:FINGERPRINT

    alt Cache HIT
        CACHE-->>BIO_SVC: abis_gallery_id: "ABIS-001-7392"
        BIO_SVC->>ABIS: 1:1 Match (gallery_id, template_query)
    else Cache MISS
        BIO_SVC->>ABIS: Lookup + 1:1 Match
        BIO_SVC->>CACHE: SET (30s TTL)
    end

    ABIS-->>BIO_SVC: Match Score: 94.7 (seuil: 80)
    BIO_SVC->>AUDIT: Publish: BIOMETRIC.VERIFIED\n{niu, agency, score, timestamp}

    BIO_SVC-->>BANK: 200 VERIFIED\n{match_score: 94.7, decision: "MATCH", niu_status: "ACTIVE"}

    Note over BANK: Latence totale: ~85ms (P99)\nNIU status vérifié en temps réel
```

---

## 5. PROTECTION CRYPTOGRAPHIQUE DES TEMPLATES

### 5.1 Modèle "No Raw Storage"

```
Capture biométrique brute (WSQ/JPEG2000)
    ↓
Feature Extraction (algorithme propriétaire ABIS)
    ↓
Template mathématique IRRÉVERSIBLE (vecteur numérique)
    ↓
Chiffrement AES-256-GCM avec clé dérivée du NIU
    ↓
Stockage dans ABIS Gallery (LUKS + TPM hardware)
    ↓
Images brutes → Archivage WORM (Vault + air-gap)
             → Suppression locale après extraction
```

**Propriétés de sécurité :**
- ❌ **Aucune image brute** en mémoire après extraction du template
- ❌ **Aucun template en clair** dans les bases de données
- ✅ **Template irréversible** : impossible de reconstruire l'empreinte
- ✅ **Clé de chiffrement** stockée dans HSM FIPS 140-2 Level 3
- ✅ **LUKS** sur tous les volumes de stockage ABIS
- ✅ **TPM 2.0** pour attestation d'intégrité matérielle

### 5.2 Hiérarchie de Clés Biométriques

```
HSM National (Root Key — Air-gapped)
    └── ABIS Master Key (HSM Online FIPS 140-2 L3)
            ├── Template Encryption Key per Modality
            │       ├── TEK-FINGERPRINT (AES-256)
            │       ├── TEK-IRIS (AES-256)
            │       └── TEK-FACE (AES-256)
            └── Gallery Index Key
```

---

## 6. API CONTRACT BIOMÉTRIQUE (OpenAPI 3.1)

```yaml
openapi: "3.1.0"
info:
  title: SNISID Biometric Service API
  version: "1.0.0"

paths:
  /v1/biometrics/enroll:
    post:
      operationId: enrollBiometrics
      summary: Enrôler les données biométriques d'un citoyen
      security:
        - oauth2: [biometric:enroll]
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [niu, templates]
              properties:
                niu:
                  type: string
                  pattern: '^\d{10}$'
                templates:
                  type: object
                  properties:
                    fingerprints:
                      type: string
                      format: byte
                      description: "ISO 19794-2 template, Base64, AES-256 encrypted"
                    iris:
                      type: string
                      format: byte
                      description: "ISO 19794-6 dual-iris, Base64"
                    face:
                      type: string
                      format: byte
                      description: "ISO 19794-5 3D template, Base64"
                pad_scores:
                  type: object
                  properties:
                    fingerprint_pad: { type: number, minimum: 0, maximum: 100 }
                    iris_pad: { type: number }
                    face_pad: { type: number }
                    combined_pad: { type: number }
                quality_scores:
                  type: object
                  properties:
                    nfiq2_mean: { type: number }
                    iris_quality: { type: number }
                    icao_compliance: { type: boolean }
      responses:
        '202':
          description: "Accepté — déduplication ABIS en cours (asynchrone)"
          content:
            application/json:
              schema:
                properties:
                  dedup_job_id: { type: string }
                  status: { type: string, enum: [PENDING_DEDUP] }
                  eta_seconds: { type: integer, example: 30 }
        '409':
          description: "Doublon biométrique détecté"
          content:
            application/json:
              schema:
                properties:
                  conflict_case_id: { type: string }
                  conflicting_niu: { type: string }
                  match_score: { type: number }
                  status: { type: string, enum: [CONFLICT_PENDING_ADJUDICATION] }

  /v1/biometrics/verify:
    post:
      operationId: verifyBiometric
      summary: Vérification 1:1 — authentifier un citoyen qui se présente
      security:
        - oauth2: [biometric:verify]
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [niu, modality, template_data]
              properties:
                niu: { type: string, pattern: '^\d{10}$' }
                modality: { type: string, enum: [FINGERPRINT, IRIS, FACE] }
                template_data:
                  type: string
                  format: byte
                  description: "Template capturé en direct (Base64)"
                liveness_score: { type: number, minimum: 0, maximum: 100 }
      responses:
        '200':
          description: "Résultat de la vérification"
          content:
            application/json:
              schema:
                properties:
                  decision: { type: string, enum: [MATCH, NO_MATCH, ERROR] }
                  match_score: { type: number, example: 94.7 }
                  niu_status: { type: string, enum: [ACTIVE, SUSPENDED, DECEASED] }
                  response_time_ms: { type: integer }

  /v1/biometrics/{niu}/status:
    get:
      operationId: getBiometricStatus
      summary: Obtenir le statut d'enrôlement biométrique d'un citoyen
      security:
        - oauth2: [biometric:read]
      parameters:
        - name: niu
          in: path
          required: true
          schema: { type: string }
      responses:
        '200':
          content:
            application/json:
              schema:
                properties:
                  niu: { type: string }
                  fingerprints_enrolled: { type: boolean }
                  iris_enrolled: { type: boolean }
                  face_enrolled: { type: boolean }
                  last_verified_at: { type: string, format: date-time }
                  verification_count: { type: integer }
                  abis_gallery_id: { type: string }

  /v1/biometrics/{niu}/revoke:
    delete:
      operationId: revokeBiometrics
      summary: Révoquer les données biométriques (fraude confirmée)
      security:
        - oauth2: [biometric:revoke]
      parameters:
        - name: niu
          in: path
          required: true
          schema: { type: string }
      requestBody:
        content:
          application/json:
            schema:
              properties:
                reason: { type: string }
                legal_reference: { type: string }
                authorized_by: { type: string }
      responses:
        '200':
          description: "Biométrie révoquée (templates supprimés de la galerie active, archivés pour prévention futur ré-enregistrement)"
```

---

## 7. CACHE BIOMÉTRIQUE OFFLINE (EDGE NODES)

### 7.1 Architecture

Chaque node edge de département maintient un cache chiffré des templates biométriques des citoyens de la commune, permettant la vérification 1:1 **sans connectivité** avec le datacenter central.

```mermaid
graph TD
    subgraph CENTRAL["🏢 SNISID Core (PaP)"]
        ABIS_CENTRAL[ABIS Gallery\n15M+ templates]
        SYNC_MGR[Cache Sync Manager\n(nightly or on-demand)]
    end

    subgraph EDGE_DEPT["📡 Edge Node Département"]
        EDGE_DB[(SQLite Encrypted\nLUKS + TPM)]
        EDGE_SVC[Biometric Verify Service\nK3s Pod]
        CACHE[~50K templates\n commune locale]
    end

    subgraph KIT["🎒 Kit Terrain"]
        KIT_DB[(Local SQLite\nTPM-bound)]
        KIT_SVC[Offline Verify\nAndroid App]
        COMMUNE_CACHE[~5K templates\nsection communale]
    end

    ABIS_CENTRAL -->|Delta sync\n(chiffré, signé)| EDGE_DEPT
    EDGE_DEPT -->|Sync partiel\n(commune ciblée)| KIT

    note1["🔐 Clés de chiffrement:\njamais transférées via réseau\nDérivées localement du HSM USB"]
```

**Contraintes du cache offline :**

| Niveau | Capacité | Population couverte | Connectivité requise | Durée autonomie |
|--------|---------|---------------------|----------------------|-----------------|
| Edge Département | 50K templates | Chef-lieu + environ | 4G/VSAT (sync) | 30 jours |
| Kit Terrain | 5K templates | Section communale | Aucune | 30+ jours |
| Vérif locale | 1:1 uniquement | Commune assignée | Aucune | Permanente |

---

## 8. ACCOMMODATIONS SPÉCIALES

### 8.1 Cas d'Exception

| Condition | Accommodation | Documentation | Impact ABIS |
|-----------|-------------|---------------|-------------|
| **Amputation doigts** | Enregistrer doigts disponibles, flag "PARTIAL_CAPTURE" | Certificat médical requis | Réduction capacité déduplication |
| **Cataractes sévères / Cécité** | Skip capture iris, relayer sur empreintes + visage | Certificat médical | Déduplication moins robuste |
| **Nourrissons (0-5 ans)** | Photo uniquement, mise à jour biométrique planifiée à 5 ans | Acte de naissance | Déduplication différée |
| **Personnes âgées (empreintes usées)** | Seuil NFIQ2 abaissé à ≥ 25, priorité iris + visage | Vérification âge | Taux FAR légèrement plus élevé |
| **Handicap moteur** | Temps de capture étendu, assistance agent, position adaptée | Attestation agent | Sans impact qualité |
| **Anomalie médicale cutanée** | Flag "MEDICAL_CONDITION", documentation MSPP | Dossier médical | Alerte superviseur |

---

## 9. GOUVERNANCE DES DONNÉES BIOMÉTRIQUES

### 9.1 Principes RGPD & Convention 108+

| Principe | Implémentation SNISID |
|---------|----------------------|
| **Finalité** | Biométrie utilisée UNIQUEMENT pour déduplication et vérification identité |
| **Minimisation** | Templates irréversibles — jamais les images brutes accessibles externement |
| **Rétention** | Templates actifs : durée de vie du citoyen. Archives : 7 ans post-décès |
| **Droit d'accès** | Citoyen peut demander confirmation de l'enrôlement via portail SNISID |
| **Droit à l'effacement** | Anonymisation post-archivage (RGPD Art. 17 — sauf obligation légale) |
| **Sécurité** | HSM FIPS 140-2 L3, LUKS, TPM, audit WORM |

### 9.2 Comité de Gouvernance Biométrique

| Rôle | Responsabilité |
|------|---------------|
| **NDPA (Directeur)** | Validation des politiques biométriques nationales |
| **CISO SNISID** | Sécurité et chiffrement des données biométriques |
| **DG-ONI** | Opérations d'enrôlement et qualité des captures |
| **Comité Éthique IA** | Supervision des algorithmes de reconnaissance |
| **Représentant citoyens (CNN)** | Défense des droits fondamentaux |

---

*Document ID : SNISID-BIO-001 v1.0.0 — Mai 2026*  
*Approuvé par : CISO National | DG-ONI | NDPA | Comité Éthique IA*  
*Classification : SOUVERAIN / BIOMÉTRIQUE SENSIBLE — République d'Haïti*
