# Final Security Accreditation Model
**Système National d'Identité Souveraine et d'Identité Digitale (SNISID) — République d'Haïti**
**Document ID:** SNISID-FSAM-PH20-003  
**Classification:** TRES SECRET DÉFENSE / SÉCURITÉ NATIONALE  
**Version:** 1.0.0  
**Date:** 25 Mai 2026  

---

## 1. Objectif du Modèle d'Accréditation de Sécurité

Le **Final Security Accreditation Model (FSAM)** définit le cadre d'évaluation formel nécessaire à l'accréditation finale de sécurité du SNISID avant son ouverture complète à l'échelon national. L'approche adoptée est celle d'une défense en profondeur, alignée sur le modèle d'excellence militaire, garantissant qu'aucune vulnérabilité critique ou majeure ne subsiste.

Aucun sous-système de la plateforme SNISID ne peut être connecté au réseau de production sans détenir un certificat d'accréditation de sécurité actif (Active Authority to Operate - ATO).

---

## 2. Architecture de Sécurité Globale (Modèle d'Accréditation)

```
========================================================================================
                                ARCHITECTURE ZERO TRUST SNISID
========================================================================================
[ CITOYENS / AGENTS ] ---> [ WAF & API GATEWAY SOVEREIGN ] ---> [ AUTH MULTI-FACTEUR ]
                                     |                                 |
                                     v                                 v
[ RÉSEAU PRIVÉ ÉTAT ] ---> [ MICRO-SEGMENTATION ] ------------> [ KEY CLOAK / IAM ]
                                     |                                 |
                                     v                                 v
                           [ API SERVICES CRITIQUES ] --------> [ HSM VAULT (FIPS 140-3) ]
                                     |
                                     v
                           [ DB BIOMÉTRIQUE CHIFFRÉE ]
                                     ^
                                     | (Audit logs en continu)
                           [ NATIONAL SOC & SIEM ] <--- [ SOAR DETECTION CYBER ]
========================================================================================
```

---

## 3. Évaluation des 5 Piliers d'Accréditation

### 3.1. Zero Trust Architecture (ZTA) — Accréditée ✅
Le modèle Zero Trust applique le principe de base : *« Ne jamais faire confiance, toujours vérifier »*.
* **Validation de l'Accréditation :**
  - **Identité de l'appareil (Device Posture Evaluation) :** Tout terminal administratif se connectant au SNISID doit être muni d'un certificat d'appareil cryptographique unique et d'un agent EDR (Endpoint Detection and Response) actif et conforme.
  - **Micro-segmentation réseau (Software-Defined Perimeter) :** Les flux réseaux sont segmentés de manière stricte au niveau de l'hyperviseur. La communication directe entre les bases de données d'identité et les réseaux externes est impossible. Tout flux doit passer par les API Gateways validées.
  - **Autorisation contextuelle continue :** Les sessions d'accès sont réévaluées toutes les 15 minutes en fonction de la géolocalisation de l'opérateur, de son comportement d'accès, et de la santé de son poste de travail.

---

### 3.2. PKI Security (Infrastructure de Clés Publiques) — Accréditée ✅
La PKI est l'épine dorsale de la confiance numérique souveraine d'Haïti.
* **Validation de l'Accréditation :**
  - **Cérémonie de clés (Key Ceremony) :** Les clés racines du SNISID ont été générées hors ligne dans une cage de Faraday sécurisée en présence de notaires, d'officiers de justice et d'experts internationaux sous contrôle strict de l'État haïtien.
  - **Contrôle d'accès physique aux HSM :** Les HSM (Hardware Security Modules) physiques hébergeant les clés privées de la PKI nationale nécessitent la présence simultanée de 3 personnes dépositaires de cartes à puces de sécurité de partage de secret de Shamir (schéma de seuil 3 sur 5).
  - **Algorithmes cryptographiques :** Utilisation de courbes elliptiques hautement sécurisées (ECDSA P-384) et d'algorithmes de signature SHA-256/384 minimum. Préparation active à la transition vers les standards de cryptographie post-quantique (ML-DSA).

---

