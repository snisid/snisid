# SNISID NATIONAL APPLICATION ECOSYSTEM ARCHITECTURE
## Document d'Architecture Applicative Nationale

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-ARCH-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |
| Classification | CONFIDENTIEL |

---

## 1. VISION ARCHITECTURALE

### 1.1 La Couche Visible du Gouvernement Numérique Haïtien

```
┌─────────────────────────────────────────────────────────────┐
│                 CITIZEN EXPERIENCE LAYER                      │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │  Super   │ │  Wallet  │ │  Mobile  │ │  Kiosk   │       │
│  │   App    │ │   App    │ │  Field   │ │   App    │       │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │
├─────────────────────────────────────────────────────────────┤
│                 GOVERNMENT OPERATIONS LAYER                   │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │   Gov    │ │  Admin   │ │  Police  │ │ Justice  │       │
│  │Super App │ │  Portal  │ │   App    │ │   App    │       │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │
├─────────────────────────────────────────────────────────────┤
│                 NOTIFICATION & COMMUNICATION LAYER             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │        National Notification Platform                │    │
│  │   SMS  •  Email  •  Push  •  In-App Secure Messages  │    │
│  └─────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────┤
│                 CORE SERVICES LAYER                           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │ Identity │ │  Civil   │ │   Case   │ │ Payment  │       │
│  │  Hub     │ │ Registry │ │  Mgmt    │ │ Gateway  │       │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │
├─────────────────────────────────────────────────────────────┤
│                 OFFLINE & SYNC LAYER                          │
│  ┌─────────────────────────────────────────────────────┐    │
│  │        Offline Mobile Sync Engine                     │    │
│  │   Caching  •  Conflict Res •  Secure Sync •  Buffer  │    │
│  └─────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────┤
│                 SECURITY & GOVERNANCE LAYER                   │
│  ┌─────────────────────────────────────────────────────┐    │
│  │        Application Security Framework                │    │
│  │   Auth • Attestation • Anti-tamper • Cert Pinning   │    │
│  └─────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────┤
│                 OBSERVABILITY LAYER                           │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Metrics (Prometheus) • Analytics (OpenTelemetry)   │    │
│  │  Logs (Loki) • Dashboards (Grafana)                 │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 Principes Fondamentaux

| Principe | Description |
|----------|-------------|
| **Offline-First** | Toutes les applications doivent fonctionner sans connexion internet permanente |
| **Souveraineté** | Hébergement national, contrôle total des données |
| **Sécurité by Design** | Sécurité intégrée dès la conception |
| **Multilingue** | Français, Créole Haïtien, Anglais |
| **Auditabilité** | Toutes les actions tracées et horodatées |
| **Mobile-First** | Conçues d'abord pour le mobile, adaptées au desktop |
| **Accessibilité** | Conformes WCAG 2.1 AA minimum |

---

## 2. DOMAINES APPLICATIFS

### 2.1 Citizen Applications — Population

| Application | Public | Usage | Offline |
|-------------|--------|-------|---------|
| Citizen Super App | Tous citoyens | Identité, État civil, Services | Partiel |
| Digital Identity Wallet | Tous citoyens | Certificats, Signatures, QR | Complet |
| Citizen Self-Service | Citoyens connectés | Demandes, Paiements, Suivi | Partiel |

### 2.2 Government Applications — Agents

| Application | Public | Usage | Offline |
|-------------|--------|-------|---------|
| Government Super App | Tous agents | Approbations, Cas, Messagerie | Partiel |
| National Admin Portal | Administrateurs | Gestion, Gouvernance, Audit | Non |
| Civil Registry App | Agents d'état civil | Enregistrement, Certificats | Partiel |
| Workflow Manager | Agents décideurs | Approbations, Escalades | Partiel |

### 2.3 Field Applications — Terrain

| Application | Public | Usage | Offline |
|-------------|--------|-------|---------|
| Mobile Field App | Agents terrain | Enrôlement, Vérification | Complet |
| Offline Enrollment Kit | Agents déploiement | Inscription hors-ligne | Complet |
| Biometric Field Tool | Agents identification | Biométrie terrain | Complet |
| Event Buffer App | Tous terrains | Collecte et synchronisation | Complet |

### 2.4 Super Apps — Centralisation

| Application | Rôle |
|-------------|------|
| Citizen Super App | Point d'entrée unique citoyen |
| Government Super App | Point d'entrée unique agent |
| National Admin Portal | Point d'entrée administrateurs |

### 2.5 Security & Justice Applications

| Application | Public | Usage | Sécurité |
|-------------|--------|-------|----------|
| Police Criminal Records | Police | Cas criminels | Ultra |
| Judicial Workflow | Justice | Dossiers judiciaires | Ultra |
| Penitentiary Operations | Prison | Gestion pénitentiaire | Ultra |
| DCPJ Investigations | DCPJ | Enquêtes | Ultra |
| Immigration Border Control | Immigration | Contrôle frontalier | Ultra |

---

## 3. ECOSYSTEME APPLICATIF COMPLET

```
                        ┌─────────────────────┐
                        │   NATIONAL ADMIN     │
                        │      PORTAL          │
                        └──────────┬──────────┘
                                   │
            ┌──────────────────────┼──────────────────────┐
            │                      │                      │
   ┌────────▼────────┐   ┌────────▼────────┐   ┌────────▼────────┐
   │    CITIZEN      │   │   GOVERNMENT    │   │  POLICE/JUSTICE │
   │   SUPER APP     │   │   SUPER APP     │   │     APPS        │
   │                 │   │                 │   │                 │
   │ • ID Wallet     │   │ • Workflow      │   │ • Criminal      │
   │ • Birth Cert    │   │ • Case Mgmt     │   │ • Judicial      │
   │ • Civil Registry│   │ • ID Verify     │   │ • Penitentiary  │
   │ • QR Verify     │   │ • Escalation    │   │ • DCPJ          │
   │ • Notifications │   │ • Messaging     │   │ • Immigration   │
   └────────┬────────┘   └────────┬────────┘   └────────┬────────┘
            │                      │                      │
   ┌────────▼──────────────────────▼──────────────────────▼────────┐
   │                    NOTIFICATION PLATFORM                       │
   │           SMS  •  Email  •  Push  •  Secure In-App            │
   └──────────────────────────────┬───────────────────────────────┘
                                  │
   ┌──────────────────────────────▼───────────────────────────────┐
   │                   CORE SERVICES LAYER                         │
   │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
   │  │ Identity │ │  Civil   │ │   Case   │ │ Payment  │       │
   │  │  Service │ │ Registry │ │  Mgmt    │ │ Gateway  │       │
   │  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │
   └──────────────────────────────┬───────────────────────────────┘
                                  │
   ┌──────────────────────────────▼───────────────────────────────┐
   │                    OFFLINE SYNC ENGINE                        │
   │       Caching  •  Delay  •  Conflict  •  Secure Sync         │
   └──────────────────────────────┬───────────────────────────────┘
                                  │
   ┌──────────────────────────────▼───────────────────────────────┐
   │              MULTI-CHANNEL ACCESS MODEL                       │
   │   Android  •  iOS  •  Web  •  Kiosks  •  Offline Devices     │
   └──────────────────────────────────────────────────────────────┘
