# VOLUME 4 : Résilience Nationale et Continuité (Disaster Recovery)
## Infrastructure Souveraine de Continuité de l'État — SNISID

L'État d'Haïti ne peut pas confier l'intégralité de son identité citoyenne à un seul bâtiment à Port-au-Prince. Cette architecture définit les mécanismes de survie à l'échelle macroscopique.

---

## 🏛️ CHAPITRE 1 : TOPOLOGIE ACTIVE-ACTIVE INTER-DATACENTERS

Le système est réparti sur deux centres de données majeurs distants géographiquement (éloignement sismique).

### Datacenter Primaire (DC1 - Port-au-Prince)
*   Héberge le cœur décisionnel de l'ABIS et de l'autorité de certification (AN-PKI).

### Datacenter Secondaire (DC2 - Cap-Haïtien)
*   Fonctionne en mode **Active-Active** (pas de délai de basculement, le trafic est équilibré en temps réel).
*   Gère 100% de la charge du Nord et de l'Artibonite.

```mermaid
graph TD
    subgraph Traffic_Routing [Global Load Balancing]
        DNS[National DNS Anycast]
        DNS -->|Poids: 60%| API1[Ingress DC1]
        DNS -->|Poids: 40%| API2[Ingress DC2]
    end

    subgraph DC1 [Port-au-Prince Datacenter]
        DB1[(CockroachDB Node)]
        K1[Kafka Broker]
    end

    subgraph DC2 [Cap-Haïtien Datacenter]
        DB2[(CockroachDB Node)]
        K2[Kafka Broker]
    end

    API1 --> DB1
    API2 --> DB2
    DB1 <==>|Réplication Synchrone| DB2
    K1 <==>|MirrorMaker 2.0 (Asynchrone)| K2
```

---

## 💾 CHAPITRE 2 : SYSTÈMES DE SAUVEGARDE SOUVERAINS (BACKUP)

Les sauvegardes du SNISID sont soumises au secret d'État.

### 2.1 Backups Immuables (WORM)
Toutes les nuits, un instantané (snapshot) crypté de la base de données est envoyé vers un cluster de stockage S3 MinIO configuré en mode **Object Lock (Compliance)**.
*   **Protection Ransomware:** Une fois écrit, le backup ne peut mathématiquement pas être supprimé ni modifié par quiconque, y compris l'administrateur système, pendant une durée de 20 ans.

### 2.2 Cold Storage (Bandes Magnétiques LTO-9)
Une fois par mois, un robot copie l'état de la nation sur des bandes magnétiques LTO-9 chiffrées.

---

## 🛡️ CHAPITRE 3 : AIR-GAPPED SYSTEMS ET BUNKER DE SURVIE

Face à une menace de cyberattaque militaire destructrice (Wiper / Ransomware d'État) :

1.  **Air-Gapped Vault (Chambre Forte Déconnectée):** Un troisième Datacenter "Bunker" (situé dans un lieu tenu secret par la défense nationale) se connecte au réseau principal de manière éphémère (30 minutes par jour à une heure aléatoire) pour aspirer les flux de Kafka.
2.  **Pull-Only:** Le bunker n'accepte aucune connexion entrante (Firewall Deny-All IN). Il initie lui-même la connexion (TCP Outbound).
3.  **Survie Ultime:** Si l'ensemble du réseau haïtien (DC1, DC2, Départements) est détruit ou compromis, les serveurs de ce bunker contiennent la copie intacte de l'identité de chaque Haïtien, prête à être restaurée sur du nouveau matériel.

---

## ⚡ CHAPITRE 4 : OPÉRATIONS DE CRISE (REGIONAL FAILOVER)

Si Port-au-Prince tombe (catastrophe sismique) :
*   Le BGP (Border Gateway Protocol) des FAI haïtiens (Natcom/Digicel) réachemine automatiquement les préfixes IP gouvernementaux vers Cap-Haïtien.
*   Le cluster K3s du Cap-Haïtien scale horizontalement (Auto-Scaling HPA) pour absorber les 60% de trafic manquants.
*   Le mode "Dégradé" s'active : les contrôles ABIS complets sont suspendus pour accélérer le traitement vital des réfugiés, remplacés par une validation PKI stricte.
