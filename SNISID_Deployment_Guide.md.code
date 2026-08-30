# SNISID SOVEREIGN OS - GUIDE DE DÉPLOIEMENT COMPLET
## Système d'Exploitation National Indépendant et Bootable
### Version: 1.0 | Classification: TOP SECRET / Gouvernement Haïtien

---

## 📋 TABLE DES MATIÈRES

1. [Vue d'ensemble](#1-vue-densemble)
2. [Architecture technique](#2-architecture-technique)
3. [Prérequis matériels](#3-prérequis-matériels)
4. [Construction de l'ISO](#4-construction-de-liso)
5. [Installation sur matériel](#5-installation-sur-matériel)
6. [Configuration post-installation](#6-configuration-post-installation)
7. [Intégration des modules NSA/Chine](#7-intégration-des-modules-nsachine)
8. [Procédures de sécurité](#8-procédures-de-sécurité)
9. [Maintenance et mises à jour](#9-maintenance-et-mises-à-jour)
10. [Dépannage](#10-dépannage)

---

## 1. VUE D'ENSEMBLE

### 1.1 Qu'est-ce que SNISID OS?

**SNISID OS** est un système d'exploitation souverain basé sur **Debian 12 "Bookworm"**, spécialement conçu pour les besoins de la Sécurité Nationale Haïtienne. Il intègre nativement:

- ✅ Les outils de renseignement de la **NSA** (DataWave, Emissary, MAAT)
- ✅ Les standards de sécurité **chinois** (MLPS 2.0, cryptographie SM2/SM3/SM4)
- ✅ Une architecture **Zero Trust** complète
- ✅ Des modules **OSINT** avancés (PhoneInfoga, GeoSpy, Storm-Breaker, etc.)
- ✅ Un chiffrement de bout en bout (LUKS2 AES-512)

### 1.2 Caractéristiques principales

| Caractéristique | Description |
|-----------------|-------------|
| **Bootable** | ISO autonome, installation bare-metal ou VM |
| **Indépendant** | Aucune télémétrie vers l'étranger |
| **Immutable** | Rootfs en lecture seule avec vérification d'intégrité |
| **Multi-BDD** | Support natif PostgreSQL, CockroachDB, MongoDB, Redis |
| **Air-Gap Ready** | Fonctionne complètement hors ligne |
| **TPM 2.0** | Support natif du Trusted Platform Module |

---

## 2. ARCHITECTURE TECHNIQUE

### 2.1 Stack technologique

```
┌─────────────────────────────────────────────────────────────┐
│                    COUCHE UTILISATEUR                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐   │
│  │ Browser  │ │ Terminal │ │ Dashboard│ │ Apps Métier  │   │
│  │ Sécurisé │ │ Tactique │ │ SIEM     │ │ (Police, etc)│   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────┘   │
├─────────────────────────────────────────────────────────────┤
│              MODULES DE RENSEIGNEMENT (OPT/SNISID)           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐   │
│  │ DataWave │ │ Emissary │ │  MAAT    │ │ OSINT Tools  │   │
│  │  (NSA)   │ │  (NSA)   │ │  (NSA)   │ │ (PhoneInfoga)│   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                 SYSTÈME DURCI (DEBIAN 12)                    │
│  • Kernel 6.1 LTS patché KSPP                              │
│  • AppArmor + SELinux en mode strict                       │
│  • nftables Zero Trust                                     │
│  • Auditd + AIDE                                           │
├─────────────────────────────────────────────────────────────┤
│              CHIFFREMENT ET SÉCURITÉ MATÉRIELLE             │
│  • LUKS2 AES-512-XTS                                       │
│  • TPM 2.0 + Secure Boot                                   │
│  • HSM pour clés maîtresses                                │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Flux de données sécurisé

```
[Capteurs Terrain] → [WireGuard VPN] → [Proxy Air-Gap] → [DataWave]
                                                        ↓
[Base Biométrique] ← [MAAT Anonymization] ← [Emissary Stream]
```

---

## 3. PRÉREQUIS MATÉRIELS

### 3.1 Configuration minimale (Poste Agent)

| Composant | Spécification | Notes |
|-----------|---------------|-------|
| CPU | Intel i7 / AMD Ryzen 7 (8 cœurs) | AES-NI requis |
| RAM | 32 Go DDR4 ECC | Minimum vital |
| Stockage | 1 To NVMe SSD | Chiffrement matériel recommandé |
| TPM | TPM 2.0 | Obligatoire pour Secure Boot |
| Réseau | 1 GbE Intel | Cartes approuvées uniquement |
| SmartCard | Lecteur CCID | Pour authentification forte |

### 3.2 Configuration recommandée (Serveur Datacenter)

| Composant | Spécification | Notes |
|-----------|---------------|-------|
| CPU | Dual Xeon Gold / EPYC Rome (64+ cœurs) | Pour IA/DataWave |
| RAM | 512 Go - 1 To DDR4/5 ECC | Cluster HA |
| Stockage | 50 To NVMe RAID 10 + 100 To HDD | Archive WORM |
| GPU | NVIDIA A100 / H100 | Pour ACE-T (IA prédictive) |
| HSM | Thales Luna / YubiHSM2 | Gestion des clés |
| Réseau | 10/25 GbE + Fibre Optique | Redondance requise |

### 3.3 Matériel certifié (Liste restreinte)

**✅ Approuvé:**
- Serveurs: Dell PowerEdge R750, HPE ProLiant DL380
- Workstations: Lenovo ThinkStation P620, Dell Precision 7865
- Réseau: Cisco ISR 4000, Fortinet FortiGate 100F
- HSM: Thales Luna Network HSM A750

**❌ Interdit:**
- Tout matériel Huawei, ZTE, Hikvision
- Composants non-certifiés TPM
- Disques SSD sans chiffrement matériel

---

## 4. CONSTRUCTION DE L'ISO

### 4.1 Préparation de l'environnement de build

```bash
# Sur une machine Debian 12 sécurisée
sudo apt-get update
sudo apt-get install -y \
    genisoimage xorriso debootstrap squashfs-tools \
    grub-pc-bin grub-efi-amd64-bin mtools isolinux

# Vérification de l'intégrité du dépôt
cd /workspace
git status
git log -n 1 --oneline
```

### 4.2 Lancement du build

```bash
# Exécution du script de construction
cd /workspace/debian-sovereign
sudo ./build-iso.sh

# Sortie attendue:
# ========================================
#   SNISID Sovereign OS Builder v1.0
#   Gouvernement Haïtien - TOP SECRET
# ========================================
# [INFO] Vérification des prérequis...
# [SUCCÈS] Toutes les dépendances sont installées
# ...
# [SUCCÈS] ISO créée: /workspace/snisid-sovereign-os-v1.0.iso
# [SUCCÈS] Checksum SHA256 généré
```

### 4.3 Vérification de l'ISO

```bash
# Vérification du checksum
cd /workspace
sha256sum -c snisid-sovereign-os-v1.0.iso.sha256

# Doit afficher: snisid-sovereign-os-v1.0.iso: OK

# Test en machine virtuelle (optionnel)
virt-manager &
# Créer une nouvelle VM avec l'ISO SNISID
```

---

## 5. INSTALLATION SUR MATÉRIEL

### 5.1 Création du support bootable

```bash
# Méthode 1: Gravure USB avec dd (Linux)
sudo dd if=snisid-sovereign-os-v1.0.iso of=/dev/sdX bs=4M status=progress conv=fsync
sync

# Méthode 2: Outil graphique (Windows/Mac)
# Utiliser Rufus (Windows) ou balenaEtcher (Mac/Linux)
# Mode: DD Image (pas ISO Image)
```

### 5.2 Procédure d'installation

1. **Démarrage sur USB**
   - Insérer la clé USB certifiée
   - Démarrer le poste avec Secure Boot activé
   - Appuyer sur F12 pour le menu de boot
   - Sélectionner "UEFI: SNISID_SOVEREIGN_V1.0"

2. **Menu de boot SNISID**
   ```
   → SNISID Sovereign OS v1.0 (Normal)
     SNISID Secure Mode (No Network)
     SNISID Rescue Mode
   ```

3. **Processus d'installation automatique**
   - Détection du TPM 2.0 ✓
   - Initialisation LUKS2 (AES-512) ✓
   - Partitionnement guidé ✓
   - Installation des paquets ✓
   - Configuration réseau ✓
   - Injection des clés HSM ✓

4. **Premier boot**
   - Écran de bienvenue officiel
   - Enrôlement biométrique (empreinte + visage)
   - Configuration carte SmartCard
   - Définition du PIN administrateur

### 5.3 Post-installation immédiate

```bash
# Connexion SSH initiale (depuis console de management)
ssh admin@snisid-local.ht
# Mot de passe: celui défini lors de l'installation

# Exécution du script de durcissement final
sudo /root/security-hardening.sh

# Vérification de l'état de sécurité
sudo lynis audit system
sudo aide --check
```

---

## 6. CONFIGURATION POST-INSTALLATION

### 6.1 Configuration réseau Zero Trust

```bash
# Éditer la configuration WireGuard
sudo nano /etc/wireguard/wg0.conf

[Interface]
PrivateKey = <clé privée générée par HSM>
Address = 10.snisid.x.x/24
DNS = 10.snisid.0.53

[Peer]
# Datacenter Principal
PublicKey = <clé publique DC1>
Endpoint = dc1.snisid.ht:51820
AllowedIPs = 10.snisid.0.0/16

[Peer]
# Site Secondaire (Fort National)
PublicKey = <clé publique DC2>
Endpoint = dc2.snisid.ht:51820
AllowedIPs = 10.snisid.1.0/24

# Activer le tunnel
sudo systemctl enable wg-quick@wg0
sudo systemctl start wg-quick@wg0
```

### 6.2 Activation des modules NSA

```bash
# Démarrage des services DataWave
sudo systemctl enable datawave-query
sudo systemctl start datawave-query

# Vérification du statut
curl -k https://localhost:8443/datawave/version

# Configuration Emissary
sudo nano /opt/snisid/modules/emissary/config/application.yml
# Ajouter les connecteurs Kafka haïtiens

sudo systemctl enable emissary-ingest
sudo systemctl start emissary-ingest
```

### 6.3 Configuration des outils OSINT

```bash
# PhoneInfoga - Service d'investigation téléphonique
cd /opt/snisid/modules/phoneinfoga
./phoneinfoga serve -addr 127.0.0.1:9000

# OSIRIS - Reconnaissance faciale
cd /opt/snisid/modules/osiris
python3 app.py --host 127.0.0.1 --port 9001

# Storm-Breaker - Honeypot tactique (Niveau 4 uniquement)
cd /opt/snisid/modules/storm-breaker
sudo ./storm-breaker --stealth-mode
```

### 6.4 Intégration bases de données nationales

```yaml
# /etc/snisid/database-config.yml
postgresql:
  host: db-primary.snisid.ht
  port: 5432
  database: snisid_identity
  sslmode: verify-full
  sslrootcert: /etc/pki/ca.crt
  
cockroachdb:
  hosts: 
    - crdb-node1.snisid.ht:26257
    - crdb-node2.snisid.ht:26257
    - crdb-node3.snisid.ht:26257
  database: snisid_biometric
  
mongodb:
  uri: mongodb://replica-set.snisid.ht:27017/?ssl=true&replicaSet=rs0
  database: snisid_documents
  
redis:
  hosts:
    - redis-cluster-1.snisid.ht:6379
    - redis-cluster-2.snisid.ht:6379
    - redis-cluster-3.snisid.ht:6379
  password: <mot_de_passe_HSM>
```

---

## 7. INTÉGRATION DES MODULES NSA/CHINE

### 7.1 DataWave (NSA) - Configuration avancée

```xml
<!-- /opt/snisid/modules/datawave/query/conf/query-config.xml -->
<queryConfiguration>
    <federatedSources>
        <source name="postgresql_identity" type="JDBC">
            <url>jdbc:postgresql://db-primary.snisid.ht:5432/snisid_identity</url>
            <driver>org.postgresql.Driver</driver>
        </source>
        <source name="cockroach_biometric" type="JDBC">
            <url>jdbc:postgresql://crdb-node1.snisid.ht:26257/snisid_biometric</url>
        </source>
        <source name="mongodb_documents" type="MongoDB">
            <uri>mongodb://replica-set.snisid.ht:27017</uri>
        </source>
    </federatedSources>
    
    <security>
        <encryption algorithm="AES-256-GCM"/>
        <audit enabled="true" destination="elasticsearch"/>
    </security>
</queryConfiguration>
```

### 7.2 Cryptographie chinoise (SM2/SM3/SM4)

```bash
# Installation du module OpenSSL GM
sudo apt-get install -y openssl-engine-gm

# Configuration pour utiliser SM4 au lieu d'AES
cat >> /etc/ssl/openssl.cnf << EOF

# SNISID Chinese Cryptography Module
[gm_module]
engine_id = gm
dynamic_path = /usr/lib/engines-1.1/gm.so
init = 1
EOF

# Génération de clés SM2 pour certificats nationaux
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:SM2 -out sm2-key.pem
openssl req -new -key sm2-key.pem -out sm2-csr.pem -subj "/C=HT/O=Gouvernement/CN=SNISID"
```

### 7.3 Conformité MLPS 2.0 (Chine)

```bash
# Script de vérification MLPS Level 3
cat > /usr/local/bin/mlps-audit.sh << 'EOF'
#!/bin/bash
echo "=== MLPS 2.0 Level 3 Audit ==="

# 1. Contrôle d'accès
echo "[1/5] Vérification contrôle d'accès..."
grep -E "^(auth|account|password|session)" /etc/pam.d/common-*

# 2. Journalisation
echo "[2/5] Vérification journalisation..."
systemctl status auditd
ls -la /var/log/audit/

# 3. Chiffrement
echo "[3/5] Vérification chiffrement..."
cryptsetup status /dev/mapper/vg-root

# 4. Protection réseau
echo "[4/5] Vérification protection réseau..."
nft list ruleset | head -20

# 5. Intégrité
echo "[5/5] Vérification intégrité..."
aide --check | tail -10

echo "=== Audit terminé ==="
EOF
chmod +x /usr/local/bin/mlps-audit.sh
```

---

## 8. PROCÉDURES DE SÉCURITÉ

### 8.1 Key Ceremony (Cérémonie des clés)

**Classification:** TOP SECRET - Présence Ministre requise

```
PROTOCOLE DE GÉNÉRATION DES CLÉS MAÎTRESSES

Date: _______________
Lieu: Centre de Données National (Bunker)
Participants requis:
  ☐ Ministre de l'Intérieur
  ☐ Directeur Général DNP
  ☐ Responsable Sécurité SI
  ☐ Notaire d'État

ÉTAPES:

1. Isolement de la salle (brouilleur RF activé)
2. Initialisation HSM Thales Luna
3. Génération clé maîtresse LUKS (AES-512)
4. Division Shamir (3 sur 5):
   - Part 1: Ministre Intérieur (coffre-fort Ministère)
   - Part 2: Directeur DNP (coffre-fort DNP)
   - Part 3: Président République (Palais National)
   - Part 4: Premier Ministre (Primature)
   - Part 5: Archives Nationales (site secondaire)
5. Signature du procès-verbal
6. Scellement des tokens HSM

SIGNATURES:
_________________  _________________  _________________
```

### 8.2 Procédure d'urgence (Compromission)

```bash
# EN CAS DE COMPROMISSION SUSPECTÉE:

# 1. Isoler immédiatement le système
sudo systemctl stop network.target
sudo ip link set eth0 down

# 2. Verrouiller toutes les sessions
sudo pkill -KILL -u utilisateur
sudo pam_tally2 --reset

# 3. Révoquer les certificats
sudo openssl ca -revoke /etc/pki/certs/compromised.crt

# 4.Notifier le SOC National
echo "URGENT: Compromission détectée $(hostname)" | \
    mail -s "ALERTE ROUGE SNISID" soc@sni.ht

# 5. Démarrer forensique
sudo ftk imager /dev/nvme0n1 /mnt/evidence/image.dd
```

### 8.3 Rotation des clés (Quarterly)

```bash
# Script automatisé de rotation
cat > /etc/cron.quarterly/snisid-key-rotation << 'EOF'
#!/bin/bash
# Rotation trimestrielle des clés de chiffrement

DATE=$(date +%Y%m%d)
BACKUP_DIR="/backup/keys/${DATE}"

mkdir -p "$BACKUP_DIR"

# Sauvegarde anciennes clés
cp /etc/luks/*.key "$BACKUP_DIR/"
cp /etc/pki/private/*.key "$BACKUP_DIR/"

# Génération nouvelles clés
openssl rand -base64 64 > /etc/luks/master.key.new
chmod 600 /etc/luks/master.key.new

# Re-chiffrement LUKS avec nouvelle clé
echo "ancienne_clé" | lukschangekey /dev/nvme0n1p3 \
    --key-slot 0 --master-key-file /etc/luks/master.key.new

# Notification
echo "Rotation clés effectuée ${DATE}" | mail -s "Key Rotation OK" admin@sni.ht
EOF
```

---

## 9. MAINTENANCE ET MISES À JOUR

### 9.1 Cycle de mise à jour

| Type | Fréquence | Fenêtre | Approbation |
|------|-----------|---------|-------------|
| Sécurité critique | Immédiate | 24h | Directeur DNP |
| Correctifs mineurs | Hebdomadaire | Dimanche 02:00 | Admin SI |
| Majeure (version) | Trimestrielle | Maintenance planifiée | Comité Stratégique |

### 9.2 Procédure de mise à jour sécurisée

```bash
# 1. Vérification de la signature
cd /var/cache/apt/archives
gpg --verify Release.gpg Release

# 2. Simulation avant application
apt-get update
apt-get dist-upgrade --dry-run

# 3. Snapshot avant mise à jour
btrfs subvolume snapshot / /snapshots/pre-update-$(date +%Y%m%d)

# 4. Application
apt-get dist-upgrade -y

# 5. Vérification post-update
systemctl --failed
aide --check
lynis audit system

# 6. Rollback si nécessaire
# Au boot: sélectionner snapshot précédent dans GRUB
```

### 9.3 Monitoring de santé

```bash
# Tableau de bord de santé système
cat > /usr/local/bin/snisid-health << 'EOF'
#!/bin/bash
echo "╔════════════════════════════════════════╗"
echo "║     SNISID OS - Health Dashboard       ║"
echo "╚════════════════════════════════════════╝"

echo ""
echo "[CPU/Memory]"
top -bn1 | grep "Cpu(s)"
free -h

echo ""
echo "[Disk Encryption]"
cryptsetup status /dev/mapper/vg-root | grep -E "cipher|key"

echo ""
echo "[Services Critiques]"
systemctl is-active datawave-query emissary-ingest auditd wireguard

echo ""
echo "[Security Status]"
echo "  Firewall: $(nft list ruleset | wc -l) règles actives"
echo "  Failed Logins: $(lastb | wc -l) tentatives"
echo "  AIDE Status: $(aide --check 2>&1 | tail -1)"

echo ""
echo "[Last Update]"
cat /etc/apt/sources.list.d/snisid.list
EOF
chmod +x /usr/local/bin/snisid-health
```

---

## 10. DÉPANNAGE

### 10.1 Problèmes courants

| Symptôme | Cause probable | Solution |
|----------|----------------|----------|
| Boot échoue après update | Kernel incompatible | Boot sur snapshot précédent |
| WireGuard ne se connecte pas | Clés expirées | Régénérer avec `wg genkey` |
| DataWave timeout | BDD inaccessible | Vérifier connectivité `pg_isready` |
| TPM error 0x80280013 | TPM verrouillé | Clear TPM + re-enroll keys |
| SmartCard non reconnue | Driver manquant | `apt-get install pcscd libccid` |

### 10.2 Mode Rescue

```bash
# Au boot, sélectionner "SNISID Rescue Mode"

# Montage du système chiffré
cryptsetup luksOpen /dev/nvme0n1p3 rescue
vgchange -ay
mount /dev/vg/root /mnt

# Chroot pour réparation
mount --bind /dev /mnt/dev
mount --bind /proc /mnt/proc
mount --bind /sys /mnt/sys
chroot /mnt

# Commands utiles
fsck /dev/vg/root
apt-get install --reinstall linux-image-amd64
update-grub
```

### 10.3 Contact support

```
┌─────────────────────────────────────────────────────────────┐
│  SUPPORT TECHNIQUE SNISID - GOUVERNEMENT HAÏTIEN            │
├─────────────────────────────────────────────────────────────┤
│  Urgence 24/7: +509-XX-XXX-XXXX (SOC National)              │
│  Email: support-tech@sni.ht                                 │
│  Ticketing: https://tickets.sni.ht                          │
│                                                             │
│  Classification des incidents:                              │
│  P1 - Critique: Arrêt total (< 15 min réponse)              │
│  P2 - Majeur: Dégradation sévère (< 1h réponse)             │
│  P3 - Mineur: Impact limité (< 4h réponse)                  │
│  P4 - Information: Demande générale (< 24h réponse)         │
└─────────────────────────────────────────────────────────────┘
```

---

## ANNEXES

### A. Glossaire technique

- **CSfC**: Commercial Solutions for Classified (NSA)
- **MLPS**: Multi-Level Protection Scheme (Chine)
- **STIG**: Security Technical Implementation Guide
- **HSM**: Hardware Security Module
- **WORM**: Write Once Read Many (archive immuable)

### B. Références réglementaires

- Constitution de la République d'Haïti, Article 291
- Loi sur la Protection des Données Personnelles (202X)
- Décret Présidentiel créant le SNISID (202X)
- NSA Linux STIG v4, RHEL 9

### C. Checklist de déploiement

- [ ] Validation décret présidentiel
- [ ] Acquisition matériel certifié
- [ ] Formation équipe technique (2 semaines)
- [ ] Key Ceremony (présence Ministre)
- [ ] Installation site pilote (Port-au-Prince)
- [ ] Tests de pénétration (équipe rouge)
- [ ] Certification finale
- [ ] Déploiement national progressif

---

**Document classifié TOP SECRET**
**Distribution limitée aux autorités habilitées**
**Gouvernement de la République d'Haïti © 2024**
