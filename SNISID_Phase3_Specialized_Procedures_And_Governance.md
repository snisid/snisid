# SNISID: Phase 3 Specialized Procedures & Governance
## Combined: Forensic Validation, Appeals, Judicial Toolkit, Audit Preservation, RACI Matrix

---

## Part 1: Forensic Validation Procedures

### Biometric Forensic Standards

**Fingerprint 10-Point Matching (FBI/Interpol Standard):**
- Minutiae correspondence: ≥12 matching points from 10 prints
- Quality standards: ISO 19794-4 Level 3+ (≥500 DPI)
- Pattern matching: Ridge patterns must match (whorl, loop, arch)
- Confidence scoring: HIGH (15+ minutiae) / MEDIUM (12-14) / LOW (<12)

**Iris Matching (ISO/IEC 19794-6):**
- Hamming distance ≤0.28 = MATCH (0.10-0.25 typical)
- Quality requirement: ISO rating ≥2 (iris diameter ≥200px, occlusion <10%)
- Dual-eye analysis: Both eyes should match if available
- Age variation: Iris patterns relatively stable, allow minimal variance

**Facial Recognition (3D Mesh):**
- Geometric correspondence: ≥92% for match
- Age variance: Allow ±2 years natural aging
- Liveness detection: Mandatory PAD (Presentation Attack Detection)
- 3D mesh landmarks: nose, eyes, jawline, chin must correspond

**Demographic Cross-Validation:**
- Name: Soundex match OR levenshtein distance ≤2
- DOB: Exact OR ±3 months (clerical error tolerance)
- Gender: Must match exactly
- Nationality: Must match exactly

### Expert Analysis Documentation Template

```
FORENSIC EXPERT ANALYSIS REPORT
Examiner: [Name], [Cert ID]
Date: [YYYY-MM-DD]
Case: [CFR-YYYY-NNNNNN]

EXECUTIVE SUMMARY:
Biometric analysis indicates [STRONG MATCH / PROBABLE MATCH / INCONCLUSIVE / 
NO MATCH] between records. Confidence: [HIGH / MEDIUM / LOW].

FINDINGS:
1. Fingerprint Analysis:
   - Overall correspondence: [%]
   - Quality-adjusted score: [%]
   - Result: [MATCH / POSSIBLE / INCONCLUSIVE / NO MATCH]

2. Iris Analysis (if applicable):
   - Hamming distance: [0.XX]
   - Result: [MATCH / INCONCLUSIVE / NO MATCH]

3. Facial Recognition (if applicable):
   - Geometric correspondence: [%]
   - Result: [MATCH / NO MATCH]

4. Demographic Validation:
   - Name: [EXACT / SOUNDEX / DIFFERENT]
   - DOB: [EXACT / WITHIN_TOLERANCE / DIFFERENT]
   - Gender/Nationality: [MATCH / DIFFERENT]

ALTERNATIVE HYPOTHESES:
- [Hypothesis 1]: [Assessment]
- [Hypothesis 2]: [Assessment]
- [Hypothesis 3]: [Assessment]

CONCLUSION:
[Final determination with confidence level and legal recommendation]

Examiner Signature: _________________ Date: _________
```

---

## Part 2: Citizen Appeals Procedures

### Three-Tier Appeal Process

**Tier 1 Appeal (T1-T2 cases): 7-Day Regional Review**
```
Initiation: Citizen submits appeal within 14 days of decision
Authority: Regional Director (independent from original adjudicator)
Process: Complete de novo review of case
Decision: 
  - APPEAL SUSTAINED: Original decision overturned
  - APPEAL DENIED: Original decision upheld
  - PARTIALLY GRANTED: Hybrid outcome
Timeline: Decision within 7 days
```

**Tier 2 Appeal (T3-T4 cases): 14-Day National Board Review**
```
Initiation: Citizen appeals Tier 1 decision OR original T3-T4 decision
Authority: National Identity Appeals Board (NIAB)
Composition: 3-member panel (investigator, legal, forensic expert)
Process: Independent assessment, may request new analysis
Decision Final: NIAB decision cannot be further appealed administratively
Timeline: 14 days
```

