# 📦 BPMN REPOSITORY — Référentiel National

> **Phase 3 / Étape 13** — Centraliser tous les workflows nationaux.
> Version : 1.0.0

---

## 1. Mission

> Un État, un référentiel.
> Aucun BPMN ne vit ailleurs que dans `BPMN/`.

---

## 2. Structure Officielle

```
BPMN/
├── Civil-Registry/        ← Naissance, décès, mariage, divorce, adoption, disparition
├── Identity/              ← Enrollment, verification, recovery, revocation...
├── Judicial/              ← Validation, suspension, fraude, appel, court integration
├── Elections/             ← Voter registration / validation
├── Immigration/           ← Entry, visa, deportation
├── Tax/                   ← Taxpayer registration, declaration, audit
├── Fraud/                 ← Detection, investigation
├── Audit/                 ← Workflow audit record
├── Offline/               ← Offline-first (enrollment, validation, sync, bio, audit)
├── Escalation/            ← SLA breach, crisis national
└── Health/                ← Vaccination, epidemic alert
```

---

## 3. Inventaire Actuel (v1.0.0)

### 3.1 Civil-Registry/ (16 BPMN)

| Fichier | Workflow ID | Criticité | SLA |
|---------|-------------|-----------|-----|
| `birth-simple.v1.0.0.bpmn` | `civil-registry.birth.simple` | CRITIQUE | 24h |
| `birth-recognition.v1.0.0.bpmn` | `civil-registry.birth.recognition` | CRITIQUE | 72h |
| `birth-late-declaration.v1.0.0.bpmn` | `civil-registry.birth.late-declaration` | CRITIQUE | 30j |
| `birth-executive-decree.v1.0.0.bpmn` | `civil-registry.birth.executive-decree` | CRITIQUE | 90j |
| `birth-judicial-judgment.v1.0.0.bpmn` | `civil-registry.birth.judicial-judgment` | CRITIQUE | 60j |
| `death-standard.v1.0.0.bpmn` | `civil-registry.death.standard` | CRITIQUE | 24h |
| `death-judicial.v1.0.0.bpmn` | `civil-registry.death.judicial` | CRITIQUE | 30j |
| `death-disaster.v1.0.0.bpmn` | `civil-registry.death.disaster` | CRITIQUE | 7j |
| `disappearance-administrative.v1.0.0.bpmn` | `civil-registry.disappearance.administrative` | CRITIQUE | 1 an |
| `marriage-civil.v1.0.0.bpmn` | `civil-registry.marriage.civil` | ÉLEVÉE | 7j |
| `marriage-judicial.v1.0.0.bpmn` | `civil-registry.marriage.judicial` | ÉLEVÉE | 30j |
| `marriage-religious-recognized.v1.0.0.bpmn` | `civil-registry.marriage.religious-recognized` | ÉLEVÉE | 30j |
| `divorce-administrative.v1.0.0.bpmn` | `civil-registry.divorce.administrative` | ÉLEVÉE | 30j |
| `divorce-judicial.v1.0.0.bpmn` | `civil-registry.divorce.judicial` | ÉLEVÉE | 180j |
| `adoption-national.v1.0.0.bpmn` | `civil-registry.adoption.national` | ÉLEVÉE | 180j |
| `adoption-international.v1.0.0.bpmn` | `civil-registry.adoption.international` | ÉLEVÉE | 365j |

### 3.2 Identity/ (8 BPMN)

| Fichier | Workflow ID | SLA |
|---------|-------------|-----|
| `enrollment.v1.0.0.bpmn` | `identity.enrollment.standard` | 24h |
| `verification.v1.0.0.bpmn` | `identity.verification.online` | 5 min |
| `recovery.v1.0.0.bpmn` | `identity.recovery.standard` | 72h |
| `revocation.v1.0.0.bpmn` | `identity.revocation.administrative` | 24h |
| `correction.v1.0.0.bpmn` | `identity.correction` | 7j / 30j |
| `duplicate-resolution.v1.0.0.bpmn` | `identity.duplicate.resolution` | 14j |
| `citizen-appeal.v1.0.0.bpmn` | `identity.appeal.citizen` | 30j |
| `judicial-suspension.v1.0.0.bpmn` | `identity.suspension.judicial` | 24h |

