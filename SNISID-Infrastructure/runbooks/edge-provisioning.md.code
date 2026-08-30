# SNISID RUNBOOK — Edge Node Provisioning (Regional / Mobile / Offline)
**Classification:** RESTREINT DEFENSE  
**Version:** 4.0.0  
**RTO cible:** < 30 minutes par nœud  
**Fréquence:** À chaque déploiement terrain

---

## 1. Types de provisioning

| Type | Méthode | Connectivité | Durée |
|------|---------|------------|-------|
| Regional | PXE + Cloud-init + Vault auto-join | Fibre/MPLS | ~15 min |
| Mobile | USB boot + Ansible + Vault token physique | 4G/5G/SAT | ~20 min |
| Offline | USB bundle complet (air-gap) | Aucune | ~25 min |
| Emergency | SD card pré-imagée + solar kit | Burst satellite | ~10 min |

## 2. Provisioning Regional (standard)

### 2.1 Préparation centralisée
```bash
# Génération token K3s edge (valide 24h, usage unique)
VAULT_TOKEN=$(vault read -field=token snisid/edge/tokens/regional-01)
K3S_TOKEN=$(k3s token create --ttl=24h --description="node-regional-01-$(date +%s)")

# Encodage GPG pour transport sécurisé
echo "${K3S_TOKEN}" | gpg --encrypt --armor --recipient snisid-edge-regional-01 > /secure/tokens/regional-01.token.gpg
```

### 2.2 PXE Boot
```bash
# Sur serveur PXE national (Core DC)
# TFTP + iPXE + kickstart cloud-init
# Image: Ubuntu 22.04 CIS-hardened + containerd + K3s agent

# Post-boot cloud-init
cat << 'EOF' > /var/lib/cloud/scripts/per-once/01-k3s-join.sh
#!/bin/bash
GPG_TOKEN=/secure/tokens/$(hostname -s).token.gpg
K3S_TOKEN=$(gpg --decrypt "${GPG_TOKEN}" 2>/dev/null)
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.28.6+k3s2 \
  K3S_URL=https://api.edge-regional-01.snisid.gouv.local:6443 \
  K3S_TOKEN="${K3S_TOKEN}" \
  INSTALL_K3S_EXEC="agent --node-label snisid.gov/region=regional-01 --node-taint snisid.gov/tier=tier-4:NoSchedule" \
  sh -
EOF
```

### 2.3 Validation
```bash
# Bastion national
ssh edge-regional-01 "sudo k3s kubectl get nodes"
ssh edge-regional-01 "sudo systemctl status k3s-agent"
# Vérifier labels, taints, Cilium connectivity test
```

## 3. Provisioning Mobile (unité terrain)

### 3.1 Kit de déploiement

Contenu du **SNISID Mobile Deployment Kit** (mallette sécurisée niveau 2):
- 1x Laptop durci (admin provisioning, LUKS boot)
- 2x Clés USB 3.2 chiffrées (OS image + secrets)
- 1x Token Vault physique (YubiKey HSM, PIN + biométrie opérateur)
- 1x Imprimante thermique (journal de provisioning, papier continu)
- 1x Kit solaire/batterie (72h autonomie)

### 3.2 Procédure sur site

```bash
# 1. Authentification opérateur
vault login -method=cert -client-cert=/secure/op-$(id -u).crt
vault read -field=provision-auth snisid/operators/$(whoami)

# 2. Boot laptop admin depuis USB secure
# BIOS vérifié (Secure Boot, TPM2.0, PCR7 validé)

# 3. Connexion nœud mobile (câble Ethernet direct)
ansible-playbook -i inventory/mobile.yml provision-mobile.yml \
  --extra-vars "target=mobile-node-01 vault_token=${VAULT_TOKEN}"

# 4. K3s agent join (via Bastion VPN satellite)
# Note: Satellite link activé uniquement pour join initial, puis standby
```

### 3.3 Ansible playbook (extrait)
```yaml
# provision-mobile.yml
- hosts: mobile_nodes
  become: yes
  vars:
    k3s_version: "v1.28.6+k3s2"
    registry_mirror: "registry-mobile-cache.snisid.gouv.local"
  tasks:
    - name: Harden OS (CIS Level 2)
      include_role:
        name: snisid-hardening
    
    - name: Install containerd + K3s airgap
      include_role:
        name: k3s-airgap
    
    - name: Configure Cilium agent (join cluster mesh)
      template:
        src: cilium-mobile-config.yaml.j2
        dest: /etc/cilium/config.yaml
    
    - name: Inject Vault agent token (wrapped, single-use)
      command: >
        vault write -field=wrapping_token 
        auth/token/create 
        policies=mobile-edge-policy 
        ttl=72h 
        num_uses=10
      register: vault_wrapped_token
      delegate_to: bastion
    
    - name: Start Vault agent (cache mode)
      systemd:
        name: vault-agent
        state: started
        enabled: yes
```

## 4. Provisioning Offline (zone sans réseau)

### 4.1 Bundle USB
Voir `../offline/k3s-offline-bundle.md` pour la construction du bundle.

### 4.2 Procédure
```bash
# Sur site (opérateur authentifié biométrie + smartcard)
# 1. Insérer USB chiffrée, déverrouiller avec token Vault mobile
./offline-verify-bundle.sh /dev/sda1

# 2. Installation automatique K3s + manifests
./offline-install.sh --target /dev/nvme0n1 --bundle /mnt/snisid-sync/

# 3. Validation
k3s kubectl get nodes
k3s kubectl get pods -A

# 4. Retrait USB, scellement nœud (tamper-evident stickers)
# 5. Journal papier signé opérateur + timestamp GPS
```

## 5. Post-provisioning obligatoire

- [ ] Nœud visible dans `kubectl get nodes` avec labels corrects
- [ ] Cilium CNI connecté (pas de pods Pending > 5 min)
- [ ] Vault Agent injecte secrets (test: lire un secret test)
- [ ] Falco runtime actif (vérifier logs `falco-events`)
- [ ] Prometheus scrape actif (metric endpoint accessible)
- [ ] Local storage test (write/read Ceph RBD edge ou local-path)
- [ ] Sync initial terrain → Core validé (Kafka test message)
- [ ] Rapport provisioning signé dans SIEM national

---

*Tout provisioning non documenté est considéré comme incident de sécurité.*
