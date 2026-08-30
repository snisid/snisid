# Runbook — Pipeline Failure Resolution

## Déclencheurs
- Job Airflow/Spark échoué.
- Retard ingestion supérieur SLO.
- Écart de volume anormal.

## Procédure
1. Vérifier logs Loki et métriques Prometheus.
2. Identifier type : source, réseau, schéma, qualité, stockage, permission.
3. Si données critiques, notifier owner et NOC/SOC selon impact.
4. Corriger cause racine ou rollback version pipeline.
5. Relancer en mode idempotent.
6. Valider DQ, lineage, metadata.
7. Fermer incident avec timeline et mesures préventives.
