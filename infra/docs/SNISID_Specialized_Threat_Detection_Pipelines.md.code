# SNISID: Specialized Threat Detection Pipelines (171–180)

This document provides deep-dive architectures for specialized cyber threat detection modules, integrating real-time telemetry with AI-assisted analysis and automated SOC orchestration.

---

## 1. Phishing & Insider Threat Detection (171, 172)

### 1.1. Phishing Intelligence Pipeline
- **Architecture**: Distributed mail/message analyzers connected to the Kafka `telemetry.ingress.v1` topic.
- **AI Integration**: NLP-based sentiment analysis and URL reputation scoring.
- **Mitigation Hook**: Automatically quarantine suspicious messages and invalidate active sessions for users who interacted with high-risk links.

### 1.2. Insider Threat Detection (UEBA)
- **Architecture**: Flink-based behavioral baseline tracking.
- **Detection**: Identifies "Low and Slow" data access, unauthorized privilege escalation attempts, and anomalous out-of-hours activity.
- **Mitigation Hook**: Automatically downgrades the user's **ISTS Trust Score**, triggering mandatory Biometric MFA for all subsequent actions.

---

## 2. Lateral Movement & Network Scanning (173, 174)

### 2.1. Lateral Movement Detection
- **Architecture**: Service-mesh (Istio) telemetry analyzer.
- **Detection**: Identifies "Unusual Service Paths" (e.g., a Web-Frontend service attempting to call the KMS directly).
- **Mitigation Hook**: Dynamic mTLS certificate revocation via SPIRE for the source workload.

### 2.2. Network Scanning Detection
- **Architecture**: eBPF-based (Cilium) port-scan detectors.
- **Detection**: Rapid connection attempts to multiple ports/services from a single internal or external IP.
- **Mitigation Hook**: Instant IP-level blocking at the kernel level via Cilium NetworkPolicies.

---

## 3. Exploit & Ransomware Detection (175, 176)

### 3.1. Exploit Detection Module
- **Architecture**: Runtime syscall monitor (Falco).
- **Detection**: Execution of anomalous binaries, buffer overflow signatures, and kernel-level shell escapes.
- **Mitigation Hook**: Instant pod termination and forensic memory snapshotting.

### 3.2. Ransomware Kill-Switch
- **Architecture**: File-system entropy and deletion rate monitor.
- **Detection**: High-frequency encryption or deletion of files in Sovereign Object Storage.
- **Mitigation Hook**: Instant "Write-Lock" on the target storage bucket and administrative account suspension.

---

## 4. Exfiltration & Zero-Day Detection (177, 178)

### 4.1. Data Exfiltration Pipeline
- **Architecture**: Ingress/Egress traffic analyzer.
- **Detection**: Bandwidth spikes to non-sovereign IP ranges or unauthorized large-scale data transfers.
- **Mitigation Hook**: Automated egress throttling and session termination.

### 4.2. Zero-Day Anomaly Engine
- **Architecture**: Unsupervised AI clustering (Isolation Forests).
- **Detection**: Identification of traffic patterns that have no historical precedent in the SCES (Common Event Schema).
- **Mitigation Hook**: Promotion to "Level 1 Investigator" for immediate manual review and shadow-mode containment.

---

## 5. Threat Scoring & Classification (179, 180)

- **Threat Scoring Aggregation**: Flink aggregates scores from all the above pipelines to create a **Unified Threat Vector** for every identity/workload.
- **Attack Vector Classification**: AI-assisted mapping of incidents to the **MITRE ATT&CK** framework, providing analysts with instant context on the attack's stage (Recon, Initial Access, Lateral Movement, etc.).

---

## 6. Implementation Strategy
- **Deployment**: Each pipeline is a micro-service consuming from specific Kafka topics.
- **Traceability**: Every detection event includes a **Forensic ID** linked to the raw telemetry in the Sovereign Audit Ledger.
- **Observability**: Metrics (Latency, Throughput, Precision/Recall) are exported to the National SOC Grafana dashboard.
