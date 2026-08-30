# SNISID: SOC Orchestration & Automated Playbooks

The Security Orchestration, Automation, and Response (SOAR) layer provides the standardized workflows for managing national-scale cyber incidents, ensuring a consistent and rapid response to complex threats.

---

## 1. SOAR Orchestration Engine

The SOAR engine integrates the **Kafka SOC Stream** with automated fulfillment systems and human-in-the-loop dashboards.

- **Workflow Orchestrator**: Uses **Temporal** or **Airflow** to manage long-running, multi-stage response workflows.
- **Integration Hub**: A library of "Connectors" that allow the SOAR engine to interact with:
  - **Kubernetes API**: For pod management and scaling.
  - **Vault API**: For credential rotation.
  - **MinIO API**: For bucket locking and WORM policy enforcement.
  - **Notification Services**: (SMS, Email, SOC Dashboard).

---

## 2. Automated Incident Playbooks

SNISID utilizes pre-approved **Standard Operating Playbooks (SOPs)** for common high-impact scenarios.

### 2.1. Playbook: Ransomware Containment (SOP-001)
1.  **Detection**: AI Anomaly engine detects mass encryption behavior (high entropy) in Object Storage.
2.  **Snapshot**: Trigger an immediate **Sovereign Immutable Backup** of the affected namespace.
3.  **Lockdown**: Switch the storage bucket to **"Strict WORM"** mode (Deny All Deletes).
4.  **Isolation**: Revoke SVIDs of all workloads with write access to that bucket.
5.  **Audit**: Export a forensic evidence package to the **Sovereign Audit Ledger**.

### 2.2. Playbook: Data Exfiltration Response (SOP-002)
1.  **Detection**: Flink monitors egress bandwidth spikes to unknown jurisdictions.
2.  **Throttling**: Automatically reduce egress bandwidth for the affected CIDR range by 90% (Cilium).
3.  **Verification**: Coordinator tasks an Investigator Agent to verify the data sensitivity.
4.  **Shutdown**: If data is PII-sensitive, the egress route is permanently severed until manual review.

### 2.3. Playbook: Admin Account Takeover (SOP-003)
1.  **Detection**: MFA bypass or "Impossible Travel" detected for a privileged identity.
2.  **Revocation**: Global invalidation of all active JWT and Refresh tokens for that identity.
3.  **Credential Reset**: Trigger an automated **Secret Rotation** for all credentials held by that identity in Vault.
4.  **Re-Verify**: Force a "Four-Eyes" physical biometric verification for the user to regain access.

---

## 3. Incident Correlation & Case Management

- **The Incident Store**: Every SOAR execution creates a "Case" in the **Sovereign SIEM**.
- **Evidence Stitching**: The engine automatically attaches Neo4j graph snapshots and Flink telemetry to the case file.
- **Reporting**: Generates automated **ISO 27001 / GDPR** compliance reports for every national-scale incident.

---

## 4. Human-in-the-Loop Guardrails

- **Manual Override**: SOC Analysts can "Pause" or "Reverse" any automated playbook action from the master dashboard.
- **Approval Gates**: For "National Level" containment (e.g., shutting down an entire regional spoke), the SOAR engine pauses and waits for an **Authorized Signature** from the National Security Director.

---

## 5. Performance Targets (SLA)

| Action | Target Latency |
| :--- | :--- |
| **Initial Containment (Quarantine)** | < 30 seconds |
| **Credential Revocation** | < 10 seconds |
| **Case Initialization & Evidence Fetch** | < 2 minutes |
| **National Escalation Notification** | < 1 minute |
