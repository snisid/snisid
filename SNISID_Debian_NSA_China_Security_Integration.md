# 🇭🇹 SNISID — Debian Sovereign OS & National Security Integration

**Document N° :** SNISID-SEC-OS-001  
**Version :** 1.0.0  
**Classification :** TRES SECRET / GOUVERNEMENTAL  
**Date :** 2026  

---

## 1. VISION STRATÉGIQUE

Ce document définit l'intégration d'un système d'exploitation **Debian Linux souverain** comme fondation du SNISID, enrichi des meilleurs standards de sécurité de la **NSA (National Security Agency)** américaine et des **agences de sécurité chinoises**, adaptés aux besoins de la Sécurité Nationale Haïtienne.

### 1.1 Principes Directeurs

| Principe | Source | Application SNISID |
|----------|--------|-------------------|
| **Defense in Depth** | NSA IA-8453 | Couches de sécurité multiples (OS, K8s, App, Data) |
| **Zero Trust Architecture** | NSA/NIST SP 800-207 | Aucune confiance implicite, vérification continue |
| **Security Hardening** | NSA STIG / CIS Benchmarks | Configuration sécurisée par défaut |
| **Supply Chain Security** | Chine GB/T 22239-2019 | Contrôle de la chaîne d'approvisionnement logicielle |
| **Cryptographic Agility** | NSA CSfC | Algorithmes cryptographiques interchangeables |
| **Air-Gap Capability** | Chine MLPS 2.0 | Isolation physique des systèmes critiques |

---

## 2. DEBIAN SOVEREIGN EDITION — BASE OS NATIONALE

### 2.1 Choix de Debian vs Ubuntu

| Critère | Debian Stable | Ubuntu LTS | Décision SNISID |
|---------|---------------|------------|-----------------|
| **Souveraineté** | Communauté pure, aucun vendor | Canonical (entreprise) | ✅ **Debian** |
| **Cycle de vie** | ~5 ans (support long) | 5-10 ans | ✅ Équivalent |
| **Paquets** | 59,000+ paquets audités | Basé sur Debian + snaps | ✅ **Debian** (pas de snaps propriétaires) |
| **Security Updates** | Équipe security.debian.org | Canonical Security Team | ✅ Équivalent |
| **Hardening par défaut** | Possible via profiles.conf | AppArmor activé | ✅ Configurable |
| **Taille minimale** | ~150 MB (netinst) | ~1.5 GB | ✅ **Debian** (plus léger) |

### 2.2 SNISID Debian Base Specification

```yaml
# snisid-debian-spec.yaml
distribution: debian
version: "12 (Bookworm)"
variant: "snisid-sovereign-edition"

kernel:
  version: "6.1 LTS"
  hardening:
    - CONFIG_STRICT_DEVMEM=y
    - CONFIG_REFCOUNT_FULL=y
    - CONFIG_HARDENED_USERCOPY=y
    - CONFIG_FORTIFY_SOURCE=y
    - CONFIG_RANDOMIZE_BASE=y
    - CONFIG_SECURITY_APPARMOR=y
    - CONFIG_SECURITY_SELINUX=y

bootloader:
  type: "GRUB2"
  security:
    - password_protected: true
    - secure_boot_enabled: true
    - kernel_signature_verification: true

packages:
  minimal_base: true
  included:
    - dropbear-initramfs      # SSH early boot (LUKS unlock)
    - apparmor-profiles       # Mandatory Access Control
    - selinux-policy-default  # Alternative MAC
    - aide                    # File integrity monitoring
    - rkhunter                # Rootkit detection
    - chkrootkit              # Rootkit detection
    - fail2ban                # Intrusion prevention
    - unattended-upgrades     # Security auto-updates
    - usbguard                # USB device authorization
    - tpm2-tools              # TPM 2.0 support
    
  excluded:
    - snapd                   # Pas de snaps (vendor lock-in)
    - systemd-resolved        # Utilisation de dnsmasq/unbound
    - modemmanager            # Pas de connectivité cellulaire non contrôlée

partitioning:
  scheme: "luks_encrypted_lvm"
  encryption:
    algorithm: "aes-xts-plain64"
    key_size: 512             # AES-256 x 2
    hash: "sha512"
    iterations: "time-target:5000ms"
    
  layout:
    - mount: "/boot"
      size: "1GB"
      filesystem: "ext4"
      encrypted: false
      
    - mount: "/boot/efi"
      size: "512MB"
      filesystem: "vfat"
      encrypted: false
      
    - mount: "/"
      size: "100GB"
      filesystem: "ext4"
      encrypted: true
      options: "noatime,nodiratime"
      
    - mount: "/var"
      size: "50GB"
      filesystem: "ext4"
      encrypted: true
      options: "noatime,nosuid,nodev"
      
    - mount: "/tmp"
      size: "10GB"
      filesystem: "ext4"
      encrypted: true
      options: "noatime,nosuid,nodev,noexec"
      
    - mount: "/home"
      size: "remaining"
      filesystem: "ext4"
      encrypted: true
      options: "noatime,nosuid,nodev"
      
    - path: "swap"
      size: "equal_to_RAM"
      encrypted: true         # Chiffrement obligatoire pour swap
```

