# PROMPT 297: AUTOMATED IDENTITY & ACCESS GOVERNANCE

This architecture defines the unified strategy for governing human and machine identities within the SNISID platform, ensuring that access is granted based on the "Principle of Least Privilege" (PoLP) and is always verifiable.

---

## 1. Identity Governance Architecture (Zero Trust Identity)

SNISID utilizes a centralized identity orchestration stack that separates identity providers (IDPs) from access enforcement.

- **Identity Provider (Keycloak/LDAP)**: Sovereign national identity store for human users, integrated with MFA and biometric verification.
- **Identity Orchestrator (HashiCorp Vault)**: The core engine for machine identities, dynamic secrets, and short-lived credentials.
- **Access Boundary (HashiCorp Boundary)**: Provides secure, identity-aware remote access to infrastructure without exposing internal networks.
- **Identity-Aware Proxy (IAP)**: Enforces identity-based access to internal web portals (e.g., Developer Portal, Grafana).

---

## 2. Access Workflows (The JIT Loop)

1.  **Request**: An engineer requests temporary "Admin" access to a specific production database via the **Developer Portal** (Prompt 282).
2.  **Verification**: The system automatically verifies the user's current security clearance and the presence of an active incident ticket.
3.  **Approval**: For high-tier access, the system requires a "Dual-Approval" from another authorized team member.
4.  **Provisioning (Just-In-Time)**: Vault generates short-lived, single-use credentials that expire automatically after 4 hours.
5.  **Revocation**: Access is automatically revoked as soon as the session expires or the incident ticket is closed.

---

## 3. Integration Strategy (Machine Identity)

- **SPIFFE/SPIRE Integration**: Every microservice is assigned a unique, cryptographically verifiable identity (SVID) that is used for mTLS and internal API authorization.
- **Service Account Lifecycle**: Automated creation and rotation of Kubernetes Service Accounts, managed via GitOps to prevent "Ghost Identities."
- **Credential Rotation**: Vault automatically rotates database and cloud API keys every 24 hours, ensuring that even if a secret is leaked, its window of utility is minimal.

---

## 4. Security & Privacy

- **Anonymous Access Prohibited**: All access to any system component, including read-only observability data, requires a verified identity.
- **PII Masking in Identity**: The system utilizes "Identity Aliases" in logs and forensic ledgers to protect the real names of intelligence officers from unauthorized internal disclosure.
- **HSM-Backed Roots**: The root keys for the identity orchestrator and CA are stored in national HSMs (Hardware Security Modules).

---

## 5. Governance Model

- **Access Review Campaigns**: Automated monthly "Recertification Campaigns" where managers must review and re-approve the access levels of their team members.
- **Identity Drift Detection**: AI identifies "Anomalous Identity Behavior" (e.g., "A machine identity is suddenly requesting secrets for a service it hasn't interacted with before") and triggers an automatic lockout.
- **Audit Ledger**: Every identity creation, access request, approval, and credential rotation is cryptographically signed and stored in the forensic ledger.

---

**PROMPT 297 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 298 — AUTOMATED INFRASTRUCTURE COMPLIANCE DRIFT DETECTION.**
