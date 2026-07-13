# SNISID: National Digital Identity BPMN Workflow Architecture
## Enterprise-Grade Government Orchestration

This document details the exact Business Process Model and Notation (BPMN) workflows for the **Système National d’Identification et d’Interopérabilité Sécurisée des Identités et des Données (SNISID)**. These stateful, saga-driven workflows are designed to execute on an enterprise workflow orchestration engine (e.g., Temporal.io or Camunda Zeebe).

---

## 1. Citizen Enrollment Workflow
**Actors:** Citizen, Registration Agent, ONI Supervisor.
**Triggers:** Citizen physical appearance at an ONI center.
**SLA:** 15 mins (processing), 24h (supervisor approval). **Escalation:** Escalates to Regional Director if unapproved >48h.
**API Calls:** `POST /api/v1/enrollments`, `POST /api/v1/fraud/deduplicate`.
**Validation/Security:** ABAC check on Agent, MRZ scan of supporting documents.
**Exceptions/Rollback:** If fraud check fails, rollback enrollment record, trigger Fraud Investigation (Workflow 13).

```mermaid
sequenceDiagram
    participant Citizen
    participant Agent
    participant Engine as Workflow Engine
    participant DB as Identity Registry
    participant Fraud as Fraud Service

    Citizen->>Agent: Present Docs
    Agent->>Engine: Start Enrollment Saga
    Engine->>DB: Save Temporary Demographics
    Engine->>Fraud: Background Velocity Check
    Fraud-->>Engine: Status: Clear
    Engine->>Agent: Prompt Biometrics
    Agent->>Engine: Submit Biometrics (Triggers WF2)
```

## 2. Biometric Enrollment Workflow
**Actors:** Automated Biometric Identification System (ABIS), Registration Agent.
**Triggers:** Invoked as a sub-process of WF1.
**SLA:** 10s matching speed. **Escalation:** Manual AFIS reviewer intervention if match score is ambiguous.
**API Calls:** `POST /api/v1/biometrics/templates`.
**Audit:** Hash of raw template logged to WORM storage.
**Security:** Images encrypted at rest (AES-256) and transit (mTLS).

```mermaid
stateDiagram-v2
    [*] --> Capture
    Capture --> QualityCheck: ISO 19794 Compliance
    QualityCheck --> Extraction: Pass
    QualityCheck --> Capture: Fail (Retry max 3)
    Extraction --> Deduplication: 1:N Search
    Deduplication --> Hit: Duplicate Found
    Deduplication --> NoHit: Unique
    Hit --> ManualReview
    NoHit --> StoreTemplate
    StoreTemplate --> [*]
```

## 3. Identity Verification Workflow
**Actors:** Verifying Agency (e.g., DGI, Bank), Citizen.
**Triggers:** API request via X-Road or NFC scan of eID.
**API Calls:** `POST /api/v1/verify`.
**Exceptions/Retries:** Max 3 PIN attempts before card lockout.

```mermaid
sequenceDiagram
    participant Agency
    participant Gateway
    participant Auth as Auth Service
    
    Agency->>Gateway: Submit eID Cert + Cryptogram
    Gateway->>Auth: Validate Challenge Signature
    Auth->>Auth: Check OCSP/CRL
    Auth-->>Gateway: Result (Valid/Revoked)
    Gateway-->>Agency: Verification Token
```

