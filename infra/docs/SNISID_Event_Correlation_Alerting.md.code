# SNISID: Event Correlation & SOC Alerting

The Correlation Engine stitches together disparate security signals into cohesive attack stories, enabling the National SOC to react to threats in seconds.

---

## 1. Distributed Event Correlation (Prompt 112)

Correlation occurs across multiple agencies and infrastructure layers.

### 1.1. Threat Chain Reconstruction
The engine identifies "Multi-Stage Attacks" by linking events via:
- **Identity Correlation**: Linking `login` events with `database_query` events via `identity_id`.
- **Trace Correlation**: Following a single `trace_id` from the API Gateway through the Service Mesh to the Backend.
- **Entity Correlation**: Linking events sharing a common IP address, Device ID, or national infrastructure node.

### 1.2. Intelligence Fusion
- **Cross-System Linking**: Flink joins Kafka topics from the `Auth`, `Network`, and `Database` domains.
- **Pattern Match**: Detects "Lateral Movement" (e.g., service A calling service B then service C in an unusual sequence).

---

## 2. Real-Time Alert Generation (Prompt 113)

Alerts are promoted from the stream into the **Sovereign SIEM (Elastic/Wazuh)**.

### 2.1. Alert Scoring & Deduplication
- **AI Alert Scorer**: A ML model evaluates the "Contextual Severity" of a correlated event chain. A failed login is minor; a failed login followed by a successful SQL injection attempt is **Critical**.
- **Deduplication**: Prevents "Alert Fatigue" by grouping 1000 identical port scans into a single "Network Reconnaissance" incident.

### 2.2. Escalation Workflows
1.  **Low/Med Alerts**: Routed to the standard SOC dashboard for manual triage.
2.  **High/Critical Alerts**: Triggers automated **PagerDuty/Slack/Email** notifications to the On-Call Security Officer.
3.  **Automatic Containment**: For P0 threats (e.g., active data exfiltration), the engine triggers an **OPA Lockdown** of the affected namespace before a human even reads the alert.

---

## 3. SOC Integration & Dashboards

- **Live Incident Stream**: A real-time, low-latency WebSocket feed into the SOC UI.
- **Attack Path Visualization**: Integrates with **Neo4j** to show the real-time graph of the attack (Source -> Entry Point -> Targets).
- **Threat Classification**: Maps every alert to the **MITRE ATT&CK** framework for national-level threat reporting.

---

## 4. Monitoring & Feedback

- **Alert-to-Incident Ratio**: Monitoring the effectiveness of the AI Scorer in reducing false positives.
- **Response Latency**: Tracking the time from `Threat_Detected` to `Action_Taken` (Target: < 1 minute for automated containment).
