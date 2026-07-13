# SNISID RUNBOOK — Emergency Certificate Revocation
**Classification:** TOP-SECRET  
**Version:** 5.0.0  
**Trigger:** Compromission certificat, vol d'équipement, licenciement sécurité, CVE critique, fraude détectée, erreur administrative grave  
**RTO cible:** < 15 minutes (infrastructure), < 1 heure (citoyen), < 4 heures (judicial/government)

---

## Types de révocation d'urgence

| Code | Type | Déclencheur | Approbation | Timeline |
|------|------|-------------|-------------|----------|
| REV-E1 | **Infrastructure auto** | Falco/Wazuh/Tetragon compromission workload/node | SOC auto | < 5 min |
| REV-E2 | **Infrastructure manual** | CVE certificat / erreur admin / rotation urgence | IGC + Infra CA Manager | < 15 min |
| REV-E3 | **Citizen urgent** | Vol wallet / fraude identité / décès / demande citoyen | SOC + Citizen CA Manager | < 1h |
| REV-E4 | **Government urgent** | Licenciement / compromission poste / vol smartcard | IGC + Gov CA Manager + HR | < 2h |
| REV-E5 | **Judicial urgent** | Discipline / compromission / fraude judiciaire | Judicial Council + IGC | < 4h |
| REV-E6 | **Device urgent** | Vol station / compromission firmware / anomalie réseau | SOC + Device CA Manager | < 15 min |
| REV-E7 | **Mass revocation** | Compromission Intermediate CA / APT large | IGC + Présidence | < 1h |
| REV-E8 | **Root CA emergency** | Compromission Root / catastrophe | Cellule Crise PKI | < 6h |

---

## REV-E1 : Infrastructure Auto-Revocation (Compromission Runtime)

### 1.1 Déclencheur automatique
```bash
# Falco rule triggered: "SNISID_Falco_EmergencyEvent"
# Or: Tetragon detects unauthorized network egress from Tier-0 pod
# Or: Wazuh detects rootkit on K8s node
# Or: certificat présenté dans TLS handshake correspond à serial révoqué (replay?)
```

### 1.2 Automated response (SOC playbook)
```bash
# 1.2.1 Isolate workload
kubectl taint node ${COMPROMISED_NODE} snisid.gov.compromised=true:NoSchedule
cilium endpoint disconnect ${POD_ENDPOINT_ID}  # Isolate Cilium-level

# 1.2.2 Revoke certificate(s)
# If cert is managed by cert-manager:
kubectl annotate certificate ${CERT_NAME} -n ${NAMESPACE} cert-manager.io/revoke-certificates="true"

# If cert is managed by Vault PKI directly:
vault write snisid-pki-infra/revoke serial_number=${CERT_SERIAL}

# 1.2.3 Force CRL regeneration (out-of-cycle)
# Trigger CronJob immediately
kubectl create job --from=cronjob/snisid-crl-generator emergency-crl-$(date +%s) -n snisid-identity

# 1.2.4 Update OCSP responder (if not auto-updating from Vault)
curl -X POST http://snisid-ocsp-core.snisid-identity.svc.cluster.local:8080/admin/reload

# 1.2.5 Notify SIEM + IGC (automated)
```

### 1.3 Post-auto (human verification within 1h)
- [ ] SOC analyst reviews forensics
- [ ] If false positive → unrevoke (Vault PKI supports unrevoke if within grace period + policy allows)
- [ ] If confirmed → permanent revocation, investigate scope, rotate all affected certs

---

## REV-E4 : Government Urgent Revocation

### 4.1 Workflow BPMN (simplified)
```
Déclencheur (HR / Manager / SOC / IGC)
  │
  ▼
[Ticketing SOC] — Classification REV-E4
  │
  ▼
[Validation IGC] — Vérification identité + autorité déclarante
  │ (si non validé → rejet, audit log)
  │ (si validé → continue)
  ▼
[Government CA Manager] — Revue contexte + approuve / escalade
  │ (si escalation nécessaire → Judicial ou Présidence)
  ▼
[Vault PKI Revoke] — serial_number + reason_code (affiliationChanged / cessationOfOperation / keyCompromise)
  │
  ▼
[CRL Regeneration] — Out-of-cycle generation + publication
  │
  ▼
[OCSP Update] — Responder reload + cache invalidation
  │
  ▼
[Notification] — HR + Manager + Agent + SIEM
  │
  ▼
[Physical Recovery] — Smartcard confiscation / désactivation / destruction
  │
  ▼
[Audit IGC] — Rapport 48h
```

### 4.2 Commandes Vault PKI
```bash
# Révocation avec reason code (RFC 5280)
vault write snisid-pki-gov/revoke \
  serial_number="2026-05-25-XXXX" \
  reason="cessationOfOperation" \
  comment="REV-E4: Employee termination — ID: AGT-5678 — Approved by Gov CA Manager + IGC"

# Vérification
vault read snisid-pki-gov/cert/2026-05-25-XXXX
# Should show: revocation_time, revocation_reason

# Vérification CRL
openssl crl -in <(curl -s https://crl.snisid.gouv.local/gov/ca.crl) -text -noout | grep "2026-05-25-XXXX"
```

