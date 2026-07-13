# Runbook — Metadata Recovery

## Déclencheurs
- Catalogue indisponible.
- Metadata supprimée/corrompue.
- Lineage manquant massif.

## Procédure
1. Déclarer incident plateforme metadata.
2. Basculer en lecture seule si corruption partielle.
3. Restaurer backup catalogue le plus récent.
4. Rejouer événements OpenLineage depuis Audit/queue.
5. Scanner Lakehouse pour réconcilier datasets.
6. Vérifier owners/classifications/rétentions.
7. Bloquer datasets orphelins jusqu'à validation.
8. Rapport de récupération et tests.
