# SNISID — Core DC Infrastructure (National Core Datacenter)
# Environment: production / Tier-0..Tier-3

terraform {
  backend "s3" {
    bucket                      = "snisid-terraform-state"
    key                         = "core/infrastructure.tfstate"
    region                      = "national"
    endpoint                    = "https://s3.interne.snisid.gouv.local"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    force_path_style            = true
  }
}

locals {
  region = "core"
  env    = "prod"
}

# ───────────────────────────────────────────────
# Cluster Kubernetes Core (API Publique SNISID)
# ───────────────────────────────────────────────
module "core_masters" {
  source = "../../modules/proxmox-k8s"
  
  for_each = toset(["01", "02", "03"])
  
  region       = local.region
  environment  = local.env
  tier         = "1"
  role         = "master"
  index        = tonumber(each.value)
  proxmox_host = "pve-core-${each.value}"
  
  cores       = 8
  memory_mb   = 16384
  disk_os_gb  = 100
  disk_data_gb = 200
  vlan_id     = 30
  
  ip_address   = "10.1.10.${10 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

module "core_workers" {
  source = "../../modules/proxmox-k8s"
  
  for_each = toset(["01", "02", "03", "04", "05"])
  
  region       = local.region
  environment  = local.env
  tier         = "1"
  role         = "worker"
  index        = tonumber(each.value)
  proxmox_host = "pve-core-${each.value}"
  
  cores       = 16
  memory_mb   = 32768
  disk_os_gb  = 100
  disk_data_gb = 500
  vlan_id     = 30
  
  ip_address   = "10.1.10.${50 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

# ───────────────────────────────────────────────
# Cluster Kubernetes Identity (Tier-0 — Isolation critique)
# ───────────────────────────────────────────────
module "identity_masters" {
  source = "../../modules/proxmox-k8s"
  
  for_each = toset(["01", "02", "03"])
  
  region       = local.region
  environment  = local.env
  tier         = "0"
  role         = "master"
  index        = tonumber(each.value)
  proxmox_host = "pve-core-identity-${each.value}"
  
  cores       = 8
  memory_mb   = 16384
  disk_os_gb  = 100
  disk_data_gb = 200
  vlan_id     = 20  # VLAN isolé Tier-0
  
  ip_address   = "10.1.20.${10 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

module "identity_workers" {
  source = "../../modules/proxmox-k8s"
  
  for_each = toset(["01", "02", "03"])
  
  region       = local.region
  environment  = local.env
  tier         = "0"
  role         = "worker"
  index        = tonumber(each.value)
  proxmox_host = "pve-core-identity-${each.value}"
  
  cores       = 16
  memory_mb   = 32768
  disk_os_gb  = 100
  disk_data_gb = 500
  vlan_id     = 20
  
  ip_address   = "10.1.20.${50 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

# ───────────────────────────────────────────────
# Cluster Kubernetes Data (Bases, Kafka, Search)
# ───────────────────────────────────────────────
module "data_masters" {
  source = "../../modules/proxmox-k8s"
  
  for_each = toset(["01", "02", "03"])
  
  region       = local.region
  environment  = local.env
  tier         = "1"
  role         = "master"
  index        = tonumber(each.value)
  proxmox_host = "pve-core-data-${each.value}"
  
  cores       = 8
  memory_mb   = 16384
  disk_os_gb  = 100
  disk_data_gb = 500
  vlan_id     = 30
  
  ip_address   = "10.1.30.${10 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

module "data_workers" {
  source = "../../modules/proxmox-k8s"
  
  for_each = toset(["01", "02", "03", "04", "05", "06"])
  
  region       = local.region
  environment  = local.env
  tier         = "1"
  role         = "worker"
  index        = tonumber(each.value)
  proxmox_host = "pve-core-data-${each.value}"
  
  cores       = 16
  memory_mb   = 65536
  disk_os_gb  = 100
  disk_data_gb = 2000
  vlan_id     = 30
  
  ip_address   = "10.1.30.${50 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

# ───────────────────────────────────────────────
# Cluster Observability (Tier-3 — Isolation métrologie)
# ───────────────────────────────────────────────
module "obs_masters" {
  source = "../../modules/proxmox-k8s"
  
  for_each = toset(["01", "02", "03"])
  
  region       = local.region
  environment  = local.env
  tier         = "3"
  role         = "master"
  index        = tonumber(each.value)
  proxmox_host = "pve-core-obs-${each.value}"
  
  cores       = 8
  memory_mb   = 16384
  disk_os_gb  = 100
  disk_data_gb = 1000
  vlan_id     = 50
  
  ip_address   = "10.1.50.${10 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}

module "obs_workers" {
  source = "../../modules/proxmox-k8s"
  
  for_each = toset(["01", "02", "03", "04"])
  
  region       = local.region
  environment  = local.env
  tier         = "3"
  role         = "worker"
  index        = tonumber(each.value)
  proxmox_host = "pve-core-obs-${each.value}"
  
  cores       = 16
  memory_mb   = 65536
  disk_os_gb  = 100
  disk_data_gb = 4000
  vlan_id     = 50
  
  ip_address   = "10.1.50.${50 + tonumber(each.value)}"
  ssh_public_key = file("../../../.secrets/snisid-admin-ssh.pub")
}
