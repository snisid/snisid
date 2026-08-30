# PROMPT 255: AUTOSCALING SYSTEM ARCHITECTURE

This architecture defines the multi-layered autoscaling strategy for SNISID, combining reactive and predictive mechanisms to handle national-scale workload surges.

---

## 1. Scaling Architecture (Layered)

SNISID employs a three-tier autoscaling model:

- **Pod Layer (Reactive)**: Horizontal Pod Autoscaler (HPA) and Vertical Pod Autoscaler (VPA).
- **Event Layer (Dynamic)**: **KEDA** (Kubernetes Event-driven Autoscaling) for non-CPU/RAM triggers.
- **Node Layer (Infrastructure)**: **Karpenter** for just-in-time node provisioning.

---

## 2. Metrics Pipelines

Autoscaling decisions are powered by a high-resolution metrics pipeline:

- **Standard Metrics**: CPU/RAM usage via `metrics-server`.
- **Custom Metrics**: Kafka consumer lag, Flink checkpoint duration, and Neo4j query latency via Prometheus.
- **External Metrics**: National event triggers (e.g., sudden spike in biometric verification requests).

---

## 3. Scaling Policies

### Microservices
- **HPA Policy**: Scale up when `average_cpu_utilization > 60%` or `request_latency > 150ms`.
- **VPA Policy**: "Recommendation" mode for production to prevent OOM events; "Auto" mode for development environments.

### Data & AI Workloads
- **KEDA Kafka Scaler**: Automatically scales `fraud-engine` pods based on the number of unprocessed messages in the national Kafka topics.
- **GPU Scaling**: Karpenter provisions GPU nodes only when pods requiring `nvidia.com/gpu` are in a pending state.

---

## 4. Predictive Algorithms (AI-Driven)

SNISID integrates a **Predictive Scaling Engine** that analyzes historical patterns to stay ahead of demand.

- **Temporal Analysis**: Predicts surges during national events or business hours.
- **Look-Ahead Provisioning**: If a surge is predicted in 15 minutes, the engine pre-warms node pools and scales up pods before the traffic hits.
- **Anomaly Detection**: Filters out synthetic spikes (DDoS) from legitimate organic surges to prevent over-scaling.

---

## 5. Runtime Orchestration Model

1.  **Metric Collection**: Prometheus scrapes high-frequency data from all clusters.
2.  **Decision Loop**: HPA/KEDA evaluates the current state against policies.
3.  **Instruction**: Karpenter receives instructions to provision new nodes if existing ones are saturated.
4.  **Reconciliation**: The cluster state is updated within seconds, ensuring **Zero Downtime** for citizens.

---

**PROMPT 255 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 256 — POD SECURITY ARCHITECTURE.**
