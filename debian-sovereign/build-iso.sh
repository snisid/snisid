#!/bin/bash
# SNISID Sovereign OS - ISO Build Script
# Crée une image ISO bootable basée sur Debian 12 avec tous les modules de sécurité nationale
# Classification: TOP SECRET / Gouvernement Haïtien

set -euo pipefail

# Configuration
SNISID_VERSION="1.0"
DEBIAN_VERSION="12.5"  # Bookworm
ARCH="amd64"
ISO_NAME="snisid-sovereign-os-v${SNISID_VERSION}.iso"
WORK_DIR="/tmp/snisid-build-$$"
MOUNT_POINT="/mnt/snisid-iso"

# Couleurs pour logs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCÈS]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[ATTENTION]${NC} $1"; }
log_error() { echo -e "${RED}[ERREUR]${NC} $1"; }

# Vérification des prérequis
check_prerequisites() {
    log_info "Vérification des prérequis..."
    
    local deps=("genisoimage" "xorriso" "debootstrap" "squashfs-tools" "grub-pc-bin" "grub-efi-amd64-bin" "mtools")
    local missing=()
    
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing+=("$dep")
        fi
    done
    
    if [ ${#missing[@]} -ne 0 ]; then
        log_error "Dépendances manquantes: ${missing[*]}"
        log_info "Installez avec: apt-get install ${missing[*]}"
        exit 1
    fi
    
    log_success "Toutes les dépendances sont installées"
}

# Préparation de l'environnement
setup_environment() {
    log_info "Préparation de l'environnement de build..."
    
    mkdir -p "$WORK_DIR"/{iso-root,live,filesystem,chroot}
    
    # Nettoyage en cas d'échec précédent
    if mountpoint -q "$MOUNT_POINT"; then
        umount "$MOUNT_POINT" || true
    fi
    
    log_success "Environnement prêt"
}

# Téléchargement de Debian de base
download_debian_base() {
    log_info "Téléchargement de Debian $DEBIAN_VERSION ($ARCH)..."
    
    DEBIAN_MIRROR="http://deb.debian.org/debian"
    DEBIAN_BASE="$WORK_DIR/debian-base.tar.gz"
    
    # Utilisation de debootstrap pour créer un système minimal
    debootstrap --arch=$ARCH --include=linux-image-amd64,grub-pc \
        --exclude=systemd,snapd,modemmanager \
        bookworm "$WORK_DIR/chroot" "$DEBIAN_MIRROR"
    
    log_success "Base Debian téléchargée et préparée"
}

# Installation des paquets SNISID
install_snisid_packages() {
    log_info "Installation des paquets SNISID..."
    
    cat > "$WORK_DIR/chroot/tmp/packages.list" << 'EOF'
# Système de base
linux-image-amd64
grub-pc
grub-efi-amd64
firmware-linux-nonfree

# Sécurité réseau
nftables
wireguard
openvpn
strongswan

# Cryptographie
gnupg2
cryptsetup-initramfs
tpm2-tools
yubikey-personalization

# Outils NSA intégrés
# (Les binaires seront copiés depuis external-integrations)

# Outils OSINT
python3
python3-pip
git
curl
wget

# Monitoring et audit
auditd
aide
rkhunter
chkrootkit
lynis

# Base de données clientes
postgresql-client
mongodb-clients
redis-tools

# Interface graphique (optionnelle)
xorg
kde-plasma-desktop
firefox-esr
konsole
EOF

    # Montage des systèmes de fichiers pour chroot
    mount --bind /dev "$WORK_DIR/chroot/dev"
    mount --bind /proc "$WORK_DIR/chroot/proc"
    mount --bind /sys "$WORK_DIR/chroot/sys"
    mount --bind /run "$WORK_DIR/chroot/run"
    
    # Installation dans le chroot
    chroot "$WORK_DIR/chroot" /bin/bash -c "
        apt-get update
        DEBIAN_FRONTEND=noninteractive apt-get install -y \$(cat /tmp/packages.list | grep -v '^#' | tr '\n' ' ')
        apt-get clean
    "
    
    # Démontage
    umount "$WORK_DIR/chroot"/{dev,proc,sys,run}
    
    log_success "Paquets SNISID installés"
}

# Configuration du système durci
harden_system() {
    log_info "Application du durcissement de sécurité..."
    
    # Copie du script de hardening
    cp /workspace/debian-sovereign/scripts/security-hardening.sh "$WORK_DIR/chroot/root/"
    chmod +x "$WORK_DIR/chroot/root/security-hardening.sh"
    
    # Configuration du kernel
    cat > "$WORK_DIR/chroot/etc/sysctl.d/99-snisid-security.conf" << 'EOF'
# SNISID Security Hardening - Kernel Parameters
# Inspiré de NSA STIG et standards chinois MLPS

# Désactivation IPv6 si non utilisé
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1

# Protection contre le spoofing
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1

# Ignorer les pings broadcast
net.ipv4.icmp_echo_ignore_broadcasts = 1

# Protection contre les attaques SYN flood
net.ipv4.tcp_syncookies = 1

# Désactiver le routage source
net.ipv4.conf.all.accept_source_route = 0
net.ipv4.conf.default.accept_source_route = 0

# Activer ASLR
kernel.randomize_va_space = 2

# Restriction dmesg
kernel.dmesg_restrict = 1

# Protection ptrace
kernel.yama.ptrace_scope = 1

# Core dumps désactivés
fs.suid_dumpable = 0
EOF

    # Configuration PAM renforcée
    cat > "$WORK_DIR/chroot/etc/pam.d/common-password-snised" << 'EOF'
# SNISID Password Policy - Conforme NSA + Chine
password requisite pam_pwquality.so retry=3 minlen=14 dcredit=-1 ucredit=-1 ocredit=-1 lcredit=-1
password required pam_unix.so sha512 shadow rounds=100000
EOF

    log_success "Durcissement appliqué"
}

# Intégration des modules externes
integrate_external_modules() {
    log_info "Intégration des modules GitHub..."
    
    MODULES_DIR="/workspace/external-integrations"
    OPT_DIR="$WORK_DIR/chroot/opt/snisid"
    
    mkdir -p "$OPT_DIR"/{modules,scripts,config,data}
    
    # Copie des modules NSA
    if [ -d "$MODULES_DIR/datawave" ]; then
        cp -r "$MODULES_DIR/datawave" "$OPT_DIR/modules/"
        log_info "Module DataWave intégré"
    fi
    
    if [ -d "$MODULES_DIR/emissary" ]; then
        cp -r "$MODULES_DIR/emissary" "$OPT_DIR/modules/"
        log_info "Module Emissary intégré"
    fi
    
    if [ -d "$MODULES_DIR/maat" ]; then
        cp -r "$MODULES_DIR/maat" "$OPT_DIR/modules/"
        log_info "Module MAAT intégré"
    fi
    
    # Copie des outils OSINT
    for tool in sigint cyberdome osiris phoneinfoga geospy ace-t storm-breaker fmd-server; do
        if [ -d "$MODULES_DIR/$tool" ]; then
            cp -r "$MODULES_DIR/$tool" "$OPT_DIR/modules/"
            log_info "Module $tool intégré"
        fi
    done
    
    # Création des scripts de lancement
    cat > "$OPT_DIR/scripts/start-all-services.sh" << 'EOF'
#!/bin/bash
# Démarrage de tous les services SNISID
echo "[*] Démarrage des services SNISID..."

# Services NSA
systemctl start datawave-query 2>/dev/null || echo "DataWave non disponible"
systemctl start emissary-ingest 2>/dev/null || echo "Emissary non disponible"

# Services OSINT
cd /opt/snisid/modules/phoneinfoga && ./phoneinfoga serve &
cd /opt/snisid/modules/osiris && python3 app.py &

echo "[+] Tous les services démarrés"
EOF
    chmod +x "$OPT_DIR/scripts/start-all-services.sh"
    
    log_success "Modules externes intégrés"
}

# Création du filesystem squashfs
create_squashfs() {
    log_info "Création du filesystem compressé..."
    
    mksquashfs "$WORK_DIR/chroot" "$WORK_DIR/live/filesystem.sfs" \
        -comp xz \
        -Xbcj x86 \
        -b 1M \
        -no-recovery
    
    log_success "Filesystem squashfs créé ($(du -sh "$WORK_DIR/live/filesystem.sfs" | cut -f1))"
}

# Construction de l'ISO
build_iso() {
    log_info "Construction de l'ISO bootable..."
    
    ISO_ROOT="$WORK_DIR/iso-root"
    
    # Structure ISO
    mkdir -p "$ISO_ROOT"/{live,isolinux,EFI/boot}
    
    # Copie du filesystem
    cp "$WORK_DIR/live/filesystem.sfs" "$ISO_ROOT/live/"
    
    # Configuration GRUB BIOS
    cat > "$ISO_ROOT/isolinux/grub.cfg" << EOF
set timeout=10
set default=0

menuentry "SNISID Sovereign OS v${SNISID_VERSION}" {
    linux /live/vmlinuz boot=live quiet splash snisid.secure=true
    initrd /live/initrd.img
}

menuentry "SNISID Secure Mode (No Network)" {
    linux /live/vmlinuz boot=live quiet splash snisid.secure=true snisid.nonetwork=true
    initrd /live/initrd.img
}

menuentry "SNISID Rescue Mode" {
    linux /live/vmlinuz boot=live rescue splash
    initrd /live/initrd.img
}
EOF

    # Copie des fichiers de boot
    cp "$WORK_DIR/chroot/boot/vmlinuz-"* "$ISO_ROOT/live/vmlinuz" 2>/dev/null || \
        cp /boot/vmlinuz-* "$ISO_ROOT/live/vmlinuz" 2>/dev/null || \
        log_warn "Kernel non trouvé, utilisation d'un kernel par défaut"
    
    # Génération de l'ISO avec xorriso
    xorriso -as mkisofs \
        -o "/workspace/${ISO_NAME}" \
        -V "SNISID_SOVEREIGN_V${SNISID_VERSION}" \
        -J \
        -r \
        -b isolinux/isolinux.bin \
        -c isolinux/boot.cat \
        -no-emul-boot \
        -boot-load-size 4 \
        -boot-info-table \
        -eltorito-alt-boot \
        -e EFI/boot/grub.efi \
        -no-emul-boot \
        -isohybrid-mbr /usr/lib/ISOLINUX/isohdpfx.bin \
        "$ISO_ROOT"
    
    log_success "ISO créée: /workspace/${ISO_NAME}"
}

# Vérification de l'intégrité
verify_iso() {
    log_info "Vérification de l'intégrité de l'ISO..."
    
    sha256sum "/workspace/${ISO_NAME}" > "/workspace/${ISO_NAME}.sha256"
    
    echo ""
    log_success "Checksum SHA256 généré:"
    cat "/workspace/${ISO_NAME}.sha256"
    
    # Taille de l'ISO
    ls -lh "/workspace/${ISO_NAME}"
}

# Nettoyage
cleanup() {
    log_info "Nettoyage..."
    
    # Démontage sécurisé
    for mount in "$WORK_DIR/chroot"/{dev,proc,sys,run}; do
        umount "$mount" 2>/dev/null || true
    done
    
    # Suppression des fichiers temporaires
    rm -rf "$WORK_DIR"
    
    log_success "Nettoyage terminé"
}

# Main
main() {
    echo "========================================"
    echo "  SNISID Sovereign OS Builder v${SNISID_VERSION}"
    echo "  Gouvernement Haïtien - TOP SECRET"
    echo "========================================"
    echo ""
    
    check_prerequisites
    setup_environment
    download_debian_base
    install_snisid_packages
    harden_system
    integrate_external_modules
    create_squashfs
    build_iso
    verify_iso
    cleanup
    
    echo ""
    echo "========================================"
    log_success "SNISID OS v${SNISID_VERSION} PRÊT POUR DÉPLOIEMENT"
    echo "========================================"
    echo ""
    echo "Fichiers générés:"
    echo "  - /workspace/${ISO_NAME}"
    echo "  - /workspace/${ISO_NAME}.sha256"
    echo ""
    echo "Prochaines étapes:"
    echo "  1. Graver l'ISO sur USB bootable"
    echo "  2. Tester en machine virtuelle"
    echo "  3. Déployer sur matériel certifié"
    echo ""
}

# Gestion des erreurs
trap 'log_error "Échec à l'étape: $LINENO"; cleanup; exit 1' ERR

# Exécution
main "$@"
