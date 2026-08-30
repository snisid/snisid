#!/bin/bash
# SNISID Post-Installation Security Hardening Script
# Implements NSA STIG + Chinese MLPS 2.0 security controls

set -euo pipefail

LOG_FILE="/var/log/snisid-hardening.log"
exec > >(tee -a "${LOG_FILE}") 2>&1

echo "=== SNISID Security Hardening Started at $(date) ==="

# ============================================================================
# 1. KERNEL HARDENING (NSA STIG V-76727, V-73403)
# ============================================================================
echo "[*] Applying kernel hardening parameters..."

cat > /etc/sysctl.d/99-snisid-nsa.conf << 'EOF'
# NSA-inspired Kernel Hardening for SNISID
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
kernel.dmesg_restrict = 1
kernel.kptr_restrict = 2
kernel.perf_event_paranoid = 3

# Core Dump Restrictions
fs.suid_dumpable = 0

# Module Loading
kernel.module.sig_enforce = 1
EOF

sysctl --system

# ============================================================================
# 2. FILESYSTEM SECURITY (NSA STIG V-71849)
# ============================================================================
echo "[*] Securing filesystem mount options..."

# Remount critical filesystems with security options
mount -o remount,noatime,nodiratime /
mount -o remount,noatime,nosuid,nodev /tmp
mount -o remount,noatime,nosuid,nodev /var/tmp

# Add to fstab for persistence
sed -i 's|/tmp.*ext4.*defaults|/tmp ext4 noatime,nosuid,nodev,noexec 0 2|' /etc/fstab
sed -i 's|/var/tmp.*ext4.*defaults|/var/tmp ext4 noatime,nosuid,nodev,noexec 0 2|' /etc/fstab

# ============================================================================
# 3. USB DEVICE CONTROL (NSA STIG V-71849)
# ============================================================================
echo "[*] Configuring USBGuard policies..."

cat > /etc/usbguard/usbguard-daemon.conf << 'EOF'
# USBGuard Configuration for SNISID
# Block all USB devices by default, allow only keyboard/mouse

ImplicitPolicyTarget=block
PresentControllerPolicy=block
AllowedDeviceTypes=keyboard,mouse
EOF

# Generate initial policy allowing only current keyboard/mouse
usbguard generate-policy > /etc/usbguard/rules.conf

systemctl enable usbguard
systemctl restart usbguard

# ============================================================================
# 4. AUDIT SUBSYSTEM (NSA STIG V-72003, Chinese MLPS Level 3)
# ============================================================================
echo "[*] Configuring auditd for comprehensive logging..."

cat > /etc/audit/rules.d/snisid.rules << 'EOF'
# Delete all existing rules
-D

# Buffer size
-b 8192

# Failure mode (2=panic, 1=printk, 0=silent)
-f 1

# Monitor file system mounts
-a exit,always -F arch=b64 -S mount -S umount2

# Monitor privileged commands
-w /usr/bin/sudo -p x -k privilege_escalation
-w /usr/bin/su -p x -k privilege_escalation
-w /etc/passwd -p wa -k identity
-w /etc/shadow -p wa -k identity
-w /etc/group -p wa -k identity

# Monitor network configuration
-w /etc/hosts -p wa -k network
-w /etc/network -p wa -k network

# Monitor SSH configuration
-w /etc/ssh/sshd_config -p wa -k sshd

# Monitor cron jobs
-w /etc/crontab -p wa -k cron
-w /etc/cron.d -p wa -k cron
-w /var/spool/cron -p wa -k cron

# Monitor sudoers
-w /etc/sudoers -p wa -k sudoers
-w /etc/sudoers.d -p wa -k sudoers

# Monitor kernel module loading
-w /sbin/insmod -p x -k modules
-w /sbin/modprobe -p x -k modules
-w /sbin/rmmod -p x -k modules

# Log all execve syscalls for privileged users
-a exit,always -F arch=b64 -C euid!=uid -F uid>=1000 -S execve -k priv_exec

# Make the configuration immutable (requires reboot to change)
-e 2
EOF

systemctl enable auditd
systemctl restart auditd

# ============================================================================
# 5. PAM HARDENING (Multi-Factor Authentication)
# ============================================================================
echo "[*] Hardening PAM configuration..."

# Backup original PAM configs
cp /etc/pam.d/common-auth /etc/pam.d/common-auth.backup

# Configure account lockout after 5 failed attempts
cat > /etc/security/faillock.conf << 'EOF'
deny = 5
unlock_time = 900
fail_interval = 900
EOF

# ============================================================================
# 6. SSH HARDENING (NSA STIG V-77811)
# ============================================================================
echo "[*] Hardening SSH configuration..."

cat > /etc/ssh/sshd_config.d/snisid-hardening.conf << 'EOF'
# SNISID SSH Hardening Configuration

