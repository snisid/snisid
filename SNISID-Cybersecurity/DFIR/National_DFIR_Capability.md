# SNISID National DFIR Capability

## 1. Objective
To establish a rigorous national capability for Digital Forensics and Incident Response, ensuring that cyber attacks can be reconstructed and evidence preserved for legal proceedings.

## 2. Core Capabilities

| Function | Support | Description |
| :--- | :---: | :--- |
| **Disk Forensics** | Yes | Imaging and analyzing physical/virtual disks. |
| **Memory Analysis** | Yes | Analyzing RAM dumps to find fileless malware and active connections. |
| **Network Forensics** | Yes | Analyzing PCAP files and Netflow to track data exfiltration. |
| **Timeline Reconstruction** | Yes | Creating a master timeline of events across multiple systems. |
| **Evidence Preservation** | Yes | Maintaining chain of custody and using write-blockers. |

## 3. The Forensics Process
1. **Identification:** Identifying affected systems and volatile data.
2. **Preservation:** Creating bit-for-bit images of disks and memory dumps.
3. **Analysis:** Searching for artifacts (Registry, Event Logs, Prefetch, MFT).
4. **Documentation:** Creating a detailed report of findings.
5. **Presentation:** Providing expert testimony or reports for legal action.

## 4. DFIR Tooling Stack

| Category | Tools |
| :--- | :--- |
| **Imaging** | FTK Imager, dd, Guymager |
| **Memory Analysis** | Volatility 3, Rekall |
| **Disk Analysis** | Autopsy, Sleuth Kit, Axiom |
| **Network Analysis** | Wireshark, Zeek, Brim |
| **Timeline** | Plaso (log2timeline) |

## 5. Legal Admissibility Requirements
To ensure evidence is legally admissible:
- **Chain of Custody:** Detailed log of who handled the evidence, when, and why.
- **Hashing:** Using SHA-256/512 to prove that the image has not been altered.
- **Write Blockers:** Using hardware write-blockers during acquisition.
- **Standardized Procedures:** Following internationally recognized forensic standards (ISO/IEC 27037).

## 6. Integration with SOC
The SOC triggers the DFIR process when an incident is escalated to Tier 3 or when legal action is anticipated.
