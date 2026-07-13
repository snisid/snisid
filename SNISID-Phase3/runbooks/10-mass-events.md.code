# 🆘 Runbook 10 — Catastrophe Nationale / Événements de Masse

**Severity :** Sev0
**Owner :** WGO + Direction Présidence ONI + DGPC (Protection Civile)

## 1. Quand l'utiliser
- Séisme majeur (M ≥ 6.0)
- Ouragan / cyclone (cat. 3+)
- Inondations massives
- Conflits / déplacement massif de population
- Épidémie déclarée (cf. health.epidemic.alert)
- Cyberattaque massive avec impact citoyen

## 2. Activation immédiate (T+0)

### 2.1 Déclencher le workflow national
```bash
zbctl create instance escalation.crisis.national \
  --variables '{
    "trigger":"EARTHQUAKE",
    "intensity":"7.0",
    "regions":["Sud","Sud-Est","Grand-Anse"],
    "approver":"DG_ONI"
  }'
```

### 2.2 Auto-actions exécutées par le workflow
- Activation **mode offline-first national** (`offline.mode.enable`)
- Notification CIMO (Centre Interministériel de Mobilisation Opérationnelle)
- Pré-bascule DR (DC3 prêt)
- Surveillance renforcée toutes alertes Sev1/2 → Sev0
- Override SLA temporaire (auto-extension x 5) sur workflows non-vitaux

## 3. Workflows critiques à protéger (priorité absolue)

| Workflow | Priorité | Raison |
|----------|----------|--------|
| `civil-registry.death.disaster` | 🔴 #1 | Identification victimes |
| `civil-registry.death.standard` | 🔴 #2 | Déclarations courantes |
| `identity.verification.online` | 🔴 #3 | Vérifs urgences (hôpitaux, sauveteurs) |
| `identity.recovery.standard` | 🟠 #4 | Cartes perdues |
| `health.epidemic.alert` | 🟠 #5 | Surveillance sanitaire |
| `escalation.crisis.national` | 🔴 #6 | Coordination |

Workflows **temporairement suspendus** (économie ressources) :
- Mariages, divorces non-urgents
- Adoptions internationales
- Audits fiscaux

## 4. Mode décès en masse (procédure dédiée)

```bash
# 1. Activer la chaîne batch (DGPC → MoH → ONI)
zbctl create instance civil-registry.death.disaster --variables '{
  "disasterId":"HT-2026-EQ-SUD-001",
  "estimatedCount":1500,
  "sources":["DGPC","MoH","ICRC"]
}'

# 2. Batch d'ingestion (CSV/HL7)
./scripts/death-batch-ingest.sh --file /data/dgpc/victims-202605.csv \
  --sign-with hardware-kit-disaster-001
```

Garanties :
- Anti-doublons activés (un mort n'est déclaré qu'une fois)
- Croisement biométrie si dispo (ABIS)
- Signature batch PKI groupée (1 sig pour N actes pour vitesse)
- Audit complet préservé

## 5. Mode offline amplifié

- Tous les kits terrain passent en mode **pure offline** (pas d'attente réseau)
- Outbox local devient le buffer principal (capacité augmentée à 30 jours)
- Sync différée toutes les 6h via points relais (DGPC, MINUSTAH, etc.)
- Audit local Merkle (cf. BPMN `offline.audit.logs`) garantit l'immuabilité

## 6. Communication publique

| Heure | Action |
|-------|--------|
| T+15 min | Communiqué officiel Présidence ONI sur `transparence.snisid.ht` |
| T+30 min | Hotline citoyens : `+509 8XXX-XXXX` (24/7) |
| T+1 h | Briefing presse |
| T+24 h | Premier bilan |
| Quotidien | Tableau de bord public anonymisé |

## 7. Coordination inter-agences

| Agence | Rôle | API |
|--------|------|-----|
| **DGPC** | Décompte victimes | `civil-registry.death.disaster.v1` → producer |
| **MoH / OMS** | Surveillance épidémie | `health.epidemic.alert.v1` |
| **Police / MINUSTAH** | Sécurité, déplacés | `security.alert.raised.v1` |
| **Diaspora / MAE** | Recherche de personnes | `identity.verification.online` |
| **ICRC / ONG** | Aide humanitaire | `identity.verification.biometric` |

## 8. Pendant la crise — checklist quotidienne

- [ ] Backlog `offline.batch.uploaded.v1` traité
- [ ] Sync DC1↔DC2 OK (lag < 30 s)
- [ ] Audit Merkle vérifié (`audit.chain.verify`)
- [ ] Cellule fraude monitore les pics suspects
- [ ] Aucun workflow critique en breach SLA non escaladé
- [ ] DR DC3 prêt (`./scripts/dr-readiness.sh`)
- [ ] Communication du jour publiée
- [ ] Briefing WGO quotidien

## 9. Sortie de crise

1. **Désactivation graduelle** (sur décision DG ONI + CIMO).
2. Restauration des SLA standards.
3. Reprise des workflows suspendus.
4. **Audit national post-crise** (Cour des Comptes + indépendant).
5. **Post-mortem** présenté au Parlement.
6. Mise à jour des plans + nouveaux drills.

## 10. Mémorandum permanent

> En crise, la **continuité de l'identité** est la première dignité de l'État.
> Aucun citoyen ne doit perdre son identité parce qu'un serveur est tombé.
> SNISID DOIT continuer à fonctionner, fût-ce en mode dégradé, signé, et différé.
