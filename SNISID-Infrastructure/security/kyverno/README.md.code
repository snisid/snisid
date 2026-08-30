# SNISID — Kyverno Policies : Canonical Location & Security Context
**Classification:** SECRET  
**Note:** Les policies Kyverno d'admission et de conformité cluster-wide sont stockées dans le référentiel canonique :

```
kubernetes/policies/kyverno/
├── istio-mtls.yaml          (Enforcement mTLS STRICT mesh-wide)
├── require-labels.yaml      (Standards labeling national + sécurité pods)
```

Ces policies sont appliquées via ArgoCD sous le project `snisid-national` et validées en CI avant tout merge sur `main`.

Le dossier `security/kyverno/` est réservé aux politiques spécifiquement liées à la sécurité runtime et compliance (ex: CIS benchmarks, Pod Security Admission fallback) qui viennent en complément des ClusterPolicies canoniques.

## Règles de déploiement
- Toute ClusterPolicy est `validationFailureAction: Enforce` en production.
- Les exceptions nécessitent approbation IGC + PR dédiée.
- Pas de `audit` seul sur les règles Tier-0/Tier-1.