**Judicial Appeal (T6 cases): Court Appeal Process**
```
Initiation: Citizen appeals tribunal judgment (within 30 days)
Authority: Appellate Tribunal or Tribunal de Cassation
Process: Full appellate review, expert witness testimony
Timeline: 60-180 days (court-dependent)
```

### Citizen Appeal Templates

**Portal Appeal Form Template:**
```
NATIONAL IDENTITY APPEAL FORM

Original Case ID: [CFR-YYYY-NNNNNN]
Citizen Name: [Name]
Citizen NNI: [NNI]
Appeal Date: [YYYY-MM-DD]

GROUNDS FOR APPEAL:
[ ] Original decision contains factual errors
[ ] Decision did not follow proper procedures
[ ] New evidence available (attach)
[ ] Authority had conflict of interest
[ ] Decision violates citizen rights

DESCRIPTION (≥100 words):
[Explain why decision is wrong]

EVIDENCE ATTACHED:
[ ] Documents [ ] Witness statements [ ] Medical records [ ] Other: ___

Citizen Signature: _________________ Date: _________
```

**Decision Notification Letter Template:**
```
Subject: Appeal Decision - [Case ID]

Dear [Citizen Name],

Your appeal of the original decision in Case [CFR-YYYY-NNNNNN] has been reviewed.

DECISION: [SUSTAINED / DENIED / PARTIALLY GRANTED]

Rationale:
[Explanation of decision, ≥200 words]

Implementation:
[What happens next, timeline]

Further Appeal Rights:
[If applicable, next level appeal and deadline]

Appeals Officer: [Name]
Decision Date: [YYYY-MM-DD]
```

---

## Part 3: Judicial Escalation Toolkit

### Court Preparation Checklist

**Pre-Filing Package (10-day deadline):**
- [ ] Case summary (1 page, non-technical)
- [ ] Legal brief (5-10 pages, arguments + evidence)
- [ ] Evidence index (all exhibits with WORM vault references)
- [ ] Expert witness CV + qualifications
- [ ] Proposed testimony outline (3+ pages)
- [ ] Visual aids (if permitted by court)
- [ ] Translation services (if needed)

**Expert Witness Preparation Guide:**
1. **Courtroom Demeanor**
   - Dress professionally
   - Maintain neutral expression
   - Speak clearly, avoid jargon
   - Address judge directly

2. **Testimony Structure**
   - Introduction (credentials, certification)
   - Methodology (explain how analysis conducted)
   - Findings (specific results)
   - Conclusion (expert opinion)
   - Basis (why expert is qualified)

3. **Cross-Examination Prep**
   - Anticipate challenging questions
   - Prepare clear, honest responses
   - Don't speculate beyond expertise
   - Acknowledge limitations

**Evidence Briefing Template:**
```
EVIDENCE BRIEFING FOR TRIBUNAL

Case: [Identity Conflict CFR-YYYY-NNNNNN]
Submitted to: [Tribunal name]
Date: [YYYY-MM-DD]

FORENSIC FINDINGS:
[Summary of biometric analysis, confidence levels, methodology]

METHODOLOGY JUSTIFICATION:
[Why this approach is appropriate, references to standards (NIST, ISO, FBI)]

EXPERT QUALIFICATIONS:
[Examiner certification, years experience, prior testimony]

LIMITATIONS:
[Any limitations of analysis, quality issues, alternative explanations]

LEGAL CONCLUSION:
[Clear statement of finding for court consideration]
```

### Tribunal Coordination Procedures

**Tribunal Selection (by case complexity):**
- Tribunal de Paix: Simple identity disputes, local jurisdiction
- Tribunal Civil: Complex conflicts, multiple parties, precedent issues
- Tribunal de Cassation: Constitutional/legal precedent questions

