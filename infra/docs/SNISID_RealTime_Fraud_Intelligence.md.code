# SNISID: Real-Time Fraud & Anomaly Intelligence

The Fraud Intelligence layer utilizes Flink's Complex Event Processing (CEP) and stateful analytics to identify and contain threats at national-scale speeds.

---

## 1. Fraud Detection Streaming Job (Prompt 107)

The `Fraud_Engine_Job` is a stateful Flink application that correlates cross-agency event streams.

### 1.1. Processing Pipeline
1.  **Ingestion**: Consumes from `national.auth.v1` and `national.banking.v1` Kafka topics.
2.  **Enrichment**: Side-lookup to the **Identity Confidence Cache**.
3.  **CEP Evaluation**: Runs the **Pattern Matcher** (Section 3).
4.  **Risk Scoring**: Calculates a real-time `fraud_score` (0-1000).
5.  **Alerting**: High-risk events (score > 800) are emitted to `soc.alert.fraud`.

---

## 2. Sliding Window Analytics (Prompt 111)

To detect temporal fraud patterns (e.g., "Credential Stuffing"), the engine employs **Multi-Window Analysis**.

- **Tumbling Window (1m)**: Measures volume spikes (e.g., "1000 failed logins in 1 minute").
- **Sliding Window (5m / 10s step)**: Tracks velocity and behavioral drift (e.g., "Gradual increase in session duration from anomalous geolocations").
- **Global Window (Identity-Based)**: Aggregates lifetime risk flags for a specific `identity_id`.

---

## 3. Complex Event Processing (CEP) Engine (Prompt 114)

The CEP engine detects **Behavioral Chain Attacks** through pattern matching.

```java
// Example: Detect suspicious elevation chain
Pattern<Event, ?> pattern = Pattern.<Event>begin("login")
    .where(new SimpleCondition<Event>() {
        @Override
        public boolean filter(Event event) { return event.type == LOGIN_SUCCESS; }
    })
    .followedBy("elevation_request")
    .where(new SimpleCondition<Event>() {
        @Override
        public boolean filter(Event event) { return event.type == PRIVILEGE_ESC; }
    })
    .within(Time.minutes(2));
```

- **Pattern States**: `Begin` -> `FollowedBy` -> `NotFollowedBy` -> `Within`.
- **Match Logic**: Successfully matched chains are instantly promoted to **High-Priority Incidents**.

---

## 4. Anomaly Detection Pipeline (Prompt 108)

The Anomaly Detection service identifies "Unknown Unknowns" through statistical deviation.

### 4.1. Adaptive Threshold Engine
- **Baseline**: The engine maintains a rolling baseline of "Normal Behavior" for each agency and identity class (stored in Flink state).
- **Z-Score Calculation**: Calculates the standard deviation of the current event's metrics (e.g., payload size, frequency) against the baseline.
- **AI Scoring**: If the Z-Score exceeds a threshold (e.g., 3.0), the event is sent to the **AI Threat Engine** for deep neural analysis.

---

## 5. Runtime Monitoring & Feedback Loop

- **Backpressure Monitoring**: If the Fraud Job slows down, the **API Gateway** automatically increases rate-limiting on non-critical ingestion sources.
- **Model Feedback**: Real-time results are fed back into the **Graph Intelligence (Neo4j)** to update relationship weights and fraud ring clusters.
