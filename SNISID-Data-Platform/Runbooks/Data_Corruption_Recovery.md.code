# Runbook — Data Corruption Recovery

## Déclencheurs
- DQ score critique sous seuil.
- Données incohérentes dans Gold/Silver.
- Checksum ou signature invalide.
- Signalement agence ou utilisateur.

## Procédure
1. Déclarer incident data et classifier impact.
2. Geler publications/consommations du dataset affecté.
3. Identifier fenêtre temporelle et lineage amont/aval.
4. Basculer consommateurs critiques vers dernière version saine.
5. Restaurer depuis snapshot Iceberg/backup si nécessaire.
6. Rejouer pipelines depuis source fiable.
7. Exécuter contrôles qualité renforcés.
8. Publier dataset corrigé après approbation owner/steward.
9. Enregistrer preuves dans Audit Fabric.
10. Post-mortem et règle préventive.
