# SNISID — DR DC Infrastructure (Secondary Disaster Recovery Datacenter)
# Environment: production (mirror) / Tier-0..Tier-3
# Topology: Active-Active avec Core DC (réplication synchrone/async)

terraform {
  backend "s3" {
    bucket                      = "snisid-terraform-state"
    key                         = "dr/infrastructure.tfstate"
    region                      = "national"
    endpoint                    = "https://s3.interne.snisid.gouv.local"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    force_path_style            = true
  }
}

locals {
  region = "dr"
  env    = "prod"
}

# ───────────────────────────────────────────────
# Cluster Kubernetes Core-DR (Mirror API Publique)
# ───────────────────────────────────────────────
module "dr_masters" {
  source = "../../modules/proxmox-k8s"

  for_each = toset(["01", "02", "03"])

  region       = local.region
  environment  = local.env
  tier         = "1"
  role         = "master"
  index        = tonumber(each.value)
  proxmox_host = "pve-dr-${each.value}"

  cores        = 8
  memory_mb    = 16384
  disk_os_gb   = 100
  disk_data_gb = 200
  vlan_id      = 130  # VLAN DR Tier-1

  ip_address     = "10.2.10.${10 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

module "dr_workers" {
  source = "../../modules/proxmox-k8s"

  for_each = toset(["01", "02", "03", "04", "05"])

  region       = local.region
  environment  = local.env
  tier         = "1"
  role         = "worker"
  index        = tonumber(each.value)
  proxmox_host = "pve-dr-${each.value}"

  cores        = 16
  memory_mb    = 32768
  disk_os_gb   = 100
  disk_data_gb = 500
  vlan_id      = 130

  ip_address     = "10.2.10.${50 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

# ───────────────────────────────────────────────
# Cluster Kubernetes Identity-DR (Tier-0 — HSM DR)
# ───────────────────────────────────────────────
module "dr_identity_masters" {
  source = "../../modules/proxmox-k8s"

  for_each = toset(["01", "02", "03"])

  region       = local.region
  environment  = local.env
  tier         = "0"
  role         = "master"
  index        = tonumber(each.value)
  proxmox_host = "pve-dr-identity-${each.value}"

  cores        = 8
  memory_mb    = 16384
  disk_os_gb   = 100
  disk_data_gb = 200
  vlan_id      = 120  # VLAN DR Tier-0

  ip_address     = "10.2.20.${10 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

module "dr_identity_workers" {
  source = "../../modules/proxmox-k8s"

  for_each = toset(["01", "02", "03"])

  region       = local.region
  environment  = local.env
  tier         = "0"
  role         = "worker"
  index        = tonumber(each.value)
  proxmox_host = "pve-dr-identity-${each.value}"

  cores        = 16
  memory_mb    = 32768
  disk_os_gb   = 100
  disk_data_gb = 500
  vlan_id      = 120

  ip_address     = "10.2.20.${50 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

# ───────────────────────────────────────────────
# Cluster Kubernetes Data-DR (Mirror Kafka/Ceph/PostgreSQL)
# ───────────────────────────────────────────────
module "dr_data_masters" {
  source = "../../modules/proxmox-k8s"

  for_each = toset(["01", "02", "03"])

  region       = local.region
  environment  = local.env
  tier         = "1"
  role         = "master"
  index        = tonumber(each.value)
  proxmox_host = "pve-dr-data-${each.value}"

  cores        = 8
  memory_mb    = 16384
  disk_os_gb   = 100
  disk_data_gb = 500
  vlan_id      = 130

  ip_address     = "10.2.30.${10 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

module "dr_data_workers" {
  source = "../../modules/proxmox-k8s"

  for_each = toset(["01", "02", "03", "04", "05", "06"])

  region       = local.region
  environment  = local.env
  tier         = "1"
  role         = "worker"
  index        = tonumber(each.value)
  proxmox_host = "pve-dr-data-${each.value}"

  cores        = 16
  memory_mb    = 65536
  disk_os_gb   = 100
  disk_data_gb = 2000
  vlan_id      = 130

  ip_address     = "10.2.30.${50 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

# ───────────────────────────────────────────────
# Cluster Observability-DR (Tier-3 — Mirror métrologie)
# ───────────────────────────────────────────────
module "dr_obs_masters" {
  source = "../../modules/proxmox-k8s"

  for_each = toset(["01", "02", "03"])

  region       = local.region
  environment  = local.env
  tier         = "3"
  role         = "master"
  index        = tonumber(each.value)
  proxmox_host = "pve-dr-obs-${each.value}"

  cores        = 8
  memory_mb    = 16384
  disk_os_gb   = 100
  disk_data_gb = 1000
  vlan_id      = 150

  ip_address     = "10.2.50.${10 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

module "dr_obs_workers" {
  source = "../../modules/proxmox-k8s"

  for_each = toset(["01", "02", "03", "04"])

  region       = local.region
  environment  = local.env
  tier         = "3"
  role         = "worker"
  index        = tonumber(each.value)
  proxmox_host = "pve-dr-obs-${each.value}"

  cores        = 16
  memory_mb    = 65536
  disk_os_gb   = 100
  disk_data_gb = 4000
  vlan_id      = 150

  ip_address     = "10.2.50.${50 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

# ───────────────────────────────────────────────
# HSM DR — Thales Luna 7 (réplication clés Core→DR)
# ───────────────────────────────────────────────
resource "proxmox_vm_qemu" "dr_hsm_thales" {
  name        = "dr-hsm-thales-01"
  target_node = "pve-dr-hsm-01"
  clone       = "ubuntu-22-04-cis-hardened"
  agent       = 1
  cores       = 4
  sockets     = 1
  cpu         = "host"
  memory      = 8192
  disk {
    type    = "scsi"
    storage = "ceph-rbd-dr"
    size    = 50
    discard = "on"
    ssd     = 1
  }
  network {
    model    = "virtio"
    bridge   = "vmbr0"
    tag      = 120
    firewall = true
  }
  os_type   = "cloud-init"
  ipconfig0 = "ip=10.2.0.10/24,gw=10.2.0.1"
  ciuser    = "snisid-hsm-admin"
  sshkeys   = file("../../../.secrets/snisid-hsm-admin-ssh.pub")
  tags      = "snisid;tier-0;region-dr;env-prod;role=hsm;managed-by-terraform"
}

# ───────────────────────────────────────────────
# Réplication Ceph multi-site (RGW pools DR)
# ───────────────────────────────────────────────
# Note: Déployé via Rook manifests ArgoCD — Terraform provisionne hosts uniquement
