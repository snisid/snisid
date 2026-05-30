# SNISID National Platform Engineering Framework

## Version 1.0.0
## Date: 2026-05-25
## Classification: SOVEREIGN NATIONAL INFRASTRUCTURE

---

## 1. FRAMEWORK OVERVIEW

The SNISID National Platform Engineering Framework defines the technical standards, governance models, and operational procedures for building and maintaining Haiti's sovereign digital identity platform.

### 1.1 Core Principles

| Principle | Description |
|-----------|-------------|
| **Sovereignty First** | All infrastructure must be nationally controlled and auditable |
| **Zero Trust** | No implicit trust; verify every request, every time |
| **Cloud-Native** | Container-first, declarative, automated infrastructure |
| **Offline-Ready** | Must function in disconnected/low-connectivity environments |
| **Security by Design** | Security integrated at every layer, not bolted on |
| **GitOps Native** | All infrastructure declarative, versioned, and auditable |
| **Observability First** | Full visibility into system health, performance, and security |

### 1.2 Scope

This framework covers:
- Kubernetes Platform Architecture
- API Gateway Standards
- Identity & Access Management
- Public Key Infrastructure
- Event-Driven Architecture
- Offline Synchronization
- GitOps & CI/CD Pipelines
- Observability Stack
- Database Architecture
- Security Baselines
- Testing & Validation

---

## 2. KUBERNETES PLATFORM STANDARDS

### 2.1 Distribution Requirements

| Requirement | Standard |
|-------------|----------|
| Distribution | RKE2 or Talos Linux |
| Kubernetes Version | >= 1.29 |
| Container Runtime | containerd >= 1.7 |
| CNI | Cilium or Calico |
| CSI | Rook-Ceph or Longhorn |

### 2.2 Cluster Architecture

- **Multi-region deployment** with at least 3 control plane nodes
- **HA control plane** with etcd quorum (3 or 5 nodes)
- **Worker node autoscaling** via Karpenter or Cluster Autoscaler
- **Service mesh** (Istio or Linkerd) for mTLS and traffic management
- **Network policies** enforcing default-deny posture

### 2.3 Namespace Strategy

```
snisid-system/          - Platform core services
snisid-identity/        - Identity services
snisid-api-gateway/     - API gateway and routing
snisid-event-bus/       - Kafka/Redpanda cluster
snisid-pki/             - Certificate management
snisid-offline/         - Offline sync engine
snisid-databases/       - Database services
snisid-observability/   - Monitoring stack
snisid-security/        - Security tools
snisid-applications/    - Future application workloads
```

---

## 3. GITOPS STANDARDS

### 3.1 Tooling

| Component | Technology |
|-----------|------------|
| GitOps Engine | ArgoCD |
| SCM | GitLab (self-hosted) |
| Secret Management | Sealed Secrets / Vault |
| Policy Engine | OPA/Gatekeeper |

### 3.2 Repository Structure

```
snisid-gitops/
├── clusters/
│   ├── production/
│   ├── staging/
│   └── development/
├── applications/
│   ├── core-platform/
│   ├── identity/
│   ├── api-gateway/
│   └── ...
├── infrastructure/
│   ├── terraform/
│   └── helm-charts/
└── policies/
    ├── network/
    └── security/
```

### 3.3 Deployment Flow

1. Developer commits to feature branch
2. CI pipeline runs tests, security scans, builds artifacts
3. PR merged to main triggers image build and Helm chart update
4. ArgoCD detects Git changes and syncs to target cluster
5. Automated health checks and rollback on failure

---

## 4. CI/CD STANDARDS

### 4.1 Pipeline Stages

| Stage | Tool | Requirement |
|-------|------|-------------|
| Source | GitLab CI | Triggers on commit/merge |
| Build | GitLab Runner | Multi-arch container builds |
| SAST | GitLab SAST / Semgrep | Must pass, no critical findings |
| DAST | OWASP ZAP | Must pass, no critical findings |
| Container Scan | Trivy / Snyk | No critical/high CVEs |
| SBOM | Syft | Generate and archive |
| Sign | Cosign | Sign all container images |
| Push | Harbor Registry | Push to sovereign registry |
| Deploy | ArgoCD | GitOps sync |

### 4.2 Artifact Standards

- All container images signed with Cosign
- SBOM generated in SPDX format
- Images stored in Harbor sovereign registry
- Image provenance verified before deployment

---

## 5. INFRASTRUCTURE AS CODE STANDARDS

### 5.1 Tooling

| Purpose | Technology |
|---------|------------|
| Provisioning | Terraform |
| Configuration | Ansible |
| Kubernetes Manifests | Helm + Kustomize |
| Validation | Conftest / OPA |

