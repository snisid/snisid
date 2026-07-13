# SNISID Security & Incident Response Runbook

## 🛡️ Security Architecture
- **Zero-Trust Networking**: All service communication is restricted via Kubernetes `NetworkPolicies`.
- **RBAC Enforcement**: Fine-grained access control at both the Gateway and Service levels.
- **Data Encryption**: 
    - At-rest: AES-256 (via storage provider).
    - In-transit: TLS 1.3 enforced for all external and internal traffic.

## 🚨 Incident Response (IR)

### 1. Unauthorized Access Detected
- **Detection**: Grafana alert "Unauthorized Access Spike" or Loki log "Invalid JWT signature".
- **Containment**: 
    - Rotate JWT signing keys (via Vault/Secrets).
    - Scale down the `api-gateway` to 1 replica and enable "Maintenance Mode".
    - Revoke suspicious tokens in Redis.

### 2. Massive Data Ingestion Anomaly
- **Detection**: Fraud Engine alert "Abnormal Identity Cluster".
- **Action**:
    - Trigger `scripts/freeze_agency.sh <AGENCY_ID>`.
    - Initiate manual audit of the GNN fraud graph for the specific cluster.

### 3. Database Integrity Failure
- **Detection**: Health check failure or checksum mismatch.
- **Action**:
    - Restore from latest snapshot via Velero: `velero restore create --from-backup snisid-daily-xxx`.

## 🛠️ Hardening Checklist
- [ ] No pods running as `root`.
- [ ] All container images signed via `Cosign`.
- [ ] `readOnlyRootFilesystem` enabled for all microservices.
- [ ] Automated vulnerability scanning in CI/CD (Trivy/Snyk).
