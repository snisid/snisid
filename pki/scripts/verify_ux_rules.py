#!/usr/bin/env python3
"""
SNISID - UI/UX Guidelines and Color Palette Compliance Auditor
Parses the centralized index.css file to verify that the official SNISID
color tokens, accessibility classes, and report watermarks are deployed.
"""

import sys
import json
import os

def verify_ux_rules():
    print("=========================================================")
    print("      SNISID AUTOMATED UI/UX RULES COMPLIANCE AUDITOR    ")
    print("=========================================================")
    
    css_path = "frontend/src/index.css"
    
    required_colors = {
        "#0d1b2a": "Background Principal",
        "#1565c0": "Accent Primaire",
        "#00bcd4": "Accent Secondaire",
        "#2e7d32": "Statut OK",
        "#e65100": "Alerte",
        "#c62828": "Critique",
        "#eceff1": "Texte Principal",
        "#1e3a5f": "Surface Cards"
    }
    
    report = {
        "evaluation_timestamp": "2026-05-24T21:46:00Z",
        "css_file_found": False,
        "colors_audited": {},
        "high_contrast_class_found": False,
        "watermark_class_found": False,
        "compliance_status": "FAILED"
    }
    
    if not os.path.isfile(css_path):
        print(f"[ERROR] CSS file not found at {css_path}!")
        return False
        
    report["css_file_found"] = True
    
    with open(css_path, "r", encoding="utf-8") as f:
        css_content = f.read().lower()
        
    print("[*] Auditing Color Palette variables in CSS...")
    all_colors_pass = True
    for hex_code, desc in required_colors.items():
        found = hex_code in css_content
        report["colors_audited"][hex_code] = {
            "description": desc,
            "found": found
        }
        status_str = "PASSED" if found else "FAILED"
        if not found:
            all_colors_pass = False
        print(f"  - Color {hex_code.upper()} ({desc}) : {status_str}")
        
    # Check classes
    high_contrast_found = ".high-contrast" in css_content
    report["high_contrast_class_found"] = high_contrast_found
    print(f"  - Class .high-contrast (Accessibility Override) : {'PASSED' if high_contrast_found else 'FAILED'}")
    
    watermark_found = ".report-watermark" in css_content
    report["watermark_class_found"] = watermark_found
    print(f"  - Class .report-watermark (Confidential Watermark) : {'PASSED' if watermark_found else 'FAILED'}")
    
    compliance_passed = all_colors_pass and high_contrast_found and watermark_found
    report["compliance_status"] = "PASSED" if compliance_passed else "FAILED"
    
    print("\n=========================================================")
    print(f"    UX COMPLIANCE STATUS : {report['compliance_status']}")
    print("=========================================================")
    
    # Write JSON report to disk
    report_path = "pki/scripts/ux_compliance_report.json"
    with open(report_path, "w") as f:
        json.dump(report, f, indent=2)
    print(f"[+] UX compliance audit report written to {report_path}")
    
    return compliance_passed

if __name__ == "__main__":
    success = verify_ux_rules()
    sys.exit(0 if success else 1)
