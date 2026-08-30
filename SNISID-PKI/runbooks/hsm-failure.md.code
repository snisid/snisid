# SNISID RUNBOOK — HSM Failure & Recovery
**Classification:** TOP-SECRET  
**Version:** 5.0.0  
**Trigger:** HSM offline, partition inaccessible, tamper detection, firmware corruption, network isolation  
**RTO cible:** < 15 minutes (failover DR), < 4 heures (reconstruction Core)

---

## Scénarios de défaillance HSM

| Code | Scénario | Symptômes | Impact |
|------|----------|-----------|--------|
| HSM-001 | **HSM Core hardware failure** | LED rouge, pas de ping, erreur POST, ventilateur bruit | CA operations halted on Core, DR active |
| HSM-002 | **Partition locked / M-of-N failure** | Authentication failures, partition disabled after max login attempts | Specific CA domain halted |
| HSM-003 | **Firmware corruption / upgrade failure** | HSM boots but rejects operations, checksum fail | All operations halted |
| HSM-004 | **Tamper detection (physical)** | Tamper-evident seal broken, temperature extreme, voltage anomaly | HSM auto-zeroizes (keys destroyed) |
| HSM-005 | **Network isolation** | HSM ping OK, no client connection, firewall block, switch failure | Operations halted but HSM healthy |
| HSM-006 | **HSM DR sync failure** | DR HSM not receiving replication, stale keys | Risk if Core fails simultaneously |

---

## HSM-001 : Core Hardware Failure

### 1.1 Detection
```bash
# Prometheus alert: hsm_status{hsm="HSM-CORE-01"} == 0
# Or: ping HSM-CORE-01.mgmt (10.1.0.10) timeout
```

### 1.2 Immediate actions (H+0)
```bash
# 1.2.1 Confirm failure (not network blip)
ping -c 10 10.1.0.10
nmap -p 1792 10.1.0.10
# If no response after 3 attempts → declare HSM-001

# 1.2.2 Failover to HSM-CORE-02 (hot standby)
# Update Vault PKI configuration to use HSM-CORE-02
kubectl edit configmap vault-config -n snisid-identity
# Change HSM IP: 10.1.0.10 → 10.1.0.11
kubectl rollout restart statefulset vault -n snisid-identity

# 1.2.3 Verify Vault unseal + CA operations
vault status
vault read snisid-pki-infra/cert/ca
vault write -f snisid-pki-infra/issue/cluster-local common_name=test-failover.snisid.gouv.local ttl=1h

# 1.2.4 Alert IGC + HSM vendor (Thales support contract)
# 1.2.5 Order replacement HSM (expedite 24h if available, otherwise 72h)
```

### 1.3 Root cause & replacement (H+1h to H+4h)
- [ ] Technician enters salle blanche (IGC + guard present)
- [ ] Visual inspection HSM-CORE-01 (power, fans, LEDs, temperature)
- [ ] If recoverable (power supply, network cable) → fix, restart, verify
- [ ] If hardware failure (motherboard, HSM module) → replace unit
- [ ] Replacement HSM:
  - [ ] Verify firmware hash matches approved version
  - [ ] Verify attestation certificate from Thales
  - [ ] Zeroize before use (factory reset)
  - [ ] Restore from DR replication (partition-by-partition)
  - [ ] OR restore from backup blob + Shamir if replication unavailable
- [ ] Test all partitions: sign test data, verify, validate M-of-N
- [ ] Return HSM-CORE-02 to standby, HSM-CORE-01 primary (or vice versa)

---

## HSM-004 : Tamper Detection (PHYSICAL COMPROMISE)

### 4.1 Detection
```bash
# HSM auto-zeroizes (keys destroyed internally)
# Alert: hsm_tamper_status{hsm="HSM-CORE-01"} == 1
# Salle blanche intrusion alarm active
```

