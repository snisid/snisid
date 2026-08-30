# PROMPT 281: AUTOMATED CERTIFICATE MANAGEMENT

This architecture defines the automated TLS certificate lifecycle management for the SNISID platform, ensuring secure, encrypted communication between all national services and users.

---

## 1. Certificate Architecture (PKI Integration)

SNISID utilizes a multi-tier Public Key Infrastructure (PKI) to manage internal and external identities.

- **Certificate Authority (CA)**: **HashiCorp Vault** serves as the root and intermediate CA, providing dynamic PKI secrets.
- **Controller (Cert-Manager)**: Deployed in every Kubernetes cluster to manage the lifecycle of `Certificate` and `Issuer` resources.
- **Issuers**:
    - **Self-Signed**: Used for bootstrapping initial mesh communication.
    - **Vault Intermediate**: Used for all internal mTLS communication within the mesh.
    - **External ACME (Let's Encrypt)**: Used for public-facing ingress gateways (where applicable).

---

## 2. Issuance Workflows (Automated & Secure)

1.  **Request**: A microservice or ingress gateway requests a certificate via a standard Kubernetes `Certificate` manifest.
2.  **Challenge**: Cert-Manager validates the request against the configured `Issuer` (e.g., DNS-01 or HTTP-01 challenges).
3.  **Generation**: Upon validation, Cert-Manager requests a new certificate from Vault or the ACME provider.
4.  **Injection**: The resulting certificate and private key are stored as a Kubernetes Secret and automatically mounted into the pod.

---

## 3. Rotation & Lifecycle Management

- **Automated Renewal**: Cert-Manager automatically renews certificates when they reach 2/3 of their lifespan (e.g., every 60 days for a 90-day certificate).
- **Graceful Reload**: Services are configured to watch the certificate secret and reload without restart, or Istio handles the rotation transparently at the Envoy proxy level.
- **Revocation**: Vault CRL (Certificate Revocation List) and OCSP (Online Certificate Status Protocol) are used to instantly invalidate compromised certificates across the entire federation.

---

## 4. Security & Governance

- **Short-Lived Certificates**: Production certificates have a maximum validity of 24–48 hours to minimize the window of opportunity for attackers using stolen keys.
- **Hardware Security (HSM)**: The Root CA keys in Vault are protected by a FIPS 140-2 Level 3 Hardware Security Module.
- **Policy Enforcement**: Kyverno policies ensure that only authorized services can request certificates for specific domains.

---

## 5. Audit & Compliance

- **Certificate Transparency**: All issued certificates are logged in the national transparency ledger.
- **Audit Ledger**: Every issuance, renewal, and revocation event is recorded in the forensic ledger, including the requesting identity and approval chain.
- **Expiry Monitoring**: AI-driven alerts monitor all certificates across the national infrastructure and trigger a high-priority "Severity-1" alert if any certificate is within 24 hours of expiry without renewal.

---

**PROMPT 281 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 282 — SERVICE CATALOG & DEVELOPER PORTAL.**
