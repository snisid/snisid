# SNISID National XDR/EDR Strategy

## 1. Objective
To provide deep visibility and response capabilities at the endpoint and workload level across the national infrastructure.

## 2. Protection Scope

| Domain | Protection | Priority | Focus |
| :--- | :---: | :---: | :--- |
| **Kubernetes Nodes** | Yes | Critical | Container escapes, rogue pods, Kubelet compromise. |
| **Endpoints** | Yes | High | Phishing, credential theft, local privilege escalation. |
| **Servers** | Yes | Critical | Lateral movement, database exfiltration, web shell deployment. |
| **Mobile Enrollment Units** | Yes | Medium | Device compromise, unauthorized app installation. |
| **Edge Nodes** | Yes | High | DDoS sources, unauthorized access to local networks. |

## 3. Technology Stack

| Domain | Recommended Technology | Role |
| :--- | :--- | :--- |
| **EDR/XDR** | Wazuh / CrowdStrike / SentinelOne | Telemetry collection, detection, and response. |
| **Runtime Security** | Falco | Real-time Kubernetes runtime monitoring. |
| **Malware Detection** | YARA | Pattern matching for malicious files/memory. |

## 4. Key Capabilities
- **Process Monitoring:** Tracking every process created on a server or node.
- **File Integrity Monitoring (FIM):** Detecting unauthorized changes to critical system files.
- **Network Visibility:** Monitoring socket connections from each process.
- **Remote Response:** Ability to isolate a host from the network via the SOC console.

## 5. Deployment Model
- **Agent-Based:** Lightweight agents installed on all OS instances.
- **Sidecar Pattern:** Security containers running alongside application pods in K8s.
- **Central Management:** A single pane of glass for managing policies and viewing alerts.

## 6. Response Actions (Automated/Manual)
- **Isolate Host:** Disconnect the device from the network.
- **Kill Process:** Immediately stop a malicious binary.
- **Quarantine File:** Move a suspicious file to a secure area for analysis.
- **Snapshot Memory:** Dump RAM for DFIR analysis before rebooting.
