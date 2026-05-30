# SNISID: Stream Resilience & Recovery Architecture

The Sovereign Resilience Layer ensures that SNISID's real-time pipelines remain stable during load spikes and recover instantly from infrastructure failures.

---

## 1. Backpressure Handling & Flow Control (Prompt 97)

To prevent cascading failures during national-scale events, SNISID employs a **Multi-Level Backpressure System**.

### 1.1. Dynamic Throttling at Ingestion
- **Leaky Bucket Rate Limiting**: The API Gateway throttles incoming events per agency/source based on real-time Kafka partition lag.
- **Priority-Based Shedding**: If the Kafka backbone enters a "Critical Congestion" state, low-priority events (e.g., `system.heartbeat`) are dropped to preserve bandwidth for high-priority events (e.g., `identity.revocation`).

### 1.2. Consumer-Side Flow Control
- **Auto-Scaling (HPA)**: Kubernetes TaskManagers (Flink) and Pods (Kafka Streams) scale horizontally based on the `consumer_lag` metric.
- **Poll Management**: Consumers utilize `max.poll.records` and `fetch.max.bytes` to prevent memory exhaustion when processing sudden bursts of large payloads.

---

## 2. Stream Failure Recovery (Prompt 105)

The platform is designed to resume processing with absolute state consistency after a crash.

### 2.1. Offset & State Recovery
- **Distributed Checkpointing (Flink)**: Flink snapshots the entire state (including RocksDB and Kafka offsets) every 60 seconds to **Sovereign Object Storage**.
- **Changelog Replay (Kafka Streams)**: Local state stores are backed by Kafka changelog topics. If a node fails, the new instance re-hydrates its state from the changelog.

### 2.2. Automatic Failover Orchestration
- **Active-Active Standby**: High-criticality streams (e.g., Border Monitoring) maintain a standby instance in a different Availability Zone with a replicated state store.
- **Failover Trigger**: If the Prometheus/Alertmanager detects a service is down for > 30 seconds, the Global Load Balancer redirects traffic, and the standby instance promoted to `Active`.

---

## 3. Buffer Management & Priority Queues

- **Tiered Buffering**: 
  - **Memory (L1)**: Instant processing.
  - **SSD (L2)**: Kafka broker local disk (7-day buffer).
  - **S3 (L3)**: Tiered Storage for infinite retention.
- **Isolation**: Critical security topics have dedicated Kafka broker node-pools to prevent "Noisy Neighbor" interference from heavy analytical workloads.

---

## 4. Operational Recovery Tooling

- **Rebalance Monitor**: Real-time visualization of Kafka partition rebalancing.
- **Checkpoint Inspector**: Tooling to verify the integrity and timestamp of the last successful Flink checkpoint before manual restoration.
- **Stream Integrity Validator**: Background process that compares source Kafka offsets with processed Sink timestamps to detect missing events.
