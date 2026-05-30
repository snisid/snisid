# SNISID RUNBOOK — Certificate Rotation (PKI Nationale)
**Classification:** TOP-SECRET  
**Version:** 4.0.0  
**RTO cible:** < 2 heures  
**Fréquence test:** Semestriel

---

## 1. PKI Nationale — Architecture

| Composant | Type | Durée | Renouvellement |
|-----------|------|-------|----------------|
| Root CA SNISID | Offline HSM | 10 ans | Cérémonie physique IGC |
| Intermediate CA | Vault PKI | 5 ans | Cérémonie physique IGC |
| Mesh/Service certs | Vault PKI | 90 jours | Auto cert-manager |
| Ingress certs | Vault PKI | 90 jours | Auto cert-manager |
| Admin client certs | Vault PKI | 30 jours | Auto + renouvellement manuel |
| Edge node certs | Vault PKI | 7 jours | Auto (courtes pour sécurité terrain) |

## 2. Rotation automatique (cert-manager + Vault)

```bash
# Vérifier l'état des certificats
kubectl get certificates -A
kubectl get certificaterequests -A

# Vérifier que cert-manager est sain
kubectl get pods -n cert-manager

# Trigger manuel d'un renouvellement anticipé (urgence compromission)
kubectl annotate certificate snisid-ingress-wildcard -n istio-system \
  cert-manager.io/next-private-key="$(openssl rand -hex 32)"
# Ou plus propre :
cmctl renew snisid-ingress-wildcard -n istio-system
```

## 3. Rotation d'urgence — Compromission suspectée

### 3.1 Isolation immédiate
```bash
# Révoquer toutes les issued certs sous l'intermediate compromise
vault write snisid-pki-int/revoke -format=json \
  serial_number=$(vault read -field=serial snisid-pki-int/cert/ca)

# Ou CRL bulk si compromission massive
vault write snisid-pki-int/config/crl \
  expiry="72h" \
  disable=false

# Notifier SOC + IGC (incident TOP-SECRET)
```

### 3.2 Nouvelle intermediate CA
```bash
# Cérémonie en salle blanche (3 agents IGC minimum, HSM Thales)
# 1. Générer nouvelle clé intermediate sur HSM (PKCS#11)
# 2. Signer avec Root CA offline (air-gapped laptop IGC)
vault write -format=json snisid-pki/root/sign-intermediate \
  csr=@/secure/new-intermediate.csr \
  common_name="SNISID Intermediate CA v$(date +%Y)" \
  ttl="43800h"  # 5 ans

# 3. Installer nouvelle chaîne
cat new-intermediate.crt root-ca.crt > ca-chain.crt
kubectl create secret generic snisid-pki-ca-new \
  --from-file=ca.crt=ca-chain.crt \
  --from-file=tls.crt=new-intermediate.crt \
  --from-file=tls.key=/dev/null  # Clé jamais hors HSM

# 4. Rolling update cert-manager issuers
kubectl patch clusterissuer snisid-vault-pki-issuer \
  --type=merge -p '{"spec":{"vault":{"path":"snisid-pki-int-v2/sign/cluster-local"}}}'

# 5. Forcer renouvellement massif
kubectl get certificates -A -o json | \
  jq -r '.items[] | [.metadata.namespace,.metadata.name] | @tsv' | \
  while read ns name; do
    cmctl renew "${name}" -n "${ns}"
  done

# 6. Vérification mesh-wide
istioctl proxy-status | grep SYNCED  # Tous les sidecars doivent recharger
```

### 3.3 Post-rotation
- [ ] Mettre à jour le `caBundle` dans tous les webhooks (Kyverno, ArgoCD, etc.)
- [ ] Mettre à jour les trust stores des edge nodes (air-gap bundle rebuild)
- [ ] Mettre à jour les laptops admin IGC
- [ ] Détruire l'ancienne chaîne (shred physiques, suppression Vault path)
- [ ] Rapport IGC + présidence

## 4. Rotation Root CA (tous les 10 ans — événement majeur)

> **Cérémonie nationale.** 6 mois de préparation. Toutes les identités nationales re-signées.

```
Phase 1 (M-6): Préparation salle blanche, HSM neuf, audit IGC externe
Phase 2 (M-3): Génération nouvelle root key (HSM), CRL de l'ancienne préparée
Phase 3 (M-1): Distribution nouvelles trust anchors (edge, mobile, admin, diplomatique)
Phase 4 (M-0): Activation nouvelle root, cross-sign 30j
Phase 5 (M+1): Révocation ancienne root, mise à jour bundles nationaux
Phase 6 (M+3): Validation totale, fermeture cérémonie
```

---

*Test semestriel sur environnement isolé. Root CA testée annuellement.*
