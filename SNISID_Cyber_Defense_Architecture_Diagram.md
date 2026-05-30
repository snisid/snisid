# SNISID Sovereign National Cyber Defense Architecture

Below is the complete enterprise-grade Mermaid diagram representing the unified Cyber Defense and Security Operations Center (SOC) architecture for SNISID. 

It integrates Zero Trust perimeters, Kubernetes container security, AI-driven UEBA (User and Entity Behavior Analytics) for insider threats, and automated SOAR response mechanisms to ensure national resilience.

```mermaid
graph TD
    %% Styling and Themes
    classDef edge fill:#e1bee7,stroke:#6a1b9a,stroke-width:2px;
    classDef k8s fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef zt fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef ai fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef soc fill:#ffebee,stroke:#c62828,stroke-width:2px;
    classDef res fill:#f3e5f5,stroke:#4a148c,stroke-width:2px;

    %% 1. National Boundary & API Security
    subgraph L1 [National Boundary & API Security]
        WAF[Web Application Firewall <br/> DDoS & L7 Mitigation]:::edge
        AGW[Kong API Gateway <br/> Rate Limiting & JWT Enforcement]:::edge
        WAF --> AGW
    end

    %% 2. Kubernetes Security & Workloads
    subgraph L2 [Kubernetes Security & Workloads]
        KYV[Kyverno Admission <br/> Image Signature Verification]:::k8s
        FALCO[Falco eBPF <br/> Container Runtime Security]:::k8s
        ISTIO[Istio Service Mesh <br/> Strict mTLS]:::k8s
        APP[SNISID Microservices]:::k8s
        
        KYV -->|Verifies at Deploy| APP
        ISTIO <-->|Encrypts Pod-to-Pod| APP
        APP -.->|Inspects Syscalls| FALCO
        AGW -->|Routes to Mesh| ISTIO
    end

    %% 3. Zero Trust & Cryptography
    subgraph L3 [Zero Trust Identity & PKI]
        PKI[National PKI / HSM]:::zt
        VAULT[HashiCorp Vault]:::zt
        OPA[Open Policy Agent]:::zt
        IAM[Keycloak IDP]:::zt
        
        PKI -.->|Roots of Trust| VAULT
        VAULT -.->|Injects Dynamic Secrets| APP
        OPA -.->|Authorizes API Calls| AGW
        IAM -.->|Authenticates Users| AGW
    end

    %% 4. AI Threat Detection & UEBA
    subgraph L4 [AI Threat & Insider Detection]
        UEBA[UEBA Engine <br/> Insider Threat Detection]:::ai
        AI_NET[AI Network Anomaly Engine]:::ai
        
        IAM -.->|Login Patterns / Context| UEBA
        ISTIO -.->|Traffic Flows / Latency| AI_NET
        APP -.->|DB Query Volumes| UEBA
    end

    %% 5. National SOC (SIEM & SOAR)
    subgraph L5 [National SOC Core]
        SIEM[Wazuh / Splunk SIEM <br/> Central Log Aggregation]:::soc
        SOAR[Cortex SOAR <br/> Automated Playbooks]:::soc
        
        FALCO -->|Runtime Alerts| SIEM
        WAF -->|Block Events| SIEM
        UEBA -->|Risk Scores > 90| SIEM
        AI_NET -->|Network Anomalies| SIEM
        SIEM -->|Triggers Playbook| SOAR
    end

    %% 6. National Resilience & Automated Response
    subgraph L6 [National Resilience & Automated Response]
        K_ISO[Kubernetes Network Isolation]:::res
        BGP[BGP Route Blackholing]:::res
        REV[PKI Certificate Revocation]:::res
        
        SOAR -->|Quarantine Compromised Pod| K_ISO
        SOAR -->|Block Attacker IP at Edge| BGP
        SOAR -->|Revoke Compromised Cert| REV
        
        K_ISO -.->|Applies Cilium Network Policy| APP
        REV -.->|Updates CRL/OCSP| PKI
    end
```
