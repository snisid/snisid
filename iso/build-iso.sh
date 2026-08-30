#!/bin/bash
# SNISID ISO Master Builder
# Packages the custom user-data, offline assets, and boot configurations into a sovereign Ubuntu ISO

set -e

UBUNTU_ISO="ubuntu-22.04.3-live-server-amd64.iso"
OUTPUT_ISO="snisid-sovereign-appliance.iso"
WORKING_DIR="./iso-build"

echo "[*] Checking prerequisites..."
command -v xorriso >/dev/null 2>&1 || { echo >&2 "I require xorriso but it's not installed. Aborting."; exit 1; }

echo "[*] Preparing build environment..."
mkdir -p "$WORKING_DIR"
# Assuming UBUNTU_ISO is already downloaded locally
if [ ! -f "$UBUNTU_ISO" ]; then
    echo "[!] $UBUNTU_ISO not found. Please download it first."
    exit 1
fi

echo "[*] Extracting base ISO..."
xorriso -osirrox on -indev "$UBUNTU_ISO" -extract / "$WORKING_DIR/ubuntu-base"
chmod -R u+w "$WORKING_DIR/ubuntu-base"

echo "[*] Injecting Autoinstall (Preseed) Configuration..."
mkdir -p "$WORKING_DIR/ubuntu-base/nocloud/"
cp ./autoinstall/user-data.yaml "$WORKING_DIR/ubuntu-base/nocloud/user-data"
touch "$WORKING_DIR/ubuntu-base/nocloud/meta-data"

echo "[*] Injecting Offline K3s Assets & Scripts..."
mkdir -p "$WORKING_DIR/ubuntu-base/scripts"
cp ./scripts/first-boot-init.sh "$WORKING_DIR/ubuntu-base/scripts/"

# Modify GRUB to point to our nocloud data source
sed -i 's/---/autoinstall ds=nocloud;s=\/cdrom\/nocloud\/ ---/' "$WORKING_DIR/ubuntu-base/boot/grub/grub.cfg"

echo "[*] Repacking Sovereign ISO..."
cd "$WORKING_DIR/ubuntu-base"
xorriso -as mkisofs -r -V "SNISID_INSTALL" \
  -J -l -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot \
  -boot-load-size 4 -boot-info-table -eltorito-alt-boot \
  -e boot/grub/efi.img -no-emul-boot -isohybrid-gpt-basdat -isohybrid-apm-hfsplus \
  -o "../../$OUTPUT_ISO" .

cd ../../
echo "[*] Generating SHA-256 Signature (Requirement #131)..."
sha256sum "$OUTPUT_ISO" > "$OUTPUT_ISO.sha256"

echo "[*] Build Complete. Output: $OUTPUT_ISO"
cat "$OUTPUT_ISO.sha256"