---

## REV-E7 : Mass Revocation (Intermediate CA Compromise / APT)

### 7.1 Context
Attaque massive : clé privée d'un workload compromis, ou erreur CA émettant certificats à mauvais destinataires, ou APT ayant obtenu accès à Vault PKI issuance.

### 7.2 Procedure
```bash
# 7.2.1 Identify scope
# Query Vault PKI for all certificates issued in compromised window
vault list snisid-pki-infra/certs | while read serial; do
  INFO=$(vault read -format=json snisid-pki-infra/cert/${serial} 2>/dev/null)
  ISSUANCE=$(echo "$INFO" | jq -r '.data.certificate' | openssl x509 -noout -startdate 2>/dev/null | cut -d= -f2)
  ISSUANCE_EPOCH=$(date -d "$ISSUANCE" +%s 2>/dev/null || echo 0)
  if [ "$ISSUANCE_EPOCH" -gt "$COMPROMISE_START_EPOCH" ] && [ "$ISSUANCE_EPOCH" -lt "$COMPROMISE_END_EPOCH" ]; then
    echo "${serial},${ISSUANCE}" >> /tmp/compromised-certs.csv
  fi
done

# 7.2.2 Mass revoke (batch)
while IFS=, read -r serial issuance; do
  vault write snisid-pki-infra/revoke \
    serial_number="${serial}" \
    reason="keyCompromise" \
    comment="REV-E7: Mass revocation — APT compromise window ${COMPROMISE_START} to ${COMPROMISE_END}"
done < /tmp/compromised-certs.csv

# 7.2.3 Emergency CRL generation
kubectl create job --from=cronjob/snisid-crl-generator emergency-mass-crl-$(date +%s) -n snisid-identity

# 7.2.4 Force cert-manager renewal for ALL certificates in affected namespaces
kubectl get certificates -A -o json | \
  jq -r '.items[] | select(.spec.issuerRef.name | contains("snisid-infra")) | [.metadata.namespace,.metadata.name] | @tsv' | \
  while read ns cert; do
    kubectl annotate certificate "${cert}" -n "${ns}" cert-manager.io/revoke-certificates="true"
    kubectl annotate certificate "${cert}" -n "${ns}" cert-manager.io/force-renew="true"
  done

# 7.2.5 Istio SDS global push (restart sidecars if necessary)
# Istio should auto-reload on cert change, but force if needed:
# istioctl proxy-status | grep NOT SENT → investigate
# If needed: kubectl rollout restart deployment -n snisid-core

# 7.2.6 Edge nodes CRL sync via emergency burst
# Trigger satellite / USB burst for edge CRL update
```

### 7.3 Post-mass revocation
- [ ] IGC audit all revoked certificates
- [ ] Forensics determine root cause (APT vector, insider, software bug)
- [ ] If Intermediate CA key compromised → NEW Intermediate CA (rotation ceremony)
- [ ] If Vault PKI compromised → full Vault cluster rebuild + new intermediate keys
- [ ] If HSM partition compromised → HSM tamper procedure (hsm-failure.md HSM-004)
- [ ] Report Présidence if > 1000 certificates revoked

---

## Revocation Reason Codes (RFC 5280 Mapping)

| Reason | Code | Usage SNISID |
|--------|------|--------------|
| unspecified | 0 | Default (avoid if possible) |
| keyCompromise | 1 | Private key compromised / extracted / stolen |
| cACompromise | 2 | CA key compromised (Intermediate or Root) |
| affiliationChanged | 3 | Employee role change / department transfer |
| superseded | 4 | Certificate replaced by newer one (renewal) |
| cessationOfOperation | 5 | Device decommissioned / employee terminated / service shut down |
| certificateHold | 6 | Temporary suspension (investigation ongoing) — REVERSIBLE |
| privilegeWithdrawn | 7 | Citizen rights withdrawn / judicial privilege revoked |
| aACompromise | 8 | Attribute authority compromise (rarely used) |

---

## Certificate Hold (Réversible)

```bash
# Temporary suspension — used during investigation
vault write snisid-pki-gov/revoke \
  serial_number="2026-05-25-XXXX" \
  reason="certificateHold" \
  comment="REV-E4-HOLD: Investigation ongoing — theft suspected — Agent: AGT-5678"

# To unhold (if investigation clears the certificate):
# Note: Not all PKI systems support unhold. Vault PKI does not support unrevoke.
# For SNISID: certificateHold is implemented by NOT adding to CRL but adding to "hold list"
# Unhold = remove from hold list, certificate becomes valid again
# Implementation: custom Vault PKI plugin or separate hold list in national registry

# SNISID-specific: Use OCSP "UNKNOWN" status during hold
# If cleared: remove from hold, OCSP returns "GOOD"
# If confirmed compromise: change revocation reason to keyCompromise, add to CRL
```

---

*Testé mensuellement via drill de révocation sur certificats test (namespace snisid-staging).*
