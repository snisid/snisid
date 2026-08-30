# SNISID: Observability Intelligence Stack

The Observability Intelligence Stack provides the "Sensory Nervous System" for SNISID, ensuring that every operational metric, log, and trace is collected and analyzed to maintain system health and security.

---

## 1. Federated Monitoring (Prometheus)

Monitoring is sharded across regions to ensure scalability and local visibility.

- **Regional Prometheus**: Each regional cluster runs a dedicated Prometheus instance for local metric collection and alerting.
- **Thanos/Cortex**: A centralized "Global View" layer that aggregates metrics from all regional Prometheus instances for national-scale reporting.
- **Node Exporter & Kube-State-Metrics**: Collecting deep infrastructure metrics, including GPU utilization, disk I/O, and pod lifecycle events.

---

## 2. Distributed Logging (Grafana Loki)

Logging is optimized for high volume and multi-tenancy.

- **Promtail Agents**: Shipped as a DaemonSet on every node to collect logs from container stdout/stderr.
- **Loki Indexing**: Logs are indexed by labels (e.g., `agency`, `region`, `app`) allowing for sub-second searches across billions of lines without the cost of full-text indexing.
- **Log Retention Policies**: Regional clusters retain detailed logs for 30 days, while summarized "Security Events" are streamed to the **National Intelligence Fusion Center** for long-term retention.

---

## 3. Distributed Tracing (Jaeger/Tempo)

To debug performance bottlenecks in the complex microservices mesh, we use distributed tracing.

- **Istio Integration**: Istio sidecars automatically inject trace headers into all cross-service requests.
- **Jaeger/Tempo Collectors**: Collecting span data from the mesh and storing it in the **Sovereign Object Storage**.
- **Service Graph Visualization**: Real-time mapping of service dependencies and latency, identifying slow components in critical paths like biometric verification.

---

## 4. AI-Native Observability

- **Anomaly Detection**: Using the **SNISID Adaptive AI** to monitor metric streams and automatically flag departures from normal operating baselines (e.g., an unexpected spike in 5xx errors at a specific border gateway).
- **Automated Root Cause Analysis**: Correlation of traces and logs during an incident to provide analysts with a "Timeline of Failure."

---

## 5. Security & Sovereignty

- **Log Redaction**: Automated filtering of PII (National IDs, Biometric hashes) from logs before they are written to disk.
- **Sovereign Access Control**: Access to Grafana dashboards and Prometheus queries is restricted via **Zero Trust Identity**, with full auditing of every query.
- **Audit Ledger Integration**: Critical alerts and system health reports are cryptographically signed and stored in the **Sovereign Audit Ledger**.