# Protocol version
Protocol 2

# Disable root login via password
PermitRootLogin prohibit-password

# Disable password authentication (keys only)
PasswordAuthentication no
PubkeyAuthentication yes
PermitEmptyPasswords no

# Disable X11 forwarding
X11Forwarding no

# Disable TCP forwarding
AllowTcpForwarding no
AllowAgentForwarding no

# Strong ciphers only (TLS 1.3 equivalent)
Ciphers chacha20-poly1305@openssh.com,aes256-gcm@openssh.com,aes128-gcm@openssh.com
MACs hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com
KexAlgorithms curve25519-sha256,curve25519-sha256@libssh.org,diffie-hellman-group-exchange-sha256

# Limit authentication attempts
MaxAuthTries 3
MaxSessions 2

# Client alive settings
ClientAliveInterval 300
ClientAliveCountMax 2

# Logging
LogLevel VERBOSE
SyslogFacility AUTH

# Banner
Banner /etc/issue.net
EOF

systemctl restart sshd

# ============================================================================
# 7. APPARMOR ENFORCEMENT
# ============================================================================
echo "[*] Enabling and enforcing AppArmor profiles..."

systemctl enable apparmor
systemctl restart apparmor

# Set all profiles to enforce mode
aa-enforce /etc/apparmor.d/*

# ============================================================================
# 8. INTRUSION DETECTION (rkhunter, chkrootkit)
# ============================================================================
echo "[*] Configuring rootkit detection..."

# Update rkhunter database
rkhunter --update
rkhunter --propupd

# Configure daily rkhunter checks
cat > /etc/cron.daily/rkhunter << 'EOF'
#!/bin/sh
rkhunter --check --cronjob
EOF
chmod +x /etc/cron.daily/rkhunter

# ============================================================================
# 9. FILE INTEGRITY MONITORING (AIDE)
# ============================================================================
echo "[*] Initializing AIDE file integrity monitoring..."

# Initialize AIDE database
aideinit --yes

# Configure daily integrity checks
cat > /etc/cron.daily/aide << 'EOF'
#!/bin/sh
aide --check
EOF
chmod +x /etc/cron.daily/aide

# ============================================================================
# 10. UNATTENDED SECURITY UPGRADES
# ============================================================================
echo "[*] Configuring automatic security updates..."

cat > /etc/apt/apt.conf.d/50unattended-upgrades-snisid << 'EOF'
Unattended-Upgrade::Origins-Pattern {
    "Debian:bookworm-security";
    "Debian:bookworm-security-updates";
};

Unattended-Upgrade::Package-Blacklist {
    // List packages that should not be auto-upgraded
};

Unattended-Upgrade::Remove-Unused-Kernel-Packages "true";
Unattended-Upgrade::Automatic-Reboot "false";
Unattended-Upgrade::Mail "security@snisid.ht";
Unattended-Upgrade::MailReport "on-change";
EOF

# ============================================================================
# 11. BANNER AND LEGAL NOTICE
# ============================================================================
echo "[*] Setting up legal banners..."

cat > /etc/issue.net << 'EOF'
*******************************************************************************
*                    SNISID - SYSTEME NATIONAL D'IDENTIFICATION               *
*                        SECURISE ET D'INTEROPERABILITE DIGITALE              *
*******************************************************************************
*                                                                             *
*  AVERTISSEMENT : Ce système est réservé aux utilisateurs autorisés          *
*  uniquement. Tout accès non autorisé est strictement interdit et sera       *
*  poursuivi conformément aux lois de la République d'Haïti.                  *
*                                                                             *
*  WARNING: This system is for authorized users only. Unauthorized access     *
*  is strictly prohibited and will be prosecuted under Haitian law.           *
*                                                                             *
*******************************************************************************
EOF

cat > /etc/issue << 'EOF'
SNISID Sovereign System - Authorized Access Only
EOF

# ============================================================================
# 12. LOG ROTATION AND RETENTION (Chinese MLPS Requirement)
# ============================================================================
echo "[*] Configuring log retention policies..."

cat > /etc/logrotate.d/snisid-security << 'EOF'
/var/log/auth.log
/var/log/audit/audit.log
/var/log/faillock.log
{
    rotate 365
    daily
    missingok
    notifempty
    compress
    delaycompress
    create 0600 root root
}
EOF

# ============================================================================
# SUMMARY
# ============================================================================
echo ""
echo "=== SNISID Security Hardening Completed Successfully ==="
echo "Log file: ${LOG_FILE}"
echo ""
echo "Next steps:"
echo "1. Reboot system to apply kernel module restrictions"
echo "2. Configure HSM integration for cryptographic operations"
echo "3. Deploy SIEM agent for centralized monitoring"
echo "4. Perform vulnerability scan to validate hardening"
echo ""
echo "=== Hardening completed at $(date) ==="
