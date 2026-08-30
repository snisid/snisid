#!/usr/bin/env python3
import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
GOV = ROOT / "Governance"

def check_file(path_str):
    p = GOV / path_str
    if not p.exists():
        print(f"❌ Missing file: {path_str}")
        return False
    return True

def j(rel):
    with open(GOV / rel, 'r', encoding='utf-8') as f:
        return json.load(f)

success = True

# 1. tous les fichiers obligatoires existent
files_to_check = [
    "governance.index.json",
    "National-Charters/national-snisid-governance-charter.json",
    "Authority-SNISID/authority.structure.json",
    "Inter-Agencies/interagency.governance.json",
    "RACI/national.raci.json",
    "Legal/legal.governance.json",
    "Workflow-Governance/workflow.governance.json",
    "Cybersecurity/cyber.governance.json",
    "PKI/pki.governance.json",
    "Operations/operations.24x7.json",
    "Citizen-Governance/citizen.governance.json",
    "Audit/audit.risk.governance.json",
    "Standards/national.standards.json",
    "Compliance/compliance.controls.json"
]

for f in files_to_check:
    if not check_file(f):
        success = False

if not success:
    sys.exit(1)

# Load jsons for checks
raci = j("RACI/national.raci.json")
inter = j("Inter-Agencies/interagency.governance.json")
legal = j("Legal/legal.governance.json")
pki = j("PKI/pki.governance.json")
standards = j("Standards/national.standards.json")

# 2. tous les domaines RACI critiques sont présents
critical_domains = ["Identity", "Biometrics", "PKI", "Cybersecurity", "Workflow", "Elections", "Judicial"]
raci_domains = [d["domain"] for d in raci["domains"]]
for cd in critical_domains:
    if cd not in raci_domains:
        print(f"❌ Missing critical RACI domain: {cd}")
        success = False

# 3. toutes les agences sont intégrées
required_agencies = ["ONI", "ANH", "DGI", "DCPJ", "DGIE", "CEP", "Justice", "Santé", "Éducation"]
present_agencies = [a["agency"] for a in inter["agencies"]]
for ra in required_agencies:
    if ra not in present_agencies:
        print(f"❌ Missing required agency: {ra}")
        success = False

# 4. toutes les lois nécessaires sont déclarées
required_laws = ["Digital Identity Law", "Biometric Governance Law", "PKI Law", "Data Protection Law", "Interoperability Law", "Cybersecurity Law", "Audit Law", "Consent Law"]
present_laws = [l["name"] for l in legal["laws_to_create"]]
for rl in required_laws:
    if rl not in present_laws:
        print(f"❌ Missing required law: {rl}")
        success = False

# 5. la règle PKI Root CA offline est présente
if pki.get("absolute_rule") != "Le Root CA doit être offline.":
    print("❌ PKI absolute rule 'Root CA offline' is missing or incorrect")
    success = False

# 6. tous les standards nationaux sont définis
required_standards = ["API Standards", "Kafka Standards", "IAM Standards", "BPMN Standards", "Security Standards", "Kubernetes Standards", "Observability Standards"]
present_standards = [s["standard"] for s in standards["standards"]]
for rs in required_standards:
    if rs not in present_standards:
        print(f"❌ Missing national standard: {rs}")
        success = False

if success:
    print("✅ SNISID Phase 1 governance validation passed.")
    sys.exit(0)
else:
    sys.exit(1)
