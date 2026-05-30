#!/bin/bash
HOST=$1
echo "=== Rapport Hardening Physique SNISID - $HOST ==="
SB=$(ssh root@$HOST "mokutil --sb-state 2>/dev/null")
[[ "$SB" == *"enabled"* ]] && echo "OK Secure Boot: ACTIF" || echo "ERREUR Secure Boot: INACTIF"
TPM=$(ssh root@$HOST "ls /dev/tpm0 2>/dev/null && echo present || echo absent")
[[ "$TPM" == "present" ]] && echo "OK TPM 2.0: PRESENT" || echo "ERREUR TPM 2.0: ABSENT"
LUKS=$(ssh root@$HOST "lsblk -o FSTYPE | grep crypto_LUKS | wc -l")
[[ "$LUKS" -gt 0 ]] && echo "OK LUKS FDE: $LUKS volume(s) chiffres" || echo "ERREUR LUKS FDE: AUCUN"
IOMMU=$(ssh root@$HOST "dmesg | grep -ci 'iommu enabled'")
[[ "$IOMMU" -gt 0 ]] && echo "OK IOMMU: ACTIF" || echo "ERREUR IOMMU: INACTIF"
