# PHASE_13_IMPLEMENTATION.md

## Nom de la phase
Phase 13 - Executive Operations (Parapheur Électronique)

## Objectif
Dématérialisation des documents officiels de l'État (Présidence, Primature, Ministères) via des signatures cryptographiques qualifiées (PKI) et un workflow d'approbation strict, avec stockage dans une base PostgreSQL.

## Fonctionnalités ajoutées
- Création d'un backend robuste en Go avec GORM pour l'entité `Document`.
- Remplacement du mock API par un système connecté à PostgreSQL.
- Implémentation de l'émulation PKI : Validation de code PIN et génération de signature cryptographique (hachage SHA-256).
- Gestion des états des documents (DRAFT, PENDING_SIG, SIGNED).

## Fichiers créés
- Aucun nouveau fichier, mais restructuration massive des existants.

## Fichiers modifiés
- `National-Executive-Operations/api/main.go`
- `National-Executive-Operations/api/main_test.go`
- `National-Executive-Operations/docker-compose.yml`

## Dépendances ajoutées
- `github.com/go-chi/chi/v5`
- `gorm.io/gorm`
- `gorm.io/driver/postgres`

## Variables d’environnement
- `DB_DSN` : Chaîne de connexion à la base de données PostgreSQL. Par défaut : `host=db user=snisid password=snisid dbname=executive_ops port=5432 sslmode=disable`

## Changements de base de données
- Intégration d'un conteneur `postgres:15` dans le `docker-compose.yml`.
- Utilisation de `AutoMigrate(&Document{})` pour créer/mettre à jour la table des documents.

## Commandes exécutées
- `go get github.com/go-chi/chi/v5 gorm.io/gorm gorm.io/driver/postgres`

## Instructions de test
- Les tests unitaires peuvent être lancés via :
```bash
cd National-Executive-Operations/api
go test -v
```
*(Cependant, le manque d'espace disque n'a pas permis d'exécuter les tests).*

## Instructions de build
```bash
cd National-Executive-Operations
docker-compose build
```

## Instructions de déploiement
```bash
cd National-Executive-Operations
docker-compose up -d
```

## Procédure de rollback
Pour annuler, restaurer les versions initiales du dossier `National-Executive-Operations/api` depuis le contrôle de source (Git).

## Risques connus
- La signature cryptographique implémentée est une émulation (hachage SHA-256 local au lieu d'une signature hardware HSM/Smartcard). L'intégration finale avec la PKI Nationale (Phase 1) reste à faire en production.

## Points à valider manuellement
- Interfaçage avec les véritables lecteurs de Smartcard sur les terminaux des Ministres.
- Le problème de saturation du disque sur l'environnement de développement.
