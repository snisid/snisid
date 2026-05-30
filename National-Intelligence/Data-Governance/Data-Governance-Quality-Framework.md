# 📐 DATA GOVERNANCE & QUALITY FRAMEWORK

> **Objectif** : Garantir que les décisions reposent sur des données fiables.

---

## 1. PILIERS

| Fonction | Support | Outil |
|----------|:-------:|-------|
| Data lineage | ✅ | OpenLineage + Marquez |
| Data quality scoring | ✅ | Great Expectations / Soda |
| Data ownership | ✅ | Data Mesh — domain owners |
| Metadata governance | ✅ | DataHub / OpenMetadata |

---

## 2. ORGANISATION

### Rôles

| Rôle | Responsabilité |
|------|----------------|
| **Chief Data Officer (CDO)** | Stratégie données SNISID |
| **Data Steward (par domaine)** | Qualité & règles métier |
| **Data Owner (métier)** | Décisions usage / partage |
| **Data Custodian (technique)** | Mise en œuvre technique |
| **Data Consumer** | Utilisateur analytique |

### Domaines (data mesh)

Identité · Population · Fraude · Opérations · Crise · GEOINT · Sécurité.

Chaque domaine = produits données documentés + SLA + qualité.

---

## 3. DATA LINEAGE

OpenLineage instrumenté dans :
- Airflow (DAGs)
- Spark / Flink jobs
- dbt transformations
- API ingestion

Visualisation Marquez : du dataset source au dashboard final.

```
[snisid_db.enrollments]
        ↓ (Debezium CDC)
[kafka:snisid.enrollments]
        ↓ (Flink job)
[delta:bronze.enrollments]
        ↓ (Spark pipeline)
[delta:silver.enrollments]
        ↓ (dbt model)
[iceberg:gold.enrollments_daily]
        ↓ (Superset dataset)
[Dashboard: Cockpit Présidentiel]
```

Toute donnée d'un cockpit est traçable jusqu'à sa source.

---

## 4. DATA QUALITY SCORING

```yaml
# great_expectations/expectations/gold.enrollments.yml
expectation_suite_name: gold.enrollments
expectations:
  - expect_column_values_to_not_be_null:
      column: enrollment_id
  - expect_column_values_to_be_unique:
      column: enrollment_id
  - expect_column_values_to_match_regex:
      column: nin_hash
      regex: "^[a-f0-9]{64}$"
  - expect_column_values_to_be_in_set:
      column: region
      value_set: [Ouest, Sud-Est, Nord, Nord-Est, Artibonite, Centre,
                  Sud, Grand-Anse, Nord-Ouest, Nippes]
  - expect_column_value_lengths_to_be_between:
      column: agent_id
      min_value: 6
      max_value: 20
  - expect_table_row_count_to_be_between:
      min_value: 100
      max_value: 1000000
```

### Score qualité (0-100)

```
DQ_Score = 0.4·Completeness + 0.3·Validity + 0.15·Uniqueness
         + 0.1·Timeliness   + 0.05·Consistency
```

Seuil minimum production : **DQ_Score ≥ 90**.
Dataset < 80 → bloqué pour usage décisionnel.

---

## 5. DATA OWNERSHIP — DATA PRODUCT CONTRACT

```yaml
# data_products/gold.enrollments_daily.yml
name: gold.enrollments_daily
domain: identite
owner: dg-snisid-identite
steward: jane.exemple@snisid.ht
description: "Enrôlements agrégés quotidiens par région"
sla:
  freshness_minutes: 30
  availability: 99.9
  dq_score_min: 92
schema:
  - {name: event_date, type: date}
  - {name: region, type: string}
  - {name: enrollments_count, type: long}
  - {name: avg_quality, type: double}
classification: confidential
consumers:
  - cockpit_presidentiel
  - bi_superset
  - ml_demand_forecast
retention_years: 10
```

---

## 6. METADATA GOVERNANCE

OpenMetadata / DataHub centralise :
- Catalogue datasets
- Glossaire métier (terme → définition officielle)
- Tags (PII, confidentialité, domaine)
- Politiques d'accès
- Notations qualité
- Lineage

---

## 7. CLASSIFICATION

| Niveau | Description | Exemple |
|--------|-------------|---------|
| `PUBLIC` | Diffusion libre | Statistiques agrégées >10k |
| `INTERNAL` | Usage interne SNISID | Métriques opérationnelles |
| `CONFIDENTIAL` | Restreint domaine | Enrôlements détaillés |
| `SECRET` | Restreint NRIC + direction | Alertes fraude actives |
| `TOP_SECRET` | Présidence/sécurité nat. | Analyses menaces |

Étiquetage automatique via tags + revue manuelle.

---

## 8. PROTECTION PII

- Tokenisation NIN à l'ingestion
- Pseudonymisation pour gold/platinum
- Vues sécurisées avec masking par rôle
- Right-to-erasure documenté (légal)
- Logs accès PII immuables

---

## 9. POLITIQUES & PROCESSUS

- **Politique d'accès** : least privilege, approbation steward
- **Politique de rétention** : 10 ans minimum (légal), TTL Bronze/Silver paramétrable
- **Politique de qualité** : pas de prod sans suite Great Expectations
- **Politique de modification** : tout changement schéma = PR + tests
- **Revue trimestrielle** : comité data governance

---

## 10. KPI GOUVERNANCE

| KPI | Cible |
|-----|-------|
| Datasets avec owner | 100 % |
| Datasets avec lineage | > 95 % |
| Datasets avec tests qualité | > 90 % |
| DQ score moyen Gold | > 95 |
| Incidents qualité (mois) | < 5 |
| Délai résolution incident DQ | < 24h |