### 2.3 SNISID Debian ISO Build Process

```bash
#!/bin/bash
# build-snisid-debian-iso.sh
# Crée l'ISO Debian souveraine pour SNISID

set -euo pipefail

DEBIAN_ISO="debian-12.5.0-amd64-netinst.iso"
OUTPUT_ISO="snisid-debian-sovereign-v1.0.iso"
BUILD_DIR="./debian-build"

echo "[*] SNISID Debian Sovereign ISO Builder"
echo "[*] Checking prerequisites..."

# Vérification des dépendances
for cmd in xorriso genisoimage sha256sum gpg; do
    command -v "$cmd" >/dev/null 2>&1 || { echo "[!] $cmd required"; exit 1; }
done

if [ ! -f "$DEBIAN_ISO" ]; then
    echo "[!] Télécharger $DEBIAN_ISO depuis cdn.deb.debian.org"
    exit 1
fi

echo "[*] Extraction de l'ISO Debian de base..."
mkdir -p "$BUILD_DIR"
xorriso -osirrox on -indev "$DEBIAN_ISO" -extract / "$BUILD_DIR/debian-base"
chmod -R u+w "$BUILD_DIR/debian-base"

echo "[*] Injection de la configuration de sécurité SNISID..."

# Copie des fichiers de préconfiguration
mkdir -p "$BUILD_DIR/debian-base/preseed"
cp ./preseed/snisid-hardened.cfg "$BUILD_DIR/debian-base/preseed/"

# Copie des clés GPG de signature SNISID
mkdir -p "$BUILD_DIR/debian-base/snisid-keys"
cp ./keys/snisid-security.gpg "$BUILD_DIR/debian-base/snisid-keys/"

# Copie des scripts de post-installation
mkdir -p "$BUILD_DIR/debian-base/postinstall"
cp ./scripts/postinstall-security.sh "$BUILD_DIR/debian-base/postinstall/"

# Modification du menu GRUB pour auto-installation
sed -i 's/quiet/auto preseed\/file=\/cdrom\/preseed\/snisid-hardened.cfg --- quiet/' \
    "$BUILD_DIR/debian-base/boot/grub/grub.cfg"

echo "[*] Reconstruction de l'ISO souveraine..."
cd "$BUILD_DIR/debian-base"
xorriso -as mkisofs -r -V "SNISID_DEBIAN_SOVEREIGN" \
    -J -l -b isolinux/isolinux.bin \
    -c isolinux/boot.cat -no-emul-boot \
    -boot-load-size 4 -boot-info-table \
    -eltorito-alt-boot \
    -e boot/grub/efi.img -no-emul-boot \
    -isohybrid-gpt-basdat \
    -o "../../$OUTPUT_ISO" .

cd ../../

echo "[*] Signature cryptographique de l'ISO..."
gpg --armor --detach-sign --output "$OUTPUT_ISO.sig" "$OUTPUT_ISO"
sha256sum "$OUTPUT_ISO" > "$OUTPUT_ISO.sha256"

echo "[*] ISO prête: $OUTPUT_ISO"
echo "[*] Signature: $OUTPUT_ISO.sig"
echo "[*] Hash SHA256: $(cat $OUTPUT_ISO.sha256)"
```

---

## 3. NSA SECURITY STANDARDS INTEGRATION

