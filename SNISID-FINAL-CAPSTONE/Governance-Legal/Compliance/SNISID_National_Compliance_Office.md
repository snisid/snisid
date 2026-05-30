# BUREAU NATIONAL DE CONFORMITÉ (BNC)

**Classification :** DOCUMENT INSTITUTIONNEL — CONFORMITÉ NATIONALE
**Référence :** SNISID-COMP-001
**Version :** 1.0
**Date :** 25 mai 2026

---

## 1. MISSION

Le Bureau National de Conformité (BNC) est l'organe permanent chargé d'industrialiser, coordonner et surveiller la conformité de l'ensemble des composantes du SNISID aux normes légales, sécuritaires et opérationnelles applicables.

---

## 2. RATTACHEMENT

Le BNC est rattaché au Conseil National de Gouvernance Numérique et rend compte directement au Conseil. Il opère de manière transversale sur l'ensemble des agences SNISID.

---

## 3. DOMAINES DE CONFORMITÉ

### 3.1 Legal Compliance (Conformité Légale)

| Domaine | Norme de Référence | Contrôle |
|---------|-------------------|----------|
| Identité numérique | Loi sur l'identité numérique | Continu |
| Protection des données | Loi sur la protection des données | Continu |
| Cybersécurité | Loi sur la cybersécurité | Continu |
| Signatures électroniques | Loi sur les signatures électroniques | Continu |
| Preuve numérique | Loi sur la preuve numérique | Continu |
| Droit constitutionnel | Constitution haïtienne | Continu |
| Droit administratif | Code administratif | Continu |

**Processus :**
1. Cartographie des obligations légales applicables
2. Matrice de conformité par agence et par système
3. Évaluation continue de conformité
4. Rapport de gaps et plan de remédiation
5. Suivi des remédiations
6. Certification de conformité

### 3.2 Security Compliance (Conformité Sécuritaire)

| Standard | Application | Fréquence de Vérification |
|----------|-------------|--------------------------|
| Politique de sécurité SNISID | Tous systèmes | Continue |
| Standards cryptographiques | PKI, chiffrement | Trimestrielle |
| Contrôle d'accès | IAM, RBAC | Mensuelle |
| Gestion des vulnérabilités | Infrastructure | Hebdomadaire |
| Tests de pénétration | Systèmes critiques | Semestrielle |
| Journalisation | Tous systèmes | Continue |
| Gestion des incidents | SOC, CERT-HT | Continue |
| Continuité d'activité | PCA/PRA | Semestrielle |

### 3.3 Operational Compliance (Conformité Opérationnelle)

| Domaine | Standard | Vérification |
|---------|---------|-------------|
| Gestion des changements | Procédure SNISID-CHG | Par changement |
| Gestion des incidents | Procédure SNISID-INC | Par incident |
| Gestion des problèmes | Procédure SNISID-PRB | Continue |
| Gestion des capacités | Procédure SNISID-CAP | Mensuelle |
| Gestion des niveaux de service | SLA/SLO définis | Continue |
| Gestion des configurations | CMDB | Continue |
| Gestion des mises en production | Procédure SNISID-REL | Par release |

### 3.4 Audit Coordination (Coordination des Audits)

| Type d'Audit | Responsable | Fréquence | Coordination BNC |
|-------------|------------|-----------|-----------------|
| Audit interne sécurité | Équipe audit interne | Trimestriel | Planification + suivi |
| Audit externe sécurité | Auditeur qualifié | Annuel | Sélection + coordination |
| Audit conformité légale | ANPD / Juridique | Semestriel | Préparation + suivi |
| Audit ANPD (protection données) | ANPD | Annuel | Liaison + préparation |
| Audit financier | Cour des Comptes | Annuel | Coordination |
| Audit technique | ANC | Annuel | Préparation technique |
| Audit indépendant souveraineté | Conseil National | Annuel | Organisation |

**Processus de coordination :**
```
Planification annuelle des audits
    → Notification aux agences (60 jours avant)
    → Préparation documentaire
    → Facilitation de l'audit
    → Réception du rapport
    → Plan d'action correctif
    → Suivi des remédiations
    → Clôture et capitalisation
```

### 3.5 Regulatory Reporting (Reporting Réglementaire)