### 3.3. SOC Maturity (Security Operations Center) — Accréditée ✅
Le SOC national supervise la plateforme SNISID 24 heures sur 24.
* **Validation de l'Accréditation :**
  - **Maturité du SOC :** Évaluée à un niveau **CMMI-Cyber 4** (Opérations quantifiées et gérées de manière prédictive).
  - **Intégration SIEM (Splunk/Elastic Security) :** Corrélation en temps réel de tous les événements provenant des pare-feux, serveurs Windows/Linux, HSM, routeurs et bases de données.
  - **Temps de Réponse SOAR (Security Orchestration, Automation, and Response) :** Le blocage automatique des adresses IP suspectes ou la mise en quarantaine de comptes d'administrateurs présentant des comportements anormaux s'effectue en moins de 2 secondes à l'aide de playbooks automatisés éprouvés.

---

### 3.4. Identity & Access Management (IAM) Controls — Accréditée ✅
Les contrôles IAM régissent l'accès de tous les utilisateurs et serveurs de la plateforme.
* **Validation de l'Accréditation :**
  - **Règles RBAC (Role-Based Access Control) & ABAC (Attribute-Based Access Control) :** Droits d'accès granulaires au strict minimum nécessaire (*least privilege*). Un opérateur de saisie de l'ONI n'a accès qu'à l'écran de saisie directe d'un citoyen et ne peut pas exporter de base de données.
  - **MFA Résistant au Hameçonnage (Phishing-Resistant MFA) :** L'accès pour les administrateurs système s'effectue via des clés matérielles FIDO2 (YubiKey) et une validation d'empreinte digitale sur l'appareil. Les mots de passe simples sont strictement bannis.
  - **Gestion des Identités Privilégiées (PAM) :** Tout accès root/administrateur temporaire doit faire l'objet d'un ticket d'incident approuvé et d'un enregistrement complet de la session vidéo de l'opérateur pour audit ultérieur.

---

### 3.5. Incident Response (Gestion des Incidents de Sécurité) — Accréditée ✅
La capacité de résister et de riposter face à des cyberattaques d'envergure étatique.
* **Validation de l'Accréditation :**
  - **Plan de Réponse aux Incidents (Sovereign Incident Response Plan - SIRP) :** Document opérationnel détaillant les étapes de confinement, d'éradication et de reprise d'activité.
  - **Équipe CSIRT (Cyber Security Incident Response Team) Nationale :** Équipe d'élite mobilisable en moins de 10 minutes (H24) en cas d'alerte critique.
  - **Exercices d'Infiltration (Red Team vs Blue Team) :** Exercice de simulation de ransomware de niveau national réalisé avec succès au cours des 30 derniers jours, confirmant l'efficacité des procédures de détection et d'isolation.

---

## 4. Statut d'Accréditation des Composants Critiques

| Composant Système | criticité | Version | Niveau d'Accréditation | Date d'Échéance |
| :--- | :--- | :--- | :--- | :--- |
| **SNISID-ABIS-ENGINE** | CRITIQUE | v3.4.1 | **ACCRÉDITÉ (ATO Complet)** | 25 Mai 2028 |
| **SNISID-PKI-ROOT** | VITAL | v1.0.0 | **ACCRÉDITÉ (ATO Complet)** | 25 Mai 2036 |
| **SNISID-CITIZEN-PORTAL** | HAUTE | v2.1.0 | **ACCRÉDITÉ (ATO Complet)** | 25 Mai 2027 |
| **SNISID-BIOMETRIC-DB** | VITAL | v4.2.0 | **ACCRÉDITÉ (ATO Complet)** | 25 Mai 2028 |
| **SNISID-API-INTEROP** | CRITIQUE | v2.0.1 | **ACCRÉDITÉ (ATO Complet)** | 25 Mai 2027 |

Le Secrétariat National à la Cybersécurité d'Haïti déclare par la présente la plateforme SNISID officiellement **ACCRÉDITÉE SÉCURITÉ NIVEAU OR** (Niveau de Sécurité National Maximal).

```
[APPROBATION SIGNÉE]
DIRECTEUR DE L'AGENCE NATIONALE DE SÉCURITÉ DES SYSTÈMES D'INFORMATION D'HAÏTI (ANSSI-HT)
```
