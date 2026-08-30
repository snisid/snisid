# Runbook — Analytics Outage Recovery

## Déclencheurs
- Trino/Grafana/Superset indisponible.
- Dashboards nationaux inaccessibles.
- Latence requêtes critique.

## Procédure
1. Identifier composant affecté : query engine, catalog, storage, dashboard, network.
2. Vérifier santé MinIO, Trino workers, metastore/catalog, auth.
3. Activer capacité de secours ou réduire charges non critiques.
4. Prioriser dashboards opérations nationales et sécurité.
5. Redémarrer composants selon ordre contrôlé.
6. Vérifier accès, requêtes critiques, audit logging.
7. Communiquer rétablissement et RCA.
