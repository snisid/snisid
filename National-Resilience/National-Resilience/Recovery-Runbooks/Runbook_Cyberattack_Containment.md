# Runbook — Cyberattack Containment

## Objectif
Contenir une cyberattaque majeure, empêcher propagation, préserver preuves et restaurer depuis une base propre.

## Déclencheurs
Ransomware, compte privilégié compromis, exfiltration suspectée, DDoS critique, altération IAM/registre.

## Étapes
1. Activer Cyber Resilience Cell et NRCC si P0 impacté.
2. Isoler segments compromis.
3. Suspendre réplication à risque.
4. Révoquer sessions/tokens suspects.
5. Activer comptes break-glass propres.
6. Préserver logs/images/preuves.
7. Identifier dernier backup propre.
8. Lancer clean room recovery.
9. Rotation secrets et clés exposés.
10. Restaurer P0/P1 progressivement.
11. Surveiller indicateurs de compromission.
12. Valider retour normal par sécurité + DR Lead.

## Validation
Aucun IOC actif, backups propres, secrets rotatés, journaux préservés, P0 validés.
