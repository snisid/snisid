variable "region" {
  description = "Région nationale (core, dr, regional-01...)"
  type        = string
}

variable "tier" {
  description = "Tier criticité"
  type        = string
  default     = "tier-0"
}

variable "namespace" {
  description = "Namespace Cilium"
  type        = string
  default     = "kube-system"
}

variable "create_namespace" {
  description = "Créer le namespace"
  type        = bool
  default     = false
}

variable "cilium_version" {
  description = "Version image Cilium"
  type        = string
  default     = "v1.15.0"
}

variable "cilium_helm_version" {
  description = "Version chart Helm Cilium"
  type        = string
  default     = "1.15.0"
}

variable "registry" {
  description = "Registry conteneurs nationale"
  type        = string
  default     = "registry.interne.snisid.gouv.local"
}

variable "k8s_api_host" {
  description = "Hostname API K8s (ex: api.core.snisid.gouv.local)"
  type        = string
}

variable "k8s_api_port" {
  description = "Port API K8s"
  type        = number
  default     = 6443
}

variable "pod_cidr" {
  description = "CIDR Pods K8s"
  type        = string
  default     = "10.244.0.0/16"
}

variable "cluster_id" {
  description = "ID cluster Cilium (mesh)"
  type        = number
  default     = 1
}

variable "operator_replicas" {
  description = "Replicas operator Cilium"
  type        = number
  default     = 3
}

variable "enable_hubble" {
  description = "Activer Hubble observabilité"
  type        = bool
  default     = true
}

variable "enable_hubble_ui" {
  description = "Activer UI Hubble"
  type        = bool
  default     = true
}

variable "enable_tetragon" {
  description = "Activer Tetragon runtime security"
  type        = bool
  default     = true
}

variable "kubeconfig_path" {
  description = "Chemin kubeconfig"
  type        = string
  default     = "~/.kube/config"
}
