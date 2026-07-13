# SNISID: Real-Time vs. Batch Processing Strategy

To guarantee that heavy analytical queries do not compromise the millisecond latency required for national security operations at the border, SNISID enforces strict isolation between Real-Time (Hot Path) and Batch (Cold Path) processing.

---

## 1. Real-Time vs. Batch Architecture Diagram (Lambda/Kappa Architecture)

SNISID utilizes an event-driven architecture that naturally bifurcates into a "Hot Path" for immediate action and a "Cold Path" for deep analytics and training.

```mermaid
graph TD
    %% Event Source
    Kafka[Apache Kafka Event Backbone]

    %% HOT PATH (Real-Time)
    subgraph "Hot Path (Real-Time / Streaming)"
        Flink[Apache Flink]
        Redis[(Redis - Low Latency State)]
        Neo4j[(Neo4j - Graph Queries)]
        SOC[SOC Alerting Engine]
        Fraud[Fraud Scoring Engine]
        
        Flink -->|Sub-second Aggregation| Redis
        Flink -->|Critical Anomalies| SOC
        Fraud -->|Query| Neo4j
    end

    %% COLD PATH (Batch)
    subgraph "Cold Path (Batch / Analytics)"
        Spark[Apache Spark]
        S3[(MinIO / S3 Data Lake)]
        DWH[(Data Warehouse / BI)]
        MLOps[AI Retraining Pipeline]
        
        Kafka_Connect[Kafka Connect (Sink)] -->|Tiered Storage| S3
        Spark -->|Hourly/Daily Jobs| S3
        Spark -->|Aggregations| DWH
        S3 -->|Historical Replay| MLOps
    end

    Kafka -->|Consumer Group A (High Priority)| Flink
    Kafka -->|Consumer Group B (Low Priority)| Fraud
    Kafka -->|Consumer Group C (Bulk)| Kafka_Connect
```

---

## 2. Processing Decision Matrix

The following matrix dictates exactly which processing paradigm must be used based on the workload type.

| Workload Type | Classification | Processing Engine | Target SLA | Primary Storage Target |
| :--- | :--- | :--- | :--- | :--- |
| **Authentication Validation** | Real-Time | Microservice / Redis | `< 50ms` | Redis / Keycloak DB |
| **Fraud Scoring (Transactional)** | Real-Time | Kafka Streams / Flink | `< 500ms` | Redis / Neo4j |
| **SOC Critical Alerting** | Real-Time | Flink (CEP) | `< 1s` | ElasticSearch (SIEM) |
| **Identity Relationship Graphing** | Near Real-Time | Flink -> Neo4j | `< 5s` | Neo4j |
| **Analytics Aggregation** | Batch | Apache Spark | `Hourly / Daily` | Data Warehouse |
| **Compliance Reporting** | Batch | Apache Spark | `End of Month` | PDF / Data Warehouse |
| **Historical AI Retraining** | Batch | Kubeflow / Spark | `Weekly` | MinIO (Object Storage) |

---

## 3. Pipeline Isolation Strategy

To prevent a massive "End of Month" compliance report from stealing CPU cycles and delaying a live facial recognition check at an airport, resources are strictly segregated.

### 3.1. Compute Isolation (Kubernetes Node Pools)
*   **Real-Time Workloads:** Deployed to dedicated high-performance K8s Node Pools with `Guaranteed` QoS (Quality of Service) classes. CPU resources are explicitly reserved (`requests` == `limits`) to prevent CPU throttling.
*   **Batch Workloads:** Deployed to distinct Node Pools utilizing preemptible/spot instances to reduce costs. They run as Kubernetes `Jobs` or `CronJobs` with `Burstable` QoS, meaning they can consume excess CPU but will be throttled before impacting the host.

### 3.2. Kafka Streaming Isolation
*   **Consumer Group Segregation:** Hot path consumers and Cold path consumers belong to entirely different `group.id`s. If the Spark batch job falls 6 hours behind processing historical data, it has zero impact on the Flink job processing live data at the head of the partition.
*   **Network QoS:** Kafka brokers prioritize network I/O for Real-Time consumer groups during periods of extreme network congestion.

---

## 4. Storage Segregation Strategy

Operational databases (OLTP) must never be queried for analytics (OLAP).

*   **Primary DB Protection:** No human analyst or BI tool is allowed to run a `SELECT COUNT(*)` or complex `JOIN` against the primary PostgreSQL Identity database or the primary Neo4j instance.
*   **Read Replicas:** If near real-time queries are necessary, they are executed exclusively against asynchronous Read Replicas.
*   **The Data Lake (S3):** All analytical, AI retraining, and compliance workloads operate against the Data Lake. Kafka Connect continuously streams all events from Kafka into Parquet files in MinIO/S3. Spark jobs read these Parquet files, entirely isolating the heavy disk I/O from the operational databases.

---

## 5. Flink vs. Spark Boundaries

*   **Apache Flink** is used exclusively for unbounded streams requiring low-latency Complex Event Processing (CEP). For example, detecting if the *same* National ID attempts to log in from Paris and Port-au-Prince within a 5-minute window.
*   **Apache Spark** is used exclusively for bounded, high-volume batch processing. For example, recalculating the baseline demographic statistics of the entire national registry at midnight.