```

---

## 4. OFFLINE-FIRST ARCHITECTURE

### 4.1 Modes de Fonctionnement

| Mode | Connectivité | Fonctionnalités |
|------|-------------|-----------------|
| **Online** | Internet | Toutes les fonctionnalités |
| **Degraded** | Intermittent | Fonctionnalités essentielles, sync différée |
| **Offline** | Aucune | Opérations terrain, buffer local |
| **Reconnect** | Retour connexion | Synchronisation, résolution conflits |

### 4.2 Stratégie de Cache

```
┌─────────────────────────────────────┐
│          APPLICATION LAYER           │
├─────────────────────────────────────┤
│  ┌───────────────────────────────┐  │
│  │    Offline Status Indicator   │  │
│  └───────────────────────────────┘  │
├─────────────────────────────────────┤
│          DATA ACCESS LAYER           │
├─────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐        │
│  │  Local   │  │  Remote  │        │
│  │  Cache   │  │   API    │        │
│  └────┬─────┘  └────┬─────┘        │
│       │              │              │
│  ┌────▼──────────────▼─────┐        │
│  │     Sync Coordinator    │        │
│  └─────────────────────────┘        │
├─────────────────────────────────────┤
│  ┌───────────────────────────────┐  │
│  │   Conflict Resolution Engine  │  │
│  └───────────────────────────────┘  │
├─────────────────────────────────────┤
│  ┌───────────────────────────────┐  │
│  │   Secure Local Storage        │  │
│  │   (Encrypted SQLite/CouchDB) │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

