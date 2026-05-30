# SNISID — IAM NATIONAL

> **Source Officielle de Vérité Identitaire Nationale**
> Infrastructure d'Identité Numérique Souveraine, Sécurisée, Auditée et Exploitable 24/7

---

## Vision

**Avant** → **Après**

| Problème | Solution |
|----------|----------|
| Identités fragmentées | Identité numérique nationale unifiée |
| Faux documents | Vérification biométrique + PKI |
| Doublons | Duplicate Detection Engine |
| Absence d'identité fiable | NNU souverain |
| Contrôle accès faible | Zero Trust + RBAC + ABAC + MFA |

---

## Architecture du Repository

```
IAM/
├── Identity-Model/          → Domaine identité national
│   ├── National-Identity-Domain-Model.md
│   └── Duplicate-Detection-Engine.md
├── NNU/                     → Numéro National Unique
│   └── NNU-Specification.md
├── RBAC/                    → Contrôle accès par rôles
│   └── National-RBAC-Specification.md
├── ABAC/                    → Contrôle accès contextuel
│   └── National-ABAC-Specification.md
├── Biometrics/              → Plateforme biométrique
│   └── Biometric-Platform-Specification.md
├── Wallet/                  → Portefeuille identité citoyen
│   └── Citizen-Identity-Wallet-Specification.md
├── Consent/                 → Moteur de consentement
│   └── Consent-Engine-Specification.md
├── Federation/              → Fédération d'identité
│   └── Identity-Federation-Specification.md
├── Workflows/               → Workflows identité
│   └── Identity-Workflows-Specification.md
├── Events/                  → Mesh d'événements
│   └── Identity-Event-Mesh-Specification.md
├── Observability/           → Observabilité identité
│   └── Identity-Observability-Stack.md
├── Runbooks/                → Procédures opérationnelles
│   └── Identity-Runbooks.md
├── Governance/              → Gouvernance identité
│   └── Identity-Governance-Specification.md
└── README.md                ← Ce fichier
```

---

## Principe Absolu

> **Dans SNISID : aucune opération sans identité vérifiée.**

---

## Livrables Phase 6 — Récapitulatif

| Élément | Fichier | Statut |
|---------|---------|--------|
| IAM National | README.md | ✅ |
| NNU | NNU/NNU-Specification.md | ✅ |
| Biométrie | Biometrics/Biometric-Platform-Specification.md | ✅ |
| Wallet Citoyen | Wallet/Citizen-Identity-Wallet-Specification.md | ✅ |
| RBAC | RBAC/National-RBAC-Specification.md | ✅ |
| ABAC | ABAC/National-ABAC-Specification.md | ✅ |
| Federation | Federation/Identity-Federation-Specification.md | ✅ |
| Consent Engine | Consent/Consent-Engine-Specification.md | ✅ |
| Identity Workflows | Workflows/Identity-Workflows-Specification.md | ✅ |
| Identity Events | Events/Identity-Event-Mesh-Specification.md | ✅ |
| Identity Governance | Governance/Identity-Governance-Specification.md | ✅ |
| Observability | Observability/Identity-Observability-Stack.md | ✅ |
| Runbooks | Runbooks/Identity-Runbooks.md | ✅ |
| Duplicate Detection | Identity-Model/Duplicate-Detection-Engine.md | ✅ |
| Identity Domain Model | Identity-Model/National-Identity-Domain-Model.md | ✅ |

---

## Règles de Sécurité

| Règle | Application |
|-------|-------------|
| Zero Trust | Natif, par défaut |
| MFA | Obligatoire pour tous les accès |
| Audit | Chaque action est journalisée |
| Chiffrement | Biométrie et données sensibles |
| Consentement | Explicite et révocable |
| Dé-duplication | Systématique |
| Gouvernance | Centralisée et traçable |

---

## Références Croisées

```
                    ┌─────────────────────┐
                    │    GOVERNANCE       │
                    │  (Politiques, PKI,  │
                    │   Audit, Droits)    │
                    └──────────┬──────────┘
                               │
           ┌───────────────────┼───────────────────┐
           │                   │                   │
    ┌──────▼──────┐     ┌──────▼──────┐     ┌──────▼──────┐
    │ IDENTITY    │     │ ACCESS      │     │ OPERATIONS  │
    │ DOMAIN      │     │ CONTROL     │     │             │
    │             │     │             │     │             │
    │ - Domain    │     │ - RBAC      │     │ - Workflows │
    │   Model     │     │ - ABAC      │     │ - Events    │
    │ - NNU       │     │ - MFA       │     │ - Observ.   │
    │ - Duplicate │     │ - Federation│     │ - Runbooks  │
    │   Detection │     │             │     │             │
    │ - Biometrics│     │ - Consent   │     │ - Wallet    │
    │             │     │   Engine    │     │             │
    └─────────────┘     └─────────────┘     └─────────────┘
```

---

> **SNISID IAM est la source officielle de vérité identitaire nationale.**
