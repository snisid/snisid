provider "kubernetes" {
  config_path = "~/.kube/config"
}

# 1. Network & Isolation Module
module "network" {
  source = "./modules/network"
}

# 2. Sovereign Kubernetes Cluster (Multi-Node)
module "k8s" {
  source     = "./modules/k8s"
  node_count = 10
  cpu        = "24"
  memory     = "128GB"
}

# 3. Persistent Storage Mesh (Neo4j + MinIO)
module "storage" {
  source        = "./modules/storage"
  neo4j_enabled = true
  minio_enabled = true
}

# 4. Zero-Trust Security (Istio + SPIFFE)
module "security" {
  source       = "./modules/security"
  enable_istio = true
}

output "cluster_endpoint" {
  value = module.k8s.endpoint
}
