# SNISID Workflow Engine Architecture
## Distributed Sagas & BPMN Orchestration

This document details the **Workflow Engine Architecture** for SNISID. Because SNISID operates as a distributed microservices ecosystem, complex business processes (like Citizen Enrollment) span multiple independent databases. To guarantee data consistency across these services, SNISID abandons fragile distributed transactions (two-phase commit) in favor of the **Saga Pattern** orchestrated by a dedicated, stateful Workflow Engine (e.g., **Temporal.io** or **Camunda Zeebe**).

---

## 1. Core Workflow Capabilities

### Distributed Workflow Execution
The Workflow Engine acts as the central conductor. It coordinates calls to the Identity Service, Biometric Service, and Notification Service, maintaining the durable state of the overall transaction. If a pod crashes in the middle of a workflow, the engine seamlessly resumes execution from the exact line of code upon recovery.

### The Saga Pattern & Compensating Transactions
When executing a multi-step workflow across distributed databases, failure at any step requires a rollback.
- **Example:** If Step 1 (Save Biometrics) succeeds, but Step 2 (Create Identity Record) fails due to a DB timeout, the engine automatically executes a **Compensating Transaction**—it calls the Biometric Service to explicitly delete the orphaned template, guaranteeing eventual consistency.

### Human-in-the-Loop & SLA Enforcement
- **Human Approvals:** Workflows can be "paused" indefinitely while awaiting human input. For example, if the ABIS flags a citizen for suspected fraud, the workflow suspends and alerts a DCPJ investigator. Once the investigator clicks "Override" in the UI, the workflow resumes.
- **SLA Escalations:** The engine natively supports timers. If the human investigator does not resolve the fraud flag within 48 hours, an SLA timer fires, automatically escalating the ticket to a Senior Supervisor.

---

## 2. Kafka & Event-Driven Integration

The Workflow Engine integrates seamlessly with the Kafka backbone.
- **Event Listeners:** The engine can wait for external events to proceed. (e.g., The workflow pauses and waits until the `snisid.agency.tax_clearance.received` event is published to Kafka).
- **Event Publishers:** Upon successfully completing a Saga, the engine emits terminal events to Kafka (e.g., `snisid.workflow.enrollment.completed`), allowing downstream analytics and notification services to react.

---

## 3. Kubernetes Deployment Strategy

- **Stateless Workers:** The actual workflow code (e.g., Temporal Worker processes) runs in standard, stateless Kubernetes pods. They can be scaled infinitely via Horizontal Pod Autoscalers (HPA) based on queue depth.
- **Stateful Core:** The engine's core orchestration cluster is stateful and backed by highly available Cassandra or PostgreSQL clusters, ensuring the persistence of every workflow's execution history.

---

## 4. Architecture Diagrams (Mermaid)

### 1. Workflow Engine Deployment Topology
This diagram illustrates how the stateful engine orchestrates stateless workers and microservices.

```mermaid
graph TD
    classDef engine fill:#e1bee7,stroke:#6a1b9a,stroke-width:2px;
    classDef worker fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef svc fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef db fill:#fff3e0,stroke:#e65100,stroke-width:2px;

    API[API Gateway]:::svc

    subgraph Workflow_Engine_Core [Temporal / Camunda]
        Server[Orchestration Server]:::engine
        History[(History DB / Cassandra)]:::db
        Server <-->|Store Workflow State| History
    end

    subgraph Kubernetes_Workers
        W1[Enrollment Worker Pod]:::worker
        W2[Fraud Review Worker Pod]:::worker
    end

    subgraph Microservices
        ID[Identity Service]:::svc
        BIO[Biometric Service]:::svc
        NOT[Notification Service]:::svc
    end

    API -->|Start Enrollment Workflow| Server
    Server -->|Dispatch Task| W1
    
    W1 -->|Execute Activity 1| BIO
    W1 -->|Execute Activity 2| ID
    W1 -->|Execute Activity 3| NOT
```

### 2. Saga Pattern with Compensating Transactions
This sequence demonstrates a failure scenario during enrollment, where the engine automatically triggers the rollback mechanisms.

```mermaid
sequenceDiagram
    participant Engine as Workflow Engine
    participant BIO as Biometric Service
    participant ID as Identity Service
    participant NOT as Notification Service

    Note over Engine: Start Enrollment Saga

    Engine->>BIO: Activity 1: Store Biometric Template
    BIO-->>Engine: 200 OK (Template_ID: 999)
    Note over Engine: State Saved: Activity 1 Complete

    Engine->>ID: Activity 2: Create Identity Profile
    ID-->>Engine: 504 Gateway Timeout (Database Down)
    
    Note over Engine: Retry Policy Triggered...
    Engine->>ID: Retry Activity 2
    ID-->>Engine: 504 Gateway Timeout
    
    Note over Engine: Max Retries Exceeded. Initiating Rollback (Compensation)
    
    Engine->>BIO: Compensating Activity 1: Delete Template_ID 999
    BIO-->>Engine: 200 OK (Deleted)
    
    Note over Engine: Saga Failed Cleanly. Consistency Maintained.
    Engine-->>NOT: Send SMS: "Enrollment Failed. Please return to kiosk."
```

### 3. BPMN Human-in-the-Loop & SLA Escalation
This flowchart visualizes the process when the Biometric ABIS detects a duplicate template, triggering a manual human review with an SLA timer.

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white,shape:circle;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white,shape:circle;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px,shape:rhombus;
    classDef timer fill:#ffebee,stroke:#c62828,stroke-width:2px,shape:circle;

    S1((Start)):::startEvent --> Flag[ABIS Flags Suspected Duplicate]:::task
    
    Flag --> Suspend[Suspend Workflow & Assign to Level 1 Investigator]:::task
    
    Suspend --> Wait{Wait for <br/> Human Action}:::gateway
    
    Wait -- Approves / Rejects --> Decision{Investigator <br/> Decision}:::gateway
    
    %% SLA Enforcement
    Wait -- 48 Hour Timeout --> SLA((Timer Fired)):::timer
    SLA --> Escalate[Auto-Escalate to Level 2 Supervisor]:::task
    Escalate --> Decision
    
    Decision -- Flag is False Positive --> Resume[Resume Enrollment Flow]:::task
    Decision -- Flag is True Positive --> Reject[Reject & Revoke Identity]:::task
    
    Resume --> End1((Complete)):::startEvent
    Reject --> End2((Terminated)):::endEvent
```

---
*Prepared by the SNISID Cloud Infrastructure & Resilience Board.*
