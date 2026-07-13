# SNISID: Autoscaling & Resilience Strategy

This document defines the autoscaling and self-healing mechanisms for SNISID, ensuring that the national infrastructure can handle population-scale bursts while maintaining continuous availability.

---

## 1. Predictive & Event-Driven Autoscaling (KEDA)

Traditional CPU/Memory based autoscaling is insufficient for the bursty nature of national intelligence workloads. We use **KEDA (Kubernetes Event-driven Autoscaling)** for proactive scaling.

- **Kafka-Driven Scaling**: Scaling the `identity-ingestion` and `fraud-scoring` services based on the number of pending messages (lag) in specific Kafka topics.
- **Flink Elasticity**: Automatically adjusting the number of Flink TaskManagers based on the current data throughput and watermarks.
- **AI/GPU Scaling**: Proactively scaling GPU node pools based on the volume of biometric verification requests at border kiosks during peak travel times.
- **Predictive Algorithms**: Using historical traffic data to "Pre-Scale" the cluster 30 minutes before known high-demand periods (e.g., national election days or major public service renewals).

---

## 2. Cluster Self-Healing Mechanisms

- **Automated Node Repair**: Integration with **Cluster Autoscaler** and cloud-native health checks to automatically terminate and replace unresponsive or degraded worker nodes.
- **Liveness & Readiness Probes**: Every SNISID container includes granular probes that check the health of internal components (e.g., the Biometric Engine checks its connection to the local matcher before reporting as "Ready").
- **Automatic Pod Restart**: Kubernetes automatically restarts pods that fail their liveness checks, with **Exponential Backoff** to prevent "CrashLoopBackOff" storms.

---

## 3. High Availability (HA) & Disaster Recovery

- **Multi-AZ / Multi-Region Deployment**: Critical services are spread across different Availability Zones and geographic regions to prevent single points of failure.
- **Pod Disruption Budgets (PDB)**: Ensuring that no more than 10% of critical pods (e.g., Auth/API Gateway) are offline at any given time during a maintenance window or node upgrade.
- **Stateful Persistence**: Using distributed storage (e.g., **Ceph** or **Portworx**) with synchronous replication to ensure no data loss during a storage node failure.

---

## 4. Performance & Resource Optimization

- **Vertical Pod Autoscaler (VPA)**: Automatically adjusting the CPU and Memory requests/limits for microservices based on actual usage patterns over time.
- **Priority & Preemption**: Assigning `Critical` priority to the core identity and border services, allowing them to preempt less critical batch processing workloads during a resource crunch.
- **Node Affinity & Tainting**: Ensuring that high-performance AI workloads only run on GPU-equipped nodes, and that sensitive cryptographic services run on specialized hardened nodes.

---

## 5. Observability & SLIs

- **Service Level Indicators (SLIs)**: Monitoring latency, error rate, throughput, and saturation for every critical path.
- **Automated Alerts**: Real-time alerting via Prometheus/Grafana if an SLI exceeds its **Service Level Objective (SLO)**.
- **Post-Mortem Integration**: Every major scaling or failure event triggers an automated "Incident Report" stored in the **Sovereign Audit Ledger**.
