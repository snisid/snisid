# SNISID — Architecture Haute Disponibilité & Reprise sur Sinistre
# SNISID — High Availability & Disaster Recovery Architecture

---

| Métadonnée | Valeur |
|---|---|
| **Document ID** | SNISID-ARC-HA-DR-001 |
| **Version** | 1.0.0 |
| **Date** | 2026-05-25 |
| **Statut** | APPROUVÉ — Production |
| **Classification** | CONFIDENTIEL / CONFIDENTIAL |
| **Auteur** | Architecture & Infrastructure Team — SNISID |
| **Révisé par** | Chief Architect, Infrastructure Lead, CISO |
| **Approuvé par** | DG SNISID, Ministère concerné |
| **Standards** | ISO 22301:2019, ISO/IEC 27031, NIST SP 800-34, TIA-942 |

---

## Table des Matières

1. [Vue d'Ensemble HA/DR](#1-vue-densemble-hadr)
2. [Topologie Active-Active](#2-topologie-active-active)
3. [Objectifs RTO/RPO par Niveau de Service](#3-objectifs-rtorpo-par-niveau-de-service)
4. [Pyramide de Résilience à 5 Niveaux](#4-pyramide-de-résilience-à-5-niveaux)
5. [Stratégie de Réplication de Bases de Données](#5-stratégie-de-réplication-de-bases-de-données)
6. [Kafka MirrorMaker 2 — Streaming Cross-DC](#6-kafka-mirrormaker-2--streaming-cross-dc)
7. [Procédures de Basculement](#7-procédures-de-basculement)
8. [Architecture DR Offshore](#8-architecture-dr-offshore)
9. [Procédures de Drill DR](#9-procédures-de-drill-dr)
10. [Monitoring HA/DR](#10-monitoring-hadr)

---

## 1. Vue d'Ensemble HA/DR

### 1.1 Stratégie Globale

Le SNISID adopte une stratégie de **Haute Disponibilité Active-Active** entre les deux datacenters souverains haïtiens (Port-au-Prince et Cap-Haïtien), complétée par un **Coffre-fort Numérique Offshore** (Islande ou Suisse) pour la récupération de dernier recours.

```mermaid
graph TB
    subgraph "Mode Nominal — Active-Active"
        direction LR
        PAP_DC["DC Port-au-Prince\n🟢 ACTIF\n60% du trafic\nCapacité: 12 nœuds K8s"]
        CAP_DC["DC Cap-Haïtien\n🟢 ACTIF\n40% du trafic\nCapacité: 8 nœuds K8s"]
        SYNC_ARROW["← Réplication Synchrone →\nPostgreSQL Patroni\nHashiCorp Vault\nKafka MirrorMaker 2"]
    end

    subgraph "Scénario 1 — Perte DC PAP"
        CAP_FULL["DC Cap-Haïtien\n🔴 ACTIF 100%\n100% du trafic\nCapacité: +4 nœuds (burst)"]
        PAP_OFFLINE["DC Port-au-Prince\n⬛ HORS LIGNE"]
        PAP_ARROW["Basculement automatique < 5 minutes"]
    end

    subgraph "Scénario 2 — Perte DC CAP"
        PAP_FULL["DC Port-au-Prince\n🔴 ACTIF 100%\n100% du trafic"]
        CAP_OFFLINE["DC Cap-Haïtien\n⬛ HORS LIGNE"]
    end

    subgraph "Scénario 3 — Catastrophe nationale"
        OFFSHORE["🌍 DR Offshore\nIslande OU Suisse\n⚠️ Actif en urgence\nServices réduits"]
        HAITI_OFFLINE["🇭🇹 Haïti\nInfrastructure indisponible"]
    end

    PAP_DC <-->|"MPLS 10Gbps + Dark Fiber"| CAP_DC
```

### 1.2 Principes HA/DR

| Principe | Description | Implémentation |
|---|---|---|
| **Zéro perte de données** | RPO = 0 pour données critiques | Réplication synchrone PostgreSQL + Kafka |
| **Basculement transparent** | Aucune intervention manuelle pour panne de composant | Patroni auto-failover, Istio circuit breaker |
| **Dégradation gracieuse** | Services essentiels maintenus même en capacité réduite | Mode offline agents terrain, cache distribué |
| **Test continu** | La capacité de récupération prouvée régulièrement | Drill mensuel composants, drill DR complet semestriel |
| **Documentation vivante** | Procédures à jour et validées | RevQ trimestrielle obligatoire |
| **Souveraineté maintenue** | Aucune dépendance externe pour opérations critiques | Infrastructure 100% on-premise Haïti |

---

## 2. Topologie Active-Active

### 2.1 Diagramme Active-Active Détaillé

```mermaid
graph TB
    subgraph "Couche DNS / Load Balancing Global"
        GDNS["DNS Global — Anycast\nsnisid.gouv.ht\nGSLB: F5 DNS + PowerDNS"]
        CDN_GLOBAL["CDN Haïti (anycast)\nPortail + Assets statiques"]
    end

    subgraph "DC Port-au-Prince — ACTIF PRIMAIRE"
        direction TB
        subgraph "Edge PAP"
            LB_PAP["F5 BIG-IP HA Pair\nVIP: 203.x.x.100\nPoids: 60"]
            GW_PAP["Kong Gateway Cluster\n3 instances\nACTIF"]
        end
        subgraph "Kubernetes PAP"
            K8S_PAP["Cluster K8s PAP\n3 Control Plane + 12 Workers\nNamespaces: tous actifs"]
            SERVICES_PAP["Services: Identity, Biometric,\nAuth, Enrollment, Audit,\nDocuments, Interop"]
        end
        subgraph "Data PAP"
            PG_PRIMARY_PAP["PostgreSQL PRIMARY\n10.30.0.10 (RW)\nMaster actif"]
            PG_REPLICA_PAP["PostgreSQL Replica\n10.30.0.11 (RO)\nLecture locale"]
            REDIS_PAP["Redis Cluster\nMaster PAP"]
            VAULT_PAP["Vault Cluster\nActive Node PAP"]
        end
    end

    subgraph "Liaison Inter-DC"
        WAN1["MPLS Primary\n10Gbps — Natcom\nLatence: ~8ms"]
        WAN2["MPLS Secondary\n10Gbps — Digicel\nLatence: ~10ms"]
        DARKFIBER["Dark Fiber\n100Gbps — Futur\nLatence: ~5ms"]
    end

    subgraph "DC Cap-Haïtien — ACTIF SECONDAIRE"
        direction TB
        subgraph "Edge CAP"
            LB_CAP["F5 BIG-IP HA Pair\nVIP: 203.x.x.200\nPoids: 40"]
            GW_CAP["Kong Gateway Cluster\n2 instances\nACTIF"]
        end
        subgraph "Kubernetes CAP"
            K8S_CAP["Cluster K8s CAP\n3 Control Plane + 8 Workers\nNamespaces: tous actifs"]
            SERVICES_CAP["Services: Identity, Biometric,\nAuth, Enrollment, Audit"]
        end
        subgraph "Data CAP"
            PG_STANDBY_CAP["PostgreSQL STANDBY\n10.130.0.10\nPromouvable en PRIMARY"]
            PG_REPLICA_CAP["PostgreSQL Replica\n10.130.0.11 (RO)"]
            REDIS_CAP["Redis Cluster\nSlave → Élu Primary si panne PAP"]
            VAULT_CAP["Vault Cluster\nPerf Standby"]
        end
    end

    subgraph "DR Offshore"
        DR_SITE["Coffre-fort Numérique\nIslande (IS) ou Suisse (CH)\nActif sur décision executive"]
    end

    GDNS -->|"60% trafic"| LB_PAP
    GDNS -->|"40% trafic"| LB_CAP
    CDN_GLOBAL --> GDNS

    LB_PAP --> GW_PAP --> K8S_PAP
    K8S_PAP --> SERVICES_PAP
    SERVICES_PAP --> PG_PRIMARY_PAP
    SERVICES_PAP --> REDIS_PAP
    SERVICES_PAP --> VAULT_PAP

    LB_CAP --> GW_CAP --> K8S_CAP
    K8S_CAP --> SERVICES_CAP
    SERVICES_CAP --> PG_STANDBY_CAP
    SERVICES_CAP --> REDIS_CAP
    SERVICES_CAP --> VAULT_CAP

    PG_PRIMARY_PAP <-->|"Réplication sync\nWAL streaming"| PG_STANDBY_CAP
    REDIS_PAP <-->|"Redis replication\nasync"| REDIS_CAP
    VAULT_PAP <-->|"Vault Enterprise\nReplication"| VAULT_CAP

    WAN1 -.->|"Lien actif"| LB_PAP
    WAN1 -.->|"Lien actif"| LB_CAP
    WAN2 -.->|"Lien backup"| LB_PAP
    WAN2 -.->|"Lien backup"| LB_CAP

    PG_PRIMARY_PAP -->|"Backup async chiffré"| DR_SITE
    PG_STANDBY_CAP -->|"Backup async chiffré"| DR_SITE
```

### 2.2 Distribution du Trafic — GSLB

```yaml
# GSLB Configuration — F5 DNS + PowerDNS avec GTM
gslb_config:
  zone: "api.snisid.gouv.ht"
  ttl: 30  # TTL court pour basculement rapide

  pools:
    - name: POOL_PAP
      members:
        - ip: "203.x.x.100"
          port: 443
          weight: 60
          health_check: "/health"
          datacenter: port-au-prince

    - name: POOL_CAP
      members:
        - ip: "203.x.x.200"
          port: 443
          weight: 40
          health_check: "/health"
          datacenter: cap-haitien

  strategy:
    normal: "weighted_round_robin"  # 60% PAP, 40% CAP
    pap_failure: "all_to_cap"       # Tout vers CAP si PAP tombe
    cap_failure: "all_to_pap"       # Tout vers PAP si CAP tombe
    detection_interval: 10s
    failover_threshold: 3           # 3 échecs consécutifs = failover
    recovery_threshold: 5           # 5 succès = retour au poids normal

  health_checks:
    interval: 10s
    timeout: 5s
    expected_codes: [200]
    expected_body: '"status":"healthy"'
```

---

## 3. Objectifs RTO/RPO par Niveau de Service

### 3.1 Classification des Services

```mermaid
graph TD
    subgraph "Niveau 1 — CRITIQUE (RTO: 2min, RPO: 0)"
        L1A[Identity Service — Vérification identité]
        L1B[Authentication Service — Tokens/Sessions]
        L1C[API Gateway — Point d'entrée]
        L1D[Audit Service — Traçabilité légale]
        L1E[Identity Database — Données civiles]
        L1F[Biometric Service — Matching 1:1]
    end

    subgraph "Niveau 2 — ÉLEVÉ (RTO: 15min, RPO: 5min)"
        L2A[Enrollment Service — Enrôlement]
        L2B[Biometric Database — Templates]
        L2C[Kafka Cluster — Event streaming]
        L2D[HashiCorp Vault — Secrets]
        L2E[Redis Cluster — Cache/Sessions]
    end

    subgraph "Niveau 3 — IMPORTANT (RTO: 1h, RPO: 15min)"
        L3A[Document Service — Génération]
        L3B[Interop Gateway — Synchronisation]
        L3C[Notification Service — SMS/Email]
        L3D[Search Service — Elasticsearch]
    end

    subgraph "Niveau 4 — STANDARD (RTO: 4h, RPO: 1h)"
        L4A[Admin Portal — Interface admin]
        L4B[Analytics Service — Rapports]
        L4C[CI/CD Pipeline — GitLab/ArgoCD]
        L4D[Monitoring Stack — Prometheus/Grafana]
    end

    subgraph "Niveau 5 — BEST EFFORT (RTO: 24h, RPO: 4h)"
        L5A[Archive Service — Archivage long terme]
        L5B[Audit Reports — Rapports d'audit]
        L5C[Training Systems — Systèmes de formation]
    end
```

### 3.2 Tableau RTO/RPO Complet

| Service | Tier | RTO | RPO | Stratégie HA | Stratégie DR |
|---|---|---|---|---|---|
| **API Gateway (Kong)** | L1 | 2 min | N/A | Active-Active 3+2 instances | DNS failover auto |
| **Identity Service** | L1 | 2 min | 0 | K8s Deployment 3+ replicas, 2 DC | Kubernetes failover |
| **Auth Service (Keycloak)** | L1 | 2 min | 30 s | Cluster 3 nodes, 2 DC | Session replication |
| **Biometric Service** | L1 | 5 min | 0 | 2+ replicas GPU, 2 DC | K8s failover |
| **Audit Service** | L1 | 2 min | 0 | 3+ replicas, 2 DC | CockroachDB multi-DC |
| **Identity DB (PostgreSQL)** | L1 | 5 min | 0 | Patroni HA, sync replication | Async replication CAP |
| **Enrollment Service** | L2 | 15 min | 5 min | 2+ replicas, 2 DC | K8s restart |
| **Biometric Vault (DB)** | L2 | 10 min | 5 min | Patroni HA locale, async CAP | Restore from backup |
| **Kafka Cluster** | L2 | 5 min | 0 | 5 brokers RF=3, 2 DC via MM2 | MirrorMaker 2 |
| **HashiCorp Vault** | L2 | 5 min | 0 | Cluster 3 nodes, Performance Standby CAP | Vault Replication Enterprise |
| **Redis Cluster** | L2 | 10 min | 5 min | 3 masters + replicas, 2 DC | Redis replication |
| **Document Service** | L3 | 1 h | 15 min | 2 replicas, 1 DC | K8s restart + PVC restore |
| **Interop Gateway** | L3 | 1 h | 15 min | 2 replicas | K8s restart |
| **Notification Service** | L3 | 1 h | 30 min | 2 replicas + queue persistence | Queue replay |
| **Search (Elasticsearch)** | L3 | 2 h | 1 h | 3 nodes cluster | Snapshot restore |
| **Admin Portal** | L4 | 4 h | 1 h | 2 replicas | K8s restart |
| **Analytics (Spark/Trino)** | L4 | 4 h | 1 h | Stateless, restart | Redémarrage |
| **Monitoring Stack** | L4 | 4 h | 2 h | 2 Prometheus replicas | Backup Grafana + Prometheus |
| **Archive Service** | L5 | 24 h | 4 h | Single replica | Restore depuis Ceph |

---

## 4. Pyramide de Résilience à 5 Niveaux

```mermaid
graph TB
    subgraph "🔺 Pyramide de Résilience SNISID"
        direction TB
        L5_TOP["🏔️ NIVEAU 5 — DR OFFSHORE CATASTROPHE\nIslande / Suisse — Récupération totale\nActivation: perte totale Haïti\nDélai: 4-24 heures\nPerte données max: 24h"]
        L4["🌊 NIVEAU 4 — DR RÉGIONAL HAÏTI\nBasculement DC2 → DC1 ou DC1 → DC2\nActivation: perte complète d'un DC\nDélai: < 5 minutes automatique\nPerte données max: 0 (sync)"]
        L3["🏗️ NIVEAU 3 — FAILOVER ZONE\nBasculement entre zones réseau\nActivation: panne réseau ou zone\nDélai: < 2 minutes\nPerte données max: 0"]
        L2["⚙️ NIVEAU 2 — HAUTE DISPONIBILITÉ CLUSTER\nRéplication et basculement intra-cluster\nActivation: panne nœud/pod\nDélai: 30 secondes (Patroni/K8s)\nPerte données max: 0"]
        L1_BOT["💾 NIVEAU 1 — REDONDANCE COMPOSANT\nRedondance matérielle (PSU, disques, NIC)\nActivation: panne hardware\nDélai: Instantané (RAID, bonding)\nPerte données max: 0"]
    end

    L5_TOP --> L4 --> L3 --> L2 --> L1_BOT

    style L5_TOP fill:#8B0000,color:#fff
    style L4 fill:#cc4400,color:#fff
    style L3 fill:#d4a017,color:#000
    style L2 fill:#2d6a2d,color:#fff
    style L1_BOT fill:#1a4a8a,color:#fff
```

### 4.1 Niveau 1 — Redondance Composant

| Composant | Redondance | Mécanisme | MTTR |
|---|---|---|---|
| Alimentation serveur | 2 PSU (A+B) | Dual PSU, PDU séparés | Immédiat |
| Stockage | RAID 10 NVMe | Hotswap | < 30 min |
| Réseau serveur | Bonding 2×25GbE | LACP (802.3ad) | Immédiat |
| Réseau DC | Switches MLAG | MLAG Arista (< 1s) | < 1 s |
| Alimentation DC | UPS N+1 + Générateurs N+1 | ATS (Automatic Transfer Switch) | < 30 s |

### 4.2 Niveau 2 — HA Cluster

```mermaid
sequenceDiagram
    participant MON as Patroni Monitor
    participant PRIMARY as PostgreSQL Primary
    participant REPLICA1 as PostgreSQL Replica 1
    participant ETCD as etcd DCS
    participant APP as Application Services

    Note over MON,APP: Détection panne PostgreSQL Primary

    PRIMARY -x MON: [TIMEOUT — Primary non responsive]
    MON->>ETCD: Primary health check failed (3 consecutive)
    ETCD->>MON: Release primary lock
    MON->>REPLICA1: Candidate pour promotion
    REPLICA1->>REPLICA1: pg_ctl promote
    REPLICA1->>ETCD: Acquire primary lock
    ETCD->>MON: Replica 1 est nouveau Primary
    MON->>APP: Mise à jour endpoint RW → Replica 1
    Note right of APP: Downtime: ~30 secondes
    Note right of APP: Perte données: 0 (sync replication)
    APP->>REPLICA1: Connexions rétablies
```

### 4.3 Niveau 3 — Failover Zone

- **Trigger** : Perte réseau d'une zone VLAN (ex. VLAN 20 Application inaccessible)
- **Mécanisme** : Kubernetes Node Not Ready → Pods eviction → Rescheduling sur nœuds disponibles
- **Délai** : 1-2 minutes (kubectl Node timeout: 40s + pod scheduling)
- **Pré-requis** : Capacité suffisante dans zones survivantes

### 4.4 Niveau 4 — DR Régional HAïti

Documenté en détail dans la section [Procédures de Basculement](#7-procédures-de-basculement).

### 4.5 Niveau 5 — DR Offshore Catastrophe

Documenté en détail dans la section [Architecture DR Offshore](#8-architecture-dr-offshore).

---

## 5. Stratégie de Réplication de Bases de Données

### 5.1 Architecture de Réplication PostgreSQL

```mermaid
graph TB
    subgraph "FLUX DE RÉPLICATION — MULTI-NIVEAUX"
        direction TB

        subgraph "Site PAP — Réplication Locale"
            PAP_PRI["PostgreSQL Primary PAP\nWAL producer\nmode: synchronous"]
            PAP_R1["Replica 1 PAP\nStreaming SYNC\nlag_max: 0\nWAL receiver"]
            PAP_R2["Replica 2 PAP\nStreaming ASYNC\nlag_max: 30s\nRead-only"]
        end

        subgraph "Liaison WAN — Réplication Cross-DC"
            WAL_SENDER["WAL Sender Process\nPAP Primary"]
            WAL_RECEIVER["WAL Receiver Process\nCAP Standby"]
        end

        subgraph "Site CAP — Réplication Distante"
            CAP_STANDBY["PostgreSQL Standby CAP\nHot standby\nStreaming ASYNC depuis PAP\nPromovable"]
            CAP_REPLICA["Replica CAP\nStreaming SYNC depuis CAP Standby\nRead-only local"]
        end

        subgraph "Archivage WAL"
            MINIO_PAP["MinIO PAP\nWAL Archive\nRétention: 30 jours"]
            MINIO_CAP["MinIO CAP\nWAL Archive miroir\nRétention: 30 jours"]
            OFFSHORE_BACKUP["Backup Offshore\nChiffré AES-256\nRétention: 7 ans"]
        end
    end

    PAP_PRI -->|"Sync WAL streaming\nconfirmation avant commit"| PAP_R1
    PAP_PRI -->|"Async WAL streaming"| PAP_R2
    PAP_PRI -->|"WAL Sender → WAN MPLS"| WAL_SENDER
    WAL_SENDER -->|"WAL Receiver"| WAL_RECEIVER
    WAL_RECEIVER -->|"Async apply"| CAP_STANDBY
    CAP_STANDBY -->|"Sync local"| CAP_REPLICA

    PAP_PRI -->|"archive_command"| MINIO_PAP
    MINIO_PAP -->|"Replication MinIO"| MINIO_CAP
    MINIO_CAP -->|"Backup quotidien chiffré"| OFFSHORE_BACKUP
```

### 5.2 Configuration Réplication Synchrone/Asynchrone

```sql
-- postgresql.conf — Paramètres de réplication
-- Sur le Primary PAP:

-- Réplication synchrone vers Replica 1 PAP (même DC)
-- Réplication asynchrone vers Standby CAP (cross-DC)
synchronous_commit = on
synchronous_standby_names = 'FIRST 1 (replica1_pap)'
-- "FIRST 1" = au moins 1 replica sync doit confirmer avant commit
-- CAP Standby est async: pas de risque de latence WAN sur les commits

-- Slots de réplication pour éviter suppression WAL avant consommation
-- PAP Replicas:
-- SELECT pg_create_physical_replication_slot('replica1_pap');
-- SELECT pg_create_physical_replication_slot('replica2_pap');
-- CAP Standby:
-- SELECT pg_create_physical_replication_slot('cap_standby');

-- Monitoring réplication:
-- SELECT client_addr, state, sent_lsn, write_lsn, flush_lsn, replay_lsn,
--        write_lag, flush_lag, replay_lag, sync_state
-- FROM pg_stat_replication;
```

```yaml
# Patroni Configuration — patroni.yml
scope: snisid-identity-cluster
namespace: /snisid/
name: postgres-primary-pap

restapi:
  listen: 10.30.0.10:8008
  connect_address: 10.30.0.10:8008

etcd3:
  hosts: 10.30.0.50:2379,10.30.0.51:2379,10.30.0.52:2379
  protocol: https
  cacert: /etc/ssl/etcd/ca.crt
  cert: /etc/ssl/etcd/client.crt
  key: /etc/ssl/etcd/client.key

bootstrap:
  dcs:
    ttl: 30
    loop_wait: 10
    retry_timeout: 30
    maximum_lag_on_failover: 1048576  # 1 MB — failover si lag > 1MB
    maximum_lag_on_syncnode: -1       # Pas de limite sur sync node
    synchronous_mode: true            # Mode synchrone global
    synchronous_mode_strict: false    # Tolérer 0 sync si replica down
    postgresql:
      use_pg_rewind: true
      use_slots: true
      parameters:
        wal_level: replica
        hot_standby: "on"
        max_wal_senders: 10
        max_replication_slots: 10
        synchronous_commit: "on"
        archive_mode: "on"
        archive_timeout: 300

  initdb:
    - encoding: UTF8
    - locale: fr_HT.UTF-8
    - data-checksums

  pg_hba:
    - host replication replicator 10.30.0.0/24 md5
    - host replication replicator 10.130.0.0/24 md5  # CAP Standby
    - hostssl all all 10.20.0.0/20 md5

postgresql:
  listen: 10.30.0.10:5432
  connect_address: 10.30.0.10:5432
  data_dir: /data/postgresql/16/main
  bin_dir: /usr/lib/postgresql/16/bin
  config_dir: /etc/postgresql/16/main

  authentication:
    superuser:
      username: postgres
      password: '{VAULT_SECRET:postgresql/data/superuser}'
    replication:
      username: replicator
      password: '{VAULT_SECRET:postgresql/data/replicator}'
    rewind:
      username: rewind_user
      password: '{VAULT_SECRET:postgresql/data/rewind}'

  callbacks:
    on_start: /etc/patroni/callbacks/on_start.sh
    on_stop: /etc/patroni/callbacks/on_stop.sh
    on_role_change: /etc/patroni/callbacks/on_role_change.sh  # Notifie monitoring

tags:
  nofailover: false
  noloadbalance: false
  clonefrom: false
  nosync: false
```

### 5.3 Stratégie de Backup PITR

```yaml
backup_strategy:
  method: "WAL-G + pgBackRest"
  
  full_backup:
    frequency: "Quotidien — 02:00 AM HT"
    retention: "30 jours full"
    storage: ["MinIO PAP", "MinIO CAP", "Offshore chiffré"]
    compression: "zstd niveau 3"
    encryption: "AES-256-CBC (clé dans Vault)"

  wal_archiving:
    mode: continu
    destination: "s3://snisid-wal/$(hostname)/%Y/%m/%d/"
    endpoint: "https://minio.snisid.gouv.ht"
    retention: "30 jours"
    
  point_in_time_recovery:
    granularity: "1 seconde (WAL continu)"
    max_rpo: "5 minutes (archive_timeout=300)"
    test_restore: "Mensuel sur environnement DR test"

  verification:
    checksum_validation: true
    test_restore_schedule: "Dimanche 04:00 AM — DC test"
    alerting: "Si backup échoue: PagerDuty + email DBA on-call"
```

---

## 6. Kafka MirrorMaker 2 — Streaming Cross-DC

### 6.1 Architecture MirrorMaker 2

```mermaid
graph LR
    subgraph "Kafka Cluster PAP — Source"
        KB_PAP1[Broker 1 PAP\n:9092]
        KB_PAP2[Broker 2 PAP\n:9092]
        KB_PAP3[Broker 3 PAP\n:9092]
        KB_PAP4[Broker 4 PAP\n:9092]
        KB_PAP5[Broker 5 PAP\n:9092]

        TOPICS_PAP["Topics PAP:\nidentity.events\nbiometric.events\nenrollment.events\naudit.events"]
    end

    subgraph "MirrorMaker 2 — Active-Active"
        direction TB
        MM2_PAP["MM2 Connector — PAP→CAP\nSource: pap-cluster\nTarget: cap-cluster\nTopics: .*\nReplication Factor: 3\nOffset sync: enabled"]

        MM2_CAP["MM2 Connector — CAP→PAP\nSource: cap-cluster\nTarget: pap-cluster\nTopics: interop.*\nReplication Factor: 3\nOffset sync: enabled"]

        HEARTBEAT["MirrorHeartbeat\nLiveness check\nInterval: 1s"]
        CHECKPOINT["MirrorCheckpoint\nConsumer offset sync\nInterval: 60s"]
    end

    subgraph "Kafka Cluster CAP — Target"
        KB_CAP1[Broker 1 CAP\n:9092]
        KB_CAP2[Broker 2 CAP\n:9092]
        KB_CAP3[Broker 3 CAP\n:9092]

        TOPICS_CAP["Topics CAP (miroirs):\npap.identity.events\npap.biometric.events\npap.enrollment.events\npap.audit.events"]
    end

    KB_PAP1 & KB_PAP2 & KB_PAP3 -->|"Consume"| MM2_PAP
    MM2_PAP -->|"Produce (SSL)"| KB_CAP1 & KB_CAP2 & KB_CAP3
    MM2_PAP --> HEARTBEAT
    MM2_PAP --> CHECKPOINT

    KB_CAP1 -->|"Consume local"| MM2_CAP
    MM2_CAP -->|"Produce back"| KB_PAP1
```

### 6.2 Configuration MirrorMaker 2

```yaml
# mm2.properties — MirrorMaker 2 Configuration
# Clusters
clusters = pap, cap

pap.bootstrap.servers = kafka-b1.snisid.internal:9093,kafka-b2.snisid.internal:9093,kafka-b3.snisid.internal:9093
pap.security.protocol = SSL
pap.ssl.truststore.location = /etc/kafka/ssl/kafka.truststore.jks
pap.ssl.keystore.location = /etc/kafka/ssl/kafka.pap.keystore.jks
pap.ssl.keystore.password = ${KAFKA_MM2_KEYSTORE_PASS}

cap.bootstrap.servers = kafka-c1.snisid.internal:9093,kafka-c2.snisid.internal:9093,kafka-c3.snisid.internal:9093
cap.security.protocol = SSL
cap.ssl.truststore.location = /etc/kafka/ssl/kafka.truststore.jks
cap.ssl.keystore.location = /etc/kafka/ssl/kafka.cap.keystore.jks

# Flows de réplication
pap->cap.enabled = true
cap->pap.enabled = true

# Topics à répliquer PAP→CAP (tous sauf les miroirs eux-mêmes)
pap->cap.topics = identity\.events, biometric\.events, enrollment\.events, audit\.events, notification\.commands, interop\.sync
pap->cap.topics.exclude = .*\.MirrorHeartbeat

# Topics à répliquer CAP→PAP (uniquement interop pour éviter boucle)
cap->pap.topics = interop\.sync
cap->pap.topics.exclude = pap\..*

# Replication settings
replication.factor = 3
tasks.max = 8
offset-syncs.topic.replication.factor = 3
heartbeats.topic.replication.factor = 3
checkpoints.topic.replication.factor = 3

# Sync consumer group offsets (pour failover transparent)
sync.group.offsets.enabled = true
sync.group.offsets.interval.seconds = 60

# Failover consumer configuration
offset.lag.max = 100000

# Monitoring
metrics.enabled = true
metrics.reporter.classes = org.apache.kafka.common.metrics.JmxReporter

# Compression
producer.compression.type = lz4

# Performance
producer.batch.size = 32768
producer.linger.ms = 5
consumer.max.poll.records = 5000
```

---

## 7. Procédures de Basculement

### 7.1 Procédure Basculement DC PAP → DC CAP (Manuel)

```mermaid
flowchart TD
    A([🚨 INCIDENT DÉTECTÉ — DC PAP INACCESSIBLE]) --> B{Type d'incident?}

    B -->|Panne automatiquement détectée| C[Patroni auto-failover\nPostgreSQL CAP → PRIMARY\n≈ 30 secondes]
    B -->|Décision manuelle| D[Activation Plan de Basculement\nAutorisation: DG ou CTO SNISID]

    C --> E[Vérification Patroni:\npatronictl -c config.yml list]
    D --> E

    E --> F{PostgreSQL CAP est Primary?}
    F -->|Non| G[Promotion manuelle:\npatronictl failover cap-cluster --master cap-node\n⚠️ Risque perte données si async]
    F -->|Oui| H[DNS Update:\napi.snisid.gouv.ht → VIP CAP\nTTL 30s → propagation ~1min]

    G --> H

    H --> I[Vault CAP: Unseal si nécessaire\nvault operator unseal -address=https://vault-cap:8200]

    I --> J[GSLB Update:\nF5 DNS: poids PAP=0, CAP=100%]

    J --> K[Vérification Services K8s CAP:\nkubectl get pods -A -n snisid-* — contexte cap]

    K --> L{Tous les services RUNNING?}
    L -->|Non| M[Restart pods défaillants:\nkubectl rollout restart deployment X\nVérifier logs: kubectl logs]
    L -->|Oui| N[Test Fonctionnel:\ncurl https://api.snisid.gouv.ht/health]

    M --> N

    N --> O{Tests OK?}
    O -->|Non| P[Escalade: Architecture & DBA on-call\nActivation War Room Teams]
    O -->|Oui| Q[Communication:\n- Agences gouvernementales\n- Tableau de bord status\n- Email direction]

    Q --> R[Monitoring renforcé:\nAlerts threshold ×50% plus sensibles\nRevue toutes les 15 minutes]

    R --> S[Post-Incident:\nAnalyse root cause\nBilan DR\nMise à jour runbooks]

    style A fill:#8B0000,color:#fff
    style D fill:#cc4400,color:#fff
    style Q fill:#2d6a2d,color:#fff
```

### 7.2 Runbook Complet — Basculement PostgreSQL

```bash
#!/bin/bash
# RUNBOOK: PostgreSQL Failover PAP → CAP
# Fichier: /opt/snisid/runbooks/postgresql-failover-pap-to-cap.sh
# Auteur: SNISID Infrastructure Team
# Révision: v1.0 — 2026-05-25

set -euo pipefail
LOG_FILE="/var/log/snisid/dr-$(date +%Y%m%d-%H%M%S).log"
exec > >(tee -a "$LOG_FILE") 2>&1

echo "=== SNISID DR RUNBOOK: PostgreSQL Failover PAP → CAP ==="
echo "Date: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo "Opérateur: $USER"
echo ""

# ÉTAPE 0: Autorisation
echo "[ÉTAPE 0] Vérification autorisation..."
echo "⚠️  Ce runbook nécessite autorisation DG SNISID ou CTO."
echo "Entrez le code d'autorisation DR (format: SNISID-DR-YYYYMMDD-XXX):"
read -r AUTH_CODE
if [[ ! "$AUTH_CODE" =~ ^SNISID-DR-[0-9]{8}-[A-Z0-9]{3}$ ]]; then
    echo "❌ Code d'autorisation invalide. Arrêt."
    exit 1
fi
echo "✅ Code d'autorisation: $AUTH_CODE"

# ÉTAPE 1: Vérification état actuel
echo ""
echo "[ÉTAPE 1] Vérification état Patroni..."
patronictl -c /etc/patroni/patroni.yml list

echo ""
echo "Lag de réplication actuel:"
psql -h 10.130.0.10 -U monitor -c "SELECT now() - pg_last_xact_replay_timestamp() AS replication_lag;"

# ÉTAPE 2: Arrêt gracieux du trafic vers PAP
echo ""
echo "[ÉTAPE 2] Réduction trafic vers PAP dans GSLB..."
# Modifier F5 GSLB — poids PAP → 0
curl -sk -u admin:${F5_API_PASS} \
  -X PATCH "https://f5-pap.snisid.internal/mgmt/tm/gtm/pool/members" \
  -H "Content-Type: application/json" \
  -d '{"ratio": 0, "member": "POOL_PAP:VS_PAP_EXT"}'
echo "✅ Trafic PAP réduit à 0%"

sleep 30  # Attendre vidange des connexions actives

# ÉTAPE 3: Promotion PostgreSQL CAP
echo ""
echo "[ÉTAPE 3] Promotion PostgreSQL Standby CAP en Primary..."
patronictl -c /etc/patroni/patroni-cap.yml failover snisid-identity-cluster \
  --master cap-postgres-standby \
  --force

echo "Attente 60s pour stabilisation..."
sleep 60

patronictl -c /etc/patroni/patroni-cap.yml list
echo "✅ Vérification: CAP Standby promu Primary"

# ÉTAPE 4: Mise à jour DNS interne
echo ""
echo "[ÉTAPE 4] Mise à jour DNS interne..."
nsupdate -k /etc/bind/tsig.key << EOF
server ns1-int.snisid.gouv.ht
zone snisid.internal.
update delete postgres-rw.snisid.internal. A
update add postgres-rw.snisid.internal. 30 A 10.130.0.10
send
EOF
echo "✅ DNS postgres-rw → 10.130.0.10 (CAP)"

# ÉTAPE 5: Déblocage Vault CAP si nécessaire
echo ""
echo "[ÉTAPE 5] Vérification état Vault CAP..."
VAULT_STATUS=$(vault status -address=https://vault-cap.snisid.internal:8200 -format=json | jq -r .sealed)
if [ "$VAULT_STATUS" == "true" ]; then
    echo "⚠️  Vault CAP scellé — Unseal requis (5/9 key shards nécessaires)"
    echo "Contacter les détenteurs de clés (voir procédure Key Ceremony)"
    # vault operator unseal -address=https://vault-cap:8200 <SHARD_1>
    # vault operator unseal -address=https://vault-cap:8200 <SHARD_2>
    # ... (5 fois)
else
    echo "✅ Vault CAP opérationnel"
fi

# ÉTAPE 6: Vérification Kubernetes CAP
echo ""
echo "[ÉTAPE 6] Vérification pods Kubernetes CAP..."
kubectl --context=snisid-cap get pods -A -l app.kubernetes.io/part-of=snisid \
  --field-selector=status.phase!=Running

FAILED_PODS=$(kubectl --context=snisid-cap get pods -A -l app.kubernetes.io/part-of=snisid \
  --field-selector=status.phase!=Running --no-headers | wc -l)

if [ "$FAILED_PODS" -gt 0 ]; then
    echo "⚠️  $FAILED_PODS pods défaillants. Restart en cours..."
    kubectl --context=snisid-cap rollout restart deployment \
      -n snisid-identity identity-service
    kubectl --context=snisid-cap rollout restart deployment \
      -n snisid-auth auth-service
    kubectl --context=snisid-cap rollout status deployment identity-service \
      -n snisid-identity --timeout=120s
fi

# ÉTAPE 7: Tests fonctionnels
echo ""
echo "[ÉTAPE 7] Tests fonctionnels de validation..."
HTTP_STATUS=$(curl -sk -o /dev/null -w "%{http_code}" \
  https://api.snisid.gouv.ht/v1/health)
if [ "$HTTP_STATUS" == "200" ]; then
    echo "✅ API Gateway: OK ($HTTP_STATUS)"
else
    echo "❌ API Gateway: ÉCHEC ($HTTP_STATUS)"
    exit 1
fi

IDENTITY_STATUS=$(curl -sk -o /dev/null -w "%{http_code}" \
  -H "Authorization: Bearer $TEST_TOKEN" \
  https://api.snisid.gouv.ht/v1/identity/health)
echo "Identity Service health: $IDENTITY_STATUS"

# ÉTAPE 8: GSLB Update final
echo ""
echo "[ÉTAPE 8] Mise à jour GSLB finale — 100% vers CAP..."
curl -sk -u admin:${F5_API_PASS} \
  -X PATCH "https://f5-pap.snisid.internal/mgmt/tm/gtm/pool" \
  -d '{"members": [{"name": "POOL_CAP", "ratio": 100}]}'
echo "✅ 100% trafic vers CAP"

# ÉTAPE 9: Notification
echo ""
echo "[ÉTAPE 9] Notifications..."
# PagerDuty — DR activé
curl -s -X POST https://events.pagerduty.com/v2/enqueue \
  -H 'Content-Type: application/json' \
  -d "{\"routing_key\": \"$PD_ROUTING_KEY\",
       \"event_action\": \"trigger\",
       \"payload\": {
         \"summary\": \"SNISID DR ACTIVÉ: Basculement PAP→CAP complet\",
         \"severity\": \"critical\",
         \"source\": \"snisid-dr-runbook\",
         \"custom_details\": {\"auth_code\": \"$AUTH_CODE\", \"operator\": \"$USER\"}
       }}"

echo ""
echo "=== BASCULEMENT TERMINÉ ==="
echo "DC Actif: Cap-Haïtien"
echo "PostgreSQL Primary: 10.130.0.10"
echo "API: https://api.snisid.gouv.ht → VIP CAP"
echo "Log: $LOG_FILE"
echo ""
echo "⚠️  ACTIONS REQUISES:"
echo "  1. Monitoring renforcé actif (alertes ×2)"
echo "  2. Informer toutes les agences partenaires"
echo "  3. Planifier retour vers PAP quand infrastructure réparée"
echo "  4. Ouvrir incident post-mortem dans JIRA/Confluence"
```

### 7.3 Procédure Retour en Service PAP (Failback)

```mermaid
flowchart TD
    A([DC PAP Rétabli]) --> B[Validation infrastructure PAP:\n- Réseau opérationnel\n- Électricité stable > 2h\n- Tests PostgreSQL PAP]

    B --> C[Resynchronisation données:\nMise à jour PostgreSQL PAP\ndepuis Primary CAP\nDurée estimée: 1-4h selon volume WAL]

    C --> D{Lag < 1 Mo?}
    D -->|Non| E[Attendre resynchronisation\nMonitorer: pg_stat_replication]
    D -->|Oui| F[Fenêtre de maintenance\n00:00-04:00 HT\nAvis aux agences: 48h avant]

    E --> D

    F --> G[Basculement progressif:\nGSLB PAP: 10% → 25% → 50% → 60%\nPar tranches de 15 minutes]

    G --> H{Métriques OK à chaque palier?}
    H -->|Non| I[Rollback: 100% CAP\nAnalyse logs]
    H -->|Oui| J[Configuration finale:\nPAP=60%, CAP=40%\nPostgreSQL PAP → PRIMARY\nPostgreSQL CAP → Standby]

    J --> K[Post-Failback:\nRapport incident\nMise à jour runbook\nDebriefing équipe]

    style A fill:#2d6a2d,color:#fff
    style J fill:#1a4a8a,color:#fff
```

---

## 8. Architecture DR Offshore

### 8.1 Concept — Coffre-fort Numérique Souverain

```mermaid
graph TB
    subgraph "Haïti — Infrastructure Primaire"
        PAP[DC Port-au-Prince]
        CAP[DC Cap-Haïtien]
        PAP <-->|Active-Active| CAP
    end

    subgraph "DR Offshore — Islande ou Suisse"
        direction TB
        OFFSHORE_DC["Centre de Données Certifié\nISO 27001 + TIER III\nJuridiction favorable\nNeutre diplomatiquement"]

        subgraph "Données Chiffrées"
            BACKUP_COLD["Sauvegardes Froides\nChiffrées AES-256-GCM\nClés conservées EN HAÏTI\nAucun accès sans autorisation haïtienne"]
            WAL_ARCHIVE["Archives WAL\nPostgreSQL PITR\n30 jours minimum"]
            SCHEMA_BACKUP["Sauvegardes Schémas\nConfigurations Kubernetes\nCharts Helm + ArgoCD manifests"]
            PKI_BACKUP["Sauvegarde PKI\nRoot CA Offline copy\nHSM key material backup\nChiffré Shamir 5/9"]
        end

        OFFSHORE_K8S["Infrastructure K8s\nCapacité réduite (30%)\nActivée seulement en catastrophe"]

        OFFSHORE_OPS["Équipe Opérations DR\n2 ingénieurs haïtiens\nAccordés par l'État haïtien\nAccès biométrique + HSM token"]
    end

    subgraph "Contrôles Souveraineté Offshore"
        KEY_VAULT["Clés de déchiffrement\nConservées EN HAÏTI\nJamais exportées offshore\nHashiCorp Vault HSM-backed"]
        LEGAL_FW["Cadre Légal\nAccord bilatéral IS/CH — HT\nDroit haïtien applicable\nAudit annuel"]
        ACCESS_CTRL["Contrôle Accès\nDeux personnes autorisées haïtiennes\nAccès physique et logique\nEnregistrement vidéo"]
    end

    PAP -->|"Backup chiffré quotidien\nAES-256-GCM + signature ECDSA"| BACKUP_COLD
    PAP -->|"WAL streaming async"| WAL_ARCHIVE
    CAP -->|"Backup chiffré quotidien"| BACKUP_COLD

    KEY_VAULT -.->|"Clés restent en Haïti\nOffshore n'a JAMAIS les clés"| OFFSHORE_DC
    LEGAL_FW -.->|"Accord de souveraineté"| OFFSHORE_DC
    ACCESS_CTRL -.->|"Contrôle accès physique"| OFFSHORE_DC

    style KEY_VAULT fill:#8B0000,color:#fff
    style LEGAL_FW fill:#4a4a8a,color:#fff
```

### 8.2 Procédure Activation DR Offshore

```yaml
procedure_activation_dr_offshore:
  declencheur: "Catastrophe nationale — Double perte DC PAP + CAP"
  autorisation_requise:
    - "Décision Conseil des Ministres OU"
    - "Autorisation Premier Ministre + Ministre MTIC OU"
    - "Protocole urgence: DG SNISID + CISO + Architecte en Chef (3/3)"

  etapes:
    "00:00 - Décision":
      - "Activation du Plan de Continuité d'Activité National"
      - "Notification équipe offshore (2 ingénieurs)"
      - "Préparation des clés Shamir (5/9 détenteurs)"

    "00:30 - Transport clés":
      - "Transport sécurisé des shards Shamir par valise diplomatique ou personne de confiance"
      - "Vérification identité biométrique à l'arrivée offshore"
      - "Reconstruction clé maîtresse: reconstruct_key(shards[5..9])"

    "01:00 - Déchiffrement backup":
      - "Identification du dernier backup valide et complet"
      - "Déchiffrement: gpg --decrypt backup_YYYYMMDD.tar.gz.gpg"
      - "Vérification checksum SHA-512"
      - "Restauration PostgreSQL depuis backup + WAL PITR"

    "02:00 - Activation infrastructure":
      - "Démarrage cluster Kubernetes offshore (capacité réduite)"
      - "Restauration Vault depuis PKI backup"
      - "Initialisation nouveaux certificats TLS (domaine de crise)"
      - "Configuration DNS de crise: api-dr.snisid.gouv.ht"

    "03:00 - Services essentiels":
      - "Activation: Identity Service (read-only), Auth Service, API Gateway"
      - "Services désactivés: Biometric capture (équipements non disponibles)"
      - "Mode dégradé: vérification identité par NIN uniquement"

    "04:00 - Opérationnel":
      - "Services essentiels fonctionnels"
      - "Communication officielle: portail gouvernemental, radio nationale"
      - "SLA dégradé: RTO 4h, RPO 24h maximum"

  services_en_mode_degrade:
    actifs: [identity_read, auth, api_gateway, citizen_portal_read_only]
    inactifs: [biometric_capture, enrollment_new, document_generation]
    mode: "Vérification identité existante uniquement, pas de nouveaux enrôlements"

  duree_max_operation_offshore: "6 mois"
  retour_haiti:
    condition: "Infrastructure haïtienne rétablie et testée"
    procedure: "Resync données offshore → Haïti, validation, basculement progressif"
```

---

## 9. Procédures de Drill DR

### 9.1 Calendrier des Exercices

```mermaid
gantt
    title Calendrier Drills DR SNISID — Annuel
    dateFormat MM-DD
    section Niveau 1 — Composants
    Test UPS + Générateurs PAP           :01-15, 1d
    Test UPS + Générateurs CAP           :01-22, 1d
    Test basculement Redis               :02-05, 1d
    Test failover Vault                  :02-19, 1d
    Test backup PostgreSQL (restore)     :03-05, 1d
    Test failover PostgreSQL Patroni     :03-19, 1d

    section Niveau 2 — Cluster
    DR Drill Kubernetes PAP (panne 2 workers) :04-09, 1d
    DR Drill Kafka (perte 2 brokers)          :04-23, 1d
    DR Drill réseau VLAN (panne inter-zones)  :05-14, 1d

    section Niveau 3 — Inter-DC
    DR Drill PAP→CAP (basculement complet)    :06-11, 2d
    Post-drill review                         :06-13, 1d
    Correctifs identifiés                     :06-14, 14d

    section Niveau 4 — DR Offshore (Simulé)
    Simulation DR offshore (tabletop)         :09-10, 1d
    Test restauration backup offshore         :09-17, 1d

    section Niveau 5 — Full DR
    DR Drill complet PAP→CAP + offshore sim   :11-12, 3d
    Post-drill report                         :11-15, 1d
    Bilan annuel DR                           :12-03, 1d
```

### 9.2 Checklist Drill DR — Basculement PAP → CAP

```markdown
# CHECKLIST DRILL DR — Basculement PAP → CAP
# Version: 1.0 | Date: [DATE_DRILL] | Référence: DRILL-[YYYY-NNN]

## PRÉ-DRILL (J-7)
- [ ] Notification parties prenantes (agences, direction)
- [ ] Revue procédures runbooks — dernière mise à jour < 30 jours
- [ ] Validation configuration Patroni et MirrorMaker 2
- [ ] Vérification disponibilité équipe (au moins 4 ingénieurs)
- [ ] Préparation environnement monitoring renforcé
- [ ] Backup état complet avant drill
- [ ] Confirmation fenêtre de maintenance (impact réduit)

## EXÉCUTION (Jour du Drill)

### Phase 1 — Simulation panne PAP (T+0)
- [ ] Isolation réseau simulée PAP (blocage VIPs au niveau F5)
- [ ] Observation comportement automatique Patroni (< 30s)
- [ ] Chronométrage: détection panne → promotion PostgreSQL CAP
- [ ] Vérification: Vault CAP opérationnel
- [ ] Vérification: Kafka MirrorMaker 2 (lag monitoring)

### Phase 2 — Basculement Applications (T+5min)
- [ ] Mise à jour GSLB → 100% CAP
- [ ] Vérification pods Kubernetes CAP: tous Running
- [ ] Test API health: GET /v1/health → 200
- [ ] Test identity verification: POST /v1/identity/verify → 200
- [ ] Test authentication: POST /auth/token → 200
- [ ] Chronométrage: T+0 → Services opérationnels

### Phase 3 — Validation Données (T+15min)
- [ ] Vérification intégrité: checksum PostgreSQL CAP vs dernier commit PAP
- [ ] Test Kafka: vérification que aucun message perdu (offsets)
- [ ] Test Vault: rotation d'un secret test réussie
- [ ] Vérification audit trail: tous événements du drill enregistrés

### Phase 4 — Retour Nominal (T+60min)
- [ ] Rétablissement PAP simulé (déblocage réseau)
- [ ] Resynchronisation PostgreSQL (monitoring lag)
- [ ] Test GSLB progressif: 10% PAP → 25% → 60%
- [ ] Validation métriques à chaque palier
- [ ] Retour configuration nominale 60/40

## POST-DRILL

### Métriques à Capturer
- [ ] RTO effectif: __ minutes (objectif: < 5 min)
- [ ] RPO effectif: __ secondes de données perdues (objectif: 0)
- [ ] Disponibilité services pendant basculement: __%
- [ ] Nombre d'erreurs 5xx pendant basculement: __
- [ ] Temps resynchronisation données: __ minutes

### Rapport de Drill
- [ ] Rapport écrit dans les 48h
- [ ] Liste des problèmes identifiés
- [ ] Tickets correctifs créés
- [ ] Mise à jour runbooks si nécessaire
- [ ] Présentation au COMEX si RTO > objectif

## CRITÈRES DE SUCCÈS
| Métrique | Objectif | Réel | Statut |
|---|---|---|---|
| RTO total | < 5 minutes | __ | |
| RPO | 0 | __ | |
| Erreurs 5xx | < 50 | __ | |
| Kafka lag fin basculement | < 1000 msgs | __ | |
| Services opérationnels | 100% | __ | |
```

---

## 10. Monitoring HA/DR

### 10.1 Tableau de Bord HA/DR — Métriques Clés

```yaml
monitoring_hadr:
  dashboards:
    dr_status:
      panels:
        - title: "Statut DC PAP"
          query: "up{job='snisid-pap-health'} == 1"
          alert: "if absent > 30s → PagerDuty CRITICAL"

        - title: "Statut DC CAP"
          query: "up{job='snisid-cap-health'} == 1"
          alert: "if absent > 30s → PagerDuty CRITICAL"

        - title: "PostgreSQL Replication Lag"
          query: "pg_replication_lag_seconds{cluster='identity'}"
          thresholds:
            warning: 10   # secondes
            critical: 60  # secondes
            alert_action: "PagerDuty HIGH si > 10s"

        - title: "Kafka MirrorMaker Lag"
          query: "kafka_mirrormaker_record_lag"
          thresholds:
            warning: 10000
            critical: 100000

        - title: "Vault HA Status"
          query: "vault_core_active"
          description: "1=active, 0=standby"

        - title: "Patroni Leader PAP"
          query: "patroni_master{cluster='snisid-identity-cluster', dc='pap'}"
          alert: "if 0 for > 60s → auto-failover check"

        - title: "RTO SLA Conformance"
          query: "snisid_failover_duration_seconds"
          sla: "< 300s (5 minutes)"
          review: "mensuel"

  alerting_rules:
    - name: "DR_THRESHOLD_CROSSED"
      condition: "replication_lag > 60s OR kafka_lag > 100000"
      action: "Page DBA + Infrastructure on-call"
      escalation: "Si non résolu en 15min → CTO"

    - name: "DC_DOWN_COMPLETE"
      condition: "all health checks fail on DC for 60s"
      action: "AUTO: Patroni failover\nMANUAL: Run dr-failover runbook"
      escalation: "Immédiat DG + CISO + CTO"

    - name: "VAULT_SEALED"
      condition: "vault_seal_status == 1"
      action: "Page sécurité on-call — unseal requis"
      priority: CRITICAL

  on_call_rotation:
    primary: "DBA + Infrastructure — rotation hebdomadaire"
    secondary: "Architecture team — escalade"
    management: "CTO/DG — incidents severity 1"
    contact: "PagerDuty + WhatsApp sécurisé Signal"
```

---

## Bloc d'Approbation / Approval Block

| Rôle | Nom | Signature | Date |
|---|---|---|---|
| **Architecte en Chef** | [À compléter] | [Signature] | 2026-05-25 |
| **Directeur Infrastructure** | [À compléter] | [Signature] | 2026-05-25 |
| **CISO** | [À compléter] | [Signature] | 2026-05-25 |
| **Directeur Général SNISID** | [À compléter] | [Signature] | 2026-05-25 |
| **Ministère Responsable** | [À compléter] | [Signature] | 2026-05-25 |

---

*Document SNISID-ARC-HA-DR-001 v1.0.0 — CONFIDENTIEL — © République d'Haïti, Programme SNISID, 2026*
*Révision obligatoire trimestrielle. Tests DR semi-annuels obligatoires.*