### 4.3 Cache Tiers

| Tier | Contenu | Priorité | Taille Max |
|------|---------|----------|------------|
| **L1** | Données utilisateur, profils | Critique | 50 MB |
| **L2** | Données de référence, catalogues | Haute | 200 MB |
| **L3** | Contenus, notifications | Normale | 500 MB |

---

## 5. SÉCURITÉ APPLICATIVE

### 5.1 Périmètre de Sécurité

| Couche | Protection |
|--------|-----------|
| **Application** | Anti-tampering, Certificate pinning, App attestation |
| **Transport** | TLS 1.3, Mutual TLS, Certificate pinning |
| **Données** | Encryption AES-256, Keychain/Keystore |
| **Authentification** | Biometric + PIN + OTP (MFA obligatoire) |
| **Session** | JWT avec rotation, Fingerprint device |
| **Audit** | Toutes les actions horodatées et signées |

### 5.2 Classification des Applications

| Classe | Exemples | Requis |
|--------|----------|--------|
| **Standard** | Citizen App, Notifications | MFA optionnel |
| **Sensible** | Government App, Wallet | MFA obligatoire |
| **Critique** | Police, Justice, Admin | MFA + Attestation + Audit temps réel |

---

## 6. COMMUNICATION INTER-APPLICATIONS

### 6.1 Flux de Données

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│ Citizen  │────▶│  Gov    │────▶│  Admin  │
│  App     │     │  App    │     │  Portal │
└────┬────┘     └────┬────┘     └────┬────┘
     │               │               │
     └───────────────┼───────────────┘
                     │
            ┌────────▼────────┐
            │   Identity Hub  │
            └────────┬────────┘
                     │
            ┌────────▼────────┐
            │  Notification   │
            │   Platform      │
            └─────────────────┘
```

### 6.2 Protocoles

| Protocole | Usage | Sécurité |
|-----------|-------|----------|
| REST API | Opérations synchrones | TLS + JWT |
| WebSocket | Notifications temps réel | TLS + Token |
| Message Queue | Opérations asynchrones | TLS + Signature |
| Secure Sync | Synchronisation offline | Encryption bout-en-bout |

---

## 7. MULTILINGUE

### 7.1 Langues Supportées

| Langue | Code | Priorité |
|--------|------|----------|
| Français | fr | Primaire |
| Créole Haïtien | ht | Nationale |
| Anglais | en | Internationale |

### 7.2 Architecture de Traduction

```
┌─────────────────────────────────────┐
│      i18n Translation Engine         │
├─────────────────────────────────────┤
│  ┌──────────┐ ┌──────────┐          │
│  │  Static  │ │ Dynamic  │          │
│  │   Texts  │ │  Content │          │
│  └──────────┘ └──────────┘          │
├─────────────────────────────────────┤
│  ┌───────────────────────────────┐  │
│  │  Translation Management       │  │
│  └───────────────────────────────┘  │
├─────────────────────────────────────┤
│  ┌──────────┐ ┌──────────┐          │
│  │  JSON    │ │  Server  │          │
│  │  Locales │ │  Render  │          │
│  └──────────┘ └──────────┘          │
└─────────────────────────────────────┘
```

---

## 8. MULTI-CANAL

| Canal | Citizen | Government | Police | Field |
|-------|---------|------------|--------|-------|
| **Android** | ✅ | ✅ | ✅ | ✅ |
| **iOS** | ✅ | ✅ | ✅ | ✅ |
| **Web** | ✅ | ✅ | ✅ | ❌ |
| **Kiosks** | ✅ | ❌ | ❌ | ❌ |
| **Offline Devices** | ✅ | ❌ | ❌ | ✅ |

---

## 9. INTERFACES & INTÉGRATIONS

### 9.1 APIs Externes

| API | Fournisseur | Usage |
|-----|-------------|-------|
| ID Verification API | SNISID | Vérification identité |
| Civil Registry API | SNISID | État civil |
| Case Management API | SNISID | Gestion de cas |
| Notification API | SNISID | Notifications |
| Payment Gateway | Partenaire | Paiements |

### 9.2 APIs Internes

| API | Usage |
|-----|-------|
| Sync API | Synchronisation offline |
| Audit API | Traçabilité |
| Governance API | Gouvernance applicative |
| Observability API | Monitoring |

---

## 10. GOUVERNANCE APPLICATIVE

### 10.1 Cycle de Vie

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│  Design  │───▶│  Review  │───▶│  Build   │───▶│  Audit   │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
                                                     │
┌──────────┐    ┌──────────┐    ┌──────────┐         │
│  Retire  │◀───│ Monitor  │◀───│  Deploy  │◀────────┘
└──────────┘    └──────────┘    └──────────┘
```

