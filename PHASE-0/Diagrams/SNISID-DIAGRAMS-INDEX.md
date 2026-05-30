# 🗺️ SNISID — DIAGRAMS INDEX
## Index Central des Diagrammes de l'Architecture Nationale

**Document ID :** SNISID-DIA-001  
**Version :** 1.0.0  
**Date :** Mai 2026  
**Classification :** Usage Gouvernemental  

---

## DIAGRAMMES D'ARCHITECTURE PRINCIPALE

### 1. Vue Macro Système — 5 Couches

```mermaid
graph TB
    subgraph L5["L5 — EXPÉRIENCE CITOYENNE"]
        MOB[📱 Mobile SNISID Flutter]
        WEB[🌐 Portail Web Next.js]
        USSD[📞 USSD Gateway]
        KIOSK[🖥️ Kiosques terrain]
        AGT[👤 App Agents]
    end

    subgraph L4["L4 — API GATEWAY SOUVERAIN"]
        KONG[Kong/APISIX + WAF]
        AUTH_GW[OAuth 2.1 Validator]
        RATE[Rate Limiter]
    end

    subgraph L3["L3 — MICROSERVICES DOMAINES"]
        ID_SVC[🆔 Identity Service]
        CIV_SVC[📜 État Civil Service]
        BIO_SVC[🔬 Biometric AFIS]
        JUS_SVC[⚖️ Justice Service]
        POL_SVC[👮 Police Service]
        IMM_SVC[✈️ Immigration Service]
        NOT_SVC[🔔 Notification Service]
        AUD_SVC[🔍 Audit Service]
        DOC_SVC[📄 Document Service]
        KYC_SVC[🏦 KYC Service]
    end

    subgraph L2["L2 — PLATEFORME"]
        K8S[☸️ Kubernetes RKE2]
        KAFKA[📨 Apache Kafka]
        PG[🗄️ PostgreSQL HA]
        VAULT[🔐 HashiCorp Vault]
        KEYCLOAK[🔑 Keycloak OIDC]
        MINIO[💾 MinIO S3]
        ISTIO[🕸️ Istio Service Mesh]
    end

    subgraph L1["L1 — INFRASTRUCTURE SOUVERAINE"]
        DC_PAP[🏢 Datacenter PaP Tier III]
        DC_CAP[🏢 DR Cap-Haïtien]
        EDGE[📡 10 Edge Nodes Départements]
        KIT[🎒 100+ Offline Kits Terrain]
    end

    L5 --> L4
    L4 --> L3
    L3 --> L2
    L2 --> L1
    EDGE --> L3
    KIT --> EDGE
```

---

### 2. Architecture Zero Trust — Flux de Requête

```mermaid
sequenceDiagram
    participant User as 👤 Agent DGI
    participant WAF as 🛡️ WAF ModSecurity
    participant GW as 🚪 API Gateway Kong
    participant KC as 🔑 Keycloak
    participant OPA as ⚖️ OPA Rego
    participant SVC as 🆔 Identity Service
    participant VAULT as 🔐 Vault
    participant DB as 🗄️ PostgreSQL

    User->>WAF: HTTPS Request + JWT
    WAF->>WAF: Sanitization, SQLi/XSS check
    WAF->>GW: Forward (Clean)
    GW->>KC: Validate JWT token
    KC-->>GW: Token Valid (sub, scope, claims)
    GW->>OPA: Check ABAC policy
    Note over OPA: Vérifie: rôle, commune,<br/>horaires, consentement,<br/>device trust
    OPA-->>GW: ALLOW / DENY
    GW->>SVC: Forward via mTLS (Istio)
    SVC->>VAULT: Get DB credentials (dynamic)
    VAULT-->>SVC: short-lived creds (TTL=5min)
    SVC->>DB: Query (TLS)
    DB-->>SVC: Data
    SVC-->>User: Response (encrypted)
    Note over SVC: Audit event → Kafka → WORM
```

---

### 3. Architecture Multi-Datacenter (HA/DR)

