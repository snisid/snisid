# SNISID — INDEX MAÎTRE DES MASTER PROMPTS
## Modules MP-19 à MP-56 — Système National d'Identification Sécurisée et Intégrée d'Haïti

```
Classification   : RESTREINT / USAGE OFFICIEL
Document         : SNISID-INDEX-MP-19-56
Version          : 1.0.0
Date             : Juin 2026
Modules couverts : 38 modules (MP-19 à MP-56)
Modules existants: MP-15 (FOVeS/SIV), MP-16 (LAPI), MP-17 (FPR), MP-18 (SIVC-HT)
Stack technique  : Go 1.22, Python 3.12, TypeScript, PostgreSQL 16, CockroachDB,
                   Redis 7, Neo4j 5, Milvus 2, ClickHouse 24, Kafka 3, RKE2/K8s
```

---

## 1. CARTOGRAPHIE COMPLÈTE DES 38 MODULES

### DOMAINE 1 — Justice Criminelle et Identification (MP-19 à MP-23)

| MP    | Code         | Module                                      | Priorité  | Port  | Dépend de               |
|-------|--------------|---------------------------------------------|-----------|-------|-------------------------|
| MP-19 | AFIS-HT      | Automated Fingerprint Identification System | 🔴 Critique | 8091  | FIR-HT, FPR-HT, BIO-ADN|
| MP-20 | FIR-HT       | Fichier Individuel des Renseignements       | 🔴 Critique | 8093  | AFIS-HT, SIPEP-HT       |
| MP-21 | SIPEP-HT     | Système Information Pénitentiaire           | 🔴 Critique | 8092  | FIR-HT, AFIS-HT         |
| MP-22 | RDEP-HT      | Registre Déportés et Extradés               | 🟠 Haute    | 8094  | FIR-HT, GANG-HT         |
| MP-23 | OPR-HT       | Ordonnances de Protection                   | 🟡 Normale  | 8096  | FIR-HT, FPR-HT          |

### DOMAINE 2 — Organisations Criminelles et Gangs (MP-24 à MP-28)

| MP    | Code         | Module                                      | Priorité  | Port  | Dépend de               |
|-------|--------------|---------------------------------------------|-----------|-------|-------------------------|
| MP-24 | GANG-HT      | Registre Organisations Criminelles          | 🔴 Critique | 8095  | CHEF-HT, TERR-HT        |
| MP-25 | CHEF-HT      | Leaders et Membres Identifiés               | 🔴 Critique | 8097  | GANG-HT, AFIS-HT        |
| MP-26 | TERR-HT      | Cartographie Territoires Contrôlés          | 🔴 Critique | 8098  | GANG-HT, PostGIS        |
| MP-27 | SANC-HT      | Interface Sanctions ONU/OFAC/UE             | 🟠 Haute    | 8100  | GANG-HT, CHEF-HT        |
| MP-28 | RESO-HT      | Analyse Réseaux Criminels                   | 🟠 Haute    | 8101  | GANG-HT, Neo4j, GNN     |

### DOMAINE 3 — Armes et Munitions (MP-29 à MP-32)

| MP    | Code         | Module                                      | Priorité  | Port  | Dépend de               |
|-------|--------------|---------------------------------------------|-----------|-------|-------------------------|
| MP-29 | SIAR-HT      | Armes à Feu Légales et Licences             | 🔴 Critique | 8102  | BIAR-HT, iARMS          |
| MP-30 | BIAR-HT      | Armes Illicites — Interface iARMS INTERPOL  | 🔴 Critique | 8103  | SIAR-HT, ATF eTRACE     |
| MP-31 | EXPL-HT      | Explosifs, IED et Munitions                 | 🟠 Haute    | 8104  | BIAR-HT, BIO-ADN        |
| MP-32 | TRAF-AR      | Routes Trafic d'Armes                       | 🟠 Haute    | 8105  | BIAR-HT, MAR-HT         |

### DOMAINE 4 — Frontières, Maritime et Aviation (MP-33 à MP-38)

