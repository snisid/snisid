# SNISID Final Production Runbooks
**Système National d'Identité Souveraine et d'Identité Digitale (SNISID) — République d'Haïti**
**Document ID:** SNISID-FPRB-PH20-017  
**Classification:** SECRET DE L'ÉTAT / EXPLOITATION TECHNIQUE  
========================================================================================

---

## 1. Introduction et Normes d'Exploitation

Ce document réunit les cinq **Runbooks Finaux de Production** obligatoires pour l'exploitation nationale du SNISID. Chaque procédure décrit de manière séquentielle et rigoureuse les actions à mener par les équipes de la War Room et du Command Center (NOCC) face aux événements majeurs de production.

---

## RUNBOOK 1 : National GoLive (Activation Finale de la Plateforme)

* **Objectif :** Activer la plateforme SNISID à l'échelon national pour l'ensemble des 12 millions de citoyens haïtiens.
* **Acteurs Clés :** Commander de la War Room, Directeur Technique (CTO), Administrateur de la PKI, DBA Principal, Analyste SOC.

### Séquence d'Activation (Étape par Étape) :

1. **Vérification de l'homologation de sécurité :**
   - S'assurer que le certificat d'accréditation (ATO) est signé et valide.
2. **Cérémonie de libération des clés PKI (HSM Cérémonie) :**
   - Insérer les cartes d'administrateurs PKI pour initialiser le HSM et libérer la clé privée de production nationale.
3. **Mise sous tension des bases de données de production :**
   - Établir la liaison de synchronisation bidirectionnelle synchrone entre DC-1 (Port-au-Prince) et DC-2 (Nord).
   - Commande de vérification de l'état de la grappe de bases de données :
     ```bash
     patronictl -c /etc/patroni/patroni.yml list
     ```
4. **Activation de la Gateway API de Production :**
   - Déployer les configurations de pare-feu et libérer les routes d'accès pour les applications citoyennes et les API d'interopérabilité (ONI, Police, DGI, Immigration, Justice).
5. **Vérification de la connectivité réseau Anycast :**
   - Libérer l'adresse IP globale souveraine sur les routeurs BGP d'Haïti.
6. **Notification aux Administrations Publiques :**
   - Envoyer le signal de connectivité complète ("Green Signal") à tous les bureaux locaux de l'ONI et aux forces de sécurité sur le terrain.

---

## RUNBOOK 2 : Production Rollback (Plan de Repli et Restauration)

* **Objectif :** Restaurer un état antérieur stable et sécurisé de la plateforme en cas d'anomalie critique majeure détectée immédiatement après une mise à jour ou un GoLive infructueux.
* **Acteurs Clés :** Administrateur Système (SRE), DBA Principal, Responsable d'Infrastructure.

### Séquence de Repli (Étape par Étape) :

1. **Détection et Décision de Repli (Rollback Decision) :**
   - Si le taux de succès des requêtes de production chute en dessous de **95%** pendant plus de 5 minutes après une mise à jour, la décision de rollback est automatique.
2. **Activation du Circuit Breaker d'Urgence :**
   - Rediriger instantanément tout le trafic utilisateur vers une page d'assistance statique ("Maintenance Planifiée").
3. **Rollback des Conteneurs Applicatifs (Kubernetes) :**
   - Restaurer l'image de conteneur stable précédente via Helm :
     ```bash
     helm rollback snisid-core-deployment 2
     ```
4. **Vérification de l'Intégrité de la Base de Données :**
   - Exécuter le script de diagnostic de cohérence des schémas PostgreSQL et Cassandra.
   - Si une corruption de données de schéma est détectée, initier le point de restauration transactionnel immédiat (Point-in-Time Recovery - PITR) à la minute précédant la mise à jour :
     ```bash
     pg_backrest --stanza=snisid-main --type=time --target="2026-05-25 14:29:59" restore
     ```
5. **Validation Post-Rollback :**
   - Exécuter la suite de tests de fumée (Smoke Tests) automatiques.
6. **Réactivation du Trafic National :**
   - Rétablir la Gateway API et annuler la page de maintenance.

---

## RUNBOOK 3 : National Cyber Crisis (Escalade et Riposte Cyber)

* **Objectif :** Contenir et neutraliser une cyberattaque de niveau étatique (DDoS massif, tentative de vol de données biométriques, infiltration par ransomware).
* **Acteurs Clés :** Directeur du SOC National, Analyste Sécurité Niveau 3, Administrateur PKI, Conseiller à la Cybersécurité d'État.

### Séquence de Riposte (Étape par Étape) :

1. **Détection et Identification de l'Attaque :**
   - Alerte critique SIEM (ex: injection de requêtes de consultation massives d'identités ou tentative d'accès physique non autorisé aux HSM).
