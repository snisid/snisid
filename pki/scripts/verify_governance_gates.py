#!/usr/bin/env python3
"""
SNISID - Governance Transition Gates Auditor
Validates compliance with legislative checklists and international MoUs
before allowing progressive ArgoCD deployments for next phases.
"""

import sys
import json

def verify_governance_gates():
    print("=========================================================")
    print("      SNISID AUTOMATED GOVERNANCE GATE AUDITOR           ")
    print("=========================================================")
    
    # 1. Legislative and Agreement Statuses (simulating DB query outputs)
    legal_status = {
        "loi_donnees_biometriques": True,
        "cadre_legal_nni": True,
        "reglementation_xroad": True,
        "loi_cybersecurite": False,
        "accord_diaspora": False
    }
    
    mou_status = {
        "mou_estonia": True,
        "mou_singapore": True,
        "mou_caricom": True,
        "mou_interpol": False  # Pending signature, blocks Phase 3 transition
    }
    
    # 2. Evaluate Phase Gates
    print("[*] Evaluating Legislative & Partnership Transition Gates:")
    
    # Phase 2 Gate: Requires Biometric Data Law & NNI Legal Framework
    phase2_gate = legal_status["loi_donnees_biometriques"] and legal_status["cadre_legal_nni"]
    print(f"  - Gate Phase 2 [Loi Biométrie + Cadre NNI] : {'PASSED' if phase2_gate else 'FAILED'}")
    
    # Phase 3 Gate: Requires X-Road Regulation & all 4 International MoUs (Estonia, Singapore, CARICOM, Interpol)
    mou_all_signed = mou_status["mou_estonia"] and mou_status["mou_singapore"] and mou_status["mou_caricom"] and mou_status["mou_interpol"]
    phase3_gate = legal_status["reglementation_xroad"] and mou_all_signed
    print(f"  - Gate Phase 3 [Règlementation X-Road + 4 MoUs] : {'PASSED' if phase3_gate else 'FAILED'}")
    if not mou_all_signed:
        pending_mous = [k for k, v in mou_status.items() if not v]
        print(f"    [!] Pending MoUs: {pending_mous}")
        
    # Phase 4 Gate: Requires Cybersécurité Law
    phase4_gate = legal_status["loi_cybersecurite"]
    print(f"  - Gate Phase 4 [Loi Cybersécurité] : {'PASSED' if phase4_gate else 'FAILED'}")
    
    # Phase 5 Gate: Requires Diaspora international agreement
    phase5_gate = legal_status["accord_diaspora"]
    print(f"  - Gate Phase 5 [Accord Diaspora DR] : {'PASSED' if phase5_gate else 'FAILED'}")
    
    # 3. Overall Governance Decision
    next_action = "PROCEED TO PHASE 2 DEVELOPMENT. BLOCKED BEFORE PHASE 3 INTEGRATION (Pending Interpol MoU)."
    if phase2_gate and phase3_gate:
        next_action = "PROCEED TO PHASE 3 INTEGRATION."
    
    report = {
        "evaluation_timestamp": "2026-05-24T21:40:00Z",
        "gates": {
            "phase_2_core_snisid": {
                "status": "APPROVED" if phase2_gate else "REJECTED",
                "checks": {
                    "loi_donnees_biometriques": legal_status["loi_donnees_biometriques"],
                    "cadre_legal_nni": legal_status["cadre_legal_nni"]
                }
            },
            "phase_3_integration": {
                "status": "APPROVED" if phase3_gate else "REJECTED",
                "checks": {
                    "reglementation_xroad": legal_status["reglementation_xroad"],
                    "international_mous": mou_status
                }
            },
            "phase_4_soc_national": {
                "status": "APPROVED" if phase4_gate else "REJECTED",
                "checks": {
                    "loi_cybersecurite": legal_status["loi_cybersecurite"]
                }
            },
            "phase_5_national": {
                "status": "APPROVED" if phase5_gate else "REJECTED",
                "checks": {
                    "accord_diaspora": legal_status["accord_diaspora"]
                }
            }
        },
        "next_action": next_action
    }
    
    print("\n=========================================================")
    print(f"    GOVERNANCE STATE: {next_action}")
    print("=========================================================")
    
    # Write JSON report to disk
    report_path = "pki/scripts/governance_gates_report.json"
    with open(report_path, "w") as f:
        json.dump(report, f, indent=2)
    print(f"[+] Governance gates audit report written to {report_path}")
    
    return True

if __name__ == "__main__":
    verify_governance_gates()
    sys.exit(0)
