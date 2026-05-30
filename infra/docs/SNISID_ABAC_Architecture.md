# SNISID: Advanced ABAC Architecture

In a national-scale Zero Trust environment, roles (RBAC) provide the baseline, but Attribute-Based Access Control (ABAC) provides the intelligence. ABAC evaluates the context, risk, and attributes of every request to make high-fidelity authorization decisions.

---

## 1. The Multidimensional Attribute Model

SNISID evaluates five distinct dimensions (SREAK) for every authorization request.

| Dimension | Variable | Examples | Source |
| :--- | :---: | :--- | :--- |
| **Subject** | $S$ | `role`, `agency_id`, `clearance_level`, `biometric_confidence` | National IdP / JWT |
| **Resource** | $R$ | `classification` (Secret/PII), `data_owner_agency`, `record_status` | Database / API |
| **Environment**| $E$ | `geo_location`, `network_zone`, `time_of_day`, `device_id` | API Gateway / MDM |
| **Action** | $A$ | `read`, `write`, `delete`, `bulk_export`, `forensic_audit` | Application Logic |
| **Risk** | $K$ | `trust_score` (0-100), `anomaly_detected` (true/false) | Threat Intel Plane |

---

## 2. Policy Structure (Rego Framework)

Policies are written in **Rego** and executed by the **Open Policy Agent (OPA)**.

### 2.1. Example: Context-Aware Intelligence Query
"Allow an Intelligence Officer to query PII data *only if* they are in a secure government zone, using an MDM-attested device, and their real-time trust score is high."

```rego
package snisid.authz.abac

default allow = false

allow {
    # Subject Baseline
    input.subject.role == "intelligence_officer"
    
    # Contextual Logic (ABAC)
    input.environment.network_zone == "government_secure_intranet"
    input.environment.device_posture.is_attested == true
    
    # Risk-Awareness (Dynamic)
    input.risk.trust_score < 20 # Low risk
    
    # Resource Constraints
    input.resource.classification == "PII"
    input.action == "read"
}
```

---

## 3. Runtime Evaluation Architecture

1. **Request Interception**: The PEP (Envoy Proxy or API Gateway) intercepts the request.
2. **Context Enrichment**: The PEP gathers environmental attributes (IP, DeviceID, Time).
3. **Attribute Fetching**: The PEP calls the **Attribute Provider** (Keycloak/SPIRE/Threat Intel) to populate the request context.
4. **PDP Query**: The PEP sends a JSON payload to the OPA Sidecar (PDP).
5. **Decision**: OPA evaluates the Rego policies against the attributes and returns `ALLOW` or `DENY`.

---

## 4. AI-Driven Adaptive Evaluation Loop

The ABAC system is not static. It adapts via a feedback loop from the **Threat Intelligence Plane**.

- **Anomalous Behavioral Attributes**: If a user's API request pattern deviates from their historical "Archetype Profile," the AI engine injects a `behavioral_anomaly=true` attribute into the user's session context.
- **Dynamic Policy Tightening**: OPA policies can be configured to automatically require MFA if `behavioral_anomaly` is detected, even if the user has the correct role.

---

## 5. Decision Engine (PDP) & Enforcement (PEP)

- **PDP (Policy Decision Point)**: OPA running as a sidecar in every pod for sub-5ms local evaluation.
- **PEP (Policy Enforcement Point)**: 
    - **Istio Envoy**: For service-to-service L7 calls.
    - **API Gateway**: For external ingress traffic.
    - **Middleware**: Custom Go middleware for deep application-level data filtering.

---

## 6. Audit & Non-Repudiation

Every ABAC decision is a rich data event.
- **Decision Logs**: OPA logs the full input (all attributes evaluated) and the policy result.
- **Traceability**: These logs are streamed to the **Sovereign Audit Ledger**, allowing a forensic analyst to answer: *"Why was this specific query allowed at 2:00 AM on a Tuesday?"* by reviewing the exact attributes present at that moment.

---

## 7. Governance & Policy Lifecycle

- **GitOps for Policies**: All Rego policies are stored in a hardened Git repository and deployed via an automated CI/CD pipeline (ArgoCD).
- **Policy Testing**: Every policy change must pass a suite of unit tests (`opa test`) to ensure no "Permissive Drift" occurs.
- **National Policy Committee**: High-level cross-agency access policies must be cryptographically signed by the National Security Council before being deployed to the mesh.