| MP    | Code         | Module                                      | Priorité  | Port  | Dépend de               |
|-------|--------------|---------------------------------------------|-----------|-------|-------------------------|
| MP-33 | SIFR-HT      | Frontières et Routes Terrestres             | 🔴 Critique | 8106  | BLKL-HT, SLTD-HT, FPR  |
| MP-34 | MAR-HT       | Surveillance Maritime Nationale             | 🔴 Critique | 8107  | PORT-HT, AIS, INTERPOL  |
| MP-35 | SLTD-HT      | Documents Voyage Perdus/Volés (INTERPOL)   | 🔴 Critique | 8108  | SIFR-HT, ONI            |
| MP-36 | BLKL-HT      | Liste Noire Entrées/Sorties Territoire      | 🔴 Critique | 8110  | FPR-HT, SANC-HT         |
| MP-37 | AERO-HT      | Aéronefs Illicites et Pistes Clandestines  | 🟠 Haute    | 8109  | TRAF-AR, MAR-HT         |
| MP-38 | PORT-HT      | Sécurité Portuaire et Ciblage Conteneurs   | 🟠 Haute    | 8111  | MAR-HT, BLAR-HT         |

### DOMAINE 5 — Renseignement Financier (MP-39 à MP-42)

| MP    | Code         | Module                                      | Priorité  | Port  | Dépend de               |
|-------|--------------|---------------------------------------------|-----------|-------|-------------------------|
| MP-39 | UCREF-INT    | Interface UCREF / FIU National              | 🔴 Critique | 8112  | GANG-HT, SANC-HT        |
| MP-40 | BLAN-HT      | Transactions Suspectes et Blanchiment       | 🟠 Haute    | 8115  | UCREF-INT, GANG-HT      |
| MP-41 | EXTORS-HT    | Extorsions, Péages Illicites et Rançons     | 🟠 Haute    | 8116  | GANG-HT, UCREF-INT      |
| MP-42 | CRYPT-HT     | Cryptomonnaies à Usage Criminel             | 🟡 Normale  | 8117  | UCREF-INT, BLAN-HT      |

### DOMAINE 6 — Personnes Disparues et Traite (MP-43 à MP-47)

| MP    | Code         | Module                                      | Priorité  | Port  | Dépend de               |
|-------|--------------|---------------------------------------------|-----------|-------|-------------------------|
| MP-43 | DIPE-HT      | Registre Personnes Disparues                | 🔴 Critique | 8118  | AFIS-HT, BIO-ADN, RVIN  |
| MP-44 | TRAIT-HT     | Traite des Personnes et Migration           | 🟠 Haute    | 8122  | MAR-HT, SIFR-HT         |
| MP-45 | ENFL-HT      | Enfants Disparus et à Risque                | 🟠 Haute    | 8119  | DIPE-HT, GANG-HT, IBESR |
| MP-46 | DPIDE-HT     | Déplacés Internes (IDPs)                    | 🟠 Haute    | 8121  | SIGDC-HT, SIGEO-HT      |
| MP-47 | VICT-HT      | Registre Victimes de Crimes Graves          | 🟡 Normale  | 8123  | FIR-HT, BIO-ADN, RVIN   |

### DOMAINE 7 — Géo-Intelligence et Gestion de Catastrophes (MP-48 à MP-52)

| MP    | Code         | Module                                      | Priorité  | Port  | Dépend de               |
|-------|--------------|---------------------------------------------|-----------|-------|-------------------------|
| MP-48 | SIGEO-HT     | Géo-Intelligence Criminelle et Cartographie | 🔴 Critique | 8125  | TERR-HT, PostGIS, H3    |
| MP-49 | SIGDC-HT     | Gestion des Désastres Civils                | 🔴 Critique | 8126  | SIGEO-HT, DPIDE-HT, ADN |
| MP-50 | RVIN-HT      | Victimes Non Identifiées                    | 🟠 Haute    | 8120  | AFIS-HT, BIO-ADN, DIPE  |
| MP-51 | MVSM-HT      | Surveillance Rassemblements de Masse        | 🟡 Normale  | 8127  | SIGEO-HT, TERR-HT       |
| MP-52 | SISAL-HT     | Alerte Précoce Multi-Risques                | 🟠 Haute    | 8128  | SIGDC-HT, SIGEO-HT      |

### DOMAINE 8 — Infrastructure, Cyber et Gouvernance (MP-53 à MP-56)