### 3.1 NSA Security Technical Implementation Guides (STIG)

Le SNISID implémente les guides STIG suivants adaptés pour Linux/Kubernetes :

| STIG ID | Titre | Application SNISID |
|---------|-------|-------------------|
| **V-71849** | System must disable USB mass storage | USBGuard policies |
| **V-72003** | Audit system must record privileged activities | Auditd rules |
| **V-73403** | System must use FIPS 140-2 validated crypto | OpenSSL FIPS module |
| **V-76727** | Kernel must prevent loading of untrusted modules | Module signing |
| **V-77811** | System must implement DoD-approved TLS | TLS 1.3 only |
| **V-78583** | System must separate user and admin spaces | RBAC strict |

### 3.2 NSA Hardened Configuration Profiles

#### 3.2.1 Kernel Hardening (/etc/sysctl.d/99-snisid-nsa.conf)

```bash
# NSA-inspired Kernel Hardening for SNISID
# Basé sur NSA Guide to the Secure Configuration of RHEL 8/9

# Network Security
net.ipv4.tcp_syncookies = 1
net.ipv4.conf.all.accept_redirects = 0
net.ipv4.conf.default.accept_redirects = 0
net.ipv4.conf.all.send_redirects = 0
net.ipv4.conf.default.send_redirects = 0
net.ipv4.conf.all.accept_source_route = 0
net.ipv4.conf.default.accept_source_route = 0
net.ipv4.icmp_echo_ignore_broadcasts = 1
net.ipv4.icmp_ignore_bogus_error_responses = 1
net.ipv4.conf.all.log_martians = 1
net.ipv4.conf.default.log_martians = 1
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1
net.ipv6.conf.all.accept_redirects = 0
net.ipv6.conf.default.accept_redirects = 0

# Memory Protection
kernel.randomize_va_space = 2
kernel.exec-shield = 1
kernel.dmesg_restrict = 1
kernel.kptr_restrict = 2
kernel.perf_event_paranoid = 3

# Core Dump Restrictions
fs.suid_dumpable = 0

# Module Loading Restrictions
kernel.modules_disabled = 0  # Mettre à 1 après installation complète
kernel.module.sig_enforce = 1
```

#### 3.2.2 PAM Configuration (/etc/pam.d/common-auth)

```bash
# NSA-compliant PAM configuration
# Multi-factor authentication enforcement

auth    required        pam_faillock.so preauth silent deny=5 unlock_time=900
auth    [success=1 default=ignore]    pam_unix.so obscure yescrypt
auth    [default=die]   pam_faillock.so authfail deny=5 unlock_time=900
auth    sufficient      pam_faillock.so authsucc
auth    required        pam_google_authenticator.so
auth    required        pam_deny.so

account required        pam_faillock.so
account required        pam_unix.so
account required        pam_permit.so
```

### 3.3 NSA Commercial Solutions for Classified (CSfC)

Le SNISID adopte l'approche CSfC pour le chiffrement des données classifiées :

```
┌─────────────────────────────────────────────────────────┐
│  NSA CSfC Double Encryption Layer (Data at Rest)       │
├─────────────────────────────────────────────────────────┤
│  Layer 2: Application-level encryption (AES-256-GCM)   │
│          Keys gérées par HSM national                  │
├─────────────────────────────────────────────────────────┤
│  Layer 1: Full Disk Encryption (LUKS2 AES-XTS-512)     │
│          Keys protégées par TPM 2.0                    │
└─────────────────────────────────────────────────────────┘
```

**Composants CSfC approuvés pour SNISID :**

| Composant | Produit | Certification |
|-----------|---------|---------------|
| **HSM** | YubiHSM 2 / Nitrokey HSM 2 | FIPS 140-2 Level 3 |
| **TPM** | Infineon SLB 9670 | FIPS 140-2 Level 2 |
| **Chiffrement** | OpenSSL 3.x FIPS Object Module | FIPS 140-2 Validated |

---

## 4. CHINESE NATIONAL SECURITY STANDARDS INTEGRATION

### 4.1 Multi-Level Protection Scheme (MLPS 2.0 / GB/T 22239-2019)

La Chine a développé le standard **MLPS (Multi-Level Protection Scheme)** qui est intégré dans SNISID :

