# SNISID: National KPI Engine & Temporal Analysis

The KPI Engine provides the Government of SNISID with a real-time "Pulse" of the nation, while the Temporal Analysis engine reconstructs the chronology of complex security events.

---

## 1. National KPI Computation Engine (Prompt 118)

The KPI Engine aggregates billions of events into high-level metrics for the **National Executive Dashboard**.

### 1.1. Metric Categories
- **Identity Metrics**: Total active citizens, enrollment velocity, regional distribution.
- **Security Metrics**: Active threat alerts, successful block rate, average response time.
- **Fraud Metrics**: Prevented fraud value, active fraud ring detections, biometric mismatch rates.
- **Operational Metrics**: Kafka throughput, Flink lag, API latency across all agencies.

### 1.2. Aggregation Workflows
- **Real-Time Tier**: Flink processes sliding windows (1m, 5m, 1h) to provide live dashboard updates via WebSockets.
- **Historical Tier**: Daily aggregates are moved to the **Sovereign Analytics Warehouse** for trend analysis over years.

---

## 2. Temporal Event Analysis (Prompt 119)

To understand the "How" and "When" of an attack, SNISID uses **Chronological Reconstruction**.

### 2.1. Sequence Analysis
- **Event Chronology**: The engine uses the `Sovereign_Timestamp` to re-order events arriving from different agencies, ensuring that "Cause" is always placed before "Effect" in the analysis.
- **Time-Based Correlation**: Identifying events that occur within a specific temporal window (e.g., "MFA failure in Region A followed by a privilege elevation in Region B within 30 seconds").

### 2.2. Delayed-Event Handling
- **Graceful Arrival**: Flink's **Watermark** mechanism allows for late-arriving data (up to 5 seconds) to be integrated into the temporal window before the state is finalized.
- **Retrospective Correction**: For extremely late data, the engine triggers a "Correction Event" that updates the historical KPI counts.

---

## 3. Temporal Graph Correlation

By linking the Temporal Engine with **Neo4j**, SNISID can visualize the "Evolution of a Threat":
1.  **T0**: Initial phish detection.
2.  **T+1h**: Account compromise.
3.  **T+4h**: Lateral movement to the database.
4.  **T+5h**: Data exfiltration attempt detected by Flink.

---

## 4. Scalability & Performance

- **Pre-Aggregation**: To handle national-scale queries, KPIs are pre-aggregated at the **Regional Spoke** level before being sent to the National Center.
- **State Partitioning**: KPI state is partitioned by `agency_id` and `region_id` to allow for massive parallelism in the Flink cluster.
