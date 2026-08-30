# PROMPT 282: SERVICE CATALOG & DEVELOPER PORTAL

This architecture defines the central "Front-Door" for all SNISID engineering and operations, providing a unified view of the entire national intelligence ecosystem.

---

## 1. Portal Architecture (Backstage-Based)

SNISID utilizes **Backstage** (by Spotify) as the foundational framework for the Internal Developer Portal (IDP).

- **Software Catalog**: A centralized registry of all microservices, libraries, and infrastructure components.
- **TechDocs**: Automated aggregation of Markdown documentation from every repository into a searchable, central UI.
- **Software Templates**: Standardized "Golden Paths" for creating new services (e.g., "Create Hardened Go Microservice").
- **Plugins**: Custom integrations for Grafana, ArgoCD, and Jira to show real-time health and deployment status directly in the portal.

---

## 2. Discovery Workflows

1.  **Onboarding**: New services are automatically registered via a `catalog-info.yaml` file located in the service's Git repository.
2.  **Metadata Enrichment**: Services are tagged with `lifecycle` (production/deprecated), `owner` (agency-team), and `system` (national-identity).
3.  **API Explorer**: Centralized view of all internal and external APIs (gRPC/Protobuf/OpenAPI), allowing developers to discover and test interfaces without leaving the portal.
4.  **Relationship Mapping**: Visual graph showing service dependencies, helping architects identify high-impact failure points.

---

## 3. Documentation Integration (TechDocs-as-Code)

- **Source of Truth**: Documentation is stored alongside code in Markdown format.
- **Continuous Aggregation**: The CI pipeline for each repo triggers a Backstage build that publishes the rendered documentation to the central portal.
- **Global Search**: A unified, AI-enhanced search engine allows officers to find technical specifications, operational runbooks, and national security directives instantly.

---

## 4. Governance Model

- **RBAC Integration**: National SSO (Keycloak) controls which services and documentation are visible to which agencies.
- **Security Scorecards**: Automated grading of services based on test coverage, open vulnerabilities, and documentation completeness (displayed prominently in the portal).
- **Audit Ledger**: Every access to the service catalog and documentation is recorded in the forensic ledger to prevent unauthorized reconnaissance.

---

## 5. Lifecycle Management (Scaffolding)

- **One-Click Provisioning**: Authorized developers can trigger the creation of a new, fully compliant microservice stack (Git Repo, CI/CD, Infra, Monitoring) via a single portal template.
- **Deprecation Workflow**: Services marked as `deprecated` trigger automated notifications to downstream consumers and are eventually hidden from the primary catalog after a grace period.
- **Standards Enforcement**: Templates ensure that every new service includes mandatory SNISID libraries for authentication, logging, and tracing by default.

---

**PROMPT 282 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 283 — AUTOMATED INFRASTRUCTURE BENCHMARKING.**
