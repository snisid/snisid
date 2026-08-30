variable "region" {
  description = "Région nationale (core, dr, regional)"
  type        = string
}

variable "namespace" {
  description = "Namespace Kubernetes"
  type        = string
  default     = "snisid-identity"
}

variable "deployment_name" {
  description = "Nom de déploiement Vault"
  type        = string
  default     = "snisid-national"
}

variable "vault_version" {
  description = "Version Vault"
  type        = string
  default     = "1.15.4"
}

variable "registry" {
  description = "Registry conteneurs nationale"
  type        = string
  default     = "registry.interne.snisid.gouv.local"
}

variable "replicas" {
  description = "Nombre de replicas Vault (recommandé 5 pour HA)")
  type        = number
  default     = 5
}

variable "storage_class" {
  description = "StorageClass pour PVC Vault (Ceph RBD Tier-0)"
  type        = string
  default     = "ceph-rbd-tier0"
}

variable "hsm_lib_path" {
  description = "Chemin lib HSM Thales Luna 7 dans l'image"
  type        = string
  default     = "/usr/lib/libCryptoki2_64.so"
}

variable "hsm_slot" {
  description = "Slot HSM PKCS#11"
  type        = string
  default     = "0"
}

variable "hsm_key_label" {
  description = "Label clé auto-unseal HSM"
  type        = string
  default     = "snisid-vault-auto-unseal"
}

variable "hsm_hmac_key_label" {
  description = "Label clé HMAC HSM"
  type        = string
  default     = "snisid-vault-hmac"
}

variable "vault_tls_cert_pem" {
  description = "Certificat TLS Vault (PEM)"
  type        = string
  sensitive   = true
}

variable "vault_tls_key_pem" {
  description = "Clé privée TLS Vault (PEM)"
  type        = string
  sensitive   = true
}

variable "vault_ca_cert_pem" {
  description = "CA certificat TLS (PEM)"
  type        = string
  sensitive   = true
}

variable "kubeconfig_path" {
  description = "Chemin kubeconfig"
  type        = string
  default     = "~/.kube/config"
}