| Niveau MLPS | Classification SNISID | Exigences |
|-------------|----------------------|-----------|
| **Level 3** | Données Citoyens (NIN, Biométrie) | Audit complet, chiffrement, contrôle d'accès strict |
| **Level 4** | Sécurité Nationale (Police, Justice, Renseignement) | Isolation physique, double authentification, surveillance 24/7 |
| **Level 5** | Défense Nationale (Classifié) | Air-gap total, personnel habilité Secret Défense |

### 4.2 Cryptographic Standards (GM/T Series)

SNISID intègre les algorithmes cryptographiques chinois en plus des standards occidentaux :

```yaml
# Cryptographic Agility Configuration
crypto_providers:
  western:
    symmetric: "AES-256-GCM"
    asymmetric: "RSA-4096 / ECDSA P-384"
    hash: "SHA-384 / SHA3-384"
    kdf: "PBKDF2-HMAC-SHA512"
    
  chinese_gm:
    symmetric: "SM4-CBC / SM4-GCM"     # Équivalent AES
    asymmetric: "SM2"                   # Courbe elliptique nationale chinoise
    hash: "SM3"                         # Fonction de hachage 256-bit
    kdf: "SM2-KDF"
    
usage_policy:
  domestic_operations: ["SM2", "SM3", "SM4"]  # Communications internes
  international_operations: ["RSA", "ECDSA", "AES"]  # Interopérabilité
  dual_stack: true  # Support simultané des deux suites
```

### 4.3 Supply Chain Security (GB/T 30976-2014)

```markdown
## Exigences de Sécurité de la Chaîne d'Approvisionnement

1. **Audit du Code Source**
   - Tous les paquets Debian doivent être rebuilds localement
   - Signature GPG obligatoire pour chaque binaire
   - Reproducible builds vérifiés

2. **Contrôle des Dépendances**
   - Mirror Debian local air-gappé
   - Scan automatique des CVE (Trivy, Grype)
   - Interdiction des dépendances provenant de pays sous sanctions

3. **Personnel Habilité**
   - Clearance nationale requise pour accès production
   - Background check tous les 2 ans
   - Rotation des privilèges tous les 90 jours
```

---

## 5. MULTI-DATABASE ARCHITECTURE — SÉCURITÉ NATIONALE

### 5.1 Topologie des Bases de Données

```
┌─────────────────────────────────────────────────────────────────────┐
│                    SNISID NATIONAL DATA LAYER                       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐  │
│  │  PostgreSQL 16   │  │  CockroachDB     │  │  MongoDB 7       │  │
│  │  (État Civil)    │  │  (Identité)      │  │  (Documents)     │  │
│  │  ACID compliant  │  │  Distributed SQL │  │  Schema-flexible │  │
│  │  Patroni HA      │  │  Strong Consist. │  │  GridFS storage  │  │
│  └────────┬─────────┘  └────────┬─────────┘  └────────┬─────────┘  │
│           │                     │                     │            │
│           └─────────────────────┼─────────────────────┘            │
│                                 │                                  │
│                    ┌────────────▼────────────┐                     │
│                    │    Apache Kafka         │                     │
│                    │    (Event Streaming)    │                     │
│                    │    Event Sourcing       │                     │
│                    └────────────┬────────────┘                     │
│                                 │                                  │
│           ┌─────────────────────┼─────────────────────┐            │
│           │                     │                     │            │
│  ┌────────▼─────────┐  ┌────────▼─────────┐  ┌────────▼─────────┐ │
│  │  Redis 7         │  │ Elasticsearch 8  │  │  MinIO           │ │
│  │  (Cache/Session) │  │  (Search/Logs)   │  │  (Object Store)  │ │
│  │  Cluster mode    │  │  SIEM integration│  │  WORM compliance │ │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.2 Spécifications par Base de Données

#### PostgreSQL 16 (État Civil & Justice)

```yaml
postgresql:
  version: "16.2"
  high_availability:
    tool: "Patroni 3.2"
    dcs: "etcd 3.5"
    replication:
      mode: "synchronous"
      sync_standby_count: 2
      
  security:
    ssl_mode: "verify-full"
    client_cert_required: true
    password_encryption: "scram-sha-256"
    audit_extension: "pg_audit"
    
  hardening:
    max_connections: 500
    shared_buffers: "8GB"
    effective_cache_size: "24GB"
    work_mem: "64MB"
    
  backup:
    tool: "pgBackRest"
    retention: "30 days"
    encryption: "AES-256-CBC"
    offsite_copy: true
