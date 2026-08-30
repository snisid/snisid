# SNISID Senior Auditor: Risk & Maturity Scoring configurations

**Classification:** RESTRICTED / SOVEREIGN AUDIT
**Compliance:** ISO 27001 / ISO 22301 / SLA Governance

This playbook defines the dynamic maturity scoring scripts, automated escalation triggers, and GitOps integration gates used by the SNISID Senior Audit Directorate.

---

## 1. Dynamic Maturity Score Calculator

This Python script calculates the overall system maturity score by parsing the active audit states of critical domains and applying penalties for unresolved high/critical risks.

```python
# File: /opt/snisid/audit/calculate_maturity.py
import json
import sys

def compute_security_maturity(audit_results_path, output_report_path):
    print(f"[*] Parsing active audit results from {audit_results_path}...")
    
    with open(audit_results_path, 'r') as f:
        audit_data = json.load(f)
        
    # Define baseline domain weights (Total: 100%)
    weights = {
        "cybersecurity_zero_trust": 0.25,
        "biometric_abis": 0.20,
        "pki_hsm_governance": 0.15,
        "offline_dr": 0.15,
        "gitops_devsecops": 0.15,
        "data_governance": 0.10
    }
    
    # Extract domain scores (Default: 0.0 if missing)
    domain_scores = audit_data.get("domain_scores", {})
    weighted_sum = 0.0
    for domain, weight in weights.items():
        score = domain_scores.get(domain, 0.0)
        weighted_sum += score * weight
        
    print(f"[*] Computed Weighted Base Score: {weighted_sum:.2f}/100")
    
    # Parse active risk items and calculate penalties
    risk_states = audit_data.get("active_risks", {})
    penalties = 0.0
    
    # Critical Risks (Flat 15 point penalty)
    if risk_states.get("absent_pen_test", True):
        print("[!] Risk Penalty: Overdue Penetration Testing (-10 points)")
        penalties += 10.0
    if risk_states.get("root_ca_ceremony_unfinalized", True):
        print("[!] Risk Penalty: Root CA Ceremony Not Finalized (-15 points)")
        penalties += 15.0
        
    # High Risks (Flat 10 point penalty)
    if risk_states.get("soc_staffing_insufficient", True):
        print("[!] Risk Penalty: SOC Staffing Insufficient (-5 points)")
        penalties += 5.0
    if risk_states.get("starlink_sla_missing", True):
        print("[!] Risk Penalty: Starlink/VSAT SLA Missing (-5 points)")
        penalties += 5.0
    if risk_states.get("cap_haitien_budget_unallocated", True):
        print("[!] Risk Penalty: Cap-Haïtien budget unallocated (-10 points)")
        penalties += 10.0
        
    final_maturity_score = max(0.0, weighted_sum - penalties)
    print(f"[+] Final Recalculated Maturity Score: {final_maturity_score:.2f}/100")
    
    # Generate structured audit report
    report = {
        "calculated_timestamp": int(time_stamp_now()),
        "base_weighted_score": weighted_sum,
        "total_risk_penalties": penalties,
        "final_maturity_score": final_maturity_score,
        "gating_status": "BLOCKED" if final_maturity_score < 75.0 or risk_states.get("root_ca_ceremony_unfinalized") else "APPROVED"
    }
    
    with open(output_report_path, 'w') as f:
        json.dump(report, f, indent=2)
        
    print(f"[+] Report written to {output_report_path}. Gating Status: {report['gating_status']}")
    return report

def time_stamp_now():
    import time
    return time.time()

if __name__ == "__main__":
    # Test Mock data simulating active unresolved risks
    mock_audit = {
        "domain_scores": {
            "cybersecurity_zero_trust": 95.0,
            "biometric_abis": 90.0,
            "pki_hsm_governance": 88.0,
            "offline_dr": 92.0,
            "gitops_devsecops": 85.0,
            "data_governance": 80.0
        },
        "active_risks": {
            "absent_pen_test": True,                 # Penalty: -10
            "root_ca_ceremony_unfinalized": False,     # Penalty: 0 (Ceremony completed)
            "soc_staffing_insufficient": True,        # Penalty: -5
            "starlink_sla_missing": False,             # Starlink contract active
            "cap_haitien_budget_unallocated": True     # Penalty: -10
        }
    }
    
    import os
    with open("temp_audit_check.json", 'w') as f:
        json.dump(mock_audit, f)
        
    compute_security_maturity("temp_audit_check.json", "maturity_report.json")
    
    # Cleanup temp check file
    if os.path.exists("temp_audit_check.json"):
        os.remove("temp_audit_check.json")
```

