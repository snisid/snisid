# Runbook — Identity Corruption Recovery

## Objectif
Restaurer l'intégrité des services IAM, registres identité ou données civiles après corruption accidentelle ou malveillante.

## Déclencheurs
Incohérences massives, modification non autorisée, corruption base, alerte intégrité, compromission IAM.

## Étapes
1. Déclarer incident critique P0.
2. Suspendre écritures non essentielles.
3. Préserver snapshots et journaux forensic.
4. Identifier fenêtre de corruption.
5. Stopper réplication si propagation possible.
6. Sélectionner dernier point sain vérifié.
7. Restaurer en environnement isolé.
8. Comparer checksums, échantillons, logs.
9. Rejouer transactions propres si possible.
10. Promouvoir registre restauré après approbation.
11. Réconcilier opérations offline/pending.
12. Produire rapport d'intégrité.

## Validation
Checksums cohérents, tests IAM réussis, registre validé, transactions suspectes revues, RPO documenté.
