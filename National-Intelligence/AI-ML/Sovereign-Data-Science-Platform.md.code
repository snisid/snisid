# 🧪 NATIONAL DATA SCIENCE PLATFORM (souveraine)

> **Objectif** : Industrialiser la data science souveraine, modèles auditables.

---

## 1. CAPACITÉS

| Fonction | Support |
|----------|:-------:|
| ML experimentation | ✅ |
| Model training | ✅ |
| GPU workloads | ✅ |
| Secure datasets | ✅ |

---

## 2. STACK

| Domaine | Outil |
|---------|-------|
| Notebooks multi-user | **JupyterHub** (sur K8s) |
| ML Platform / Pipelines | **Kubeflow Pipelines** |
| Experiment tracking | **MLflow** |
| Feature Store | Feast |
| Model registry | MLflow Registry |
| Model serving | KServe |
| GPU orchestration | NVIDIA GPU Operator |
| Dataset versioning | DVC + lakeFS |
| Explainability | SHAP, LIME, Captum |
| Bias/Fairness | Fairlearn, AIF360 |

---

## 3. ARCHITECTURE

```
┌────────────────────────────────────────────────┐
│ Data Scientists / Analystes (browser)          │
└──────────────────┬─────────────────────────────┘
                   │  OIDC SSO
                   ▼
┌────────────────────────────────────────────────┐
│ JupyterHub (K8s)  ←→  Kubeflow Pipelines       │
│                  ←→  MLflow Tracking            │
└──────────────────┬─────────────────────────────┘
                   │
        ┌──────────┼────────────┐
        ▼          ▼            ▼
   [Spark on K8s] [Ray]    [GPU pools]
        │          │            │
        └──────────┴────────────┘
                   │
                   ▼
        ┌────────────────────┐
        │ Lakehouse (Delta)  │
        │ Feature Store      │
        └────────────────────┘
                   │
                   ▼
        ┌────────────────────┐
        │ Model Registry      │
        │ → KServe deploy     │
        └────────────────────┘
```

---

## 4. SÉCURITÉ DATASETS

| Mécanisme | Description |
|-----------|-------------|
| Workspace tenanté | Namespace K8s par équipe |
| RBAC datasets | Lecture restreinte par niveau de classification |
| Pseudonymisation auto | Vue gold avec NIN haché |
| DLP | Détection PII en notebook (pre-commit hook) |
| Audit | Toutes lectures de datasets sensibles loguées |
| Sortie contrôlée | Export datasets nécessite approbation |

---

## 5. WORKFLOW STANDARD

```
1. Demande projet → comité d'éthique IA
2. Création namespace + accès datasets approuvés
3. Expérimentation en JupyterHub
4. Tracking MLflow obligatoire (autolog)
5. Tests biais + explainability
6. Revue par pair + comité
7. Promotion vers staging (Kubeflow pipeline)
8. Shadow mode 30 jours
9. Production via KServe
10. Monitoring drift continu
```

---

## 6. EXEMPLE PIPELINE KUBEFLOW

```python
# pipeline_fraud_v3.py
import kfp
from kfp import dsl

@dsl.component
def ingest(): ...

@dsl.component
def feature_eng(): ...

@dsl.component
def train(params: dict): ...

@dsl.component
def evaluate(): ...

@dsl.component
def fairness_check(): ...

@dsl.component
def register_if_passed(): ...

@dsl.pipeline(name="fraud-v3-pipeline")
def fraud_pipeline():
    d  = ingest()
    f  = feature_eng().after(d)
    t  = train({"max_depth": 8, "eta": 0.05}).after(f)
    e  = evaluate().after(t)
    fc = fairness_check().after(e)
    register_if_passed().after(fc)

kfp.Client().create_run_from_pipeline_func(fraud_pipeline, arguments={})
```

---

## 7. GPU WORKLOADS

- Pool GPU dédié (NVIDIA A100 / H100 selon disponibilité souveraine)
- Quota par équipe
- Scheduling MIG (Multi-Instance GPU)
- Priorité aux entraînements production
- Monitoring conso (Prometheus DCGM exporter)

---

## 8. AUDITABILITÉ DES MODÈLES

Chaque modèle promu = 1 entrée registre avec :
- Hash du code (Git commit)
- Version des datasets (DVC/lakeFS)
- Paramètres entraînement
- Métriques perf + équité
- Model card YAML
- Approbation comité (signature)
- Lien explainability artifacts

---

## 9. ENVIRONNEMENTS

| Env | Usage |
|-----|-------|
| `sandbox` | Exploration libre, datasets pseudonymisés |
| `research` | Projets validés, GPU partagé |
| `staging` | Tests production-like, shadow mode |
| `prod` | KServe, supervision NRIC |

---

## 10. KPI DATA SCIENCE

| KPI | Cible |
|-----|-------|
| Time-to-first-model | < 2 semaines |
| Models en prod auditables | 100 % |
| Couverture explainability | 100 % |
| Re-training automatisé | > 90 % |
| Drift détecté < 7j | > 95 % |
