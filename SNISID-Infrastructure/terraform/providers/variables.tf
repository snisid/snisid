# SNISID — Terraform Shared Variables (Providers)
# Classification: SECRET
# Variables sensibles injectées via Vault / CI / SOPS

variable "dns_tsig_secret_core" {
  description = "TSIG secret pour DNS dynamic updates Core DC"
  type        = string
  sensitive   = true
}

variable "dns_tsig_secret_dr" {
  description = "TSIG secret pour DNS dynamic updates DR DC"
  type        = string
  sensitive   = true
}

variable "proxmox_api_token_core" {
  description = "Token API Proxmox Core DC"
  type        = string
  sensitive   = true
}

variable "proxmox_api_token_dr" {
  description = "Token API Proxmox DR DC"
  type        = string
  sensitive   = true
}

variable "ceph_admin_keyring_b64" {
  description = "Ceph admin keyring (base64 encodé)"
  type        = string
  sensitive   = true
}

variable "hsm_pin_core" {
  description = "PIN HSM Thales Luna 7 Core"
  type        = string
  sensitive   = true
}

variable "hsm_pin_dr" {
  description = "PIN HSM Thales Luna 7 DR"
  type        = string
  sensitive   = true
}
