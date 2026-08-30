# SNISID Zero Trust Security Model

## 1. Objective
To eliminate implicit trust and ensure that every request for access is verified, regardless of where it originates.

## 2. Core Principles

| Principle | Obligatory | Implementation Detail |
| :--- | :---: | :--- |
| **Verify Explicitly** | Yes | Always authenticate and authorize based on all available data points. |
| **Least Privilege** | Yes | Grant only the minimum access required for the task (Just-in-Time / Just-Enough Admin). |
| **Continuous Validation** | Yes | Re-verify identity and device health throughout the session. |
| **Device Trust** | Yes | Only managed and compliant devices can access sensitive resources. |
| **Session Monitoring** | Yes | Log and analyze every action taken within a session. |

## 3. Architecture Pillars

### Identity (The New Perimeter)
- **Strong Authentication:** Mandatory MFA (FIDO2/WebAuthn) for all users.
- **Adaptive Auth:** Triggering higher friction (additional MFA) based on risk (location, time, device).

### Devices
- **Device Inventory:** Only known devices are allowed.
- **Posture Check:** Checking for OS updates, disk encryption, and antivirus status before access.

### Network (Micro-segmentation)
- **Software Defined Perimeter (SDP):** Hiding resources from the public internet.
- **Micro-segmentation:** Dividing the network into small zones to prevent lateral movement.

### Applications & Workloads
- **Service Mesh (Istio/Linkerd):** Enforcing mTLS between microservices.
- **API Security:** Validating every API call via an API Gateway with OAuth2/OIDC.

## 4. Trust Engine Logic
`Request` $\rightarrow$ `Trust Engine (Identity + Device + Context)` $\rightarrow$ `Policy Decision Point (PDP)` $\rightarrow$ `Policy Enforcement Point (PEP)` $\rightarrow$ `Access Granted/Denied`

## 5. Transition Path
1. **Phase 1:** Implement MFA and Device Inventory.
2. **Phase 2:** Implement Micro-segmentation in Kubernetes.
3. **Phase 3:** Remove all VPNs in favor of Identity-Aware Proxies.
4. **Phase 4:** Implement continuous risk-based session evaluation.
