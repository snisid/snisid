# SNISID Runbook — National Deployment Crisis Management Playbook
**Code de Procédure :** SNISID-RB-05  
**Statut :** Approuvé  
**Audience :** Incident Commander, Comité Exécutif de Crise, Responsable Sécurité Nationale  

---

## 1. Objectif

Ce runbook régit la gestion des événements de crise de niveau S1 (Critique) impactant la sécurité physique, logique, géopolitique ou climatique du déploiement national du SNISID en Haïti.

---

## 2. Procédures Spécifiques par Scénario de Crise

### SCÉNARIO A : Vol ou Compromission Physique d'un Local Edge Node (LEN)
*   *Déclencheur :* Signalement du vol d'un serveur d'Edge Node dans un bureau de liaison ou agression physique d'un site par des bandes armées.
*   *Protocole d'Action Immédiate :*
    1. **Révocation Cryptographique :** L'ingénieur SecOps de la War Room révoque immédiatement le certificat de sécurité mTLS rattaché au LEN volé dans l'infrastructure centrale :
       `kubectl exec -it central-ca-pod -- revoke-client-certificate --node-id=LEN-JEREMIE-04`
    2. **Autodestruction de Clé à Distance (Remote Wipe/Zeroisation) :** Envoyer une commande d'effacement de secours si le serveur est encore allumé et connecté en satellite :
       `ansible -i production_nodes.ini -m shell -a "dd if=/dev/urandom of=/dev/mapper/luks_root bs=1M count=100 && reboot" LEN-JEREMIE-04`
       *Remarque : Même sans cette commande, le disque SSD crypté par LUKS (AES-256) reste totalement illisible sans la YubiKey matérielle de l'opérateur et le mot de passe dynamique lié au TPM 2.0 physique.*
    3. **Alerte Forces de l'Ordre :** Contacter la Direction Générale de la PNH pour coordonner la récupération du matériel et la sécurisation du périmètre communal.

### SCÉNARIO B : Catastrophe Naturelle (Séisme, Cyclone Majeur)
*   *Déclencheur :* Passage d'un ouragan de catégorie 4 ou séisme majeur dans une péninsule d'Haïti.
*   *Protocole d'Action Immédiate :*
    1. **Mise à l'abri du matériel :** Ordonner aux opérateurs locaux de démonter les antennes Starlink et les terminaux biométriques mobiles et de les stocker dans des caisses étanches prévues à cet effet.
    2. **Sauvegarde Finale Interne :** Effectuer une sauvegarde manuelle de la base de données locale du LEN sur une clé USB durcie scellée.
    3. **Bascule en Mode Secours Régional :** Les opérations d'enrôlement physique sont suspendues sur les communes sinistrées. L'enrôlement est redirigé vers des structures temporaires d'urgence (tentes de secours SNISID dotées d'unités mobiles alimentées par kits solaires) déployées par la Protection Civile.

### SCÉNARIO C : Détection de Tentative d'Intrusion Logique Massive (Cyberattaque)
*   *Déclencheur :* Tentative d'intrusion SSH non autorisée ou attaque par déni de service distribué (DDoS) sur l'API Gateway centrale.
*   *Protocole d'Action Immédiate :*
    1. **Isolation des Edge Nodes :** Configurer les pare-feux centraux pour n'autoriser les flux gRPC entrants que depuis la liste blanche IP stricte des routeurs Starlink du SNISID.
    2. **Suspension des Connexions Externes :** Couper temporairement les API de consultation des tiers (banques commerciales, ministères partenaires) pour soulager la charge du serveur central et isoler la menace.
    3. **Analyse Forensic :** Extraire les logs de connexion suspectés de l'API Gateway via Loki pour analyse de l'origine de l'attaque :
       `logcli query '{job="snisid-api-gateway"} |= "invalid signature" |= "unauthorized"' --limit=1000`
