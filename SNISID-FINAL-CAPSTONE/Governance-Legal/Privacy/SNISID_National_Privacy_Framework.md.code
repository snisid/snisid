# SNISID — CADRE NATIONAL DE PROTECTION DE LA VIE PRIVÉE

**Classification :** CADRE STRATÉGIQUE — VIE PRIVÉE
**Référence :** SNISID-PRIV-002
**Version :** 1.0
**Date :** 25 mai 2026

---

## 1. OBJECTIF

Intégrer la protection de la vie privée dans l'architecture même du SNISID (Privacy by Design), garantissant que les droits fondamentaux des citoyens sont respectés à chaque étape du cycle de vie des données.

---

## 2. PRINCIPES DIRECTEURS

| Principe | Description | Application SNISID |
|----------|-------------|-------------------|
| Privacy by Design | La vie privée intégrée dès la conception | Architecture, code, processus |
| Privacy by Default | Protection maximale par défaut | Paramètres par défaut restrictifs |
| Minimisation | Ne collecter que le strict nécessaire | Chaque champ justifié |
| Transparence | Informer clairement les citoyens | Notices claires, registre public |
| Contrôle citoyen | Les citoyens maîtrisent leurs données | Portail citoyen, droits effectifs |
| Responsabilité | Prouver la conformité | Documentation, audits, DPD |

---

## 3. FONCTIONS DE PROTECTION

### 3.1 Consent Management (Gestion du Consentement)

**Principes du consentement SNISID :**

| Exigence | Description |
|----------|-------------|
| Libre | Pas de conséquence négative en cas de refus (sauf obligation légale) |
| Spécifique | Un consentement par finalité |
| Éclairé | Information claire, complète, compréhensible |
| Univoque | Action positive et non ambiguë |
| Révocable | Retrait aussi facile que l'octroi |
| Documenté | Preuve du consentement conservée |

**Catégories de consentement :**

| Catégorie | Base légale | Consentement requis |
|-----------|-----------|-------------------|
| Enregistrement d'identité | Obligation légale | Non (base légale) |
| Collecte biométrique | Obligation légale | Non (base légale) — information obligatoire |
| Partage inter-agences (obligatoire) | Mission d'intérêt public | Non (base légale) — information obligatoire |
| Partage inter-agences (optionnel) | Consentement | Oui, spécifique |
| Services numériques optionnels | Consentement | Oui |
| Communication d'informations | Consentement | Oui |
| Recherche / statistiques | Consentement ou anonymisation | Oui (si non anonymisé) |

**Portail de gestion du consentement citoyen :**
```
Espace Citoyen > Mes Consentements
├── Consentements actifs
│   ├── [Service] — Accordé le [date] — [Révoquer]
│   ├── [Service] — Accordé le [date] — [Révoquer]
│   └── ...
├── Historique des consentements
│   ├── [Service] — Révoqué le [date]
│   └── ...
├── Demandes en attente
│   ├── [Agence] demande accès à [données] pour [finalité]
│   │   [Accepter] [Refuser] [En savoir plus]
│   └── ...
└── Paramètres
    ├── Notifications de demandes d'accès : [Activé]
    ├── Rapport mensuel d'accès : [Activé]
    └── Langue préférée : [Français/Créole]
```

### 3.2 Purpose Limitation (Limitation de Finalité)

| Règle | Application |
|-------|------------|
| Finalité déclarée | Chaque traitement a une finalité explicite, documentée et approuvée |
| Non-détournement | Les données collectées pour une finalité ne peuvent être utilisées pour une autre sans nouvelle base légale |
| Contrôle technique | Le système vérifie automatiquement que l'usage correspond à la finalité autorisée |
| Journalisation | Chaque accès est journalisé avec la finalité invoquée |
| Audit de finalité | Vérification régulière de la cohérence usage/finalité |

**Matrice de finalités autorisées :**

| Donnée | Identification | Sécurité | Fiscal | Santé | Recherche |
|--------|---------------|----------|--------|-------|-----------|
| Nom/Prénoms | ✅ | ✅ | ✅ | ✅ | ❌ (anonymisé) |
| NNI | ✅ | ✅ | ✅ | ✅ | ❌ |
| Biométrie | ✅ | ✅ | ❌ | ❌ | ❌ |
| Adresse | ✅ | ✅ | ✅ | ❌ | ❌ |
| Photo | ✅ | ✅ | ❌ | ❌ | ❌ |
| État civil | ✅ | ❌ | ✅ | ❌ | ❌ (anonymisé) |

### 3.3 Data Minimization (Minimisation des Données)

| Principe | Mise en œuvre |
|----------|-------------|
| Collecte minimale | Chaque champ justifié par une analyse de nécessité |
| Accès minimal | Accès uniquement aux données nécessaires à la mission |
| Conservation minimale | Durée de conservation la plus courte possible |
| Granularité d'accès | Accès par attribut, pas par profil complet |
| Pseudonymisation | Par défaut pour les usages non nominatifs |
| Anonymisation | Pour les statistiques et la recherche |

**Niveaux d'accès granulaires :**

