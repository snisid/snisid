# Stratégie de Versioning et Release — SNISID

## Versioning Sémantique
Format: `vMAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]`

- **MAJOR**: Changements cassants (breaking changes)
- **MINOR**: Nouvelles fonctionnalités rétrocompatibles
- **PATCH**: Corrections de bugs et sécurité
- **PRERELEASE**: `-alpha.N`, `-beta.N`, `-rc.N`

## Branches

```
main          ─── Produit stable, déploiement production
remediation-main  ─── Branche de remédiation sécurité
develop       ─── Intégration continue
feature/*     ─── Nouvelles fonctionnalités
fix/*         ─── Corrections
release/v*   ─── Préparation de release
```

## Workflow de Release

### 1. Développement
```
feature/auth-vault → PR → develop
```

### 2. Intégration
```
develop → tests automatisés → lint → security scan
```

### 3. Préparation de Release
```
develop → release/v1.2.0
  - Mise à jour CHANGELOG.md
  - Mise à jour des versions
  - Tests de régression
  - Revue sécurité
```

### 4. Release Candidate
```
release/v1.2.0 → v1.2.0-rc.1 → déploiement staging
  - Tests d'intégration
  - Tests de performance
  - Tests de pénétration
  - Validation Security Owner
```

### 5. Production
```
release/v1.2.0 → main → git tag v1.2.0
  - ArgoCD sync automatique
  - Canary deployment (10% → 50% → 100%)
  - Monitoring post-deploy (15 min)
```

### 6. Hotfix
```
fix/critical-bug → main (PR directe approuvée par Security Owner)
  - git tag v1.2.1
  - Merge back vers develop
```

## Release Cadence

| Type | Fréquence | Responsable |
|------|-----------|-------------|
| Security Patch | 24-48h | Security Owner |
| Bug Fix | Hebdomadaire | Engineering Manager |
| Feature Release | Bi-mensuel | Release Manager |
| Major Release | Trimestriel | Programme Director |

## Artefacts de Release

Chaque release produit:
- Image Docker taguée (`registry.snisid.ht/snisid-core:v1.2.0`)
- SBOM (Software Bill of Materials) (`snisid-sbom-v1.2.0.json`)
- Rapport de scan de vulnérabilités
- Release notes + CHANGELOG
- Tag Git signé GPG

## Politique de Rollback

1. **Détection** : Alertes Prometheus (erreur rate > 1%, latency > 1s)
2. **Décision** : Security Owner + Lead Architecte
3. **Action** : `argocd app rollback snisid-core --to-revision X`
4. **Post-mortem** : Analyse des causes dans les 24h

## Signataires Requis par Phase du Projet

| Phase | Approbations requises |
|-------|----------------------|
| Phase 0 (Stratégie) | CNN + AND |
| Phase 1 (Sécurité) | Security Owner + AND |
| Phase 2-3 (Infra) | Lead Architecte |
| Phase 4-7 (Services) | Engineering Manager |
| Production Go | Security Owner + 1 reviewer |