```mermaid
graph LR
    subgraph DC_PAP["🏢 Port-au-Prince (PRIMARY — Active)"]
        K8S_PAP[K8s Cluster A\n3 Control + 12 Workers]
        KAFKA_PAP[Kafka Cluster A\n5 Brokers]
        PG_PAP[PostgreSQL Patroni A\n1 Primary + 2 Replicas]
        CEPH_PAP[Ceph Storage A\n500 TB]
        HSM_PAP[HSM Array A\nFIPS 140-2 L3]
    end

    subgraph DC_CAP["🏢 Cap-Haïtien (DR — Hot Standby)"]
        K8S_CAP[K8s Cluster B\n3 Control + 8 Workers]
        KAFKA_CAP[Kafka Cluster B\n3 Brokers]
        PG_CAP[PostgreSQL Patroni B\n1 Primary + 2 Replicas]
        CEPH_CAP[Ceph Storage B\n300 TB]
        HSM_CAP[HSM Array B\nFIPS 140-2 L3]
    end

    subgraph EDGE_NET["📡 Edge Network (10 Départements)"]
        EDGE1[Edge Ouest]
        EDGE2[Edge Nord]
        EDGE3[Edge Artibonite]
        EDGEN[... 7 autres]
    end

    subgraph OFFLINE["🎒 Offline Kits Terrain"]
        KIT1[Kit Section A]
        KIT2[Kit Section B]
        KITN[... 100+ kits]
    end

    KAFKA_PAP <-->|MirrorMaker 2\nSync bidirectionnelle| KAFKA_CAP
    PG_PAP <-->|Patroni Streaming\nRPO < 1 min| PG_CAP
    CEPH_PAP <-->|RBD Mirroring\nAsync| CEPH_CAP

    DC_PAP <-->|Liaison fibre dédiée\n10 Gbps + Microwave backup| DC_CAP
    
    DC_PAP -->|4G/VSAT sync| EDGE1 & EDGE2 & EDGE3 & EDGEN
    
    EDGE1 -->|WiFi/BLE/USB| KIT1 & KIT2 & KITN
```

---

### 4. Architecture PKI Nationale

```mermaid
graph TD
    ROOT["🔒 ROOT CA NATIONALE\n(Offline — Air-gapped — Faraday)\nHSM FIPS 140-2 L4\nECDSA P-384\nValidité: 25 ans"]

    ROOT -->|Signed| POL_CIT["Policy CA — Citoyens\nValidité: 10 ans"]
    ROOT -->|Signed| POL_GOV["Policy CA — Officials\nValidité: 10 ans"]
    ROOT -->|Signed| POL_DEV["Policy CA — Devices\nValidité: 10 ans"]
    ROOT -->|Signed| POL_DOC["Policy CA — Documents\nValidité: 10 ans"]

    POL_CIT -->|Signed| ISS_EID["Issuing CA — eID\n(HSM L3 — Online)\nValidité: 5 ans"]
    POL_GOV -->|Signed| ISS_GOV["Issuing CA — GovSign\n(HSM L3 — Online)"]
    POL_DEV -->|Signed| ISS_TLS["Issuing CA — TLS\n(HSM L3 — Online)"]
    POL_DEV -->|Signed| ISS_IOT["Issuing CA — IoT/Edge\n(HSM L3 — Online)"]
    POL_DOC -->|Signed| ISS_SIGN["Issuing CA — DocSign\n(HSM L3 — Online)"]

    ISS_EID -->|Issues| CERT_CIN[Certificats CIN Citoyens\nECDSA P-256 — 5 ans]
    ISS_GOV -->|Issues| CERT_AGENT[Certificats Agents d'État\n3 ans]
    ISS_TLS -->|Issues| CERT_TLS[Certificats TLS Services\nAuto-renouvelés 90 jours]
    ISS_IOT -->|Issues| CERT_EDGE[Certificats Edge Nodes\n1 an]
    ISS_SIGN -->|Issues| CERT_DOC[Certificats Signature Documents\n1 an]
```

---

### 5. Architecture SOC National

