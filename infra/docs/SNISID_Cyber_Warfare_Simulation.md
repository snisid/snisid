# SNISID: Cyber Warfare Simulation & Red-Teaming

The Cyber Warfare Simulation layer provides a safe, isolated "Sovereign War Room" for continuously testing SNISID's resilience against advanced persistent threats (APTs) and autonomous cyber-attacks.

---

## 1. Breach Propagation Simulator

The Simulator uses the **National Relationship Graph (Neo4j)** to model how an attack could spread through the system.

- **Monte Carlo Modeling**: Runs thousands of simulations to identify the "Shortest Path to Compromise" for critical assets (e.g., the Biometric Root or KMS).
- **Probability Vectors**: Each node relationship in the graph is assigned a "Compromise Probability" based on its current **ISTS (Service Trust Score)** and authentication strength.
- **Blast Radius Analysis**: Visualizes the potential impact of a single service compromise on the wider national infrastructure.

---

## 2. AI Red-Team Swarm (Prompt 141)

The Red-Team Swarm is a collection of "Adversarial Agents" designed to find and exploit weaknesses in the platform's defense.

### 2.1. Agent Personas
- **The Phisher**: Simulates social engineering and credential theft attempts.
- **The Lateral Shifter**: Attempts to bypass Istio/Cilium policies to move between namespaces.
- **The Exfiltrator**: Models low-and-slow data exfiltration patterns to test Flink's anomaly detection.
- **The Poisoner**: Attempts to inject malicious data into the AI training pipeline or the Graph mesh.

### 2.2. Training & Feedback
- The Red-Team Swarm uses **Reinforcement Learning** to improve its tactics. It "wins" by successfully reaching a target asset without triggering a SOAR playbook.
- **Continuous Improvement**: Successful Red-Team tactics are automatically converted into new **Detection Rules** for the Blue-Team (SOC Agent Swarm).

---

## 3. Sovereign War Room Orchestration

The War Room provides the operational environment for these simulations.

- **Isolation (The Sandbox)**: Simulations are strictly confined to the **`snisid-war-room`** namespace. Hardware-level isolation is enforced to ensure simulation traffic never touches the real national backbone.
- **Digital Twin**: The simulator creates a "Digital Twin" of the SNISID infrastructure—mirroring the graph topology and service mesh configuration without accessing real citizen PII.
- **Replay-to-Sim**: Real historical incidents from the **Forensic Replay Engine** are injected into the War Room to test if the current system state would have prevented the historical attack.

---

## 4. Operational Workflows

1.  **Objective Definition**: National Security Director sets a goal (e.g., "Simulate a ransomware attack on the Banking Regional Spoke").
2.  **Simulation Launch**: The Digital Twin is hydrated, and the Red-Team Swarm is activated.
3.  **Autonomous Conflict**: The AI Red-Team (Offense) and AI SOC Swarm (Defense) compete in real-time within the sandbox.
4.  **Reporting**: A **Resilience Score** is generated, detailing which defense layers failed and providing specific remediation steps.

---

## 5. Security & Safety Controls

- **Kill-Switch**: A physical and software-based "Big Red Button" that instantly terminates all simulation agents and wipes the Digital Twin memory.
- **Guardrails**: Rego-based constraints that prevent simulation agents from making external network calls or accessing production-grade HSM keys.
- **Audit**: Every action taken by a Red-Team agent is recorded with a **"Simulation-Only"** flag in the Sovereign Audit Ledger for post-mission analysis.
