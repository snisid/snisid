# SNISID National SOC Architecture

## 1. Overview
The National Security Operations Center (SOC) is the central nervous system for cyberdefense of the SNISID infrastructure. Its primary goal is to provide 24/7 monitoring, detection, analysis, and response to cyber threats targeting national digital assets.

## 2. Mission Statement
To ensure the continuous availability, integrity, and confidentiality of SNISID services through proactive monitoring and rapid incident response, maintaining national digital sovereignty.

## 3. Operational Model (24/7/365)
The SOC operates on a continuous rotation basis to ensure zero gaps in surveillance.

### Shift Structure
- **Shift A (Day):** 06:00 - 14:00
- **Shift B (Evening):** 14:00 - 22:00
- **Shift C (Night):** 22:00 - 06:00
- **On-Call/Escalation:** Senior Analysts and Incident Commanders available 24/7.

## 4. Core Functional Components

| Function | Description | Key Activities |
| :--- | :--- | :--- |
| **Monitoring** | Real-time surveillance | Log analysis, alert triage, dashboard monitoring. |
| **Incident Response** | Containment and Eradication | Triage, isolation, remediation, recovery. |
| **Threat Hunting** | Proactive search for threats | Hypothesis-based searching for undetected intruders. |
| **DFIR** | Digital Forensics & Incident Response | Root cause analysis, memory/disk forensics, legal preservation. |
| **SIEM** | Correlation and Alerting | Log aggregation, rule-based correlation, alert generation. |
| **Crisis Operations** | High-level coordination | Government liaison, strategic decision making during major breaches. |

## 5. Logical Architecture
The SOC is integrated into the SNISID ecosystem as follows:
- **Data Collection Layer:** Agents (Fluentbit, Wazuh) on all nodes.
- **Analysis Layer:** SIEM (OpenSearch) and Correlation engines.
- **Human Layer:** Tier 1 (Triage) $\rightarrow$ Tier 2 (Analysis) $\rightarrow$ Tier 3 (Hunting/Forensics).
- **Response Layer:** SOAR (Security Orchestration, Automation, and Response) and Manual intervention.

## 6. KPI & Metrics
- **MTTD (Mean Time To Detect):** Target < 15 minutes for critical alerts.
- **MTTR (Mean Time To Respond):** Target < 1 hour for critical containment.
- **False Positive Rate:** Target < 20%.
- **Coverage:** 100% of critical SNISID assets.

## 7. Governance & Integration
The SOC reports directly to the National Cyber Governance body and coordinates closely with the National CERT.
