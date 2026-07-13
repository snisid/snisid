# SNISID: Secure Token Lifecycle Architecture

Tokens are the primary currency of authorization within SNISID. To prevent credential theft and replay attacks, the platform employs a hardened, high-assurance token model that binds identity to both hardware and real-time risk.

---

## 1. JWT Architecture: The Sovereign Claim Set

SNISID uses **JWT (JSON Web Tokens)** with a standardized payload that enables context-aware authorization at the edge.

- **Header**: Standard RS256/ES256 signing with `kid` (Key ID).
- **Payload**:
  - `sub`: The National Identity UUID (NID).
  - `nid`: Primary Identity Document Number.
  - `tenant_id`: The Agency ID (Jurisdictional boundary).
  - `trust_score`: The real-time trust score from the CAS.
  - `jkt`: JWK Thumbprint (used for **DPoP Binding**).
  - `nonce`: Unique request identifier (optional for high-risk actions).

---

## 2. Token Types & Lifecycle

### 2.1. Access Tokens (The Bearer)
*   **Standard TTL**: 15 minutes to 4 hours (Dynamic).
*   **Binding**: Strictly bound to the client's public key (DPoP).
*   **Usage**: Sent in the `Authorization: DPoP <token>` header for all API calls.

### 2.2. Refresh Tokens (The Key)
*   **Standard TTL**: 12 hours (Standard) to 30 days (Long-lived).
*   **Storage**: Must be stored in the device's Secure Enclave or TPM.
*   **Rotation**: Every time a refresh token is used, it is rotated. The old token is invalidated.

---

## 3. Device-Binding (DPoP)

To eliminate the risk of "Token Theft," SNISID enforces **DPoP (Demonstrating Proof-of-Possession)**.

1. **Generation**: The client generates an ephemeral asymmetric key pair on the device.
2. **Binding**: When requesting a token, the client sends its public key. SNISID embeds the public key thumbprint (`jkt`) into the JWT.
3. **Usage**: For every API call, the client must sign a "DPoP Proof" (containing a timestamp and the HTTP method/URL) using its private key.
4. **Validation**: The API Gateway verifies both the JWT signature AND the DPoP Proof signature against the `jkt` claim. 
   - **Benefit**: If an attacker steals the JWT, they cannot use it without the private key from the victim's device.

---

## 4. Revocation & Introspection

Revocation must be global and instant (< 30 seconds).

### 4.1. The Bloom Filter Strategy
For high-speed validation at the Gateway, SNISID uses **Redis-backed Bloom Filters**.
- **Process**: When a session is killed, the `jti` (JWT ID) is added to a global Bloom Filter.
- **Validation**: Gateways check the Bloom Filter for every request. If a match is found, they perform a secondary hard-check against Redis.

### 4.2. OAuth2 Introspection (RFC 7662)
For high-sensitivity services (e.g., National Root CA access), microservices perform a synchronous **Token Introspection** call to the Identity Hub to verify the token's status in real-time.

---

## 5. Risk-Adaptive Expiration

Token TTL is not static; it is dictated by the **Continuous Authentication (CAS)** score.

| Trust Score | Access Token TTL | Refresh Token TTL | MFA Requirement |
| :---: | :--- | :--- | :--- |
| **0.9 - 1.0** | 4 Hours | 30 Days | Standard (Biometric at Login) |
| **0.7 - 0.9** | 1 Hour | 12 Hours | Periodic Re-auth |
| **0.5 - 0.7** | 15 Minutes | 2 Hours | Mandatory Step-up MFA |
| **< 0.5** | **REVOKED** | **REVOKED** | Account Locked |

---

## 6. Cryptographic Signing (JWKS)

- **HSM-Backed Keys**: All token signing keys are stored and executed inside a FIPS 140-2 Level 3 HSM.
- **Rotation**: Signing keys are rotated every 24 hours. The old keys are kept in the **JWKS (JSON Web Key Set)** endpoint for a 1-hour overlap to prevent verification failures.
- **Algorithm**: Defaulting to **EdDSA (Ed25519)** for high performance and strong security.

---

## 7. Threat Mitigation Strategy

| Threat Scenario | Mitigation Mechanism |
| :--- | :--- |
| **Token Theft** | DPoP binding ensures the token is useless without the hardware-bound private key. |
| **Token Replay** | Short TTLs and DPoP Proof timestamps minimize the window for replay. |
| **Privilege Escalation** | OPA re-evaluates permissions based on the `trust_score` claim in the JWT. |
| **Session Hijacking** | Continuous trust scoring detects anomalous IP/Geo shifts and triggers revocation. |
| **IdP Compromise** | HSM-backed signing keys ensure that even a software breach cannot result in the extraction of signing material. |