### 10.2 Validation Requise

| Phase | Validation | Responsable |
|-------|-----------|-------------|
| Design | Architecture Review | Architecte National |
| Review | Security Review | CSO |
| Build | UX Validation | Design Authority |
| Pre-Deploy | Offline Certification | Test Lab |
| Deploy | Release Governance | Change Board |

---

## 11. OBSERVABILITÉ

### 11.1 Métriques Clés

| Métrique | Seuil | Alerte |
|----------|-------|--------|
| App Uptime | >99.9% | <99.5% |
| API Latency p95 | <500ms | >1s |
| Crash Rate | <0.1% | >0.5% |
| Offline Sync Success | >98% | <95% |
| User Satisfaction | >4.0/5 | <3.5 |

### 11.2 Stack Technique

| Outil | Usage |
|-------|-------|
| Prometheus | Métriques |
| OpenTelemetry | Tracing |
| Loki | Logs |
| Grafana | Dashboards |
| Sentry | Crash Reporting |

---

## 12. DIMENSIONNEMENT

### 12.1 Utilisateurs Estimés

| Type | Estimation | Croissance |
|------|------------|------------|
| Citizens | 12M (3 ans) | 30% an |
| Government Agents | 50,000 | 10% an |
| Police/Justice | 15,000 | 5% an |
| Field Agents | 5,000 | 20% an |

### 12.2 Volume Transactions

| Opération | Volume/Jour | Peak |
|-----------|-------------|------|
| Authentifications | 500,000 | 1,000,000 |
| Vérifications QR | 200,000 | 500,000 |
| Notifications | 1,000,000 | 5,000,000 |
| Sync Offline | 50,000 | 200,000 |

---

## 13. DÉPLOIEMENT

### 13.1 Phases

| Phase | Timeline | Livrables |
|-------|----------|-----------|
| **P12.1** | J0-J30 | Citizen Super App v1 |
| **P12.2** | J15-J45 | Government Super App v1 |
| **P12.3** | J30-J60 | Police & Justice Apps v1 |
| **P12.4** | J45-J75 | Mobile Field Apps v1 |
| **P12.5** | J60-J90 | Wallet & Notifications v1 |
| **P12.6** | J75-J100 | Admin Portal & Governance |
| **P12.7** | J90-J120 | Full Integration |

### 13.2 Infrastructure

| Composant | Spécification |
|-----------|---------------|
| App Servers | 8 vCPU, 16 GB RAM (x6) |
| Sync Servers | 4 vCPU, 8 GB RAM (x4) |
| Database | PostgreSQL 16 (Cluster) |
| Cache | Redis Cluster |
| CDN | National Points of Presence |
| Kiosks | 500 units nationwide |

---

## 14. DOCUMENTS ASSOCIÉS

| Document | Référence |
|----------|-----------|
| Citizen Super App Spec | SNISID-PH12-CSA-001 |
| Government App Spec | SNISID-PH12-GOV-001 |
| Security Framework | SNISID-PH12-SEC-001 |
| UX Design System | SNISID-PH12-UX-001 |
| Offline Sync Engine | SNISID-PH12-OFF-001 |
| App Governance Model | SNISID-PH12-GOV-001 |

---

## 15. APPROBATIONS

| Rôle | Nom | Date | Signature |
|------|-----|------|-----------|
| Architecte National | DNI | 2026-05-25 | — |
| CSO | DNI | 2026-05-25 | — |
| UX Director | DNI | 2026-05-25 | — |

---
*Fin du document — SNISID National Application Ecosystem Architecture v1.0*
