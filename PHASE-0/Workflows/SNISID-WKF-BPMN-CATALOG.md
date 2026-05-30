# SNISID-WKF-BPMN-CATALOG
## Catalogue Complet des Workflows BPMN — SNISID Phase 0

---

| Champ         | Valeur                                      |
|---------------|---------------------------------------------|
| **Document ID** | SNISID-WKF-BPMN-CATALOG-v1.0             |
| **Version**     | 1.0.0                                     |
| **Statut**      | APPROUVÉ — Production                     |
| **Date**        | 2026-05-25                                |
| **Auteur**      | Direction Technique SNISID / ONI           |
| **Révision**    | Comité Technique SNISID                   |
| **Classification** | CONFIDENTIEL — Usage Interne ONI       |

---

### BLOC D'APPROBATION

| Rôle                          | Nom                | Signature       | Date       |
|-------------------------------|--------------------|-----------------|------------|
| Directeur Général ONI         | [DG ONI]           | ____________    | __________ |
| Directeur Technique SNISID    | [DT SNISID]        | ____________    | __________ |
| Architecte Principal          | [Arch. Principal]  | ____________    | __________ |
| Officier Sécurité (CISO)      | [CISO ONI]         | ____________    | __________ |
| Représentant MJSP             | [Rep. MJSP]        | ____________    | __________ |

---

## TABLE DES MATIÈRES

