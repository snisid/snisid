# SNISID — Terraform Module : Proxmox Kubernetes Node Provisioning
# Classification: RESTREINT DEFENSE
# Role: Provision VMs souveraines pour clusters Kubernetes nationaux

terraform {
  required_providers {
    proxmox = {
      source  = "Telmate/proxmox"
      version = ">= 2.9.0"
    }
  }
}

locals {
  node_name = "${var.region}-${var.environment}-${var.tier}-${var.role}-${format("%02d", var.index + 1)}"
  
  common_tags = [
    "snisid",
    "tier-${var.tier}",
    "region-${var.region}",
    "env-${var.environment}",
    "role-${var.role}",
    "managed-by-terraform"
  ]
}

# VM Master / Worker Kubernetes
resource "proxmox_vm_qemu" "k8s_node" {
  name        = local.node_name
  target_node = var.proxmox_host
  clone       = var.template_name
  agent       = 1
  
  # Hardware souverain
  cores   = var.cores
  sockets = 1
  cpu     = "host"  # Pas de virtualisation imbriquée sauf si nécessaire
  memory  = var.memory_mb
  
  # Disque OS sur stockage Ceph RBD local Proxmox
  disk {
    type    = "scsi"
    storage = "ceph-rbd"
    size    = var.disk_os_gb
    discard = "on"
    ssd     = 1
  }
  
  # Disque données (etcd, containerd, logs)
  disk {
    type    = "scsi"
    storage = "ceph-rbd"
    size    = var.disk_data_gb
    discard = "on"
    ssd     = 1
  }
  
  # Réseau souverain — VLAN isolé par tier
  network {
    model    = "virtio"
    bridge   = "vmbr0"
    tag      = var.vlan_id
    firewall = true
  }
  
  # Cloud-init
  os_type    = "cloud-init"
  ipconfig0  = "ip=${var.ip_address}/24,gw=${var.gateway}"
  ciuser     = "snisid-admin"
  sshkeys    = var.ssh_public_key
  
  tags = join(";", local.common_tags)
  
  # Sécurité: pas de console VNC externe, SPICE local uniquement
  args = "-vga none -serial null"
  
  lifecycle {
    prevent_destroy = var.tier == "0" ? true : false
    ignore_changes  = [network]
  }
}

# Enregistrement DNS interne souverain
resource "dns_a_record_set" "node_dns" {
  count = var.register_dns ? 1 : 0
  
  zone      = "snisid.gouv.local."
  name      = local.node_name
  addresses = [var.ip_address]
  ttl       = 300
}

# Firewall Proxmox (nftables backend) — Deny all par défaut
resource "proxmox_firewall_rules" "node_firewall" {
  count = var.tier == "0" ? 1 : 0
  
  node_name = local.node_name
  
  # SSH depuis bastions uniquement (CIDR management)
  rule {
    type  = "in"
    action = "ACCEPT"
    macro = "SSH"
    source = var.management_cidr
    log    = "info"
  }
  
  # Kubernetes API depuis management + autres masters
  rule {
    type   = "in"
    action = "ACCEPT"
    dport  = "6443"
    proto  = "tcp"
    source = var.k8s_api_allowed_cidr
  }
  
  # etcd peer (masters uniquement)
  dynamic "rule" {
    for_each = var.role == "master" ? [1] : []
    content {
      type   = "in"
      action = "ACCEPT"
      dport  = "2379-2380"
      proto  = "tcp"
      source = var.etcd_peer_cidr
    }
  }
  
  # Deny all inbound par défaut (implicite, mais explicite ici pour audit)
  rule {
    type   = "in"
    action = "DROP"
    log    = "alert"
  }
}