| MP    | Code         | Module                                      | Priorité  | Port  | Dépend de               |
|-------|--------------|---------------------------------------------|-----------|-------|-------------------------|
| MP-53 | CORR-HT      | Anti-Corruption et Intégrité PNH            | 🔴 Critique | 8130  | FIR-HT, GANG-HT, UCREF  |
| MP-54 | SIPCI-HT     | Protection Infrastructures Critiques        | 🟠 Haute    | 8131  | SIGEO-HT, GANG-HT       |
| MP-55 | CYBRE-HT     | Cybercriminalité Nationale                  | 🟠 Haute    | 8132  | CRYPT-HT, UCREF-INT     |
| MP-56 | ONG-HT       | Registre ONGs et Acteurs Humanitaires       | 🟡 Normale  | 8133  | SIGEO-HT, SANC-HT       |

---

## 2. GRAPHE DE DÉPENDANCES INTER-MODULES

```
SNISID IDENTITÉ (CORE)
│
├── AFIS-HT (MP-19) ──────────────────┐
│       │                             │
│   FIR-HT (MP-20) ───────────────── │──────────────┐
│       │                             │              │
│   SIPEP-HT (MP-21) ──────── FPR-HT (MP-17)        │
│       │                             │              │
│   RDEP-HT (MP-22) ─────── GANG-HT (MP-24) ──────┐ │
│                               │                  │ │
│   CHEF-HT (MP-25) ────────── │                  │ │
│   TERR-HT (MP-26) ─ PostGIS ─┤                  │ │
│   SANC-HT (MP-27) ─ OFAC/ONU ┤                  │ │
│   RESO-HT (MP-28) ─ Neo4j ───┘                  │ │
│                                                  │ │
├── SIAR-HT (MP-29) ─── iARMS ──────────────────  │ │
│   BIAR-HT (MP-30) ─── ATF eTRACE ────────────── │ │
│   EXPL-HT (MP-31) ───────────────────────────── │ │
│   TRAF-AR (MP-32) ───────────────────────────── │ │
│                                                  │ │
├── SIFR-HT (MP-33) ─── INTERPOL PISCES ────────  │ │
│   MAR-HT  (MP-34) ─── AIS / INTERPOL SVD ─────  │ │
│   SLTD-HT (MP-35) ─── INTERPOL SLTD ──────────  │ │
│   BLKL-HT (MP-36) ─── Redis Hotlist ───────────  │ │
│   AERO-HT (MP-37) ─── FAA / OACI ─────────────  │ │
│   PORT-HT (MP-38) ─── CBP / ISPS ─────────────  │ │
│                                                  │ │
├── UCREF-INT (MP-39) ─ Egmont / FATF ──────────  │ │
│   BLAN-HT  (MP-40) ──────────────────────────── │ │
│   EXTORS-HT(MP-41) ──────────────────────────── │ │
│   CRYPT-HT (MP-42) ─ Chainalysis ─────────────  │ │
│                                                  │ │
├── DIPE-HT (MP-43) ─── INTERPOL MP ─────────── ◄─┘ │
│   TRAIT-HT(MP-44) ─── IOM ─────────────────────    │
│   ENFL-HT (MP-45) ─── NCMEC / ICSE ───────────    │
│   DPIDE-HT(MP-46) ─── IOM DTM ────────────────    │
│   VICT-HT (MP-47) ─────────────────────────────    │
│                                                     │
├── SIGEO-HT(MP-48) ─── PostGIS + H3 + MapLibre ─    │
│   SIGDC-HT(MP-49) ─── USGS / NHC / OCHA ──────  ◄─┘
│   RVIN-HT (MP-50) ─── ADN + AFIS + DVI ────────
│   MVSM-HT (MP-51) ─────────────────────────────
│   SISAL-HT(MP-52) ─── SMS / FCM / Radio ───────
│
├── CORR-HT (MP-53) ─── IGPNH ISOLÉ ─────────────
│   SIPCI-HT(MP-54) ─────────────────────────────
│   CYBRE-HT(MP-55) ─── MISP / CONATEL ──────────
└── ONG-HT  (MP-56) ─── OCHA / ULCC ────────────
```

---

## 3. TABLEAU DE BORD D'IMPLÉMENTATION

### Phase 1 — Fondations judiciaires (Semaines 1-8)

