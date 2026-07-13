# SNISID: Secure Encrypted Backup & Restore Architecture

The Sovereign Recovery System ensures that the nation's identity ledger can survive catastrophic physical destruction, regional network outages, or sophisticated ransomware campaigns.

---

## 1. Backup Architecture: Multi-Tiered Protection

SNISID implements a hierarchy of backups to balance performance with absolute durability.

### 1.1. Data Source Strategies
- **PostgreSQL**: Continuous archiving of WAL (Write Ahead Logs) to S3, with daily full snapshots.
- **Neo4j**: Weekly full database dumps + incremental transaction log backups.
- **MinIO (Object Storage)**: Cross-Region Replication (CRR) + Versioning.
- **Kubernetes (State)**: Velero snapshots of all Persistent Volumes (PVs) and ETCD state.

### 1.2. Storage Tiers
- **Tier 1 (On-Site Snapshots)**: Fast recovery from accidental deletions.
- **Tier 2 (Off-Site Vault)**: Synchronous replication to the Disaster Recovery (DR) Region.
- **Tier 3 (Air-Gapped Sovereign Archive)**: Weekly physical export to a disconnected vault.

---

## 2. Encryption Model & Key Isolation

Backup data is useless to an attacker due to strict cryptographic isolation.

- **Encryption at Source**: Data is encrypted using AES-256-GCM *before* it leaves the database host or storage controller.
- **Key Separation**: The **Backup Encryption Keys (BEK)** are managed by a dedicated, physically isolated HSM partition. These keys are never stored on the same network as the primary identity database.
- **Envelope Encryption**: Each backup set has its own unique DEK, which is wrapped by the Master BEK.

---

## 3. Ransomware Resilience: WORM & Air-Gap

To prevent a "Wipeout" attack, SNISID enforces two physical and logical barriers:

- **Immutable WORM Buckets**: All backups are written to S3 buckets in **Compliance-mode Object Lock**. Even a compromised Root Admin cannot delete or overwrite a backup until the retention period (e.g., 30 days for daily, 10 years for weekly) expires.
- **Physical Air-Gap**: The Sovereign Archive is stored on encrypted LTO-9 Tapes or SSDs. The "Restore" process from this tier requires physical access to a high-security bunker and manual ingestion into an isolated "Clean Room" network.

---

## 4. Automated Restoration Validation

A backup is only a backup if it can be restored.

- **The Validation Sandbox**: Every 24 hours, the system automatically pulls a random backup set and restores it into a temporary, isolated Kubernetes namespace.
- **Integrity Testing**: Automated scripts run a battery of tests (e.g., "Can I query Identity X?", "Is the Biometric Hash consistent?") to ensure the restored data is functional.
- **Health Reporting**: Any failure in the validation drill triggers a **Critical P0 Alert** to the National SOC.

---

## 5. Disaster Recovery (DR) Strategy

SNISID maintains a **Warm Standby** configuration.

- **Failover Logic**: If the Primary Region is declared "Down" by the Global Load Balancer, traffic is routed to the DR Region.
- **Data Consistency**: Synchronous Cross-Region Replication (CRR) ensures that the DR database is identical to the Primary (RPO < 1 minute).
- **RTO (Recovery Time Objective)**: 4 hours to achieve 100% capacity in the DR region.

---

## 6. Recovery SLA Framework

| Component | RPO (Data Loss) | RTO (Recovery Time) |
| :--- | :--- | :--- |
| **Core Identity API** | < 1 Minute | 15 Minutes |
| **Biometric Matcher** | < 1 Minute | 1 Hour |
| **Audit Ledger** | 0 (Sync) | 30 Minutes |
| **Graph Analytics** | 1 Hour | 4 Hours |
| **Historical Archives** | 24 Hours | 48 Hours (Physical) |

---

## 7. Integrity Verification Model

- **Continuous Checksumming**: A background process in the storage cluster continuously re-verifies the SHA-256 hashes of all stored backup objects to detect bit-rot.
- **Digital Signatures**: Every backup manifest is signed by the **Recovery Service Identity**. During a restore, the system verifies the signature to ensure the backup has not been tampered with or replaced with a malicious payload.
