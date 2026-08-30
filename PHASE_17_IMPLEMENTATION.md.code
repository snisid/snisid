# PHASE_17_IMPLEMENTATION.md

## Nom de la phase
Phase 17 - Sustainability & Financial Autonomy

## Objectif
Résoudre le risque de dette technologique et d'usure des serveurs à 5 ans, en instaurant un modèle économique "Sovereign-Sustained". Un fonds de roulement matériel doit percevoir structurellement 30% des revenus liés aux API (monétisation KYC/B2B et passeports accélérés).

## Fonctionnalités ajoutées
1. **API de Facturation (Billing API)** : Création d'un microservice de comptabilité capable de traiter des requêtes de paiement (`TransactionRequest`).
2. **Smart Split Engine** : Implémentation de la fonction `ProcessTransaction()` qui prélève mathématiquement 30% du montant payé pour alimenter le `HardwareRefreshFund`, tout en transférant les 70% restants au compte usuel du Trésor Public (`TreasuryAccount`).

## Fichiers créés
- `National-Sustainability-Billing/main.go` (Point d'entrée du serveur - simulé via billing_logic)
- `National-Sustainability-Billing/billing_logic.go` (Logique de répartition)
- `National-Sustainability-Billing/billing_test.go` (Tests de mathématiques de split)

## Fichiers modifiés
- `PHASE_17_IMPLEMENTATION.md` (Création).

## Dépendances ajoutées
- Langage `Go`.
- Bibliothèque `testing`.

## Variables d’environnement
- `HardwareRefreshRatio = 0.30` (Défini en constante inaltérable dans le code source).

## Changements de base de données
- Migration de la mémoire vive vers une base de données **SQLite** transactionnelle (`snisid_billing.db`) en utilisant `modernc.org/sqlite`. 
- Implémentation du support complet **ACID** (Atomicity, Consistency, Isolation, Durability) via `BEGIN TRANSACTION`, `COMMIT` et `ROLLBACK` pour s'assurer que la répartition 70/30 ne puisse jamais créer d'asymétrie financière même en cas de crash du serveur.

## Commandes exécutées
```bash
# Initialisation des dépendances et lancement des tests
cd National-Sustainability-Billing
go mod init billing
go get modernc.org/sqlite
go test
```

## Instructions de déploiement
L'API Billing doit être déployée de manière asynchrone (Kafka ou Pub/Sub) derrière l'API Gateway. Lors de chaque succès de paiement provenant d'une banque pour l'API B2B, un événement de facturation doit appeler ce microservice.

## Procédure de rollback
Pour annuler l'implémentation de la Phase 17 :
```bash
rm -rf National-Sustainability-Billing
git checkout HEAD -- PHASE_17_IMPLEMENTATION.md
```

## Risques connus
- **Concurrency** : La structure `sync.Mutex` fonctionne pour une seule instance serveur. En cas de scaling horizontal (Kubernetes), ce registre en mémoire doit être migré vers une base de données transactionnelle avec verrouillage pessimiste (ex: PostgreSQL `FOR UPDATE`).

## Points à valider manuellement
- Exécuter la commande `go test` dans le répertoire pour valider publiquement que 100 HTG génèrent bien 30 HTG de provisions BRH et 70 HTG pour le Trésor Public.
