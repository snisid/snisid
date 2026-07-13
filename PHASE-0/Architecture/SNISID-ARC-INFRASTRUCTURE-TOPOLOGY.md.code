# SNISID — Topologie d'Infrastructure Complète
# SNISID — Complete Infrastructure Topology

---

| Métadonnée | Valeur |
|---|---|
| **Document ID** | SNISID-ARC-INF-001 |
| **Version** | 1.0.0 |
| **Date** | 2026-05-25 |
| **Statut** | APPROUVÉ — Production |
| **Classification** | CONFIDENTIEL / CONFIDENTIAL |
| **Auteur** | Infrastructure Team — SNISID Programme |
| **Révisé par** | Chief Architect, Network Lead, Security Architect |
| **Approuvé par** | Directeur Infrastructure, DG SNISID |
| **Standard** | ISO/IEC 27001:2022, NIST SP 800-53, TIA-942 Tier III |

---

## Table des Matières

1. [Vue Générale de la Topologie](#1-vue-générale-de-la-topologie)
2. [Topologie Réseau Complète](#2-topologie-réseau-complète)
3. [Topologie Kubernetes](#3-topologie-kubernetes)
4. [Service Mesh — Istio](#4-service-mesh--istio)
5. [Topologie Bases de Données](#5-topologie-bases-de-données)
6. [Topologie Apache Kafka](#6-topologie-apache-kafka)
7. [Topologie Stockage](#7-topologie-stockage)
8. [Architecture DNS](#8-architecture-dns)
9. [Load Balancers](#9-load-balancers)
10. [Zones Pare-feu et Règles](#10-zones-pare-feu-et-règles)
11. [Plan d'Adressage IP](#11-plan-dadressage-ip)
12. [Topologie Physique des Datacenters](#12-topologie-physique-des-datacenters)

---

## 1. Vue Générale de la Topologie

### 1.1 Macro-Topologie SNISID

```mermaid
graph TB
    subgraph "Internet / WAN"
        INT[Internet mondial]
        CARICOM[CARICOM Region Network]
        DIASPORA[Diaspora Users]
    end

    subgraph "Périmètre HAïtien — DMZ Externe"
        direction TB
        CDN[CDN / WAF — Cloudflare Souverain]
        DDOS[Anti-DDoS Scrubbing Center]
        EXT_LB[External Load Balancer — F5 BIG-IP]
        BORDER_FW[Border Firewall — Fortinet FortiGate 6000F]
    end

    subgraph "DC Primaire — Port-au-Prince"
        direction TB
        subgraph "DMZ-PAP"
            API_GW_PAP[API Gateway Cluster — PAP]
            JUMP_PAP[Jump Server PAP]
        end
        subgraph "Application Zone PAP"
            K8S_PAP[Kubernetes Cluster PAP — 12 Worker Nodes]
        end
        subgraph "Data Zone PAP"
            DB_PAP[PostgreSQL HA Cluster — PAP]
            KAFKA_PAP[Kafka Cluster — PAP — 5 Brokers]
            VAULT_PAP[HashiCorp Vault — PAP]
            HSM_PAP[HSM Luna Network 7 — PAP]
        end
        subgraph "Storage Zone PAP"
            CEPH_PAP[Ceph Cluster — PAP — 6 OSD Nodes]
            MINIO_PAP[MinIO Distributed — PAP]
        end
        INT_FW_PAP[Internal Firewall — PAP]
    end

    subgraph "DC Secondaire — Cap-Haïtien"
        direction TB
        subgraph "DMZ-CAP"
            API_GW_CAP[API Gateway Cluster — CAP]
        end
        subgraph "Application Zone CAP"
            K8S_CAP[Kubernetes Cluster CAP — 8 Worker Nodes]
        end
        subgraph "Data Zone CAP"
            DB_CAP[PostgreSQL HA Cluster — CAP]
            KAFKA_CAP[Kafka Cluster — CAP — 3 Brokers]
            VAULT_CAP[HashiCorp Vault — CAP]
            HSM_CAP[HSM Luna Network 7 — CAP]
        end
        subgraph "Storage Zone CAP"
            CEPH_CAP[Ceph Cluster — CAP — 4 OSD Nodes]
        end
        INT_FW_CAP[Internal Firewall — CAP]
    end

    subgraph "DR Offshore — Coffre-fort Numérique"
        DR_VAULT[DR Vault — Islande/Suisse]
        DR_BACKUP[Backup Chiffré Hors-ligne]
    end

    subgraph "Sites Terrain — Réseau Gouvernemental"
        OEC_SITES[Sites OEC — 145 Communes]
        MEL_SITES[Sites MEL — Bureaux Électoraux]
        POLICE_SITES[PNH — Postes de Police]
    end

    INT --> DDOS --> CDN
    DIASPORA --> CDN
    CDN --> BORDER_FW
    BORDER_FW --> EXT_LB
    EXT_LB --> API_GW_PAP
    EXT_LB --> API_GW_CAP

    API_GW_PAP --> INT_FW_PAP --> K8S_PAP
    K8S_PAP --> DB_PAP
    K8S_PAP --> KAFKA_PAP
    K8S_PAP --> VAULT_PAP
    VAULT_PAP --> HSM_PAP
    K8S_PAP --> CEPH_PAP
    CEPH_PAP --> MINIO_PAP

    API_GW_CAP --> INT_FW_CAP --> K8S_CAP
    K8S_CAP --> DB_CAP
    K8S_CAP --> KAFKA_CAP
    K8S_CAP --> VAULT_CAP
    VAULT_CAP --> HSM_CAP

    DB_PAP <-->|"Streaming réplication synchrone"| DB_CAP
    KAFKA_PAP <-->|"MirrorMaker 2"| KAFKA_CAP
    VAULT_PAP <-->|"Vault Replication Enterprise"| VAULT_CAP

    DB_PAP -->|"Backup chiffré quotidien"| DR_VAULT
    DB_CAP -->|"Backup chiffré quotidien"| DR_VAULT

    OEC_SITES <-->|"VPN IPsec / MPLS"| API_GW_PAP
    MEL_SITES <-->|"VPN IPsec / MPLS"| API_GW_PAP
    POLICE_SITES <-->|"VPN IPsec / MPLS"| API_GW_PAP
    CARICOM <-->|"Interconnexion sécurisée"| BORDER_FW
```

---

## 2. Topologie Réseau Complète

### 2.1 Diagramme Réseau Détaillé — DC Port-au-Prince

```mermaid
graph TB
    subgraph "EDGE — Niveau Périmètre"
        direction LR
        ISP1[ISP 1 — Natcom Fibre 10Gbps]
        ISP2[ISP 2 — Digicel Business 10Gbps]
        ISP3[ISP 3 — Link-Up 1Gbps Backup]
        BGP_ROUTER[BGP Router — Cisco ASR 1001-X Dual]
    end

    subgraph "CORE — Niveau Agrégation"
        CORE_SW1[Core Switch 1 — Arista 7280 100G MLAG Primary]
        CORE_SW2[Core Switch 2 — Arista 7280 100G MLAG Secondary]
        FORTIGATE_PRIMARY[FortiGate 6000F — Primary FGCP HA]
        FORTIGATE_SECONDARY[FortiGate 6000F — Secondary FGCP HA]
    end

    subgraph "VLAN 10 — DMZ Externe"
        direction LR
        F5_EXT_1[F5 BIG-IP i4000 — VIP Externe Primary]
        F5_EXT_2[F5 BIG-IP i4000 — VIP Externe Secondary]
        KONG_1[Kong Gateway Node 1 — 16vCPU/32GB]
        KONG_2[Kong Gateway Node 2 — 16vCPU/32GB]
        KONG_3[Kong Gateway Node 3 — 16vCPU/32GB]
        IDS_IPS[Suricata IDS/IPS — Inline Mode]
    end

    subgraph "VLAN 20 — Application"
        K8S_CP1[K8s Control Plane 1 — 8vCPU/32GB]
        K8S_CP2[K8s Control Plane 2 — 8vCPU/32GB]
        K8S_CP3[K8s Control Plane 3 — 8vCPU/32GB]
        K8S_W1[K8s Worker 1 — 32vCPU/128GB]
        K8S_W2[K8s Worker 2 — 32vCPU/128GB]
        K8S_W3[K8s Worker 3 — 32vCPU/128GB]
        K8S_W4[K8s Worker 4 — 32vCPU/128GB — GPU]
        K8S_W5[K8s Worker 5 — 32vCPU/128GB — GPU]
        K8S_W6[K8s Worker 6 — 16vCPU/64GB]
        K8S_W7[K8s Worker 7 — 16vCPU/64GB]
        K8S_W8[K8s Worker 8 — 16vCPU/64GB — Biometric]
        K8S_W9[K8s Worker 9 — 16vCPU/64GB — Biometric]
    end

    subgraph "VLAN 30 — Data Tier"
        PG_PRIMARY[PostgreSQL Primary — 32vCPU/256GB/4TB NVMe]
        PG_REPLICA1[PostgreSQL Replica 1 — 32vCPU/256GB/4TB NVMe]
        PG_REPLICA2[PostgreSQL Replica 2 — 16vCPU/128GB/2TB NVMe]
        REDIS_1[Redis Master — 16vCPU/64GB]
        REDIS_2[Redis Replica 1 — 16vCPU/64GB]
        REDIS_3[Redis Replica 2 — 8vCPU/32GB]
    end

    subgraph "VLAN 40 — Kafka / Messaging"
        KAFKA_B1[Kafka Broker 1 — 16vCPU/64GB/10TB SSD]
        KAFKA_B2[Kafka Broker 2 — 16vCPU/64GB/10TB SSD]
        KAFKA_B3[Kafka Broker 3 — 16vCPU/64GB/10TB SSD]
        KAFKA_B4[Kafka Broker 4 — 16vCPU/64GB/10TB SSD]
        KAFKA_B5[Kafka Broker 5 — 16vCPU/64GB/10TB SSD]
    end

    subgraph "VLAN 50 — Storage"
        CEPH_MON1[Ceph Monitor 1]
        CEPH_MON2[Ceph Monitor 2]
        CEPH_MON3[Ceph Monitor 3]
        CEPH_OSD1[Ceph OSD 1 — 4x16TB SAS]
        CEPH_OSD2[Ceph OSD 2 — 4x16TB SAS]
        CEPH_OSD3[Ceph OSD 3 — 4x16TB SAS]
        CEPH_OSD4[Ceph OSD 4 — 4x16TB SAS]
        CEPH_OSD5[Ceph OSD 5 — 4x16TB SAS]
        CEPH_OSD6[Ceph OSD 6 — 4x16TB SAS]
    end

    subgraph "VLAN 60 — Sécurité / HSM"
        VAULT_1[Vault Node 1 — 8vCPU/32GB]
        VAULT_2[Vault Node 2 — 8vCPU/32GB]
        VAULT_3[Vault Node 3 — 8vCPU/32GB]
        HSM_1[HSM Luna 7 Network — Primary]
        HSM_2[HSM Luna 7 Network — Backup]
    end

    subgraph "VLAN 70 — Management / OOB"
        BASTION[Bastion Host — Jump Server]
        MGMT_SW[Management Switch — OOB]
        IPMI[iDRAC/iLO — Out-of-Band Management]
    end

    ISP1 --> BGP_ROUTER
    ISP2 --> BGP_ROUTER
    ISP3 --> BGP_ROUTER
    BGP_ROUTER --> CORE_SW1
    BGP_ROUTER --> CORE_SW2
    CORE_SW1 <-->|MLAG| CORE_SW2
    CORE_SW1 --> FORTIGATE_PRIMARY
    CORE_SW2 --> FORTIGATE_SECONDARY
    FORTIGATE_PRIMARY <-->|FGCP HA| FORTIGATE_SECONDARY
    FORTIGATE_PRIMARY --> F5_EXT_1
    FORTIGATE_SECONDARY --> F5_EXT_2
    F5_EXT_1 --> KONG_1
    F5_EXT_1 --> KONG_2
    F5_EXT_2 --> KONG_3
    KONG_1 --> K8S_W1
    KONG_2 --> K8S_W2
    KONG_3 --> K8S_W3
```

### 2.2 Interconnexion Inter-Datacenter

```mermaid
graph LR
    subgraph "DC Port-au-Prince"
        PAP_CORE[Core Network PAP]
        PAP_WAN[WAN Edge PAP]
    end

    subgraph "Liaisons WAN"
        direction TB
        MPLS_PRIMARY[MPLS Privé — 10Gbps — Natcom]
        MPLS_SECONDARY[MPLS Privé — 10Gbps — Digicel]
        DARK_FIBER[Dark Fiber — 100Gbps — TBD]
        SATELLITE[VSAT Backup — Starlink Enterprise 500Mbps]
    end

    subgraph "DC Cap-Haïtien"
        CAP_WAN[WAN Edge CAP]
        CAP_CORE[Core Network CAP]
    end

    PAP_WAN <-->|Primary| MPLS_PRIMARY
    PAP_WAN <-->|Secondary| MPLS_SECONDARY
    PAP_WAN <-->|Tertiary Backup| SATELLITE
    MPLS_PRIMARY <-->|eBGP| CAP_WAN
    MPLS_SECONDARY <-->|eBGP Failover| CAP_WAN
    SATELLITE <-->|Emergency| CAP_WAN
    PAP_CORE <-->|Future — 100G| DARK_FIBER
    DARK_FIBER <-->|Future| CAP_CORE

    style MPLS_PRIMARY fill:#2d7a2d,color:#fff
    style MPLS_SECONDARY fill:#d4a017,color:#000
    style SATELLITE fill:#cc4444,color:#fff
```

---

## 3. Topologie Kubernetes

### 3.1 Architecture du Cluster Kubernetes PAP

```mermaid
graph TB
    subgraph "Control Plane — Haute Disponibilité"
        direction LR
        CP1["Control Plane 1\nkube-apiserver\netcd leader\nkube-scheduler\nkube-controller-manager\n8vCPU / 32GB"]
        CP2["Control Plane 2\nkube-apiserver\netcd follower\nkube-scheduler (standby)\nkube-controller-manager (standby)\n8vCPU / 32GB"]
        CP3["Control Plane 3\nkube-apiserver\netcd follower\n8vCPU / 32GB"]
        ETCD_LB[etcd Load Balancer — HAProxy]
        API_LB[kube-apiserver VIP — keepalived]
    end

    subgraph "Node Pool: General Purpose"
        W1["Worker-GP-01\n32vCPU / 128GB\nLabels: pool=general"]
        W2["Worker-GP-02\n32vCPU / 128GB\nLabels: pool=general"]
        W3["Worker-GP-03\n32vCPU / 128GB\nLabels: pool=general"]
        W4["Worker-GP-04\n16vCPU / 64GB\nLabels: pool=general"]
    end

    subgraph "Node Pool: Biometric / GPU"
        W5["Worker-BIO-01\n32vCPU / 128GB\nNVIDIA A10 GPU\nLabels: pool=biometric, gpu=true"]
        W6["Worker-BIO-02\n32vCPU / 128GB\nNVIDIA A10 GPU\nLabels: pool=biometric, gpu=true"]
    end

    subgraph "Node Pool: Database (DaemonSet-like)"
        W7["Worker-DB-01\n32vCPU / 256GB\n4TB NVMe Local\nLabels: pool=database"]
        W8["Worker-DB-02\n32vCPU / 256GB\n4TB NVMe Local\nLabels: pool=database"]
        W9["Worker-DB-03\n16vCPU / 128GB\n2TB NVMe Local\nLabels: pool=database"]
    end

    subgraph "Node Pool: System / Infra"
        W10["Worker-SYS-01\n16vCPU / 64GB\nLabels: pool=system"]
        W11["Worker-SYS-02\n16vCPU / 64GB\nLabels: pool=system"]
        W12["Worker-SYS-03\n8vCPU / 32GB\nLabels: pool=system"]
    end

    subgraph "Add-ons Système (DaemonSets)"
        CNI[Cilium CNI — eBPF networking]
        CSI[Ceph CSI — Storage driver]
        METRICS[metrics-server]
        LOGGING[Fluentbit DaemonSet]
        OTEL[OpenTelemetry Collector]
    end

    CP1 <-->|Raft consensus| CP2
    CP2 <-->|Raft consensus| CP3
    CP1 <-->|Raft consensus| CP3
    ETCD_LB --> CP1
    ETCD_LB --> CP2
    ETCD_LB --> CP3
    API_LB --> CP1
    API_LB --> CP2
    API_LB --> CP3

    CP1 --> W1
    CP1 --> W2
    CP1 --> W5
    CP1 --> W7
    CP1 --> W10
```

### 3.2 Namespaces et Organisation Logique

```yaml
# Organisation des Namespaces Kubernetes SNISID
namespaces:
  production:
    - name: snisid-identity
      description: "Identity Service et composants associés"
      resource_quota:
        requests.cpu: "40"
        requests.memory: "160Gi"
        limits.cpu: "80"
        limits.memory: "320Gi"
      node_selector: {pool: general}
      labels: {tier: critical, data-class: personal, env: production}

    - name: snisid-biometric
      description: "Biometric Service — isolation maximale"
      resource_quota:
        requests.cpu: "20"
        requests.memory: "80Gi"
        limits.cpu: "40"
        limits.memory: "160Gi"
      node_selector: {pool: biometric}
      labels: {tier: critical, data-class: biometric, env: production}
      network_policy: "deny-all-ingress-except-gateway"

    - name: snisid-auth
      description: "Authentication Service, Keycloak"
      resource_quota:
        requests.cpu: "16"
        requests.memory: "64Gi"
      labels: {tier: critical, env: production}

    - name: snisid-enrollment
      description: "Enrollment orchestration"
      resource_quota:
        requests.cpu: "8"
        requests.memory: "32Gi"
      labels: {tier: high, env: production}

    - name: snisid-documents
      description: "Document generation service"
      resource_quota:
        requests.cpu: "8"
        requests.memory: "32Gi"
      labels: {tier: high, env: production}

    - name: snisid-gateway
      description: "API Gateway (Kong) — DMZ interne"
      resource_quota:
        requests.cpu: "12"
        requests.memory: "48Gi"
      labels: {tier: critical, env: production}

    - name: snisid-audit
      description: "Audit service — write-heavy"
      resource_quota:
        requests.cpu: "8"
        requests.memory: "32Gi"
      labels: {tier: critical, env: production}

    - name: snisid-interop
      description: "Interoperability gateway"
      labels: {tier: high, env: production}

    - name: snisid-notifications
      description: "SMS/Email notification service"
      labels: {tier: medium, env: production}

  infrastructure:
    - name: monitoring
      description: "Prometheus, Grafana, Loki, Tempo"
    - name: istio-system
      description: "Istio service mesh control plane"
    - name: cert-manager
      description: "Certificate automation"
    - name: vault-agent
      description: "HashiCorp Vault agent injectors"
    - name: kafka-system
      description: "Kafka Strimzi operator"
    - name: ceph-csi
      description: "Storage CSI driver"
    - name: argocd
      description: "GitOps deployment"
    - name: security-system
      description: "OPA Gatekeeper, Falco, Trivy Operator"
```

### 3.3 Configuration RKE2 Production

```yaml
# /etc/rancher/rke2/config.yaml — Control Plane
cluster-name: snisid-prod-pap
tls-san:
  - k8s-api.snisid.gouv.ht
  - 10.10.0.100  # VIP keepalived

cni: cilium

# FIPS 140-2 compliance
fips: true

# Disable cloud providers (on-premise)
cloud-provider-name: ""

# etcd snapshot
etcd-snapshot-schedule-cron: "0 */6 * * *"
etcd-snapshot-retention: 30
etcd-snapshot-dir: /var/lib/etcd-snapshots

# Audit logging
kube-apiserver-arg:
  - "--audit-log-path=/var/log/kubernetes/audit.log"
  - "--audit-log-maxage=30"
  - "--audit-log-maxbackup=10"
  - "--audit-log-maxsize=100"
  - "--audit-policy-file=/etc/kubernetes/audit-policy.yaml"
  - "--enable-admission-plugins=NodeRestriction,PodSecurityAdmission,AlwaysPullImages"
  - "--encryption-provider-config=/etc/kubernetes/encryption-config.yaml"
  - "--anonymous-auth=false"
  - "--tls-min-version=VersionTLS12"
  - "--tls-cipher-suites=TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
  - "--oidc-issuer-url=https://auth.snisid.gouv.ht/realms/snisid"
  - "--oidc-client-id=kubernetes"
  - "--oidc-username-claim=preferred_username"
  - "--oidc-groups-claim=groups"

kube-controller-manager-arg:
  - "--terminated-pod-gc-threshold=100"
  - "--use-service-account-credentials=true"

kube-scheduler-arg:
  - "--config=/etc/kubernetes/scheduler-config.yaml"

# Node labels for pools
node-label:
  - "datacenter=port-au-prince"
  - "site=snisid-dc-pap"
  - "tier=control-plane"
```

---

## 4. Service Mesh — Istio

### 4.1 Topologie Istio

```mermaid
graph TB
    subgraph "Istio Control Plane — istiod"
        PILOT[Pilot — Service Discovery & xDS]
        CITADEL[Citadel — Certificate Authority\nSPIFFE compatible]
        GALLEY[Galley — Config Validation]
        ISTIOD[istiod — Unified Control Plane Pod]
    end

    subgraph "Data Plane — Envoy Sidecars"
        subgraph "snisid-gateway namespace"
            GW_PROXY[Envoy Proxy\nIngress Gateway]
        end
        subgraph "snisid-identity namespace"
            ID_APP[Identity Service App]
            ID_PROXY[Envoy Sidecar\nmTLS auto-inject]
        end
        subgraph "snisid-biometric namespace"
            BIO_APP[Biometric Service App]
            BIO_PROXY[Envoy Sidecar\nmTLS auto-inject]
        end
        subgraph "snisid-auth namespace"
            AUTH_APP[Auth Service App]
            AUTH_PROXY[Envoy Sidecar\nmTLS auto-inject]
        end
        subgraph "snisid-audit namespace"
            AUDIT_APP[Audit Service App]
            AUDIT_PROXY[Envoy Sidecar\nmTLS auto-inject]
        end
    end

    subgraph "Observabilité Mesh"
        KIALI[Kiali — Service Graph]
        JAEGER[Jaeger — Distributed Tracing]
        PROM_ISTIO[Prometheus — Mesh Metrics]
    end

    ISTIOD --> PILOT
    ISTIOD --> CITADEL
    ISTIOD --> GALLEY

    PILOT -->|xDS config push| ID_PROXY
    PILOT -->|xDS config push| BIO_PROXY
    PILOT -->|xDS config push| AUTH_PROXY
    PILOT -->|xDS config push| AUDIT_PROXY
    PILOT -->|xDS config push| GW_PROXY

    CITADEL -->|SVID certificates| ID_PROXY
    CITADEL -->|SVID certificates| BIO_PROXY
    CITADEL -->|SVID certificates| AUTH_PROXY
    CITADEL -->|SVID certificates| AUDIT_PROXY

    ID_APP <-->|localhost| ID_PROXY
    BIO_APP <-->|localhost| BIO_PROXY
    AUTH_APP <-->|localhost| AUTH_PROXY
    AUDIT_APP <-->|localhost| AUDIT_PROXY

    GW_PROXY <-->|mTLS ISTIO_MUTUAL| ID_PROXY
    ID_PROXY <-->|mTLS ISTIO_MUTUAL| BIO_PROXY
    ID_PROXY <-->|mTLS ISTIO_MUTUAL| AUTH_PROXY
    ID_PROXY <-->|mTLS ISTIO_MUTUAL| AUDIT_PROXY

    ID_PROXY -->|traces| JAEGER
    BIO_PROXY -->|traces| JAEGER
    ID_PROXY -->|metrics| PROM_ISTIO
    KIALI --> PROM_ISTIO
```

### 4.2 PeerAuthentication — mTLS Strict Mode

```yaml
# PeerAuthentication — STRICT mTLS pour tous les namespaces SNISID
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: snisid-mtls-strict-global
  namespace: istio-system
spec:
  mtls:
    mode: STRICT  # Aucune communication en clair autorisée
---
# AuthorizationPolicy — Biometric Service: accès très restreint
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: biometric-service-authz
  namespace: snisid-biometric
spec:
  action: ALLOW
  rules:
    - from:
        - source:
            principals:
              - "cluster.local/ns/snisid-identity/sa/identity-service"
              - "cluster.local/ns/snisid-enrollment/sa/enrollment-service"
              - "cluster.local/ns/snisid-gateway/sa/api-gateway"
      to:
        - operation:
            methods: ["POST", "GET"]
            paths:
              - "/v1/biometric/capture/*"
              - "/v1/biometric/match/*"
              - "/v1/biometric/verify/*"
      when:
        - key: request.headers[x-purpose]
          values: ["enrollment", "verification", "emergency"]
---
# VirtualService — Canary deployment Identity Service
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: identity-service-vs
  namespace: snisid-identity
spec:
  hosts:
    - identity-service
  http:
    - name: canary-route
      match:
        - headers:
            x-canary:
              exact: "true"
      route:
        - destination:
            host: identity-service
            subset: v2-canary
          weight: 100
    - name: stable-route
      route:
        - destination:
            host: identity-service
            subset: v1-stable
          weight: 95
        - destination:
            host: identity-service
            subset: v2-canary
          weight: 5
```

---

## 5. Topologie Bases de Données

### 5.1 PostgreSQL Patroni — Haute Disponibilité

```mermaid
graph TB
    subgraph "PostgreSQL Identity Cluster — PAP"
        direction TB
        PG_PRIMARY_PAP["PostgreSQL Primary — PAP\n32vCPU / 256GB RAM\n4TB NVMe SSD\nWAL archiving: S3/MinIO\nIP: 10.30.0.10"]
        PG_REPLICA1_PAP["PostgreSQL Replica 1 — PAP\n32vCPU / 256GB RAM\n4TB NVMe SSD\nStreaming replication sync\nIP: 10.30.0.11"]
        PG_REPLICA2_PAP["PostgreSQL Replica 2 — PAP\n16vCPU / 128GB RAM\n2TB NVMe SSD\nStreaming replication async\nIP: 10.30.0.12"]

        PATRONI1[Patroni Agent 1]
        PATRONI2[Patroni Agent 2]
        PATRONI3[Patroni Agent 3]

        ETCD_CLUSTER["etcd Cluster (3 nodes)\nLeader Election\nDCS — Distributed Config Store"]

        PG_VIP_RW["Virtual IP RW\n10.30.0.5 — Primary\nHAProxy / keepalived"]
        PG_VIP_RO["Virtual IP RO\n10.30.0.6 — Replicas\nRead-only connections"]
    end

    subgraph "PostgreSQL Biometric Cluster — PAP (Isolé)"
        BIO_PRIMARY["Biometric Primary — PAP\n32vCPU / 256GB\n8TB NVMe\nChiffrement colonne AES-256\nIP: 10.30.1.10"]
        BIO_REPLICA["Biometric Replica — PAP\n32vCPU / 256GB\n8TB NVMe\nIP: 10.30.1.11"]
        BIO_VIP["VIP Biometric\n10.30.1.5"]
    end

    subgraph "PostgreSQL Identity Cluster — CAP (DR Actif)"
        PG_PRIMARY_CAP["PostgreSQL Primary — CAP\n32vCPU / 256GB\n4TB NVMe\nIP: 10.130.0.10"]
        PG_REPLICA1_CAP["PostgreSQL Replica 1 — CAP\n16vCPU / 128GB\n2TB NVMe\nIP: 10.130.0.11"]
    end

    subgraph "pgBouncer — Connection Pooling"
        PGBOUNCER_PAP1[pgBouncer Instance 1 — PAP\nMax: 5000 conn]
        PGBOUNCER_PAP2[pgBouncer Instance 2 — PAP\nMax: 5000 conn]
    end

    PG_PRIMARY_PAP -->|Sync streaming replication| PG_REPLICA1_PAP
    PG_PRIMARY_PAP -->|Async streaming replication| PG_REPLICA2_PAP

    PATRONI1 --> PG_PRIMARY_PAP
    PATRONI2 --> PG_REPLICA1_PAP
    PATRONI3 --> PG_REPLICA2_PAP

    PATRONI1 <-->|DCS| ETCD_CLUSTER
    PATRONI2 <-->|DCS| ETCD_CLUSTER
    PATRONI3 <-->|DCS| ETCD_CLUSTER

    PG_VIP_RW -->|Leader| PG_PRIMARY_PAP
    PG_VIP_RO -->|Round-robin| PG_REPLICA1_PAP
    PG_VIP_RO -->|Round-robin| PG_REPLICA2_PAP

    PGBOUNCER_PAP1 --> PG_VIP_RW
    PGBOUNCER_PAP2 --> PG_VIP_RW

    PG_PRIMARY_PAP -->|Logical replication async MPLS| PG_PRIMARY_CAP
    PG_PRIMARY_CAP -->|Streaming sync local| PG_REPLICA1_CAP

    BIO_PRIMARY -->|Sync replication| BIO_REPLICA
    BIO_VIP -->|Leader| BIO_PRIMARY
```

### 5.2 Configuration PostgreSQL Production

```sql
-- postgresql.conf — Configuration Production SNISID
-- Fichier: /etc/postgresql/16/main/postgresql.conf

-- CONNEXIONS ET AUTHENTIFICATION
max_connections = 200              -- pgBouncer gère le pooling
superuser_reserved_connections = 5
ssl = on
ssl_cert_file = '/etc/ssl/postgresql/server.crt'
ssl_key_file = '/etc/ssl/postgresql/server.key'
ssl_ca_file = '/etc/ssl/postgresql/ca.crt'
ssl_min_protocol_version = 'TLSv1.2'
ssl_ciphers = 'ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384'

-- MÉMOIRE ET PERFORMANCE
shared_buffers = 64GB              -- 25% RAM (256GB total)
effective_cache_size = 192GB       -- 75% RAM
work_mem = 256MB
maintenance_work_mem = 4GB
huge_pages = on
wal_buffers = 64MB

-- WAL ET RÉPLICATION
wal_level = replica
max_wal_senders = 10
wal_keep_size = 10GB
hot_standby = on
hot_standby_feedback = on
synchronous_commit = on
synchronous_standby_names = 'FIRST 1 (replica1_pap)'

-- ARCHIVAGE WAL vers MinIO
archive_mode = on
archive_command = 'aws s3 cp %p s3://snisid-wal-archive/%f --endpoint-url=https://minio.snisid.gouv.ht'
archive_timeout = 300              -- Archive toutes les 5 minutes au max

-- CHIFFREMENT DONNÉES AU REPOS (pgcrypto)
-- Note: chiffrement applicatif via AES-256-GCM pour colonnes sensibles

-- LOGGING ET AUDIT
log_destination = 'jsonlog'
logging_collector = on
log_directory = '/var/log/postgresql'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_rotation_age = 1d
log_rotation_size = 1GB
log_min_duration_statement = 1000  -- Log requêtes > 1s
log_connections = on
log_disconnections = on
log_duration = off
log_lock_waits = on
log_checkpoints = on
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
log_statement = 'ddl'             -- Log tous les DDL (CREATE, ALTER, DROP)

-- AUTOVACUUM
autovacuum = on
autovacuum_max_workers = 4
autovacuum_naptime = 60

-- MAINTENANCE
checkpoint_completion_target = 0.9
checkpoint_timeout = 15min
max_wal_size = 20GB
min_wal_size = 2GB
```

### 5.3 CockroachDB — Audit Ledger Distribué

```mermaid
graph LR
    subgraph "CockroachDB Cluster — Audit Ledger"
        CR1["CockroachDB Node 1 — PAP\n10.30.2.10:26257\nRegion: haiti-pap"]
        CR2["CockroachDB Node 2 — PAP\n10.30.2.11:26257\nRegion: haiti-pap"]
        CR3["CockroachDB Node 3 — PAP\n10.30.2.12:26257\nRegion: haiti-pap"]
        CR4["CockroachDB Node 4 — CAP\n10.130.2.10:26257\nRegion: haiti-cap"]
        CR5["CockroachDB Node 5 — CAP\n10.130.2.11:26257\nRegion: haiti-cap"]
    end

    subgraph "Consensus Raft"
        RAFT[Consensus Raft — Quorum 3/5]
    end

    subgraph "Applications"
        AUDIT_SVC[Audit Service]
        BLOCKCHAIN[Hyperledger Fabric\nAncre des hashes]
    end

    CR1 <-->|Raft| CR2
    CR1 <-->|Raft| CR3
    CR2 <-->|Raft| CR4
    CR3 <-->|Raft| CR5
    CR4 <-->|Raft| CR5

    AUDIT_SVC -->|Write audit records| CR1
    AUDIT_SVC -->|Write audit records| CR4
    BLOCKCHAIN -->|Anchor block hashes| CR1
```

---

## 6. Topologie Apache Kafka

### 6.1 Architecture Kafka Cluster

```mermaid
graph TB
    subgraph "Kafka Cluster PAP — KRaft Mode (sans ZooKeeper)"
        direction TB
        subgraph "Controllers (KRaft)"
            KC1["Kafka Controller 1\nKRaft Controller\n10.40.0.10:9093"]
            KC2["Kafka Controller 2\nKRaft Controller\n10.40.0.11:9093"]
            KC3["Kafka Controller 3\nKRaft Controller\n10.40.0.12:9093"]
        end

        subgraph "Brokers (Combined Mode)"
            KB1["Kafka Broker 1\n16vCPU / 64GB\n10TB NVMe SSD\n10.40.0.10:9092"]
            KB2["Kafka Broker 2\n16vCPU / 64GB\n10TB NVMe SSD\n10.40.0.11:9092"]
            KB3["Kafka Broker 3\n16vCPU / 64GB\n10TB NVMe SSD\n10.40.0.12:9092"]
            KB4["Kafka Broker 4\n16vCPU / 64GB\n10TB NVMe SSD\n10.40.0.13:9092"]
            KB5["Kafka Broker 5\n16vCPU / 64GB\n10TB NVMe SSD\n10.40.0.14:9092"]
        end

        subgraph "Schema Registry"
            SR[Schema Registry — Confluent\n10.40.0.20:8081\nAvro / Protobuf / JSON Schema]
        end

        subgraph "Kafka Connect"
            KC_POSTGRES[Kafka Connect — PostgreSQL Sink]
            KC_ELASTIC[Kafka Connect — Elasticsearch Sink]
            KC_S3[Kafka Connect — MinIO/S3 Sink (archivage)]
        end

        subgraph "Topics Principaux"
            T1[identity.events\nPartitions: 12\nRF: 3\nRetention: 365j]
            T2[biometric.events\nPartitions: 6\nRF: 3\nRetention: 90j\nChiffré]
            T3[enrollment.events\nPartitions: 12\nRF: 3\nRetention: 365j]
            T4[audit.events\nPartitions: 24\nRF: 5\nRetention: 7ans\nImmuable]
            T5[notification.commands\nPartitions: 6\nRF: 3\nRetention: 7j]
            T6[interop.sync\nPartitions: 12\nRF: 3\nRetention: 30j]
        end
    end

    subgraph "Kafka Cluster CAP — MirrorMaker 2"
        MM2[MirrorMaker 2\nActive-Active replication\n10.140.0.10]
        KB_CAP1[Kafka Broker 1 CAP]
        KB_CAP2[Kafka Broker 2 CAP]
        KB_CAP3[Kafka Broker 3 CAP]
    end

    subgraph "Producteurs"
        ID_SVC_P[Identity Service Producer]
        BIO_SVC_P[Biometric Service Producer]
        ENROLL_P[Enrollment Service Producer]
    end

    subgraph "Consommateurs"
        AUDIT_C[Audit Service Consumer Group: audit-writers]
        NOTIF_C[Notification Service Consumer Group: notifiers]
        SEARCH_C[Search Service Consumer Group: indexers]
        INTEROP_C[Interop Gateway Consumer Group: syncs]
    end

    KC1 <-->|Raft| KC2
    KC1 <-->|Raft| KC3

    ID_SVC_P -->|Produce TLS SASL| KB1
    BIO_SVC_P -->|Produce TLS SASL| KB2
    ENROLL_P -->|Produce TLS SASL| KB3

    KB1 <-->|Replication RF=3| KB2
    KB2 <-->|Replication RF=3| KB3
    KB3 <-->|Replication RF=3| KB4
    KB4 <-->|Replication RF=3| KB5

    KB1 --> T1
    KB2 --> T2
    KB3 --> T3
    KB4 --> T4
    KB5 --> T5

    T4 --> AUDIT_C
    T5 --> NOTIF_C
    T1 --> SEARCH_C
    T6 --> INTEROP_C

    KB1 <-->|MirrorMaker 2 replication| MM2
    MM2 --> KB_CAP1
    MM2 --> KB_CAP2
    MM2 --> KB_CAP3

    KB1 --> KC_POSTGRES
    KB2 --> KC_ELASTIC
    T4 --> KC_S3
```

### 6.2 Configuration Kafka Production

```yaml
# kafka-strimzi-cluster.yaml — Déploiement via Strimzi Operator
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: snisid-kafka
  namespace: kafka-system
spec:
  kafka:
    version: "3.7.0"
    replicas: 5
    listeners:
      - name: tls
        port: 9093
        type: internal
        tls: true
        authentication:
          type: tls  # mTLS pour communication interne
      - name: external
        port: 9094
        type: loadbalancer
        tls: true
        authentication:
          type: scram-sha-512  # SCRAM pour agents terrain
    authorization:
      type: opa
      url: http://opa.security-system.svc.cluster.local:8181/v1/data/kafka/authz
    config:
      default.replication.factor: 3
      min.insync.replicas: 2
      offsets.topic.replication.factor: 3
      transaction.state.log.replication.factor: 3
      transaction.state.log.min.isr: 2
      log.retention.hours: 8760  # 1 an par défaut
      log.segment.bytes: 1073741824  # 1 GB
      log.retention.check.interval.ms: 300000
      compression.type: lz4
      auto.create.topics.enable: false  # Topics créés explicitement
      delete.topic.enable: true
      unclean.leader.election.enable: false  # IMPORTANT: évite perte données
      message.max.bytes: 10485760  # 10 MB max message
    storage:
      type: persistent-claim
      size: 10Ti
      class: ceph-rbd-performance
      deleteClaim: false
    metricsConfig:
      type: jmxPrometheusExporter
    jvmOptions:
      -Xms: 32g
      -Xmx: 32g

  zookeeper:  # KRaft mode: zookeeper section absent en production KRaft
    replicas: 0

  entityOperator:
    topicOperator: {}
    userOperator: {}
```

---

## 7. Topologie Stockage

### 7.1 Architecture Ceph

```mermaid
graph TB
    subgraph "Ceph Cluster PAP — 6 Nœuds OSD"
        direction TB
        subgraph "Monitors (MON)"
            MON1[Ceph Monitor 1\n10.50.0.10]
            MON2[Ceph Monitor 2\n10.50.0.11]
            MON3[Ceph Monitor 3\n10.50.0.12]
        end

        subgraph "Managers (MGR)"
            MGR1[Ceph Manager 1\nDashboard, Prometheus]
            MGR2[Ceph Manager 2\nStandby]
        end

        subgraph "OSD Nodes — Données"
            OSD1["OSD Node 1\n4× 16TB SAS 12G\n64vCPU, 128GB RAM\nOSD 0-3"]
            OSD2["OSD Node 2\n4× 16TB SAS 12G\n64vCPU, 128GB RAM\nOSD 4-7"]
            OSD3["OSD Node 3\n4× 16TB SAS 12G\n64vCPU, 128GB RAM\nOSD 8-11"]
            OSD4["OSD Node 4\n4× 16TB SAS 12G\n64vCPU, 128GB RAM\nOSD 12-15"]
            OSD5["OSD Node 5\n4× 16TB SAS 12G\n64vCPU, 128GB RAM\nOSD 16-19"]
            OSD6["OSD Node 6\n4× 16TB SAS 12G\n64vCPU, 128GB RAM\nOSD 20-23"]
        end

        subgraph "CRUSH Map — Failure Domains"
            CRUSH[CRUSH Algorithm\nFailure domain: host\nReplicas: 3\nCapacité totale: ~384 TB raw\n~128 TB utilisable]
        end

        subgraph "Pools Ceph"
            POOL_RBD["Pool: rbd-identity\nType: Replicated 3x\nUsage: PVC PostgreSQL Identity\nPerformance: SSD cache tier"]
            POOL_RBD_BIO["Pool: rbd-biometric\nType: Replicated 3x\nUsage: PVC PostgreSQL Biometric\nChiffrement: dmcrypt"]
            POOL_RBD_KAFKA["Pool: rbd-kafka\nType: Replicated 2x\nUsage: PVC Kafka Brokers\nPerformance: NVMe tier"]
            POOL_CFS["Pool: cephfs-documents\nType: Replicated 3x\nUsage: CephFS Documents\nMulti-client mount"]
            POOL_RGW["Pool: rgw-objects\nType: Replicated 3x + EC 6+3\nUsage: Object Storage (S3 API)"]
        end
    end

    subgraph "Interfaces d'Accès"
        RBD_CSI[RBD CSI Driver\nk8s StorageClass]
        CEPHFS_CSI[CephFS CSI Driver\nk8s StorageClass]
        RADOSGW[RADOS Gateway\nS3 + Swift API compatible]
        MINIO[MinIO Gateway\nhttps://minio.snisid.gouv.ht]
    end

    MON1 <-->|Quorum| MON2
    MON1 <-->|Quorum| MON3
    MGR1 --> MON1
    OSD1 <-->|CRUSH| CRUSH
    OSD2 <-->|CRUSH| CRUSH
    OSD3 <-->|CRUSH| CRUSH
    OSD4 <-->|CRUSH| CRUSH
    OSD5 <-->|CRUSH| CRUSH
    OSD6 <-->|CRUSH| CRUSH
    CRUSH --> POOL_RBD
    CRUSH --> POOL_RBD_BIO
    CRUSH --> POOL_RBD_KAFKA
    CRUSH --> POOL_CFS
    CRUSH --> POOL_RGW
    POOL_RBD --> RBD_CSI
    POOL_CFS --> CEPHFS_CSI
    POOL_RGW --> RADOSGW
    RADOSGW --> MINIO
```

### 7.2 StorageClasses Kubernetes

```yaml
# StorageClasses SNISID — Ceph CSI
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: ceph-rbd-performance
  annotations:
    storageclass.kubernetes.io/is-default-class: "false"
provisioner: rbd.csi.ceph.com
parameters:
  clusterID: snisid-ceph-pap
  pool: rbd-identity
  imageFormat: "2"
  imageFeatures: layering
  csi.storage.k8s.io/provisioner-secret-name: rook-csi-rbd-provisioner
  csi.storage.k8s.io/provisioner-secret-namespace: rook-ceph
  csi.storage.k8s.io/controller-expand-secret-name: rook-csi-rbd-provisioner
  csi.storage.k8s.io/node-stage-secret-name: rook-csi-rbd-node
  csi.storage.k8s.io/fstype: ext4
  encrypted: "true"  # Chiffrement LUKS2 côté volume
  encryptionKMSID: "vault-tokens-snisid"  # Clés via Vault
reclaimPolicy: Retain  # JAMAIS Delete pour données production
allowVolumeExpansion: true
volumeBindingMode: WaitForFirstConsumer
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: ceph-rbd-biometric
  labels:
    data-classification: biometric
    compliance: restricted
provisioner: rbd.csi.ceph.com
parameters:
  clusterID: snisid-ceph-pap
  pool: rbd-biometric
  encrypted: "true"
  encryptionKMSID: "vault-tokens-biometric"  # Clés séparées pour biométrie
reclaimPolicy: Retain
allowVolumeExpansion: false  # Expansion manuelle après approbation
```

---

## 8. Architecture DNS

### 8.1 Hiérarchie DNS SNISID

```mermaid
graph TB
    subgraph "DNS Externe — Publiquement Résolvable"
        ROOT_DNS[Root DNS — IANA\n. root]
        HT_TLD[.ht TLD\nAHAITI — Autorité .ht]
        GOUV_HT[gouv.ht\nServeurs DNS Gouvernementaux\nns1.gouv.ht, ns2.gouv.ht]
        SNISID_EXT[snisid.gouv.ht\nZone Externe\nns1-ext.snisid.gouv.ht\nns2-ext.snisid.gouv.ht]
    end

    subgraph "Enregistrements DNS Externes"
        EXT_RECORDS["www.snisid.gouv.ht → CDN/WAF IP
        api.snisid.gouv.ht → EXT LB VIP
        portail.snisid.gouv.ht → CDN
        enroll.snisid.gouv.ht → CDN
        diaspora.snisid.gouv.ht → CDN"]
    end

    subgraph "DNS Interne — Privé"
        CORE_DNS["CoreDNS — Kubernetes\nIn-cluster DNS\n.cluster.local"]
        INT_DNS1["BIND 9 — DNS Interne PAP\nns1-int.snisid.gouv.ht\n10.10.0.53"]
        INT_DNS2["BIND 9 — DNS Interne CAP\nns2-int.snisid.gouv.ht\n10.110.0.53"]
    end

    subgraph "Zones Internes"
        ZONE_SVC["snisid.internal\nZone services internes"]
        INT_RECORDS["identity-svc.snisid.internal → 10.20.x.x
        biometric-svc.snisid.internal → 10.20.x.x
        auth-svc.snisid.internal → 10.20.x.x
        kafka.snisid.internal → 10.40.0.10-14
        vault.snisid.internal → 10.60.0.10-12
        postgres-rw.snisid.internal → 10.30.0.5
        postgres-ro.snisid.internal → 10.30.0.6"]
    end

    subgraph "DNSSEC"
        DNSSEC[DNSSEC Signing\nAlgorithme: ECDSAP384SHA384\nZSK rotation: 30j\nKSK rotation: 1an]
    end

    ROOT_DNS --> HT_TLD --> GOUV_HT --> SNISID_EXT
    SNISID_EXT --> EXT_RECORDS
    SNISID_EXT -.->|Délégation split-DNS| INT_DNS1
    INT_DNS1 <-->|Zone transfer TSIG| INT_DNS2
    INT_DNS1 --> ZONE_SVC --> INT_RECORDS
    CORE_DNS -->|Forward non-.cluster| INT_DNS1
    SNISID_EXT --> DNSSEC
```

---

## 9. Load Balancers

### 9.1 Architecture F5 BIG-IP

```yaml
# Configuration F5 BIG-IP — Virtual Servers SNISID
virtual_servers:
  vs_api_external:
    name: "VS_SNISID_API_EXT"
    ip: "203.x.x.100"  # IP publique
    port: 443
    protocol: TCP
    ssl_profile_client: "SSL_PROFILE_SNISID_TLS13"
    ssl_profile_server: "SSL_PROFILE_BACKEND_MTLS"
    pool: "POOL_KONG_GATEWAYS"
    irule: ["IRULE_GEOBLOCKING", "IRULE_HEADER_INSERT"]
    persistence: cookie_insert
    health_monitor: "MONITOR_KONG_HTTPS"

  vs_api_mtls:
    name: "VS_SNISID_API_MTLS"
    ip: "203.x.x.101"
    port: 8443
    protocol: TCP
    ssl_profile_client: "SSL_PROFILE_MTLS_GOV_AGENCIES"
    pool: "POOL_KONG_INTERNAL"
    irule: ["IRULE_CLIENT_CERT_VALIDATION", "IRULE_AGENCY_ROUTING"]

pools:
  POOL_KONG_GATEWAYS:
    load_balancing_method: least_connections
    health_monitor: MONITOR_KONG_HTTPS
    members:
      - {ip: "10.10.0.21", port: 8443, priority: 1}  # Kong 1 PAP
      - {ip: "10.10.0.22", port: 8443, priority: 1}  # Kong 2 PAP
      - {ip: "10.10.0.23", port: 8443, priority: 1}  # Kong 3 PAP
      - {ip: "10.110.0.21", port: 8443, priority: 2}  # Kong 1 CAP (failover)
    slow_ramp_time: 30

ssl_profiles:
  SSL_PROFILE_SNISID_TLS13:
    ciphers: "ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384"
    protocols: TLSv1.3
    options: "no-sslv3 no-tlsv1 no-tlsv1.1 no-tlsv1.2"
    ocsp_stapling: enabled
    hsts: "max-age=63072000; includeSubDomains; preload"

irules:
  IRULE_GEOBLOCKING: |
    when HTTP_REQUEST {
      set country [whereis [IP::client_addr] country]
      # Autoriser: HT (Haïti), FR (France diaspora), CA, US, ...
      if { not ($country equals "HT" or $country equals "FR" or
                $country equals "CA" or $country equals "US" or
                $country equals "DO" or $country equals "") } {
        HTTP::respond 403 content "Accès géo-restreint / Geo-restricted access"
      }
    }
```

---

## 10. Zones Pare-feu et Règles

### 10.1 Architecture des Zones de Sécurité

```mermaid
graph LR
    subgraph "ZONE 0 — INTERNET"
        INET[Internet / WAN]
    end

    subgraph "ZONE 1 — DMZ EXTERNE"
        direction TB
        CDN_WAF[CDN / WAF]
        BORDER[Border Firewall]
    end

    subgraph "ZONE 2 — DMZ APPLICATIVE"
        LB[Load Balancer]
        API_GW[API Gateway]
    end

    subgraph "ZONE 3 — APPLICATION"
        K8S_APP[Kubernetes Workers\nGeneral Purpose]
    end

    subgraph "ZONE 4 — DATA"
        DB_ZONE[Serveurs BD]
        KAFKA_ZONE[Kafka Brokers]
    end

    subgraph "ZONE 5 — BIOMETRIQUE ISOLÉE"
        BIO_ZONE[Biometric Servers\nAir-gap partiel]
        BIO_DB[Biometric Vault DB]
    end

    subgraph "ZONE 6 — SÉCURITÉ HSM"
        VAULT_ZONE[HashiCorp Vault]
        HSM_ZONE[HSM Luna 7]
    end

    subgraph "ZONE 7 — MANAGEMENT OOB"
        MGMT[Jump Server\niDRAC/iLO\nSNMP]
    end

    INET -->|HTTPS 443, 8443 seul| CDN_WAF
    CDN_WAF -->|Trafic filtré| BORDER
    BORDER -->|→ DMZ App| LB
    LB -->|HTTPS| API_GW
    API_GW -->|gRPC mTLS 9090-9100| K8S_APP
    K8S_APP -->|PostgreSQL 5432, Redis 6379| DB_ZONE
    K8S_APP -->|Kafka 9092-9093| KAFKA_ZONE
    K8S_APP -->|gRPC 9090 UNIQUEMENT| BIO_ZONE
    BIO_ZONE -->|PostgreSQL 5432| BIO_DB
    K8S_APP -->|Vault API 8200| VAULT_ZONE
    BIO_ZONE -->|Vault API 8200| VAULT_ZONE
    VAULT_ZONE -->|PKCS#11/HSM| HSM_ZONE
    MGMT -->|SSH 22 Bastion ONLY| K8S_APP
    MGMT -->|iDRAC/iLO 443| DB_ZONE

    style ZONE_5 fill:#ff6b6b,color:#fff
    style ZONE_6 fill:#4a4a8a,color:#fff
    style HSM_ZONE fill:#2d2d5f,color:#fff
```

### 10.2 Règles Pare-feu Conceptuelles

```yaml
# Matrice de règles pare-feu SNISID
firewall_rules:
  zone_internet_to_dmz_ext:
    - id: FW-001
      src: "0.0.0.0/0"
      dst: "CDN-WAF-VIP"
      port: [443, 80]
      protocol: TCP
      action: ALLOW
      description: "HTTPS public vers CDN"
      log: true

    - id: FW-002
      src: "GOV-AGENCY-IPs"  # IPs des agences gouvernementales
      dst: "API-MTLS-VIP"
      port: [8443]
      protocol: TCP
      action: ALLOW
      condition: "client_cert_required=true"
      description: "API mTLS agences gouvernementales"
      log: true

    - id: FW-DEFAULT
      src: any
      dst: any
      action: DENY
      log: true
      description: "Deny all — default drop"

  dmz_app_to_application:
    - id: FW-010
      src: "API-GATEWAY-SUBNET"
      dst: "K8S-WORKER-SUBNET"
      port: [9090, 9091, 9092, 9100, 8080]
      protocol: TCP
      action: ALLOW
      description: "Kong → Kubernetes Services (gRPC, HTTP)"

  application_to_data:
    - id: FW-020
      src: "K8S-WORKER-SUBNET"
      dst: "DB-SUBNET"
      port: [5432, 5433, 6379, 6380]
      protocol: TCP
      action: ALLOW
      description: "Application → PostgreSQL, Redis"
      log: true

    - id: FW-021
      src: "K8S-WORKER-SUBNET"
      dst: "KAFKA-SUBNET"
      port: [9092, 9093, 9094]
      protocol: TCP
      action: ALLOW
      description: "Application → Kafka"

  application_to_biometric:
    - id: FW-030
      src: ["IDENTITY-SVC-IP", "ENROLLMENT-SVC-IP", "API-GW-IP"]
      dst: "BIOMETRIC-ZONE-SUBNET"
      port: [9090]  # gRPC UNIQUEMENT
      protocol: TCP
      action: ALLOW
      description: "Accès restreint Biometric Service"
      log: true
      alert_on: "source_not_in_whitelist"

    - id: FW-031
      src: any
      dst: "BIOMETRIC-ZONE-SUBNET"
      action: DENY
      log: true
      alert: true
      description: "Block all unauthorized biometric access"

  application_to_hsm_vault:
    - id: FW-040
      src: "K8S-ALL-SUBNET"
      dst: "VAULT-CLUSTER-IPs"
      port: [8200, 8201]
      protocol: TCP
      action: ALLOW
      description: "Services → Vault API"

    - id: FW-041
      src: "VAULT-CLUSTER-IPs"
      dst: "HSM-IPs"
      port: [1792]  # Luna NTLS
      protocol: TCP
      action: ALLOW
      description: "Vault → HSM Luna Network"
      log: true
```

---

## 11. Plan d'Adressage IP

### 11.1 Table d'Adressage Complète

| Réseau / Subnet | CIDR | Plage Hôtes | Usage | Datacenter |
|---|---|---|---|---|
| **Infrastructure PAP** | | | | |
| DMZ Externe | 10.10.0.0/24 | 10.10.0.1-254 | LB, API Gateway, CDN reverse | PAP |
| Application Kubernetes | 10.20.0.0/20 | 10.20.0.1-4094 | Kubernetes Pod Network | PAP |
| Services K8s | 10.21.0.0/16 | 10.21.0.1-65534 | Kubernetes Service Network | PAP |
| Data — PostgreSQL | 10.30.0.0/24 | 10.30.0.1-254 | PostgreSQL clusters | PAP |
| Data — Biométrique | 10.30.1.0/24 | 10.30.1.1-254 | Biometric DB — isolé | PAP |
| Data — Audit/Cockroach | 10.30.2.0/24 | 10.30.2.1-254 | CockroachDB audit | PAP |
| Messagerie Kafka | 10.40.0.0/24 | 10.40.0.1-254 | Kafka Brokers + Controllers | PAP |
| Stockage Ceph | 10.50.0.0/24 | 10.50.0.1-254 | Ceph Monitors, MGR, OSD | PAP |
| Ceph Storage Network | 10.51.0.0/24 | 10.51.0.1-254 | Réseau de réplication Ceph | PAP |
| Sécurité HSM/Vault | 10.60.0.0/24 | 10.60.0.1-254 | Vault cluster, HSM | PAP |
| Management OOB | 10.70.0.0/24 | 10.70.0.1-254 | iDRAC, iLO, IPMI, Bastion | PAP |
| **Infrastructure CAP** | | | | |
| DMZ Externe CAP | 10.110.0.0/24 | 10.110.0.1-254 | LB, API Gateway CAP | CAP |
| Application K8s CAP | 10.120.0.0/20 | 10.120.0.1-4094 | Kubernetes Pod Network CAP | CAP |
| Services K8s CAP | 10.121.0.0/16 | 10.121.0.1-65534 | Kubernetes Service Network CAP | CAP |
| Data — PostgreSQL CAP | 10.130.0.0/24 | 10.130.0.1-254 | PostgreSQL HA CAP | CAP |
| Data — Biométrique CAP | 10.130.1.0/24 | 10.130.1.1-254 | Biometric DB CAP | CAP |
| Messagerie Kafka CAP | 10.140.0.0/24 | 10.140.0.1-254 | Kafka Cluster CAP | CAP |
| Stockage Ceph CAP | 10.150.0.0/24 | 10.150.0.1-254 | Ceph CAP | CAP |
| **Réseau Gouvernemental** | | | | |
| Sites OEC Terrain | 172.16.0.0/16 | 172.16.0.1-65534 | 145 bureaux OEC (sous-réseaux /24) | Terrain |
| Sites MEL Terrain | 172.17.0.0/16 | 172.17.0.1-65534 | Bureaux électoraux | Terrain |
| Sites Police PNH | 172.18.0.0/16 | 172.18.0.1-65534 | Postes de police | Terrain |
| Sites MAECI | 172.19.0.0/16 | 172.19.0.1-65534 | Ambassades et consulats | International |
| **Inter-Datacenter** | | | | |
| Lien WAN PAP-CAP | 192.168.100.0/30 | 192.168.100.1-2 | eBGP point-to-point | WAN |
| Lien WAN Backup | 192.168.101.0/30 | 192.168.101.1-2 | eBGP backup | WAN |

### 11.2 IPAM et Gestion

```yaml
# IPAM Configuration — NetBox SNISID
ipam:
  tool: NetBox
  version: "3.7"
  url: "https://ipam.snisid.gouv.ht"
  auth: LDAP + MFA

  prefixes_auto_allocation:
    kubernetes_pods: "10.20.0.0/20"
    kubernetes_services: "10.21.0.0/16"
    algorithm: "first-available"

  reservations:
    - {ip: "10.10.0.1", role: "Gateway DMZ — FortiGate"}
    - {ip: "10.10.0.5", role: "VIP F5 externe"}
    - {ip: "10.30.0.5", role: "VIP PostgreSQL Read-Write"}
    - {ip: "10.30.0.6", role: "VIP PostgreSQL Read-Only"}
    - {ip: "10.70.0.100", role: "Bastion Host"}

  dns_integration:
    provider: BIND 9
    auto_register: true
    zones: ["snisid.internal", "10.in-addr.arpa"]
```

---

## 12. Topologie Physique des Datacenters

### 12.1 DC Port-au-Prince — Spécifications

```yaml
datacenter_pap:
  nom: "Centre de Données Primaire — Port-au-Prince"
  localisation: "Port-au-Prince, Haïti"
  classification: "TIA-942 Tier III (objectif Tier III+)"
  superficie: "400 m²"

  alimentation:
    utilitaire_primaire: "EDH — Électricité d'Haïti 20kV"
    utilitaire_secondaire: "Natcom Power Grid 20kV"
    ups_systeme: "APC Galaxy VM 200kVA × 2 (N+1)"
    autonomie_ups: "4 heures pleine charge"
    generateurs:
      - {modele: "Caterpillar C175 500kVA", qty: 2, config: "N+1"}
      - {carburant: "Diesel", reservoir: "10000L", autonomie: "72h"}
    pdu: "Schneider iPDU 3-phase redondant A/B"

  refrigeration:
    type: "Air conditionné de précision CRAC + Refroidissement liquide rangées"
    capacite: "300kW refroidissement"
    redondance: "N+1 (5 CRAC unités, 1 standby)"
    temperature_cible: "21°C ± 2°C"
    humidite_cible: "45-55% RH"
    free_cooling: "Economiseur pour saison fraîche"

  connectivite:
    operateurs:
      - {nom: "Natcom", capacite: "10Gbps Fibre", type: "MPLS + Internet"}
      - {nom: "Digicel Business", capacite: "10Gbps Fibre", type: "MPLS + Internet"}
      - {nom: "Link-Up", capacite: "1Gbps Fibre", type: "Internet backup"}
      - {nom: "Starlink Enterprise", capacite: "500Mbps", type: "Satellite backup"}
    diversite_physique: "Entrées câbles séparées, côtés opposés du bâtiment"

  securite_physique:
    perimetre: "Clôture électrifiée + caméras CCTV 4K"
    acces: "Biométrique (empreintes + badge) à 3 niveaux"
    gardiennage: "24/7 — 4 agents minimum"
    cctv: "128 caméras — rétention 90 jours"
    cage_faraday: "Salle HSM/PKI en cage Faraday certifiée"
    anti_intrusion: "Capteurs vibration, mouvement infrarouge, portes magnétiques"

  serveurs:
    total_serveurs: 45
    total_racks: 15
    puissance_installee: "180 kW"
    switches_core: "Arista 7280 × 2 (MLAG)"
    switches_access: "Arista 7050 × 8 (25GbE)"
    patch_panels: "Fibre LC/LC OM4 multi-mode + monomode"

  conformite:
    certifications: ["ISO 27001", "TIA-942 Tier III audit (en cours)"]
    tests: "Test basculement générateur mensuel"
    dr_drill: "Semestriel"
```

---

## Bloc d'Approbation / Approval Block

| Rôle | Nom | Signature | Date |
|---|---|---|---|
| **Directeur Infrastructure** | [À compléter] | [Signature] | 2026-05-25 |
| **Architecte Réseau** | [À compléter] | [Signature] | 2026-05-25 |
| **Architecte Sécurité** | [À compléter] | [Signature] | 2026-05-25 |
| **Directeur Général SNISID** | [À compléter] | [Signature] | 2026-05-25 |

---

*Document SNISID-ARC-INF-001 v1.0.0 — CONFIDENTIEL — © République d'Haïti, Programme SNISID, 2026*
