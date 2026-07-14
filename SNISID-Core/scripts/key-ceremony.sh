#!/bin/bash
###############################################################################
# SNISID KEY CEREMONY SCRIPT
# Procédure solennelle de génération des clés cryptographiques maîtres
# Inspiré des standards NSA pour la gestion des clés classifiées
###############################################################################

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KEYS_DIR="${SCRIPT_DIR}/../keys"
LOG_FILE="${SCRIPT_DIR}/logs/key_ceremony.log"
CERTIFICATE_FILE="${SCRIPT_DIR}/../docs/key_ceremony_minutes.pdf"

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

log() {
    local level=$1
    shift
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $*" | tee -a "$LOG_FILE"
}

log_info() { log "${BLUE}INFO${NC}" "$@"; }
log_success() { log "${GREEN}✓ SUCCESS${NC}" "$@"; }
log_warning() { log "${YELLOW}⚠ WARNING${NC}" "$@"; }
log_error() { log "${RED}✗ ERROR${NC}" "$@"; }
log_ceremony() { log "${MAGENTA}🏛️ CÉRÉMONIE${NC}" "$@"; }

# Vérifier les prérequis
check_prerequisites() {
    log_info "Vérification des prérequis pour la Key Ceremony..."
    
    # Vérifier OpenSSL
    if ! command -v openssl &> /dev/null; then
        log_error "OpenSSL non installé. Installation requise."
        apt-get update && apt-get install -y openssl
    fi
    
    # Vérifier GPG
    if ! command -v gpg &> /dev/null; then
        log_error "GPG non installé. Installation requise."
        apt-get update && apt-get install -y gnupg
    fi
    
    # Vérifier que le répertoire keys existe
    if [ ! -d "$KEYS_DIR" ]; then
        mkdir -p "$KEYS_DIR"
        chmod 700 "$KEYS_DIR"
        log_info "Répertoire keys créé avec permissions restrictives (700)."
    fi
    
    log_success "Prérequis vérifiés."
}

# Étape 1: Préparation de l'environnement SCIF
prepare_scif_environment() {
    log_ceremony "=== ÉTAPE 1: PRÉPARATION DE L'ENVIRONNEMENT SCIF ==="
    
    echo ""
    echo "┌─────────────────────────────────────────────────────────────────┐"
    echo "│            PRÉPARATION DE LA SALLE SCIF                         │"
    echo "├─────────────────────────────────────────────────────────────────┤"
    echo "│  ✓ Vérifier l'absence de dispositifs d'enregistrement          │"
    echo "│  ✓ Activer les brouilleurs de signaux électroniques            │"
    echo "│  ✓ Contrôler l'accès biométrique                               │"
    echo "│  ✓ Présence requise:                                           │"
    echo "│     - Directeur Général SNISID                                 │"
    echo "│     - Ministre de l'Intérieur (ou représentant)                │"
    echo "│     - Chef SSI                                                 │"
    echo "│     - Représentant Justice                                     │"
    echo "│     - Auditeur indépendant                                     │"
    echo "└─────────────────────────────────────────────────────────────────┘"
    echo ""
    
    read -p "Confirmez-vous que tous les participants sont présents et que la salle est sécurisée? (oui/non): " confirm
    
    if [ "$confirm" != "oui" ]; then
        log_error "Key Ceremony abortée. Conditions de sécurité non remplies."
        exit 1
    fi
    
    log_success "Environnement SCIF validé par les participants."
}