2. **Isolement Réseau Immédiat (Isolation Mode) :**
   - Couper tous les accès VPN administratifs et isoler l'environnement de production bureautique de la zone d'infrastructure critique SNISID.
   - Règle de blocage IP globale instantanée via l'API Gateway :
     ```bash
     kubectl apply -f /etc/kubernetes/network-policies/deny-all-ingress.yaml
     ```
3. **Révocation des Identifiants Suspects :**
   - Révoquer immédiatement les certificats d'administration ou les jetons d'accès (JWT) des comptes suspectés de compromission dans Keycloak/IAM.
4. **Analyse Forensique et Nettoyage :**
   - Prendre des instantanés (Snapshots) mémoire des machines virtuelles suspectes pour analyse ultérieure, puis détruire et recréer les instances applicatives via l'infrastructure déclarative (GitOps).
5. **Restauration de la Confiance PKI :**
   - Si une clé intermédiaire de signature d'identité est suspectée d'être compromise, révoquer immédiatement le certificat via la Liste de Révocation de Certificats (CRL) de la PKI souveraine et distribuer le nouveau certificat.
6. **Rétablissement Graduel sous Haute Surveillance :**
   - Reconnecter les agences d'intégration l'une après l'autre après audit de conformité de sécurité de leur point d'accès.

---

## RUNBOOK 4 : Mass Citizen Incident (Stabilisation en Surcharge de Trafic)

* **Objectif :** Stabiliser la plateforme face à un afflux extrême et inattendu de requêtes de citoyens (par exemple, suite à une rumeur publique ou une urgence nationale).
* **Acteurs Clés :** Ingénieur Cloud SRE, Duty Manager du NOCC, Directeur de la Communication.

### Séquence de Stabilisation (Étape par Étape) :

1. **Activation de l'Auto-scaling Horizontal :**
   - Augmenter les limites de réplication des pods de traitement de l'identité :
     ```bash
     kubectl scale deployment snisid-identity-api --replicas=50
     ```
2. **Activation de la Triage Queue (Rate Limiting) :**
   - Configurer l'API Gateway pour limiter les requêtes par utilisateur (Rate Limiting) afin de protéger la base de données centrale :
     ```nginx
     limit_req zone=citizen_ip burst=10 delay=5;
     ```
3. **Mise en cache agressive (Caching Mode) :**
   - Forcer la mise en cache des requêtes de vérification d'identité statiques non sensibles dans Redis pour soulager le moteur PostgreSQL principal.
4. **Activation de la Dégradation Gracieuse (Graceful Degradation) :**
   - Désactiver temporairement les modules non vitaux du portail citoyen (comme la consultation de l'historique des connexions ou le téléchargement de documents d'aide de grande taille) pour allouer l'intégralité du calcul aux authentifications critiques et enrôlements.
5. **Communication Publique Immédiate :**
   - Diffuser un message d'information sur les réseaux nationaux et les radios pour expliquer la situation et rassurer la population sur la résilience du système.

---

## RUNBOOK 5 : Emergency Offline Mode (Mode de Continuité Souverain)

* **Objectif :** Configurer les terminaux et guichets nationaux pour fonctionner en mode totalement déconnecté (Hors Ligne) lors d'une perte totale des liaisons réseau nationales.
* **Acteurs Clés :** Administrateur d'Infrastructure de Terrain, Agents ONI locaux, Force POLIFRONT.

### Séquence de Continuité (Étape par Étape) :

1. **Détection de la Perte de Liaison Réseau :**
   - Le terminal mobile d'enrôlement ou de vérification détecte la perte de liaison avec l'API Gateway centrale.
2. **Bascule Automatique en Mode Local Sécurisé :**
   - Le système local charge la base de données locale chiffrée contenant les listes de révocation d'identité et de signatures biométriques clés.
3. **Exécution des Vérifications d'Identité Hors Ligne :**
   - Le terminal effectue l'authentification 1:1 directement entre la carte d'identité physique du citoyen (contenant la puce NFC cryptée) et l'empreinte digitale ou faciale lue par le capteur de l'appareil.
4. **Stockage Chiffré des Nouvelles Données d'Enrôlement :**
   - Les nouvelles données d'enrôlement collectées localement sont signées numériquement par la clé privée du terminal de l'agent et stockées dans l'enclave matérielle inviolable de l'appareil (Secure Element).
5. **Rétablissement du Réseau et Synchronisation des Données :**
   - Une fois la connexion rétablie, le terminal établit une session mTLS sécurisée et téléverse par petits paquets les transactions hors ligne accumulées. Le système central intègre les données après vérification d'intégrité et déduplication ABIS globale.

---

```
========================================================================================
       [RUNBOOKS FINAUX D'EXPLOITATION DU SNISID - ADOPTÉS ET CERTIFIÉS]
       
       Signataires :
       - Le Chef des Opérations Techniques du SNISID
       - Le Directeur de la Continuité d'Activité de l'État d'Haïti
========================================================================================
```