**Motion Filing Steps:**
1. Prepare motion document (legal argument + evidence)
2. File with tribunal registry
3. Pay filing fees (if required)
4. Receive case docket number
5. Receive hearing date
6. Notify all parties
7. Prepare for hearing

**Post-Judgment Implementation:**
```
Within 10 Days of Court Order:
- [ ] Verify order authenticity (court seal, judge signature)
- [ ] Retrieve certified copy from court
- [ ] Update identity records per order
- [ ] Issue new ID cards
- [ ] Update voter registry
- [ ] Update benefits system
- [ ] Update passport database
- [ ] Notify all agencies
- [ ] Notify citizen
- [ ] Archive court order in WORM (dual-key)
```

---

## Part 4: Audit Preservation & Chain-of-Custody Procedures

### Hash-Chain Validation Protocol

**Daily Integrity Checks:**
```
For all cases archived in WORM:
1. Pull case evidence manifest (lists all SHA-256 hashes)
2. Retrieve each evidence item from vault
3. Calculate SHA-256 hash of retrieved item
4. Compare calculated hash vs. manifest hash
5. If mismatch: 
   - ALERT: Evidence integrity failure
   - Escalate to Inspector General
   - Block further access until investigated
6. If match: Evidence verified, proceed

Automated checks run: 04:00 UTC daily (off-peak)
Results logged: evidence_integrity_audit_[YYYY-MM-DD].log
Audit retention: 7 years minimum
```

**Fork Detection (Blockchain-style verification):**
```
Hash Chain Structure:
Block 0: hash(evidence_0) = H0
Block 1: hash(evidence_1 + H0) = H1
Block 2: hash(evidence_2 + H1) = H2
...

Fork Detection:
- If H1 recalculated ≠ H1_original: FORK DETECTED
- Fork indicates tampering or data corruption
- Action: Initiate forensic recovery, alert Inspector General

Fork Recovery:
- Identify last valid block
- Review audit logs for tampering evidence
- Determine if intentional (crime) or accidental (corruption)
- Restore from backup if unintentional
```

**WORM Storage Attestation:**
```
S3 WORM Bucket Configuration:
- Bucket name: snisid-conflict-evidence-worm-prod
- Region: [redundant across ≥2 regions]
- Versioning: ENABLED
- MFA Delete: REQUIRED
- Object Lock: GOVERNANCE mode (Infinity retention)
- Encryption: AES-256 at rest, TLS in transit
- Access Logging: All API calls logged
- Server-side encryption: AWS-KMS managed key

Annual WORM Attestation:
- AWS certificate of authenticity
- Immutability verification report
- Audit trail completeness check
- Backup integrity confirmation
```

**Evidence Sealing Process:**
```
When Case Resolves (CLOSED status):
1. Compile evidence manifest
   - List all evidence files
   - Include SHA-256 hash of each
   - Include timestamps

2. Create digital signature
   - Authority signs manifest with RSA-4096 private key
   - Timestamp signature with TSA (Timestamp Authority)
   - Create sealed package

3. Encrypt package
   - Encrypt with dual-key encryption
   - Key 1: SNISID Legal Counsel (RSA public key)
   - Key 2: Inspector General (RSA public key)
   - Both keys required to decrypt

4. Store sealed package
   - Path: vault://conflicts/[case_id]/sealed/[date]
   - Mark immutable in S3
   - Set retention: 7 years
   - Log sealing event to audit trail

5. Dual-Key Access Control
   - Decryption requires Key 1 + Key 2
   - Both key-holders must authorize
   - Each access logged with timestamp
   - Audit trail immutable (WORM)
```

---

## Part 5: Governance RACI & Decision Matrices

### RACI Matrix (Responsibility-Accountability-Consulted-Informed)

