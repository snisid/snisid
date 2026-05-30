# 🇭🇹 SNISID — PHASE 3 : NATIONAL WORKFLOW FACTORY

> **Système National d'Information et de Sécurité Identitaire d'Haïti**
> Industrialisation de l'administration publique haïtienne via workflows orchestrés, événementiels et auditables.

> ✅ **Phase 3 livrée à 100 %** — voir [`PHASE3-COMPLETE.md`](PHASE3-COMPLETE.md)

---

## 🎯 Mission de la Phase 3

Transformer l'administration publique haïtienne en **machine opérationnelle orchestrée, standardisée et industrialisée**.

| Avant | Après |
|-------|-------|
| Procédures manuelles | Workflows BPMN orchestrés |
| Chaos administratif | Standards nationaux versionnés |
| Dépendance humaine | Automatisation événementielle |
| Absence de traçabilité | Audit immuable + Event Sourcing |
| Lenteur | SLA mesurables |
| Corruption potentielle | PKI + signature + 4-eyes principle |
| Workflows non standardisés | Catalogue BPMN national (43 BPMN) |

---

## 📦 Livrables — 13 étapes (toutes ✅)

| # | Étape | Document | BPMN/Code |
|---|-------|----------|-----------|
| 1 | National Workflow Factory | [`docs/01-National-Workflow-Factory-Architecture.md`](docs/01-National-Workflow-Factory-Architecture.md) | `workflow-engine/` |
| 2 | Catalogue BPMN National | [`docs/02-National-BPMN-Catalog.md`](docs/02-National-BPMN-Catalog.md) | `BPMN/` |
| 3 | **BPMN État Civil** (16) | [`docs/03-BPMN-Etat-Civil.md`](docs/03-BPMN-Etat-Civil.md) | `BPMN/Civil-Registry/` |
| 4 | **Workflows Identité** (8) | [`docs/04-Workflows-Identite.md`](docs/04-Workflows-Identite.md) | `BPMN/Identity/` |
| 5 | **Workflows Judiciaires** (5) | [`docs/05-Workflows-Judiciaires.md`](docs/05-Workflows-Judiciaires.md) | `BPMN/Judicial/` |
| 6 | Modèle Event-Driven | [`docs/06-National-Event-Architecture.md`](docs/06-National-Event-Architecture.md) | — |
| 7 | Kafka Event Mesh | [`docs/07-Kafka-Event-Mesh.md`](docs/07-Kafka-Event-Mesh.md) | `kafka/` |
| 8 | Workflow Governance Office | [`docs/08-Workflow-Governance-Office.md`](docs/08-Workflow-Governance-Office.md) | — |
| 9 | SLA/SLO Nationaux | [`docs/09-SLA-SLO-Nationaux.md`](docs/09-SLA-SLO-Nationaux.md) | `sla/` |
| 10 | Runbooks Opérationnels (10) | [`docs/10-Runbooks-Operationnels.md`](docs/10-Runbooks-Operationnels.md) | `runbooks/` |
| 11 | Workflows Offline-First (5) | [`docs/11-Offline-Workflows.md`](docs/11-Offline-Workflows.md) | `BPMN/Offline/` |
| 12 | Workflow Observability | [`docs/12-Workflow-Observability-Model.md`](docs/12-Workflow-Observability-Model.md) | `observability/` |
| 13 | Référentiel BPMN centralisé | [`docs/13-BPMN-Repository.md`](docs/13-BPMN-Repository.md) | `scripts/lint-bpmn.sh` |

---

## 🗂 Arborescence Complète

```
SNISID-Phase3/
├── README.md
├── PHASE3-COMPLETE.md             ← Récapitulatif final
├── docs/                          (13 documents stratégiques)
├── BPMN/                          (43 BPMN nationaux)
│   ├── Civil-Registry/  (16)
│   ├── Identity/        (8)
│   ├── Judicial/        (5)
│   ├── Elections/       (2)
│   ├── Immigration/     (1)
│   ├── Tax/             (1)
│   ├── Health/          (1)
│   ├── Fraud/           (1)
│   ├── Audit/           (1)
│   ├── Offline/         (5)
│   └── Escalation/      (2)
├── workflow-engine/               (Temporal + Zeebe, TypeScript)
├── kafka/                         (topics, schemas Avro, ACLs)
├── sla/                           (sla-catalog.yaml)
├── runbooks/                      (10 runbooks + index)
├── observability/                 (Prometheus, Grafana, Loki, Tempo, OTel, Alertmanager)
└── scripts/                       (lint-bpmn.sh, sign-bpmn.sh)
```

---

## ⚖ Règle Absolue

> **Dans SNISID : tout devient workflow.**
> Aucune action administrative n'est exécutée hors d'un workflow versionné, audité, signé et gouverné.

---

## 🚨 Erreurs Absolument Évitées

| Erreur | Garde-fou implémenté |
|--------|----------------------|
| BPMN trop simples | Lint automatique CI ; 10 règles ; **0 erreurs sur 43 BPMN** |
| Pas d'escalade | `escalation.sla.breach` partout ; SLA boundary events |
| Pas de Kafka | Tout workflow émet via `kafka.emit` |
| Pas de versioning | `zeebe:versionTag` SemVer obligatoire |
| Pas d'audit | Merkle chain + WORM 30 ans + alerte Sev0 si rupture |
| Pas de rollback | Runbook 03 + WGO Rollback Cell |

---

## 🚀 Démarrage Rapide

```bash
# 1. Vérifier les BPMN
./scripts/lint-bpmn.sh

# 2. Démarrer la stack d'observabilité localement
cd observability/deploy
docker compose -f docker-compose.observability.yml up -d
open http://localhost:3000   # Grafana (admin / snisid-dev)

# 3. Démarrer le workflow engine
cd ../../workflow-engine
npm install && npm run worker:zeebe
```

---

## 🔗 Phases

| Phase | Statut |
|-------|--------|
| Phase 1 (Identité + PKI + biométrie) | ✅ Préalable |
| Phase 2 (Cybersécurité + souveraineté) | ✅ Préalable |
| **Phase 3 (National Workflow Factory)** | ✅ **LIVRÉ** |
| Phase 4 (Interop internationale + ICAO + diaspora) | ➡ À venir |

---

**Maintenu par :** Workflow Governance Office — Office National d'Identification (ONI)
