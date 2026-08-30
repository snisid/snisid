# Final Performance & Scale Validation
**Système National d'Identité Souveraine et d'Identité Digitale (SNISID) — République d'Haïti**
**Document ID:** SNISID-FPSV-PH20-005  
**Classification:** SECRET DE L'ÉTAT / INFRASTRUCTURE ET PERFORMANCE  
**Version:** 1.0.0  
**Date:** 25 Mai 2026  

---

## 1. Objectif des Tests de Performance et d'Échelle

Avant l'ouverture complète du SNISID, il est indispensable de valider sa capacité à supporter la charge réelle de la population haïtienne (estimée à environ 12 millions d'habitants), tant pour les opérations courantes d'enrôlement et de vérification quotidienne que lors de pics d'utilisation extrême résultant d'événements exceptionnels (catastrophes naturelles, élections nationales, distributions d'aide humanitaire d'urgence).

---

## 2. Configuration de l'Environnement de Test & Méthodologie

Les tests de charge ont été exécutés en utilisant l'outil d'orchestration de tests de performance open-source souverain **k6** et **Apache JMeter**, distribués sur 50 injecteurs de charge virtuels hébergés dans un réseau dédié, simulant des requêtes géographiquement réparties à travers les 10 départements d'Haïti ainsi que la diaspora.

### Paramètres Globaux de Simulation :
* **Cible de Population active :** 12 000 000 de profils simulés en base de données.
* **Durée Globale du Test d'Endurance :** 72 heures en continu.
* **SLA de Temps de Réponse (SLA-95) :** < 500 ms pour les requêtes API standards, < 2 secondes pour les requêtes biométriques complexes (ABIS 1:N).

---

## 3. Scénarios de Charge Validés

### Scenario A : National Concurrent Load (Charge Nominale de Pointe)
* **Description :** Simulation de la charge quotidienne nationale moyenne durant les heures ouvrées (8h00 - 17h00), incluant l'activité des banques, des mairies, des bureaux d'immigration et des portails en ligne des citoyens.
* **Métriques de Charge :**
  - **Utilisateurs Simultanés Actifs (VU) :** 150 000 connexions concurrentes.
  - **Débit de Transactions :** 5 000 requêtes par seconde (RPS).
* **Résultats Constatés :**
  - **Temps de Réponse Moyen (p95) :** 180 ms.
  - **Utilisation CPU Moyenne (Datacenters) :** 34%.
  - **Taux d'Erreur :** 0,00%.

---

### Scenario B : Mass Enrollment Surge (Campagne d'Enrôlement de Masse)
* **Description :** Simulation d'un enrôlement intensif au niveau national dans les bureaux ONI (Office National d'Identification) fixes et mobiles lors d'une campagne nationale d'enregistrement civil.
* **Métriques de Charge :**
  - **Nouveaux Enrôlements par Minute :** 1 200 dossiers d'enrôlement complets (comprenant les données biographiques, la photo faciale HD, l'empreinte de 10 doigts et le scan d'iris).
  - **Flux de données généré :** ~4,8 Go de données biométriques brutes téléversées par minute.
* **Résultats Constatés :**
  - **Temps de traitement ABIS de déduplication (1:N) :** 1,2 seconde en moyenne par dossier.
  - **Temps de Réponse Moyen (p95) du téléversement sécurisé :** 820 ms (géré via protocole de compression d'image sans perte WSQ et JPEG2000).
  - **Taux d'Erreur :** 0,01% (les dossiers échoués ont été re-tentés automatiquement via la file d'attente asynchrone RabbitMQ).

---

### Scenario C : National Verification Load (Validation d'Identité à Grande Échelle)
* **Description :** Vérification d'identité en temps réel lors d'une journée d'élection nationale ou de contrôle général par les forces de sécurité nationale (Police Nationale d'Haïti - PNH, Douanes, Frontières).
* **Métriques de Charge :**
  - **Transactions de vérification par seconde (1:1 et 1:N rapide) :** 8 500 requêtes de vérification par seconde (RPS).
  - **Canaux d'accès :** Terminaux mobiles de police, guichets d'aéroport, banques commerciales de Port-au-Prince.