---

## 2. IF/THEN/ESCALATE Decision Routing Engine

This script maps simulated audit warnings to specific escalation playbooks.

```python
# File: /opt/snisid/audit/escalate_risks.py
import sys

def evaluate_and_route_risk(risk_key, days_unresolved):
    """
    IF/THEN/ESCALATE risk routing logic.
    """
    print(f"[*] Evaluating risk incident: {risk_key} (Unresolved: {days_unresolved} days)")
    
    if risk_key == "root_ca_ceremony_unfinalized":
        if days_unresolved >= 1:
            print("[!] CRITICAL: Root CA procedure unfinalized. Gating active.")
            print("[*] ACTION: Halt all PKI certificate signers immediately.")
            print("[*] ESCALATION: Alert BRH Governor & Ministry of Justice.")
            return "HALT_PKI_OPERATIONS"
            
    elif risk_key == "absent_pen_test":
        if days_unresolved >= 30:
            print("[!] CRITICAL: Pen test overdue by 30 days.")
            print("[*] ACTION: Apply pipeline push blocks on main branch.")
            print("[*] ESCALATION: Page CISO & Request Red Team Allocation.")
            return "BLOCK_GITOPS_PROMOTION"
            
    elif risk_key == "soc_staffing_insufficient":
        if days_unresolved >= 60:
            print("[!] HIGH: SOC staffing below threshold.")
            print("[*] ACTION: Default all intrusion alerts to automatic isolation.")
            print("[*] ESCALATION: Route HR requisition query.")
            return "FORCE_AUTO_QUARANTINE"
            
    print("[+] Risk level within acceptable threshold. Continue monitoring.")
    return "CONTINUE_MONITORING"

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python escalate_risks.py <risk_key> <days_unresolved>")
        sys.exit(1)
    evaluate_and_route_risk(sys.argv[1], int(sys.argv[2]))
```

---

## 3. GitOps Roadmap Promotion Gate (GitLab CI Rule)

This YAML file integrates into the infrastructure repository pipeline, running the compliance calculator and blocking promotions to `Phase 2` if critical risks are unresolved.

```yaml
# File: /deployments/gitops/pipeline-gates.yaml
stages:
  - audit-check
  - deploy-promote

validate_maturity_gate:
  stage: audit-check
  image: python:3.10-slim
  script:
    - python /opt/snisid/audit/calculate_maturity.py /var/log/audit/active_status.json /tmp/report.json
    - |
      STATUS=$(python -c "import json; print(json.load(open('/tmp/report.json'))['gating_status'])")
      if [ "$STATUS" = "BLOCKED" ]; then
        echo "[!] ERROR: Maturity Score below compliance threshold or critical risks open."
        echo "[!] Promotion to next phase is BLOCKED."
        exit 1
      else
        echo "[+] Audit gate passed. Continuing promotion."
      fi
  artifacts:
    paths:
      - /tmp/report.json
```

---

## 4. Operational Risk Registry Schema

The active status files parsed by the auditing calculators follow this format:

```json
{
  "timestamp": 1779672000,
  "domain_scores": {
    "cybersecurity_zero_trust": 95.0,
    "biometric_abis": 90.0,
    "pki_hsm_governance": 88.0,
    "offline_dr": 92.0,
    "gitops_devsecops": 85.0,
    "data_governance": 80.0
  },
  "active_risks": {
    "absent_pen_test": true,
    "root_ca_ceremony_unfinalized": false,
    "soc_staffing_insufficient": true,
    "starlink_sla_missing": false,
    "cap_haitien_budget_unallocated": true
  }
}
```

---

*Verified and signed by the SNISID Senior Audit Directorate.*
