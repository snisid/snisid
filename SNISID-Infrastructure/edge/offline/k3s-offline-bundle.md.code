# SNISID — K3s Offline Bundle : Emergency & Isolated Zones
**Classification:** RESTREINT DEFENSE  
**Scope:** Zones sans connectivité réseau (rural, montagne, catastrophe)  
**Autonomie:** 7 jours offline complet

---

## 1. Offline Node Architecture

```
[Emergency Offline Node]
├── K3s single-node (SQLite)
├── containerd with pre-loaded images
├── Local Vault Agent (token cache, pas de master key)
├── Local PostgreSQL (réplica partiel identités)
├── Local Kafka (buffer sortant, 7j retention)
├── Biometric SDK offline
├── Synchronisation : USB/SD chiffrés (LUKS + GPG national)
└── Power : batteries/solaire (72h-7j)
```

## 2. Air-gap Image Bundle

```bash
#!/bin/bash
# SNISID — Script de construction du bundle offline
# Exécuté sur bastion national avant déploiement terrain

OFFLINE_DIR="/mnt/secure/snisid-offline-bundle"
REGISTRY="registry.interne.snisid.gouv.local"
IMAGES=(
  "${REGISTRY}/rancher/k3s:v1.28.6-k3s2"
  "${REGISTRY}/cilium/cilium:v1.15.0"
  "${REGISTRY}/cilium/operator-generic:v1.15.0"
  "${REGISTRY}/vault:1.15.4"
  "${REGISTRY}/postgres:15.5-alpine-hardened"
  "${REGISTRY}/bitnami/kafka:3.5.1"
  "${REGISTRY}/snisid/biometric-agent:v2026.05.25"
  "${REGISTRY}/snisid/api-offline-cache:v2026.05.25"
  "${REGISTRY}/prometheus/node-exporter:v1.7.0"
)

mkdir -p "${OFFLINE_DIR}/images" "${OFFLINE_DIR}/helm" "${OFFLINE_DIR}/manifests"

# Pull & save images
for img in "${IMAGES[@]}"; do
  ctr image pull "${img}"
  name=$(echo "${img}" | tr '/' '_' | tr ':' '_')
  ctr image export "${OFFLINE_DIR}/images/${name}.tar" "${img}"
  sha256sum "${OFFLINE_DIR}/images/${name}.tar" >> "${OFFLINE_DIR}/SHA256SUMS"
done

# Helm charts (packaged)
helm repo add snisid https://charts.interne.snisid.gouv.local
helm pull snisid/k3s-edge-offline --version 4.0.0 -d "${OFFLINE_DIR}/helm/"
helm pull snisid/vault-agent-offline --version 4.0.0 -d "${OFFLINE_DIR}/helm/"

# Manifests K8s (policies réduites offline)
cp ../../security/kyverno/require-labels.yaml "${OFFLINE_DIR}/manifests/"
cp ../../networking/cilium/cilium-clusterwide-policies.yaml "${OFFLINE_DIR}/manifests/"

# Scripts sync
install -m 700 ../sync/offline-sync-usb.sh "${OFFLINE_DIR}/"
install -m 700 ../sync/offline-verify-bundle.sh "${OFFLINE_DIR}/"

# Chiffrement du bundle pour transport physique
# Clé GPG : SNISID_OFFLINE_TRANSPORT_KEY (4096-bit RSA, subkeys matérielles)
gpg --batch --yes --cipher-algo AES256 --compress-algo 0 \
    --symmetric --output "${OFFLINE_DIR}.gpg" "${OFFLINE_DIR}"

# Destruction du bundle clair
shred -u -n 35 -z "${OFFLINE_DIR}"
```

## 3. Synchronisation par média physique sécurisé

### USB / SD Card
- **Format:** LUKS2 + ext4
- **Clé:** Décryptée par Vault Agent offline (token wrapping)
- **Contenu:** Différentiel Kafka + nouvelles identités + mises à jour images
- **Fréquence:** Quotidien (zone rurale) / À chaque passage (zone isolée)
- **Vérification:** Cosign signature des images + SHA256SUMS signé GPG

### Protocole de réception
```bash
# offline-verify-bundle.sh
MOUNTPOINT="/mnt/snisid-sync"
DEVICE="/dev/sda1"  # Auto-detected by udev rule

# Vérification LUKS header integrity
cryptsetup luksDump "${DEVICE}" | grep -q "SNISID_OFFLINE"

# Unlock
VAULT_TOKEN=$(vault read -field=token snisid/offline/media-token)
cryptsetup open --type luks "${DEVICE}" snisid-sync --key-file <(echo -n "${VAULT_TOKEN}" | sha256sum | awk '{print $1}')

mount /dev/mapper/snisid-sync "${MOUNTPOINT}"

# Vérification GPG
gpg --verify "${MOUNTPOINT}/SHA256SUMS.sig" "${MOUNTPOINT}/SHA256SUMS" || exit 1
sha256sum -c "${MOUNTPOINT}/SHA256SUMS" || exit 1

# Import images
cd "${MOUNTPOINT}/images"
for tar in *.tar; do
  ctr image import "${tar}"
done

# Apply manifests
kubectl apply -f "${MOUNTPOINT}/manifests/"

# Sync Kafka outbound buffer (messages pending to Core)
kafka-console-producer ...  # Append to local buffer

# Cleanup
umount "${MOUNTPOINT}"
cryptsetup close snisid-sync
sync
```

## 4. Emergency Node — Spécificités catastrophe

| Paramètre | Valeur |
|-----------|--------|
| Alimentation | Batteries LiFePO4 + panneau solaire 100W |
| Autonomie élec | 72h sans soleil, illimité avec soleil |
| Connectivité | Satellite Iridium / Starlink souverain (si disponible) |
| Compute | Ruggedized box, -20°C à +60°C, IP67 |
| K3s | Single-node, SQLite, pas d'etcd multi-node |
| Identités | Cache LRU 50 000 identités locales |
| Biométrie | Match offline, template chiffré AES-256 |
| Sync | USB chiffré ou satellite burst |

---

*Ce bundle est reconstruit automatiquement chaque semaine par CI/CD national et validé par signature GPG IGC.*
