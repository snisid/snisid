# SNISID — Documentation des Workflows de Gouvernance
## Système National d'Identité et de Services d'Identité Digitale

---

| Métadonnée | Valeur |
|---|---|
| **Document ID** | SNISID-GOV-WRK-001 |
| **Version** | 1.0.0 |
| **Statut** | APPROUVÉ — EN VIGUEUR |
| **Date de création** | 2026-05-25 |
| **Date de révision** | 2026-11-25 |
| **Classification** | GOUVERNANCE / USAGE INTERNE |
| **Propriétaire** | AND — Directeur Architecture & Innovation |
| **Révisé par** | CCB + IGB |
| **Approuvé par** | DG AND / CNN |
| **Références** | SNISID-GOV-ORG-001, SNISID-GOV-RACI-001, ISO 27001, ITIL v4 |

---

> **RAPPEL** : Ces workflows sont des procédures officielles de gouvernance. Leur non-respect constitue une faute de gestion documentée. En situation de crise, les workflows d'urgence prévalent sur les workflows standards.

---

## TABLE DES MATIÈRES

1. [Workflow 1 — Décision d'Architecture et Changements Majeurs](#workflow-1--décision-darchitecture-et-changements-majeurs)
2. [Workflow 2 — Protocole d'Urgence Gouvernance (Crise Cyber / Catastrophe)](#workflow-2--protocole-durgence-gouvernance-crise-cyber--catastrophe)
3. [Workflow 3 — Approbation Budgétaire](#workflow-3--approbation-budgétaire)
4. [Workflow 4 — Partenariat International](#workflow-4--partenariat-international)
5. [Workflow 5 — Chaîne d'Escalade d'Incidents (L1 → L2 → L3 → CNN)](#workflow-5--chaîne-descalade-dincidents-l1--l2--l3--cnn)
6. [Workflow 6 — Validation de Conformité Légale](#workflow-6--validation-de-conformité-légale)
7. [Workflow 7 — Changement Technique (RFC Process)](#workflow-7--changement-technique-rfc-process)
8. [Workflow 8 — Traitement de Plainte Citoyenne](#workflow-8--traitement-de-plainte-citoyenne)
9. [Workflow 9 — Accréditation de Fournisseur](#workflow-9--accréditation-de-fournisseur)
10. [Workflow 10 — Révocation d'Identité pour Fraude](#workflow-10--révocation-didentité-pour-fraude)
11. [Indicateurs de Performance des Workflows (KPIs)](#indicateurs-de-performance-des-workflows-kpis)

---

## WORKFLOW 1 — Décision d'Architecture et Changements Majeurs

### 1.1 Périmètre

Ce workflow s'applique à toute modification impactant :
- L'architecture technique fondamentale du SNISID
- Les protocoles de communication inter-systèmes
- L'infrastructure PKI, X-Road, DataCenter
- Les algorithmes biométriques ou cryptographiques
- Les API publiques et protocoles d'interopérabilité
- Les schémas de données du registre national d'identité

**Seuil d'activation** : Tout changement classifié **Majeur (M)** ou **Stratégique (ST)** par le CCB.

### 1.2 Acteurs

| Acteur | Rôle |
|---|---|
| **Demandeur** | Tout organe, agence, ou fournisseur autorisé |
| **CCB** | Évaluation technique et première décision |
| **AND** | Validation stratégique |
| **IGB** | Impact sur l'identité |
| **NDPA** | Impact sur les données personnelles |
| **ETH** | Impact éthique IA |
| **CNN** | Décision finale pour changements stratégiques |

### 1.3 Diagramme Mermaid — Workflow Architecture

```mermaid
flowchart TD
    INIT["📋 Demande de Changement\nArchitectural Soumise\n[RFC Formulaire CCB]"]
    
    TRIAGE["🔍 CCB — Triage Initial\n(72h max)\nClassification: N / M / ST"]
    
    Q_CLASS{"Classification?"}
    
    NORMAL["Workflow Normal\n→ SNISID-GOV-WRK-007\nChangement Technique"]
    
    MAJEUR_PHASE["📊 Phase Analyse Majeure\n(21 jours)"]
    
    ARCH_REVIEW["🏗️ Architecture Review\nDirecteur Architecture AND\n(5 j)"]
    
    IGB_REVIEW["🆔 Revue IGB\nImpact identité\n(5 j)"]
    
    NDPA_REVIEW["🔒 Revue NDPA\nImpact données\n(5 j)"]
    
    ETH_REVIEW["⚖️ Revue Éthique\nImpact droits fondamentaux\n(5 j)"]
    
    CCB_VOTE["⚙️ Vote CCB\n(7/10 requis pour Majeur)"]
    
    Q_CCB{"CCB Approuve?"}
    
    REJECT_CCB["❌ Rejet CCB\nNotification demandeur\nDocumentation motifs"]
    
    AND_VALIDATE["🎯 Validation AND\nDG AND signature\n(5 j)"]
    
    Q_STRAT{"Classifié\nStratégique?"}
    
    AND_DECIDE["✅ Décision AND\nArrêté AND signé\nNotification organes"]
    
    CNN_BRIEF["📑 Brief CNN Préparé\n14 jours avant session"]
    
    CNN_VOTE["⭐ CNN Vote\n(majorité 2/3 requise)\nSession plénière"]
    
    Q_CNN{"CNN Approuve?"}
    
    REJECT_CNN["❌ Rejet CNN\nDocument de raisons\nPossibilité révision"]
    
    CNN_RESOLUTION["✅ Résolution CNN\nNumérotée et signée"]
    
    IMPLEM["⚡ Implémentation\nSous supervision CCB\n+ Tests + Validation"]
    
    PIR["📝 Post-Implementation\nReview (PIR)\n30 jours après déploiement"]
    
    ARCHIVE["🗂️ Archivage\nRegistre des changements\n30 ans"]

    INIT --> TRIAGE
    TRIAGE --> Q_CLASS
    Q_CLASS -->|"Normale (N)"| NORMAL
    Q_CLASS -->|"Majeure (M)"| MAJEUR_PHASE
    Q_CLASS -->|"Stratégique (ST)"| MAJEUR_PHASE
    MAJEUR_PHASE --> ARCH_REVIEW
    MAJEUR_PHASE --> IGB_REVIEW
    MAJEUR_PHASE --> NDPA_REVIEW
    MAJEUR_PHASE --> ETH_REVIEW
    ARCH_REVIEW --> CCB_VOTE
    IGB_REVIEW --> CCB_VOTE
    NDPA_REVIEW --> CCB_VOTE
    ETH_REVIEW --> CCB_VOTE
    CCB_VOTE --> Q_CCB
    Q_CCB -->|"NON"| REJECT_CCB
    Q_CCB -->|"OUI"| AND_VALIDATE
    AND_VALIDATE --> Q_STRAT
    Q_STRAT -->|"NON — Majeure"| AND_DECIDE
    Q_STRAT -->|"OUI — Stratégique"| CNN_BRIEF
    CNN_BRIEF --> CNN_VOTE
    CNN_VOTE --> Q_CNN
    Q_CNN -->|"NON"| REJECT_CNN
    Q_CNN -->|"OUI"| CNN_RESOLUTION
    CNN_RESOLUTION --> IMPLEM
    AND_DECIDE --> IMPLEM
    IMPLEM --> PIR
    PIR --> ARCHIVE

    style INIT fill:#16213e,color:#fff
    style CNN_RESOLUTION fill:#533483,color:#fff
    style AND_DECIDE fill:#0f3460,color:#fff
    style REJECT_CCB fill:#e94560,color:#fff
    style REJECT_CNN fill:#e94560,color:#fff
```

### 1.4 Vue ASCII — Flux Simplifié

```
[DEMANDEUR] → Formulaire RFC
      ↓
[CCB] → Triage (72h) → Classification N / M / ST
      ↓
    [M/ST] ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
      ↓                                             ║
[Revues parallèles — 5j max chacune]                ║
  ├─ Arch Review (AND)                               ║
  ├─ IGB Review                                      ║
  ├─ NDPA Review                                     ║
  └─ ETH Review                                      ║
      ↓                                             ║
[CCB Vote] ← 7/10 membres → REJETÉ ←━━━━━━━━━━━━━━━┛
      ↓ (Approuvé)
[AND Validation] → DG AND signe
      ↓
    [ST?]
  OUI ↓        NON ↓
[CNN Plénière]  [Arrêté AND]
  2/3 vote         ↓
      ↓         [Implémentation CCB]
[Résolution CNN]    ↓
      ↓         [PIR 30j]
[Implémentation]    ↓
      ↓         [Archivage 30 ans]
[PIR 30j]
      ↓
[Archivage 30 ans]
```

### 1.5 Délais et SLA

| Étape | SLA Standard | SLA Urgent |
|---|---|---|
| Soumission → Triage CCB | 3 jours ouvrables | 24h |
| Revues parallèles | 5 jours ouvrables | 48h |
| Vote CCB | 2 jours après revues | 4h |
| Validation AND | 5 jours après CCB | 24h |
| Brief CNN | 14 jours avant session | — |
| Vote CNN | Session suivante (max 90j) | Session extraordinaire (72h) |
| Délai total (Majeure) | ~21 jours | ~5 jours |
| Délai total (Stratégique) | ~60 jours | ~15 jours |

---

## WORKFLOW 2 — Protocole d'Urgence Gouvernance (Crise Cyber / Catastrophe)

### 2.1 Périmètre et Déclencheurs

Ce protocole s'active en cas de :
- **Cyberattaque** classifiée P1 ou P0 par le SOC/CERT-HT
- **Catastrophe naturelle** (séisme >7.0, ouragan catégorie 4+) affectant l'infrastructure SNISID
- **Compromission** du registre national d'identité
- **Défaillance généralisée** des services SNISID (>50% indisponible)
- **Décision Présidentielle** d'activation de l'état d'urgence numérique

### 2.2 Phases du Protocole d'Urgence

```
PHASE ALPHA (T+0 à T+2h) — DÉTECTION ET ALERTE
PHASE BRAVO (T+2h à T+6h) — ACTIVATION ET MOBILISATION
PHASE CHARLIE (T+6h à T+24h) — OPÉRATIONS D'URGENCE
PHASE DELTA (T+24h à Résolution) — RÉPONSE PROLONGÉE
PHASE ECHO (Post-incident) — RÉTABLISSEMENT ET LEÇONS
```

### 2.3 Diagramme Mermaid — Protocole d'Urgence

```mermaid
sequenceDiagram
    autonumber
    participant DET as 🔍 Détecteur<br/>(SOC Tier 1)
    participant SOC as 🛡️ SOC/CERT-HT<br/>Directeur
    participant CISO as 🔐 CISO AND
    participant DG as 🎯 DG AND
    participant CELL as 🆘 Cellule de Crise<br/>AND
    participant CNN as ⭐ CNN
    participant PRES as 🏛️ Président
    participant PUB as 📢 Public

    Note over DET,PUB: ═══ PHASE ALPHA — DÉTECTION (T+0 à T+2h) ═══

    DET->>SOC: Alerte P1/P0 détectée (T+0)
    SOC->>SOC: Confirmation et classification (T+15min)
    SOC->>CISO: Notification critique (T+15min)
    CISO->>DG: Briefing d'urgence (T+30min)
    DG->>CELL: Convocation Cellule de Crise (T+45min)
    CELL->>CELL: Premier bilan de situation (T+1h)
    DG->>CNN: Notification et convocation (T+1h30)

    Note over DET,PUB: ═══ PHASE BRAVO — ACTIVATION (T+2h à T+6h) ═══

    CNN->>CNN: Session extraordinaire (T+2h)
    CNN->>PRES: Recommandation État d'Urgence (T+2h30)
    PRES->>CNN: Décret d'Urgence Numérique (T+3h)
    CNN->>DG: Pouvoirs d'urgence délégués (T+3h15)
    DG->>SOC: Activation Plan CRISNUM-01 (T+3h30)
    CELL->>PUB: Première communication publique (T+4h)

    Note over DET,PUB: ═══ PHASE CHARLIE — OPÉRATIONS (T+6h à T+24h) ═══

    SOC->>DG: Rapport situation toutes les 4h
    DG->>CNN: Rapport situation toutes les 4h
    CNN->>CNN: Session permanente
    CELL->>CELL: Gestion ressources + coordination

    Note over DET,PUB: ═══ PHASE DELTA — RÉPONSE PROLONGÉE (>T+24h) ═══

    DG->>CNN: Rapport quotidien 08h00
    CNN->>PRES: Mise à jour quotidienne
    SOC->>CELL: Rapport opérationnel continu

    Note over DET,PUB: ═══ PHASE ECHO — RÉTABLISSEMENT ═══

    SOC->>DG: Déclaration "fin d'incident" proposée
    DG->>CNN: Demande levée État d'Urgence
    CNN->>PRES: Recommandation levée
    PRES->>CNN: Levée Décret d'Urgence
    CELL->>CELL: Rapport Post-Incident (PIR 30j)
    CELL->>CNN: PIR final + recommandations
```

### 2.4 Vue ASCII — Chronologie d'Urgence

```
T+0      ┌─────────────────────────────────────────────┐
         │ DÉTECTION : Alerte P1/P0 SOC Tier 1         │
T+15min  │ CONFIRMATION : SOC Directeur classifie       │
T+15min  │ NOTIFICATION : CISO AND alerté              │
T+30min  │ BRIEFING : DG AND informé                   │
T+45min  │ CELLULE : Convocation Cellule de Crise       │
T+1h     │ BILAN : Premier état de situation            │
T+1h30   │ CNN : Notification + convocation urgente     │
         └─────────────────────────────────────────────┘
T+2h     ┌─────────────────────────────────────────────┐
         │ SESSION CNN EXTRAORDINAIRE                    │
T+2h30   │ RECOMMANDATION : État d'Urgence Numérique   │
T+3h     │ DÉCRET PRÉSIDENTIEL : État d'Urgence         │
T+3h15   │ DÉLÉGATION : Pouvoirs d'urgence → DG AND    │
T+3h30   │ ACTIVATION : Plan CRISNUM-01                 │
T+4h     │ COMMUNICATION : Premier message public       │
         └─────────────────────────────────────────────┘
T+6h+    ┌─────────────────────────────────────────────┐
         │ OPÉRATIONS CONTINUES                          │
         │ Rapports SOC → DG → CNN : toutes les 4h     │
         │ CNN en session permanente                     │
         │ Communications publiques régulières          │
         └─────────────────────────────────────────────┘
```

### 2.5 Équipe de Crise et Contacts d'Urgence

| Rôle | Responsable | Backup | Canal Urgence |
|---|---|---|---|
| **Coordinateur Crise** | DG AND | DGA AND | Téléphone sécurisé H24 |
| **Responsable Technique** | CISO | Directeur Architecture | SOC Hotline |
| **Porte-Parole** | DGA AND (désigné) | Directeur Juridique | — |
| **Liaison CNN** | Directeur Juridique | DGA AND | Messagerie sécurisée |
| **Coordination Internationale** | Dir. Partenariats | — | Canal diplomatique |
| **Opérations Terrain** | COO AND | — | Radio sécurisée |

### 2.6 Plans de Continuité Activés

| Plan | Déclencheur | Durée max |
|---|---|---|
| **CRISNUM-01** (Cyber Critique) | P1/P0 SOC | Jusqu'à résolution |
| **CRISNUM-02** (Catastrophe Naturelle) | Séisme / Ouragan | Jusqu'à rétablissement |
| **CRISNUM-03** (Pandémie) | Impossibilité opérationnelle | Durée crise |
| **CRISNUM-04** (Défaillance Infra) | >50% services down | 72h max pour rétablissement |

---

## WORKFLOW 3 — Approbation Budgétaire

### 3.1 Périmètre

Applicable à tous les engagements budgétaires du SNISID selon les seuils suivants :

| Seuil | Approbateur |
|---|---|
| < 500,000 HTG | COO AND (délégation) |
| 500K – 5M HTG | CFO AND |
| 5M – 50M HTG | DG AND |
| 50M – 500M HTG | DG AND + CNN Information |
| > 500M HTG | CNN Vote + Parlement |
| Budget pluriannuel (tout montant) | CNN Vote + Parlement |

### 3.2 Diagramme Mermaid — Chaîne Budgétaire

```mermaid
flowchart TD
    BESOIN["💡 Identification du Besoin\nPar organe ou direction AND"]
    
    DOSSIER["📋 Constitution du Dossier\nBusinessCase + Estimation\n+ RACI + Risques"]
    
    Q_SEUIL{"Montant estimé?"}
    
    ROUTINE["🟢 Procédure Délégation\n< 5M HTG\nCFO/COO AND"]
    
    AND_BUDGET["🎯 Revue AND\nDG + CFO + COO\n10 jours ouvrables"]
    
    Q_AND{"AND Approuve\n(5M-50M)?"}
    
    ARR_AND["✅ Arrêté AND\nEngagement budgétaire\nNotification MEF"]
    
    CNN_BUDGET["⭐ Dossier CNN\nPréparation 14j\nBudget > 50M HTG"]
    
    MEF_AVIS["🏦 Avis MEF\n(si >50M HTG)\n15 jours"]
    
    CNN_VOTE["⭐ Vote CNN\nMajorité simple\nSession ordinaire"]
    
    Q_CNN_B{"Montant >\n500M HTG?"}
    
    CNN_APPROVE["✅ Résolution CNN\nEngagement approuvé"]
    
    PARLEMENT["🏛️ Parlement\nLoi de Finances\nou Rectificative"]
    
    APPEL_OFFRE["📢 Appel d'Offres\nPublic ou Restreint\nSelon seuils ARMP"]
    
    EXECUTION["⚡ Exécution Budgétaire\nSuivi MEF\nRapport trimestriel AND"]
    
    AUDIT["🔍 Audit\nCour Supérieure des Comptes\nAnnuel"]

    BESOIN --> DOSSIER
    DOSSIER --> Q_SEUIL
    Q_SEUIL -->|"< 5M HTG"| ROUTINE
    Q_SEUIL -->|"5M - 50M HTG"| AND_BUDGET
    Q_SEUIL -->|"> 50M HTG"| CNN_BUDGET
    AND_BUDGET --> Q_AND
    Q_AND -->|"NON"| BESOIN
    Q_AND -->|"OUI 5M-50M"| ARR_AND
    CNN_BUDGET --> MEF_AVIS
    MEF_AVIS --> CNN_VOTE
    CNN_VOTE --> Q_CNN_B
    Q_CNN_B -->|"NON (<500M)"| CNN_APPROVE
    Q_CNN_B -->|"OUI (>500M)"| PARLEMENT
    PARLEMENT --> CNN_APPROVE
    CNN_APPROVE --> APPEL_OFFRE
    ARR_AND --> APPEL_OFFRE
    ROUTINE --> APPEL_OFFRE
    APPEL_OFFRE --> EXECUTION
    EXECUTION --> AUDIT

    style BESOIN fill:#16213e,color:#fff
    style CNN_APPROVE fill:#533483,color:#fff
    style PARLEMENT fill:#1a1a2e,color:#fff
    style AUDIT fill:#2c7873,color:#fff
```

### 3.3 Vue ASCII — Seuils Budgétaires

```
MONTANT          APPROBATEUR          DÉLAI      PROCÉDURE
─────────────────────────────────────────────────────────
< 500K HTG   → COO/CFO AND          3 j       Délégation interne
500K-5M HTG  → CFO AND              5 j       Arrêté CFO
5M-50M HTG   → DG AND               10 j      Arrêté DG AND
50M-500M HTG → DG AND + CNN (info)  30 j      Résolution CNN
> 500M HTG   → CNN + Parlement      90+ j     Loi de Finances
Pluriannuel  → CNN + Parlement      Session   Loi de Finances
─────────────────────────────────────────────────────────
```

---

## WORKFLOW 4 — Partenariat International

### 4.1 Périmètre

Applicable à tout accord avec :
- États étrangers (bilatéral)
- Organisations internationales (ONU, PNUD, BID, BM, OEA, CARICOM)
- Agences de coopération technique (GIZ, USAID, AFD, etc.)
- Entités privées internationales (accord stratégique)

### 4.2 Classification des Partenariats

| Catégorie | Description | Approbateur |
|---|---|---|
| **Technique** | Assistance technique, formation, logiciels | AND DG |
| **Opérationnel** | Fourniture de services intégrés dans SNISID | AND + CNN information |
| **Stratégique** | Accord cadre, traité bilatéral, financement majeur | CNN vote |
| **Souverain** | Traité impliquant souveraineté des données nationales | CNN 3/4 + Parlement |

### 4.3 Diagramme Mermaid — Workflow Partenariat

```mermaid
flowchart TD
    INITIAT["🌐 Initiative de Partenariat\n(Interne AND ou Proposition Externe)"]
    
    EVAL_INIT["📊 Évaluation Initiale\nDir. Partenariats AND\n(10 jours)"]
    
    Q_CAT{"Catégorie du\nPartenariat?"}
    
    TECH["🔧 Partenariat Technique\nNégociation AND\nSignature DG AND"]
    
    DOSSIER_CNN["📑 Dossier Complet\npour CNN\n- Analyse juridique\n- Analyse risques\n- Analyse données (NDPA)\n- Avis éthique"]
    
    NDPA_CHECK["🔒 NDPA — Revue\nTransfert données\n(30 jours)"]
    
    MJP_AVIS["⚖️ Ministère Justice\nAvis légal\n(21 jours)"]
    
    MAEC_COORD["🌍 MAEC — Coordination\nDiplomatique\n(14 jours)"]
    
    AND_NEGOC["🎯 Négociation AND\n+ Parties prenantes\n(Variable)"]
    
    CNN_STRAT["⭐ CNN — Partenariat\nStratégique\nMajorité simple"]
    
    CNN_SOUV["⭐ CNN — Partenariat\nSouverain\n3/4 + Parlement"]
    
    SIGNATURE["✍️ Signature Officielle\nCNN Président ou\nDG AND (selon niveau)"]
    
    PUBLI["📢 Publication\nJournal Officiel\n(si requis)"]
    
    SUIVI["📈 Suivi Annuel\nRapport AND → CNN"]
    
    REVUE["🔄 Revue Triennale\nou à la demande"]

    INITIAT --> EVAL_INIT
    EVAL_INIT --> Q_CAT
    Q_CAT -->|"Technique"| TECH
    Q_CAT -->|"Opérationnel / Stratégique / Souverain"| DOSSIER_CNN
    DOSSIER_CNN --> NDPA_CHECK
    DOSSIER_CNN --> MJP_AVIS
    DOSSIER_CNN --> MAEC_COORD
    NDPA_CHECK --> AND_NEGOC
    MJP_AVIS --> AND_NEGOC
    MAEC_COORD --> AND_NEGOC
    AND_NEGOC --> CNN_STRAT
    AND_NEGOC --> CNN_SOUV
    CNN_STRAT --> SIGNATURE
    CNN_SOUV --> SIGNATURE
    TECH --> SIGNATURE
    SIGNATURE --> PUBLI
    PUBLI --> SUIVI
    SUIVI --> REVUE
```

### 4.4 Délais Standard

| Étape | Délai |
|---|---|
| Évaluation initiale AND | 10 jours |
| Revue NDPA | 30 jours |
| Avis MJP | 21 jours |
| Coordination MAEC | 14 jours |
| Négociation (variable) | 30-180 jours |
| Brief CNN | 14 jours avant session |
| Vote CNN (session ordinaire) | ~90 jours max |
| **TOTAL Partenariat Stratégique** | **~6-12 mois** |

---

## WORKFLOW 5 — Chaîne d'Escalade d'Incidents (L1 → L2 → L3 → CNN)

### 5.1 Définition des Niveaux

| Niveau | Nom | Acteur | Périmètre |
|---|---|---|---|
| **L1** | Support Standard | SOC Tier 1 / OJRNH | Incidents mineurs, requêtes courantes |
| **L2** | Investigation | SOC Tier 2 / CISO | Incidents confirmés, impact limité |
| **L3** | Réponse Avancée | SOC Tier 3 / AND | Incidents majeurs, impact systémique |
| **CNN** | Crise Nationale | CNN / Président | Crise nationale, état d'urgence |

### 5.2 Diagramme Mermaid — Chaîne d'Escalade

```mermaid
flowchart TD
    DETECTION["🔍 Détection Événement\nMonitoring SOC H24\nAlerte automatique ou manuelle"]
    
    L1["📊 NIVEAU L1\nSOC TIER 1\n• Triage initial\n• Classification P5-P3\n• Documentation\n• Résolution si possible\nDélai: < 30 min"]
    
    Q_L1{"Résolu\nen L1?"}
    
    L1_CLOSE["✅ Ticket Fermé L1\nDocumentation\nRapport quotidien"]
    
    L2["🔬 NIVEAU L2\nSOC TIER 2 + CISO\n• Investigation approfondie\n• Analyse forensique\n• Classification P3-P2\n• Notification AND\nDélai: < 2h"]
    
    Q_L2{"Résolu\nen L2?"}
    
    L2_CLOSE["✅ Ticket Fermé L2\nRapport incident\nRetour expérience"]
    
    L3["🚨 NIVEAU L3\nSOC TIER 3 + AND CELLULE\n• Crise active\n• Classification P2-P1\n• Pouvoirs étendus\n• Coordination nationale\nDélai: Continu"]
    
    Q_L3{"Impact\nnational?"}
    
    L3_CLOSE["✅ Ticket Fermé L3\nPIR 30 jours\nRecommandations"]
    
    CNN_LEVEL["⭐ NIVEAU CNN\nSESSION EXTRAORDINAIRE\n• État d'Urgence Numérique\n• Classification P1-P0\n• Décret Présidentiel\n• Coordination gouvernementale"]
    
    CNN_CLOSE["✅ Résolution Crise\nRapport CNN\nLeçons apprises\nMise à jour procédures"]
    
    NOTIF_AND["📢 Notification AND\n(30 min si P3+)"]
    NOTIF_CNN["📢 Notification CNN\n(2h si P1/P0)"]
    
    PIR["📝 Post-Incident Review\nMandatoire tous niveaux\nL1: 24h\nL2: 7j\nL3/CNN: 30j"]

    DETECTION --> L1
    L1 --> Q_L1
    Q_L1 -->|"OUI"| L1_CLOSE
    Q_L1 -->|"NON (P3+)"| L2
    L1 --> NOTIF_AND
    L2 --> Q_L2
    Q_L2 -->|"OUI"| L2_CLOSE
    Q_L2 -->|"NON (P2+)"| L3
    L3 --> Q_L3
    Q_L3 -->|"NON"| L3_CLOSE
    Q_L3 -->|"OUI (P1/P0)"| CNN_LEVEL
    L3 --> NOTIF_CNN
    CNN_LEVEL --> CNN_CLOSE
    L1_CLOSE --> PIR
    L2_CLOSE --> PIR
    L3_CLOSE --> PIR
    CNN_CLOSE --> PIR

    style L1 fill:#2d6a4f,color:#fff
    style L2 fill:#f6a623,color:#000
    style L3 fill:#e94560,color:#fff
    style CNN_LEVEL fill:#533483,color:#fff
```

### 5.3 Vue ASCII — Chaîne d'Escalade

```
ÉVÉNEMENT DÉTECTÉ
      │
      ▼
┌─────────────────────────────────┐
│         L1 — SOC Tier 1         │ P5-P4 : Nominal / Vigilance
│ Délai résolution : < 30 min     │ ← 80% des incidents résolus ici
│ Rapport : Quotidien 07h00       │
└──────────────┬──────────────────┘
               │ Non résolu ou P3+
               ▼ Notification AND 30min
┌─────────────────────────────────┐
│     L2 — SOC Tier 2 + CISO     │ P3 : Alerte
│ Délai investigation : < 2h      │ ← 15% des incidents résolus ici
│ Rapport : Immédiat au DG AND   │
└──────────────┬──────────────────┘
               │ Non résolu ou P2+
               ▼
┌─────────────────────────────────┐
│    L3 — SOC Tier 3 + AND       │ P2-P1 : Grave / Critique
│ Cellule de crise activée        │ ← 4% des incidents
│ Coordination nationale          │ Notification CNN si P1
└──────────────┬──────────────────┘
               │ P1/P0 — Impact national
               ▼ Notification CNN 2h
┌─────────────────────────────────┐
│      CNN — SESSION URGENCE      │ P1-P0 : Critique / Guerre Cyber
│ État d'Urgence Numérique        │ ← 1% des incidents
│ Décret Présidentiel             │
└─────────────────────────────────┘
```

### 5.4 SLA par Niveau

| Niveau | Délai Première Réponse | Délai Résolution | Notification |
|---|---|---|---|
| **L1** | 5 minutes | 30 minutes | Rapport quotidien |
| **L2** | 15 minutes | 2 heures | AND dans les 30 min |
| **L3** | 30 minutes | Indéfini (continu) | AND immédiat + CNN si P1 (2h) |
| **CNN** | 2 heures | Indéfini | Parlement si >72h |

---

## WORKFLOW 6 — Validation de Conformité Légale

### 6.1 Périmètre

Applicable à :
- Tout nouveau système ou processus SNISID
- Toute modification législative ou réglementaire
- Tout accord ou contrat impliquant des données personnelles
- Toute nouvelle fonctionnalité IA ou traitement automatisé

### 6.2 Diagramme Mermaid — Conformité Légale

```mermaid
flowchart TD
    SOUMIS["📋 Soumission pour Validation\nLégale\n(Système, processus, contrat, loi)"]
    
    SCREENING["🔍 Screening Initial\nDirecteur Juridique AND\n(5 jours)"]
    
    Q_SCOPE{"Scope de\nvalidation?"}
    
    LEGAL_INT["⚖️ Revue Légale Interne\nDirecteur Juridique AND\n(10 jours)\n- Constitution\n- Lois nationales\n- Décrets en vigueur"]
    
    LEGAL_EXT["🌐 Revue Légale Externe\n+ Internationale\n(15 jours)\n- Convention 108+\n- Budapest Convention\n- ICAO 9303\n- Traités ratifiés"]
    
    NDPA_LEGAL["🔒 Revue NDPA\n(obligatoire si données)\n(30 jours)\nDPIA si requis"]
    
    ETH_LEGAL["⚖️ Revue Éthique\n(si IA / traitement automatisé)\n(30 jours)\nAvis formel"]
    
    MJP_CONSULT["⚖️ Consultation MJP\n(si impact légal majeur)\n(21 jours)"]
    
    SYNTH["📊 Rapport de Conformité\nDirecteur Juridique AND\nGaps + Recommandations"]
    
    Q_CONFORM{"Conforme\nsans réserve?"}
    
    Q_GAPS{"Gaps\ncritiques?"}
    
    APPRO["✅ Attestation de\nConformité SNISID\nNo. SNISID-LEG-ATT-YYYY-NNN"]
    
    COND["⚠️ Conformité Sous\nConditions\nPlan d'action 90j"]
    
    BLOCAGE["❌ Blocage Légal\nSuspension du projet\nPlan de remédiation requis"]
    
    CNN_LEGAL["⭐ CNN — Si Gap légal\nfondamental\nInitiative législative"]
    
    REGISTRE["🗂️ Registre de Conformité\nMis à jour\nAudit annuel NDPA"]

    SOUMIS --> SCREENING
    SCREENING --> Q_SCOPE
    Q_SCOPE -->|"National"| LEGAL_INT
    Q_SCOPE -->|"International"| LEGAL_EXT
    Q_SCOPE -->|"Données personnelles"| NDPA_LEGAL
    Q_SCOPE -->|"IA / Automatisé"| ETH_LEGAL
    LEGAL_INT --> SYNTH
    LEGAL_EXT --> SYNTH
    NDPA_LEGAL --> SYNTH
    ETH_LEGAL --> SYNTH
    MJP_CONSULT --> SYNTH
    SYNTH --> Q_CONFORM
    Q_CONFORM -->|"OUI"| APPRO
    Q_CONFORM -->|"NON"| Q_GAPS
    Q_GAPS -->|"Non critiques"| COND
    Q_GAPS -->|"Critiques"| BLOCAGE
    BLOCAGE --> CNN_LEGAL
    APPRO --> REGISTRE
    COND --> REGISTRE
```

### 6.3 Checklist de Conformité SNISID

| Exigence | Texte de Référence | Statut |
|---|---|---|
| Protection données personnelles | Convention 108+ / Loi SNISID | À valider |
| Droits du citoyen (accès, rectification) | Art. 36 Constitution / Loi Protection Données | À valider |
| Sécurité des systèmes | ISO 27001 / Budapest Convention | À valider |
| Biométrie — consentement | Loi Biométrie (en cours) | En attente loi |
| Signatures électroniques | Loi Signature Électronique (en cours) | En attente loi |
| Interopérabilité | Loi Interopérabilité (en cours) | En attente loi |
| Continuité de service | ISO 22301 / Plan BCP | À valider |
| Documents de voyage | ICAO 9303 | À valider |
| ODD 16.9 — Identité légale pour tous | Agenda 2030 ONU | En cours |

---

## WORKFLOW 7 — Changement Technique (RFC Process)

### 7.1 Vue ASCII — RFC Process Complet

```
INITIATION
  [Demandeur soumet RFC via portail CCB]
  Formulaire: Titre, Description, Impact, Catégorie estimée
  Délai dépôt: > 10j avant implémentation souhaitée (changement Normal)
       │
       ▼
TRIAGE (CCB — 72h)
  [Président CCB + Architecte Principal]
  Classification: Standard (S) / Normal (N) / Urgent (U) / Majeur (M) / Stratégique (ST)
  S → Pré-approuvé, peut être implémenté
  N → Agenda prochaine session CCB
  U → Session CCB d'urgence (48h)
  M → Workflow Architecture (WRK-001)
  ST → Workflow Architecture (WRK-001)
       │ (N ou U)
       ▼
ANALYSE TECHNIQUE (CCB)
  [Tous membres CCB]
  - Analyse impact technique
  - Analyse impact sécurité (SOC)
  - Analyse impact données (NDPA si applicable)
  - Test en environnement staging
  - Plan de rollback documenté
       │
       ▼
VOTE CCB
  Quorum: 6/10 | Majorité simple (N) | 7/10 (M)
  Options: APPROUVÉ / APPROUVÉ SOUS CONDITIONS / REJETÉ / DIFFÉRÉ
       │ (APPROUVÉ)
       ▼
PLANIFICATION
  [Responsable Changement désigné]
  - Fenêtre de maintenance définie
  - Notification parties prenantes (72h avant)
  - Plan de communication
  - Plan de rollback confirmé
       │
       ▼
IMPLÉMENTATION
  [Équipe technique + SOC en surveillance]
  - Exécution selon plan
  - Go/No-Go au début de la fenêtre
  - SOC en mode alerte pendant implémentation
  - Test smoke post-déploiement
       │
       ▼
VALIDATION
  [Responsable Changement + SOC]
  - Tests de validation fonctionnelle
  - Confirmation pas d'impact négatif
  - Métriques de performance
  - Décision: SUCCÈS / ROLLBACK
       │ (SUCCÈS)
       ▼
CLÔTURE
  [Président CCB]
  - Ticket fermé avec statut SUCCÈS
  - Documentation mise à jour (CMDB)
  - Notification parties prenantes
  - PIR planifié (si changement N, U, M, ST)
```

### 7.2 Fenêtres de Maintenance Standard

| Système | Fenêtre Standard | Fenêtre Urgence |
|---|---|---|
| **Production Core Identity** | Dimanche 02h00-06h00 | Selon incident, préavis 2h |
| **X-Road** | Samedi 22h00-02h00 | Selon incident |
| **PKI / OCSP** | Dimanche 01h00-04h00 | Cérémonie formelle requise |
| **DataCenter** | Mensuel — Dimanche 00h00-06h00 | Selon plan BCP |
| **Portails citoyens** | Quotidien 02h00-04h00 | N/A |

---

## WORKFLOW 8 — Traitement de Plainte Citoyenne

### 8.1 Diagramme Mermaid — Plainte Citoyenne

```mermaid
flowchart TD
    PLAINTE["📝 Soumission Plainte\nCitoyen\n(En ligne / Guichet / Courrier)"]
    
    ACCUSE["📬 Accusé de Réception\nAutomatique\nNo. de suivi SNISID-PLT-YYYY-NNN\n(24h max)"]
    
    TRIAGE_PLT["🔍 Triage\nService Clients AND L1\n(3 jours ouvrables)"]
    
    Q_TYPE{"Type de\nplainte?"}
    
    TECHNIQUE["🔧 Plainte Technique\n(accès service, bug)\nL1 Support AND → 10j"]
    
    IDENTITE["🆔 Plainte Identitaire\n(erreur donnée, refus inscription)\nIGB → 30j"]
    
    DONNEES["🔒 Plainte Données\n(accès, rectification, suppression)\nNDPA → 30j"]
    
    SECURITE["🛡️ Plainte Sécurité\n(fraude, usurpation)\nSOC + PNH → 24h"]
    
    ETHIQUE["⚖️ Plainte Éthique\n(discrimination, biais IA)\nComité Éthique → 45j"]
    
    TRAITEMENT["⚙️ Traitement par\nOrgane Compétent"]
    
    Q_RESOL{"Plainte\nrésolue?"}
    
    NOTIF_CIT["📢 Notification Citoyen\nRésultat + Explications\nVoies de recours"]
    
    RECOURS["⚖️ Voies de Recours\nHierarchie AND\nou\nRecours Judiciaire"]
    
    REPORT_AND["📊 Rapport Mensuel\nStatistiques plaintes\nTendances → AND"]

    PLAINTE --> ACCUSE
    ACCUSE --> TRIAGE_PLT
    TRIAGE_PLT --> Q_TYPE
    Q_TYPE -->|"Technique"| TECHNIQUE
    Q_TYPE -->|"Identitaire"| IDENTITE
    Q_TYPE -->|"Données"| DONNEES
    Q_TYPE -->|"Sécurité"| SECURITE
    Q_TYPE -->|"Éthique"| ETHIQUE
    TECHNIQUE --> TRAITEMENT
    IDENTITE --> TRAITEMENT
    DONNEES --> TRAITEMENT
    SECURITE --> TRAITEMENT
    ETHIQUE --> TRAITEMENT
    TRAITEMENT --> Q_RESOL
    Q_RESOL -->|"OUI"| NOTIF_CIT
    Q_RESOL -->|"NON"| RECOURS
    RECOURS --> NOTIF_CIT
    NOTIF_CIT --> REPORT_AND
```

### 8.2 SLA Traitement des Plaintes

| Type | Accusé Réception | Résolution Standard | Escalade si Non Résolu |
|---|---|---|---|
| **Technique** | 24h | 10 jours | AND L3 |
| **Identitaire** | 24h | 30 jours | IGB + AND |
| **Données (DSAR)** | 24h | 30 jours (prorogeable 30j) | NDPA |
| **Sécurité / Fraude** | 2h | 24-72h | SOC + PNH |
| **Éthique** | 24h | 45 jours | Comité Éthique + AND |

---

## WORKFLOW 9 — Accréditation de Fournisseur

### 9.1 Vue ASCII

```
CANDIDATURE FOURNISSEUR
  [Soumission dossier via Portail AND]
  Documents requis:
  - Statuts légaux + K-bis ou équivalent
  - Certifications ISO (27001, 9001, etc.)
  - Références clients gouvernementaux
  - Capacités financières (bilan 3 ans)
  - Plan de continuité
  - Politique sécurité + données
  - Déclarations de conformité légale
       │
       ▼
REVUE ADMINISTRATIVE (AND — 15 jours)
  [Directeur Juridique AND]
  - Vérification complétude dossier
  - Vérification légalité entité
  - Vérification absence de sanctions
       │
       ▼
ÉVALUATION TECHNIQUE (IGB + SOC + CCB — 21 jours)
  - Capacités techniques
  - Sécurité des systèmes
  - Interopérabilité
  - Audit sur site si requis
       │
       ▼
REVUE NDPA (30 jours — si données personnelles)
  - Politique protection données
  - Mesures techniques
  - Localisation des données
  - Sous-traitants
       │
       ▼
DÉCISION (AND DG — 5 jours)
  ACCRÉDITÉ / ACCRÉDITÉ SOUS CONDITIONS / REFUSÉ
       │
       ▼
CONTRACTUALISATION
  [MJP avis si >50M HTG]
  - Contrat cadre AND
  - SLA annexé
  - Clauses de sécurité (obligatoires)
  - Clauses données (DPA annexé)
  - Clause de résiliation
       │
       ▼
SUIVI ANNUEL
  [AND + IGB + SOC]
  - Audit annuel
  - Revue des incidents
  - Renouvellement accréditation (tous les 3 ans)
```

---

## WORKFLOW 10 — Révocation d'Identité pour Fraude

### 10.1 Diagramme Mermaid — Révocation pour Fraude

```mermaid
sequenceDiagram
    autonumber
    participant DET as Détecteur<br/>(SOC/PNH/Citoyen)
    participant SOC as SOC/CERT-HT
    participant IGB as Identity<br/>Governance Board
    participant AND as AND<br/>(DG)
    participant NDPA as NDPA
    participant OJRNH as OJRNH<br/>(État Civil)
    participant CIT as Citoyen<br/>Concerné
    participant MJP as Min. Justice
    
    DET->>SOC: Signalement fraude identitaire
    SOC->>IGB: Alerte — investigation identité
    IGB->>IGB: Investigation préliminaire (72h)
    
    alt Fraude Confirmée
        IGB->>AND: Demande suspension préventive
        AND->>AND: Suspension temporaire (DG AND)
        AND->>CIT: Notification suspension + droit réponse (48h)
        AND->>NDPA: Notification suspension
        
        CIT->>IGB: Réponse citoyen (si conteste)
        IGB->>IGB: Instruction complète (21 jours)
        
        alt Fraude prouvée
            IGB->>AND: Recommandation révocation définitive
            AND->>NDPA: Consultation obligatoire
            NDPA->>AND: Avis NDPA (15 jours)
            AND->>AND: Décision révocation (DG AND signé)
            AND->>OJRNH: Instruction révocation NIN
            OJRNH->>OJRNH: Révocation dans registre
            AND->>MJP: Transmission dossier judiciaire
            AND->>CIT: Notification révocation + voies de recours
        else Fraude infirmée
            AND->>AND: Levée suspension immédiate
            AND->>CIT: Notification + excuses
            IGB->>IGB: Rapport d'incident
        end
    else Fraude Non Confirmée
        SOC->>IGB: Classement sans suite
        IGB->>IGB: Documentation
    end
```

---

## Indicateurs de Performance des Workflows (KPIs)

### Tableau de Bord Gouvernance

| KPI | Cible | Fréquence Mesure | Responsable |
|---|---|---|---|
| **Taux de respect SLA changements N** | ≥ 95% | Mensuel | CCB |
| **Délai moyen décision AND** | ≤ 5 jours | Mensuel | AND DG |
| **Délai moyen traitement plainte citoyen** | ≤ 20 jours | Mensuel | AND COO |
| **Taux d'incidents P1+ sur 12 mois** | ≤ 2 | Annuel | SOC |
| **Taux de conformité légale (attestations)** | 100% | Trimestriel | Dir. Juridique |
| **Délai moyen accréditation fournisseur** | ≤ 60 jours | Trimestriel | AND DG |
| **Taux de sessions CNN avec quorum** | 100% | Trimestriel | CNN Président |
| **Délai notification CNN (incident P1)** | ≤ 2h | Par incident | SOC |
| **Taux approbation budget dans délais** | ≥ 90% | Trimestriel | AND CFO |
| **NPS citoyen (satisfaction service)** | ≥ 70 | Semestriel | AND COO |

---

## Bloc de Signature

```
APPROUVÉ PAR LE DIRECTEUR GÉNÉRAL AND

Nom            : ___________________________
Qualité        : Directeur Général, Autorité Nationale Numérique
Signature      : ___________________________  Date : __________
Cachet AND     : [CACHET AND]

VALIDÉ PAR LE DIRECTEUR ARCHITECTURE & INNOVATION

Nom            : ___________________________
Qualité        : Directeur Architecture & Innovation, AND
Signature      : ___________________________  Date : __________

VALIDÉ PAR LE PRÉSIDENT CCB

Nom            : ___________________________
Qualité        : Président, Change Control Board SNISID
Signature      : ___________________________  Date : __________
```

**HISTORIQUE DES RÉVISIONS**

| Version | Date | Modifications | Approuvé par |
|---|---|---|---|
| 0.1 | 2026-03-01 | Workflows 1, 2, 5 initiaux | AND Architecture |
| 0.5 | 2026-04-15 | Ajout workflows 3, 4, 6, 7 | DGA AND |
| 0.8 | 2026-05-10 | Ajout workflows 8, 9, 10 + KPIs | DG AND |
| 1.0 | 2026-05-25 | Version finale, revue CCB + IGB | DG AND |

---

*Document SNISID-GOV-WRK-001 v1.0.0 — Propriété de l'Autorité Nationale Numérique de la République d'Haïti*

*© 2026 République d'Haïti — SNISID Phase 0 — Gouvernance Nationale*
