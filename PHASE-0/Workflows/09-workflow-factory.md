# ⚙️ SNISID — National Workflow Factory

**Document N° :** SNISID-WKF-009
**Étape Phase 0 :** 9/16
**Principe :** *Les workflows doivent être industrialisés nationalement.*

---

## 1. Vision

Mettre fin à la **dispersion artisanale** des procédures gouvernementales. Chaque acte administratif suit un **workflow modélisé, versionné, audité, mesuré** — produit par une **usine à workflows** nationale.

---

## 2. Plateforme Cible

| Composant | Choix recommandé |
|-----------|------------------|
| **BPMN Engine** | Camunda 8 (Zeebe) OSS ou Flowable |
| **DMN Engine** | Camunda DMN |
| **Form Builder** | Form.io ou Camunda Forms |
| **Workflow designer** | Camunda Modeler (BPMN 2.0) |
| **Tasklist** | Camunda Tasklist + intégration agents |
| **Notifications** | Service notifications central (SMS/Email/Push/USSD) |
| **Audit** | Event store immuable Kafka + ELK |

> Standards : **BPMN 2.0**, **DMN 1.3**, **CMMN 1.1** (cas complexes).

---

## 3. Domaines Industrialisés (5)

| Domaine | Volume estimé / an | Criticité |
|---------|---------------------|-----------|
| **État civil** | ~600 000 actes naissance + mariage + décès | CRITIQUE |
| **Justice** | ~150 000 procédures (casier, jugements) | CRITIQUE |
| **Police** | ~80 000 procédures (plaintes, recherches) | CRITIQUE |
| **Immigration** | ~500 000 passages frontière + 50 000 passeports | CRITIQUE |
| **Fiscalité** | ~2M déclarations + paiements | HAUTE |

---

## 4. Cycle de Vie d'un Workflow

```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ Analyse  │──▶│  Design  │──▶│ Validate │──▶│ Deploy   │
│ métier   │   │  BPMN    │   │ + Légal  │   │ versionné│
└──────────┘   └──────────┘   └──────────┘   └────┬─────┘
                                                    │
                                              ┌─────▼────┐
                                              │ Exécuter │
                                              └─────┬────┘
                                                    │
                                              ┌─────▼────┐
                                              │ Monitor  │
                                              │ + KPI    │
                                              └─────┬────┘
                                                    │
                                              ┌─────▼────┐
                                              │ Amélior. │
                                              │ continue │
                                              └──────────┘
```

---

## 5. Standards de Modélisation

Chaque workflow doit :
1. Avoir un **propriétaire métier** identifié
2. Référencer la **base légale** (loi, article)
3. Définir **SLA** par étape (temps max)
4. Avoir un **DMN** pour les règles de décision
5. Inclure les **cas d'exception** (refus, recours, escalade)
6. Émettre des **événements** sur le bus national
7. Être **observable** (KPI temps réel)
8. Avoir une version **offline** (formulaire mobile dégradé)

---

## 6. Bibliothèque de Workflows Standard

| Code | Workflow | Domaine |
|------|----------|---------|
| EC-N01 | Déclaration naissance simple | État civil |
| EC-N02 | Naissance par reconnaissance | État civil |
| EC-N03 | Naissance par déclaration tardive | État civil |
| EC-N04 | Naissance par décret | État civil |
| EC-N05 | Naissance par jugement (minutes) | État civil |
| EC-M01 | Mariage civil | État civil |
| EC-M02 | Mariage religieux transcrit | État civil |
| EC-D01 | Divorce par consentement mutuel | État civil |
| EC-D02 | Divorce contentieux | État civil |
| EC-X01 | Déclaration de décès | État civil |
| EC-A01 | Adoption simple | État civil |
| EC-A02 | Adoption plénière | État civil |
| ID-01 | Enrôlement CIN | Identité |
| ID-02 | Renouvellement CIN | Identité |
| ID-03 | Duplicata CIN | Identité |
| ID-04 | Rectification données | Identité |
| JU-01 | Demande casier judiciaire | Justice |
| JU-02 | Notification jugement | Justice |
| PO-01 | Dépôt plainte | Police |
| PO-02 | Recherche personne | Police |
| IM-01 | Demande passeport | Immigration |
| IM-02 | Passage frontière | Immigration |
| FI-01 | Déclaration fiscale | Fiscalité |

---

## 7. Pattern Standard d'un Workflow

```bpmn
Start (Citoyen / Agent)
   │
   ▼
[Vérification identité] (service SNISID Identity)
   │
   ▼
[Saisie formulaire] (Form.io)
   │
   ▼
[Vérifications automatiques] (DMN règles)
   │
   ├── KO ──▶ [Retour pour correction] ──▶ loop
   │
   ▼ OK
[Validation officier] (humaine)
   │
   ├── Refus ──▶ [Motivation + notification + recours]
   │
   ▼ Accepté
[Génération document signé] (PKI)
   │
   ▼
[Émission événement bus] (Kafka)
   │
   ▼
[Notification citoyen] (SMS/Email)
   │
   ▼
End
```

---

## 8. Industrialisation — Métriques

| KPI | Cible |
|-----|-------|
| Délai moyen acte naissance | < 24h (vs 30+ jours actuellement) |
| Délai casier judiciaire | < 48h |
| Taux d'automatisation décisions | ≥ 60 % |
| Taux d'erreur (re-traitement) | < 2 % |
| Satisfaction usagers | ≥ 80 % |

---

## 9. Gouvernance Workflow Factory

- **Comité Workflow** mensuel (DSI agences + AND)
- **Catalogue national** des workflows versionnés (Git)
- **Approbation 4-yeux** pour mise en prod
- **Rollback** automatisé en cas d'incident
- **A/B testing** possible pour amélioration

---
*Fin du document — Étape 9/16*
