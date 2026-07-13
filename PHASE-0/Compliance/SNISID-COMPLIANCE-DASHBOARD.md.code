# 📊 SNISID — COMPLIANCE DASHBOARD
## Tableau de Bord de Conformité Continue — Phase 0

**Document ID :** SNISID-CMP-001  
**Version :** 1.0.0  
**Date :** Mai 2026  
**Propriétaire :** Autorité d'Audit et de Conformité (AAC) + NDPA  
**Classification :** Usage Gouvernemental — Restreint  

---

## TABLEAU DE BORD EXÉCUTIF DE CONFORMITÉ

```
╔══════════════════════════════════════════════════════════════════════╗
║         SNISID — COMPLIANCE STATUS DASHBOARD — Phase 0              ║
║                    Mai 2026 — Baseline                               ║
╠══════════════════════════════════════════════════════════════════════╣
║  ISO 27001  │  🟡 EN COURS  │ Target: Q3 2030 (après build)         ║
║  ISO 22301  │  🟡 PLANIFIÉ  │ Target: Q4 2028 (après DR validé)     ║
║  NIST CSF   │  🟡 PARTIEL   │ 40% implémenté (architecture Phase 0) ║
║  Convention 108+ │ 🔴 REQUIS │ Ratification traité requise           ║
║  Budapest    │  🔴 REQUIS   │ Ratification traité requise            ║
║  ODD 16.9   │  🔴 EN COURS  │ Cible 2030 (≥95% enrôlement)          ║
║  ICAO 9303  │  🟡 PLANIFIÉ  │ Phase 2 (passeport biométrique)       ║
║  eIDAS-like │  🟡 PLANIFIÉ  │ Loi signature électronique requise    ║
╚══════════════════════════════════════════════════════════════════════╝
```

---

## 1. MAPPING ISO/IEC 27001:2022

### Domaine A.5 — Politiques de Sécurité de l'Information

| Contrôle | Exigence | Statut Phase 0 | Action requise | Owner | Échéance |
|---------|----------|----------------|----------------|-------|----------|
| A.5.1 | Politique SI définie et approuvée | 🟡 Défini architecturalement | Valider formellement par CNN | CISO | Q3 2026 |
| A.5.2 | Révision annuelle politique | 🟡 Processus défini | Mettre en place dès AND créée | AND | Q4 2026 |

### Domaine A.6 — Organisation de la Sécurité

| Contrôle | Exigence | Statut Phase 0 | Action requise | Owner | Échéance |
|---------|----------|----------------|----------------|-------|----------|
| A.6.1 | RSSI/CISO défini | 🟡 Rôle défini | Recrutement Phase 1 | AND | Q1 2027 |
| A.6.7 | Télétravail sécurisé | 🟡 Architecturé | Implémentation Phase 2 | CISO | 2028 |

### Domaine A.8 — Gestion des Actifs

| Contrôle | Exigence | Statut Phase 0 | Action requise | Owner | Échéance |
|---------|----------|----------------|----------------|-------|----------|
| A.8.1 | Inventaire actifs | 🔴 À créer | CMDB Phase 1 | Infra | Q2 2027 |
| A.8.2 | Classification données | ✅ Définie | Implémenter Apache Atlas | Data Steward | Phase 2 |
| A.8.3 | Retrait des médias | 🟡 Politique définie | Procédures à formaliser | CISO | Q4 2026 |

### Domaine A.9 — Contrôle d'Accès

| Contrôle | Exigence | Statut Phase 0 | Action requise | Owner | Échéance |
|---------|----------|----------------|----------------|-------|----------|
| A.9.1 | Politique contrôle accès | ✅ Définie (RBAC + ABAC) | Implémenter Keycloak + OPA | IAM | Phase 2 |
| A.9.2 | Gestion accès utilisateurs | 🟡 Architecturé | SCIM 2.0 provisioning | IAM | Phase 2 |
| A.9.3 | Responsabilités utilisateurs | ✅ Défini | Manuel utilisateur | BGW | Q1 2027 |
| A.9.4 | Contrôle accès systèmes | ✅ Zero Trust défini | Implémentation Phase 2 | CISO | Phase 2 |

### Domaine A.10 — Cryptographie

| Contrôle | Exigence | Statut Phase 0 | Action requise | Owner | Échéance |
|---------|----------|----------------|----------------|-------|----------|
| A.10.1 | Politique cryptographique | ✅ Définie (Standards doc) | Formaliser décret | AND | Q4 2026 |

### Domaine A.11 — Sécurité Physique

| Contrôle | Exigence | Statut Phase 0 | Action requise | Owner | Échéance |
|---------|----------|----------------|----------------|-------|----------|
| A.11.1 | Périmètre physique DC | 🟡 Spécifié (Tier III) | Construction DC Phase 1-2 | Infra | 2027 |
| A.11.2 | Sécurité équipements | 🟡 Spécifiée | Mise en œuvre Phase 2 | Infra | 2027 |

