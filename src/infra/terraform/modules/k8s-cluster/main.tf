variable "cluster_name" {
  description = "Name of the bare-metal K8s cluster"
  type        = string
}

variable "kubernetes_version" {
  description = "Target Kubernetes Version"
  type        = string
  default     = "v1.28.0"
}

# Simulated Bare-Metal provisioning module for SNISID
resource "null_resource" "k8s_control_plane" {
  provisioner "local-exec" {
    command = "echo 'Provisioning K8s Control Plane for ${var.cluster_name}'"
  }
}
