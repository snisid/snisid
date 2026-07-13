# NATIONAL ADMIN PORTAL — SNISID
## Cockpit Administration Nationale

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-ADM-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |
| Classification | CONFIDENTIEL |

---

## 1. PRÉSENTATION

Portail central d'administration pour la gestion de tout l'écosystème applicatif SNISID. Point d'accès unique pour les administrateurs système, les gestionnaires de flux, la gouvernance des API et la supervision des incidents.

### 1.1 Public Cible

| Rôle | Effectif | Accès |
|------|----------|-------|
| Super Admin | 3 | Total |
| System Admin | 10 | Technique |
| Workflow Manager | 5 | Processus |
| API Governance | 4 | API |
| Security Admin | 5 | Sécurité |
| Audit Viewer | 8 | Lecture seule |

---

## 2. FONCTIONNALITÉS

### 2.1 User Administration

| Fonction | Support | Audit |
|----------|---------|-------|
| User Management | ✅ Création, Modification, Désactivation | ✅ |
| Role-Based Access (RBAC) | ✅ | ✅ |
| Permission Groups | ✅ | ✅ |
| Multi-Tenant Management | ✅ | ✅ |
| Account Recovery | ✅ | ✅ |
| Session Management | ✅ Force logout | ✅ |
| Bulk Operations | ✅ Import/Export | ✅ |

### 2.2 Workflow Management

| Fonction | Support |
|----------|---------|
| Workflow Designer | ✅ Drag & Drop |
| Approval Chain Config | ✅ Multi-niveaux |
| Escalation Rules | ✅ SLA-based |
| Workflow Templates | ✅ |
| Workflow Monitoring | ✅ Real-time |
| Workflow Analytics | ✅ Reports |
| Workflow Versioning | ✅ |

### 2.3 API Governance

| Fonction | Support |
|----------|---------|
| API Registry | ✅ |
| Rate Limiting | ✅ Configurable |
| API Key Management | ✅ Rotation, Expiration |
| Usage Analytics | ✅ |
| API Versioning | ✅ |
| Documentation Portal | ✅ |
| API Health Monitoring | ✅ |
| Throttling Rules | ✅ |

### 2.4 Audit Access

| Fonction | Support |
|----------|---------|
| Audit Log Viewer | ✅ Filtrable, Exportable |
| Immutable Audit Trail | ✅ Blockchain-based |
| User Activity Timeline | ✅ |
| Anomaly Detection | ✅ Automated |
| Audit Reports | ✅ PDF, CSV, JSON |
| Retention Management | ✅ Configurable (1-10 ans) |
| Compliance Reports | ✅ |

### 2.5 Incident Visibility

| Fonction | Support |
|----------|---------|
| Incident Dashboard | ✅ Real-time |
| Incident Timeline | ✅ |
| SLA Monitoring | ✅ |
| Escalation Tracking | ✅ |
| Incident Reports | ✅ |
| Root Cause Analysis | ✅ |
| Communication Log | ✅ |
| Resolution Tracking | ✅ |

---

## 3. TABLEAU DE BORD ADMIN

```
┌──────────────────────────────────────────────────────────┐
│  🇭🇹 SNISID Admin Portal  │  Super Admin  │  ⚙️         │
├──────────────────────────────────────────────────────────┤
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐            │
│  │Users   │ │ APIs   │ │Workflow│ │Audit   │            │
│  │ 1,234  │ │  156   │ │   23   │ │ 4,567  │            │
│  └────────┘ └────────┘ └────────┘ └────────┘            │
│                                                          │
│  ┌──────────────────┐  ┌──────────────────┐             │
│  │ Active Incidents  │  │ API Health       │             │
│  │  🔴 3 Critical   │  │ 🟢 98.5% uptime  │             │
│  │  🟡 5 Warning    │  │ 🟡 45 req/s      │             │
│  │  🔵 12 Info      │  │ 🔴 2% error rate │             │
│  └──────────────────┘  └──────────────────┘             │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Audit Feed (Real-time)                          │   │
│  │  14:32:15  admin@dni  Approved API key rotation  │   │
│  │  14:30:01  sys@dni   Deployed CitizenApp v2.1   │   │
│  │  14:28:44  sec@dni   Blocked suspicious API call │   │
│  └──────────────────────────────────────────────────┘   │
├──────────────────────────────────────────────────────────┤
│  Users │ APIs │ Workflows │ Audit │ Incidents │ Settings │
└──────────────────────────────────────────────────────────┘
```

---

## 4. SÉCURITÉ

### 4.1 Accès

| Niveau | Accès | MFA |
|--------|-------|-----|
| **Viewer** | Lecture seule | 2FA |
| **Operator** | Opérations standard | 3FA |
| **Manager** | Gestion configuration | 3FA + Device Attestation |
| **Admin** | Administration complète | 3FA + HSM Token |
| **Super Admin** | Accès total | 3FA + HSM + Physical key |

### 4.2 Auditabilité

```
┌──────────────────────────────────────────────────┐
│            IMMUTABLE AUDIT TRAIL                  │
├──────────────────────────────────────────────────┤
│                                                  │
│  Chaque action admin est :                       │
│  ✅ Horodatée (temps atomique NTP)               │
│  ✅ Signée numériquement                         │
│  ✅ Chaînée (hash du bloc précédent)             │
│  ✅ Stockée en écriture seule                    │
│  ✅ Répliquée (3 copies géo-distribuées)        │
│  ✅ Non-modifiable (append-only ledger)          │
│  ✅ Exportable pour enquête                      │
│                                                  │
└──────────────────────────────────────────────────┘
```

---

## 5. PERFORMANCE

| Métrique | Cible |
|----------|-------|
| Dashboard Load | < 1s |
| User Search | < 500ms |
| Audit Query | < 2s |
| API Analytics | < 3s |
| Incident Load | < 1s |
| Bulk Operations | < 5s per 1000 |
| Export Size Limit | 100,000 records |

---

## 6. INTÉGRATIONS

| Système | Intégration |
|---------|-------------|
| Identity Hub | User sync, Roles |
| API Gateway | API Registry, Metrics |
| Monitoring Stack | Prometheus, Grafana |
| SIEM | Audit logs export |
| Notification Platform | Alerting |
| All Apps | Deployment management |

---

## 7. DÉPLOIEMENT

| Version | Date | Fonctionnalités |
|---------|------|-----------------|
| v1.0-Beta | J+30 | User Mgmt, Audit Viewer |
| v1.0 | J+45 | Workflow, API Governance |
| v1.1 | J+60 | Incident Dashboard |
| v2.0 | J+90 | Full Analytics, Reports |

---
*Fin du document — National Admin Portal v1.0*