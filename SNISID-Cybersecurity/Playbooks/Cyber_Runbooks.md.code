# SNISID National Cyber Runbooks

## 1. Objective
To provide operational instructions for the maintenance and recovery of the security infrastructure itself, ensuring the defense system remains resilient.

## 2. Operational Runbooks

| Runbook | Description | Trigger |
| :--- | :--- | :--- |
| **SOC Outage Recovery** | Restoring the SOC monitoring capability after a failure. | SIEM/Dashboard Unavailability. |
| **SIEM Overload** | Managing log ingestion spikes to prevent data loss. | Log queue saturation / Disk pressure. |
| **PKI Compromise Emergency** | Step-by-step guide to rotating the National Root CA. | Root Key theft/compromise. |
| **Data Breach Containment** | Technical steps to shut down data egress points. | Confirmed mass exfiltration. |
| **Nation-State Attack Crisis** | High-level switch to "War Footing" operational mode. | C4 Level 4 Activation. |

## 3. Runbook vs. Playbook
- **Playbook:** "How to deal with an attacker" (Threat-centric).
- **Runbook:** "How to fix the tool" (System-centric).

## 4. Example: SIEM Overload Runbook
1. **Identify Source:** Use Grafana to find the log source causing the spike.
2. **Apply Rate Limiting:** Implement temporary throttling at the Fluentbit layer.
3. **Expand Storage:** Dynamically increase OpenSearch data nodes.
4. **Drop Low-Priority Logs:** Temporarily disable "Debug" or "Info" level logs for non-critical systems.
5. **Verify Stability:** Monitor CPU/Memory until levels return to baseline.

## 5. Maintenance Schedule
- **Quarterly Review:** All runbooks are reviewed and updated.
- **Dry Runs:** Once a year, a "Tool Failure" is simulated to test the runbooks.
