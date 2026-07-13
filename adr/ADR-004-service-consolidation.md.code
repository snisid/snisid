# ADR-004: Consolidation des Microservices par Bounded Context

**Statut:** Accepté
**Date:** 2026-06-18
**Décideurs:** Lead Architecte, Engineering Manager

## Contexte
Le projet contient ~80 répertoires de services Go, dont ~38 sont vides ou squelettiques, avec une séparation excessive qui complexifie la maintenance et le déploiement.

## Décision
Consolider les services par bounded context (Domain-Driven Design):
- 6-8 services cibles au lieu de 80: gateway, identity-provider, biometrics, fraud-engine, audit, iam-adapter, bio-adn, workflow-engine
- Architecture hexagonale pour chaque service maintenu
- Les services vides sont supprimés ou documentés comme "à intégrer"
- Les responsabilités partagées sont déplacées dans `internal/` ou `pkg/`

## Conséquences
Positives:
- Réduction massive de la complexité opérationnelle
- Équipes alignées sur les bounded contexts
- Déploiements plus rapides et fiables

Négatives:
- Services plus gros (mais plus cohérents)
- Migration nécessaire pour les services existants
- Coordination interfacace plus critique

## Alternatives considérées
1. Status quo (80+ services): Rejeté (ingérable)
2. Monolithe unique: Rejeté (pas d'isolation, pas de scaling indépendant)
