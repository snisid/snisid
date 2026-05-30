# SNISID SIEM Architecture
## National Security Information and Event Management

This document details the **SIEM (Security Information and Event Management) Architecture** for SNISID. As the central nervous system of the National SOC, the SIEM is responsible for ingesting, correlating, and analyzing terabytes of telemetry across the entire sovereign ecosystem in real-time. It transforms raw logs into actionable intelligence to detect advanced persistent threats (APTs), insider threats, and system anomalies.

---

## 1. Centralized Telemetry Ingestion

The SIEM ingests telemetry from every layer of the SNISID stack.

### Ingestion Sources
1. **Kubernetes & eBPF (Falco/Cilium):** Real-time kernel syscalls, network flows, and container lifecycle events.
2. **Microservices (OpenTelemetry):** Distributed traces, RED metrics, and application logs (Identity, Consent, Biometrics).
3. **API Gateway & Mesh (Istio/Envoy):** L7 access logs, mTLS termination events, and rate-limiting triggers.
4. **IAM & PKI (Keycloak/HSM):** Authentication successes/failures, JWT issuance, OPA authorization denials, and cryptographic signing events.
5. **Immutable Audit Service:** Non-repudiable business events (e.g., "Agent X queried Citizen Y").

### OpenTelemetry (OTel) Pipeline
- **OTel Collectors:** Deployed as DaemonSets on every Kubernetes node. They scrape metrics and receive traces/logs via gRPC.
- **Data Normalization:** The OTel collectors parse raw logs, extract fields (TraceID, UserID, IP), and format them into the unified Elastic Common Schema (ECS) before forwarding.

---

## 2. SIEM Data Pipelines & Correlation Engine

### High-Throughput Buffering (Kafka)
To prevent the SIEM from buckling under a massive DDoS attack or a sudden spike in logs, all telemetry is first buffered in a dedicated **Kafka Security Cluster**. The SIEM consumption engine reads from this buffer at its own pace.

### Real-Time Correlation & Detection
The SIEM uses a stream processing engine to correlate events across different domains in real-time.
- **Rule-Based Detection (Sigma Rules):** Triggers instantly on known bad behavior (e.g., "5 failed logins followed by a successful login from a new IP").
- **Anomaly Detection (Machine Learning):** Establishes baselines for entity behavior (UEBA). For example, it detects if an API token suddenly queries 100x more records than its 30-day historical average.

---

## 3. Alerting & SOC Integration

- **Tiered Alerting:** High-fidelity alerts (P1/P2) are routed directly to the SOC Analysts via SOAR (Security Orchestration, Automation, and Response) platform integrations (e.g., TheHive or Cortex XSOAR).
- **Automated Mitigation:** For critical, unambiguous threats (e.g., Falco detects a shell spawned in a container), the SIEM triggers a SOAR playbook to instantly quarantine the node or revoke the IAM session without human intervention.

---

## 4. Immutable Retention & Compliance

- **Hot Tier (30 Days):** Data is stored in high-performance SSD-backed Elasticsearch indices for rapid threat hunting and dashboarding.
- **Cold Vault (10 Years):** To comply with national legal requirements, logs are aged out into an S3-compatible Ceph object store with strict **Write-Once-Read-Many (WORM)** policies. These logs are mathematically chained and cannot be altered, ensuring forensic non-repudiation in Haitian courts.

---

## 5. Architecture Diagrams (Mermaid)

### 1. SIEM Data Pipeline & Aggregation Topology
This diagram illustrates the flow of telemetry from edge devices to the core SIEM correlation engine.

```mermaid
graph TD
    classDef source fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef ingest fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef siem fill:#e1bee7,stroke:#6a1b9a,stroke-width:2px;
    classDef storage fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef action fill:#ffebee,stroke:#c62828,stroke-width:2px;

    subgraph Telemetry_Sources
        K8S[Falco / eBPF Kernel Events]:::source
        MESH[Istio Envoy Access Logs]:::source
        IAM[Keycloak Auth Logs]:::source
        PKI[HSM Cryptographic Logs]:::source
    end

    subgraph Ingestion_Layer
        OTEL[OpenTelemetry Collectors <br/> Normalize to ECS Format]:::ingest
        KAFKA[(Kafka Security Buffer)]:::ingest
    end

    subgraph SIEM_Core
        CORR[Real-Time Correlation Engine <br/> Sigma Rules + UEBA]:::siem
        HOT[(Elasticsearch <br/> 30-Day Hot Index)]:::siem
    end

    subgraph Long_Term_Storage
        COLD[(WORM Object Lock <br/> 10-Year Cold Vault)]:::storage
    end

    subgraph SOC_Operations
        SOAR[SOAR Platform <br/> Automated Playbooks]:::action
        DASH[SOC Dashboards <br/> Threat Hunting]:::action
    end

    K8S --> OTEL
    MESH --> OTEL
    IAM --> OTEL
    PKI --> OTEL

    OTEL -->|Stream Normalized Logs| KAFKA
    KAFKA -->|Consume & Analyze| CORR
    
    CORR -->|Index Data| HOT
    HOT -.->|Age Out Policy| COLD
    
    CORR -->|Trigger High-Fidelity Alert| SOAR
    HOT --> DASH
```

### 2. Threat Detection Workflow (Insider Threat Scenario)
This sequence diagram details how the SIEM detects and reacts to a compromised API credential.

```mermaid
sequenceDiagram
    participant API as API Gateway
    participant OTel as OpenTelemetry
    participant SIEM as SIEM Engine
    participant SOAR as SOAR Playbook
    participant IAM as Keycloak (IdP)

    API->>OTel: Log: "Token 123 queried Citizen A"
    API->>OTel: Log: "Token 123 queried Citizen B"
    Note over API, OTel: Attacker scripts 5,000 rapid queries.
    
    OTel->>SIEM: Stream 5,000 logs in 30 seconds
    
    SIEM->>SIEM: UEBA Engine analyzes velocity
    Note over SIEM: Baseline: 10 queries/min. Current: 10,000/min.
    
    SIEM->>SOAR: Fire P1 Alert: "Anomalous Data Exfiltration"
    
    SOAR->>IAM: API Call: Revoke Token 123 & Lock Account
    IAM-->>SOAR: 200 OK (Account Locked)
    
    SOAR->>SIEM: Log: "Automated Mitigation Applied"
```

---
*Prepared by the SNISID Cloud Infrastructure & Resilience Board.*