### 3.3 Judicial/ (5 BPMN)

| Fichier | Workflow ID | SLA |
|---------|-------------|-----|
| `judicial-validation.v1.0.0.bpmn` | `judicial.validation.act` | 48h |
| `identity-suspension.v1.0.0.bpmn` | `judicial.suspension.identity` | 4h |
| `fraud-investigation.v1.0.0.bpmn` | `judicial.investigation.fraud` | 90j |
| `appeal-management.v1.0.0.bpmn` | `judicial.appeal.management` | 90j |
| `court-integration.v1.0.0.bpmn` | `judicial.court.integration` | 24h |

### 3.4 Elections/ (2 BPMN)
- `voter-registration.v1.0.0.bpmn` — 7j
- `voter-validation.v1.0.0.bpmn` — 1h (p95 < 5min)

### 3.5 Immigration/ (1 BPMN)
- `entry-standard.v1.0.0.bpmn` — 1 min

### 3.6 Tax/ (1 BPMN)
- `taxpayer-registration.v1.0.0.bpmn` — 24h

### 3.7 Fraud/ (1 BPMN)
- `fraud-detection-automated.v1.0.0.bpmn` — 1 min

### 3.8 Audit/ (1 BPMN)
- `audit-workflow-record.v1.0.0.bpmn` — temps réel

### 3.9 Offline/ (5 BPMN)
- `offline-enrollment.v1.0.0.bpmn`
- `offline-validation.v1.0.0.bpmn`
- `offline-biometrics.v1.0.0.bpmn`
- `offline-audit-logs.v1.0.0.bpmn`
- `delayed-sync.v1.0.0.bpmn`

### 3.10 Escalation/ (2 BPMN)
- `sla-breach.v1.0.0.bpmn` — 5 min
- `crisis-national.v1.0.0.bpmn` — 1 min

### 3.11 Health/ (1 BPMN)
- `vaccination-record.v1.0.0.bpmn` — 24h

**TOTAL : 43 BPMN nationaux** (v1.0.0)

---

## 4. Convention de Nommage Fichier

```
<workflow-name>.v<MAJOR>.<MINOR>.<PATCH>.bpmn
```

Exemples valides :
- ✅ `birth-simple.v1.0.0.bpmn`
- ✅ `voter-validation.v2.1.3.bpmn`
- ❌ `birth_simple.bpmn` (sans version, refusé)
- ❌ `Birth-Simple.v1.bpmn` (PATCH manquant)

## 5. Convention de Nommage ID Processus

```
<domain>.<entity>.<action>
```

Avec `versionTag` Camunda en SemVer.
Exemples :
- ✅ `civil-registry.birth.simple` + `versionTag="1.0.0"`
- ✅ `identity.verification.online`
- ❌ `BirthSimple` (camelCase refusé)

---

## 6. Garde-fous (toujours présents dans chaque BPMN)

| Garde-fou | Réalisation |
|-----------|-------------|
| SLA | `boundaryEvent` avec `timerEventDefinition` |
| Escalade | `callActivity` vers `escalation.sla.breach` |
| Audit trail | `serviceTask` avec `audit.emit` au start et end |
| Human validation | au moins un `userTask` (4-eyes si CRITIQUE) |
| PKI validation | `serviceTask` avec `pki.sign.qualified` |
| Fraud detection | `callActivity` vers `fraud.detection.automated` |
| Event sourcing | `serviceTask` avec `kafka.emit` |
| Notifications | `serviceTask` avec `notification.send` |
| Compensation | saga pattern (compensation handlers) |
| Versioning | `zeebe:versionTag` obligatoire |

