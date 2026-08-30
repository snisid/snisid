# SNISID: Identity Verification & Risk Scoring Engine

The Identity Risk Engine (IRE) provides a continuous, real-time "Trust Temperature" for every human and machine identity in the national ecosystem.

---

## 1. Identity Verification Stream Processor (Prompt 109)

The IRE processes verification events from border crossings, mobile app logins, and administrative portals.

### 1.1. Multi-Source Validation
- **Biometric Consistency**: Compares live match results against the **Sovereign Biometric Vault** (SHA-256 hash).
- **Device Correlation**: Verifies the `Device_SVID` and hardware fingerprint against the registered session context.
- **Behavioral Verification**: Checks if the current verification action matches the user's **Temporal Profile** (e.g., "Is this person normally at this border crossing at this time?").

---

## 2. Streaming Risk Scoring Engine (Prompt 110)

Risk scores are dynamic, not static. Every event re-calculates the `Confidence_Score`.

### 2.1. Scoring Formula (The Trust Vector)
The risk score $R$ is a weighted aggregate of multiple vectors:
$$R = (W_{auth} \times C_{auth}) + (W_{geo} \times C_{geo}) + (W_{behav} \times C_{behav}) + (W_{threat} \times C_{threat})$$

| Vector | Description | Influence |
| :--- | :--- | :--- |
| **Auth ($C_{auth}$)** | MFA strength and freshness. | High |
| **Geo ($C_{geo}$)** | Location risk and impossible travel detection. | Medium |
| **Behav ($C_{behav}$)** | Statistical deviation from established norms. | Medium |
| **Threat ($C_{threat}$)** | Source IP risk and global blacklist status. | High |

### 2.2. Adaptive Scoring Strategy
- **Decay**: Risk scores "cool down" over time if no suspicious activity is detected.
- **Boost**: A single "High-Severity Anomaly" (e.g., login from a sanctioned region) can instantly spike the risk score to 1.0 (Critical).

---

## 3. Runtime Decision Engine

The IRE interfaces directly with the **OPA Policy Plane**.

- **Decision**: `Risk_Score > Threshold` triggers an automatic **TAAC (Threat-Aware Access Control)** reaction.
- **Actions**:
  - **Re-Auth**: Force a new biometric MFA challenge.
  - **Step-Down**: Automatically reduce user permissions to "Read-Only".
  - **Quarantine**: Revoke the SPIFFE SVID and kill all active sessions.

---

## 4. AI-Assisted Scoring

- **Neural Evaluator**: A Flink job that runs a lightweight **TensorFlow/PyTorch** model on the event stream to detect subtle fraud signatures.
- **Continuous Learning**: The model is periodically retrained using replayed data from the **Forensic Replay Engine**.
