#!/bin/bash
# SNISID JUDAS-PROTOCOL DAEMON
# Version: 1.0 (TOP SECRET)
# Function: Surveillance continue des techniciens et alerte automatique

set -euo pipefail

# Configuration
ALERT_INTERVAL=900 # 15 minutes en secondes
CENTRAL_SERVER="central.snisid.gouv.ht"
SATELLITE_UPLINK="/dev/sat0"
LOG_FILE="/var/log/snisid/judas_watch.log"
EVIDENCE_DIR="/var/evidence/judas"

# Couleurs pour console
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_alert() {
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo -e "${RED}[ALERT]$NC ${timestamp} - $1" | tee -a "$LOG_FILE"
}

log_info() {
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo -e "${GREEN}[INFO]$NC ${timestamp} - $1" | tee -a "$LOG_FILE"
}

# Fonction pour capturer les détails de l'utilisateur
get_user_identity() {
    local user_id="$1"
    
    # Récupération depuis la base centrale (chiffrée)
    local user_data=$(snisid-cli user get "$user_id" --format=json --decrypt)
    
    echo "$user_data" | jq -r '{
        nom: .full_name,
        prenom: .first_name,
        nif: .fiscal_id,
        photo_base64: .biometric_photo,
        role: .role,
        region: .assigned_region
    }'
}

# Fonction pour détecter les tentatives de contournement
detect_breach_attempt() {
    log_info "Démarrage de la surveillance Judas Protocol..."
    
    while true; do
        # Vérifier les tentatives d'élévation de privilèges
        local priv_esc=$(grep -i "privilege escalation\|sudo abuse\|root access denied" /var/log/auth.log | tail -n 1)
        
        # Vérifier les accès hors juridiction
        local jurisdiction_violation=$(snisid-cli audit check-jurisdiction --violations-only)
        
        # Vérifier les exports massifs
        local bulk_export=$(snisid-cli audit check-bulk-export --threshold=5 --last=15min)
        
        # Vérifier tentative modification logs
        local log_tamper=$(grep -i "audit_log.*modified\|history.*deleted" /var/log/audit/audit.log | tail -n 1)
        
        if [[ -n "$priv_esc" || -n "$jurisdiction_violation" || -n "$bulk_export" || -n "$log_tamper" ]]; then
            handle_breach "$priv_esc" "$jurisdiction_violation" "$bulk_export" "$log_tamper"
        fi
        
        sleep "$ALERT_INTERVAL"
    done
}

# Fonction de gestion de brèche
handle_breach() {
    local priv_esc="$1"
    local jur_viol="$2"
    local bulk_exp="$3"
    local log_tamp="$4"
    
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    local alert_id=$(uuidgen)
    
    # Identifier l'utilisateur fautif
    local culprit_user=$(whoami)
    local user_identity=$(get_user_identity "$culprit_user")
    
    log_alert "⚠️  TENTATIVE DE CONTORNNEMENT DÉTECTÉE - ID: $alert_id"
    log_alert "Utilisateur: $(echo $user_identity | jq -r '.nom + " " + .prenom')"
    log_alert "NIF: $(echo $user_identity | jq -r '.nif')"
    log_alert "Rôle: $(echo $user_identity | jq -r '.role')"
    
    # Créer le dossier de preuve
    local evidence_path="$EVIDENCE_DIR/$alert_id"
    mkdir -p "$evidence_path"
    
    # Capturer screenshot et session
    scrot "$evidence_path/screen_evidence.png"
    snisid-cli session record --output="$evidence_path/session_recording.mkv" --last=15min
    
    # Générer le rapport juridique
    cat > "$evidence_path/legal_report.json" << EOF
{
    "alert_id": "$alert_id",
    "timestamp": "$timestamp",
    "culprit": $user_identity,
    "violations": {
        "privilege_escalation": "$(echo "$priv_esc" | base64)",
        "jurisdiction_violation": "$(echo "$jur_viol" | base64)",
        "bulk_export_attempt": "$(echo "$bulk_exp" | base64)",
        "log_tampering": "$(echo "$log_tamp" | base64)"
    },
    "evidence_files": [
        "screen_evidence.png",
        "session_recording.mkv"
    ],
    "legal_basis": "Code Pénal Haïtien Art. 247-251, Loi SNISID Art. 18",
    "status": "ACTIVE_INVESTIGATION"
}
EOF
    
    # Transmission au central via satellite (bypass réseau local)
    transmit_to_central "$evidence_path" "$alert_id"
    
    # Actions automatiques immédiates
    execute_countermeasures "$culprit_user"
}