| Ordre | Module   | Justification                                    |
|-------|----------|--------------------------------------------------|
| 1     | AFIS-HT  | Identification criminelle — base de tout le reste|
| 2     | FIR-HT   | Casier judiciaire — lien légal fondamental       |
| 3     | SIPEP-HT | Prisons — 12,000 détenus sans registre           |
| 4     | GANG-HT  | Menace N°1 à la sécurité nationale               |
| 5     | CHEF-HT  | Leaders gangs — cibles opérationnelles           |
| 6     | SIAR-HT  | Armes légales — base du contrôle                 |
| 7     | BIAR-HT  | Armes illicites — 500,000 en circulation         |
| 8     | CORR-HT  | Intégrité — protéger SNISID lui-même             |

### Phase 2 — Frontières et renseignement (Semaines 9-16)

| Ordre | Module    | Justification                                    |
|-------|-----------|--------------------------------------------------|
| 9     | SIFR-HT   | Contrôle frontière HT-DR — priorité douanière    |
| 10    | SLTD-HT   | Documents volés — fraud passeports               |
| 11    | BLKL-HT   | Liste noire — effectivité frontière             |
| 12    | MAR-HT    | Maritime — drogue et armes par la mer            |
| 13    | UCREF-INT | Blanchiment — économie criminelle gangs          |
| 14    | RDEP-HT   | Déportés — 400 Mawozo composé de déportés USA    |
| 15    | TERR-HT   | Cartographie gang — planification opérationnelle |
| 16    | SANC-HT   | Sanctions ONU/OFAC — obligation internationale   |

### Phase 3 — Protection populations (Semaines 17-24)

| Ordre | Module    | Justification                                    |
|-------|-----------|--------------------------------------------------|
| 17    | DIPE-HT   | Disparus — kidnapping endémique                  |
| 18    | SIGEO-HT  | Cartographie — image opérationnelle commune      |
| 19    | SIGDC-HT  | Catastrophes — risque sismique permanent         |
| 20    | ENFL-HT   | Enfants — 225,000 restaveks, recrutement gangs   |
| 21    | SISAL-HT  | Alertes — population et agents terrain           |
| 22    | SIPCI-HT  | Infra critiques — EDH, routes, ports             |
| 23    | EXTORS-HT | Péages et rançons — économie gang documentée     |
| 24    | PORT-HT   | Conteneurs — transit cocaïne documenté           |

### Phase 4 — Maturité et spécialisation (Semaines 25-32)

| Ordre | Module    | Justification                             |
|-------|-----------|-------------------------------------------|
| 25    | RESO-HT   | Analyse réseau — GNN PyTorch              |
| 26    | TRAIT-HT  | Traite — migrants et exploitation         |
| 27    | DPIDE-HT  | IDPs — 580,000 déplacés                   |
| 28    | RVIN-HT   | Victimes non identifiées                  |
| 29    | VICT-HT   | Documentation crimes graves               |
| 30    | BLAN-HT   | Dossiers blanchiment                      |
| 31    | EXPL-HT   | Explosifs et IED                          |
| 32    | TRAF-AR   | Routes trafic armes                       |

### Phase 5 — Complétion (Semaines 33-40)

| Ordre | Module    | Justification                             |
|-------|-----------|-------------------------------------------|
| 33    | CYBRE-HT  | Fraudes MonCash — cybercriminalité        |
| 34    | CRYPT-HT  | Cryptomonnaies — rançons et blanchiment   |
| 35    | AERO-HT   | Pistes clandestines drogue                |
| 36    | MVSM-HT   | Rassemblements — peyi lòk                 |
| 37    | OPR-HT    | Ordonnances protection                    |
| 38    | ONG-HT    | Acteurs humanitaires — 10,000 ONGs        |

---

## 4. INFRASTRUCTURE PARTAGÉE

### Bases de données — Attribution des schémas

