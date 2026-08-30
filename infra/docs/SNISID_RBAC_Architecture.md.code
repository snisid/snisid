# SNISID: Production-Grade RBAC Architecture

The SNISID RBAC system provides the structural foundation for access control, defining a clear hierarchy of roles and permissions that ensure every actor—human or machine—operates with the absolute minimum privilege required.

---

## 1. Role Hierarchy (The 5-Level Model)

SNISID uses a strictly hierarchical and additive role model.

| Level | Role Name | Scope | Key Capabilities |
| :--- | :--- | :--- | :--- |
| **L0** | **Global Root Admin** | System-Wide | PKI Root management, Global configuration. (MFA + HSM required). |
| **L1** | **National SOC Analyst** | Infrastructure | Platform-wide audit review, Threat hunting, Emergency lockdown. |
| **L2** | **Agency Administrator** | Tenant-Specific | Local user provisioning, SCIM management, Local policy overrides. |
| **L3** | **Intelligence Officer** | Functional | Data query, enrollment processing, case management. |
| **L4** | **Standard User / Citizen** | Personal | Self-service profile updates, credential renewal. |

**Complete Hierarchy Model**: See the [SNISID Role Hierarchy Model](file:///c:/Users/sopil/Desktop/SNISID/SNISID_Role_Hierarchy_Model.md) for multi-domain inheritance (SOC, Audit, AI), Mutually Exclusive Roles (MER), and escalation controls.

---

## 2. Fine-Grained Permission Matrix

Permissions are granular and resource-specific.

| Function | L3: Officer | L2: Agency Admin | L1: SOC Analyst |
| :--- | :---: | :---: | :---: |
| **Identity:Create** | ✅ | ✅ | ❌ |
| **Identity:Read_PII** | ✅ (Assigned) | ❌ | ❌ |
| **Identity:Delete** | ❌ | ✅ (4-Eyes) | ❌ |
| **Audit:Read_Local** | ✅ | ✅ | ✅ |
| **Audit:Read_Global** | ❌ | ❌ | ✅ |
| **Security:Lockdown** | ❌ | ✅ (Agency) | ✅ (National) |

---

## 3. Delegated Administration & Tenancy

To maintain agency sovereignty, SNISID enforces **Tenant Isolation**.

- **Scoped Authority**: An `Agency Administrator` for the Tax Authority (DGI) can only view and manage identities belonging to `tenant_id: dgi`.
- **RBAC Delegation**: Agency Admins can create sub-roles within their tenant (e.g., `role:tax_auditor_l2`) but cannot elevate their own privileges to National-level roles.

---

## 4. SOC-Specific & Emergency Privileges

### 4.1. The SOC Role Suite
- **Threat Hunter**: Read-only access to all network flow logs and audit trails across all agencies.
- **Incident Responder**: Ability to trigger the "Kill Switch" for specific identities or entire subnets.
- **Audit Reviewer**: Special access to the **Sovereign Audit Ledger** for cryptographic proof of non-repudiation.

### 4.2. Emergency "Break-Glass" Roles
- **Emergency Responder**: A dormant role that grants temporary access to life-safety data (e.g., blood type, allergies).
- **Trigger**: Requires an explicit justification and triggers a high-priority SOC alert.

---

## 5. Just-In-Time (JIT) Privilege Elevation

SNISID eliminates the risk of "Standing Privileges" for high-risk actions.

1. **Request**: An officer needs to perform a bulk data export (High Risk).
2. **Elevation**: They request temporary assignment to the `role:data_exporter` role.
3. **Approval**: Requires a second officer's cryptographic approval (Four-Eyes Principle).
4. **Duration**: The role is automatically stripped after 60 minutes or upon logout.
5. **Re-Auth**: Mandatory biometric check required at the moment of elevation.

---

## 6. Enforcement Workflow

The RBAC model is integrated into the **Policy Plane**.

1. **Identity Issuance**: The National IdP embeds the user's `roles` and `tenant_id` into the signed JWT.
2. **PEP Evaluation**: The API Gateway/Sidecar sends the JWT to OPA.
3. **Rego Policy**: OPA evaluates the RBAC level. 
   * *Logic*: `allow if user.role == "agency_admin" and user.tenant_id == resource.tenant_id`.
4. **Context Enrichment**: OPA adds ABAC context (Time, Geo, Risk) to the final decision.

---

## 7. Role Governance Strategy

- **Recertification**: All L1 and L2 roles must be manually recertified every 90 days.
- **Separation of Duties (SoD)**: The system prevents a single user from holding conflicting roles (e.g., `Agency Admin` and `Audit Reviewer`).
- **Cryptographic Audit**: Every change to the RBAC mapping (User -> Role) is logged as a signed event to prevent administrative tampering.
