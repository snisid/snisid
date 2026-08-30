variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
}

variable "pod_cidr" {
  description = "CIDR range for Kubernetes Pods"
  type        = string
  default     = "10.244.0.0/16"
}

variable "service_cidr" {
  description = "CIDR range for Kubernetes Services"
  type        = string
  default     = "10.96.0.0/12"
}

output "pod_cidr" {
  value = var.pod_cidr
}

output "service_cidr" {
  value = var.service_cidr
}
