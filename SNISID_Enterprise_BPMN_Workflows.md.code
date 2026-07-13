# SNISID Enterprise BPMN Workflows

This document contains a comprehensive suite of enterprise-grade Business Process Model and Notation (BPMN) workflows (represented via Mermaid flowcharts) governing the core operations, inter-agency interactions, and disaster resilience of the SNISID platform.

---

## 1. Unified Enrollment, Biometrics & Fraud Escalation Workflow

This workflow covers the core citizen onboarding process, integrating biometric ABIS deduplication and real-time fraud detection.

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white,shape:circle;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white,shape:circle;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px,shape:rhombus;
    classDef dataStore fill:#FFF3E0,stroke:#E65100,stroke-width:2px;
    classDef subProcess fill:#F3E5F5,stroke:#7B1FA2,stroke-width:2px,stroke-dasharray: 5 5;

    subgraph Agency Edge
        S1((Start)):::startEvent --> Cap[Capture Demographics & Fingerprints/Face]:::task
        Cap --> Submit[Submit Enrollment Payload]:::task
    end

    subgraph Core Identity Saga
        Submit --> Val1[Format & Syntax Validation]:::task
        Val1 --> ABIS[Route to ABIS 1:N Search]:::task
    end

    subgraph Biometric Subsystem
        ABIS --> CheckDup{Duplicate <br/> Found?}:::gateway
        CheckDup -- No --> BioOk[Biometric Validated]:::task
        CheckDup -- Yes --> FraudScore[Calculate Fraud Score]:::task
    end

    subgraph Fraud & Escalation
        FraudScore --> Assess{Score > <br/> Threshold?}:::gateway
        Assess -- Yes --> AutoReject[Auto-Reject & Flag Identity]:::task
        Assess -- No --> ManualRev[Queue for Human Investigation]:::task
        
        ManualRev --> ReviewDecision{Investigator <br/> Approved?}:::gateway
        ReviewDecision -- No --> AutoReject
        ReviewDecision -- Yes --> Override[Manual Override Approved]:::task
    end

    BioOk --> WriteDB
    Override --> WriteDB
    
    subgraph Finalization
        WriteDB[(Commit to Citizen Registry)]:::dataStore --> TriggerPKI[[Trigger PKI Issuance]]:::subProcess
        TriggerPKI --> Issue((End: Enrolled)):::endEvent
        AutoReject --> Reject((End: Rejected)):::endEvent
    end
```

---

## 2. Citizen Verification & Cryptographic Consent Workflow

This workflow dictates how third-party entities (e.g., Banks, Telecoms) verify citizen identities using explicit, mobile-based cryptographic consent.

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white,shape:circle;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white,shape:circle;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px,shape:rhombus;
    classDef dataStore fill:#FFF3E0,stroke:#E65100,stroke-width:2px;

    subgraph Third-Party Service
        S2((Start)):::startEvent --> Req[Request Citizen Profile]:::task
    end

    subgraph SNISID API Gateway
        Req --> GWCheck[Verify OAuth2 Client Credentials]:::task
        GWCheck --> ConsentDB[Check Active Consent Ledger]:::task
    end

    subgraph Consent Engine
        ConsentDB --> HasConsent{Valid Consent <br/> Exists?}:::gateway
        HasConsent -- Yes --> FetchData[Retrieve Citizen Data]:::task
        HasConsent -- No --> TriggerPush[Trigger Mobile Push Request]:::task
    end

    subgraph Citizen Mobile App
        TriggerPush --> PromptUser[Prompt Citizen for Approval]:::task
        PromptUser --> UserDecide{Citizen <br/> Approves?}:::gateway
        UserDecide -- No --> Deny[Return HTTP 403 Forbidden]:::task
        UserDecide -- Yes --> BioAuth[FIDO2 Biometric Auth]:::task
        BioAuth --> SignGrant[Cryptographically Sign Consent Grant]:::task
    end

    SignGrant --> StoreGrant[(Update Consent Ledger)]:::dataStore
    StoreGrant --> FetchData
    
    FetchData --> ReturnPayload[Return Minimally Required JSON]:::task
    ReturnPayload --> E2((End: Verified)):::endEvent
    Deny --> E3((End: Denied)):::endEvent
```

