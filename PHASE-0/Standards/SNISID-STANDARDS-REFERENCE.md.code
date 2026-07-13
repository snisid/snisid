# 📐 SNISID — STANDARDS TECHNIQUES NATIONAUX
## Référentiel National des Standards de l'Infrastructure Numérique Souveraine

**Document ID :** SNISID-STD-001  
**Version :** 1.0.0  
**Date :** Mai 2026  
**Propriétaire :** Autorité Nationale Numérique (AND) — Bureau de Gouvernance de l'Interopérabilité (BGI)  
**Classification :** Usage Gouvernemental  

---

## 1. PRÉAMBULE

Le présent référentiel fixe les **standards techniques obligatoires** pour tout composant du SNISID, toute agence connectée au bus national, et tout prestataire contractant avec l'État haïtien dans le cadre du programme SNISID.

Tout écart à ces standards doit être formellement approuvé par le **Change Control Board (CCB)** et documenté.

---

## 2. STANDARDS DE SÉCURITÉ

### 2.1 Cryptographie
| Standard | Algorithme | Usage | Niveau |
|----------|-----------|-------|--------|
| Chiffrement symétrique | AES-256-GCM | Données at-rest + transport | Obligatoire |
| Asymétrique | RSA-4096 ou ECDSA P-384 | PKI, signatures | Obligatoire |
| Hachage | SHA-256 / SHA-384 | Intégrité, signatures | Obligatoire |
| KDF | PBKDF2 / Argon2id | Dérivation clés | Obligatoire |
| TLS | TLS 1.3 (TLS 1.2 toléré temporairement) | Transport | Obligatoire |
| Post-quantum | CRYSTALS-Kyber (migration 2028+) | Forward security | Planifié |

### 2.2 Authentification
| Standard | Usage | Conformité |
|----------|-------|------------|
| FIDO2 / WebAuthn | Admin & Agents (MFA hardware) | Obligatoire pour accès privilégié |
| TOTP (RFC 6238) | MFA logiciel | Accepté pour agents terrain |
| OIDC (OpenID Connect 1.0) | SSO inter-agences | Obligatoire |
| SAML 2.0 | Fédération legacy | Toléré (migration OIDC) |
| OAuth 2.1 | API Authorization | Obligatoire |
| PKI / X.509 | Certificats services + agents | Obligatoire |
| mTLS | Service-to-service | Obligatoire (via Istio) |

### 2.3 Biométrie
| Standard | Usage |
|----------|-------|
| ISO/IEC 19794-2 | Templates empreintes digitales |
| ISO/IEC 19794-5 | Photo faciale (ICAO 9303 compliant) |
| ISO/IEC 19794-6 | Templates iris |
| NFIQ 2.0 | Score qualité empreintes |
| ISO/IEC 30107 | Liveness detection (Presentation Attack Detection) |

### 2.4 HSM
| Standard | Niveau requis |
|----------|--------------|
| FIPS 140-2 | Niveau 3 minimum (Niveau 4 pour Root CA) |
| Common Criteria | EAL4+ pour HSM Root CA |
| PKCS#11 | Interface logicielle standard HSM |

---

## 3. STANDARDS API

### 3.1 REST APIs
```yaml
standard: OpenAPI 3.1.0
format: JSON (primary) + CBOR (edge/offline)
versioning: URI path (/v1, /v2) — semver
charset: UTF-8
pagination: cursor-based (next_cursor param)
errors: RFC 7807 (Problem Details for HTTP APIs)
security: OAuth 2.1 + mTLS + JWT (RS256/ES256)
content-type: application/json; charset=utf-8
```

### 3.2 Événements Asynchrones
```yaml
standard: CloudEvents 1.0
transport: Apache Kafka (topics structurés)
encoding: Apache Avro (schema versioning via Schema Registry)
schemas: centralisés dans Schema Registry national
naming: {domaine}.{entité}.{action}.{version}
  examples:
    - etat-civil.naissance.declaree.v1
    - identite.citoyen.cree.v1
    - securite.incident.detecte.v1
```

