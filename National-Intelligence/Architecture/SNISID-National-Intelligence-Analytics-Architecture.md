# 🏛️ SNISID — NATIONAL INTELLIGENCE & ANALYTICS ARCHITECTURE

> Document d'architecture officiel — Phase 18
> Classification : **SOUVERAIN — USAGE GOUVERNEMENTAL**

---

## 1. VISION ARCHITECTURALE

SNISID Intelligence Platform = **architecture analytique souveraine** combinant :

- Ingestion massive multi-source
- Lakehouse national (Delta Lake / Iceberg)
- Pipelines analytiques industrialisés (Spark, Flink)
- Stack IA/ML auditable (Kubeflow, MLflow)
- Decision Intelligence en temps réel
- 100 % hébergée dans les datacenters souverains haïtiens

---

## 2. PRINCIPES DIRECTEURS

| Principe | Description |
|----------|-------------|
| **Souveraineté absolue** | Aucune donnée critique ne sort du territoire |
| **Temps réel** | Latence décisionnelle < 5 secondes pour événements critiques |
| **Auditabilité** | Tous modèles IA traçables (lineage, versioning, explainability) |
| **Supervision humaine** | L'IA assiste, l'humain décide |
| **Scalabilité massive** | Architecture x1000 sans refonte |
| **Résilience** | Multi-AZ, multi-région, mode dégradé offline |
| **Qualité by design** | Data quality scoring obligatoire à l'ingestion |

---

## 3. DOMAINES ARCHITECTURAUX

### 3.1 DATA INGESTION LAYER

| Source | Mode | Technologie |
|--------|------|-------------|
| Bases SNISID (PostgreSQL) | CDC | Debezium → Kafka |
| Événements applicatifs | Streaming | Kafka |
| Logs nationaux | Batch + Stream | Fluent Bit → Loki / Kafka |
| Capteurs IoT (biométrie, terrain) | Streaming | MQTT → Kafka |
| Partenaires gouvernementaux | API REST / SFTP | NiFi |
| Archives historiques | Batch | Airflow → Spark |
| GEOINT (cartes, drones, satellite) | Batch + Stream | NiFi + Kafka |

**Garantie** : ingestion 100 % chiffrée TLS 1.3, signée, journalisée.

---

### 3.2 ANALYTICS PIPELINES

```
[Sources] → [Kafka] → [Flink Streaming] ─┐
                                          ├─→ [Lakehouse Delta/Iceberg]
[Batch APIs] → [Airflow] → [Spark Batch] ─┘
                                          ↓
                              [Feature Store] → [ML/AI Stack]
                                          ↓
                                  [BI / Dashboards / Decision]
```

**Orchestration** : Apache Airflow (DAGs gouvernementaux versionnés Git)
**Stream processing** : Apache Flink (analytics temps réel)
**Batch processing** : Apache Spark sur Kubernetes

---

### 3.3 LAKEHOUSE LAYER

**Architecture en médaillon** :

| Couche | Description | Format | Rétention |
|--------|-------------|--------|-----------|
| **Bronze** | Données brutes ingérées | Delta Lake | 10 ans |
| **Silver** | Données nettoyées, conformes | Delta Lake | 10 ans |
| **Gold** | Données métier agrégées | Delta Lake / Iceberg | Permanent |
| **Platinum** | Datasets décisionnels | Iceberg | Permanent |

**Stockage** : MinIO (S3 souverain) + Ceph (haute disponibilité)
**Catalogue** : Apache Hive Metastore + Unity Catalog équivalent souverain

---

### 3.4 AI / ML STACK

| Composant | Outil |
|-----------|-------|
| Notebooks | JupyterHub multi-tenant |
| ML Platform | Kubeflow |
| Experiment tracking | MLflow |
| Feature Store | Feast |
| Model Serving | KServe / Seldon Core |
| GPU orchestration | NVIDIA GPU Operator on K8s |
| Model registry | MLflow Registry |
| Explainability | SHAP, LIME (obligatoire) |
| Bias detection | Fairlearn, AIF360 |

**Règle absolue** : Aucun modèle en production sans :
- Carte modèle signée
- Tests d'équité
- Audit de biais
- Validation humaine d'un comité d'éthique

---

### 3.5 DECISION INTELLIGENCE LAYER

| Capacité | Description |
|----------|-------------|
| **Real-time scoring** | Décisions fraude / risque < 100ms |
| **Recommendations** | Suggestions workflow gouvernemental |
| **Forecasting** | Prédictions population, infra, demande |
| **Anomaly detection** | Détection immédiate événements anormaux |
| **What-if simulations** | Simulations politiques publiques |
| **Decision logging** | Traçabilité totale des décisions assistées IA |

