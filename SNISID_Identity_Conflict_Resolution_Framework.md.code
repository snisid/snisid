# SNISID: National Identity Conflict Resolution Framework
## Enterprise Conflict Adjudication, Forensic Validation & Escalation Governance

This document defines the **complete conflict resolution architecture** for the Système National d'Identification et d'Interopérabilité Sécurisée des Identités et des Données (SNISID). When Haiti's 15+ million citizens are enrolled under the principle of **"One Citizen, One Identity,"** conflicts *will* arise — duplicate biometric hits, contradictory civil records, contested identities, and deliberate fraud. This framework codifies every pathway from automated detection to judicial finality, ensuring no conflict is lost, no citizen is denied due process, and every decision is preserved in an immutable forensic chain.

---

## Table of Contents

1. [Conflict Taxonomy & Detection Architecture](#1-conflict-taxonomy--detection-architecture)
2. [Master Conflict Resolution BPMN Workflow](#2-master-conflict-resolution-bpmn-workflow)
3. [Duplicate Biometrics Resolution](#3-duplicate-biometrics-resolution)
4. [Conflicting Identity Records Resolution](#4-conflicting-identity-records-resolution)
5. [Fraud Investigation Workflow](#5-fraud-investigation-workflow)
6. [Manual Verification Procedures](#6-manual-verification-procedures)
7. [Judicial Escalation Workflow](#7-judicial-escalation-workflow)
8. [Citizen Appeals Process](#8-citizen-appeals-process)
9. [Forensic Validation Framework](#9-forensic-validation-framework)
10. [Audit Preservation Architecture](#10-audit-preservation-architecture)
11. [Escalation Governance Model](#11-escalation-governance-model)
12. [Operational Procedures & SLA Definitions](#12-operational-procedures--sla-definitions)
13. [Governance RACI Matrix](#13-governance-raci-matrix)

---

## 1. Conflict Taxonomy & Detection Architecture

### 1.1 Conflict Classification

Every identity conflict detected by SNISID is classified into one of six severity tiers, which directly determine the resolution pathway, required authority level, and SLA timeline.

| Tier | Classification | Example | Resolution Authority | SLA |
|------|---------------|---------|---------------------|-----|
| **T1** | Administrative Duplicate | Typo creates two records for same person | Enrollment Agent | 24h |
| **T2** | Biometric Near-Match | ABIS returns similarity score 85–94% | Biometric Adjudicator | 48h |
| **T3** | Biometric Hard-Match | ABIS returns similarity score ≥95% | Senior Forensic Examiner | 72h |
| **T4** | Multi-Record Conflict | Multiple NNIs with conflicting demographics for confirmed same person | Regional Director + Legal | 7 days |
| **T5** | Suspected Fraud | Deliberate duplicate enrollment, synthetic biometrics, or identity theft | DCPJ Fraud Unit | 30 days |
| **T6** | Judicial Conflict | Court-contested identity, inheritance disputes, contested nationality | Tribunal de Paix / Tribunal Civil | 90 days |

### 1.2 Detection Sources

```mermaid
graph TD
    classDef detect fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef engine fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef conflict fill:#ffebee,stroke:#c62828,stroke-width:2px;

    subgraph Detection_Sources
        D1[ABIS 1:N Deduplication<br/>During Enrollment]:::detect
        D2[Cross-Registry<br/>Reconciliation Batch]:::detect
        D3[Citizen Self-Report<br/>via Portal]:::detect
        D4[Inter-Agency Data<br/>Discrepancy Alert]:::detect
        D5[AI Fraud Engine<br/>Velocity / UEBA Anomaly]:::detect
        D6[Judicial Court Order<br/>Identity Challenge]:::detect
        D7[Death Registry<br/>Conflict with Active NNI]:::detect
    end

    subgraph Conflict_Engine
        CE[Conflict Classification<br/>Engine]:::engine
        CE --> T1[Tier 1: Administrative]:::conflict
        CE --> T2[Tier 2: Biometric Near]:::conflict
        CE --> T3[Tier 3: Biometric Hard]:::conflict
        CE --> T4[Tier 4: Multi-Record]:::conflict
        CE --> T5[Tier 5: Fraud]:::conflict
        CE --> T6[Tier 6: Judicial]:::conflict
    end

    D1 --> CE
    D2 --> CE
    D3 --> CE
    D4 --> CE
    D5 --> CE
    D6 --> CE
    D7 --> CE
```

### 1.3 Conflict Record Data Model

Every detected conflict generates an immutable `ConflictCase` record:

```yaml
ConflictCase:
  case_id: "CFR-2026-0001847"       # Unique sequential case ID
  detected_at: "2026-05-23T20:01:00Z"
  detection_source: "ABIS_DEDUP"     # Enum: ABIS_DEDUP | BATCH_RECON | CITIZEN_REPORT | INTER_AGENCY | AI_FRAUD | JUDICIAL | DEATH_REGISTRY
  tier: "T3"                         # Enum: T1-T6
  status: "OPEN"                     # Enum: OPEN | ASSIGNED | UNDER_REVIEW | ESCALATED | PENDING_JUDICIAL | RESOLVED | APPEALED | CLOSED
  primary_nni: "HT-2026-00391847"
  conflicting_nni: "HT-2019-00128374"
  abis_match_score: 97.3             # Nullable
  assigned_to: "AGT-00482"           # Nullable — assigned adjudicator
  escalation_chain: []               # Array of escalation records
  resolution:
    decision: null                   # MERGE | DEACTIVATE_PRIMARY | DEACTIVATE_CONFLICTING | REFER_JUDICIAL | FRAUD_CONFIRMED | CLEARED
    decided_by: null
    decided_at: null
    rationale: null
  audit_hash: "sha256:a4f8c2e1..."   # Cryptographic chain hash
  evidence_vault_ref: "vault://conflicts/CFR-2026-0001847"
```

---

## 2. Master Conflict Resolution BPMN Workflow

This is the **top-level orchestration workflow** governing all conflict resolution. It acts as the primary saga coordinator, dispatching to specialized sub-processes based on conflict tier.

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px;
    classDef subProc fill:#F3E5F5,stroke:#7B1FA2,stroke-width:2px,stroke-dasharray: 5 5;
    classDef danger fill:#ffebee,stroke:#c62828,stroke-width:2px;
    classDef audit fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;

    S((Conflict<br/>Detected)):::startEvent

    S --> LOG[Create ConflictCase<br/>Record in Ledger]:::audit
    LOG --> CLASSIFY[Classify Conflict Tier<br/>T1–T6]:::task
    CLASSIFY --> FREEZE[Freeze Affected NNIs<br/>Prevent Mutations]:::danger
    FREEZE --> EVIDENCE[Snapshot All Evidence<br/>to Forensic Vault]:::audit

    EVIDENCE --> TIER_GW{Conflict<br/>Tier?}:::gateway

    TIER_GW -- T1 --> AUTO_RESOLVE[[Auto-Resolve<br/>Administrative Duplicate]]:::subProc
    TIER_GW -- T2 --> BIO_NEAR[[Biometric Near-Match<br/>Adjudication]]:::subProc
    TIER_GW -- T3 --> BIO_HARD[[Biometric Hard-Match<br/>Forensic Exam]]:::subProc
    TIER_GW -- T4 --> MULTI_REC[[Multi-Record<br/>Reconciliation]]:::subProc
    TIER_GW -- T5 --> FRAUD[[Fraud Investigation<br/>DCPJ Referral]]:::subProc
    TIER_GW -- T6 --> JUDICIAL[[Judicial Escalation<br/>Court Process]]:::subProc

    AUTO_RESOLVE --> RESOLUTION
    BIO_NEAR --> RESOLUTION
    BIO_HARD --> RESOLUTION
    MULTI_REC --> RESOLUTION
    FRAUD --> RESOLUTION
    JUDICIAL --> RESOLUTION

    RESOLUTION{Resolution<br/>Accepted?}:::gateway
    RESOLUTION -- Yes --> EXECUTE[Execute Resolution<br/>Merge / Deactivate / Clear]:::task
    RESOLUTION -- No --> APPEAL[[Citizen Appeals<br/>Process]]:::subProc
    RESOLUTION -- Escalate --> ESCALATE[Escalate to<br/>Next Authority]:::danger

    APPEAL --> APPEAL_DEC{Appeal<br/>Outcome}:::gateway
    APPEAL_DEC -- Upheld --> EXECUTE
    APPEAL_DEC -- Overturned --> REOPEN[Reopen Case at<br/>Higher Tier]:::task
    REOPEN --> TIER_GW

    ESCALATE --> TIER_GW

    EXECUTE --> UNFREEZE[Unfreeze<br/>Surviving NNI]:::task
    UNFREEZE --> NOTIFY[Notify Affected<br/>Citizens & Agencies]:::task
    NOTIFY --> SEAL[Seal Case in<br/>WORM Audit Vault]:::audit
    SEAL --> E((Case<br/>Closed)):::endEvent
```

---

## 3. Duplicate Biometrics Resolution

### 3.1 ABIS Deduplication Conflict Workflow

When the ABIS 1:N search during enrollment returns a candidate match, the conflict enters this specialized sub-process.

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px;
    classDef forensic fill:#F3E5F5,stroke:#7B1FA2,stroke-width:2px;
    classDef danger fill:#ffebee,stroke:#c62828,stroke-width:2px;

    S((ABIS Returns<br/>Match Candidate)):::startEvent

    S --> SCORE_GW{Match<br/>Score?}:::gateway

    SCORE_GW -- "85–89%<br/>(Low Confidence)" --> LC[Low-Confidence Path]:::task
    SCORE_GW -- "90–94%<br/>(Medium Confidence)" --> MC[Medium-Confidence Path]:::task
    SCORE_GW -- "≥95%<br/>(High Confidence)" --> HC[High-Confidence Path]:::danger

    LC --> MODALITY_CHECK[Request Additional<br/>Modality Capture]:::task
    MODALITY_CHECK --> SECOND_SEARCH[ABIS Re-search with<br/>Iris + Face Fusion]:::forensic
    SECOND_SEARCH --> FUSED_SCORE{Fused Score<br/>≥90%?}:::gateway
    FUSED_SCORE -- No --> CLEARED[Clear: Unique Identity<br/>Proceed with Enrollment]:::task
    FUSED_SCORE -- Yes --> MC

    MC --> ADJUDICATOR[Assign to Certified<br/>Biometric Adjudicator]:::task
    ADJUDICATOR --> SIDE_BY_SIDE[Side-by-Side Minutiae<br/>Comparison Interface]:::forensic
    SIDE_BY_SIDE --> ADJ_DEC{Adjudicator<br/>Decision}:::gateway
    ADJ_DEC -- "Same Person" --> MERGE_CANDIDATES[Prepare Identity<br/>Merge Proposal]:::task
    ADJ_DEC -- "Different Persons" --> CLEARED
    ADJ_DEC -- "Inconclusive" --> HC

    HC --> FREEZE_BOTH[Freeze Both NNIs<br/>Immediately]:::danger
    FREEZE_BOTH --> SENIOR_EXAM[Assign Senior<br/>Forensic Examiner]:::forensic
    SENIOR_EXAM --> MULTI_MODAL[Full Multi-Modal<br/>Forensic Examination]:::forensic
    MULTI_MODAL --> PHYSICAL[Request In-Person<br/>Appearance of Both Citizens]:::task
    PHYSICAL --> LIVE_CAPTURE[Live Biometric Capture<br/>Under Observation]:::forensic
    LIVE_CAPTURE --> FORENSIC_DEC{Forensic<br/>Determination}:::gateway

    FORENSIC_DEC -- "Confirmed Same Person" --> MERGE_CANDIDATES
    FORENSIC_DEC -- "Different Persons<br/>ABIS False Positive" --> CLEAR_FP[Document False Positive<br/>Retrain Model]:::task
    FORENSIC_DEC -- "Suspected Fraud" --> FRAUD_REF[Refer to DCPJ<br/>Fraud Investigation]:::danger

    MERGE_CANDIDATES --> MERGE_APPROVAL{Regional Director<br/>Approves Merge?}:::gateway
    MERGE_APPROVAL -- Yes --> EXECUTE_MERGE[Execute NNI Merge<br/>Preserve Audit Trail]:::task
    MERGE_APPROVAL -- No --> ESCALATE_LEGAL[Escalate to<br/>Legal Review]:::danger

    CLEARED --> E1((Enrollment<br/>Proceeds)):::endEvent
    CLEAR_FP --> E1
    EXECUTE_MERGE --> E2((Conflict<br/>Resolved)):::endEvent
    FRAUD_REF --> E3((Fraud Case<br/>Opened)):::endEvent
    ESCALATE_LEGAL --> E4((Legal<br/>Review)):::endEvent
```

### 3.2 Biometric Adjudication Scoring Matrix

| Modality | Weight | Threshold (Same Person) | Threshold (Inconclusive) |
|----------|--------|------------------------|--------------------------|
| Fingerprint (10-print) | 40% | ≥95% individual, ≥12 minutiae match | 85–94% |
| Iris (dual) | 35% | Hamming distance ≤ 0.28 | 0.28–0.35 |
| Facial (3D mesh) | 20% | ≥92% similarity | 80–91% |
| Demographic cross-check | 5% | Name + DOB + Commune match | Partial match |

### 3.3 Multi-Modal Fusion Decision Logic

```mermaid
stateDiagram-v2
    [*] --> FingerCheck
    
    FingerCheck --> IrisCheck: Score Recorded
    IrisCheck --> FaceCheck: Score Recorded
    FaceCheck --> DemoCheck: Score Recorded
    DemoCheck --> FusionEngine: All Scores Collected
    
    FusionEngine --> ConfirmedSame: Weighted Score ≥ 92%
    FusionEngine --> ConfirmedDifferent: Weighted Score < 70%
    FusionEngine --> ManualReview: Weighted Score 70–91%
    
    ManualReview --> ForensicExam: Examiner Requests
    ManualReview --> ConfirmedSame: Examiner Confirms
    ManualReview --> ConfirmedDifferent: Examiner Rejects
    
    ConfirmedSame --> MergeProcess
    ConfirmedDifferent --> ClearEnrollment
    ForensicExam --> ConfirmedSame
    ForensicExam --> ConfirmedDifferent
    ForensicExam --> FraudReferral
    
    MergeProcess --> [*]
    ClearEnrollment --> [*]
    FraudReferral --> [*]
```

---

## 4. Conflicting Identity Records Resolution

### 4.1 Record Conflict Detection Sources

Conflicting identity records arise when two or more NNIs contain contradictory authoritative data for what appears to be the same person. Common scenarios:

- **Civil registry mismatch:** Birth certificate says "Jean-Pierre BAPTISTE" born 1985-03-15, but NNI record says "Jean Pierre BATISTE" born 1985-05-13.
- **Cross-agency discrepancy:** DGI tax records under NIF-A link to NNI-X, but Immigration passport records under the same biometrics link to NNI-Y.
- **Historical migration:** Pre-SNISID records from multiple paper-based registries merged incorrectly during the initial digitization campaign.

### 4.2 Record Conflict Resolution BPMN

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px;
    classDef evidence fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef danger fill:#ffebee,stroke:#c62828,stroke-width:2px;

    S((Record Conflict<br/>Detected)):::startEvent

    S --> GATHER[Gather All Records<br/>Across Registries]:::task
    GATHER --> TIMELINE[Construct Chronological<br/>Identity Timeline]:::evidence
    TIMELINE --> COMPARE[Automated Field-by-Field<br/>Discrepancy Analysis]:::task

    COMPARE --> SEVERITY{Discrepancy<br/>Severity}:::gateway

    SEVERITY -- "Minor: Typo / Accent<br/>Formatting difference" --> AUTO_FIX[Apply Automated<br/>Normalization Rules]:::task
    SEVERITY -- "Moderate: Name variant<br/>or DOB off by ≤30 days" --> AGENT_REVIEW[Assign to Data<br/>Quality Agent]:::task
    SEVERITY -- "Major: Core identity<br/>fields contradict" --> LEGAL_REVIEW[Assign to Legal<br/>Adjudicator]:::danger

    AUTO_FIX --> VERIFY_FIX{Citizen<br/>Confirms?}:::gateway
    VERIFY_FIX -- Yes --> APPLY[Apply Correction<br/>to Master Record]:::task
    VERIFY_FIX -- No --> AGENT_REVIEW

    AGENT_REVIEW --> DOCS[Request Original<br/>Source Documents]:::evidence
    DOCS --> CROSS_REF[Cross-Reference with<br/>Civil Registry Archives]:::task
    CROSS_REF --> AGENT_DEC{Agent<br/>Determination}:::gateway
    AGENT_DEC -- Resolved --> APPLY
    AGENT_DEC -- Cannot Determine --> LEGAL_REVIEW

    LEGAL_REVIEW --> COURT_DOCS[Require Sworn<br/>Affidavit from Citizen]:::evidence
    COURT_DOCS --> WITNESS[Gather Witness<br/>Statements if Needed]:::task
    WITNESS --> LEGAL_DEC{Legal<br/>Determination}:::gateway
    LEGAL_DEC -- "Administrative Fix" --> APPLY
    LEGAL_DEC -- "Court Order Required" --> JUDICIAL_REF[Refer to Tribunal<br/>de Paix]:::danger

    APPLY --> PROPAGATE[Propagate Correction<br/>to All Linked Agencies]:::task
    PROPAGATE --> AUDIT_SEAL[Seal Previous Record<br/>Version in Audit Vault]:::evidence
    AUDIT_SEAL --> E1((Conflict<br/>Resolved)):::endEvent
    JUDICIAL_REF --> E2((Judicial<br/>Process)):::endEvent
```

### 4.3 Source Document Authority Hierarchy

When conflicting documents exist, SNISID follows this strict authority precedence:

| Priority | Document Source | Authority Level | Notes |
|----------|---------------|-----------------|-------|
| 1 | Court Judgment (Jugement Supplétif) | Absolute | Overrides all other sources |
| 2 | Original Birth Certificate (Acte de Naissance) | Primary | Official civil registry extract |
| 3 | Baptismal Certificate | Secondary | Accepted where civil records destroyed |
| 4 | Hospital Birth Record | Secondary | Medical institution attestation |
| 5 | Passport / Travel Document | Tertiary | Immigration-issued identity |
| 6 | National ID Card (CIN) | Tertiary | Pre-SNISID legacy card |
| 7 | Electoral Card | Quaternary | CEP-issued voter registration |
| 8 | Sworn Affidavit (Affidavit Notarié) | Lowest | Requires corroborating evidence |

---

## 5. Fraud Investigation Workflow

### 5.1 Fraud Classification

| Code | Fraud Type | Detection Method | Penalty Severity |
|------|-----------|------------------|------------------|
| **F1** | Duplicate Enrollment | ABIS 1:N hard match + different demographic alias | Criminal |
| **F2** | Identity Theft | Citizen reports someone enrolled using their biometrics | Criminal |
| **F3** | Synthetic Biometrics | Liveness/PAD detection failure, deepfake artifacts | Criminal |
| **F4** | Ghost Worker Fraud | Salary disbursement to non-existent or deceased NNI | Criminal |
| **F5** | Agent Collusion | Insider deliberately bypasses ABIS checks | Criminal + Administrative |
| **F6** | Document Forgery | Altered birth certificates or court orders submitted | Criminal |
| **F7** | Identity Laundering | Serial identity changes to evade legal obligations | Criminal |

### 5.2 Fraud Investigation BPMN

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px;
    classDef danger fill:#ffebee,stroke:#c62828,stroke-width:2px;
    classDef forensic fill:#F3E5F5,stroke:#7B1FA2,stroke-width:2px;
    classDef audit fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;

    S((Fraud<br/>Suspected)):::startEvent

    S --> FREEZE[Immediate: Freeze<br/>All Suspect NNIs]:::danger
    FREEZE --> PRESERVE[Forensic Snapshot<br/>All Related Records]:::audit
    PRESERVE --> ASSIGN[Assign Lead<br/>Fraud Investigator]:::task

    ASSIGN --> INVESTIGATE[Conduct Investigation]:::forensic
    
    subgraph Investigation_Phase
        INVESTIGATE --> BIO_EXAM[Forensic Biometric<br/>Examination]:::forensic
        INVESTIGATE --> DOC_EXAM[Document Authenticity<br/>Verification]:::forensic
        INVESTIGATE --> AGENT_AUDIT[Agent Activity<br/>Audit Trail Review]:::audit
        INVESTIGATE --> GEO_CHECK[Geo-Temporal<br/>Correlation Analysis]:::task
        INVESTIGATE --> NETWORK[Social Network<br/>Analysis of NNIs]:::task
    end

    BIO_EXAM --> FINDINGS
    DOC_EXAM --> FINDINGS
    AGENT_AUDIT --> FINDINGS
    GEO_CHECK --> FINDINGS
    NETWORK --> FINDINGS

    FINDINGS[Compile Investigation<br/>Findings Report]:::task

    FINDINGS --> FRAUD_DEC{Fraud<br/>Confirmed?}:::gateway

    FRAUD_DEC -- "No: False Positive" --> CLEAR[Clear All Flags<br/>Unfreeze NNIs]:::task
    FRAUD_DEC -- "Yes: Administrative<br/>Fraud (F4–F5)" --> ADMIN_ACTION[Administrative Sanctions<br/>+ Record Correction]:::danger
    FRAUD_DEC -- "Yes: Criminal<br/>Fraud (F1–F3, F6–F7)" --> DCPJ_REF[Refer to DCPJ<br/>Criminal Division]:::danger

    ADMIN_ACTION --> DEACTIVATE[Deactivate<br/>Fraudulent NNI]:::danger
    DEACTIVATE --> REVOKE_CERTS[Revoke All<br/>PKI Certificates]:::danger
    REVOKE_CERTS --> NOTIFY_AGENCIES[Broadcast Revocation<br/>to All Agencies via X-Road]:::task

    DCPJ_REF --> WARRANT[DCPJ Obtains<br/>Judicial Warrant]:::task
    WARRANT --> PROSECUTE[Criminal<br/>Prosecution]:::danger
    PROSECUTE --> COURT_DEC{Court<br/>Verdict}:::gateway
    COURT_DEC -- Guilty --> DEACTIVATE
    COURT_DEC -- Not Guilty --> CLEAR

    CLEAR --> E1((Case<br/>Cleared)):::endEvent
    NOTIFY_AGENCIES --> E2((Fraud Case<br/>Closed)):::endEvent
    COURT_DEC --> E3((Judicial<br/>Finality)):::endEvent
```

### 5.3 Fraud Investigation Sequence — Agent Collusion Scenario

```mermaid
sequenceDiagram
    participant AI as AI Fraud Engine
    participant SOC as National SOC
    participant INV as Lead Investigator
    participant WORM as Audit Vault
    participant DCPJ as DCPJ Fraud Unit
    participant COURT as Tribunal Correctionnel

    AI->>SOC: Alert: Agent AGT-00482 bypassed<br/>ABIS check 47 times in 72h
    SOC->>SOC: SOAR Playbook: Suspend Agent Account
    SOC->>INV: Assign Case CFR-2026-0004821

    INV->>WORM: Extract Agent Activity Logs<br/>(Cryptographic Chain Verified)
    WORM-->>INV: 47 Enrollment Records<br/>All Missing Biometric Steps

    INV->>INV: Cross-reference: 47 NNIs<br/>all link to same commune,<br/>same day, same terminal

    INV->>INV: Forensic Biometric Review:<br/>12 of 47 NNIs share identical<br/>fingerprint templates

    INV->>DCPJ: Submit Case File with<br/>Cryptographic Evidence Chain
    
    DCPJ->>COURT: File Criminal Charges:<br/>Identity Fraud, Abuse of Authority
    
    Note over COURT: Trial Process
    
    COURT-->>DCPJ: Verdict: Guilty
    DCPJ->>SOC: Order: Deactivate 47 Fraudulent NNIs
    SOC->>SOC: Execute Mass Revocation
```

---

## 6. Manual Verification Procedures

### 6.1 In-Person Verification Protocol

When automated systems cannot resolve a conflict, SNISID triggers a mandatory in-person verification. This is a controlled, audited process.

**Operational Procedure: MV-001 — Citizen In-Person Verification**

| Step | Action | Actor | System | Evidence |
|------|--------|-------|--------|----------|
| 1 | Issue summons letter (bilingual FR/HT) with case reference | System | Notification Service | Delivery receipt logged |
| 2 | Citizen presents at designated ONI verification center | Citizen | — | Sign-in recorded |
| 3 | Verify original documents against scanned copies in vault | Verification Agent | Document Viewer | Side-by-side comparison screenshot |
| 4 | Capture fresh biometrics under direct camera observation | Biometric Technician | ABIS Workstation | Liveness verified, video recorded |
| 5 | Run fresh biometrics against both conflicting NNIs (1:1) | System | ABIS | Match scores logged |
| 6 | Record sworn verbal declaration from citizen | Verification Agent | Audio Recorder | Timestamped recording, consent noted |
| 7 | Agent submits determination with evidence attachments | Verification Agent | Case Management | Decision + rationale logged |
| 8 | Supervisor reviews and countersigns within 24h | Supervisor | Workflow Engine | Maker-checker pattern enforced |

### 6.2 Manual Verification BPMN

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px;
    classDef evidence fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;

    S((Manual Verification<br/>Triggered)):::startEvent

    S --> SUMMONS[Issue Bilingual<br/>Summons Letter]:::task
    SUMMONS --> WAIT_APPEAR{Citizen<br/>Appears?}:::gateway

    WAIT_APPEAR -- "No (30-day<br/>deadline)" --> SECOND_SUMMONS[Issue Second Summons<br/>via Registered Mail]:::task
    SECOND_SUMMONS --> WAIT2{Citizen<br/>Appears?}:::gateway
    WAIT2 -- No --> DEFAULT[Default Judgment:<br/>Deactivate Newer NNI]:::task
    WAIT2 -- Yes --> VERIFY_ID

    WAIT_APPEAR -- Yes --> VERIFY_ID[Verify Original<br/>Source Documents]:::evidence

    VERIFY_ID --> FRESH_BIO[Capture Fresh<br/>Biometrics Under Observation]:::task
    FRESH_BIO --> MATCH_TEST[Run 1:1 Against<br/>Both Conflicting NNIs]:::task

    MATCH_TEST --> MATCH_RESULT{Match<br/>Result}:::gateway
    MATCH_RESULT -- "Matches NNI-A only" --> RESOLVE_A[Confirm NNI-A<br/>Deactivate NNI-B]:::task
    MATCH_RESULT -- "Matches NNI-B only" --> RESOLVE_B[Confirm NNI-B<br/>Deactivate NNI-A]:::task
    MATCH_RESULT -- "Matches Both" --> MERGE_EVAL[Evaluate Merge<br/>Eligibility]:::task
    MATCH_RESULT -- "Matches Neither" --> ESCALATE[Escalate: Possible<br/>Third-Party Fraud]:::task

    RESOLVE_A --> SUPERVISOR[Supervisor<br/>Countersign]:::task
    RESOLVE_B --> SUPERVISOR
    MERGE_EVAL --> SUPERVISOR

    SUPERVISOR --> E1((Verification<br/>Complete)):::endEvent
    ESCALATE --> E2((Fraud<br/>Referral)):::endEvent
    DEFAULT --> E3((Default<br/>Decision)):::endEvent
```

---

## 7. Judicial Escalation Workflow

### 7.1 Escalation Triggers

Judicial escalation occurs when:
- A core identity field (name, nationality, parentage) is contested and cannot be resolved administratively
- A citizen disputes the administrative resolution through formal legal challenge
- A court order is required to merge, split, or revoke an identity
- Criminal fraud charges are necessary
- Inheritance or property disputes hinge on identity determination

### 7.2 Judicial Process BPMN

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px;
    classDef legal fill:#F3E5F5,stroke:#7B1FA2,stroke-width:2px;
    classDef audit fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef danger fill:#ffebee,stroke:#c62828,stroke-width:2px;

    S((Judicial<br/>Escalation)):::startEvent

    S --> PREPARE[Prepare Legal<br/>Case File]:::legal
    PREPARE --> EVIDENCE[Compile Digital<br/>Evidence Package]:::audit

    EVIDENCE --> SIGN_PKG[Digitally Sign Evidence<br/>with SNISID Institutional Key]:::audit
    SIGN_PKG --> FILE_COURT[File with Appropriate<br/>Court Jurisdiction]:::legal

    FILE_COURT --> JURISDICTION{Court<br/>Type}:::gateway

    JURISDICTION -- "Identity Correction<br/>(Minor)" --> TDP[Tribunal<br/>de Paix]:::legal
    JURISDICTION -- "Identity Determination<br/>(Major)" --> TC[Tribunal<br/>Civil]:::legal
    JURISDICTION -- "Criminal Fraud" --> TCORR[Tribunal<br/>Correctionnel]:::legal
    JURISDICTION -- "National Security" --> TSPEC[Tribunal<br/>Spécial]:::danger

    TDP --> HEARING
    TC --> HEARING
    TCORR --> HEARING
    TSPEC --> HEARING

    HEARING[Court Hearing<br/>SNISID Expert Testimony]:::legal

    HEARING --> COURT_ORDER{Court<br/>Order}:::gateway

    COURT_ORDER -- "Identity Confirmed" --> CONFIRM[Confirm Primary NNI<br/>Deactivate Duplicate]:::task
    COURT_ORDER -- "Identity Corrected" --> CORRECT[Apply Court-Ordered<br/>Corrections to Record]:::task
    COURT_ORDER -- "Identity Revoked" --> REVOKE[Full Identity<br/>Revocation + PKI]:::danger
    COURT_ORDER -- "New Identity Ordered<br/>(Jugement Supplétif)" --> CREATE[Create New NNI<br/>Per Court Specifications]:::task
    COURT_ORDER -- "Case Dismissed" --> DISMISS[Restore Original<br/>Status]:::task

    CONFIRM --> REGISTER_ORDER
    CORRECT --> REGISTER_ORDER
    REVOKE --> REGISTER_ORDER
    CREATE --> REGISTER_ORDER
    DISMISS --> REGISTER_ORDER

    REGISTER_ORDER[Register Court Order<br/>in SNISID Legal Ledger]:::audit
    REGISTER_ORDER --> PROPAGATE[Propagate Changes<br/>to All Agencies]:::task
    PROPAGATE --> NOTIFY[Notify Citizen<br/>of Final Determination]:::task
    NOTIFY --> SEAL[Seal Case in<br/>WORM Vault with Court Reference]:::audit
    SEAL --> E((Judicial<br/>Finality)):::endEvent
```

### 7.3 Court-to-SNISID Integration Sequence

```mermaid
sequenceDiagram
    participant CITIZEN as Citizen / Attorney
    participant COURT as Tribunal Civil
    participant SNISID_LEGAL as SNISID Legal Liaison
    participant ENGINE as Workflow Engine
    participant DB as Identity Registry
    participant PKI as PKI Service
    participant AUDIT as Audit Vault

    CITIZEN->>COURT: File Identity Correction Petition
    COURT->>SNISID_LEGAL: Request Digital Evidence Package

    SNISID_LEGAL->>ENGINE: Generate Court Evidence Export
    ENGINE->>AUDIT: Retrieve Conflict Case History<br/>(Cryptographic Chain Verified)
    AUDIT-->>ENGINE: Full Audit Trail + Hash Proof
    ENGINE-->>SNISID_LEGAL: Signed Evidence Package

    SNISID_LEGAL->>COURT: Submit Evidence + Expert Declaration

    Note over COURT: Hearing & Deliberation

    COURT->>SNISID_LEGAL: Issue Judgment:<br/>"Correct DOB to 1985-03-15,<br/>Merge NNI-X into NNI-Y"
    
    SNISID_LEGAL->>ENGINE: Execute Court Order CFR-2026-0001847
    ENGINE->>DB: Update DOB, Merge Records
    ENGINE->>PKI: Reissue Certificates with Corrected Data
    ENGINE->>AUDIT: Log Court Order Execution<br/>(Reference: Judgment No. 2026-TC-0847)
    ENGINE-->>SNISID_LEGAL: Execution Confirmed
    
    SNISID_LEGAL->>COURT: File Execution Report
    SNISID_LEGAL->>CITIZEN: Notify: Identity Updated
```

---

## 8. Citizen Appeals Process

### 8.1 Appeal Rights

Every citizen affected by a conflict resolution decision has the constitutional right to appeal. The appeals process is multi-tiered:

| Level | Appeal Body | Scope | Timeline | Filing Method |
|-------|------------|-------|----------|--------------|
| **L1** | ONI Regional Director | Administrative decisions (T1–T3) | 15 business days | Online portal or in-person |
| **L2** | National Identity Appeals Board (NIAB) | All administrative decisions, DCPJ referrals | 30 business days | Formal written petition |
| **L3** | Tribunal Administratif | Government agency disputes | 60 calendar days | Legal filing via attorney |
| **L4** | Cour de Cassation | Constitutional rights challenges | 90 calendar days | Supreme Court petition |

### 8.2 Citizen Appeals BPMN

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px;
    classDef citizen fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef legal fill:#F3E5F5,stroke:#7B1FA2,stroke-width:2px;

    S((Citizen Files<br/>Appeal)):::startEvent

    S --> VALIDATE[Validate Appeal<br/>Eligibility & Timeliness]:::task
    VALIDATE --> ELIGIBLE{Appeal<br/>Eligible?}:::gateway

    ELIGIBLE -- "No: Out of Time<br/>or Not Appealable" --> REJECT[Reject Appeal<br/>with Explanation]:::task
    ELIGIBLE -- Yes --> ACCEPT[Accept Appeal<br/>Stay Original Decision]:::citizen

    ACCEPT --> LEVEL{Appeal<br/>Level}:::gateway

    LEVEL -- L1 --> REGIONAL[ONI Regional<br/>Director Review]:::task
    LEVEL -- L2 --> NIAB[National Identity<br/>Appeals Board Hearing]:::legal
    LEVEL -- L3 --> TRIBUNAL[Tribunal<br/>Administratif Filing]:::legal
    LEVEL -- L4 --> CASSATION[Cour de Cassation<br/>Constitutional Review]:::legal

    REGIONAL --> GATHER_L1[Gather Additional<br/>Evidence from Citizen]:::citizen
    GATHER_L1 --> REVIEW_L1[Complete Case Review<br/>with Fresh Eyes]:::task
    REVIEW_L1 --> DEC_L1{L1<br/>Decision}:::gateway
    DEC_L1 -- "Uphold Original" --> CITIZEN_NEXT{Citizen Accepts<br/>or Escalates?}:::gateway
    DEC_L1 -- "Overturn: Correct Error" --> OVERTURN[Apply Corrected<br/>Decision]:::task

    CITIZEN_NEXT -- Accepts --> FINALIZE
    CITIZEN_NEXT -- Escalates --> NIAB

    NIAB --> PANEL[Convene 3-Member<br/>Review Panel]:::legal
    PANEL --> HEARING_L2[Formal Hearing:<br/>Citizen Presents Case]:::citizen
    HEARING_L2 --> DEC_L2{NIAB<br/>Decision}:::gateway
    DEC_L2 -- Uphold --> CITIZEN_NEXT2{Citizen Accepts<br/>or Escalates?}:::gateway
    DEC_L2 -- Overturn --> OVERTURN

    CITIZEN_NEXT2 -- Accepts --> FINALIZE
    CITIZEN_NEXT2 -- Escalates --> TRIBUNAL

    TRIBUNAL --> COURT_PROC[Full Judicial<br/>Proceedings]:::legal
    COURT_PROC --> DEC_L3{Court<br/>Decision}:::gateway
    DEC_L3 -- Final --> FINALIZE

    CASSATION --> CONST_REVIEW[Constitutional<br/>Rights Review]:::legal
    CONST_REVIEW --> DEC_L4{Supreme Court<br/>Decision}:::gateway
    DEC_L4 -- Final --> FINALIZE

    OVERTURN --> FINALIZE[Finalize Decision<br/>Update Records & Notify]:::task
    REJECT --> E1((Appeal<br/>Rejected)):::endEvent
    FINALIZE --> E2((Appeal<br/>Concluded)):::endEvent
```

### 8.3 Appeal Notification Template (Bilingual)

```
═══════════════════════════════════════════════════════════════
RÉPUBLIQUE D'HAÏTI — SNISID
NOTIFICATION DE RÉSOLUTION DE CONFLIT D'IDENTITÉ
NOTIFIKASYON REZOLISYON KONFLI IDANTITE
═══════════════════════════════════════════════════════════════

Réf: CFR-2026-0001847
NNI: HT-2026-00391847
Date: 2026-05-23

FR: Nous vous informons qu'une décision a été rendue 
    concernant le conflit d'identité enregistré sous la 
    référence ci-dessus. Vous disposez de 15 jours ouvrables 
    pour faire appel de cette décision.

HT: Nou enfòme w ke yo te pran yon desizyon konsènan 
    konfli idantite ki anrejistre anba referans ki pi wo a. 
    Ou gen 15 jou travay pou fè apèl kont desizyon sa a.

Décision / Desizyon: [MERGE | DEACTIVATE | CORRECTED]
Motif / Rezon: [Rationale in both languages]

Pour faire appel / Pou fè apèl:
  → En ligne / Sou entènèt: https://portal.snisid.gouv.ht/appeals
  → En personne / An pèsòn: Bureau ONI le plus proche
  → Par courrier / Pa lapòs: ONI, Rue [X], Port-au-Prince

═══════════════════════════════════════════════════════════════
```

---

## 9. Forensic Validation Framework

### 9.1 Forensic Examiner Certification Requirements

| Certification | Issuing Body | Validity | Required For |
|--------------|-------------|----------|-------------|
| Certified Latent Print Examiner (CLPE) | IAI (International Association for Identification) | 5 years | Fingerprint adjudication (T3+) |
| ABIS Operator Certification | SNISID Academy | 2 years | All biometric adjudication |
| ISO/IEC 19795 Biometric Testing | ISO | 3 years | ABIS threshold calibration |
| Forensic Document Examiner (FDE) | ABFDE | 5 years | Document forgery detection (F6) |
| Digital Forensics Certification | SNISID Academy | 2 years | Digital evidence handling |

### 9.2 Forensic Examination Workflow

```mermaid
flowchart TD
    classDef startEvent fill:#4CAF50,stroke:#388E3C,stroke-width:2px,color:white;
    classDef endEvent fill:#F44336,stroke:#D32F2F,stroke-width:2px,color:white;
    classDef task fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef gateway fill:#FFC107,stroke:#FFA000,stroke-width:2px;
    classDef forensic fill:#F3E5F5,stroke:#7B1FA2,stroke-width:2px;
    classDef chain fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;

    S((Forensic Exam<br/>Requested)):::startEvent

    S --> CHAIN_START[Initialize Forensic<br/>Chain of Custody]:::chain
    CHAIN_START --> RETRIEVE[Retrieve Evidence from<br/>Sealed Vault via M-of-N Keys]:::forensic

    RETRIEVE --> EXAM_TYPE{Examination<br/>Type}:::gateway

    EXAM_TYPE -- Biometric --> BIO_FORENSIC[Multi-Modal<br/>Biometric Analysis]:::forensic
    EXAM_TYPE -- Document --> DOC_FORENSIC[Physical & Digital<br/>Document Analysis]:::forensic
    EXAM_TYPE -- Digital --> DIG_FORENSIC[System Log &<br/>Metadata Forensics]:::forensic

    subgraph Biometric_Forensics
        BIO_FORENSIC --> MINUTIAE[Detailed Minutiae<br/>Point-by-Point Analysis]:::forensic
        BIO_FORENSIC --> IRIS_COMP[Iris Code<br/>Hamming Analysis]:::forensic
        BIO_FORENSIC --> FACE_3D[3D Facial<br/>Geometry Comparison]:::forensic
        BIO_FORENSIC --> LIVENESS[Historical Liveness<br/>Telemetry Review]:::forensic
    end

    subgraph Document_Forensics
        DOC_FORENSIC --> INK_PAPER[Ink Age &<br/>Paper Analysis]:::forensic
        DOC_FORENSIC --> SEAL_CHECK[Official Seal &<br/>Stamp Verification]:::forensic
        DOC_FORENSIC --> FONT_ANALYSIS[Typography &<br/>Printing Analysis]:::forensic
        DOC_FORENSIC --> METADATA_CHECK[Digital Metadata<br/>Examination]:::forensic
    end

    subgraph Digital_Forensics
        DIG_FORENSIC --> CHAIN_VERIFY[Verify Cryptographic<br/>Audit Chain Integrity]:::chain
        DIG_FORENSIC --> TIMELINE_RECON[Reconstruct Event<br/>Timeline from Logs]:::forensic
        DIG_FORENSIC --> IP_ANALYSIS[IP & Device<br/>Fingerprint Correlation]:::forensic
    end

    MINUTIAE & IRIS_COMP & FACE_3D & LIVENESS --> BIO_REPORT[Biometric<br/>Forensic Report]:::task
    INK_PAPER & SEAL_CHECK & FONT_ANALYSIS & METADATA_CHECK --> DOC_REPORT[Document<br/>Forensic Report]:::task
    CHAIN_VERIFY & TIMELINE_RECON & IP_ANALYSIS --> DIG_REPORT[Digital<br/>Forensic Report]:::task

    BIO_REPORT --> COMPILE[Compile Unified<br/>Forensic Report]:::task
    DOC_REPORT --> COMPILE
    DIG_REPORT --> COMPILE

    COMPILE --> PEER_REVIEW[Independent Peer<br/>Review by Second Examiner]:::forensic
    PEER_REVIEW --> AGREE{Examiners<br/>Agree?}:::gateway
    AGREE -- Yes --> SIGN_REPORT[Digitally Sign Report<br/>with Examiner PKI Credentials]:::chain
    AGREE -- No --> THIRD[Assign Third<br/>Examiner Tiebreaker]:::forensic
    THIRD --> SIGN_REPORT

    SIGN_REPORT --> SEAL[Seal Report in<br/>Forensic Evidence Vault]:::chain
    SEAL --> E((Forensic Report<br/>Complete)):::endEvent
```

### 9.3 Forensic Chain of Custody Record

```yaml
ForensicChainOfCustody:
  case_ref: "CFR-2026-0001847"
  evidence_items:
    - item_id: "EV-001"
      type: "BIOMETRIC_TEMPLATE"
      description: "10-print fingerprint template, NNI HT-2026-00391847"
      retrieved_from: "ABIS Gallery Vault"
      retrieved_by: "EXAM-00127 (Marie-Claire DESROSIERS)"
      retrieved_at: "2026-05-23T14:30:00Z"
      retrieval_witness: "EXAM-00089 (Jean FRANÇOIS)"
      integrity_hash: "sha256:7f3c8a2b..."
      chain:
        - action: "RETRIEVED"
          actor: "EXAM-00127"
          timestamp: "2026-05-23T14:30:00Z"
          location: "Forensic Lab, ONI Central"
        - action: "EXAMINED"
          actor: "EXAM-00127"
          timestamp: "2026-05-23T15:45:00Z"
          location: "Forensic Workstation WS-04"
        - action: "PEER_REVIEWED"
          actor: "EXAM-00089"
          timestamp: "2026-05-23T17:00:00Z"
          location: "Forensic Workstation WS-07"
        - action: "SEALED"
          actor: "EXAM-00127"
          timestamp: "2026-05-23T17:30:00Z"
          location: "Evidence Vault V-02"
          seal_hash: "sha256:b9d4e1f7..."
```

---

## 10. Audit Preservation Architecture

### 10.1 Conflict Resolution Audit Requirements

Every conflict resolution action generates an immutable audit record that satisfies three requirements:

1. **Legal Admissibility:** Records must be cryptographically signed and chain-linked to be admissible in Haitian courts
2. **Forensic Completeness:** Every decision, evidence item, and state transition is captured
3. **Tamper Evidence:** Any modification to historical records is mathematically detectable

### 10.2 Audit Event Taxonomy for Conflict Resolution

| Event Code | Event Description | Data Captured | Retention |
|-----------|-------------------|---------------|-----------|
| `CR.DETECTED` | Conflict initially detected | Source, tier, affected NNIs, ABIS score | 10 years |
| `CR.CLASSIFIED` | Conflict tier assigned | Tier level, classification rationale | 10 years |
| `CR.FROZEN` | NNI(s) frozen | Frozen NNIs, freeze timestamp, authority | 10 years |
| `CR.ASSIGNED` | Case assigned to adjudicator | Assignee ID, tier, SLA deadline | 10 years |
| `CR.EVIDENCE.ADDED` | Evidence item attached | Evidence hash, type, source | 10 years |
| `CR.ESCALATED` | Case escalated to higher tier | From-tier, to-tier, escalation reason | 10 years |
| `CR.DECISION` | Resolution decision made | Decision type, rationale, decided-by | **Permanent** |
| `CR.APPEAL.FILED` | Citizen files appeal | Appeal level, grounds, filing method | **Permanent** |
| `CR.APPEAL.DECIDED` | Appeal decided | Upheld/overturned, rationale | **Permanent** |
| `CR.EXECUTED` | Resolution executed in registry | Changes applied, before/after snapshot | **Permanent** |
| `CR.JUDICIAL.FILED` | Case filed with court | Court reference, jurisdiction | **Permanent** |
| `CR.JUDICIAL.ORDER` | Court order received | Order text, judge, court reference | **Permanent** |
| `CR.SEALED` | Case sealed in WORM vault | Final hash, seal timestamp | **Permanent** |

### 10.3 Audit Preservation Pipeline

```mermaid
sequenceDiagram
    participant WF as Workflow Engine
    participant AS as Audit Service
    participant HSM as PKI HSM
    participant HOT as OpenSearch (Hot)
    participant WORM as Ceph WORM Vault
    participant IG as Inspector General<br/>(Read-Only)

    WF->>AS: Emit CR.DECISION Event<br/>(Case CFR-2026-0001847)

    AS->>AS: Fetch Previous Event Hash<br/>(CR.ASSIGNED hash)
    AS->>AS: Concatenate: Payload + Previous Hash
    AS->>HSM: Request Institutional Digital Signature
    HSM-->>AS: Return Signature (Hash N)

    AS->>HOT: Index for Real-Time Query
    AS->>WORM: Write Immutable Record<br/>[Payload, PrevHash, Signature]
    
    Note over WORM: Object Lock: 10-Year<br/>Retention Enforced.<br/>Even root cannot delete.

    IG->>HOT: Periodic Audit Query:<br/>"Show all CR.DECISION events<br/>where tier ≥ T4 this month"
    HOT-->>IG: 47 Cases Returned
    
    IG->>WORM: Verify Chain Integrity:<br/>Recalculate all hashes
    WORM-->>IG: Chain Valid ✓
```

---

## 11. Escalation Governance Model

### 11.1 Escalation Authority Matrix

```mermaid
graph TD
    classDef l1 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef l2 fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef l3 fill:#fce4ec,stroke:#c62828,stroke-width:2px;
    classDef l4 fill:#f3e5f5,stroke:#6a1b9a,stroke-width:2px;
    classDef l5 fill:#ffebee,stroke:#b71c1c,stroke-width:3px;

    subgraph "Level 1 — Operational (T1–T2)"
        L1A[Enrollment Agent<br/>Auto-resolve typos]:::l1
        L1B[Biometric Adjudicator<br/>Near-match review]:::l1
    end

    subgraph "Level 2 — Supervisory (T3)"
        L2A[Senior Forensic Examiner<br/>Hard-match adjudication]:::l2
        L2B[ONI Center Supervisor<br/>Manual verification oversight]:::l2
    end

    subgraph "Level 3 — Regional (T4)"
        L3A[Regional ONI Director<br/>Multi-record reconciliation]:::l3
        L3B[SNISID Legal Counsel<br/>Legal determination]:::l3
    end

    subgraph "Level 4 — National (T5)"
        L4A[DCPJ Fraud Unit<br/>Criminal investigation]:::l4
        L4B[National Identity<br/>Appeals Board]:::l4
        L4C[National CISO<br/>Insider threat cases]:::l4
    end

    subgraph "Level 5 — Judicial (T6)"
        L5A[Tribunal de Paix<br/>Minor corrections]:::l5
        L5B[Tribunal Civil<br/>Major determinations]:::l5
        L5C[Tribunal Correctionnel<br/>Criminal fraud]:::l5
        L5D[Cour de Cassation<br/>Constitutional challenges]:::l5
    end

    L1A --> L1B
    L1B --> L2A
    L2A --> L2B
    L2B --> L3A
    L3A --> L3B
    L3B --> L4A
    L3B --> L4B
    L4A --> L5C
    L4B --> L5A
    L4B --> L5B
    L5A --> L5D
    L5B --> L5D
```

### 11.2 Escalation Trigger Rules

```mermaid
flowchart LR
    classDef rule fill:#E3F2FD,stroke:#1565C0,stroke-width:2px;
    classDef trigger fill:#FFC107,stroke:#FFA000,stroke-width:2px;

    subgraph Auto-Escalation Rules
        R1[SLA Breach:<br/>Case exceeds tier SLA<br/>by >50%]:::rule --> T1[Auto-escalate<br/>to next level]:::trigger
        R2[Adjudicator Conflict<br/>of Interest Detected]:::rule --> T2[Reassign to<br/>different region]:::trigger
        R3[3+ Cases with<br/>Same Agent Involvement]:::rule --> T3[Escalate to<br/>DCPJ for Pattern Analysis]:::trigger
        R4[Citizen Files<br/>Formal Complaint]:::rule --> T4[Escalate to<br/>Regional Director]:::trigger
        R5[Forensic Examiners<br/>Disagree After Tiebreaker]:::rule --> T5[Escalate to<br/>National Appeals Board]:::trigger
        R6[Case Involves<br/>Government Official's NNI]:::rule --> T6[Mandatory escalation<br/>to National CISO]:::trigger
    end
```

### 11.3 Conflict of Interest Controls

To prevent corruption in the adjudication process:

| Control | Implementation | Verification |
|---------|---------------|-------------|
| **Geographic Isolation** | Adjudicators cannot review cases from their home commune | System enforces via commune-of-origin check |
| **Relationship Screening** | System cross-references adjudicator's family NNIs against case participants | Automated pre-assignment check |
| **Random Assignment** | Cases assigned via weighted random algorithm, not manual selection | Algorithm audited quarterly |
| **Rotation Policy** | Adjudicators rotate regions every 6 months | HR system integration |
| **Dual-Control** | All T3+ decisions require independent countersignature | Workflow engine enforced |
| **Audit Trail** | Every case view, action, and decision is logged to WORM | Continuous monitoring by Inspector General |

---

## 12. Operational Procedures & SLA Definitions

### 12.1 SLA Escalation Timeline

```mermaid
gantt
    title Conflict Resolution SLA Timeline
    dateFormat  YYYY-MM-DD
    axisFormat  %d

    section Tier 1 (Administrative)
    Detection & Classification       :t1a, 2026-01-01, 2h
    Agent Resolution                  :t1b, after t1a, 20h
    Supervisor Approval               :t1c, after t1b, 2h

    section Tier 2 (Near-Match)
    Detection & Classification       :t2a, 2026-01-01, 2h
    Adjudicator Assignment           :t2b, after t2a, 4h
    Side-by-Side Review              :t2c, after t2b, 36h
    Decision & Countersign           :t2d, after t2c, 6h

    section Tier 3 (Hard-Match)
    Detection & Freeze               :t3a, 2026-01-01, 1h
    Forensic Examiner Assignment     :t3b, after t3a, 4h
    Multi-Modal Forensic Exam        :t3c, after t3b, 48h
    Peer Review                      :t3d, after t3c, 16h
    Regional Director Approval       :t3e, after t3d, 4h

    section Tier 4 (Multi-Record)
    Full Record Gathering            :t4a, 2026-01-01, 24h
    Legal Analysis                   :t4b, after t4a, 72h
    Resolution Proposal              :t4c, after t4b, 48h
    Director + Legal Approval        :t4d, after t4c, 24h

    section Tier 5 (Fraud)
    Immediate Freeze & DCPJ Notify   :t5a, 2026-01-01, 2h
    Full Investigation               :t5b, after t5a, 25d
    Prosecution Decision             :t5c, after t5b, 3d

    section Tier 6 (Judicial)
    Court Filing                     :t6a, 2026-01-01, 7d
    Court Proceedings                :t6b, after t6a, 60d
    Order Execution                  :t6c, after t6b, 7d
```

### 12.2 Key Performance Indicators (KPIs)

| KPI | Target | Measurement | Reporting |
|-----|--------|-------------|-----------|
| Mean Time to Detect (MTTD) | < 1 minute (automated) | Time from enrollment to conflict flag | Real-time dashboard |
| Mean Time to Assign (MTTA) | < 4 hours | Time from detection to adjudicator assignment | Daily report |
| Mean Time to Resolve (MTTR) | Within SLA per tier | Time from detection to case closure | Weekly report |
| SLA Compliance Rate | ≥ 95% | % of cases resolved within SLA | Monthly report |
| False Positive Rate (ABIS) | < 2% | % of biometric conflicts that are false positives | Quarterly report |
| Appeal Overturn Rate | < 10% | % of decisions overturned on appeal | Monthly report |
| Fraud Conviction Rate | ≥ 80% | % of referred fraud cases resulting in conviction | Annual report |
| Audit Chain Integrity | 100% | % of audit records passing chain verification | Continuous |

### 12.3 Operational Procedure Index

| Procedure Code | Title | Scope |
|---------------|-------|-------|
| **OP-CR-001** | Conflict Detection & Initial Classification | All tiers |
| **OP-CR-002** | NNI Freeze & Evidence Preservation | T2+ |
| **OP-CR-003** | Biometric Adjudication (Near-Match) | T2 |
| **OP-CR-004** | Forensic Biometric Examination | T3+ |
| **OP-CR-005** | Multi-Modal Fusion Analysis | T3+ |
| **OP-CR-006** | In-Person Manual Verification | T2+ |
| **OP-CR-007** | Record Conflict Reconciliation | T4 |
| **OP-CR-008** | Fraud Investigation Launch | T5 |
| **OP-CR-009** | DCPJ Criminal Referral | T5 |
| **OP-CR-010** | Judicial Case Filing | T6 |
| **OP-CR-011** | Court Order Execution | T6 |
| **OP-CR-012** | Citizen Appeal Processing | All tiers |
| **OP-CR-013** | NNI Merge Execution | T2–T4 |
| **OP-CR-014** | NNI Deactivation Execution | T3–T6 |
| **OP-CR-015** | Post-Resolution Agency Notification | All tiers |
| **OP-CR-016** | Forensic Chain of Custody Management | T3+ |
| **OP-CR-017** | Escalation Trigger Review | All tiers |
| **OP-CR-018** | Conflict of Interest Screening | All tiers |
| **OP-CR-019** | Case Sealing & WORM Archival | All tiers |
| **OP-CR-020** | Inspector General Quarterly Audit | Governance |

---

## 13. Governance RACI Matrix

### 13.1 Conflict Resolution RACI

| Activity | Enrollment Agent | Biometric Adjudicator | Forensic Examiner | Regional Director | DCPJ Fraud Unit | Legal Counsel | NIAB | Court | Inspector General | National CISO |
|----------|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| Detect Conflict | **R** | I | I | I | I | — | — | — | I | I |
| Classify Tier | **R** | C | C | I | — | — | — | — | I | — |
| Freeze NNI | A | **R** | — | I | — | — | — | — | I | I |
| T1 Resolution | **R/A** | — | — | I | — | — | — | — | I | — |
| T2 Adjudication | I | **R/A** | C | I | — | — | — | — | I | — |
| T3 Forensic Exam | — | C | **R** | **A** | I | — | — | — | I | — |
| T4 Reconciliation | — | — | C | **R/A** | — | **R** | I | — | I | — |
| T5 Fraud Investigation | — | — | C | I | **R/A** | C | — | — | I | C |
| T6 Judicial Process | — | — | C | I | C | **R** | — | **A** | I | — |
| Citizen Appeal (L1) | — | — | — | **R/A** | — | C | I | — | I | — |
| Citizen Appeal (L2) | — | — | — | C | — | C | **R/A** | — | I | — |
| Court Order Execution | — | — | — | I | I | **R** | — | **A** | I | I |
| Audit Preservation | I | I | I | I | I | I | I | I | **R/A** | C |
| Escalation Governance | — | — | — | C | C | C | C | — | **R** | **A** |

**Legend:** R = Responsible, A = Accountable, C = Consulted, I = Informed

### 13.2 Governance Review Cadence

| Review | Frequency | Participants | Output |
|--------|-----------|-------------|--------|
| Operational Case Review | Daily | Adjudicators, Supervisors | Active case status update |
| SLA Compliance Review | Weekly | Regional Directors, Ops Manager | SLA breach report + remediation |
| Fraud Pattern Analysis | Bi-weekly | DCPJ, AI Team, CISO | Emerging fraud vector assessment |
| Appeals Trend Analysis | Monthly | NIAB, Legal Counsel | Systemic error identification |
| ABIS Threshold Calibration | Quarterly | Forensic Examiners, ABIS Vendor | FAR/FRR optimization |
| Inspector General Audit | Quarterly | IG, National CISO | Independence & integrity report |
| National Governance Review | Annually | All stakeholders | Framework revision & policy update |

---

## Appendix A: Conflict Resolution State Machine

```mermaid
stateDiagram-v2
    [*] --> DETECTED: Conflict Identified

    DETECTED --> CLASSIFIED: Tier Assigned
    CLASSIFIED --> FROZEN: NNIs Frozen

    FROZEN --> ASSIGNED: Adjudicator Assigned
    ASSIGNED --> UNDER_REVIEW: Investigation Begins

    UNDER_REVIEW --> DECISION_PENDING: Evidence Complete
    UNDER_REVIEW --> ESCALATED: Tier Upgrade Required

    ESCALATED --> ASSIGNED: Reassigned at Higher Tier

    DECISION_PENDING --> RESOLVED: Decision Rendered
    DECISION_PENDING --> PENDING_JUDICIAL: Court Referral

    PENDING_JUDICIAL --> RESOLVED: Court Order Received

    RESOLVED --> APPEALED: Citizen Files Appeal
    RESOLVED --> EXECUTED: No Appeal / Appeal Window Closed

    APPEALED --> UNDER_REVIEW: Appeal Reopens Review
    APPEALED --> EXECUTED: Appeal Upheld

    EXECUTED --> SEALED: Changes Applied & Verified
    SEALED --> [*]: Case Permanently Closed
```

---

## Appendix B: API Contracts for Conflict Resolution

```yaml
paths:
  /v1/conflicts:
    post:
      summary: Report a new identity conflict
      security: [bearer_token, mtls]
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [detection_source, primary_nni]
              properties:
                detection_source:
                  type: string
                  enum: [ABIS_DEDUP, BATCH_RECON, CITIZEN_REPORT, INTER_AGENCY, AI_FRAUD, JUDICIAL, DEATH_REGISTRY]
                primary_nni: { type: string }
                conflicting_nni: { type: string }
                abis_match_score: { type: number, minimum: 0, maximum: 100 }
                evidence: 
                  type: array
                  items:
                    type: object
                    properties:
                      type: { type: string }
                      data: { type: string, format: base64 }
      responses:
        '201':
          description: Conflict case created
          content:
            application/json:
              schema:
                properties:
                  case_id: { type: string, example: "CFR-2026-0001847" }
                  tier: { type: string, example: "T3" }
                  status: { type: string, example: "OPEN" }
                  sla_deadline: { type: string, format: date-time }

  /v1/conflicts/{case_id}/resolution:
    post:
      summary: Submit a resolution decision for a conflict case
      security: [bearer_token, mtls, rbac_adjudicator]
      parameters:
        - name: case_id
          in: path
          required: true
          schema: { type: string }
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [decision, rationale]
              properties:
                decision:
                  type: string
                  enum: [MERGE, DEACTIVATE_PRIMARY, DEACTIVATE_CONFLICTING, REFER_JUDICIAL, FRAUD_CONFIRMED, CLEARED]
                rationale: { type: string, minLength: 50 }
                evidence_refs:
                  type: array
                  items: { type: string }
                countersigned_by: { type: string }
      responses:
        '200':
          description: Resolution accepted, pending execution
        '409':
          description: Conflict case already resolved or under appeal

  /v1/conflicts/{case_id}/appeal:
    post:
      summary: File a citizen appeal against a conflict resolution
      security: [citizen_auth, fido2]
      parameters:
        - name: case_id
          in: path
          required: true
          schema: { type: string }
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [appeal_level, grounds]
              properties:
                appeal_level:
                  type: string
                  enum: [L1_REGIONAL, L2_NIAB, L3_TRIBUNAL, L4_CASSATION]
                grounds: { type: string, minLength: 100 }
                supporting_documents:
                  type: array
                  items:
                    type: object
                    properties:
                      filename: { type: string }
                      data: { type: string, format: base64 }
      responses:
        '201':
          description: Appeal filed, original decision stayed
        '400':
          description: Appeal window expired or invalid level
```

---

*Ratified by the National Digital Identity & Interoperability Steering Committee (Comité National d'Identité Numérique et d'Interopérabilité).*

*Document Classification: SNISID-GOV-CR-001 | Version 1.0 | Date: 2026-05-23*

*Prepared by the SNISID Enterprise Architecture & Legal Governance Division.*
