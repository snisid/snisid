variable "worker_node_count" {
  description = "Number of standard stateless worker nodes"
  type        = number
  default     = 3
}

variable "gpu_node_count" {
  description = "Number of GPU nodes for ABIS and Fraud Detection"
  type        = number
  default     = 2
}

resource "null_resource" "standard_node_pool" {
  count = var.worker_node_count
  provisioner "local-exec" {
    command = "echo 'Provisioning Standard Worker Node ${count.index}'"
  }
}

resource "null_resource" "gpu_node_pool" {
  count = var.gpu_node_count
  provisioner "local-exec" {
    command = "echo 'Provisioning GPU Worker Node ${count.index} with NVIDIA Taints'"
  }
}