### 4.2 IMMEDIATE — H+0 to H+5min (CATASTROPHE)
```bash
# 4.2.1 ACTIVER CELLULE CRISE PKI (voir root-ca-compromise.md Phase 1)
# 4.2.2 Gardes armés → salle blanche. Personne n'entre ni ne sort.
# 4.2.3 Biométrie salle blanche → check last 24h access logs
# 4.2.4 Vidéo surveillance → live + last 72h extraction
# 4.2.5 Isolation réseau Cilium PKI (emergency lockdown policy)

# 4.2.6 Determine: accidental (environmental) or malicious (intrusion)
# If environmental (fire suppression triggered, water leak, temperature spike):
#   → HSM-001 hardware replacement procedure
# If malicious (broken door, unauthorized access, device planted):
#   → ROOT CA COMPROMISE PROCEDURE (root-ca-compromise.md)
```

### 4.3 Post-tamper (if accidental)
- [ ] HSM zeroized = keys destroyed. Cannot recover HSM-CORE-01.
- [ ] Activate HSM-CORE-02 (standby) immediately.
- [ ] If HSM-CORE-02 also affected (same environmental cause) → activate HSM-DR-01.
- [ ] If all HSMs Core destroyed → proceed to Root CA reconstruction.
- [ ] If HSM-CORE-02 healthy:
  - [ ] Verify no tamper on HSM-CORE-02 (visual, seal, video)
  - [ ] Verify partitions intact (sign test)
  - [ ] Order 2 replacement HSMs (new primary + new standby)
  - [ ] Reconstruct HSM-CORE-01 replacement from DR replication (NOT from Shamir unless necessary)
  - [ ] Full audit IGC + forensics + report Présidence

---

## HSM-005 : Network Isolation

### 5.1 Diagnosis
```bash
# HSM pings, but Vault cannot connect via PKCS#11
# Check: switch port, VLAN config, firewall nftables, Cilium policy

# 5.1.1 Check network path
hping3 -S -p 1792 10.1.0.10
# If SYN but no ACK → firewall / Cilium / switch ACL

# 5.1.2 Check Cilium policies
cilium bpf policy get | grep 10.1.0.10
kubectl get cnp -A | grep hsm

# 5.1.3 Check nftables host firewall
nft list ruleset | grep 1792

# 5.1.4 Check switch port status
ssh switch-admin@core-switch "show interface status | include HSM"
```

### 5.2 Resolution
- [ ] If Cilium policy error → restore valid policy, verify connectivity
- [ ] If nftables rule error → restore known-good ruleset from backup
- [ ] If switch failure → failover to redundant switch, verify VLAN 20
- [ ] If network card HSM → HSM-001 hardware failure
- [ ] If DDoS / flood → activate DDoS mitigation, rate-limit, SOC investigation

---

## HSM DR Sync Failure (HSM-006)

### 6.1 Detection
```bash
# Alert: hsm_dr_sync_lag_seconds{hsm="HSM-DR-01"} > 3600
# Or: partitions on DR missing recent keys
```

### 6.2 Actions
```bash
# 6.2.1 Check DR network connectivity
ping 10.2.0.10

# 6.2.2 Check DR HSM health
vault read -format=json sys/health  # Via DR bastion

# 6.2.3 If DR HSM healthy but sync broken:
# Trigger manual replication from Core HSM
lunash:> partition sync -partition SNISID_GOV_CA -source 10.1.0.10 -destination 10.2.0.10
lunash:> partition sync -partition SNISID_CITIZEN_CA -source 10.1.0.10 -destination 10.2.0.10
lunash:> partition sync -partition SNISID_INFRA_CA -source 10.1.0.10 -destination 10.2.0.10
lunash:> partition sync -partition SNISID_JUDICIAL_CA -source 10.1.0.10 -destination 10.2.0.10
lunash:> partition sync -partition SNISID_DEVICE_CA -source 10.1.0.10 -destination 10.2.0.10
# Note: Root CA partition is NEVER synced automatically (offline)

# 6.2.4 Verify sync integrity
# Sign test data on DR HSM, verify against Core public key
```

---

## Post-recovery

- [ ] Full audit log HSM (Core + DR) reviewed by IGC
- [ ] All CA operations validated (issue test certs for each intermediate)
- [ ] OCSP/CRL operations validated
- [ ] Digital Signature Service test signing
- [ ] Vault PKI unseal + operations validated
- [ ] Report to Présidence (if tamper or hardware failure)
- [ ] Vendor RMA processed (if hardware defect)
- [ ] Spare HSM inventory checked (minimum 1 cold spare on-site)

---

*Testé trimestriellement via simulation HSM failover sur environnement staging.*
