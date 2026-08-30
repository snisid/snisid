# SNISID: Biometric Fusion Engine

The Biometric Fusion Engine is the high-assurance verification core of SNISID, responsible for combining multiple biometric modalities into a single, tamper-proof **Sovereign Biometric Template**.

---

## 1. Multi-Modal Fusion Architecture

The engine utilizes a "Score-Level Fusion" approach to combine signals from disparate sensors.

- **Modality Handlers**:
  - **Face Handler**: High-resolution 3D facial capture and point-cloud generation.
  - **Fingerprint Handler**: Multi-finger latent and live-scan processing.
  - **Iris Handler**: Dual-iris near-infrared capture.
- **Fusion Logic**: The system calculates a **Unified Confidence Score**. If one modality is obstructed (e.g., face mask), the engine automatically increases the weight of the Iris or Fingerprint signals.
- **Template Generation**: All biometrics are transformed into a **Biometric Vector (Embedding)**. Raw images are never stored in the primary matching database; only the cryptographically signed embeddings are persisted.

---

## 2. AI-Driven Liveness Detection

To prevent presentation attacks (spoofing), the engine performs multi-layered liveness checks.

- **Face Liveness**: Analysis of micro-expressions, skin texture, and eye movement using **CNN-based temporal analysis**.
- **Fingerprint Liveness**: Detection of artificial materials (silicone/gelatin) using thermal and multi-spectral sensors.
- **Iris Liveness**: Pupillary light reflex verification.
- **Cross-Modality Liveness**: Ensuring the "Identity Signal" is coherent across all sensors simultaneously (preventing deepfake/replay fusion).

---

## 3. Secure Processing & Confidential Computing

- **Hardware-Rooted Security**: Biometric matching and template generation occur within **Trusted Execution Environments (TEE)** such as Intel SGX or AMD SEV.
- **Biometric Vault**: The master biometric templates are stored in an HSM-backed vault, accessible only via a one-way matching API.
- **Privacy-Preserving Matching**: Matching is performed on **Homomorphically Encrypted** templates where possible, ensuring the matching server never sees the raw biometric vector in plain text.

---

## 4. Operational Workflows

- **Enrollment Flow**: Multi-sensor capture -> Liveness verification -> Feature extraction -> Deduplication (via ABIS) -> Sovereign Template signing -> Persistence.
- **Verification Flow**: Real-time capture -> Liveness check -> Template generation -> 1:1 Match against the claimed identity -> Confidence score generation -> Access Decision.
- **Template Refresh**: Automated notification for citizens to update their biometric profiles if the "Template Aging" score indicates significant drift (e.g., child growth or facial surgery).

---

## 5. National Scalability

- **Distributed Matchers**: Deployed at regional clusters for low-latency 1:1 verification.
- **GPU Acceleration**: Matcher instances use GPU acceleration for high-speed vector similarity calculations.
- **Fail-Open/Fail-Closed Policy**: Configurable security policies (e.g., *"In the event of a total Biometric Cluster failure, fall back to high-assurance Cryptographic OTP verification"*).
