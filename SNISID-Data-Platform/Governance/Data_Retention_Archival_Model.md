# Data Retention & Archival Model

## Objectif
Gérer le cycle de vie des données nationales : conservation, legal hold, archivage et suppression contrôlée.

## Politiques

| Domaine | Rétention | Archivage | Suppression |
|---|---:|---|---|
| Citizens golden records | Vie + 10 ans | Archive WORM | Interdite sauf correction légale |
| Civil Registry | À vie | Archive permanente | Interdite |
| Identity credentials | Durée validité + 10 ans | Archive sécurisée | Révocation conservée |
| Audit critical events | 10 ans minimum | WORM | Selon loi et audit |
| API logs standard | 2 ans | Archive froide | Purge approuvée |
| Analytics aggregates | 7 ans | Lakehouse archive | Purge contrôlée |
| Security incidents | 10 ans | WORM | Legal approval |
| Temporary landing files | 90 jours | Non | Auto purge si ingéré |

## Legal hold

Un legal hold suspend toute suppression lorsqu'une donnée est liée à :
- enquête judiciaire,
- audit officiel,
- litige,
- incident sécurité,
- demande d'autorité compétente.

## Exigences d'archivage

- Chiffrement fort.
- Stockage immutable pour données critiques.
- Checksum périodique.
- Test de restauration annuel.
- Metadata conservée avec l'archive.
- Traçabilité des consultations d'archive.

## Processus de suppression

1. Identification données éligibles.
2. Vérification absence legal hold.
3. Approbation owner + legal + security si sensible.
4. Suppression cryptographique ou purge certifiée.
5. Enregistrement preuve dans Audit Fabric.
