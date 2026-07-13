# SNISID National SIEM Architecture

## 1. Objective
To centralize the detection of cyber threats across the entire national digital infrastructure by correlating telemetry from disparate sources.

## 2. Technology Stack

| Domain | Recommended Technology | Role |
| :--- | :--- | :--- |
| **SIEM Core** | OpenSearch / Splunk | Indexing, Search, and Storage. |
| **Log Ingestion** | Fluentbit / Logstash | Collection and shipping of logs. |
| **Correlation** | Sigma Rules | Standardized detection patterns across different SIEMs. |
| **Analytics** | Detection Engine | Real-time anomaly detection. |
| **Dashboards** | Kibana / Grafana | Visualization of security posture and alerts. |

## 3. Ingestion Matrix (Mandatory Sources)

| Source | Obligatory | Priority | Logic/Focus |
| :--- | :---: | :---: | :--- |
| **Kubernetes Logs** | Yes | High | Pod crashes, unauthorized API calls, kubelet errors. |
| **IAM Logs** | Yes | Critical | Login failures, privilege escalations, MFA bypass attempts. |
| **Kafka Logs** | Yes | Medium | Message tampering, unauthorized topic access. |
| **BPMN Logs** | Yes | Medium | Process manipulation, unauthorized state transitions. |
| **Firewall Logs** | Yes | High | Port scanning, blocked connections, egress to C2 servers. |
| **Endpoint Logs** | Yes | High | Process creation, registry changes, file integrity. |

## 4. Data Pipeline Architecture
`Log Source` $\rightarrow$ `Fluentbit (Collector)` $\rightarrow$ `Kafka (Buffer)` $\rightarrow$ `Logstash (Normalization)` $\rightarrow$ `OpenSearch (Storage/Index)` $\rightarrow$ `Kibana (Visualization)`

## 5. Correlation & Detection Logic
- **Rule-Based Detection:** Using Sigma rules to detect known attack patterns (e.g., Brute Force, Pass-the-Hash).
- **Behavioral Detection:** Baseline creation for normal user activity; alerts on deviations.
- **Cross-Source Correlation:** Linking a firewall alert (external IP) to a Kubernetes pod log (internal activity).

## 6. Retention Policy
- **Hot Storage (30 Days):** Instant search for incident response.
- **Warm Storage (90 Days):** Slower search for trend analysis.
- **Cold Storage (1 Year+):** Compressed archival for legal/compliance requirements.