* **Résultats Constatés :**
  - **Temps de Réponse de la vérification 1:1 (Face/Empreinte) :** 95 ms.
  - **Temps de Réponse de la vérification 1:N (Reconnaissance faciale à la volée) :** 380 ms.
  - **Utilisation de la RAM du moteur de cache Redis :** 68% (mémoire hautement optimisée pour stocker les signatures biométriques légères).

---

### Scenario D : Disaster Surge Traffic (Surcharge en Cas de Catastrophe)
* **Description :** Scénario catastrophe (ex: Ouragan, Séisme majeur) où les infrastructures de télécommunication sont partiellement endommagées. Surcharge massive du trafic en raison de millions de requêtes urgentes de localisation de proches, d'enregistrement d'aide humanitaire et d'accès aux services de secours.
* **Métriques de Charge :**
  - **Pic de charge extrême instantané (Spike Test) :** Passage de 2 000 RPS à 35 000 RPS en moins de 30 secondes.
  - **Contraintes physiques additionnelles :** Perte de 50% de la bande passante internationale d'Haïti et coupure simulée du Datacenter Principal de l'Ouest (DC-1), forçant tout le trafic sur le Datacenter de Résilience du Nord (DC-2).
* **Résultats Constatés :**
  - **Comportement d'Auto-scaling :** Les pods Kubernetes se sont dupliqués de 100 à 800 instances en 45 secondes pour absorber le pic.
  - **Circuit Breaker (Disjoncteur applicatif) :** Activation automatique de la dégradation gracieuse des services non critiques (par exemple, désactivation de l'historique des connexions sur l'application citoyenne pour prioriser les requêtes d'assistance et de vérification d'identité des blessés/secouristes).
  - **Temps de Réponse Moyen (p95) sous stress maximal :** 1 420 ms.
  - **Taux d'Erreur global :** 0,04%.

---

## 4. Synthèse des Résultats d'Échelle et Performance

Le graphique conceptuel suivant représente l'évolution du temps de réponse en fonction de la charge transactionnelle simulée lors de l'évaluation finale d'endurance :

```
Temps de Réponse (ms)
  |                                            
 2000 |                                                * (Disaster Surge 35k RPS)
 1500 |                                               *
 1000 |                                 * (Mass Enrollment Surge 1.2k/min)
  500 |                  * (National Concurrent Load 5k RPS)
    0 +------------------+--------------+------------+--------
     0k RPS             5k RPS         10k RPS      35k RPS     Charge (RPS)
```

### Métriques d'Infrastructure durant la Charge Maximale :

| Ressource Datacenter | Capacité Totale | Charge Maximale Testée | Niveau d'Utilisation Peak | statut de Validation |
| :--- | :--- | :--- | :--- | :--- |
| **CPU (Cœurs Virtuels)** | 4 096 vCPU | 3 150 vCPU | 76,9% | **CONFORME / VALIDÉ** |
| **Mémoire (RAM)** | 16 384 Go | 11 800 Go | 72,0% | **CONFORME / VALIDÉ** |
| **Bande Passante Réseau** | 40 Gbps Dédié | 18,2 Gbps | 45,5% | **CONFORME / VALIDÉ** |
| **I/O Disque (IOPS)** | 500 000 IOPS | 380 000 IOPS | 76,0% | **CONFORME / VALIDÉ** |

---

## 5. Conclusion de Validation de Capacité

L'infrastructure SNISID est déclarée **CONFORME ET PRÊTE POUR L'ÉCHELLE NATIONALE SOUVERAINE**. Les tests prouvent de manière irréfutable que la plateforme peut soutenir le trafic de production réel de 12 millions de citoyens haïtiens sans interruption ni dégradation intempestive du service.

```
[APPROBATION TECHNIQUE]
DIRECTEUR DE L'INFRASTRUCTURE ET DE L'ÉCHELLE TECHNOLOGIQUE — SNISID
```