### 5.2 Requirements

- All infrastructure defined as code
- State stored in encrypted backend (Terraform Cloud or self-hosted)
- Peer review required for all IaC changes
- Automated drift detection via ArgoCD
- Secrets managed via Vault, never in Git

---

## 6. SECURITY BASELINES

### 6.1 CIS Benchmarks

| Level | Target |
|-------|--------|
| CIS Kubernetes Benchmark | Level 1 + Level 2 |
| CIS Docker Benchmark | Level 1 + Level 2 |
| CIS Ubuntu/CentOS Benchmark | Level 1 + Level 2 |
| CIS NGINX/Kong Benchmark | Level 1 |

### 6.2 Container Security

- Non-root containers by default
- Read-only root filesystem where possible
- No privilege escalation
- Resource limits enforced
- Image scanning in CI/CD pipeline
- Admission controllers for policy enforcement

### 6.3 Network Security

- Default-deny network policies
- mTLS for all service-to-service communication
- WAF at API gateway layer
- Network segmentation by trust zones
- Egress filtering on all namespaces

### 6.4 Hardening Requirements

- Minimal base images (distroless/alpine)
- Kernel hardening (sysctl configurations)
- Audit logging enabled on all nodes
- Regular vulnerability scanning
- Automated patch management

---

## 7. API GATEWAY STANDARDS

### 7.1 Technology Stack

| Component | Technology |
|-----------|------------|
| Gateway | Kong Gateway |
| Authentication | Keycloak (OAuth2/OIDC) |
| Certificate Management | cert-manager |
| Rate Limiting | Kong rate-limiting plugin |
| Audit Logging | Kong file-log + ELK |

### 7.2 Routing Requirements

- All external traffic through gateway
- JWT/OAuth2 validation at gateway level
- mTLS for internal service communication
- Rate limiting per consumer/endpoint
- Request/response transformation
- Circuit breaker patterns

---

## 8. IDENTITY & ACCESS MANAGEMENT STANDARDS

### 8.1 Technology Stack

| Component | Technology |
|-----------|------------|
| IAM Provider | Keycloak |
| PAM | HashiCorp Vault + Teleport |
| MFA | TOTP + FIDO2/WebAuthn |
| Biometric Integration | Custom module |

### 8.2 Access Control Model

- RBAC for coarse-grained access
- ABAC for fine-grained attribute-based decisions
- PAM for privileged access with just-in-time provisioning
- MFA mandatory for all administrative access
- Session monitoring and anomaly detection

---

## 9. PUBLIC KEY INFRASTRUCTURE STANDARDS

### 9.1 PKI Hierarchy

```
National Root CA (Offline)
├── Intermediate CA - TLS Certificates
├── Intermediate CA - Code Signing
├── Intermediate CA - Client Authentication
└── Intermediate CA - Document Signing
```

### 9.2 Requirements

- Root CA kept offline, air-gapped
- HSM for key generation and storage
- Automated certificate lifecycle via cert-manager
- CRL/OCSP for revocation checking
- Certificate transparency logging
- Maximum certificate validity: 1 year

---

## 10. EVENT-DRIVEN ARCHITECTURE STANDARDS

### 10.1 Technology Stack

| Component | Technology |
|-----------|------------|
| Message Broker | Apache Kafka / Redpanda |
| Schema Registry | Confluent Schema Registry |
| Stream Processing | Kafka Streams / ksqlDB |
| Connect | Kafka Connect |

### 10.2 Event Standards

- CloudEvents format for all events
- Schema validation via Schema Registry
- Event versioning strategy
- Dead letter queues for failed processing
- Exactly-once semantics where required
- Event replay capability for disaster recovery

---

## 11. OFFLINE SYNCHRONIZATION STANDARDS

### 11.1 Architecture

- Local SQLite/LevelDB cache on edge devices
- Delta synchronization protocol
- Conflict resolution: last-write-wins with audit trail
- Store-and-forward messaging queue
- Automatic sync when connectivity restored

### 11.2 Requirements

- Full functionality offline for critical operations
- Automatic conflict detection and resolution
- Encrypted local storage
- Bandwidth-efficient delta sync
- Sync status monitoring and alerting

---

## 12. OBSERVABILITY STANDARDS

### 12.1 Technology Stack

| Component | Technology |
|-----------|------------|
| Metrics | Prometheus + VictoriaMetrics |
| Dashboards | Grafana |
| Logs | Loki + Promtail |
| Traces | Tempo + OpenTelemetry |
| Alerts | Alertmanager + PagerDuty |

### 12.2 Requirements

