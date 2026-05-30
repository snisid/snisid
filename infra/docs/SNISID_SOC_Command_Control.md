# SNISID: SOC Command & Control Architecture

The SOC Command & Control (C2) layer provides the unified interface and operational backbone for the National Security Operations Center, integrating human analysts with the autonomous defense swarm.

---

## 1. National SOC Topology (Prompt 151)

SNISID utilizes a **Federated SOC Model** to ensure regional autonomy and national oversight.

- **National Command Center (NCC)**: Aggregates high-severity incidents from all regions and coordinates national-scale responses.
- **Regional SOC Nodes (RSN)**: Perform local triage and autonomous containment for regional infrastructure.
- **Co-Located Agency SOCs**: Specialized nodes for the Police, Border Control, and Finance agencies that consume filtered intelligence streams.

---

## 2. SOC Dashboard & API Backend (Prompt 155)

The Dashboard provides real-time visibility into the national threat posture.

### 2.1. API Architecture
- **GraphQL Mesh**: Aggregates data from Kafka (Live Alerts), Neo4j (Graph Explorer), and PostgreSQL (Case History).
- **Streaming WebSockets**: Provides a sub-100ms feed of "Live Incidents" to the operator UI.
- **Query Optimization**: Utilizes **Materialized Views** and **Redis Caching** for high-speed dashboard rendering during massive ingestion bursts.

---

## 3. Case Management & Investigation (Prompt 157)

Investigations are managed as **Collaborative Security Objects**.

- **Incident Lifecycle**: `Detected` -> `Triaged` -> `Investigating` -> `Containing` -> `Remediated` -> `Post-Mortem`.
- **Evidence Management**: Every case is linked to an **Evidence Bundle** in Sovereign Object Storage, including packet captures (PCAPs), graph snapshots, and AI reasoning logs.
- **Collaboration**: Real-time analyst chat and investigation notebooks integrated with the **Forensic Replay Engine**.

---

## 4. Role-Based Operator System (Prompt 164)

Access is strictly governed by the **Sovereign Operator Hierarchy**.

| Role | Access Level | Responsibilities |
| :--- | :--- | :--- |
| **Tier 1 Analyst** | Read-Only (Masked) | Triage and basic evidence gathering. |
| **Tier 2 Investigator** | Read/Write (Regional) | Deep forensic analysis and incident containment. |
| **Tier 3 Commander** | Full Access (National) | Strategic response and national-level policy overrides. |
| **Agency Liaison** | Restricted (Agency-Specific) | Monitoring agency-specific threats and compliance. |

---

## 5. Multi-Channel Notification System (Prompt 161)

The Notification Engine ensures that critical alerts reach the right humans instantly.

- **Priority Routing**:
  - **P0 (Critical)**: SMS + Push + Dedicated SOC Alarm + Secure Voice Call.
  - **P1 (High)**: Push + Email + SOC Dashboard Alert.
  - **P2 (Standard)**: Email + Dashboard Log.
- **Escalation Logic**: If a P0 alert is not "Acknowledged" within 180 seconds, the engine automatically escalates to the **National Security Duty Officer**.

---

## 6. SOC Session Recording & Audit (Prompt 162, 165)

Every action taken by an operator is recorded for forensic integrity.
- **Session Replay**: Captures analyst screen and command-line interactions during an investigation.
- **Immutable Audit Trail**: Analyst actions are signed by their **SVID** and stored in the **Sovereign Audit Ledger**, making it impossible to "Delete the Logs" of an unauthorized investigation.
