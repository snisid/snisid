# SNISID National Cyber Standards

## 1. Objective
To standardize security configurations across all government and critical infrastructure entities to ensure a consistent defense baseline.

## 2. Core Standards

| Domain | Standard Requirement | Compliance Level |
| :--- | :--- | :---: |
| **IAM Security** | Mandatory MFA for all accounts; No shared accounts; 90-day password rotation (if not MFA). | Mandatory |
| **Kubernetes Security** | No privileged containers; NetworkPolicies by default (Deny All); Image scanning in CI/CD. | Mandatory |
| **API Security** | All APIs must use HTTPS; OAuth2/OIDC for auth; Rate limiting on all public endpoints. | Mandatory |
| **Logging** | All logs must be sent to National SIEM in JSON format; Timestamp in UTC. | Mandatory |
| **Incident Severity** | Standardized levels (Low, Medium, High, Critical) based on impact and urgency. | Mandatory |
| **Threat Classification** | Use of MITRE ATT&CK for all incident reporting. | Mandatory |

## 3. Incident Severity Matrix

| Level | Impact | Response Time | Notification |
| :--- | :--- | :---: | :--- |
| **Low** | Minor anomaly; No data loss. | < 24 Hours | SOC Internal |
| **Medium** | Single system affected; Service degradation. | < 4 Hours | SOC $\rightarrow$ CERT |
| **High** | Multiple systems affected; Sensitive data leak. | < 1 Hour | CERT $\rightarrow$ C4 |
| **Critical** | National service outage; Root compromise. | Immediate | C4 $\rightarrow$ Government |

## 4. Compliance Auditing
- **Automated Audits:** Weekly scans using OpenSCAP or similar tools to check configuration drift.
- **Manual Audits:** Annual third-party security audits.
- **Penetration Testing:** Bi-annual mandatory Red Team exercises.

## 5. Enforcement
Non-compliance with these standards results in the immediate revocation of network access for the offending entity until the vulnerability is remediated.