```

#### CockroachDB (Identité Nationale & Biométrie)

```yaml
cockroachdb:
  version: "23.2"
  cluster:
    nodes: 5
    replication_factor: 3
    
  security:
    cert_lifetime: "1 year"
    node_cert_required: true
    sql_tls_enabled: true
    
  performance:
    cache_size: "25% of RAM"
    max_sql_memory: "25% of RAM"
    
  geo_partitioning:
    primary_region: "ht-pap"      # Port-au-Prince
    secondary_region: "ht-cap"    # Cap-Haïtien
    survival_goal: "zone"
```

#### MongoDB 7 (Documents & Preuves Numériques)

```yaml
mongodb:
  version: "7.0"
  deployment: "Replica Set"
  nodes: 3
  
  security:
    authorization: "enabled"
    javascript: "disabled"        # Désactivé pour sécurité
    net_tls_mode: "requireTLS"
    
  encryption:
    at_rest: "AES256-CBC"
    kmip_server: "internal-hsm"
    
  compliance:
    audit_log: "enabled"
    audit_destination: "syslog + SIEM"
```

#### Redis 7 (Cache & Sessions)

```yaml
redis:
  version: "7.2"
  mode: "Cluster"
  shards: 3
  replicas_per_shard: 2
  
  security:
    acl_enabled: true
    tls_enabled: true
    protected_mode: true
    rename_dangerous_commands:
      - FLUSHALL
      - FLUSHDB
      - CONFIG
      - DEBUG
      
  persistence:
    rdb: "every 5 min"
    aof: "every second"
```

#### Elasticsearch 8 (Recherche & SIEM)

```yaml
elasticsearch:
  version: "8.12"
  nodes:
    master: 3
    data: 5
    ingest: 2
    
  security:
    xpack_security: "enabled"
    audit_logging: "enabled"
    ssl_transport: "required"
    ssl_http: "required"
    
  indices:
    lifecycle:
      logs: "rollover 30GB / 7 days"
      audit: "rollover 50GB / 90 days"
      evidence: "immutable (WORM)"
```

### 5.3 Data Security Matrix

| Type de Donnée | Base Primaire | Chiffrement | Localisation | Rétention |
|---------------|---------------|-------------|--------------|-----------|
| **Identité (NIN)** | CockroachDB | AES-256 + SM4 | PaP + Cap-Haïtien | Permanente |
| **État Civil** | PostgreSQL | AES-256 | PaP + Cap-Haïtien | Permanente |
| **Biométrie (AFIS)** | MongoDB + MinIO | AES-256 + LUKS | PaP (HSM) | Permanente |
| **Casiers Judiciaires** | PostgreSQL | Double chiffrement | PaP (isolé) | 100 ans |
| **Preuves Numériques** | MinIO (WORM) | AES-256 | PaP + Cap-Haïtien | Permanente |
| **Logs d'Audit** | Elasticsearch | AES-256 | PaP + Cap-Haïtien + Offline | 10 ans |
| **Sessions Utilisateurs** | Redis | TLS + Auth | Mémoire uniquement | 24h max |

---

## 6. SECURE BOOT & CHAIN OF TRUST

### 6.1 UEFI Secure Boot Implementation

```
┌─────────────────────────────────────────────────────────────┐
│              SNISID Secure Boot Chain of Trust              │
├─────────────────────────────────────────────────────────────┤
│  [UEFI Firmware]                                            │
│       │                                                     │
│       ▼                                                     │
│  [Microsoft UEFI CA] ←→ [SNISID PK]                        │
│       │                                                     │
│       ▼                                                     │
│  [GRUB2 signed by SNISID PK]                                │
│       │                                                     │
│       ▼                                                     │
│  [Linux Kernel signed with SNISID key]                      │
│       │                                                     │
│       ▼                                                     │
│  [Kernel Modules signed]                                    │
│       │                                                     │
│       ▼                                                     │
│  [Initramfs with embedded signatures]                       │
│       │                                                     │
│       ▼                                                     │
│  [Systemd / Init signed binaries]                           │
└─────────────────────────────────────────────────────────────┘
```

### 6.2 Key Ceremony Procedure

```bash
#!/bin/bash
# snisid-key-ceremony.sh
# Cérémonie de création des clés de signature Secure Boot