### Domaine A.12 — Sécurité Exploitation

| Contrôle | Exigence | Statut Phase 0 | Action requise | Owner | Échéance |
|---------|----------|----------------|----------------|-------|----------|
| A.12.1 | Procédures exploitation | 🟡 Runbooks définis | Documentation complète Phase 2 | NOC | Phase 2 |
| A.12.3 | Sauvegardes | ✅ 3-2-1-1-0 défini | Mise en œuvre Phase 2 | Infra | Phase 2 |
| A.12.6 | Gestion vulnérabilités | 🟡 Processus défini | Outils Phase 2 | CISO | Phase 2 |

### Domaine A.13 — Sécurité Communications

| Contrôle | Exigence | Statut Phase 0 | Action requise | Owner | Échéance |
|---------|----------|----------------|----------------|-------|----------|
| A.13.1 | Sécurité réseaux | ✅ Zero Trust + mTLS | Implémentation Phase 2 | Infra | Phase 2 |
| A.13.2 | Transferts information | ✅ X-Road + CloudEvents | Gouvernance APIs | BGI | Phase 2 |

---

## 2. MAPPING NIST CSF 2.0

### GOVERN (GV)

| Catégorie | Sous-catégorie | Statut | Action |
|-----------|---------------|--------|--------|
| GV.OC | Contexte organisationnel | ✅ Vision nationale définie | Valider CNN |
| GV.RM | Stratégie gestion risques | ✅ Risk register 30 risques | Comité Risques mensuel |
| GV.RR | Rôles & responsabilités | ✅ RACI matrix définie | Recrutement Phase 1 |
| GV.PO | Politiques | 🟡 Architecturales | Formaliser lois et décrets |
| GV.OV | Oversight | 🟡 CNN + AND définis | Activation légale |
| GV.SC | Supply chain | 🟡 Politiques définies | Processus achat Phase 1 |

### IDENTIFY (ID)

| Catégorie | Statut | Action |
|-----------|--------|--------|
| ID.AM — Assets | 🔴 CMDB à créer | Phase 1 |
| ID.RA — Risk Assessment | ✅ EBIOS RM appliqué | Révisions trimestrielles |
| ID.IM — Improvement | 🟡 Processus définis | Post-mortem formels Phase 2 |

### PROTECT (PR)

| Catégorie | Statut | Action |
|-----------|--------|--------|
| PR.AA — Authentication | ✅ FIDO2 + OIDC définis | Implémentation Phase 2 |
| PR.AT — Awareness | 🔴 Programme à créer | Phase 1 (formation) |
| PR.DS — Data Security | ✅ Chiffrement + WORM définis | Phase 2 |
| PR.PS — Platform Security | 🟡 Talos/CIS définis | Phase 2 |
| PR.IR — Incident Response | 🟡 CSIRT-HT défini | Activation Phase 1 |

### DETECT (DE)

| Catégorie | Statut | Action |
|-----------|--------|--------|
| DE.AE — Adverse Events | 🟡 SOC + SIEM définis | Déploiement Phase 2 |
| DE.CM — Continuous Monitoring | 🟡 OTel + Prometheus | Phase 2 |

### RESPOND (RS)

| Catégorie | Statut | Action |
|-----------|--------|--------|
| RS.MA — Incident Management | ✅ Procédures PRDR | Activation SOC Phase 2 |
| RS.AN — Analysis | 🟡 SOAR playbooks définis | Outils Phase 2 |
| RS.CO — Communication | 🟡 Plan crise défini | Exercices Phase 2 |

### RECOVER (RC)

| Catégorie | Statut | Action |
|-----------|--------|--------|
| RC.RP — Recovery Plan | ✅ DRP défini (PaP→CapH) | Tests Phase 2 |
| RC.CO — Communications | 🟡 Status page planifiée | Phase 2 |

---

## 3. MAPPING ODD 16.9

**ODD 16.9 :** *D'ici 2030, garantir à tous une identité juridique, notamment grâce à l'enregistrement des naissances.*

| Indicateur ODD | Baseline 2026 | Cible 2027 | Cible 2028 | Cible 2029 | Cible 2030 |
|---------------|--------------|-----------|-----------|-----------|-----------|
| % population avec identité légale | ~65 % (CIN existantes) | 70 % | 80 % | 90 % | ≥ 95 % |
| % naissances enregistrées | ~60-70 % | 75 % | 85 % | 90 % | ≥ 95 % |
| Délai enregistrement naissance | 30+ jours | 7 jours | 48h | 24h | <24h |
| Couverture géographique | Urbaine principalement | +2 dépts | +5 dépts | 10 dépts | National |
| Inclusion groupes vulnérables | Faible | Programme lancé | En cours | 80% | ≥ 95% |

---

## 4. MAPPING CONVENTION 108+ DU CONSEIL DE L'EUROPE

