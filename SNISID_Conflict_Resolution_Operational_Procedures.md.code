# SNISID: Conflict Resolution Operational Procedures
## Tier-Specific Runbooks, Decision Matrices & Checklists

---

## Table of Contents
1. [Overview & Common Procedures](#overview--common-procedures)
2. [Tier 1: Administrative Duplicate Resolution (24h SLA)](#tier-1-administrative-duplicate-resolution)
3. [Tier 2: Biometric Near-Match Adjudication (48h SLA)](#tier-2-biometric-near-match-adjudication)
4. [Tier 3: Biometric Hard-Match Forensic Examination (72h SLA)](#tier-3-biometric-hard-match-forensic-examination)
5. [Tier 4: Multi-Record Conflict Resolution (7-Day SLA)](#tier-4-multi-record-conflict-resolution)
6. [Tier 5: Fraud Investigation (30-Day SLA)](#tier-5-fraud-investigation)
7. [Tier 6: Judicial Escalation (90-Day SLA)](#tier-6-judicial-escalation)
8. [Common Operational Tasks](#common-operational-tasks)
9. [Quality Assurance Checklists](#quality-assurance-checklists)

---

## Overview & Common Procedures

### 1.1 Case Lifecycle States & Transitions

All conflict cases flow through the following immutable state machine:

```
OPEN → ASSIGNED → UNDER_REVIEW → ESCALATED → RESOLVED/APPEAL/REFERRED_JUDICIAL → CLOSED
                                ↓
                          PENDING_JUDICIAL → CLOSED
```

**State Definitions:**
- **OPEN** (0-4h): Case auto-created by detection engine, awaiting assignment
- **ASSIGNED** (4h+): Assigned to appropriate adjudicator/examiner
- **UNDER_REVIEW** (6h+): Active investigation/examination in progress
- **ESCALATED** (varies): Case escalated to higher authority (SLA breach, conflict of interest, pattern detection)
- **PENDING_JUDICIAL** (30d+): Case referred to court, awaiting tribunal decision
- **RESOLVED** (varies): Decision made, implementation in progress
- **APPEALED** (varies): Citizen filed appeal, case re-opened at Tier N+1
- **CLOSED** (final): All appeals exhausted, final decision implemented, evidence sealed

### 1.2 Universal Case Management Procedures

**Every case must follow:**
1. **Evidence Intake**: All evidence immediately logged to WORM storage with SHA-256 hash
2. **Chain-of-Custody**: Every access recorded (who, when, duration, purpose)
3. **Conflict of Interest Check**: Automated system + manual verification
4. **SLA Tracking**: Real-time dashboard; alerts at 50%, 75%, 100% elapsed
5. **Escalation Rules**: Auto-escalate if SLA breach >50% or conflict detected
6. **Documentation**: All decisions timestamped with rationale (≥100 words for T3+)
7. **Notification**: Citizen notified within 24h of resolution (T1-T4), 48h (T5), per court order (T6)
8. **Appeal Rights**: Explicit notification of appeal rights with appeal deadline and procedures

**Decision Authority Rule:**
- Every decision requires approval by the assigned authority
- T3+ decisions require dual-control: primary investigator + supervisor counter-signature
- Overrides of auto-decisions require T3+ approval + written justification
- No authority can override tribunal decisions without appellate court order

### 1.3 Common Required Documentation

**For all tiers, case file must contain:**
- **Case Summary** (1-2 paragraphs: parties, conflict nature, evidence summary)
- **Decision Rationale** (≥50 words: why this decision was made)
- **Evidence References** (WORM vault URLs with hash verification)
- **Authority Approval** (digital signature + timestamp)
- **Appeal Rights Notice** (required template)
- **Next Steps** (if applicable: implementation procedures, appeal timeline)

---

## Tier 1: Administrative Duplicate Resolution

### Scope
Automated or near-automatic resolution of clerical duplicates created during enrollment. **Authority:** Enrollment Agent with supervisor review. **SLA:** 24 hours.

### Detection Scenarios
- **Name + DOB exact match**: Same biographic data, different NNI
- **Name typo + DOB match**: Obvious clerical error (e.g., "Jean Marie" vs. "Jean-Marie")
- **Soundex match + DOB ±3 months**: Name variation (e.g., "Jon" vs. "John"), DOB data entry error
- **Document number exact match**: Multiple enrollments with same identity document number

### Procedural Steps

#### Step 1: Auto-Detection & Assignment (0-2h)
**System automatically:**
1. Detects T1 candidate via fuzzy matching rules (soundex + levenshtein distance)
2. Calculates confidence score:
   - Exact name + exact DOB match = 100%
   - Soundex match + DOB ±7 days = 95%
   - Soundex match + DOB ±3 months = 85%
   - Document number match = 100%
3. Creates ConflictCase with `tier: T1`, `status: OPEN`, `detected_at: [timestamp]`
4. If confidence ≥95%, automatically assigns to primary NNI's enrollment agent
5. Routes to queue: `SNISID.T1.Pending` (Kafka topic)

#### Step 2: Enrollment Agent Review (2-12h)
**Enrollment agent must:**

**Checklist: T1 Duplicate Verification**
- [ ] Access case in conflict resolution portal
- [ ] Pull both NNI records (side-by-side comparison)
- [ ] Verify matching fields:
  - [ ] Name (exact or obvious typo?)
  - [ ] Date of birth (exact or clerical error ±3 months?)
  - [ ] Gender
  - [ ] Nationality
  - [ ] Document issuing dates (same source document?)
- [ ] Cross-check enrollment location/date (different agents? minutes apart?)
- [ ] If biometrics exist: verify 1:1 match (quality score ≥0.95)
- [ ] Verify no judicial holds or flags on either NNI

**Decision Logic:**
```
IF all identifying fields match OR document number matches:
  DECISION = MERGE
ELSE IF fields mostly match but 1-2 discrepancies:
  DECISION = REQUEST_CLARIFICATION (escalate to supervisor for phone call to citizen)
ELSE:
  DECISION = NOT_A_DUPLICATE (refer to T2 or close)
```

#### Step 3: Auto-Resolution (12-20h)
**If decision = MERGE:**
1. Supervisor approves agent decision (dual-control)
2. System executes merge:
   - Designate primary NNI (usually earlier enrollment)
   - Archive secondary NNI with reason: `ADMINISTRATIVE_DUPLICATE_MERGED`
   - Redirect all secondary NNI queries to primary
   - Copy all valid documents from secondary to primary
   - Biometric template consolidation (if both have biometrics, keep higher-quality)
3. Log events to audit trail:
   - `identity.merged` (source: secondary_nni, target: primary_nni, reason: administrative_duplicate)
4. Notify both enrolling agents + regional supervisor
5. Set case status → `RESOLVED`

**If decision = REQUEST_CLARIFICATION:**
1. Generate citizen contact form with questions
2. Contact citizen by phone/SMS/portal message
3. Require response within 5 days
4. Escalate to supervisor if citizen unavailable
5. Adjust case SLA +5 days during investigation

**If decision = NOT_A_DUPLICATE:**
1. Close case: `status: CLOSED`
2. Flag for potential T2/T3 escalation if biometric score >85%
3. Notify citizen (if contact exists)

#### Step 4: Citizen Notification (20-24h)
**Template: Administrative Merge Notification**
```
Subject: Your National Identity Number Consolidation

Dear [Name],

Our records indicate that you were enrolled with two identity numbers 
on [dates]. After verification, these records refer to the same person.

We have consolidated your records under NNI: [PRIMARY_NNI]
Your previous NNI [SECONDARY_NNI] is archived and will no longer be used.

All your documents and services now use [PRIMARY_NNI].

This consolidation does not affect your identity status or document validity.

You have 7 days to appeal this decision by contacting [regional ONI center].

Regards,
Office of National Identity
```

### Quality Assurance for Tier 1

**Supervisor Review Checklist (before final approval):**
- [ ] Enrollment agent followed decision logic correctly
- [ ] Both NNI records reviewed and documented
- [ ] Conflict of interest: neither agent had family relationship with citizen
- [ ] Documentation complete (summary + rationale ≥50 words)
- [ ] SLA on track (<24h from case creation)
- [ ] Citizen notification ready to send

**False Positive Rate Target:** <1% of T1 cases overturned on appeal

---

## Tier 2: Biometric Near-Match Adjudication

### Scope
ABIS match score 85–94% (suspicious but not conclusive). Requires human adjudicator review and two-person consensus. **Authority:** Biometric Adjudicator (two adjudicators required for consensus). **SLA:** 48 hours.

### Biometric Quality Standards

**Fingerprint (40% weight in final score):**
- Capture quality: ISO 19794-4 Level 3+ (≥500 DPI, <3 false minutiae per 100 mm²)
- Minutiae count: ≥20 usable minutiae per print
- Match threshold: ≥95% minutiae correspondence
- Special cases: Worn/scarred prints require ≥3 additional minutiae matching

**Iris (35% weight):**
- Capture quality: ISO/IEC 19794-6 Level 3+ (iris diameter ≥200 pixels, <3% occlusion)
- Hamming distance threshold: ≤0.28 (indicates match)
- Cross-eye matching: optional if same person may have different iris characteristics
- Cataract/surgery cases: noted in record, requires manual examination

**Facial (20% weight):**
- Capture quality: ISO/IEC 19794-5 Level 3 (frontal, neutral expression, no occlusion)
- 3D mesh comparison: ≥92% geometric correspondence
- Liveness detection: must pass PAD (Presentation Attack Detection) module
- Age variation: allow ±2 years due to natural aging

**Demographics (5% weight):**
- Name consistency: soundex match OR levenshtein distance ≤2
- DOB consistency: exact match preferred, allow ±3 months for data entry errors
- Gender consistency: must match exactly

### Procedural Steps

#### Step 1: Queue & Assignment (0-4h)
**System automatically:**
1. Routes T2 case to `SNISID.T2.Pending` (Kafka topic, SLA timer starts)
2. Randomly selects first adjudicator from available pool (weighted by geographic region)
3. Assigns case: `assigned_to: [ADJUDICATOR_ID_1]`
4. Sends notification: email + portal dashboard alert
5. Case status → `ASSIGNED`

#### Step 2: First Adjudicator Review (4-20h)

**Biometric Adjudicator Decision Matrix:**

| ABIS Score | Quality Issues | Biometric Confidence | Preliminary Decision |
|------------|----------------|----------------------|----------------------|
| 85–87% | None | Low-Medium | INCONCLUSIVE |
| 88–90% | Minor (worn print) | Medium | LIKELY_MATCH |
| 91–93% | Moderate (photo angle) | Medium-High | LIKELY_MATCH |
| 94%+ | None or minor | High | PROBABLE_MATCH |
| 94%+ | Major (PAD failure, old iris) | Medium | INCONCLUSIVE |

**Adjudicator Checklist:**
- [ ] Access ABIS match result (template hashes, NOT raw images per privacy)
- [ ] Review biometric quality scores for both records
- [ ] Examine demographic fields (name, DOB, gender, doc type)
- [ ] Cross-check document numbers (same issuing authority?)
- [ ] Check enrollment dates (recent? months apart?)
- [ ] Pull age-progression report (if age >10 years different)
- [ ] Review any prior notes/flags on either NNI
- [ ] Determine: LIKELY_MATCH, INCONCLUSIVE, or NOT_MATCH
- [ ] Document rationale (≥100 words: which biometrics matched, which didn't, confidence level)

**Decision Form (Adjudicator 1):**
```
Case: [CFR-YYYY-NNNNNN]
Adjudicator: [ID], [Name]
Review Time: [start_time] → [end_time] (duration in minutes)

Biometric Analysis:
  - Fingerprint Match: [YES/NO] (score: X%, minutiae: NN/NN)
  - Iris Match: [YES/NO] (Hamming distance: 0.XX)
  - Face Match: [YES/NO] (mesh correspondence: XX%)
  - Liveness Check: [PASS/FAIL]

Demographic Analysis:
  - Name: [EXACT_MATCH / SOUNDEX / DIFFERENT]
  - DOB: [EXACT / WITHIN_3MO / DIFFERENT]
  - Gender: [MATCH / DIFFERENT]
  - Nationality: [MATCH / DIFFERENT]

Overall Assessment:
  Preliminary Decision: [LIKELY_MATCH / INCONCLUSIVE / NOT_MATCH]
  Confidence Level: [HIGH / MEDIUM / LOW]
  
Rationale (≥100 words):
[Detailed explanation of analysis, confidence factors, alternative hypotheses]

Secondary Adjudicator Assignment: [ID] [auto-assigned by system]
```

#### Step 3: Second Adjudicator Independent Review (20-40h)
**System automatically:**
1. After Adjudicator 1 submits, randomly assigns Adjudicator 2 (different region)
2. **Masks** Adjudicator 1's decision (blind review)
3. Routes to Adjudicator 2: case status → `UNDER_REVIEW` (second review)

**Adjudicator 2 performs identical review (independent):**
- Same checklist
- Same decision matrix
- Same documentation requirements
- **No visibility** into Adjudicator 1's findings

#### Step 4: Consensus & Resolution (40-46h)

**System compares decisions:**

| Adj 1 | Adj 2 | Consensus Result | Action |
|-------|-------|------------------|--------|
| MATCH | MATCH | **CONSENSUS: MATCH** | Escalate to Senior Examiner (→ T3) |
| NOT_MATCH | NOT_MATCH | **CONSENSUS: NOT_MATCH** | Close case: status → CLOSED |
| INCONCLUSIVE | INCONCLUSIVE | **CONSENSUS: INCONCLUSIVE** | Escalate to Senior Examiner (→ T3) |
| MATCH | NOT_MATCH | **DISAGREEMENT** | Route to Senior Examiner (tiebreaker) |
| MATCH | INCONCLUSIVE | Weighted as MATCH | Escalate to Senior Examiner (→ T3) |
| NOT_MATCH | INCONCLUSIVE | Weighted as NOT_MATCH | Close case: status → CLOSED |

**Escalation to T3 Procedure (if consensus = MATCH or INCONCLUSIVE or DISAGREEMENT):**
1. Combine both adjudicator reports into unified brief
2. Highlight points of agreement/disagreement
3. Route to Senior Forensic Examiner with note: "T2 escalation"
4. Upgrade case: `tier: T3`, `status: ESCALATED`
5. Reset SLA timer to 72-hour T3 SLA
6. Notify both adjudicators + examiner

**Close Case Procedure (if consensus = NOT_MATCH):**
1. Generate case closure report (both adjudicator findings + consensus)
2. Case status → `CLOSED`
3. Archive to WORM: `vault://conflicts/CFR-YYYY-NNNNNN/closure`
4. Send citizen notification (if applicable):
   ```
   Your national identity record was reviewed as part of routine 
   quality assurance. No conflicts were identified. Your identity 
   status remains ACTIVE.
   ```

### Quality Assurance for Tier 2

**Supervisor Spot-Check (20% of cases):**
- [ ] Both adjudicators properly trained (recert valid)
- [ ] Decision matrix applied correctly
- [ ] Documentation ≥100 words per reviewer
- [ ] No evidence of bias (one adjudicator always choosing same option)
- [ ] SLA compliance (<48h)

**Monthly Metrics:**
- Mean consensus rate: target ≥90% (adjudicators agreeing in first two reviews)
- Mean time to consensus: target <40 hours
- Appeal overturn rate: target <5%

---

## Tier 3: Biometric Hard-Match Forensic Examination

### Scope
ABIS match score ≥95% or escalated from T2 disagreement. Requires expert forensic examiner with advanced training. **Authority:** Senior Forensic Examiner (dual-control with supervisor counter-signature). **SLA:** 72 hours.

### Forensic Examiner Qualifications

**Minimum Requirements:**
- ≥3 years biometric identification experience
- NIST-certified training in fingerprint matching (10-point AFIS standard)
- ISO/IEC 17024 certifications:
  - Fingerprint expert
  - Iris biometric expert (OR facial recognition expert)
- Annual recertification: ≥40 hours continuing education
- Annual audit of case accuracy: ≥95% agreement with reviewed cases
- Conflict-of-interest clearance: no cases within own region of residence

### Procedural Steps

#### Step 1: Case Assignment & Preparation (0-6h)

**Senior Examiner Intake Checklist:**
- [ ] Access case in forensic workstation (isolated network, audit-logged)
- [ ] Review case history (detection source, both adjudicator reports if T2 escalation)
- [ ] Request raw biometric templates from ABIS vault (audit logged)
- [ ] Verify chain-of-custody on all evidence (hash verification)
- [ ] Cross-check demographic data from both NNI records
- [ ] Review any prior criminal/fraud flags on either identity
- [ ] Determine examination strategy (which biometrics to prioritize)

**Forensic Examination Plan (document in case file):**
```
Examiner: [Name], [Cert ID], [Exam Region]
Date: [YYYY-MM-DD]
Priority Sequence:
  1. [PRIMARY_BIOMETRIC] (e.g., "Right Index Fingerprint")
  2. [SECONDARY_BIOMETRIC] (e.g., "Right Thumb Fingerprint")
  3. [TERTIARY_BIOMETRIC] (e.g., "Iris Comparison")
  4. Demographic Reconciliation
```

#### Step 2: Fingerprint Forensic Analysis (8-24h)

**10-Point Fingerprint Matching Standard (FBI/Interpol):**

For each finger comparison:
- **Minutiae Extraction**: Identify ridge endings and bifurcations
- **Minutiae Comparison**: Match ≥12 minutiae pairs from same fingers
- **Pattern Analysis**: Verify ridge patterns (whorl, loop, arch) match
- **Quality Assessment**: Score each print for usability (0-10 scale)

**Examiner Procedure:**

1. **Phase 1: Independent Pattern Assessment (4h)**
   - Without seeing ABIS score, examine fingerprint patterns visually
   - Document: "Print on NNI-A shows [pattern type], Print on NNI-B shows [pattern type]"
   - Compare ridge counts, core position, delta position
   - Preliminary opinion: LIKELY_MATCH, INCONCLUSIVE, DIFFERENT

2. **Phase 2: Minutiae-Level Analysis (6h)**
   - Extract all usable minutiae points
   - Create minutiae matrix for each print
   - Manual matching: identify corresponding minutiae pairs
   - Document each match with confidence level (HIGH: clear match, MEDIUM: reasonable match, LOW: questionable)

3. **Phase 3: Quality Verification (2h)**
   - Assign quality score to each print (0-10 scale):
     - 9-10: Excellent (clear ridge definition, no artifacts)
     - 7-8: Good (minor smudging, most minutiae visible)
     - 5-6: Fair (significant smudging, some minutiae obscured)
     - 3-4: Poor (heavy wear/scarring, ≥1/3 unclear)
     - <3: Unacceptable (unusable for matching)
   - If quality <5 on either print: note in report as "Quality-Compromised Match"

4. **Phase 4: Documentation (2h)**
   - Document all findings in forensic report
   - Include fingerprint quality comparison table
   - Final determination per finger:
     - MATCH (≥12 minutiae, patterns consistent)
     - POSSIBLE MATCH (9-11 minutiae, requires caution)
     - INCONCLUSIVE (quality poor, <9 minutiae)
     - NOT A MATCH (patterns different, <5 minutiae)

**Sample Fingerprint Comparison Table:**
| Finger | Quality (Rec) | Quality (Claim) | Minutiae Match | Pattern Match | Determination |
|--------|---------------|-----------------|----------------|---------------|--------------|
| R. Index | 8 | 7 | 15/18 (83%) | LOOP/LOOP | MATCH |
| R. Middle | 7 | 6 | 12/16 (75%) | WHORL/WHORL | MATCH |
| R. Ring | 6 | 5 | 10/14 (71%) | ARCH/ARCH | POSSIBLE |
| R. Little | 5 | 4 | 8/12 (67%) | LOOP/LOOP | INCONCLUSIVE |
| L. Thumb | 8 | 8 | 16/19 (84%) | WHORL/WHORL | MATCH |
| ... | ... | ... | ... | ... | ... |
| **Overall** | - | - | **61/79 (77%)** | **4/5 Match** | **STRONG MATCH** |

#### Step 3: Iris Biometric Analysis (6-12h, if applicable)

**Iris Matching Protocol:**

1. **Quality Assessment (2h)**
   - Check iris diameter: ≥200 pixels required
   - Measure occlusion: <10% acceptable
   - Verify: both images frontal, similar lighting
   - Assess: image sharpness, contrast, eyelash interference

2. **Feature Extraction (2h)**
   - Extract iris feature code (Daugman algorithm or equivalent)
   - Document: cryptogram bits, bit strength, occlusion masks

3. **Comparison (2h)**
   - Calculate Hamming distance between extracted codes
   - HD ≤0.28: MATCH (typically 0.10-0.25 range)
   - HD 0.28-0.35: INCONCLUSIVE
   - HD >0.35: NOT A MATCH
   - Document: "Hamming distance: 0.XX (Match)" or similar

**Iris Matching Report Section:**
```
Iris Analysis:
  Left Iris Quality: [ISO rating 0-4]
  Right Iris Quality: [ISO rating 0-4]
  
  Left Iris HD: 0.XX → [MATCH/INCONCLUSIVE/NO_MATCH]
  Right Iris HD: 0.XX → [MATCH/INCONCLUSIVE/NO_MATCH]
  
  Note: If quality <2, flag for manual re-capture recommended
```

#### Step 4: Facial Recognition Analysis (4-8h, if applicable)

**3D Facial Mesh Comparison:**
1. Extract 3D facial landmarks (nose, eyes, jawline, etc.)
2. Compare geometric distances between landmarks
3. Allow ±2 years for natural aging effects
4. Require ≥92% geometric correspondence for MATCH

#### Step 5: Demographic Reconciliation (2-4h)

**Cross-Check Name, DOB, Documents:**

| Field | Record A | Record B | Match | Reconciliation |
|-------|----------|----------|-------|-----------------|
| First Name | Jean | Jean | ✓ | EXACT |
| Last Name | Desroches | Desrosches | ✗ | Spelling variation (1 letter diff) - **Acceptable** |
| DOB | 1985-05-15 | 1985-05-18 | ~ | 3-day difference - **Possible data entry error** |
| Gender | M | M | ✓ | EXACT |
| Nationality | Haiti | Haiti | ✓ | EXACT |
| Doc Type (Primary) | National ID | National ID | ✓ | EXACT |
| Doc Issuer | ONI Port-au-Prince | ONI Cap-Haïtien | ✗ | **Different regions - requires explanation** |

**Reconciliation Decision:**
```
Demographic finding: Mostly aligned, with 3-month DOB discrepancy 
and 2-letter surname variation. Variations consistent with 
administrative data entry errors (clerical transcription).
Confidence in biometric match: CONFIRMED
```

#### Step 6: Final Determination & Report (4-6h)

**Forensic Examiner Report (≥300 words, template):**

```
FORENSIC EXAMINATION REPORT
SNISID Case: CFR-2026-001847
Examiner: [Name], [Cert ID]
Exam Date: 2026-05-23
Examination Type: T3 Biometric Hard-Match Forensic

EXECUTIVE SUMMARY:
Biometric analysis of NNI [A] and NNI [B] indicates [STRONG MATCH / 
PROBABLE MATCH / INCONCLUSIVE / NOT A MATCH]. Overall confidence: [HIGH / 
MEDIUM / LOW].

METHODOLOGY:
[Describe techniques used: fingerprint minutiae analysis, iris Hamming 
distance, facial mesh comparison, etc. Reference standards: FBI 10-print, 
ISO/IEC 19794 series]

FINDINGS:

Fingerprint Analysis:
  - Overall minutiae correspondence: 77% (61 of 79 minutiae matched)
  - Quality-adjusted score: 82% (accounting for quality discrepancies)
  - Determina: STRONG MATCH (≥5 fingers with ≥12 minutiae each)
  
Iris Analysis (if applicable):
  - Left Iris Hamming Distance: 0.18 → MATCH
  - Right Iris Hamming Distance: 0.22 → MATCH
  - Determination: STRONG MATCH

Facial Recognition (if applicable):
  - Geometric correspondence: 94% → MATCH
  - Liveness detection: PASS
  - Age variance: +1 year (within natural aging ±2 years) → ACCEPTABLE

Demographic Analysis:
  - Name: Soundex match with 1-letter surname variation (clerical error)
  - DOB: 3-day variance (data entry error, within ±3-month threshold)
  - Identity documents: Consistent issuer patterns, numbers validate

ALTERNATIVE HYPOTHESES CONSIDERED:
1. Could the fingerprints belong to close relatives? 
   → No: pattern types differ in 2/10 prints, minutiae don't align
2. Could biometric spoofing be involved?
   → No: liveness checks passed, multiple biometric modalities consistent
3. Could this be a near-identical twin?
   → Unlikely: iris HD ≤0.22 indicates likely same individual

CONCLUSION:
With [HIGH / MEDIUM / LOW] confidence, these biometric records refer to the 
same individual. Recommendation: [MERGE / DEACTIVATE_SECONDARY / 
REFER_JUDICIAL / REQUIRE_ADDITIONAL_EVIDENCE]

Examiner Signature: _________________ Date: _________
```

#### Step 7: Supervisor Counter-Signature (2-4h)

**Supervisor Review Checklist (T3 dual-control):**
- [ ] Examiner has valid certifications (not expired)
- [ ] Examiner has no conflict of interest in case
- [ ] Forensic methodology matches NIST/FBI standards
- [ ] Report ≥300 words with clear reasoning
- [ ] Confidence level justified by evidence
- [ ] Alternative hypotheses considered
- [ ] Quality-adjusted scoring applied correctly
- [ ] No apparent bias in analysis
- [ ] SLA on track (<72h from assignment)

**Supervisor Approval:**
```
Reviewed by: [Supervisor Name], [Cert ID], [Date/Time]
Approval: ✓ APPROVED / ⚠ REQUIRES REVISION / ✗ REJECTED

Comments:
[If revision required, specific issues to address]

Supervisor Signature: _________________ Digital Signature Hash: _______
```

### Resolution Decisions (T3)

**Based on examiner confidence level:**

| Confidence | Recommendation | Next Step |
|------------|-----------------|-----------|
| HIGH (≥95%) | MERGE | Immediate merge authorization (dual-control approval) |
| MEDIUM (85-94%) | MERGE WITH CONDITIONS | Merge approved; secondary record flagged for audit |
| LOW (<85%) | INCONCLUSIVE | Escalate to Regional Director (→ T4) |
| Disagreement | JUDICIAL REFERRAL | Escalate to legal counsel (→ T6) |

---

## Tier 4: Multi-Record Conflict Resolution

### Scope
Multiple NNIs confirmed to refer to same person, but with contradictory biographic data (marriage name changes, address conflicts, nationality ambiguity). Requires regional director + legal counsel coordination. **Authority:** Regional ONI Director + SNISID Legal Counsel. **SLA:** 7 days.

### Procedural Steps

#### Step 1: Case Intake & Multi-Record Consolidation (0-12h)

**Regional Director Checklist:**
- [ ] Access all conflicting NNI records (typically 2-3, sometimes ≥4)
- [ ] Extract biographic data into reconciliation matrix
- [ ] Identify source of each data element (enrollment dates, documents cited)
- [ ] Assess which record is "authoritative" (most recent comprehensive enrollment)
- [ ] Flag any marriage/divorce records, name change documents
- [ ] Flag any address history discrepancies
- [ ] Cross-check with Civil Registry (marriage certificates, divorce decrees)

#### Step 2: Data Reconciliation Matrix

**Sample Multi-Record Conflict:**

| Field | NNI-A (2019) | NNI-B (2023) | Source NNI-A | Source NNI-B | Reconciliation |
|-------|--------------|--------------|--------------|--------------|----------------|
| First Name | Jean | Jean | Cert of Birth | Birth Cert | ✓ Match |
| Last Name | Desroches | Rousseau | Cert of Birth | Marriage Cert (2021) | **Different - married name change** |
| DOB | 1985-05-15 | 1985-05-15 | Birth Cert | Birth Cert | ✓ Match |
| Gender | M | M | Birth Cert | Birth Cert | ✓ Match |
| Address | Port-au-Prince | Cap-Haïtien | 2020 Reg. Update | 2023 Re-enroll | Different - moved |
| Nationality | Haiti | Haiti | Birth Cert | Birth Cert | ✓ Match |
| Status | ACTIVE | PENDING | - | Recent enroll | **Conflict** |

#### Step 3: Legal Counsel Review (12-48h)

**Legal Analysis:**
1. Review marriage/divorce documents (verify authenticity with DGAE - civil registry)
2. Assess: Is name change legally valid? Does it justify separate NNI?
3. Assess: Address change legitimate (relocation documented)?
4. Assess: Which NNI should be primary (older typically, or more recent if comprehensive)?
5. Identify: Any legal issues (e.g., disputed marriage, ongoing divorce proceedings)

**Legal Counsel Decision Matrix:**

| Scenario | Recommendation |
|----------|-----------------|
| Legal marriage w/ name change cert → Merge to newer record (spouse's surname) |
| Divorce (no name change) → Merge to primary NNI |
| Address relocation (utility bills confirm) → Merge, update address to latest |
| Name change without legal documentation → Escalate to T5 (fraud investigation) |
| Disputed marriage (conflicting documents) → Refer to T6 (judicial) |
| Multiple nationalities claimed → Refer to T6 (judicial) |

#### Step 4: Merge Decision & Implementation (48-120h)

**Regional Director Final Decision (dual-control: Director + Legal Counsel signature):**

```
MERGE AUTHORIZATION - TIER 4 MULTI-RECORD CONFLICT

Case: CFR-2026-001847
Primary Record (SURVIVING): [NNI-A] [Name]
Secondary Record (ARCHIVED): [NNI-B] [Name]
Reason: [Legal marriage certificate / Address relocation / Administrative consolidation]

Implementation Plan:
  1. Designate [NNI-A] as surviving record
  2. Consolidate name: [NNI-A] = [Official name per legal documents]
  3. Consolidate address: [Latest confirmed address]
  4. Consolidate documents: Copy valid docs from NNI-B to NNI-A
  5. Archive NNI-B with status: MERGED_TO_[NNI-A]
  6. Update all service points (voter registry, benefits, etc.) within 5 days

Approved by:
  Regional Director: _________________ Date: _______
  Legal Counsel: _________________ Date: _______
```

#### Step 5: Citizen Notification (120-168h)

**Template: Multi-Record Consolidation Notice**

```
Subject: Your National Identity Records Have Been Consolidated

Dear [Name],

Our records indicate that you held multiple national identity numbers 
due to administrative updates (name change, address relocation, etc.):
  - Previous NNI: [NNI-B]
  - Current NNI: [NNI-A]

These records have been consolidated under a single NNI for your convenience.

Your name in the system is now: [Official name per legal documents]
Your registered address: [Current address]

All services and documents now use NNI [NNI-A]. The previous NNI [NNI-B] 
is archived and will no longer be used.

You have 14 days to appeal this decision by contacting [Regional ONI Director].

Regards,
Office of National Identity
```

---

## Tier 5: Fraud Investigation

### Scope
Deliberate duplicate enrollment, synthetic/spoofed biometrics, or identity theft. Requires DCPJ (Directorate for Combating Organized Crime) investigation with potential judicial referral. **Authority:** DCPJ Fraud Unit, National Identity Appeals Board (NIAB). **SLA:** 30 days.

### Detection Scenarios
- **Synthetic biometrics**: AI-generated fingerprints (GAN artifacts, statistical anomalies)
- **Biometric spoofing**: Liveness detection failures (2D photo, 3D mask, screen replay)
- **Identity theft**: Fraudulent use of another person's documents
- **Velocity abuse**: Multiple enrollments by same biometric in short timeframe
- **UEBA anomalies**: Behavior patterns inconsistent with normal citizen usage

### Procedural Steps

#### Step 1: Fraud Suspicion Flagging & DCPJ Intake (0-24h)

**Fraud Detection Triggers (AI Engine):**
- ABIS score ≥99.5% (suspiciously perfect match)
- Biometric quality markers: statistical anomalies in feature extraction
- Liveness detection: PAD (Presentation Attack Detection) failure after ≥2 retries
- Velocity check: >3 enrollment attempts with different demographics in 7 days
- Document analysis: MRZ checksum failure, UV watermark anomalies

**DCPJ Intake Process:**
1. Flagged case routed to: `DCPJ.Fraud.Intake` queue
2. DCPJ Investigator assigned (random rotation)
3. Case status: `ESCALATED`, `tier: T5`
4. SLA timer: 30-day investigation window

**DCPJ Investigator Initial Checklist:**
- [ ] Access case file (including flagged biometric analysis)
- [ ] Interview enrollment agent (what did citizen say? any unusual behavior?)
- [ ] Pull metadata: IP addresses, device fingerprints, geolocation at enrollment
- [ ] Check for prior fraud flags in SNISID or law enforcement databases
- [ ] Preliminary assessment: Is this likely fraud or false positive?

#### Step 2: Investigation Protocols

**Sub-Protocol A: Synthetic Biometric Detection**

```
Technical Analysis (2-3 days):
  1. Extract raw biometric template features
  2. Run statistical anomaly detection (compare to known-good gallery)
  3. Analyze: ridge flow discontinuities, impossible minutiae patterns
  4. PAD analysis: color channel inconsistencies, frequency domain artifacts
  5. Generate: "Synthetic Confidence Score" (0-100%)
  
  Documentation:
    - Statistical anomalies found: [yes/no]
    - GAN fingerprint artifacts detected: [yes/no]
    - Synthetic confidence: [0-100%]
    - Recommendation: [LIKELY_SYNTHETIC / INCONCLUSIVE / LIKELY_REAL]
```

**Sub-Protocol B: Identity Theft Investigation**

```
Person-to-Person Interview (3-5 days):
  1. DCPJ visits residence at address on primary NNI
  2. Conducts in-person interview with claimed citizen
  3. Documents: "Does citizen recognize the secondary NNI / enrollment date / 
     location? Is this person aware of this enrollment?"
  4. If citizen denies: likely identity theft
     If citizen admits: likely voluntary duplicate or testing enrollment
  5. Interview documented: video (with consent) + written statement
  
  Follow-up:
    - If identity theft: file complaint with national police
    - Preserve evidence: interview recording, signed statement
    - Notify primary identity of potential theft
```

**Sub-Protocol C: Velocity & UEBA Analysis**

```
Behavioral Analysis (2-3 days):
  1. Pull enrollment logs: dates, times, locations, agents involved
  2. Check: Are enrollments by same person? Same documents? Different documents?
  3. Geolocation check: Are enrollment locations geographically impossible? 
     (e.g., Port-au-Prince + Cap-Haïtien 3 hours apart, same day)
  4. Device fingerprinting: Same device used for multiple enrollments?
  5. Agent involvement: Is one enrollment agent involved in multiple fraud cases?
  
  Decision Tree:
    IF same documents + same biometric + impossible travel
      → Likely fraud, recommend criminal investigation
    IF different documents + different biometric + same velocity
      → Possible credential mill, recommend surveillance
    IF one enrollment agent in multiple cases
      → Possible insider threat, recommend internal affairs review
```

#### Step 3: Warrant Request (if applicable, 5-10 days)

**Grounds for Judicial Warrant:**
- Sufficient evidence of identity theft OR synthetic biometric
- Need to access citizen's devices (phone, computer) for digital evidence
- Need to subpoena third parties (enrollment agents, documents issuers)

**Warrant Package to Public Prosecutor:**
1. Case summary (what fraud is suspected)
2. Technical evidence (synthetic biometric analysis, velocity data)
3. Interview findings (if applicable)
4. Legal grounds (which criminal statute)
5. Specific items to be seized/examined

**Prosecutor Decision:** Warrant granted → DCPJ can conduct search/seizure. Warrant denied → Refer to T6 (judicial appeal).

#### Step 4: Evidence Collection & Preservation (10-25 days)

**Digital Evidence Protocol:**
- Access encryption keys in Vault (RSA-4096 sealed)
- Extract evidence: interview recordings, biometric forensics, metadata
- Maintain chain-of-custody: every access logged to WORM
- Hash all evidence files (SHA-256) for integrity verification

**Evidence Vault Sealing:**
```
Evidence Package: CFR-2026-001847
Sealed: 2026-05-30T14:00:00Z
Hash: sha256:a4f8c2e1...
Dual-Key Required for Access: 
  Key 1 (DCPJ): [key-id-1]
  Key 2 (SNISID Legal): [key-id-2]
Expiry: 2033-05-30 (7-year retention for potential appeal)
```

#### Step 5: DCPJ Final Decision (25-30 days)

**DCPJ Investigator Report & Recommendation:**

```
FRAUD INVESTIGATION REPORT

Case: CFR-2026-001847
Investigator: [DCPJ Agent Name], Badge [ID]
Investigation Period: 2026-05-23 → 2026-06-20

FINDINGS:

[Technical analysis results: synthetic biometric confidence score, 
 velocity analysis, device fingerprinting results]

[Interview results: citizen statement, enrollment agent statement, 
 witness statements if applicable]

[Warrant results: devices seized, digital evidence extracted, 
 documents obtained]

CONCLUSION:

Fraud Probability: [HIGH / MEDIUM / LOW]
Type of Fraud: [Synthetic Biometric / Identity Theft / Credential Mill / Other]
Criminal Statute Violated: [if applicable]

RECOMMENDATION:

[ ] FRAUD CONFIRMED - Recommend criminal prosecution
    Secondary NNI to be REVOKED with status: FRAUD_CONFIRMED
    Refer to Public Prosecutor for charging

[ ] FRAUD SUSPECTED - Insufficient evidence for prosecution
    Secondary NNI to be REVOKED with status: FRAUD_SUSPECTED
    Case closed; evidence retained for 7 years

[ ] NOT FRAUD - Likely legitimate explanation
    Case closed; both NNIs remain ACTIVE
    Recommend re-enrollment at proper ONI center

Investigator Signature: _________________ Date: _________
```

#### Step 6: NIAB Appeal Board Review (if T5 decision appealed)

**National Identity Appeals Board (NIAB) Composition:**
- 3 members: Senior DCPJ investigator, Legal counsel, Independent forensic expert
- Meets monthly or as needed
- 14-day review period for appeals

**NIAB Authority:**
- Override DCPJ recommendation (rarely)
- Request additional investigation
- Approve fraud confirmation + revocation
- Order restitution if applicable

---

## Tier 6: Judicial Escalation

### Scope
Court-contested identity, inheritance disputes, nationality challenges. Escalated to civil/family tribunals. **Authority:** Tribunals de Paix, Tribunaux Civil, Tribunal de Cassation. **SLA:** 90 days (provisional; court's timeline is binding).

### Procedural Steps

#### Step 1: Judicial Referral Package Preparation (0-5 days)

**Legal Counsel Preparation:**
1. Draft Court Referral Motion (identifies specific legal issue)
2. Compile evidence package (all case files, forensic reports, witness statements)
3. Secure expert witness availability (if needed)
4. Determine applicable tribunal (Paix for simple cases, Civil for complex)

**Court Referral Motion Template:**
```
MOTION FOR JUDICIAL DETERMINATION
Submitted to: Tribunal Civil, District of Port-au-Prince
Case: National Identity Conflict CFR-2026-001847

PARTIES:
  Plaintiff: Office of National Identity (SNISID)
  Interested Parties: [Citizen A], [Citizen B]
  
FACTS:
[Narrative of conflict: duplicate NNIs detected, forensic analysis results, 
 administrative resolution attempts made]

LEGAL ISSUE:
Whether identity records NNI-[A] and NNI-[B] refer to the same person, 
and if so, how to legally consolidate them while preserving citizen rights.

EVIDENCE SUBMITTED:
  - Forensic examination report (Dr. [Name], Senior Examiner)
  - Biometric ABIS analysis
  - Witness statements (enrollment agents, family members if applicable)
  - Civil registry documents (birth certificate, marriage certificate)
  
REQUESTED RELIEF:
The Court order consolidation of NNI-[B] into NNI-[A], with primary 
name to be determined per legal documentation.

Respectfully submitted,
[SNISID Legal Counsel]
Date: [YYYY-MM-DD]
```

#### Step 2: Court Coordination & Expert Witness Preparation (5-20 days)

**Expert Witness Guide:**

**Pre-Trial Preparation:**
- Review all case materials (forensic reports, ABIS analysis)
- Prepare testimony outline (≥3 pages, non-technical summary)
- Anticipate cross-examination questions
- Practice testimony delivery (clear, plain language, no jargon)

**Courtroom Testimony (Typical Q&A):**
```
Prosecutor: "Dr. Examiner, in your professional opinion, do the fingerprints 
            on NNI-A and NNI-B match?"
Witness:    "Yes. I found 15 corresponding minutiae points on the right 
            index finger, 12 on the right middle, and 14 on the right ring 
            finger. The pattern types also matched across all fingers. With 
            high confidence, these are the same individual's fingerprints."

Defense:    "Doctor, could these be from a twin?"
Witness:    "No. Fingerprints are unique even to identical twins. I also 
            examined iris biometrics, which show the same uniqueness. The 
            likelihood of two unrelated individuals having these exact patterns 
            is less than 1 in 100 billion."
```

#### Step 3: Tribunal Hearing & Court Determination (20-80 days)

**Tribunal Decision Authority:**
- **Tribunal de Paix**: Small identity disputes, simple factual issues
- **Tribunal Civil**: Complex identity conflicts, name change disputes, inheritance implications
- **Tribunal de Cassation**: Appeals of Tribunal Civil decisions

**Tribunal Judgment Components:**
```
COURT JUDGMENT
Case: Identity Conflict CFR-2026-001847
Tribunal: Tribunal Civil, Port-au-Prince
Judge: [Honorable Judge Name]
Date: 2026-07-15

HOLDING:
It is hereby ORDERED that national identity records NNI-[A] and NNI-[B] 
refer to the same natural person, namely [Official Name per court determination].

The two identity records shall be consolidated under NNI-[A], with the 
individual's official name recorded as [Court-Determined Name].

National identity record NNI-[B] is hereby CANCELLED and archived with 
annotation: "Consolidated into NNI-[A] by Court Order [Case Number]."

[Court Seal + Judge Signature]
```

#### Step 4: Implementation of Court Order (80-90 days)

**SNISID Compliance Procedure:**
1. Receive court order (certified copy)
2. Verify order authenticity (court seal, judge signature, case docket)
3. Execute consolidation per court order terms
4. Update all public registries (voter registry, benefits, passport database)
5. Issue new national ID card (with court-ordered name) within 10 days
6. Notify citizen (in writing) of consolidation and new NNI
7. Preserve court order in WORM: `vault://judicial/[case_number]/order`

**Notification to Citizen:**
```
Subject: Court Order - Identity Record Consolidation

Dear [Official Name],

The Tribunal Civil of Port-au-Prince has issued a Court Order 
(Case No. [####]) consolidating your national identity records.

Your official national identity number is: [NNI-A]
Your official name (as ordered by court): [Official Name]

A new national ID card reflecting this order will be issued within 10 days. 
You can collect it at [nearest ONI center].

This consolidation does not affect your rights, benefits, or status.

Regards,
Office of National Identity
```

---

## Common Operational Tasks

### Task: Conflict of Interest Screening

**Automated Screening (real-time):**
1. When case assigned to adjudicator/examiner, system checks:
   - Geographic home commune vs. case location (must be ≥50 km apart)
   - Family relationship: NNI cross-checked against system relationships
   - Professional relationship: prior cases with either NNI
2. If conflict detected: case **automatically** reassigned
3. If no automatic reassignment available: escalate to supervisor

**Manual Verification (adjudicator):**
Before reviewing case, adjudicator must certify:
```
CONFLICT OF INTEREST CERTIFICATION
Adjudicator: [Name], [ID]
Case: [CFR-YYYY-NNNNNN]

I certify that:
  [ ] I have no family relationship with either party to this case
  [ ] I reside >50 km from both case locations
  [ ] I have no prior professional relationship with either party
  [ ] I have no financial interest in outcome of this case

Signature: _________________ Date: _________
```

### Task: SLA Tracking & Auto-Escalation

**Real-Time SLA Monitoring:**
- Dashboard shows: Case created time, current time, SLA deadline, % elapsed
- Alerts: 50% elapsed → notify supervisor; 75% elapsed → escalate to next authority

**Escalation Automation:**
```
IF time_elapsed > SLA_deadline * 0.50:
  NOTIFY (supervisor, "Case overdue, escalation imminent")
  
IF time_elapsed > SLA_deadline * 0.75:
  ESCALATE (case to next authority level)
  NOTIFY (citizen, "Your case is being escalated due to processing delays")
  INCIDENT_REPORT (SLA breach detected)

IF time_elapsed > SLA_deadline * 1.00:
  AUTO_ESCALATE (case to next authority)
  NOTIFY (supervisor + director, "Critical SLA breach")
  GOVERNANCE_REVIEW (flag for audit)
```

### Task: Evidence Sealing & Chain-of-Custody

**When Case is RESOLVED or CLOSED:**
1. Compile all case evidence (forensic reports, biometric templates, interview recordings)
2. Generate manifest: list all evidence items with SHA-256 hashes
3. Create sealed evidence package: encrypt with dual-key (RSA-4096)
   - Key 1: SNISID Records Management
   - Key 2: Office of Inspector General
4. Store in WORM vault: `vault://conflicts/[case_id]/sealed/[date]`
5. Access is logged: every retrieval requires both key-holders
6. Retention period: 7 years (allows for appeals)

**Chain-of-Custody Access Log (example):**
```
Case: CFR-2026-001847
Sealed: 2026-06-15T10:00:00Z

Access History:
  2026-07-01 10:30 - DCPJ-Investigator-003 (Warrant execution, 45 min)
  2026-07-15 14:15 - Court-Clerk (Evidence presentation, 2h)
  2026-08-20 09:00 - Appeals-Board-Reviewer (Appeal review, 1.5h)
  
No unauthorized access detected.
Hash integrity: VERIFIED (all hashes match original manifest)
```

---

## Quality Assurance Checklists

### General QA: Every Tier

**Pre-Decision Checklist (before any resolution decision):**
- [ ] Case evidence: all files present and integrity verified (hash match)
- [ ] Chain-of-custody: complete access log with no gaps
- [ ] Documentation: decision rationale ≥50 words (T1-T2) or ≥100 words (T3+)
- [ ] Conflict of interest: verified no relationship exists
- [ ] SLA compliance: decision within tier timeframe
- [ ] Authority: decision-maker has required certifications/authority
- [ ] Dual-control (T3+): supervisor counter-signature present
- [ ] Appeal rights: citizen notification includes appeal deadline + procedure
- [ ] No regulatory violations: decision complies with applicable law

**Post-Decision Checklist (before case closure):**
- [ ] Resolution implemented (merge, revocation, referral, etc.)
- [ ] Citizen notified (within 24-48h per tier)
- [ ] Public registries updated (voter, benefits, etc.)
- [ ] Evidence sealed and archived
- [ ] Case file complete (all documentation present)
- [ ] Audit trail recorded (immutable event log)

### Tier-Specific QA

**T1 Monthly Audit:**
- Review 20% of closed T1 cases
- Verify: no false merges (citizens were actually same person)
- Measure: false positive rate target <1%

**T2 Quarterly Audit:**
- Review all T2 → T3 escalations
- Verify: adjudicators properly applied quality standards
- Measure: mean consensus rate ≥90%

**T3 Annual Forensic Review:**
- Chief Examiner reviews 10% of all T3 cases
- Verify: methodology matched NIST standards
- Measure: examiner accuracy ≥95% (cases consistent with final outcomes)

**T5 Fraud Conviction Rate:**
- Track cases referred to prosecution vs. convictions
- Target: ≥80% of fraud-confirmed cases result in conviction

---

**Document Version:** 2.0
**Last Updated:** 2026-05-23
**Next Review:** 2026-08-23