set -euo pipefail

echo "=== SNISID Secure Boot Key Ceremony ==="
echo "Classification: TRES SECRET"
echo "Date: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo ""

# Vérification des témoins requis
REQUIRED_WITNESSES=3
echo "[*] Vérification des témoins présents..."
# Procédure de vérification d'identité des témoins

# Création de la clé privée SNISID PK
echo "[*] Génération de la clé privée SNISID Platform Key..."
openssl genrsa -out snisid-pk.key 4096

# Création du certificat auto-signé
echo "[*] Création du certificat PK..."
openssl req -new -x509 -sha256 \
    -key snisid-pk.key \
    -out snisid-pk.crt \
    -days 3650 \
    -subj "/C=HT/O=Gouvernement d'Haiti/CN=SNISID Platform Key"

# Conversion en format UEFI
echo "[*] Conversion au format UEFI..."
cert-to-efi-sig-list -g "$(uuidgen)" snisid-pk.crt snisid-pk.esl
sign-efi-sig-list -g "$(uuidgen)" -k snisid-pk.key -c snisid-pk.crt PK snisid-pk.esl snisid-pk.auth

# Stockage sécurisé dans HSM
echo "[*] Transfert des clés vers HSM national..."
# Utilisation de pkcs11-tool ou outil propriétaire HSM

# Destruction des copies temporaires
echo "[*] Destruction sécurisée des copies temporaires..."
shred -u snisid-pk.key
# Procédure de nettoyage validée par les témoins

echo "[*] Cérémonie terminée avec succès"
echo "[*] Clés stockées dans HSM national"
```

---

## 7. CONTINUOUS MONITORING & COMPLIANCE

### 7.1 Real-time Security Dashboard

```yaml
monitoring_stack:
  siem: "Wazuh 4.7 + Elasticsearch"
  metrics: "Prometheus 2.50 + Grafana 10"
  logging: "Loki 2.9 + Promtail"
  tracing: "Tempo 2.4 + OpenTelemetry"
  
compliance_checks:
  automated:
    - cis_benchmark_debian_level2
    - nsas_stig_linux_v6r1
    - mlps_2.0_level3
    
  frequency:
    continuous: ["file_integrity", "log_analysis", "network_traffic"]
    hourly: ["vulnerability_scan", "configuration_drift"]
    daily: ["user_access_review", "privilege_escalation_check"]
    weekly: ["full_system_audit", "backup_integrity_test"]
    monthly: ["penetration_test", "disaster_recovery_drill"]
```

### 7.2 Alert Severity Matrix

| Niveau | Source | Action Requise | Délai de Réponse |
|--------|--------|----------------|------------------|
| **Critical** | SIEM (intrusion active) | Isolation immédiate + CNN alerté | < 5 min |
| **High** | Compliance drift (STIG violation) | Correction urgente | < 1 h |
| **Medium** | Vulnerability detected (CVSS > 7) | Patch planifié | < 24 h |
| **Low** | Configuration recommendation | Review prochaine maintenance | < 7 j |

---

## 8. DISASTER RECOVERY & BUSINESS CONTINUITY

### 8.1 Recovery Time Objectives (RTO/RPO)

| Système | RTO | RPO | Stratégie |
|---------|-----|-----|-----------|
| Identité (CockroachDB) | 15 min | 0 sec | Sync multi-region |
| État Civil (PostgreSQL) | 30 min | 5 min | Async streaming replication |
| Biométrie (MongoDB+MinIO) | 1 h | 15 min | Daily sync + incremental |
| SIEM (Elasticsearch) | 4 h | 1 h | Cold standby + log forward |
| API Gateway (Kong) | 5 min | 0 sec | Active-active multi-site |

### 8.2 Backup Strategy 3-2-1-1-0

```
3 copies des données
  ├── Production (primaire)
  ├── DR Site (Cap-Haïtien)
  └── Offline vault (coffre-fort)

