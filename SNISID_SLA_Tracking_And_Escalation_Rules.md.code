# SNISID: SLA Tracking & Auto-Escalation Rules
## Real-Time Monitoring, Escalation Triggers, Metrics & Compliance Reporting

---

## Table of Contents
1. [SLA Definitions & Timelines](#sla-definitions--timelines)
2. [Auto-Escalation Rules Engine](#auto-escalation-rules-engine)
3. [Real-Time Monitoring Dashboard](#real-time-monitoring-dashboard)
4. [SLA Compliance Reporting](#sla-compliance-reporting)
5. [Escalation Audit Trail](#escalation-audit-trail)
6. [SLA Variance & Root-Cause Analysis](#sla-variance--root-cause-analysis)

---

## SLA Definitions & Timelines

### Tier-Based SLA Matrix

| Tier | Classification | SLA Duration | Status Transitions | Escalation Timeline |
|------|-----------------|--------------|-------------------|-------------------|
| **T1** | Admin Duplicate | **24 hours** | OPEN (0-4h) → ASSIGNED (4h) → UNDER_REVIEW (6h) → RESOLVED (20h) → CLOSED (24h) | 50% at 12h, 75% at 18h, 100% at 24h |
| **T2** | Biometric Near-Match | **48 hours** | OPEN (0-4h) → ASSIGNED (4h) → UNDER_REVIEW (20h) → ESCALATED/CLOSED (40h+) | 50% at 24h, 75% at 36h, 100% at 48h |
| **T3** | Biometric Hard-Match | **72 hours** | OPEN (0-4h) → ASSIGNED (4h) → UNDER_REVIEW (30h) → RESOLVED/ESCALATED (60h+) | 50% at 36h, 75% at 54h, 100% at 72h |
| **T4** | Multi-Record Conflict | **7 days** | OPEN (0-12h) → ASSIGNED (12h) → UNDER_REVIEW (48h) → RESOLVED/APPEALED (168h+) | 50% at 3.5d, 75% at 5.25d, 100% at 7d |
| **T5** | Fraud Investigation | **30 days** | OPEN (0-24h) → ESCALATED (24h) → UNDER_REVIEW (10d) → RESOLVED/PENDING_JUDICIAL (25d+) | 50% at 15d, 75% at 22.5d, 100% at 30d |
| **T6** | Judicial Escalation | **90 days** | OPEN (0-24h) → ESCALATED (24h) → PENDING_JUDICIAL (30d+) → RESOLVED (80d+) | Court timeline overrides (not subject to ONI SLA) |

### SLA Clock Mechanics

**SLA Timer Starts:**
- Case created by automated detection engine: `case_created_at = [timestamp]`
- Clock is suspended for: citizen information requests (paused until response), court proceedings (suspended until judgment), inter-agency coordination (paused until response)
- All paused time does NOT count against SLA (only active investigation time counts)

**SLA Timer Stops:**
- Final decision made and approved: `decision_completed_at = [timestamp]`
- Case status transitions to RESOLVED or CLOSED
- All notifications sent to citizen
- WORM archive complete

**Elapsed Time Calculation:**
```
SLA_elapsed_hours = (current_time - case_created_at) - paused_time_total
SLA_elapsed_percentage = (SLA_elapsed_hours / SLA_duration_hours) * 100
SLA_status = {
  "healthy": SLA_elapsed_percentage < 50%
  "warning": 50% ≤ SLA_elapsed_percentage < 75%
  "critical": 75% ≤ SLA_elapsed_percentage < 100%
  "breach": SLA_elapsed_percentage ≥ 100%
}
```

---

## Auto-Escalation Rules Engine

### Rule Set 1: Time-Based Escalation

**Rule 1A: 50% SLA Elapsed → Supervisor Alert**

```
TRIGGER:
  IF (case.elapsed_percentage ≥ 50%) AND (case.status ≠ RESOLVED):
    AND (no_prior_50_alert_sent):
      ACTION_TRIGGERED = true

ACTION:
  1. Send alert notification to case supervisor
  2. Message: "Case [CFR-YYYY-NNNNNN] has reached 50% SLA. 
              Current status: [status]. Estimated completion: [hours_remaining]."
  3. Alert medium: Email + Portal dashboard (red badge)
  4. Supervisor action: Review case, provide status update, or escalate if blocked
  5. Set flag: no_prior_50_alert_sent = false (prevent duplicate alerts)

AUDIT_LOG:
  - Event: SLA_50_ALERT_SENT
  - Case: [case_id]
  - Time: [timestamp]
  - Supervisor notified: [supervisor_id]
  - Status before alert: [status]
  - Elapsed time: [hours:minutes]

NO_AUTOMATIC_ESCALATION (supervisor must act if escalation needed)
```

**Rule 1B: 75% SLA Elapsed → Auto-Escalate to Next Authority**

```
TRIGGER:
  IF (case.elapsed_percentage ≥ 75%) AND (case.status ≠ RESOLVED):
    AND (case.escalation_count < MAX_ESCALATIONS_ALLOWED):
      ACTION_TRIGGERED = true

MAX_ESCALATIONS_ALLOWED = {
  "T1": 1,   // T1 can escalate once (to Level 2)
  "T2": 2,   // T2 can escalate twice (Level 2 → Level 3)
  "T3": 1,   // T3 can escalate once (to Level 4)
  "T4": 1,   // T4 can escalate once (to Level 4)
  "T5": 2,   // T5 can escalate twice (DCPJ → NIAB)
  "T6": 1    // T6 cases don't escalate (court is final)
}

ACTION:
  1. Identify NEXT_AUTHORITY (based on current authority + tier)
  2. Reassign case: case.assigned_to = [NEXT_AUTHORITY_ID]
  3. Reset SLA timer: case.reassignment_time = [timestamp]
  4. Increment escalation counter: case.escalation_count += 1
  5. Update case status: case.status = ESCALATED
  6. Add escalation record to audit trail with reason: "SLA_75_PERCENT_BREACH"

NOTIFICATIONS:
  - To NEXT_AUTHORITY: "Case inherited: 75% of previous SLA elapsed. 
                        Complete review and decision ASAP. 
                        Remaining SLA time: [X hours]"
  - To CITIZEN: "Your case is being escalated to ensure timely resolution. 
                 Your new case contact: [next_authority_name]"
  - To SUPERVISORY_CHAIN: "SLA escalation triggered - Case [CFR-YYYY-NNNNNN] 
                           escalated from [previous_authority] to [next_authority]"

AUDIT_LOG:
  - Event: SLA_75_AUTO_ESCALATION
  - Case: [case_id]
  - Time: [timestamp]
  - From authority: [previous_authority_id]
  - To authority: [next_authority_id]
  - Reason: SLA breach (75% elapsed)
  - Remaining SLA: [hours:minutes]
  - Escalation count: [N of MAX]
```

**Rule 1C: 100% SLA Breach → Force Escalation + Incident Report**

```
TRIGGER:
  IF (case.elapsed_percentage ≥ 100%) AND (case.status ≠ RESOLVED):
      ACTION_TRIGGERED = true

ACTION:
  1. Force escalate case (even if MAX_ESCALATIONS_ALLOWED exceeded)
  2. Case status: case.status = ESCALATED
  3. Create incident report: [INCIDENT-YYYY-NNNNNN]
  4. SLA breach severity: CRITICAL

INCIDENT_REPORT CONTENTS:
  - Case ID: [CFR-YYYY-NNNNNN]
  - Case tier: [T1-T6]
  - Original SLA deadline: [YYYY-MM-DD HH:MM:SS]
  - Days overdue: [X days Y hours]
  - Authority responsible: [authority_id + name]
  - Timeline of escalations: [all prior escalations with times]
  - Reason for delay (if documented): [reason]
  - Current case status: [status]
  - Remedial action taken: [escalated to next authority]
  - Notification sent to: [citizen + supervisory chain]

NOTIFICATIONS:
  - To NATIONAL_DIRECTOR: "CRITICAL: Case [CFR-YYYY-NNNNNN] SLA breach 
                           (>100% elapsed). Incident [INCIDENT-YYYY-NNNNNN] 
                           created. Authority [name] overdue. Escalated to next level."
  - To CITIZEN: "Your case exceeded our processing timeline. We are prioritizing 
                 your case for immediate resolution. Updated status will be 
                 provided within 24 hours."
  - To GOVERNANCE_TEAM: "SLA breach incident created for governance audit"

ESCALATION (if still not resolved):
  - Next authority receives urgency flag: [ESCALATION_PRIORITY = HIGH]
  - Dashboard highlights case in red: "OVERDUE - REQUIRES IMMEDIATE ACTION"
  - Supervisor receives escalation confirmation requirement (must acknowledge)

AUDIT_LOG:
  - Event: SLA_100_CRITICAL_BREACH
  - Case: [case_id]
  - Time: [timestamp]
  - Incident created: [incident_id]
  - Days overdue: [X days]
  - Escalated to: [next_authority_id]
  - National Director notified: [yes/no]

SUBSEQUENT_ACTION (within 24 hours):
  IF case still unresolved after escalation:
    1. National Director personally reviews case
    2. If authority at fault: Performance review / Disciplinary action
    3. If resource constraint: Allocate additional staff / Extend SLA justification
    4. If blocked by external factor: Coordination with external party
```

---

### Rule Set 2: Escalation Due to Conflict of Interest

**Rule 2A: Automated Conflict Detection → Reassignment**

```
TRIGGER:
  IF (case assigned to authority) AND (automated conflict-of-interest 
      check returns CONFLICT):
      ACTION_TRIGGERED = true

CONFLICT_TYPES:
  1. Geographic: Authority home commune = case location (distance <50 km)
  2. Family relationship: Either conflicting NNI in authority's family tree
  3. Prior employment: Authority has prior case with either party (<2 years)
  4. Financial interest: Authority has undisclosed financial interest

ACTION:
  1. Case reassignment: Select next available authority (different region)
  2. Original authority: Blocked from accessing case
  3. Case status: ESCALATED (due to conflict of interest)
  4. No SLA extension: Original SLA remains in effect

AUDIT_LOG:
  - Event: CONFLICT_OF_INTEREST_DETECTED_AUTO_REASSIGNMENT
  - Case: [case_id]
  - Time: [timestamp]
  - Original authority: [authority_id] (marked as conflicted)
  - Conflict type: [geographic / family / prior / financial]
  - New authority: [new_authority_id]
  - Reason: [specific conflict description]
```

**Rule 2B: Manual Conflict Discovery (post-decision) → Decision Rescission**

```
TRIGGER:
  IF (case decision already signed) AND (conflict of interest discovered):
      ACTION_TRIGGERED = true

ACTION:
  1. Void decision: case.resolution.status = VOID
  2. Case status: ESCALATED → requires new review
  3. Audit note: "Decision rescinded due to undisclosed conflict of interest"
  4. New authority assigned: Different region, conflict-cleared
  5. Case requires de novo (from-scratch) review

NOTIFICATIONS:
  - To CITIZEN: "We have discovered that your case decision was made by 
                an authority with a conflict of interest. Per SNISID standards, 
                this decision is being rescinded and your case will be reviewed 
                by a different authority. We apologize for this delay."
  - To ORIGINAL_AUTHORITY: "Conflict of interest discovered in Case [CFR]. 
                            Your decision has been rescinded. Disciplinary 
                            review may follow."
  - To SUPERVISORY_CHAIN: "Conflict-of-interest failure - decision rescinded"

AUDIT_LOG:
  - Event: CONFLICT_OF_INTEREST_POST_DECISION_RESCISSION
  - Case: [case_id]
  - Time: [timestamp]
  - Original authority: [authority_id] (with conflict)
  - Conflict type: [description]
  - Decision voided: [original decision text]
  - New authority assigned: [new_authority_id]
  - Remedial action: De novo review required
```

---

### Rule Set 3: Pattern Detection → Escalation for Quality Review

**Rule 3A: Adjudicator Consensus Rate Below Threshold**

```
MONITOR:
  - Adjudicator ID: [ADJ-NNNN]
  - Sample period: Last 30 days, minimum 10 cases
  - Consensus metric: % of cases where adjudicator agrees with 2nd adjudicator

TRIGGER:
  IF (adjudicator_consensus_rate < 85%) AND (case_count_in_period ≥ 10):
      ACTION_TRIGGERED = true

ACTION:
  1. Flag adjudicator for quality review
  2. Supervisor pulls adjudicator's last 10 cases
  3. Senior Examiner independently reviews all 10 cases
  4. Compare: Examiner findings vs. Adjudicator findings
  5. If disagreement rate >20%:
     - Adjudicator receives retraining (mandatory 40+ hours)
     - Adjudicator placed on probation (3 months)
     - Re-certification exam required before returning to T2 assignments
  6. If disagreement rate ≤20%:
     - Adjudicator receives coaching (no penalties)
     - Case review continues

NOTIFICATIONS:
  - To ADJUDICATOR: "Your case consensus rate has fallen below target. 
                     Quality review in progress. [Retraining/Coaching] required."
  - To SUPERVISOR: "Adjudicator [name] flagged for quality review. 
                    Consensus rate: [%]."
  - To SENIOR_EXAMINER: "Review and comparison required for adjudicator [name]"

AUDIT_LOG:
  - Event: PATTERN_DETECTION_ADJUDICATOR_CONSENSUS_BELOW_THRESHOLD
  - Adjudicator: [adjudicator_id]
  - Time: [timestamp]
  - Consensus rate: [%]
  - Sample size: [N cases]
  - Action recommended: [retraining / coaching]
```

**Rule 3B: Examiner Reversal Rate Above Threshold**

```
MONITOR:
  - Examiner ID: [EXAM-NNNN]
  - Sample period: Last 90 days, minimum 10 cases
  - Reversal metric: % of cases where examiner decision was overturned on appeal

TRIGGER:
  IF (examiner_reversal_rate > 15%) AND (case_count_in_period ≥ 10):
      ACTION_TRIGGERED = true

ACTION:
  1. Flag examiner for quality audit
  2. National Director reviews examiner's disputed cases
  3. Analyze: Examiner findings vs. appeals board findings
  4. If systematic error pattern detected:
     - Examiner receives mandatory retraining (40+ hours)
     - Probationary period: 3 months (all decisions double-checked by supervisor)
     - Re-certification exam required
     - Performance improvement plan (PIP)
  5. If isolated errors:
     - Coaching and mentoring (no penalties)
     - Continued monitoring

NOTIFICATIONS:
  - To EXAMINER: "Your case reversal rate has exceeded target. 
                  Quality audit in progress. [Retraining/Coaching] required."
  - To NATIONAL_DIRECTOR: "Senior Examiner [name] flagged. Reversal rate: [%]"
  - To SUPERVISORY_CHAIN: "Performance issue identified and remedial plan in progress"

AUDIT_LOG:
  - Event: PATTERN_DETECTION_EXAMINER_REVERSAL_ABOVE_THRESHOLD
  - Examiner: [examiner_id]
  - Time: [timestamp]
  - Reversal rate: [%]
  - Sample size: [N cases]
  - Action recommended: [retraining / coaching]
```

**Rule 3C: One Agent Involved in ≥3 Fraud Cases**

```
MONITOR:
  - Enrollment Agent ID: [AGT-NNNN]
  - Metric: Number of fraud cases involving this agent (as party or witness)
  - Threshold: ≥3 cases with confirmed fraud

TRIGGER:
  IF (fraud_case_count ≥ 3) AND (all ≥2 cases: fraud_confirmed):
      ACTION_TRIGGERED = true

ACTION:
  1. Escalate to DCPJ Fraud Unit
  2. DCPJ initiates insider threat investigation
  3. Possible scenarios:
     - Agent assisting fraudsters (corruption)
     - Agent conducting unauthorized re-enrollments
     - Agent running credential mill with criminal syndicate
  4. Remedial actions:
     - Agent immediately suspended pending investigation
     - Audit all agent's cases from last 12 months
     - Digital forensics: Agent's computer/phone access logs
     - Potential criminal charges if fraud confirmed
  5. If insider threat confirmed:
     - Termination of employment
     - Referral to Public Prosecutor for prosecution
     - Possible asset seizure

NOTIFICATIONS:
  - To AGENT: "Your employment is suspended pending investigation 
               into suspected fraud. You will be contacted by DCPJ."
  - To DCPJ_DIRECTOR: "Insider threat referral - Enrollment Agent [name] 
                       involved in ≥3 fraud cases. Investigation authorized."
  - To SNISID_DIRECTOR: "Critical: Insider threat investigation launched. 
                        Agent [name] suspended."
  - To SUPERVISORY_CHAIN: "Confidential: Insider threat investigation in progress"

AUDIT_LOG:
  - Event: PATTERN_DETECTION_INSIDER_THREAT_REFERRAL
  - Agent: [agent_id]
  - Time: [timestamp]
  - Fraud cases involved: [case_count]
  - Referral to DCPJ: [yes]
  - Investigation status: [initiated / in_progress / concluded]
```

---

### Rule Set 4: Citizen Complaint → Automatic Escalation

**Rule 4A: Citizen Files Formal Appeal**

```
TRIGGER:
  IF (citizen submits formal complaint via portal / mail / regional office):
      ACTION_TRIGGERED = true

COMPLAINT_CHANNELS:
  1. Online portal: Appeal form submitted + case ID
  2. Physical mail: Letter to regional director
  3. Phone call: Citizens can call regional ONI center to initiate appeal
  4. In-person: Walk into ONI center, file complaint with staff

ACTION:
  1. Create appeal record: [APPEAL-YYYY-NNNNNN]
  2. Link to original case: [CFR-YYYY-NNNNNN]
  3. Case automatically escalated to next authority level
  4. Case status: APPEALED
  5. Appeal grounds assessed: 
     - Is complaint substantive or frivolous?
     - Was original decision within authority? 
     - Is there new evidence?
  6. If substantive: Case review scheduled within 7 days
     If frivolous: Appeal denied with written explanation

ESCALATION_TARGET:
  - If appeal filed during T1-T2: Escalate to Level 3 (Regional Director)
  - If appeal filed during T3-T4: Escalate to Level 4 (NIAB / National Director)
  - If appeal filed during T5: Escalate to Level 4 (NIAB Appeal Board)
  - If appeal filed during T6: Escalate to appellate tribunal

NOTIFICATIONS:
  - To CITIZEN: "Your appeal has been received and assigned case ID [APPEAL-YYYY-NNNNNN]. 
                 Your appeal will be reviewed within [X days]. You will receive 
                 a decision notice by [date]."
  - To ORIGINAL_AUTHORITY: "Your decision in Case [CFR-YYYY-NNNNNN] has been 
                            appealed. You may submit a written rebuttal within 3 days."
  - To NEXT_AUTHORITY: "Appeal received for Case [CFR-YYYY-NNNNNN]. 
                        Complete de novo review required within [X days]."

APPEAL_DECISION_OPTIONS:
  1. APPEAL_SUSTAINED: Original decision overturned, new decision issued
  2. APPEAL_DENIED: Original decision upheld
  3. APPEAL_PARTIALLY_GRANTED: Original decision modified (hybrid outcome)

AUDIT_LOG:
  - Event: CITIZEN_APPEAL_FILED
  - Case: [case_id]
  - Appeal: [appeal_id]
  - Time: [timestamp]
  - Citizen: [citizen_nni]
  - Grounds: [brief reason]
  - Escalation target: [next_authority_id]
  - Substantive assessment: [substantive / frivolous]
```

**Rule 4B: Frivolous Complaint Detection**

```
FRIVOLOUS_CRITERIA:
  1. Same complaint filed ≥2 times previously (duplicate appeal)
  2. Complaint submitted after appeal deadline passed (untimely)
  3. Complaint contains no factual basis (adjudicator properly applied rules)
  4. Outcome remained unchanged after 2+ prior appeals (res judicata)
  5. Complaint explicitly states no basis ("I just don't like the decision")

IF frivolous_complaint_detected:
  ACTION:
    1. Deny appeal with written explanation
    2. Provide citation to decision rule that was applied
    3. Notify citizen of appeal deadline (if still within window)
    4. No escalation (stays with original authority)
    5. Case remains CLOSED

NOTIFICATION_TO_CITIZEN:
  "Your appeal has been reviewed and determined to be frivolous 
   [reason: duplicate / untimely / no factual basis / res judicata]. 
   Your original case decision remains in effect. 
   You have [X days] to file a further appeal if new evidence is available."
```

---

## Real-Time Monitoring Dashboard

### Dashboard 1: Executive SLA Status (National Director View)

```
SNISID CONFLICT RESOLUTION DASHBOARD
Updated: [real-time, every 5 minutes]
Report Date: [YYYY-MM-DD]

TIER SUMMARY (All Tiers Combined):
├─ Total Open Cases: 847
├─ Total Overdue (100%+ SLA): 12 cases [RED - CRITICAL]
├─ Total At-Risk (75-99% SLA): 34 cases [ORANGE - WARNING]
├─ Total Healthy (0-74% SLA): 801 cases [GREEN]
├─ SLA Compliance Rate: 98.6% (target: ≥95%)
└─ Mean Case Resolution Time: 42 hours [Target: Within tier SLA]

BY TIER:
┌─────────┬──────────┬─────────┬──────────┬────────────┐
│ Tier    │ Total    │ Overdue │ At-Risk  │ SLA %      │
├─────────┼──────────┼─────────┼──────────┼────────────┤
│ T1 (24h)│ 120      │ 0       │ 5        │ 99.6%  ✓   │
│ T2 (48h)│ 180      │ 1       │ 8        │ 99.4%  ✓   │
│ T3 (72h)│ 145      │ 2       │ 12       │ 98.6%  ✓   │
│ T4 (7d) │ 85       │ 5       │ 7        │ 94.1%  ⚠   │
│ T5 (30d)│ 215      │ 4       │ 2        │ 98.1%  ✓   │
│ T6 (90d)│ 102      │ 0       │ 0        │ 100%   ✓   │
└─────────┴──────────┴─────────┴──────────┴────────────┘

CRITICAL INCIDENTS:
├─ SLA Breaches (24h): 2 cases
│  ├─ CFR-2026-001847: T4 case, 2.5 days overdue [Regional Director responsible]
│  └─ CFR-2026-001850: T5 case, 1.2 days overdue [DCPJ Fraud Unit responsible]
├─ Conflict-of-Interest Detections: 0 (this week)
├─ Pattern-Based Escalations: 1
│  └─ Adjudicator ADJ-0142 (consensus rate 82%) [Retraining recommended]
└─ Citizen Complaints (open): 3 [All within SLA]

REGIONAL PERFORMANCE:
├─ Port-au-Prince: 98.9% SLA compliance
├─ Cap-Haïtien: 97.2% SLA compliance
├─ Santiago: 99.1% SLA compliance
└─ Gonaïves: 94.5% SLA compliance [Below target, manager meeting scheduled]

RECOMMENDED ACTIONS:
1. Escalate 2 overdue cases to National Director for personal review
2. Schedule retraining for adjudicator ADJ-0142
3. Address regional bottleneck in Gonaïves (insufficient staffing?)
```

### Dashboard 2: Authority-Level Performance (Regional Director View)

```
REGIONAL ONI CENTER: Port-au-Prince
Updated: [real-time, every 5 minutes]
Period: Last 7 days

CURRENT WORKLOAD:
├─ T1 Cases: 20 (healthy: 20, at-risk: 0, overdue: 0)
├─ T2 Cases: 25 (healthy: 24, at-risk: 1, overdue: 0)
├─ T3 Cases: 18 (healthy: 18, at-risk: 0, overdue: 0)
├─ T4 Cases: 12 (healthy: 11, at-risk: 1, overdue: 0)
├─ T5 Cases: 8 (healthy: 8, at-risk: 0, overdue: 0)
└─ Total: 83 cases

INDIVIDUAL AUTHORITY PERFORMANCE:
┌──────────────────────┬────────┬────────┬──────────┐
│ Authority            │ Tier   │ Cases  │ SLA %    │
├──────────────────────┼────────┼────────┼──────────┤
│ Enrollment Agent 01  │ T1     │ 15     │ 100%  ✓  │
│ Biometric Adj. 01    │ T2     │ 12     │ 98%   ✓  │
│ Biometric Adj. 02    │ T2     │ 13     │ 99%   ✓  │
│ Senior Examiner 01   │ T3     │ 18     │ 100%  ✓  │
│ Regional Director 01 │ T4/T5  │ 12     │ 98%   ✓  │
└──────────────────────┴────────┴────────┴──────────┘

CASE QUEUE (Oldest First):
1. CFR-2026-001923: T1, 18h elapsed (75% SLA), created by Agent-01
2. CFR-2026-001924: T2, 36h elapsed (75% SLA), adjudicators: ADJ-01, ADJ-02
3. CFR-2026-001925: T3, 48h elapsed (67% SLA), examiner: EXAM-01

UPCOMING DEADLINES (Next 24 Hours):
├─ 3 cases due within 12 hours [Prioritize]
├─ 7 cases due within 24 hours [Monitor]
└─ 2 cases require supervisor sign-off today
```

### Dashboard 3: Escalation Trigger History (Compliance Officer View)

```
ESCALATION TRIGGER LOG
Period: Last 7 days
Updated: [real-time]

TIME-BASED ESCALATIONS:
├─ 50% SLA Alerts Sent: 42
├─ 75% SLA Auto-Escalations: 8
├─ 100% SLA Critical Breaches: 2
│  ├─ CFR-2026-001847: T4 case, escalated from Regional Director to National Director
│  └─ CFR-2026-001850: T5 case, escalated from DCPJ to NIAB

CONFLICT-OF-INTEREST ESCALATIONS:
├─ Automated Conflict Detections: 3
│  ├─ Geographic conflict: 1 (agent home commune <50 km)
│  ├─ Family relationship: 1 (authority related to citizen)
│  └─ Prior employment: 1 (authority had prior case)
├─ Cases Reassigned: 3
└─ Post-Decision Rescissions: 0

PATTERN-BASED ESCALATIONS:
├─ Adjudicator Consensus Rate <85%: 1 [ADJ-0142]
├─ Examiner Reversal Rate >15%: 0
├─ Insider Threat Referrals: 0

CITIZEN COMPLAINTS & APPEALS:
├─ Formal Appeals Filed: 3
│  ├─ APPEAL-2026-001923: Substantive (escalated to next authority)
│  ├─ APPEAL-2026-001924: Frivolous (duplicate complaint, denied)
│  └─ APPEAL-2026-001925: Under review (decision pending)
├─ Appeals Sustained: 1 (28 days ago)
├─ Appeals Denied: 2 (last month)
└─ Appeal Overturn Rate: 25% (below 10% target)

TOTAL ESCALATIONS (This Week): 14
├─ SLA-based: 10
├─ Conflict-based: 3
├─ Pattern-based: 1
└─ Citizen-initiated: 0
```

---

## SLA Compliance Reporting

### Weekly SLA Compliance Report (for National Director)

```
SNISID WEEKLY SLA COMPLIANCE REPORT
Week Ending: 2026-05-23
Report Generated: 2026-05-24

EXECUTIVE SUMMARY:
SLA Compliance Rate: 98.6%
Cases Overdue: 12 of 847 (1.4%)
Target Compliance: ≥95%
Status: ON TRACK ✓

COMPLIANCE BY TIER:
┌─────────┬─────────────┬──────────────┬────────────┐
│ Tier    │ Cases Closed│ On-Time (%)  │ Status     │
├─────────┼─────────────┼──────────────┼────────────┤
│ T1 (24h)│ 85          │ 99.6%        │ EXCELLENT  │
│ T2 (48h)│ 62          │ 99.4%        │ EXCELLENT  │
│ T3 (72h)│ 48          │ 98.6%        │ GOOD       │
│ T4 (7d) │ 35          │ 94.1%        │ FAIR       │
│ T5 (30d)│ 28          │ 98.1%        │ GOOD       │
│ T6 (90d)│ 12          │ 100%         │ EXCELLENT  │
└─────────┴─────────────┴──────────────┴────────────┘

REGIONAL PERFORMANCE:
┌────────────────┬────────────┬───────────┐
│ Region         │ On-Time %  │ Status    │
├────────────────┼────────────┼───────────┤
│ Port-au-Prince │ 98.9%      │ GOOD      │
│ Cap-Haïtien    │ 97.2%      │ GOOD      │
│ Santiago       │ 99.1%      │ EXCELLENT │
│ Gonaïves       │ 94.5%      │ FAIR      │
│ Jérémie        │ 98.7%      │ GOOD      │
│ Les Cayes      │ 99.2%      │ EXCELLENT │
└────────────────┴────────────┴───────────┘

OVERDUE CASES (12 Total):
1. CFR-2026-001847 (T4): 2.5 days overdue [Regional Director, Port-au-Prince]
2. CFR-2026-001850 (T5): 1.2 days overdue [DCPJ Fraud Unit, National]
3. CFR-2026-001852 (T3): 0.8 days overdue [Senior Examiner, Santiago]
... [9 more cases with hours overdue]

ROOT CAUSES (Preliminary Analysis):
├─ Staffing shortage (Gonaïves region): 3 cases affected
├─ Complex case requiring additional evidence: 4 cases affected
├─ Waiting for external agency response (court, civil registry): 3 cases
├─ System outage (database): 1 case [2 hours - quickly resolved]
└─ Authority illness/leave: 1 case [coverage in progress]

RECOMMENDED ACTIONS:
1. Allocate temporary staff to Gonaïves (3-month assignment)
2. Fast-track 3 cases waiting on external responses
3. Review case complexity criteria for T4/T5 (possible SLA extension criteria)
4. Monitor database reliability (1 outage too many)
5. Implement mutual backup coverage for authority illness/leave

FORECAST (Next 30 Days):
├─ Projected compliance rate: 97.8% (slight decrease due to seasonal staffing)
├─ Anticipated bottleneck: T4 cases (longer resolution time)
└─ Planned interventions: Hire 2 temporary Senior Examiners (June-August)
```

### Monthly Governance Audit Report (for Inspector General)

```
SNISID MONTHLY GOVERNANCE AUDIT REPORT
Month: May 2026
Report Generated: 2026-06-01

AUDIT SCOPE:
├─ Cases reviewed: 20% random sample (169 of 847 open cases)
├─ Time period: May 1-31, 2026
├─ Auditors: Governance Team (3 auditors, 40 hours total)
└─ Findings: 2 minor issues, 0 major issues

CASE AUDIT RESULTS:

1. Decision Authority Compliance:
   - 169 cases reviewed for proper authority making decision
   - Cases with proper authority: 169/169 (100%)
   - Improper overrides detected: 0
   - Unauthorized decisions: 0
   Status: COMPLIANT ✓

2. Conflict-of-Interest Screening:
   - 169 cases reviewed for conflict-of-interest procedures
   - Automated conflict checks performed: 169/169 (100%)
   - Cases with undisclosed conflicts: 0
   - Post-decision conflict detections: 1 (resolved, decision rescinded)
   Status: COMPLIANT ✓ [Minor: 1 rescission required]

3. Documentation Quality:
   - Cases with proper decision rationales: 167/169 (98.8%)
   - Rationales meeting word count (≥50 words T1-2, ≥100 words T3+): 167/169
   - Two cases with insufficient documentation:
     ├─ CFR-2026-001823 (T1): Rationale only 42 words [CORRECTED]
     └─ CFR-2026-001824 (T2): Rationale only 87 words [CORRECTED]
   Status: COMPLIANT ✓ [Minor: 2 cases required correction]

4. Evidence Chain-of-Custody:
   - Cases with complete access logs: 169/169 (100%)
   - WORM storage integrity verified: 169/169 (100%)
   - Hash mismatches detected: 0
   - Evidence missing or corrupted: 0
   Status: COMPLIANT ✓

5. Citizen Notification Compliance:
   - Closed cases with proper citizen notification: 92/92 (100%)
   - Notification timeliness (≤24h for T1-4, ≤48h for T5): 92/92 (100%)
   - Appeal rights clearly stated: 92/92 (100%)
   - Appeal deadline provided: 92/92 (100%)
   Status: COMPLIANT ✓

AUTHORITY PERFORMANCE REVIEW:

Individual Adjudicator Consensus Rates (last 30 days):
├─ ADJ-0131: 92% (excellent)
├─ ADJ-0132: 89% (good)
├─ ADJ-0133: 94% (excellent)
├─ ADJ-0141: 87% (acceptable)
├─ ADJ-0142: 82% (below threshold) ← FLAGGED FOR RETRAINING
├─ ADJ-0143: 91% (excellent)
└─ Average across team: 89.3% (target: ≥90%) [Slight miss]

Senior Examiner Reversal Rates (last 90 days):
├─ EXAM-0151: 8% (excellent)
├─ EXAM-0152: 12% (good)
├─ EXAM-0153: 10% (excellent)
├─ EXAM-0161: 14% (at threshold, monitor)
└─ Average across team: 11% (target: ≤10%) [Slight miss]

PATTERN DETECTION FINDINGS:

1. One-Off Staffing Issue:
   - Adjudicator ADJ-0142 consensus rate: 82% (below 85% threshold)
   - Likely cause: Personal circumstances (confirmed with staff)
   - Remedial action: 40-hour retraining program scheduled (June 1-15)
   - Monitoring: Weekly check-ins with supervisor

2. Regional Bottleneck (Gonaïves):
   - SLA compliance: 94.5% (below 95% target)
   - Likely cause: Understaffing (2 senior examiners for region size)
   - Remedial action: Temporary staffing request approved
   - Timeline: 2 examiners allocated (June-August)

3. T4 Case Complexity:
   - Average T4 case resolution time: 5.2 days (vs. 7-day SLA)
   - Complexity drivers: Multi-record reconciliation, legal review needed
   - Trend: Stable (no degradation)
   - Recommendation: Monitor for future SLA extension needs

COMPLIANCE SUMMARY:
✓ Decision authority: 100% compliant
✓ Conflict-of-interest: 99.4% compliant (1 issue resolved)
✓ Documentation: 98.8% compliant (2 minor corrections)
✓ Evidence integrity: 100% compliant
✓ Citizen notification: 100% compliant

OVERALL GOVERNANCE HEALTH: GOOD (Minor issues resolved, no major concerns)

RECOMMENDED ACTIONS:
1. Complete ADJ-0142 retraining by June 15
2. Deploy temporary examiners to Gonaïves by June 5
3. Continue monitoring EXAM-0161 (reversal rate at threshold)
4. Document lessons learned from post-decision conflict rescission
5. Consider SLA adjustment policy for exceptionally complex cases
```

---

## Escalation Audit Trail

### Immutable Escalation Log (WORM Archive)

```
ESCALATION AUDIT TRAIL
Case: CFR-2026-001847
Tier: T4 (Multi-Record Conflict)
Created: 2026-05-15 09:30:00Z
Current Status: ESCALATED
Last Updated: 2026-05-23 14:15:00Z

TIMELINE:

Event 1: CASE_CREATED
├─ Timestamp: 2026-05-15T09:30:00Z
├─ Source: Automated conflict detection
├─ Status: OPEN
├─ SLA Deadline: 2026-05-22 09:30:00Z (7 days)
└─ Audit hash: sha256:4a8f...

Event 2: CASE_ASSIGNED
├─ Timestamp: 2026-05-15T10:15:00Z
├─ Assigned to: Regional Director 01 (Port-au-Prince)
├─ Status: ASSIGNED
├─ Conflict-of-interest cleared: YES
└─ Audit hash: sha256:5b2c...

Event 3: CASE_UNDER_REVIEW
├─ Timestamp: 2026-05-16T08:00:00Z
├─ Authority: Regional Director 01
├─ Status: UNDER_REVIEW
├─ Estimated completion: 2026-05-20
└─ Audit hash: sha256:6d4e...

Event 4: SLA_50_PERCENT_ALERT
├─ Timestamp: 2026-05-18T21:30:00Z
├─ SLA elapsed: 50.6%
├─ Supervisor notified: YES
├─ Message: "Case approaching SLA deadline. Status update needed."
└─ Audit hash: sha256:7e5f...

Event 5: SLA_75_PERCENT_ESCALATION
├─ Timestamp: 2026-05-21T15:45:00Z
├─ SLA elapsed: 76.2%
├─ Escalated from: Regional Director 01 → National Director
├─ Status: ESCALATED
├─ Escalation reason: SLA_BREACH_75_PERCENT
├─ New SLA deadline: 2026-05-23 15:45:00Z (48-hour extended SLA)
├─ Notifications sent: YES (National Director, citizen, supervisory chain)
└─ Audit hash: sha256:8f6g...

Event 6: CASE_STATUS_UPDATE
├─ Timestamp: 2026-05-23T09:00:00Z
├─ Authority: National Director
├─ Status: UNDER_REVIEW (continued by National Director)
├─ Progress note: "Reviewing case documentation. Decision expected by 2026-05-24."
└─ Audit hash: sha256:9g7h...

ESCALATION CHAIN INTEGRITY VERIFICATION:
├─ All hashes verified: ✓
├─ Chronological order verified: ✓
├─ No gaps or missing events: ✓
├─ No evidence of tampering: ✓
└─ WORM storage immutability: ✓

Chain-of-Custody Access Log:
├─ Regional Director 01: Accessed 2026-05-15-2026-05-21 (multiple times)
├─ National Director: Accessed 2026-05-21-2026-05-23 (ongoing)
├─ Governance Auditor: Accessed 2026-05-24 (audit sample review)
└─ No unauthorized access: ✓
```

---

## SLA Variance & Root-Cause Analysis

### Variance Analysis (Cases Exceeding SLA)

**Question: Why do some cases exceed SLA when auto-escalation should prevent breaches?**

**Answer: Intentional Pauses & External Dependencies**

```
SLA CLOCK PAUSE RULES:
1. Citizen Information Requested:
   - Agency pauses clock and awaits citizen response
   - Example: T1 duplicate requires citizen clarification
   - Clock paused until: Citizen provides response (max 5 days)
   - Reason: Cannot make decision without information

2. External Agency Coordination:
   - Clock paused while waiting for response from:
     • Civil registry (marriage certificate verification)
     • Court (judicial order issuance)
     • Other government agencies (inter-agency verification)
   - Clock paused until: Response received
   - Reason: Decision dependent on external input

3. Judicial Proceedings:
   - T6 cases: Clock paused during court proceedings
   - Example: Waiting for tribunal hearing date
   - Clock paused until: Court judgment issued
   - Reason: Court timeline overrides ONI SLA

4. Evidence Collection (T5):
   - Clock paused while DCPJ investigates T5 fraud case
   - Warrant execution, digital forensics, witness interviews
   - Clock paused until: Investigation complete
   - Reason: Quality investigation requires time

WHEN PAUSED TIME IS EXCLUDED:
- SLA clock only counts active authority work time
- Paused time (waiting for citizen/court/external response) does NOT count
- Cases may show "wall clock" delay but no SLA breach

EXAMPLE:
Case CFR-2026-001847 (T4, 7-day SLA):
├─ Created: 2026-05-15 (Monday 09:30)
├─ Assigned to Regional Director: 2026-05-15 (0.33 days elapsed)
├─ Clock paused: 2026-05-16 → 2026-05-20 (pending civil registry response)
│  [Paused time: 4 days - does NOT count against SLA]
├─ Clock resumed: 2026-05-20 (civil registry responded)
├─ Escalated (75% SLA): 2026-05-21 at 3.5 active days
│  [SLA calculation: 3.5 / 7 = 50% of 7-day SLA elapsed]
└─ Wall-clock time: 6 days; Active time: 3.5 days (within SLA)

Result: No SLA breach (active work time is 3.5 days < 7 days)
BUT: Case appears overdue on wall-clock timeline (6 days elapsed)
```

**Variance Analysis by Tier:**

```
Tier 1 (24h SLA):
├─ Mean resolution time: 18 hours (3 hours faster than SLA)
├─ Variance: ±2 hours (low variance, predictable)
├─ Cause of variance: Enrollment agent availability, duplicate certainty
└─ Action: None needed (performing above target)

Tier 2 (48h SLA):
├─ Mean resolution time: 40 hours (8 hours faster than SLA)
├─ Variance: ±4 hours (low variance)
├─ Cause of variance: Adjudicator consensus speed
└─ Action: None needed (performing above target)

Tier 3 (72h SLA):
├─ Mean resolution time: 60 hours (12 hours faster than SLA)
├─ Variance: ±8 hours (moderate variance)
├─ Cause of variance: Biometric complexity, quality issues, age progression analysis
└─ Action: Monitor for deterioration (currently healthy)

Tier 4 (7-day SLA):
├─ Mean resolution time: 5.2 days (1.8 days faster than SLA)
├─ Variance: ±1.5 days (moderate variance, increased by external dependencies)
├─ Cause of variance:
│  ├─ Civil registry response time (0.5-2 days)
│  ├─ Legal counsel availability (0-1 day)
│  └─ Case complexity (0-1 day)
└─ Action: Monitor civil registry SLA compliance; consider SLA extension criteria

Tier 5 (30-day SLA):
├─ Mean resolution time: 18 days (12 days faster than SLA)
├─ Variance: ±5 days (high variance due to investigation complexity)
├─ Cause of variance:
│  ├─ Warrant need (0-2 days)
│  ├─ Digital forensics (2-7 days)
│  ├─ Witness interviews (1-3 days)
│  └─ DCPJ staffing availability (0-2 days)
└─ Action: Monitor DCPJ resource allocation; current performance good

Tier 6 (90-day SLA):
├─ Mean resolution time: N/A (court timeline independent of ONI SLA)
├─ Variance: Court-dependent (±30 days or more)
├─ Cause: Tribunal scheduling, appeal processes, expert witness availability
└─ Action: None needed (court authority overrides SLA)
```

---

**Document Version:** 2.0
**Last Updated:** 2026-05-23
**Next Review:** 2026-08-23
