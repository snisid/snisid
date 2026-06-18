# ADR-005: Stratégie Offline-First avec CRDT pour la Résolution de Conflits

**Statut:** Accepté
**Date:** 2026-06-18
**Décideurs:** Lead Architecte, Domain Experts

## Contexte
Haïti a une connectivité internet intermittente. Le système doit fonctionner hors-ligne et se synchroniser à la reconnexion, avec une résolution correcte des conflits.

## Décision
- Architecture Offline-First: toutes les opérations critiques fonctionnent en local
- CRDT (Conflict-free Replicated Data Types) pour la résolution automatique des conflits
- WatermelonDB (React Native) et SQLite pour le stockage local
- Sync Engine asynchrone avec queue de réconciliation priorisée
- Horizon temporel des conflits: Last-Writer-Wins avec horodatage NTP

## Conséquences
Positives:
- Résilience aux pannes réseau
- Expérience utilisateur fluide (pas de mode "hors-ligne" visible)
- Scalabilité horizontale (les agents de terrain peuvent travailler indépendamment)

Négatives:
- Complexité de la logique de réconciliation
- Stockage local augmenté
- Gestion des conflits métier complexes (ex: double enrôlement)
