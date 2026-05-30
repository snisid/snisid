# SNISID: Cyber Threat Detection Backbone

The Detection Backbone provides the specialized engines that identify specific attack patterns, ranging from low-level network intrusions to sophisticated behavioral anomalies.

---

## 1. Multi-Layer Detection Architecture

SNISID utilizes a **Defense-in-Depth Detection Strategy** across the entire national infrastructure.

### 1.1. Real-Time IDS/IPS (Prompt 166)
- **Network IDS (NIDS)**: Cilium-based eBPF sensors monitor traffic between services, detecting protocol anomalies and known exploit patterns.
- **Host IDS (HIDS)**: Falco-based sensors on Kubernetes nodes monitor system calls (e.g., unexpected shell executions or file modifications).
- **Alerting**: Detections are instantly emitted to the `national.soc.ids` Kafka topic.

---

## 2. Detection Engines

### 2.1. Signature-Based Detection (Prompt 168)
- **IOC Matching**: High-speed matching of IPs, Domains, and Hashes against the **Sovereign Threat Intel** feed.
- **Malware Signatures**: Integration with Yara for scanning uploaded objects and forensic memory dumps in real-time.
- **Update Pipeline**: Signatures are updated hourly across the national mesh via the **Sovereign Distribution System**.

### 2.2. Anomaly-Based Detection (Prompt 167)
- **Behavioral Baselines**: Every identity and service has a "Normal" behavioral profile (e.g., "Standard API calling frequency").
- **Statistical Deviation**: Flink identifies Z-score deviations in throughput, latency, and payload size.
- **Adaptive Thresholds**: AI models adjust thresholds based on time-of-day, national holidays, and regional load patterns.

### 2.3. Behavioral Analysis Engine (Prompt 169)
- **User Behavior (UBA)**: Detects credential abuse, insider threats, and lateral movement by analyzing the sequence of an identity's actions.
- **Entity Behavior (EBA)**: Monitors service-to-service communication patterns to detect compromised workloads.

---

## 3. Specialized Detection Pipelines (Prompts 170-180)

| Pipeline | Detection Mechanism |
| :--- | :--- |
| **Malware Pattern** | Heuristic analysis of system call sequences and process behavior. |
| **Phishing Detection** | AI-based analysis of email/SMS headers and link reputation. |
| **Lateral Movement** | Detection of anomalous cross-namespace service mesh requests. |
| **Ransomware** | Real-time monitoring of file entropy and mass-deletion actions in S3. |
| **Data Exfiltration** | Outbound bandwidth anomaly detection and DPI (Deep Packet Inspection). |
| **Zero-Day Anomaly** | Identification of novel traffic patterns that deviate from SCES norms. |

---

## 4. Integration & Orchestration

- **Kafka/Flink Integration**: All engines are producers to the **Kafka Event Backbone**. Flink performs multi-engine correlation (e.g., "HIDS alert + NIDS alert = High Confidence Incident").
- **SOC Orchestration**: Every high-confidence detection triggers an **Autonomous SOC Agent Swarm** investigation.
- **Forensic Traceability**: Every detection includes the raw packet/log snippet as evidence for post-incident review.

---

## 5. Deployment & Resilience

- **Kubernetes Strategy**: Detection sensors are deployed as **DaemonSets** on every node to ensure 100% coverage.
- **Resource Isolation**: Detection engines run in dedicated namespaces with resource quotas to prevent them from interfering with production identity services.
- **Runtime Resilience**: If a sensor fails, the **Sovereign Monitoring System** triggers an automatic restart and alerts the Regional SOC of a "Visibility Gap."