### 3.3 gRPC (services internes haute performance)
```yaml
standard: Protocol Buffers 3 (proto3)
transport: HTTP/2 + TLS 1.3
usage: Biometric matching (AFIS), internal microservices
```

### 3.4 GraphQL (analytics, dashboards)
```yaml
standard: GraphQL June 2018+
usage: Analytics APIs, reporting APIs
transport: HTTP/1.1 + HTTP/2
```

---

## 4. STANDARDS DONNÉES

### 4.1 Identifiants
| Identifiant | Format | Standard |
|------------|--------|----------|
| NIN (National Identity Number) | 13 chiffres + checksum Luhn | Standard national SNISID |
| UUID | v4 (RFC 4122) | Entités internes |
| Codes géographiques | CNIGS officiel | Standard national |
| Codes pays | ISO 3166-1 alpha-3 | International |
| Codes langue | BCP 47 | International |
| Date/Heure | ISO 8601 + timezone | International |

### 4.2 Données Santé
```
Standard: HL7 FHIR R4 (R4B acceptable)
Usage: Naissances (MSPP), Décès, Dossiers médicaux interconnectés
Format: JSON FHIR
```

### 4.3 Données Géographiques
```
Vector: GeoJSON (RFC 7946)
Services: OGC WMS, WFS, WMTS
Projection: WGS84 (EPSG:4326) par défaut
Topographic: ISO 19115 metadata
Source: CNIGS (Centre National de l'Information Géo-Spatiale)
```

### 4.4 Documents Officiels
```
Format archivage: PDF/A-3 (ISO 19005-3)
Signature électronique: XAdES-LTA (Long-Term Archive)
Timestamping: RFC 3161 (TSA nationale)
QR Code vérification: ISO/IEC 18004 + URL vérification SNISID
```

### 4.5 Cartes & eID
```
Carte physique: ISO/IEC 7816 (contact)
NFC: ISO/IEC 14443 (sans contact)
Chip OS: ISO/IEC 24727 (Government eID)
Passeport: ICAO Doc 9303 (MRTD biométrique)
Crypto carte: RSA-2048 minimum (ECDSA P-256 recommandé)
```

---

## 5. STANDARDS INFRASTRUCTURE

### 5.1 Kubernetes
```yaml
distribution: RKE2 (production) / k3s (edge)
version: K8s 1.29+ (latest stable)
CNI: Cilium (eBPF) — v1.14+
CSI: Rook-Ceph (storage) / local-path (edge)
ingress: NGINX Ingress Controller + ModSecurity WAF
service_mesh: Istio 1.20+ (mTLS, traffic management)
gitops: ArgoCD 2.9+ (declarative state)
secrets: External Secrets Operator + HashiCorp Vault
monitoring: kube-prometheus-stack (Prometheus + Grafana)
```

### 5.2 Observabilité
```yaml
standard: OpenTelemetry (OTEL)
metrics: Prometheus (scrape model)
traces: Jaeger / Tempo (distributed tracing)
logs: Fluent Bit → Loki / Elasticsearch
dashboards: Grafana 10+
alerting: Alertmanager → PagerDuty/OpsGenie
SLO tracking: Pyrra ou Sloth
```

### 5.3 CI/CD
```yaml
vcs: GitLab CE (self-hosted, souverain)
ci_engine: GitLab CI
image_registry: Harbor (self-hosted)
artifact_signing: Sigstore/Cosign
sast: Semgrep + SonarQube
sca: Trivy (container + deps)
dast: OWASP ZAP
secrets_scanning: Gitleaks
iac_scanning: Checkov (Terraform)
gitops_deploy: ArgoCD
```

### 5.4 IaC
```yaml
provisioning: Terraform 1.6+
configuration: Ansible 8+
k8s_manifests: Helm 3.12+ + Kustomize
policy_as_code: OPA (Open Policy Agent) + Rego
network_policy: Cilium NetworkPolicy
```

### 5.5 OS & Hardening
```yaml
os_production: Talos Linux (immutable, readonly FS)
os_edge: Ubuntu Server LTS 24.04 (hardened)
benchmarks: CIS Level 2 (Linux, K8s, PostgreSQL, Docker)
tpm: TPM 2.0 obligatoire pour edge nodes
disk_encryption: LUKS2 + AES-XTS-512
boot: Secure Boot (UEFI) + Measured Boot
```

