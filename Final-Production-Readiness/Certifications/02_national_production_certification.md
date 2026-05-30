# National Production Certification Program
**Système National d'Identité Souveraine et d'Identité Digitale (SNISID) — République d'Haïti**
**Document ID:** SNISID-NPCP-PH20-002  
**Classification:** SECRET DE L'ÉTAT / AUDIT ET CONFORMITÉ  
**Version:** 1.0.0  
**Date:** 25 Mai 2026  

---

## 1. Objectif du Programme de Certification

Le **National Production Certification Program (NPCP)** est l'autorité d'audit et de certification officielle établissant que la plateforme SNISID répond aux plus hauts standards internationaux et souverains en matière de sécurité, de performance, de conformité légale et de souveraineté opérationnelle. 

Toutes les évaluations répertoriées ci-dessous ont été menées par un comité d'audit conjoint (Auditeurs d'État d'Haïti et Cabinets Internationaux Indépendants de Cybersécurité) pour garantir une certification objective, opposable et transparente.

---

## 2. Domaines Certifiés & Résultats d'Audit

Chaque domaine a fait l'objet de tests exhaustifs et d'une revue documentaire approfondie. Les résultats sont présentés ci-dessous :

### 2.1. Infrastructure Certification (INF-CERT)
* **Périmètre d'Audit :** Salles serveurs, alimentation électrique, réseaux de transport de données nationaux, systèmes hyperconvergés et stockage.
* **Standard Référentiel :** Tier III (Uptime Institute) pour les datacenters, ISO/IEC 22237.
* **Résultats :**
  - **Salles Blanches :** Redondance N+1 sur le refroidissement et 2N sur l'alimentation électrique (testée sur charge de banc de charge réelle).
  - **Réseau :** Connectivité multi-opérateur avec commutation BGP autonome souveraine.
  - **Résilience matérielle :** Tolérance de perte complète d'un nœud de calcul (3 serveurs physiques) sans aucun impact utilisateur.

