# SNISID: Quantified Risk Matrix
## Institutional & Operational Risk Assessment (WB/IMF Standard)

---

| Risk Category | Risk Description | Probability | Impact | Mitigation Strategy |
| :--- | :--- | :---: | :---: | :--- |
| **Cybersecurity** | Breach of the central identity registry by state or non-state actors. | High | Critical | Zero-Trust Architecture, mTLS everywhere, Hardware Security Modules (HSM), 24/7 Autonomous SOC monitoring. |
| **Institutional** | Resistance from existing agencies to adopt a unified identity layer. | Medium | High | Phased rollout with incentive structures, mandatory inter-agency laws, and direct ministerial oversight. |
| **Data Privacy** | Misuse of biometric data or unauthorized surveillance. | Medium | Critical | GDPR alignment, OPA-enforced ABAC policies, immutable audit logs with Merkle-tree hashing, and judicial oversight. |
| **Operational** | System downtime during critical government operations. | Low | High | Multi-region active-active deployment, K8s self-healing, and Disaster Recovery (DR) sites with 15-min RPO. |
| **Political** | Changes in government leading to project abandonment. | Medium | High | Independent SNISID Authority with multi-year non-partisan mandate and international financing treaty. |
| **Financial** | Misallocation of funds or procurement delays. | Low | Medium | World Bank STEP procurement portal usage, independent financial audits, and PMU capacity building. |

---

### Risk Rating Summary
- **Overall Implementation Risk:** Substantial
- **Political and Governance:** High
- **Technical/Security:** High
- **Fiduciary:** Low (with WB controls)
