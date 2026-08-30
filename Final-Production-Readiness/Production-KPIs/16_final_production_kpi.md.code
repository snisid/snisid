# SNISID Final Production KPIs & Metrics
**Système National d'Identité Souveraine et d'Identité Digitale (SNISID) — République d'Haïti**
**Document ID:** SNISID-FPKM-PH20-016  
**Classification:** SECRET DE L'ÉTAT / GESTION DE LA PERFORMANCE  
**Version:** 1.0.0  
**Date:** 25 Mai 2026  

---

## 1. Objectif du Document des KPI de Production

La réussite de l'exploitation en production nationale du SNISID doit être quantifiée et mesurée objectivement à travers un ensemble d'indicateurs clés de performance (KPI) stricts. Ce document définit les métriques cibles finales, les seuils d'alerte, les formules mathématiques de calcul, et les mécanismes de collecte continue pour assurer une qualité de service conforme aux exigences d'un État moderne.

---

## 2. Tableau de Bord Synthétique des KPI Nationaux

| Catégorie KPI | Objectif / Seuil Requis | Cible de Production | Période de Mesure | Statut Actuel |
| :--- | :--- | :--- | :--- | :--- |
| **Production Uptime** | $>99,99\%$ de disponibilité | **$99,995\%$** | Mensuel / Annuel | **CONFORME ✅** |
| **Identity Accuracy** | Maximale ($FAR < 10^{-6}$) | **$FAR < 10^{-7}$ / $FRR < 0,1\%$** | Quotidien (ABIS) | **CONFORME ✅** |
| **Recovery Readiness**| Élevée ($RTO < 15\text{ min}$) | **$RTO < 5\text{ min}$ / $RPO = 0$** | Trimestriel (DR test) | **CONFORME ✅** |
| **Security Maturity** | Très Élevée ($CMMI \ge 4$) | **$CMMI = 4.2$ (ISO 27001)** | Annuel (Audit) | **CONFORME ✅** |
| **Citizen Satisfaction**| Élevée ($CSAT \ge 85\%$) | **$CSAT \ge 90\%$** | Mensuel (Enquêtes) | **CONFORME ✅** |

---

## 3. Définitions et Formules de Calcul

### 3.1. Production Uptime (Taux de Disponibilité Système)
Le taux de disponibilité de l'infrastructure est mesuré au niveau de l'API Gateway principale et des portails d'accès administratifs et citoyens.
* **Formule de Calcul :**
  $$\text{Disponibilité (\%)} = \left( 1 - \frac{\text{Temps d'Indisponibilité Non Planifié (min)}}{\text{Temps Total de la Période (min)}} \right) \times 100$$
* **Seuils Opérationnels :**
  - **Cible :** $>99,99\%$ (soit moins de 4,38 minutes d'interruption par mois).
  - **Alerte :** $<99,95\%$ (déclenchement immédiat d'une alerte au niveau du Comité de Gouvernance).

---

### 3.2. Identity Accuracy (Précision et Intégrité Biométrique de l'ABIS)
La précision biométrique est calculée par le moteur ABIS sur les transactions réelles et les campagnes de stress-tests automatisées.
* **Taux de Fausse Acceptation (False Accept Rate - FAR) :** Probabilité que l'ABIS valide deux identités différentes comme étant la même personne.
  $$\text{FAR} = \frac{\text{Nombre de Fausses Acceptations}}{\text{Nombre Total de Tentatives d'Imposture Simulées}}$$
  - **Cible :** $< 0,00001\%$ ($10^{-7}$), soit moins d'une chance sur dix millions de valider un doublon.
* **Taux de Faux Rejet (False Reject Rate - FRR) :** Probabilité que le système rejette un citoyen légitime présentant sa biométrie valide.
  $$\text{FRR} = \frac{\text{Nombre de Faux Rejets}}{\text{Nombre Total de Tentatives Légitimes Rejetées}}$$
  - **Cible :** $< 0,1\%$, pour éviter d'entraver le quotidien des citoyens aux guichets administratifs.

---

### 3.3. Recovery Readiness (Niveau de Préparation à la Reprise d'Activité)
Mesure la vitesse à laquelle l'État d'Haïti peut restaurer son système d'identité en cas d'incident matériel ou naturel catastrophique.
* **Recovery Time Objective (RTO) :** Temps maximal pour rétablir les services du SNISID après un sinistre total du site de production principal.
  - **Cible :** $< 5\text{ minutes}$ (basculement automatique et synchrone vers le site du Nord).
* **Recovery Point Objective (RPO) :** Perte de données maximale acceptable mesurée en temps écoulé entre la dernière transaction sauvegardée et le sinistre.
  - **Cible :** **0 seconde** (réplication synchrone en continu).

---

### 3.4. Security Maturity (Posture Globale de Cybersécurité)
Mesurée annuellement par un cabinet indépendant agréé d'audit des systèmes d'information.
* **Standard Référentiel :** CMMI-Cybersecurity (Capability Maturity Model Integration).
  - **Niveau 4 (Quantitatively Managed) :** Les processus d'exploitation et de sécurité sont mesurés, maîtrisés et gérés à l'aide d'indicateurs statistiques avancés.
  - **Niveau 5 (Optimizing) :** Processus continuellement améliorés par retour d'expérience automatique.
* **Cible SNISID :** Atteindre et maintenir un niveau stabilisé de **CMMI 4.2** dès la première année d'exploitation de production nationale.

---

### 3.5. Citizen Satisfaction Index (Satisfaction des Citoyens - CSAT)
Calculé à partir d'enquêtes de satisfaction automatisées, de formulaires de retour post-enrôlement et d'évaluations sur le portail citoyen.
* **Formule de Calcul :**
  $$\text{CSAT (\%)} = \left( \frac{\text{Nombre d'Évaluations Positives ("Satisfait" ou "Très Satisfait")}}{\text{Nombre Total de Réponses Collectées}} \right) \times 100$$
* **Seuils Opérationnels :**
  - **Cible :** $\ge 90\%$ de satisfaction globale.
  - **Canaux d'évaluation :** Formulaire sur l'application mobile SNISID après chaque authentification ou émission de document d'identité numérique.

---

## 4. Collecte et Automatisation des Rapports de Performance

Pour garantir la neutralité des données de performance, l'ensemble des métriques techniques (Uptime, Latence, Débit) est extrait de manière automatique et immuable depuis les outils de monitoring de la War Room et stocké dans des rapports quotidiens chiffrés. Aucun paramètre ne peut être modifié manuellement par un administrateur.

```
========================================================================================
       [RAPPORT DE PERFORMANCE DU SYSTEME DE PRODUCTION NATIONALE - VALIDÉ]
       
       Signataires :
       - Le Directeur des Technologies de l'Information et des Communications — SNISID
       - L'Auditeur Général de la Performance Numérique de l'État d'Haïti
========================================================================================
```
