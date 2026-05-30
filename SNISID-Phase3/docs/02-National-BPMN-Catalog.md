# 📚 NATIONAL BPMN CATALOG

> **Phase 3 / Étape 2** — Référentiel unique des workflows nationaux.
> Version : 1.0.0

---

## 1. Objectif

Standardiser **tous** les workflows de l'État haïtien.
Chaque workflow doit être :
- ✅ **Versionné** (SemVer)
- ✅ **Audité** (event sourcing + WORM)
- ✅ **Signé** (PKI nationale)
- ✅ **Validé juridiquement** (Legal Validation Board)

---

## 2. Domaines & Criticité

| Domaine | Criticité | RPO | RTO | Disponibilité |
|---------|-----------|-----|-----|----------------|
| **État Civil** | 🔴 CRITIQUE | 0 s | 5 min | 99,99 % |
| **Identité** | 🔴 CRITIQUE | 0 s | 5 min | 99,99 % |
| **Justice** | 🟠 ÉLEVÉE | 60 s | 15 min | 99,95 % |
| **Élections** | 🟠 ÉLEVÉE | 60 s | 15 min | 99,95 % (99,99 % en période électorale) |
| **Immigration** | 🟠 ÉLEVÉE | 60 s | 15 min | 99,95 % |
| **Fiscalité** | 🟠 ÉLEVÉE | 5 min | 1 h | 99,9 % |
| **Santé** | 🟡 MOYENNE | 15 min | 4 h | 99,5 % |

---

## 3. Convention de Nommage

```
<domaine>.<sous-domaine>.<action>.v<MAJOR>.<MINOR>.<PATCH>
```

Exemples :
- `civil-registry.birth.simple.v1.0.0`
- `identity.enrollment.standard.v2.1.0`
- `judicial.suspension.identity.v1.0.0`

---

## 4. Catalogue Complet

### 4.1 État Civil (`civil-registry.*`)

| ID Workflow | Description | Criticité | SLA |
|-------------|-------------|-----------|-----|
| `civil-registry.birth.simple` | Naissance standard | CRITIQUE | 24 h |
| `civil-registry.birth.recognition` | Naissance par reconnaissance | CRITIQUE | 72 h |
| `civil-registry.birth.late-declaration` | Déclaration tardive | CRITIQUE | 30 j |
| `civil-registry.birth.executive-decree` | Naissance par décret | CRITIQUE | 90 j |
| `civil-registry.birth.judicial-judgment` | Naissance par jugement | CRITIQUE | 60 j |
| `civil-registry.death.standard` | Décès standard | CRITIQUE | 24 h |
| `civil-registry.death.judicial` | Décès judiciaire | CRITIQUE | 30 j |
| `civil-registry.death.disaster` | Décès catastrophe | CRITIQUE | 7 j |
| `civil-registry.disappearance.administrative` | Disparition administrative | CRITIQUE | 1 an |
| `civil-registry.marriage.civil` | Mariage civil | ÉLEVÉE | 7 j |
| `civil-registry.marriage.judicial` | Mariage judiciaire | ÉLEVÉE | 30 j |
| `civil-registry.marriage.religious-recognized` | Mariage religieux reconnu | ÉLEVÉE | 30 j |
| `civil-registry.divorce.administrative` | Divorce administratif | ÉLEVÉE | 30 j |
| `civil-registry.divorce.judicial` | Divorce judiciaire | ÉLEVÉE | 180 j |
| `civil-registry.adoption.national` | Adoption nationale | ÉLEVÉE | 180 j |
| `civil-registry.adoption.international` | Adoption internationale | ÉLEVÉE | 365 j |

### 4.2 Identité (`identity.*`)

| ID Workflow | Description | SLA |
|-------------|-------------|-----|
| `identity.enrollment.standard` | Création identité | 24 h |
| `identity.verification.online` | Vérification en ligne | 5 min |
| `identity.verification.biometric` | Vérification biométrique | 30 s |
| `identity.recovery.standard` | Récupération identité | 72 h |
| `identity.revocation.administrative` | Révocation | 24 h |
| `identity.correction.minor` | Correction mineure | 7 j |
| `identity.correction.major` | Correction majeure | 30 j |
| `identity.duplicate.resolution` | Anti-duplication | 14 j |
| `identity.appeal.citizen` | Contestation citoyen | 30 j |
| `identity.suspension.judicial` | Suspension judiciaire | 24 h |

### 4.3 Justice (`judicial.*`)

| ID Workflow | Description | SLA |
|-------------|-------------|-----|
| `judicial.validation.act` | Validation acte | 48 h |
| `judicial.suspension.identity` | Blocage identité | 4 h |
| `judicial.investigation.fraud` | Enquête fraude | 90 j |
| `judicial.appeal.management` | Gestion appel | 90 j |
| `judicial.court.integration` | Intégration tribunaux | 24 h |