---

## 6. STANDARDS BASES DE DONNÉES

| Usage | Technologie | Standard |
|-------|-------------|----------|
| Core Transactionnel | PostgreSQL 16+ (Patroni HA) | ACID, MVCC |
| Core distribué multi-région | CockroachDB 23+ | ACID distributed |
| Cache / Sessions | Redis 7+ (Cluster mode) | — |
| Recherche | OpenSearch 2.x | — |
| Analytics | ClickHouse 24+ | OLAP |
| Object Storage | MinIO (S3-compatible) | AWS S3 API |
| Time-series | TimescaleDB (sur PG) | — |
| Bloc storage | Ceph 18+ (via Rook) | RADOS |

---

## 7. STANDARDS RÉSEAUX

### 7.1 Protocoles
| Couche | Standard |
|--------|---------|
| Transport | TCP/IP (IPv4 + IPv6 dual-stack) |
| DNS | BIND 9 / CoreDNS (interne K8s) |
| NTP | NTPv4 (RFC 5905) — serveurs souverains |
| BGP | BGP-4 (inter-DC) |
| VPN | WireGuard (edge-to-DC) / IPsec IKEv2 (legacy) |
| SD-WAN | Cilium Cluster Mesh (K8s multi-cluster) |

### 7.2 Segmentation
```
Zone Internet publique ──[WAF + DDoS protection]──▶
Zone DMZ (API Gateway) ──[Firewall L7]──▶  
Zone Application (K8s) ──[Cilium Network Policy]──▶
Zone Data (Bases de données) ──[Microsegmentation]──▶
Zone HSM/PKI (air-gappée)
```

---

## 8. STANDARDS PROJETS & DOCUMENTATION

### 8.1 Documentation
```yaml
format: Markdown (CommonMark + GFM)
diagrams: Mermaid (embedded in markdown)
api_docs: OpenAPI 3.1 (YAML)
architecture: C4 model (PlantUML / Mermaid)
versionning: Git (branching model GitFlow adapté)
review: Pull Request obligatoire (4 yeux min)
```

### 8.2 Nomenclature
```
Services: snisid-{domain}-{function} (ex: snisid-identity-enrollment)
Namespaces K8s: snisid-{env}-{domain} (ex: snisid-prod-identity)
Topics Kafka: {domain}.{entity}.{action}.v{n}
Secrets Vault: snisid/{env}/{service}/{key}
Repos Git: snisid/{domain}/{component}
```

### 8.3 Conventions Code
```
Java/Spring: Google Java Style Guide
Go: Effective Go + gofmt
Python: PEP 8 + Black
TypeScript: ESLint Airbnb + Prettier
Commits: Conventional Commits (feat/fix/docs/chore)
```

---

## 9. CERTIFICATIONS REQUISES PAR COMPOSANT

| Composant | Certification |
|-----------|--------------|
| Datacenter primaire PaP | Uptime Institute Tier III |
| Datacenter DR Cap-Haïtien | Uptime Institute Tier III |
| HSM Root CA | FIPS 140-2 Level 4 + CC EAL4+ |
| HSM Issuing CA | FIPS 140-2 Level 3 |
| Plateforme SNISID globale | ISO/IEC 27001 |
| Continuité d'activité | ISO 22301 |
| Biometric AFIS | FBI Appendix F (empreintes) + ICAO 9303 (photo) |
| Processus développement | ISO/IEC 27034 |

---

## 10. CYCLE DE VIE DES STANDARDS

- **Révision** : annuelle minimum, ou sur incident majeur
- **Approbation modifications** : CCB + AND + CISO
- **Communication** : publiée au portail BGI (intranet État)
- **Dérogation** : formulaire de dérogation BGI, validé CCB, durée max 12 mois
- **Archivage** : versions précédentes conservées 5 ans

---

*Document approuvé par le Bureau de Gouvernance de l'Interopérabilité (BGI)*  
*SNISID — République d'Haïti — Classification : Usage Gouvernemental*
