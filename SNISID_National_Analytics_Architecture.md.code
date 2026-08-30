# SNISID National Analytics Architecture
## Sovereign Data Lake & Predictive Intelligence

This document details the architectural design for the **National Analytics Service**. While the core SNISID microservices are heavily optimized for extreme transactional throughput (OLTP), the Haitian government also requires deep, long-term analytical capabilities (OLAP) to drive policy decisions, monitor national demographics, and retrain fraud-detection AI models.

Crucially, this analytics architecture strictly enforces privacy; all citizen Personally Identifiable Information (PII) is cryptographically tokenized or anonymized before entering the Data Lake.

---

## 1. The Sovereign Data Lakehouse

SNISID utilizes an open-source Data Lakehouse architecture (combining the flexibility of a Data Lake with the ACID guarantees of a Data Warehouse).

### Storage & Compute Separation
- **Storage Layer:** All analytical data is stored in **Apache Iceberg** table formats on S3-compatible, on-premise Ceph Object Storage.
- **Compute Layer:** **Trino** (formerly PrestoSQL) and **Apache Spark** are used to run massively parallel, distributed SQL queries across the petabytes of stored data.

### The ELT Pipeline (Extract, Load, Transform)
1. **Extract/Load:** A Kafka Connect sink continuously streams raw, immutable events (e.g., `CitizenRegistered`, `FraudFlagged`) from the Kafka Event Bus directly into the "Bronze" (Raw) zone of the Data Lake.
2. **Transform (Anonymization):** Spark jobs run periodically to clean the data, stripping out PII (like exact names or raw biometric templates), mapping exact birthdates to birth *years*, and moving the cleaned data to the "Silver" and "Gold" zones for querying.

---

## 2. Analytics Domains & Integration

### 1. Citizen & Demographic Analytics (Business Intelligence)
- **Use Case:** The CEP (Electoral Council) needs to predict how many citizens will turn 18 before the next election cycle to allocate voting booths.
- **Integration:** Data is visualized using **Apache Superset** or **Metabase**, providing rich, interactive dashboards for authorized government ministers and analysts.

### 2. Fraud & Predictive AI Analytics
- **Use Case:** The ML models in the Fraud Detection Service must be periodically retrained on historical data to adapt to new attack vectors.
- **Integration:** Data scientists use Jupyter Notebooks connected directly to the Spark clusters to analyze millions of past fraud flags and train new TensorFlow models.

### 3. Operational Analytics & Observability
- **Use Case:** The SRE team needs to monitor the real-time latency of the Identity Service and track API error rates.
- **Integration:** This bypasses the Data Lake entirely. **OpenTelemetry** traces, metrics, and logs are scraped by Prometheus and visualized in real-time **Grafana** dashboards.

---

## 3. Architecture Diagrams (Mermaid)

### 1. Global Data Pipeline & Analytics Topology
This diagram illustrates how data flows from the transactional microservices, through the anonymization pipelines, and into the BI dashboards.

```mermaid
graph TD
    classDef oltp fill:#ffebee,stroke:#c62828,stroke-width:2px;
    classDef broker fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef lake fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef compute fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef bi fill:#e1bee7,stroke:#6a1b9a,stroke-width:2px;

    subgraph Transactional_Systems_OLTP
        ID[Identity Service]:::oltp
        F[Fraud Service]:::oltp
    end

    KAFKA[(Kafka Event Bus)]:::broker

    subgraph Data_Lakehouse_Storage
        BRONZE[(Bronze Zone <br/> Raw Kafka Events)]:::lake
        SILVER[(Silver Zone <br/> Cleaned & Anonymized)]:::lake
        GOLD[(Gold Zone <br/> Aggregated Views)]:::lake
    end

    subgraph Distributed_Compute
        SPARK[Apache Spark <br/> ETL / Anonymization]:::compute
        TRINO[Trino SQL Engine]:::compute
    end

    subgraph Consumption_Layer
        BI[Apache Superset <br/> Demographics Dashboards]:::bi
        JUP[Jupyter Notebooks <br/> AI Model Training]:::bi
    end

    ID -->|Publish Events| KAFKA
    F -->|Publish Events| KAFKA
    
    KAFKA -->|Kafka Connect Sink| BRONZE
    
    BRONZE --> SPARK
    SPARK -->|Strip PII & Transform| SILVER
    SPARK -->|Create Materialized Views| GOLD
    
    SILVER --> TRINO
    GOLD --> TRINO
    
    TRINO --> BI
    TRINO --> JUP
```

### 2. Operational Observability & KPI Monitoring Flow
This flowchart isolates the operational (SRE) analytics, distinct from the demographic business intelligence.

```mermaid
graph LR
    classDef app fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef otel fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef storage fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef viz fill:#e1bee7,stroke:#6a1b9a,stroke-width:2px;

    App1[Citizen Registry Pod]:::app
    App2[API Gateway Pod]:::app
    
    OTEL[OpenTelemetry Collector]:::otel
    
    PROM[(Prometheus <br/> Metrics TSDB)]:::storage
    TEMPO[(Tempo <br/> Distributed Tracing)]:::storage
    LOKI[(Loki <br/> Application Logs)]:::storage
    
    GRAF[Grafana Dashboards <br/> Real-Time SRE Alerts]:::viz

    App1 -->|Export RED Metrics, Logs, Traces| OTEL
    App2 -->|Export RED Metrics, Logs, Traces| OTEL
    
    OTEL --> PROM
    OTEL --> TEMPO
    OTEL --> LOKI
    
    PROM --> GRAF
    TEMPO --> GRAF
    LOKI --> GRAF
```

---
*Prepared by the SNISID Cloud Infrastructure & Resilience Board.*