| Niveau | Données accessibles | Exemple d'usage |
|--------|-------------------|-----------------|
| Vérification seule | Oui/Non (la personne existe-t-elle ?) | Vérification en ligne |
| Attributs basiques | Nom, photo | Identification visuelle |
| Profil limité | Nom, date de naissance, photo | Services publics courants |
| Profil complet | Toutes données non biométriques | Administration complète |
| Profil étendu | Profil complet + biométrie | Enquête autorisée |

### 3.4 Access Rights (Droits d'Accès)

**Portail Citoyen — Mes Données :**

| Droit | Mécanisme | Délai |
|-------|-----------|-------|
| Consultation | Portail en ligne + bureaux | Immédiat (en ligne) |
| Copie | Export PDF/JSON signé | 24 heures |
| Historique d'accès | Qui a consulté mes données, quand, pourquoi | Immédiat |
| Demande de rectification | Formulaire en ligne | 15 jours |
| Opposition | Formulaire motivé | 15 jours |
| Plainte | Portail ANPD intégré | 5 jours (accusé réception) |

**Journal d'accès citoyen (extrait) :**
```
╔══════════════════════════════════════════════════════╗
║  MON HISTORIQUE D'ACCÈS — NNI: ****-****-**12       ║
╠══════════════════════════════════════════════════════╣
║ 25/05/2026 09:15 — ONI Bureau PAP                   ║
║   Agent: ****42 — Motif: Renouvellement CNI         ║
║   Données: Profil complet — Durée: 8 min            ║
║                                                      ║
║ 20/05/2026 14:30 — PNH Commissariat Nord             ║
║   Agent: ****87 — Motif: Vérification identité       ║
║   Données: Nom + Photo — Durée: 2 min               ║
║                                                      ║
║ 15/05/2026 11:00 — DGI En ligne                      ║
║   Système: Auto — Motif: Déclaration fiscale         ║
║   Données: NNI + Nom — Durée: < 1 min               ║
╚══════════════════════════════════════════════════════╝
```

### 3.5 Right to Correction (Droit de Rectification)

**Processus de rectification :**
```
Citoyen soumet demande de rectification
    → Vérification d'identité (LOA 3 minimum)
    → Enregistrement de la demande (ticket)
    → Pièces justificatives requises ?
        → Oui → Demande de pièces
        → Non → Traitement direct
    → Vérification par agent habilité
    → Décision
        → Acceptée → Modification + notification + journal
        → Refusée → Notification motivée + voies de recours
    → Délai maximum : 15 jours ouvrés
    → Recours possible devant l'ANPD
```

**Données rectifiables vs non-rectifiables :**

| Donnée | Rectifiable | Condition |
|--------|-------------|-----------|
| Nom / Prénoms | Oui | Décision d'état civil |
| Date de naissance | Oui | Jugement supplétif |
| Adresse | Oui | Justificatif |
| Photo | Oui | Nouvelle capture |
| Sexe | Oui | Décision judiciaire |
| NNI | Non | Identifiant permanent |
| Empreintes digitales | Mise à jour | Nouvelle capture si qualité insuffisante |
| Filiation | Oui | Acte d'état civil |

---

## 4. ANALYSE D'IMPACT SUR LA VIE PRIVÉE (AIVP)

### 4.1 Quand Réaliser une AIVP ?

| Situation | AIVP obligatoire |
|-----------|-----------------|
| Nouveau traitement de données personnelles | Oui |
| Modification significative d'un traitement existant | Oui |
| Nouvelle interconnexion de bases de données | Oui |
| Nouveau système de surveillance | Oui |
| Utilisation de biométrie | Oui |
| Utilisation d'IA sur des données personnelles | Oui |
| Profilage | Oui |

### 4.2 Contenu de l'AIVP

| Section | Contenu |
|---------|---------|
| Description du traitement | Finalité, données, personnes concernées, flux |
| Base légale | Justification du traitement |
| Nécessité et proportionnalité | Pourquoi ces données, pourquoi pas moins |
| Risques identifiés | Risques pour les droits et libertés |
| Mesures de mitigation | Mesures techniques et organisationnelles |
| Avis du DPD | Recommandations du Délégué |
| Décision | Acceptation, modification ou refus |
| Consultation ANPD | Si risque résiduel élevé |

---

## 5. FORMATION ET SENSIBILISATION

| Public | Formation | Fréquence |
|--------|-----------|-----------|
| Agents SNISID | Protection des données et vie privée | Annuelle obligatoire |
| Développeurs | Privacy by Design, sécurité des données | Semestrielle |
| Managers | Responsabilités légales | Annuelle |
| Nouveaux agents | Module vie privée dans l'intégration | À l'embauche |
| Citoyens | Campagnes de sensibilisation | Continue |

---

## 6. INDICATEURS VIE PRIVÉE

| KPI | Cible | Mesure |
|-----|-------|--------|
| Taux de demandes d'accès traitées dans les délais | 100% | Mensuel |
| Taux de rectification dans les délais | 100% | Mensuel |
| Nombre de violations de données | 0 | Continu |
| Taux de réalisation des AIVP | 100% | Continu |
| Taux de formation complétée | 100% | Annuel |
| Taux de satisfaction citoyenne (vie privée) | ≥ 80% | Semestriel |

---

*Document cadre préparé dans le cadre de la Phase 14 — SNISID National Legal Framework*
