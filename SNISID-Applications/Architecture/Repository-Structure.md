# APPLICATION REPOSITORY — SNISID
## Référentiel Applicatif National

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-REPO-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |

---

## 1. STRUCTURE DU RÉFÉRENTIEL

```
SNISID/
└── Phase-12/
    ├── Architecture/
    │   ├── SNISID-National-Application-Ecosystem-Architecture.md   ✅
    │   ├── Multi-Channel-Access-Model.md                          ✅
    │   └── Repository-Structure.md                                ✅
    │
    ├── Applications/
    │   ├── Citizen-Super-App/
    │   │   └── Citizen-Super-App-Spec.md                          ✅
    │   │
    │   ├── Government-App/
    │   │   └── Government-Super-App-Spec.md                       ✅
    │   │
    │   ├── Police/
    │   │   └── Police-Justice-Apps-Spec.md                        ✅
    │   │
    │   ├── Justice/
    │   │   └── [Référence Police-Justice-Apps-Spec.md]            ✅
    │   │
    │   ├── Mobile-Field/
    │   │   └── Mobile-Field-Apps-Spec.md                          ✅
    │   │
    │   ├── Wallet/
    │   │   └── Digital-Identity-Wallet-Spec.md                    ✅
    │   │
    │   ├── Notifications/
    │   │   └── National-Notification-Platform-Spec.md             ✅
    │   │
    │   ├── Admin-Portal/
    │   │   └── National-Admin-Portal-Spec.md                      ✅
    │   │
    │   ├── UX-Design-System/
    │   │   └── National-UX-Design-System.md                       ✅
    │   │
    │   ├── Security/
    │   │   └── Application-Security-Framework.md                  ✅
    │   │
    │   ├── Offline/
    │   │   └── Offline-Mobile-Sync-Engine.md                      ✅
    │   │
    │   ├── Observability/
    │   │   └── App-Observability-Stack.md                         ✅
    │   │
    │   ├── Governance/
    │   │   └── App-Governance-Model.md                            ✅
    │   │
    │   ├── Runbooks/
    │   │   └── Application-Runbooks.md                            ✅
    │   │
    │   ├── KPIs/
    │   │   └── Application-KPIs.md                                ✅
    │   │
    │   └── [Apps_List.md]                                         ✅
    │
    └── README.md                                                  ✅
```

---

## 2. DESCRIPTION DES DOSSIERS

| Dossier | Contenu | Responsable |
|---------|---------|-------------|
| **Citizen-Super-App/** | App citoyenne nationale | Équipe Mobile |
| **Government-App/** | App gouvernementale | Équipe Web |
| **Police/** | Application Police | Équipe Sécurité |
| **Justice/** | Application Justice | Équipe Sécurité |
| **Mobile-Field/** | Apps terrain offline | Équipe Mobile |
| **Wallet/** | Portefeuille identité | Équipe Sécurité |
| **Notifications/** | Plateforme notification | Équipe Backend |
| **UX-Design-System/** | Système de design | Équipe Design |
| **Security/** | Framework sécurité | Équipe Securité |
| **Offline/** | Moteur sync offline | Équipe Mobile |
| **Admin-Portal/** | Portail admin | Équipe Web |
| **Observability/** | Stack monitoring | Équipe Ops |
| **Governance/** | Modèle gouvernance | Équipe Produit |
| **Runbooks/** | Procédures ops | Équipe Ops |
| **KPIs/** | Indicateurs perf | Équipe Produit |
| **Architecture/** | Architecture globale | Équipe Architecture |

---

## 3. STANDARDS DE RÉFÉRENCEMENT

### 3.1 Conventions de Nommage

| Type | Format | Exemple |
|------|--------|---------|
| Document | `{Nom}-{Type}-{Version}.md` | `Citizen-Super-App-Spec-v1.0.md` |
| Source Code | `{app-name}-{module}` | `citizen-app-auth` |
| Configuration | `{app}.{env}.config` | `citizen-app.prod.config` |
| Docker | `{app}-{service}` | `citizen-app-api` |
| Kubernetes | `{app}-{deploy}` | `citizen-app-deployment` |

### 3.2 IDs de Documents

```
SNISID-PH12-{DOMAINE}-{NUMÉRO}

Domaines:
CSA  → Citizen Super App
GOV  → Government App
PJ   → Police & Justice
MFA  → Mobile Field App
WAL  → Wallet
NOT  → Notification
UX   → UX Design System
SEC  → Security
OFF  → Offline
ADM  → Admin Portal
OBS  → Observability
KPI  → KPIs
RUN  → Runbooks
ARCH → Architecture
REPO → Repository
```

---

## 4. GIT WORKFLOW

### 4.1 Branches

```
main
  └── develop
       ├── feature/citizen-app-auth
       ├── feature/government-app-dashboard
       ├── feature/police-app-cases
       └── release/v1.0.0
```

### 4.2 Commit Convention

```
{type}({scope}): {description}

Types: feat, fix, docs, style, refactor, test, chore
Scopes: citizen, gov, police, wallet, sync, ux, sec

Exemple:
feat(citizen): add offline QR verification
fix(sync): resolve conflict on enrollment data
docs(wallet): update recovery procedure
```

---

## 5. INTÉGRATION CONTINUE

### 5.1 Pipeline

```
Code Push
    │
    ▼
Lint & Format
    │
    ▼
Unit Tests
    │
    ▼
Build APK/IPA
    │
    ▼
SAST Security Scan
    │
    ▼
Offline Test Suite
    │
    ▼
Accessibility Check
    │
    ▼
Performance Benchmark
    │
    ▼
Release Artifact
```

### 5.2 Qualité Requise

| Check | Seuil | Bloquant |
|-------|-------|----------|
| Tests Coverage | > 80% | < 70% |
| Lint Errors | 0 | > 0 |
| Security Vulnerabilities | 0 | > 0 (Critical) |
| Accessibility Score | > 90% | < 80% |
| Bundle Size | < 30 MB | > 40 MB |
| Build Time | < 10 min | > 20 min |

---

## 6. ACCÈS AU RÉFÉRENTIEL

| Rôle | Accès | Branches |
|------|-------|----------|
| Developer | Read/Write | feature/* |
| Senior Dev | Read/Write | develop |
| Tech Lead | Read/Write | release/* |
| Architect | Read only | main (PR review) |
| Security | Read only | All (audit) |
| CI/CD | Read/Write | All (automated) |

---
*Fin du document — Repository Structure v1.0*