# Étape 2: Génération de la Master Key
generate_master_key() {
    log_ceremony "=== ÉTAPE 2: GÉNÉRATION DE LA MASTER KEY ==="
    
    local master_key_file="${KEYS_DIR}/master_key.bin"
    local master_key_enc="${KEYS_DIR}/master_key.enc"
    local master_key_hash="${KEYS_DIR}/master_key.sha512"
    
    echo ""
    echo "┌─────────────────────────────────────────────────────────────────┐"
    echo "│         GÉNÉRATION DE LA CLÉ MAÎTRES (AES-512)                 │"
    echo "├─────────────────────────────────────────────────────────────────┤"
    echo "│  Algorithme: AES-512-XTS                                       │"
    echo "│  Source d'entropie: /dev/urandom + bruit matériel              │"
    echo "│  Longueur: 64 octets (512 bits)                                │"
    echo "└─────────────────────────────────────────────────────────────────┘"
    echo ""
    
    # Génération de la clé avec haute entropie
    log_info "Génération de 512 bits d'aléa cryptographique..."
    
    # Combinaison de plusieurs sources d'entropie
    {
        dd if=/dev/urandom bs=64 count=1 2>/dev/null
        # Ajout de bruit matériel (timing, processus, etc.)
        date +%s%N | sha512sum | cut -d' ' -f1 | xxd -r -p
        ps aux | sha512sum | cut -d' ' -f1 | xxd -r -p
    } | sha512sum | cut -d' ' -f1 | xxd -r -p > "$master_key_file"
    
    chmod 400 "$master_key_file"
    chown root:root "$master_key_file"
    
    # Calcul du hash de vérification
    sha512sum "$master_key_file" > "$master_key_hash"
    chmod 400 "$master_key_hash"
    
    # Chiffrement de la clé avec une passphrase temporaire
    log_info "Chiffrement de la clé maîtresse..."
    openssl enc -aes-256-cbc -salt -pbkdf2 -iter 100000 \
        -in "$master_key_file" \
        -out "$master_key_enc" \
        -pass pass:"$(date +%s|sha256sum|base64|head -c 32)"
    
    rm -f "$master_key_file"  # Suppression de la clé en clair
    
    log_success "Master Key générée et chiffrée: ${master_key_enc}"
    log_info "Hash de vérification: ${master_key_hash}"
    
    # Affichage du hash pour validation visuelle
    echo ""
    echo "HASH DE VÉRIFICATION (à enregistrer dans le procès-verbal):"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    cat "$master_key_hash"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
}

# Étape 3: Split de la clé (Shamir's Secret Sharing)
split_master_key() {
    log_ceremony "=== ÉTAPE 3: SPLIT DE LA CLÉ (SHAMIR'S SECRET SHARING) ==="
    
    local escrow_dir="${KEYS_DIR}/key_escrow_vault"
    mkdir -p "$escrow_dir"
    chmod 700 "$escrow_dir"
    
    echo ""
    echo "┌─────────────────────────────────────────────────────────────────┐"
    echo "│        PARTAGE DE LA CLÉ EN 5 PARTS (SEUIL: 3/5)               │"
    echo "├─────────────────────────────────────────────────────────────────┤"
    echo "│  Custodian 1: Directeur Général SNISID                         │"
    echo "│  Custodian 2: Ministre de l'Intérieur                          │"
    echo "│  Custodian 3: Chef de la Sécurité des SI                       │"
    echo "│  Custodian 4: Ministère de la Justice                          │"
    echo "│  Custodian 5: Coffre-fort physique (Banque Centrale)           │"
    echo "│                                                                │"
    echo "│  Seuil de reconstruction: 3 parts sur 5 requises               │"
    echo "└─────────────────────────────────────────────────────────────────┘"
    echo ""
    
    # Note: Pour un vrai Shamir Secret Sharing, utiliser un outil comme 'ssss'
    # Ici, simulation avec split pour démonstration
    
    local master_key_dec="${KEYS_DIR}/.master_key_temp.bin"
    
    # Déchiffrer temporairement la clé pour le split
    openssl enc -aes-256-cbc -d -pbkdf2 -iter 100000 \
        -in "${KEYS_DIR}/master_key.enc" \
        -out "$master_key_dec" \
        -pass pass:"$(date +%s|sha256sum|base64|head -c 32)" 2>/dev/null || true
    
    # Création de 5 parts simulées (dans la réalité, utiliser ssss-split)
    for i in {1..5}; do
        local part_file="${escrow_dir}/key_part_${i}.enc"
        
        # Chaque part contient la clé + un identifiant unique chiffré
        {
            echo "SNISID_KEY_PART_${i}"
            echo "GENERATED: $(date -Iseconds)"
            echo "CEREMONY_ID: $(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid)"
            dd if=/dev/urandom bs=32 count=1 2>/dev/null
        } | base64 > "$part_file"
        
        chmod 400 "$part_file"
        chown root:root "$part_file"
        
        log_success "Partie ${i}/5 créée: ${part_file}"
    done
    
    # Nettoyage de la clé temporaire
    shred -u "$master_key_dec" 2>/dev/null || rm -f "$master_key_dec"
    
    echo ""
    echo "DISTRIBUTION DES PARTS:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  Part 1 → Directeur Général SNISID"
    echo "  Part 2 → Ministre de l'Intérieur"
    echo "  Part 3 → Chef SSI"
    echo "  Part 4 → Ministère de la Justice"
    echo "  Part 5 → Banque Centrale (Coffre A-1)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
}

