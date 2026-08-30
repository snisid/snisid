#!/bin/bash
# SNISID — nftables : National Host Firewall Rules
# Classification: SECRET
# Role: Firewall host-level pour nœuds K8s / Proxmox / bastions
# Deployment: systemd service via Ansible / cloud-init
# Note: Cilium gère le L3-L7 intra-cluster ; nftables est le dernier rempart host-level.

set -euo pipefail

NFT_FILE="/etc/nftables.conf"
BACKUP_DIR="/etc/nftables.backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

backup_previous() {
    mkdir -p "${BACKUP_DIR}"
    if [ -f "${NFT_FILE}" ]; then
        cp "${NFT_FILE}" "${BACKUP_DIR}/nftables.conf.${TIMESTAMP}"
    fi
}

generate_rules() {
    cat << 'EOF'
#!/usr/sbin/nft -f

flush ruleset

table inet snisid_filter {
    # ───────────────────────────────────────────────
    # Sets (IPs autorisées nationales)
    # ───────────────────────────────────────────────
    set management_cidr {
        type ipv4_addr
        flags interval
        elements = { 10.0.0.0/24 }
    }

    set k8s_api_allowed {
        type ipv4_addr
        flags interval
        elements = { 10.1.0.0/16, 10.2.0.0/16 }
    }

    set etcd_peer_cidr {
        type ipv4_addr
        flags interval
        elements = { 10.1.1.0/24, 10.2.1.0/24 }
    }

    set bastion_hosts {
        type ipv4_addr
        elements = { 10.0.0.10, 10.0.0.11 }
    }

    set national_dns {
        type ipv4_addr
        elements = { 10.1.0.53, 10.2.0.53 }
    }

    set national_ntp {
        type ipv4_addr
        elements = { 10.0.0.123 }
    }

    set hsm_thales {
        type ipv4_addr
        elements = { 10.1.0.10, 10.2.0.10 }
    }

    # ───────────────────────────────────────────────
    # Chains
    # ───────────────────────────────────────────────
    chain input {
        type filter hook input priority 0; policy drop;

        # State tracking
        ct state established,related accept
        ct state invalid drop

        # Loopback
        iif "lo" accept

        # ICMP limited (health checks basiques)
        ip protocol icmp icmp type { echo-request, echo-reply } limit rate 5/second accept
        ip6 nexthdr icmpv6 icmpv6 type { echo-request, echo-reply } limit rate 5/second accept

        # SSH depuis bastions uniquement (management VLAN)
        ip saddr @bastion_hosts tcp dport 22 ct state new limit rate 10/minute accept

        # Kubernetes API (6443) depuis control plane & workers uniquement
        ip saddr @k8s_api_allowed tcp dport 6443 ct state new accept

        # etcd peer (masters uniquement — filtré par labels K8s mais defense in depth)
        ip saddr @etcd_peer_cidr tcp dport { 2379, 2380 } ct state new accept

        # Node exporter metrics (Prometheus scrape)
        ip saddr 10.1.50.0/24 tcp dport 9100 ct state new accept
        ip saddr 10.1.50.0/24 tcp dport 10250 ct state new accept

        # kubelet readonly port désactivé (10255) — JAMAIS ouvert
        # kubelet secure port (10250) déjà autorisé ci-dessus pour prometheus

        # DNS national (CoreDNS sur nodes si local cache)
        ip saddr @national_dns udp dport 53 ct state new accept
        ip saddr @national_dns tcp dport 53 ct state new accept

        # NTP national
        ip saddr @national_ntp udp dport 123 ct state new accept

        # HSM Thales Luna 7 (PKCS#11) — uniquement depuis Identity Cluster nodes
        ip saddr { 10.1.20.0/24, 10.2.20.0/24 } tcp dport 1792 ct state new accept

        # Ceph OSD / MON (cluster network)
        ip saddr { 10.1.30.0/24, 10.2.30.0/24 } tcp dport { 3300, 6789, 6800-7300 } ct state new accept
        ip saddr { 10.1.30.0/24, 10.2.30.0/24 } udp dport { 3300, 6789, 6800-7300 } ct state new accept

        # WireGuard (Cilium node encryption)
        udp dport 51871 ct state new accept

        # Drop & log tout le reste
        log prefix "SNISID-DROP-INPUT: " limit rate 10/minute
        drop
    }

    chain forward {
        type filter hook forward priority 0; policy drop;

        # Cilium / K8s pod-to-pod est géré par Cilium BPF et iptables legacy
        # nftables forward est minimaliste ici
        ct state established,related accept

        # Refuser forwarding non-établi (Cilium gère le reste)
        drop
    }

    chain output {
        type filter hook output priority 0; policy drop;

        # State tracking
        ct state established,related accept

        # Loopback
        oif "lo" accept

        # DNS queries vers résolveurs nationaux
        ip daddr @national_dns udp dport 53 ct state new accept
        ip daddr @national_dns tcp dport 53 ct state new accept

        # NTP queries
        ip daddr @national_ntp udp dport 123 ct state new accept

        # HTTPS vers registre interne, Vault, API core (service mesh egress contrôlé)
        ip daddr { 10.1.1.10, 10.1.20.51, 10.1.10.11, 10.1.10.12, 10.1.10.13 } tcp dport 443 ct state new accept

        # HSM PKCS#11 (clients Vault)
        ip daddr @hsm_thales tcp dport 1792 ct state new accept

        # Ceph cluster network (si node OSD)
        ip daddr { 10.1.30.0/24, 10.2.30.0/24 } tcp dport { 3300, 6789, 6800-7300 } ct state new accept

        # Syslog / SIEM forwarding (rsyslog TLS)
        ip daddr 10.1.50.20 tcp dport 6514 ct state new accept

        # SNMP monitoring (read-only, community string complexe)
        ip daddr 10.1.50.21 udp dport 161 ct state new accept

        # WireGuard
        udp dport 51871 ct state new accept

        # Drop & log tout egress non autorisé (Zero Trust host-level)
        log prefix "SNISID-DROP-OUTPUT: " limit rate 10/minute
        drop
    }
}
EOF
}

deploy_rules() {
    backup_previous
    generate_rules > "${NFT_FILE}"
    chmod 600 "${NFT_FILE}"
    nft -f "${NFT_FILE}"
    echo "[SNISID-FIREWALL] Rules deployed: ${NFT_FILE}"
    echo "[SNISID-FIREWALL] Backup: ${BACKUP_DIR}/nftables.conf.${TIMESTAMP}"
}

# Main
if [ "$(id -u)" -ne 0 ]; then
    echo "ERR: Must run as root."
    exit 1
fi

deploy_rules

# Persist via systemd
systemctl enable nftables || true
systemctl restart nftables || true

# Validate
nft list ruleset | grep -q "snisid_filter" && echo "[SNISID-FIREWALL] Validation OK."
