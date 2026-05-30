# ⏪ Runbook 03 — Rollback BPMN

**Severity :** Sev1
**Owner :** WGO Rollback Cell + Astreinte Workflow Engine
**Délai d'exécution cible :** < 30 minutes

## 1. Quand l'utiliser
- Bug critique introduit par une nouvelle version
- Conformité juridique remise en cause (signalement LVB)
- Fuite de données (PII) liée à une étape ajoutée
- SLA dégradé > 30 % suite à un déploiement

## 2. Pré-requis
- Identifier la version stable précédente (Git tag, ex. `v1.2.0` → cible `v1.1.4`)
- Décision **DG WGO** (approbation 2-yeux ou téléphone si crise)

## 3. Procédure

### A. Geler les nouvelles instances
```bash
# Marquer la version courante "frozen" → plus de nouvelles créations
zbctl set process-instance-creation --process-id civil-registry.birth.simple \
  --version-tag 1.2.0 --frozen=true
```

### B. Redéployer la version précédente
```bash
# Récupérer le BPMN signé de la version précédente depuis Git
git checkout tags/civil-registry.birth.simple-v1.1.4 -- BPMN/Civil-Registry/birth-simple.v1.1.4.bpmn

# Redéploiement via le job WGO (vérifie la signature)
npm --prefix workflow-engine run deploy:bpmn -- --file BPMN/Civil-Registry/birth-simple.v1.1.4.bpmn

# Mettre à jour ArgoCD pour pointer vers la version précédente
argocd app set snisid-bpmn --revision v1.1.4
argocd app sync snisid-bpmn
```

### C. Migrer les instances actives
```bash
# Option 1 : Laisser les instances en cours sur la nouvelle version se terminer (compatibilité backward)
# Option 2 : Migrer les instances actives
zbctl migrate process-instance --from civil-registry.birth.simple:1.2.0 \
  --to civil-registry.birth.simple:1.1.4 \
  --mapping mappings.json
```

### D. Auditer & signaler
```bash
./scripts/emit-event.sh governance.bpmn.rollback.v1 \
  '{"workflow":"civil-registry.birth.simple","from":"1.2.0","to":"1.1.4","reason":"..."}'
```

## 4. Vérification
- Nouvelles instances créées sur `v1.1.4` (vérif via Operate UI)
- Aucune nouvelle alerte SLA
- Tableau de bord SLO redevient vert sous 15 min

## 5. Communication
| Audience | Délai |
|----------|-------|
| WGO + Direction | immédiat |
| Métier impacté | 15 min |
| Citoyens (si impact perçu) | 1 h |

## 6. Suivi
- Issue GitLab post-mortem ouverte sous 24 h
- LVB doit re-valider la version v1.2.0 corrigée avant ré-essai
- Nouvelle version proposée doit passer test régression spécifique au bug