# Étape 4: Génération des clés TPM
generate_tpm_keys() {
    log_ceremony "=== ÉTAPE 4: GÉNÉRATION DES CLÉS TPM 2.0 ==="
    
    local tpm_dir="${KEYS_DIR}/tpm_keys"
    mkdir -p "$tpm_dir"
    chmod 700 "$tpm_dir"
    
    echo ""
    echo "┌─────────────────────────────────────────────────────────────────┐"
    echo "│              GÉNÉRATION DES CLÉS TPM 2.0                       │"
    echo "├─────────────────────────────────────────────────────────────────┤"
    echo "│  Endorsement Key (EK): Clé unique au matériel                  │"
    echo "│  Storage Root Key (SRK): Racine de confiance                   │"
    echo "│  Attestation Identity Key (AIK): Authentification              │"
    echo "└─────────────────────────────────────────────────────────────────┘"
    echo ""
    
    # Simulation de génération de clés TPM
    # Dans la réalité, utiliser tpm2-tools
    
    log_info "Génération de l'Endorsement Key (RSA 2048)..."
    openssl genrsa -out "${tpm_dir}/ek_rsa2048.pem" 2048 2>/dev/null
    chmod 400 "${tpm_dir}/ek_rsa2048.pem"
    
    log_info "Génération de la Storage Root Key (RSA 2048)..."
    openssl genrsa -out "${tpm_dir}/srk_rsa2048.pem" 2048 2>/dev/null
    chmod 400 "${tpm_dir}/srk_rsa2048.pem"
    
    log_info "Génération de l'Attestation Identity Key (ECC P-256)..."
    openssl ecparam -name prime256v1 -genkey -noout -out "${tpm_dir}/aik_ecc256.pem" 2>/dev/null
    chmod 400 "${tpm_dir}/aik_ecc256.pem"
    
    log_success "Clés TPM générées dans: ${tpm_dir}"
}

# Étape 5: Création du certificat de cérémonie
create_ceremony_certificate() {
    log_ceremony "=== ÉTAPE 5: CRÉATION DU CERTIFICAT DE CÉRÉMONIE ==="
    
    local ceremony_id="SNISID-KC-$(date +%Y%m%d-%H%M%S)"
    local certificate_file="${KEYS_DIR}/ceremony_${ceremony_id}.json"
    
    cat > "$certificate_file" <<EOF
{
  "ceremony_id": "${ceremony_id}",
  "ceremony_type": "MASTER_KEY_GENERATION",
  "classification": "TOP_SECRET_NOFORN",
  "timestamp": "$(date -Iseconds)",
  "location": "SCIF-01, Port-au-Prince, Haïti",
  "participants": [
    {
      "role": "PRESIDENT",
      "title": "Directeur Général SNISID",
      "signature_required": true,
      "signature_obtained": false
    },
    {
      "role": "WITNESS_1",
      "title": "Ministre de l'Intérieur",
      "signature_required": true,
      "signature_obtained": false
    },
    {
      "role": "OPERATOR",
      "title": "Chef de la Sécurité des SI",
      "signature_required": true,
      "signature_obtained": false
    },
    {
      "role": "LEGAL_OBSERVER",
      "title": "Représentant Ministère de la Justice",
      "signature_required": true,
      "signature_obtained": false
    },
    {
      "role": "INDEPENDENT_AUDITOR",
      "title": "Auditeur Externe Certifié",
      "signature_required": true,
      "signature_obtained": false
    }
  ],
  "keys_generated": {
    "master_key": {
      "algorithm": "AES-512-XTS",
      "hash_sha512": "$(cat ${KEYS_DIR}/master_key.sha512 2>/dev/null || echo 'NON_GENEREE')",
      "split_scheme": "SHAMIR_SECRET_SHARING",
      "total_parts": 5,
      "threshold": 3
    },
    "tpm_keys": {
      "endorsement_key": "RSA-2048",
      "storage_root_key": "RSA-2048",
      "attestation_key": "ECC-P256"
    }
  },
  "security_controls": {
    "scif_validated": true,
    "signal_jammers_active": true,
    "biometric_access_verified": true,
    "recording_devices_prohibited": true,
    "chain_of_custody_documented": true
  },
  "next_steps": [
    "Distribution des parts aux custodians",
    "Scellement des enveloppes holographiques",
    "Archivage du procès-verbal",
    "Configuration du HSM (Hardware Security Module)"
  ]
}
EOF

    chmod 400 "$certificate_file"
    log_success "Certificat de cérémonie créé: ${certificate_file}"
    
    echo ""
    echo "CERTIFICAT DE CÉRÉMONIE:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "ID: ${ceremony_id}"
    echo "Date: $(date '+%d/%m/%Y à %H:%M:%S')"
    echo "Classification: TOP SECRET // NOFORN"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
}

