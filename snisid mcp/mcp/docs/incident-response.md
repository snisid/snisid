# Incident Response SNISID

## Niveaux

- LOW : anomalie sans donnée sensible.
- MEDIUM : tentative refusée répétée.
- HIGH : suspicion compromission compte/device.
- CRITICAL : fuite potentielle, bypass, accès non autorisé sensible.

## Runbook

1. Ouvrir incident via `incidentResponse.openIncident`.
2. Geler sessions et devices suspects.
3. Exporter audit chain et vérifier intégrité.
4. Rotation des clés/API keys si nécessaire.
5. Analyse forensique, rapport, remédiation.
6. Revue des droits RBAC et règles Gateway.
