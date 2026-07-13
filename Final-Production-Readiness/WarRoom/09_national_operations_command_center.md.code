# National Operations Command Center (NOCC)
**Système National d'Identité Souveraine et d'Identité Digitale (SNISID) — République d'Haïti**
**Document ID:** SNISID-NOCC-PH20-009  
**Classification:** SECRET DE L'ÉTAT / OPÉRATIONS 24-7  
**Version:** 1.0.0  
**Date:** 25 Mai 2026  

---

## 1. Mission du National Operations Command Center (NOCC)

Le **National Operations Command Center (NOCC)** d'Haïti est le centre névralgique de supervision physique, technologique et sécuritaire de la plateforme SNISID. Opérationnel 24 heures sur 24, 7 jours sur 7 et 365 jours par an, le NOCC assure la surveillance en temps réel de l'état de santé du système national, prévient les pannes de services, identifie les cybermenaces et coordonne les interventions techniques immédiates.

```
========================================================================================
                          ORGANISATION ET STRUCTURE DU NOCC
========================================================================================
                      [ CO-DIRECTEURS OPÉRATIONNELS NOCC ]
                                       |
       +-------------------------------+-------------------------------+
       |                               |                               |
       v                               v                               v
[ UNITÉ INFRASTRUCTURE ]     [ UNITÉ CYBERSÉCURITÉ ]      [ UNITÉ INCIDENTS CITOYENS ]
    - Supervision Systèmes       - SOC Niveau 1, 2 & 3        - Support Bureaux ONI
    - Réseaux & Datacenters       - Analyse d'Alertes SIEM     - Retours d'Expérience
    - Base de données/ABIS       - Réponse à l'Incident       - Gestion des Files
========================================================================================
```

---

## 2. Piliers Opérationnels du NOCC

### 2.1. 24/7/365 Operations (Régime d'Exploitation Continue)
Le NOCC fonctionne sans aucune interruption. L'équipe d'ingénieurs et d'analystes d'Haïti est structurée selon un modèle de rotation militaire rigoureux :
* **Organisation des Rotations (Shifts) :** 
  - Trois quarts (shifts) quotidiens de 8 heures (Shift A : 06h00 - 14h00 | Shift B : 14h00 - 22h00 | Shift C : 22h00 - 06h00).
  - Chaque shift comprend un chef de quart (Duty Manager), 2 ingénieurs système/réseau, 2 analystes SOC de sécurité, et 2 spécialistes support applicatif.
* **Protocole de Passation de Consignes (Shift Handover) :**
  - Une réunion de passation formelle de 15 minutes est obligatoire à la fin de chaque shift. Elle fait l'objet d'un rapport numérique signé dans le registre du NOCC détaillant les incidents en cours, les maintenances prévues, et les niveaux d'alerte.

---

### 2.2. National Monitoring (Supervision Technique et Métier)
La surveillance technique est assurée par une suite d'observabilité souveraine (Prometheus, Grafana, OpenTelemetry) affichant en temps réel les indicateurs clés sur le mur d'écrans géant du NOCC :
* **Métriques Système (Santé de l'Infrastructure) :**
  - Température et humidité des salles serveurs des datacenters.
  - Taux de charge CPU, RAM et I/O de l'ABIS et des bases de données.
  - Disponibilité des liens réseau d'IP Transit et de transport de données inter-datacenters.
* **Métriques Métier (Activité Citoyenne) :**
  - Volume d'enrôlements en temps réel par département d'Haïti.
  - Temps de traitement des transactions biométriques.
  - Taux de réussite de génération des identités numériques uniques.

---

### 2.3. Security Monitoring (Supervision de la Cybersécurité)
Intégrée au SOC central, cette cellule surveille en continu la surface d'attaque nationale :
* **Indicateurs de Surveillance de la Sécurité :**
  - Volume de requêtes suspectes interceptées par les pare-feux applicatifs (WAF).
  - Tentatives de connexions administratives ou VPN non autorisées ou suspectes.
  - Anomalies de flux de données (ex: tentatives d'exfiltration de base de données).
  - Détection de signatures de codes malveillants ou d'activités suspectes sur les terminaux d'enrôlement ou postes administratifs par l'EDR.

---

### 2.4. Crisis Escalation Matrix (Matrice d'Escalade de Crise)
En cas d'incident technique ou de sécurité, le NOCC applique immédiatement des protocoles d'escalade prédéfinis selon la gravité (Severity Level - SEV) de l'événement :

| Sévérité de l'Incident | Critères d'Impact | Temps d'Alerte Maximal | Destinataire de l'Escalade | Plan d'Action Immédiat |
| :--- | :--- | :--- | :--- | :--- |
| **SEV-1 (CRITIQUE)** | Arrêt complet du SNISID, de la PKI ou de l'ABIS nationale. | < 2 minutes | Directeur Technique, Conseil de Sécurité Nationale, Ministres. | Activation de la cellule de crise, basculement manuel ou automatique du Datacenter (DR). |
| **SEV-2 (MAJEUR)** | Arrêt partiel (ex: 1 département déconnecté ou portail citoyen indisponible). | < 10 minutes | Directeur des Opérations, Responsable d'Infrastructure. | Intervention des ingénieurs d'astreinte, routage de secours satellite activé. |
| **SEV-3 (MODÉRÉ)** | Ralentissement des temps de réponse applicatifs (> 2 secondes). | < 30 minutes | Chef d'équipe d'astreinte technique. | Optimisation des index de base de données, allocation temporaire de vCPU additionnels. |
| **SEV-4 (MINEUR)** | Anomalie cosmétique sur l'application ou besoin d'assistance utilisateur. | < 4 heures | Équipe Support Applicatif Niveau 1/2. | Résolution via le processus standard de ticketing (Jira Service Management). |

---

## 3. Infrastructures Physiques et Logiques du Centre

* **Emplacement Physique :** Bâtiment hautement fortifié et sécurisé par la Police Nationale, muni de vitrages pare-balles, d'une cage de Faraday, de contrôles d'accès biométriques stricts, d'une autonomie électrique totale (générateurs + panneaux solaires), et de liaisons de communication par satellite (Starlink militaire / VSAT).
* **Salle de Contrôle (Control Room) :** Équipée d'un mur d'images de 12 écrans LED géants pour la visualisation en temps réel des graphiques de performance et des cartes d'alertes géospatiales d'Haïti.

---

## 4. Attestation d'Opérationnalité Pré-GoLive

Le Comité d'Exploitation du SNISID atteste par la présente que le **National Operations Command Center (NOCC)** est entièrement construit, équipé, doté d'un personnel haïtien certifié et opérationnel à 100%. Il supervise activement la plateforme SNISID 24h/24 en mode pilote et est déclaré prêt pour le GoLive national complet.

```
[SIGNÉ ÉLECTRONIQUEMENT]
DIRECTEUR GÉNÉRAL DE L'OFFICE NATIONAL D'OPÉRATIONS ET DE COMMANDEMENT (NOCC-HT)
```
