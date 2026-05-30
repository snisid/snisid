# 🕵️ NATIONAL FRAUD ANALYTICS PLATFORM

> **Objectif** : Détecter rapidement la fraude nationale sous toutes ses formes.

---

## 1. PÉRIMÈTRE DE DÉTECTION

| Domaine | Support | Mécanisme |
|---------|:-------:|-----------|
| Identity fraud | ✅ | ML scoring + règles + biométrie |
| Duplicate identities | ✅ | Déduplication biométrique 1:N |
| Suspicious enrollments | ✅ | CEP Flink + anomaly detection |
| Insider abuse | ✅ | UEBA agents |

---

## 2. ARCHITECTURE

```
[Events Kafka] ─┬─► [Flink CEP rules]      ──► fraud_alerts (Kafka)
                ├─► [ML Scoring Service]    ──► fraud_alerts
                ├─► [Biometric 1:N Engine]  ──► duplicate_clusters
                └─► [UEBA Agent Analytics]  ──► insider_alerts
                              │
                              ▼
                ┌────────────────────────────┐
                │ Fraud Case Management UI   │
                │  (analystes NRIC)          │
                └────────────────────────────┘
                              │
                              ▼
                ┌────────────────────────────┐
                │ Lakehouse — Gold fraud     │
                │ Reporting + ML retraining  │
                └────────────────────────────┘
```

---

## 3. DÉTECTION DOUBLONS BIOMÉTRIQUES

- ABIS national (Automated Biometric Identification System)
- Comparaison 1:N sur empreintes + visage + iris (selon dispo)
- Clustering périodique (graph-based)
- Score de similarité > 0.92 → alerte fraude

```python
# duplicate_cluster.py
import networkx as nx

G = nx.Graph()
for (id_a, id_b, score) in similarity_pairs:
    if score >= 0.92:
        G.add_edge(id_a, id_b, weight=score)

clusters = list(nx.connected_components(G))
suspicious = [c for c in clusters if len(c) >= 2]
```

---

## 4. RÈGLES CEP (extraits)

| ID | Règle | Action |
|----|-------|--------|
| R001 | > 50 enrôlements/agent/heure | Alerte orange |
| R002 | Même appareil > 3 NIN/jour | Alerte rouge |
| R003 | Coord. GPS instables (téléportation) | Alerte rouge |
| R004 | Match biométrique > 0.95 sur 2 NIN différents | Bloquer + alerte critique |
| R005 | Agent consulte > 100 dossiers sans action | Alerte insider |
| R006 | Modification batch hors horaires | Alerte insider |
| R007 | Enrôlement avec doc identité signalé volé | Bloquer immédiat |

---

## 5. MODÈLE ML — FRAUD SCORING

Features clés :
- Vitesse enrôlement agent
- Historique agent (anomalies passées)
- Qualité capture biométrique
- Distance géographique vs résidence déclarée
- Cohérence cross-référentiel (état civil, fiscal)
- Heure de l'opération
- Device fingerprint

```python
# fraud_model_training.py
import xgboost as xgb, mlflow

with mlflow.start_run(run_name="fraud_scoring_v3.3"):
    params = dict(
        objective="binary:logistic",
        max_depth=8, eta=0.05,
        scale_pos_weight=12,    # déséquilibre classes
        eval_metric=["auc","aucpr"]
    )
    booster = xgb.train(params, dtrain, num_boost_round=600,
                        evals=[(dval,"val")], early_stopping_rounds=30)

    mlflow.log_metrics({"auc": auc(booster), "pr_auc": pr_auc(booster)})
    mlflow.xgboost.log_model(booster, "model",
        registered_model_name="fraud_scoring")
```

---

## 6. UEBA — INSIDER ABUSE

| Signal | Détection |
|--------|-----------|
| Pic d'accès | Baseline + écart-type |
| Heures atypiques | Profil horaire par rôle |
| Données sensibles consultées sans dossier | Corrélation tickets |
| Export volumineux | DLP intégré |
| Tentative privilege escalation | SIEM corrélation |

---

## 7. CASE MANAGEMENT

UI dédiée aux analystes NRIC :
- File d'attente alertes priorisées par score
- Vue 360° du suspect / dossier
- Lien direct biométrie, GEOINT, historique
- Workflow : investiguer → confirmer / rejeter → action (blocage, judiciaire)
- Toutes actions loggées + signature analyste

---

## 8. KPI FRAUDE

| KPI | Cible |
|-----|-------|
| Précision détection (top 1%) | > 85 % |
| Recall fraudes confirmées | > 75 % |
| Time-to-alert (P95) | < 60 s |
| Taux faux positifs | < 5 % |
| Doublons biométriques détectés | > 99 % |
| Cas traités sous SLA (48h) | > 95 % |

---

## 9. RÉTROACTION ML

Cycle vertueux :
1. Analyste qualifie une alerte (vrai/faux positif)
2. Label remonte dans Lakehouse
3. Re-training mensuel modèle
4. Validation shadow mode
5. Promotion production via runbook