| Article | Exigence | Statut | Action |
|---------|----------|--------|--------|
| Art. 1 | Objet et finalité de la protection | 🟡 Dans Loi NDPA | Vote loi Q3 2026 |
| Art. 4 | Qualité des données | ✅ Data Model défini | Mise en œuvre Phase 2 |
| Art. 5 | Légitimité du traitement | 🟡 Dans Loi NDPA | Vote loi Q3 2026 |
| Art. 6 | Données sensibles (biométriques) | 🟡 Dans Loi Biométrie | Vote loi Q4 2026 |
| Art. 7 | Sécurité des données | ✅ Framework défini | Phase 2 |
| Art. 8 | Droits des personnes | ✅ 8 droits définis | Portail Confiance Phase 2 |
| Art. 9 | Exceptions | 🟡 Défini pour sécurité nationale | Loi NDPA |
| Art. 15 | Autorité de contrôle | 🟡 NDPA définie | Loi NDPA |

**Action requise :** Ratification officielle Convention 108+ par Haïti (démarche diplomatique — Ministère des Affaires Étrangères)

---

## 5. MAPPING CONVENTION DE BUDAPEST (Cybercriminalité)

| Titre | Exigence | Statut | Action |
|-------|----------|--------|--------|
| Titre 1 | Infractions informatiques | 🟡 Dans Loi Cyber | Vote Q1 2027 |
| Titre 2 | Infractions liées au contenu | 🟡 Dans Loi Cyber | Vote Q1 2027 |
| Titre 3 | Coopération internationale | 🟡 Protocoles définis | Ratification + bilatéraux |
| Titre 4 | Procédures (conservation, production) | 🟡 Procédures SOC | Formaliser CSIRT-HT |

**Action requise :** Ratification officielle Convention de Budapest (démarche diplomatique + MJSP)

---

## 6. CALENDRIER DE MISE EN CONFORMITÉ

```
2026
Q2  Dépôt Loi-cadre SNISID
Q3  Dépôt Loi Protection Données (→ NDPA)
Q4  Dépôt Loi Biométrie + Signature Électronique
Q4  AND créée par décret

2027
Q1  Dépôt Loi Cybersécurité + Interopérabilité
Q2  Vote Loi-cadre SNISID ← JALON CRITIQUE
Q3  Ratification Convention 108+ (démarche lancée)
Q4  Ratification Convention Budapest
Q4  Décrets d'application publiés

2028
Q1  Premières certifications (pilote département Ouest)
Q2  Audit ISO 27001 préliminaire
Q4  ISO 27001 audit complet (après déploiement Phase 2)

2029
Q2  NDPA opérationnelle et indépendante
Q4  Premier rapport NDPA public

2030
Q3  Certification ISO 27001 OBTENUE
Q4  Rapport ODD 16.9 au Parlement (cible ≥95%)
```

---

## 7. AUDITS PLANIFIÉS

| Audit | Fréquence | Auditeur | Prochain |
|-------|-----------|----------|---------|
| ISO 27001 surveillance | Annuel | Cabinet certifié externe | Q3 2028 (après Phase 2) |
| ISO 27001 certification | Initial + renouvellement 3 ans | Bureau de certification | Q3 2030 |
| Pentest infrastructure | Annuel | Cabinet externe agréé | Q4 2028 |
| Pentest applicatif | Bi-annuel | Externe + Red Team interne | Phase 2 |
| Audit PKI | Annuel | Cabinet spécialisé PKI | Q4 2027 (après key ceremony) |
| Audit NDPA | Annuel | NDPA (indépendante) | Q2 2029 (dès NDPA opé.) |
| Audit financier SNISID | Annuel | Cour des Comptes + externe | Q1 2028 |
| Code review sécurité | Continu | SAST + interne | En continu dès Phase 2 |
| Tabletop exercises | Trimestriel | Interne + externe facilitateur | Q1 2028 |
| Disaster Recovery Drill | Semestriel | NOC + SOC + AND | H1 et H2 annuellement |

---

## 8. INDICATEURS DE CONFORMITÉ

| KPI Conformité | Cible | Fréquence mesure |
|----------------|-------|-----------------|
| Lois SNISID votées | 7/7 d'ici Q2 2027 | Trimestriel |
| Contrôles ISO 27001 implémentés | ≥ 80% d'ici 2029 | Semestriel |
| Incidents non-conformité | 0 critique / an | Mensuel |
| Plaintes NDPA non-résolues | 0 au-delà de 90 jours | Mensuel |
| Personnel formé sécurité | 100% | Annuel |
| Audits sans réserve majeure | 100% | Annuel |
| DRP testé et validé | 2x/an | Semestriel |
| Vulnérabilités critiques ouvertes | 0 >7 jours | Hebdomadaire |

---

*Document approuvé par l'Autorité d'Audit et de Conformité (AAC)*  
*SNISID — République d'Haïti*  
*Prochaine révision : Trimestre suivant la création de l'AND*