- OpenTelemetry instrumentation on all services
- Golden signals: latency, traffic, errors, saturation
- SLO/SLI definitions for all critical services
- Alert routing with escalation policies
- Log retention: 90 days minimum
- Trace sampling: adaptive based on error rate

---

## 13. DATABASE STANDARDS

### 13.1 Technology Stack

| Purpose | Technology |
|---------|------------|
| Relational | PostgreSQL 16+ |
| Cache | Redis 7+ |
| Search | OpenSearch |
| Time Series | TimescaleDB (PostgreSQL extension) |

### 13.2 Requirements

- HA replication (minimum 3 nodes)
- Encryption at rest (AES-256)
- Automated backups with point-in-time recovery
- Connection pooling (PgBouncer)
- Read replicas for query offloading
- Database monitoring and alerting

---

## 14. TESTING & VALIDATION STANDARDS

### 14.1 Testing Types

| Test Type | Tool | Frequency |
|-----------|------|-----------|
| Unit Testing | JUnit/pytest | Every commit |
| Integration Testing | Testcontainers | Every PR |
| Load Testing | k6 / Locust | Weekly |
| Chaos Testing | Chaos Mesh | Monthly |
| Security Testing | OWASP ZAP / Trivy | Every PR + Monthly |
| DR Simulation | Custom runbooks | Quarterly |

### 14.2 Requirements

- All tests automated and integrated in CI/CD
- Minimum 80% code coverage
- Performance benchmarks established
- Chaos engineering for resilience validation
- Disaster recovery tested quarterly

---

## 15. GOVERNANCE & COMPLIANCE

### 15.1 Change Management

- All changes through GitOps pull requests
- Mandatory code review (minimum 2 approvers)
- Automated policy checks before deployment
- Change advisory board for production changes
- Rollback procedures documented and tested

### 15.2 Audit Requirements

- All infrastructure changes logged and immutable
- Access logs retained for 1 year minimum
- Security audit trails for all sensitive operations
- Regular third-party security assessments
- Compliance reporting automated

### 15.3 Data Sovereignty

- All data stored within national borders
- Encryption keys managed nationally
- No dependency on foreign cloud providers for core services
- Data classification and handling procedures defined
- Cross-border data transfer requires explicit approval

---

## APPENDIX A: TECHNOLOGY MATRIX

| Domain | Primary | Secondary | Notes |
|--------|---------|-----------|-------|
| Container Orchestration | RKE2 | Talos Linux | Sovereign deployment |
| Service Mesh | Istio | Linkerd | mTLS required |
| API Gateway | Kong | APISIX | Rate limiting mandatory |
| IAM | Keycloak | - | OAuth2/OIDC |
| PKI | cert-manager + Vault | - | HSM integration |
| Event Bus | Redpanda | Kafka | Schema registry |
| GitOps | ArgoCD | Flux | Declarative |
| CI/CD | GitLab CI | - | DevSecOps |
| Metrics | Prometheus | VictoriaMetrics | - |
| Logs | Loki | - | - |
| Traces | Tempo | Jaeger | OpenTelemetry |
| Secrets | HashiCorp Vault | Sealed Secrets | - |
| Database | PostgreSQL | - | HA required |
| Cache | Redis | - | Cluster mode |
| Search | OpenSearch | - | - |

---

## APPENDIX B: REFERENCE ARCHITECTURE

```
┌─────────────────────────────────────────────────────────────┐
│                    NATIONAL PLATFORM LAYER                   │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐ │
│  │ Observab │  │ Security │  │  GitOps  │  │  DevSecOps   │ │
│  │  ility   │  │ Baseline │  │ Platform │  │  Pipeline    │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                   CORE SERVICES LAYER                        │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐ │
│  │   API    │  │ Identity │  │    PKI   │  │ Event-Driven │ │
│  │ Gateway  │  │   Core   │  │  Found.  │  │     Bus      │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────┘ │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                   │
│  │  Zero    │  │ Offline  │  │ Database │                   │
│  │ Trust    │  │   Sync   │  │ Found.   │                   │
│  │   IAM    │  │  Engine  │  │          │                   │
│  └──────────┘  └──────────┘  └──────────┘                   │
├─────────────────────────────────────────────────────────────┤
│                  INFRASTRUCTURE LAYER                        │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐│
│  │           Kubernetes National Platform (RKE2)            ││
│  │  Multi-region │ HA Control Plane │ Service Mesh │ CNI    ││
│  └──────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

---

*Document Owner: SNISID National Platform Engineering Team*  
*Review Cycle: Quarterly*  
*Classification: SOVEREIGN - RESTRICTED*
