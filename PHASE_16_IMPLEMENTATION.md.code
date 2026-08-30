# PHASE_16_IMPLEMENTATION.md

## Nom de la phase
Phase 16 - Human Capacity & Brain Drain

## Objectif
Mettre en œuvre les mécanismes technico-RH pour fidéliser l'Unité d'Élite Hors-Grille du SNISID. Cela implique la structuration des accréditations de sécurité ("SECRET Clearance") exigées pour les ingénieurs système, et le suivi de leurs contrats d'engagement.

## Fonctionnalités ajoutées
1. **Accréditation OPA (Policy-as-Code)** : Création d'une politique `engineering_clearance.rego` qui verrouille l'accès aux environnements de `PRODUCTION` du SNISID, en le limitant exclusivement aux ingénieurs possédant la clearance `"SECRET"`.
2. **API Human Capacity** : Création d'un microservice Go (`National-Human-Capacity/main.go`) permettant d'enrôler les ingénieurs. Le service attribue automatiquement le statut "HORS-GRILLE", accorde l'accréditation "SECRET", et configure une date de fin de contrat imposant un engagement de 3 ans (contre l'octroi de certifications financées par l'État).

## Fichiers créés
- `governance/opa-policies/engineering_clearance.rego`
- `governance/opa-policies/engineering_clearance_test.rego`
- `National-Human-Capacity/main.go`

## Fichiers modifiés
- `PHASE_16_IMPLEMENTATION.md` (Création).

## Dépendances ajoutées
- Open Policy Agent (Rego).
- Go (Standard Library `net/http` pour l'API HR).

## Variables d’environnement
- `PORT=8081` (Port de l'API Human Capacity).

## Changements de base de données
- Aucun (la base de données de l'API Go est simulée en mémoire par une `map`).

## Commandes exécutées
```bash
# Lancement de l'API Human Capacity
cd National-Human-Capacity
go run main.go
```

## Instructions de déploiement
Le binaire Go doit être déployé derrière l'API Gateway, et l'API d'authentification IAM (Phase 8) devra s'y synchroniser pour associer la `clearance` du token JWT aux politiques OPA.

## Procédure de rollback
Pour annuler l'implémentation de la Phase 16 :
```bash
rm governance/opa-policies/engineering_clearance*
rm -rf National-Human-Capacity
```

## Risques connus
- **Bloquage de production** : Si la politique OPA est mal déployée ou si l'IAM ne renvoie pas le bon attribut "SECRET" dans le token des ingénieurs, plus personne ne pourra patcher l'infrastructure en production.

## Points à valider manuellement
- Procéder à l'enrôlement d'un ingénieur de test via `POST /api/hr/enroll` et vérifier l'apparition du statut "HORS-GRILLE".
