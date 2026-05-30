# Runbook — Datacenter Recovery

## Objectif
Restaurer les services SNISID critiques après perte ou indisponibilité d'un datacenter national.

## Déclencheurs
Perte primary DC, destruction physique, panne énergie prolongée, cyber isolement, décision NRCC.

## Rôles
Crisis Commander, DR Lead, Network Lead, IAM Lead, Database Lead, Communications Lead.

## Étapes
1. Déclarer L3/L4 et ouvrir journal crise.
2. Geler changements non essentiels.
3. Évaluer réplication et dernier point sain.
4. Isoler site perdu/compromis.
5. Activer traffic manager/DNS crise.
6. Promouvoir Secondary National DC ou Regional DR.
7. Restaurer réseau → clés → IAM → registre → APIs P0 → P1.
8. Valider intégrité et accès.
9. Communiquer statut agences.
10. Surveiller 24h en mode crise.

## Validation
IAM fonctionnel, registre accessible, APIs P0 répondent, RTO/RPO mesurés, aucune corruption propagée.
