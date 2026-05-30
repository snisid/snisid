# Data Lineage System

## Objectif
Tracer l'origine, les transformations, les consommateurs, APIs et workflows de chaque donnée de bout en bout.

## Périmètre lineage

| Domaine | Description |
|---|---|
| Source systems | Origine officielle des données |
| Transformations | Jobs Spark, SQL Trino, Airflow, stream processing |
| Consumers | Dashboards, APIs, agences, modèles IA |
| APIs | Exposition, version, finalité |
| Workflows | Circulation BPMN et intégrations |

## Architecture

```text
Pipelines/SQL/Spark/Kafka/API -> OpenLineage events -> Lineage Store -> Metadata Catalog -> Audit Fabric
```

## Standards

- Chaque pipeline émet événements OpenLineage `START`, `COMPLETE`, `FAIL`.
- Chaque transformation référence code versionné et identité service.
- Chaque dataset dérivé conserve parent datasets.
- Chaque export/API est enregistré comme consommateur.
- Chaque correction manuelle est liée à ticket, approbateur et justification.

## Questions auxquelles le lineage doit répondre

1. D'où vient cette donnée ?
2. Qui l'a modifiée ?
3. Quelle règle a transformé la donnée ?
4. Quels dashboards/API/modèles IA l'utilisent ?
5. Quel impact si la source devient indisponible ?
6. Quelles personnes/agences ont accédé à cette donnée ?

## SLA

| Niveau dataset | Lineage requis | Tolérance |
|---|---|---:|
| CRITIQUE | Technique + métier complet | 0 trou |
| ÉLEVÉ | Technique complet | < 1h retard |
| STANDARD | Technique minimal | < 24h retard |