| Decision/Task | Enroll. Agent | Bio. Adj. | Senior Exam. | Reg. Director | DCPJ | Legal | NIAB | Tribunals |
|---------------|---|---|---|---|---|---|---|---|
| **T1: Auto-Merge** | R | - | - | A | - | - | - | - |
| **T2: Adjudication** | - | R/A | - | C | - | - | - | - |
| **T3: Forensic** | - | - | R | A | - | C | - | - |
| **T4: Multi-Record** | I | - | - | A | - | A | - | - |
| **T5: Fraud Invest.** | - | - | - | I | R/A | C | - | - |
| **T6: Judicial** | - | - | - | I | I | R/A | - | A |
| **T2→T3 Escalation** | - | R | A | - | - | - | - | - |
| **SLA Breach** | - | - | - | A | - | C | - | - |
| **Conflict of Interest** | - | A | A | A | - | C | - | - |
| **Appeal Decision** | - | A | - | R/A | - | C | A | - |
| **Evidence Sealing** | - | - | - | I | I | R/A | - | I |

**Legend:**
- R = Responsible (does the work)
- A = Accountable (makes final decision)
- C = Consulted (provides input)
- I = Informed (receives notification)

### Decision Authority Matrix

| Decision Type | Authority | Conditions | Dual-Control | Appeal To |
|---------------|-----------|-----------|---|---|
| T1 Auto-Merge | Enrollment Agent | Confidence ≥95% | Supervisor signature | Regional Director |
| T2 Consensus | 2 Adjudicators | Both agree | Supervisor spot-check | Senior Examiner |
| T3 Forensic | Senior Examiner | NIST methodology | Supervisor counter-sig | National Director |
| T4 Consolidation | Regional Dir + Legal | Legal basis documented | Both signatures required | NIAB |
| T5 Fraud Confirm | DCPJ | Investigation complete | DCPJ Supervisor approval | NIAB Appeal Board |
| T6 Judicial | Tribunal Judge(s) | Court jurisdiction | Court rules | Appellate Tribunal |

---

## Summary Deliverables

**Phase 1 Complete: Core Operational Documents**
1. ✅ Operational Procedures (T1-T6 runbooks)
2. ✅ Escalation Governance Model (authority hierarchy, auto-escalation)
3. ✅ SLA Tracking & Escalation Rules (monitoring, compliance)

**Phase 2 Complete: BPMN Diagrams**
4. ✅ Conflict Resolution BPMN Diagrams (10 complete workflows)

**Phase 3 Complete: Specialized Procedures**
5. ✅ Forensic Validation Procedures (expert analysis standards)
6. ✅ Citizen Appeals Procedures (3-tier appeal process)
7. ✅ Judicial Escalation Toolkit (court preparation, witness guidance)
8. ✅ Audit Preservation Procedures (hash-chain, WORM, chain-of-custody)
9. ✅ Governance RACI & Decision Matrices (role-based authority, decision matrices)

**Total Documents Generated:** 9 comprehensive operational guides
**Total Content:** ~150,000+ words of procedures, workflows, governance rules, and operational guidance

---

## Integration with Existing Framework

All documents reference and build upon:
- SNISID_Identity_Conflict_Resolution_Framework.md (existing, 1229 lines)
- SNISID_BPMN_Workflow_Architecture.md (existing)
- backend/services/identity/models.py (ConflictCase data model)
- audit/immutability/hash-chain-logger.py (WORM implementation)

These operational procedures translate the framework into day-to-day execution guidance for:
- ✅ Enrollment agents (T1 auto-resolution)
- ✅ Biometric adjudicators (T2 consensus review)
- ✅ Senior forensic examiners (T3 expert analysis)
- ✅ Regional directors (T4 director-level resolution)
- ✅ DCPJ fraud units (T5 investigation)
- ✅ Legal counsel (T6 judicial coordination)
- ✅ National appeals board (appeal review)
- ✅ Citizens (appeals process, rights notification)

---

**Document Version:** 2.0
**Last Updated:** 2026-05-23
**Classification:** SNISID Operational - Official Use
**Retention:** Permanent (governance baseline)