```mermaid
graph TD
    subgraph INGEST["Sources de Données"]
        SRC_K8S[K8s Audit Logs]
        SRC_NET[Network Flow Cilium eBPF]
        SRC_APP[Application Logs]
        SRC_HSM[HSM Audit Logs]
        SRC_WAF[WAF Logs]
        SRC_EDR[EDR Agents]
    end

    subgraph COLLECT["Collecte & Normalisation"]
        FB[Fluent Bit DaemonSet]
        VECTOR[Vector.dev]
    end

    subgraph SIEM_STACK["SIEM — Elastic Security / Wazuh"]
        ES[Elasticsearch Cluster]
        WAZUH[Wazuh Manager]
        KIBANA[Kibana / Grafana]
    end

    subgraph SOAR["SOAR — Orchestration"]
        SOAR_ENG[Shuffle SOAR / TheHive]
        PLAYBOOK[Playbooks Automatisés]
    end

    subgraph SOC_TIERS["SOC Tiers"]
        T1[Tier 1 — Surveillance 24/7\n8 analystes par shift]
        T2[Tier 2 — Analyse & IR\n5 experts]
        T3[Tier 3 — Threat Hunting\n3 seniors + CISO]
    end

    subgraph CTI["Threat Intelligence"]
        MISP[MISP Platform]
        TIP[TIP Feeds]
        FIRST[FIRST.org Integration]
    end

    INGEST --> COLLECT
    COLLECT --> SIEM_STACK
    SIEM_STACK --> SOAR
    SOAR --> SOC_TIERS
    CTI --> SIEM_STACK
    CTI --> SOC_TIERS
```

---

### 6. Architecture Interopérabilité X-Road

```mermaid
sequenceDiagram
    participant DGI as DGI - Système Fiscal
    participant SS_DGI as Security Server DGI
    participant CS as Central Server X-Road\n(PKI + Routing + Audit)
    participant SS_ONI as Security Server ONI
    participant ONI as ONI - Identity Registry

    Note over DGI,ONI: Vérification identité pour déclaration fiscale
    DGI->>SS_DGI: Request: GET citizen/{nin}
    SS_DGI->>SS_DGI: Sign + Encrypt (mTLS PKI)
    SS_DGI->>CS: Lookup ONI endpoint + validate cert
    CS-->>SS_DGI: Route + global config
    SS_DGI->>SS_ONI: Encrypted + Signed request
    Note over SS_ONI: Verify DGI certificate\nCheck access rights
    SS_ONI->>ONI: Forward request (internal)
    ONI-->>SS_ONI: Citizen data (minimal)
    SS_ONI->>SS_ONI: Audit log (WORM)
    SS_ONI-->>SS_DGI: Encrypted + Signed response
    SS_DGI->>SS_DGI: Audit log (WORM)
    SS_DGI-->>DGI: Citizen identity data
    Note over CS: Central audit trail\nimmutable pour les 2 parties
```

---

### 7. Workflow Naissance Simple (EC-N01)

```mermaid
flowchart TD
    START([👶 Naissance survenue]) --> DECL
    DECL[Déclarant se présente\nà l'OEC - en ligne ou guichet]
    DECL --> ID_CHK{Identification\ndéclarant\nNIN valide?}
    ID_CHK -->|Non| ID_ERR[Refus + orientation\nvers enrôlement]
    ID_CHK -->|Oui| SAISIE[Saisie données enfant\nnome, prénom, date, lieu, sexe]
    SAISIE --> PARENT_ID[Identification parents\nNIN + acte mariage si applicable]
    PARENT_ID --> TEMOIN[Identification témoins\n≥ 2 personnes majeures]
    TEMOIN --> MED[Attestation médicale?\nFHIR MSPP si dispo]
    MED --> DMN{Règles DMN\nautonmatiques}
    DMN -->|❌ Dates incohérentes| ERR1[Retour déclarant\npour correction]
    DMN -->|❌ Doublon détecté| ERR2[Alerte Data Steward\nbloc workflow]
    DMN -->|✅ Toutes règles OK| OEC_VAL[Validation\nOfficier État Civil]
    OEC_VAL -->|Refus motivé| RECOURS[Notification déclarant\n+ voie de recours]
    OEC_VAL -->|✅ Accepté| NIN_GEN[Génération NIN enfant\nby ONI]
    NIN_GEN --> SIGN[Signature électronique OEC\nPKI SNISID XAdES-LTA]
    SIGN --> DOC[Génération Acte de Naissance\nPDF/A-3 + QR code vérification]
    DOC --> EVENT[Émission événement Kafka\netat-civil.naissance.declaree.v1]
    EVENT --> NOTIF[Notifications cascadées:\n📧 Parents\n🗳️ CEP futur électeur\n🏥 MSPP suivi pédiatrique\n📚 MENFP inscription]
    NOTIF --> END([✅ NIN attribué\nActe disponible < 24h])
    
    style START fill:#22c55e,color:#fff
    style END fill:#22c55e,color:#fff
    style ERR1 fill:#ef4444,color:#fff
    style ERR2 fill:#ef4444,color:#fff
    style RECOURS fill:#f97316,color:#fff
```

