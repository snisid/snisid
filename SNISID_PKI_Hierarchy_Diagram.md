# SNISID National PKI Hierarchy Architecture

Below is the complete, enterprise-grade Mermaid diagram representing the multi-tiered Public Key Infrastructure (PKI) of the SNISID platform. 

It highlights the integration of Hardware Security Modules (HSMs), revocation validation (OCSP/CRL), end-entity citizens, and fully automated Kubernetes mTLS.

```mermaid
graph TD
    %% Styling and Themes
    classDef offline fill:#e8eaf6,stroke:#3f51b5,stroke-width:2px;
    classDef inter fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef issuing fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef hsm fill:#ffebee,stroke:#c62828,stroke-width:2px;
    classDef valid fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef entity fill:#fce4ec,stroke:#ad1457,stroke-width:2px;
    classDef k8s fill:#f3e5f5,stroke:#6a1b9a,stroke-width:2px;

    %% 1. Tier 1: Offline Root
    subgraph L1 [Tier 1: Sovereign Offline Root Authority]
        ROOT[National Root CA <br/> Air-Gapped Faraday Vault]:::offline
        HSM1[FIPS 140-3 Level 4 HSM <br/> M-of-N Dual Control]:::hsm
        ROOT <-->|Generates & Secures Core Keys| HSM1
    end

    %% 2. Tier 2: Intermediate CAs
    subgraph L2 [Tier 2: Intermediate Policy CAs]
        INT_CIT[Citizen Identity Policy CA]:::inter
        INT_GOV[Government & Infrastructure Policy CA]:::inter
        
        ROOT ==>|Offline Key Ceremony Signing| INT_CIT
        ROOT ==>|Offline Key Ceremony Signing| INT_GOV
    end

    %% 3. Tier 3: Online Issuing CAs
    subgraph L3 [Tier 3: Online Issuing CAs]
        ISS_EID[Citizen eID Issuing CA <br/> High Availability]:::issuing
        ISS_GOV[Government Device Issuing CA]:::issuing
        ISS_TLS[Infrastructure Sub-CA <br/> HashiCorp Vault]:::issuing
        HSM2[FIPS 140-3 Level 3 Network HSMs]:::hsm
        
        INT_CIT ==>|Signs| ISS_EID
        INT_GOV ==>|Signs| ISS_GOV
        INT_GOV ==>|Signs| ISS_TLS
        
        ISS_EID <-->|Protects Online Signing Keys| HSM2
        ISS_GOV <-->|Protects Online Signing Keys| HSM2
    end

    %% 4. Certificate Validation Services
    subgraph L4 [Revocation & Validation Services]
        OCSP[OCSP Responders <br/> Distributed Anycast Nodes]:::valid
        CRL[CRL Distribution Points <br/> Offline Bloom Filter Cache]:::valid
        
        ISS_EID -.->|Publishes Revocation Status| OCSP
        ISS_EID -.->|Publishes Delta Updates| CRL
        ISS_TLS -.->|Publishes Status| OCSP
    end

    %% 5. End Entities (Citizens & Gov)
    subgraph L5 [End Entities]
        CIT_CERT[Citizen eID Smart Card <br/> CC EAL6+ Secure Chip]:::entity
        GOV_CERT[Civil Servant / Admin <br/> FIDO2 Hardware Token]:::entity
        
        ISS_EID -->|Issues 10yr X.509 Cert| CIT_CERT
        ISS_GOV -->|Issues 3yr Client Cert| GOV_CERT
        
        CIT_CERT -.->|Live Verification| OCSP
    end

    %% 6. Kubernetes Automation & mTLS
    subgraph L6 [Kubernetes Cert Lifecycle Automation]
        CM[cert-manager Controller]:::k8s
        POD[SNISID Microservice Pod]:::k8s
        ISTIO[Istio Envoy Proxy <br/> Service Mesh]:::k8s
        
        ISS_TLS -->|Automated ACME/API Signing| CM
        CM -->|Automatically Mounts 24h TLS Secret| POD
        POD <-->|Strict Zero Trust mTLS Tunnel| ISTIO
    end
```
