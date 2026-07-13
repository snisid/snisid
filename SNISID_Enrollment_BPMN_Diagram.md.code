# SNISID BPMN Identity Enrollment Workflow

Below is the complete Mermaid flowchart simulating a standard Business Process Model and Notation (BPMN) diagram for the National Identity Enrollment Saga.

It covers offline edge synchronization, multi-system validation, biometric deduplication, and automated PKI issuance.

```mermaid
flowchart TD
    %% BPMN Styling
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white,shape:circle;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white,shape:circle;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px,shape:rhombus;
    classDef subprocess fill:#F3E5F5,stroke:#7B1FA2,stroke-width:2px,stroke-dasharray: 5 5;
    classDef dataStore fill:#FFF3E0,stroke:#E65100,stroke-width:2px;

    %% Citizen Swimlane
    subgraph Citizen [Citizen Lane]
        Start((Start)):::startEvent --> PresentDocs[Present ID Documents]:::task
        ReceiveID((Card Issued)):::endEvent
        ReceiveRejection((Rejected)):::endEvent
    end

    %% Remote Agency Edge Swimlane
    subgraph Agency [Remote Agency Edge]
        PresentDocs --> CaptureData[Capture Demographics & Biometrics]:::task
        CaptureData --> CheckConn{Internet <br/> Online?}:::gateway
        
        CheckConn -- No --> OfflineBuffer[(Local NATS Buffer)]:::dataStore
        OfflineBuffer -->|Offline Sync Loop| CheckConn
        
        CheckConn -- Yes --> SendPayload[Submit Enrollment Payload]:::task
        DeliverCard[Deliver eID Smart Card]:::task --> ReceiveID
    end

    %% SNISID Core Saga Swimlane
    subgraph SNISID_Core [SNISID Core Identity Saga]
        SendPayload --> LogStart[Write WORM Audit Log]:::task
        LogStart --> VerifyCivil[Agency Verification: Civil Registry]:::task
        VerifyCivil --> FraudCheck[Fraud Detection Velocity Check]:::task
        
        FraudCheck --> FraudGate{Fraud <br/> Detected?}:::gateway
        FraudGate -- Yes --> Exception[Exception: Flag Account]:::task
        
        FraudGate -- No --> TriggerBio[Route to ABIS]:::task
        
        WaitBio[Wait for Biometric Match]:::task
        ApproveChain[Supervisor Human Approval]:::task
        
        WaitBio --> ApproveChain
        
        ApproveChain --> ApprovalGate{Approved?}:::gateway
        ApprovalGate -- No --> RejectNotify[Send Rejection Notification]:::task
        RejectNotify --> ReceiveRejection
        
        ApprovalGate -- Yes --> WriteDB[(Write to Identity Registry)]:::dataStore
        WriteDB --> TriggerPKI[Request eID Certificates]:::task
    end

    %% Biometric ABIS Swimlane
    subgraph Biometric [ABIS Subsystem]
        TriggerBio --> BioProcess[[Biometric 1:N Deduplication]]:::subprocess
        BioProcess --> BioGate{Duplicate <br/> Found?}:::gateway
        BioGate -- Yes --> Exception
        BioGate -- No --> BioValid[Biometric Validated]:::task
        BioValid --> WaitBio
    end

    %% PKI & Notification Swimlane
    subgraph PKI_SOC [PKI & SOC Systems]
        Exception --> SOCAlert[Trigger SOC Analyst Alert]:::task
        
        TriggerPKI --> IssueCert[[Generate eID X.509 Certs]]:::subprocess
        IssueCert --> SMSNotify[Send SMS 'Card Ready']:::task
        SMSNotify --> DeliverCard
    end
```
