# SNISID National Threat Intelligence Platform (TIP)

## 1. Objective
To transform the SOC from a reactive posture to a proactive posture by collecting, analyzing, and operationalizing cyber threat intelligence.

## 2. Functional Capabilities

| Function | Support | Description |
| :--- | :---: | :--- |
| **IOC Management** | Yes | Centralized database of IPs, domains, hashes, and file paths. |
| **Threat Feeds** | Yes | Integration of OSINT, Commercial, and Government feeds. |
| **Campaign Tracking** | Yes | Grouping individual incidents into larger adversary campaigns. |
| **Adversary Profiling** | Yes | Mapping attacker TTPs (Tactics, Techniques, and Procedures) using MITRE ATT&CK. |
| **Risk Scoring** | Yes | Prioritizing threats based on relevance to SNISID assets. |

## 3. Intelligence Cycle
1. **Planning:** Identifying "Priority Intelligence Requirements" (PIRs).
2. **Collection:** Gathering raw data from feeds, honey-pots, and internal logs.
3. **Processing:** Normalizing data into STIX/TAXII formats.
4. **Analysis:** Turning raw data into actionable intelligence.
5. **Dissemination:** Pushing IOCs to SIEM, Firewalls, and EDR.

## 4. Feed Integration
- **OSINT:** AlienVault OTX, MISP communities, Twitter/X security researchers.
- **Commercial:** Mandiant, CrowdStrike, etc.
- **Internal:** IOCs discovered during DFIR investigations.

## 5. Operationalization
The TIP must automatically feed the following:
- **SIEM:** Triggering alerts when an IOC is seen in logs.
- **Firewall/WAF:** Automatically blocking high-confidence malicious IPs.
- **EDR:** Scanning endpoints for specific malicious hashes.

## 6. Tooling Recommendations
- **MISP (Malware Information Sharing Platform):** For sharing and storing IOCs.
- **OpenCTI:** For mapping adversary relationships and campaign tracking.
