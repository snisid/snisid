# 🛰️ SNISID — National Operations Model (24/7)

**Document N° :** SNISID-OPS-012
**Étape Phase 0 :** 12/16
**Principe :** *Les opérations doivent être permanentes.*

---

## 1. Vision

SNISID est une **infrastructure critique 24/7/365**. Une indisponibilité = blocage de naissances, mariages, vérifications d'identité, contrôles frontaliers, etc. Le modèle opérationnel doit garantir un service continu, observé, gouverné et résilient.

---

## 2. Architecture Opérationnelle

```
┌─────────────────────────────────────────────────────────┐
│                    AUTORITÉ NATIONALE NUMÉRIQUE         │
├─────────────────────────────────────────────────────────┤
│  ┌───────────┐    ┌───────────┐    ┌──────────────┐    │
│  │   NOC     │    │   SOC     │    │  Service Desk│    │
│  │ Network   │    │ Security  │    │   (L1/L2)    │    │
│  │ Operations│    │ Operations│    │              │    │
│  └─────┬─────┘    └─────┬─────┘    └──────┬───────┘    │
│        │                │                  │            │
│  ┌─────▼────────────────▼──────────────────▼──────┐    │
│  │         ÉQUIPES SRE / DEVOPS / DATA            │    │
│  └────────────────────────────────────────────────┘    │
│        │                                                │
│  ┌─────▼──────────────────────────────────────────┐    │
│  │      CRISIS MANAGEMENT (escalade direction)    │    │
│  └────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

---

## 3. NOC — Network Operations Center

**Mission :** surveiller la disponibilité, la performance, la capacité.

- **Effectif :** 3 shifts × 4 personnes = 12 ingénieurs + 1 lead NOC
- **Outils :** Grafana, Prometheus, Zabbix, NetBox, ThousandEyes (ou alt OSS)
- **Périmètre :** datacenters, réseau, K8s clusters, bases, edge nodes
- **SLO suivis :** voir Architecture (§9)

---

## 4. SOC — Security Operations Center

Voir document **Cybersecurity (07)**. Synthèse :
- 3 shifts × 8 analystes (Tier 1) + Tier 2/3 en heures ouvrables (astreinte)
- SIEM + SOAR + EDR + threat intel
- Liaison CSIRT-HT national

---

## 5. Service Desk (Support)

| Tier | Mission | Canaux |
|------|---------|--------|
| **L1** | Premier contact, qualification, FAQ | Tél, WhatsApp, USSD, email |
| **L2** | Diagnostic technique, résolution standard | Ticketing |
| **L3** | Expertise produit, escalade éditeur | — |

**Cibles SLA support :**
- Réponse L1 : < 5 min (vocal) / < 15 min (écrit)
- Résolution P1 (bloquant national) : < 1 h
- Résolution P2 (bloquant agence) : < 4 h
- Résolution P3 : < 24 h
- Résolution P4 : < 5 j

---

## 6. Incident Response — Cycle PRDR

```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ Préparer │──▶│ Détecter │──▶│ Répondre │──▶│ Restaurer│
└──────────┘   └──────────┘   └──────────┘   └────┬─────┘
      ▲                                            │
      └──────────── Post-mortem / RCA ─────────────┘
```

**Sévérités :**
| Sev | Description | Cellule activée |
|-----|-------------|------------------|
| SEV-1 | Indispo service critique national | NOC + SOC + Direction AND |
| SEV-2 | Dégradation majeure ou indispo agence | NOC + SOC |
| SEV-3 | Dégradation localisée | NOC ou SOC selon nature |
| SEV-4 | Incident mineur | Service Desk |

---

## 7. Crisis Escalation

**Cellule de crise** activable sous **1h** :
- Directeur AND (chef de cellule)
- CISO
- Head NOC
- Head SOC
- Responsable communication
- Représentant CNN si nécessaire
- Représentant ministre concerné

**Communication crise :**
- Status page publique (status.snisid.gouv.ht)
- Communiqué presse < 4h pour SEV-1
- Information citoyens via SMS + radio nationale
- Liaison directe avec ministères impactés

---

## 8. Astreintes & Garde

- Astreinte 24/7 par rôle critique
- Rotation hebdomadaire
- Prime d'astreinte conforme grille publique
- Outils on-call : OpsGenie / PagerDuty (OSS alt : Grafana OnCall)

---

## 9. Change Management

- **CAB** (Change Advisory Board) hebdomadaire
- Fenêtres de maintenance définies (nuits, dimanches)
- **Freeze** électoral (mois précédant scrutin)
- Tout changement = ticket + plan rollback + tests
- Approbation 4-yeux pour prod critique

---

## 10. Capacity Management

- Suivi trimestriel
- Forecast à 12 mois glissants
- Trigger d'achat à 70 % d'utilisation soutenue
- Plan d'extension capacité datacenter défini

---

## 11. Documentation & Runbooks

- Tous les runbooks dans wiki interne versionné (Git)
- Runbook obligatoire par alerte
- Tests trimestriels (game days)
- Onboarding documenté (J1 à J30 nouveau collaborateur)

---

## 12. KPI Opérationnels

| KPI | Cible |
|-----|-------|
| Disponibilité plateforme | ≥ 99,95 % |
| MTTD | < 5 min (alerte auto) |
| MTTR P1 | < 1 h |
| Taux de changements réussis | ≥ 98 % |
| Satisfaction utilisateurs internes | ≥ 85 % |
| Couverture runbook par alerte | 100 % |

---

## 13. Modèle RH cible

| Rôle | Effectif 2028 |
|------|---------------|
| Directeur AND + COMEX | 5 |
| NOC | 13 |
| SOC | 28 |
| Service Desk | 20 |
| SRE / DevOps | 25 |
| Développeurs core | 40 |
| Data / IA | 15 |
| Sécurité (hors SOC) | 8 |
| Légal / Conformité | 6 |
| Communication | 5 |
| Administration | 15 |
| **Total** | **~180** |

> Mix : haut potentiel haïtien (≥ 80 %) + 20 % expertise internationale temporaire avec transfert de compétences obligatoire.

---
*Fin du document — Étape 12/16*