---

### 8. Architecture Offline-First

```mermaid
graph TB
    subgraph CENTRAL["🏢 CENTRAL CORE (PaP + Cap-Haïtien)"]
        API_CORE[API Gateway Core]
        KAFKA_CORE[Kafka Central]
        DB_CORE[PostgreSQL Master]
    end

    subgraph EDGE_DEPT["📡 EDGE NODE — Chef-lieu Département"]
        API_EDGE[API Edge local\nk3s/RKE2]
        KAFKA_EDGE[NATS JetStream\nBuffer sync]
        DB_EDGE[SQLite/PostgreSQL local\n50 TB chiffré]
        SYNC_ENG[Sync Engine\nStore & Forward]
    end

    subgraph OFFLINE_KIT["🎒 OFFLINE KIT — Section Communale"]
        TABLET[Tablette Opérateur Android]
        BIO_CAP[Capteurs Biométriques]
        PRINTER[Imprimante portable]
        HSM_USB[HSM USB YubiHSM]
        LOCAL_DB[SQLite local chiffré]
        SOLAR[Panneau solaire 100W\nBatterie 200Wh]
    end

    CENTRAL <-->|4G/VSAT/Fibre\nSync différée| EDGE_DEPT
    EDGE_DEPT <-->|WiFi/BLE/USB\nSync locale| OFFLINE_KIT

    note1["💡 Mode dégradé:\nKit fonctionne 30+ jours\nsans aucune connectivité"]
    note2["🔐 Sécurité:\nLUKS + HSM + Zeroization\nen cas d'intrusion physique"]
```

---

### 9. Cycle de Vie Identité Nationale

```mermaid
stateDiagram-v2
    [*] --> PRE_ENROLLED: Naissance déclarée\n(Event: birth.declared)
    
    PRE_ENROLLED --> ENROLLED: Capture biométrique\n(Event: biometric.captured)
    
    ENROLLED --> ACTIVE: Validation AFIS 1:N\n+ approbation OEC\n(Event: identity.activated)
    
    ACTIVE --> SUSPENDED: Signalement fraude\nou carte perdue\n(Event: identity.suspended)
    
    SUSPENDED --> ACTIVE: Résolution incident\nrevalidation\n(Event: identity.restored)
    
    SUSPENDED --> REVOKED: Fraude confirmée\nou décision judiciaire\n(Event: identity.revoked)
    
    ACTIVE --> DECEASED: Acte décès reçu\n(Event: death.registered)
    
    DECEASED --> ARCHIVED: Après 10 ans\n(Event: identity.archived)
    
    note right of ACTIVE
        État normal
        Accès tous services
        CIN valide
    end note
    
    note right of DECEASED
        CIN invalide immédiatement
        Cascade notifications
        CEP, DGI, OFATMA...
    end note
```

---

## DIAGRAMMES OPÉRATIONNELS

### 10. Escalade Incident P1

```mermaid
flowchart LR
    DETECT[🔍 Détection\nSIEM/SOC\n≤ 5 min] --> TRIAGE
    TRIAGE[⚡ Triage P1\nT1 Analyst\n≤ 10 min] --> AUTO
    AUTO[🤖 SOAR Auto-Response\nIsolation pod/réseau\n≤ 5 min] --> ESCALATE
    ESCALATE[📞 Escalade\nIncident Commander\n≤ 15 min] --> CONTAIN
    CONTAIN{Containment\nréussi < 30 min?}
    CONTAIN -->|✅ Oui| RECOVER[Recovery\nRCA post-mortem]
    CONTAIN -->|❌ Non| CISO_ALERT[CISO notification\nurgence]
    CISO_ALERT --> FULL_IR[Full IR Team\nactivée]
    FULL_IR --> DR_CHECK{Impact\nInfrastructure\nglobale?}
    DR_CHECK -->|Oui| CNN_ALERT[Cellule crise\nCNN activée]
    CNN_ALERT --> DR_DEPLOY[Déploiement PRA\nbasculement DC]
    DR_CHECK -->|Non| FULL_IR
    
    style DETECT fill:#ef4444,color:#fff
    style CNN_ALERT fill:#7c3aed,color:#fff
    style DR_DEPLOY fill:#1d4ed8,color:#fff
```

---

*Index des Diagrammes — SNISID Phase 0*  
*SNISID — République d'Haïti — Classification : Usage Gouvernemental*
