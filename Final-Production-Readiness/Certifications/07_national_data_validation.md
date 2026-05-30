# National Data Validation Program
**Système National d'Identité Souveraine et d'Identité Digitale (SNISID) — République d'Haïti**
**Document ID:** SNISID-NDVP-PH20-007  
**Classification:** SECRET DE L'ÉTAT / FIABILITÉ DES DONNÉES  
**Version:** 1.0.0  
**Date:** 25 Mai 2026  

---

## 1. Objectif du Programme de Validation des Données

La valeur d'un système d'identité nationale dépend directement de la qualité, de l'exactitude et de la pureté des données qu'il héberge. Le **National Data Validation Program (NDVP)** applique des filtres d'intégrité biométrique et biographique stricts pour éliminer définitivement les doublons, corriger les incohérences historiques d'état civil, et garantir que chaque citoyen d'Haïti est représenté par une identité unique, inaltérable et vérifiable.

---

## 2. Les Quatre Piliers de la Qualité des Données SNISID

```
========================================================================================
                         CADRE DE FIABILITÉ DES DONNÉES SNISID
========================================================================================
[1] IDENTITY UNIQUENESS  ===> Déduplication exhaustive de la population via l'ABIS.
[2] DUPLICATE RATE       ===> Maintien d'un taux de doublons à un niveau proche de zéro.
[3] DATA CONSISTENCY     ===> Validation syntaxique, sémantique et historique d'état civil.
[4] BIOMETRIC INTEGRITY  ===> Conformité aux normes internationales NIST/ISO pour la biométrie.
========================================================================================
```

---

## 3. Protocoles Techniques de Validation

### 3.1. Identity Uniqueness (Garantie d'Unicité)
* **Méthodologie :**
  - Chaque nouvel enrôlement est soumis à une comparaison de type **1:N (un contre tous)** contre l'intégralité du registre national par le moteur ABIS (Automated Biometric Identification System).
  - La comparaison croise simultanément la biométrie faciale (Face Match), les empreintes digitales (10 doigts) et, le cas échéant, la reconnaissance de l'iris.
* **Seuils d'Acceptation :**
  - Score de correspondance d'empreintes digitales basé sur le standard **NIST MINEX III** certifié.
  - Seuil de décision fixé à un taux d'erreur de faux positif égal à **1 sur 10 000 000** (FAR = $10^{-7}$).
  - Tout dossier présentant un score de correspondance supérieur à ce seuil est mis en quarantaine pour arbitrage manuel par une équipe d'experts légistes d'Haïti.

---

### 3.2. Duplicate Rate Minimization (Contrôle et Résolution des Doublons)
* **Méthodologie :**
  - Nettoyage en amont de la base historique de l'ONI (contenant des décennies d'enregistrements papier ou semi-numériques parfois fragmentaires).
  - Utilisation d'algorithmes de comparaison textuelle floue (**Jaro-Winkler** et **Levenshtein**) ajustés pour les spécificités patronymiques haïtiennes (par exemple, inversion fréquente des noms et prénoms, orthographes variables de "Jean", omissions de traits d'union).
* **Indicateurs de Performance (KPI de Pureté) :**
  - **Taux de Doublons Résiduels après Traitement :** < 0,001%.
  - **Dossiers en Quarantaine d'Arbitrage :** Tous résolus et arbitrés par les agents ONI assermentés avant le GoLive de production.

---

### 3.3. Data Consistency (Cohérence et Structuration des Données)
* **Méthodologie :**
  - Validation syntaxique stricte des formulaires d'état civil au niveau des terminaux d'enrôlement (empêchant l'entrée de caractères invalides ou de dates incohérentes comme des naissances futures ou antérieures à 120 ans).
  - Liens logiques d'état civil : Validation parentale automatique (les enfants enregistrés doivent être liés à des parents valides déjà identifiés dans le SNISID).
* **Règles d'Intégrité de Base de Données (PostgreSQL / Cassandra) :**
  - Schémas JSON validant en temps réel la structure de chaque profil citoyen.
  - Contraintes d'unicité sur le numéro de registre de naissance d'origine.

---

### 3.4. Biometric Integrity (Conformité de la Biométrie aux Standards Internationaux)
* **Méthodologie :**
  - **Qualité Faciale :** Validation automatique en temps réel lors de la capture d'image selon la norme **ISO/IEC 19794-5 (ICAO Frontal Face compliance)**. Vérification de l'éclairage, du centrage, des yeux ouverts, de la neutralité de l'expression faciale, et de l'absence de reflets ou d'occultations.
  - **Qualité Empreintes Digitales :** Évaluation instantanée via l'algorithme **NFIQ 2.0** (NIST Fingerprint Image Quality). Toute empreinte digitale de qualité inférieure à un score de 60 sur 100 est automatiquement rejetée lors de l'enrôlement, forçant l'opérateur à refaire la capture ou à justifier une exception médicale de manière formelle.

---

## 4. Résultats Statistiques de la Campagne de Purification

Une campagne globale d'audit de la base de données de pré-production (contenant 6,2 millions de profils d'enrôlement historiques) a été menée du 15 Janvier 2026 au 20 Mai 2026 :

| Paramètre de Qualité | Valeur Initiale (Janvier 2026) | Valeur Finale (Mai 2026) | Seuil de Tolérance GoLive | Statut de Validation |
| :--- | :--- | :--- | :--- | :--- |
| **Profils orphelins (Sans biométrie)**| 4,2% | **0,00%** (Tous régularisés) | 0,05% | **VALIDÉ** |
| **Taux de doublons identifiés (ABIS)**| 1,8% | **0,00%** (Fusionnés/Supprimés)| 0,01% | **VALIDÉ** |
| **Incohérences de dates d'état civil** | 2,5% | **0,02%** (Corrigés par greffier) | 0,05% | **VALIDÉ** |
| **Conformité ISO Faciale (19794-5)** | 78,4% | **99,8%** | 99,0% | **VALIDÉ** |
| **Qualité Empreinte NFIQ 2.0 (>=60)**  | 81,2% | **99,2%** | 98,0% | **VALIDÉ** |

---

## 5. Conclusion de Validation de la Qualité des Données

La base de données nationale du SNISID est déclarée **CONFORME, PURIFIÉE ET EXTRÊMEMENT FIABLE**. L'unicité de chaque citoyen au sein de l'État haïtien est désormais garantie par des preuves cryptographiques et biométriques indiscutables, interdisant toute forme de fraude à l'identité ou de double inscription électorale.

```
[APPROBATION TECHNIQUE]
DIRECTEUR DE LA SÉCURITÉ ET DE LA FIABILITÉ DES DONNÉES BIOMÉTRIQUES — SNISID
```