```
PostgreSQL 16 (Principal) :
  snisid_afis     → AFIS-HT (MP-19)
  snisid_fir      → FIR-HT (MP-20)
  snisid_sipep    → SIPEP-HT (MP-21)
  snisid_rdep     → RDEP-HT (MP-22)
  snisid_opr      → OPR-HT (MP-23)
  snisid_gang     → GANG-HT (MP-24)
  snisid_chef     → CHEF-HT (MP-25)
  snisid_terr     → TERR-HT (MP-26) + PostGIS
  snisid_sanc     → SANC-HT (MP-27)
  snisid_sivc     → SIVC-HT (MP-18)
  snisid_siar     → SIAR-HT (MP-29)
  snisid_biar     → BIAR-HT (MP-30)
  snisid_expl     → EXPL-HT (MP-31)
  snisid_trafar   → TRAF-AR (MP-32)
  snisid_sifr     → SIFR-HT (MP-33)
  snisid_mar      → MAR-HT (MP-34)
  snisid_sltd     → SLTD-HT (MP-35)
  snisid_blkl     → BLKL-HT (MP-36)
  snisid_aero     → AERO-HT (MP-37)
  snisid_port     → PORT-HT (MP-38)
  snisid_ucref    → UCREF-INT (MP-39)
  snisid_blan     → BLAN-HT (MP-40)
  snisid_extors   → EXTORS-HT (MP-41)
  snisid_crypt    → CRYPT-HT (MP-42)
  snisid_dipe     → DIPE-HT (MP-43)
  snisid_trait    → TRAIT-HT (MP-44)
  snisid_enfl     → ENFL-HT (MP-45)
  snisid_dpide    → DPIDE-HT (MP-46)
  snisid_vict     → VICT-HT (MP-47)
  snisid_sigeo    → SIGEO-HT (MP-48) + PostGIS
  snisid_sigdc    → SIGDC-HT (MP-49)
  snisid_rvin     → RVIN-HT (MP-50)
  snisid_mvsm     → MVSM-HT (MP-51)
  snisid_sisal    → SISAL-HT (MP-52)
  snisid_corr     → CORR-HT (MP-53) [ISOLÉ]
  snisid_sipci    → SIPCI-HT (MP-54)
  snisid_cybre    → CYBRE-HT (MP-55)
  snisid_ong      → ONG-HT (MP-56)

Milvus (Vecteurs biométriques) :
  afis_fingerprints  → AFIS-HT
  face_embeddings    → SNISID CORE

Neo4j (Graphe criminologique) :
  Graphe central RESO-HT (MP-28)
  Alimenté par GANG, CHEF, SIVC, FIR, BLAN, DIPE

Redis (Hotlists temps réel) :
  sivc:plate:*       → SIVC-HT (MP-18)
  blkl:person:*      → BLKL-HT (MP-36)
  sltd:doc:*         → SLTD-HT (MP-35)
  fpr:warrant:*      → FPR-HT (MP-17)

ClickHouse (Analytiques) :
  Toutes les tables de faits *_facts/* _events
```

### Kafka Topics centraux

```
sivc.alerts.*       → Alertes véhiculaires
gang.*.created      → Gangs, membres, incidents
fir.record.*        → Casiers judiciaires
sipep.inmate.*      → Mouvements pénitentiaires
sifr.crossing.*     → Passages frontaliers
mar.vessel.*        → Événements maritimes
sigeo.incidents.*   → Géolocalisation unifiée
sigdc.disaster.*    → Urgences et désastres
sisal.alert.*       → Alertes précoces multi-risques
corr.behavioral.*   → Anomalies intégrité (isolé)
*.audit             → Audit trail immuable (tous modules)
```

---

## 5. STANDARDS D'IMPLÉMENTATION OBLIGATOIRES

```
1. Authentification    : SPIFFE/SPIRE (identité service), JWT RS256 (utilisateurs)
2. Transport           : TLS 1.3 obligatoire, mTLS inter-services
3. Secrets             : HashiCorp Vault (HSM Luna en production)
4. Conteneurs          : Docker multi-stage, images Alpine hardened
5. Orchestration       : RKE2/Kubernetes (SUSE Rancher) — pas de K8s public cloud
6. Migrations SQL      : golang-migrate — séquentielles et réversibles
7. Logs                : Structurés JSON, niveau INFO en prod, WARN sur perfs
8. Métriques           : Prometheus + Grafana — SLO définis par module
9. Tests               : coverage ≥ 80% (unité), ≥ 60% (intégration)
10. Audit              : Chaque accès données sensibles → Kafka `*.audit` immuable
11. HSM                : Luna Network HSM 7 pour toutes clés cryptographiques
12. Biométrie          : ISO/IEC 19794 (empreintes), ISO/IEC 19785 (BioAPI)
13. RLS PostgreSQL     : Activé sur toutes tables contenant données personnelles
14. Backup             : RPO ≤ 1h, RTO ≤ 4h, chiffré, hors-site, testé mensuellement
15. Classification     : Chaque document tagué RESTREINT / SECRET / TOP SECRET
```

