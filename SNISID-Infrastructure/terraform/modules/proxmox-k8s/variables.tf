variable "region" {
  description = "Région nationale (core, dr, regional-01..N, edge-01..N)"
  type        = string
}

variable "environment" {
  description = "Environnement (prod, staging, dr)"
  type        = string
  default     = "prod"
}

variable "tier" {
  description = "Tier criticité (0, 1, 2, 3, 4)"
  type        = string
}

variable "role" {
  description = "Rôle K8s (master, worker, etcd, bastion)"
  type        = string
}

variable "index" {
  description = "Index numérique du nœud"
  type        = number
}

variable "proxmox_host" {
  description = "Nœud Proxmox cible"
  type        = string
}

variable "template_name" {
  description = "Template cloud-init Ubuntu 22.04 LTS durci CIS"
  type        = string
  default     = "ubuntu-22-04-cis-hardened"
}

variable "cores" {
  type    = number
  default = 4
}

variable "memory_mb" {
  type    = number
  default = 8192
}

variable "disk_os_gb" {
  type    = number
  default = 50
}

variable "disk_data_gb" {
  type    = number
  default = 100
}

variable "vlan_id" {
  description = "VLAN segment réseau (10=Management, 20=Tier-0, 30=Tier-1, 40=Tier-2, 50=Observability, 60=Edge)"
  type        = number
}

variable "ip_address" {
  type = string
}

variable "gateway" {
  type    = string
  default = "10.0.0.1"
}

variable "ssh_public_key" {
  type = string
}

variable "register_dns" {
  type    = bool
  default = true
}

variable "management_cidr" {
  type    = string
  default = "10.0.0.0/24"
}

variable "k8s_api_allowed_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "etcd_peer_cidr" {
  type    = string
  default = "10.0.1.0/24"
}
