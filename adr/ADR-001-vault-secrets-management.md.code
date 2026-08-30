# ADR-001: HashiCorp Vault pour la Gestion des Secrets

**Statut:** Accepté
**Date:** 2026-06-18
**Décideurs:** Security Owner, Lead Architecte

## Contexte
Le code source contient des secrets hardcodés (mots de passe, clés API, tokens JWT) dans plusieurs services Go et Python, exposant le système à des risques critiques de sécurité.

## Décision
Utiliser HashiCorp Vault comme source unique de vérité pour tous les secrets, avec:
- Vault Agent pour l'injection de secrets dans les pods Kubernetes
- External Secrets Operator pour synchroniser les secrets vers Kubernetes Secrets
- Authentification Kubernetes (Kubernetes Auth) pour l'accès des workloads

## Conséquences
Positives:
- Élimination des secrets dans le code source
- Rotation automatique des secrets
- Audit trail complet des accès aux secrets
- Politiques d'accès granulaires (ABAC)

Négatives:
- Complexité opérationnelle ajoutée
- Dépendance sur la disponibilité de Vault
- Courbe d'apprentissage pour les équipes

## Alternatives considérées
1. AWS Secrets Manager / GCP Secret Manager: Rejeté pour souveraineté des données
2. Kubernetes Secrets seuls: Rejeté (pas de rotation, pas d'audit)
3. SOPS + Git: Rejeté (pas de gestion dynamique des secrets)
