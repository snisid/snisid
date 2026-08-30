#!/usr/bin/env python3
"""
Génère un rapport Markdown à partir des artefacts Governance-as-Code.
Sans dépendance externe.
"""
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
GOV = ROOT / "Governance"
BUILD = ROOT / "build"
BUILD.mkdir(exist_ok=True)

def j(rel):
    with open(GOV / rel, 'r', encoding='utf-8') as f:
        return json.load(f)

index = j("governance.index.json")
charter = j("National-Charters/national-snisid-governance-charter.json")
authority = j("Authority-SNISID/authority.structure.json")
inter = j("Inter-Agencies/interagency.governance.json")
raci = j("RACI/national.raci.json")
legal = j("Legal/legal.governance.json")
workflow = j("Workflow-Governance/workflow.governance.json")
cyber = j("Cybersecurity/cyber.governance.json")
pki = j("PKI/pki.governance.json")
ops = j("Operations/operations.24x7.json")
citizen = j("Citizen-Governance/citizen.governance.json")
audit = j("Audit/audit.risk.governance.json")
standards = j("Standards/national.standards.json")
compliance = j("Compliance/compliance.controls.json")

report_content = f"""# PHASE 1 — Gouvernance Nationale SNISID — Rapport d'implémentation
**Objectif** : {index['objective']}

**Règle absolue** : {index['absolute_rule']}

## 1. Conseil National SNISID
| Organe | Fonction | Autorité | Veto |
|--------|----------|----------|------|
"""
for body in charter['governing_bodies']:
    report_content += f"| {body['name']} | {body['function']} | {body['authority']} | {'Oui' if body['veto_right'] else 'Non'} |\n"

report_content += f"""
## 2. Autorité Nationale SNISID 24/7
| Département | Fonction | Disponibilité |
|-------------|----------|---------------|
"""
for dept in authority['departments']:
    report_content += f"| {dept['name']} | {dept['function']} | {dept['availability']} |\n"

report_content += f"""
## 3. Gouvernance inter-agences
| Agence | Rôle | APIs | SLA critique |
|--------|------|------|--------------|
"""
for agency in inter['agencies']:
    apis_str = ", ".join(agency['apis'])
    report_content += f"| {agency['agency']} | {agency['role']} | {apis_str} | {agency['sla']['critical_response_minutes']} min |\n"

report_content += f"""
## 4. RACI nationale
| Domaine | Criticité | Responsible | Accountable | Consulted | Informed |
|---------|-----------|-------------|-------------|-----------|----------|
"""
for d in raci['domains']:
    r = ", ".join(d['Responsible']) if isinstance(d['Responsible'], list) else d['Responsible']
    a = ", ".join(d['Accountable']) if isinstance(d['Accountable'], list) else d['Accountable']
    c = ", ".join(d['Consulted']) if isinstance(d['Consulted'], list) else d['Consulted']
    i = ", ".join(d['Informed']) if isinstance(d['Informed'], list) else d['Informed']
    report_content += f"| {d['domain']} | {d['criticality']} | {r} | {a} | {c} | {i} |\n"

report_content += """
## 5. Gouvernance juridique
"""
for law in legal['laws_to_create']:
    must_define = ", ".join(law['must_define'])
    report_content += f"- **{law['name']}** — {law['objective']} ; définit : {must_define}\n"

report_content += f"""
## 6. Gouvernance workflows
**Règle** : {workflow['rule']}

**Approvals obligatoires** : {", ".join(workflow['production_gate']['required_approvals'])}

## 7. Gouvernance cybersécurité
| Organe | Fonction | Disponibilité |
|--------|----------|---------------|
"""
for s in cyber['structure']:
    report_content += f"| {s['body']} | {s['function']} | {s['availability']} |\n"

report_content += f"""
## 8. Gouvernance PKI
**Règle absolue** : {pki['absolute_rule']}

**Hiérarchie** : {" → ".join(pki['hierarchy'])}

## 9. Gouvernance opérationnelle 24/7
| Élément | Description | Owner | Disponibilité |
|---------|-------------|-------|---------------|
"""
for e in ops['elements']:
    report_content += f"| {e['name']} | {e['description']} | {e['owner']} | {e['availability']} |\n"

report_content += """
## 10. Gouvernance citoyenne
"""
for right in citizen['citizen_rights']:
    report_content += f"- {right}\n"

report_content += """
## 11. Audit & risque
"""
for control in audit['controls']:
    report_content += f"- {control}\n"

report_content += """
## 12. Standards nationaux
| Standard | Domaine | Owner |
|----------|---------|-------|
"""
for s in standards['standards']:
    report_content += f"| {s['standard']} | {s['domain']} | {s['owner']} |\n"

report_content += """
## 13. Contrôles de conformité
| ID | Contrôle | Domaine | Sévérité |
|----|----------|---------|----------|
"""
for c in compliance['controls']:
    report_content += f"| {c['id']} | {c['name']} | {c['domain']} | {c['severity']} |\n"

report_content += """
## Résultat final Phase 1
✅ Gouvernance nationale
✅ Gouvernance juridique
✅ Gouvernance cyber
✅ Gouvernance PKI
✅ Gouvernance workflows
✅ Gouvernance opérationnelle
✅ Gouvernance citoyenne
✅ Gouvernance audit
✅ Standards nationaux

SNISID est maintenant représenté comme une institution numérique gouvernable, pas seulement comme architecture ou logiciel.
"""

output_file = BUILD / "PHASE1_GOVERNANCE_REPORT.md"
output_file.write_text(report_content, encoding='utf-8')
print(f"✅ Rapport généré: build/PHASE1_GOVERNANCE_REPORT.md")