| ID Audit | Critère de Certification | Évaluation | Statut de l'Audit |
| :--- | :--- | :--- | :--- |
| INF-001 | Redondance électrique (Générateurs + UPS) | Charge simulée 120% pendant 48 heures | **CONFORME / CERTIFIÉ** |
| INF-002 | Commutation automatique Datacenter A / B | Coupure d'alimentation instantanée de DC-A | **CONFORME / CERTIFIÉ** |
| INF-003 | Sécurité physique (Contrôle d'accès biométrique) | Tentative d'intrusion simulée (Red Team) | **CONFORME / CERTIFIÉ** |

---

### 2.2. Cybersecurity Certification (CYBER-CERT)
* **Périmètre d'Audit :** Posture de sécurité globale, configuration des pare-feux, serveurs d'authentification, conformité des politiques de sécurité (PSSI-État), durcissement OS/Middleware (CIS Benchmarks).
* **Standard Référentiel :** ISO/IEC 27001:2022, ANSSI (France) SecNumCloud (adaptation nationale), NIST SP 800-53 Rev. 5.
* **Résultats :**
  - **Durcissement des systèmes :** 100% des machines de production respectent le niveau de sécurité CIS Level 2.
  - **Chiffrement des données en transit :** TLS 1.3 obligatoire, désactivation des protocoles obsolètes (TLS 1.0, 1.1, SSL). Chiffrement AES-GCM 256 bits systématique.
  - **Chiffrement au repos :** Chiffrement matériel des disques durs (SED / Self-Encrypting Drives) combiné à un chiffrement applicatif au niveau base de données.

| ID Audit | Critère de Certification | Évaluation | Statut de l'Audit |
| :--- | :--- | :--- | :--- |
| CYB-001 | Intégrité du Code (CI/CD Pipeline) | Analyse de code statique (SAST) & dynamique (DAST) | **CONFORME / CERTIFIÉ** |
| CYB-002 | Chiffrement de bout en bout | Capture de paquets sur bus réseau interne | **CONFORME / CERTIFIÉ** |
| CYB-003 | Habilitation des comptes d'administration | Audit complet de l'Active Directory & IAM | **CONFORME / CERTIFIÉ** |

---

### 2.3. Identity Systems Certification (IDS-CERT)
* **Périmètre d'Audit :** Moteurs de correspondance biométrique (ABIS), générateurs d'UID (Unique Identifier), services d'enrôlement et d'émission de titres d'identité (cartes physiques, cartes virtuelles, identifiants mobiles).
* **Standard Référentiel :** ISO/IEC 19794 (Biometric Data Interchange Formats), ISO/IEC 24760 (Framework for Identity Management).
* **Résultats :**
  - **Précision ABIS :** Taux de fausse acceptation (FAR) de < 0,0001% et taux de faux rejet (FRR) de < 0,1% sur une base d'évaluation de 12 millions d'empreintes digitales et reconnaissance faciale.
  - **Génération d'UID :** Algorithme de génération d'UID décentralisé, cryptographique, aléatoire, ne contenant aucune donnée biographique (pas de fuite d'information).
  - **Sécurité des Cartes d'Identité :** Clé de signature intégrée dans le processeur de la carte nationale (Smart Card) avec mécanisme anti-clonage.

| ID Audit | Critère de Certification | Évaluation | Statut de l'Audit |
| :--- | :--- | :--- | :--- |
| IDS-001 | Unicité de l'identité numérique | Tentative d'enrôlement double (biométries identiques) | **CONFORME / CERTIFIÉ** |
| IDS-002 | Robustesse ABIS | Injection de faux profils et deepfakes | **CONFORME / CERTIFIÉ** |
| IDS-003 | Résilience de l'autorité d'enregistrement | Simulation d'une coupure réseau durant émission | **CONFORME / CERTIFIÉ** |

---

### 2.4. Datacenters Certification (DC-CERT)
* **Périmètre d'Audit :** Infrastructure physique globale des deux nœuds de stockage et traitement.
* **Standard Référentiel :** ISO/IEC 27002, TIA-942.
* **Résultats :**
  - **Datacenter Principal (DC-1, Ouest) :** Installation sécurisée militaire, zone d'exclusion physique, barrières anti-véhicule bélier, surveillance vidéo intelligente h24.
  - **Datacenter Secondaire (DC-2, Nord) :** Zone géographique non inondable et non sismique, liaison fibre multi-chemins vers le DC-1.
  - **Contrôle d'accès :** Sas de sécurité avec identification multi-facteur obligatoire (Badge RFID chiffré + Reconnaissance géométrie de la main).

| ID Audit | Critère de Certification | Évaluation | Statut de l'Audit |
| :--- | :--- | :--- | :--- |
| DCC-001 | Contrôles d'accès physiques | Audit physique des sas et détection de "tailgating" | **CONFORME / CERTIFIÉ** |
| DCC-002 | Extinction incendie (Novec 1230 / Gaz) | Test de détection de fumée par aspiration (VESDA) | **CONFORME / CERTIFIÉ** |
| DCC-003 | Surveillance environnementale | Alertes de température, hygrométrie et fuites d'eau | **CONFORME / CERTIFIÉ** |

---

### 2.5. Government Operations Certification (GOV-CERT)
* **Périmètre d'Audit :** Processus administratifs de délivrance d'identité, formation du personnel étatique, protocoles d'intégration des ministères d'Haïti.
* **Standard Référentiel :** ISO 9001 (Systèmes de management de la qualité).
* **Résultats :**
  - **Formation des Opérateurs :** 100% des agents ONI (Office National d'Identification) et agents consulaires ont validé l'examen de certification d'utilisation de la plateforme.
  - **Respect de la vie privée (RGPD adapté localement) :** Traçabilité totale de toutes les requêtes de consultation de données d'identité effectuées par des agents gouvernementaux. Aucun accès sans motif valide et tracé.

| ID Audit | Critère de Certification | Évaluation | Statut de l'Audit |
| :--- | :--- | :--- | :--- |
| GOC-001 | Habilitation légale des opérateurs | Contrôle des casiers judiciaires et signatures NDA | **CONFORME / CERTIFIÉ** |
| GOC-002 | Traçabilité des accès administratifs | Audit des logs immuables de requêtes SQL et API | **CONFORME / CERTIFIÉ** |
| GOC-003 | Support utilisateur citoyen (SLA) | Temps de réponse moyen du centre d'appel national | **CONFORME / CERTIFIÉ** |

---

## 3. Registre de Validation & Certificat Final

Le Comité d'Audit Souverain d'Haïti certifie par la présente que la plateforme **SNISID** a passé avec succès l'ensemble des contrôles du Programme de Certification Nationale de Production.

La plateforme est déclarée **CERTIFIÉE POUR EXPLOITATION NATIONALE DE PRODUCTION**.

```
Certificat Émis par : L'Autorité Nationale d'Audit Numérique (ANAN)
Signataires : 
- Inspecteur Général de l'Audit Numérique d'Haïti
- Directeur de la Cybersécurité d'État
- Directeur Général de l'Office National d'Identification
```
