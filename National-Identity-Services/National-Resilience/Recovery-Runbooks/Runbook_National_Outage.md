# Runbook — National Outage

## Objectif
Maintenir continuité SNISID lors d'une panne nationale réseau, énergie ou infrastructure.

## Déclencheurs
Panne Internet nationale, blackout, indisponibilité multi-agences, crise L3/L4.

## Étapes
1. Activer NRCC.
2. Classer outage : réseau, énergie, DC, multi-cause.
3. Prioriser P0/P1.
4. Activer satellite/radio.
5. Mode économie énergie si nécessaire.
6. Activer edge/offline nodes.
7. Suspendre traitements non essentiels.
8. SITREP toutes les 2h ou selon L4.
9. Préparer sync post-outage.
10. Retour normal après validation stabilité.

## Validation
P0 disponible ou offline, communication active, régions prioritaires rapportent statut, logs conservés.
