# SNISID: SOC Intelligence & Alert Correlation

The Correlation & Prioritization Engine (CPE) acts as the "Decision Filter" for the National SOC, reducing noise and highlighting critical threat chains.

---

## 1. Real-Time Alert Correlation Engine (Prompt 152)

CPE reconstructs attack chains by fusing disparate security events in real-time.

### 1.1. Correlation Algorithms
- **Temporal Linking**: Grouping events from different sources that occur within a specific time window and share an `identity_id` or `source_ip`.
- **Relationship Linking (Neo4j)**: Correlating alerts that occur on nodes within 2 hops of each other in the national graph.
- **AI-Assisted Fusion**: An ML model analyzes the "Semantic Similarity" of alerts (e.g., a "SQL Injection" alert and an "Anomalous Data Egress" alert) to identify a cohesive attack story.

---

## 2. Alert Deduplication & Clustering (Prompt 158)

To prevent "Alert Storms" and analyst fatigue:

- **Similarity Clustering**: Identifies identical alerts from multiple sensors (e.g., 5 firewalls detecting the same scan) and collapses them into a single **Threat Group**.
- **Fingerprinting**: Every alert is assigned a **SHA-256 Hash** of its core attributes. Subsequent identical alerts simply increment the `Occurrence_Count` of the original incident.
- **Noise Reduction**: Automatically silences "Known Safe" patterns (e.g., scheduled backup traffic) while maintaining the underlying log for forensic audit.

---

## 3. Incident Classification & Scoring (Prompt 153)

Every correlated incident is assigned a **Sovereign Severity Score (SSS)**.

### 3.1. Severity Scoring Formula
$$SSS = (Base\_Severity \times Asset\_Criticality) + National\_Impact$$

- **Base Severity**: Technical severity (e.g., RCE = 10, Failed Login = 1).
- **Asset Criticality**: The importance of the target node (e.g., National ID Vault = 10x multiplier).
- **National Impact**: A dynamic modifier based on the current national security level (e.g., Election Day = +20).

---

## 4. Event Prioritization Engine (Prompt 163)

Prioritization ensures that the most dangerous threats are handled first.

- **Dynamic Queue Management**: The SOC Dashboard sorts incidents based on their SSS.
- **Autonomous Escalation**: If an incident's score exceeds 800 (Critical), the engine automatically triggers the **SOAR Response Agent** for instant containment.
- **Impact Estimation**: The AI engine predicts the "Potential Damage" (e.g., estimated data loss) if the incident is not remediated within the next 10 minutes.

---

## 5. Workflow Integration

1.  **Ingestion**: Raw alerts enter the Kafka SOC topic.
2.  **Deduplication**: Cluster identical alerts.
3.  **Correlation**: Link disparate alerts into a "Threat Story."
4.  **Scoring**: Calculate SSS and National Impact.
5.  **Prioritization**: Route to the correct Analyst Tier or Autonomous Agent.
