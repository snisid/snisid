#!/bin/bash
# SNISID Airgap & K3s First Boot Initialization
# Executed automatically by systemd upon first successful LUKS unlock

set -e

LOG_FILE="/var/log/snisid-firstboot.log"
exec > >(tee -a ${LOG_FILE}) 2>&1

echo "[*] Starting SNISID Sovereign Initialization (Airgap Mode)"

# 1. Enforce Airgap Posture (Requirement #129)
echo "[*] Dropping default routes to external internet..."
# Find the primary interface (excluding loopback)
PRIMARY_IF=$(ip route show default | awk '/default/ {print $5}')
if [ -n "$PRIMARY_IF" ]; then
    ip route del default
    echo "[!] Network isolated. Internal routing only."
fi

# 2. Offline K3s Installation
echo "[*] Bootstrapping K3s cluster from offline assets..."
# The iso build process placed the binary and images in /opt/snisid-offline
K3S_BIN="/opt/snisid-offline/k3s"
INSTALL_SCRIPT="/opt/snisid-offline/install.sh"
IMAGES_DIR="/var/lib/rancher/k3s/agent/images/"

if [ -f "$INSTALL_SCRIPT" ] && [ -f "$K3S_BIN" ]; then
    chmod +x $K3S_BIN
    cp $K3S_BIN /usr/local/bin/
    
    # Run the installer telling it to skip downloading binaries
    INSTALL_K3S_SKIP_DOWNLOAD=true INSTALL_K3S_EXEC="server --disable traefik --disable servicelb" $INSTALL_SCRIPT
    echo "[*] K3s Installation complete."
else
    echo "[!] Offline assets missing. Skipping K3s bootstrap."
fi

# 3. Wait for K3s to be ready
echo "[*] Waiting for Kubernetes API..."
until /usr/local/bin/kubectl get node &> /dev/null; do
    sleep 2
done

# 4. Load SNISID Baseline Manifests
echo "[*] Applying SNISID Core Manifests..."
# Apply the network policies (Cilium), Vault, Keycloak, etc.
if [ -d "/opt/snisid-offline/manifests" ]; then
    /usr/local/bin/kubectl apply -f /opt/snisid-offline/manifests/
fi

echo "[*] First boot initialization completed in $(cat /proc/uptime | awk '{print $1}') seconds."

# 5. Disable this script from running again
systemctl disable snisid-firstboot.service