## 4. Birth Registration Workflow
**Actors:** Hospital Admin, Civil Registry (National Archives).
**Triggers:** Birth at registered medical facility.
**Compliance:** Automatic generation of a unique NNI (Numéro National d'Identification) pre-enrollment ID.

```mermaid
sequenceDiagram
    participant Hospital
    participant Engine
    participant Archives
    Hospital->>Engine: Submit Birth Record
    Engine->>Archives: Create Civil Record
    Archives-->>Engine: Civil ID
    Engine->>Engine: Generate NNI (Pending Bio)
    Engine-->>Hospital: Print Birth Certificate
```

## 5. Death Registration Workflow
**Actors:** Medical Examiner, Civil Registry.
**Triggers:** Issuance of Death Certificate.
**Security Validation:** Maker-checker principle. Medical examiner drafts, Civil Registry official approves.

```mermaid
stateDiagram-v2
    [*] --> Drafted
    Drafted --> Approved: Civil Registry Check
    Approved --> UpdateDB: Set Status=Deceased
    UpdateDB --> RevokeCerts: Trigger PKI Revocation
    RevokeCerts --> NotifyAgencies: X-Road Event
    NotifyAgencies --> [*]
```

## 6. Marriage Registration Workflow
**Actors:** Civil Officer, Citizen A, Citizen B.
**Triggers:** Civil marriage ceremony.
**API Calls:** `PATCH /api/v1/citizens/marital-status`.

```mermaid
sequenceDiagram
    participant Officer
    participant Engine
    participant DB
    Officer->>Engine: Submit Marriage License
    Engine->>DB: Link NNI_A to NNI_B
    Engine->>DB: Update Marital Status
    Engine-->>Officer: Registry Confirmation
```

## 7. Address Changes Workflow
**Actors:** Citizen, Utility Provider (Oracle).
**Triggers:** Citizen web portal submission.
**Validation:** Cross-references utility bill APIs to verify physical existence of the address.

```mermaid
sequenceDiagram
    participant Citizen
    participant Portal
    participant Engine
    Citizen->>Portal: Upload Proof of Address
    Portal->>Engine: Address Update Saga
    Engine->>Portal: Request Utility Validation
    Engine->>Engine: Geocode Address
    Engine->>Engine: Update Record
```

## 8. Identity Correction Workflow
**Actors:** Citizen, High-Level Adjudicator.
**Triggers:** Found error in name/DOB.
**SLA:** 30 days. **Audit:** Mandatory preservation of historical records.

```mermaid
stateDiagram-v2
    [*] --> ClaimSubmitted
    ClaimSubmitted --> DocumentReview
    DocumentReview --> Rejected: Insufficient Proof
    DocumentReview --> CourtOrderRequired: Core Field Change
    DocumentReview --> Approved: Minor Typo
    Approved --> UpdateRecord
    UpdateRecord --> [*]
```

## 9. Identity Revocation Workflow
**Actors:** Judicial Court, ONI Director.
**Triggers:** Court order of severe identity fraud.
**Rollback:** Restores to previous NNI if merged incorrectly.

```mermaid
sequenceDiagram
    participant Court
    participant Director
    participant Engine
    participant PKI
    Court->>Director: Order Revocation
    Director->>Engine: Execute Revocation
    Engine->>PKI: Revoke all eID Certificates
    Engine->>Engine: Flag NNI as Invalid
```

## 10. Lost Credential Replacement Workflow
**Actors:** Citizen, Registration Agent.
**Triggers:** Citizen reports card lost.
**Security:** Instant certificate suspension to prevent misuse.

```mermaid
stateDiagram-v2
    [*] --> ReportedLost
    ReportedLost --> SuspendCert
    SuspendCert --> BioVerification: Citizen Appears in Person
    BioVerification --> Reissue: Match Confirmed
    Reissue --> IssueNewCard
    IssueNewCard --> [*]
```

## 11. Consent Management Workflow
**Actors:** Citizen, Third-Party Agency.
**Triggers:** Agency requests data, Citizen approves via Mobile Push.
**Compliance:** GDPR/Privacy-by-design. Consent expires after TTL.

```mermaid
sequenceDiagram
    participant Agency
    participant Engine
    participant App as Mobile App
    Agency->>Engine: Request Medical Access
    Engine->>App: Push Notification
    App-->>Engine: Citizen Approves (Cryptographic Signed)
    Engine->>Engine: Write Consent Grant to Ledger
    Engine-->>Agency: Provide Temporary Access Token
```

## 12. Inter-Agency Verification Workflow
**Actors:** Agency A, Agency B, X-Road Central Server.
**Triggers:** Agency A requires data from Agency B.
**Exception Handling:** If Agency B is down, Circuit Breaker trips, returns gracefully to A.

```mermaid
sequenceDiagram
    participant DGI
    participant XRoad
    participant ONI
    DGI->>XRoad: Request NNI Data
    XRoad->>XRoad: Validate Interoperability Policy
    XRoad->>ONI: Route Request via mTLS
    ONI-->>XRoad: Return Data
    XRoad-->>DGI: Data Payload
```

## 13. Fraud Investigation Workflow
**Actors:** Fraud Analyst, ABIS, SOC.
**Triggers:** Velocity check failure or Biometric 1:N duplicate hit.

```mermaid
stateDiagram-v2
    [*] --> SuspiciousActivity
    SuspiciousActivity --> FreezeAccount
    FreezeAccount --> AnalystReview
    AnalystReview --> Cleared: False Positive
    AnalystReview --> ConfirmedFraud: True Positive
    ConfirmedFraud --> LawEnforcementNotify
    Cleared --> UnfreezeAccount
    UnfreezeAccount --> [*]
```

## 14. Judicial Verification Workflow
**Actors:** Police (DCPJ), Judge.
**Triggers:** Warrant issued for data extraction.
**Security:** Requires M-of-N split key approval to export bulk records.

```mermaid
sequenceDiagram
    participant DCPJ
    participant Engine
    participant Judge
    DCPJ->>Engine: Submit Warrant Request
    Engine->>Judge: Request Cryptographic Approval
    Judge-->>Engine: Approve (Key 1 of 2)
    Engine->>Engine: Decrypt Vault
    Engine-->>DCPJ: Export Audit Trail
```

## 15. Immigration Verification Workflow
**Actors:** Border Control Agent.
**Triggers:** Citizen crosses border.
**API Calls:** `POST /api/v1/border/crossings`.

```mermaid
sequenceDiagram
    participant Agent
    participant ePassport
    participant Engine
    Agent->>ePassport: Scan NFC (BAC/EAC)
    ePassport-->>Agent: Signed Identity Data
    Agent->>Engine: Verify Document Signer CA
    Engine-->>Agent: Validated
```

## 16. Voter Verification Workflow
**Actors:** CEP (Electoral Council), Polling Station.
**Triggers:** Election day check-in.
**Offline:** Polling stations use pre-cached Bloom filters to verify offline.

```mermaid
stateDiagram-v2
    [*] --> CitizenArrives
    CitizenArrives --> OfflineBioCheck
    OfflineBioCheck --> Match: Authenticated
    Match --> MarkVotedLocal
    MarkVotedLocal --> SyncToCentral: Post-Election
    SyncToCentral --> [*]
```

## 17. Social Benefits Verification Workflow
**Actors:** Social Affairs, Citizen.
**Triggers:** Application for welfare.
**Validation:** Checks Death Registry and Identity Registry to prevent "Ghost" payments.

```mermaid
sequenceDiagram
    participant Social
    participant Engine
    participant DB
    Social->>Engine: Check Eligibility
    Engine->>DB: Verify NNI Alive Status
    Engine->>DB: Verify Income Bracket Consent
    Engine-->>Social: Eligibility Result
```

## 18. Offline Synchronization Workflow
**Actors:** Edge Node, Central Core.
**Triggers:** Internet connectivity restored at remote site.
**Retries:** Exponential backoff if central core is busy.

```mermaid
sequenceDiagram
    participant EdgeNode
    participant SyncService
    participant CoreKafka
    EdgeNode->>SyncService: Ping (Online)
    SyncService->>EdgeNode: Request Batch
    EdgeNode->>SyncService: Push Local JetStream Events
    SyncService->>CoreKafka: Publish to Main Topic
    SyncService-->>EdgeNode: ACK
    EdgeNode->>EdgeNode: Purge Local Cache
```

## 19. Disaster Recovery Operations Workflow
**Actors:** Infrastructure Admin.
**Triggers:** Region A goes offline entirely.
**SLA:** RTO < 15 mins.

```mermaid
stateDiagram-v2
    [*] --> RegionAFailure
    RegionAFailure --> GTM_Failover: Redirect DNS
    GTM_Failover --> PromoteRegionB: DB Promoted to Active
    PromoteRegionB --> TrafficRestored
    TrafficRestored --> [*]
```

## 20. Security Incident Workflows
**Actors:** SOC SOAR, Threat Intel.
**Triggers:** Falco runtime alert (e.g., unauthorized shell in Identity Service).
**Escalation:** Isolates pod instantly.

```mermaid
sequenceDiagram
    participant Falco
    participant SOAR
    participant K8s
    participant CISO
    Falco->>SOAR: Critical Alert (Container Breach)
    SOAR->>K8s: Apply NetworkPolicy (Isolate Pod)
    SOAR->>CISO: Send SMS Alert
    CISO->>SOAR: Authorize Pod Termination
    SOAR->>K8s: Delete Pod
```

---
*Enterprise Architect Note: All workflows are designed to emit OpenTelemetry spans. Exceptions at any node automatically trigger compensating transactions (Saga pattern) ensuring eventual consistency across the entire national infrastructure.*
