# SNISID — Matrice RACI Nationale Complète
## Système National d'Identité et de Services d'Identité Digitale

---

| Métadonnée | Valeur |
|---|---|
| **Document ID** | SNISID-GOV-RACI-001 |
| **Version** | 1.0.0 |
| **Statut** | APPROUVÉ — EN VIGUEUR |
| **Date de création** | 2026-05-25 |
| **Date de révision** | 2026-11-25 |
| **Classification** | GOUVERNANCE / USAGE INTERNE |
| **Propriétaire** | AND — Directeur Juridique & Conformité |
| **Révisé par** | IGB + NDPA |
| **Approuvé par** | DG AND / CNN |
| **Référence** | SNISID-GOV-ORG-001 / SNISID-GOV-WRK-001 |

---

> **IMPORTANT** : Cette matrice RACI est le document de référence pour toute question de responsabilité dans le SNISID. En cas de contradiction avec un autre document, ce RACI prévaut pour les questions de responsabilité. Toute modification requiert l'approbation du DG AND et notification au CNN.

---

## TABLE DES MATIÈRES

1. [Légende et Guide d'Utilisation](#1-légende-et-guide-dutilisation)
2. [Acteurs Institutionnels — Référentiel Complet](#2-acteurs-institutionnels--référentiel-complet)
3. [Matrice RACI — Gouvernance et Politique](#3-matrice-raci--gouvernance-et-politique)
4. [Matrice RACI — Gestion de l'Identité](#4-matrice-raci--gestion-de-lidentité)
5. [Matrice RACI — Données et Protection](#5-matrice-raci--données-et-protection)
6. [Matrice RACI — Sécurité et Infrastructure](#6-matrice-raci--sécurité-et-infrastructure)
7. [Matrice RACI — Opérations et Services](#7-matrice-raci--opérations-et-services)
8. [Matrice RACI — Conformité Légale et Audit](#8-matrice-raci--conformité-légale-et-audit)
9. [Matrice RACI — Incidents et Crises](#9-matrice-raci--incidents-et-crises)
10. [Matrice RACI — Partenariats et International](#10-matrice-raci--partenariats-et-international)
11. [Matrice d'Escalade pour Décisions Bloquées](#11-matrice-descalade-pour-décisions-bloquées)
12. [Processus de Révision Trimestrielle](#12-processus-de-révision-trimestrielle)

---

## 1. Légende et Guide d'Utilisation

### 1.1 Définition des Rôles RACI

| Code | Rôle | Définition | Obligations |
|---|---|---|---|
| **R** | **Responsible** (Responsable) | L'acteur qui **exécute** la tâche. Il est en charge du travail effectif. | Exécuter, documenter, rendre compte |
| **A** | **Accountable** (Redevable) | L'acteur qui **répond** de la tâche devant l'autorité supérieure. Il ne peut y en avoir qu'un seul par tâche. | Signer, approuver, rendre des comptes |
| **C** | **Consulted** (Consulté) | L'acteur dont l'**avis est requis** avant toute décision ou action. Communication bidirectionnelle. | Donner un avis formel, être disponible |
| **I** | **Informed** (Informé) | L'acteur qui **doit être notifié** des résultats. Communication unidirectionnelle. | Accusé de réception, action si nécessaire |
| **V** | **Veto** (Droit de veto) | L'acteur dont l'**opposition bloque** la décision. Spécifique SNISID. | Exercer le veto par écrit avec justification |
| **—** | **Non impliqué** | L'acteur n'est pas concerné par cette activité. | Aucune |

### 1.2 Règles d'Utilisation de la Matrice

| Règle | Description |
|---|---|
| **Règle du A unique** | Chaque ligne ne peut avoir qu'un seul **A** (Accountable). S'il y en a plusieurs, escalader vers AND pour clarification. |
| **Règle du R minimum** | Chaque ligne doit avoir au moins un **R** (Responsible). |
| **Règle du conflit RACI** | En cas de désaccord sur la responsabilité, AND arbitre dans les 48h. |
| **Règle du veto NDPA** | Le veto NDPA est suspensif et immédiat. Seul le CNN peut lever un veto NDPA. |
| **Règle de l'escalade** | Si le A ne peut décider dans les délais prévus, escalade automatique à AND. |
| **Règle de documentation** | Toute exécution de tâche majeure doit être documentée avec référence RACI. |

### 1.3 Codes Couleur de Criticité

| Criticité | Code | Description |
|---|---|---|
| 🔴 **CRITIQUE** | [C] | Activité critique nationale — toute défaillance a impact direct citoyen |
| 🟠 **ÉLEVÉE** | [E] | Activité à fort impact — délais stricts, supervision AND |
| 🟡 **MODÉRÉE** | [M] | Activité importante — processus standard |
| 🟢 **NORMALE** | [N] | Activité de routine — délégation possible |

---

## 2. Acteurs Institutionnels — Référentiel Complet

### 2.1 Corps de Gouvernance SNISID

| Code | Institution | Rôle Principal |
|---|---|---|
| **CNN** | Conseil National du Numérique | Autorité politique suprême |
| **AND** | Autorité Nationale Numérique | Direction exécutive SNISID |
| **IGB** | Identity Governance Board | Gouvernance identitaire |
| **NDPA** | National Data Protection Authority | Protection données personnelles |
| **SOC** | National SOC / CERT-HT | Sécurité opérationnelle |
| **CCB** | Change Control Board | Contrôle des changements |
| **XROAD** | Comité Interministériel X-Road | Interopérabilité nationale |
| **PKI** | National PKI Authority | Infrastructure à clé publique |
| **ETH** | Comité Éthique IA & Données | Éthique et droits fondamentaux |

### 2.2 Agences Gouvernementales

| Code | Institution | Rôle dans SNISID |
|---|---|---|
| **OJRNH** | Office de l'État Civil National | Registre état civil, source primaire identité |
| **MJP** | Ministère de la Justice | Cadre légal, judiciaire |
| **MI** | Ministère de l'Intérieur | Sécurité intérieure, territorialité |
| **MEF** | Ministère de l'Économie et Finances | Financement, budget |
| **MSPP** | Ministère de la Santé | Données de santé liées à l'identité |
| **MENFP** | Ministère de l'Éducation | Identité scolaire et académique |
| **PNH** | Police Nationale d'Haïti | Application de la loi, biométrie judiciaire |
| **BRH** | Banque de la République d'Haïti | KYC bancaire, identité financière |
| **AGS** | Administration Générale des Douanes | Identité aux frontières |
| **DGI** | Direction Générale des Impôts | NIF lié à l'identité |
| **CONATEL** | Conseil National des Télécommunications | Infrastructure télécom |
| **MEH** | Ministère de l'Environnement | Registres fonciers liés à l'identité |
| **MPCE** | Ministère Planification | Statistiques nationales |
| **MAST** | Ministère Affaires Sociales | Protection sociale, vulnérables |
| **MAEC** | Ministère Affaires Étrangères | Diaspora, consulats |

---

## 3. Matrice RACI — Gouvernance et Politique

### Colonne Headers (voir tableau ci-dessous)
CNN=Conseil National du Numérique | AND=Autorité Nationale Numérique | IGB=Identity Governance Board | NDPA=National Data Protection Authority | ETH=Comité Éthique | MJP=Min. Justice | MI=Min. Intérieur

| # | Activité | Criticité | CNN | AND | IGB | NDPA | ETH | MJP | MI |
|---|---|---|---|---|---|---|---|---|---|
| G-01 | Adoption de la politique nationale d'identité digitale | 🔴[C] | **A** | R | C | C | C | C | I |
| G-02 | Révision annuelle de la stratégie SNISID | 🟠[E] | A | **R** | C | C | C | I | I |
| G-03 | Approbation du budget pluriannuel SNISID | 🔴[C] | **A** | R | I | I | — | I | I |
| G-04 | Nomination / révocation DG AND | 🔴[C] | **A** | — | I | I | — | — | — |
| G-05 | Adoption standards de gouvernance SNISID | 🟠[E] | A | **R** | C | C | C | I | I |
| G-06 | Publication rapport annuel public SNISID | 🟡[M] | A | **R** | C | C | C | I | I |
| G-07 | Activation État d'Urgence Numérique | 🔴[C] | **A** | R | I | I | I | C | C |
| G-08 | Définition des seuils de sanctions SNISID | 🟠[E] | A | **R** | C | C | I | C | I |
| G-09 | Validation partenariats internationaux stratégiques | 🔴[C] | **A** | R | I | C | I | C | C |
| G-10 | Révision de l'organigramme de gouvernance | 🟠[E] | **A** | R | C | C | C | I | I |

---

## 4. Matrice RACI — Gestion de l'Identité

| # | Activité | Criticité | CNN | AND | IGB | NDPA | SOC | OJRNH | PNH | MEF | MSPP | MENFP | MAST | MAEC |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| ID-01 | Définition des standards de création d'identité nationale | 🔴[C] | I | A | **R** | C | I | C | I | I | I | I | I | I |
| ID-02 | Enregistrement nouveau citoyen (naissance) | 🔴[C] | — | A | C | I | I | **R** | — | — | R | — | — | — |
| ID-03 | Enregistrement citoyen adulte non enregistré | 🔴[C] | — | A | C | I | I | **R** | C | — | — | — | C | — |
| ID-04 | Enregistrement citoyen de la diaspora | 🟠[E] | — | A | C | I | I | **R** | — | — | — | — | — | R |
| ID-05 | Capture et validation données biométriques | 🔴[C] | — | A | **R** | C | C | R | I | — | — | — | — | — |
| ID-06 | Émission identité digitale (NIN / carte numérique) | 🔴[C] | — | **A** | R | I | I | R | — | — | — | — | — | — |
| ID-07 | Mise à jour des données identitaires (changement état civil) | 🟠[E] | — | A | **R** | C | I | R | — | — | — | — | — | — |
| ID-08 | Révocation identité (décès) | 🟠[E] | — | A | **R** | I | I | R | — | — | R | — | — | — |
| ID-09 | Révocation identité (fraude prouvée) | 🔴[C] | I | **A** | R | C | R | C | C | — | — | — | — | — |
| ID-10 | Suspension temporaire d'identité digitale | 🟠[E] | I | **A** | R | C | R | I | C | — | — | — | — | — |
| ID-11 | Procédure de réclamation / contestation identitaire | 🟡[M] | — | A | **R** | C | — | R | — | — | — | — | — | — |
| ID-12 | Gestion des identités des réfugiés et apatrides | 🟠[E] | I | **A** | R | C | I | C | I | — | C | — | R | — |
| ID-13 | Accréditation d'un centre d'enregistrement | 🟡[M] | — | A | **R** | C | C | I | I | — | — | — | — | — |
| ID-14 | Révocation d'un centre d'enregistrement | 🟠[E] | I | **A** | R | C | I | I | — | — | — | — | — | — |
| ID-15 | Accréditation d'un opérateur biométrique | 🟠[E] | — | A | **R** | C | C | — | I | — | — | — | — | — |
| ID-16 | Audit des registres d'identité nationaux | 🟠[E] | I | A | **R** | C | I | R | — | — | — | — | — | — |
| ID-17 | Définition des procédures identité pour mineurs | 🟡[M] | — | A | **R** | C | — | R | — | — | R | R | C | — |
| ID-18 | Définition procédures identité handicapés / vulnérables | 🟡[M] | — | A | **R** | C | — | R | — | — | R | — | R | — |

---

## 5. Matrice RACI — Données et Protection

| # | Activité | Criticité | CNN | AND | IGB | NDPA | ETH | SOC | CCB | OJRNH | MEF | MSPP |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| D-01 | Définition de la politique de classification des données | 🔴[C] | I | A | C | **R** | C | C | I | I | I | I |
| D-02 | Autorisation d'un nouveau traitement de données personnelles | 🔴[C] | — | C | C | **A** | C | I | I | — | — | — |
| D-03 | Veto sur un traitement de données (NDPA) | 🔴[C] | V | C | C | **A/R** | C | I | I | I | I | I |
| D-04 | Levée du veto NDPA | 🔴[C] | **A/R** | C | — | — | C | — | — | — | — | — |
| D-05 | Gestion des droits d'accès aux données | 🟠[E] | — | A | C | **R** | I | C | C | I | I | I |
| D-06 | Réponse aux demandes d'accès citoyen (DSAR) | 🟡[M] | — | A | C | **R** | I | — | — | R | R | R |
| D-07 | Traitement des plaintes citoyennes données | 🟡[M] | — | I | — | **A/R** | C | — | — | I | I | I |
| D-08 | Notification de violation de données (Data Breach) | 🔴[C] | I | **A** | I | R | I | R | I | I | I | I |
| D-09 | Procédure de suppression de données (droit à l'oubli) | 🟡[M] | — | A | C | **R** | I | C | C | R | R | R |
| D-10 | Audit de conformité protection des données | 🟠[E] | I | A | C | **R** | I | C | — | R | R | R |
| D-11 | Évaluation d'impact (DPIA / PIA) sur nouveau système | 🟠[E] | — | A | C | **R** | C | C | C | — | — | — |
| D-12 | Politique de rétention et archivage des données | 🟠[E] | I | A | C | **R** | I | C | C | I | I | I |
| D-13 | Destruction sécurisée des données biométriques | 🔴[C] | I | **A** | R | R | I | C | C | — | — | — |
| D-14 | Transfert international de données personnelles | 🔴[C] | I | **A** | I | R | C | I | — | — | — | — |
| D-15 | Publication du registre des traitements | 🟡[M] | — | A | I | **R** | I | — | — | — | — | — |

---

## 6. Matrice RACI — Sécurité et Infrastructure

| # | Activité | Criticité | CNN | AND | SOC | CCB | PKI | NDPA | IGB | CONATEL |
|---|---|---|---|---|---|---|---|---|---|---|
| S-01 | Définition de la politique de sécurité SNISID | 🔴[C] | I | **A** | R | C | C | C | C | I |
| S-02 | Gestion des clés PKI de la CA Racine | 🔴[C] | I | **A** | C | I | R | I | — | — |
| S-03 | Émission de certificats numériques citoyens | 🔴[C] | — | A | I | I | **R** | I | C | — |
| S-04 | Révocation de certificats numériques | 🟠[E] | — | A | C | I | **R** | C | C | — |
| S-05 | Cérémonie Root CA (génération clé racine) | 🔴[C] | I | **A** | C | C | R | I | — | — |
| S-06 | Gestion des HSM (Hardware Security Modules) | 🟠[E] | — | A | C | C | **R** | — | — | — |
| S-07 | Tests de pénétration (pentest) du SNISID | 🟠[E] | — | **A** | R | C | C | C | — | — |
| S-08 | Gestion des vulnérabilités (patch management) | 🟠[E] | — | A | **R** | R | C | — | — | — |
| S-09 | Surveillance temps réel (monitoring H24) | 🔴[C] | — | A | **R** | I | I | — | — | — |
| S-10 | Gestion des identités et accès (IAM) système | 🟠[E] | — | **A** | R | C | C | C | C | — |
| S-11 | Plan de reprise d'activité (DR/BCP) | 🟠[E] | I | **A** | R | C | C | I | I | C |
| S-12 | Test DRP / BCP annuel | 🟡[M] | I | **A** | R | C | C | — | — | — |
| S-13 | Gestion du DataCenter Principal | 🔴[C] | — | **A** | C | C | C | — | — | C |
| S-14 | Gestion du DataCenter Secondaire | 🔴[C] | — | **A** | C | C | C | — | — | C |
| S-15 | Audit de sécurité annuel ISO 27001 | 🟠[E] | I | **A** | R | C | C | C | — | — |
| S-16 | Gestion du réseau SNISID (backbone) | 🟠[E] | — | **A** | C | R | — | — | — | C |
| S-17 | Accès d'urgence aux systèmes (break glass) | 🔴[C] | I | **A** | R | C | C | — | — | — |

---

## 7. Matrice RACI — Opérations et Services

| # | Activité | Criticité | CNN | AND | IGB | SOC | CCB | OJRNH | MI | MEF | MSPP | MENFP | PNH | BRH | DGI |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| O-01 | Déploiement de nouveaux centres d'enregistrement | 🟠[E] | I | **A** | R | I | C | R | C | C | — | — | — | — | — |
| O-02 | Formation des agents d'enregistrement | 🟡[M] | — | A | **R** | I | — | R | — | — | — | — | — | — | — |
| O-03 | Maintenance des kiosques d'enregistrement | 🟡[M] | — | **A** | I | C | R | R | — | — | — | — | — | — | — |
| O-04 | Gestion des files d'attente et SLA service citoyen | 🟡[M] | — | **A** | C | I | — | R | — | — | — | — | — | — | — |
| O-05 | Déploiement d'une nouvelle fonctionnalité X-Road | 🟠[E] | — | A | C | C | **R** | I | C | C | C | C | I | I | I |
| O-06 | Onboarding d'un nouveau ministère sur X-Road | 🟡[M] | — | A | I | I | **R** | — | C | C | C | C | I | I | I |
| O-07 | Gestion des APIs d'accès tiers | 🟡[M] | — | **A** | C | C | R | — | — | — | — | — | — | — | — |
| O-08 | Supervision du SLA des services SNISID | 🟠[E] | — | **A** | I | R | I | I | I | I | I | I | I | I | I |
| O-09 | Reporting opérationnel mensuel | 🟡[M] | I | **A** | R | R | R | R | — | — | — | — | — | — | — |
| O-10 | Gestion des demandes de support niveau 1 | 🟢[N] | — | A | I | I | — | **R** | — | — | — | — | — | — | — |
| O-11 | Gestion des demandes de support niveau 2 | 🟡[M] | — | **A** | C | I | — | R | — | — | — | — | — | — | — |
| O-12 | Gestion des demandes de support niveau 3 | 🟠[E] | — | **A** | C | C | C | R | — | — | — | — | — | — | — |
| O-13 | Intégration KYC bancaire (service BRH) | 🟠[E] | — | **A** | C | C | C | — | — | — | — | — | — | R | — |
| O-14 | Intégration service NIF fiscal (DGI) | 🟡[M] | — | **A** | C | I | C | — | — | — | — | — | — | — | R |
| O-15 | Intégration service PNH (biométrie judiciaire) | 🔴[C] | I | **A** | C | C | C | — | C | — | — | — | R | — | — |

---

## 8. Matrice RACI — Conformité Légale et Audit

| # | Activité | Criticité | CNN | AND | IGB | NDPA | ETH | MJP | MEF | CCB |
|---|---|---|---|---|---|---|---|---|---|---|
| L-01 | Audit de conformité ISO 27001 | 🟠[E] | I | **A** | C | C | — | — | — | C |
| L-02 | Audit de conformité ISO 22301 (continuité) | 🟡[M] | I | **A** | I | I | — | — | — | C |
| L-03 | Évaluation conformité NIST CSF 2.0 | 🟡[M] | I | **A** | C | C | — | — | — | C |
| L-04 | Conformité Convention 108+ (Conseil de l'Europe) | 🟠[E] | I | A | — | **R** | C | C | — | — |
| L-05 | Conformité ICAO 9303 (documents voyage biométriques) | 🟠[E] | I | **A** | R | C | — | — | — | I |
| L-06 | Rapport de conformité Convention de Budapest | 🟡[M] | I | A | — | **R** | — | C | — | — |
| L-07 | Revue légale de tout nouveau système / processus | 🟠[E] | — | **A** | C | C | C | R | — | C |
| L-08 | Réponse aux injonctions judiciaires | 🔴[C] | I | **A** | I | C | — | R | — | — |
| L-09 | Préparation du rapport annuel au Parlement | 🟠[E] | **A** | R | C | C | C | I | I | — |
| L-10 | Certification WebTrust PKI | 🟡[M] | — | **A** | — | I | — | — | — | I |
| L-11 | Coordination avec Cour Supérieure des Comptes | 🟡[M] | I | **A** | — | — | — | I | C | — |
| L-12 | Revue par le Parlement (Commission Numérique) | 🟠[E] | **A** | R | I | I | I | C | I | — |
| L-13 | Mise à jour du registre des conformités | 🟢[N] | — | A | I | **R** | — | — | — | — |
| L-14 | Traitement des litiges avec fournisseurs | 🟡[M] | — | **A** | — | I | — | R | C | — |
| L-15 | Gestion de la propriété intellectuelle SNISID | 🟡[M] | — | **A** | — | — | — | R | — | — |

---

## 9. Matrice RACI — Incidents et Crises

| # | Activité | Criticité | CNN | AND | SOC | CCB | PKI | NDPA | IGB | MI | MJP | PNH |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| INC-01 | Détection d'incident de sécurité (P5-P4) | 🟡[M] | — | I | **A/R** | I | I | — | — | — | — | — |
| INC-02 | Gestion incident P3 (alerte) | 🟠[E] | — | **A** | R | C | C | I | I | — | — | — |
| INC-03 | Gestion incident P2 (grave) | 🔴[C] | I | **A** | R | C | C | C | I | I | — | — |
| INC-04 | Gestion incident P1 (critique) | 🔴[C] | **A** | R | R | C | C | C | I | C | C | C |
| INC-05 | Gestion incident P0 (guerre cyber) | 🔴[C] | **A** | R | R | C | C | C | I | C | C | C |
| INC-06 | Notification de violation de données aux citoyens | 🔴[C] | I | **A** | R | — | — | R | — | — | — | — |
| INC-07 | Investigation forensique post-incident | 🟠[E] | I | **A** | R | C | C | C | — | C | C | C |
| INC-08 | Coordination avec INTERPOL / autorités étrangères | 🟠[E] | I | **A** | C | — | — | I | — | C | R | R |
| INC-09 | Activation Plan de Continuité (BCP) | 🔴[C] | I | **A** | R | C | C | I | I | — | — | — |
| INC-10 | Exercice de simulation d'incident annuel | 🟡[M] | I | **A** | R | C | C | C | C | — | — | — |
| INC-11 | Rapport post-incident (PIR) | 🟡[M] | I | A | **R** | C | C | C | I | — | — | — |
| INC-12 | Révocation d'urgence de l'identité compromise | 🔴[C] | — | **A** | C | I | R | C | R | — | — | — |
| INC-13 | Notification au CNN (incident P1/P0) | 🔴[C] | I | **R** | R | — | — | — | — | — | — | — |
| INC-14 | Communication publique de crise | 🔴[C] | **A** | R | I | — | — | I | — | C | I | — |
| INC-15 | Gestion crise catastrophe naturelle (séisme, ouragan) | 🔴[C] | **A** | R | R | C | C | I | I | R | — | C |

---

## 10. Matrice RACI — Partenariats et International

| # | Activité | Criticité | CNN | AND | IGB | NDPA | ETH | MJP | MAEC | MEF |
|---|---|---|---|---|---|---|---|---|---|---|
| P-01 | Négociation accord bilatéral identité numérique | 🔴[C] | **A** | R | C | C | I | C | R | I |
| P-02 | Ratification accord international sur données | 🔴[C] | **A** | R | I | C | C | C | R | I |
| P-03 | Coopération technique avec Estonie (X-Road) | 🟠[E] | I | **A** | C | C | I | — | C | — |
| P-04 | Coopération ICAO (passeports biométriques) | 🟠[E] | I | **A** | R | C | I | I | C | — |
| P-05 | Partenariat avec organismes de financement (BID, BM) | 🟠[E] | **A** | R | I | I | I | — | I | R |
| P-06 | Accréditation d'un fournisseur de services tiers | 🟡[M] | — | **A** | R | C | C | — | — | — |
| P-07 | Révocation d'accréditation fournisseur | 🟠[E] | I | **A** | R | C | I | C | — | — |
| P-08 | Rapport à l'ONU (ODD 16.9) | 🟡[M] | I | **A** | C | I | C | — | R | — |
| P-09 | Participation aux forums internationaux (ID4D, OGP) | 🟡[M] | I | **A** | C | I | I | — | C | — |
| P-10 | Coopération avec CARICOM sur identité | 🟠[E] | I | **A** | C | C | I | I | R | — |

---

## 11. Matrice d'Escalade pour Décisions Bloquées

### 11.1 Procédure d'Escalade Standard

```
NIVEAU 0 (Organe N3)
│   Délai : selon procédure de l'organe
│   Si bloqué 48h → NIVEAU 1
│
▼
NIVEAU 1 (AND — Médiation)
│   AND arbitre entre les organes N3 en conflit
│   Délai : 5 jours ouvrables
│   Si bloqué ou insuffisant → NIVEAU 2
│
▼
NIVEAU 2 (AND — Décision)
│   DG AND prend la décision de dernier ressort N3
│   Délai : 5 jours ouvrables supplémentaires
│   Si impact stratégique national → NIVEAU 3
│
▼
NIVEAU 3 (CNN — Arbitrage Final)
│   CNN saisi en session ordinaire ou extraordinaire
│   Délai : 21 jours (ordinaire) / 72h (extraordinaire)
│   Décision finale et souveraine — non contestable
│
▼
NIVEAU ULTIME (Parlement)
    Si question constitutionnelle ou légale fondamentale
    Saisine de la Commission Parlementaire compétente
```

### 11.2 Matrice de Résolution des Conflits Inter-organes

| Conflit Entre | Nature du Conflit | Arbitre | Délai |
|---|---|---|---|
| IGB vs NDPA | Tension identité / protection données | AND DG | 5 j |
| SOC vs CCB | Urgence sécurité vs procédure changement | CISO AND | 24h |
| CCB vs XROAD | Changement X-Road bloqué | AND CTO | 5 j |
| NDPA vs AND | Veto NDPA sur décision AND | CNN | 30 j |
| ETH vs AND | Moratoire éthique sur fonctionnalité | CNN | 30 j |
| PKI vs CCB | Procédure PKI vs changement demandé | CISO AND | 48h |
| IGB vs ETH | Politique identitaire vs éthique | AND DG | 10 j |
| Ministère vs AND | Refus connexion X-Road | Comité XROAD | 30 j |
| AND vs Fournisseur | Litige contractuel | MJP / Arbitrage | Selon contrat |
| Tout organe vs CNN | N/A — CNN est autorité suprême | Parlement | — |

### 11.3 Situations de Blocage Critique

| Situation | Procédure d'Exception | Autorité |
|---|---|---|
| **Veto NDPA bloquant une opération critique** | CNN vote en session extraordinaire (72h) | CNN 3/4 |
| **Absence de quorum CCB pour changement urgent P1** | DG AND approuve provisoirement, CCB ratifie sous 72h | DG AND |
| **IGB ne peut statuer sur cas exceptionnel** | AND décide, IGB ratifie rétrospectivement | DG AND |
| **CNN en impossibilité de se réunir (catastrophe)** | Président de la République décide, CNN ratifie | Président |
| **Conflit total paralysant le SNISID** | Plan CRISNUM-01 activé, AND assume tous pouvoirs | DG AND |

---

## 12. Processus de Révision Trimestrielle

### 12.1 Calendrier de Révision

| Révision | Période | Responsable | Livrables |
|---|---|---|---|
| **Q1 — Révision Ordinaire** | Janvier-Février | AND + tous organes N3 | RACI révisé, delta report |
| **Q2 — Révision Ordinaire** | Avril-Mai | AND + tous organes N3 | RACI révisé, delta report |
| **Q3 — Révision Approfondie** | Juillet-Septembre | AND + CNN | RACI v.x+1 soumis CNN |
| **Q4 — Révision Annuelle** | Octobre-Décembre | AND + CNN + Parlement | RACI vx+1 approuvé, rapport annuel |

### 12.2 Procédure de Révision

**Étape 1 : Collecte des Propositions (J-45 avant publication)**
- Chaque organe N3 soumet ses propositions de modification au format standardisé
- Formulaire : SNISID-RACI-MOD-YYYY-NNN
- Justification obligatoire pour chaque modification proposée

**Étape 2 : Analyse de l'Impact (J-30)**
- AND analyse les propositions et leur impact croisé
- Consultation des organes concernés
- Rapport d'impact préliminaire

**Étape 3 : Consultation (J-21)**
- Session de consultation avec tous les présidents d'organes N3
- NDPA et Comité Éthique donnent avis formel
- Réponses documentées

**Étape 4 : Rédaction (J-14)**
- AND rédige la version révisée du RACI
- Revue juridique par Directeur Juridique AND
- Version draft soumise à IGB pour validation

**Étape 5 : Approbation (J-7)**
- Revue finale par DG AND
- Pour révisions majeures (>10% des lignes) : soumission au CNN
- Signature et publication

**Étape 6 : Communication (J-0)**
- Distribution officielle à tous les organes et agences
- Formation si nouvelles responsabilités majeures
- Mise à jour du registre de gouvernance

### 12.3 Critères de Révision Urgente (Hors Calendrier)

Une révision urgente du RACI peut être déclenchée si :
- Un incident majeur a révélé une lacune de responsabilité
- Un nouveau texte légal modifie les attributions d'un organe
- CNN émet une résolution modifiant les responsabilités
- NDPA identifie un conflit de responsabilité en matière de données
- Un nouvel organe est créé ou supprimé

**Délai d'exécution révision urgente : 15 jours ouvrables maximum**

---

## Bloc de Signature

**APPROUVÉ PAR LE DIRECTEUR GÉNÉRAL AND**

```
Nom            : ___________________________
Qualité        : Directeur Général, Autorité Nationale Numérique
Signature      : ___________________________
Date           : ___________________________
Cachet AND     : [CACHET AND]
```

**VALIDÉ PAR LE DIRECTEUR JURIDIQUE AND**

```
Nom            : ___________________________
Qualité        : Directeur Juridique & Conformité, AND
Signature      : ___________________________
Date           : ___________________________
```

**NOTIFICATION CNN — Résolution d'Information**

```
Résolution     : CNN-2026-002
Date           : ___________________________
Note           : Approuvé par CNN en session ordinaire Q2 2026
```

---

**HISTORIQUE DES RÉVISIONS**

| Version | Date | Modifications | Approuvé par |
|---|---|---|---|
| 0.1 | 2026-03-15 | Première ébauche — gouvernance et identité | AND Juridique |
| 0.5 | 2026-04-20 | Ajout sécurité, opérations, conformité | DGA AND |
| 0.8 | 2026-05-10 | Consultation organes N3, révisions | DG AND |
| 1.0 | 2026-05-25 | Version finale complète | DG AND / CNN |

---

*Document SNISID-GOV-RACI-001 v1.0.0 — Propriété de l'Autorité Nationale Numérique de la République d'Haïti*

*© 2026 République d'Haïti — SNISID Phase 0 — Gouvernance Nationale*
