# SNISID National CERT Framework

## 1. Objective
The National Computer Emergency Response Team (CERT) is the operational arm for incident response and national coordination during cyber crises.

## 2. Core Functions

| Function | Support | Description |
| :--- | :---: | :--- |
| **Incident Response** | Yes | Lead the containment and eradication of active threats. |
| **Malware Analysis** | Yes | Reverse engineering of payloads found within SNISID. |
| **National Coordination** | Yes | Orchestrating response between government agencies and infrastructure providers. |
| **Crisis Communication** | Yes | Managing internal and external messaging during a breach. |
| **Threat Sharing** | Yes | Disseminating indicators of compromise (IOCs) to stakeholders. |

## 3. Incident Lifecycle Management (NIST based)
1. **Preparation:** Establishing tools, training, and communication channels.
2. **Detection & Analysis:** Validating alerts from the SOC.
3. **Containment, Eradication, & Recovery:** Stopping the bleed, removing the threat, restoring services.
4. **Post-Incident Activity:** Lessons learned and forensic reporting.

## 4. Coordination Capabilities
- **Inter-Agency Liaison:** Direct line to Ministry of Interior, Defense, and Telecommunications.
- **Private Sector Integration:** Collaboration with ISPs and Cloud providers.
- **International Cooperation:** Coordination with global CERTs (e.g., FIRST).

## 5. Malware Analysis Lab
- **Isolated Sandboxes:** Air-gapped environments for executing suspicious binaries.
- **Static Analysis Tools:** Ghidra, IDA Pro.
- **Dynamic Analysis Tools:** Cuckoo Sandbox, Any.run.

## 6. Crisis Communication Protocol
- **Triage Level 1 (Low):** Internal CERT notification.
- **Triage Level 2 (Medium):** SOC and IT Management notification.
- **Triage Level 3 (High):** National Cyber Governance and Executive notification.
- **Triage Level 4 (Critical):** Full Government notification and public disclosure (if necessary).
