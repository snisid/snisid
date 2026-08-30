variable "region" {
  description = "Région nationale (core, dr, regional-01..N)"
  type        = string
}

variable "tier" {
  description = "Tier criticité du stockage"
  type        = string
  default     = "tier-1"
}

variable "namespace" {
  description = "Namespace Kubernetes pour Rook Ceph"
  type        = string
  default     = "rook-ceph"
}

variable "registry" {
  description = "Registry conteneurs nationale"
  type        = string
  default     = "registry.interne.snisid.gouv.local"
}

variable "mon_count" {
  description = "Nombre de MONs Ceph (impair, 3 minimum)"
  type        = number
  default     = 3
}

variable "storage_nodes" {
  description = "Liste des nœuds de stockage avec devices"
  type = list(object({
    name    = string
    devices = list(object({ name = string }))
  }))
  default = []
}

variable "pools" {
  description = "Liste des pools Ceph à créer"
  type = list(object({
    name               = string
    replicated_size    = number
    failure_domain     = string
    device_class       = optional(string, "ssd")
    compression_mode   = optional(string, "aggressive")
    storage_class_name = string
    reclaim_policy     = optional(string, "Delete")
    is_default         = optional(bool, false)
  }))
  default = []
}

variable "enable_filesystem" {
  description = "Activer CephFS partagé"
  type        = bool
  default     = true
}

variable "enable_object_store" {
  description = "Activer Ceph Object Gateway (S3)"
  type        = bool
  default     = true
}

variable "kubeconfig_path" {
  description = "Chemin kubeconfig pour provider K8s"
  type        = string
  default     = "~/.kube/config"
}