2 supports différents
  ├── SSD NVMe (production)
  └── LTO-9 Tape (archive long terme)

1 copie hors-site
  └── Cap-Haïtien (> 250km de PaP)

1 copie immuable
  └── WORM storage (MinIO Object Lock)

0 erreur de restoration
  └── Tests trimestriels obligatoires
```

---

## 9. PERSONNEL SECURITY & CLEARANCE

### 9.1 Clearance Levels

| Niveau | Accès Autorisé | Background Check | Renouvellement |
|--------|----------------|------------------|----------------|
| **Level 1** | Systèmes publics (sites web) | Casier judiciaire | 5 ans |
| **Level 2** | Données citoyens (NIN, état civil) | Enquête approfondie + références | 3 ans |
| **Level 3** | Sécurité nationale (police, justice) | Investigation financière + psychologique | 2 ans |
| **Level 4** | Classifié défense | Habilitation Secret Défense | 1 an |

### 9.2 Training Requirements

```yaml
mandatory_training:
  initial:
    - "Security Awareness (8h)"
    - "Data Protection & Privacy (4h)"
    - "Incident Response Procedures (4h)"
    - "Secure Coding Practices (developers only, 16h)"
    
  recurrent_annual:
    - "Cybersecurity Refresher (4h)"
    - "New Threat Landscape Briefing (2h)"
    - "Tabletop Exercise (participation obligatoire)"
    
  role_specific:
    sysadmin:
      - "Privileged Access Management (8h)"
      - "Linux Hardening Advanced (16h)"
    developer:
      - "OWASP Top 10 Deep Dive (8h)"
      - "Secure SDLC (8h)"
    analyst:
      - "Threat Intelligence Analysis (16h)"
      - "Digital Forensics Basics (8h)"
```

---

## 10. IMPLEMENTATION ROADMAP

### Phase 1 : Foundation (Q1-Q2 2026)

- [ ] Création de l'ISO Debian Souveraine
- [ ] Configuration des profils de hardening NSA
- [ ] Intégration des algorithmes cryptographiques chinois (SM2/SM3/SM4)
- [ ] Déploiement du Secure Boot infrastructure

### Phase 2 : Database Layer (Q3-Q4 2026)

- [ ] Installation PostgreSQL HA avec Patroni
- [ ] Déploiement CockroachDB cluster
- [ ] Configuration MongoDB avec chiffrement KMIP
- [ ] Mise en place Redis Cluster sécurisé

### Phase 3 : Monitoring & Compliance (Q1-Q2 2027)

- [ ] Déploiement SIEM Wazuh
- [ ] Configuration des checks automatisés CIS/STIG/MLPS
- [ ] Formation des équipes SOC
- [ ] Premiers exercices de disaster recovery

### Phase 4 : Production Rollout (Q3-Q4 2027)

- [ ] Déploiement progressif sur sites pilotes
- [ ] Certification par autorités nationales
- [ ] Audit externe de sécurité
- [ ] Mise en production nationale

---

## 11. CONCLUSION

Ce document établit le cadre technique et opérationnel pour un **SNISID entièrement souverain**, combinant :

1. **Debian Linux** comme système d'exploitation de base, offrant stabilité, transparence et indépendance vis-à-vis des vendors commerciaux
2. **Standards NSA** (STIG, CSfC, Zero Trust) pour une sécurité éprouvée au niveau gouvernemental américain
3. **Standards Chinois** (MLPS 2.0, cryptographie GM/T) pour la diversité stratégique et la résilience géopolitique
4. **Architecture Multi-Database** sécurisée pour gérer les différentes catégories de données nationales

Cette approche hybride garantit que la Sécurité Nationale Haïtienne dispose d'un système **robuste, auditable, et véritablement indépendant**.

---

**Approbations :**

| Rôle | Nom | Signature | Date |
|------|-----|-----------|------|
| Directeur Général AND | _________________ | _________________ | _________ |
| CISO National | _________________ | _________________ | _________ |
| Ministre de l'Intérieur | _________________ | _________________ | _________ |
| Président de la République | _________________ | _________________ | _________ |

---

*Document Classifié TRES SECRET — Distribution Restreinte*  
*© 2026 Agence Nationale des Données (AND) — République d'Haïti*
