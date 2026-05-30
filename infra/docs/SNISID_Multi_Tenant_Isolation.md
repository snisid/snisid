# SNISID: Multi-Tenant Isolation Architecture

SNISID is a sovereign platform shared by multiple government agencies (Tenants). To ensure national security and agency autonomy, the platform enforces "Hard Multi-Tenancy" through a multi-layered isolation model.

---

## 1. The Tenant Architecture

Every Agency (Tenant) is encapsulated within a **Sovereign Boundary**.

- **Namespace Isolation**: Each agency operates in a dedicated Kubernetes Namespace (e.g., `agency-police`, `agency-tax`).
- **Identity Realm**: Tenants have dedicated Keycloak Realms and OPA Policy Bundles.
- **Resource Quotas**: Hard limits on CPU, Memory, and Storage to prevent "Noisy Neighbor" effects.

---

## 2. The 5-Layer Isolation Model

### 2.1. Network Isolation (The Moat)
- **Cilium L3/L4**: Default-deny policies prevent traffic between agency namespaces.
- **Istio L7**: mTLS is mandated for all traffic. AuthorizationPolicies ensure only `Principal: agency-tax/*` can talk to `Service: tax-db`.

### 2.2. Data Isolation (The Vault)
- **Database-per-Tenant**: Critical PII databases are physically or logically separated (separate RDS instances or isolated schemas).
- **Storage Encryption**: Per-tenant AWS EBS or On-Prem LUN encryption.

### 2.3. Encryption Isolation (The Keys)
- **Per-Tenant KMS**: Every agency has its own HSM-backed Root Key.
- **Envelope Encryption**: Data is encrypted with a Data Encryption Key (DEK) which is wrapped by the Agency's Master Key (KEK). Even the Platform Admin cannot decrypt Agency data without the Agency's KEK.

### 2.4. Policy Isolation (The Rules)
- **Scoped OPA Bundles**: Agency A cannot view or modify Agency B's authorization rules.
- **Admin Delegation**: Agency Admins manage their own staff roles within their tenant boundary.

### 2.5. Audit Isolation (The Ledger)
- **Scoped Audit Streams**: Audit events are tagged with `tenant_id`.
- **Visibility**: Agency Admins can only view their own agency's logs; only the National SOC can view the global cross-tenant stream.

---

## 3. Tenant-Aware Authorization

The system enforces isolation at every API call.

1. **JWT Claim**: Every token contains a `tenant_id` claim.
2. **OPA Enforcement**: 
   ```rego
   # Prevent cross-tenant data access
   allow {
       input.jwt.tenant_id == input.resource.tenant_id
       input.jwt.roles[_] == "agency_officer"
   }
   ```
3. **Storage Enforcement**: Database drivers automatically append `WHERE tenant_id = '...'` to all queries (Row-Level Security).

---

## 4. Cross-Agency Federation Controls

In specific cases (e.g., Police investigating Tax fraud), cross-tenant access is permitted under strict governance.

- **Bilateral Trust**: Both agencies must cryptographically sign a "Data Sharing Agreement" in the Policy Plane.
- **Purpose-Bound Access**: The `active_case_id` must be present in the request to authorize a cross-tenant query.
- **Visibility**: The "Owned" agency receives a real-time notification whenever an external agency queries their data.

---

## 5. Shared Infrastructure Security

- **Kafka Multi-Tenancy**: Topics are prefixed by tenant (e.g., `dgi.events.identity`) and protected by ACLs.
- **AI Engine Isolation**: Models are trained on tenant-specific data silos. Cross-tenant AI inference is blocked unless explicitly federated.
- **Control Plane Hardening**: The Kubernetes API and Istio Control Plane are protected by the **National SOC** and utilize hardware-attested identities.

---

## 6. Breach Containment (Tenant Quarantine)

If an agency tenant is compromised (e.g., a massive credential leak in the Tax Authority):
1. **Quarantine Trigger**: The National SOC activates the `TENANT_LOCKDOWN` state for `tenant_id: dgi`.
2. **Immediate Block**:
   - **Ingress**: All external traffic to the `tax-ns` is dropped at the API Gateway.
   - **Egress**: All outgoing calls from the `tax-ns` to the rest of the mesh are blocked.
   - **Token Revocation**: All JWTs issued for that tenant are globally invalidated.
3. **Isolation**: Other agencies (Police, Immigration) continue operating unaffected.
