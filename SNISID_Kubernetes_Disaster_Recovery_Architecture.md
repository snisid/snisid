# SNISID Kubernetes Disaster Recovery Architecture
## Velero Backups, Etcd Resilience & Air-Gapped Recovery

This document details the **Kubernetes Disaster Recovery (DR) Architecture** for SNISID. Beyond standard active-active replication, this architecture focuses on "Black Swan" events—such as catastrophic ransomware wiping out the primary storage arrays, or a Category 5 hurricane physically destroying the Port-au-Prince datacenter. It ensures that the entire SNISID Kubernetes ecosystem can be restored from scratch.

---

## 1. Etcd Backup & Cluster State

The `etcd` key-value store is the absolute "brain" of Kubernetes. If `etcd` is lost, the cluster is dead.
- **Automated Snapshots:** CronJobs running on the master nodes execute `etcdctl snapshot save` every 15 minutes.
- **Off-Node Storage:** Snapshots are immediately pushed via an encrypted tunnel to an immutable S3 bucket. They are never kept solely on the local master node's disk.

---

## 2. Workload & Persistent Volume Backups (Velero)

SNISID utilizes **VMware Velero** to backup cluster resources, namespaces, and StatefulSet persistent volumes.

### Kubernetes Resources (Stateless)
- Velero takes hourly snapshots of all Kubernetes objects (Deployments, Services, ConfigMaps, Secrets, Istio VirtualServices, OPA Policies) and stores them as JSON/tarballs in an S3-compatible backend (e.g., MinIO or Ceph).

### Persistent Volumes (Stateful)
- **CSI Snapshots:** For databases (CockroachDB, Kafka), Velero utilizes the Ceph Container Storage Interface (CSI) to trigger instant storage-level snapshots.
- **Restic/Kopia:** For volumes that do not support native CSI snapshots, Velero uses Kopia to copy the filesystem data block-by-block.

---

## 3. Immutable & Air-Gapped Storage

To survive targeted ransomware attacks or malicious insiders attempting to wipe backups:
1. **Object Lock (WORM):** The primary Ceph S3 backup bucket enforces Write-Once-Read-Many (WORM). Once Velero writes a backup, it cannot be deleted or modified by anyone—even the storage administrator—for 30 days.
2. **Air-Gapped Cold Vault:** Once a day, the backups are pushed over a unidirectional data diode to an air-gapped, offline storage vault (or LTO tape drives in an offsite bunker). This ensures that even if the entire SNISID network is logically compromised, an untouched baseline exists offline.

---

## 4. Disaster Recovery Scenarios & Automation

### 1. Accidental Namespace Deletion (Targeted Restore)
- **Scenario:** A Junior admin accidentally runs `kubectl delete namespace snisid-identity`.
- **Response:** An SRE runs `velero restore create --from-backup <latest> --include-namespaces snisid-identity`. The Identity microservices are pulled from S3 and fully restored within 3 minutes.

### 2. Complete Cluster Loss (DC1 Destroyed)
- **Scenario:** Port-au-Prince is hit by a massive earthquake. The bare-metal cluster is gone.
- **Response:** 
  1. The DR team spins up a vanilla Kubernetes cluster in the DC2 (Cap-Haïtien) location via automated Terraform.
  2. Velero is installed and pointed to the replicated (or air-gapped) S3 backup bucket.
  3. `velero restore create --from-backup <latest>` is executed. All namespaces, policies, and persistent volumes are rehydrated into the new cluster.
  4. Global DNS automatically flips to point to the DC2 ingress.

---

## 5. Architecture & Recovery Diagrams (Mermaid)

### 1. Velero Backup & Air-Gapped Topology
This diagram illustrates how data flows from the live cluster into immutable and air-gapped storage.

```mermaid
graph TD
    classDef k8s fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef velero fill:#e1bee7,stroke:#6a1b9a,stroke-width:2px;
    classDef storage fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef secure fill:#ffebee,stroke:#c62828,stroke-width:2px;

    subgraph Port_au_Prince_DC1
        K8S[Live Kubernetes Cluster]:::k8s
        ETCD[(etcd Control Plane)]:::k8s
        PV[(Ceph Persistent Volumes)]:::k8s
        
        V[Velero Controller]:::velero
        K8S <--> V
        ETCD -.->|Snapshot CronJob| V
        PV -.->|CSI / Restic| V
    end

    subgraph Immutability_Tier [Networked Storage]
        S3[(S3 Object Storage <br/> Ceph RADOS)]:::storage
        WORM[WORM / Object Lock <br/> 30-Day Retention]:::storage
        S3 <--> WORM
    end

    subgraph Air_Gapped_Vault [Off-Grid DR Bunker]
        DIODE[Unidirectional Data Diode]:::secure
        TAPE[(LTO Cold Tape Storage <br/> Disconnected)]:::secure
    end

    V -->|Hourly Encrypted Push| S3
    S3 -->|Daily Sync| DIODE
    DIODE -->|One-Way Write| TAPE
```

### 2. Complete Cluster Recovery Workflow (DC1 to DC2)
This sequence maps the automated procedure to restore SNISID after a catastrophic physical loss.

```mermaid
sequenceDiagram
    participant SRE as Site Reliability Engineer
    participant TF as Terraform (GitOps)
    participant DC2 as Cap-Haïtien Cluster (New)
    participant S3 as DC2 S3 Backup Replica
    participant DNS as Sovereign DNS

    Note over SRE, DNS: DC1 is completely destroyed.
    
    SRE->>TF: 1. Execute `terraform apply -var="region=dc2"`
    TF->>DC2: 2. Provision vanilla K8s Bare-Metal nodes
    DC2-->>TF: 3. Cluster Ready
    
    TF->>DC2: 4. Install Velero & Configure S3 Credentials
    
    SRE->>DC2: 5. Execute `velero restore create --from-backup latest-dc1`
    DC2->>S3: 6. Fetch Kubernetes Manifests & PV Snapshots
    S3-->>DC2: 7. Download JSON/Tarballs
    
    DC2->>DC2: 8. Rehydrate Namespaces, Secrets, Deployments, PVs
    Note over DC2: Wait for Pods to report Ready (5 minutes)
    
    SRE->>DNS: 9. Update BGP / Global Load Balancer
    DNS->>DC2: 10. Route all national traffic to Cap-Haïtien
    Note over SRE, DNS: System Restored. RTO < 1 Hour.
```

---
*Prepared by the SNISID Cloud Infrastructure & Resilience Board.*