### 4.4 Élections (`elections.*`)

| ID Workflow | Description | SLA |
|-------------|-------------|-----|
| `elections.voter.registration` | Inscription électeur | 7 j |
| `elections.voter.validation` | Validation électeur | 1 h |
| `elections.candidate.registration` | Inscription candidat | 14 j |
| `elections.results.publication` | Publication résultats | 24 h |

### 4.5 Immigration (`immigration.*`)

| ID Workflow | Description | SLA |
|-------------|-------------|-----|
| `immigration.entry.standard` | Entrée standard | 1 min |
| `immigration.visa.request` | Demande visa | 30 j |
| `immigration.deportation.judicial` | Expulsion judiciaire | 30 j |

### 4.6 Fiscalité (`tax.*`)

| ID Workflow | Description | SLA |
|-------------|-------------|-----|
| `tax.taxpayer.registration` | Enregistrement contribuable | 24 h |
| `tax.declaration.annual` | Déclaration annuelle | 90 j |
| `tax.audit.standard` | Contrôle fiscal | 180 j |

### 4.7 Santé (`health.*`)

| ID Workflow | Description | SLA |
|-------------|-------------|-----|
| `health.vaccination.record` | Carnet vaccinal | 24 h |
| `health.epidemic.alert` | Alerte épidémie | 15 min |

### 4.8 Fraude (`fraud.*`)

| ID Workflow | Description | SLA |
|-------------|-------------|-----|
| `fraud.detection.automated` | Détection auto | 1 min |
| `fraud.investigation.case` | Enquête | 60 j |
| `fraud.case.closure` | Clôture dossier | 7 j |

### 4.9 Audit (`audit.*`)

| ID Workflow | Description | SLA |
|-------------|-------------|-----|
| `audit.workflow.record` | Enregistrement audit | temps réel |
| `audit.legal.export` | Export juridique | 24 h |

### 4.10 Offline (`offline.*`)

| ID Workflow | Description | SLA |
|-------------|-------------|-----|
| `offline.enrollment.field` | Enrôlement terrain | tolérant |
| `offline.sync.delayed` | Sync différée | tolérant |
| `offline.biometric.capture` | Biométrie hors ligne | tolérant |

### 4.11 Escalation (`escalation.*`)

| ID Workflow | Description | SLA |
|-------------|-------------|-----|
| `escalation.sla.breach` | Dépassement SLA | 5 min |
| `escalation.crisis.national` | Crise nationale | 1 min |
| `escalation.fraud.critical` | Fraude critique | 5 min |

---

## 5. Exigences Communes à Tout BPMN

Chaque BPMN du catalogue **DOIT** inclure les éléments suivants :

| Élément | Obligatoire | Détail |
|---------|-------------|--------|
| **SLA** | ✅ | Timer boundary event |
| **Escalade** | ✅ | Escalation boundary event |
| **Audit trail** | ✅ | Service task → topic `audit.<domain>` |
| **Human validation** | ✅ | Au moins une User Task (4-eyes si CRITIQUE) |
| **PKI validation** | ✅ | Service task signature qualifiée |
| **Fraud detection** | ✅ | Call activity → moteur anti-fraude |
| **Event sourcing** | ✅ | Émission Kafka à chaque transition |
| **Notifications** | ✅ | SMS/Email/Push citoyen + agent |
| **Compensation** | ✅ | Saga pattern si appels distribués |
| **Versioning** | ✅ | `versionTag` Camunda + SemVer |
| **Legal signature** | ✅ | TSA RFC 3161 sur acte final |

---

## 6. Processus d'Inscription au Catalogue

```
1. Issue GitLab "BPMN proposal"
2. Modélisation Camunda Modeler
3. Revue WGO (Workflow Governance Office)
4. Revue LVB (Legal Validation Board)
5. Tests : SLA, chaos engineering, fraud
6. Signature PKI du fichier .bpmn
7. Merge dans /BPMN/<domaine>/
8. Déploiement GitOps (ArgoCD → Zeebe)
9. Surveillance Prometheus + alertes
```

---

## 7. Matrice de Compatibilité Versions

| Niveau | Impact | Action requise |
|--------|--------|----------------|
| **PATCH** (x.x.+1) | Bugfix | Revue WGO express |
| **MINOR** (x.+1.0) | Nouveau chemin compatible | Revue WGO + tests |
| **MAJOR** (+1.0.0) | Breaking change | Revue WGO + LVB + plan de migration + double-run |

---

**Référentiel maintenu par :** Workflow Governance Office
