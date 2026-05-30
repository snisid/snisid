#!/bin/bash
# === scripts/enable_iommu_snisid.sh ===
# SNISID - Activation IOMMU Protection DMA
# Reference: SNISID-AUDIT-2025 R-008
set -euo pipefail
 
LOG="/var/log/snisid/iommu_$(date +%Y%m%d).log"
mkdir -p "$(dirname $LOG)"
log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG"; }
 
# Detecter CPU
CPU_VENDOR=$(grep -m1 'vendor_id' /proc/cpuinfo | awk '{print $3}')
log "CPU detecte: $CPU_VENDOR"
 
if [[ "$CPU_VENDOR" == "GenuineIntel" ]]; then
    IOMMU_PARAM="intel_iommu=on iommu=pt"
elif [[ "$CPU_VENDOR" == "AuthenticAMD" ]]; then
    IOMMU_PARAM="amd_iommu=on iommu=pt"
else
    log "ERREUR: CPU non supporte: $CPU_VENDOR"; exit 1
fi
 
# Modifier GRUB
GRUB_FILE="/etc/default/grub"
CURRENT=$(grep '^GRUB_CMDLINE_LINUX=' "$GRUB_FILE" | cut -d'"' -f2)
 
if ! echo "$CURRENT" | grep -q "iommu"; then
    NEW_CMD="$CURRENT $IOMMU_PARAM lockdown=confidentiality"
    sed -i "s|^GRUB_CMDLINE_LINUX=.*|GRUB_CMDLINE_LINUX=\"$NEW_CMD\"|" "$GRUB_FILE"
    update-grub
    log "GRUB mis a jour: $IOMMU_PARAM + lockdown=confidentiality"
fi
 
# Regles udev protection DMA
cat > /etc/udev/rules.d/99-snisid-dma.rules << 'EOF'
# SNISID - Bloquer Thunderbolt par defaut
ACTION=="add", SUBSYSTEM=="thunderbolt", ATTR{authorized}="0"
# Alerter si nouveau dispositif PCIe branche
ACTION=="add", SUBSYSTEM=="pci", RUN+="/usr/local/bin/snisid-pci-alert.sh %k"
EOF
 
# Script alerte SIEM
cat > /usr/local/bin/snisid-pci-alert.sh << 'EOF'
#!/bin/bash
MSG="ALERTE SECURITE SNISID: Nouveau PCIe sur $(hostname): $1"
logger -p security.warning -t SNISID_IOMMU "$MSG"
curl -sf -X POST "https://soc.snisid.gouv.ht/api/alerts" \
  -H "Content-Type: application/json" \
  --data '{"severity":"HIGH","message":"'"$MSG"'"}' || true
EOF
chmod +x /usr/local/bin/snisid-pci-alert.sh
 
# Script de verification post-reboot
cat > /usr/local/bin/verify_iommu.sh << 'EOF'
#!/bin/bash
echo "=== Verification IOMMU SNISID ==="
if dmesg | grep -qi "iommu enabled"; then
    echo "OK IOMMU: ACTIF"
    dmesg | grep -i "iommu" | head -3
else
    echo "ERREUR IOMMU: INACTIF - reboot requis ou non supporte par BIOS"
    exit 1
fi
echo "Groupes IOMMU: $(find /sys/kernel/iommu_groups/ -type l | wc -l) groupes"
grep -o 'iommu[^ ]*' /proc/cmdline | head -3
LOCKDOWN=$(cat /sys/kernel/security/lockdown 2>/dev/null || echo "non_supporte")
echo "Kernel lockdown: $LOCKDOWN"
EOF
chmod +x /usr/local/bin/verify_iommu.sh
 
log "Configuration IOMMU terminee. REBOOT REQUIS puis executer: verify_iommu.sh"
log "Test DMA en lab: pcileech --device fpga --cmd probe (apres reboot)"
log "Resultat attendu: echec de l'attaque DMA si IOMMU actif"
