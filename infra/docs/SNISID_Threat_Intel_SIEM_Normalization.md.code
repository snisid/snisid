# SNISID: Threat Intel & SIEM Normalization

The Normalization and Intelligence layer provides the "Unified Language" for the SNISID SOC, ensuring that data from thousands of diverse sources can be analyzed and correlated instantly.

---

## 1. Security Event Normalization (Prompt 156)

SNISID uses a **Sovereign Common Event Schema (SCES)** to normalize multi-format logs.

- **Normalization Workflow**:
  1. **Capture**: Ingest raw logs (JSON, Syslog, Netflow, CEF).
  2. **Parse**: Extract core fields (Timestamp, Source, Target, Action, Result).
  3. **Map**: Align fields to SCES (e.g., mapping `src_ip`, `source`, and `ip_src` all to `network.source.ip`).
  4. **Enrich**: Inject metadata (Geo-IP, Asset Criticality, Threat Intel status).
- **Format**: All normalized events are emitted as **Protobuf** messages to the Kafka backbone.

---

## 2. Threat Intelligence Ingestion (Prompt 154)

The Threat Intel Platform (TIP) aggregates global and national signals to provide proactive defense.

- **Feed Integration**:
  - **Internal**: Indicators of Compromise (IOCs) from regional SOC nodes.
  - **External**: STIX/TAXII feeds from international partners and commercial providers.
  - **Open Source**: Automated scraping of threat reports and vulnerability databases.
- **AI-Assisted Scoring**: AI evaluates the "National Relevance" of an IOC. A malware hash targeting European banking is scored higher than a generic IoT botnet signature.
- **Real-Time Injection**: Validated IOCs are instantly pushed to the **Flink Correlation Engine** for live stream matching.

---

## 3. SOC Metrics Engine (Prompt 160)

The Metrics Engine calculates the operational effectiveness of the national defense.

- **KPIs**:
  - **MTTD (Mean Time To Detect)**: Time from `Threat_Injected` to `Alert_Generated`.
  - **MTTR (Mean Time To Respond)**: Time from `Alert_Generated` to `Incident_Contained`.
  - **True Positive Rate**: Percentage of alerts that were actual security incidents.
  - **Autonomous Coverage**: Percentage of incidents handled by the AI Swarm without human intervention.
- **Dashboard**: Real-time temporal trend analysis for the National Security Director.

---

## 4. Forensic Audit Trail (Prompt 162)

To ensure the integrity of security operations, every SOC action is cryptographically preserved.

- **Immutable Traceability**: Every investigation step, from "Case Opened" to "Pod Quarantined," is recorded in the **Sovereign Audit Ledger**.
- **Session Correlation**: Links analyst commands and dashboard clicks to the specific `Investigation_ID` and `Timestamp`.
- **Tamper Detection**: Periodic Merkle-tree validation checks to ensure historical audit logs haven't been altered by a compromised administrator.

---

## 5. Intelligence Distribution

- **Feed Distribution**: Validated intelligence is shared back with regional agencies (Police, Banks) via a "Sovereign Threat Broadcast" topic.
- **Sovereign Blacklist**: A high-speed, mTLS-secured list of blocked IPs/Hashes consumed by the **Sovereign API Gateway** and **Istio Mesh**.