| Rapport | Destinataire | Fréquence | Contenu |
|---------|-------------|-----------|---------|
| Tableau de bord conformité | Conseil de Gouvernance | Mensuel | KPIs, gaps, risques |
| Rapport de conformité légale | Conseil National | Trimestriel | Statut par loi |
| Rapport de conformité sécurité | ANC | Trimestriel | Statut sécuritaire |
| Rapport ANPD | ANPD | Semestriel | Protection données |
| Rapport annuel de conformité | Parlement | Annuel | Bilan complet |
| Rapport d'incidents de conformité | Conseil | Sur événement | Détail incident |

---

## 4. STRUCTURE ORGANISATIONNELLE

```
Bureau National de Conformité (BNC)
├── Directeur de la Conformité
│   ├── Adjoint Conformité Légale
│   └── Adjoint Conformité Technique
├── Division Conformité Légale
│   ├── Cellule Identité & État Civil
│   ├── Cellule Protection des Données
│   ├── Cellule Cybersécurité
│   └── Cellule Droit Administratif
├── Division Conformité Sécuritaire
│   ├── Cellule Standards de Sécurité
│   ├── Cellule Cryptographie
│   └── Cellule Tests & Audits
├── Division Conformité Opérationnelle
│   ├── Cellule Processus
│   ├── Cellule Niveaux de Service
│   └── Cellule Configuration
├── Division Audit & Coordination
│   ├── Cellule Audit Interne
│   ├── Cellule Relations Auditeurs Externes
│   └── Cellule Suivi Remédiations
└── Division Reporting & Analyse
    ├── Cellule Tableaux de Bord
    ├── Cellule Analyse de Risques
    └── Cellule Documentation
```

---

## 5. PROCESSUS CLÉS

### 5.1 Cycle de Conformité Continue

```
┌──────────────────────────────────────────┐
│           CYCLE DE CONFORMITÉ            │
│                                          │
│  1. IDENTIFIER → Obligations applicables │
│  2. ÉVALUER   → Niveau de conformité     │
│  3. PLANIFIER → Actions de remédiation   │
│  4. EXÉCUTER  → Mise en conformité       │
│  5. VÉRIFIER  → Contrôle de conformité   │
│  6. REPORTER  → Rapport aux autorités    │
│  7. AMÉLIORER → Leçons apprises          │
│                                          │
│  ← Cycle continu, jamais terminé →       │
└──────────────────────────────────────────┘
```

### 5.2 Gestion des Non-Conformités

| Sévérité | Description | Délai de Remédiation | Escalade |
|----------|-------------|---------------------|----------|
| Critique | Violation légale active | 24 heures | Conseil National immédiat |
| Majeure | Risque de violation | 7 jours | Directeur BNC |
| Significative | Écart par rapport au standard | 30 jours | Chef de division |
| Mineure | Amélioration recommandée | 90 jours | Chef de cellule |

### 5.3 Processus d'Exception

Lorsqu'une obligation ne peut être respectée temporairement :
1. Demande d'exception documentée
2. Justification et analyse de risque
3. Mesures compensatoires proposées
4. Approbation par le Directeur BNC (ou Conseil si critique)
5. Durée limitée (maximum 90 jours)
6. Suivi renforcé
7. Revue à expiration

---

## 6. OUTILS DE CONFORMITÉ

| Outil | Fonction | Utilisateur |
|-------|---------|-------------|
| Registre de conformité | Base de données des obligations | BNC |
| Matrice de conformité | Mapping obligations ↔ contrôles | Agences |
| Tableau de bord temps réel | Monitoring conformité | Direction |
| Système de tickets non-conformité | Suivi des remédiations | Tous |
| Bibliothèque de politiques | Référentiel documentaire | Tous |
| Outil d'auto-évaluation | Check-lists par agence | Agences |

---

## 7. INDICATEURS DE PERFORMANCE

| KPI | Cible | Mesure |
|-----|-------|--------|
| Taux de conformité légale global | ≥ 95% | Mensuel |
| Taux de conformité sécuritaire | ≥ 98% | Mensuel |
| Délai moyen de remédiation | ≤ 15 jours | Mensuel |
| Taux de clôture des audits | 100% dans les délais | Trimestriel |
| Nombre de non-conformités critiques | 0 | Continu |
| Taux de reporting dans les délais | 100% | Mensuel |
| Couverture d'audit | 100% des systèmes/an | Annuel |

---

*Document institutionnel préparé dans le cadre de la Phase 14 — SNISID National Legal Framework*
