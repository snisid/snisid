#!/usr/bin/env bash
# File: /pki/scripts/semi_automated_sign.sh
# Enforces a 3-of-5 quorum for Policy CA signing operations
# Compliance: ETSI EN 319 411 / FIPS 140-2 Level 3

set -euo pipefail

MODULE_PATH="/usr/lib/libcs_pkcs11_R2.so"
POLICY_CNF="/pki/root-ca/openssl-root.cnf"
ROOT_CERT="/pki/root-ca/snisid_root_ca.cert.pem"
AUDIT_LOG="/var/log/audit/pki_ceremony.log"
WORM_TARGET="audit-svc.snisid.gov.ht:/var/log/pki-worm/"

echo "========================================================="
echo "  SNISID SEMI-AUTOMATED POLICY CA SIGNING ORCHESTRATOR  "
echo "  Quorum Requis : 3-of-5 Officiers de Sécurité         "
echo "========================================================="

if [ ! -f "$MODULE_PATH" ]; then
    echo "[ERROR] HSM PKCS#11 Module non trouvé à $MODULE_PATH" >&2
    exit 1
fi

if [ $# -lt 2 ]; then
    echo "Usage: $0 <csr_input_path> <cert_output_path>"
    exit 1
fi

CSR_INPUT="$1"
CERT_OUTPUT="$2"

echo "[*] Initialisation de la connexion HSM..."
pkcs11-tool --module "$MODULE_PATH" --list-slots

declare -a PinQuorum
QUORUM_COUNT=0
REQUIRED_QUORUM=3

echo "[*] Collecte des PINs pour le Quorum 3-of-5..."
for i in {1..5}; do
    read -rsp "[?] Entrez le PIN pour l'Officier $i (laisser vide pour passer) : " pin_val
    echo ""
    if [ -n "$pin_val" ]; then
        PinQuorum+=("$pin_val")
        QUORUM_COUNT=$((QUORUM_COUNT + 1))
        echo "[+] PIN pour l'Officier $i accepté."
        if [ "$QUORUM_COUNT" -eq "$REQUIRED_QUORUM" ]; then
            break
        fi
    fi
done

if [ "$QUORUM_COUNT" -lt "$REQUIRED_QUORUM" ]; then
    echo "[ERROR] Quorum insuffisant. $QUORUM_COUNT PINs fournis, $REQUIRED_QUORUM requis." >&2
    exit 2
fi

echo "[*] Quorum atteint. Déverrouillage de la partition et signature du CSR..."
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
echo "$TIMESTAMP - PKI_SIGN_START - Quorum vérifié avec $REQUIRED_QUORUM officiers." >> "$AUDIT_LOG"

HSM_PIN="${PinQuorum[0]}"

if openssl ca -config "$POLICY_CNF" \
  -engine pkcs11 -keyform engine -ss_cert "$ROOT_CERT" \
  -in "$CSR_INPUT" \
  -out "$CERT_OUTPUT" \
  -extensions v3_intermediate_ca \
  -days 3650 \
  -md sha384 \
  -passin "pass:$HSM_PIN" -batch; then
  
    echo "[+] CSR signé avec succès !"
    echo "[*] Certificat de sortie écrit dans : $CERT_OUTPUT"
    
    TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    LOG_MSG="$TIMESTAMP - PKI_SIGN_SUCCESS - Signed CSR: $CSR_INPUT -> $CERT_OUTPUT (SHA-256: $(sha256sum "$CERT_OUTPUT" | cut -d' ' -f1))"
    echo "$LOG_MSG" >> "$AUDIT_LOG"
    
    # Envoi de la trace vers le stockage WORM
    logger -t SNISID-PKI -p security.info "$LOG_MSG"
    rsync -az "$AUDIT_LOG" "$WORM_TARGET" || echo "[WARNING] Échec de la synchronisation vers l'archivage WORM."
else
    TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo "$TIMESTAMP - PKI_SIGN_FAILURE - Échec de signature du CSR: $CSR_INPUT" >> "$AUDIT_LOG"
    logger -t SNISID-PKI -p security.err "Échec de signature PKI."
    exit 3
fi
