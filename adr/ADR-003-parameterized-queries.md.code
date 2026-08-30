# ADR-003: Parameterized Queries pour Toutes les Opérations SQL

**Statut:** Accepté
**Date:** 2026-06-18
**Décideurs:** Security Owner, Lead Architecte

## Contexte
De nombreuses requêtes SQL sont construites via `fmt.Sprintf` en Go et f-strings en Python, créant des risques d'injection SQL dans des services critiques (bio-adn, identity, criminal records).

## Décision
- Go: Utiliser `pgx` avec des paramètres positionnels ($1, $2, ...) ou `sqlc` pour la génération de code type-safe
- Python: Utiliser SQLAlchemy ORM ou des requêtes paramétrées avec `?` ou `%(name)s`
- Interdire formellement toute construction de requête par concaténation de chaînes
- Ajouter une règle staticcheck/golangci-lint pour détecter `fmt.Sprintf` avec des motifs SQL

## Conséquences
Positives:
- Élimination des risques d'injection SQL
- Code plus lisible et maintenable
- Meilleure performance (plan de requêtes préparées)

Négatives:
- Refactoring important des couches d'accès aux données
- Cas particuliers (ORDER BY dynamique, IN clauses) nécessitent des helpers contrôlés

## Alternatives considérées
1. ORM complet (GORM): Rejeté (perte de contrôle sur les requêtes complexes)
2. Validation d'entrée uniquement: Rejeté (défense insuffisante)