1. [Introduction et Conventions BPMN](#1-introduction)
2. [Catalogue des Workflows — État Civil Naissance (EC-N)](#2-etat-civil-naissance)
   - EC-N01 Naissance Simple
   - EC-N02 Naissance par Reconnaissance
   - EC-N03 Naissance par Déclaration Tardive
   - EC-N04 Naissance par Décret
   - EC-N05 Naissance par Jugement au rang des Minutes
3. [Catalogue des Workflows — État Civil Mariage (EC-M)](#3-etat-civil-mariage)
4. [Catalogue des Workflows — État Civil Divorce (EC-D)](#4-etat-civil-divorce)
5. [Catalogue des Workflows — État Civil Décès (EC-X)](#5-etat-civil-deces)
6. [Catalogue des Workflows — Adoption (EC-A)](#6-adoption)
7. [Catalogue des Workflows — Identité Nationale (ID)](#7-identite-nationale)
8. [Catalogue des Workflows — Passeport (PP)](#8-passeport)
9. [Catalogue des Workflows — Fiscalité (FI)](#9-fiscalite)
10. [Tables DMN — Règles de Décision](#10-dmn-tables)
11. [Intégration Bus National (Kafka Events)](#11-kafka-events)
12. [Procédures Mode Hors-Ligne](#12-offline-mode)
13. [Gestion des Exceptions et Escalades](#13-exceptions)

---

## 1. INTRODUCTION ET CONVENTIONS BPMN

### 1.1 Objectif du Catalogue

Ce catalogue définit de manière exhaustive les 23 workflows opérationnels du Système National d'Identification et d'État Civil d'Haïti (SNISID). Chaque workflow est spécifié selon la norme **BPMN 2.0** (Business Process Model and Notation) avec:
- Diagramme narratif (description textuelle BPMN)
- Acteurs (Swimlanes)
- SLA (Service Level Agreement) par étape
- Règles DMN (Decision Model and Notation)
- Événements Kafka publiés/consommés
- Procédure mode hors-ligne (Offline)
- Gestion des exceptions et escalades

### 1.2 Conventions de Nommage

| Préfixe | Domaine                          |
|---------|----------------------------------|
| EC-N    | État Civil — Naissance           |
| EC-M    | État Civil — Mariage             |
| EC-D    | État Civil — Divorce             |
| EC-X    | État Civil — Décès               |
| EC-A    | État Civil — Adoption            |
| ID      | Identité Nationale (NIN)         |
| PP      | Passeport                        |
| FI      | Fiscalité / NIF                  |

### 1.3 Notation des SLA

| Code   | Signification              |
|--------|---------------------------|
| SLA-C  | Critique (< 4 heures)     |
| SLA-H  | Haut (< 24 heures)        |
| SLA-M  | Moyen (< 5 jours ouvrés)  |
| SLA-B  | Bas (< 30 jours)          |
| SLA-X  | Légal (défini par loi)    |

### 1.4 Acteurs Systémiques

| Code Actor  | Entité                                              |
|-------------|-----------------------------------------------------|
| DECL        | Déclarant (père, mère, responsable légal)           |
| OFCE-EC     | Officier d'État Civil                               |
| CHEF-EC     | Chef de la Section d'État Civil                     |
| SUPV-DEP    | Superviseur Départemental                           |
| TRIBUNAL    | Tribunal de Première Instance                       |
| PARQUET     | Parquet (Ministère Public)                          |
| ONI-CENTR   | ONI — Direction Centrale                            |
| SYS-SNISID  | Système SNISID (automatisé)                         |
| SYS-MSPP    | Système MSPP (FHIR HL7)                             |
| SYS-CEP     | Conseil Électoral Permanent                         |
| SYS-DGI     | Direction Générale des Impôts                       |
| SYS-OFATMA  | OFATMA (assurance)                                  |
| SYS-BRH     | Banque de la République d'Haïti                     |

---

## 2. ÉTAT CIVIL — NAISSANCE (EC-N)

### ═══════════════════════════════════════════════════════════
### WORKFLOW EC-N01 : NAISSANCE SIMPLE
### ═══════════════════════════════════════════════════════════

#### 2.1.1 Identification du Workflow

| Attribut              | Valeur                                    |
|-----------------------|-------------------------------------------|
| **ID Workflow**       | EC-N01                                    |
| **Nom**               | Déclaration de Naissance Simple           |
| **Version BPMN**      | 2.0                                       |
| **SLA Global**        | 24 heures (délai légal : 30 jours)        |
| **Priorité**          | HAUTE                                     |
| **Fréquence estimée** | 650 actes/jour (national)                 |
| **Déclencheur**       | Naissance vivante en établissement        |
| **Loi de référence**  | Code Civil Haïtien Art. 56, Loi 2011 ONI |

#### 2.1.2 Acteurs et Swimlanes

```
SWIMLANE 1: DÉCLARANT (Père/Mère/Responsable)
SWIMLANE 2: OFFICIER D'ÉTAT CIVIL (OFCE-EC)
SWIMLANE 3: SYSTÈME SNISID (SYS-SNISID)
SWIMLANE 4: MSPP / ÉTABLISSEMENT DE SANTÉ (SYS-MSPP)
SWIMLANE 5: ONI CENTRAL (ONI-CENTR)
```

#### 2.1.3 Description BPMN — Flux Nominal

```
[START EVENT: Naissance constatée]
    │
    ▼ SWIMLANE: MSPP
[TASK T01] Établissement de santé génère Attestation de Naissance (FHIR R4)
    │ SLA: 2h  │ Acteur: SYS-MSPP  │ Output: FHIR Bundle/Patient
    ▼
[GATEWAY G01: Accouchement en établissement ?]
    │ YES ──────────────────────────────────────────────────┐
    │ NO (domicile) → [TASK T01b] Certificat sage-femme     │
    └───────────────────────────────────────────────────────┘
    ▼
[MESSAGE EVENT: Notification SNISID ← MSPP]
    │ Kafka Topic: mspp.naissance.attestation.v1
    ▼ SWIMLANE: DÉCLARANT
[TASK T02] Déclarant se présente au Bureau d'État Civil
    │ SLA: 30 jours  │ Acteur: DECL
    │ Documents: Attestation MSPP, Pièces identité parents
    ▼ SWIMLANE: OFFICIER D'ÉTAT CIVIL
[TASK T03] Réception et vérification des documents
    │ SLA: 30 min  │ Acteur: OFCE-EC
    │ DMN: DMN-N01-01 (Vérification complétude dossier)
    ▼
[GATEWAY G02: Documents complets ?]
    │ NON → [TASK T03b] Délivrance Reçu de Dossier Incomplet
    │         → [END EVENT: Dossier incomplet — 7j pour compléter]
    │ OUI ──────────────────────────────────────────────────┐
    ▼                                                       │
[TASK T04] Saisie des données dans SNISID Offline/Online   │
    │ SLA: 15 min  │ Acteur: OFCE-EC                        │
    │ Input: Formulaire EC-N01-F001                         │
    ▼ SWIMLANE: SYSTÈME SNISID
[TASK T05] Déduplication biographique automatique
    │ SLA: 30 sec  │ Acteur: SYS-SNISID
    │ DMN: DMN-N01-02 (Score déduplication > 0.95 → alerte)
    ▼
[GATEWAY G03: Doublon potentiel détecté ?]
    │ OUI → [TASK T05b] Alerte Superviseur — Vérification manuelle
    │         SLA: 4h  │ Acteur: CHEF-EC
    │ NON ──────────────────────────────────────────────────┐
    ▼                                                       │
[TASK T06] Génération Numéro d'Acte (format: EC-AAAA-DDD-NNNNNN)
    │ SLA: 1 sec  │ Acteur: SYS-SNISID (automatique)
    ▼
[TASK T07] Calcul et réservation NIN provisoire
    │ SLA: 2 sec  │ Acteur: SYS-SNISID
    │ Output: NIN-PROVISOIRE-[UUID]
    ▼
[TASK T08] Génération PDF/A-3 + signature XAdES-LTA
    │ SLA: 10 sec  │ Acteur: SYS-SNISID
    │ Output: Acte de naissance PDF/A-3 avec QR Code
    ▼ SWIMLANE: OFFICIER D'ÉTAT CIVIL
[TASK T09] Validation et signature électronique OFCE-EC
    │ SLA: 5 min  │ Acteur: OFCE-EC
    │ Certificat: PKI-SNISID X.509 v3
    ▼
[TASK T10] Remise acte original au déclarant + copie certifiée
    │ SLA: 5 min  │ Acteur: OFCE-EC
    ▼ SWIMLANE: SYSTÈME SNISID
[TASK T11] Publication événement Kafka
    │ SLA: 5 sec  │ Acteur: SYS-SNISID
    │ Kafka Topics: 
    │   snisid.naissance.enregistree.v1 (fanout)
    │   cep.electeur.pre-enregistrement.v1
    │   dgi.contribuable.creation.v1
    │   ofatma.assure.creation.v1
    ▼
[TASK T12] Mise à jour registre central ONI
    │ SLA: 30 sec  │ Acteur: SYS-SNISID (async)
    ▼
[END EVENT: Naissance enregistrée — NIN provisoire assigné]
```

#### 2.1.4 Procédure Mode Hors-Ligne (EC-N01 Offline)

```
PROCÉDURE OFFLINE EC-N01-OFF
═══════════════════════════════════════════════════════

Pré-conditions:
  - Kit terrain certifié SNISID avec firmware v2.x+
  - Batterie > 40% ou alimentation solaire active
  - Clé HSM locale connectée et déverrouillée
  - Certificat offline valide (< 30 jours d'expiration)

Étapes:
  1. Officier sélectionne Mode Hors-Ligne dans l'application
  2. Saisie formulaire EC-N01-F001 (validation locale JSON Schema)
  3. Signature avec clé HSM locale (PKI offline)
  4. Stockage dans base SQLite chiffrée (AES-256-GCM)
  5. Génération QR Code local avec hash SHA-3 de l'acte
  6. Impression sur imprimante thermique portative
  7. Numéro provisoire: OFF-[CODEDEP]-[DATE]-[SEQ6]
  
Sync au retour connectivité:
  1. Moteur sync détecte connectivité (WiFi/4G/satellite)
  2. Upload delta vers SNISID Central API
  3. Validation et déduplication centralisée
  4. Attribution NIN définitif et numéro d'acte définitif
  5. Invalidation QR Code provisoire → nouveau QR Code
  6. Notification SMS au déclarant (si numéro fourni)

Délai maximum hors-ligne: 30 jours
Capacité stockage offline: 10 000 actes par kit
```

#### 2.1.5 Gestion des Exceptions EC-N01

| Exception                    | Code | Délai Résolution | Escalade       |
|------------------------------|------|------------------|----------------|
| Documents manquants          | E01  | 7 jours          | OFCE-EC        |
| Doublon détecté              | E02  | 4 heures         | CHEF-EC        |
| Désaccord nom parents        | E03  | 24 heures        | TRIBUNAL       |
| Attestation MSPP invalide    | E04  | 2 heures         | SYS-MSPP       |
| Signature PKI échouée        | E05  | 30 min           | ONI-CENTR      |
| NIN épuisé (espace)          | E06  | 1 heure          | ONI-CENTR DBA  |
| Fraude documentaire suspectée| E07  | 1 heure          | PARQUET/CISO   |

---

### ═══════════════════════════════════════════════════════════
### WORKFLOW EC-N02 : NAISSANCE PAR RECONNAISSANCE
### ═══════════════════════════════════════════════════════════

#### 2.2.1 Identification du Workflow

| Attribut              | Valeur                                      |
|-----------------------|---------------------------------------------|
| **ID Workflow**       | EC-N02                                      |
| **Nom**               | Naissance par Acte de Reconnaissance        |
| **SLA Global**        | 5 jours ouvrés                              |
| **Priorité**          | HAUTE                                       |
| **Déclencheur**       | Reconnaissance volontaire père non marié    |
| **Loi de référence**  | Code Civil Haïtien Art. 302-327             |

#### 2.2.2 Description BPMN — Flux Nominal

```
[START EVENT: Demande de reconnaissance déposée]
    │
    ▼ SWIMLANE: DÉCLARANT (Père reconnaissant)
[TASK T01] Déclarant se présente avec enfant et mère
    │ SLA: N/A (volontaire)  │ Acteur: DECL
    │ Documents requis:
    │   - Acte naissance enfant (si déjà enregistré)
    │   - Pièces identité père et mère
    │   - Déclaration sous serment de paternité
    ▼ SWIMLANE: OFFICIER D'ÉTAT CIVIL
[TASK T02] Réception demande et vérification identité
    │ SLA: 30 min  │ Acteur: OFCE-EC
    │ Vérification biométrique père (empreinte digitale)
    ▼
[GATEWAY G01: Enfant déjà enregistré ?]
    │ OUI → [TASK T02b] Récupération acte existant dans SNISID
    │ NON → [TASK T02c] Lancement parallèle EC-N01
    ▼
[TASK T03] Vérification consentement de la mère
    │ SLA: 1 heure  │ Acteur: OFCE-EC
    │ Si mère absente → procédure judiciaire EC-N04/N05
    ▼
[GATEWAY G02: Mère consentante présente ?]
    │ NON → [TASK T03b] Convocation mère (délai 15 jours)
    │         → [GATEWAY: Mère comparaît ?]
    │           NON → Transfert TRIBUNAL
    │ OUI ──────────────────────────────────────────────────┐
    ▼                                                       │
[TASK T04] Rédaction Acte de Reconnaissance
    │ SLA: 1 heure  │ Acteur: OFCE-EC
    ▼
[TASK T05] Lecture publique de l'acte et signatures
    │ SLA: 30 min  │ Acteur: OFCE-EC + DECL
    │ Signatures: Père + Mère + 2 témoins + OFCE-EC
    ▼ SWIMLANE: SYSTÈME SNISID
[TASK T06] Mise à jour registre — Ajout filiation paternelle
    │ SLA: 5 min  │ Acteur: SYS-SNISID
    │ Modification: ENFANT.pere_id ← NIN_PERE
    ▼
[TASK T07] Génération Acte de Reconnaissance (PDF/A-3 + XAdES-LTA)
    │ SLA: 30 sec  │ Acteur: SYS-SNISID
    ▼
[TASK T08] Mention marginale sur acte de naissance original
    │ SLA: 2 heures  │ Acteur: OFCE-EC
    ▼
[TASK T09] Publication événements Kafka
    │ Topics: snisid.filiation.reconnaissance.v1
    │         snisid.naissance.update.v1
    ▼
[END EVENT: Reconnaissance enregistrée — Filiation mise à jour]
```

#### 2.2.3 Procédure Mode Hors-Ligne EC-N02

```
PROCÉDURE OFFLINE EC-N02-OFF
═══════════════════════════════════════════════════════

Contraintes spécifiques:
  - Vérification biométrique père OBLIGATOIRE (ne peut être différée)
  - Si biométrie indisponible → procédure papier + validation centrale obligatoire
  - La mise à jour de filiation est différée jusqu'à sync

Étapes offline:
  1. Saisie déclaration reconnaissance (formulaire EC-N02-F001)
  2. Capture biométrique père (empreinte + photo) locale
  3. Génération acte provisoire signé par HSM local
  4. Stockage transactions: [RECON-OFF-ID, NIN_ENFANT_PROV, NIN_PERE_PROV]
  5. Remise document papier avec mention "PROVISOIRE — EN ATTENTE VALIDATION"

Sync priorité: HAUTE (délai sync max: 48h pour reconnaissance)
```

---

### ═══════════════════════════════════════════════════════════
### WORKFLOW EC-N03 : NAISSANCE PAR DÉCLARATION TARDIVE
### ═══════════════════════════════════════════════════════════

#### 2.3.1 Identification du Workflow

| Attribut              | Valeur                                      |
|-----------------------|---------------------------------------------|
| **ID Workflow**       | EC-N03                                      |
| **Nom**               | Déclaration de Naissance Tardive            |
| **SLA Global**        | 60 jours civils (instruction complète)      |
| **Priorité**          | MOYENNE                                     |
| **Déclencheur**       | Déclaration après 30 jours de la naissance  |
| **Loi de référence**  | Code Civil Haïtien Art. 58-60, Décret 2013 |

#### 2.3.2 Description BPMN — Flux Nominal

```
[START EVENT: Demande de déclaration tardive]
    │
    ▼ SWIMLANE: DÉCLARANT
[TASK T01] Dépôt dossier avec justificatifs retard
    │ SLA: N/A  │ Acteur: DECL
    │ Documents:
    │   - Attestation médicale (si naissance en établissement)
    │   - Déclarations sous serment de 2 témoins
    │   - Justification du retard (écrite + signée)
    │   - Pièces identité parents
    ▼ SWIMLANE: OFFICIER D'ÉTAT CIVIL
[TASK T02] Réception et enregistrement dossier
    │ SLA: 1 heure  │ Acteur: OFCE-EC
    │ Numéro de dossier: TAR-AAAA-DDD-NNNN
    ▼
[TASK T03] Instruction préliminaire et vérifications
    │ SLA: 15 jours  │ Acteur: OFCE-EC
    │   - Vérification base MSPP (naissances non déclarées)
    │   - Recherche acte existant (déduplication SNISID)
    │   - Vérification témoins (NIN valides et vivants)
    ▼ SWIMLANE: CHEF DE SECTION
[TASK T04] Rapport d'instruction et avis chef section
    │ SLA: 5 jours  │ Acteur: CHEF-EC
    │ DMN: DMN-N03-01 (Critères recevabilité tardive)
    ▼
[GATEWAY G01: Dossier recevable ?]
    │ NON → [TASK G01a] Notification refus motivé
    │         Délai appel: 30 jours → TRIBUNAL
    │ OUI ──────────────────────────────────────────────────┐
    ▼                                                       │
[TASK T05] Transmission dossier au Parquet
    │ SLA: 3 jours  │ Acteur: CHEF-EC → PARQUET
    ▼ SWIMLANE: PARQUET
[TASK T06] Réquisitions du Parquet
    │ SLA: 20 jours  │ Acteur: PARQUET
    │ Enquête complémentaire si nécessaire
    ▼
[GATEWAY G02: Parquet favorable ?]
    │ NON → [TASK G02a] Ordonnance de refus → Appel possible
    │ OUI ──────────────────────────────────────────────────┐
    ▼                                                       │
[TASK T07] Ordonnance du Parquet favorable
    │ SLA: 5 jours  │ Acteur: PARQUET
    ▼ SWIMLANE: OFFICIER D'ÉTAT CIVIL
[TASK T08] Établissement de l'acte de naissance tardive
    │ SLA: 3 jours  │ Acteur: OFCE-EC
    │ Mention obligatoire: "Déclaration tardive — Ordonnance N°XXX"
    ▼ SWIMLANE: SYSTÈME SNISID
[TASK T09] Enregistrement avec flag TARDIF
    │ SLA: 5 min  │ Acteur: SYS-SNISID
    │ Attributs: tardif=true, ordonnance_parquet=[ref]
    ▼
[TASK T10] Génération acte PDF/A-3 avec mention tardive
    │ SLA: 30 sec
    ▼
[TASK T11] Publication Kafka
    │ Topics: snisid.naissance.tardive.enregistree.v1
    ▼
[END EVENT: Acte tardif émis — Dossier archivé]
```

#### 2.3.3 Table de Décision DMN-N03-01 — Recevabilité Tardive

```
DMN-N03-01: Décision de recevabilité déclaration tardive
═══════════════════════════════════════════════════════════

INPUT 1: delai_jours (entier, jours depuis naissance)
INPUT 2: documents_medicaux (boolean)
INPUT 3: temoins_valides (entier, nb témoins NIN vérifiés)
INPUT 4: retard_justifie (boolean)

OUTPUT: decision (RECEVABLE | RECEVABLE_ENQUETE | IRRECEVABLE)
OUTPUT: motif (string)

RÈGLES:
┌─────────────────┬──────────────┬───────────────┬────────────────┬──────────────────────┬─────────────────────────────┐
│ delai_jours     │ docs_med     │ temoins_valid │ retard_justif  │ decision             │ motif                       │
├─────────────────┼──────────────┼───────────────┼────────────────┼──────────────────────┼─────────────────────────────┤
│ 31-365          │ true         │ >= 2          │ true           │ RECEVABLE            │ Dossier complet             │
│ 31-365          │ false        │ >= 2          │ true           │ RECEVABLE_ENQUETE    │ Documents médicaux absents  │
│ 31-365          │ -            │ < 2           │ -              │ IRRECEVABLE          │ Témoins insuffisants        │
│ > 365           │ true         │ >= 2          │ true           │ RECEVABLE_ENQUETE    │ Délai très long, enquête    │
│ > 365           │ false        │ -             │ false          │ IRRECEVABLE          │ Retard excessif non justifié│
│ 31-365          │ -            │ >= 2          │ false          │ RECEVABLE_ENQUETE    │ Justification insuffisante  │
└─────────────────┴──────────────┴───────────────┴────────────────┴──────────────────────┴─────────────────────────────┘
```

---

### ═══════════════════════════════════════════════════════════
### WORKFLOW EC-N04 : NAISSANCE PAR DÉCRET
### ═══════════════════════════════════════════════════════════

#### 2.4.1 Identification du Workflow

| Attribut              | Valeur                                         |
|-----------------------|------------------------------------------------|
| **ID Workflow**       | EC-N04                                         |
| **Nom**               | Naissance par Décret Présidentiel              |
| **SLA Global**        | Variable (selon décret) — nominal 30 jours     |
| **Priorité**          | TRÈS HAUTE                                     |
| **Déclencheur**       | Décret présidentiel ou ordonnance ministérielle |
| **Loi de référence**  | Constitution 1987 Art. 44, Loi ONI Art. 12     |

#### 2.4.2 Description BPMN — Flux Nominal

```
[START EVENT: Réception Décret / Ordonnance officielle]
    │
    ▼ SWIMLANE: ONI CENTRAL
[TASK T01] Réception et authentification document officiel
    │ SLA: 4 heures  │ Acteur: ONI-CENTR
    │ Vérification: Signature Premier Ministre + contre-seing MJSP
    │ DMN: DMN-N04-01 (Validité formelle décret)
    ▼
[TASK T02] Numérisation et enregistrement dans SNISID-GOV
    │ SLA: 2 heures  │ Acteur: ONI-CENTR
    │ Référence: DECRET-AAAA-NNNNN
    ▼
[TASK T03] Extraction données biographiques du décret
    │ SLA: 4 heures  │ Acteur: ONI-CENTR (analyste)
    ▼
[TASK T04] Recherche acte antérieur dans SNISID
    │ SLA: 30 min  │ Acteur: SYS-SNISID
    │ DMN: DMN-N04-02 (Acte existant à modifier ou créer)
    ▼
[GATEWAY G01: Acte existant à corriger ?]
    │ OUI → [TASK G01a] Procédure de rectification (EC-N04-RECT)
    │ NON → Création nouvel acte
    ▼
[TASK T05] Établissement acte par décret (avec référence officielle)
    │ SLA: 5 jours  │ Acteur: ONI-CENTR
    │ Mention obligatoire: "Par Décret N° [X] du [DATE]"
    ▼ SWIMLANE: SYSTÈME SNISID
[TASK T06] Enregistrement avec flag DECRET
    │ SLA: 5 min  │ Acteur: SYS-SNISID
    │ Attributs: type=DECRET, ref_decret=[N°], validite_decret=[date]
    ▼
[TASK T07] NIN définitif (si non existant) ou mise à jour NIN existant
    │ SLA: 1 min
    ▼
[TASK T08] Génération acte + notification MJSP + Journal Officiel
    │ SLA: 24 heures
    ▼
[TASK T09] Publication Kafka haute priorité
    │ Topics: snisid.naissance.decret.v1 (priority=HIGH)
    │         mjsp.decret.execution.confirme.v1
    ▼
[END EVENT: Acte par décret émis — Publication Journal Officiel]
```

---

### ═══════════════════════════════════════════════════════════
### WORKFLOW EC-N05 : NAISSANCE PAR JUGEMENT AU RANG DES MINUTES
### ═══════════════════════════════════════════════════════════

#### 2.5.1 Identification du Workflow

| Attribut              | Valeur                                          |
|-----------------------|-------------------------------------------------|
| **ID Workflow**       | EC-N05                                          |
| **Nom**               | Jugement supplétif d'acte de naissance          |
| **SLA Global**        | 6 mois (délai judiciaire)                       |
| **Priorité**          | MOYENNE                                         |
| **Déclencheur**       | Absence totale d'acte de naissance (adulte)     |
| **Loi de référence**  | CPC Art. 452-459, Code Civil Art. 62, Loi ONI  |

#### 2.5.2 Description BPMN — Flux Nominal

```
[START EVENT: Requête jugement supplétif]
    │
    ▼ SWIMLANE: DÉCLARANT (Requérant)
[TASK T01] Dépôt requête au Greffe du Tribunal
    │ SLA: N/A  │ Acteur: DECL (+ Avocat optionnel)
    │ Documents: Requête écrite, déclaration 2 témoins, 
    │            tout document prouvant l'identité
    ▼ SWIMLANE: TRIBUNAL
[TASK T02] Enregistrement au rôle et fixation audience
    │ SLA: 30 jours  │ Acteur: TRIBUNAL (Greffe)
    ▼
[TASK T03] Notification Parquet
    │ SLA: 15 jours avant audience  │ Acteur: TRIBUNAL
    ▼ SWIMLANE: PARQUET
[TASK T04] Réquisitions du Parquet
    │ SLA: 10 jours  │ Acteur: PARQUET
    │ Enquête de moralité et vérification identité
    ▼ SWIMLANE: TRIBUNAL
[TASK T05] Audience en chambre du conseil
    │ SLA: (selon rôle)  │ Acteur: TRIBUNAL
    │ Audition requérant + témoins + représentant Parquet
    ▼
[GATEWAY G01: Tribunal prononce jugement ?]
    │ REJET → [TASK G01a] Notification rejet + délai appel 30j
    │ ACCORD ───────────────────────────────────────────────┐
    ▼                                                       │
[TASK T06] Rédaction et prononcé du jugement supplétif
    │ SLA: 30 jours après audience  │ Acteur: TRIBUNAL
    ▼
[TASK T07] Délai d'opposition (30 jours)
    │ SLA: 30 jours
    ▼
[GATEWAY G02: Opposition formée ?]
    │ OUI → [TASK G02a] Procédure contradictoire → Nouveau jugement
    │ NON ──────────────────────────────────────────────────┐
    ▼                                                       │
[TASK T08] Signification jugement passé en force de chose jugée
    │ SLA: 10 jours  │ Acteur: TRIBUNAL
    ▼
[TASK T09] Transmission jugement à l'Officier d'État Civil
    │ SLA: 5 jours  │ Acteur: TRIBUNAL → OFCE-EC
    ▼ SWIMLANE: OFFICIER D'ÉTAT CIVIL
[TASK T10] Transcription jugement au registre d'état civil
    │ SLA: 5 jours  │ Acteur: OFCE-EC
    │ Mention: "Jugement supplétif N° [X] du [DATE] — Tribunal [Y]"
    ▼ SWIMLANE: SYSTÈME SNISID
[TASK T11] Enregistrement avec flag JUGEMENT_SUPLETIF
    │ SLA: 5 min  │ Acteur: SYS-SNISID
    │ Attributs: type=JUGEMENT, ref_jugement=[N°], tribunal=[code]
    ▼
[TASK T12] Attribution NIN définitif
    │ SLA: 1 min
    ▼
[TASK T13] Publication Kafka
    │ Topics: snisid.naissance.jugement.v1
    │         snisid.identite.nin.attribue.v1
    ▼
[END EVENT: Acte par jugement — NIN définitif attribué]
```

#### 2.5.3 Procédure Mode Hors-Ligne EC-N05

```
NOTE IMPORTANTE: EC-N05 NE PEUT PAS ÊTRE TRAITÉ EN MODE HORS-LIGNE
La transcription d'un jugement exige une connexion centrale obligatoire
pour vérification du jugement (API Tribunal) avant enregistrement.

Procédure dégradée:
  1. Réception jugement en main propre (copie certifiée)
  2. Photocopie numérisée et chiffrement local
  3. Transmission différée lors de la reconnexion
  4. NIN attribué uniquement après validation centrale
  5. Document intermédiaire: ATTESTATION D'ATTENTE NIN (valide 30j)
```

---

## 3. ÉTAT CIVIL — MARIAGE (EC-M)

### WORKFLOW EC-M01 : MARIAGE CIVIL

| Attribut        | Valeur                              |
|-----------------|-------------------------------------|
| **ID**          | EC-M01                              |
| **Nom**         | Mariage Civil                       |
| **SLA Global**  | 30 jours (publication bans + acte)  |
| **Loi**         | Code Civil Art. 142-173             |

```
[START EVENT: Demande de mariage déposée]
    ▼
[TASK T01] Dépôt dossier mariage (futurs époux)
    │ Documents: pièces identité, actes naissance, certificats médicaux,
    │            autorisation parentale si mineur, résidence
    ▼
[TASK T02] Vérification empêchements légaux (SNISID)
    │ SLA: 1 heure  │ DMN: DMN-M01-01
    │ Vérifie: célibat, âge légal, liens consanguinité, tutelle
    ▼
[GATEWAY G01: Empêchements détectés ?]
    │ OUI → Notification refus + recours TRIBUNAL
    │ NON ▼
[TASK T03] Publication des bans (21 jours)
    │ SLA: 21 jours  │ Affichage physique + publication SNISID public
    ▼
[GATEWAY G02: Opposition aux bans ?]
    │ OUI → Instruction opposition (PARQUET + TRIBUNAL)
    │ NON ▼
[TASK T04] Célébration du mariage
    │ SLA: J+7 après bans  │ Acteur: OFCE-EC
    │ Présence: 2 époux + 2 témoins + OFCE-EC
    ▼
[TASK T05] Rédaction et lecture de l'acte de mariage
    │ Signatures: 2 époux + 2 témoins + OFCE-EC
    ▼
[TASK T06] Enregistrement SNISID + mise à jour statut civil époux
    │ Topics: snisid.mariage.celebre.v1, snisid.etat-civil.update.v1
    ▼
[TASK T07] Mentions marginales sur actes naissance des époux
    ▼
[END EVENT: Mariage enregistré — Livret de famille généré]
```

### WORKFLOW EC-M02 : MARIAGE RELIGIEUX À EFFET CIVIL

| Attribut  | Valeur                              |
|-----------|-------------------------------------|
| **ID**    | EC-M02                              |
| **Nom**   | Transcription mariage religieux     |
| **SLA**   | 30 jours après cérémonie religieuse |

```
[START EVENT: Certificat mariage religieux reçu]
    ▼
[TASK T01] Vérification validité culte reconnu (liste ONI)
    │ DMN: DMN-M02-01 (Cultes reconnus par l'État haïtien)
    ▼
[TASK T02] Vérification condition préalable mariage civil (si exigée)
    ▼
[TASK T03] Transcription acte — mêmes procédures EC-M01 T02 à T07
    ▼
[END EVENT: Mariage religieux transcrit au registre civil]
```

---

## 4. ÉTAT CIVIL — DIVORCE (EC-D)

### WORKFLOW EC-D01 : DIVORCE CONTENTIEUX

| Attribut  | Valeur                          |
|-----------|---------------------------------|
| **ID**    | EC-D01                          |
| **Nom**   | Divorce par décision judiciaire |
| **SLA**   | Variable (6 mois - 2 ans)       |
| **Loi**   | Code Civil Art. 208-245         |

```
[START EVENT: Jugement de divorce définitif reçu]
    ▼
[TASK T01] Réception copie exécutoire du jugement
    │ Acteur: OFCE-EC
    ▼
[TASK T02] Vérification authenticité jugement (API Tribunal SNISID)
    │ SLA: 2 heures
    ▼
[TASK T03] Transcription jugement divorce au registre
    │ SLA: 5 jours  │ Acteur: OFCE-EC
    ▼
[TASK T04] Mentions marginales sur acte de mariage original
    │ Mention: "Dissous par jugement N° [X] du [DATE]"
    ▼
[TASK T05] Mise à jour statut civil des ex-époux dans SNISID
    │ Statut: DIVORC(É/ÉE)  │ SLA: 1 heure
    ▼
[TASK T06] Publication événement Kafka
    │ Topics: snisid.mariage.dissous.v1, snisid.etat-civil.update.v1
    │         dgi.contribuable.statut-civil.update.v1
    ▼
[END EVENT: Divorce transcrit — Statuts mis à jour]
```

### WORKFLOW EC-D02 : DIVORCE PAR CONSENTEMENT MUTUEL

| Attribut  | Valeur                          |
|-----------|---------------------------------|
| **ID**    | EC-D02                          |
| **Nom**   | Divorce par consentement mutuel |
| **SLA**   | 90 jours (délai légal)          |

```
[START EVENT: Convention de divorce déposée chez notaire]
    ▼
[TASK T01] Enregistrement convention (Notaire → SNISID)
    │ Acteur: NOTAIRE (accès SNISID via API-GW)
    ▼
[TASK T02] Délai de réflexion légal (30 jours)
    ▼
[GATEWAY G01: Rétractation dans délai ?]
    │ OUI → Annulation convention → END
    │ NON ▼
[TASK T03] Homologation par le Tribunal
    │ SLA: 60 jours
    ▼
[TASK T04] → Même flux EC-D01 T01 à T06]
    ▼
[END EVENT: Divorce par consentement transcrit]
```

---

## 5. ÉTAT CIVIL — DÉCÈS (EC-X)

### WORKFLOW EC-X01 : DÉCLARATION DE DÉCÈS

| Attribut  | Valeur                                |
|-----------|---------------------------------------|
| **ID**    | EC-X01                                |
| **Nom**   | Déclaration et enregistrement du décès |
| **SLA**   | 24 heures (légal: 24h)               |
| **Loi**   | Code Civil Art. 77-89                |

```
[START EVENT: Décès constaté]
    ▼
[TASK T01] Établissement certificat médical de décès (MSPP)
    │ SLA: 2h  │ Acteur: SYS-MSPP (médecin)
    │ FHIR Resource: Observation/cause-of-death
    ▼
[TASK T02] Déclaration par famille au bureau d'état civil
    │ SLA: 24h  │ Documents: Certificat médical, pièce identité déclarant
    ▼
[TASK T03] Vérification identité du défunt dans SNISID
    │ SLA: 15 min  │ Recherche par NIN, nom, date naissance
    ▼
[TASK T04] Enregistrement acte de décès
    │ SLA: 30 min  │ Acteur: OFCE-EC
    ▼ SWIMLANE: SYSTÈME SNISID
[TASK T05] Mise à jour statut NIN: VIVANT → DÉCÉDÉ
    │ SLA: 5 min  │ CRITIQUE — action irréversible
    │ Validation double: OFCE-EC + CHEF-EC obligatoire
    ▼
[TASK T06] Cascade de notifications (FANOUT)
    │ SLA: 30 min (async)
    │ Kafka Topics publiés simultanément:
    │   snisid.deces.enregistre.v1 (master event)
    │   ├── cep.electeur.radiation.v1 (radiation liste électorale)
    │   ├── dgi.contribuable.deces.v1 (clôture dossier fiscal)
    │   ├── ofatma.assure.deces.v1 (déclenchement assurance décès)
    │   ├── brh.compte.gel.v1 (notification système bancaire)
    │   ├── mspp.patient.deces.v1 (clôture dossier médical)
    │   ├── csc.pension.suspension.v1 (ONA/CSC pension)
    │   └── passeport.invalidation.v1 (invalidation document voyage)
    ▼
[TASK T07] Génération acte de décès PDF/A-3 + QR Code
    ▼
[TASK T08] Mention marginale sur acte de naissance (si accessible)
    │ SLA: 48 heures
    ▼
[TASK T09] Archivage dossier défunt (statut ARCHIVÉ)
    ▼
[END EVENT: Décès enregistré — Cascade notifications complétée]
```

#### 5.1 Gestion Cascade Décès — Procédures d'Exception

| Notification      | Délai max | Retry max | Escalade si échec |
|-------------------|-----------|-----------|-------------------|
| CEP Radiation     | 24h       | 3         | CEP Administrateur |
| DGI Clôture       | 48h       | 5         | DGI Chef Division  |
| OFATMA Décès      | 4h        | 10        | OFATMA Direction   |
| BRH Gel           | 1h        | 5         | BRH Risk Officer   |
| MSPP Clôture      | 72h       | 3         | MSPP DI            |
| ONA Pension       | 24h       | 5         | ONA Direction      |
| Passeport Inv.    | 1h        | 10        | ONI Passeports     |

---

## 6. ADOPTION (EC-A)

### WORKFLOW EC-A01 : ADOPTION SIMPLE

| Attribut  | Valeur                          |
|-----------|---------------------------------|
| **ID**    | EC-A01                          |
| **Nom**   | Adoption simple                 |
| **SLA**   | 90 jours (voie judiciaire)      |
| **Loi**   | Code Civil Art. 353-370         |

```
[START EVENT: Jugement d'adoption prononcé]
    ▼
[TASK T01] Réception copie jugement d'adoption (Tribunal → ONI)
    ▼
[TASK T02] Vérification authenticité + réf. IBESR si enfant mineur
    │ SLA: 4h  │ Vérification API IBESR (Institut Bien-Être Social)
    ▼
[TASK T03] Transcription jugement au registre d'état civil
    │ SLA: 5 jours  │ Mention: adoption simple, lien biologique conservé
    ▼
[TASK T04] MISE À JOUR SNISID:
    │   - Ajout lien adoptif (type=ADOPTIF_SIMPLE)
    │   - Conservation lien biologique (type=BIOLOGIQUE)
    │   - Mise à jour NOM si changement accordé
    ▼
[TASK T05] Génération nouveau extrait de naissance (mention adoption)
    ▼
[END EVENT: Adoption simple enregistrée]
```

### WORKFLOW EC-A02 : ADOPTION PLÉNIÈRE

| Attribut  | Valeur                              |
|-----------|-------------------------------------|
| **ID**    | EC-A02                              |
| **Nom**   | Adoption plénière                   |
| **SLA**   | 6 mois (incluant délais judiciaires) |

```
[START EVENT: Jugement adoption plénière définitif]
    ▼
[TASK T01] → T02 identiques à EC-A01
    ▼
[TASK T03] Annulation acte de naissance biologique (SCELLÉ)
    │ CRITIQUE: Acte biologique scellé — inaccessible sauf tribunal
    │ Nouvelle filiation: adoptants = parents légaux EXCLUSIFS
    ▼
[TASK T04] Création NOUVEAU acte de naissance (adoptant comme parents)
    │ Mention: Acte de naissance établi suite adoption plénière
    ▼
[TASK T05] Suppression lien biologique dans SNISID
    │ Conservation: ID interne pour statistiques (anonymisé)
    ▼
[TASK T06] Génération nouveau NIN si enfant étranger
    ▼
[END EVENT: Adoption plénière — Nouveau acte de naissance émis]
```

---

## 7. IDENTITÉ NATIONALE (ID)

### WORKFLOW ID-01 : ÉMISSION CARTE NATIONALE D'IDENTITÉ (CNI)

| Attribut  | Valeur                              |
|-----------|-------------------------------------|
| **ID**    | ID-01                               |
| **Nom**   | Émission Carte Nationale d'Identité |
| **SLA**   | 30 jours                            |
| **Loi**   | Loi ONI 2011, Décret CNI 2012       |

```
[START EVENT: Demande de CNI reçue]
    ▼
[TASK T01] Vérification existence acte de naissance dans SNISID
    │ Prérequis: NIN provisoire ou NIN définitif existant
    ▼
[TASK T02] Capture biométrique (empreintes 10 doigts + iris + photo)
    │ SLA: 15 min  │ Équipement: station biométrique certifiée
    ▼
[TASK T03] Déduplication biométrique AFIS (Automated Fingerprint ID)
    │ SLA: 2 min  │ Seuil: score > 250 → doublon
    ▼
[GATEWAY G01: Doublon biométrique ?]
    │ OUI → Investigation fraude (CISO + OFCE-EC)
    │ NON ▼
[TASK T04] Attribution NIN définitif (si provisoire)
    ▼
[TASK T05] Personnalisation carte (chip + données biométriques)
    │ SLA: 24h (centre de personnalisation)
    ▼
[TASK T06] Contrôle qualité et vérification chip
    ▼
[TASK T07] Remise au demandeur avec signature
    ▼
[END EVENT: CNI émise — Biométrie enregistrée AFIS]
```

### WORKFLOW ID-02 : RENOUVELLEMENT CNI

| Attribut  | Valeur                    |
|-----------|---------------------------|
| **ID**    | ID-02                     |
| **Nom**   | Renouvellement CNI         |
| **SLA**   | 15 jours                  |

```
[Similaire à ID-01 avec T01: Vérification CNI précédente + invalidation]
[+ Archivage ancienne biométrie avec timestamp invalidation]
```

### WORKFLOW ID-03 : CORRECTION / RECTIFICATION

| Attribut  | Valeur           |
|-----------|------------------|
| **ID**    | ID-03            |
| **Nom**   | Rectification CNI |
| **SLA**   | 10 jours         |

```
[START EVENT: Demande de correction déposée]
    ▼
[TASK T01] Identification erreur (typographique vs substantielle)
    │ DMN: DMN-ID03-01 (Nature de l'erreur et procédure applicable)
    ▼
[GATEWAY G01: Erreur typographique seule ?]
    │ OUI → Correction administrative (Chef Section)
    │ NON → Procédure judiciaire (rectification d'acte)
    ▼
[TASK T02] Correction dans SNISID + journal d'audit complet
    │ OBLIGATOIRE: Traçabilité avant/après correction
    ▼
[END EVENT: CNI corrigée — Ancienne invalidée]
```

---

## 8. PASSEPORT (PP)

### WORKFLOW PP-01 : ÉMISSION PASSEPORT BIOMÉTRIQUE

| Attribut  | Valeur               |
|-----------|----------------------|
| **ID**    | PP-01                |
| **Nom**   | Émission Passeport   |
| **SLA**   | 21 jours             |
| **Loi**   | Décret Passeport 2013 |

```
[START EVENT: Demande passeport déposée]
    ▼
[TASK T01] Vérification NIN + CNI valide dans SNISID
    ▼
[TASK T02] Vérification liste noire / interdits de sortie (MJSP)
    │ DMN: DMN-PP01-01 (Statut légal voyageur)
    ▼
[TASK T03] Capture biométrique (si première demande ou biométrie outdatée)
    ▼
[TASK T04] Vérification FBI Watchlist / INTERPOL (API sécurisée)
    │ SLA: 30 min
    ▼
[TASK T05] Personnalisation passeport ICAO 9303 (chip RFID + MRZ)
    │ SLA: 5 jours (impression + personnalisation)
    ▼
[TASK T06] Contrôle qualité (lecture MRZ + test chip BAC/PACE)
    ▼
[TASK T07] Remise sécurisée au demandeur
    ▼
[END EVENT: Passeport émis — Enregistrement INTERPOL]
```

### WORKFLOW PP-02 : RENOUVELLEMENT PASSEPORT

| Attribut  | Valeur              |
|-----------|---------------------|
| **ID**    | PP-02               |
| **SLA**   | 15 jours            |

### WORKFLOW PP-03 : PASSEPORT D'URGENCE

| Attribut  | Valeur                       |
|-----------|------------------------------|
| **ID**    | PP-03                        |
| **SLA**   | 48 heures                    |
| **Note**  | Validité limitée (6 mois)    |

---

## 9. FISCALITÉ (FI)

### WORKFLOW FI-01 : ATTRIBUTION NIF (NUMÉRO D'IDENTIFICATION FISCALE)

| Attribut  | Valeur                      |
|-----------|-----------------------------|
| **ID**    | FI-01                       |
| **Nom**   | Attribution NIF via SNISID  |
| **SLA**   | 5 jours ouvrés              |
| **Loi**   | Loi DGI, Accord SNISID-DGI  |

```
[START EVENT: Événement Kafka reçu de SNISID]
    │ Topic: dgi.contribuable.creation.v1 (déclenché par naissance)
    │ Topic: dgi.contribuable.adulte.v1 (18 ans atteint)
    ▼ SWIMLANE: DGI SYSTÈME
[TASK T01] Réception et validation event SNISID
    │ Vérification: signature Kafka JWS, schéma Avro
    ▼
[TASK T02] Recherche NIF existant (cross-référence NIN)
    │ DMN: DMN-FI01-01 (NIF existant ou création)
    ▼
[GATEWAY G01: NIF existant ?]
    │ OUI → Mise à jour cross-référence NIN↔NIF
    │ NON ▼
[TASK T03] Attribution NIF (algorithme DGI)
    │ Format: HTI-[DEPT]-[YYYY]-[NNNNNNNN]
    ▼
[TASK T04] Notification SNISID (callback API)
    │ Topic: snisid.nin.nif.lie.v1
    ▼
[TASK T05] Création dossier contribuable dans SYGEF (DGI)
    ▼
[END EVENT: NIF attribué — NIN↔NIF liés dans les deux systèmes]
```

---

## 10. TABLES DMN — RÈGLES DE DÉCISION

### 10.1 DMN-N01-01 : Complétude Dossier Naissance Simple

```
Table: Vérification complétude dossier EC-N01
═══════════════════════════════════════════════

Inputs:
  - attestation_mspp: boolean
  - identite_pere: boolean
  - identite_mere: boolean
  - formulaire_ec_n01: boolean

Output:
  - complet: boolean
  - manquants: list<string>

RÈGLES:
R1: attestation_mspp=true AND id_pere=true AND id_mere=true AND form=true
    → complet=true, manquants=[]

R2: attestation_mspp=false AND autres=true
    → complet=false, manquants=["attestation_mspp"]

R3: (id_pere=false OR id_mere=false)
    → complet=false, manquants=["identite_parent"]

R4: formulaire_ec_n01=false
    → complet=false, manquants=["formulaire_ec_n01"]

Priority: Highest priority rule wins (R2 > R3 > R4 > R1)
```

### 10.2 DMN-N01-02 : Déduplication

```
Table: Alerte doublon naissance
════════════════════════════════

Input:
  - score_similarite: decimal (0.0 - 1.0)
  - nom_exact: boolean
  - date_naissance_exacte: boolean
  - lieu_naissance_exact: boolean

Output:
  - action: PASSER | ALERTER | BLOQUER
  - niveau: INFO | WARNING | CRITICAL

RÈGLES:
R1: score < 0.70 → action=PASSER
R2: score 0.70-0.84 → action=ALERTER, niveau=WARNING
R3: score 0.85-0.94 → action=ALERTER, niveau=CRITICAL
R4: score >= 0.95 → action=BLOQUER, niveau=CRITICAL
R5: nom_exact=true AND date_naissance_exacte=true → action=BLOQUER
R6: R4 OR R5 → Transmission immédiate CHEF-EC
```

### 10.3 DMN-M01-01 : Empêchements au Mariage

```
Table: Vérification empêchements légaux mariage
════════════════════════════════════════════════

Inputs:
  - statut_civil_epoux1: CÉLIBATAIRE | MARIÉ | DIVORCÉ | VEUF
  - statut_civil_epoux2: CÉLIBATAIRE | MARIÉ | DIVORCÉ | VEUF
  - age_epoux1: entier
  - age_epoux2: entier
  - lien_parente: AUCUN | ASCENDANT | DESCENDANT | FRERE_SOEUR | COUSIN
  - tuteur_consent: boolean (si mineur)

Output:
  - autorise: boolean
  - empechement: string

RÈGLES:
R1: statut_civil_epoux1=MARIÉ OR statut_civil_epoux2=MARIÉ
    → autorise=false, empechement="Bigamie"

R2: age_epoux1 < 18 OR age_epoux2 < 18
    → autorise=false, empechement="Minorité (sauf autorisation TPI)"
    [Exception: tuteur_consent=true AND âge >= 16 → RECEVABLE_EXCEPTION]

R3: lien_parente IN [ASCENDANT, DESCENDANT, FRERE_SOEUR]
    → autorise=false, empechement="Inceste — Empêchement absolu"

R4: TOUTES autres conditions → autorise=true
```

### 10.4 DMN-PP01-01 : Statut Légal Voyageur

```
Table: Vérification statut voyageur pour passeport
════════════════════════════════════════════════════

Inputs:
  - interdiction_sortie: boolean (MJSP)
  - mandat_arret_actif: boolean (Parquet)
  - dette_fiscale_bloquante: boolean (DGI)
  - document_requisition: boolean (Tribunal)

Output:
  - autorisation: AUTORISÉ | SUSPENDU | REFUSÉ
  - motif: string

RÈGLES:
R1: mandat_arret_actif=true → REFUSÉ, "Mandat d'arrêt actif"
R2: interdiction_sortie=true → SUSPENDU, "Interdiction sortie territoire"
R3: document_requisition=true → SUSPENDU, "Réquisition judiciaire"
R4: dette_fiscale_bloquante=true → SUSPENDU, "Dette fiscale > seuil"
R5: TOUT false → AUTORISÉ
Priority: R1 > R2 > R3 > R4 > R5
```

---

## 11. INTÉGRATION BUS NATIONAL — KAFKA EVENTS

### 11.1 Topologie des Topics SNISID

```
KAFKA CLUSTER: snisid-kafka-prod (3 brokers, RF=3, min.insync=2)
══════════════════════════════════════════════════════════════════

NAMESPACE: snisid.*  (publié par SNISID)
NAMESPACE: mspp.*   (publié par MSPP FHIR)
NAMESPACE: cep.*    (publié par CEP)
NAMESPACE: dgi.*    (publié par DGI)
NAMESPACE: ofatma.* (publié par OFATMA)
NAMESPACE: brh.*    (publié par BRH)
```

### 11.2 Catalogue des Events par Workflow

| Topic Kafka                           | Workflow   | Schema     | Retention | Partitions |
|---------------------------------------|------------|------------|-----------|------------|
| snisid.naissance.enregistree.v1       | EC-N01-05  | Avro       | 365j      | 12         |
| snisid.naissance.tardive.enregistree.v1 | EC-N03   | Avro       | 365j      | 6          |
| snisid.naissance.decret.v1            | EC-N04     | Avro       | 365j      | 3          |
| snisid.naissance.jugement.v1          | EC-N05     | Avro       | 365j      | 3          |
| snisid.filiation.reconnaissance.v1    | EC-N02     | Avro       | 365j      | 6          |
| snisid.mariage.celebre.v1             | EC-M01     | Avro       | 365j      | 6          |
| snisid.mariage.dissous.v1             | EC-D01/D02 | Avro       | 365j      | 6          |
| snisid.deces.enregistre.v1            | EC-X01     | Avro       | forever   | 12         |
| snisid.adoption.enregistree.v1        | EC-A01/A02 | Avro       | 365j      | 3          |
| snisid.identite.nin.attribue.v1       | ID-01      | Avro       | forever   | 12         |
| snisid.identite.cni.emise.v1          | ID-01      | Avro       | 365j      | 6          |
| snisid.passeport.emis.v1              | PP-01      | Avro       | 365j      | 6          |
| snisid.nin.nif.lie.v1                 | FI-01      | Avro       | 365j      | 6          |
| mspp.naissance.attestation.v1         | EC-N01     | FHIR/JSON  | 90j       | 6          |
| cep.electeur.radiation.v1             | EC-X01     | Avro       | 365j      | 6          |
| dgi.contribuable.creation.v1          | EC-N01/FI  | Avro       | 365j      | 6          |

### 11.3 Schéma Avro — Event Naissance (exemple)

```json
{
  "namespace": "ht.gov.snisid.events",
  "type": "record",
  "name": "NaissanceEnregistree",
  "version": "1.0.0",
  "fields": [
    {"name": "event_id", "type": "string", "doc": "UUID v4"},
    {"name": "event_time", "type": "long", "logicalType": "timestamp-millis"},
    {"name": "workflow_id", "type": "string", "doc": "EC-N01 à EC-N05"},
    {"name": "acte_numero", "type": "string"},
    {"name": "nin_provisoire", "type": ["null", "string"], "default": null},
    {"name": "nin_definitif", "type": ["null", "string"], "default": null},
    {"name": "nom", "type": "string"},
    {"name": "prenoms", "type": {"type": "array", "items": "string"}},
    {"name": "date_naissance", "type": "string", "logicalType": "date"},
    {"name": "lieu_naissance_code", "type": "string", "doc": "Code commune ONI"},
    {"name": "sexe", "type": {"type": "enum", "name": "Sexe", "symbols": ["M", "F"]}},
    {"name": "nin_pere", "type": ["null", "string"], "default": null},
    {"name": "nin_mere", "type": ["null", "string"], "default": null},
    {"name": "officier_ec_id", "type": "string"},
    {"name": "bureau_ec_code", "type": "string"},
    {"name": "offline_origin", "type": "boolean", "default": false},
    {"name": "hash_acte", "type": "string", "doc": "SHA3-256 du PDF/A-3"},
    {"name": "signature_jws", "type": "string", "doc": "JWS compact du payload"}
  ]
}
```

---

## 12. PROCÉDURES MODE HORS-LIGNE (SYNTHÈSE)

### 12.1 Tableau de Compatibilité Offline par Workflow

| Workflow | Offline Possible | Mode Dégradé | Délai Sync Max | Priorité Sync |
|----------|-----------------|--------------|----------------|---------------|
| EC-N01   | OUI (complet)   | Non requis   | 30 jours       | NORMALE       |
| EC-N02   | OUI (partiel)   | Biométrie obligatoire | 48h | HAUTE    |
| EC-N03   | NON — instruction | Réception docs | N/A          | N/A           |
| EC-N04   | NON — décret    | Réception scannée | N/A         | N/A           |
| EC-N05   | NON — jugement  | Attestation d'attente | N/A       | N/A           |
| EC-M01   | OUI (partiel)   | Bans en ligne requis | 7 jours   | HAUTE         |
| EC-D01   | NON — judiciaire | N/A          | N/A            | N/A           |
| EC-X01   | OUI (critique)  | Sync 48h max | 48 heures      | CRITIQUE      |
| EC-A01   | NON — judiciaire | N/A          | N/A            | N/A           |
| ID-01    | OUI (biométrie locale) | N/A   | 24 heures      | HAUTE         |
| PP-01    | NON — chip ICAO | N/A          | N/A            | N/A           |
| FI-01    | AUTO (event-driven) | N/A     | 72 heures      | NORMALE       |

### 12.2 Procédure Générale Retour Connectivité

```
PROCÉDURE SYNC-RETOUR-CONN
═══════════════════════════

1. DÉTECTION: Moteur sync détecte connectivité (WiFi/4G/VSAT)
   - Ping test: api.snisid.ht (timeout 5s, 3 essais)
   - Test SSL: validation certificat serveur
   - Test auth: vérification token JWT applicatif

2. AUTHENTIFICATION: Handshake mutuel TLS 1.3
   - Certificat client kit terrain (PKI-SNISID)
   - Vérification CRL/OCSP du certificat kit

3. NÉGOCIATION: Échange manifest de sync
   - Kit envoie: {kit_id, last_sync_ts, nb_actes_pending, merkle_root}
   - Serveur répond: {accept, expected_merkle, conflicts_detected}

4. TRANSMISSION: Upload par lots (batch de 100 actes max)
   - Ordre: CRITIQUE → HAUTE → NORMALE
   - Compression: zstd niveau 3
   - Chiffrement: AES-256-GCM (clé de session négociée)

5. VALIDATION: Serveur valide chaque acte
   - Déduplication biographique centrale
   - Détection conflits (vector clocks)
   - Attribution NIN définitifs

6. CONFIRMATION: Serveur retourne résultats
   - ACK par acte (numéro définitif ou erreur)
   - Mise à jour Merkle tree local
   - Mise à jour last_sync_ts

7. NOTIFICATION: SMS aux bénéficiaires (si numéro disponible)
   - "Votre acte [TYPE] a été validé. NIN: [NNNNN]"
```

---

## 13. GESTION DES EXCEPTIONS ET ESCALADES

### 13.1 Matrice d'Escalade Globale

| Niveau | Délai non-résolution | Escalade vers        | Canal               |
|--------|---------------------|---------------------|---------------------|
| L1     | > SLA étape         | Chef Section         | Notification app    |
| L2     | > 2x SLA étape      | Superviseur Dép.     | SMS + Email         |
| L3     | > 3x SLA étape      | Direction ONI        | Appel + Ticket P1   |
| L4     | > SLA global        | Ministre MJSP        | Rapport officiel    |
| FRAUDE | Immédiat            | PARQUET + CISO SOC  | Appel direct + P1   |

### 13.2 Procédure de Traitement des Fraudes Documentaires

```
PROC-FRAUDE-001: Suspicion fraude documentaire
═══════════════════════════════════════════════

DÉCLENCHEURS:
  - Score déduplication > 0.95 (doublon probable)
  - Document présentant signes altération (UV check)
  - Signalement officier (intuition professionnelle)
  - Alerte SNISID (algorithme ML fraude)

ACTIONS IMMÉDIATES (< 15 minutes):
  1. Gel de la transaction en cours
  2. Notification CISO SOC (alerte FR-[TICKET])
  3. Conservation de toutes les preuves (photos, docs scannés)
  4. Isolation du Kit terrain si offline

ACTIONS DANS L'HEURE:
  1. Rapport incident (formulaire FR-001)
  2. Transmission PARQUET (si présomption crime)
  3. Audit trail complet exporté (signé HSM)
  4. Alerte bureaux adjacents (réseau intranet ONI)

ACTIONS DANS 24 HEURES:
  1. Analyse forensique documents
  2. Rapport CISO à Direction ONI
  3. Si réseau de fraude → alerte nationale
```

### 13.3 SLA Tableau de Bord — KPIs Workflows

| KPI                              | Cible    | Alerte   | Critique  | Fréquence mesure |
|----------------------------------|----------|----------|-----------|-----------------|
| Taux de complétion EC-N01        | > 99.5%  | < 99%    | < 95%     | Temps réel       |
| Délai moyen enregistrement       | < 30 min | > 1h     | > 4h      | Horaire          |
| Taux doublon détecté/résolu      | > 99%    | < 98%    | < 95%     | Quotidien        |
| Taux sync offline dans délai     | > 98%    | < 95%    | < 90%     | Quotidien        |
| Disponibilité bus Kafka          | > 99.9%  | < 99.5%  | < 99%     | Temps réel       |
| Temps génération PDF/A-3         | < 10 sec | > 30 sec | > 60 sec  | Temps réel       |
| Taux signature PKI réussie       | > 99.9%  | < 99.5%  | < 99%     | Temps réel       |

---

## ANNEXES

### Annexe A : Référence Légale Complète

| Article / Loi              | Workflow(s) concerné(s)     |
|----------------------------|-----------------------------|
| Code Civil Art. 56-62      | EC-N01 à EC-N05              |
| Code Civil Art. 142-173    | EC-M01, EC-M02               |
| Code Civil Art. 208-245    | EC-D01, EC-D02               |
| Code Civil Art. 77-89      | EC-X01                       |
| Code Civil Art. 302-327    | EC-N02                       |
| Code Civil Art. 353-370    | EC-A01, EC-A02               |
| CPC Art. 452-459           | EC-N05                       |
| Loi ONI 2011               | Tous workflows               |
| Décret CNI 2012            | ID-01 à ID-03                |
| Décret Passeport 2013      | PP-01 à PP-03                |
| Loi DGI + Accord SNISID    | FI-01                        |

### Annexe B : Formulaires Référencés

| Formulaire     | Description                           | Format    |
|----------------|---------------------------------------|-----------|
| EC-N01-F001    | Déclaration naissance simple          | PDF/A-3   |
| EC-N02-F001    | Demande reconnaissance paternité      | PDF/A-3   |
| EC-N03-F001    | Déclaration naissance tardive         | PDF/A-3   |
| EC-M01-F001    | Demande de mariage civil              | PDF/A-3   |
| EC-X01-F001    | Déclaration de décès                  | PDF/A-3   |
| ID-01-F001     | Demande CNI (1ère émission)           | PDF/A-3   |
| PP-01-F001     | Demande de passeport                  | PDF/A-3   |
| FR-001         | Rapport d'incident fraude             | PDF/A-3   |

---

*Document généré par: Système de Gestion Documentaire SNISID*
*Classification: CONFIDENTIEL — Usage Interne ONI*
*Révision suivante prévue: 2026-11-25*
*Hash document: [Calculé à la validation]*