# Transmission sécurisée au central
transmit_to_central() {
    local evidence_path="$1"
    local alert_id="$2"
    
    log_info "Transmission de l'alerte $alert_id au Bureau Central..."
    
    # Chiffrement AES-512-GCM
    tar czf - "$evidence_path" | openssl enc -aes-512-gcm -K "$CENTRAL_PUBKEY" -out "/tmp/alert_$alert_id.enc"
    
    # Envoi via lien satellite direct
    if [[ -e "$SATELLITE_UPLINK" ]]; then
        cat "/tmp/alert_$alert_id.enc" > "$SATELLITE_UPLINK"
        log_info "Alerte transmise via satellite avec succès"
    else
        # Fallback: HTTPS vers central avec certificat pinning
        curl --cert /etc/snisid/certs/client.crt \
             --key /etc/snisid/certs/client.key \
             --cacert /etc/snisid/certs/ca-central.crt \
             -X POST "https://$CENTRAL_SERVER/api/v1/alerts/judas" \
             -H "Content-Type: application/octet-stream" \
             --data-binary "@/tmp/alert_$alert_id.enc"
        log_info "Alerte transmise via HTTPS sécurisé"
    fi
    
    rm -f "/tmp/alert_$alert_id.enc"
}

# Contre-mesures automatiques
execute_countermeasures() {
    local user="$1"
    
    log_alert "🔒 EXÉCUTION DES CONTRE-MESURES AUTOMATIQUES"
    
    # 1. Désactiver le compte immédiatement
    snisid-cli user disable "$user" --reason="JUDAS_PROTOCOL_TRIGGERED" --force
    
    # 2. Révoquer toutes les sessions actives
    snisid-cli session revoke-all --user="$user"
    
    # 3. Verrouiller la station de travail à distance
    snisid-cli device lock --user="$user" --message="SYSTÈME VERROUILLÉ - ENQUÊTE EN COURS"
    
    # 4. Notifier l'IGPN et la DIC
    snisid-cli notify send \
        --to="igpn@snisid.gouv.ht,dic@snisid.gouv.ht" \
        --priority="CRITICAL" \
        --subject="TENTATIVE CORRUPTION SNISID - $user" \
        --body="Un technicien a tenté de contourner le système. Dossier envoyé."
    
    # 5. Ajouter à la liste noire biométrique
    snisid-cli blacklist add --user="$user" --scope="BIOMETRIC_ALL"
    
    log_info "Contre-mesures exécutées avec succès contre $user"
}

# Vérification d'intégrité du daemon lui-même
self_check() {
    local daemon_hash=$(sha512sum "$0" | awk '{print $1}')
    local expected_hash=$(cat /etc/snisid/checksums/judas-daemon.sha512)
    
    if [[ "$daemon_hash" != "$expected_hash" ]]; then
        echo "⚠️  ALERTE: Le daemon Judas a été modifié! Auto-destruction..."
        shred -u "$0"
        exit 1
    fi
}

# Main
main() {
    echo "╔════════════════════════════════════════════════════════╗"
    echo "║     SNISID JUDAS-PROTOCOL DAEMON v1.0                  ║"
    echo "║     Système Anti-Corruption & Détection Intrusion      ║"
    echo "║     TOP SECRET - GOUVERNEMENT HAITIEN                  ║"
    echo "╚════════════════════════════════════════════════════════╝"
    
    self_check
    
    if [[ $EUID -ne 0 ]]; then
        echo "Ce script doit être exécuté en root"
        exit 1
    fi
    
    # Créer les répertoires nécessaires
    mkdir -p "$EVIDENCE_DIR" "$(dirname $LOG_FILE)"
    
    log_info "Démarrage du service de surveillance permanente..."
    detect_breach_attempt
}

main "$@"
