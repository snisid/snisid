# Identity Runbooks

> **Procédures Opérationnelles Identité — SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Industrialiser les opérations identité. Les procédures doivent être :
- **Répétables** — Exécutables par tout agent formé
- **Auditables** — Chaque action est tracée
- **Documentées** — Pas d'ambiguïté
- **Testées** — Validation régulière

---

## 2. RUNBOOKS

### 2.1 Identity Fraud — Investigation

```yaml
runbook: identity-fraud-investigation
title: Investigation de Fraude Identitaire
severity: critical
sla:
  detection_to_response: 15_minutes
  investigation_start: 30_minutes
  resolution: 24_hours
  escalation: 72_hours
team: fraud_investigator

trigger:
  event: fraud.detected
  source: SIEM / Duplicate Engine / Manual report
```

**Étapes :**

| # | Action | Détail | Commande |
|---|--------|--------|----------|
| 1 | Confirmer la détection | Vérifier l'alerte dans le SIEM OpenSearch. Confirmer que ce n'est pas un faux positif. | `snisid-cli fraud verify --alert-id <id>` |
| 2 | Isoler l'identité suspecte | Suspendre temporairement le NNU. Bloquer les accès en cours. Notifier les partenaires. | `snisid-cli identity suspend --nnu <nnu> --reason fraud_investigation` |
| 3 | Collecter preuves | Logs d'accès, historique biométrique, documents d'enrôlement, corrélations. | `snisid-cli investigation collect --nnu <nnu> --from 30d --output /evidence/` |
| 4 | Analyser preuves | Vérifier authenticité documents, analyser biométrie, corréler avec base nationale. | `snisid-cli investigation analyze --case /evidence/<case>/` |
| 5 | Décider | Si fraude confirmée → Révoquer + Escalade judiciaire. Si faux positif → Rétablir. | `snisid-cli investigation decide --case <case> --decision revoke\|restore\|escalate` |
| 6 | Exécuter décision | Révocation: NNU→revoked, certs→CRL. Rétablissement: NNU→active. Escalade: Dossier→Parquet. | `snisid-cli investigation execute --case <case>` |
| 7 | Documenter et clôturer | Rapport final, leçons apprises, mise à jour règles détection. | `snisid-cli investigation close --case <case> --report /reports/<case>.md` |

**Post-actions :** update_detection_rules, notify_stakeholders, archive_case, generate_report, lessons_learned_session

---

### 2.2 Enrollment Failure — Récupération

```yaml
runbook: enrollment-failure-recovery
title: Récupération d'Échec d'Enrôlement
severity: high
sla:
  detection_to_response: 30_minutes
  recovery_start: 1_hour
  resolution: 4_hours
team: enrollment_officer
```

**Étapes :**

| # | Action | Détail | Commande |
|---|--------|--------|----------|
| 1 | Diagnostiquer l'échec | Vérifier les logs du workflow. Identifier l'étape ayant échoué. | `snisid-cli workflow logs --workflow enrollment --id <id>` |
| 2 | Corriger le problème | Technique: redémarrer service. Document: demander correct. Biométrie: re-capturer. | `snisid-cli enrollment fix --id <id> --type technical\|document\|biometric\|network` |
| 3 | Reprendre l'enrôlement | Reprendre le workflow à l'étape suivante. | `snisid-cli enrollment resume --id <id> --from-step <step>` |
| 4 | Vérifier résultat | Confirmer identité créée, NNU généré, biométrie ok, consentements ok. | `snisid-cli enrollment verify --id <id>` |
| 5 | Notifier le citoyen | Envoyer notification de succès + instructions wallet. | `snisid-cli notification send --nnu <nnu> --type enrollment_success` |

---

### 2.3 Duplicate Resolution — Escalade

```yaml
runbook: duplicate-resolution-escalation
title: Résolution de Doublons et Escalade
severity: high
sla:
  detection_to_response: 1_hour
  investigation_start: 2_hours
  resolution: 24_hours
  escalation_judicial: 72_hours
team: identity_analyst
```

**Étapes :**

| # | Action | Détail | Commande |
|---|--------|--------|----------|
| 1 | Identifier les identités dupliquées | Récupérer candidates, comparer données, déterminer source. | `snisid-cli duplicate identify --trigger-id <event_id>` |
| 2 | Analyser la nature | Erreur saisie → Correction. Fraude involontaire → Investigation. Fraude intentionnelle → Judiciaire. | `snisid-cli duplicate analyze --identities <nnu_list>` |
| 3 | Résoudre le doublon | Fusionner dans source, révoquer doublons, mettre à jour références. | `snisid-cli duplicate resolve --source <nnu> --duplicates <list> --action merge\|revoke` |
| 4 | Escalader si nécessaire | Si fraude intentionnelle → Dossier au Parquet avec preuves. | `snisid-cli duplicate escalate --case <id> --to judicial` |
| 5 | Clôturer et archiver | Documenter résolution, archiver, mettre à jour règles. | `snisid-cli duplicate close --case <id>` |