---

## 3. Inter-Agency Approval: Judicial & Tax Validation Workflow

This workflow represents highly privileged cross-agency data verification over the X-Road protocol, specifically for background checks spanning Justice (DCPJ) and Tax (DGI) authorities.

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white,shape:circle;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white,shape:circle;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px,shape:rhombus;

    subgraph Requesting Agency
        S3((Start)):::startEvent --> Init[Initiate Background Check]:::task
        Init --> XRoadOut[Sign SOAP/REST Envelope]:::task
    end

    subgraph Interoperability Gateway
        XRoadOut --> API[Route to SNISID X-Road Gateway]:::task
        API --> MTLSCheck[Verify Agency mTLS Certificate]:::task
    end

    subgraph DGI (Tax)
        MTLSCheck --> Fork1[Parallel Dispatch]:::task
        Fork1 --> TaxCheck[Query DGI NIF Status]:::task
        TaxCheck --> TaxStatus{Tax <br/> Compliant?}:::gateway
        TaxStatus -- Yes --> TaxOk[Tax Clearance = True]:::task
        TaxStatus -- No --> TaxFail[Tax Clearance = False]:::task
    end

    subgraph DCPJ (Judicial)
        Fork1 --> JudCheck[Query Criminal Records API]:::task
        JudCheck --> JudStatus{Warrant <br/> Active?}:::gateway
        JudStatus -- Yes --> JudFail[Judicial Clearance = False]:::task
        JudStatus -- No --> JudOk[Judicial Clearance = True]:::task
    end

    subgraph Aggregation
        TaxOk & TaxFail & JudFail & JudOk --> Aggregate[Aggregate Agency Responses]:::task
        Aggregate --> Audit[Write Immutable Audit Log]:::task
        Audit --> Return[Return Aggregate Response Payload]:::task
        Return --> E4((End)):::endEvent
    end
```

---

## 4. Disaster Recovery (PRA/PCA) Execution Workflow

This automated workflow triggers when a catastrophic event (e.g., Category 5 Hurricane, Earthquake) severs communication to the primary Port-au-Prince datacenter.

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white,shape:circle;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white,shape:circle;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px,shape:rhombus;
    classDef alert fill:#ffebee,stroke:#c62828,stroke-width:2px;

    subgraph Global Monitoring
        S5((Start: Ping Failure)):::startEvent --> Monitor[Thanos/Prometheus Detects DC1 Offline]:::task
        Monitor --> Wait1[Wait 60s for Transient Recovery]:::task
        Wait1 --> ConfirmOffline{DC1 Still <br/> Offline?}:::gateway
    end

    subgraph Executive Alerting
        ConfirmOffline -- Yes --> Page[Page National CISO & DR Team]:::alert
        ConfirmOffline -- No --> FalseAlarm((End: False Alarm)):::endEvent
    end

    subgraph Automated Failover Engine
        Page --> AutoFail{Auto-Failover <br/> Enabled?}:::gateway
        AutoFail -- No --> ManualExec[Wait for Manual Executive Trigger]:::task
        ManualExec --> GSLB
        
        AutoFail -- Yes --> GSLB[Update GSLB: Route 100% Traffic to DC2]:::task
        GSLB --> DB[Promote CockroachDB Follower to Leader in DC2]:::task
        DB --> Argo[Trigger ArgoCD Scale-Up in DC2]:::task
        Argo --> VerifyState[Verify APIs Responding 200 OK]:::task
    end

    subgraph SOC Notification
        VerifyState --> NotifyGov[Broadcast 'Operating in Region B' Alert]:::task
        NotifyGov --> E5((End: DR Active)):::endEvent
    end
```
