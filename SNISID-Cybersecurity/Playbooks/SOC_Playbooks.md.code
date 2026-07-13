# SNISID SOC Playbooks

## 1. Objective
To industrialize the response to cyber incidents by providing standardized, step-by-step procedures for SOC analysts.

## 2. Playbook Catalog

| Playbook | Description | Priority | Trigger |
| :--- | :--- | :---: | :--- |
| **Ransomware Response** | Detection, Isolation, and Recovery from encryption attacks. | Critical | EDR Alert: File Encryption / Ransom Note. |
| **Identity Compromise** | Isolating compromised accounts and resetting credentials. | High | IAM Alert: Impossible Travel / MFA Failure. |
| **Insider Threat** | Investigation and containment of unauthorized internal activity. | High | DLP Alert: Mass Data Exfiltration. |
| **Kubernetes Breach** | Containing pod escapes and securing the K8s API. | Critical | Falco Alert: Shell in Container. |
| **PKI Compromise** | Revoking certificates and emergency root rotation. | Critical | Audit Log: Unauthorized Root Key Access. |
| **DDoS Attack** | Traffic scrubbing and edge filtering. | Medium | Network Alert: Inbound Traffic Spike. |

## 3. General Playbook Structure
Each playbook must follow this structure:
1. **Triage:** How to verify the alert is not a false positive.
2. **Containment:** Immediate steps to stop the spread (e.g., isolate host).
3. **Eradication:** Removing the threat (e.g., delete malicious files, kill processes).
4. **Recovery:** Restoring services from clean backups.
5. **Post-Mortem:** Documenting the timeline and identifying root causes.

## 4. Example: Ransomware Playbook (Simplified)
- **Step 1: Triage** $\rightarrow$ Verify file extensions changed and check for known ransomware signatures.
- **Step 2: Containment** $\rightarrow$ Use EDR to isolate affected endpoints from the network immediately.
- **Step 3: Analysis** $\rightarrow$ Identify the entry vector (e.g., Phishing, RDP) to prevent re-entry.
- **Step 4: Eradication** $\rightarrow$ Wipe affected systems; do not attempt to "clean" encrypted files.
- **Step 5: Recovery** $\rightarrow$ Restore from most recent offline backup.
- **Step 6: Notification** $\rightarrow$ Inform the National CERT and relevant government agencies.
