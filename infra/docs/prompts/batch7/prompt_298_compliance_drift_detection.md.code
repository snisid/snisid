# PROMPT 298: AUTOMATED INFRASTRUCTURE COMPLIANCE DRIFT DETECTION

This architecture defines the continuous, automated strategy for identifying and reporting deviations from the platform's compliance baseline, ensuring that the SNISID infrastructure always adheres to national security and regulatory mandates.

---

## 1. Detection Architecture (Continuous Auditing)

SNISID utilizes a "Query-First" auditing stack that treats the infrastructure as a searchable database.

- **Auditor (Steampipe)**: Periodically queries all cloud and Kubernetes APIs using SQL to verify compliance against predefined security frameworks (e.g., CIS, NIST).
- **Rule Engine (CloudCustodian)**: Enforces "Compliance Policies" (e.g., "All EBS volumes must be encrypted", "No security group can allow 0.0.0.0/0 on port 22").
- **State Analyzer**: Correlates live infrastructure state with the "Signed Compliance Baseline" stored in the security repository.
- **Risk Scorer**: AI-driven engine that assigns a "Compliance Severity" to every detected drift based on the classification of the affected assets.

---

## 2. Workflows (The Auditing Cycle)

1.  **Scheduled Scans**: Full infrastructure compliance scans executed every 6 hours across all regional clusters.
2.  **Event-Driven Detection**: Cloud-provider events (e.g., "CreateBucket", "ModifySecurityGroup") trigger near-instant compliance checks for the affected resource.
3.  **Gap Analysis**: The system identifies the "Drift" (e.g., "A database was created without an automated backup policy").
4.  **Alerting & Notification**: High-severity compliance drifts trigger immediate alerts in the **Alerting System** (Prompt 280) and log an entry in the forensic ledger.

---

## 3. Analysis & Reporting

- **Compliance Posture Dashboard**: Real-time view in the **Developer Portal** (Prompt 282) showing the compliance percentage for each agency and project.
- **Historical Drift Timeline**: Visualizes how the platform's compliance has evolved over time, identifying recurring "Compliance Weak Spots."
- **Executive Audit Report**: Automated weekly summary for national security leadership, highlighting critical gaps and remediation progress.

---

## 4. Integration Model (The Remediation Bridge)

- **Auto-Correction Integration**: Connected to Prompt 288 (Drift Auto-Correction) to automatically revert non-compliant changes where safe and authorized.
- **Ticket Generation**: Lower-priority compliance issues automatically generate "Remediation Tasks" in the project's issue tracker.
- **Policy Feedback Loop**: Recurring compliance drifts are used as signals to update the "Secure Scaffolding" (Golden Paths) in the Developer Portal to prevent future occurrences.

---

## 5. Governance Strategy

- **Compliance Thresholds**: Projects that fall below a 95% compliance score are automatically restricted from deploying new resources until the gaps are remediated.
- **Audit Ledger**: Every scan result, detected drift, and remediation action is cryptographically signed and stored for long-term regulatory verification.
- **Sovereign Regulatory Mapping**: The system automatically maps technical drifts to specific national laws and international standards (e.g., GDPR, National Security Directive 12).

---

**PROMPT 298 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 299 — AUTOMATED INFRASTRUCTURE SECURITY TESTING (DAST).**
