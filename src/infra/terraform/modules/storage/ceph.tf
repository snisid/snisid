variable "ceph_storage_capacity_tb" {
  description = "Capacity of the Ceph cluster in Terabytes"
  type        = number
  default     = 100
}

# Simulated Ceph Object Storage Configuration
resource "null_resource" "ceph_s3_worm_bucket" {
  provisioner "local-exec" {
    command = "echo 'Configuring Ceph S3 Bucket with WORM (Object Lock) for Audit Logs'"
  }
}

resource "null_resource" "ceph_csi_driver" {
  provisioner "local-exec" {
    command = "echo 'Deploying Ceph CSI Driver for StatefulSet Persistent Volumes'"
  }
}