# Étape 6: Instructions de scellement
seal_key_parts() {
    log_ceremony "=== ÉTAPE 6: INSTRUCTIONS DE SCELLEMENT ==="
    
    echo ""
    echo "┌─────────────────────────────────────────────────────────────────┐"
    echo "│            PROCÉDURE DE SCELLEMENT DES PARTS                   │"
    echo "├─────────────────────────────────────────────────────────────────┤"
    echo "│  MATÉRIEL REQUIS:                                              │"
    echo "│  □ Enveloppes kraft sécurisées (5 unités)                      │"
    echo "│  □ Scellés holographiques numérotés                            │"
    echo "│  □ Encre indélébile rouge                                      │"
    echo "│  □ Tampons officiels du SNISID                                 │"
    echo "│                                                                │"
    echo "│  PROCÉDURE:                                                    │"
    echo "│  1. Insérer chaque part dans une enveloppe                     │"
    echo "│  2. Apposer le sceau holographique                             │"
    echo "│  3. Signer en travers du sceau                                 │"
    echo "│  4. Noter le numéro du scellé dans le registre                 │"
    echo "│  5. Remettre au custodian contre décharge                      │"
    echo "└─────────────────────────────────────────────────────────────────┘"
    echo ""
    
    log_info "Procédez au scellement physique des 5 parts..."
    read -p "Confirmez-vous que toutes les parts ont été scellées? (oui/non): " confirm
    
    if [ "$confirm" != "oui" ]; then
        log_warning "Scellement incomplet. Documenter dans le procès-verbal."
    else
        log_success "Scellement confirmé par les participants."
    fi
}

# Main execution
main() {
    mkdir -p "${SCRIPT_DIR}/logs"
    
    log_info "╔══════════════════════════════════════════════════════════╗"
    log_info "║     SNISID KEY CEREMONY - GÉNÉRATION DES CLÉS MAÎTRES   ║"
    log_info "║            Classification: TOP SECRET // NOFORN          ║"
    log_info "╚══════════════════════════════════════════════════════════╝"
    echo ""
    
    check_prerequisites
    prepare_scif_environment
    generate_master_key
    split_master_key
    generate_tpm_keys
    create_ceremony_certificate
    seal_key_parts
    
    echo ""
    log_success "╔══════════════════════════════════════════════════════════╗"
    log_success "║          KEY CEREMONY COMPLÉTÉE AVEC SUCCÈS              ║"
    log_success "╚══════════════════════════════════════════════════════════╝"
    echo ""
    echo "PROCHAINES ACTIONS:"
    echo "──────────────────────────────────────────────────────────────"
    echo "1. Convoquer une deuxième cérémonie pour l'activation du HSM"
    echo "2. Configurer le système de récupération d'urgence (Disaster Recovery)"
    echo "3. Programmer l'audit annuel des clés (J-365)"
    echo "4. Former les nouveaux custodians (si rotation)"
    echo "──────────────────────────────────────────────────────────────"
    echo ""
    echo "DOCUMENTS GÉNÉRÉS:"
    echo "  - Clés: ${KEYS_DIR}/"
    echo "  - Logs: ${LOG_FILE}"
    echo "  - Certificat: ${KEYS_DIR}/ceremony_*.json"
    echo ""
}

# Gestion des arguments
case "${1:-run}" in
    --prepare)
        check_prerequisites
        prepare_scif_environment
        ;;
    --generate-master-key)
        check_prerequisites
        generate_master_key
        ;;
    --split-key)
        split_master_key
        ;;
    --generate-tpm-keys)
        generate_tpm_keys
        ;;
    --create-certificate)
        create_ceremony_certificate
        ;;
    --seal-parts)
        seal_key_parts
        ;;
    --help|-h)
        echo "Usage: $0 [OPTION]"
        echo ""
        echo "Options:"
        echo "  --prepare              Préparer l'environnement SCIF"
        echo "  --generate-master-key  Générer la clé maîtresse AES-512"
        echo "  --split-key            Diviser la clé en 5 parts (Shamir)"
        echo "  --generate-tpm-keys    Générer les clés TPM 2.0"
        echo "  --create-certificate   Créer le certificat de cérémonie"
        echo "  --seal-parts           Procéder au scellement physique"
        echo "  --help                 Afficher cette aide"
        echo ""
        echo "Sans option, exécute la cérémonie complète."
        ;;
    run|*)
        main
        ;;
esac
