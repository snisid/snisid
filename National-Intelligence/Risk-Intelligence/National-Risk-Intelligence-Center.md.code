# 🛡️ NATIONAL RISK INTELLIGENCE CENTER (NRIC)

> **Objectif** : Anticiper les risques nationaux par fusion analytique multi-domaines.

---

## 1. MISSION

Le NRIC est l'organe central d'analyse stratégique des risques pour SNISID et l'État haïtien.
Il fusionne données opérationnelles, sécurité, fraude, infrastructure et menaces pour :
- **Anticiper** plutôt que réagir
- **Quantifier** le risque (scoring continu)
- **Alerter** les niveaux ministériels et présidentiel
- **Recommander** des mesures préventives

---

## 2. DOMAINES D'ANALYSE

| Domaine | Analyse | Sources |
|---------|---------|---------|
| **Cyber risks** | Tentatives intrusion, exfiltration, ransomware | SIEM, EDR, NDR, Kafka security topics |
| **Identity fraud** | Doublons, usurpations, faux dossiers | Lakehouse Gold fraud, ML scoring |
| **Operational instability** | Pannes services, dégradations SLO | Prometheus, traces, logs |
| **Infrastructure failures** | Datacenters, réseau, énergie | Observability stack, capteurs IoT |
| **National threats** | Menaces sécuritaires, troubles civils | OSINT, partenaires gouvernementaux |

---

## 3. ARCHITECTURE NRIC

```
┌──────────────────────────────────────────────────────┐
│           RISK FUSION DASHBOARD (Cockpit NRIC)       │
└──────────────────────────┬───────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │   RISK SCORING ENGINE   │
              │ (ML + règles + experts) │
              └────────────┬────────────┘
                           │
      ┌────────┬───────────┼───────────┬────────────┐
      ▼        ▼           ▼           ▼            ▼
   [Cyber] [Fraude] [Opérations]  [Infra]    [Menaces nat.]
      │        │           │           │            │
      └────────┴───────────┴───────────┴────────────┘
                           │
                ┌──────────┴──────────┐
                │  LAKEHOUSE Gold     │
                │  + Streaming Kafka  │
                └─────────────────────┘
```

---

## 4. RISK SCORING ENGINE

### 4.1 Formule composite

```
Risk_Score(domain, t) =
    α · ML_anomaly_score(t)
  + β · rule_based_score(t)
  + γ · historical_trend(t)
  + δ · expert_weight
```

Avec α+β+γ+δ = 1, calibrés par domaine et révisés trimestriellement.

### 4.2 Niveaux

| Score | Niveau | Action |
|-------|--------|--------|
| 0–0.3 | 🟢 Vert | Monitoring routine |
| 0.3–0.6 | 🟡 Jaune | Surveillance accrue |
| 0.6–0.8 | 🟠 Orange | Notification ministère |
| 0.8–1.0 | 🔴 Rouge | Alerte présidentielle + cellule crise |

---

## 5. INDICATEURS CLÉS DE RISQUE (KRI)

| KRI | Description | Seuil rouge |
|-----|-------------|-------------|
| `kri.cyber.intrusion_attempts_per_hour` | Tentatives intrusion détectées | > 1000 |
| `kri.fraud.duplicate_identity_rate` | Taux de doublons biométriques | > 0.5 % |
| `kri.ops.critical_services_down` | Services régaliens KO | ≥ 1 |
| `kri.infra.datacenter_health` | Score santé DC | < 0.7 |
| `kri.threat.osint_severity` | Sévérité menaces OSINT | ≥ HIGH |

---

## 6. ALERTING & ESCALADE

```yaml
escalation_matrix:
  green:
    notify: [nric_analyst]
    channel: dashboard
  yellow:
    notify: [nric_analyst, domain_lead]
    channel: [dashboard, email]
  orange:
    notify: [nric_director, ministry_lead]
    channel: [dashboard, email, sms]
    sla_response_minutes: 30
  red:
    notify: [nric_director, ministry_lead, presidential_advisor]
    channel: [dashboard, email, sms, push, hotline]
    sla_response_minutes: 10
    auto_actions:
      - open_crisis_room
      - snapshot_lakehouse
      - lock_high_risk_workflows
```

---

## 7. MODÈLE ML — RISK FORECAST

```python
# risk_forecast.py
from sklearn.ensemble import GradientBoostingRegressor
import mlflow

with mlflow.start_run(run_name="nric_cyber_risk_forecast_v1"):
    model = GradientBoostingRegressor(n_estimators=300, max_depth=6)
    model.fit(X_train, y_train)

    mlflow.log_metric("rmse", rmse(model, X_val, y_val))
    mlflow.log_metric("mae", mae(model, X_val, y_val))
    mlflow.sklearn.log_model(model, "model",
        registered_model_name="nric_cyber_risk_forecast")
```

Tous modèles : **audit obligatoire**, explainability SHAP, validation comité.

---

## 8. WAR-ROOM NRIC

Pièce physique + virtuelle équipée :
- Mur d'écrans Grafana / Superset
- Carte GEOINT temps réel
- Ligne directe ministères + Présidence
- Procédures runbooks à portée
- Journal d'incidents WORM

---

## 9. RAPPORTS STRATÉGIQUES

| Cadence | Destinataire | Contenu |
|---------|--------------|---------|
| Quotidien | Direction SNISID | État des risques 24h |
| Hebdomadaire | Ministères clés | Tendances + actions recommandées |
| Mensuel | Présidence | Vue stratégique nationale |
| Trimestriel | Conseil sécurité | Évolution menaces structurelles |
| Ad hoc | Crise | Briefings temps réel |