---

## 4. ARCHITECTURE DE RÉFÉRENCE (DIAGRAMME LOGIQUE)

```
┌──────────────────────────────────────────────────────────────────────┐
│                    PRESIDENTIAL / MINISTRY UI                         │
│  Dashboards • Cockpit • Strategic Intelligence Center                 │
└────────────────────────────┬─────────────────────────────────────────┘
                             │
┌────────────────────────────┴─────────────────────────────────────────┐
│                      DECISION INTELLIGENCE LAYER                      │
│   Recommandations • Scoring • Forecasting • Anomalies • Simulations  │
└────────────────────────────┬─────────────────────────────────────────┘
                             │
       ┌─────────────────────┼─────────────────────┐
       │                     │                     │
┌──────▼──────┐      ┌───────▼───────┐    ┌────────▼────────┐
│  BI / DASH  │      │   AI/ML STACK │    │  REAL-TIME      │
│ Superset    │      │  Kubeflow     │    │  ANALYTICS      │
│ Grafana     │      │  MLflow       │    │  Flink + Druid  │
│ Metabase    │      │  JupyterHub   │    │                 │
└──────┬──────┘      └───────┬───────┘    └────────┬────────┘
       │                     │                     │
       └─────────────────────┼─────────────────────┘
                             │
┌────────────────────────────┴─────────────────────────────────────────┐
│                  NATIONAL LAKEHOUSE (Delta/Iceberg)                   │
│        Bronze → Silver → Gold → Platinum  •  Hive Metastore           │
└────────────────────────────┬─────────────────────────────────────────┘
                             │
┌──────────────┬─────────────┴─────────────┬──────────────┐
│              │                           │              │
▼              ▼                           ▼              ▼
[MinIO]    [Ceph]                  [Spark / Flink]   [Airflow]
                             │
┌────────────────────────────┴─────────────────────────────────────────┐
│              INGESTION LAYER (Kafka + NiFi + Debezium)                │
└────────────────────────────┬─────────────────────────────────────────┘
                             │
   ┌────────┬──────────┬─────┴─────┬──────────┬────────────┐
   ▼        ▼          ▼           ▼          ▼            ▼
[SNISID DBs] [Logs] [Events] [GEOINT]  [IoT/Sensors]  [Partner APIs]
```

---

## 5. SÉCURITÉ ARCHITECTURALE

| Couche | Contrôle |
|--------|----------|
| Réseau | Zero-Trust, mTLS, segmentation par zone |
| Identité | OIDC + RBAC fin granulaire (cf. Phase IAM) |
| Données | Chiffrement at-rest (AES-256) + in-transit (TLS 1.3) |
| Modèles IA | Signature cryptographique + registre auditable |
| Accès analytique | Row-level + column-level security |
| PII | Tokenisation / pseudonymisation systématique |
| Audit | Logs immuables (WORM) sur 10 ans |

---

## 6. SOUVERAINETÉ DES DONNÉES

> **Règle non négociable** : aucune donnée stratégique n'est traitée ou stockée hors du territoire haïtien.

- Datacenters nationaux primaires + secondaires
- Modèles IA entraînés sur infrastructure souveraine
- Pas de dépendance cloud étranger pour les workloads critiques
- Clés cryptographiques détenues par HSM souverain
- Code source des composants critiques auditable par l'État

---

## 7. INTÉGRATION AVEC LES PHASES SNISID ANTÉRIEURES

| Phase | Apport à Phase 18 |
|-------|-------------------|
| Phase IAM | Authentification analystes / décideurs |
| Phase Data Platform | Sources canoniques alimentent le lakehouse |
| Phase Workflows | Événements workflow → analytics temps réel |
| Phase Operations | Métriques opérationnelles → observability |
| Phase Security | Logs sécurité → Risk Intelligence Center |

---

## 8. ROADMAP D'IMPLÉMENTATION

| Trimestre | Livrable |
|-----------|----------|
| T1 | Lakehouse Bronze/Silver + ingestion Kafka |
| T2 | BI Superset + premiers dashboards exécutifs |
| T3 | AI/ML stack + premiers modèles fraude |
| T4 | GEOINT + Crisis Analytics + Decision Cockpit |
| T+1 | Predictive Governance opérationnel |

---

**Document signé** : Direction Nationale SNISID
**Version** : 1.0 — Phase 18
**Classification** : Souverain
