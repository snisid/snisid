module "network" {
  source       = "../../modules/network"
  cluster_name = var.region
}

module "k8s_cluster" {
  source            = "../../modules/k8s-cluster"
  cluster_name      = var.region
  worker_node_count = var.worker_count
  gpu_node_count    = 2
}

module "storage" {
  source                   = "../../modules/storage"
  ceph_storage_capacity_tb = 250
}
