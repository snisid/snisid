# Final Observability & War Room Protocol
**Système National d'Identité Souveraine et d'Identité Digitale (SNISID) — République d'Haïti**
**Document ID:** SNISID-FOWR-PH20-013  
**Classification:** SECRET DE L'ÉTAT / GOLIVE NATIONAL OPERATIONAL  
**Version:** 1.0.0  
**Date:** 25 Mai 2026  

---

## 1. Rôle et Objectif de la War Room Nationale

Le **GoLive National** d'une infrastructure de l'envergure du SNISID est une opération délicate qui requiert une supervision continue et ultra-précise. La **National War Room** est une cellule opérationnelle temporaire hautement stratégique, activée 15 jours avant le lancement et maintenue au moins 30 jours après celui-ci (Phase de Stabilisation et Transition).

La War Room réunit physiquement et virtuellement les meilleurs experts techniques, spécialistes sécurité, directeurs de la communication, et officiers de liaison ministériels sous un commandement unique pour piloter le déploiement minute par minute.

---

## 2. Organisation de la War Room Nationale

```
========================================================================================
                          ORGANISATION DE LA WAR ROOM SNISID
========================================================================================
                   [ DIRECTEUR DU GOLIVE NATIONAL (COMMANDER) ]
                                       |
       +-------------------------------+-------------------------------+
       |                               |                               |
       v                               v                               v
[ CELLULE TELEMETRY ]        [ CELLULE SECURITY ALERTS ]     [ CELLULE DEPLOYMENT & FIELD ]
  - Production Health          - Surveillance SOC/WAF          - Suivi des Départements
  - Latence API & ABIS         - Tentatives de sabotage        - Enrôlements mobiles
  - Taux d'Erreur              - Analyse EDR & Intégrité       - Incidents Bureaux ONI
========================================================================================
```

---

## 3. Les Quatre Domaines de Supervision de la War Room

### 3.1. Production Health (Supervision Applicative et Système)
* **Objectif :** S'assurer que le système d'identité fonctionne avec une fluidité absolue sous la charge réelle de production.
* **Métriques Surveillées en Continu (Mise à jour toutes les 10 secondes) :**
  - **Taux de Succès des Requêtes (Success Rate) :** Doit rester supérieur à **99,99%**.
  - **Latence au Centile 95 (p95 Latency) :** Doit rester inférieure à **300 ms** pour les requêtes d'authentification simples.
  - **État d'utilisation des CPU/RAM/Disque** de tous les nœuds de calcul des deux datacenters d'Haïti.
* **Outils d'Observabilité :** Tableaux de bord de production Grafana temps réel alimentés par Prometheus, logs unifiés via Elasticsearch, traçage des requêtes distribuées via OpenTelemetry.

---

### 3.2. Security Alerts (Supervision de la Cybersécurité)
* **Objectif :** Détecter immédiatement toute tentative de déstabilisation cybernétique nationale orchestrée par des acteurs hostiles internes ou externes lors de l'annonce officielle du lancement.
* **Métriques Surveillées en Continu :**
  - Nombre de requêtes HTTP malveillantes bloquées par le pare-feu applicatif (WAF).
  - Tentatives d'attaques par force brute sur l'IdP national ou sur le portail des agents administratifs.
  - Alertes comportementales d'accès réseau en provenance des bastions d'administration de la PKI ou du moteur ABIS.
* **Outils d'Observabilité :** Écrans de contrôle SIEM (Elastic Security), flux de threat intelligence mis à jour en continu, console d'administration centralisée de l'EDR d'État.

---

### 3.3. Regional Rollout (Suivi Départemental du Déploiement)
* **Objectif :** Suivre la mise en service géographique progressive de la plateforme sur l'ensemble du territoire national d'Haïti.
* **Métriques Surveillées en Continu :**
  - Nombre de bureaux de l'ONI ouverts et connectés au SNISID par département (Ouest, Nord, Sud, Artibonite, etc.).
  - Nombre total d'identités numériques créées et validées par région en temps réel.
  - Cartographie thermique d'enrôlement et d'utilisation de la plateforme.
* **Visualisation :** Carte interactive géographique d'Haïti affichée sur l'écran central de la War Room avec code couleur par région (Vert : Déploiement complet | Jaune : Déploiement en cours avec perturbations | Rouge : Non ouvert ou hors ligne).

---

### 3.4. Citizen Incidents & Field Support (Suivi des Problèmes Terrain)
* **Objectif :** Identifier et résoudre immédiatement les anomalies bloquantes remontées par les citoyens d'Haïti ou par les opérateurs de terrain.
* **Métriques Surveillées en Continu :**
  - Temps moyen d'attente sur la ligne nationale d'assistance (800-SNISID).
  - Volume de requêtes de support ouvertes et résolues.
  - Taux d'incidents signalés par les citoyens sur l'application mobile d'identité numérique (ex: problèmes de reconnaissance faciale ou d'envoi d'OTP SMS par les opérateurs télécom d'Haïti comme Digicel ou Natcom).

---

## 4. Protocole d'Alerte et Procédure "Triage-War-Room"

En cas d'anomalie détectée par l'un des contrôleurs de la War Room, la procédure de triage d'urgence suivante est immédiatement appliquée :

```
[ANOMALIE DÉTECTÉE] 
        |
        v
[TRIAGE EN < 60 SECONDES] ===> Détermination de la gravité (SEV-1 à SEV-4).
        |
        +---> Si SEV-1 (Critique) : Le "Commander" de la War Room arrête les opérations 
        |     non critiques, active la cellule d'ingénieurs d'urgence et ordonne 
        |     l'application du Runbook d'Incident Majeur.
        |
        +---> Si SEV-2 (Majeur) : Isolement du composant affecté (ex: suspension temporaire 
        |     de la connexion d'un département) et routage du trafic d'urgence.
        |
        v
[RÉSOLUTON & RETOUR À LA NORMALE] ===> Analyse post-mortem obligatoire en moins de 12 heures.
```

---

## 5. Attestation de Disponibilité de la War Room

Le Directeur du Déploiement National certifie par le présent document que la **National War Room** est entièrement construite, configurée avec l'ensemble des écrans de télémétrie en temps réel, dotée des équipes techniques et des liaisons sécurisées, et déclarée **PRÊTE POUR L'ACTIVATION ET LE CONTRÔLE DU GOLIVE NATIONAL**.

```
[SIGNÉ ÉLECTRONIQUEMENT]
CONSEILLER SPÉCIAL AU COMMANDE ET AUX OPÉRATIONS DE DÉPLOIEMENT DU SNISID
```
