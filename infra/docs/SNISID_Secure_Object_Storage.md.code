# SNISID: Secure Encrypted Object Storage

The Sovereign Object Vault is designed to store the nation's most sensitive binary data—biometric profiles, scanned evidence, and system backups—with absolute integrity, cryptographic protection, and forensic non-repudiation.

---

## 1. Storage Architecture: Distributed & Resilient

SNISID utilizes a distributed **MinIO / S3-compatible** cluster deployed across three national data center zones.

- **Erasure Coding**: Objects are striped across multiple nodes and drives, allowing the cluster to survive the loss of up to 50% of its physical hardware without data loss.
- **Malware Scanning (Pre-Upload)**:
  1.  Microservice uploads object to a **Quarantine Bucket**.
  2.  A sidecar service (ClamAV/ICAP) scans the file for malicious payloads.
  3.  Only "Clean" files are moved to the **Production Vault**.
  4.  "Infected" files are moved to a **Forensic Isolation Bucket** for SOC analysis.

---

## 2. Encryption Model: Defense-in-Depth

### 2.1. Server-Side Encryption (SSE-KMS)
All objects are encrypted by default using AES-256-GCM.
- **Master Key**: Managed by the **Sovereign KMS (Vault)**.
- **Granularity**: Keys are managed at the **Bucket** or **Tenant Prefix** level to ensure agency isolation.

### 2.2. Client-Side Encryption (SSE-C)
For ultra-sensitive biometrics (e.g., facial embeddings), the microservice encrypts the data *locally* using a unique record-level key before sending it to the storage cluster. This ensures that even a storage administrator cannot view the data.

---

## 3. Immutability & Integrity (WORM)

To protect forensic evidence and national records, SNISID enforces **WORM (Write Once, Read Many)** storage.

- **Object Lock (Compliance Mode)**: Once a document (e.g., a criminal record) is written, it **cannot be deleted or overwritten** by anyone, including the `root` user, until the retention period (e.g., 10 years) has expired.
- **Content-Addressable Storage (CAS)**: Every object is indexed by its SHA-256 hash. If the hash does not match the content during a `GET` request, the system triggers a critical integrity alert.
- **Versioning**: Every modification creates a new version; older versions are preserved to prevent accidental or malicious overwriting.

---

## 4. Access Workflows: Least Privilege

- **No Public Access**: All buckets are strictly private and not reachable via public internet.
- **Presigned URLs**: Microservices request a short-lived (e.g., 5-minute) **Presigned URL** from the Storage Controller. The application then uses this URL to perform the specific `PUT` or `GET` operation.
- **ABAC Enforcement**: Access to specific folders (e.g., `agency/police/evidence/`) is governed by OPA policies that check the `agency_id` and `clearance_level` of the requesting service.

---

## 5. Replication & Archival Strategy

### 5.1. Multi-Region Replication
- **Synchronous CRR**: Data written to the Primary region is synchronously replicated to the Disaster Recovery (DR) region.
- **Availability**: During a regional outage, services automatically fail over to the DR storage cluster.

### 5.2. Air-Gapped Sovereign Archival
Critical national datasets (Master Identity Ledger, Root CA backups) are periodically replicated to a physically isolated **Air-Gapped Vault**.

**Detailed Recovery Strategy**: See the [SNISID Backup & Restore Architecture](file:///c:/Users/sopil/Desktop/SNISID/SNISID_Backup_Restore_Architecture.md) for ransomware resilience, RTO/RPO SLAs, and automated restoration validation drills.

---

## 6. Access Auditing & Monitoring

Every storage event is tracked and routed to the **Sovereign Audit Ledger**:
- `s3.object.put`: Who uploaded what, and was it scanned?
- `s3.object.get`: Who accessed a sensitive record?
- `s3.object.lock`: When was a retention period set?
- `s3.integrity.fail`: Critical alert if a hash mismatch is detected.