---

### 2.4 IAM Outage — Continuité

```yaml
runbook: iam-outage-continuity
title: Continuité d'Activité — Panne IAM
severity: critical
sla:
  detection_to_response: 5_minutes
  rto: 1_hour
  rpo: 15_minutes
team: iam_admin + cyber_analyst
```

**Étapes :**

| # | Action | Détail | Commande |
|---|--------|--------|----------|
| 1 | Confirmer la panne | Vérifier dashboards Grafana. Confirmer étendue. Identifier services affectés. | `snisid-cli health check --all` |
| 2 | Activer plan de continuité | Basculer vers site de secours. Activer fallback. Notifier parties. | `snisid-cli continuity activate --site secondary` |
| 3 | Diagnostiquer la cause | Analyser logs système. Vérifier infrastructure. Identifier cause racine. | `snisid-cli incident diagnose --service iam --since 1h` |
| 4 | Résoudre le problème | Appliquer correctif. Redémarrer services. Vérifier santé système. | `snisid-cli incident fix --cause <root_cause> --action <fix>` |
| 5 | Basculer vers site principal | Vérifier stabilité. Re-basculer services. Synchroniser données. | `snisid-cli continuity switch-back --verify-stability true` |
| 6 | Post-incident | Documenter incident. Analyser leçons apprises. Mettre à jour procédures. | `snisid-cli incident close --report /reports/<incident>.md` |

**Post-actions :** review_infrastructure, update_procedures, test_continuity_plan (within 7 days), generate_report

---

### 2.5 Biometric Mismatch — Validation

```yaml
runbook: biometric-mismatch-validation
title: Validation de Non-Correspondance Biométrique
severity: high
sla:
  detection_to_response: 15_minutes
  validation_start: 30_minutes
  resolution: 4_hours
team: enrollment_officer + identity_analyst
```

**Étapes :**

| # | Action | Détail | Commande |
|---|--------|--------|----------|
| 1 | Vérifier qualité capture | Contrôler scores qualité. Vérifier conditions de capture. | `snisid-cli biometric quality-check --capture-id <id>` |
| 2 | Re-capturer si nécessaire | Si qualité insuffisante → re-capturer. Autre capteur si disponible. | `snisid-cli biometric recapture --nnu <nnu> --modality fingerprint\|face` |
| 3 | Re-vérifier matching | Comparer nouvelle capture avec registre. Vérifier scores. | `snisid-cli biometric verify --capture-id <new_id>` |
| 4 | Investiger si mismatch confirmé | Vérifier identité citoyen. Analyser historique biométrique. | `snisid-cli biometric investigate --nnu <nnu>` |
| 5 | Décider | Erreur technique → Corriger. Changement physiologique → Mettre à jour template. Usurpation → Escalader fraude. | `snisid-cli biometric decide --case <id> --decision correct\|update\|escalate` |
| 6 | Exécuter et clôturer | Appliquer décision. Documenter résultat. Archiver dossier. | `snisid-cli biometric close --case <id>` |

---

## 3. MATRICE DES RUNBOOKS

| Runbook | Sévérité | SLA | Équipe | Test |
|---------|----------|-----|--------|------|
| Identity Fraud | Critique | 24h | Fraud Investigator | Trimestriel |
| Enrollment Failure | Haute | 4h | Enrollment Officer | Mensuel |
| Duplicate Resolution | Haute | 24h | Identity Analyst | Trimestriel |
| IAM Outage | Critique | 1h | IAM Admin + Cyber | Mensuel |
| Biometric Mismatch | Haute | 4h | Enrollment + Analyst | Trimestriel |

---

## 4. TESTS DES RUNBOOKS

```yaml
runbook_testing:
  frequency:
    critical: monthly
    high: quarterly
    medium: semi_annually
    low: annually
  
  method:
    - tabletop_exercise
    - simulation
    - red_team_exercise
    - chaos_engineering
  
  validation:
    - sla_met: true
    - steps_followed: true
    - communication_effective: true
    - documentation_complete: true
    - lessons_learned_captured: true
```

---

> **Les opérations identité sont répétables, auditées et testées.**
