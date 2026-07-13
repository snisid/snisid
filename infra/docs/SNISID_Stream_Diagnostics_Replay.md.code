# SNISID: Stream Diagnostics & Debugging Tooling

To maintain the reliability of the national event mesh, SNISID provides specialized tooling for live diagnostics, event tracing, and historical investigation.

---

## 1. Stream Debugging Architecture

The **Diagnostics Proxy** allows engineers to inspect live streams without affecting production consumers.

- **Shadow Consumption**: A dedicated "Shadow Consumer" group reads from production topics but does not commit offsets or trigger side-effects.
- **Event Sampling**: To handle massive throughput, the proxy can be configured to sample 0.1% of traffic or filter for a specific `identity_id` or `trace_id`.

---

## 2. Event Lineage & Tracing (Prompt 120)

Every event is part of a larger story.

- **W3C Trace Context**: Every Kafka message includes a `trace_id` and `parent_id` in the headers.
- **Lineage Tracking**: The **Sovereign Lineage Service** records the "Hop Count" of an event (e.g., Ingestion -> Normalization -> Fraud_Engine -> Alerting).
- **Lineage Visualization**: A graphical tool for developers to see the path a specific "Poison Pill" event took through the mesh.

---

## 3. Offset Inspection & Time-Travel Debugging

- **Offset Explorer**: A web interface to inspect the current offsets of all consumer groups across all national and regional clusters.
- **Time-Travel Debugging**: Integration with the **Forensic Replay Engine** allows a developer to "Seek" to a specific timestamp in the past and replay events into a local **Debug Sidecar** to reproduce a reported bug.

---

## 4. Operational Investigation Workflows

### 4.1. "Poison Pill" Handling
1.  **Detection**: A consumer fails to process an event after 3 retries.
2.  **Quarantine**: The event is moved to the `*.dlq` topic.
3.  **Inspection**: A developer uses the **Diagnostics Proxy** to view the raw payload and schema headers.
4.  **Fix & Re-Inject**: After fixing the downstream service, the event is re-injected into the main topic via a secure "Re-Injection API".

### 4.2. Stream Latency Profiling
- **Lag Heatmaps**: Visualizing which Kafka partitions are experiencing the highest consumer lag.
- **Processing Time Tracing**: Measuring the exact time (ms) an event spent in each Flink operator or microservice.

---

## 5. Security & Isolation

- **Read-Only Enforcement**: Debugging tools are strictly read-only and cannot write to production topics.
- **SVID Authentication**: Only developers with a valid **Sovereign Developer Identity** and "Diagnostic" clearance level can access the Shadow Consumer logs.
- **Masking**: PII is automatically masked in all debugging views, even for authorized engineers.