---

## 6. CONVENTIONS DE CODE

```go
// Pattern service SNISID standard
type {Module}Service struct {
    repo    {Module}Repository
    kafka   EventPublisher
    logger  *zap.Logger
    metrics *prometheus.CounterVec
}

// Pattern handler REST standard
// GET /api/v1/{module}/{resource}
// POST /api/v1/{module}/{resource}
// PATCH /api/v1/{module}/{resource}/:id
// DELETE /api/v1/{module}/{resource}/:id (avec motif obligatoire)

// Réponse erreur standard
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    TraceID string `json:"trace_id"`
}

// Réponse pagination standard
type PaginatedResponse[T any] struct {
    Data       []T    `json:"data"`
    Total      int64  `json:"total"`
    Page       int    `json:"page"`
    PageSize   int    `json:"page_size"`
    HasMore    bool   `json:"has_more"`
}
```

---

## 7. CONTACTS ET INTÉGRATIONS INTERNATIONALES

| Organisation       | Module(s)              | Type d'intégration          | Protocole        |
|--------------------|------------------------|-----------------------------|------------------|
| INTERPOL BCN HTI   | SIVC, AFIS, SLTD, iARMS| BCN Port-au-Prince liaison  | I-24/7 / MAPI    |
| FBI (USA)          | AFIS-HT, RDEP-HT       | IAFIS, eCRIMS               | API bilatérale   |
| ATF (USA)          | SIAR-HT, BIAR-HT       | eTRACE                      | API bilatérale   |
| DEA (USA)          | TRAF-AR, MAR-HT        | Liaison opérationnelle      | Canal sécurisé   |
| DHS/ICE (USA)      | RDEP-HT, SIFR-HT       | Déportations, biométrie     | API bilatérale   |
| JIATF-South        | MAR-HT                 | Opérations maritimes conjointes | Canal militaire|
| UNODC              | TRAF-AR, TRAIT-HT      | Rapports et analyses        | API REST          |
| OFAC (USA)         | SANC-HT                | SDN List                    | XML Feed          |
| ONU CSNU           | SANC-HT                | Liste 2653                  | XML Feed          |
| Egmont Group       | UCREF-INT              | FIU Network                 | GoAML / REST      |
| OIM                | TRAIT-HT, DPIDE-HT     | DTM, Counter-Trafficking    | API REST          |
| OCHA               | SIGDC-HT, DPIDE-HT     | ReliefWeb, 5W Matrix        | API REST          |
| Chainalysis        | CRYPT-HT               | Blockchain analysis         | API KYT           |
| INTERPOL ICSE      | ENFL-HT                | Child exploitation DB       | I-24/7            |
| NCMEC              | ENFL-HT, DIPE-HT       | Missing Children            | API REST          |
| USGS               | SIGDC-HT, SISAL-HT     | Earthquake feeds            | FDSN REST         |
| NHC (NOAA)         | SIGDC-HT, SISAL-HT     | Hurricane advisories        | RSS/JSON          |

---

## 8. STATISTIQUES DU PROJET SNISID COMPLET

```
Modules totaux              : 42 (MP-15 à MP-56)
Fichiers à créer (estimé)   : 1,800+
Migrations SQL (estimé)     : 280+
Microservices Go            : 42+
Services Python ML          : 8+
Tables PostgreSQL (estimé)  : 380+
Indexes créés (estimé)      : 1,100+
Endpoints REST (estimé)     : 520+
Topics Kafka                : 95+
Collections Milvus          : 4
Graphes Neo4j               : 1 (central RESO-HT)
Vues ClickHouse             : 45+
Intégrations internationales: 17 organisations
Ports de service            : 8080-8133
```

---

*SNISID — Système National d'Identification Sécurisée et Intégrée d'Haïti*
*Index Maître MP-19 à MP-56 — Version 1.0.0 — Juin 2026*
*République d'Haïti — Ministère de la Justice et de la Sécurité Publique*
*Classification : RESTREINT / USAGE OFFICIEL SEULEMENT*
