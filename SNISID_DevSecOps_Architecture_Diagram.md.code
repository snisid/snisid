# SNISID Sovereign DevSecOps Architecture

Below is the complete enterprise-grade Mermaid diagram representing the DevSecOps and Secure Software Supply Chain architecture for SNISID. 

It integrates GitOps pull-based deployments, Policy-as-Code admission controllers, immutable infrastructure provisioning, and automated security scanning (SAST/DAST/SBOM) at every phase of the pipeline.

```mermaid
graph TD
    %% Styling and Themes
    classDef dev fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef ci fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef cd fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef sec fill:#ffebee,stroke:#c62828,stroke-width:2px;
    classDef k8s fill:#f3e5f5,stroke:#4a148c,stroke-width:2px;
    classDef iac fill:#e1bee7,stroke:#6a1b9a,stroke-width:2px;

    %% 1. Developer & IaC
    subgraph L1 [Developer & Infrastructure as Code]
        DEV[Developer]:::dev
        TF[Terraform IaC Pipelines]:::iac
        GIT_CODE[Source Code Repository]:::dev
        GIT_OPS[GitOps Manifest Repository <br/> Helm/Kustomize]:::dev
        
        DEV -->|Commits Application Code| GIT_CODE
        TF -->|Provisions Immutable Clusters| K8S_API
        TF -->|Configures Vault Auth Methods| VAULT
    end

    %% 2. Continuous Integration & Security (CI)
    subgraph L2 [Continuous Integration & Secure Supply Chain]
        CI[CI Pipeline / Runner]:::ci
        SAST[SAST & Secret Scanning <br/> SonarQube / Gitleaks]:::sec
        SBOM[SBOM Generation <br/> Syft / Trivy]:::sec
        BUILD[OCI Container Build <br/> Kaniko / Buildah]:::ci
        SIGN[Image Signing <br/> Sigstore / Cosign]:::sec
        HARBOR[Harbor OCI Registry <br/> Vulnerability Scanning]:::ci
        
        GIT_CODE -->|Trigger Build| CI
        CI --> SAST
        SAST --> SBOM
        SBOM --> BUILD
        BUILD --> SIGN
        SIGN -->|Push Signed Image & SBOM| HARBOR
        CI -->|Automated PR: Update Image Tag| GIT_OPS
    end

    %% 3. Continuous Delivery & GitOps (CD)
    subgraph L3 [Continuous Delivery & Policy-as-Code]
        ARGO[ArgoCD Controller <br/> Pull-Based GitOps]:::cd
        VAULT[HashiCorp Vault <br/> Secret Provider Class]:::sec
        KYV[Kyverno Admission Controller <br/> Policy-as-Code]:::sec
        
        GIT_OPS -.->|Webhooks / Polls State| ARGO
        ARGO -->|Applies Desired State| K8S_API
    end

    %% 4. Kubernetes Production Environment
    subgraph L4 [Kubernetes Production Environment]
        K8S_API[Kubernetes API Server]:::k8s
        APP[SNISID Microservice Pod]:::k8s
        DAST[Dynamic Application <br/> Security Testing OWASP ZAP]:::sec
        
        K8S_API -->|Validates Creation| KYV
        KYV -.->|Verifies Cosign Image Signature| HARBOR
        KYV -->|Allows Pod Admission| APP
        VAULT -.->|Injects Dynamic Secrets via CSI| APP
        DAST -.->|Scans Running APIs| APP
    end

    %% 5. Compliance Automation
    subgraph L5 [Compliance & Audit Automation]
        COMP[Compliance & Audit Dashboard]:::sec
        SBOM -.->|Feeds VEX Vulnerability Data| COMP
        KYV -.->|Audit Mode Policy Violations| COMP
    end
```
