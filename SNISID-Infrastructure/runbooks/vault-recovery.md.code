# SNISID RUNBOOK — Vault Recovery (National Secrets & PKI)
**Classification:** TOP-SECRET  
**Version:** 4.0.0  
**RTO cible:** < 15 minutes  
**Fréquence test:** Mensuel (drill unseal + shamir)

---

## 1. Architecture Vault Nationale

- **Mode:** HA Raft (5 nœuds : 3 Core + 2 DR)
- **Seal:** Auto-unseal HSM Thales Luna 7 (PKCS#11) + Shamir backup (6 parts, threshold 4)
- **Secrets critiques:** Identifiants citoyens, clés biométriques, credentials DB, certificats

## 2. Scénarios

| Code | Scénario | Déclencheur |
|------|----------|-------------|
| VLT-001 | Vault sealed (HSM disconnect) | `vault_core_unsealed == 0` |
| VLT-002 | Raft quorum perdu | `vault_raft_voters < 3` |
| VLT-003 | HSM completement indisponible | Auto-unseal impossible, fallback Shamir |
| VLT-004 | Corruption raft storage | Vault crash-loop, logs `raft index mismatch` |

## 3. Procédure VLT-001 : HSM reconnect + auto-unseal

```bash
# 1. Vérifier état HSM
ssh hsm-admin@hsm-thales-01 "lunash:> partition show -partition SNISID_VAULT"
# Si HSM down : basculer vers HSM DR (hsm-thales-02)

# 2. Vérifier connectivité réseau HSM
nc -zv 10.1.0.10 1792

# 3. Si HSM DR disponible
# Modifier vault config pour pointer HSM DR (déployé via ArgoCD)
kubectl edit configmap vault-config -n snisid-identity
# Changer slot/lib vers HSM DR
kubectl rollout restart statefulset vault -n snisid-identity

# 4. Auto-unseal devrait fonctionner automatiquement après redémarrage
vault status
```

## 4. Procédure VLT-003 : Shamir Recovery (HSM total loss)

> **Cérémonie Shamir :** 6 parts physiques dans 6 coffres distincts (présidence, IGC, SOC, DR-physique, ambassade A, ambassade B).
> Threshold 4 pour reconstruction. JAMAIS moins de 2 personnes présentes par part.

```bash
# 1. Rassembler 4 parts Shamir minimum
# 2. Assembler via CLI Vault (air-gapped laptop, jamais connecté au réseau)
vault operator unseal -tls-skip-verify $(cat part1.key)
vault operator unseal -tls-skip-verify $(cat part2.key)
vault operator unseal -tls-skip-verify $(cat part3.key)
vault operator unseal -tls-skip-verify $(cat part4.key)

# 3. Générer immédiatement un nouveau master key
vault operator generate-root -generate-otp -init

# 4. Reconfigurer auto-unseal avec nouveau HSM (procédure IGC)
# 5. Répartir nouvelles parts Shamir dans nouveaux coffres
```

## 5. Procédure VLT-002 : Raft Quorum Recovery

```bash
# 1. Identifier les voters actifs
vault operator raft list-peers

# 2. Si 2 nœuds perdus (quorum 3/5 = OK, mais risque)
# Si 3 nœuds perdus (quorum rompu) :

# Sur le nœud avec les données les plus récentes :
vault operator raft remove-peer vault-core-04
vault operator raft remove-peer vault-core-05

# Reprovisionner les nœuds via ArgoCD + Helm
# Les nouveaux nœuds rejoignent automatiquement le cluster Raft

# 3. Vérifier
vault operator raft list-peers
vault status
```

## 6. Post-recovery obligatoire

- [ ] Rotation immédiate des root credentials transit
- [ ] Révocation de tous les tokens dynamiques
- [ ] Audit complet des logs Vault (si possible)
- [ ] Notification IGC + Présidence (VLT-003)
- [ ] Mise à jour du disaster recovery register physique

---

*Test mensuel. Tout accès à Vault recovery est loggé et alerté SOC en temps réel.*
