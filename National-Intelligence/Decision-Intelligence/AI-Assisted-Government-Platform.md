# 🤖 AI-ASSISTED GOVERNMENT PLATFORM

> **Objectif** : Assister la gouvernance par l'IA, sous supervision humaine permanente.

---

## 1. PRINCIPE FONDAMENTAL

> **L'IA propose, l'humain décide.**
> Aucune décision affectant les droits d'un citoyen ne peut être prise par une IA sans validation humaine traçable.

---

## 2. CAPACITÉS

| Fonction | Support | Description |
|----------|:-------:|-------------|
| Workflow recommendations | ✅ | Suggestion optimisation processus |
| Fraud scoring | ✅ | Probabilité fraude par événement |
| Risk predictions | ✅ | Anticipation incidents |
| Operational optimization | ✅ | Routing, allocation ressources |
| Citizen service intelligence | ✅ | Personnalisation services publics |

---

## 3. ARCHITECTURE

```
[Decideur/Agent UI]
     │
     ▼
┌─────────────────────┐
│ AI Assistant Layer  │  (API REST + chat ops)
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ Model Serving       │  KServe / Seldon
│ (versionned models) │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ Feature Store       │  Feast
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ Lakehouse Gold      │
└─────────────────────┘
```

---

## 4. MODÈLES DÉPLOYÉS (CATALOGUE)

| Modèle | Type | Usage | Supervisé par |
|--------|------|-------|---------------|
| `fraud_scoring_v3` | XGBoost | Score fraude enrôlement | Risk Center |
| `workflow_optim_v1` | RL | Routing dossiers | Direction Ops |
| `service_demand_forecast_v2` | LSTM | Prévision charge services | Planning |
| `citizen_intent_classifier_v1` | Transformer | Classification demandes | Service citoyen |
| `agent_risk_score_v1` | LightGBM | Score risque agent (insider) | Sécurité interne |

---

## 5. CONTRAT MODÈLE (MODEL CARD)

Chaque modèle en prod possède une **Model Card** signée :

```yaml
model_id: fraud_scoring_v3
version: 3.2.1
owner: nric_data_science_team
purpose: "Évaluer probabilité fraude lors enrôlement biométrique"
training_data:
  source: gold.enrollments_labeled
  range: 2023-01-01 .. 2026-03-31
  rows: 4_812_345
  pii_handling: pseudonymisation SHA-256 NIN
performance:
  auc: 0.94
  precision_at_top_1pct: 0.87
  recall: 0.79
fairness:
  evaluated_groups: [region, age_bracket, gender]
  max_disparity: 0.06
  status: PASS
explainability:
  method: SHAP global + local
  top_features: [enrollment_speed, agent_history, geo_distance, biometric_quality]
human_oversight:
  decision_threshold: 0.7
  human_review_required_above: 0.5
  escalation: nric_analyst_queue
audit:
  approved_by: ethics_committee_2026Q1
  next_review: 2026Q3
  drift_monitoring: enabled
```

---

## 6. INTERFACE DÉCIDEUR (EXEMPLE)

```
┌────────────────────────────────────────────────────────┐
│ DOSSIER #482931 — Enrôlement Cap-Haïtien                │
├────────────────────────────────────────────────────────┤
│ 🤖 Recommandation IA : VÉRIFICATION HUMAINE REQUISE     │
│    Score fraude : 0.72  (seuil revue : 0.50)            │
│                                                         │
│ Top facteurs (SHAP) :                                   │
│   • Vitesse enrôlement anormale (+0.34)                 │
│   • Empreinte similaire à dossier #392011 (+0.21)       │
│   • Agent à risque modéré (+0.09)                       │
│                                                         │
│ ⚠️ Cette recommandation est INDICATIVE.                  │
│    Décision finale = agent SNISID habilité.             │
│                                                         │
│ [Valider]  [Rejeter]  [Escalader NRIC]  [Justifier]    │
└────────────────────────────────────────────────────────┘
```

Chaque clic est **loggé immuablement** (utilisateur, timestamp, justification).

---

## 7. GOUVERNANCE IA

| Mécanisme | Description |
|-----------|-------------|
| **Comité d'éthique** | Valide tout déploiement modèle |
| **Audit trimestriel** | Performance, biais, drift |
| **Right to explanation** | Tout citoyen peut demander explication |
| **Kill switch** | Désactivation immédiate possible |
| **Shadow mode** | Nouveau modèle tourne en parallèle 30 jours min |
| **Rollback** | Versioning MLflow + procédure runbook |

---

## 8. DRIFT DETECTION

```python
# drift_monitor.py
from evidently.report import Report
from evidently.metric_preset import DataDriftPreset, TargetDriftPreset

report = Report(metrics=[DataDriftPreset(), TargetDriftPreset()])
report.run(reference_data=ref_df, current_data=prod_df)

if report.as_dict()["metrics"][0]["result"]["dataset_drift"]:
    alert("MODEL_DRIFT_DETECTED", model="fraud_scoring_v3")
    trigger_runbook("model_rollback_review")
```

---

## 9. SUPERVISION HUMAINE OBLIGATOIRE

| Type de décision | Niveau supervision |
|------------------|--------------------|
| Affectation droits citoyen | 100 % humain |
| Refus document officiel | 100 % humain |
| Score fraude > seuil | Revue analyste obligatoire |
| Recommandation workflow | Validation manager |
| Forecast capacité | Validation planificateur |

---

## 10. CONFORMITÉ

- Alignée principes OCDE pour IA gouvernementale
- Conforme cadre national haïtien de protection des données
- Documentation publique des cas d'usage (transparence)