Une CI lint vérifie ces garde-fous (cf. `scripts/lint-bpmn.sh`).

---

## 7. GitOps Workflow

```
docs:    BPMN/                        ← source de vérité unique
sigs:    .bpmn-signatures/            ← signatures PKI WGO (1:1 avec BPMN)
ci:      .gitlab-ci.yml               ← lint + tests + sign verify
cd:      argocd/applications/bpmn-*.yml ← déploiement Zeebe
```

Pipeline :
```
[PR] → [lint XSD + checks garde-fous] → [tests unitaires]
     → [revue WGO] → [revue LVB] → [sign PKI]
     → [merge main] → [ArgoCD sync staging] → [tests E2E]
     → [approval WGO] → [ArgoCD sync production]
     → [smoke tests] → [monitoring SLO]
```

---

## 8. Lint BPMN — règles automatisées

| Règle | Sévérité | Description |
|-------|----------|-------------|
| `BPMN-001` | ERROR | Pas de `versionTag` |
| `BPMN-002` | ERROR | Pas d'audit start/end |
| `BPMN-003` | ERROR | Pas d'émission Kafka |
| `BPMN-004` | ERROR | Pas de signature PKI (sauf utilitaires) |
| `BPMN-005` | ERROR | Pas d'escalation (SLA boundary) |
| `BPMN-006` | WARN  | Pas de notification |
| `BPMN-007` | ERROR | UserTask sans `candidateGroups` |
| `BPMN-008` | ERROR | ServiceTask sans `taskDefinition.type` |
| `BPMN-009` | ERROR | Pas de signature `.sig` dans `.bpmn-signatures/` |
| `BPMN-010` | WARN  | Pas de Fraud Detection sur workflow CRITIQUE |

Script : `scripts/lint-bpmn.sh` (à exécuter en CI).

---

## 9. Gestion des Versions

| Type | Exemple | Action requise |
|------|---------|----------------|
| **PATCH** (x.x.+1) | bugfix interne | revue WGO express |
| **MINOR** (x.+1.0) | nouveau chemin compatible | revue WGO + LVB + tests |
| **MAJOR** (+1.0.0) | breaking change | revue WGO + LVB + plan migration + double-run minimum 90j |

Chaque MAJOR :
- crée un **nouveau topic Kafka** avec suffixe `.vN`
- migration progressive instances actives (`zbctl migrate process-instance`)
- annonce publique sur portail développeurs

---

## 10. Outils

| Outil | Usage |
|-------|-------|
| **Camunda Modeler 5.x** | Édition `.bpmn` |
| **Zeebe Play** | Test local |
| **Operate** | Visualisation instances |
| **Tasklist** | Tâches humaines |
| **`scripts/lint-bpmn.sh`** | Vérification garde-fous |
| **`scripts/sign-bpmn.sh`** | Signature PKI WGO |
| **`workflow-engine deploy:bpmn`** | Déploiement Zeebe (vérifie sign) |

---

## 11. Politique de Conservation

- **Repo Git** : conservation ad vitam (historique complet)
- **Fichiers BPMN** : tags Git `<workflow>-vMAJ.MIN.PATCH` immortels
- **Signatures `.sig`** : conservées 30 ans (preuve juridique)
- **Versions déployées** : trace dans topic `governance.bpmn.deployed.v1` (30 ans WORM)

---

## 12. Roadmap

| Version | Contenu |
|---------|---------|
| v1.0.0 | 43 BPMN nationaux (livré Phase 3) |
| v1.1.0 | + Visa, Déportation, Audit fiscal |
| v1.2.0 | + Diaspora, services consulaires |
| v2.0.0 | Refonte Civil-Registry avec nouveau code civil HT (à venir) |

---

**Référentiel maintenu par :** Workflow Governance Office
